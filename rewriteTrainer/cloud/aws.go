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

package cloud

import (
	"gatoor/orca/base"
	Logger "gatoor/orca/rewriteTrainer/log"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"sync"
	"sort"
	"gatoor/orca/rewriteTrainer/installer"
	"gatoor/orca/client/types"
	"fmt"
	"gatoor/orca/rewriteTrainer/state/configuration"
	"os"
	"gatoor/orca/rewriteTrainer/db"
	"strconv"
)

var AWSLogger = Logger.LoggerWithField(Logger.Logger, "module", "aws")

type SpawnLog []base.HostId

var spawnLogMutex = &sync.Mutex{}

func (s SpawnLog) Add(hostId base.HostId) {
	spawnLogMutex.Lock()
	defer spawnLogMutex.Unlock()
	s = append(s, hostId)
}

func (s SpawnLog) Remove(hostId base.HostId) {
	spawnLogMutex.Lock()
	defer spawnLogMutex.Unlock()
	i := -1
	for iter, host := range s {
		if host == hostId {
			i = iter
		}
	}
	if i >= 0 {
		s = append(s[:i], s[i + 1:]...)
	}
}

type AWSProvider struct {
	ProviderConfiguration base.ProviderConfiguration
	Type                  base.ProviderType
}

func (a *AWSProvider) CheckCredentials() bool {
	if a.ProviderConfiguration.AWSConfiguration.Key == "" || a.ProviderConfiguration.AWSConfiguration.Secret == "" {
		AWSLogger.Errorf("No AWS Credentials set")
		return false
	}
	AWSLogger.Infof("Checking AwsCredentials: Key='%s' Secret='%s...'", a.ProviderConfiguration.AWSConfiguration.Key, a.ProviderConfiguration.AWSConfiguration.Secret[:4])

	sess, err := session.NewSession()
	if err != nil {
		AWSLogger.Errorf("AwsCredentials fail: %s", err)
		return false
	}
	svc := ec2.New(sess, &aws.Config{Region: aws.String(a.ProviderConfiguration.AWSConfiguration.Region)})

	_, err = svc.DescribeInstances(nil)
	if err != nil {
		AWSLogger.Errorf("AwsCredentials fail: %s", err)
		return false
	}

	return true
}

func (a *AWSProvider) Init() {
	a.Type = PROVIDER_AWS
	if a.ProviderConfiguration.AWSConfiguration.Key == "" || a.ProviderConfiguration.AWSConfiguration.Secret == "" {
		AWSLogger.Errorf("Missing AWS credential environment variables AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY")
	}

	//TODO: This is amazingly shitty, but because the aws api sucks and I have no patience its the approach for now
	os.Setenv("AWS_ACCESS_KEY_ID", a.ProviderConfiguration.AWSConfiguration.Key)
	os.Setenv("AWS_SECRET_ACCESS_KEY", a.ProviderConfiguration.AWSConfiguration.Secret)

	////TODO: When the cloud provider init is called, we use the aws api based on the credentials set to populate below:
	//t2_micro := base.ProviderInstanceType{
	//	InstanceCost: 0.2,
	//	InstanceResources:base.InstanceResources{
	//		TotalCpuResource:10000,
	//		TotalNetworkResource:10000,
	//		TotalMemoryResource:10000,
	//	},
	//	SpotInstanceCost:0.02,
	//	SupportsSpotInstance:true,
	//	Type:"t2.micro",
	//}
	//
	//a.ProviderConfiguration.AvailableInstanceTypes = make(map[base.InstanceType]base.ProviderInstanceType)
	//a.ProviderConfiguration.AvailableInstanceTypes["t2.micro"] = t2_micro
}

func (a *AWSProvider) GetAvailableInstances(instanceType base.InstanceType) base.ProviderInstanceType {
	return a.ProviderConfiguration.AvailableInstanceTypes[instanceType]
}

func (a *AWSProvider) UpdateAvailableInstances(instanceType base.InstanceType, instance base.ProviderInstanceType) {
	a.ProviderConfiguration.AvailableInstanceTypes[instanceType] = instance
}

func (a *AWSProvider) GetAllAvailableInstanceTypes() map[base.InstanceType]base.ProviderInstanceType {
	return a.ProviderConfiguration.AvailableInstanceTypes
}

func (a *AWSProvider) SpawnInstance(ty base.InstanceType) base.HostId {
	db.Audit.Insert__AuditEvent(db.AuditEvent{
		Details:map[string]string{
			"message": fmt.Sprintf("Trying to spawn a single instance of type '%s' in region %s with AMI %s",
				ty, a.ProviderConfiguration.AWSConfiguration.Region, a.ProviderConfiguration.AWSConfiguration.AMI),
			"subsystem": "cloud.aws",
			"level": "info",
		}})

	svc := ec2.New(session.New(&aws.Config{Region: aws.String(a.ProviderConfiguration.AWSConfiguration.Region)}))

	runResult, err := svc.RunInstances(&ec2.RunInstancesInput{
		ImageId:      aws.String(a.ProviderConfiguration.AWSConfiguration.AMI),
		InstanceType: aws.String(string(ty)),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
		KeyName:      &a.ProviderConfiguration.SSHKey,
		SecurityGroupIds: aws.StringSlice([]string{string(a.ProviderConfiguration.AWSConfiguration.SecurityGroupId)}),
	})

	if err != nil {
		db.Audit.Insert__AuditEvent(db.AuditEvent{Details:map[string]string{
			"message": fmt.Sprintf("Could not spawn instance of type %s: %s", ty, err),
			"subsystem": "cloud.aws",
			"level": "error",
		}})

		return ""
	}

	id := base.HostId(*runResult.Instances[0].InstanceId)
	db.Audit.Insert__AuditEvent(db.AuditEvent{Details:map[string]string{
		"message": fmt.Sprintf("Spawned a single instance of type '%s'. Id=%s", ty, id),
		"subsystem": "cloud.aws",
		"level": "info",
	}})
	return id
}

func (a *AWSProvider) SpawnSpotInstance(ty base.InstanceType, price float64) base.HostId {
	db.Audit.Insert__AuditEvent(db.AuditEvent{Details:map[string]string{
		"message": fmt.Sprintf("Trying to spawn a single spot instance of type '%s' in region %s with AMI %s",
			ty, a.ProviderConfiguration.AWSConfiguration.Region, a.ProviderConfiguration.AWSConfiguration.AMI),
		"subsystem": "cloud.aws",
		"level": "info",
	}})

	svc := ec2.New(session.New(&aws.Config{Region: aws.String(a.ProviderConfiguration.AWSConfiguration.Region)}))

	runResult, err := svc.RunInstances(&ec2.RequestSpotInstancesInput{
		LaunchSpecification: &ec2.RequestSpotLaunchSpecification{
			ImageId:      aws.String(a.ProviderConfiguration.AWSConfiguration.AMI),
			InstanceType: aws.String(string(ty)),
			KeyName:      &a.ProviderConfiguration.SSHKey,
			SecurityGroupIds: aws.StringSlice([]string{string(a.ProviderConfiguration.AWSConfiguration.SecurityGroupId)}),
		},

		Type: aws.String("one-time"),
		InstanceCount: aws.Int64(1),
		SpotPrice: aws.String(strconv.FormatFloat(price, "f", 4, 64)),
	})

	if err != nil {
		db.Audit.Insert__AuditEvent(db.AuditEvent{Details:map[string]string{
			"message": fmt.Sprintf("Could not spawn instance of type %s: %s", ty, err),
			"subsystem": "cloud.aws",
			"level": "error",
		}})

		return ""
	}

	id := base.HostId(*runResult.Instances[0].InstanceId)
	db.Audit.Insert__AuditEvent(db.AuditEvent{Details:map[string]string{
		"message": fmt.Sprintf("Spawned a single instance of type '%s'. Id=%s", ty, id),
		"subsystem": "cloud.aws",
		"level": "info",
	}})
	return id
}

func (a *AWSProvider) waitOnInstanceReady(hostId base.HostId) bool {
	svc := ec2.New(session.New(&aws.Config{Region: aws.String(a.ProviderConfiguration.AWSConfiguration.Region)}))

	if err := svc.WaitUntilInstanceRunning(&ec2.DescribeInstancesInput{InstanceIds: aws.StringSlice([]string{string(hostId)}), }); err != nil {
		AWSLogger.Errorf("WaitOnInstanceReady for %s failed: %s", hostId, err)
	}
	return true
}

func installOrcaClient(hostId base.HostId, ip base.IpAddr, trainerIp base.IpAddr, sshKey string, sshUser string) {
	clientConf := types.Configuration{
		Type: types.DOCKER_CLIENT,
		TrainerPollInterval: 30,
		AppStatusPollInterval: 10,
		MetricsPollInterval: 10,
		TrainerUrl: fmt.Sprintf("http://%s:5000/push", trainerIp),
		Port: 5001,
		HostId: hostId,
	}
	installer.InstallNewInstance(clientConf, ip, sshKey, sshUser)
}

func (a AWSProvider) SpawnInstanceSync(ty base.InstanceType, spot bool) base.HostId {
	AWSLogger.Infof("Spawning Instance synchronously, type %s", ty)
	id := ""
	if spot {
		id = a.SpawnSpotInstance(ty, a.GetAvailableInstances(ty).SpotInstanceCost)
		if id == "" {
			return ""
		}
		if !a.waitOnInstanceReady(id) {
			return ""
		}
	} else {
		id = a.SpawnInstance(ty)
		if id == "" {
			return ""
		}
		if !a.waitOnInstanceReady(id) {
			return ""
		}

	}

	ipAddr := a.GetIp(id)
	sshKeyPath := state_configuration.GlobalConfigurationState.ConfigurationRootPath + "/" + a.ProviderConfiguration.SSHKey + ".pem"
	installOrcaClient(id, ipAddr, state_configuration.GlobalConfigurationState.Trainer.Ip, sshKeyPath, a.ProviderConfiguration.SSHUser)

	return id
}

func (a *AWSProvider) UpdateLoadBalancers(hostId base.HostId, app base.AppName, version base.Version, event string) {
	//var app_configuration, _ = state_configuration.GlobalConfigurationState.GetApp(app, version)
	//if app_configuration.Type == base.APP_HTTP {
	//	if event == base.STATUS_DEAD {
	//		svc := elb.New(session.New(&aws.Config{Region: aws.String(a.ProviderConfiguration.AWSConfiguration.Region)}))
	//
	//		params := &elb.DeregisterInstancesFromLoadBalancerInput{
	//			Instances: []*elb.Instance{{InstanceId: aws.String(string(hostId))}},
	//			LoadBalancerName: aws.String(string(app_configuration.LoadBalancer)),
	//		}
	//		_, err := svc.DeregisterInstancesFromLoadBalancer(params)
	//		if err != nil {
	//			db.Audit.Insert__AuditEvent(db.AuditEvent{Details:map[string]string{
	//				"message": fmt.Sprintf("Could not deregister instance %s from elb %s. Reason was %s", hostId, app_configuration.LoadBalancer, err.Error()),
	//				"subsystem": "cloud.aws",
	//				"level": "error",
	//			}})
	//
	//			return
	//		}
	//
	//		db.Audit.Insert__AuditEvent(db.AuditEvent{Details:map[string]string{
	//			"message": fmt.Sprintf("Deregistered instance %s from elb %s", hostId, app_configuration.LoadBalancer),
	//			"subsystem": "cloud.aws",
	//			"level": "info",
	//		}})
	//
	//	} else if event == base.STATUS_RUNNING {
	//		svc := elb.New(session.New(&aws.Config{Region: aws.String(a.ProviderConfiguration.AWSConfiguration.Region)}))
	//
	//		params := &elb.RegisterInstancesWithLoadBalancerInput{
	//			Instances: []*elb.Instance{
	//				{
	//					InstanceId: aws.String(string(hostId)),
	//				},
	//			},
	//			LoadBalancerName: aws.String(string(app_configuration.LoadBalancer)),
	//		}
	//		_, err := svc.RegisterInstancesWithLoadBalancer(params)
	//		if err != nil {
	//			db.Audit.Insert__AuditEvent(db.AuditEvent{Details:map[string]string{
	//				"message": fmt.Sprintf("Error linking instance %s from elb %s. Reason was %s", hostId, app_configuration.LoadBalancer, err.Error()),
	//				"subsystem": "cloud.aws",
	//				"level": "error",
	//			}})
	//
	//			return
	//		}
	//
	//		db.Audit.Insert__AuditEvent(db.AuditEvent{Details:map[string]string{
	//			"message": fmt.Sprintf("Linked instance %s to elb %s", hostId, app_configuration.LoadBalancer),
	//			"subsystem": "cloud.aws",
	//			"level": "info",
	//		}})
	//	}
	//}
}

func (a *AWSProvider) SpawnInstanceLike(hostId base.HostId) base.HostId {
	AWSLogger.Errorf("NOT IMPLEMENTED")
	AWSLogger.Errorf("NOT IMPLEMENTED")
	AWSLogger.Errorf("NOT IMPLEMENTED")
	return hostId
}

func (a *AWSProvider) SpawnInstances(tys []base.InstanceType) bool {
	AWSLogger.Infof("Trying to spawn %d instances", len(tys))
	AWSLogger.Errorf("NOT IMPLEMENTED")
	AWSLogger.Errorf("NOT IMPLEMENTED")
	AWSLogger.Errorf("NOT IMPLEMENTED")
	return true
}

func (a *AWSProvider) getInstanceInfo(hostId base.HostId) (*ec2.Instance, error) {
	svc := ec2.New(session.New(&aws.Config{Region: aws.String(a.ProviderConfiguration.AWSConfiguration.Region)}))
	res, err := svc.DescribeInstances(&ec2.DescribeInstancesInput{InstanceIds: aws.StringSlice([]string{string(hostId)}), })
	if err != nil {
		return &ec2.Instance{}, err
	}
	if len(res.Reservations) != 1 || len(res.Reservations[0].Instances) != 1 {
		return &ec2.Instance{}, errors.New("Wrong instance count")
	}
	return res.Reservations[0].Instances[0], nil
}

func (a *AWSProvider) GetIp(hostId base.HostId) base.IpAddr {
	AWSLogger.Infof("Getting IpAddress of instance %s", hostId)
	info, err := a.getInstanceInfo(hostId)
	if err != nil {
		AWSLogger.Infof("Got IpAddress for instance %s failed: %s", hostId, err)
		return ""
	}
	ip := base.IpAddr(*info.PublicIpAddress)
	AWSLogger.Infof("Got IpAddress %s for instance %s", ip, hostId)
	return ip
}

func (a *AWSProvider) GetIsSpotInstance(hostId base.HostId) bool {
	AWSLogger.Infof("Getting GetIsSpotInstance of instance %s", hostId)
	info, err := a.getInstanceInfo(hostId)
	if err != nil {
		AWSLogger.Infof("Got GetIsSpotInstance for instance %s failed: %s", hostId, err)
		return false
	}
	return info.SpotInstanceRequestId != nil
}

func (a *AWSProvider) GetInstanceType(hostId base.HostId) base.InstanceType {
	AWSLogger.Infof("Getting InstanceType of instance %s", hostId)
	info, err := a.getInstanceInfo(hostId)
	if err != nil {
		AWSLogger.Infof("Got InstanceType for instance %s failed: %s", hostId, err)
		return ""
	}
	ty := base.InstanceType(*info.InstanceType)
	AWSLogger.Infof("Got InstanceType %s for instance %s", ty, hostId)
	return ty
}

func (a *AWSProvider) GetResources(ty base.InstanceType) base.InstanceResources {
	return a.GetAvailableInstances(ty).InstanceResources
}

func (a *AWSProvider) CheckInstance(hostId base.HostId) InstanceStatus {
	AWSLogger.Infof("Checking instance %s", hostId)
	svc := ec2.New(session.New(&aws.Config{Region: aws.String(a.ProviderConfiguration.AWSConfiguration.Region)}))
	res, err := svc.DescribeInstanceStatus(&ec2.DescribeInstanceStatusInput{InstanceIds: aws.StringSlice([]string{string(hostId)})})
	if err != nil {
		AWSLogger.Infof("Checking instance %s failed:%s", hostId, err)
		return INSTANCE_STATUS_DEAD
	}
	if len(res.InstanceStatuses) != 1 {
		return INSTANCE_STATUS_DEAD
	}
	status := *res.InstanceStatuses[0].InstanceState.Name
	AWSLogger.Info(status)
	return INSTANCE_STATUS_HEALTHY
}

func (a *AWSProvider) TerminateInstance(hostId base.HostId) bool {
	db.Audit.Insert__AuditEvent(db.AuditEvent{Details:map[string]string{
		"message": fmt.Sprintf("Trying to terminate instance %s", hostId),
		"subsystem": "cloud.aws",
		"level": "error",
	}})

	svc := ec2.New(session.New(&aws.Config{Region: aws.String(a.ProviderConfiguration.AWSConfiguration.Region)}))
	_, err := svc.TerminateInstances(&ec2.TerminateInstancesInput{
		InstanceIds: aws.StringSlice([]string{string(hostId)}),
	})

	if err != nil {
		db.Audit.Insert__AuditEvent(db.AuditEvent{Details:map[string]string{
			"message": fmt.Sprintf("Could not terminate instance %s: %s", hostId, err),
			"subsystem": "cloud.aws",
			"level": "error",
		}})

		return false
	}
	db.Audit.Insert__AuditEvent(db.AuditEvent{Details:map[string]string{
		"message": fmt.Sprintf("Terminated instance %s", hostId),
		"subsystem": "cloud.aws",
		"level": "error",
	}})

	return true

}

