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
	"time"
	"gatoor/orca/base"
	"sync"
	"errors"
	"sort"
)

var GlobalAppsStatusTracker AppsStatusTracker

func init() {
	GlobalAppsStatusTracker = AppsStatusTracker{}
}

type AppCrash struct {
	HostId base.HostId
	AppName base.AppName
	AppVersion base.Version
	Time time.Time
	Cause string
}

type VersionRating string
type CrashCount int
type RunningCount int

type AppTrack struct {
	Rating VersionRating
	CrashDetails []AppCrash
	RunningCount RunningCount
}

type AppsStatusTracker map[base.AppName]map[base.Version]AppTrack



var appsTrackerMutex = &sync.Mutex{}

const (
	APP_EVENT_SUCCESSFUL_UPDATE = "APP_EVENT_SUCCESSFUL_UPDATE"
	APP_EVENT_ROLLBACK = "APP_EVENT_ROLLBACK"
	APP_EVENT_CRASH = "APP_EVENT_CRASH"
	APP_EVENT_CHECKIN = "APP_EVENT_CHECKIN"

	RATING_STABLE = "RATING_STABLE"
	RATING_CRASHED = "RATING_CRASHED"
)

func (a *AppsStatusTracker) Init() {
	(*a) = make(map[base.AppName]map[base.Version]AppTrack)
}

func (a *AppsStatusTracker) GetRating(app base.AppName, version base.Version) (VersionRating, error) {
	appsTrackerMutex.Lock()
	defer appsTrackerMutex.Unlock()
	if elem, exists := (*a)[app]; exists {
		if e, ex := elem[version]; ex {
			return e.Rating, nil
		}
	}
	return "", errors.New("App does not exist")
}

func (a *AppsStatusTracker) Get(app base.AppName, version base.Version) (AppTrack, error) {
	appsTrackerMutex.Lock()
	defer appsTrackerMutex.Unlock()
	if elem, exists := (*a)[app]; exists {
		if e, ex := elem[version]; ex {
			return e, nil
		}
	}
	return AppTrack{}, errors.New("App does not exist")
}

func (a *AppsStatusTracker) Update(hostId base.HostId, app base.AppName, version base.Version, event string) {
	appsTrackerMutex.Lock()
	defer appsTrackerMutex.Unlock()
	newElem := AppTrack{}

	crashDetail := AppCrash{HostId: hostId, AppName: app, AppVersion: version, Time: time.Now().UTC(), Cause: event,}
	if elem, exists := (*a)[app]; exists {
		if v, ex := elem[version]; ex {
			if event == APP_EVENT_CRASH || event == APP_EVENT_ROLLBACK {
				v.Rating = RATING_CRASHED
				v.CrashDetails = append(v.CrashDetails, crashDetail)

			} else {
				v.RunningCount++
			}
			(*a)[app][version] = v
			return
		}
		if event == APP_EVENT_CRASH || event == APP_EVENT_ROLLBACK {
			newElem.Rating = RATING_CRASHED
			newElem.CrashDetails = []AppCrash{crashDetail}
		} else {
			newElem.Rating = RATING_STABLE
			newElem.RunningCount++
		}
		(*a)[app][version] = newElem
		return
	}
	if event == APP_EVENT_CRASH || event == APP_EVENT_ROLLBACK {
		newElem.Rating = RATING_CRASHED
		newElem.CrashDetails = []AppCrash{crashDetail}
	} else {
		newElem.Rating = RATING_STABLE
		newElem.RunningCount++
	}
	(*a)[app] = make(map[base.Version]AppTrack)
	(*a)[app][version] = newElem
}

func (a *AppsStatusTracker) UpdateAll(hostInfo base.HostInfo, time time.Time) {
	for _, appObj := range hostInfo.Apps {
		if appObj.Status == base.STATUS_RUNNING {
			a.Update(hostInfo.HostId, appObj.Name, appObj.Version, APP_EVENT_CHECKIN)
		} else if appObj.Status == base.STATUS_DEAD {
			a.Update(hostInfo.HostId, appObj.Name, appObj.Version, APP_EVENT_CRASH)
		}
	}
}

func (a *AppsStatusTracker) LastStable(app base.AppName) base.Version {
	appsTrackerMutex.Lock()
	defer appsTrackerMutex.Unlock()
	if _, exists := (*a)[app]; exists {
		var versions base.Versions
		for version := range (*a)[app] {
			versions = append(versions, version)
		}
		sort.Sort(sort.Reverse(versions))
		for _, ver := range versions {
			if (*a)[app][ver].Rating == RATING_STABLE {
				return ver
			}
		}
	}
	return 1
}
