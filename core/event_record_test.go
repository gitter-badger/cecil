package core

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/tleyden/cecil/eventrecord"
)

func TestEventRecord(t *testing.T) {

	eventRecord, err := eventrecord.NewMossEventRecord(false, "")
	if err != nil {
		t.Fatalf("Error setting up event recording: %v", err)
	}

	numMessages := 5

	for i := 0; i < numMessages; i++ {
		sqsMessage := &sqs.Message{}
		msgID := fmt.Sprintf("hello-%v", i)
		sqsMessage.MessageId = &msgID
		err = eventRecord.StoreSQSMessage(sqsMessage)
		if err != nil {
			t.Fatalf("Error storing sqs message: %v", err)
		}

	}

	sqsMessages, err := eventRecord.GetStoredSQSMessages()
	if err != nil {
		t.Fatalf("Error getting sqs messages: %v", err)
	}
	if len(sqsMessages) != numMessages {
		t.Fatalf("Expected %v messages, got %v", numMessages, len(sqsMessages))
	}

}
