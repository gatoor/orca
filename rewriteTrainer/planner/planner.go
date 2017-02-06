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

package planner

import (
	"gatoor/orca/base"
	"gatoor/orca/rewriteTrainer/state/cloud"
	"errors"
	"sync"
	Logger "gatoor/orca/rewriteTrainer/log"
	"reflect"
	"fmt"
	"gatoor/orca/rewriteTrainer/state/needs"
	"gatoor/orca/rewriteTrainer/state/configuration"
	"gatoor/orca/rewriteTrainer/db"
	"gatoor/orca/rewriteTrainer/cloud"
	"gatoor/orca/rewriteTrainer/needs"
	"sort"
	"time"
)

var PlannerLogger = Logger.LoggerWithField(Logger.Logger, "module", "planner")
var QueueLogger = Logger.LoggerWithField(PlannerLogger, "object", "queue")

type HostDesiredConfig struct {

}

type LayoutDiff map[base.HostId]map[base.AppName]state_cloud.AppsVersion

type UpdateState string

var Queue PlannerQueue

const (
	STATE_QUEUED = "STATE_QUEUED"
	STATE_APPLYING = "STATE_APPLYING"
	STATE_SUCCESS = "STATE_SUCCESS"
	STATE_FAIL = "STATE_FAIL"
	STATE_UNKNOWN = "STATE_UNKNOWN"
)

type AppsUpdateState struct {
	State UpdateState
	Version state_cloud.AppsVersion
}

type PlannerQueue struct {
	Queue map[base.HostId]map[base.AppName]AppsUpdateState
	lock *sync.Mutex
}

func NewPlannerQueue() *PlannerQueue{
	QueueLogger.Info("Initializing")
	p := &PlannerQueue{}
	p.Queue = make(map[base.HostId]map[base.AppName]AppsUpdateState)
	p.lock = &sync.Mutex{}
	QueueLogger.Info("Initialized")
	return p
}

func (p PlannerQueue) AllEmpty() bool {
	p.lock.Lock()
	defer p.lock.Unlock()
	for _, host := range p.Queue {
		if len(host) != 0 {
			return false
		}
	}
	return true
}

//sorting is done in here to make it testable
func (p PlannerQueue) Apply(diff LayoutDiff) {
	QueueLogger.Info("Applying LayoutDiff")

	var hosts []string
	for k := range diff {
		hosts = append(hosts, string(k))
	}
	sort.Strings(hosts)

	for _, host := range hosts {
		var apps []string
		for a := range diff[base.HostId(host)] {
			apps = append(apps, string(a))
		}
		sort.Strings(apps)
		for _, app := range apps {
			p.Add(base.HostId(host), base.AppName(app), diff[base.HostId(host)][base.AppName(app)])
		}
	}
}

func (p PlannerQueue) Empty(hostId base.HostId) bool {
	p.lock.Lock()
	defer p.lock.Unlock()
	if _, exists := p.Queue[hostId]; exists {
		return len(p.Queue[hostId]) == 0
	}
	return true
}

func (p PlannerQueue) Add(hostId base.HostId, appName base.AppName, elem state_cloud.AppsVersion) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if _, exists := p.Queue[hostId]; !exists {
		p.Queue[hostId] = make(map[base.AppName]AppsUpdateState)
	}
	if _, exists := p.Queue[hostId][appName]; !exists {
		QueueLogger.Infof("Adding to host '%s' app '%s': '%v'", hostId, appName, elem)
		p.Queue[hostId][appName] = AppsUpdateState{STATE_QUEUED, elem}
	}
}

func (p PlannerQueue) Get(hostId base.HostId) (map[base.AppName]AppsUpdateState, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if _, exists := p.Queue[hostId]; exists {
		if len(p.Queue[hostId]) == 0 {
			return make(map[base.AppName]AppsUpdateState), errors.New(fmt.Sprintf("no element in queue for host %s", hostId))
		}

		elem := p.Queue[hostId]
		return elem, nil
	}
	return make(map[base.AppName]AppsUpdateState), errors.New("failed to get correct AppsVersion")
}

func (p PlannerQueue) GetState(hostId base.HostId, appName base.AppName) (UpdateState, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if _, exists := p.Queue[hostId]; exists {
		if len(p.Queue[hostId]) == 0 {
			return STATE_UNKNOWN, errors.New(fmt.Sprintf("no element in queue for host %s", hostId))
		}

		if elem, ex := p.Queue[hostId][appName]; ex {
			return elem.State, nil
		}
	}
	return STATE_UNKNOWN, errors.New("failed to get correct AppsVersion")
}

func (p PlannerQueue) Remove(hostId base.HostId, appName base.AppName) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if _, exists := p.Queue[hostId]; exists {
		if _, exists := p.Queue[hostId][appName]; exists {
			//if val.State == STATE_SUCCESS || val.State == STATE_FAIL {
				QueueLogger.Infof("Removing QueueElement host '%s' app '%s'. State was %s", hostId, appName, p.Queue[hostId][appName].State)
				delete(p.Queue[hostId], appName)
			//}
		}
	}
}

func (p PlannerQueue) RemoveApp(appName base.AppName, version base.Version) {
	p.lock.Lock()
	defer p.lock.Unlock()
	for host, apps:= range p.Queue {
		for appN, appObj := range apps {
			if appN == appName && version == appObj.Version.Version {
				QueueLogger.Infof("Removing %s:%d from Queue of host '%s'", appN, version, host)
				delete(apps, appN)
			}

		}
	}
}


func (p PlannerQueue) RollbackApp(appName base.AppName, version base.Version, stableVersion base.Version) {
	p.lock.Lock()
	defer p.lock.Unlock()
	for host, apps:= range p.Queue {
		for appN, appObj := range apps {
			if appN == appName && version == appObj.Version.Version {
				appObj.Version.Version = stableVersion
				QueueLogger.Infof("Rollback %s:%d to %d on host '%s'", appN, version, stableVersion, host)
				apps[appN] = appObj
			}

		}
	}
}


func (p PlannerQueue) SetState(hostId base.HostId, appName base.AppName, state UpdateState) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if _, exists := p.Queue[hostId]; exists {
		if _, exists := p.Queue[hostId][appName]; !exists {
			return
		}
		elem := p.Queue[hostId][appName]
		elem.State = state
		p.Queue[hostId][appName] = elem
		QueueLogger.Infof("Set state of '%s' '%s' to '%s'", hostId, appName, state)
	}
}

func (p PlannerQueue) Snapshot() map[base.HostId]map[base.AppName]AppsUpdateState{
	p.lock.Lock()
	defer p.lock.Unlock()
	res := make(map[base.HostId]map[base.AppName]AppsUpdateState)
	for k,v := range p.Queue {
		res[k] = v
	}
	QueueLogger.Infof("Created snapshot")
	return res
}

func (p PlannerQueue) RemoveHost(hostId base.HostId) {
	p.lock.Lock()
	defer p.lock.Unlock()
	delete(p.Queue, hostId)
	QueueLogger.Infof("Removed host '%s'", hostId)
}


func init() {
	PlannerLogger.Info("Initializing Planner")
	Queue = *NewPlannerQueue()
	PlannerLogger.Info("Initialized Planner")
}

func Diff(master state_cloud.CloudLayout, slave state_cloud.CloudLayout) LayoutDiff {
	diff := LayoutDiff{}
	PlannerLogger.Infof("Generating diff from master: '%v' and slave: '%v'", master, slave)
	for hostId, layoutElem := range master.Layout {
		if _, exists := slave.Layout[hostId]; !exists {
			diff[hostId] = make(map[base.AppName]state_cloud.AppsVersion)
			continue
		}
		tmp := diff[hostId]
		tmp = appsDiff(layoutElem.Apps, slave.Layout[hostId].Apps)
		diff[hostId] = tmp
	}
	PlannerLogger.Infof("Generated diff: '%v'", diff)
	return diff
}

func appsDiff(master map[base.AppName]state_cloud.AppsVersion, slave map[base.AppName]state_cloud.AppsVersion) map[base.AppName]state_cloud.AppsVersion {
	diff := make(map[base.AppName]state_cloud.AppsVersion)
	for appName, versionElem := range master {
		if !reflect.DeepEqual(versionElem, slave[appName]) {
			diff[appName] = master[appName]
		}
	}
	//handle removal of an app from host:
	for appName, versionElem := range slave {
		if _, exists := master[appName]; !exists {
			if !reflect.DeepEqual(versionElem, master[appName]) {
				PlannerLogger.Infof("Got app removal for app '%s'", appName)
				diff[appName] = state_cloud.AppsVersion{versionElem.Version, 0, base.AppStats{}, time.Time{}}
			}
		}
	}
	return diff
}

func Plan() {
	layoutBefore := state_cloud.GlobalCloudLayout.Desired
	PlannerLogger.Info("Stating Plan()")
	PlannerLogger.Infof("Before planning AvailableInstances: %+v", state_cloud.GlobalAvailableInstances)
	PlannerLogger.Infof("Before planning Layout: %+v", state_cloud.GlobalCloudLayout.Desired.Layout)
	doPlan()
	PlannerLogger.Infof("After planning AvailableInstances: %+v", state_cloud.GlobalAvailableInstances)
	PlannerLogger.Infof("After planning Layout: %+v", state_cloud.GlobalCloudLayout.Desired.Layout)
	PlannerLogger.Infof("After planning Diff: %+v", Diff(state_cloud.GlobalCloudLayout.Desired, layoutBefore))
	PlannerLogger.Info("Finished Plan()")
}

func getGlobalMissingResources() base.InstanceResources {
	neededCpu, neededMem, neededNet := getGlobalMinNeeds()
	availableCpu, availableMem, availableNet := getGlobalResources()

	res := base.InstanceResources{
		TotalCpuResource: base.CpuResource(int(neededCpu) - int(availableCpu)),
		TotalMemoryResource: base.MemoryResource(int(neededMem) - int(availableMem)),
		TotalNetworkResource: base.NetworkResource(int(neededNet) - int(availableNet)),
	}
	return res
}

func InitialPlan() {
	PlannerLogger.Info("Stating initialPlan()")
	neededCpu, neededMem, neededNet := getGlobalMinNeeds()
	availableCpu, availableMem, availableNet := getGlobalResources()

	if int(neededCpu) > int(availableCpu) {
		PlannerLogger.Warnf("Not enough Cpu resources available (needed=%d - available=%d) - spawning new instance TODO", neededCpu, availableCpu)
		cloud.CurrentProvider.SpawnInstances(cloud.CurrentProvider.SuitableInstanceTypes(getGlobalMissingResources()))
		doPlan()
		return
	}
	if int(neededMem) > int(availableMem) {
		PlannerLogger.Warnf("Not enough Memory resources available (needed=%d - available=%d) - spawning new instance TODO", neededMem, availableMem)
		cloud.CurrentProvider.SpawnInstances(cloud.CurrentProvider.SuitableInstanceTypes(getGlobalMissingResources()))
		doPlan()
		return
	}
	if int(neededNet) > int(availableNet) {
		PlannerLogger.Warnf("Not enough Network resources available (needed=%d - available=%d) - spawning new instance TODO", neededNet, availableNet)
		cloud.CurrentProvider.SpawnInstances(cloud.CurrentProvider.SuitableInstanceTypes(getGlobalMissingResources()))
		doPlan()
		return
	}

	doPlan()
	PlannerLogger.Info("Finished initialPlan()")
}


type FailedAssign struct {
	TargetHost base.HostId
	AppName base.AppName
	AppVersion base.Version
	DeploymentCount base.DeploymentCount
}

type MissingAssign struct {
	AppName base.AppName
	AppVersion base.Version
	AppType base.AppType
	DeploymentCount base.DeploymentCount
}

var FailedAssigned []FailedAssign
var failedAssignMutex = &sync.Mutex{}
var MissingAssigned []MissingAssign
var missingssignMutex = &sync.Mutex{}

//TODO fancier optimizations
func doPlan() {
	wipeDesired()
	doPlanInternal()
	handleFailedAssign()
	handleMissingAssign()
	//assignSurplusResources()
}


func appPlanningOrder(allApps map[base.AppName]base.AppConfiguration) ([]base.AppName, []base.AppName) {
	//apps := make([]base.AppName, len(allApps), len(allApps))
	httpApps := make(map[base.AppName]base.AppConfiguration)
	workerApps := make(map[base.AppName]base.AppConfiguration)

	for appName, appObj := range allApps {
		if appObj.Type == base.APP_HTTP {
			httpApps[appName] = appObj
		} else {
			workerApps[appName] = appObj
		}
	}

	httpOrdered := sortByTotalNeeds(httpApps)
	workersOrdered := sortByTotalNeeds(workerApps)

	PlannerLogger.Debugf("Sorted Apps http:%+v, worker:%+v", httpOrdered, workersOrdered)

	//apps = append(httpOrdered, workersOrdered...)
	return httpOrdered, workersOrdered
}



func createChunks(apps []base.AppName) [][]base.AppName{
	const MAX_CONCURRENT = 3
	iter := int(len(apps) / MAX_CONCURRENT)
	if iter < 1 {
		iter = 1
	}
	res := [][]base.AppName{}
	for i:= 0; i < len(apps); i += iter {
		current := iter
		if i + current >= len(apps) {
			current = len(apps) - i
		}
		res = append(res, apps[i:(i+current)])
	}
	return res
}


func doPlanInternal() {
	apps := state_configuration.GlobalConfigurationState.AllAppsLatest()
	httpOrder, workerOrder := appPlanningOrder(apps)

	httpChunks := createChunks(httpOrder)
	for _, chunk := range httpChunks {
		var wg sync.WaitGroup
		wg.Add(len(chunk))
		for _, appName := range chunk {
			appObj := apps[appName]
			appObj.TargetDeploymentCount = appObj.GetDeploymentCount()
			PlannerLogger.Infof("Assigning HttpApp %s:%d. Need to do this %d times", appObj.Name, appObj.Version, appObj.TargetDeploymentCount)
			go func () {
				defer wg.Done()
				planHttp(appObj, findHttpHostWithResources, false)
			}()

		}
		wg.Wait()
	}

	for _, appName := range workerOrder {
		appObj := apps[appName]
		appObj.TargetDeploymentCount = appObj.GetDeploymentCount()
		PlannerLogger.Infof("Assigning WorkerApp %s:%d. Need to do this %d times", appObj.Name, appObj.Version, appObj.TargetDeploymentCount)
		planWorker(appObj, findHostWithResources, false)

	}
}


func handleMissingAssign() {
	PlannerLogger.Infof("Starting handleMissingAssign for %d elements", len(MissingAssigned))
	PlannerLogger.Error("HANDLE MISSING ASSIGN NOT IMPLEMENTED -- TODO call CloudProvider to spawn instances")
	PlannerLogger.Error("HANDLE MISSING ASSIGN NOT IMPLEMENTED -- TODO call CloudProvider to spawn instances")
	PlannerLogger.Error("HANDLE MISSING ASSIGN NOT IMPLEMENTED -- TODO call CloudProvider to spawn instances")
	//cloud.CurrentProvider.SpawnInstance("m1.xlarge")
	PlannerLogger.Infof("handleMissingAssign complete")
}

func handleFailedAssign() {
	PlannerLogger.Infof("Starting handleFailed Assign for %d elements", len(FailedAssigned))
	for _, failed := range FailedAssigned {
		appObj, err := state_configuration.GlobalConfigurationState.GetApp(failed.AppName, failed.AppVersion)
		if err == nil {
			PlannerLogger.Infof("Retrying failed assignment of app '%s', DeploymentCount: %d", failed.AppName, failed.DeploymentCount)
			appObj.TargetDeploymentCount = failed.DeploymentCount
			if appObj.Type == base.APP_HTTP {
				planHttp(appObj, findHttpHostWithResources, true)
			} else {
				planWorker(appObj, findHostWithResources, true)
			}
		}
	}
	PlannerLogger.Info("handleFailed Assign complete")
}

func assignSurplusResources() {

}

type HostFinderFunc func (ns base.AppNeeds, app base.AppName, sortedHosts []base.HostId, goodHosts map[base.HostId]bool) base.HostId
type DeploymentCountFunc func (resources base.InstanceResources, ns base.AppNeeds) base.DeploymentCount

func planApp(appObj base.AppConfiguration, hostFinderFunc HostFinderFunc, deploymentCountFunc DeploymentCountFunc, ignoreFailures bool) bool {
	success := true
	ns, err := state_needs.GlobalAppsNeedState.Get(appObj.Name, appObj.Version)
	if err != nil {
		return false
	}
	var deployed base.DeploymentCount
	deployed = 0
	sortedHosts := sortByAvailableResources()
	goodHosts := state_cloud.GlobalCloudLayout.Current.FindHostsWithApp(appObj.Name)

	for deployed <= appObj.TargetDeploymentCount {
		hostId := hostFinderFunc(ns, appObj.Name, sortedHosts, goodHosts)
		if hostId == "" {
			db.Audit.Insert__AuditEvent(db.AuditEvent{Details:map[string]string{
				"message": fmt.Sprintf("App %s:%d could not find suitable host", appObj.Name, appObj.Version),
				"subsystem": "planner",
				"application": string(appObj.Name),
				"application.version": string(appObj.Version),
				"level": "warning",
			}})

			success = false
			break
		}
		var depl base.DeploymentCount
		if appObj.Type == base.APP_HTTP {
			depl = 1
		} else {
			resources, err := state_cloud.GlobalAvailableInstances.GetResources(hostId)
			if err != nil {
				break
			}
			depl = deploymentCountFunc(resources, ns)
		}

		if deployed == appObj.TargetDeploymentCount {
			PlannerLogger.Infof("Assigned all deployments of App %s:%d (%d)", appObj.Name, appObj.Version, appObj.TargetDeploymentCount)
			return success
		}
		if depl > appObj.TargetDeploymentCount - deployed {
			depl = appObj.TargetDeploymentCount - deployed
		}
		if !assignAppToHost(hostId, appObj, depl) {
			if !ignoreFailures {
				addFailedAssign(hostId, appObj.Name, appObj.Version, depl)
			} else {
				db.Audit.Insert__AuditEvent(db.AuditEvent{Details:map[string]string{
					"message": fmt.Sprintf("Assign of App %s:%d failed again. Will not try again.", appObj.Name, appObj.Version),
					"subsystem": "planner",
					"application": string(appObj.Name),
					"application.version": string(appObj.Version),
					"level": "warning",
				}})
			}
			success = false
		}
		deployed += depl
	}

	if deployed < appObj.TargetDeploymentCount {
		db.Audit.Insert__AuditEvent(db.AuditEvent{Details:map[string]string{
			"message": fmt.Sprintf("App %s:%d could not deploy TargetDeploymentCount %d, only deployed %d", appObj.Name, appObj.Version, appObj.TargetDeploymentCount, deployed),
			"subsystem": "planner",
			"application": string(appObj.Name),
			"application.version": string(appObj.Version),
			"level": "warning",
		}})

		addMissingAssign(appObj.Name, appObj.Version, appObj.Type, appObj.TargetDeploymentCount - deployed)
		success = false
	}
	return success
}

func planWorker(appObj base.AppConfiguration, hostFinderFunc HostFinderFunc, ignoreFailures bool) bool {
	return planApp(appObj, hostFinderFunc, maxDeploymentOnHost, ignoreFailures)
}

func planHttp(appObj base.AppConfiguration, hostFinderFunc HostFinderFunc, ignoreFailures bool) bool {
	httpDeploymentCountFunc := func(resources base.InstanceResources, needs base.AppNeeds) base.DeploymentCount {
		return 1
	}
	return planApp(appObj, hostFinderFunc, httpDeploymentCountFunc, ignoreFailures)
}

func maxDeploymentOnHost(resources base.InstanceResources, ns base.AppNeeds) base.DeploymentCount {
	availCpu := int(resources.TotalCpuResource - resources.UsedCpuResource)
	availMem := int(resources.TotalMemoryResource - resources.UsedMemoryResource)
	availNet := int(resources.TotalNetworkResource - resources.UsedNetworkResource)
	if int(ns.CpuNeeds) == 0 || int(ns.MemoryNeeds) == 0 || int(ns.NetworkNeeds) == 0 {
		return 0
	}
	maxCpu := int(availCpu / int(ns.CpuNeeds))
	maxMem := int(availMem / int(ns.MemoryNeeds))
	maxNet := int(availNet / int(ns.NetworkNeeds))
	if maxCpu <= maxMem && maxCpu <= maxNet {
		return base.DeploymentCount(maxCpu)
	}
	if maxMem <= maxCpu && maxMem <= maxNet {
		return base.DeploymentCount(maxMem)
	}
	if maxNet <= maxCpu && maxNet <= maxMem {
		return base.DeploymentCount(maxNet)
	}
	return 0
}


func findHostWithResources(ns base.AppNeeds, app base.AppName, sortedHosts []base.HostId, goodHosts map[base.HostId]bool) base.HostId{
	var backUpHost base.HostId = ""
	PlannerLogger.Infof("findHostWithResources for app %s. goodHosts=%+v; sortedHosts=%+v", app, goodHosts, sortedHosts)
	for host := range goodHosts {
		if state_cloud.GlobalAvailableInstances.HostHasResourcesForApp(host, ns) {
			PlannerLogger.Infof("Found suitable host '%s'. It already has app '%s' installed", host, app)
			return host
		}
	}

	for _, hostId := range sortedHosts {
		if state_cloud.GlobalAvailableInstances.HostHasResourcesForApp(hostId, ns) {
			PlannerLogger.Infof("Found suitable host '%s'", hostId)
			return hostId
		} else {
			PlannerLogger.Infof("Host '%s' has insufficient resources foor app %s", hostId, app)
		}
	}
	PlannerLogger.Infof("Found no suitable host which has app '%s' installed. Returning backup host '%s'", app, backUpHost)
	return backUpHost
}

var TotalIter int = 0

func findHttpHostWithResources(ns base.AppNeeds, app base.AppName, sortedHosts []base.HostId, goodHosts map[base.HostId]bool) base.HostId {
	var backUpHost base.HostId = ""
	PlannerLogger.Infof("findHttpHostWithResources for app %s. goodHosts=%+v; sortedHosts=%+v", app, goodHosts, sortedHosts)
	for host := range goodHosts {
		if state_cloud.GlobalCloudLayout.Desired.HostHasApp(host, app) {
			continue
		}
		if state_cloud.GlobalAvailableInstances.HostHasResourcesForApp(host, ns) {
			PlannerLogger.Infof("Found suitable host '%s'", host)
			return host
		}
	}

	for _, hostId := range sortedHosts {
		TotalIter += 1
		if state_cloud.GlobalCloudLayout.Desired.HostHasApp(hostId, app) {
			continue
		}
		if state_cloud.GlobalAvailableInstances.HostHasResourcesForApp(hostId, ns) {
			PlannerLogger.Infof("Found suitable host '%s'", hostId)
			return hostId
		} else {
			PlannerLogger.Infof("Host '%s' has insufficient resources foor app %s", hostId, app)
		}
	}
	return backUpHost
}


func wipeDesired() {
	PlannerLogger.Info("Wiping Desired layout")
	FailedAssigned = []FailedAssign{}
	MissingAssigned = []MissingAssign{}
	state_cloud.GlobalCloudLayout.Desired.Wipe()
	for hostId := range state_cloud.GlobalAvailableInstances {
		state_cloud.GlobalCloudLayout.Desired.AddEmptyHost(hostId)
	}
	state_cloud.GlobalAvailableInstances.WipeUsage()
}

func forceAssignToHost(hostId base.HostId, app base.AppConfiguration, count base.DeploymentCount) {

}

func addFailedAssign(host base.HostId, name base.AppName, version base.Version, count base.DeploymentCount) {
	failedAssignMutex.Lock()
	FailedAssigned = append(FailedAssigned, FailedAssign{
		TargetHost: host, AppName: name, AppVersion: version, DeploymentCount: count,
	})
	failedAssignMutex.Unlock()
}

func addMissingAssign(name base.AppName, version base.Version, ty base.AppType, count base.DeploymentCount) {
	missingssignMutex.Lock()
	MissingAssigned = append(MissingAssigned, MissingAssign {
		AppName: name, AppVersion: version, AppType: ty, DeploymentCount: count,
	})
	missingssignMutex.Unlock()
}

func assignAppToHost(hostId base.HostId, app base.AppConfiguration, count base.DeploymentCount) bool {
	PlannerLogger.Infof("Assign %s:%d to host '%s' %d times", app.Name, app.Version, hostId, count)
	ns, err := state_needs.GlobalAppsNeedState.Get(app.Name, app.Version)
	if err != nil {
		PlannerLogger.Warnf("App %s:%d on host '%s': GlobalAppsNeedState.Get failed", app.Name, app.Version, hostId)
		addFailedAssign(hostId, app.Name, app.Version, count)
		return false
	}
	deployedNeeds := base.AppNeeds{
		CpuNeeds: base.CpuNeeds(int(ns.CpuNeeds) * int(count)),
		MemoryNeeds: base.MemoryNeeds(int(ns.MemoryNeeds) * int(count)),
		NetworkNeeds: base.NetworkNeeds(int(ns.NetworkNeeds) * int(count)),
	}
	if !state_cloud.GlobalAvailableInstances.HostHasResourcesForApp(hostId, ns) {
		db.Audit.Insert__AuditEvent(db.AuditEvent{Details:map[string]string{
			"message": fmt.Sprintf("App %s:%d on host '%s': Instance resources are insufficient, needed: %+v", app.Name, app.Version, hostId, ns),
			"subsystem": "planner",
			"application": string(app.Name),
			"application.version": string(app.Version),
			"host": string(hostId),
			"level": "warning",
		}})

		addFailedAssign(hostId, app.Name, app.Version, count)
		return false
	}
	updateInstanceResources(hostId, deployedNeeds)
	state_cloud.GlobalCloudLayout.Desired.AddApp(hostId, app.Name, app.Version, count)

	db.Audit.Insert__AuditEvent(db.AuditEvent{Details:map[string]string{
		"message": fmt.Sprintf("Assign %s:%d to host '%s' %d times successful", app.Name, app.Version, hostId, count),
		"subsystem": "planner",
		"application": string(app.Name),
		"application.version": string(app.Version),
		"host": string(hostId),
		"level": "info",
	}})

	return true
}

func updateInstanceResources(hostId base.HostId, needs base.AppNeeds)  {
	current, err := state_cloud.GlobalAvailableInstances.GetResources(hostId)
	if err != nil {
		return
	}
	current.UsedCpuResource += base.CpuResource(needs.CpuNeeds)
	current.UsedMemoryResource += base.MemoryResource(needs.MemoryNeeds)
	current.UsedNetworkResource += base.NetworkResource(needs.NetworkNeeds)
	state_cloud.GlobalAvailableInstances.Update(hostId, current)
}


func getGlobalResources() (base.CpuResource, base.MemoryResource, base.NetworkResource) {
	var totalCpuResources base.CpuResource
	var totalMemoryResources base.MemoryResource
	var totalNetworkResources base.NetworkResource

	for _, resources := range state_cloud.GlobalAvailableInstances {
		totalCpuResources += resources.TotalCpuResource
		totalMemoryResources += resources.TotalMemoryResource
		totalNetworkResources+= resources.TotalNetworkResource
	}
	PlannerLogger.Infof("Total available resources: Cpu: %d, Memory: %d, Network: %d", totalCpuResources, totalMemoryResources, totalNetworkResources)
	return totalCpuResources, totalMemoryResources, totalNetworkResources
}


func getGlobalMinNeeds() (base.CpuNeeds, base.MemoryNeeds, base.NetworkNeeds){
	var totalCpuNeeds base.CpuNeeds
	var totalMemoryNeeds base.MemoryNeeds
	var totalNetworkNeeds base.NetworkNeeds

	for appName, appObj := range state_configuration.GlobalConfigurationState.Apps {
		version := appObj.LatestVersion()
		appNeeds , err := state_needs.GlobalAppsNeedState.Get(appName, version)
		if err != nil {
			PlannerLogger.Warnf("Missing needs for app %s:%d", appName, version)
			continue
		}
		elem := appObj[version]
		cpu := int(elem.GetDeploymentCount()) * int(appNeeds.CpuNeeds)
		mem := int(elem.GetDeploymentCount()) * int(appNeeds.MemoryNeeds)
		net := int(elem.GetDeploymentCount()) * int(appNeeds.NetworkNeeds)
		PlannerLogger.Infof("AppMinNeeds for %s:%d: Cpu=%d, Memory=%d, Network=%d", appName, version, cpu, mem, net)
		totalCpuNeeds += base.CpuNeeds(cpu)
		totalMemoryNeeds += base.MemoryNeeds(mem)
		totalNetworkNeeds += base.NetworkNeeds(net)
	}
	PlannerLogger.Infof("GlobalAppMinNeeds: Cpu=%d, Memory=%d, Network=%d", totalCpuNeeds, totalMemoryNeeds, totalNetworkNeeds)
	return totalCpuNeeds, totalMemoryNeeds, totalNetworkNeeds
}


func getGlobalCurrentNeeds() (base.CpuNeeds, base.MemoryNeeds, base.NetworkNeeds) {
	var totalCpuNeeds base.CpuNeeds
	var totalMemoryNeeds base.MemoryNeeds
	var totalNetworkNeeds base.NetworkNeeds

	for hostId, hostObj := range state_cloud.GlobalCloudLayout.Current.Layout {
		for appName, appObj := range hostObj.Apps {
			appNeeds , err := state_needs.GlobalAppsNeedState.Get(appName, appObj.Version)
			if err != nil {
				PlannerLogger.Warnf("Missing needs for app %s:%d", appName, appNeeds)
				continue
			}
			cpu := int(appObj.DeploymentCount) * int(appNeeds.CpuNeeds)
			mem := int(appObj.DeploymentCount) * int(appNeeds.MemoryNeeds)
			net := int(appObj.DeploymentCount) * int(appNeeds.NetworkNeeds)
			PlannerLogger.Infof("AppNeeds on host '%s': App %s:%d deployed %d times: Cpu=%d, Memory=%d, Network=%d", hostId, appName, appObj.Version, appObj.DeploymentCount, cpu, mem, net)
			totalCpuNeeds += base.CpuNeeds(cpu)
			totalMemoryNeeds += base.MemoryNeeds(mem)
			totalNetworkNeeds += base.NetworkNeeds(net)
		}
	}
	PlannerLogger.Infof("GlobalAppCurrentNeeds: Cpu=%d, Memory=%d, Network=%d", totalCpuNeeds, totalMemoryNeeds, totalNetworkNeeds)
	return totalCpuNeeds, totalMemoryNeeds, totalNetworkNeeds
}



