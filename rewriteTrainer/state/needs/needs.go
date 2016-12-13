package state_needs

import (
	"gatoor/orca/base"
	"sync"
	"errors"
	Logger "gatoor/orca/rewriteTrainer/log"
	"gatoor/orca/rewriteTrainer/needs"
	"sort"
)

var needsStateMutex = &sync.Mutex{}
var GlobalAppsNeedState AppsNeedState = make(map[base.AppName]AppNeedVersion)
var StateNeedsLogger = Logger.LoggerWithField(Logger.Logger, "module", "state_needs")

type AppsNeedState map[base.AppName]AppNeedVersion

//TODO use WeeklyNeeds here
type AppNeedVersion map[base.Version]needs.AppNeeds

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
func (a AppsNeedState) Get(app base.AppName, version base.Version) (needs.AppNeeds, error) {
	needsStateMutex.Lock()
	defer needsStateMutex.Unlock()
	if _, exists := a[app]; !exists {
		StateNeedsLogger.Warnf("App '%s' does not exist", app)
		return needs.AppNeeds{}, errors.New("No such App")
	}
	if _, exists := a[app][version]; !exists {
		StateNeedsLogger.Warnf("App '%s' does not exist", app)
		return needs.AppNeeds{}, errors.New("No such Version")
	}
	res := a[app][version]
	StateNeedsLogger.Debugf("Get for %s:%d: %+v", app, version, res)
	return res, nil
}

func (a AppsNeedState) lastValidNeeds(app base.AppName) needs.AppNeeds {
	if _, exists := a[app]; !exists {
		return needs.AppNeeds{1, 1, 1}
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
	return needs.AppNeeds{1, 1, 1}
}

//TODO use WeeklyNeeds
func (a AppsNeedState) UpdateNeeds(app base.AppName, version base.Version, ns needs.AppNeeds) {
	needsStateMutex.Lock()
	defer needsStateMutex.Unlock()
	if _, exists := a[app]; !exists {
		a[app] = make(map[base.Version]needs.AppNeeds)
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




