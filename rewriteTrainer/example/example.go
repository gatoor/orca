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

package example

import (
	"gatoor/orca/rewriteTrainer/state/cloud"
	"gatoor/orca/rewriteTrainer/config"
	"gatoor/orca/base"
)

func ExampleCloudState() {
	state := state_cloud.CloudLayoutAll{}
	state.Init()
	state.Current.AddEmptyHost("host1")
	state.Current.AddEmptyHost("host2")
	state.Current.AddEmptyHost("host3")
	state.Current.AddApp("host1", "app1", 1, 1)
	state.Current.AddApp("host1", "app11", 1, 2)
	state.Current.AddApp("host2", "app2", 2, 10)

	state.Desired.AddEmptyHost("host1")
	state.Desired.AddEmptyHost("host2")
	state.Desired.AddEmptyHost("host3")
	state.Desired.AddApp("host1", "app1", 1, 1)
	state.Desired.AddApp("host1", "app11", 1, 2)
	state.Desired.AddApp("host2", "app2", 2, 5)
	state.Desired.AddApp("host3", "app2", 2, 5)
}

func ExampleJsonConfig() config.JsonConfiguration {
	conf := config.JsonConfiguration{}

	conf.Trainer.Port = 5000
	conf.Trainer.Policies.TRY_TO_REMOVE_HOSTS = true

	conf.Apps = []base.AppConfiguration{
		{
			Name: "http1",
			Version: 1,
			Type: base.APP_HTTP,
			MinDeploymentCount: 2,
			TargetDeploymentCount: 2,
			Needs: base.AppNeeds{
				MemoryNeeds: base.MemoryNeeds(5),
				CpuNeeds: base.CpuNeeds(5),
				NetworkNeeds: base.NetworkNeeds(5),
			},
		},{
			Name: "app1",
			Version: 1,
			Type: base.APP_WORKER,
			MinDeploymentCount: 2,
			TargetDeploymentCount: 2,
			Needs: base.AppNeeds{
				MemoryNeeds: base.MemoryNeeds(5),
				CpuNeeds: base.CpuNeeds(5),
				NetworkNeeds: base.NetworkNeeds(5),
			},
		},
		{
			Name: "app11",
			Version: 1,
			Type: base.APP_WORKER,
			MinDeploymentCount: 2,
			TargetDeploymentCount: 2,
			Needs: base.AppNeeds{
				MemoryNeeds: base.MemoryNeeds(5),
				CpuNeeds: base.CpuNeeds(5),
				NetworkNeeds: base.NetworkNeeds(5),
			},
		},
		{
			Name: "app2",
			Version: 2,
			Type: base.APP_WORKER,
			MinDeploymentCount: 2,
			TargetDeploymentCount: 2,
			Needs: base.AppNeeds{
				MemoryNeeds: base.MemoryNeeds(5),
				CpuNeeds: base.CpuNeeds(5),
				NetworkNeeds: base.NetworkNeeds(5),
			},
		},
	}
	return conf
}

