// Licensed to the Apache Software Foundation (ASF) under one or more contributor license agreements;
// and to You under the Apache License, Version 2.0.  See LICENSE in project root for full license + copyright.

package models

import (
	"time"

	"github.com/inconshreveable/log15"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/mattn/go-sqlite3"
	uuid "github.com/satori/go.uuid"
)

// Logger is the logger used in this package; it is initialized by the core package (see core/core-init.go)
var Logger log15.Logger

type DBService struct {
	db *gorm.DB
}

func NewDBService(db *gorm.DB) *DBService {
	return &DBService{db: db}
}

type DBServiceInterface interface {
	GetAllAccounts() ([]Account, error)
	GetAccountByID(accountID int) (*Account, error)
	AccountByEmailExists(accountEmail string) (*Account, bool, error)
	LeasesForAccount(accountID int, terminated *bool) ([]Lease, error)
	LeasesForCloudaccount(cloudaccountID int, terminated *bool) ([]Lease, error)
	LeasesForOwner(ownerEmail string) ([]Lease, error)
	LeaseByGroupUID(accountID uint, cloudaccountID *uint, groupUID string) (*Lease, error)
	LeaseByInstanceID(accountID uint, cloudaccountID *uint, instanceID string) (*Lease, error)
	ActiveInstancesForGroup(accountID uint, cloudaccountID *uint, groupUID string) ([]*Instance, error)
	GetOwnerByID(ownerID uint) (*Owner, error)
	GetOwnerByEmail(email string, cloudaccountID uint) (*Owner, error)
	GetOwnerByKeyName(keyName string, cloudaccountID uint) (*Owner, error)

	GetLeaseByID(leaseID int) (*Lease, error)
	GetCloudaccountByID(cloudaccountID int) (*Cloudaccount, error)
	CloudaccountByAWSIDExists(AWSID string) (bool, error)
	GetSlackConfigForAccount(accountID uint) (*SlackConfig, error)
	GetMailerConfigForAccount(accountID uint) (*MailerConfig, error)
	GetLeaseByUUID(leaseUUID uuid.UUID) (*Lease, error)
	GetInstanceByAWSInstanceID(accountID uint, cloudaccountID *uint, instanceID string) (*Instance, error)
}

// GetAllAccounts fetches all accounts from the DB
func (s *DBService) GetAllAccounts() ([]Account, error) {
	var accounts []Account
	err := s.db.
		Table("accounts").
		Find(&accounts).
		Error
	if err != nil {
		return accounts, err
	}
	return accounts, nil
}

// GetAccountByID fetches from DB a specific account selected by ID
func (s *DBService) GetAccountByID(accountID int) (*Account, error) {

	var account Account

	err := s.db.
		Table("accounts").
		Where("id = ?", uint(accountID)).
		First(&account).
		Error

	return &account, err
}

// AccountByEmailExists returns true if an account with that email address
// exists in the DB. Returns error in case of error.
func (s *DBService) AccountByEmailExists(accountEmail string) (*Account, bool, error) {
	var account Account
	err := s.db.
		Where(&Account{Email: accountEmail}).
		First(&account).
		Error
	if err == gorm.ErrRecordNotFound {
		return nil, false, nil
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, false, err
	}
	return &account, true, nil
}

// GetOwnerByID fetches from DB a specific owner selected by ID
func (s *DBService) GetOwnerByID(ownerID uint) (*Owner, error) {
	var owner Owner
	err := s.db.
		Table("owners").
		Where("id = ?", uint(ownerID)).
		First(&owner).
		Error
	return &owner, err
}

// GetOwnerByEmail fetches from DB a specific owner selected by email
func (s *DBService) GetOwnerByEmail(email string, cloudaccountID uint) (*Owner, error) {
	var owner Owner
	err := s.db.
		Table("owners").
		Where("email = ?", email).
		Where("cloudaccount_id = ?", cloudaccountID).
		First(&owner).
		Error
	return &owner, err
}

// GetOwnerByKeyName fetches from DB a specific owner selected by email
func (s *DBService) GetOwnerByKeyName(keyName string, cloudaccountID uint) (*Owner, error) {
	var owner Owner
	err := s.db.
		Table("owners").
		Where("key_name = ?", keyName).
		Where("cloudaccount_id = ?", cloudaccountID).
		First(&owner).
		Error
	return &owner, err
}

// LeasesForAccount fetches from DB all the leases for a specific
// account.
func (s *DBService) LeasesForAccount(accountID int, terminated *bool) ([]Lease, error) {
	var leases []Lease
	query := s.db.
		Table("leases").
		Where("account_id = ?", uint(accountID))

	if terminated != nil {
		if *terminated {
			// select only terminated leases
			query = query.Where("terminated_at > ?", time.Time{})
		} else {
			// select only non-terminated leases
			query = query.Where("terminated_at IS NULL")
		}
	}

	err := query.Find(&leases).Error
	if err == gorm.ErrRecordNotFound {
		return []Lease{}, err
	}
	return leases, err
}

// LeasesForCloudaccount retches from DB all the leases for a specific
// cloudaccount.
func (s *DBService) LeasesForCloudaccount(cloudaccountID int, terminated *bool) ([]Lease, error) {
	var leases []Lease
	query := s.db.
		Table("leases").
		Where("cloudaccount_id = ?", uint(cloudaccountID))

	if terminated != nil {
		if *terminated {
			// select only terminated leases
			query = query.Where("terminated_at > ?", time.Time{})
		} else {
			// select only non-terminated leases
			query = query.Where("terminated_at IS NULL")
		}
	}
	err := query.Find(&leases).Error
	if err == gorm.ErrRecordNotFound {
		return []Lease{}, err
	}
	return leases, err
}

// LeasesForOwner retches from DB all the leases for a specific
// owner (by owner email address).
func (s *DBService) LeasesForOwner(ownerEmail string) ([]Lease, error) {
	// TODO: verify email

	var owner Owner
	err := s.db.
		Table("owners").
		Where("email = ?", ownerEmail).
		Find(&owner).
		Error

	if err == gorm.ErrRecordNotFound {
		return []Lease{}, err
	}

	var leases []Lease
	err = s.db.
		Table("leases").
		Where("owner_id = ?", owner.ID).
		Find(&leases).
		Error

	if err == gorm.ErrRecordNotFound {
		return []Lease{}, err
	}
	return leases, err
}

// GetLeaseByID fetches from DB a specific lease selected by ID
func (s *DBService) GetLeaseByID(leaseID int) (*Lease, error) {

	var lease Lease
	err := s.db.
		Table("leases").
		Where("id = ?", uint(leaseID)).
		First(&lease).
		Error

	return &lease, err
}

// GetInstanceByAWSInstanceID fetches from DB a specific instance selected by instanceID
func (s *DBService) GetInstanceByAWSInstanceID(accountID uint, cloudaccountID *uint, instanceID string) (*Instance, error) {
	var err error

	search := Instance{
		AccountID:  accountID,
		InstanceID: instanceID,
	}
	if cloudaccountID != nil {
		search.CloudaccountID = *cloudaccountID
	}

	var inst Instance
	err = s.db.
		Table("instances").
		Where(&search).
		Where("terminated_at IS NULL").
		First(&inst).
		Error

	if err != nil {
		return nil, err
	}

	return &inst, err
}

// ActiveInstancesForGroup fetches from DB all non-terminated instances for a group
func (s *DBService) ActiveInstancesForGroup(accountID uint, cloudaccountID *uint, groupUID string) ([]*Instance, error) {
	var err error

	search := Instance{
		AccountID: accountID,
		GroupUID:  groupUID,
	}
	if cloudaccountID != nil {
		search.CloudaccountID = *cloudaccountID
	}

	var inst []*Instance
	err = s.db.
		Table("instances").
		Where(&search).
		Where("terminated_at IS NULL").
		Find(&inst).
		Error

	if err != nil {
		return nil, err
	}

	return inst, err
}

// LeaseByInstanceID fetches a lease by instanceID and other parameters; cloudaccountID is optional.
func (s *DBService) LeaseByInstanceID(accountID uint, cloudaccountID *uint, instanceID string) (*Lease, error) {
	ins, err := s.GetInstanceByAWSInstanceID(accountID, cloudaccountID, instanceID)
	if err != nil {
		return nil, err
	}

	lease, err := s.LeaseByGroupUID(accountID, cloudaccountID, ins.GroupUID)
	if err != nil {
		return nil, err
	}

	return lease, nil
}

// LeaseByGroupUID fetches from DB a specific lease selected by groupUID and other parameters; cloudaccountID is optional.
func (s *DBService) LeaseByGroupUID(accountID uint, cloudaccountID *uint, groupUID string) (*Lease, error) {
	var err error

	search := Lease{
		AccountID: accountID,
		GroupUID:  groupUID,
	}
	if cloudaccountID != nil {
		search.CloudaccountID = *cloudaccountID
	}

	var lease Lease
	err = s.db.
		Table("leases").
		Where(&search).
		Where("terminated_at IS NULL").
		First(&lease).
		Error

	if err != nil {
		return nil, err
	}

	return &lease, err
}

// GetCloudaccountByID fetches a cloudaccount from DB selected by ID.
func (s *DBService) GetCloudaccountByID(cloudaccountID int) (*Cloudaccount, error) {

	var cloudaccount Cloudaccount
	err := s.db.
		Table("cloudaccounts").
		Where("id = ?", uint(cloudaccountID)).
		First(&cloudaccount).
		Error

	return &cloudaccount, err
}

// CloudaccountByAWSIDExists returns true in case a cloudaccount with that AWSID
// exists in the DB.
func (s *DBService) CloudaccountByAWSIDExists(AWSID string) (bool, error) {
	var cloudaccount Cloudaccount
	err := s.db.
		Where(&Cloudaccount{AWSID: AWSID}).
		Find(&cloudaccount).
		Error

	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	return true, nil
}

// GetSlackConfigForAccount fetches the slack config for an account
func (s *DBService) GetSlackConfigForAccount(accountID uint) (*SlackConfig, error) {
	var slackConf SlackConfig
	err := s.db.
		Table("slack_configs").
		Where("account_id = ?", uint(accountID)).
		Find(&slackConf).
		Error

	return &slackConf, err
}

// GetMailerConfigForAccount fetches from DB the custom mailer config for a specific account.
func (s *DBService) GetMailerConfigForAccount(accountID uint) (*MailerConfig, error) {
	var mailerConf MailerConfig
	err := s.db.
		Table("mailer_configs").
		Where("account_id = ?", uint(accountID)).
		Find(&mailerConf).
		Error

	return &mailerConf, err
}

// GetLeaseByUUID returns a lease selected by leaseUUID
func (s *DBService) GetLeaseByUUID(leaseUUID uuid.UUID) (*Lease, error) {
	var lease Lease
	err := s.db.
		Table("leases").
		Where(&Lease{
			UUID: leaseUUID.String(),
		}).
		Where("terminated_at IS NULL").
		First(&lease).Error

	return &lease, err
}
