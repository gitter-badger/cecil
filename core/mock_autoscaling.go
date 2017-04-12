package core

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface"
)

type MockAutoScaling struct {

	// Everytime a method is invoked on this MockAutoScaling, a new message will be pushed
	// into this channel with the primary argument of the method invocation (eg,
	// it will be a *autoscaling.DescribeAutoScalingInstancesInput if DescribeAutoScalingInstancesOutput is invoked)
	recordedEvents chan AWSInputOutput

	// Queue of describe instances output
	DescribeAutoScalingInstancesOutput       chan *autoscaling.DescribeAutoScalingInstancesOutput
	DescribeAutoScalingInstancesOutputErrors chan error

	// Embed the AutoScalingAPI interface.  No idea what will happen if unimplemented methods are called.
	autoscalingiface.AutoScalingAPI
}

func NewMockAutoScaling() *MockAutoScaling {
	return &MockAutoScaling{
		recordedEvents:                           make(chan AWSInputOutput, 100),
		DescribeAutoScalingInstancesOutput:       make(chan *autoscaling.DescribeAutoScalingInstancesOutput, 100),
		DescribeAutoScalingInstancesOutputErrors: make(chan error, 100),
	}
}

func (m *MockAutoScaling) DescribeAutoScalingInstances(dii *autoscaling.DescribeAutoScalingInstancesInput) (output *autoscaling.DescribeAutoScalingInstancesOutput, err error) {

	Logger.Info("MockAutoScaling DescribeAutoScalingInstances", "DescribeAutoScalingInstancesInput", dii)
	defer func() {
		recordEvent(m.recordedEvents, dii, output)
	}()

	select {
	case output = <-m.DescribeAutoScalingInstancesOutput:
		return output, nil
	case err = <-m.DescribeAutoScalingInstancesOutputErrors:
		return nil, err
	case <-time.After(time.Second):
		return nil, fmt.Errorf("Timed out since there are no DescribeAutoScalingInstances queued")
	}

	return nil, fmt.Errorf("Unexpected DescribeAutoScalingInstances behavior")

}

func (m *MockAutoScaling) DeleteAutoScalingGroup(*autoscaling.DeleteAutoScalingGroupInput) (*autoscaling.DeleteAutoScalingGroupOutput, error) {
	return &autoscaling.DeleteAutoScalingGroupOutput{}, nil
}

func DescribeAutoScalingInstancesOutput(stackID, stackName, instanceId string) *autoscaling.DescribeAutoScalingInstancesOutput {

	newOutput := autoscaling.DescribeAutoScalingInstancesOutput{}
	newOutput.SetAutoScalingInstances([]*autoscaling.InstanceDetails{&autoscaling.InstanceDetails{
	// TODO: complete fields
	}})

	/*	newOutput.StackResources = []*autoscaling.StackResource{
		&autoscaling.StackResource{
			LogicalResourceId:    aws.String("myGlusterLC"),
			PhysicalResourceId:   aws.String(instanceId),
			ResourceStatus:       aws.String("CREATE_IN_PROGRESS"),
			ResourceStatusReason: aws.String("Resource creation Initiated"),
			ResourceType:         aws.String("AWS::EC2::Instance"),
			StackId:              aws.String(stackID),
			StackName:            aws.String(stackName),
			Timestamp:            aws.Time(time.Now().Add(-time.Second * 10)),
		},
	}*/
	return &newOutput
}
