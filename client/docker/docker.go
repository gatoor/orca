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

package docker

import (
	Logger "gatoor/orca/rewriteTrainer/log"
	"gatoor/orca/client/types"
	"gatoor/orca/base"
	DockerClient "github.com/fsouza/go-dockerclient"
	"time"
	"bytes"
	"fmt"
)

var DockerLogger = Logger.LoggerWithField(Logger.Logger, "module", "docker")
var dockerCli *DockerClient.Client
var dockerAuth DockerClient.AuthConfiguration

type Client struct {

}


func (c *Client) Init() {
	var err error
	dockerCli, err = DockerClient.NewClient("unix:///var/run/docker.sock")

	if err != nil {
		DockerLogger.Fatalf("Docker client could not be instantiated: %v", err)
	}
}

func DockerCli() *DockerClient.Client {
	if dockerCli == nil {
		DockerLogger.Infof("DockerClient was nil, instantiating again.")
		var err error
		dockerCli, err = DockerClient.NewClient("unix:///var/run/docker.sock")

		if err != nil {
			DockerLogger.Fatalf("Docker client could not be instantiated: %v", err)
		}
	}
	return dockerCli
}

func (c *Client) Type() types.ClientType {
	return types.DOCKER_CLIENT
}

func (c *Client) InstallApp(appConf base.AppConfiguration, appsState *types.AppsState, conf *types.Configuration) bool {
	DockerLogger.Infof("Installing docker app %s:%d with tag %s from %s", appConf.Name, appConf.Version, appConf.DockerConfig.Tag, appConf.DockerConfig.Repository)
	var buf bytes.Buffer
	imageOpt := DockerClient.PullImageOptions{
		Repository: appConf.DockerConfig.Repository,
		Tag: appConf.DockerConfig.Tag,
		OutputStream: &buf,
	}
	err := DockerCli().PullImage(imageOpt, DockerClient.AuthConfiguration{})
	if err != nil {
		DockerLogger.Errorf("Install of app %s:%d failed: %s", appConf.Name, appConf.Version, err)
		return false
	}
	DockerLogger.Infof("Install of app %s:%d successful", appConf.Name, appConf.Version)
	return true
}


func (c *Client) RunApp(appId base.AppId, appConf base.AppConfiguration, appsState *types.AppsState, conf *types.Configuration) bool {
	DockerLogger.Infof("Running docker app %s - %s:%d with tag %s", appId, appConf.Name, appConf.Version, appConf.DockerConfig.Tag)

	bindings := make(map[DockerClient.Port][]DockerClient.PortBinding)
	ports := make(map[DockerClient.Port]struct{})
	for _, v := range appConf.PortMappings {
		bindings[DockerClient.Port(v.ContainerPort)] = []DockerClient.PortBinding{DockerClient.PortBinding{HostPort: v.HostPort}}
		ports[DockerClient.Port(v.ContainerPort)] = struct{}{}
	}
	DockerLogger.Warnf("Bindinds are %+v", bindings)

	hostConfig := DockerClient.HostConfig{PortBindings: bindings, PublishAllPorts:true}
	config := DockerClient.Config{AttachStdout: true, AttachStdin: true, Image: fmt.Sprintf("%s:%s", appConf.Name, appConf.DockerConfig.Tag), ExposedPorts:ports}
	opts := DockerClient.CreateContainerOptions{Name: string(appId), Config: &config, HostConfig:&hostConfig}
	container, containerErr := DockerCli().CreateContainer(opts)
	if containerErr != nil {
		DockerLogger.Errorf("Running docker app %s - %s:%d with tag %s container creation failed: %s", appId, appConf.Name, appConf.Version, appConf.DockerConfig.Tag, containerErr)
		return false
	}

	err := DockerCli().StartContainer(container.ID, &hostConfig)
	if err != nil {
		DockerLogger.Warnf("Running docker app %s - %s:%d failed: %s", appId, appConf.Name, appConf.Version, err)
		return false
	}
	DockerLogger.Infof("Running docker app %s - %s:%d successful", appId, appConf.Name, appConf.Version)
	return true
}


func (c *Client) QueryApp(appId base.AppId, appConf base.AppConfiguration, appsState *types.AppsState, conf *types.Configuration) bool {
	DockerLogger.Debugf("Query docker app %s - %s:%d", appId, appConf.Name, appConf.Version)
	resp, err := DockerCli().InspectContainer(string(appId))
	if err != nil {
		DockerLogger.Debugf("Query docker app %s - %s:%d failed: %s", appId, appConf.Name, appConf.Version, err)
		return false
	}
	DockerLogger.Debugf("Query docker app %s - %s:%d successful %+v", appId, appConf.Name, appConf.Version, resp)
	return resp.State.Running
}

func (c *Client) StopApp(appId base.AppId, appConf base.AppConfiguration, appsState *types.AppsState, conf *types.Configuration) bool {
	DockerLogger.Infof("Stopping docker app %s - %s:%d", appId, appConf.Name, appConf.Version)
	err := DockerCli().StopContainer(fmt.Sprintf("%s", appId), 0)
	fail := false
	if err != nil {
		DockerLogger.Infof("Stopping docker app %s - %s:%d failed: %s", appId, appConf.Name, appConf.Version, err)
		fail = true
	}
	opts := DockerClient.RemoveContainerOptions{ID: string(appId)}
	err = DockerCli().RemoveContainer(opts)
	if err != nil {
		DockerLogger.Infof("Stopping docker app %s - %s:%d failed: %s", appId, appConf.Name, appConf.Version, err)
		fail = true
	}
	if fail {
		return false
	}
	DockerLogger.Infof("Stopping docker app %s - %s:%d successful", appId, appConf.Name, appConf.Version)
	return true
}

func (c *Client) DeleteApp(appConf base.AppConfiguration, appsState *types.AppsState, conf *types.Configuration) bool {
	DockerLogger.Infof("Deleting docker app %s:%d with tag %s", appConf.Name, appConf.Version, appConf.DockerConfig.Tag)
	err := DockerCli().RemoveImage(fmt.Sprintf("%s:%s", appConf.Name, appConf.DockerConfig.Tag))
	if err != nil {
		DockerLogger.Infof("Deleting docker app %s:%d with tag %s failed: %s", appConf.Name, appConf.Version, appConf.DockerConfig.Tag, err)
		return false
	}
	DockerLogger.Infof("Deleting docker app %s:%d successful", appConf.Name, appConf.Version)
	return true
}

func (c *Client) AppMetrics(appId base.AppId, appConf base.AppConfiguration, appsState *types.AppsState, conf *types.Configuration, metrics *types.AppsMetricsById) bool {
	DockerLogger.Debugf("Getting AppMetrics for app %s %s:%d", appId, appConf.Name, appConf.Version)
	errC := make(chan error, 1)
	statsC := make(chan *DockerClient.Stats)
	done := make(chan bool)

	go func() {
		errC <- DockerCli().Stats(DockerClient.StatsOptions{ID: string(appId), Stats: statsC, Stream: true, Done: done})
		close(errC)
	}()
	var resultStats []*DockerClient.Stats
	count := 0
	for {
		count++
		stats, ok := <-statsC
		if !ok || count > 2 {
			close(done)
			break
		}
		resultStats = append(resultStats, stats)
	}
	//err := <-errC
	DockerLogger.Infof("%+v", resultStats)

	//if (err != nil){
	//	DockerLogger.Infof("Getting AppMetrics for app %s %s:%d failed: %s. Only %d results", appId, appConf.Name, appConf.Version, err, len(resultStats))
	//	return false
	//}
	appMet := parseDockerStats(resultStats[0], resultStats[1])
	//TODO test performance of app and add this to appMet.ResponsePerformance
	DockerLogger.Infof("PRE PARSE: %+v", resultStats[0])
	DockerLogger.Infof("DOCKER METRICS: %+v", appMet)
	metrics.Add(appId, time.Now().UTC().Format(time.RFC3339Nano), appMet)
	DockerLogger.Debugf("Getting AppMetrics for app %s %s:%d successful", appId, appConf.Name, appConf.Version)
	return true
}

func parseDockerStats(stat0 *DockerClient.Stats, stat1 *DockerClient.Stats) base.AppStats {
	var (
		cpuPercent = uint64(0)
		cpuDelta = float64(stat1.CPUStats.CPUUsage.TotalUsage) - float64(stat0.CPUStats.CPUUsage.TotalUsage)
		systemDelta = float64(stat1.CPUStats.SystemCPUUsage) - float64(stat0.CPUStats.SystemCPUUsage)
	)

	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = uint64((cpuDelta / systemDelta) * float64(len(stat1.CPUStats.CPUUsage.PercpuUsage)) * 100.0)
	}

	appStats := base.AppStats{}
	appStats.CpuUsage = base.Usage(cpuPercent)

	averageMem := (stat1.MemoryStats.Usage + stat0.MemoryStats.Usage) / 2
	appStats.MemoryUsage = base.Usage(averageMem)
	averageNet := (stat1.Network.RxBytes + stat0.Network.RxBytes) / 2
	if (stat1.Network.TxBytes + stat0.Network.TxBytes) / 2  > averageNet {
		averageNet = (stat1.Network.TxBytes + stat0.Network.TxBytes) / 2
	}
	appStats.NetworkUsage = base.Usage(averageNet)

	return appStats
}
