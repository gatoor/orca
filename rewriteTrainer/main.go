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

package main

import (
	"gatoor/orca/rewriteTrainer/state/needs"
	"gatoor/orca/rewriteTrainer/state/cloud"
	"gatoor/orca/rewriteTrainer/state/configuration"
	Logger "gatoor/orca/rewriteTrainer/log"
	"gatoor/orca/rewriteTrainer/config"
	"gatoor/orca/rewriteTrainer/api"
	"gatoor/orca/rewriteTrainer/cloud"
	"gatoor/orca/rewriteTrainer/db"
	"gatoor/orca/rewriteTrainer/scheduler"
	"gatoor/orca/rewriteTrainer/planner"
	"time"
	"flag"
)


const CHECKIN_WAIT_TIME = 30

func main() {
	var configurationRoot = flag.String("configroot", "/orca/config/", "Configuration Root Directory")
	flag.Parse()

	db.Audit.Init("localhost")
	db.Audit.Insert__AuditEvent(db.AuditEvent{Details:map[string]string{
		"message": "Orca Trainer Started",
	}})

	var baseConfiguration config.JsonConfiguration

	Logger.InitLogger.Info("Starting trainer...")
	initState()
	initConfig(&baseConfiguration, *configurationRoot)
	cloud.Init()
	db.Audit.Init("")
	initApi(&baseConfiguration)

	state_cloud.GlobalCloudLayout.InitBaseInstances()

	waitForCheckin()
	scheduler.Start()
	planner.InitialPlan()
	Logger.InitLogger.Info("Trainer started")
	ticker := time.NewTicker(time.Second * 60)
	for {
		<- ticker.C
	}
}

func waitForCheckin() {
	Logger.InitLogger.Infof("Waiting %ds for existsing clients to check in", CHECKIN_WAIT_TIME)
	time.Sleep(time.Duration(CHECKIN_WAIT_TIME * time.Second))
	Logger.InitLogger.Info("Done waiting")
}

func initState() {
	state_configuration.GlobalConfigurationState.Init()
	state_cloud.GlobalCloudLayout.Init()
	state_needs.GlobalAppsNeedState = state_needs.AppsNeedState{}
}

func initConfig(baseConfiguration *config.JsonConfiguration, configurationRoot string) {
	baseConfiguration.Init(configurationRoot)
	baseConfiguration.Load()
	baseConfiguration.Check()
	baseConfiguration.ApplyToState()

	db.Audit.Insert__AuditEvent(db.AuditEvent{Details:map[string]string{
		"message": "Configuration has been loaded from filesystem",
	}})
}


func initApi(baseConfiguration *config.JsonConfiguration) {
	var a api.Api
	a.ConfigManager = baseConfiguration
	a.Init()
}