package core

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
)

type MockCloudFormation struct {

	// Everytime a method is invoked on this MockCloudFormation, a new message will be pushed
	// into this channel with the primary argument of the method invocation (eg,
	// it will be a *cloudformation.DescribeStackResourcesInput if DescribeStackResourcesResponses is invoked)
	recordedEvents chan AWSInputOutput

	// Queue of describe instances output
	DescribeStackResourcesResponses chan *cloudformation.DescribeStackResourcesOutput
	DescribeStackResourcesErrors    chan error

	// Embed the CloudFormationAPI interface.  No idea what will happen if unimplemented methods are called.
	cloudformationiface.CloudFormationAPI
}

func NewMockCloudFormation() *MockCloudFormation {
	return &MockCloudFormation{
		recordedEvents:                  make(chan AWSInputOutput, 100),
		DescribeStackResourcesResponses: make(chan *cloudformation.DescribeStackResourcesOutput, 100),
		DescribeStackResourcesErrors:    make(chan error, 100),
	}
}

func (m *MockCloudFormation) DescribeStackResources(dii *cloudformation.DescribeStackResourcesInput) (output *cloudformation.DescribeStackResourcesOutput, err error) {

	Logger.Info("MockCloudFormation DescribeStackResources", "DescribeStackResourcesInput", dii)
	defer func() {
		recordEvent(m.recordedEvents, dii, output)
	}()

	select {
	case output = <-m.DescribeStackResourcesResponses:
		return output, nil
	case err = <-m.DescribeStackResourcesErrors:
		return nil, err
	case <-time.After(time.Second):
		return nil, fmt.Errorf("Timed out since there are no DescribeStackResources queued")
	}

	return nil, fmt.Errorf("Unexpected DescribeStackResources behavior")

}

func DescribeStackResourcesOutput(stackID, stackName, instanceId string) *cloudformation.DescribeStackResourcesOutput {

	newOutput := cloudformation.DescribeStackResourcesOutput{}
	newOutput.StackResources = []*cloudformation.StackResource{
		&cloudformation.StackResource{
			LogicalResourceId:    aws.String("myGlusterLC"),
			PhysicalResourceId:   aws.String(instanceId),
			ResourceStatus:       aws.String("CREATE_IN_PROGRESS"),
			ResourceStatusReason: aws.String("Resource creation Initiated"),
			ResourceType:         aws.String("AWS::EC2::Instance"),
			StackId:              aws.String(stackID),
			StackName:            aws.String(stackName),
			Timestamp:            aws.Time(time.Now().Add(-time.Second * 10)),
		},
	}
	return &newOutput
}
