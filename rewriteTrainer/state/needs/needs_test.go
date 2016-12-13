package state_needs_test

import (
	"testing"
	"gatoor/orca/rewriteTrainer/state/needs"
	"gatoor/orca/rewriteTrainer/needs"
)


func prepareNeedsState() state_needs.AppsNeedState {
	return state_needs.AppsNeedState{}
}


func TestAppsNeedState_GetNeeds(t *testing.T) {
	ns := prepareNeedsState()

	ns.UpdateNeeds("app1", 1, needs.AppNeeds{
		CpuNeeds: needs.CpuNeeds(3),
		MemoryNeeds: needs.MemoryNeeds(10),
		NetworkNeeds: needs.NetworkNeeds(1),
	})

	_, err0 := ns.Get("unknown", 100)
	if err0 == nil {
		t.Error("found an app that's not there")
	}
	_, err1 := ns.Get("app1", 100)
	if err1 == nil {
		t.Error("found a version that's not there")
	}
	val, err2 := ns.Get("app1", 1)
	if err2 != nil {
		t.Error("did not find app/version")
	}
	if val.MemoryNeeds != 10 {
		t.Error("got wrong needs value")
	}

	ns.UpdateNeeds("app1", 2, needs.AppNeeds{})
	val2, _ := ns.Get("app1", 2)
	if val2.CpuNeeds != 3 {
		t.Error(val2)
	}
}


func TestAppsNeedState_GetAll(t *testing.T) {
	ns := prepareNeedsState()

	ns.UpdateNeeds("app1", 1, needs.AppNeeds{
		CpuNeeds: needs.CpuNeeds(1),
		MemoryNeeds: needs.MemoryNeeds(1),
		NetworkNeeds: needs.NetworkNeeds(1),
	})
	ns.UpdateNeeds("app1", 2, needs.AppNeeds{
		CpuNeeds: needs.CpuNeeds(2),
		MemoryNeeds: needs.MemoryNeeds(2),
		NetworkNeeds: needs.NetworkNeeds(2),
	})

	_, err0 := ns.GetAll("unknown")
	if err0 == nil {
		t.Error("found an app that's not there")
	}

	val, err2 := ns.GetAll("app1")
	if err2 != nil {
		t.Error("did not find app")
	}
	if len(val) != 2 {
		t.Error("didn't get all versions")
	}
	if val[1].MemoryNeeds != 1 {
		t.Error("got wrong needs value")
	}
}



