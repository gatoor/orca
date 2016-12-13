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

package metrics

import (
	"testing"
	"gatoor/orca/rewriteTrainer/db"
	"gatoor/orca/base"
	"time"
)


func TestSampler_RecordStats(t *testing.T) {
	db.Init("_test")
	tim := "sometimestamp"
	stats :=  base.MetricsWrapper{}
	stats.HostMetrics = make(map[string]base.HostStats)
	stats.AppMetrics = make(map[base.AppName]map[base.Version]map[string]base.AppStats)
	stats.HostMetrics[time.Unix(0,0).Format(time.RFC3339Nano)] = base.HostStats{3, 2, 1}
	RecordStats("host1", stats, tim)

	res := db.Audit.Get(db.BUCKET_AUDIT_RECEIVED_STATS, "host1_sometimestamp")

	if res != "{\"HostMetrics\":{\"1970-01-01T12:00:00+12:00\":{\"MemoryUsage\":3,\"CpuUsage\":2,\"NetworkUsage\":1}},\"AppMetrics\":{}}" {
		t.Error(res)
	}
	db.Close()
}


func TestSampler_RecordHostInfo(t *testing.T) {
	db.Init("_test")
	time := "sometimestamp"
	info :=  base.HostInfo{
		HostId: "host1",
		IpAddr: "0.0.0.0",
		OsInfo: base.OsInfo{},
		Apps: []base.AppInfo{},
	}
	RecordHostInfo(info, time)

	res := db.Audit.Get(db.BUCKET_AUDIT_RECEIVED_HOST_INFO, "host1_sometimestamp")

	if res != "{\"HostId\":\"host1\",\"IpAddr\":\"0.0.0.0\",\"OsInfo\":{\"Os\":\"\",\"Arch\":\"\"},\"Apps\":[]}" {
		t.Error(res)
	}
	db.Close()
}
