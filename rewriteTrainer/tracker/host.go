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

package tracker

import (
	"gatoor/orca/base"
	Logger "gatoor/orca/rewriteTrainer/log"
	"time"
	"sync"
	"errors"
	"gatoor/orca/rewriteTrainer/cloud"
	"gatoor/orca/rewriteTrainer/state/cloud"
)

var GlobalHostTracker HostTracker
var GlobalHostCrashHandler HostCrashHandler

var TrackerLogger = Logger.LoggerWithField(Logger.Logger, "module", "tracker")

func init() {
	GlobalHostTracker = HostTracker{}
	GlobalHostCrashHandler = HostCrashHandler{}
}

var hostTrackerMutex = &sync.Mutex{}
var hostCrashHandlerMutex = &sync.Mutex{}

type HostCrashReason string
type HostSpawnStatus string

//type HostSpawnQueue map[]

type HostSpawn struct {
	InitiatedTime time.Time
	NewHostId base.HostId
	OldHostId base.HostId
	Status HostSpawnStatus
}

type HostCrashHandler map[base.HostId]HostSpawn

const (
	HOST_CHECKIN_TIMEOUT = time.Duration(time.Minute * 15)
	HOST_CRASH_TIMEOUT = "HOST_CRASH_TIMEOUT"
	HOST_CRASH_CLOUD_PROVIDER_KILL = "HOST_CRASH_CLOUD_PROVIDER_KILL"

	HOST_STATUS_SPAWN_TRIGGERED = "HOST_STATUS_SPAWN_TRIGGERED"
)

type HostCrash struct {
	Time time.Time
	Reason HostCrashReason
}

type HostTrackingInfo struct {
	LastCheckin time.Time
	//Crashes []HostCrash
	//LastConfig responder.PushConfiguration
}


type HostTracker map[base.HostId]HostTrackingInfo

func (h *HostTracker) Init() {
	(*h) = make(map[base.HostId]HostTrackingInfo)
}

func (h *HostTracker) Update(hostId base.HostId, checkin time.Time) {
	TrackerLogger.Infof("Checking in host %s at %s", hostId, checkin.Format(time.RFC3339Nano))
	hostTrackerMutex.Lock()
	defer hostTrackerMutex.Unlock()
	if _, exists := (*h)[hostId]; !exists {
		(*h)[hostId] = HostTrackingInfo{}
	}
	elem := (*h)[hostId]
	elem.LastCheckin = checkin
	(*h)[hostId] = elem
	addAvailableInstance(hostId)
}

func addAvailableInstance(hostId base.HostId) {
	resources := cloud.CurrentProvider.GetResources(cloud.CurrentProvider.GetInstanceType(hostId))
	state_cloud.GlobalAvailableInstances.Update(hostId, resources)
}

func (h *HostTracker) Get(hostId base.HostId) (HostTrackingInfo, error) {
	hostTrackerMutex.Lock()
	defer hostTrackerMutex.Unlock()
	if val, exists := (*h)[hostId]; exists {
		return val, nil
	}
	return HostTrackingInfo{}, errors.New("Host does not exist")
}

func (h *HostTracker) Remove(hostId base.HostId) {
	TrackerLogger.Infof("Removing host %s from HostTracker", hostId)
	//hostTrackerMutex.Lock()
	//defer hostTrackerMutex.Unlock()
	delete(*h, hostId)
}

func (h *HostTracker) CheckCheckinTimeout() {
	TrackerLogger.Info("Checking for timed out hosts")
	hostTrackerMutex.Lock()
	defer hostTrackerMutex.Unlock()
	for hostId, hostInfo := range (*h) {
		if hostInfo.LastCheckin.Before(time.Now().UTC().Add(-HOST_CHECKIN_TIMEOUT)) {
			TrackerLogger.Warnf("Host '%s' checking timed out, last checkin was at '%s'", hostId, hostInfo.LastCheckin)
			GlobalHostCrashHandler.spawnHost(hostId)
			removeHostFromState(hostId)
		}
	}

}

func (h *HostTracker) CheckCloudProvider() {
	hostTrackerMutex.Lock()
	hosts := (*h)
	hostTrackerMutex.Unlock()
	for hostId := range hosts {
		if cloud.CurrentProvider.CheckInstance(hostId) == cloud.INSTANCE_STATUS_DEAD {
			GlobalHostCrashHandler.spawnHost(hostId)
		}
	}
}

func (h *HostTracker) HandleCloudProviderEvent(providerEvent cloud.ProviderEvent) {
	TrackerLogger.Infof("Got Provicer event %+v", providerEvent)
	elem, err := GlobalHostCrashHandler.Get(providerEvent.HostId)
	if err != nil {
		if providerEvent.Type == cloud.PROVIDER_EVENT_KILLED {
			TrackerLogger.Warnf("Cloud Provider killed host '%s', spawning new host", providerEvent.HostId)
			GlobalHostCrashHandler.spawnHost(providerEvent.HostId)
			removeHostFromState(providerEvent.HostId)
		}
	}
	TrackerLogger.Info(providerEvent)
	if providerEvent.Type == cloud.PROVIDER_EVENT_READY {
		GlobalHostCrashHandler.checkinHost(providerEvent.HostId)
	} else {
		TrackerLogger.Warnf("CloudProvider Event for already handled host: %+v", elem)
	}
}

func removeHostFromState(hostId base.HostId) {
	TrackerLogger.Infof("Host %s has died. Removing it from all state objects", hostId)
	state_cloud.GlobalAvailableInstances.Remove(hostId)
	state_cloud.GlobalCloudLayout.Current.RemoveHost(hostId)
	GlobalHostTracker.Remove(hostId)
}


func (h *HostCrashHandler) spawnHost(hostId base.HostId) {
	hostCrashHandlerMutex.Lock()
	defer hostCrashHandlerMutex.Unlock()
	if _, exists := (*h)[hostId]; !exists {
		now := time.Now().UTC()
		TrackerLogger.Warnf("Triggered Host spawn to replace host '%s' at '%s'", hostId, now.Format(time.RFC3339Nano))
		newId := cloud.CurrentProvider.SpawnInstanceLike(hostId)
		TrackerLogger.Warnf("Got replacement host %s for old host %s", newId, hostId)
		(*h)[hostId] = HostSpawn{OldHostId: hostId, NewHostId: newId, InitiatedTime: now, Status: HOST_STATUS_SPAWN_TRIGGERED}
	}
}

func (h *HostCrashHandler) Get(hostId base.HostId) (HostSpawn, error) {
	hostCrashHandlerMutex.Lock()
	defer hostCrashHandlerMutex.Unlock()
	if elem, exists := (*h)[hostId]; exists {
		return elem, nil
	}
	return HostSpawn{}, errors.New("No host")
}

func (h *HostCrashHandler) checkinHost(hostId base.HostId) {
	hostCrashHandlerMutex.Lock()
	defer hostCrashHandlerMutex.Unlock()
	var delHost base.HostId
	for oldHost, obj := range (*h) {
		if obj.NewHostId == hostId {
			TrackerLogger.Infof("New host '%s' checked in. It replaced host '%s', spawning was triggered at %s", hostId, obj.OldHostId, obj.InitiatedTime)
			delHost = oldHost
		}
	}
	delete((*h), delHost)
}



//scheduled tassk: check last checkin of hosts + check cloudprovider host status -> spawn new instance