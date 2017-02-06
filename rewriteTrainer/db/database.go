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

package db

import (
	"time"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gatoor/orca/base"
	"fmt"
)

type OrcaDb struct {
	session *mgo.Session
	db      *mgo.Database
}

var Audit OrcaDb

func (a *OrcaDb) Init(hostname string) {
	session, err := mgo.Dial(hostname)
	if err != nil {
		panic(err)
	}

	a.session = session
	a.db = session.DB("orca")
}

func (a *OrcaDb) Close() {
	a.session.Close()
}

func (a *OrcaDb) Insert__AppMetrics(host base.HostId, stats base.MetricsWrapper, _ string) {
	c := a.db.C("app_metrics")
	for appName, obj := range stats.AppMetrics {
		for appVersion, appMetrics := range obj {
			for _, metric := range appMetrics {
				entity := UtilisationStatistic{
					AppName:appName,
					AppVersion:appVersion,
					Cpu:metric.CpuUsage,
					Mbytes:metric.MemoryUsage,
					Network:metric.NetworkUsage,
					Host:host,
					Timestamp:time.Now(),
				}

				err := c.Insert(&entity)
				if err != nil {
					return
				}
			}
		}
	}
}

func (db *OrcaDb) Insert__AuditEvent(event AuditEvent) {
	event.Timestamp = time.Now()
	c := db.db.C("audit")
	err := c.Insert(&event)
	if err != nil {
		return
	}
}
func (db *OrcaDb) Insert__ApplicationUtilisationStatistic(event ApplicationUtilisationStatistic) {
	event.Timestamp = time.Now()
	c := db.db.C("app_utilisation")
	err := c.Insert(&event)
	if err != nil {
		return
	}
}

func (db *OrcaDb) Insert__ApplicationCountStatistic(count ApplicationCountStatistic) {
	c := db.db.C("app_count")
	err := c.Insert(&count)
	if err != nil {
		return
	}
}

func (db *OrcaDb) Query__ApplicationCountStatistic(application base.AppName) []ApplicationCountStatistic {
	c := db.db.C("app_count")
	var results []ApplicationCountStatistic
	err := c.Find(bson.M{"appname": application}).Sort("-Timestamp").All(&results)
	if err != nil {
		panic("error querying db")
	}

	return results
}


func (db *OrcaDb) Query__ApplicationUtilisationStatistic(application base.AppName) []ApplicationUtilisationStatistic {
	c := db.db.C("app_utilisation")
	var results []ApplicationUtilisationStatistic
	err := c.Find(bson.M{"appname": application}).Sort("-Timestamp").All(&results)
	if err != nil {
		panic("error querying db")
	}

	return results
}

func (db *OrcaDb) Query__AuditEvents(application base.AppName) []AuditEvent {
	c := db.db.C("audit")
	var results []AuditEvent
	if application != "" {
		err := c.Find(bson.M{"details.application": application}).Sort("-Timestamp").All(&results)
		if err != nil {
			panic("error querying db")
		}

	} else {
		err := c.Find(nil).Sort("-Timestamp").All(&results)
		if err != nil {
			panic("error querying db")
		}
	}

	return results
}

func (a *OrcaDb) Query__AppMetrics_Performance__ByMinute(application base.AppName) []UtilisationStatistic {
	c := a.db.C("app_metrics")

	var result = make([]UtilisationStatistic, 0)
	operations := []bson.M{
		bson.M{
			"$match" :bson.M{"appname":application},
		},
		bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"utc_year": bson.M{"$year": "$timestamp" },
					"utc_dayOfYMonth": bson.M{"$dayOfMonth": "$timestamp" },
					"utc_month": bson.M{"$month": "$timestamp" },
					"utc_hour": bson.M{"$hour": "$timestamp" },
					"utc_minute": bson.M{"$minute": "$timestamp" },
				},
				"cpu": bson.M{
					"$sum": "$cpu",
				},
				"mbytes": bson.M{
					"$sum": "$mbytes",
				},
				"network": bson.M{
					"$sum": "$network",
				},
			},
		},
		bson.M{"$sort": bson.M{"_id": 1 } },
	}

	pipe := c.Pipe(operations)
	results := []bson.M{}
	err1 := pipe.All(&results)

	if err1 != nil {
		fmt.Printf("ERROR : %s\n", err1.Error())
	}

	for _, item := range results {
		timestamp := time.Date(
			item["_id"].(bson.M)["utc_year"].(int),
			time.Month(item["_id"].(bson.M)["utc_month"].(int)),
			item["_id"].(bson.M)["utc_dayOfYMonth"].(int),
			item["_id"].(bson.M)["utc_hour"].(int),
			item["_id"].(bson.M)["utc_minute"].(int), 0, 0,
			time.UTC,
		)

		entry := UtilisationStatistic{
			Timestamp:timestamp,
			Mbytes: base.Usage(item["mbytes"].(int64)),
			Cpu: base.Usage(item["cpu"].(int64)),
			Network: base.Usage(item["network"].(int64)),
		}
		result = append(result, entry)
	}

	return result
}

func (a *OrcaDb) Query__AppMetrics_Performance__ByMinute_SingleHost(application base.AppName, host base.HostId) []UtilisationStatistic {
	c := a.db.C("app_metrics")

	var result = make([]UtilisationStatistic, 0)
	operations := []bson.M{
		bson.M{
			"$match" :bson.M{"appname":application},
		},
		bson.M{
			"$match" :bson.M{"host":host},
		},
		bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"utc_year": bson.M{"$year": "$timestamp" },
					"utc_dayOfYMonth": bson.M{"$dayOfMonth": "$timestamp" },
					"utc_month": bson.M{"$month": "$timestamp" },
					"utc_hour": bson.M{"$hour": "$timestamp" },
					"utc_minute": bson.M{"$minute": "$timestamp" },
				},
				"cpu": bson.M{
					"$sum": "$cpu",
				},
				"mbytes": bson.M{
					"$sum": "$mbytes",
				},
				"network": bson.M{
					"$sum": "$network",
				},
			},
		},
		bson.M{"$sort": bson.M{"_id": 1 } },
	}

	pipe := c.Pipe(operations)
	results := []bson.M{}
	err1 := pipe.All(&results)

	if err1 != nil {
		fmt.Printf("ERROR : %s\n", err1.Error())
	}

	for _, item := range results {
		timestamp := time.Date(
			item["_id"].(bson.M)["utc_year"].(int),
			time.Month(item["_id"].(bson.M)["utc_month"].(int)),
			item["_id"].(bson.M)["utc_dayOfYMonth"].(int),
			item["_id"].(bson.M)["utc_hour"].(int),
			item["_id"].(bson.M)["utc_minute"].(int), 0, 0,
			time.UTC,
		)

		entry := UtilisationStatistic{
			Timestamp:timestamp,
			Mbytes: base.Usage(item["mbytes"].(int64)),
			Cpu: base.Usage(item["cpu"].(int64)),
			Network: base.Usage(item["network"].(int64)),
		}
		result = append(result, entry)
	}

	return result
}

