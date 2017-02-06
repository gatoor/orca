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

package needs

//import (
//"time"
//)
//
//type Needs int
//
//type MemoryNeeds Needs
//type CpuNeeds Needs
//type NetworkNeeds Needs
//
//type AppNeeds struct {
//	MemoryNeeds MemoryNeeds
//	CpuNeeds CpuNeeds
//	NetworkNeeds NetworkNeeds
//}
//
//const (
//	NEEDS_DELTA = 15
//	NEEDS_CAUTION_INTERVAL = 2
//)
//
//
//type Minutes int
//
//type TimeBasedNeeds map[Minutes]AppNeeds
//
//type WeeklyNeeds map[time.Weekday]TimeBasedNeeds
//
//
//func (t TimeBasedNeeds) Get(minutes Minutes) AppNeeds {
//	return t[minutes]
//}
//
//func (w WeeklyNeeds) Get(day time.Weekday, minutes Minutes) AppNeeds {
//	var maxNeeds AppNeeds
//	for i := (minutes - NEEDS_CAUTION_INTERVAL * NEEDS_DELTA); i <= minutes + NEEDS_CAUTION_INTERVAL * NEEDS_DELTA; i += NEEDS_DELTA {
//		if i >= 0 && i < 1440 {
//			current := w[day][Minutes(i)]
//			if current.CpuNeeds > maxbase.CpuNeeds {
//				maxbase.CpuNeeds = current.CpuNeeds
//			}
//			if current.MemoryNeeds > maxbase.NetworkNeeds {
//				maxbase.NetworkNeeds = current.MemoryNeeds
//			}
//			if current.NetworkNeeds > maxNeeds.NetworkNeeds {
//				maxNeeds.NetworkNeeds = current.NetworkNeeds
//			}
//		} else if i < 0 {
//			var dayBefore time.Weekday
//			if day == time.Sunday {
//				dayBefore = time.Saturday
//			} else {
//				dayBefore = day - 1
//			}
//			currentBefore := w[dayBefore][Minutes(1440 + i)]
//			if currentBefore.CpuNeeds > maxbase.CpuNeeds {
//				maxbase.CpuNeeds = currentBefore.CpuNeeds
//			}
//			if currentBefore.MemoryNeeds > maxbase.NetworkNeeds {
//				maxbase.NetworkNeeds = currentBefore.MemoryNeeds
//			}
//			if currentBefore.NetworkNeeds > maxNeeds.NetworkNeeds {
//				maxNeeds.NetworkNeeds = currentBefore.NetworkNeeds
//			}
//		} else if i >= 1440 {
//			var dayAfter time.Weekday
//			if day == time.Saturday {
//				dayAfter = time.Sunday
//			} else {
//				dayAfter = day + 1
//			}
//			currentAfter := w[dayAfter][Minutes(i - 1440)]
//			if currentAfter.CpuNeeds > maxbase.CpuNeeds {
//				maxbase.CpuNeeds = currentAfter.CpuNeeds
//			}
//			if currentAfter.MemoryNeeds > maxbase.NetworkNeeds {
//				maxbase.NetworkNeeds = currentAfter.MemoryNeeds
//			}
//			if currentAfter.NetworkNeeds > maxNeeds.NetworkNeeds {
//				maxNeeds.NetworkNeeds = currentAfter.NetworkNeeds
//			}
//		}
//	}
//	return maxNeeds
//}
//
//// initialized by JSONconfig
////added by metrics.analyzer
//func (w WeeklyNeeds) Set(day time.Weekday, minutes Minutes, ns AppNeeds) {
//	if _, exists := w[day]; !exists {
//		w[day] = make(map[Minutes]AppNeeds)
//	}
//	w[day][minutes] = ns
//}
//
//// helper method for Apps that have the same needs 24/7
//func (w WeeklyNeeds) SetFlat(ns AppNeeds) {
//	 for i := 0; i < 7; i++ {
//		 w[time.Weekday(i)] = make(map[Minutes]AppNeeds)
//		 for m := 0; m < (24 * 60); m += NEEDS_DELTA {
//			 w[time.Weekday(i)][Minutes(m)] = ns
//		 }
//	 }
//}
//
////get weekday and minutes in NEEDS_DELTA increments. always rounded down
//func timeToWeekdayMinutes(t time.Time) (time.Weekday, Minutes) {
//	utcTime := t.UTC()
//	w := utcTime.Weekday()
//	m := utcTime.Hour() * 60 + utcTime.Minute()
//	if m % NEEDS_DELTA != 0 {
//		m = int(m / NEEDS_DELTA) * NEEDS_DELTA
//	}
//	return w, Minutes(m)
//}
//
//func CurrentTimeForNeeds() (time.Weekday, Minutes) {
//	return timeToWeekdayMinutes(time.Now())
//}
