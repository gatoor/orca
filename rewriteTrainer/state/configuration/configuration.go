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

package state_configuration

import (
	"sync"
	"gatoor/orca/base"
	"errors"
	"sort"
	Logger "gatoor/orca/rewriteTrainer/log"
	"gatoor/orca/rewriteTrainer/state/needs"
)

var ConfigLogger = Logger.LoggerWithField(Logger.Logger, "module", "configuration")
var GlobalConfigurationState ConfigurationState

var configurationStateMutex = &sync.Mutex{}


type ConfigurationState struct {
	ConfigurationRootPath string

	Trainer base.TrainerConfigurationState
	Apps AppsConfigurationState
	Habitats HabitatsConfigurationState
	CloudProvider base.ProviderConfiguration
}

func (c *ConfigurationState) Init() {
	configurationStateMutex.Lock()
	defer configurationStateMutex.Unlock()
	c.Apps = AppsConfigurationState{}
	c.Habitats = HabitatsConfigurationState{}
	c.Trainer = base.TrainerConfigurationState{
		Port: 5000,
		Policies: base.TrainerPolicies{
			TRY_TO_REMOVE_HOSTS: true,
		},
	}
}

func (c *ConfigurationState) Snapshot() ConfigurationState {
	configurationStateMutex.Lock()
	defer configurationStateMutex.Unlock()
	res := *c
	return res
}

func (c * ConfigurationState) AllAppsLatest() map[base.AppName]base.AppConfiguration {
	apps := make(map[base.AppName]base.AppConfiguration)
	configurationStateMutex.Lock()
	confApps := c.Apps
	configurationStateMutex.Unlock()
	for appName, appObj := range confApps {
		elem, err := c.GetApp(appName, appObj.LatestVersion())
		if err == nil {
			apps[appName] = elem
		}
	}
	ConfigLogger.Infof("AllAppsLatest: %+v", apps)
	return apps
}

func (c *ConfigurationState) GetApp (name base.AppName, version base.Version) (base.AppConfiguration, error) {
	configurationStateMutex.Lock()
	defer configurationStateMutex.Unlock()
	if _, exists := (*c).Apps[name]; !exists {
		return base.AppConfiguration{}, errors.New("No such App")
	}
	if _, exists := (*c).Apps[name][version]; !exists {
		return base.AppConfiguration{}, errors.New("No such Version")
	}
	res := (*c).Apps[name][version]
	return res, nil
}

func (c *ConfigurationState) GetHabitat (name base.HabitatName, version base.Version) (base.HabitatConfiguration, error){
	configurationStateMutex.Lock()
	defer configurationStateMutex.Unlock()
	if _, exists := (*c).Habitats[name]; !exists {
		return base.HabitatConfiguration{}, errors.New("No such Habitat")
	}
	if _, exists := (*c).Habitats[name][version]; !exists {
		return base.HabitatConfiguration{}, errors.New("No such Version")
	}
	res := (*c).Habitats[name][version]
	return res, nil
}

func (c *ConfigurationState) ConfigureApp (conf base.AppConfiguration) {
	ConfigLogger.Infof("ConfigureApp %s:%d", conf.Name, conf.Version)
	configurationStateMutex.Lock()
	defer configurationStateMutex.Unlock()
	if _, exists := c.Apps[conf.Name]; !exists {
		c.Apps[conf.Name] = AppConfigurationVersions{}
	}
	c.Apps[conf.Name][conf.Version] = conf
	state_needs.GlobalAppsNeedState.UpdateNeeds(conf.Name, conf.Version, conf.Needs)
}

func (c *ConfigurationState) ConfigureHabitat (conf base.HabitatConfiguration) {
	configurationStateMutex.Lock()
	defer configurationStateMutex.Unlock()
	if _, exists := c.Habitats[conf.Name]; !exists {
		c.Habitats[conf.Name] = HabitatConfigurationVersions{}
	}
	c.Habitats[conf.Name][conf.Version] = conf
}

type AppsConfigurationState map[base.AppName]AppConfigurationVersions
type AppConfigurationVersions map[base.Version]base.AppConfiguration

func (a AppConfigurationVersions) LatestVersion() base.Version {
	var keys []int
	for k := range a {
		keys = append(keys, int(k))
	}
	if len(keys) == 0 {
		return 0
	}

	sort.Sort(sort.Reverse(sort.IntSlice(keys)))
	return base.Version(keys[0])
}

type HabitatsConfigurationState map[base.HabitatName]HabitatConfigurationVersions
type HabitatConfigurationVersions map[base.Version]base.HabitatConfiguration
type ProviderConfigurationState base.ProviderConfiguration


