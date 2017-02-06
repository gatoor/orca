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
	"gatoor/orca/rewriteTrainer/state/cloud"
	"testing"
	"gatoor/orca/rewriteTrainer/planner"
	"gatoor/orca/rewriteTrainer/state/configuration"
	"gatoor/orca/rewriteTrainer/state/needs"
	"gatoor/orca/base"
	"gatoor/orca/rewriteTrainer/cloud"
	"gatoor/orca/rewriteTrainer/config"
	"gatoor/orca/rewriteTrainer/needs"
)

func initTrainer() {
	state_configuration.GlobalConfigurationState.Init()
	state_cloud.GlobalCloudLayout.Init()
	state_needs.GlobalAppsNeedState = state_needs.AppsNeedState{}
	applySampleConfig()
	initCloudProvider()
	planner.Queue = *planner.NewPlannerQueue()
	cloud.Init()
}

func applySampleConfig() {
	conf := config.JsonConfiguration{}

	conf.Trainer.Port = 5000
	conf.Trainer.Policies.TRY_TO_REMOVE_HOSTS = true


	httpApp1 := base.AppConfiguration{
		Name: "httpApp_1",
		Version: 1,
		Type: base.APP_HTTP,
		TargetDeploymentCount: 3,
		MinDeploymentCount: 3,
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
		Version: 2,
		Type: base.APP_HTTP,
		TargetDeploymentCount: 2,
		MinDeploymentCount: 2,
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
		Version: 3,
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
		Version: 7,
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
		Version: 8,
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
		Version: 9,
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
		Version: 10,
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

func initCloudProvider() {
	state_cloud.GlobalAvailableInstances.Update("cpuHost_1", base.InstanceResources{
		TotalCpuResource: 500,
		TotalMemoryResource: 100,
		TotalNetworkResource: 100,
	})
	state_cloud.GlobalAvailableInstances.Update("cpuHost_2", base.InstanceResources{
		TotalCpuResource: 501,
		TotalMemoryResource: 101,
		TotalNetworkResource: 101,
	})
	state_cloud.GlobalAvailableInstances.Update("memoryHost_1", base.InstanceResources{
		TotalCpuResource: 200,
		TotalMemoryResource: 300,
		TotalNetworkResource: 100,
	})
	state_cloud.GlobalAvailableInstances.Update("generalHost_1", base.InstanceResources{
		TotalCpuResource: 101,
		TotalMemoryResource: 101,
		TotalNetworkResource: 101,
	})
	state_cloud.GlobalAvailableInstances.Update("generalHost_2", base.InstanceResources{
		TotalCpuResource: 102,
		TotalMemoryResource: 102,
		TotalNetworkResource: 102,
	})
	state_cloud.GlobalAvailableInstances.Update("generalHost_3", base.InstanceResources{
		TotalCpuResource: 103,
		TotalMemoryResource: 103,
		TotalNetworkResource: 103,
	})
	state_cloud.GlobalAvailableInstances.Update("emptyHost", base.InstanceResources{
		TotalCpuResource: 1000,
		TotalMemoryResource: 1000,
		TotalNetworkResource: 1000,
	})
}

func testInitConfig(t *testing.T) {
	if len(state_configuration.GlobalConfigurationState.Apps) != 5 {
		t.Error("init state_config apps wrong len")
	}
	if len(state_configuration.GlobalConfigurationState.Apps["httpApp_1"]) != 2 {
		t.Error("init state_config apps wrong len")
	}
	if len(state_configuration.GlobalConfigurationState.Apps["httpApp_2"]) != 1 {
		t.Error("init state_config apps wrong len")
	}
	if len(state_configuration.GlobalConfigurationState.Apps["workerApp_1"]) != 2 {
		t.Error("init state_config apps wrong len")
	}
	if len(state_configuration.GlobalConfigurationState.Apps["workerApp_3"]) != 1 {
		t.Error("init state_config apps wrong len")
	}

	if len(state_cloud.GlobalCloudLayout.Current.Layout) != 0 {
		t.Error("init state_cloud current should be empty")
	}
	if len(state_cloud.GlobalCloudLayout.Desired.Layout) != 0 {
		t.Error("init state_cloud desired should be empty")
	}

	if len(state_needs.GlobalAppsNeedState) != 5 {
		t.Error("init state_needs wrong len")
	}
	elem , _ := state_needs.GlobalAppsNeedState.Get("httpApp_1", 2)
	if elem.CpuNeeds != 2 {
		t.Error("wrong needs")
	}
	elem2 , _ := state_needs.GlobalAppsNeedState.Get("workerApp_3", 10)
	if elem2.MemoryNeeds != 2 {
		t.Error("wrong needs")
	}
}

func testLayout(t *testing.T) {
	if len(state_cloud.GlobalAvailableInstances) != 7 {
		t.Error(state_cloud.GlobalAvailableInstances)
	}
}

func testInit(t * testing.T) {
	testInitConfig(t)
	testLayout(t)
}













func TestResponder_GetConfigForHost_Integration(t *testing.T) {
	initTrainer()
	testInit(t)

	state_cloud.GlobalCloudLayout.Current.AddEmptyHost("cpuHost_1")
	state_cloud.GlobalCloudLayout.Current.AddEmptyHost("cpuHost_2")
	state_cloud.GlobalCloudLayout.Current.AddEmptyHost("memoryHost_1")
	state_cloud.GlobalCloudLayout.Current.AddEmptyHost("generalHost_1")
	state_cloud.GlobalCloudLayout.Current.AddEmptyHost("generalHost_2")
	state_cloud.GlobalCloudLayout.Current.AddEmptyHost("generalHost_3")
	state_cloud.GlobalCloudLayout.Current.AddEmptyHost("emptyHost")

	//httpApp_1 already deployed MinDeploymentCount
	state_cloud.GlobalCloudLayout.Current.AddApp("cpuHost_1", "httpApp_1", 1, 1)
	state_cloud.GlobalCloudLayout.Current.AddApp("cpuHost_2", "httpApp_1", 2, 1)

	//httpApp_2 missing 1 app for MinDeploymentCount
	state_cloud.GlobalCloudLayout.Current.AddApp("generalHost_1", "httpApp_2", 3, 1)
	state_cloud.GlobalCloudLayout.Current.AddApp("generalHost_2", "httpApp_2", 3, 1)
	state_cloud.GlobalCloudLayout.Current.AddApp("generalHost_3", "httpApp_2", 3, 1)

	//workerApp_1 missing

	//workerApp_2 already deployed MinDeploymentCount all on one host
	state_cloud.GlobalCloudLayout.Current.AddApp("cpuHost_1", "workerApp_2", 9, 5)

	//workerApp_3 missing 90 for MinDeploymentCount
	state_cloud.GlobalCloudLayout.Current.AddApp("generalHost_2", "workerApp_3", 10, 5)
	state_cloud.GlobalCloudLayout.Current.AddApp("generalHost_3", "workerApp_3", 10, 5)


	if len(state_cloud.GlobalCloudLayout.Current.Layout) != 7 {
		t.Error(state_cloud.GlobalCloudLayout.Current)
	}
	if len(state_cloud.GlobalCloudLayout.Desired.Layout) != 0 {
		t.Error(state_cloud.GlobalCloudLayout.Desired)
	}

	planner.Plan()

	if len(state_cloud.GlobalCloudLayout.Current.Layout) == 0 {
		t.Error(state_cloud.GlobalCloudLayout.Current)
	}
	if len(state_cloud.GlobalCloudLayout.Desired.Layout) != 7 {
		t.Error(state_cloud.GlobalCloudLayout.Desired)
	}

	//check httpApp_1 update
	host, _ := state_cloud.GlobalCloudLayout.Desired.GetHost("cpuHost_1")
	if  host.Apps["httpApp_1"].Version != 2 ||  host.Apps["httpApp_1"].DeploymentCount != 1 {
		t.Error(state_cloud.GlobalCloudLayout.Desired)
	}

	diff := planner.Diff(state_cloud.GlobalCloudLayout.Desired, state_cloud.GlobalCloudLayout.Current)

	if diff["cpuHost_1"]["httpApp_1"].Version != 2 || diff["cpuHost_1"]["httpApp_1"].DeploymentCount != 1 {
		t.Error(diff["cpuHost_1"])
	}
	if len(diff["emptyHost"]) != 0{
		t.Error(diff["emptyHost"])
	}
	if _, exists := diff["generalHost_1"]["httpApp_2"]; exists {
		t.Error(diff["gernerHost_1"])
	}
	if diff["generalHost_3"]["workerApp_3"].DeploymentCount != 14 {
		t.Error(diff["generalHost_3"]["workerApp_3"])
	}

	instances := state_cloud.GlobalAvailableInstances

	cpuHost_1, _ := instances.GetResources("cpuHost_1")
	if cpuHost_1.TotalCpuResource != 500 || cpuHost_1.UsedCpuResource != 94 {
		t.Error(cpuHost_1)
	}
	cpuHost_2, _ := instances.GetResources("cpuHost_2")
	if cpuHost_2.TotalCpuResource != 501 || cpuHost_2.UsedCpuResource != 310 || cpuHost_2.UsedMemoryResource != 90 {
		t.Error(cpuHost_2)
	}
	memoryHost_1, _ := instances.GetResources("memoryHost_1")
	if memoryHost_1.TotalCpuResource != 200 || memoryHost_1.UsedCpuResource != 197 || memoryHost_1.UsedMemoryResource != 57 {
		t.Error(memoryHost_1)
	}
	generalHost_1, _ := instances.GetResources("generalHost_1")
	if generalHost_1.TotalCpuResource != 101 || generalHost_1.UsedCpuResource != 94 || generalHost_1.UsedMemoryResource != 64 {
		t.Error(generalHost_1)
	}
	generalHost_2, _ := instances.GetResources("generalHost_2")
	if generalHost_2.TotalCpuResource != 102 || generalHost_2.UsedCpuResource != 99 || generalHost_2.UsedMemoryResource != 29 {
		t.Error(generalHost_2)
	}
	emptyHost, _ := instances.GetResources("emptyHost")
	if emptyHost.TotalCpuResource != 1000 || emptyHost.UsedCpuResource != 0 || emptyHost.UsedMemoryResource != 0 || emptyHost.UsedNetworkResource != 0 {
		t.Error(emptyHost)
	}


	// planning worked, now test the responder
	planner.Queue.Apply(diff)


	cpuHost_1Config, _ := GetConfigForHost("cpuHost_1")
	if cpuHost_1Config.DeploymentCount != 4 && cpuHost_1Config.DeploymentCount != 1 {
		t.Errorf("%+v", cpuHost_1Config)
	}

}