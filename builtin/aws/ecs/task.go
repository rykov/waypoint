package ecs

import (
	"context"
	"crypto/rand"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint/builtin/aws/utils"
	"github.com/oklog/ulid"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
)

// TaskLauncher implements the TaskLauncher plugin interface to support
// launching on-demand tasks for the Waypoint server.
type TaskLauncher struct {
	config TaskLauncherConfig
}

// StartTaskFunc implements component.TaskLauncher
func (p *TaskLauncher) StartTaskFunc() interface{} {
	return p.StartTask
}

// StopTaskFunc implements component.TaskLauncher
func (p *TaskLauncher) StopTaskFunc() interface{} {
	return p.StopTask
}

// TaskLauncherConfig is the configuration structure for the task plugin. At
// this time all these are simply copied from what the Waypoint Server
// installation is using, with the only exception being the OdrTaskRoleName.
type TaskLauncherConfig struct {
	// Cluster is the ECS we're operating in
	Cluster string `hcl:"cluster,optional"`

	// Region is the AWS region we're operating in, e.g. us-west-2, us-east-1
	Region string `hcl:"region,optional"`

	// ExecutionRoleName is the name of the AWS IAM role to apply to the task's
	// Execution Role. At this time we reuse the same Role as the Server
	// Execution Role.
	ExecutionRoleName string `hcl:"execution_role_name,optional"`

	// TaskRoleName is the name of the AWS IAM role to apply to the task. This
	// role determines the privileges the ODR builder has, and must have the correct
	// policies in place to work with the provided registries.
	TaskRoleName string `hcl:"task_role_name,optional"`

	// Subnets are the list of subnets for the cluster. These will match the
	// subnets used for the Cluster
	Subnets string `hcl:"subnets,optional"`

	// SecurityGroupId is the security group used for the Waypoint tasks.
	SecurityGroupId string `hcl:"security_group_id,optional"`

	// LogGroup is the CloudWatch log group name to use.
	LogGroup string `hcl:"log_group,optional"`
}

func (p *TaskLauncher) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(
		docs.FromConfig(&TaskLauncherConfig{}),
		docs.FromFunc(p.StartTaskFunc()),
	)
	if err != nil {
		return nil, err
	}

	doc.Description(`
Launch an ECS task for on-demand tasks from the Waypoint server.

This will use the standard AWS environment variables and IAM Role information to
source authentication information for AWS, using the configured task role.
If no task role name is specified, Waypoint will create one with the required
permissions.
`)

	doc.SetField(
		"odr_task_role_name",
		"Task role name to be used for the task role in the On-Demand Runner task",
		docs.Summary(`
This role must have the correct IAM policies to complete its task.
If this IAM role does not already exist, a role will be created with the correct
permissions"
`,
		),
	)

	return doc, nil
}

// TaskLauncher implements Configurable
func (p *TaskLauncher) Config() (interface{}, error) {
	return &p.config, nil
}

// StopTask signals to docker to stop the container created previously. This
// method is currently unimplemented.
func (p *TaskLauncher) StopTask(
	ctx context.Context,
	log hclog.Logger,
	ti *TaskInfo,
) error {
	log.Warn("the StopTask method is not currently implemented for ECS tasks")
	return nil
}

// StartTask runs an ECS Task to perform the requested job.
func (p *TaskLauncher) StartTask(
	ctx context.Context,
	log hclog.Logger,
	tli *component.TaskLaunchInfo,
) (*TaskInfo, error) {
	// Generate an ID for our task name.
	id, err := ulid.New(ulid.Now(), rand.Reader)
	if err != nil {
		return nil, err
	}

	// Arguments here represent the command inputs used when executing the
	// waypoint runner, e.x:
	//   "runner, agent, -vvv, -id, <some id>, -odr"
	cmd := aws.StringSlice(tli.Arguments)

	sess, err := utils.GetSession(&utils.SessionConfig{
		Region: p.config.Region,
	})
	if err != nil {
		return nil, err
	}
	ecsSvc := ecs.New(sess)

	envs := []*ecs.KeyValuePair{}
	for k, v := range tli.EnvironmentVariables {
		envs = append(envs, &ecs.KeyValuePair{
			Name:  aws.String(k),
			Value: aws.String(v),
		})
	}
	// This environment variable informs Kaniko that it is operating in a
	// container. Without this environment variable, Kaniko will fail with an
	// error
	envs = append(envs, &ecs.KeyValuePair{
		Name:  aws.String("container"),
		Value: aws.String("docker"),
	})

	taskName := fmt.Sprintf("waypoint-odr-task-%s", id.String())

	logOptions := buildLoggingOptions(
		nil,
		p.config.Region,
		p.config.LogGroup,
		taskName,
	)

	def := ecs.ContainerDefinition{
		Name:        &taskName,
		Image:       &tli.OciUrl,
		Command:     cmd,
		Environment: envs,
		LogConfiguration: &ecs.LogConfiguration{
			LogDriver: aws.String(ecs.LogDriverAwslogs),
			Options:   logOptions,
		},
	}

	exRoleArn, err := roleArn(p.config.ExecutionRoleName, sess)
	if err != nil {
		return nil, err
	}

	taskRoleArn, err := roleArn(p.config.TaskRoleName, sess)
	if err != nil {
		return nil, err
	}

	registerTaskDefinitionInput := ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions:    []*ecs.ContainerDefinition{&def},
		ExecutionRoleArn:        &exRoleArn,
		TaskRoleArn:             &taskRoleArn,
		Cpu:                     aws.String("1024"),
		Memory:                  aws.String("2048"),
		Family:                  aws.String("waypoint-runner"),
		NetworkMode:             aws.String("awsvpc"),
		RequiresCompatibilities: []*string{aws.String(ecs.LaunchTypeFargate)},
	}

	taskDef, err := utils.RegisterTaskDefinition(&registerTaskDefinitionInput, ecsSvc)
	if err != nil {
		return nil, err
	}

	taskDefArn := *taskDef.TaskDefinitionArn

	subnetStrings := strings.Split(p.config.Subnets, ",")
	subnets := aws.StringSlice(subnetStrings)

	if _, err := ecsSvc.RunTask(&ecs.RunTaskInput{
		LaunchType:           aws.String(ecs.LaunchTypeFargate),
		Cluster:              &p.config.Cluster,
		Count:                aws.Int64(1),
		TaskDefinition:       &taskDefArn,
		EnableECSManagedTags: aws.Bool(true),
		NetworkConfiguration: &ecs.NetworkConfiguration{
			AwsvpcConfiguration: &ecs.AwsVpcConfiguration{
				Subnets:        subnets,
				SecurityGroups: []*string{&p.config.SecurityGroupId},
				// without a public IP we cannot reach out to ECR or other
				// registries
				AssignPublicIp: aws.String("ENABLED"),
			},
		},
	}); err != nil {
		return nil, err
	}

	return &TaskInfo{
		Id: taskName,
	}, nil
}

func roleArn(name string, sess *session.Session) (string, error) {
	iamSvc := iam.New(sess)
	input := &iam.GetRoleInput{
		RoleName: &name,
	}

	roleOut, err := iamSvc.GetRole(input)
	if err != nil {
		return "", err
	}
	if roleOut.Role == nil {
		return "", fmt.Errorf("no role found for (%s) role name", name)
	}
	// Arn is a required field of Role so we assume it will be populated and
	// prevent any nil dereference here
	return *roleOut.Role.Arn, nil
}

var _ component.TaskLauncher = (*TaskLauncher)(nil)
