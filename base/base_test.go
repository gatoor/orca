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

package base

import (
	"testing"
	"time"
)

func TestBaseB_timeToWeekdayMinutes(t *testing.T) {
	t1, _ := time.Parse(time.RFC3339Nano, "2016-11-26T00:15:44+00:00")
	w1, m1 := timeToWeekdayMinutes(t1)
	if w1 != time.Saturday || m1 != 15 {
		t.Error(w1, m1)
	}
	t2, _ := time.Parse(time.RFC3339Nano, "2016-11-27T00:05:44+00:00")
	w2, m2 := timeToWeekdayMinutes(t2)
	if w2 != time.Sunday || m2 != 0 {
		t.Error(w2, m2)
	}
	t3, _ := time.Parse(time.RFC3339Nano, "2016-11-27T10:16:44+00:00")
	w3, m3 := timeToWeekdayMinutes(t3)
	if w3 != time.Sunday || m3 != 615 {
		t.Error(w3, m3)
	}
	t4, _ := time.Parse(time.RFC3339Nano, "2016-11-27T10:29:44+00:00")
	w4, m4 := timeToWeekdayMinutes(t4)
	if w4 != time.Sunday || m4 != 615 {
		t.Error(w4, m4)
	}
}


func TestNeeds_SetFlat(t *testing.T) {
	weekly := WeekdayBasedAppNeeds{}
	weekly.SetFlat(AppNeeds{CpuNeeds:10, MemoryNeeds:20, NetworkNeeds:30})

	t1, _ := time.Parse(time.RFC3339Nano, "2016-11-26T06:00:44+00:00")

	if weekly.Get(timeToWeekdayMinutes(t1)).CpuNeeds!= 10 || weekly.Get(timeToWeekdayMinutes(t1)).MemoryNeeds != 20 {
		t.Error(weekly.Get(timeToWeekdayMinutes(t1)))
	}
	t2, _ := time.Parse(time.RFC3339Nano, "2016-11-29T16:47:44+00:00")

	if weekly.Get(timeToWeekdayMinutes(t2)).CpuNeeds != 10 || weekly.Get(timeToWeekdayMinutes(t2)).MemoryNeeds != 20 {
		t.Error(weekly.Get(timeToWeekdayMinutes(t2)))
	}

	t3 := time.Now()
	if weekly.Get(timeToWeekdayMinutes(t3)).CpuNeeds != 10 || weekly.Get(timeToWeekdayMinutes(t3)).MemoryNeeds  != 20 {
		t.Error(weekly.Get(timeToWeekdayMinutes(t3)))
	}

	weeklyDepl := WeekdayBasedDeploymentCount{}
	weeklyDepl.SetFlat(DeploymentCount(5))

	t11, _ := time.Parse(time.RFC3339Nano, "2016-11-26T06:00:44+00:00")

	if weeklyDepl.Get(timeToWeekdayMinutes(t11)) != 5  {
		t.Error(weeklyDepl.Get(timeToWeekdayMinutes(t11)))
	}
	t22, _ := time.Parse(time.RFC3339Nano, "2016-11-29T16:47:44+00:00")

	if weeklyDepl.Get(timeToWeekdayMinutes(t22)) != 5 {
		t.Error(weeklyDepl.Get(timeToWeekdayMinutes(t22)))
	}

	t33 := time.Now()
	if weeklyDepl.Get(timeToWeekdayMinutes(t33)) != 5 {
		t.Error(weeklyDepl.Get(timeToWeekdayMinutes(t33)))
	}
}

func TestNeeds_Set(t *testing.T) {
	weekly := WeekdayBasedAppNeeds{}
	weekly.Set(time.Tuesday, 600, AppNeeds{CpuNeeds:10, MemoryNeeds:5, NetworkNeeds: 2})

	if weekly.Get(time.Tuesday, 600).CpuNeeds != 10 {
		t.Error(weekly.Get(time.Tuesday, 600))
	}
	if weekly.Get(time.Tuesday, 900).CpuNeeds != 0 {
		t.Error(weekly.Get(time.Tuesday, 900))
	}

	weeklyDeploymentCount := WeekdayBasedDeploymentCount{}
	weeklyDeploymentCount.Set(time.Tuesday, 600, DeploymentCount(10))

	if weeklyDeploymentCount.Get(time.Tuesday, 600) != 10 {
		t.Error(weeklyDeploymentCount.Get(time.Tuesday, 600))
	}
	if weeklyDeploymentCount.Get(time.Tuesday, 900) != 0 {
		t.Error(weeklyDeploymentCount.Get(time.Tuesday, 900))
	}
}

func TestNeeds_GetDaySwitches(t *testing.T) {
	weekly := WeekdayBasedAppNeeds{}
	weekly.SetFlat(AppNeeds{CpuNeeds:1, MemoryNeeds:1, NetworkNeeds: 1})
	weekly.Set(time.Monday, 1410, AppNeeds{CpuNeeds:20, MemoryNeeds:5, NetworkNeeds: 2})
	weekly.Set(time.Tuesday, 0, AppNeeds{CpuNeeds:10, MemoryNeeds:5, NetworkNeeds: 2})

	if weekly.Get(time.Tuesday, 0).CpuNeeds != 20 {
		t.Error(weekly.Get(time.Tuesday, 0))
	}
	if weekly.Get(time.Tuesday, MINUTES_DELTA).CpuNeeds != 10 {
		t.Error(weekly.Get(time.Tuesday, MINUTES_DELTA))
	}
	if weekly.Get(time.Tuesday, 2 * MINUTES_DELTA).CpuNeeds != 10 {
		t.Error(weekly.Get(time.Tuesday, 2 * MINUTES_DELTA))
	}
	if weekly.Get(time.Tuesday, 20 * MINUTES_DELTA).CpuNeeds != 1 {
		t.Error(weekly.Get(time.Tuesday, 20 * MINUTES_DELTA))
	}
	weekly.Set(time.Monday, 1410, AppNeeds{CpuNeeds:1, MemoryNeeds:1, NetworkNeeds: 1})
	if weekly.Get(time.Monday, 1425).CpuNeeds != 10 {
		t.Error(weekly.Get(time.Tuesday, (24 * 60) - MINUTES_DELTA))
	}

	weekly.SetFlat(AppNeeds{CpuNeeds:1, MemoryNeeds:1, NetworkNeeds: 1})
	weekly.Set(time.Sunday, 1425, AppNeeds{CpuNeeds:20, MemoryNeeds:5, NetworkNeeds: 2})
	if weekly.Get(time.Sunday, 1440).CpuNeeds != 20 {
		t.Error(weekly.Get(time.Sunday, 1440))
	}
	if weekly.Get(time.Sunday, 1410).CpuNeeds != 20 {
		t.Error(weekly.Get(time.Sunday, 1410))
	}
	if weekly.Get(time.Monday, 0).CpuNeeds != 20 {
		t.Error(weekly.Get(time.Monday, 0))
	}
	if weekly.Get(time.Monday, 15).CpuNeeds != 20 {
		t.Error(weekly.Get(time.Monday, 15))
	}
	if weekly.Get(time.Monday, 30).CpuNeeds != 1 {
		t.Error(weekly.Get(time.Monday, 30))
	}


	weekly.SetFlat(AppNeeds{CpuNeeds:1, MemoryNeeds:1, NetworkNeeds: 1})
	weekly.Set(time.Sunday, 0, AppNeeds{CpuNeeds:20, MemoryNeeds:5, NetworkNeeds: 2})
	if weekly.Get(time.Saturday, 1425).CpuNeeds != 20 {
		t.Error(weekly.Get(time.Saturday, 30))
	}
	if weekly.Get(time.Sunday, 30).CpuNeeds != 20 {
		t.Error(weekly.Get(time.Sunday, 30))
	}
	if weekly.Get(time.Sunday, 45).CpuNeeds != 1 {
		t.Error(weekly.Get(time.Sunday, 45))
	}
	if weekly.Get(time.Saturday, 1410).CpuNeeds != 20 {
		t.Error(weekly.Get(time.Saturday, 1410))
	}
	if weekly.Get(time.Saturday, 1395).CpuNeeds != 1 {
		t.Error(weekly.Get(time.Saturday, 1410))
	}

	weekly.SetFlat(AppNeeds{CpuNeeds:1, MemoryNeeds:1, NetworkNeeds: 1})
	weekly.Set(time.Saturday, 0, AppNeeds{CpuNeeds:20, MemoryNeeds:5, NetworkNeeds: 2})
	if weekly.Get(time.Friday, 1425).CpuNeeds != 20 {
		t.Error(weekly.Get(time.Friday, 1425))
	}
	if weekly.Get(time.Saturday, 30).CpuNeeds != 20 {
		t.Error(weekly.Get(time.Saturday, 30))
	}
	if weekly.Get(time.Saturday, 45).CpuNeeds != 1 {
		t.Error(weekly.Get(time.Saturday, 45))
	}
	if weekly.Get(time.Friday, 1410).CpuNeeds != 20 {
		t.Error(weekly.Get(time.Friday, 1410))
	}
	if weekly.Get(time.Friday, 1395).CpuNeeds != 1 {
		t.Error(weekly.Get(time.Friday, 1410))
	}
	weekly.SetFlat(AppNeeds{CpuNeeds:1, MemoryNeeds:1, NetworkNeeds: 1})
	weekly.Set(time.Saturday, 1425, AppNeeds{CpuNeeds:20, MemoryNeeds:5, NetworkNeeds: 2})
	if weekly.Get(time.Saturday, 1425).CpuNeeds != 20 {
		t.Error(weekly.Get(time.Saturday, 1425))
	}
	if weekly.Get(time.Sunday, 15).CpuNeeds != 20 {
		t.Error(weekly.Get(time.Sunday, 30))
	}
	if weekly.Get(time.Sunday, 45).CpuNeeds != 1 {
		t.Error(weekly.Get(time.Sunday, 45))
	}
	if weekly.Get(time.Saturday, 1410).CpuNeeds != 20 {
		t.Error(weekly.Get(time.Saturday, 1410))
	}
	if weekly.Get(time.Saturday, 1380).CpuNeeds != 1 {
		t.Error(weekly.Get(time.Saturday, 1380))
	}
}

func TestWeekdayBasedDeploymentCount(t *testing.T) {
	weekly := WeekdayBasedDeploymentCount{}
	weekly.SetFlat(1)
	weekly.Set(time.Monday, 1410, 5)
	weekly.Set(time.Tuesday, 0, 2)

	if weekly.Get(time.Tuesday, 0) != 5 {
		t.Error(weekly.Get(time.Tuesday, 0))
	}
	if weekly.Get(time.Tuesday, MINUTES_DELTA) != 2 {
		t.Error(weekly.Get(time.Tuesday, MINUTES_DELTA))
	}
	if weekly.Get(time.Tuesday, 2 * MINUTES_DELTA) != 2 {
		t.Error(weekly.Get(time.Tuesday, 2 * MINUTES_DELTA))
	}
	if weekly.Get(time.Tuesday, 20 * MINUTES_DELTA) != 1 {
		t.Error(weekly.Get(time.Tuesday, 20 * MINUTES_DELTA))
	}
	weekly.Set(time.Monday, 1410, 100)
	if weekly.Get(time.Monday, 1425) != 100 {
		t.Error(weekly.Get(time.Tuesday, (24 * 60) - MINUTES_DELTA))
	}

	weekly.SetFlat(1)
	weekly.Set(time.Sunday, 1425, 20)
	if weekly.Get(time.Sunday, 1440) != 20 {
		t.Error(weekly.Get(time.Sunday, 1440))
	}
	if weekly.Get(time.Sunday, 1410) != 20 {
		t.Error(weekly.Get(time.Sunday, 1410))
	}
	if weekly.Get(time.Monday, 0) != 20 {
		t.Error(weekly.Get(time.Monday, 0))
	}
	if weekly.Get(time.Monday, 15) != 20 {
		t.Error(weekly.Get(time.Monday, 15))
	}
	if weekly.Get(time.Monday, 30) != 1 {
		t.Error(weekly.Get(time.Monday, 30))
	}
}