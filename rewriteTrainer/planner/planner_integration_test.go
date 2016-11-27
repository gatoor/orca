package planner

import (
	"testing"
	"gatoor/orca/rewriteTrainer/state/configuration"
	"gatoor/orca/rewriteTrainer/state/cloud"
	"gatoor/orca/rewriteTrainer/state/needs"
	"gatoor/orca/rewriteTrainer/config"
	"gatoor/orca/rewriteTrainer/cloud"
	"gatoor/orca/base"
	Logger "gatoor/orca/rewriteTrainer/log"
	"fmt"
	"github.com/Sirupsen/logrus"
	"math/rand"
	"gatoor/orca/rewriteTrainer/needs"
)


func applySampleConfig() {
	conf := config.JsonConfiguration{}

	conf.Trainer.Port = 5000

	httpApp1 := config.AppJsonConfiguration{
		Name: "httpApp_1",
		Version: "http_1.0",
		Type: base.APP_HTTP,
		MinDeploymentCount: 3,
		MaxDeploymentCount: 10,
		InstallCommands: []base.OsCommand{
			{
				Type: base.EXEC_COMMAND,
				Command: base.Command{"ls", "/home"},
			},
			{
				Type: base.FILE_COMMAND,
				Command: base.Command{"/server/app1/app1.conf", "somefilecontent as a string"},
			},
		},
			QueryStateCommand: base.OsCommand{
			Type: base.EXEC_COMMAND,
			Command: base.Command{"wget", "http://localhost:1234/check"},
		},
			RemoveCommand: base.OsCommand{
			Type: base.EXEC_COMMAND,
			Command: base.Command{"rm", "-rf /server/app1"},
		},
			Needs: needs.AppNeeds{
			MemoryNeeds: needs.MemoryNeeds(1),
			CpuNeeds: needs.CpuNeeds(1),
			NetworkNeeds: needs.NetworkNeeds(1),
		},
	}

	httpApp1_v2 := config.AppJsonConfiguration{
		Name: "httpApp_1",
		Version: "http_1.1",
		Type: base.APP_HTTP,
		MinDeploymentCount: 2,
		MaxDeploymentCount: 10,
		InstallCommands: []base.OsCommand{
			{
				Type: base.EXEC_COMMAND,
				Command: base.Command{"ls", "/home"},
			},
			{
				Type: base.FILE_COMMAND,
				Command: base.Command{"/server/app1/app1.conf", "somefilecontent as a string"},
			},
		},
			QueryStateCommand: base.OsCommand{
			Type: base.EXEC_COMMAND,
			Command: base.Command{"wget", "http://localhost:1234/check"},
		},
			RemoveCommand: base.OsCommand{
			Type: base.EXEC_COMMAND,
			Command: base.Command{"rm", "-rf /server/app1"},
		},
			Needs: needs.AppNeeds{
			MemoryNeeds: needs.MemoryNeeds(2),
			CpuNeeds: needs.CpuNeeds(2),
			NetworkNeeds: needs.NetworkNeeds(5),
		},
	}

	httpApp2 := config.AppJsonConfiguration{
		Name: "httpApp_2",
		Version: "http_2.0",
		Type: base.APP_HTTP,
		MinDeploymentCount: 4,
		MaxDeploymentCount: 10,
		InstallCommands: []base.OsCommand{
			{
				Type: base.EXEC_COMMAND,
				Command: base.Command{"ls", "/home"},
			},
			{
				Type: base.FILE_COMMAND,
				Command: base.Command{"/server/app1/app1.conf", "somefilecontent as a string"},
			},
		},
		QueryStateCommand: base.OsCommand{
			Type: base.EXEC_COMMAND,
			Command: base.Command{"wget", "http://localhost:1234/check"},
		},
		RemoveCommand: base.OsCommand{
			Type: base.EXEC_COMMAND,
			Command: base.Command{"rm", "-rf /server/app1"},
		},
		Needs: needs.AppNeeds{
			MemoryNeeds: needs.MemoryNeeds(1),
			CpuNeeds: needs.CpuNeeds(1),
			NetworkNeeds: needs.NetworkNeeds(1),
		},
	}

	workerApp1 := config.AppJsonConfiguration{
		Name: "workerApp_1",
		Version: "worker_1.0",
		Type: base.APP_WORKER,
		MinDeploymentCount: 1,
		MaxDeploymentCount: 1,
		InstallCommands: []base.OsCommand{
			{
				Type: base.EXEC_COMMAND,
				Command: base.Command{"ls", "/home"},
			},
			{
				Type: base.FILE_COMMAND,
				Command: base.Command{"/server/app1/app1.conf", "somefilecontent as a string"},
			},
		},
		QueryStateCommand: base.OsCommand{
			Type: base.EXEC_COMMAND,
			Command: base.Command{"wget", "http://localhost:1234/check"},
		},
		RemoveCommand: base.OsCommand{
			Type: base.EXEC_COMMAND,
			Command: base.Command{"rm", "-rf /server/app1"},
		},
		Needs: needs.AppNeeds{
			CpuNeeds: needs.CpuNeeds(50),
			MemoryNeeds: needs.MemoryNeeds(10),
			NetworkNeeds: needs.NetworkNeeds(10),
		},
	}

	workerApp1_v2 := config.AppJsonConfiguration{
		Name: "workerApp_1",
		Version: "worker_1.1",
		Type: base.APP_WORKER,
		MinDeploymentCount: 1,
		MaxDeploymentCount: 1,
		InstallCommands: []base.OsCommand{
			{
				Type: base.EXEC_COMMAND,
				Command: base.Command{"ls", "/home"},
			},
			{
				Type: base.FILE_COMMAND,
				Command: base.Command{"/server/app1/app1.conf", "somefilecontent as a string"},
			},
		},
		QueryStateCommand: base.OsCommand{
			Type: base.EXEC_COMMAND,
			Command: base.Command{"wget", "http://localhost:1234/check"},
		},
		RemoveCommand: base.OsCommand{
			Type: base.EXEC_COMMAND,
			Command: base.Command{"rm", "-rf /server/app1"},
		},
		Needs: needs.AppNeeds{
			CpuNeeds: needs.CpuNeeds(70),
			MemoryNeeds: needs.MemoryNeeds(40),
			NetworkNeeds: needs.NetworkNeeds(30),
		},
	}

	workerApp2 := config.AppJsonConfiguration{
		Name: "workerApp_2",
		Version: "worker_2.0",
		Type: base.APP_WORKER,
		MinDeploymentCount: 5,
		MaxDeploymentCount: 10,
		InstallCommands: []base.OsCommand{
			{
				Type: base.EXEC_COMMAND,
				Command: base.Command{"ls", "/home"},
			},
			{
				Type: base.FILE_COMMAND,
				Command: base.Command{"/server/app1/app1.conf", "somefilecontent as a string"},
			},
		},
		QueryStateCommand: base.OsCommand{
			Type: base.EXEC_COMMAND,
			Command: base.Command{"wget", "http://localhost:1234/check"},
		},
		RemoveCommand: base.OsCommand{
			Type: base.EXEC_COMMAND,
			Command: base.Command{"rm", "-rf /server/app1"},
		},
		Needs: needs.AppNeeds{
			CpuNeeds: needs.CpuNeeds(23),
			MemoryNeeds: needs.MemoryNeeds(23),
			NetworkNeeds: needs.NetworkNeeds(23),
		},
	}

	workerApp3 := config.AppJsonConfiguration{
		Name: "workerApp_3",
		Version: "worker_3.0",
		Type: base.APP_WORKER,
		MinDeploymentCount: 100,
		MaxDeploymentCount: 200,
		InstallCommands: []base.OsCommand{
			{
				Type: base.EXEC_COMMAND,
				Command: base.Command{"ls", "/home"},
			},
			{
				Type: base.FILE_COMMAND,
				Command: base.Command{"/server/app1/app1.conf", "somefilecontent as a string"},
			},
		},
		QueryStateCommand: base.OsCommand{
			Type: base.EXEC_COMMAND,
			Command: base.Command{"wget", "http://localhost:1234/check"},
		},
		RemoveCommand: base.OsCommand{
			Type: base.EXEC_COMMAND,
			Command: base.Command{"rm", "-rf /server/app1"},
		},
		Needs: needs.AppNeeds{
			CpuNeeds: needs.CpuNeeds(7),
			MemoryNeeds: needs.MemoryNeeds(2),
			NetworkNeeds: needs.NetworkNeeds(1),
		},
	}



	conf.Apps = []config.AppJsonConfiguration{
		httpApp1, httpApp1_v2, httpApp2,
		workerApp1, workerApp1_v2, workerApp2, workerApp3,
	}

	conf.ApplyToState()
}

func initTrainer() {
	state_configuration.GlobalConfigurationState.Init()
	state_cloud.GlobalCloudLayout.Init()
	state_needs.GlobalAppsNeedState = state_needs.AppsNeedState{}
	applySampleConfig()
	initCloudProvider()
	cloud.Init()
}

func initCloudProvider() {
	state_cloud.GlobalAvailableInstances.Update("cpuHost_1", state_cloud.InstanceResources{
		TotalCpuResource: 500,
		TotalMemoryResource: 100,
		TotalNetworkResource: 100,
	})
	state_cloud.GlobalAvailableInstances.Update("cpuHost_2", state_cloud.InstanceResources{
		TotalCpuResource: 501,
		TotalMemoryResource: 101,
		TotalNetworkResource: 101,
	})
	state_cloud.GlobalAvailableInstances.Update("memoryHost_1", state_cloud.InstanceResources{
		TotalCpuResource: 200,
		TotalMemoryResource: 300,
		TotalNetworkResource: 100,
	})
	state_cloud.GlobalAvailableInstances.Update("generalHost_1", state_cloud.InstanceResources{
		TotalCpuResource: 101,
		TotalMemoryResource: 101,
		TotalNetworkResource: 101,
	})
	state_cloud.GlobalAvailableInstances.Update("generalHost_2", state_cloud.InstanceResources{
		TotalCpuResource: 102,
		TotalMemoryResource: 102,
		TotalNetworkResource: 102,
	})
	state_cloud.GlobalAvailableInstances.Update("generalHost_3", state_cloud.InstanceResources{
		TotalCpuResource: 103,
		TotalMemoryResource: 103,
		TotalNetworkResource: 103,
	})
	state_cloud.GlobalAvailableInstances.Update("emptyHost", state_cloud.InstanceResources{
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
	elem , _ := state_needs.GlobalAppsNeedState.Get("httpApp_1", "http_1.1")
	if elem.CpuNeeds != 2 {
		t.Error("wrong needs")
	}
	elem2 , _ := state_needs.GlobalAppsNeedState.Get("workerApp_3", "worker_3.0")
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

func TestPlannerIntegration_initialPlan_NoCurrentLayout(t *testing.T) {
	initTrainer()
	testInit(t)

	if len(state_cloud.GlobalCloudLayout.Current.Layout) != 0 {
		t.Error(state_cloud.GlobalCloudLayout.Current)
	}
	if len(state_cloud.GlobalCloudLayout.Desired.Layout) != 0 {
		t.Error(state_cloud.GlobalCloudLayout.Desired)
	}

	if len(state_cloud.GlobalCloudLayout.Current.Layout) != 0 {
		t.Error(state_cloud.GlobalCloudLayout.Current)
	}
	if len(state_cloud.GlobalCloudLayout.Desired.Layout) != 0 {
		t.Error(state_cloud.GlobalCloudLayout.Desired)
	}

	InitialPlan()

	instances := state_cloud.GlobalAvailableInstances

	cpuHost_1, _ := instances.GetResources("cpuHost_1")
	if cpuHost_1.TotalCpuResource != 500 || cpuHost_1.UsedCpuResource != 350 {
		t.Error(cpuHost_1)
	}
	cpuHost_2, _ := instances.GetResources("cpuHost_2")
	if cpuHost_2.TotalCpuResource != 501 || cpuHost_2.UsedCpuResource != 56 || cpuHost_2.UsedMemoryResource != 16 {
		t.Error(cpuHost_2)
	}
	memoryHost_1, _ := instances.GetResources("memoryHost_1")
	if memoryHost_1.TotalCpuResource != 200 || memoryHost_1.UsedCpuResource != 197 || memoryHost_1.UsedMemoryResource != 57 {
		t.Error(memoryHost_1)
	}
	generalHost_1, _ := instances.GetResources("generalHost_1")
	if generalHost_1.TotalCpuResource != 101 || generalHost_1.UsedCpuResource != 96 || generalHost_1.UsedMemoryResource != 66 {
		t.Error(generalHost_1)
	}
	generalHost_2, _ := instances.GetResources("generalHost_2")
	if generalHost_2.TotalCpuResource != 102 || generalHost_2.UsedCpuResource != 95 || generalHost_2.UsedMemoryResource != 95 {
		t.Error(generalHost_2)
	}
	emptyHost, _ := instances.GetResources("emptyHost")
	if emptyHost.TotalCpuResource != 1000 || emptyHost.UsedCpuResource != 0 || emptyHost.UsedMemoryResource != 0 || emptyHost.UsedNetworkResource != 0 {
		t.Error(emptyHost)
	}

}



func TestPlannerIntegration_regularPlan (t *testing.T) {
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
	state_cloud.GlobalCloudLayout.Current.AddApp("cpuHost_1", "httpApp_1", "http_1.0", 1)
	state_cloud.GlobalCloudLayout.Current.AddApp("cpuHost_2", "httpApp_1", "http_1.1", 1)

	//httpApp_2 missing 1 app for MinDeploymentCount
	state_cloud.GlobalCloudLayout.Current.AddApp("generalHost_1", "httpApp_2", "http_2.0", 1)
	state_cloud.GlobalCloudLayout.Current.AddApp("generalHost_2", "httpApp_2", "http_2.0", 1)
	state_cloud.GlobalCloudLayout.Current.AddApp("generalHost_3", "httpApp_2", "http_2.0", 1)

	//workerApp_1 missing

	//workerApp_2 already deployed MinDeploymentCount all on one host
	state_cloud.GlobalCloudLayout.Current.AddApp("cpuHost_1", "workerApp_2", "worker_2.0", 5)

	//workerApp_3 missing 90 for MinDeploymentCount
	state_cloud.GlobalCloudLayout.Current.AddApp("generalHost_2", "workerApp_3", "worker_3.0", 5)
	state_cloud.GlobalCloudLayout.Current.AddApp("generalHost_3", "workerApp_3", "worker_3.0", 5)


	if len(state_cloud.GlobalCloudLayout.Current.Layout) != 7 {
		t.Error(state_cloud.GlobalCloudLayout.Current)
	}
	if len(state_cloud.GlobalCloudLayout.Desired.Layout) != 0 {
		t.Error(state_cloud.GlobalCloudLayout.Desired)
	}

	Plan()

	if len(state_cloud.GlobalCloudLayout.Current.Layout) == 0 {
		t.Error(state_cloud.GlobalCloudLayout.Current)
	}
	if len(state_cloud.GlobalCloudLayout.Desired.Layout) != 7 {
		t.Error(state_cloud.GlobalCloudLayout.Desired)
	}

	//check httpApp_1 update
	host, _ := state_cloud.GlobalCloudLayout.Desired.GetHost("cpuHost_1")
	if  host.Apps["httpApp_1"].Version != "http_1.1" ||  host.Apps["httpApp_1"].DeploymentCount != 1 {
		t.Error(state_cloud.GlobalCloudLayout.Desired)
	}

	diff := Diff(state_cloud.GlobalCloudLayout.Desired, state_cloud.GlobalCloudLayout.Current)

	/*
	The struct should look something like this (there is some randomness in Desired and Diff though):



	Current = map[
		cpuHost_1:{
			HostId:cpuHost_1 IpAddress: HabitatVersion: Apps:map[
				httpApp_1:{Version:http_1.0 DeploymentCount:1}
				workerApp_2:{Version:worker_2.0 DeploymentCount:5}
			]
		}
		cpuHost_2:{
			HostId:cpuHost_2 IpAddress: HabitatVersion: Apps:map[
				httpApp_1:{Version:http_1.1 DeploymentCount:1}
			]
		}
		memoryHost_1:{
			HostId:memoryHost_1 IpAddress: HabitatVersion: Apps:map[
			]
		}
		generalHost_1:{
			HostId:generalHost_1 IpAddress: HabitatVersion: Apps:map[
				httpApp_2:{Version:http_2.0 DeploymentCount:1}
			]
		}
		generalHost_2:{
			HostId:generalHost_2 IpAddress: HabitatVersion: Apps:map[
				httpApp_2:{Version:http_2.0 DeploymentCount:1}
				workerApp_3:{Version:worker_3.0 DeploymentCount:5}
			]
		}
		generalHost_3:{
			HostId:generalHost_3 IpAddress: HabitatVersion: Apps:map[
				workerApp_3:{Version:worker_3.0 DeploymentCount:5} httpApp_2:{Version:http_2.0 DeploymentCount:1}
			]
		}
	]

	-------------------

	Desired = map[
		cpuHost_1:{HostId:cpuHost_1 IpAddress: HabitatVersion: Apps:map[
			httpApp_1:{Version:http_1.1 DeploymentCount:1}
			workerApp_2:{Version:worker_2.0 DeploymentCount:4}
			]
		}
		cpuHost_2:{HostId:cpuHost_2 IpAddress: HabitatVersion: Apps:map[
			httpApp_1:{Version:http_1.1 DeploymentCount:1}
			workerApp_3:{Version:worker_3.0 DeploymentCount:44}
			]
		}
		memoryHost_1:{HostId:memoryHost_1 IpAddress: HabitatVersion: Apps:map[
			httpApp_2:{Version:http_2.0 DeploymentCount:1} workerApp_3:{Version:worker_3.0 DeploymentCount:28}
			]
		}
		generalHost_1:{HostId:generalHost_1 IpAddress: HabitatVersion: Apps:map[
			httpApp_2:{Version:http_2.0 DeploymentCount:1}
			workerApp_1:{Version:worker_1.1 DeploymentCount:1}
			workerApp_2:{Version:worker_2.0 DeploymentCount:1}
			]
		}
		generalHost_2:{HostId:generalHost_2 IpAddress: HabitatVersion: Apps:map[
			httpApp_2:{Version:http_2.0 DeploymentCount:1}
			workerApp_3:{Version:worker_3.0 DeploymentCount:14}
			]
		}
		generalHost_3:{HostId:generalHost_3 IpAddress: HabitatVersion: Apps:map[
			httpApp_2:{Version:http_2.0 DeploymentCount:1}
			workerApp_3:{Version:worker_3.0 DeploymentCount:14}
			]
		}
		emptyHost:{HostId:emptyHost IpAddress: HabitatVersion: Apps:map[]}
	]


	-------------------

	diff = map[
		cpuHost_1:map[
			workerApp_2:{Version:worker_2.0 DeploymentCount:4} 					<= same app, deployments count differ
			httpApp_1:{Version:http_1.1 DeploymentCount:1}						<= same app new version
		]
		cpuHost_2:map[											<= httpApp_1 is not listed here, it doesn't have to change!
			workerApp_3:{Version:worker_3.0 DeploymentCount:44}
		]
		memoryHost_1:map[										<= all new apps
			httpApp_2:{Version:http_2.0 DeploymentCount:1}
			workerApp_3:{Version:worker_3.0 DeploymentCount:28}
		]
		generalHost_1:map[										<= httpApp_2 is still there and 2 new apps
			workerApp_1:{Version:worker_1.1 DeploymentCount:1}
			workerApp_2:{Version:worker_2.0 DeploymentCount:1}
		]
		generalHost_2:map[
			workerApp_3:{Version:worker_3.0 DeploymentCount:14}
		]
		generalHost_3:map[
			workerApp_3:{Version:worker_3.0 DeploymentCount:14} 					<= added more of the same apps
		]
		emptyHost:map[]
	]

	 */


	if diff["cpuHost_1"]["httpApp_1"].Version != "http_1.1" || diff["cpuHost_1"]["httpApp_1"].DeploymentCount != 1 {
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

	// now generate the queue

	Queue.Apply(diff)

	/*
		Queue should look something like this:

		map[
			cpuHost_1:map[
				httpApp_1:{STATE_QUEUED {http_1.1 1}}
				workerApp_2:{STATE_QUEUED {worker_2.0 4}}
				]
			cpuHost_2:map[
				workerApp_3:{STATE_QUEUED {worker_3.0 44}}
				]
			memoryHost_1:map[
				workerApp_3:{STATE_QUEUED {worker_3.0 28}}
				httpApp_2:{STATE_QUEUED {http_2.0 1}}
				]
			generalHost_1:map[
				workerApp_1:{STATE_QUEUED {worker_1.1 1}}
				workerApp_2:{STATE_QUEUED {worker_2.0 1}}
				]
			generalHost_2:map[
				workerApp_3:{STATE_QUEUED {worker_3.0 14}}
				]
			]
			generalHost_3:map[
				workerApp_3:{STATE_QUEUED {worker_3.0 14}}
				]

	 */

	if len(Queue.Queue) != 6 {
		t.Error(Queue.Queue)
	}

	if Queue.Queue["cpuHost_1"]["httpApp_1"].State != STATE_QUEUED || Queue.Queue["cpuHost_1"]["httpApp_1"].Version.Version != "http_1.1" || Queue.Queue["cpuHost_1"]["httpApp_1"].Version.DeploymentCount != 1 {
		t.Error(Queue.Queue["cpuHost_1"]["httpApp_1"])
	}
	if Queue.Queue["memoryHost_1"]["workerApp_3"].State != STATE_QUEUED || Queue.Queue["memoryHost_1"]["workerApp_3"].Version.Version != "worker_3.0" || Queue.Queue["memoryHost_1"]["workerApp_3"].Version.DeploymentCount != 28 {
		t.Error(Queue.Queue["memoryHost_1"]["workerApp_3"])
	}
}


func TestPlannerIntegration_initialPlan_BigAssDeployment(t *testing.T) {
	Logger.SetLogLevel(logrus.WarnLevel)
	initTrainer()
	testInit(t)

	appCount := 0

	for i:= 1; i <= 4000; i++ {
		state_cloud.GlobalAvailableInstances.Update(base.HostId("filllerHost_" + fmt.Sprint(i)), state_cloud.InstanceResources{
			TotalCpuResource: 500,
			TotalMemoryResource: 100,
			TotalNetworkResource: 100,
		})
	}

	for i:= 1; i <= 50; i++ {
		c := rand.Intn(1000)
		d := rand.Intn(1000)
		appCount += c + d
		state_configuration.GlobalConfigurationState.ConfigureApp(base.AppConfiguration{
			Name: base.AppName("fillerHttp_" + fmt.Sprint(i)),
			Type: base.APP_HTTP,
			Version: "1.0",
			MinDeploymentCount: base.DeploymentCount(c),
		})
		state_needs.GlobalAppsNeedState.UpdateNeeds(base.AppName("fillerHttp_" + fmt.Sprint(i)), "1.0", needs.AppNeeds{
			CpuNeeds: needs.CpuNeeds(rand.Intn(10) + 1), MemoryNeeds: needs.MemoryNeeds(rand.Intn(10) + 1), NetworkNeeds: needs.NetworkNeeds(rand.Intn(10) + 1),
		})

		state_configuration.GlobalConfigurationState.ConfigureApp(base.AppConfiguration{
			Name: base.AppName("fillerWorker_" + fmt.Sprint(i)),
			Type: base.APP_WORKER,
			Version: "1.0",
			MinDeploymentCount: base.DeploymentCount(d),
		})
		state_needs.GlobalAppsNeedState.UpdateNeeds(base.AppName("fillerWorker_" + fmt.Sprint(i)), "1.0", needs.AppNeeds{
			CpuNeeds: needs.CpuNeeds(rand.Intn(10) + 1), MemoryNeeds: needs.MemoryNeeds(rand.Intn(10) + 1), NetworkNeeds: needs.NetworkNeeds(rand.Intn(10) + 1),
		})
	}

	if len(state_cloud.GlobalCloudLayout.Current.Layout) != 0 {
		t.Error(state_cloud.GlobalCloudLayout.Current)
	}
	if len(state_cloud.GlobalCloudLayout.Desired.Layout) != 0 {
		t.Error(state_cloud.GlobalCloudLayout.Desired)
	}

	fmt.Println("")

	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>")
	fmt.Println(state_cloud.GlobalCloudLayout.Desired.Layout)
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>")
	fmt.Println(state_cloud.GlobalAvailableInstances)
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>")
}