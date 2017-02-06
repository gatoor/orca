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
	"gatoor/orca/rewriteTrainer/config"
)

func AwsJsonConfig() config.JsonConfiguration {
	conf := config.JsonConfiguration{}

	conf.Trainer.Port = 5000

	//conf.Habitats = []config.HabitatJsonConfiguration{
	//	{
	//		Name: "dockerHabitat",
	//		Version: "0.1",
	//		InstallCommands: []base.OsCommand{
	//			{base.EXEC_COMMAND, base.Command{"apt-get", "update"},},
	//			{base.EXEC_COMMAND, base.Command{"apt-get", "-y install apt-transport-https ca-certificates"},},
	//			{base.EXEC_COMMAND, base.Command{"apt-key", "adv --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys 58118E89F3A912897C070ADBF76221572C52609D"},},
	//			{base.FILE_COMMAND, base.Command{"/etc/apt/sources.list.d/docker.list", "deb https://apt.dockerproject.org/repo ubuntu-xenial main"},},
	//			{base.EXEC_COMMAND, base.Command{"apt-get", "update"},},
	//			{base.EXEC_COMMAND, base.Command{"apt-get", "-y install docker-engine"},},
	//			{base.EXEC_COMMAND, base.Command{"mkdir", "/orca/apps"},},
	//			{base.EXEC_COMMAND, base.Command{"chmod", "777 /orca/apps"},},
	//		},
	//	},
	//}

	//conf.Apps = []config.AppJsonConfiguration{
	//	{
	//		Name: "nginx",
	//		Version: "0.1",
	//		Type: base.APP_HTTP,
	//		MinDeploymentCount: 2,
	//		MaxDeploymentCount: 10,
	//		//InstallCommands: []base.OsCommand{
	//		//	{base.EXEC_COMMAND, base.Command{"mkdir", "/orca/apps/nginx"},},
	//		//	{base.FILE_COMMAND, base.Command{"/orca/apps/nginx/index.html", "HELLO ORCA!"},},
	//		//	{base.EXEC_COMMAND, base.Command{"docker", "run --name orca-nginx -p 80:80 -v /orca/apps/nginx:/usr/share/nginx/html nginx"},},
	//		//},
	//		//QueryStateCommand: base.OsCommand{
	//		//	Type: base.EXEC_COMMAND,
	//		//	Command: base.Command{"wget", "http://localhost"},
	//		//},
	//		//RemoveCommand: base.OsCommand{
	//		//	Type: base.EXEC_COMMAND,
	//		//	Command: base.Command{"docker", "stop orca-nginx"},
	//		//},
	//		Needs: base.AppNeeds{
	//			MemoryNeeds: base.MemoryNeeds(5),
	//			CpuNeeds: base.CpuNeeds(5),
	//			NetworkNeeds: base.NetworkNeeds(5),
	//		},
	//	},
	//	{
	//		Name: "nginxworker",
	//		Version: "0.2",
	//		Type: base.APP_WORKER,
	//		MinDeploymentCount: 5,
	//		MaxDeploymentCount: 10,
	//		//InstallCommands: []base.OsCommand{
	//		//	{base.EXEC_COMMAND, base.Command{"mkdir", "/orca/apps/nginxworker"},},
	//		//	{base.FILE_COMMAND, base.Command{"/orca/apps/nginxworker/index.html", "HELLO ORCA as worker!"},},
	//		//	{base.EXEC_COMMAND, base.Command{"docker", "run --name orca-nginxworker -v /orca/apps/nginxworker:/usr/share/nginx/html nginx"},},
	//		//},
	//		//QueryStateCommand: base.OsCommand{
	//		//	Type: base.EXEC_COMMAND,
	//		//	Command: base.Command{"wget", "http://localhost"},
	//		//},
	//		//RemoveCommand: base.OsCommand{
	//		//	Type: base.EXEC_COMMAND,
	//		//	Command: base.Command{"docker", "stop orca-nginxworker"},
	//		//},
	//		Needs: base.AppNeeds{
	//			MemoryNeeds: base.MemoryNeeds(3),
	//			CpuNeeds: base.CpuNeeds(3),
	//			NetworkNeeds: base.NetworkNeeds(3),
	//		},
	//	},
	//}
	return conf
}