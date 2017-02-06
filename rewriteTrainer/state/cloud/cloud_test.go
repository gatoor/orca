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

package state_cloud_test

import (
	"gatoor/orca/base"
	"gatoor/orca/rewriteTrainer/state/cloud"
	"testing"
	"gatoor/orca/rewriteTrainer/needs"
	"gatoor/orca/rewriteTrainer/state/needs"
)

func prepareLayoutState() state_cloud.CloudLayoutAll {
	var GlobalCloudLayout state_cloud.CloudLayoutAll
	GlobalCloudLayout.Init()
	return GlobalCloudLayout
}

func TestLayoutAddHost(t *testing.T) {
	state := prepareLayoutState()
	state.Current.AddHost("hostId1", state_cloud.CloudLayoutElement{
		HostId: "hostId1",
		IpAddress: "0.0.0.0",
		HabitatVersion: 1,
		Apps: make(map[base.AppName]state_cloud.AppsVersion),
	})

	if _, err := state.Current.GetHost("unknown"); err == nil {
		t.Error()
	}
	if val, err := state.Current.GetHost("hostId1"); err == nil {
		if val.HostId != "hostId1" {
			t.Error()
		}
	} else {
		t.Error(err)
	}
}

func TestLayoutRemoveHost(t *testing.T) {
	state := prepareLayoutState()
	state.Current.AddHost("hostId1", state_cloud.CloudLayoutElement{
		HostId: "hostId1",
		IpAddress: "0.0.0.0",
		HabitatVersion: 1,
		Apps: make(map[base.AppName]state_cloud.AppsVersion),
	})
	state.Current.AddHost("hostId2", state_cloud.CloudLayoutElement{
		HostId: "hostId2",
		IpAddress: "0.0.0.0",
		HabitatVersion: 1,
		Apps: make(map[base.AppName]state_cloud.AppsVersion),
	})
	if val, err := state.Current.GetHost("hostId1"); err == nil {
		if val.HostId != "hostId1" {
			t.Error()
		}
	} else {
		t.Error(err)
	}
	state.Current.RemoveHost("hostId1")

	if _, err := state.Current.GetHost("hostId1"); err == nil {
		t.Error()
	}

	if val, err := state.Current.GetHost("hostId2"); err == nil {
		if val.HostId != "hostId2" {
			t.Error()
		}
	} else {
		t.Error(err)
	}
}

func TestLayoutAddApp(t *testing.T) {
	state := prepareLayoutState()
	state.Current.AddHost("hostId1", state_cloud.CloudLayoutElement{
		HostId: "hostId1",
		IpAddress: "0.0.0.0",
		HabitatVersion: 1,
		Apps: make(map[base.AppName]state_cloud.AppsVersion),
	})
	state.Current.AddApp("hostId1", "someapp", 2, 1)

	if val, err := state.Current.GetHost("hostId1"); err == nil {
		if val.HostId != "hostId1" {
			t.Error()
		}
		if val.Apps["someapp"].Version != 2 {
			t.Error()
		}
	} else {
		t.Error(err)
	}
}

func TestLayoutRemoveApp(t *testing.T) {
	state := prepareLayoutState()
	state.Current.AddHost("hostId1", state_cloud.CloudLayoutElement{
		HostId: "hostId1",
		IpAddress: "0.0.0.0",
		HabitatVersion: 1,
		Apps: make(map[base.AppName]state_cloud.AppsVersion),
	})
	state.Current.AddApp("hostId1", "someapp", 2, 1)
	state.Current.AddApp("hostId1", "someapp2", 4, 5)

	if val, err := state.Current.GetHost("hostId1"); err == nil {
		if val.HostId != "hostId1" {
			t.Error()
		}
		if val.Apps["someapp"].Version!= 2 {
			t.Error()
		}
	} else {
		t.Error(err)
	}

	state.Current.RemoveApp("hostId1", "someapp")

	if val, err := state.Current.GetHost("hostId1"); err == nil {
		if val.HostId != "hostId1" {
			t.Error()
		}
		if val.Apps["someapp"].Version == 2 {
			t.Error()
		}
	} else {
		t.Error(err)
	}

	if val, err := state.Current.GetHost("hostId1"); err == nil {
		if val.HostId != "hostId1" {
			t.Error()
		}
		if val.Apps["someapp2"].Version != 4 {
			t.Error()
		}
	} else {
		t.Error(err)
	}
}

func TestLayoutIllegalAccess(t *testing.T) {
	state := prepareLayoutState()
	state.Current.RemoveHost("somehost")
	state.Current.RemoveApp("somehost", "someapp")
}

func TestLayout_FindHostsWithApp(t *testing.T) {
	state := prepareLayoutState()
	state.Current.AddHost("hostId1", state_cloud.CloudLayoutElement{
		HostId: "hostId1",
		IpAddress: "0.0.0.0",
		HabitatVersion: 1,
		Apps: make(map[base.AppName]state_cloud.AppsVersion),
	})
	state.Current.AddHost("hostId2", state_cloud.CloudLayoutElement{
		HostId: "hostId2",
		IpAddress: "0.0.0.0",
		HabitatVersion: 1,
		Apps: make(map[base.AppName]state_cloud.AppsVersion),
	})
	state.Current.AddHost("hostId3", state_cloud.CloudLayoutElement{
		HostId: "hostId3",
		IpAddress: "0.0.0.0",
		HabitatVersion: 1,
		Apps: make(map[base.AppName]state_cloud.AppsVersion),
	})
	state.Current.AddApp("hostId1", "app1", 10, 10)
	state.Current.AddApp("hostId1", "app2", 10, 10)
	state.Current.AddApp("hostId2", "app1", 10, 10)
	state.Current.AddApp("hostId2", "app2", 10, 10)
	state.Current.AddApp("hostId3", "app1", 10, 10)

	res := state.Current.FindHostsWithApp("app1")
	if !res["hostId1"] || !res["hostId2"] || !res["hostId3"] {
		t.Error(res)
	}
	res1 := state.Current.FindHostsWithApp("unnkown")
	if res1["hostId2"] {
		t.Error(res)
	}
	res2 := state.Current.FindHostsWithApp("app2")
	if  !res2["hostId1"] || !res2["hostId2"] || res2["hostId3"]{
		t.Error(res2)
	}
}

func TestLAvailableInstances(t *testing.T) {
	instances := state_cloud.AvailableInstances{}

	if len(instances) != 0 {
		t.Error("should be empty")
	}

	instances.Update("host1", base.InstanceResources{
		TotalCpuResource: 10.0,
		TotalMemoryResource: 20.0,
		TotalNetworkResource: 30.0,
	})

	instances.Update("host2", base.InstanceResources{
		TotalCpuResource: 100.0,
		TotalMemoryResource: 200.0,
		TotalNetworkResource: 300.0,
	})

	if len(instances) != 2 {
		t.Error("len should be 2")
	}

	_, err := instances.GetResources("unknown")

	if err == nil {
		t.Error("should err")
	}

	elem, _ := instances.GetResources("host1")
	if elem.TotalMemoryResource != 20.0 {
		t.Error("wrong memory resource")
	}
	elem2, _ := instances.GetResources("host2")
	if elem2.TotalCpuResource != 100.0 {
		t.Error("wrong cpu resource")
	}

	instances.Update("host1", base.InstanceResources{
		TotalCpuResource: 15.0,
		TotalMemoryResource: 25.0,
		TotalNetworkResource: 35.0,

	})

	if len(instances) != 2 {
		t.Error("len should be 2")
	}
	elem3, _ := instances.GetResources("host1")
	if elem3.TotalMemoryResource != 25.0 {
		t.Error("wrong memory resource")
	}

	instances.Remove("host4")
	if len(instances) != 2 {
		t.Error("len should be 2")
	}

	instances.Remove("host2")
	if len(instances) != 1 {
		t.Error("len should be 1")
	}
}

func TestUpdateHost(t *testing.T) {
	state := prepareLayoutState()
	apps :=  make(map[base.AppName]state_cloud.AppsVersion)
	apps["app1"] = state_cloud.AppsVersion{
		1, 5,
	}
	apps["app2"] = state_cloud.AppsVersion{
		2, 20,
	}

	state.Current.AddHost("hostId1", state_cloud.CloudLayoutElement{
		HostId: "hostId1",
		IpAddress: "0.0.0.0",
		HabitatVersion: 1,
		Apps: apps,
	})

	if _, err := state.Current.GetHost("unknown"); err == nil {
		t.Error()
	}
	if val, err := state.Current.GetHost("hostId1"); err == nil {
		if val.HostId != "hostId1" {
			t.Error()
		}
		if len(val.Apps) != 2 {
			t.Error(val.Apps)
		}
	} else {
		t.Error(err)
	}


	appInfo1 := base.AppInfo{
		Type: base.APP_HTTP,
		Name: "app1",
		Version: 1,
		Status: base.STATUS_DEPLOYING,
	}

	appInfo2 := base.AppInfo{
		Type: base.APP_HTTP,
		Name: "app1",
		Version: 1,
		Status: base.STATUS_RUNNING,
	}
	appInfo3 := base.AppInfo{
		Type: base.APP_HTTP,
		Name: "app2",
		Version: 2,
		Status: base.STATUS_RUNNING,
	}
	appInfo4 := base.AppInfo{
		Type: base.APP_HTTP,
		Name: "app2",
		Version: 2,
		Status: base.STATUS_RUNNING,
	}

	hostInfo := base.HostInfo{
		HostId: "hostId1",
		IpAddr: "1.1.1.1",
		OsInfo: base.OsInfo{},
		Apps: []base.AppInfo{appInfo1, appInfo2, appInfo3, appInfo4},
	}

	state.Current.UpdateHost(hostInfo)

	if _, err := state.Current.GetHost("unknown"); err == nil {
		t.Error()
	}
	if val2, err := state.Current.GetHost("hostId1"); err == nil {
		if val2.HostId != "hostId1" {
			t.Error()
		}
		if val2.IpAddress != "1.1.1.1" {
			t.Error(val2.IpAddress)
		}
		if len(val2.Apps) != 2 {
			t.Error(val2.Apps)
		}
		if val2.Apps["app2"].Version != 2 {
			t.Error(val2.Apps["app2"].Version)
		}
		if val2.Apps["app2"].DeploymentCount != 2 {
			t.Error(val2.Apps["app2"].DeploymentCount)
		}
	} else {
		t.Error(err)
	}
}

func Test_HostHasResourcesForApp(t *testing.T) {
	available := state_cloud.AvailableInstances{}
	available.Update("host1", base.InstanceResources{
		TotalCpuResource: 100, TotalNetworkResource: 100, TotalMemoryResource: 100,
		UsedCpuResource: 50, UsedMemoryResource: 50, UsedNetworkResource: 50,
	})
	if available.HostHasResourcesForApp("host1", base.AppNeeds{CpuNeeds: 60, MemoryNeeds: 30, NetworkNeeds: 30}) {
		t.Error()
	}
	if available.HostHasResourcesForApp("host1", base.AppNeeds{CpuNeeds: 30, MemoryNeeds: 70, NetworkNeeds: 30}) {
		t.Error()
	}
	if available.HostHasResourcesForApp("host1", base.AppNeeds{CpuNeeds: 30, MemoryNeeds: 0, NetworkNeeds: 51}) {
		t.Error()
	}
	if !available.HostHasResourcesForApp("host1", base.AppNeeds{CpuNeeds: 30, MemoryNeeds: 0, NetworkNeeds: 30}) {
		t.Error()
	}
}

func TestAvailableInstances_GlobalResourceConsumption(t *testing.T) {
	available := state_cloud.AvailableInstances{}
	available.Update("host1", base.InstanceResources{
		TotalCpuResource: 100, TotalNetworkResource: 100, TotalMemoryResource: 100,
		UsedCpuResource: 50, UsedMemoryResource: 50, UsedNetworkResource: 50,
	})
	available.Update("host12", base.InstanceResources{
		TotalCpuResource: 100, TotalNetworkResource: 100, TotalMemoryResource: 100,
		UsedCpuResource: 20, UsedMemoryResource: 30, UsedNetworkResource: 40,
	})
	total := available.GlobalResourceConsumption()
	if total.TotalCpuResource != 200 || total.UsedCpuResource != 70 || total.UsedMemoryResource!= 80 || total.UsedNetworkResource != 90 {
		t.Error(total)
	}
}

func TestCloudLayout_AllNeeds(t *testing.T) {
	layout := state_cloud.CloudLayout{}
	layout.Wipe()
	apps := make(map[base.AppName]state_cloud.AppsVersion)
	apps["app1"] = state_cloud.AppsVersion{1, 2}
	apps["app2"] = state_cloud.AppsVersion{1, 5}
	layout.AddHost("host1", state_cloud.CloudLayoutElement{HostId: "host1", Apps: apps})
	state_needs.GlobalAppsNeedState.UpdateNeeds("app1", 1, base.AppNeeds{CpuNeeds:10, MemoryNeeds:5, NetworkNeeds:1})
	state_needs.GlobalAppsNeedState.UpdateNeeds("app2", 1, base.AppNeeds{CpuNeeds:1, MemoryNeeds:1, NetworkNeeds:1})
	ns := layout.AllNeeds()
	if ns.CpuNeeds != 25 || ns.MemoryNeeds != 15 || ns.NetworkNeeds != 7 {
		t.Error(ns)
	}
}

func TestCloudLayout_Needs(t *testing.T) {
	layout := state_cloud.CloudLayout{}
	layout.Wipe()
	apps := make(map[base.AppName]state_cloud.AppsVersion)
	apps["app1"] = state_cloud.AppsVersion{1, 2}
	apps["app2"] = state_cloud.AppsVersion{1, 5}
	layout.AddHost("host1", state_cloud.CloudLayoutElement{HostId: "host1", Apps: apps})
	state_needs.GlobalAppsNeedState.UpdateNeeds("app1", 1, base.AppNeeds{CpuNeeds:10, MemoryNeeds:5, NetworkNeeds:1})
	state_needs.GlobalAppsNeedState.UpdateNeeds("app2", 1, base.AppNeeds{CpuNeeds:1, MemoryNeeds:1, NetworkNeeds:1})
	ns := layout.Needs("app1")
	if ns.CpuNeeds != 20 || ns.MemoryNeeds != 10 || ns.NetworkNeeds != 2 {
		t.Error(ns)
	}
}
