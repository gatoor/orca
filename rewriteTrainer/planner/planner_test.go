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

package planner


import (
	"testing"
	"gatoor/orca/base"
	"gatoor/orca/rewriteTrainer/state/cloud"
	"gatoor/orca/rewriteTrainer/state/configuration"
	"gatoor/orca/rewriteTrainer/example"
	"gatoor/orca/rewriteTrainer/state/needs"
	"gatoor/orca/rewriteTrainer/cloud"
	"gatoor/orca/rewriteTrainer/tracker"
	"time"
)

func TestPlannerQueue_AllEmpty(t *testing.T) {
	queue := NewPlannerQueue()

	if queue.AllEmpty() == false {
		t.Error("should be all empty")
	}

	queue.Add("host1", "app1", state_cloud.AppsVersion{1, 2})

	if queue.AllEmpty() == true {
		t.Error("should have elements")
	}
}

func TestPlannerQueue_Empty(t *testing.T) {
	queue := NewPlannerQueue()

	if queue.Empty("somehost") == false {
		t.Error("should be empty")
	}

	queue.Add("host1", "app1", state_cloud.AppsVersion{1, 2})

	if queue.Empty("somehost") == false {
		t.Error("should be empty")
	}

	if queue.Empty("host1") == true {
		t.Error("should have elements")
	}
}

func TestPlannerQueue_Remove(t *testing.T) {
	queue := NewPlannerQueue()

	if queue.Empty("somehost") == false {
		t.Error("should be empty")
	}

	queue.Add("host1", "app1", state_cloud.AppsVersion{1, 2})

	if queue.Empty("host1") == true {
		t.Error("should not be empty")
	}
	elem, _ := queue.Get("host1")
	if elem["app1"].Version.Version != 1 {
		t.Error("wrong version")
	}
	queue.SetState("host1", "app1", STATE_FAIL)

	state, _ := queue.GetState("host1", "app1")
	if state != STATE_FAIL {
		t.Error(queue.GetState("host1", "app1"))
	}

	elem1, _ := queue.Get("host1")
	queue.Remove("host1", "app1")

	_, err2 := queue.GetState("host1", "app1")
	if err2 == nil {
		t.Error(queue.GetState("host1", "app1"))
	}

	if elem1["app1"].Version.Version == 1 {
		t.Error("wrong version")
	}

	if queue.Empty("host1") == false {
		t.Error("should not have elements")
	}
}

func TestPlannerQueue_RemoveApp(t *testing.T) {
	queue := NewPlannerQueue()
	queue.Add("host1", "app1", state_cloud.AppsVersion{1, 1})
	queue.Add("host1", "app2", state_cloud.AppsVersion{3, 1})
	queue.Add("host2", "app1", state_cloud.AppsVersion{2, 1})
	queue.Add("host2", "app2", state_cloud.AppsVersion{3, 1})
	queue.Add("host3", "app1", state_cloud.AppsVersion{1, 1})

	h1, _ := queue.Get("host1")
	if len(h1) != 2 {
		t.Error(h1)
	}
	if h1["app1"].Version.Version != 1 {
		t.Error(h1)
	}
	h2, _ := queue.Get("host2")
	if len(h2) != 2 {
		t.Error(h2)
	}
	if h2["app1"].Version.Version != 2 {
		t.Error(h2)
	}
	h3, _ := queue.Get("host3")
	if len(h3) != 1 {
		t.Error(h3)
	}

	queue.RemoveApp("app1", 1)

	h11, _ := queue.Get("host1")
	if len(h11) != 1 {
		t.Error(h11)
	}
	h22, _ := queue.Get("host2")
	if len(h22) != 2 {
		t.Error(h22)
	}
	h33, _ := queue.Get("host3")
	if len(h33) != 0 {
		t.Error(h33)
	}
}


func TestPlannerQueue_RollbackApp(t *testing.T) {
	queue := NewPlannerQueue()
	queue.Add("host1", "app1", state_cloud.AppsVersion{1, 1})
	queue.Add("host1", "app2", state_cloud.AppsVersion{3, 1})
	queue.Add("host2", "app1", state_cloud.AppsVersion{2, 1})
	queue.Add("host2", "app2", state_cloud.AppsVersion{3, 1})
	queue.Add("host3", "app1", state_cloud.AppsVersion{1, 1})

	h1, _ := queue.Get("host1")
	if len(h1) != 2 {
		t.Error(h1)
	}
	if h1["app1"].Version.Version != 1 {
		t.Error(h1)
	}
	h2, _ := queue.Get("host2")
	if len(h2) != 2 {
		t.Error(h2)
	}
	if h2["app1"].Version.Version != 2 {
		t.Error(h2)
	}
	h3, _ := queue.Get("host3")
	if len(h3) != 1 {
		t.Error(h3)
	}

	queue.RollbackApp("app1", 3, 1)

	h11, _ := queue.Get("host1")
	if len(h11) != 2 || h11["app1"].Version.Version != 1 {
		t.Error(h11)
	}
	h22, _ := queue.Get("host2")
	if len(h22) != 2 {
		t.Error(h22)
	}
	h33, _ := queue.Get("host3")
	if len(h33) != 1 || h33["app1"].Version.Version != 1 {
		t.Error(h33)
	}
}


func TestPlannerQueue_PopSuccessFailState(t *testing.T) {
	queue := NewPlannerQueue()

	queue.Add("host1", "app1", state_cloud.AppsVersion{1, 2})
	queue.Add("host2", "app1", state_cloud.AppsVersion{10, 2})
	queue.Add("host1", "app1", state_cloud.AppsVersion{3, 2})
	queue.Add("host1", "app1", state_cloud.AppsVersion{3, 3})

	if queue.Empty("host1") == true || queue.Empty("host2") == true{
		t.Error("should have elements")
	}

	queue.SetState("host2", "app1", STATE_SUCCESS)

	elem, err := queue.Get("host2")
	if err != nil {
		t.Error("unexpected error")
	}
	if elem["app1"].Version.Version != 10 {
		t.Error("wrong version")
	}

	queue.SetState("host1", "app1", STATE_FAIL)

	elem1, err1 := queue.Get("host1")
	if err1 != nil {
		t.Error("unexpected error")
	}
	if elem1["app1"].Version.Version != 1 {
		t.Error(elem1["app1"].Version.Version)
	}
}


func TestPlannerQueue_PopQueuedApplyingState(t *testing.T) {
	queue := NewPlannerQueue()

	queue.Add("host1", "app1", state_cloud.AppsVersion{1, 2})
	queue.Add("host2", "app1", state_cloud.AppsVersion{10, 2})
	queue.Add("host1", "app1", state_cloud.AppsVersion{3, 2})
	queue.Add("host1", "app1", state_cloud.AppsVersion{3, 3})

	elem1, err1 := queue.Get("host1")
	if err1 != nil {
		t.Error("unexpected error")
	}
	if elem1["app1"].Version.Version != 1 {
		t.Error("wrong version")
	}



	elem2, err2 := queue.Get("host1")
	if err2 != nil {
		t.Error("unexpected error")
	}

	if elem2["app1"].Version.Version != 1 {
		t.Error("wrong version")
	}
	queue.SetState("host1", "app1", STATE_APPLYING)

	elem4, err4 := queue.Get("host1")
	if err4 != nil {
		t.Error("unexpected error")
	}

	if elem4["app1"].Version.Version != 1 {
		t.Error("wrong version")
	}
	queue.SetState("host1", "app1", STATE_FAIL)
	elem3, err3 := queue.Get("host1")
	if err3 != nil {
		t.Error("unexpected error")
	}
	if elem3["app1"].Version.DeploymentCount != 2 {
		t.Error("wrong deployment count")
	}

	if queue.Empty("host1") == true {
		t.Error("pop removed too many elements")
	}
}


func initAppsDiff() (map[base.AppName]state_cloud.AppsVersion, map[base.AppName]state_cloud.AppsVersion) {
	return make(map[base.AppName]state_cloud.AppsVersion), make(map[base.AppName]state_cloud.AppsVersion)
}

func TestPlanner_appsDiff_Equal(t *testing.T) {
	master, slave := initAppsDiff()
	master["app1"] = state_cloud.AppsVersion{1, 1}
	master["app2"] = state_cloud.AppsVersion{2, 2}

	slave["app1"] = state_cloud.AppsVersion{1, 1}
	slave["app2"] = state_cloud.AppsVersion{2, 2}

	diff := appsDiff(master, slave)
	if len(diff) != 0 {
		t.Error("found diff in equal apps")
	}
}


func TestPlanner_appsDiff_Update(t *testing.T) {
	master, slave := initAppsDiff()

	master["app1"] = state_cloud.AppsVersion{3, 1}
	master["app2"] = state_cloud.AppsVersion{21, 2}

	slave["app1"] = state_cloud.AppsVersion{1, 1}
	slave["app2"] = state_cloud.AppsVersion{2, 2}

	diff := appsDiff(master, slave)

	if len(diff) != 2 {
		t.Error("found no diff")
	}

	if diff["app1"].Version != 3 {
		t.Error("wrong version")
	}
	if diff["app2"].Version != 21 {
		t.Error("wrong version")
	}
}


func TestPlanner_appsDiff_ScaleUp(t *testing.T) {
	master, slave := initAppsDiff()

	master["app1"] = state_cloud.AppsVersion{1, 2}
	master["app2"] = state_cloud.AppsVersion{2, 2}

	slave["app1"] = state_cloud.AppsVersion{1, 1}
	slave["app2"] = state_cloud.AppsVersion{2, 2}

	diff := appsDiff(master, slave)

	if len(diff) != 1 {
		t.Error("found no diff")
	}

	if diff["app1"].Version != 1 {
		t.Error("wrong version")
	}
	if diff["app1"].DeploymentCount != 2 {
		t.Error("wrong count")
	}
}


func TestPlanner_appsDiff_DeployNew(t *testing.T) {
	master, slave := initAppsDiff()

	master["app1"] = state_cloud.AppsVersion{1, 1}
	master["app2"] = state_cloud.AppsVersion{2, 2}

	slave["app1"] = state_cloud.AppsVersion{1, 1}

	diff := appsDiff(master, slave)

	if len(diff) != 1 {
		t.Error("found no diff")
	}

	if diff["app2"].Version != 2 {
		t.Error("wrong version")
	}
}


func TestPlanner_appsDiff_RemoveApp(t *testing.T) {
	master, slave := initAppsDiff()

	master["app1"] = state_cloud.AppsVersion{1, 1}

	slave["app1"] = state_cloud.AppsVersion{1, 1}
	slave["app2"] = state_cloud.AppsVersion{2, 2}

	diff := appsDiff(master, slave)
	if len(diff) != 1 {
		t.Error("found no diff")
	}

	if diff["app2"].Version != 2 {
		t.Error("wrong version")
	}

	if diff["app2"].DeploymentCount != 0 {
		t.Error("wrong version")
	}
}


func TestPlanner_appsDiff_RollbackVersion(t *testing.T) {
	master, slave := initAppsDiff()

	master["app1"] = state_cloud.AppsVersion{1, 1}
	master["app2"] = state_cloud.AppsVersion{2, 2}

	slave["app1"] = state_cloud.AppsVersion{1, 1}
	slave["app2"] = state_cloud.AppsVersion{3, 2}

	diff := appsDiff(master, slave)

	if len(diff) != 1 {
		t.Error("found no diff")
	}

	if diff["app2"].Version != 2 {
		t.Error("wrong version")
	}
}


func TestPlanner_appsDiff_ScaleDown(t *testing.T) {
	master, slave := initAppsDiff()

	master["app1"] = state_cloud.AppsVersion{1, 1}
	master["app2"] = state_cloud.AppsVersion{2, 1}

	slave["app1"] = state_cloud.AppsVersion{1, 1}
	slave["app2"] = state_cloud.AppsVersion{2, 2}

	diff := appsDiff(master, slave)

	if len(diff) != 1 {
		t.Error("found no diff")
	}

	if diff["app2"].Version != 2 {
		t.Error("wrong version")
	}
	if diff["app2"].DeploymentCount != 1 {
		t.Error("wrong deployment count")
	}
}

func TestPlanner_appsDiff_Combination(t *testing.T) {
	master, slave := initAppsDiff()

	//do nothing
	master["app1"] = state_cloud.AppsVersion{1, 1}
	slave["app1"] = state_cloud.AppsVersion{1, 1}

	// update
	master["app2"] = state_cloud.AppsVersion{3, 1}
	slave["app2"] = state_cloud.AppsVersion{1, 1}

	//rollback
	master["app3"] = state_cloud.AppsVersion{1, 1}
	slave["app3"] = state_cloud.AppsVersion{3, 1}

	//deploy new
	master["app4"] = state_cloud.AppsVersion{3, 1}

	//remove
	slave["app5"] = state_cloud.AppsVersion{1, 1}

	//rollback scale up
	master["app6"] = state_cloud.AppsVersion{1, 2}
	slave["app6"] = state_cloud.AppsVersion{3, 1}

	//scale down
	master["app7"] = state_cloud.AppsVersion{1, 1}
	slave["app7"] = state_cloud.AppsVersion{1, 5}


	diff := appsDiff(master, slave)

	if len(diff) != 6 {
		t.Error("found no diff")
	}

	if diff["app2"].Version != 3 {
		t.Error("wrong version")
	}
	if diff["app2"].DeploymentCount != 1 {
		t.Error("wrong deployment count")
	}

	if diff["app3"].Version != 1 {
		t.Error("wrong version")
	}
	if diff["app3"].DeploymentCount != 1 {
		t.Error("wrong deployment count")
	}

	if diff["app4"].Version != 3 {
		t.Error("wrong version")
	}
	if diff["app4"].DeploymentCount != 1 {
		t.Error("wrong deployment count")
	}

	if diff["app5"].Version != 1 {
		t.Error("wrong version")
	}
	if diff["app5"].DeploymentCount != 0 {
		t.Error("wrong deployment count")
	}

	if diff["app6"].Version != 1 {
		t.Error("wrong version")
	}
	if diff["app6"].DeploymentCount != 2 {
		t.Error("wrong deployment count")
	}

	if diff["app7"].Version != 1 {
		t.Error("wrong version")
	}
	if diff["app7"].DeploymentCount != 1 {
		t.Error("wrong deployment count")
	}
}

func initDiff() LayoutDiff {
	masterCloud := state_cloud.CloudLayout{Type: "master",Layout: make(map[base.HostId]state_cloud.CloudLayoutElement),}
	slaveCloud := state_cloud.CloudLayout{Type: "slave",Layout: make(map[base.HostId]state_cloud.CloudLayoutElement),}
	masterCloud.AddEmptyHost("host1")
	masterCloud.AddEmptyHost("host2")
	slaveCloud.AddEmptyHost("host1")
	slaveCloud.AddEmptyHost("host2")
	masterCloud.AddApp("host1", "app1", 1, 1)
	slaveCloud.AddApp("host1", "app1", 1, 1)
	masterCloud.AddApp("host1", "app2", 3, 1)
	slaveCloud.AddApp("host1", "app2", 1, 1)
	masterCloud.AddApp("host2", "app1", 1, 10)
	slaveCloud.AddApp("host2", "app1", 1, 2)
	return Diff(masterCloud, slaveCloud)
}

func TestPlanner_Diff(t *testing.T) {
	diff := initDiff()

	if len(diff) != 2 {
		t.Error("no diff")
	}

	if diff["host1"]["app2"].DeploymentCount != 1 && diff["host1"]["app2"].Version != 3 {
		t.Error("host1 wrong diff")
	}
	if diff["host2"]["app1"].DeploymentCount != 10 && diff["host2"]["app1"].Version != 1 {
		t.Error("host1 wrong diff")
	}
}


func TestPlannerQueue_Snapshot(t *testing.T) {
	queue := NewPlannerQueue()

	queue.Add("host1", "app1", state_cloud.AppsVersion{1, 2})
	queue.Add("host2", "app1", state_cloud.AppsVersion{10, 2})

	snapshot := queue.Snapshot()

	queue.RemoveHost("host2")

	if queue.Empty("host2") == false {
		t.Error("RemoveHost didn't remove element")
	}

	//sanity check of empty check below
	if _, exists := queue.Queue["host2"]; exists {
		t.Error("RemoveHost didn't remove elem")
	}

	if _, exists := snapshot["host2"]; !exists {
		t.Error("RemoveHost removed snapshotted elem")
	}
}

func TestPlannerQueue_Apply_empty(t *testing.T) {
	diff := initDiff()
	queue := NewPlannerQueue()
	queue.Apply(diff)

	if queue.Empty("host3") == false {
		t.Error("got unexpected host")
	}

	if queue.Empty("host1") == true {
		t.Error("host1 missing")
	}
	elem, _ := queue.Get("host1")
	if  elem["app2"].Version.DeploymentCount != 1 {
		t.Error("host1 wrong deployment count " + string(elem["app1"].Version.DeploymentCount))
	}
	elem2, _ := queue.Get("host2")
	if  elem2["app1"].Version.DeploymentCount != 10 {
		t.Error("host1 wrong deployment count " + string(elem2["app1"].Version.DeploymentCount))
	}
}

func TestPlannerQueue_Apply(t *testing.T) {
	diff := initDiff()
	queue := NewPlannerQueue()
	queue.Apply(diff)
	queue.Apply(diff)

	if queue.Empty("host3") == false {
		t.Error("got unexpected host")
	}

	if queue.Empty("host1") == true {
		t.Error("host1 missing")
	}
	elem, _ := queue.Get("host1")
	if  elem["app2"].Version.DeploymentCount != 1 {
		t.Error("host1 wrong deployment count " + string(elem["app1"].Version.DeploymentCount))
	}
	if len(elem) != 1 {
		t.Error("too many elements")
	}

	elem2, _ := queue.Get("host2")
	if  elem2["app1"].Version.DeploymentCount != 10 {
		t.Error("host1 wrong deployment count " + string(elem2["app1"].Version.DeploymentCount))
	}

	if len(elem2) != 1 {
		t.Error("too many elements")
	}

	e := diff["host1"]["app2"]
	e.Version = 21
	diff["host1"]["app2"] = e
	queue.Apply(diff)

	elem3, _ := queue.Get("host1")
	if  elem3["app2"].Version.Version != 3 {
		t.Error("host1 applied change even though there was already a change in the queue")
	}
	if len(elem3) != 1 {
		t.Error("too many elements")
	}

	elem4, _ := queue.Get("host2")
	if  elem4["app1"].Version.DeploymentCount != 10 {
		t.Error("host1 wrong deployment count " + string(elem4["app1"].Version.DeploymentCount))
	}

	if len(elem4) != 1 {
		t.Error("too many elements")
	}
}


func TestPlanner_initialPlan(t *testing.T) {
	state_configuration.GlobalConfigurationState.Init()
	state_cloud.GlobalCloudLayout.Init()
	state_needs.GlobalAppsNeedState = make(map[base.AppName]state_needs.AppNeedVersion)

	state_cloud.GlobalAvailableInstances.Update("host1", base.InstanceResources{TotalCpuResource: 10.0, TotalMemoryResource: 10.0, TotalNetworkResource: 10.0,})
	state_cloud.GlobalAvailableInstances.Update("host2", base.InstanceResources{TotalCpuResource: 20.0, TotalMemoryResource: 20.0, TotalNetworkResource: 20.0,})

	config := example.ExampleJsonConfig()
	config.ApplyToState()
	cloud.Init()

	if len(state_configuration.GlobalConfigurationState.Apps) != 4 {
		t.Error("init state_config apps wrong len")
	}
	//if len(state_configuration.GlobalConfigurationState.Habitats) != 2 {
	//	t.Error("init state_config habitats wrong len")
	//}

	if len(state_cloud.GlobalCloudLayout.Current.Layout) != 0 {
		t.Error("init state_cloud current should be empty")
	}
	if len(state_cloud.GlobalCloudLayout.Desired.Layout) != 0 {
		t.Error("init state_cloud desired should be empty")
	}

	if len(state_needs.GlobalAppsNeedState) != 4 {
		t.Error(state_needs.GlobalAppsNeedState)
	}

	InitialPlan()


	if len(state_cloud.GlobalCloudLayout.Current.Layout) != 0 {
		t.Error("init state_cloud current should be empty")
	}
	if len(state_cloud.GlobalCloudLayout.Desired.Layout) == 0 {
		t.Error("init state_cloud desired should have elements")
	}

}

func TestPlanner_getGlobalResources(t *testing.T) {
	state_cloud.GlobalCloudLayout.Init()

	state_cloud.GlobalAvailableInstances.Update("host1", base.InstanceResources{TotalCpuResource: 10.0, TotalMemoryResource: 10.0, TotalNetworkResource: 10.0,})
	state_cloud.GlobalAvailableInstances.Update("host2", base.InstanceResources{TotalCpuResource: 21.0, TotalMemoryResource: 22.0, TotalNetworkResource: 23.0,})

	cpu, mem, net := getGlobalResources()

	if cpu != 31.0 {
		t.Error("wrong cpu resources")
	}
	if mem != 32.0 {
		t.Error("wrong cpu resources")
	}
	if net != 33.0 {
		t.Error("wrong cpu resources")
	}
}

func TestPlanner_getGlobalMinNeeds(t *testing.T) {
	state_configuration.GlobalConfigurationState.Init()
	conf := example.ExampleJsonConfig()

	conf.Apps = append(conf.Apps, base.AppConfiguration{
		Name: "app1",
			Version: 2,
			Type: base.APP_WORKER,
			MinDeploymentCount: 2,
			TargetDeploymentCount: 2,
			Needs: base.AppNeeds{
			CpuNeeds: base.CpuNeeds(11),
			MemoryNeeds: base.MemoryNeeds(22),
			NetworkNeeds: base.NetworkNeeds(33),
		},
	},)

	conf.ApplyToState()

	cpu, mem, net := getGlobalMinNeeds()

	if cpu != 52 {
		t.Error(cpu)
	}
	if mem != 74 {
		t.Error("wrong mem resources")
	}
	if net != 96 {
		t.Error("wrong net resources")
	}
}

func TestPlanner_getGlobalCurrentNeeds(t *testing.T) {
	state_cloud.GlobalCloudLayout.Init()
	state_configuration.GlobalConfigurationState.Init()
	conf := example.ExampleJsonConfig()

	conf.Apps = append(conf.Apps, base.AppConfiguration{
		Name: "app1",
			Version: 2,
			Type: base.APP_WORKER,
			MinDeploymentCount: 2,
			TargetDeploymentCount: 2,
			Needs: base.AppNeeds{
			CpuNeeds: base.CpuNeeds(11),
			MemoryNeeds: base.MemoryNeeds(22),
			NetworkNeeds: base.NetworkNeeds(33),
		},
	},)

	conf.ApplyToState()

	state_cloud.GlobalCloudLayout.Current.AddEmptyHost("host1")
	state_cloud.GlobalCloudLayout.Current.AddEmptyHost("host2")
	state_cloud.GlobalCloudLayout.Current.AddApp("host1", "app1", 2, 10)
	state_cloud.GlobalCloudLayout.Current.AddApp("host1", "app2", 25, 1)
	state_cloud.GlobalCloudLayout.Current.AddApp("host2", "app1", 1, 4)
	state_cloud.GlobalCloudLayout.Current.AddApp("host2", "app11", 1, 20)


	cpu, mem, net := getGlobalCurrentNeeds()

	if cpu != 230 {
		t.Error(cpu)
	}
	if mem != 340 {
		t.Error(mem)
	}
	if net != 450 {
		t.Error(net)
	}
}

func TestPlanner_wipeDesired(t *testing.T) {
	state_cloud.GlobalCloudLayout.Init()
	state_cloud.GlobalCloudLayout.Desired.AddEmptyHost("somehost1")
	state_cloud.GlobalCloudLayout.Desired.AddEmptyHost("somehost2")
	state_cloud.GlobalCloudLayout.Desired.AddApp("somehost1", "app1", 1, 10)
	state_cloud.GlobalCloudLayout.Desired.AddApp("somehost2", "app1", 1, 10)

	res, _ := state_cloud.GlobalCloudLayout.Desired.GetHost("somehost1")

	if res.Apps["app1"].Version != 1 {
		t.Error(res.Apps["app1"].Version)
	}
	res2, _ := state_cloud.GlobalCloudLayout.Desired.GetHost("somehost2")

	if res2.Apps["app1"].DeploymentCount != 10 {
		t.Error(res2.Apps["app1"].DeploymentCount)
	}

	state_cloud.GlobalAvailableInstances.Update("newhost", base.InstanceResources{})
	wipeDesired()

	elem, _ := state_cloud.GlobalCloudLayout.Desired.GetHost("somehost1")

	if len(state_cloud.GlobalCloudLayout.Desired.Layout) != 1 {
		t.Error("should have only one instance")
	}

	if len(elem.Apps) != 0 {
		t.Error("should be empty")
	}

	if elem.HostId != "" {
		t.Error("wrong id, host should not exist")
	}

	elem2, _ := state_cloud.GlobalCloudLayout.Desired.GetHost("newhost")
	if elem2.HostId != "newhost" {
		t.Error("wrong id, host should not exist")
	}
}


func TestPlanner_updateInstanceResources(t *testing.T) {
	state_cloud.GlobalAvailableInstances.Update("host1", base.InstanceResources{
		TotalCpuResource: 100, TotalMemoryResource: 200, TotalNetworkResource: 300, UsedCpuResource: 10, UsedMemoryResource: 20, UsedNetworkResource: 30,
	})

	res, _ := state_cloud.GlobalAvailableInstances.GetResources("host1")

	if res.TotalCpuResource != 100 {
		t.Error(state_cloud.GlobalAvailableInstances)
	}
	if res.UsedCpuResource != 10 {
		t.Error(res.UsedCpuResource)
	}

	updateInstanceResources("host1", base.AppNeeds{CpuNeeds: 1, MemoryNeeds: 2, NetworkNeeds: 3})

	res2, _ := state_cloud.GlobalAvailableInstances.GetResources("host1")

	if res2.TotalCpuResource != 100 {
		t.Error(res.TotalCpuResource)
	}
	if res2.UsedCpuResource != 11 {
		t.Error(res.UsedCpuResource)
	}
	if res2.TotalMemoryResource != 200 {
		t.Error(res.TotalMemoryResource)
	}
	if res2.UsedMemoryResource != 22 {
		t.Error(res.UsedMemoryResource)
	}
	if res2.TotalNetworkResource != 300 {
		t.Error(res.TotalNetworkResource)
	}
	if res2.UsedNetworkResource != 33 {
		t.Error(res.UsedNetworkResource)
	}
}

func TestPlanner_assignAppToHost(t *testing.T) {
	state_cloud.GlobalAvailableInstances.Update("host1", base.InstanceResources{
		TotalCpuResource: 100, TotalMemoryResource: 200, TotalNetworkResource: 300, UsedCpuResource: 10, UsedMemoryResource: 20, UsedNetworkResource: 30,
	})
	state_needs.GlobalAppsNeedState.UpdateNeeds("app1", 1, base.AppNeeds{CpuNeeds: 50, MemoryNeeds: 60, NetworkNeeds: 70})

	appConf := base.AppConfiguration{Name: "app1", Version: 1,}

	res, _ := state_cloud.GlobalAvailableInstances.GetResources("host1")

	if res.TotalCpuResource != 100 {
		t.Error(res.TotalCpuResource)
	}
	if res.UsedCpuResource != 10 {
		t.Error(res.UsedCpuResource)
	}

	assignAppToHost("host1", appConf, 2)


	res2, _ := state_cloud.GlobalAvailableInstances.GetResources("host1")

	if res2.TotalCpuResource != 100 {
		t.Error(res.TotalCpuResource)
	}
	if res2.UsedCpuResource != 110 {
		t.Error(res.UsedCpuResource)
	}
	if res2.UsedMemoryResource != 140 {
		t.Error(res.UsedMemoryResource)
	}
	if res2.UsedNetworkResource != 170 {
		t.Error(res.UsedNetworkResource)
	}

	state_needs.GlobalAppsNeedState.UpdateNeeds("app1", 1, base.AppNeeds{CpuNeeds: 550, MemoryNeeds: 60, NetworkNeeds: 70})
	assignAppToHost("host1", appConf, 2)

	if FailedAssigned[0].AppName != "app1" {
		t.Error("should error due to insufficient resources")
	}
}

func TestPlanner_findHostWithResources_NoCurrent(t *testing.T) {
	state_cloud.GlobalCloudLayout.Init()
	state_cloud.GlobalAvailableInstances.Update("host1", base.InstanceResources{
		TotalCpuResource: 100, TotalMemoryResource: 200, TotalNetworkResource: 300, UsedCpuResource: 90, UsedMemoryResource: 190, UsedNetworkResource: 290,
	})

	state_cloud.GlobalAvailableInstances.Update("host2", base.InstanceResources{
		TotalCpuResource: 100, TotalMemoryResource: 200, TotalNetworkResource: 300, UsedCpuResource: 10, UsedMemoryResource: 20, UsedNetworkResource: 30,
	})

	res := findHostWithResources(base.AppNeeds{CpuNeeds: 20, MemoryNeeds:10, NetworkNeeds: 10}, "", []base.HostId{"host1", "host2"}, nil)

	if res != "host2" {
		t.Error("wrong host")
	}

	res2 := findHostWithResources(base.AppNeeds{CpuNeeds: 20, MemoryNeeds:10, NetworkNeeds: 10000}, "", []base.HostId{"host1", "host2"}, nil)

	if res2 != "" {
		t.Error("wrong host")
	}
}
func TestPlanner_findHostWithResources_WithCurrent(t *testing.T) {
	state_cloud.GlobalCloudLayout.Init()
	state_cloud.GlobalCloudLayout.Current.AddEmptyHost("host1")
	state_cloud.GlobalCloudLayout.Current.AddApp("host1", "app1", 1, 2)
	state_cloud.GlobalCloudLayout.Current.AddEmptyHost("host2")
	state_cloud.GlobalCloudLayout.Current.AddApp("host2", "app2", 3, 2)
	state_cloud.GlobalCloudLayout.Current.AddEmptyHost("host3")
	state_cloud.GlobalCloudLayout.Current.AddApp("host3", "app3", 3, 2)

	state_cloud.GlobalAvailableInstances.Update("host1", base.InstanceResources{
		TotalCpuResource: 10.0, TotalMemoryResource: 20.0, TotalNetworkResource: 30.0, UsedCpuResource: 1.0, UsedMemoryResource: 2.0, UsedNetworkResource: 3.0,
	})

	state_cloud.GlobalAvailableInstances.Update("host2", base.InstanceResources{
		TotalCpuResource: 10.0, TotalMemoryResource: 20.0, TotalNetworkResource: 30.0, UsedCpuResource: 1.0, UsedMemoryResource: 2.0, UsedNetworkResource: 3.0,
	})

	state_cloud.GlobalAvailableInstances.Update("host3", base.InstanceResources{
		TotalCpuResource: 1000.0, TotalMemoryResource: 2000.0, TotalNetworkResource: 3000.0, UsedCpuResource: 1.0, UsedMemoryResource: 2.0, UsedNetworkResource: 3.0,
	})

	res := findHostWithResources(base.AppNeeds{CpuNeeds: 2.0, MemoryNeeds:1.0, NetworkNeeds: 1.0}, "app1", []base.HostId{"host1", "host2", "host3"}, state_cloud.GlobalCloudLayout.Current.FindHostsWithApp("app1"))

	if res != "host1" {
		t.Error("wrong host")
	}

	res2 := findHostWithResources(base.AppNeeds{CpuNeeds: 2.0, MemoryNeeds:1.0, NetworkNeeds: 1.0}, "app2", []base.HostId{"host1", "host2", "host3"}, state_cloud.GlobalCloudLayout.Current.FindHostsWithApp("app2"))

	if res2 != "host2" {
		t.Error("wrong host")
	}
	res3 := findHostWithResources(base.AppNeeds{CpuNeeds: 200.0, MemoryNeeds:100.0, NetworkNeeds: 100.0}, "unkown", []base.HostId{"host1", "host2", "host3"}, state_cloud.GlobalCloudLayout.Current.FindHostsWithApp("unkown"))

	if res3 != "host3" {
		t.Error("wrong host")
	}
}

func TestPlanner_findHttpHostWithResources_NoCurrent(t *testing.T) {
	state_cloud.GlobalCloudLayout.Init()
	state_cloud.GlobalCloudLayout.Desired.AddEmptyHost("host1")
	state_cloud.GlobalCloudLayout.Desired.AddApp("host1", "app1", 1, 2)
	state_cloud.GlobalCloudLayout.Desired.AddEmptyHost("host2")

	state_cloud.GlobalAvailableInstances.Update("host1", base.InstanceResources{
		TotalCpuResource: 10.0, TotalMemoryResource: 20.0, TotalNetworkResource: 30.0, UsedCpuResource: 1.0, UsedMemoryResource: 1.0, UsedNetworkResource: 2.0,
	})

	state_cloud.GlobalAvailableInstances.Update("host2", base.InstanceResources{
		TotalCpuResource: 10.0, TotalMemoryResource: 20.0, TotalNetworkResource: 30.0, UsedCpuResource: 1.0, UsedMemoryResource: 2.0, UsedNetworkResource: 3.0,
	})

	res := findHttpHostWithResources(base.AppNeeds{CpuNeeds: 2.0, MemoryNeeds:1.0, NetworkNeeds: 1.0}, "app1", []base.HostId{"host1", "host2", "host3"}, nil)

	if res != "host2" {
		t.Error(res)
	}

	res2 := findHttpHostWithResources(base.AppNeeds{CpuNeeds: 2.0, MemoryNeeds:1.0, NetworkNeeds: 1000.0}, "app1", []base.HostId{"host1", "host2", "host3"}, nil)

	if res2 != "" {
		t.Error("wrong host")
	}
}

func TestPlanner_findHttpHostWithResources_WithCurrent(t *testing.T) {
	state_cloud.GlobalCloudLayout.Init()
	state_cloud.GlobalCloudLayout.Current.AddEmptyHost("host1")
	state_cloud.GlobalCloudLayout.Current.AddApp("host1", "app1", 1, 2)
	state_cloud.GlobalCloudLayout.Current.AddEmptyHost("host2")
	state_cloud.GlobalCloudLayout.Current.AddApp("host2", "app2", 3, 2)
	state_cloud.GlobalCloudLayout.Current.AddEmptyHost("host3")
	state_cloud.GlobalCloudLayout.Current.AddApp("host3", "app3", 3, 2)

	state_cloud.GlobalCloudLayout.Desired.AddEmptyHost("host1")
	state_cloud.GlobalCloudLayout.Desired.AddApp("host1", "app1", 1, 2)
	state_cloud.GlobalCloudLayout.Desired.AddEmptyHost("host2")
	state_cloud.GlobalCloudLayout.Desired.AddEmptyHost("host3")

	state_cloud.GlobalAvailableInstances.Update("host1", base.InstanceResources{
		TotalCpuResource: 10.0, TotalMemoryResource: 20.0, TotalNetworkResource: 30.0, UsedCpuResource: 1.0, UsedMemoryResource: 2.0, UsedNetworkResource: 3.0,
	})

	state_cloud.GlobalAvailableInstances.Update("host2", base.InstanceResources{
		TotalCpuResource: 10.0, TotalMemoryResource: 20.0, TotalNetworkResource: 30.0, UsedCpuResource: 1.0, UsedMemoryResource: 2.0, UsedNetworkResource: 3.0,
	})

	state_cloud.GlobalAvailableInstances.Update("host3", base.InstanceResources{
		TotalCpuResource: 1000.0, TotalMemoryResource: 2000.0, TotalNetworkResource: 3000.0, UsedCpuResource: 1.0, UsedMemoryResource: 2.0, UsedNetworkResource: 3.0,
	})

	res := findHttpHostWithResources(base.AppNeeds{CpuNeeds: 2.0, MemoryNeeds:1.0, NetworkNeeds: 1.0}, "app1", []base.HostId{"host1", "host2", "host3"}, state_cloud.GlobalCloudLayout.Current.FindHostsWithApp("app1"))

	if res == "host1" {
		t.Error("wrong host")
	}

	res2 := findHttpHostWithResources(base.AppNeeds{CpuNeeds: 2.0, MemoryNeeds:1.0, NetworkNeeds: 1.0}, "app2", []base.HostId{"host1", "host2", "host3"}, state_cloud.GlobalCloudLayout.Current.FindHostsWithApp("app2"))

	if res2 != "host2" {
		t.Error(res2)
	}
	res3 := findHttpHostWithResources(base.AppNeeds{CpuNeeds: 200.0, MemoryNeeds:100.0, NetworkNeeds: 100.0}, "unkown", []base.HostId{"host1", "host2", "host3"}, state_cloud.GlobalCloudLayout.Current.FindHostsWithApp("unkown"))

	if res3 != "host3" {
		t.Error(res3)
	}
}


func TestPlanner_maxDeploymentOnHost(t *testing.T) {
	resources := base.InstanceResources{
		TotalCpuResource: 10.0, TotalNetworkResource: 10.0, TotalMemoryResource: 10.0,
		UsedCpuResource: 1.0, UsedNetworkResource: 1.0, UsedMemoryResource: 1.0,
	}
	ns := base.AppNeeds{CpuNeeds:1.0, MemoryNeeds: 1.0, NetworkNeeds: 1.0}

	res := maxDeploymentOnHost(resources, ns)
	if res != 9 {
		t.Error("wrong res")
	}

	needs2 := base.AppNeeds{CpuNeeds:6.0, MemoryNeeds: 1.0, NetworkNeeds: 1.0}
	res2 := maxDeploymentOnHost(resources, needs2)
	if res2 != 1 {
		t.Error("wrong res")
	}

	needs3 := base.AppNeeds{CpuNeeds:60.0, MemoryNeeds: 1.0, NetworkNeeds: 1.0}
	res3 := maxDeploymentOnHost(resources, needs3)
	if res3 != 0 {
		t.Error("wrong res")
	}
	needs4 := base.AppNeeds{CpuNeeds:1.0, MemoryNeeds: 4.0, NetworkNeeds: 1.0}
	res4 := maxDeploymentOnHost(resources, needs4)
	if res4 != 2 {
		t.Error(res4)
	}
	needs5 := base.AppNeeds{CpuNeeds:1.0, MemoryNeeds: 1.0, NetworkNeeds: 2.0}
	res5 := maxDeploymentOnHost(resources, needs5)
	if res5 != 4 {
		t.Error(res5)
	}
}


func TestPlanner_planHttp_moreInstancesThanNeeded(t *testing.T) {
	state_cloud.GlobalCloudLayout.Init()
	FailedAssigned = []FailedAssign{}
	MissingAssigned = []MissingAssign{}
	resources := base.InstanceResources{
		TotalCpuResource: 10.0, TotalNetworkResource: 10.0, TotalMemoryResource: 10.0,
	}

	state_cloud.GlobalAvailableInstances.Update("host1", resources)
	state_cloud.GlobalAvailableInstances.Update("host2", resources)
	state_cloud.GlobalAvailableInstances.Update("host3", resources)

	state_needs.GlobalAppsNeedState.UpdateNeeds("app1", 1, base.AppNeeds{CpuNeeds: 1.0, MemoryNeeds: 2.0, NetworkNeeds: 3.0})

	appObj := base.AppConfiguration{
		Name: "app1", Version: 1, TargetDeploymentCount: 2,
	}

	called := 0
	hostF := func (needs base.AppNeeds, app base.AppName, hosts []base.HostId, g map[base.HostId]bool) base.HostId {
		if called == 1{
			called += 1
			return "host2"
		}
		if called >= 2 {
			return "host3"
		}
		called += 1
		return "host1"
	}

	res := planHttp(appObj, hostF, false)

	if !res {
		t.Error("should return true")
	}

	host1Res, _ := state_cloud.GlobalAvailableInstances.GetResources("host1")

	if host1Res.UsedCpuResource != 1.0 {
		t.Error(host1Res.UsedCpuResource)
	}
	if host1Res.UsedMemoryResource != 2.0 {
		t.Error(host1Res.UsedMemoryResource)
	}
	if host1Res.UsedNetworkResource != 3.0 {
		t.Error(host1Res.UsedNetworkResource)
	}

	host2Res, _ := state_cloud.GlobalAvailableInstances.GetResources("host2")

	if host2Res.UsedCpuResource != 1.0 {
		t.Error(host2Res.UsedCpuResource)
	}
	if host2Res.UsedMemoryResource != 2.0 {
		t.Error(host2Res.UsedMemoryResource)
	}
	if host2Res.UsedNetworkResource != 3.0 {
		t.Error(host2Res.UsedNetworkResource)
	}

	host3Res, _ := state_cloud.GlobalAvailableInstances.GetResources("host3")

	if host3Res.UsedCpuResource != 0.0 {
		t.Error(host3Res.UsedCpuResource)
	}
	if host3Res.UsedMemoryResource != 0.0 {
		t.Error(host3Res.UsedMemoryResource)
	}
	if host3Res.UsedNetworkResource != 0.0 {
		t.Error(host3Res.UsedNetworkResource)
	}

	if len(MissingAssigned) != 0 {
		t.Error(MissingAssigned)
	}
	if len(FailedAssigned) != 0 {
		t.Error(FailedAssigned)
	}
}

func TestPlanner_planHttp_lessInstancesThanNeeded(t *testing.T) {
	state_cloud.GlobalCloudLayout.Init()
	FailedAssigned = []FailedAssign{}
	MissingAssigned = []MissingAssign{}
	resources := base.InstanceResources{
		TotalCpuResource: 10, TotalNetworkResource: 10, TotalMemoryResource: 10,
	}

	state_cloud.GlobalAvailableInstances.Update("host1", resources)

	state_needs.GlobalAppsNeedState.UpdateNeeds("app1", 1, base.AppNeeds{CpuNeeds: 1.0, MemoryNeeds: 2.0, NetworkNeeds: 3.0})

	appObj := base.AppConfiguration{
		Name: "app1", Version: 1, TargetDeploymentCount: 8,
	}

	called := 0
	hostF := func (ns base.AppNeeds, app base.AppName, hosts []base.HostId, g map[base.HostId]bool) base.HostId {
		if called == 0 {
			called += 1
			return "host1"
		}
		return ""
	}

	res := planHttp(appObj, hostF, false)

	if res {
		t.Error("should return false")
	}

	host1Res, _ := state_cloud.GlobalAvailableInstances.GetResources("host1")

	if host1Res.UsedCpuResource != 1.0 {
		t.Error(host1Res.UsedCpuResource)
	}
	if host1Res.UsedMemoryResource != 2.0 {
		t.Error(host1Res.UsedMemoryResource)
	}
	if host1Res.UsedNetworkResource != 3.0 {
		t.Error(host1Res.UsedNetworkResource)
	}

	if MissingAssigned[0].AppName != "app1" {
		t.Error(MissingAssigned)
	}
	if MissingAssigned[0].DeploymentCount != 7 {
		t.Error(MissingAssigned)
	}

	if len(FailedAssigned) != 0 {
		t.Error(FailedAssigned)
	}
}

func TestPlanner_planWorker_moreInstancesThanNeeded(t *testing.T) {
	state_cloud.GlobalCloudLayout.Init()
	FailedAssigned = []FailedAssign{}
	MissingAssigned = []MissingAssign{}
	resources := base.InstanceResources{
		TotalCpuResource: 10.0, TotalNetworkResource: 10.0, TotalMemoryResource: 10.0,
	}


	state_cloud.GlobalAvailableInstances.Update("host1", resources)
	state_cloud.GlobalAvailableInstances.Update("host2", resources)
	state_cloud.GlobalAvailableInstances.Update("host3", resources)

	state_needs.GlobalAppsNeedState.UpdateNeeds("app1", 1, base.AppNeeds{CpuNeeds: 1.0, MemoryNeeds: 2.0, NetworkNeeds: 3.0})

	appObj := base.AppConfiguration{
		Name: "app1", Version: 1, TargetDeploymentCount: 5,
	}

	called := false
	hostF := func (ns base.AppNeeds, app base.AppName, hosts []base.HostId, g map[base.HostId]bool) base.HostId {
		if called {
			return "host2"
		}
		called = true
		return "host1"
	}

	res := planWorker(appObj, hostF, false)

	if !res {
		t.Error("should return true")
	}

	host1Res, _ := state_cloud.GlobalAvailableInstances.GetResources("host1")

	if host1Res.UsedCpuResource != 3.0 {
		t.Error(host1Res.UsedCpuResource)
	}
	if host1Res.UsedMemoryResource != 6.0 {
		t.Error(host1Res.UsedMemoryResource)
	}
	if host1Res.UsedNetworkResource != 9.0 {
		t.Error(host1Res.UsedNetworkResource)
	}

	host2Res, _ := state_cloud.GlobalAvailableInstances.GetResources("host2")

	if host2Res.UsedCpuResource != 2.0 {
		t.Error(host2Res.UsedCpuResource)
	}
	if host2Res.UsedMemoryResource != 4.0 {
		t.Error(host2Res.UsedMemoryResource)
	}
	if host2Res.UsedNetworkResource != 6.0 {
		t.Error(host2Res.UsedNetworkResource)
	}

	host3Res, _ := state_cloud.GlobalAvailableInstances.GetResources("host3")

	if host3Res.UsedCpuResource != 0.0 {
		t.Error(host3Res.UsedCpuResource)
	}
	if host3Res.UsedMemoryResource != 0.0 {
		t.Error(host3Res.UsedMemoryResource)
	}
	if host3Res.UsedNetworkResource != 0.0 {
		t.Error(host3Res.UsedNetworkResource)
	}
	if len(MissingAssigned) != 0 {
		t.Error(MissingAssigned)
	}
	if len(FailedAssigned) != 0 {
		t.Error(FailedAssigned)
	}
}

func TestPlanner_planWorker_lessInstancesThanNeeded(t *testing.T) {
	state_cloud.GlobalCloudLayout.Init()
	FailedAssigned = []FailedAssign{}
	MissingAssigned = []MissingAssign{}
	resources := base.InstanceResources{
		TotalCpuResource: 10.0, TotalNetworkResource: 10.0, TotalMemoryResource: 10.0,
	}

	state_cloud.GlobalAvailableInstances.Update("host1", resources)

	state_needs.GlobalAppsNeedState.UpdateNeeds("app1", 1, base.AppNeeds{CpuNeeds: 1.0, MemoryNeeds: 2.0, NetworkNeeds: 3.0})

	appObj := base.AppConfiguration{
		Name: "app1", Version: 1, TargetDeploymentCount: 20,
	}

	called := false
	hostF := func (ns base.AppNeeds, app base.AppName, hosts []base.HostId, g map[base.HostId]bool) base.HostId {
		if called {
			return "host2"
		}
		called = true
		return "host1"
	}

	res := planWorker(appObj, hostF, false)

	if res {
		t.Error("should return false")
	}

	host1Res, _ := state_cloud.GlobalAvailableInstances.GetResources("host1")

	if host1Res.UsedCpuResource != 3.0 {
		t.Error(host1Res.UsedCpuResource)
	}
	if host1Res.UsedMemoryResource != 6.0 {
		t.Error(host1Res.UsedMemoryResource)
	}
	if host1Res.UsedNetworkResource != 9.0 {
		t.Error(host1Res.UsedNetworkResource)
	}

	if MissingAssigned[0].AppName != "app1" {
		t.Error(MissingAssigned)
	}
	if MissingAssigned[0].DeploymentCount != 17 {
		t.Error(MissingAssigned)
	}

	if len(FailedAssigned) != 0 {
		t.Error(FailedAssigned)
	}
}


func TestPlanner_sortByTotalNeeds(t *testing.T) {
	state_needs.GlobalAppsNeedState.UpdateNeeds("app1", 1, base.AppNeeds{CpuNeeds: 1.0, MemoryNeeds: 1.0, NetworkNeeds: 1.0,})
	state_needs.GlobalAppsNeedState.UpdateNeeds("app2", 3, base.AppNeeds{CpuNeeds: 2.0, MemoryNeeds: 2.0, NetworkNeeds: 2.0,})
	state_needs.GlobalAppsNeedState.UpdateNeeds("app3", 30, base.AppNeeds{CpuNeeds: 1.0, MemoryNeeds: 2.0, NetworkNeeds: 2.0,})
	state_needs.GlobalAppsNeedState.UpdateNeeds("app4", 40, base.AppNeeds{CpuNeeds: 10.0, MemoryNeeds: 1.0, NetworkNeeds: 1.0,})

	apps := make(map[base.AppName]base.AppConfiguration)

	apps["app1"] = base.AppConfiguration{Name: "app1", Version: 1}
	apps["app2"] = base.AppConfiguration{Name: "app2", Version: 3}
	apps["app3"] = base.AppConfiguration{Name: "app3", Version: 30}
	apps["app4"] = base.AppConfiguration{Name: "app4", Version: 40}

	sorted := sortByTotalNeeds(apps)

	if len(sorted) != 4 {
		t.Error("wrong len")
	}

	if sorted[0] != "app4" || sorted[1] != "app2" || sorted[2] != "app3" || sorted[3] != "app1" {
		t.Error(sorted)
	}
}


func TestPlanner_appPlanningOrder(t *testing.T) {
	state_needs.GlobalAppsNeedState.UpdateNeeds("httpApp1", 1, base.AppNeeds{CpuNeeds: 1.0, MemoryNeeds: 1.0, NetworkNeeds: 1.0,})
	state_needs.GlobalAppsNeedState.UpdateNeeds("httpApp2", 3, base.AppNeeds{CpuNeeds: 2.0, MemoryNeeds: 2.0, NetworkNeeds: 2.0,})
	state_needs.GlobalAppsNeedState.UpdateNeeds("httpApp3", 30, base.AppNeeds{CpuNeeds: 1.0, MemoryNeeds: 2.0, NetworkNeeds: 2.0,})
	state_needs.GlobalAppsNeedState.UpdateNeeds("httpApp4", 40, base.AppNeeds{CpuNeeds: 10.0, MemoryNeeds: 1.0, NetworkNeeds: 1.0,})
	state_needs.GlobalAppsNeedState.UpdateNeeds("workerApp1", 1, base.AppNeeds{CpuNeeds: 10.0, MemoryNeeds: 10.0, NetworkNeeds: 10.0,})
	state_needs.GlobalAppsNeedState.UpdateNeeds("workerApp2", 3, base.AppNeeds{CpuNeeds: 20.0, MemoryNeeds: 20.0, NetworkNeeds: 20.0,})
	state_needs.GlobalAppsNeedState.UpdateNeeds("workerApp3", 30, base.AppNeeds{CpuNeeds: 10.0, MemoryNeeds: 20.0, NetworkNeeds: 20.0,})
	state_needs.GlobalAppsNeedState.UpdateNeeds("workerApp4", 40, base.AppNeeds{CpuNeeds: 100.0, MemoryNeeds: 10.0, NetworkNeeds: 10.0,})

	apps := make(map[base.AppName]base.AppConfiguration)

	apps["httpApp1"] = base.AppConfiguration{Name: "httpApp1", Version: 1, Type: base.APP_HTTP}
	apps["httpApp2"] = base.AppConfiguration{Name: "httpApp2", Version: 3, Type: base.APP_HTTP}
	apps["httpApp3"] = base.AppConfiguration{Name: "httpApp3", Version: 30, Type: base.APP_HTTP}
	apps["httpApp4"] = base.AppConfiguration{Name: "httpApp4", Version: 40, Type: base.APP_HTTP}
	apps["workerApp1"] = base.AppConfiguration{Name: "workerApp1", Version: 1, Type: base.APP_WORKER}
	apps["workerApp2"] = base.AppConfiguration{Name: "workerApp2", Version: 3, Type: base.APP_WORKER}
	apps["workerApp3"] = base.AppConfiguration{Name: "workerApp3", Version: 30, Type: base.APP_WORKER}
	apps["workerApp4"] = base.AppConfiguration{Name: "workerApp4", Version: 40, Type: base.APP_WORKER}

	sortedHttp, sortedWorker := appPlanningOrder(apps)

	if len(sortedHttp) != 4 {
		t.Error("wrong len")
	}

	if sortedHttp[0] != "httpApp4" || sortedHttp[1] != "httpApp2" || sortedHttp[2] != "httpApp3" || sortedHttp[3] != "httpApp1" {
		t.Error(sortedHttp)
	}

	if len(sortedWorker) != 4 {
		t.Error("wrong len")
	}

	if sortedWorker[0] != "workerApp4" || sortedWorker[1] != "workerApp2" || sortedWorker[2] != "workerApp3" || sortedWorker[3] != "workerApp1" {
		t.Error(sortedWorker)
	}

}


func TestPlanner_sortByAvailableResources(t *testing.T) {
	HOST_UPPER_WATERMARK = 0.1
	state_cloud.GlobalCloudLayout.Init()
	state_configuration.GlobalConfigurationState.Trainer.Policies.TRY_TO_REMOVE_HOSTS = true
	state_cloud.GlobalAvailableInstances.Update("host1", base.InstanceResources{
		TotalCpuResource: 100, TotalMemoryResource: 100, TotalNetworkResource: 100,
		UsedCpuResource: 100, UsedMemoryResource: 100, UsedNetworkResource: 100,
	})
	state_cloud.GlobalAvailableInstances.Update("host2", base.InstanceResources{
		TotalCpuResource: 1000, TotalMemoryResource: 100, TotalNetworkResource: 100,
		UsedCpuResource: 0, UsedMemoryResource: 0, UsedNetworkResource: 0,
	})
	state_cloud.GlobalAvailableInstances.Update("host3", base.InstanceResources{
		TotalCpuResource: 100, TotalMemoryResource: 100, TotalNetworkResource: 100,
		UsedCpuResource: 0, UsedMemoryResource: 50, UsedNetworkResource: 0,
	})
	state_cloud.GlobalAvailableInstances.Update("host4", base.InstanceResources{
		TotalCpuResource: 100, TotalMemoryResource: 100, TotalNetworkResource: 100,
		UsedCpuResource: 0, UsedMemoryResource: 50, UsedNetworkResource: 50,
	})

	state_cloud.GlobalAvailableInstances.Update("host5", base.InstanceResources{
		TotalCpuResource: 1000, TotalMemoryResource: 100, TotalNetworkResource: 100,
		UsedCpuResource: 0, UsedMemoryResource: 100, UsedNetworkResource: 0,
	})

	sorted := sortByAvailableResources()

	if len(sorted) != 3 {
		t.Error("wrong len")
	}

	if sorted[0] != "host4" || sorted[1] != "host3" || sorted[2] != "host2" {
		t.Error(sorted)
	}
}

func TestPlanner_sortByAvailableResources_reverse(t *testing.T) {
	HOST_UPPER_WATERMARK = 0.1
	state_cloud.GlobalCloudLayout.Init()
	state_configuration.GlobalConfigurationState.Trainer.Policies.TRY_TO_REMOVE_HOSTS = false
	state_cloud.GlobalAvailableInstances.Update("host1", base.InstanceResources{
		TotalCpuResource: 100, TotalMemoryResource: 100, TotalNetworkResource: 100,
		UsedCpuResource: 100, UsedMemoryResource: 100, UsedNetworkResource: 100,
	})
	state_cloud.GlobalAvailableInstances.Update("host2", base.InstanceResources{
		TotalCpuResource: 1000, TotalMemoryResource: 100, TotalNetworkResource: 100,
		UsedCpuResource: 0, UsedMemoryResource: 0, UsedNetworkResource: 0,
	})
	state_cloud.GlobalAvailableInstances.Update("host3", base.InstanceResources{
		TotalCpuResource: 100, TotalMemoryResource: 100, TotalNetworkResource: 100,
		UsedCpuResource: 0, UsedMemoryResource: 50, UsedNetworkResource: 0,
	})
	state_cloud.GlobalAvailableInstances.Update("host4", base.InstanceResources{
		TotalCpuResource: 100, TotalMemoryResource: 100, TotalNetworkResource: 100,
		UsedCpuResource: 0, UsedMemoryResource: 50, UsedNetworkResource: 50,
	})

	state_cloud.GlobalAvailableInstances.Update("host5", base.InstanceResources{
		TotalCpuResource: 1000, TotalMemoryResource: 100, TotalNetworkResource: 100,
		UsedCpuResource: 0, UsedMemoryResource: 100, UsedNetworkResource: 0,
	})

	sorted := sortByAvailableResources()

	if len(sorted) != 3 {
		t.Error("wrong len")
	}

	if sorted[0] != "host2" || sorted[1] != "host3" || sorted[2] != "host4" {
		t.Error(sorted)
	}
}

func TestPlanner_checkWatermark(t *testing.T) {
	resources := base.InstanceResources{
		TotalCpuResource: 100,
		TotalMemoryResource: 100,
		TotalNetworkResource: 100,
		UsedCpuResource: 100,
		UsedMemoryResource: 100,
		UsedNetworkResource: 100,
	}

	if checkWatermark(resources,"") {
		t.Error("should be false")
	}

	resources1 := base.InstanceResources{
		TotalCpuResource: 100,
		TotalMemoryResource: 100,
		TotalNetworkResource: 100,
		UsedCpuResource: 0,
		UsedMemoryResource: 70,
		UsedNetworkResource: 0,
	}

	if !checkWatermark(resources1,"") {
		t.Error("should be true")
	}

	resources2 := base.InstanceResources{
		TotalCpuResource: 100,
		TotalMemoryResource: 100,
		TotalNetworkResource: 100,
		UsedCpuResource: 90,
		UsedMemoryResource: 0,
		UsedNetworkResource: 0,
	}

	if checkWatermark(resources2,"") {
		t.Error("should be false")
	}

	resources3 := base.InstanceResources{
		TotalCpuResource: 100,
		TotalMemoryResource: 100,
		TotalNetworkResource: 100,
		UsedCpuResource: 0,
		UsedMemoryResource: 94,
		UsedNetworkResource: 0,
	}

	if checkWatermark(resources3,"") {
		t.Error("should be false")
	}

	resources4 := base.InstanceResources{
		TotalCpuResource: 100,
		TotalMemoryResource: 100,
		TotalNetworkResource: 1000000,
		UsedCpuResource: 0,
		UsedMemoryResource: 0,
		UsedNetworkResource: 999999,
	}

	if checkWatermark(resources4,"") {
		t.Error("should be false")
	}
}

func TestPlanner_handleFailedAssign_http(t *testing.T) {
	state_cloud.GlobalCloudLayout.Init()
	state_configuration.GlobalConfigurationState.Init()

	state_cloud.GlobalAvailableInstances.Update("host1", base.InstanceResources{
		TotalCpuResource: 8.0, TotalNetworkResource: 10.0, TotalMemoryResource: 10.0,
	})
	state_cloud.GlobalAvailableInstances.Update("host2", base.InstanceResources{
		TotalCpuResource: 9.0, TotalNetworkResource: 10.0, TotalMemoryResource: 10.0,
	})
	state_cloud.GlobalAvailableInstances.Update("host3", base.InstanceResources{
		TotalCpuResource: 10.0, TotalNetworkResource: 10.0, TotalMemoryResource: 10.0,
	})

	state_needs.GlobalAppsNeedState.UpdateNeeds("app1", 1, base.AppNeeds{CpuNeeds: 1.0, MemoryNeeds: 2.0, NetworkNeeds: 3.0})

	state_configuration.GlobalConfigurationState.ConfigureApp(base.AppConfiguration{
		Name: "app1",
		Type: base.APP_HTTP,
		Version: 1,
		TargetDeploymentCount: 100,
	})

	wipeDesired()
	FailedAssigned = []FailedAssign{
		FailedAssign{AppName: "app1", AppVersion: 1, DeploymentCount: 2, TargetHost: "somehost",},
	}

	if len(state_cloud.GlobalCloudLayout.Current.Layout) != 0 {
		t.Error(state_cloud.GlobalCloudLayout.Current.Layout)
	}

	handleFailedAssign()

	if len(state_cloud.GlobalCloudLayout.Desired.Layout) != 3 {
		t.Error(state_cloud.GlobalCloudLayout.Desired.Layout)
	}

	if !state_cloud.GlobalCloudLayout.Desired.HostHasApp("host1", "app1") {
		t.Error(state_cloud.GlobalCloudLayout.Desired.Layout)
	}
	if !state_cloud.GlobalCloudLayout.Desired.HostHasApp("host2", "app1") {
		t.Error(state_cloud.GlobalCloudLayout.Desired.Layout)
	}
	if state_cloud.GlobalCloudLayout.Desired.HostHasApp("host3", "app1") {
		t.Error(state_cloud.GlobalCloudLayout.Desired.Layout)
	}
}

func TestPlanner_handleFailedAssign_worker(t *testing.T) {
	state_cloud.GlobalCloudLayout.Init()
	state_configuration.GlobalConfigurationState.Init()
	resources := base.InstanceResources{
		TotalCpuResource: 10.0, TotalNetworkResource: 10.0, TotalMemoryResource: 10.0,
	}
	resources3 := base.InstanceResources{
		TotalCpuResource: 5.0, TotalNetworkResource: 8.0, TotalMemoryResource: 10.0,
	}
	state_cloud.GlobalAvailableInstances.Update("host1", resources)
	state_cloud.GlobalAvailableInstances.Update("host2", resources)
	state_cloud.GlobalAvailableInstances.Update("host3", resources3)


	state_needs.GlobalAppsNeedState.UpdateNeeds("app1", 1, base.AppNeeds{CpuNeeds: 1.0, MemoryNeeds: 2.0, NetworkNeeds: 3.0})

	state_configuration.GlobalConfigurationState.ConfigureApp(base.AppConfiguration{
		Name: "app1",
		Type: base.APP_WORKER,
		Version: 1,
		TargetDeploymentCount: 100,
	})

	wipeDesired()
	FailedAssigned = []FailedAssign{
		FailedAssign{AppName: "app1", AppVersion: 1, DeploymentCount: 2, TargetHost: "somehost",},
	}

	if len(state_cloud.GlobalCloudLayout.Current.Layout) != 0 {
		t.Error(state_cloud.GlobalCloudLayout.Current.Layout)
	}

	handleFailedAssign()

	if len(state_cloud.GlobalCloudLayout.Desired.Layout) != 3 {
		t.Error(state_cloud.GlobalCloudLayout.Desired.Layout)
	}

	if state_cloud.GlobalCloudLayout.Desired.HostHasApp("host1", "app1") {
		t.Error(state_cloud.GlobalCloudLayout.Desired.Layout)
	}
	if state_cloud.GlobalCloudLayout.Desired.HostHasApp("host2", "app1") {
		t.Error(state_cloud.GlobalCloudLayout.Desired.Layout)
	}
	if !state_cloud.GlobalCloudLayout.Desired.HostHasApp("host3", "app1") {
		t.Error(state_cloud.GlobalCloudLayout.Desired.Layout)
	}
	elem, _ := state_cloud.GlobalCloudLayout.Desired.GetHost("host3")
	if elem.Apps["app1"].DeploymentCount != 2 {
		t.Error(elem.Apps["app1"])
	}
}



func TestPlanner_MultiplePlanningSteps_noChanges(t *testing.T) {
	state_configuration.GlobalConfigurationState.Init()
	state_cloud.GlobalCloudLayout.Init()
	state_needs.GlobalAppsNeedState = make(map[base.AppName]state_needs.AppNeedVersion)

	state_cloud.GlobalAvailableInstances.Update("host1", base.InstanceResources{TotalCpuResource: 100, TotalMemoryResource: 100, TotalNetworkResource: 100,})
	state_cloud.GlobalAvailableInstances.Update("host2", base.InstanceResources{TotalCpuResource: 200, TotalMemoryResource: 200, TotalNetworkResource: 200,})

	config := example.ExampleJsonConfig()
	config.ApplyToState()
	cloud.Init()

	if len(state_configuration.GlobalConfigurationState.Apps) != 4 {
		t.Error("init state_config apps wrong len")
	}

	if len(state_cloud.GlobalCloudLayout.Current.Layout) != 0 {
		t.Error("init state_cloud current should be empty")
	}
	if len(state_cloud.GlobalCloudLayout.Desired.Layout) != 0 {
		t.Error("init state_cloud desired should be empty")
	}

	if len(state_needs.GlobalAppsNeedState) != 4 {
		t.Error(state_needs.GlobalAppsNeedState)
	}

	InitialPlan()


	if len(state_cloud.GlobalCloudLayout.Current.Layout) != 0 {
		t.Error("init state_cloud current should be empty")
	}
	if len(state_cloud.GlobalCloudLayout.Desired.Layout) == 0 {
		t.Error("init state_cloud desired should have elements")
	}
	for i := 0; i < 10; i++ {
		Plan()

		if len(state_cloud.GlobalCloudLayout.Current.Layout) != 0 {
			t.Error("init state_cloud current should be empty")
		}
		if len(state_cloud.GlobalCloudLayout.Desired.Layout) != 2 || state_cloud.GlobalCloudLayout.Desired.Layout["host1"].Apps["http1"].DeploymentCount != 1 || state_cloud.GlobalCloudLayout.Desired.Layout["host1"].Apps["app1"].DeploymentCount != 2 {
			t.Errorf("iteration %d - %+v", i, state_cloud.GlobalCloudLayout.Desired.Layout)
		}
	}
}



func TestPlanner_MultiplePlanningSteps_hostKilledByEventAndNewHostSpawned(t *testing.T) {
	state_configuration.GlobalConfigurationState.Init()
	state_cloud.GlobalCloudLayout.Init()
	state_needs.GlobalAppsNeedState = make(map[base.AppName]state_needs.AppNeedVersion)

	state_cloud.GlobalAvailableInstances.Update("host1", base.InstanceResources{TotalCpuResource: 100, TotalMemoryResource: 100, TotalNetworkResource: 100,})
	state_cloud.GlobalAvailableInstances.Update("host2", base.InstanceResources{TotalCpuResource: 200, TotalMemoryResource: 200, TotalNetworkResource: 200,})

	config := example.ExampleJsonConfig()
	config.ApplyToState()
	cloud.Init()

	if len(state_configuration.GlobalConfigurationState.Apps) != 4 {
		t.Error("init state_config apps wrong len")
	}

	if len(state_cloud.GlobalCloudLayout.Current.Layout) != 0 {
		t.Error("init state_cloud current should be empty")
	}
	if len(state_cloud.GlobalCloudLayout.Desired.Layout) != 0 {
		t.Error("init state_cloud desired should be empty")
	}

	if len(state_needs.GlobalAppsNeedState) != 4 {
		t.Error(state_needs.GlobalAppsNeedState)
	}

	InitialPlan()


	if len(state_cloud.GlobalCloudLayout.Current.Layout) != 0 {
		t.Error("init state_cloud current should be empty")
	}
	if len(state_cloud.GlobalCloudLayout.Desired.Layout) == 0 {
		t.Error("init state_cloud desired should have elements")
	}

	Plan()

	if len(state_cloud.GlobalCloudLayout.Current.Layout) != 0 {
		t.Error("init state_cloud current should be empty")
	}
	if len(state_cloud.GlobalCloudLayout.Desired.Layout) != 2 || state_cloud.GlobalCloudLayout.Desired.Layout["host1"].Apps["http1"].DeploymentCount != 1 || state_cloud.GlobalCloudLayout.Desired.Layout["host1"].Apps["app1"].DeploymentCount != 2 {
		t.Errorf("%+v", state_cloud.GlobalCloudLayout.Desired.Layout)
	}

	tracker.GlobalHostTracker.HandleCloudProviderEvent(cloud.ProviderEvent{"host1", cloud.PROVIDER_EVENT_KILLED})

	if len(state_cloud.GlobalAvailableInstances) != 1 {
		t.Error(state_cloud.GlobalAvailableInstances)
	}

	Plan()

	if len(state_cloud.GlobalCloudLayout.Current.Layout) != 0 {
		t.Error("init state_cloud current should be empty")
	}
	if len(state_cloud.GlobalCloudLayout.Desired.Layout) != 1 || state_cloud.GlobalCloudLayout.Desired.Layout["host2"].Apps["http1"].DeploymentCount != 1 || state_cloud.GlobalCloudLayout.Desired.Layout["host2"].Apps["app1"].DeploymentCount != 2 || state_cloud.GlobalCloudLayout.Desired.Layout["host2"].Apps["app2"].DeploymentCount != 2 {
		t.Errorf("%+v", state_cloud.GlobalCloudLayout.Desired.Layout)
	}
	//host not ready yet
	Plan()

	if len(state_cloud.GlobalCloudLayout.Current.Layout) != 0 {
		t.Error("init state_cloud current should be empty")
	}
	if len(state_cloud.GlobalCloudLayout.Desired.Layout) != 1 || state_cloud.GlobalCloudLayout.Desired.Layout["host2"].Apps["http1"].DeploymentCount != 1 || state_cloud.GlobalCloudLayout.Desired.Layout["host2"].Apps["app1"].DeploymentCount != 2 || state_cloud.GlobalCloudLayout.Desired.Layout["host2"].Apps["app2"].DeploymentCount != 2 {
		t.Errorf("%+v", state_cloud.GlobalCloudLayout.Desired.Layout)
	}

	//new host is pushing to the trainer = available for planning
	hostInfo := base.HostInfo{
		HostId: "newReplacementHost",
		IpAddr: "1.2.5.6",
		OsInfo: base.OsInfo{},
		Apps: []base.AppInfo{},
	}
	state_cloud.GlobalCloudLayout.Current.UpdateHost(hostInfo)
	state_cloud.GlobalCloudLayout.Current.Layout["host2"] = state_cloud.GlobalCloudLayout.Desired.Layout["host2"]
	tracker.GlobalHostTracker.Update(hostInfo.HostId, time.Now().UTC())

	if len(state_cloud.GlobalCloudLayout.Current.Layout) != 2 {
		t.Error(state_cloud.GlobalCloudLayout.Current.Layout)
	}
	if len(state_cloud.GlobalAvailableInstances) != 2 || state_cloud.GlobalAvailableInstances["newReplacementHost"].TotalCpuResource != 1000 {
		t.Error(state_cloud.GlobalAvailableInstances)
	}

	Plan()

	if len(state_cloud.GlobalCloudLayout.Desired.Layout) != 2 || state_cloud.GlobalCloudLayout.Desired.Layout["newReplacementHost"].Apps["http1"].DeploymentCount != 1 || state_cloud.GlobalCloudLayout.Desired.Layout["host2"].Apps["http1"].DeploymentCount != 1 || state_cloud.GlobalCloudLayout.Desired.Layout["host2"].Apps["app1"].DeploymentCount != 2 || state_cloud.GlobalCloudLayout.Desired.Layout["host2"].Apps["app2"].DeploymentCount != 2 {
		t.Errorf("%+v", state_cloud.GlobalAvailableInstances)
		t.Errorf("%+v", state_cloud.GlobalCloudLayout.Desired.Layout)
	}
}