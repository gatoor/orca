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

	ns.UpdateNeeds("app1", 1, base.AppNeeds{
		CpuNeeds: base.CpuNeeds(3),
		MemoryNeeds: base.MemoryNeeds(10),
		NetworkNeeds: base.NetworkNeeds(1),
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

	ns.UpdateNeeds("app1", 2, base.AppNeeds{})
	val2, _ := ns.Get("app1", 2)
	if val2.CpuNeeds != 3 {
		t.Error(val2)
	}
}


func TestAppsNeedState_GetAll(t *testing.T) {
	ns := prepareNeedsState()

	ns.UpdateNeeds("app1", 1, base.AppNeeds{
		CpuNeeds: base.CpuNeeds(1),
		MemoryNeeds: base.MemoryNeeds(1),
		NetworkNeeds: base.NetworkNeeds(1),
	})
	ns.UpdateNeeds("app1", 2, base.AppNeeds{
		CpuNeeds: base.CpuNeeds(2),
		MemoryNeeds: base.MemoryNeeds(2),
		NetworkNeeds: base.NetworkNeeds(2),
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



