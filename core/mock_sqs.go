package core

import (
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

type MockSQS struct {
	queuedReceiveMessages chan *sqs.ReceiveMessageOutput
	messagesReceived      *sync.WaitGroup // whenever a message is received, call Done() on this waitgroup
	messagesDeleted       *sync.WaitGroup // whenever a message is deleted, call Done() on this waitgroup
	sqsiface.SQSAPI
}

func NewMockSQS(messagesReceived *sync.WaitGroup, messagesDeleted *sync.WaitGroup) *MockSQS {
	return &MockSQS{
		queuedReceiveMessages: make(chan *sqs.ReceiveMessageOutput, 100),
		messagesReceived:      messagesReceived,
		messagesDeleted:       messagesDeleted,
	}
}

// Enqueue a message that will be returned on the next call to ReceiveMessage()
func (m *MockSQS) Enqueue(rmi *sqs.ReceiveMessageOutput) {
	m.queuedReceiveMessages <- rmi
}

func (m *MockSQS) ReceiveMessage(*sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {

	select {
	case rmi := <-m.queuedReceiveMessages:

		logger.Info("MockSQS returning queued message", "sqsmessage", fmt.Sprintf("%+v", rmi))

		// update the wait group
		m.messagesReceived.Done()

		return rmi, nil
	default:
		logger.Info("MockSQS returning empty message, since there's nothing queued")
		return &sqs.ReceiveMessageOutput{}, nil
	}

}

func (m *MockSQS) DeleteMessage(dmi *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	// update the wait group
	logger.Info("MockSQS DeleteMessage", "sqsmessage", fmt.Sprintf("%+v", dmi))
	m.messagesDeleted.Done()
	return nil, nil

}
