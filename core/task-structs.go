package core

import (
	"time"
)

// @@@@@@@@@@@@@@@ Task structs @@@@@@@@@@@@@@@

type NewLeaseTask struct {
	AWSAccountID string // message.account
	InstanceID   string // message.detail.instance-id
	Region       string // message.region

	LaunchedAt    time.Time // get from the request for tags to ec2 api, not from event
	InstanceType  string
	InstanceOwner string
	//InstanceTags []string
}

type TerminatorTask struct {
	AWSAccountID string
	InstanceID   string
	Region       string // needed? arn:aws:ec2:us-east-1:859795398601:instance/i-fd1f96cc

	Action string // default is TerminatorActionTerminate
}

type LeaseTerminatedTask struct {
}

type RenewerTask struct {
}

type NotifierTask struct {
	From     string
	To       string
	Subject  string
	BodyHTML string
	BodyText string
}
