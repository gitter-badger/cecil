// Licensed to the Apache Software Foundation (ASF) under one or more contributor license agreements;
// and to You under the Apache License, Version 2.0.  See LICENSE in project root for full license + copyright.

package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// BasicModel implements the basic fields of a DB model, along with their JSON tags
type BasicModel struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"-"`
}

// @@@@@@@@@@@@@@@ DB models @@@@@@@@@@@@@@@

// Account is a Cecil account. The owner of an account is called an "Admin"
// because he/she administers one or more cloudaccounts.
type Account struct {
	gorm.Model `json:"-"`
	Email      string `sql:"size:255;unique;index" json:"email"`

	Name    string `sql:"size:255" json:"name"`
	Surname string `sql:"size:255" json:"surname"`

	Verified          bool   `sql:"DEFAULT:false" json:"verified"`
	VerificationToken string `sql:"unique" json:"-"`
	RequestedNewToken bool   `sql:"DEFAULT:false" json:"-"`

	Disabled bool `sql:"DEFAULT:false" json:"-"`
	Deleted  bool `sql:"DEFAULT:false" json:"-"`

	DefaultLeaseDuration time.Duration `sql:"DEFAULT:0" json:"-"`

	Cloudaccounts []Cloudaccount `json:"-"`
	SlackConfig   SlackConfig    `json:"-"`
}

// Cloudaccount is a cloud service account; e.g. an AWS account.
type Cloudaccount struct {
	gorm.Model
	AccountID uint `json:"account_id"`

	DefaultLeaseDuration time.Duration `sql:"DEFAULT:0" json:"default_lease_duration"`
	Provider             string        `json:"provider"` // e.g. AWS
	AWSID                string        `sql:"size:255;unique;index" json:"aws_id"`
	ExternalID           string        `json:"-"`

	Disabled bool `sql:"DEFAULT:false" json:"disabled"`
	Deleted  bool `sql:"DEFAULT:false" json:"deleted"`

	Leases []Lease `json:"-"`
	Owners []Owner `json:"-"`
}

// Owner is a whitelisted email address that can/owns (i.e. is responsible of) leases.
type Owner struct {
	gorm.Model
	CloudaccountID uint

	Email   string
	KeyName string

	Disabled bool `sql:"DEFAULT:false"`
	Leases   []Lease
}

// Instance is the DB model for an instance
type Instance struct {
	BasicModel // using this instead of gorm.Model because gorm.Model does not have json tags
	LeaseID    uint

	AccountID      uint   `json:"account_id,omitempty"`
	CloudaccountID uint   `json:"cloudaccount_id,omitempty"`
	AWSAccountID   string `json:"aws_account_id,omitempty"`

	GroupUID  string    `json:"group_id,omitempty"`
	GroupType GroupType `json:"group_type,omitempty"`

	InstanceID       string     `json:"instance_id,omitempty"`
	AvailabilityZone string     `json:"availability_zone,omitempty"`
	InstanceType     string     `json:"instance_type,omitempty"`
	Region           string     `json:"region,omitempty"`
	LaunchedAt       time.Time  `json:"launched_at,omitempty"`
	TerminatedAt     *time.Time `json:"terminated_at,omitempty"`
}

type GroupType int

const (
	GroupCecilGroupTag GroupType = iota
	GroupASG
	GroupCF
	GroupTime
	GroupSingle
)

func (gt GroupType) String() string {
	switch gt {
	case GroupCecilGroupTag:
		return "GroupCecilGroupTag"
	case GroupASG:
		return "GroupASG"
	case GroupCF:
		return "GroupCF"
	case GroupTime:
		return "GroupTime"
	case GroupSingle:
		return "GroupSingle"
	default:
		return ""
	}
}

func (gt GroupType) EmailDisplayString() string {
	switch gt {
	case GroupCecilGroupTag:
		return "GroupTag"
	case GroupASG:
		return "AutoScalingGroup"
	case GroupCF:
		return "Cloudformation"
	case GroupTime:
		return "GroupTime"
	case GroupSingle:
		return "EC2 Instance"
	default:
		return ""
	}
}

// Leases are the core functionality of Cecil. They establish a contract between a user of AWS (usually a developer,
// tester, or sys. admin) and the AWS resources they create, to make sure that the user is aware of the resource
// and still requires it
type Lease struct {
	BasicModel // using this instead of gorm.Model because gorm.Model does not have json tags

	UUID      string `sql:"size:255;unique;index" json:"-"` // TODO: sho or hide UUID on json responses?
	TokenOnce string `json:"-"`

	AccountID      uint   `json:"account_id,omitempty"`
	CloudaccountID uint   `json:"cloudaccount_id,omitempty"`
	OwnerID        uint   `json:"owner_id,omitempty"`
	AWSAccountID   string `json:"aws_account_id,omitempty"`

	GroupUID  string    `json:"group_id,omitempty"`
	GroupType GroupType `json:"group_type,omitempty"`

	Region           string `json:"region,omitempty"` //TODO: make sure regions is set
	AwsContainerName string `json:"aws_container_name,omitempty"`

	Deleted                     bool `json:"deleted,omitempty"`
	NumTimesAllertedAboutExpiry int  `json:"-"`

	LaunchedAt   time.Time  `json:"launched_at,omitempty"`
	ExpiresAt    time.Time  `json:"expires_at,omitempty"`
	ApprovedAt   *time.Time `json:"approved_at,omitempty"`
	TerminatedAt *time.Time `json:"terminated_at,omitempty"`
}

// SlackConfig contains the configuration used to setup slack
type SlackConfig struct {
	gorm.Model
	AccountID uint
	Token     string
	ChannelID string
}

// MailerConfig contains the configuration of a custom mailer that will be used instead of
// the default one.
type MailerConfig struct {
	gorm.Model
	AccountID uint

	Domain       string
	APIKey       string
	PublicAPIKey string
	FromName     string
}
