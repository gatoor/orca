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

package cloud

import (
	"gatoor/orca/base"
	"gatoor/orca/rewriteTrainer/state/cloud"
	"strings"
	"strconv"
	"gatoor/orca/rewriteTrainer/state/configuration"
)

const (
	PROVIDER_TEST = "TEST"
	PROVIDER_AWS = "AWS"

	INSTANCE_STATUS_SPAWN_TRIGGERED = "INSTANCE_STATUS_SPAWN_TRIGGERED"
	INSTANCE_STATUS_SPAWNING = "INSTANCE_STATUS_SPAWNING"
	INSTANCE_STATUS_HEALTHY = "INSTANCE_STATUS_HEALTHY"
	INSTANCE_STATUS_DEAD = "INSTANCE_DEAD"

	PROVIDER_EVENT_KILLED = "PROVIDER_EVENT_KILLED"
	PROVIDER_EVENT_READY = "PROVIDER_EVENT_READY"
)

type ProviderEventType string
type InstanceStatus string

type ProviderEvent struct {
	HostId base.HostId
	Type ProviderEventType
}

type Provider interface {
	Init()
	SpawnInstances([]base.InstanceType) bool
	SpawnInstance(base.InstanceType) base.HostId
	SpawnInstanceSync(base.InstanceType) base.HostId
	SpawnInstanceLike(base.HostId) base.HostId
	GetIp(base.HostId) base.IpAddr
	GetResources(base.InstanceType) base.InstanceResources
	GetInstanceType(base.HostId) base.InstanceType
	SuitableInstanceTypes(base.InstanceResources) []base.InstanceType
	CheckInstance(base.HostId) InstanceStatus
	TerminateInstance(base.HostId) bool
	GetSpawnLog() []base.HostId
	RemoveFromSpawnLog(base.HostId)

	UpdateLoadBalancers(hostId base.HostId, app base.AppName, version base.Version, event string)
}

var CurrentProvider Provider

func Init() {
	AWSLogger.Infof("Initializing CloudProvider of type %s", state_configuration.GlobalConfigurationState.CloudProvider.Type)
	if state_configuration.GlobalConfigurationState.CloudProvider.Type == PROVIDER_AWS {
		CurrentProvider = &AWSProvider{ProviderConfiguration: state_configuration.GlobalConfigurationState.CloudProvider}
	} else {
		CurrentProvider = &TestProvider{ProviderConfiguration: state_configuration.GlobalConfigurationState.CloudProvider}
	}
	CurrentProvider.Init()

	spawnToMinInstances()
	AWSLogger.Infof("Instances: %+v", state_cloud.GlobalAvailableInstances)
}

func spawnToMinInstances() {
	if len(state_cloud.GlobalAvailableInstances) < int(state_configuration.GlobalConfigurationState.CloudProvider.MinInstances) {
		AWSLogger.Infof("Not enough instances available. Spawning more, available:%d min:%d", len(state_cloud.GlobalAvailableInstances), state_configuration.GlobalConfigurationState.CloudProvider.MinInstances)
		for i := len(state_cloud.GlobalAvailableInstances); i < int(state_configuration.GlobalConfigurationState.CloudProvider.MinInstances); i++ {
			CurrentProvider.SpawnInstanceSync(state_configuration.GlobalConfigurationState.CloudProvider.BaseInstanceType)
		}
	} else {
		AWSLogger.Infof("Enough instances available, going on")
	}
}


type TestProvider struct {
	ProviderConfiguration base.ProviderConfiguration

	Type base.ProviderType
	InstanceTypes []base.InstanceType
	SpawnList []base.HostId
	KillList []base.HostId
}

func (a *TestProvider) Init() {
	a.Type = PROVIDER_TEST
	a.InstanceTypes = []base.InstanceType{"test", "otherstuff"}
}

func (a *TestProvider) UpdateLoadBalancers(hostId base.HostId, app base.AppName, version base.Version, event string){

}

func (a *TestProvider) GetResources(ty base.InstanceType) base.InstanceResources {
	if strings.Contains(string(ty), "host2") {
		return base.InstanceResources{UsedMemoryResource:0, UsedCpuResource:0, UsedNetworkResource:0, TotalMemoryResource: 2000, TotalNetworkResource:2000, TotalCpuResource:2000}
	}
	if !strings.Contains(string(ty), "metrics=") {
		return base.InstanceResources{UsedMemoryResource:0, UsedCpuResource:0, UsedNetworkResource:0, TotalMemoryResource: 1000, TotalNetworkResource:1000, TotalCpuResource:1000}
	}
	arr := strings.Split(strings.Split(string(ty), "metrics=")[1], "_")
	cpu, _ := strconv.Atoi(arr[0])
	mem, _ := strconv.Atoi(arr[1])
	net, _ := strconv.Atoi(arr[2])

	return base.InstanceResources{UsedMemoryResource:0, UsedCpuResource:0, UsedNetworkResource:0, TotalMemoryResource: base.MemoryResource(mem), TotalNetworkResource:base.NetworkResource(net), TotalCpuResource:base.CpuResource(cpu)}
}

func (a *TestProvider) SpawnInstance(ty base.InstanceType) base.HostId {
	AWSLogger.Infof("Trying to spawn a single instance of type '%s'", ty)
	AWSLogger.Errorf("NOT IMPLEMENTED TEST CLOUD PROVIDER SpawnInstance")
	a.SpawnList = append(a.SpawnList, base.HostId(string(ty)))
	return "TODO"
}

func (a *TestProvider) GetSpawnLog() []base.HostId{
	AWSLogger.Errorf("NOT IMPLEMENTED TEST CLOUD PROVIDER GetSpawnLog")
	return a.SpawnList
}

func (a *TestProvider) RemoveFromSpawnLog(base.HostId) {
	AWSLogger.Errorf("NOT IMPLEMENTED TEST CLOUD PROVIDER RemoveFromSpawnLog")
}

func (a *TestProvider) TerminateInstance(hostId base.HostId) bool{
	AWSLogger.Errorf("NOT IMPLEMENTED TEST CLOUD PROVIDER TerminateInstance")
	a.KillList = append(a.KillList, hostId)
	return true
}

func (a *TestProvider) GetInstanceType(hostId base.HostId) base.InstanceType{
	return base.InstanceType(hostId)
}

func (a *TestProvider) SpawnInstanceSync(ty base.InstanceType) base.HostId {
	AWSLogger.Infof("Trying to spawn a single instance of type '%s'", ty)
	AWSLogger.Errorf("NOT IMPLEMENTED TEST CLOUD PROVIDER SpawnInstanceSync")
	return ""
}

func (a *TestProvider) SpawnInstanceLike(hostId base.HostId) base.HostId{
	AWSLogger.Errorf("NOT IMPLEMENTED TEST CLOUD PROVIDER SpawnInstanceLike")
	//a.SpawnList = append(a.SpawnList, "new_" + hostId)
	return "new_" + hostId
}

func (a *TestProvider) SpawnInstances(tys []base.InstanceType) bool{
	AWSLogger.Infof("Trying to spawn %d instances", len(tys))
	AWSLogger.Errorf("NOT IMPLEMENTED TEST CLOUD PROVIDER SpawnInstances")
	return true
}

func (a *TestProvider) GetIp(hostId base.HostId) base.IpAddr {
	return ""
}

func (a *TestProvider) SuitableInstanceTypes(resources base.InstanceResources) []base.InstanceType {
	res := []base.InstanceType{}
	return res
}

func (a *TestProvider) CheckInstance(hostId base.HostId) InstanceStatus {
	if strings.Contains(string(hostId), "healthy") {
		return INSTANCE_STATUS_HEALTHY
	}
	if strings.Contains(string(hostId), "spawning") {
		return INSTANCE_STATUS_SPAWNING
	}
	if strings.Contains(string(hostId), "spawn_triggered") {
		return INSTANCE_STATUS_SPAWN_TRIGGERED
	}
	if strings.Contains(string(hostId), "dead") {
		return INSTANCE_STATUS_DEAD
	}
	return INSTANCE_STATUS_HEALTHY
}

