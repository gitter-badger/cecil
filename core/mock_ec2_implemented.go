package core

import (
	"github.com/aws/aws-sdk-go/service/ec2"
)

type MockEc2 struct {

	// Everytime a method is invoked on this MockEc2, a new message will be pushed
	// into this channel with the primary argument of the method invocation (eg,
	// it will be a *ec2.DescribeInstancesInput if DescribeInstances is invoked)
	methodInvocationsChan chan<- interface{}
}

func NewMockEc2(mic chan<- interface{}) *MockEc2 {
	return &MockEc2{
		methodInvocationsChan: mic,
	}
}

func (m *MockEc2) DescribeInstances(input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {

	logger.Info("MockEc2 DescribeInstances", "DescribeInstancesInput", input)
	defer func() {
		m.methodInvocationsChan <- input
	}()

	az := "us-east-1a"
	instanceState := ec2.InstanceStateNamePending

	instance := ec2.Instance{
		InstanceId: input.InstanceIds[0],
		Placement: &ec2.Placement{
			AvailabilityZone: &az,
		},
		State: &ec2.InstanceState{
			Name: &instanceState,
		},
	}
	reservation := ec2.Reservation{
		Instances: []*ec2.Instance{
			&instance,
		},
	}
	output := ec2.DescribeInstancesOutput{
		Reservations: []*ec2.Reservation{
			&reservation,
		},
	}

	return &output, nil
}

func (m *MockEc2) TerminateInstances(input *ec2.TerminateInstancesInput) (*ec2.TerminateInstancesOutput, error) {

	logger.Info("MockEc2 TerminateInstances", "TerminateInstances", input)
	defer func() {
		m.methodInvocationsChan <- input
	}()

	output := ec2.TerminateInstancesOutput{}

	return &output, nil

	// panic("Not implemented")
}
