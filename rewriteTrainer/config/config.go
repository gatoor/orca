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

package config

import (
	"gatoor/orca/base"
	Logger "gatoor/orca/rewriteTrainer/log"
	"gatoor/orca/rewriteTrainer/state/needs"
	"os"
	"encoding/json"
	"gatoor/orca/util"
	"fmt"
	"gatoor/orca/rewriteTrainer/state/configuration"
	"gatoor/orca/rewriteTrainer/state/cloud"
	"io/ioutil"
)


type JsonConfiguration struct {
	configRoot string

	Trainer base.TrainerConfigurationState
	AvailableInstances []base.HostId
	Habitats []base.HabitatConfiguration
	Apps []base.AppConfiguration
	CloudProvider base.ProviderConfiguration
}

func loadConfigFromFile(filename string, conf interface{}) {
	Logger.InitLogger.Infof("Loading config file from %s", filename)
	file, err := os.Open(filename)
	if err != nil {
		Logger.InitLogger.Fatalf("Could not open config file %s - %s", filename, err)
		return
	}

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(conf); err != nil {
		extra := ""
		if serr, ok := err.(*json.SyntaxError); ok {
			line, col, highlight := util.HighlightBytePosition(file, serr.Offset)
			extra = fmt.Sprintf(":\nError at line %d, column %d (file offset %d):\n%s",
				line, col, serr.Offset, highlight)
		}
		Logger.InitLogger.Fatalf("error parsing JSON object in config file %s%s\n%v",
			file.Name(), extra, err)
	}
	file.Close()
}

func saveConfigToFile(filename string, conf interface{}) {
	res, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		Logger.InitLogger.Errorf("JsonConfiguration Derialize failed: %s; %+v", err, conf)
	}
	var result = string(res)
	err = ioutil.WriteFile(filename, []byte(result), 0644)
	if err != nil {
		panic(err)
	}
}

func (j *JsonConfiguration) Check() {
}

var (
	TRAINER_CONFIGURATION_FILE = "trainer.json"
	APPS_CONFIGURATION_FILE = "applications.json"
	AVAILABLE_INSTANCES_CONFIGURATION_FILE = "instances.json"
	CLOUD_PROVIDER_CONFIGURATION_FILE = "provider.json"
)

func (j *JsonConfiguration) Init(configurationRoot string){
	j.configRoot = configurationRoot
	state_configuration.GlobalConfigurationState.ConfigurationRootPath = configurationRoot
}

func (j *JsonConfiguration) Load() {
	loadConfigFromFile(j.configRoot + TRAINER_CONFIGURATION_FILE, &j.Trainer)
	loadConfigFromFile(j.configRoot + APPS_CONFIGURATION_FILE, &j.Apps)
	loadConfigFromFile(j.configRoot + AVAILABLE_INSTANCES_CONFIGURATION_FILE, &j.AvailableInstances)
	loadConfigFromFile(j.configRoot + CLOUD_PROVIDER_CONFIGURATION_FILE, &j.CloudProvider)
}

func (j *JsonConfiguration) Save() {
	Logger.InitLogger.Infof("Saving all configuration files")
	saveConfigToFile(j.configRoot + TRAINER_CONFIGURATION_FILE, state_configuration.GlobalConfigurationState.Trainer)

	var buffer = make([]base.AppConfiguration, 0)
	for _, application := range state_configuration.GlobalConfigurationState.AllAppsLatest() {
		buffer = append(buffer, application)
	}

	saveConfigToFile(j.configRoot + APPS_CONFIGURATION_FILE, buffer)
	//loadConfigFromFile(j.configRoot + AVAILABLE_INSTANCES_CONFIGURATION_FILE + ".saved", state_configuration.GlobalConfigurationState.)
	saveConfigToFile(j.configRoot + CLOUD_PROVIDER_CONFIGURATION_FILE, state_configuration.GlobalConfigurationState.CloudProvider)
}


func (j *JsonConfiguration)  ApplyToState() {
	Logger.InitLogger.Infof("Applying config to State")
	//applyHabitatConfig(j.Habitats)
	applyTrainerConfig(j.Trainer)
	applyAppsConfig(j.Apps)
	applyNeeds(j.Apps)
	applyAvailableInstances(j.AvailableInstances)
	applyCloudProviderConfiguration(j.CloudProvider)
	Logger.InitLogger.Infof("Config was applied to State")
}

func applyAvailableInstances(instances []base.HostId) {
	Logger.InitLogger.Infof("Applying AvailableInstances config: %+v", instances)
	for _, hostId := range instances {
		state_cloud.GlobalAvailableInstances.Update(hostId, base.InstanceResources{UsedCpuResource:0, UsedMemoryResource:0, UsedNetworkResource:0, TotalCpuResource: 20, TotalMemoryResource: 20, TotalNetworkResource: 20})
	}
}

func applyCloudProviderConfiguration(conf base.ProviderConfiguration) {
	Logger.InitLogger.Infof("Applying CloudProvider config: %+v", conf)
	state_configuration.GlobalConfigurationState.CloudProvider.SSHKey = conf.SSHKey
	state_configuration.GlobalConfigurationState.CloudProvider.SSHUser = conf.SSHUser
	state_configuration.GlobalConfigurationState.CloudProvider.Type = conf.Type
	state_configuration.GlobalConfigurationState.CloudProvider.MaxInstances = conf.MaxInstances
	state_configuration.GlobalConfigurationState.CloudProvider.MinInstances = conf.MinInstances
	state_configuration.GlobalConfigurationState.CloudProvider.AWSConfiguration = conf.AWSConfiguration
}

func applyAppsConfig(appsConfs []base.AppConfiguration) {
	for _, aConf := range appsConfs {
		app := base.AppConfiguration{
			Name: aConf.Name,
			Type: aConf.Type,
			Version: aConf.Version,
			TargetDeploymentCount: aConf.TargetDeploymentCount,
			MinDeploymentCount: aConf.MinDeploymentCount,
			DockerConfig: aConf.DockerConfig,
			LoadBalancer: aConf.LoadBalancer,
			Network: aConf.Network,
			PortMappings: aConf.PortMappings,
			Needs: aConf.Needs,
		}
		app.WeekdayBasedDeploymentCount.SetFlat(aConf.TargetDeploymentCount)
		app.WeekdayBasedAppNeeds.SetFlat(aConf.Needs)
		state_configuration.GlobalConfigurationState.ConfigureApp(app)
	}
}

func applyHabitatConfig (habitatConfs []base.HabitatConfiguration) {
	for _, hConf := range habitatConfs {
		state_configuration.GlobalConfigurationState.ConfigureHabitat(base.HabitatConfiguration{
			Name: hConf.Name,
			Version: hConf.Version,
			InstallCommands: hConf.InstallCommands,
		})
	}
}

func applyTrainerConfig (trainerConf base.TrainerConfigurationState) {
	state_configuration.GlobalConfigurationState.Trainer.Port = trainerConf.Port
	state_configuration.GlobalConfigurationState.Trainer.Ip = trainerConf.Ip
	state_configuration.GlobalConfigurationState.Trainer.Policies.TRY_TO_REMOVE_HOSTS = trainerConf.Policies.TRY_TO_REMOVE_HOSTS
}

//TODO use WeeklyNeeds
func applyNeeds(appConfs []base.AppConfiguration) {
	for _, aNeeds := range appConfs {
		state_needs.GlobalAppsNeedState.UpdateNeeds(aNeeds.Name, aNeeds.Version, aNeeds.Needs)
	}
}


func (j *JsonConfiguration) Serialize() string {
	res, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		Logger.InitLogger.Errorf("JsonConfiguration Derialize failed: %s; %+v", err, j)
	}
	return string(res)
}