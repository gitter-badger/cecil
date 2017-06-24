// Licensed to the Apache Software Foundation (ASF) under one or more contributor license agreements;
// and to You under the Apache License, Version 2.0.  See LICENSE in project root for full license + copyright.

package config

import "time"

// Config is the struct that holds the main Cecil configuration
type Config struct {
	Server struct {
		Scheme   string // http, or https
		HostName string // hostname for links back to REST API from emails, etc
		Port     string
	}
	Lease struct {
		Duration                  time.Duration
		ApprovalTimeoutDuration   time.Duration
		FirstWarningBeforeExpiry  time.Duration
		SecondWarningBeforeExpiry time.Duration
		MaxPerOwner               int
	}
	DefaultMailer struct {
		Domain       string
		APIKey       string
		PublicAPIKey string
	}
	ProductName string
}
