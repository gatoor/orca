package cloud

import (
	"os"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"errors"
	"fmt"
)

type AwsCloudEngine struct {
	awsAccessKeyId     string
	awsAccessKeySecret string
	awsRegion          string
	awsBaseAmi          string
	sshKey          string
	sshKeyPath	string
	securityGroupId          string
}

func (aws *AwsCloudEngine) Init(awsAccessKeyId string, awsAccessKeySecret string, awsRegion string, awsBaseAmi string, sshKey string,sshKeyPath string, securityGroupId string) {
	aws.awsAccessKeySecret = awsAccessKeySecret
	aws.awsAccessKeyId = awsAccessKeyId
	aws.awsRegion = awsRegion
	aws.awsBaseAmi = awsBaseAmi
	aws.sshKey = sshKey
	aws.sshKeyPath = sshKeyPath
	aws.securityGroupId = securityGroupId

	//TODO: This is amazingly shitty, but because the aws api sucks and I have no patience its the approach for now
	os.Setenv("AWS_ACCESS_KEY_ID", aws.awsAccessKeyId)
	os.Setenv("AWS_SECRET_ACCESS_KEY", aws.awsAccessKeySecret)
}

func (a *AwsCloudEngine) getInstanceInfo(hostId HostId) (*ec2.Instance, error) {
	svc := ec2.New(session.New(&aws.Config{Region: aws.String(a.awsRegion)}))
	res, err := svc.DescribeInstances(&ec2.DescribeInstancesInput{InstanceIds: aws.StringSlice([]string{string(hostId)}), })
	if err != nil {
		return &ec2.Instance{}, err
	}
	if len(res.Reservations) != 1 || len(res.Reservations[0].Instances) != 1 {
		return &ec2.Instance{}, errors.New("Wrong instance count")
	}
	return res.Reservations[0].Instances[0], nil
}

func (a *AwsCloudEngine) GetIp(hostId HostId) string {
	info, err := a.getInstanceInfo(hostId)
	if err != nil {
		return ""
	}

	return string(*info.PublicIpAddress)
}


func (a *AwsCloudEngine) waitOnInstanceReady(hostId HostId) bool {
	svc := ec2.New(session.New(&aws.Config{Region: aws.String(a.awsRegion)}))

	if err := svc.WaitUntilInstanceRunning(&ec2.DescribeInstancesInput{InstanceIds: aws.StringSlice([]string{string(hostId)}), }); err != nil {
		fmt.Println("WaitOnInstanceReady for %s failed: %s", hostId, err)
	}
	return true
}


func (engine *AwsCloudEngine) SpawnInstanceSync(instanceType InstanceType) HostId {
	fmt.Println("AwsCloudEngine SpawnInstanceSync called with ", instanceType)
	svc := ec2.New(session.New(&aws.Config{Region: aws.String(engine.awsRegion)}))

	runResult, err := svc.RunInstances(&ec2.RunInstancesInput{
		ImageId:      aws.String(engine.awsBaseAmi),
		InstanceType: aws.String(string(instanceType)),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
		KeyName:      &engine.sshKey,
		SecurityGroupIds: aws.StringSlice([]string{string(engine.securityGroupId)}),
	})

	if err != nil {
		fmt.Println("AwsCloudEngine SpawnInstanceSync encountered an error ", err)
		return ""
	}

	id := HostId(*runResult.Instances[0].InstanceId)
	fmt.Println("AwsCloudEngine SpawnInstanceSync got a new host, lets wait until its ready. HostID is ", id)
	if !engine.waitOnInstanceReady(id) {
		return ""
	}

	fmt.Println("AwsCloudEngine SpawnInstanceSync finished")
	return id
}

func (aws *AwsCloudEngine) GetInstanceType(HostId) InstanceType {
	return ""

}

func (aws *AwsCloudEngine) TerminateInstance(HostId) bool {
	return false
}

func (aws *AwsCloudEngine) GetPem() string {
	return aws.sshKeyPath
}
