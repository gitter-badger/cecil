// Licensed to the Apache Software Foundation (ASF) under one or more contributor license agreements;
// and to You under the Apache License, Version 2.0.  See LICENSE in project root for full license + copyright.

package core

import (
	"crypto/rsa"
	"sync"

	"github.com/jinzhu/gorm"
	"github.com/tleyden/cecil/awstools"
	"github.com/tleyden/cecil/config"
	"github.com/tleyden/cecil/mailers"
	"github.com/tleyden/cecil/models"
	"github.com/tleyden/cecil/queues"
	"github.com/tleyden/cecil/slackbot"
)

// Service is fundamental element of Cecil, and holds most of what is used by Cecil.
type Service struct {
	queues *queues.QueuesGroup

	config *config.Config

	DB *gorm.DB
	*models.DBService

	defaultMailer *mailers.MailerInstance
	*mailers.CustomMailerService

	*slackbot.SlackBotService

	AWS awstools.AWSServices
	rsa struct {
		publicKey  *rsa.PublicKey
		privateKey *rsa.PrivateKey
	}


	mu *sync.RWMutex
}

// NewService returns a new service
func NewService() *Service {
	service := &Service{
		mu:          &sync.RWMutex{},
	}
	return service
}

// GormDB returns *gorm.DB of Service
func (s *Service) GormDB() *gorm.DB {
	return s.DB
}


// AWSServices returns AWSServices
func (s *Service) AWSServices() *awstools.AWSServices {

	return &s.AWS
}

// DefaultMailer returns defaultMailer
func (s *Service) DefaultMailer() *mailers.MailerInstance {
	return s.defaultMailer
}

// Config returns config
func (s *Service) Config() *config.Config {
	return s.config
}

// Queues returns queues
func (s *Service) Queues() queues.QueuesGroupInterface {
	return s.queues
}
