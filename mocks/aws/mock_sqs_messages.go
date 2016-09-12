package mockaws

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
)

// TODO: move this somewhere else -- doesn't aws sdk provide this?
type SQSEnvelope struct {
	Type             string `json:"Type"`
	MessageId        string `json:"MessageId"`
	TopicArn         string `json:"TopicArn"`
	Message          string `json:"Message"`
	Timestamp        string `json:"Timestamp"`
	SignatureVersion string `json:"SignatureVersion"`
	Signature        string `json:"Signature"`
	SigningCertURL   string `json:"SigningCertURL"`
	UnsubscribeURL   string `json:"UnsubscribeURL"`
}

type SQSMessageDetail struct {
	InstanceID string `json:"instance-id"`
	State      string `json:"state"`
}

// TODO: move this somewhere else -- doesn't aws sdk provide this?
type SQSMessage struct {
	Version    string           `json:"version"`
	ID         string           `json:"id"`
	DetailType string           `json:"detail-type"`
	Source     string           `json:"source"`
	Account    string           `json:"account"`
	Time       string           `json:"time"`
	Region     string           `json:"region"`
	Resources  []string         `json:"resources"`
	Detail     SQSMessageDetail `json:"detail"`
}

func NewInstanceLaunchMessage(awsAccountID, awsRegion string, result *string) {

	// create an message
	message := SQSMessage{
		Account: awsAccountID,
		Detail: SQSMessageDetail{
			State:      ec2.InstanceStateNamePending,
			InstanceID: "i-mockinstance",
		},
	}
	messageSerialized, err := json.Marshal(message)
	if err != nil {
		panic(fmt.Sprintf("Error marshaling json: %v", err)) // TODO: return error
	}

	// create an envelope and put the message in
	envelope := SQSEnvelope{
		TopicArn: fmt.Sprintf("todo0:todo1:todo2:%v:%v", awsRegion, awsAccountID),
		Message:  string(messageSerialized),
	}

	// serialize to a string
	envelopeSerialized, err := json.Marshal(envelope)
	if err != nil {
		panic(fmt.Sprintf("Error marshaling json: %v", err)) // TODO: return error
	}

	envSerializedString := string(envelopeSerialized)
	*result = envSerializedString
}
