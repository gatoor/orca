/*
Copyright Alex Mack and Michael Lawson (michael@sphinix.com)
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

package cloud


import (
	"gatoor/orca/trainer/model"
	orcaSSh "gatoor/orca/util"
)

type CloudProvider struct {
	Engine CloudEngine
	Changes []*model.ChangeServer

	apiEndpoint string
	sshUser string
}

func (cloud* CloudProvider) Init(engine CloudEngine, sshUser string, apiEndpoint string){
	cloud.Engine = engine
	cloud.apiEndpoint = apiEndpoint
	cloud.sshUser= sshUser
}

func (cloud* CloudProvider) ActionChange(change *model.ChangeServer){
	/* First push this change onto the change queue for the cloud provider */
	cloud.AddChange(change)

	go func() {
		/* Here we can spawn a new server */
		if change.Type == "new_server" {
			newHostId := cloud.Engine.SpawnInstanceSync("t2.micro")
			if newHostId != "" {
				/* If the change times out we need to nuke it */
				change.NewHostId = string(newHostId)

				/* A new server was created, wahoo */
				/* Next we should install some stuff to it */
				ipAddr := cloud.Engine.GetIp(newHostId)
				sshKeyPath := cloud.Engine.GetPem()

				for {
					session, addr := orcaSSh.Connect(cloud.sshUser, string(ipAddr) + ":22", sshKeyPath)
					if session == nil {
						//fail
					}

					SUPERVISOR_CONFIG := "'[unix_http_server]\\nfile=/var/run/supervisor.sock\\nchmod=0770\\nchown=root:supervisor\\n[supervisord]\\nlogfile=/var/log/supervisor/supervisord.log\\npidfile=/var/run/supervisord.pid\\nchildlogdir=/var/log/supervisor\\n[rpcinterface:supervisor]\\nsupervisor.rpcinterface_factory = supervisor.rpcinterface:make_main_rpcinterface\\n[supervisorctl]\\nserverurl=unix:///var/run/supervisor.sock\\n[include]\\nfiles = /etc/supervisor/conf.d/*.conf' > /etc/supervisor/supervisord.conf"
					ORCA_SUPERVISOR_CONFIG := "'[program:orca_client]\\ncommand=/orca/bin/orcahostd --interval 30 --hostid "+string(newHostId)+" --traineruri "+cloud.apiEndpoint+"\\nautostart=true\\nautorestart=true\\nstartretries=2\\nuser=root\\nredirect_stderr=true\\nstdout_logfile=/orca/log/client.log\\nstdout_logfile_maxbytes=50MB\\n' > /etc/supervisor/conf.d/orca.conf"

					instance := []string{
						"echo orca | sudo -S addgroup --system supervisor",
						"echo orca | sudo -S apt-get update",
						"echo orca | sudo -S apt-get install -y git golang supervisor docker.io",
						"echo orca | sudo -S sh -c \"echo " + SUPERVISOR_CONFIG + "\"",
						"echo orca | sudo -S sh -c \"echo " + ORCA_SUPERVISOR_CONFIG + "\"",
						"echo orca | sudo -S rm -rf /orca",
						"echo orca | sudo -S mkdir -p /orca",
						"echo orca | sudo -S mkdir -p /orca/apps",
						"echo orca | sudo -S mkdir -p /orca/log",
						"echo orca | sudo -S mkdir -p /orca/client",
						"echo orca | sudo -S mkdir -p /orca/client/data",
						"echo orca | sudo -S mkdir -p /orca/client/config",
						"echo orca | sudo -S chmod -R 777 /orca",

						"rm -rf /orca/src/bluewhale && mkdir -p /orca/src/bluewhale && cd /orca/src/bluewhale && git clone https://github.com/bluewhale/orcahostd.git",
						"GOPATH=/orca bash -c 'cd /orca/src/bluewhale/orcahostd && go get github.com/Sirupsen/logrus && go get golang.org/x/crypto/ssh && go get github.com/fsouza/go-dockerclient && go get github.com/gorilla/mux'",
						"GOPATH=/orca bash -c 'cd /orca/src/bluewhale/orcahostd && go build && go install'",
						"echo orca | sudo -S service supervisor restart",
					}

					for _, cmd := range instance {
						res := orcaSSh.ExecuteSshCommand(session, addr, cmd)
						if !res {

						}
					}

					break
				}

				cloud.RemoveChange(change.Id)
			}
		}
	}()

}

func (cloud *CloudProvider) HasChanges() bool {
	return len(cloud.Changes) > 0;
}

func (cloud *CloudProvider) GetAllChanges() []*model.ChangeServer {
	return []*model.ChangeServer{}
}

func (cloud* CloudProvider) RemoveChange(changeId string){
	newChanges := make([]*model.ChangeServer, 0)
	for _, change := range cloud.Changes {
		if change.Id != changeId {
			newChanges = append(newChanges, change)
		}
	}
	cloud.Changes = newChanges
}

func (cloud* CloudProvider) AddChange(change *model.ChangeServer){
	cloud.Changes = append(cloud.Changes, change)
}

