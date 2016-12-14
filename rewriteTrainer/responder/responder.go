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

package responder

import (
	"gatoor/orca/base"
	"gatoor/orca/rewriteTrainer/planner"
	"errors"
	"gatoor/orca/rewriteTrainer/state/configuration"
	Logger "gatoor/orca/rewriteTrainer/log"
	"gatoor/orca/rewriteTrainer/tracker"
	"fmt"
	"gatoor/orca/rewriteTrainer/cloud"
	"gatoor/orca/rewriteTrainer/state/cloud"
	"gatoor/orca/rewriteTrainer/audit"
	"sort"
)

var ResponderLogger = Logger.LoggerWithField(Logger.Logger, "module", "responder")

func GetConfigForHost(hostId base.HostId) (base.PushConfiguration, error) {
	ResponderLogger.Infof("Getting config for host %s", hostId)
	appName, elem, err := getQueueElement(hostId)
	if err == nil {
		ResponderLogger.Infof("Got QueueElement for host %s: app %s", hostId, appName)
		config, err := state_configuration.GlobalConfigurationState.GetApp(appName, elem.Version.Version)
		if err != nil {
			ResponderLogger.Warnf("Getting config for host %s app %s failed %s", hostId, appName, err)
			return base.PushConfiguration{}, errors.New("GlobalConfiguraitonState does not have this app")
		}
		return base.PushConfiguration{DeploymentCount: elem.Version.DeploymentCount, AppConfiguration: config}, nil
	}
	ResponderLogger.Warnf("Getting config for host %s failed %s", hostId, err)
	return base.PushConfiguration{}, err
}

func getQueueElement(hostId base.HostId) (base.AppName, planner.AppsUpdateState, error) {
	ResponderLogger.Infof("Getting queue element for host %s", hostId)
	apps, err := planner.Queue.Get(hostId)
	if err == nil {

		var appNames []string
		for appName := range apps {
			appNames = append(appNames, string(appName))
		}
		sort.Strings(appNames)

		for _, appName := range appNames {
			if apps[base.AppName(appName)].State == planner.STATE_APPLYING {
				ResponderLogger.Infof("Got STATE_APPLYING queue element for host '%s' app '%s'", hostId, appName)
				return base.AppName(appName), apps[base.AppName(appName)], nil
			}
		}

		for _, appName := range appNames {
			if apps[base.AppName(appName)].State == planner.STATE_QUEUED {
				ResponderLogger.Infof("Got STATE_QUEUED queue element for host '%s' app '%s'", hostId, appName)
				planner.Queue.SetState(hostId, base.AppName(appName), planner.STATE_APPLYING)
				return base.AppName(appName), apps[base.AppName(appName)], nil
			}
		}
	}
	ResponderLogger.Infof("Get queue element for host '%s' failed: %s", hostId, err)
	return "", planner.AppsUpdateState{}, errors.New("No element for host ")
}

func CheckAppState(hostInfo base.HostInfo) {
	ResponderLogger.Infof("checking Apps state for host '%s'", hostInfo.HostId)
	queued, _ := planner.Queue.Get(hostInfo.HostId)

	if len(hostInfo.Apps) == 0 {
		handleEmptyHost(hostInfo)
		return
	}

	if !checkScaling(hostInfo, queued) {
		ResponderLogger.Infof("Host %s did not scale yet. Performing simple checks", hostInfo.HostId)
		for _, appObj := range hostInfo.Apps {
			if _, exists := queued[appObj.Name]; exists {
				simpleAppCheck(appObj, hostInfo.HostId)
			}
		}
		return
	}

	for _, appObj := range hostInfo.Apps {
		if queuedState, exists := queued[appObj.Name]; exists {
			//the app is queued for an update
			if queuedState.State == planner.STATE_QUEUED {
				//the update is waiting for other updates, perform simple check
				simpleAppCheck(appObj, hostInfo.HostId)
			} else {
				checkAppUpdate(appObj, hostInfo.HostId, queuedState)
			}
		} else {
			//no updates, just check if its running
			simpleAppCheck(appObj, hostInfo.HostId)
		}
	}
}

func checkScaling(hostInfo base.HostInfo, queued map[base.AppName]planner.AppsUpdateState) bool {
	appsCount := make(map[base.AppName]int)
	if len(queued) > 0 {
		for _, app := range hostInfo.Apps {
			appsCount[app.Name]++
		}
		ResponderLogger.Infof("Check scaling on host %s. Queued: %+v", hostInfo.HostId, queued)
		for appName, count := range appsCount {
			if queued[appName].State != planner.STATE_QUEUED && queued[appName].State != planner.UpdateState("") {
				if int(queued[appName].Version.DeploymentCount) != count {
					ResponderLogger.Warnf("Scaling up of app %s:%d on host %s is not done. Should be %d but is %d", appName, queued[appName].Version.Version, hostInfo.HostId, queued[appName].Version.DeploymentCount, count)
					return false
				} else {
					ResponderLogger.Infof("Scaling up of app %s:%d on host %s successful", appName, queued[appName].Version.Version, hostInfo.HostId)
				}
			}
		}
	}
	return true
}

func handleEmptyHost(hostInfo base.HostInfo) {
	isNew := false
	for _, hostId := range cloud.CurrentProvider.GetSpawnLog() {
		if hostId == hostInfo.HostId {
			isNew = true
			cloud.CurrentProvider.RemoveFromSpawnLog(hostId)
			state_cloud.GlobalAvailableInstances.Update(hostId, cloud.CurrentProvider.GetResources(cloud.CurrentProvider.GetInstanceType(hostId)))
		}
	}
	if !isNew {
		ResponderLogger.Warnf("TODO empty host %s, add logic to kill host if not needed", hostInfo.HostId)
		//cloud.CurrentProvider.TerminateInstance(hostInfo.HostId)
	}
}

func simpleAppCheck(appObj base.AppInfo, hostId base.HostId) {
	if appObj.Status != base.STATUS_RUNNING {
		audit.Audit.AddEvent(map[string]string{
			"message": fmt.Sprintf("App %s:%d on host '%s' is not running. Adding it to GlobalAppCrashes", appObj.Name, appObj.Version, hostId),
			"subsystem": "responder",
			"application": string(appObj.Name),
			"host": string(hostId),
			"level": "error",
		})

		tracker.GlobalAppsStatusTracker.Update(hostId, appObj.Name, appObj.Version, tracker.APP_EVENT_CRASH)
		cloud.CurrentProvider.UpdateLoadBalancers(hostId, appObj.Name, appObj.Version, base.STATUS_DEAD)

	} else {
		tracker.GlobalAppsStatusTracker.Update(hostId, appObj.Name, appObj.Version, tracker.APP_EVENT_CHECKIN)
	}
}

func checkAppUpdate(appObj base.AppInfo, hostId base.HostId, queuedState planner.AppsUpdateState) {
	ResponderLogger.Infof("Check update of App %s:%d on host '%s'", appObj.Name, appObj.Version, hostId)
	if queuedState.State != planner.STATE_APPLYING {
		audit.Audit.AddEvent(map[string]string{
			"message": fmt.Sprintf("Got illegal state %s for update of App %s:%d on host '%s'", queuedState.State, appObj.Name, appObj.Version, hostId),
			"subsystem": "responder",
			"application": string(appObj.Name),
			"host": string(hostId),
			"level": "error",
		})
		return
	}
	if appObj.Status == base.STATUS_RUNNING {
		if appObj.Version == queuedState.Version.Version {
			audit.Audit.AddEvent(map[string]string{
				"message": fmt.Sprintf("Update of App %s:%d on host '%s' successful", appObj.Name, appObj.Version, hostId),
				"subsystem": "responder",
				"application": string(appObj.Name),
				"host": string(hostId),
				"level": "info",
			})

			handleSuccessfulUpdate(hostId, appObj.Name, appObj.Version)
			cloud.CurrentProvider.UpdateLoadBalancers(hostId, appObj.Name, appObj.Version, base.STATUS_RUNNING)

		} else {
			audit.Audit.AddEvent(map[string]string{
				"message": fmt.Sprintf("Update of App %s:%d on host '%s' rolled back to version %s", appObj.Name, queuedState.Version.Version, hostId, appObj.Version),
				"subsystem": "responder",
				"application": string(appObj.Name),
				"host": string(hostId),
				"level": "error",
			})

			handleRollback(hostId, appObj.Name, queuedState.Version.Version)
			cloud.CurrentProvider.UpdateLoadBalancers(hostId, appObj.Name, appObj.Version, base.STATUS_RUNNING)
		}
	}
	if appObj.Status == base.STATUS_DEPLOYING {
		if appObj.Version == queuedState.Version.Version {
			audit.Audit.AddEvent(map[string]string{
				"message": fmt.Sprintf("Update of App %s:%d on host '%s' is still applying", appObj.Name, appObj.Version, hostId),
				"subsystem": "responder",
				"application": string(appObj.Name),
				"host": string(hostId),
				"level": "info",
			})

			return
		} else {
			audit.Audit.AddEvent(map[string]string{
				"message": fmt.Sprintf("Update of App %s:%d on host '%s' rolling back to version %d", appObj.Name, queuedState.Version, hostId, appObj.Version),
				"subsystem": "responder",
				"application": string(appObj.Name),
				"host": string(hostId),
				"level": "error",
			})

			handleRollback(hostId, appObj.Name, queuedState.Version.Version)
		}
	}
	if appObj.Status == base.STATUS_DEAD {
		audit.Audit.AddEvent(map[string]string{
			"message": fmt.Sprintf("Update of App %s:%d on host '%s' was fatal for the app, the version that died on the host is %d", appObj.Name, queuedState.Version, hostId, appObj.Version),
			"subsystem": "responder",
			"application": string(appObj.Name),
			"host": string(hostId),
			"level": "error",
		})

		handleFatalUpdate(hostId, appObj.Name, appObj.Version, queuedState.Version.DeploymentCount)
		cloud.CurrentProvider.UpdateLoadBalancers(hostId, appObj.Name, appObj.Version, base.STATUS_DEAD)
	}
}

func handleSuccessfulUpdate(hostId base.HostId, appName base.AppName, version base.Version) {
	tracker.GlobalAppsStatusTracker.Update("", appName, version, tracker.APP_EVENT_SUCCESSFUL_UPDATE)
	planner.Queue.Remove(hostId, appName)
}

func handleRollback(hostId base.HostId, appName base.AppName, failedVersion base.Version) {
	fmt.Println("ROLLBACK")
	fmt.Println(hostId)
	fmt.Println(appName)
	fmt.Println(failedVersion)
	fmt.Println("ROLLBACK")
	tracker.GlobalAppsStatusTracker.Update(hostId, appName, failedVersion, tracker.APP_EVENT_ROLLBACK)
	planner.Queue.SetState(hostId, appName, planner.STATE_SUCCESS)
	planner.Queue.RemoveApp(appName, failedVersion)
}

func handleFatalUpdate(hostId base.HostId, appName base.AppName, version base.Version, deploymentCount base.DeploymentCount) {
	tracker.GlobalAppsStatusTracker.Update(hostId, appName, version, tracker.APP_EVENT_CRASH)
	planner.Queue.RollbackApp(appName, version, tracker.GlobalAppsStatusTracker.LastStable(appName))
	planner.Queue.Add(hostId, appName, state_cloud.AppsVersion{DeploymentCount: deploymentCount, Version: tracker.GlobalAppsStatusTracker.LastStable(appName)})
}

