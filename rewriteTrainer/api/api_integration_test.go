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

package api

import (
	"testing"
	"gatoor/orca/base"
	"gatoor/orca/rewriteTrainer/state/cloud"
	"gatoor/orca/rewriteTrainer/db"
	"gatoor/orca/rewriteTrainer/state/configuration"
	"gatoor/orca/rewriteTrainer/config"
	"gatoor/orca/rewriteTrainer/tracker"
	"time"
	"gatoor/orca/rewriteTrainer/responder"
	"gatoor/orca/rewriteTrainer/planner"
	"gatoor/orca/rewriteTrainer/cloud"
)


func applySampleConfig() {
	conf := config.JsonConfiguration{}

	conf.Trainer.Port = 5000


	httpApp1 := base.AppConfiguration{
		Name: "httpApp_1",
		Version: 1,
		Type: base.APP_HTTP,
		MinDeploymentCount: 3,
		TargetDeploymentCount: 3,
		//InstallCommands: []base.OsCommand{
		//	{
		//		Type: base.EXEC_COMMAND,
		//		Command: base.Command{"ls", "/home"},
		//	},
		//	{
		//		Type: base.FILE_COMMAND,
		//		Command: base.Command{"/server/app1/app1.conf", "somefilecontent as a string"},
		//	},
		//},
		//QueryStateCommand: base.OsCommand{
		//	Type: base.EXEC_COMMAND,
		//	Command: base.Command{"wget", "http://localhost:1234/check"},
		//},
		//RemoveCommand: base.OsCommand{
		//	Type: base.EXEC_COMMAND,
		//	Command: base.Command{"rm", "-rf /server/app1"},
		//},
		Needs: base.AppNeeds{
			MemoryNeeds: base.MemoryNeeds(1),
			CpuNeeds: base.CpuNeeds(1),
			NetworkNeeds: base.NetworkNeeds(1),
		},
	}

	httpApp1_v2 := base.AppConfiguration{
		Name: "httpApp_1",
		Version: 11,
		Type: base.APP_HTTP,
		MinDeploymentCount: 2,
		TargetDeploymentCount: 3,
		//InstallCommands: []base.OsCommand{
		//	{
		//		Type: base.EXEC_COMMAND,
		//		Command: base.Command{"ls", "/home"},
		//	},
		//	{
		//		Type: base.FILE_COMMAND,
		//		Command: base.Command{"/server/app1/app1.conf", "somefilecontent as a string"},
		//	},
		//},
		//QueryStateCommand: base.OsCommand{
		//	Type: base.EXEC_COMMAND,
		//	Command: base.Command{"wget", "http://localhost:1234/check"},
		//},
		//RemoveCommand: base.OsCommand{
		//	Type: base.EXEC_COMMAND,
		//	Command: base.Command{"rm", "-rf /server/app1"},
		//},
		Needs: base.AppNeeds{
			MemoryNeeds: base.MemoryNeeds(2),
			CpuNeeds: base.CpuNeeds(2),
			NetworkNeeds: base.NetworkNeeds(5),
		},
	}

	httpApp2 := base.AppConfiguration{
		Name: "httpApp_2",
		Version: 2,
		Type: base.APP_HTTP,
		MinDeploymentCount: 4,
		TargetDeploymentCount: 4,
		//InstallCommands: []base.OsCommand{
		//	{
		//		Type: base.EXEC_COMMAND,
		//		Command: base.Command{"ls", "/home"},
		//	},
		//	{
		//		Type: base.FILE_COMMAND,
		//		Command: base.Command{"/server/app1/app1.conf", "somefilecontent as a string"},
		//	},
		//},
		//QueryStateCommand: base.OsCommand{
		//	Type: base.EXEC_COMMAND,
		//	Command: base.Command{"wget", "http://localhost:1234/check"},
		//},
		//RemoveCommand: base.OsCommand{
		//	Type: base.EXEC_COMMAND,
		//	Command: base.Command{"rm", "-rf /server/app1"},
		//},
		Needs: base.AppNeeds{
			MemoryNeeds: base.MemoryNeeds(1),
			CpuNeeds: base.CpuNeeds(1),
			NetworkNeeds: base.NetworkNeeds(1),
		},
	}

	workerApp1 := base.AppConfiguration{
		Name: "workerApp_1",
		Version: 10,
		Type: base.APP_WORKER,
		MinDeploymentCount: 1,
		TargetDeploymentCount: 1,
		//InstallCommands: []base.OsCommand{
		//	{
		//		Type: base.EXEC_COMMAND,
		//		Command: base.Command{"ls", "/home"},
		//	},
		//	{
		//		Type: base.FILE_COMMAND,
		//		Command: base.Command{"/server/app1/app1.conf", "somefilecontent as a string"},
		//	},
		//},
		//QueryStateCommand: base.OsCommand{
		//	Type: base.EXEC_COMMAND,
		//	Command: base.Command{"wget", "http://localhost:1234/check"},
		//},
		//RemoveCommand: base.OsCommand{
		//	Type: base.EXEC_COMMAND,
		//	Command: base.Command{"rm", "-rf /server/app1"},
		//},
		Needs: base.AppNeeds{
			CpuNeeds: base.CpuNeeds(50),
			MemoryNeeds: base.MemoryNeeds(10),
			NetworkNeeds: base.NetworkNeeds(10),
		},
	}

	workerApp1_v2 := base.AppConfiguration{
		Name: "workerApp_1",
		Version: 11,
		Type: base.APP_WORKER,
		MinDeploymentCount: 1,
		TargetDeploymentCount: 1,
		//InstallCommands: []base.OsCommand{
		//	{
		//		Type: base.EXEC_COMMAND,
		//		Command: base.Command{"ls", "/home"},
		//	},
		//	{
		//		Type: base.FILE_COMMAND,
		//		Command: base.Command{"/server/app1/app1.conf", "somefilecontent as a string"},
		//	},
		//},
		//QueryStateCommand: base.OsCommand{
		//	Type: base.EXEC_COMMAND,
		//	Command: base.Command{"wget", "http://localhost:1234/check"},
		//},
		//RemoveCommand: base.OsCommand{
		//	Type: base.EXEC_COMMAND,
		//	Command: base.Command{"rm", "-rf /server/app1"},
		//},
		Needs: base.AppNeeds{
			CpuNeeds: base.CpuNeeds(70),
			MemoryNeeds: base.MemoryNeeds(40),
			NetworkNeeds: base.NetworkNeeds(30),
		},
	}

	workerApp2 := base.AppConfiguration{
		Name: "workerApp_2",
		Version: 20,
		Type: base.APP_WORKER,
		MinDeploymentCount: 5,
		TargetDeploymentCount: 5,
		//InstallCommands: []base.OsCommand{
		//	{
		//		Type: base.EXEC_COMMAND,
		//		Command: base.Command{"ls", "/home"},
		//	},
		//	{
		//		Type: base.FILE_COMMAND,
		//		Command: base.Command{"/server/app1/app1.conf", "somefilecontent as a string"},
		//	},
		//},
		//QueryStateCommand: base.OsCommand{
		//	Type: base.EXEC_COMMAND,
		//	Command: base.Command{"wget", "http://localhost:1234/check"},
		//},
		//RemoveCommand: base.OsCommand{
		//	Type: base.EXEC_COMMAND,
		//	Command: base.Command{"rm", "-rf /server/app1"},
		//},
		Needs: base.AppNeeds{
			CpuNeeds: base.CpuNeeds(23),
			MemoryNeeds: base.MemoryNeeds(23),
			NetworkNeeds: base.NetworkNeeds(23),
		},
	}

	workerApp3 := base.AppConfiguration{
		Name: "workerApp_3",
		Version: 30,
		Type: base.APP_WORKER,
		MinDeploymentCount: 100,
		TargetDeploymentCount: 100,
		//InstallCommands: []base.OsCommand{
		//	{
		//		Type: base.EXEC_COMMAND,
		//		Command: base.Command{"ls", "/home"},
		//	},
		//	{
		//		Type: base.FILE_COMMAND,
		//		Command: base.Command{"/server/app1/app1.conf", "somefilecontent as a string"},
		//	},
		//},
		//QueryStateCommand: base.OsCommand{
		//	Type: base.EXEC_COMMAND,
		//	Command: base.Command{"wget", "http://localhost:1234/check"},
		//},
		//RemoveCommand: base.OsCommand{
		//	Type: base.EXEC_COMMAND,
		//	Command: base.Command{"rm", "-rf /server/app1"},
		//},
		Needs: base.AppNeeds{
			CpuNeeds: base.CpuNeeds(7),
			MemoryNeeds: base.MemoryNeeds(2),
			NetworkNeeds: base.NetworkNeeds(1),
		},
	}



	conf.Apps = []base.AppConfiguration{
		httpApp1, httpApp1_v2, httpApp2,
		workerApp1, workerApp1_v2, workerApp2, workerApp3,
	}

	conf.ApplyToState()
}



func prepare() {
	db.Init("api_integration_test_")
	state_configuration.GlobalConfigurationState.Init()
	state_cloud.GlobalCloudLayout.Init()
	tracker.GlobalHostTracker = tracker.HostTracker{}
	tracker.GlobalAppsStatusTracker = tracker.AppsStatusTracker{}
	planner.Queue = *planner.NewPlannerQueue()
	applySampleConfig()
	cloud.Init()
}

func after() {
	db.Close()
}

func TestApi_doHandlePush_NoChanges(t *testing.T) {
	prepare()
	defer after()
	app1 := base.AppInfo{
		Type: base.APP_HTTP,
		Name: "httpApp_1",
		Version: 1,
		Status: base.STATUS_RUNNING,
	}
	app2 := base.AppInfo{
		Type: base.APP_WORKER,
		Name: "workerApp_2",
		Version: 20,
		Status: base.STATUS_RUNNING,
	}

	info := base.HostInfo{
	      HostId: "host1", Apps: []base.AppInfo{app1, app2, app2},
	}

	stats := base.MetricsWrapper{}

	if len(state_cloud.GlobalCloudLayout.Current.Layout) != 0 {
		t.Error(state_cloud.GlobalCloudLayout.Current)
	}

	doHandlePush(info, stats)
	//nothing changed, so the config should be empty
	_, errC := responder.GetConfigForHost("host1")

	elem, _ := state_cloud.GlobalCloudLayout.Current.GetHost("host1")
	if elem.HostId != "host1" || elem.Apps["httpApp_1"].DeploymentCount != 1 || elem.Apps["workerApp_2"].DeploymentCount != 2 {
		t.Error(elem)
	}
	if errC == nil {
		t.Error(errC)
	}

	info2 := base.HostInfo{
		HostId: "host1", Apps: []base.AppInfo{app1, app2, app2, app2},
	}
	doHandlePush(info2, stats)

	elem2, _ := state_cloud.GlobalCloudLayout.Current.GetHost("host1")
	if elem2.HostId != "host1" || elem2.Apps["httpApp_1"].DeploymentCount != 1 || elem2.Apps["workerApp_2"].DeploymentCount != 3 {
		t.Error(elem2)
	}

	track, _ := tracker.GlobalHostTracker.Get("host1")
	if !track.LastCheckin.After(time.Now().UTC().Add(-time.Duration(time.Second * 2))) {
		t.Error(track)
	}

	appTrack := tracker.GlobalAppsStatusTracker["httpApp_1"][1]
	if appTrack.Rating != tracker.RATING_STABLE || appTrack.RunningCount != 2 || len(appTrack.CrashDetails) != 0{
		t.Errorf("%+v", appTrack)
	}
	_, errC2 := responder.GetConfigForHost("host1")

	if errC2 == nil {
		t.Error(errC2)
	}
}


func TestApi_doHandlePush_AppShouldBeUpdated(t *testing.T) {
	//check that the config contains the updated version
	prepare()
	defer after()
	app1 := base.AppInfo{
		Type: base.APP_HTTP,
		Name: "httpApp_1",
		Version: 1,
		Status: base.STATUS_RUNNING,
	}
	app2 := base.AppInfo{
		Type: base.APP_WORKER,
		Name: "workerApp_1",
		Version: 10,
		Status: base.STATUS_RUNNING,
	}

	info := base.HostInfo{
		HostId: "host1", Apps: []base.AppInfo{app1, app2, app2},
	}

	stats := base.MetricsWrapper{}

	if len(state_cloud.GlobalCloudLayout.Current.Layout) != 0 {
		t.Error(state_cloud.GlobalCloudLayout.Current)
	}

	doHandlePush(info, stats)
	//nothing changed, so the config should be empty
	_, errC := responder.GetConfigForHost("host1")

	elem, _ := state_cloud.GlobalCloudLayout.Current.GetHost("host1")
	if elem.HostId != "host1" || elem.Apps["httpApp_1"].DeploymentCount != 1 || elem.Apps["workerApp_1"].DeploymentCount != 2 {
		t.Error(elem)
	}
	if errC == nil {
		t.Error(errC)
	}


	info2 := base.HostInfo{
		HostId: "host1", Apps: []base.AppInfo{app1, app2, app2, app2},
	}

	//set an update:
	planner.Queue.Add("host1", "httpApp_1", state_cloud.AppsVersion{Version: 11, DeploymentCount: 1})
	doHandlePush(info2, stats)

	elem2, _ := state_cloud.GlobalCloudLayout.Current.GetHost("host1")
	if elem2.HostId != "host1" || elem2.Apps["httpApp_1"].DeploymentCount != 1 || elem2.Apps["workerApp_1"].DeploymentCount != 3 {
		t.Error(elem2)
	}

	track, _ := tracker.GlobalHostTracker.Get("host1")
	if !track.LastCheckin.After(time.Now().UTC().Add(-time.Duration(time.Second * 2))) {
		t.Error(track)
	}

	appTrack := tracker.GlobalAppsStatusTracker["httpApp_1"][1]
	if appTrack.Rating != tracker.RATING_STABLE || appTrack.RunningCount != 2 || len(appTrack.CrashDetails) != 0{
		t.Errorf("%+v", appTrack)
	}
	conf, _ := responder.GetConfigForHost("host1")

	if conf.DeploymentCount != 1 || conf.AppConfiguration.Name != "httpApp_1" || conf.AppConfiguration.Type != base.APP_HTTP || conf.AppConfiguration.Version != 11 {
		t.Error(conf)
	}

	//same with worker app

	planner.Queue.RemoveHost("host1")
	planner.Queue.Add("host1", "workerApp_1", state_cloud.AppsVersion{Version: 11, DeploymentCount: 10})
	doHandlePush(info2, stats)

	elem3, _ := state_cloud.GlobalCloudLayout.Current.GetHost("host1")
	if elem3.HostId != "host1" || elem3.Apps["httpApp_1"].DeploymentCount != 1 || elem3.Apps["workerApp_1"].DeploymentCount != 3 {
		t.Error(elem3)
	}

	track1, _ := tracker.GlobalHostTracker.Get("host1")
	if !track.LastCheckin.After(time.Now().UTC().Add(-time.Duration(time.Second * 2))) {
		t.Error(track1)
	}

	appTrack1 := tracker.GlobalAppsStatusTracker["httpApp_1"][1]
	if appTrack1.Rating != tracker.RATING_STABLE || appTrack1.RunningCount != 3 || len(appTrack1.CrashDetails) != 0{
		t.Errorf("%+v", appTrack1)
	}
	conf1, _ := responder.GetConfigForHost("host1")

	if conf1.DeploymentCount != 10 || conf1.AppConfiguration.Name != "workerApp_1" || conf1.AppConfiguration.Type != base.APP_WORKER || conf1.AppConfiguration.Version != 11 {
		t.Error(conf)
	}
}


func TestApi_doHandlePush_AppStillUpdating(t *testing.T) {
	// An app update is not done jet.
	// check that the CurrentLayout is NOT updated
	// check that the app tracker is NOT updated
	prepare()
	defer after()
	app1 := base.AppInfo{
		Type: base.APP_HTTP,
		Name: "httpApp_1",
		Version: 11,
		Status: base.STATUS_DEPLOYING,
	}
	app2 := base.AppInfo{
		Type: base.APP_WORKER,
		Name: "workerApp_1",
		Version: 10,
		Status: base.STATUS_RUNNING,
	}

	info := base.HostInfo{
		HostId: "host1", Apps: []base.AppInfo{app1, app2, app2},
	}

	stats := base.MetricsWrapper{}

	if len(state_cloud.GlobalCloudLayout.Current.Layout) != 0 {
		t.Error(state_cloud.GlobalCloudLayout.Current)
	}

	planner.Queue.Add("host1", "httpApp_1", state_cloud.AppsVersion{11, 1})
	planner.Queue.Add("host2", "httpApp_1", state_cloud.AppsVersion{11, 1})
	planner.Queue.SetState("host1", "httpApp_1", planner.STATE_APPLYING)
	s, _ := planner.Queue.GetState("host1", "httpApp_1")
	if s != planner.STATE_APPLYING {
		t.Error(s)
	}

	doHandlePush(info, stats)

	conf, err := responder.GetConfigForHost("host1")
	if err != nil {
		t.Error(conf)
	}
	if conf.DeploymentCount != 1 || conf.AppConfiguration.Name != "httpApp_1" {
		t.Error(conf)
	}

	//app is not rated yet:
	tA := tracker.GlobalAppsStatusTracker["httpApp_1"][11]
	if tA.RunningCount != 0 || tA.Rating != "" {
		t.Error(tA)
	}

	layout := state_cloud.GlobalCloudLayout.Current.Layout["host1"].Apps
	//current layout for host1 should show no httpApp_1, it is not running!
	if layout["workerApp_1"].Version != 10 || layout["httpApp_1"].Version != 0 {
		t.Error(layout)
	}
	//check that other queued updates still exist:
	s2, err2 := planner.Queue.Get("host2")
	if err2 != nil {
		t.Error(s2)
	}
}


func TestApi_doHandlePush_AppUpdate(t *testing.T) {
	// An app updated successfully.
	// check that the CurrentLayout is updated
	// check that the app tracker is updated


	prepare()
	defer after()
	app1 := base.AppInfo{
		Type: base.APP_HTTP,
		Name: "httpApp_1",
		Version: 11,
		Status: base.STATUS_RUNNING,
	}
	app2 := base.AppInfo{
		Type: base.APP_WORKER,
		Name: "workerApp_1",
		Version: 10,
		Status: base.STATUS_RUNNING,
	}

	info := base.HostInfo{
		HostId: "host1", Apps: []base.AppInfo{app1, app2, app2},
	}

	stats := base.MetricsWrapper{}

	if len(state_cloud.GlobalCloudLayout.Current.Layout) != 0 {
		t.Error(state_cloud.GlobalCloudLayout.Current)
	}

	planner.Queue.Add("host1", "httpApp_1", state_cloud.AppsVersion{11, 1})
	planner.Queue.Add("host2", "httpApp_1", state_cloud.AppsVersion{11, 1})
	planner.Queue.SetState("host1", "httpApp_1", planner.STATE_APPLYING)
	s, _ := planner.Queue.GetState("host1", "httpApp_1")
	if s != planner.STATE_APPLYING {
		t.Error(s)
	}

	doHandlePush(info, stats)

	conf, err := responder.GetConfigForHost("host1")
	if err == nil {
		t.Error(conf)
	}

	tA := tracker.GlobalAppsStatusTracker["httpApp_1"][11]
	if tA.RunningCount != 1 || tA.Rating != tracker.RATING_STABLE {
		t.Error(tA)
	}

	layout := state_cloud.GlobalCloudLayout.Current.Layout["host1"].Apps
	if layout["workerApp_1"].Version != 10 || layout["httpApp_1"].Version != 11 {
		t.Error(layout)
	}

	//check that other queued updates still exist:
	s2, err2 := planner.Queue.Get("host2")
	if err2 != nil {
		t.Error(s2)
	}
}


func TestApi_doHandlePush_AppRollback(t *testing.T) {
	// An app rolled back.
	// check that the CurrentLayout is NOT updated
	// check that the app tracker is updated
	// check that the planner queue is updated (all updates of the app should be removed)


	prepare()
	defer after()
	app1 := base.AppInfo{
		Type: base.APP_HTTP,
		Name: "httpApp_1",
		Version: 1,
		Status: base.STATUS_RUNNING,
	}
	app2 := base.AppInfo{
		Type: base.APP_WORKER,
		Name: "workerApp_1",
		Version: 10,
		Status: base.STATUS_RUNNING,
	}

	info := base.HostInfo{
		HostId: "host1", Apps: []base.AppInfo{app1, app2, app2},
	}

	stats := base.MetricsWrapper{}

	if len(state_cloud.GlobalCloudLayout.Current.Layout) != 0 {
		t.Error(state_cloud.GlobalCloudLayout.Current)
	}

	planner.Queue.Add("host1", "httpApp_1", state_cloud.AppsVersion{11, 1})
	planner.Queue.Add("host2", "httpApp_1", state_cloud.AppsVersion{11, 1})
	planner.Queue.SetState("host1", "httpApp_1", planner.STATE_APPLYING)
	s, _ := planner.Queue.GetState("host1", "httpApp_1")
	if s != planner.STATE_APPLYING {
		t.Error(s)
	}

	doHandlePush(info, stats)

	conf, err := responder.GetConfigForHost("host1")
	if err == nil {
		t.Error(conf)
	}

	//new app version never checked in, handle as crash:
	tA := tracker.GlobalAppsStatusTracker["httpApp_1"][11]
	if tA.RunningCount != 0 || tA.Rating != tracker.RATING_CRASHED {
		t.Error(tA)
	}

	layout := state_cloud.GlobalCloudLayout.Current.Layout["host1"].Apps
	if layout["workerApp_1"].Version != 10 || layout["httpApp_1"].Version != 1 {
		t.Error(layout)
	}

	//check that other queued updates were removed:
	s2, err2 := planner.Queue.Get("host2")
	if err2 == nil {
		t.Error(s2)
	}
}


func TestApi_doHandlePush_AppShouldBeRemovedFromHost(t *testing.T) {
	//check that the config contains the new deploymentCount 0

	prepare()
	defer after()
	app1 := base.AppInfo{
		Type: base.APP_HTTP,
		Name: "httpApp_1",
		Version: 11,
		Status: base.STATUS_RUNNING,
	}
	app2 := base.AppInfo{
		Type: base.APP_WORKER,
		Name: "workerApp_1",
		Version: 10,
		Status: base.STATUS_RUNNING,
	}

	info := base.HostInfo{
		HostId: "host1", Apps: []base.AppInfo{app1, app2, app2},
	}

	stats := base.MetricsWrapper{}

	if len(state_cloud.GlobalCloudLayout.Current.Layout) != 0 {
		t.Error(state_cloud.GlobalCloudLayout.Current)
	}

	planner.Queue.Add("host1", "httpApp_1", state_cloud.AppsVersion{11, 0})
	planner.Queue.Add("host2", "httpApp_1", state_cloud.AppsVersion{11, 1})
	planner.Queue.SetState("host1", "httpApp_1", planner.STATE_APPLYING)
	s, _ := planner.Queue.GetState("host1", "httpApp_1")
	if s != planner.STATE_APPLYING {
		t.Error(s)
	}

	doHandlePush(info, stats)

	conf, err := responder.GetConfigForHost("host1")
	if err != nil {
		t.Error(conf)
	}

	if conf.DeploymentCount != 0 {
		t.Error(conf)
	}

	tA := tracker.GlobalAppsStatusTracker["httpApp_1"][11]
	if tA.RunningCount != 1 || tA.Rating != tracker.RATING_STABLE {
		t.Error(tA)
	}

	layout := state_cloud.GlobalCloudLayout.Current.Layout["host1"].Apps
	if layout["workerApp_1"].Version != 10 || layout["httpApp_1"].Version != 11 {
		t.Error(layout)
	}

	//check that other queued updates still exist:
	s2, err2 := planner.Queue.Get("host2")
	if err2 != nil {
		t.Error(s2)
	}
}

func TestApi_doHandlePush_AppShouldBeScaled(t *testing.T) {
	//check that the config contains the new deploymentCount
	prepare()
	defer after()
	app1 := base.AppInfo{
		Type: base.APP_HTTP,
		Name: "httpApp_1",
		Version: 11,
		Status: base.STATUS_RUNNING,
	}
	app2 := base.AppInfo{
		Type: base.APP_WORKER,
		Name: "workerApp_1",
		Version: 10,
		Status: base.STATUS_RUNNING,
	}

	info := base.HostInfo{
		HostId: "host1", Apps: []base.AppInfo{app1, app2, app2},
	}

	stats := base.MetricsWrapper{}

	if len(state_cloud.GlobalCloudLayout.Current.Layout) != 0 {
		t.Error(state_cloud.GlobalCloudLayout.Current)
	}

	planner.Queue.Add("host1", "httpApp_1", state_cloud.AppsVersion{11, 5})
	planner.Queue.Add("host2", "httpApp_1", state_cloud.AppsVersion{11, 1})
	planner.Queue.SetState("host1", "httpApp_1", planner.STATE_APPLYING)
	s, _ := planner.Queue.GetState("host1", "httpApp_1")
	if s != planner.STATE_APPLYING {
		t.Error(s)
	}

	doHandlePush(info, stats)

	conf, err := responder.GetConfigForHost("host1")
	if err != nil {
		t.Error(conf)
	}

	if conf.DeploymentCount != 5 {
		t.Errorf("%+v", conf)
	}

	tA := tracker.GlobalAppsStatusTracker["httpApp_1"][11]
	if tA.RunningCount != 1 || tA.Rating != tracker.RATING_STABLE {
		t.Error(tA)
	}

	layout := state_cloud.GlobalCloudLayout.Current.Layout["host1"].Apps
	if layout["workerApp_1"].Version != 10 || layout["httpApp_1"].Version != 11 {
		t.Error(layout)
	}

	//check that other queued updates still exist:
	s2, err2 := planner.Queue.Get("host2")
	if err2 != nil {
		t.Error(s2)
	}
}


func TestApi_doHandlePush_NewInstance(t *testing.T) {
	// This is the first time the new host checks in.
	// Check that the CurrentLayout is updated.
	// Check that the HostTracker is updated.

	prepare()
	defer after()
	state_cloud.GlobalCloudLayout.Current.AddEmptyHost("somehost")
	info := base.HostInfo{
		HostId: "new_host", Apps: []base.AppInfo{},
	}

	stats := base.MetricsWrapper{}

	if len(state_cloud.GlobalCloudLayout.Current.Layout) != 1 {
		t.Error(state_cloud.GlobalCloudLayout.Current)
	}

	doHandlePush(info, stats)

	conf, err := responder.GetConfigForHost("host1")
	if err == nil {
		t.Error(conf)
	}

	layout := state_cloud.GlobalCloudLayout.Current.Layout["new_host"].Apps
	if len(layout) != 0 {
		t.Error(layout)
	}
	if len(state_cloud.GlobalCloudLayout.Current.Layout) != 2 {
		t.Error(state_cloud.GlobalCloudLayout.Current.Layout)
	}
}


func TestApi_doHandlePush_AppDied(t *testing.T) {
	// An app died.
	// check that the CurrentLayout is updated
	// check that the app tracker is updated


	prepare()
	defer after()
	app1 := base.AppInfo{
		Type: base.APP_HTTP,
		Name: "httpApp_1",
		Version: 1,
		Status: base.STATUS_RUNNING,
	}
	app2 := base.AppInfo{
		Type: base.APP_WORKER,
		Name: "workerApp_2",
		Version: 20,
		Status: base.STATUS_DEAD,
	}

	info := base.HostInfo{
		HostId: "host1", Apps: []base.AppInfo{app1, app2, app2},
	}

	stats := base.MetricsWrapper{}

	if len(state_cloud.GlobalCloudLayout.Current.Layout) != 0 {
		t.Error(state_cloud.GlobalCloudLayout.Current)
	}

	doHandlePush(info, stats)

	_, errC := responder.GetConfigForHost("host1")

	elem, _ := state_cloud.GlobalCloudLayout.Current.GetHost("host1")

	if elem.HostId != "host1" || elem.Apps["httpApp_1"].DeploymentCount != 1 {
		t.Error(elem)
	}

	if _, exists := elem.Apps["workerApp_2"]; exists {
		t.Error(elem.Apps["workerApp_2"] )
	}

	if errC == nil {
		t.Error(errC)
	}

	info2 := base.HostInfo{
		HostId: "host1", Apps: []base.AppInfo{app1, app2, app2, app2},
	}
	doHandlePush(info2, stats)

	track, _ := tracker.GlobalHostTracker.Get("host1")
	if !track.LastCheckin.After(time.Now().UTC().Add(-time.Duration(time.Second * 2))) {
		t.Error(track)
	}

	appTrack := tracker.GlobalAppsStatusTracker["httpApp_1"][1]
	if appTrack.Rating != tracker.RATING_STABLE || appTrack.RunningCount != 2 || len(appTrack.CrashDetails) != 0{
		t.Errorf("%+v", appTrack)
	}

	appTrackWorker := tracker.GlobalAppsStatusTracker["workerApp_2"][20]
	if appTrackWorker.Rating != tracker.RATING_CRASHED || appTrackWorker.CrashDetails[0].Cause != "APP_EVENT_CRASH" {
		t.Errorf("%+v", appTrackWorker)
	}

	_, errC2 := responder.GetConfigForHost("host1")

	if errC2 == nil {
		t.Error(errC2)
	}
}

func TestApi_doHandleCloudEvent_HostSpawned(t *testing.T) {
	// Host spawned.   TODO think about other stuff to check
	// check that the CurrentLayout is updated
	// check that the app tracker is updated
}

func TestApi_doHandleCloudEvent_HostDied(t *testing.T) {
	// Host died.   TODO think about other stuff to check
	// check that the CurrentLayout is updated
	// check that the app tracker is updated
}
