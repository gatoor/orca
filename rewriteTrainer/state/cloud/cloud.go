/*
Copyright Alex Mack
This file is part of Orca.

Orca is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Orca is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Orca.  If not, see <http://www.gnu.org/licenses/>.
*/

package state_cloud

import (
	"errors"
	"gatoor/orca/base"
	Logger "gatoor/orca/rewriteTrainer/log"
	"gatoor/orca/rewriteTrainer/db"
	"gatoor/orca/rewriteTrainer/state/needs"
	"sync"
	"gatoor/orca/rewriteTrainer/needs"
	"fmt"
	"sort"
	"gatoor/orca/rewriteTrainer/cloud"
	"gatoor/orca/rewriteTrainer/state/configuration"
	"time"
)

var StateCloudLogger = Logger.LoggerWithField(Logger.Logger, "module", "state_cloud")

var GlobalCloudLayout CloudLayoutAll
var GlobalAvailableInstances AvailableInstances
var cloudLayoutMutex = &sync.Mutex{}
var availableInstancesMutex = &sync.Mutex{}

type CloudLayout struct {
	Type   string
	Layout map[base.HostId]CloudLayoutElement
}

type ResourceObjList []ResourceObj

type ResourceObj struct {
	HostId                     base.HostId
	CombinedAvailableResources int
}

type AvailableInstances map[base.HostId]base.InstanceResources

type CloudLayoutAll struct {
	Current CloudLayout
	Desired CloudLayout
}

type CloudLayoutElement struct {
	HostId         base.HostId
	IpAddress      base.IpAddr
	InstanceType   base.InstanceType
	SafeInstance   base.SafeInstance
	HabitatVersion base.Version
	Apps           map[base.AppName]AppsVersion
}

type AppsVersion struct {
	Version         base.Version
	DeploymentCount base.DeploymentCount

	StatisticPoint  base.AppStats
	StatisticPointTimestamp time.Time
}

func (c *CloudLayoutAll) Init() {
	c.Current = CloudLayout{
		Type: "Current",
		Layout: make(map[base.HostId]CloudLayoutElement),
	}
	c.Desired = CloudLayout{
		Type: "Desired",
		Layout: make(map[base.HostId]CloudLayoutElement),
	}
	GlobalAvailableInstances = AvailableInstances{}

}

func (c *CloudLayoutAll) InitBaseInstances(){
	if len(GlobalAvailableInstances) < int(state_configuration.GlobalConfigurationState.CloudProvider.MinInstances) {
		for i := len(GlobalAvailableInstances); i < int(state_configuration.GlobalConfigurationState.CloudProvider.MinInstances); i++ {
			cloud.CurrentProvider.SpawnInstanceSync(state_configuration.GlobalConfigurationState.CloudProvider.BaseInstanceType)
		}
	}
}

func (c *CloudLayoutAll) Snapshot() CloudLayoutAll {
	cloudLayoutMutex.Lock()
	defer cloudLayoutMutex.Unlock()
	res := *c
	return res
}

func (c *CloudLayout) GetHost(host base.HostId) (CloudLayoutElement, error) {
	cloudLayoutMutex.Lock()
	defer cloudLayoutMutex.Unlock()
	if _, exists := (*c).Layout[host]; !exists {
		return CloudLayoutElement{}, errors.New("No such host")
	}
	return (*c).Layout[host], nil
}

func (c *CloudLayout) HostHasApp(host base.HostId, appName base.AppName) bool {
	cloudLayoutMutex.Lock()
	defer cloudLayoutMutex.Unlock()
	_, exists := c.Layout[host].Apps[appName]
	return exists
}

func (c *CloudLayout) FindHostsWithApp(appName base.AppName) map[base.HostId]bool {
	cloudLayoutMutex.Lock()
	defer cloudLayoutMutex.Unlock()
	hosts := make(map[base.HostId]bool)
	//sorting to make it testable
	var hostsSort []string
	for k := range (*c).Layout {
		hostsSort = append(hostsSort, string(k))
	}
	sort.Strings(hostsSort)

	for _, hostId := range hostsSort {
		for app := range (*c).Layout[base.HostId(hostId)].Apps {
			if app == appName {
				hosts[base.HostId(hostId)] = true
			}
		}
	}

	StateCloudLogger.WithField("type", c.Type).Debugf("Found Hosts %+v with App '%s'", hosts, appName)
	return hosts
}

func (c *CloudLayout) AddHost(host base.HostId, elem CloudLayoutElement) {
	StateCloudLogger.WithField("type", c.Type).Infof("Adding host '%s': '%+v'", host, elem)
	cloudLayoutMutex.Lock()
	defer cloudLayoutMutex.Unlock()

	if _, ok := (*c).Layout[host]; !ok {
		db.Audit.Insert__AuditEvent(db.AuditEvent{Details:map[string]string{
			"message": fmt.Sprintf("New host discovered, '%s': '%+v'", host, elem),
			"subsystem": "cloud",
			"level": "info",
		}})
	}

	if (elem.InstanceType == "") {
		elem.InstanceType = cloud.CurrentProvider.GetInstanceType(host)
	}

	if (elem.IpAddress == "") {
		elem.IpAddress = cloud.CurrentProvider.GetIp(host)
	}

	(*c).Layout[host] = elem
}

func (c *CloudLayout) DeploymentCount(app base.AppName, version base.Version) (base.DeploymentCount, base.DeploymentCount) {
	cloudLayoutMutex.Lock()
	defer cloudLayoutMutex.Unlock()
	var versionCount base.DeploymentCount = 0
	var total base.DeploymentCount = 0
	for _, host := range (*c).Layout {
		if app, exists := host.Apps[app]; exists {
			if app.Version == version {
				versionCount += app.DeploymentCount
			}
			total += app.DeploymentCount
		}
	}
	return versionCount, total
}

func (c *CloudLayout) AddEmptyHost(host base.HostId) {
	StateCloudLogger.WithField("type", c.Type).Infof("Adding empty host '%s'", host)
	entity := CloudLayoutElement{
		HostId: host,
		HabitatVersion: 0,
		Apps: make(map[base.AppName]AppsVersion),
	}
	c.AddHost(host, entity)
}

func (c *CloudLayout) Wipe() {
	StateCloudLogger.WithField("type", c.Type).Infof("Wipe")
	cloudLayoutMutex.Lock()
	defer cloudLayoutMutex.Unlock()
	(*c).Layout = make(map[base.HostId]CloudLayoutElement)
}

func (c *CloudLayout) RemoveHost(host base.HostId) {
	StateCloudLogger.WithField("type", c.Type).Infof("Removing host '%s'", host)
	cloudLayoutMutex.Lock()
	defer cloudLayoutMutex.Unlock()
	delete((*c).Layout, host)
}

func (c *CloudLayout) AddApp(host base.HostId, app base.AppName, version base.Version, count base.DeploymentCount) {
	StateCloudLogger.WithField("type", c.Type).Infof("Adding App %s:%d to host '%s' %d times", app, version, host, count)
	cloudLayoutMutex.Lock()
	defer cloudLayoutMutex.Unlock()
	if val, exists := (*c).Layout[host]; exists {
		if _, ex := val.Apps[app]; !ex {
			val.Apps[app] = AppsVersion{
				version,
				count,
				base.AppStats{},
				time.Time{},
			}
		}
	}
}

func (c *CloudLayout) RemoveApp(host base.HostId, app base.AppName) {
	StateCloudLogger.WithField("type", c.Type).Infof("Removing app '%s' from host '%s'", app, host)
	cloudLayoutMutex.Lock()
	defer cloudLayoutMutex.Unlock()
	if val, exists := (*c).Layout[host]; exists {
		delete(val.Apps, app)
	}
}

func (c *CloudLayout) Needs(app base.AppName) base.AppNeeds {
	cloudLayoutMutex.Lock()
	defer cloudLayoutMutex.Unlock()
	ns := base.AppNeeds{}
	for _, elem := range (*c).Layout {
		if appObj, exists := elem.Apps[app]; exists {
			currentNeeds, err := state_needs.GlobalAppsNeedState.Get(app, appObj.Version)
			if err == nil {
				ns.CpuNeeds += base.CpuNeeds(int(appObj.DeploymentCount) * int(currentNeeds.CpuNeeds))
				ns.MemoryNeeds += base.MemoryNeeds(int(appObj.DeploymentCount) * int(currentNeeds.MemoryNeeds))
				ns.NetworkNeeds += base.NetworkNeeds(int(appObj.DeploymentCount) * int(currentNeeds.NetworkNeeds))
			}
		}
	}
	return ns
}

func (c *CloudLayout) AllNeeds() base.AppNeeds {
	cloudLayoutMutex.Lock()
	defer cloudLayoutMutex.Unlock()
	ns := base.AppNeeds{}
	for _, elem := range (*c).Layout {
		for appName, appObj := range elem.Apps {
			currentNeeds, err := state_needs.GlobalAppsNeedState.Get(appName, appObj.Version)
			if err == nil {
				ns.CpuNeeds += base.CpuNeeds(int(appObj.DeploymentCount) * int(currentNeeds.CpuNeeds))
				ns.MemoryNeeds += base.MemoryNeeds(int(appObj.DeploymentCount) * int(currentNeeds.MemoryNeeds))
				ns.NetworkNeeds += base.NetworkNeeds(int(appObj.DeploymentCount) * int(currentNeeds.NetworkNeeds))
			}
		}
	}
	return ns
}

func handleHostWithoutApps(hostIndo base.HostInfo) {
	//TODO check cloudprovider spawnlog if it is in there remove it from there
	// else: the host was cleaned of apps to shut it down. Do this
}

func (c *CloudLayout) UpdateHost(hostInfo base.HostInfo, stats base.MetricsWrapper) {
	if len(hostInfo.Apps) == 0 {
		handleHostWithoutApps(hostInfo)
	}
	apps := make(map[base.AppName]AppsVersion)
	appCounter := make(map[base.AppName]base.DeploymentCount)
	runningVersions := make(map[base.AppName]base.Version)

	for _, val := range hostInfo.Apps {
		if val.Status != base.STATUS_RUNNING {
			continue
		}
		if _, exists := appCounter[val.Name]; !exists {
			appCounter[val.Name] = 0
			runningVersions[val.Name] = val.Version
		}
		appCounter[val.Name] += 1
	}

	for appName, count := range appCounter {
		appsVersionObject := AppsVersion{
			Version: runningVersions[appName],
			DeploymentCount: count,
		}

		version := runningVersions[appName]
		for _, stat := range stats.AppMetrics[appName][version]{
			appsVersionObject.StatisticPoint = stat
			appsVersionObject.StatisticPointTimestamp = time.Now()
			break
		}

		apps[appName] = appsVersionObject
	}

	elem := CloudLayoutElement{
		HostId: hostInfo.HostId,
		IpAddress: hostInfo.IpAddr,
		Apps: apps,
	}
	StateCloudLogger.WithField("type", c.Type).Debugf("UpdateHost for host '%s': %+v", hostInfo.HostId, hostInfo)
	c.AddHost(hostInfo.HostId, elem)
}

var AvailableInstancesLogger = StateCloudLogger.WithField("type", "AvailableInstances")

func (a *AvailableInstances) HostHasResourcesForApp(hostId base.HostId, ns needs.AppNeeds) bool {
	availableInstancesMutex.Lock()
	res := (*a)[hostId]
	if int(res.TotalCpuResource - res.UsedCpuResource) >= int(ns.CpuNeeds) &&
		int(res.TotalMemoryResource - res.UsedMemoryResource) >= int(ns.MemoryNeeds) &&
		int(res.TotalNetworkResource - res.UsedNetworkResource) >= int(ns.NetworkNeeds) {
		availableInstancesMutex.Unlock()
		return true
	}
	availableInstancesMutex.Unlock()
	return false
}

func (p ResourceObjList) Len() int {
	return len(p)
}
func (p ResourceObjList) Less(i, j int) bool {
	return p[i].CombinedAvailableResources < p[j].CombinedAvailableResources
}
func (p ResourceObjList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (a *AvailableInstances) Update(hostId base.HostId, resources base.InstanceResources) {
	AvailableInstancesLogger.Infof("Updating host '%s': '%+v'", hostId, resources)
	availableInstancesMutex.Lock()
	defer availableInstancesMutex.Unlock()
	if (*a) == nil {
		(*a) = make(map[base.HostId]base.InstanceResources)
	}
	(*a)[hostId] = resources
}

func (a *AvailableInstances) GetResources(hostId base.HostId) (base.InstanceResources, error) {
	availableInstancesMutex.Lock()
	if _, exists := (*a)[hostId]; !exists {
		AvailableInstancesLogger.Warnf("Instance '%s' does not exist", hostId)
		availableInstancesMutex.Unlock()
		return base.InstanceResources{}, errors.New("Host does not exist")
	}
	AvailableInstancesLogger.Debugf("GetResources for host '%s': %+v", hostId, (*a)[hostId])
	res := (*a)[hostId]
	availableInstancesMutex.Unlock()
	return res, nil
}

func (a *AvailableInstances) Remove(hostId base.HostId) {
	AvailableInstancesLogger.Infof("Deleting instance '%s'", hostId)
	availableInstancesMutex.Lock()
	defer availableInstancesMutex.Unlock()
	if _, exists := (*a)[hostId]; !exists {
		AvailableInstancesLogger.Warnf("Instance '%s' does not exist", hostId)
		return
	}
	AvailableInstancesLogger.Debugf("Remove host '%s'", hostId)
	delete((*a), hostId)
}

func (a *AvailableInstances) GlobalResourceConsumption() base.InstanceResources {
	availableInstancesMutex.Lock()
	defer availableInstancesMutex.Unlock()
	var availableCpu, availableMemory, availableNetwork, usedCpu, usedMemory, usedNetwork int
	for _, elem := range (*a) {
		availableCpu += int(elem.TotalCpuResource)
		availableMemory += int(elem.TotalMemoryResource)
		availableNetwork += int(elem.TotalNetworkResource)
		usedCpu += int(elem.UsedCpuResource)
		usedMemory += int(elem.UsedMemoryResource)
		usedNetwork += int(elem.UsedNetworkResource)
	}
	return base.InstanceResources{
		TotalCpuResource: base.CpuResource(availableCpu),
		TotalMemoryResource: base.MemoryResource(availableMemory),
		TotalNetworkResource: base.NetworkResource(availableNetwork),
		UsedCpuResource: base.CpuResource(usedCpu),
		UsedMemoryResource: base.MemoryResource(usedMemory),
		UsedNetworkResource: base.NetworkResource(usedNetwork),
	}
}

func (a *AvailableInstances) WipeUsage() {
	StateCloudLogger.Info("Wiping usage state of available instances")
	availableInstancesMutex.Lock()
	defer availableInstancesMutex.Unlock()
	for key, elem := range *a {
		elem.UsedCpuResource = 0
		elem.UsedNetworkResource = 0
		elem.UsedMemoryResource = 0
		(*a)[key] = elem
	}
}
