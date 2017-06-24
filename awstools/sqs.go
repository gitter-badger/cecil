// Licensed to the Apache Software Foundation (ASF) under one or more contributor license agreements;
// and to You under the Apache License, Version 2.0.  See LICENSE in project root for full license + copyright.

package awstools

import (
	"time"
	"strings"
)

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
// SQSMessage is an SQS message, marshaled from SQSEnvelope.Message
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

type SQSPolicyAttribute struct {
	Statement []SQSPolicyStatement `json:"Statement"`
}

type SQSPolicyCondition struct {
	ArnEquals SQSPolicyArnEqualsCondition `json:"ArnEquals"`
}

type SQSPolicyArnEqualsCondition struct {
	SourceArn string `json:"aws:SourceArn"`
}

// Are the arns equal taking wildcards ("*") into account?
// The following arns would be considered equal:
// arn:aws:sns:*:78861:CecilTopic
// arn:aws:sns:us-east-1:78861:CecilTopic
func CheckArnsEqual(arn1, arn2 string) bool {

	sep := ":"
	wildcard := "*"
	arn1Components := strings.Split(arn1, sep)
	arn2Components := strings.Split(arn2, sep)
	if len(arn1Components) != len(arn2Components) {
		return false
	}
	for index, arn1Component := range arn1Components {
		arn2Component := arn2Components[index]
		if arn1Component == wildcard || arn2Component == wildcard {
			continue
		}
		if arn1Component != arn2Component {
			return false
		}
	}
	return true

}