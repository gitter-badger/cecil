package config

import "time"

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
