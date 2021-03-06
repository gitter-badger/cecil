// Licensed to the Apache Software Foundation (ASF) under one or more contributor license agreements;
// and to You under the Apache License, Version 2.0.  See LICENSE in project root for full license + copyright.

package core

import (
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

type MockSQS struct {

	// Anything added to this channel will be returned to consumers of this MockSQS
	queuedReceiveMessages chan *sqs.ReceiveMessageOutput

	// All actions performed are saved in recordedEvents
	recordedEvents chan AWSInputOutput

	messagesReceived *sync.WaitGroup // whenever a message is received, call Done() on this waitgroup
	messagesDeleted  *sync.WaitGroup // whenever a message is deleted, call Done() on this waitgroup
	sqsiface.SQSAPI
}

func NewMockSQS() *MockSQS {
	return &MockSQS{
		queuedReceiveMessages: make(chan *sqs.ReceiveMessageOutput, 100),
		recordedEvents:        make(chan AWSInputOutput, 100),
	}
}

// Enqueue a message that will be returned on the next call to ReceiveMessage()
func (m *MockSQS) Enqueue(rmi *sqs.ReceiveMessageOutput) {
	m.queuedReceiveMessages <- rmi
}

func (m *MockSQS) ReceiveMessage(rmi *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {

	select {
	case rmo := <-m.queuedReceiveMessages:

		Logger.Info("MockSQS returning queued message", "sqsmessage", fmt.Sprintf("%+v", rmo))

		recordEvent(m.recordedEvents, rmi, rmo)

		return rmo, nil
	default:
		Logger.Info("MockSQS returning empty message, since there's nothing queued")
		return &sqs.ReceiveMessageOutput{}, nil
	}

}

func (m *MockSQS) DeleteMessage(dmi *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {

	recordEvent(m.recordedEvents, dmi, nil)

	Logger.Info("MockSQS DeleteMessage", "sqsmessage", fmt.Sprintf("%+v", dmi))

	return nil, nil

}

func recordEvent(dest chan<- AWSInputOutput, input interface{}, output interface{}) {

	// Record that this event happened
	event := AWSInputOutput{
		Input:  input,
		Output: output,
	}
	dest <- event

}

func (m *MockSQS) WaitForReceivedMessageInput() {
	awsInputOutput := <-m.recordedEvents
	Logger.Info("MockSQS", "recorded receive msg event", fmt.Sprintf("%+v", awsInputOutput))
	rmi, ok := awsInputOutput.Input.(*sqs.ReceiveMessageInput)
	if !ok {
		panic(fmt.Sprintf("Expected *sqs.ReceiveMessageInput, got %+v", awsInputOutput))
	}
	Logger.Info("rmi", fmt.Sprintf("%+v", rmi))

}

func (m *MockSQS) WaitForDeletedMessageInput(receiptHandle string) {

	// Wait until the SQS message is deleted by the eventinjestor
	awsInputOutput := <-m.recordedEvents
	dmi, ok := awsInputOutput.Input.(*sqs.DeleteMessageInput)
	if !ok {
		panic(fmt.Sprintf("Expected *sqs.DeleteMessageInput, got: %+v", awsInputOutput))
	}
	if *dmi.ReceiptHandle != receiptHandle {
		panic(fmt.Sprintf("Expected different receipt handle"))
	}
	Logger.Info("MockSQS", "recorded deleted event", fmt.Sprintf("%+v", awsInputOutput))

}
