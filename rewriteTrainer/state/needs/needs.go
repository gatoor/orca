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

package state_needs

import (
	"gatoor/orca/base"
	"sync"
	"errors"
	Logger "gatoor/orca/rewriteTrainer/log"
	"sort"
)

var needsStateMutex = &sync.Mutex{}
var GlobalAppsNeedState AppsNeedState = make(map[base.AppName]AppNeedVersion)
var StateNeedsLogger = Logger.LoggerWithField(Logger.Logger, "module", "state_needs")

type AppsNeedState map[base.AppName]AppNeedVersion

//TODO use WeeklyNeeds here
type AppNeedVersion map[base.Version]base.AppNeeds

func (a AppsNeedState) GetAll(app base.AppName) (AppNeedVersion, error) {
	needsStateMutex.Lock()
	defer needsStateMutex.Unlock()
	if _, exists := a[app]; !exists {
		StateNeedsLogger.Warnf("App '%s' does not exist", app)
		return AppNeedVersion{}, errors.New("No such App")
	}
	res := a[app]
	StateNeedsLogger.Debugf("GetAll for '%s': %+v", app, res)
	return res, nil
}

//TODO get them by current time with WeeklyNeeds
func (a AppsNeedState) Get(app base.AppName, version base.Version) (base.AppNeeds, error) {
	needsStateMutex.Lock()
	defer needsStateMutex.Unlock()
	if _, exists := a[app]; !exists {
		StateNeedsLogger.Warnf("App '%s' does not exist", app)
		return base.AppNeeds{}, errors.New("No such App")
	}
	if _, exists := a[app][version]; !exists {
		StateNeedsLogger.Warnf("App '%s' does not exist", app)
		return base.AppNeeds{}, errors.New("No such Version")
	}
	res := a[app][version]
	StateNeedsLogger.Debugf("Get for %s:%d: %+v", app, version, res)
	return res, nil
}

func (a AppsNeedState) lastValidNeeds(app base.AppName) base.AppNeeds {
	if _, exists := a[app]; !exists {
		return base.AppNeeds{1, 1, 1}
	}
	var versions base.Versions
	for version := range a[app] {
		versions = append(versions, version)
	}
	sort.Sort(sort.Reverse(versions))

	for _, version := range versions {
		if a[app][version].CpuNeeds != 0 && a[app][version].MemoryNeeds != 0 && a[app][version].NetworkNeeds != 0 {
			return a[app][version]
		}
	}
	return base.AppNeeds{1, 1, 1}
}

//TODO use WeeklyNeeds
func (a AppsNeedState) UpdateNeeds(app base.AppName, version base.Version, ns base.AppNeeds) {
	needsStateMutex.Lock()
	defer needsStateMutex.Unlock()
	if _, exists := a[app]; !exists {
		a[app] = make(map[base.Version]base.AppNeeds)
	}
	if ns.CpuNeeds == 0 || ns.MemoryNeeds == 0 || ns.NetworkNeeds == 0 {
		ns = a.lastValidNeeds(app)
		StateNeedsLogger.Warnf("UpdateNeeds for %s:%d: Needs are 0, using last available needs %+v", app, version, ns)
	}
	StateNeedsLogger.Debugf("UpdateNeeds for %s:%d: %+v", app, version, ns)
	a[app][version] = ns
}

func (a AppsNeedState) Snapshot() AppsNeedState {
	needsStateMutex.Lock()
	defer needsStateMutex.Unlock()
	res := a
	return res
}




