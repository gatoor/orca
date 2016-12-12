package api

import (
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	"gatoor/orca/rewriteTrainer/state/configuration"
	"encoding/json"
	"gatoor/orca/rewriteTrainer/state/needs"
	"gatoor/orca/rewriteTrainer/state/cloud"
	Logger "gatoor/orca/rewriteTrainer/log"
	"gatoor/orca/rewriteTrainer/metrics"
	"gatoor/orca/rewriteTrainer/responder"
	"gatoor/orca/rewriteTrainer/db"
	"gatoor/orca/rewriteTrainer/tracker"
	"gatoor/orca/base"
	"gatoor/orca/rewriteTrainer/config"
)

const ORCA_VERSION = "0.1"

type Api struct{}

var ApiLogger = Logger.LoggerWithField(Logger.Logger, "module", "api")

func (api Api) Init() {
	ApiLogger.Infof("Initializing Api on Port %d", state_configuration.GlobalConfigurationState.Trainer.Port)

	r := mux.NewRouter()

	r.HandleFunc("/push", pushHandler)

	r.HandleFunc("/state/config", getStateConfiguration)
	r.HandleFunc("/state/config/applications", getStateConfigurationApplications)
	r.HandleFunc("/state/cloud", getStateCloud)
	r.HandleFunc("/state/needs", getStateNeeds)

	http.Handle("/", r)

	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%d", state_configuration.GlobalConfigurationState.Trainer.Port), nil)
		if err != nil {
			ApiLogger.Fatalf("Api failed to start - %s", err)
		}
	}()
}

func returnJson(w http.ResponseWriter, obj interface{}) {
	j, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		ApiLogger.Errorf("Json serialization failed - %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func pushHandler(w http.ResponseWriter, r *http.Request) {
	hostInfo, stats, err := metrics.ParsePush(r)
	ApiLogger.Infof("Got metrics from host '%s'", hostInfo.HostId)
	if err != nil {
		json.NewEncoder(w).Encode(nil)
		return
	}

	doHandlePush(hostInfo, stats)

	config, err := responder.GetConfigForHost(hostInfo.HostId)

	if err != nil {
		ApiLogger.Infof("Sending empty response to host '%s'", hostInfo.HostId)
		config = base.PushConfiguration{}
		config.OrcaVersion = ORCA_VERSION
		returnJson(w, config)
		return
	}
	config.OrcaVersion = ORCA_VERSION
	ApiLogger.Infof("Sending new config to host '%s': '%+v'", hostInfo.HostId, config)
	returnJson(w, config)
}

func doHandlePush(hostInfo base.HostInfo, stats base.MetricsWrapper) {
	timeString, time := db.GetNow()

	metrics.RecordStats(hostInfo.HostId, stats, timeString)
	metrics.RecordHostInfo(hostInfo, timeString)

	state_cloud.UpdateCurrent(hostInfo, timeString)
	tracker.GlobalHostTracker.Update(hostInfo.HostId, time)

	responder.CheckAppState(hostInfo)
}

func getStateConfiguration(w http.ResponseWriter, r *http.Request) {
	ApiLogger.Infof("Query to getStateConfiguration")
	returnJson(w, state_configuration.GlobalConfigurationState.Snapshot())
}

func getStateConfigurationApplications(w http.ResponseWriter, r *http.Request) {
	ApiLogger.Infof("Query to getStateConfigurationApplications")
	if (r.Method == "POST") {
		/* We need to create a new configuration object */
		var object config.AppJsonConfiguration
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&object); err != nil {
			ApiLogger.Infof("An error occurred while reading the application information")
		}

		ApiLogger.Infof("Read new configuration for application %s", object.Name)
		ApiLogger.Infof("version %s", object.Version)

		var new_version = object.Version + 1
		state_configuration.GlobalConfigurationState.ConfigureApp(base.AppConfiguration{
			Name: object.Name,
			Type: object.Type,
			Version: new_version,
			TargetDeploymentCount: object.TargetDeploymentCount,
			MinDeploymentCount: object.MinDeploymentCount,
			//InstallCommands []base.OsCommand
			//QueryStateCommand base.OsCommand
			//RemoveCommand base.OsCommand
			//RunCommand base.OsCommand
			//StopCommand base.OsCommand
			DockerConfig: object.DockerConfig,
			LoadBalancer: object.LoadBalancer,
			Network: object.Network,
		})
	}

	returnJson(w, state_configuration.GlobalConfigurationState.AllAppsLatest())
}

func getStateCloud(w http.ResponseWriter, r *http.Request) {
	ApiLogger.Infof("Query to getStateCloud")
	returnJson(w, state_cloud.GlobalCloudLayout.Snapshot())
}

func getStateNeeds(w http.ResponseWriter, r *http.Request) {
	ApiLogger.Infof("Query to getStateNeeds")
	returnJson(w, state_needs.GlobalAppsNeedState.Snapshot())
}

func doHandleCloudEvent() {

}