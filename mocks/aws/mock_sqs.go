package mockaws

import (
	"log"
	"sync"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type MockSQS struct {
	queuedReceiveMessages chan *sqs.ReceiveMessageOutput
	messagesReceived      *sync.WaitGroup // whenever a message is received, call Done() on this waitgroup
	messagesDeleted       *sync.WaitGroup // whenever a message is deleted, call Done() on this waitgroup

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
		log.Printf("returning msg from m.queuedReceiveMessages")

		// update the wait group
		m.messagesReceived.Done()

		return rmi, nil
	default:
		log.Printf("returning new ReceiveMessageOutput{}")
		return &sqs.ReceiveMessageOutput{}, nil
	}

}

func (m *MockSQS) AddPermissionRequest(*sqs.AddPermissionInput) (*request.Request, *sqs.AddPermissionOutput) {
	panic("Not implemented")
}

func (m *MockSQS) AddPermission(*sqs.AddPermissionInput) (*sqs.AddPermissionOutput, error) {
	panic("Not implemented")
}

func (m *MockSQS) ChangeMessageVisibilityRequest(*sqs.ChangeMessageVisibilityInput) (*request.Request, *sqs.ChangeMessageVisibilityOutput) {
	panic("Not implemented")
}

func (m *MockSQS) ChangeMessageVisibility(*sqs.ChangeMessageVisibilityInput) (*sqs.ChangeMessageVisibilityOutput, error) {
	panic("Not implemented")
}

func (m *MockSQS) ChangeMessageVisibilityBatchRequest(*sqs.ChangeMessageVisibilityBatchInput) (*request.Request, *sqs.ChangeMessageVisibilityBatchOutput) {
	panic("Not implemented")
}

func (m *MockSQS) ChangeMessageVisibilityBatch(*sqs.ChangeMessageVisibilityBatchInput) (*sqs.ChangeMessageVisibilityBatchOutput, error) {
	panic("Not implemented")
}

func (m *MockSQS) CreateQueueRequest(*sqs.CreateQueueInput) (*request.Request, *sqs.CreateQueueOutput) {
	panic("Not implemented")
}

func (m *MockSQS) CreateQueue(*sqs.CreateQueueInput) (*sqs.CreateQueueOutput, error) {
	panic("Not implemented")
}

func (m *MockSQS) DeleteMessageRequest(*sqs.DeleteMessageInput) (*request.Request, *sqs.DeleteMessageOutput) {
	panic("Not implemented")
}

func (m *MockSQS) DeleteMessage(*sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	// update the wait group
	m.messagesDeleted.Done()
	return nil, nil

}

func (m *MockSQS) DeleteMessageBatchRequest(*sqs.DeleteMessageBatchInput) (*request.Request, *sqs.DeleteMessageBatchOutput) {
	panic("Not implemented")
}

func (m *MockSQS) DeleteMessageBatch(*sqs.DeleteMessageBatchInput) (*sqs.DeleteMessageBatchOutput, error) {
	panic("Not implemented")
}

func (m *MockSQS) DeleteQueueRequest(*sqs.DeleteQueueInput) (*request.Request, *sqs.DeleteQueueOutput) {
	panic("Not implemented")
}

func (m *MockSQS) DeleteQueue(*sqs.DeleteQueueInput) (*sqs.DeleteQueueOutput, error) {
	panic("Not implemented")
}

func (m *MockSQS) GetQueueAttributesRequest(*sqs.GetQueueAttributesInput) (*request.Request, *sqs.GetQueueAttributesOutput) {
	panic("Not implemented")
}

func (m *MockSQS) GetQueueAttributes(*sqs.GetQueueAttributesInput) (*sqs.GetQueueAttributesOutput, error) {
	panic("Not implemented")
}

func (m *MockSQS) GetQueueUrlRequest(*sqs.GetQueueUrlInput) (*request.Request, *sqs.GetQueueUrlOutput) {
	panic("Not implemented")
}

func (m *MockSQS) GetQueueUrl(*sqs.GetQueueUrlInput) (*sqs.GetQueueUrlOutput, error) {
	panic("Not implemented")
}

func (m *MockSQS) ListDeadLetterSourceQueuesRequest(*sqs.ListDeadLetterSourceQueuesInput) (*request.Request, *sqs.ListDeadLetterSourceQueuesOutput) {
	panic("Not implemented")
}

func (m *MockSQS) ListDeadLetterSourceQueues(*sqs.ListDeadLetterSourceQueuesInput) (*sqs.ListDeadLetterSourceQueuesOutput, error) {
	panic("Not implemented")
}

func (m *MockSQS) ListQueuesRequest(*sqs.ListQueuesInput) (*request.Request, *sqs.ListQueuesOutput) {
	panic("Not implemented")
}

func (m *MockSQS) ListQueues(*sqs.ListQueuesInput) (*sqs.ListQueuesOutput, error) {
	panic("Not implemented")
}

func (m *MockSQS) PurgeQueueRequest(*sqs.PurgeQueueInput) (*request.Request, *sqs.PurgeQueueOutput) {
	panic("Not implemented")
}

func (m *MockSQS) PurgeQueue(*sqs.PurgeQueueInput) (*sqs.PurgeQueueOutput, error) {
	panic("Not implemented")
}

func (m *MockSQS) ReceiveMessageRequest(*sqs.ReceiveMessageInput) (*request.Request, *sqs.ReceiveMessageOutput) {
	panic("Not implemented")
}

func (m *MockSQS) RemovePermissionRequest(*sqs.RemovePermissionInput) (*request.Request, *sqs.RemovePermissionOutput) {
	panic("Not implemented")
}

func (m *MockSQS) RemovePermission(*sqs.RemovePermissionInput) (*sqs.RemovePermissionOutput, error) {
	panic("Not implemented")
}

func (m *MockSQS) SendMessageRequest(*sqs.SendMessageInput) (*request.Request, *sqs.SendMessageOutput) {
	panic("Not implemented")
}

func (m *MockSQS) SendMessage(*sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	panic("Not implemented")
}

func (m *MockSQS) SendMessageBatchRequest(*sqs.SendMessageBatchInput) (*request.Request, *sqs.SendMessageBatchOutput) {
	panic("Not implemented")
}

func (m *MockSQS) SendMessageBatch(*sqs.SendMessageBatchInput) (*sqs.SendMessageBatchOutput, error) {
	panic("Not implemented")
}

func (m *MockSQS) SetQueueAttributesRequest(*sqs.SetQueueAttributesInput) (*request.Request, *sqs.SetQueueAttributesOutput) {
	panic("Not implemented")
}

func (m *MockSQS) SetQueueAttributes(*sqs.SetQueueAttributesInput) (*sqs.SetQueueAttributesOutput, error) {
	panic("Not implemented")
}
