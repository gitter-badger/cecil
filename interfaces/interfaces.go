// Licensed to the Apache Software Foundation (ASF) under one or more contributor license agreements;
// and to You under the Apache License, Version 2.0.  See LICENSE in project root for full license + copyright.

package interfaces

import (
	"github.com/jinzhu/gorm"
	"github.com/tleyden/cecil/awstools"
	"github.com/tleyden/cecil/config"
	"github.com/tleyden/cecil/mailers"
	"github.com/tleyden/cecil/models"
	"github.com/tleyden/cecil/queues"
)

// CoreServiceInterface is the interface of core.Service
type CoreServiceInterface interface {
	GormDB() *gorm.DB
	DefaultMailer() *mailers.MailerInstance
	AWSServices() *awstools.AWSServices
	Config() *config.Config
	models.DBServiceInterface
	Queues() queues.QueuesGroupInterface
	//queues.QueuesGroupInterface
}
