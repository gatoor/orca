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
	"gatoor/orca/rewriteTrainer/state/cloud"
	"gatoor/orca/rewriteTrainer/state/configuration"
	"gatoor/orca/base"
	"sort"
	"gatoor/orca/rewriteTrainer/state/needs"
)


/*
	Sort the available hosts of the cluster...
	TRY_TO_REMOVE_HOSTS=true: Get the host with the least amount of available resources - this allows us to shut an instance down.
	TRY_TO_REMOVE_HOSTS=false: Get the host with a lot of spare resources - balance the load on the cluster
*/

type HostPair struct {
	Key base.HostId
	Val base.InstanceResources
}

type HostList []HostPair
type HostLess func(p HostList, i, j int) bool


func (p HostList) Len() int { return len(p) }
func (p HostList) Less(i, j int) bool { return (
	int(p[i].Val.TotalCpuResource - p[i].Val.UsedCpuResource) +
		int(p[i].Val.TotalMemoryResource - p[i].Val.UsedMemoryResource) +
		int(p[i].Val.TotalNetworkResource - p[i].Val.UsedNetworkResource)) < (
	int(p[j].Val.TotalCpuResource - p[j].Val.UsedCpuResource) +
		int(p[j].Val.TotalMemoryResource - p[j].Val.UsedMemoryResource) +
		int(p[j].Val.TotalNetworkResource - p[j].Val.UsedNetworkResource))}
func (p HostList) Swap(i, j int){ p[i], p[j] = p[j], p[i] }


var HOST_UPPER_WATERMARK float32 = 0.2

func checkWatermark(resources base.InstanceResources, hostId base.HostId) bool{
	if float32(resources.TotalCpuResource - resources.UsedCpuResource) / float32(resources.TotalCpuResource) <= HOST_UPPER_WATERMARK {
		PlannerLogger.Debugf("Hit Cpu Watermark for host '%s'", hostId)
		return false
	}
	if float32(resources.TotalMemoryResource - resources.UsedMemoryResource) / float32(resources.TotalMemoryResource) <= HOST_UPPER_WATERMARK {
		PlannerLogger.Debugf("Hit Memory Watermark for host '%s'", hostId)
		return false
	}
	if float32(resources.TotalNetworkResource - resources.UsedNetworkResource) / float32(resources.TotalNetworkResource) <= HOST_UPPER_WATERMARK {
		PlannerLogger.Debugf("Hit Network Watermark for host '%s'", hostId)
		return false
	}
	return true
}

func sortByAvailableResources() []base.HostId {
	//PlannerLogger.Warn("TODO IMPLEMENT SORT BY AVAILABLE RESOURCERS BASED ON NEEDS TYPE")
	hostList := make(HostList, len(state_cloud.GlobalAvailableInstances))
	c:= 0
	for hostId, hostResources := range state_cloud.GlobalAvailableInstances {
		if checkWatermark(hostResources,hostId) {
			hostList[c] = HostPair{
				hostId, hostResources,
			}
			c++
		}
	}
	hostList = hostList[0:c]

	if !state_configuration.GlobalConfigurationState.Trainer.Policies.TRY_TO_REMOVE_HOSTS {
		sort.Sort(sort.Reverse(hostList))
	} else {
		sort.Sort(hostList)
	}

	sorted := make([]base.HostId, len(hostList))
	for i := 0; i < len(hostList); i++ {
		sorted[i] = hostList[i].Key
	}
	PlannerLogger.Infof("Sorted Hosts with TRY_TO_REMOVE_HOSTS=%t: %+v", state_configuration.GlobalConfigurationState.Trainer.Policies.TRY_TO_REMOVE_HOSTS, sorted)
	return sorted
}


/*
	Sort the apps by Needs: High needs => should be planned first
 */


type ConfPair struct {
	Key base.AppName
	Val base.AppNeeds
}

type ConfList []ConfPair

func (p ConfList) Len() int { return len(p) }
func (p ConfList) Less(i, j int) bool { return (int(p[i].Val.CpuNeeds) + int(p[i].Val.MemoryNeeds) + int(p[i].Val.NetworkNeeds)) > (int(p[j].Val.CpuNeeds) + int(p[j].Val.MemoryNeeds) + int(p[j].Val.NetworkNeeds))}
func (p ConfList) Swap(i, j int){ p[i], p[j] = p[j], p[i] }


func sortByTotalNeeds(apps map[base.AppName]base.AppConfiguration) []base.AppName {
	sorted := make([]base.AppName, len(apps), len(apps))
	appList := make(ConfList, len(apps), len(apps))
	c := 0

	for appName, appObj := range apps {
		needs, _ := state_needs.GlobalAppsNeedState.Get(appName, appObj.Version)
		appList[c] = ConfPair{
			appName, needs,
		}
		c++
	}

	sort.Sort(appList)

	for i := 0; i < len(apps); i++ {
		sorted[i] = appList[i].Key
	}

	return sorted
}