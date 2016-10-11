package core

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

type MockEc2 struct {

	// Everytime a method is invoked on this MockEc2, a new message will be pushed
	// into this channel with the primary argument of the method invocation (eg,
	// it will be a *ec2.DescribeInstancesInput if DescribeInstances is invoked)
	recordedEvents chan AWSInputOutput

	// Queue of describe instances output
	describeInstanceResponses chan *ec2.DescribeInstancesOutput

	// Embed the EC2API interface.  No idea what will happen if unimplemented methods are called.
	ec2iface.EC2API
}

func NewMockEc2() *MockEc2 {
	return &MockEc2{
		recordedEvents:            make(chan AWSInputOutput, 100),
		describeInstanceResponses: make(chan *ec2.DescribeInstancesOutput, 100),
	}
}

func (m *MockEc2) DescribeInstances(dii *ec2.DescribeInstancesInput) (output *ec2.DescribeInstancesOutput, err error) {

	logger.Info("MockEc2 DescribeInstances", "DescribeInstancesInput", dii)
	defer func() {
		recordEvent(m.recordedEvents, dii, output)
	}()

	logger.Info("DescribeInstances", "input", dii)
	output = <-m.describeInstanceResponses
	logger.Info("DescribeInstances", "output", output)

	return output, nil

}

func (m *MockEc2) TerminateInstances(tii *ec2.TerminateInstancesInput) (output *ec2.TerminateInstancesOutput, err error) {

	logger.Info("MockEc2 TerminateInstances", "TerminateInstances", tii)
	defer func() {
		recordEvent(m.recordedEvents, tii, output)
	}()

	output = &ec2.TerminateInstancesOutput{}

	return output, nil
}

func (m *MockEc2) waitForDescribeInstancesInput() {

	awsInputOutput := <-m.recordedEvents

	dii, ok := awsInputOutput.Input.(*ec2.DescribeInstancesInput)
	if !ok {
		panic(fmt.Sprintf("Expected ec2.DescribeInstancesInput"))
	}
	logger.Info("waitForDescribeInstancesInput", "dii", fmt.Sprintf("%+v", dii))

}

func (m *MockEc2) waitForTerminateInstancesInput() {

	awsInputOutput := <-m.recordedEvents

	tii, ok := awsInputOutput.Input.(*ec2.TerminateInstancesInput)
	if !ok {
		panic(fmt.Sprintf("Expected ec2.TerminateInstancesInput"))
	}
	logger.Info("waitForTerminateInstancesInput", "tii", fmt.Sprintf("%+v", tii))

}

func describeInstanceOutput(instanceState, instanceId string) *ec2.DescribeInstancesOutput {

	az := "us-east-1a"

	instance := ec2.Instance{
		InstanceId: &instanceId,
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
	output := &ec2.DescribeInstancesOutput{
		Reservations: []*ec2.Reservation{
			&reservation,
		},
	}

	return output

}
