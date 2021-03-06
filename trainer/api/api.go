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

package api

import (
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	"encoding/json"
	"gatoor/orca/trainer/configuration"
	"gatoor/orca/trainer/model"
	"gatoor/orca/trainer/state"
	log "gatoor/orca/util/log"
)

type Api struct {
	configurationStore *configuration.ConfigurationStore
	state              *state.StateStore
}

var ApiLogger = log.LoggerWithField(log.Logger, "module", "api")

func (api *Api) Init(port int, configurationStore *configuration.ConfigurationStore, state *state.StateStore) {
	api.configurationStore = configurationStore
	api.state = state
	ApiLogger.Infof("Initializing Api on Port %d", port)

	r := mux.NewRouter()

	/* Routes for the client */
	r.HandleFunc("/config", api.getAllConfiguration)
	r.HandleFunc("/config/applications", api.getAllConfigurationApplications)
	r.HandleFunc("/config/applications/configuration/latest", api.getAllConfigurationApplications_Configurations_Latest)
	r.HandleFunc("/state", api.getAllRunningState)
	r.HandleFunc("/checkin", api.hostCheckin)

	r.HandleFunc("/audit", api.getAudit)
	r.HandleFunc("/audit/application", api.getAuditApplication)

	http.Handle("/", r)

	func() {
		err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
		if err != nil {
			ApiLogger.Fatalf("Api failed to start - %s", err)
		}
	}()
}

func returnJson(w http.ResponseWriter, obj interface{}) {
	fmt.Printf("%+v", obj)

	j, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		ApiLogger.Errorf("Json serialization failed - %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func (api *Api) getAllConfiguration(w http.ResponseWriter, r *http.Request) {
	returnJson(w, api.configurationStore.GetAllConfiguration())
}

func (api *Api) getAllConfigurationApplications(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		applicationName := r.URL.Query().Get("application")

		var object model.ApplicationConfiguration
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&object); err == nil {
			application, err := api.configurationStore.GetConfiguration(applicationName)
			if err != nil {
				object.Config = make(map[string]model.VersionConfig)
				application = api.configurationStore.Add(applicationName, &object)
			}

			state.Audit.Insert__AuditEvent(state.AuditEvent{Details:map[string]string{
				"message": "Modified application " + applicationName + " in pool",
				"application": applicationName,
			}})

			application.MinDeployment = object.MinDeployment
			application.DesiredDeployment = object.DesiredDeployment
			api.configurationStore.Save()
		}

	}

	listOfApplications := []*model.ApplicationConfiguration{}
	for _, application := range api.configurationStore.GetAllConfiguration() {
		listOfApplications = append(listOfApplications, application)
	}
	returnJson(w, listOfApplications)
}

func (api *Api) getAllConfigurationApplications_Configurations_Latest(w http.ResponseWriter, r *http.Request) {
	applicationName := r.URL.Query().Get("application")
	application, err := api.configurationStore.GetConfiguration(applicationName)
	if err == nil {
		if r.Method == "POST" {
			var object model.VersionConfig
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&object); err == nil {
				newVersion := application.GetNextVersion()
				object.Version = newVersion
				application.Config[newVersion] = object

				state.Audit.Insert__AuditEvent(state.AuditEvent{Details:map[string]string{
					"message": "Modified application " + applicationName + ", created new configuration",
					"application": applicationName,
				}})

				api.configurationStore.Save()
			}
		}

		returnJson(w, application.GetLatestConfiguration())
		return
	}

	returnJson(w, nil)
}

func (api *Api) getAllRunningState(w http.ResponseWriter, r *http.Request) {
	returnJson(w, api.state.GetAllHosts())
}

func (api *Api) hostCheckin(w http.ResponseWriter, r *http.Request) {
	var apps model.HostCheckinDataPackage
	hostId := r.URL.Query().Get("host")

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&apps); err != nil {
		ApiLogger.Infof("An error occurred while reading the application information")
	}

	result, err := api.state.HostCheckin(hostId, apps)
	if err == nil {
		returnJson(w, result.Changes)
		return
	} else {
		returnJson(w, nil)
	}
}

func (api *Api) getAudit(w http.ResponseWriter, r *http.Request) {
	returnJson(w, state.Audit.Query__AuditEvents(""))
}

func (api *Api) getAuditApplication(w http.ResponseWriter, r *http.Request) {
	applicationName := r.URL.Query().Get("application")
	returnJson(w, state.Audit.Query__AuditEvents(applicationName))
}