package awstools

import "time"

// TODO: doesn't aws sdk provide this?
// SQSEnvelope is the a container of an SQS message
type SQSEnvelope struct {
	Type             string    `json:"Type"`
	MessageId        string    `json:"MessageId"`
	TopicArn         string    `json:"TopicArn"`
	Message          string    `json:"Message"`
	Timestamp        time.Time `json:"Timestamp"`
	SignatureVersion string    `json:"SignatureVersion"`
	Signature        string    `json:"Signature"`
	SigningCertURL   string    `json:"SigningCertURL"`
	SubscribeURL     string    `json:"SubscribeURL"`
	UnsubscribeURL   string    `json:"UnsubscribeURL"`
}

// SQSMessageDetail is the details of a new instance SQS message
type SQSMessageDetail struct {
	InstanceID string `json:"instance-id"`
	State      string `json:"state"`
}

// TODO: doesn't aws sdk provide this?
// SQSMessage is a message
type SQSMessage struct {
	Version    string           `json:"version"`
	ID         string           `json:"id"`
	DetailType string           `json:"detail-type"`
	Source     string           `json:"source"`
	Account    string           `json:"account"`
	Time       time.Time        `json:"time"`
	Region     string           `json:"region"`
	Resources  []string         `json:"resources"`
	Detail     SQSMessageDetail `json:"detail"`
}
