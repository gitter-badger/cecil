package models

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	uuid "github.com/satori/go.uuid"
)

type DBService struct {
	db *gorm.DB
}

func NewDBService(db *gorm.DB) *DBService {
	return &DBService{db: db}
}

type DBServiceInterface interface {
	FetchAllAccounts() ([]Account, error)
	FetchAccountByID(accountID int) (*Account, error)
	AccountByEmailExists(accountEmail string) (*Account, bool, error)
	LeasesForAccount(accountID int, terminated *bool) ([]Lease, error)
	LeasesForCloudaccount(cloudaccountID int, terminated *bool) ([]Lease, error)
	LeasesForOwner(ownerEmail string) ([]Lease, error)
	LeaseByGroupUID(accountID uint, cloudaccountID *uint, groupUID string) (*Lease, error)
	LeaseByInstanceID(accountID uint, cloudaccountID *uint, instanceID string) (*Lease, error)
	ActiveInstancesForGroup(accountID uint, cloudaccountID *uint, groupUID string) ([]*Instance, error)

	FetchLeaseByID(leaseID int) (*Lease, error)
	FetchCloudaccountByID(cloudaccountID int) (*Cloudaccount, error)
	CloudaccountByAWSIDExists(AWSID string) (bool, error)
	FetchSlackConfig(accountID uint) (*SlackConfig, error)
	FetchMailerConfig(accountID uint) (*MailerConfig, error)
	LeaseByUUID(leaseUUID uuid.UUID) (*Lease, error)
	InstanceByInstanceID(accountID uint, cloudaccountID *uint, instanceID string) (*Instance, error)
}

func (l *Lease) SetGroupType(groupType GroupType) {
	l.GroupType = groupType
}

func (l *Lease) SetGroupUID(groupUID string) {
	l.GroupUID = groupUID
}

func (s *DBService) FetchAllAccounts() ([]Account, error) {
	var accounts []Account
	err := s.db.Table("accounts").Find(&accounts).Error
	if err != nil {
		return accounts, err
	}
	return accounts, nil
}

// FetchAccountByID fetches from DB a specific account selected by ID
func (s *DBService) FetchAccountByID(accountID int) (*Account, error) {

	// TODO: figure out why it always finds one result, even if none are in the db
	// check whether the account exists
	var accountCount int64
	var account Account
	err := s.db.Table("accounts").Where("id = ?", uint(accountID)).Count(&accountCount).Find(&account).Error
	if err == gorm.ErrRecordNotFound {
		return &Account{}, err
	}

	if uint(accountID) != account.ID {
		return &Account{}, gorm.ErrRecordNotFound
	}

	return &account, err
}

// AccountByEmailExists returns true if an account with that email address
// exists in the DB. Returns error in case of error.
func (s *DBService) AccountByEmailExists(accountEmail string) (*Account, bool, error) {
	var account Account
	err := s.db.Where(&Account{Email: accountEmail}).Find(&account).Error
	if err == gorm.ErrRecordNotFound {
		return nil, false, nil
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, false, err
	}
	return &account, true, nil
}

// LeasesForAccount fetches from DB all the leases for a specific
// account.
func (s *DBService) LeasesForAccount(accountID int, terminated *bool) ([]Lease, error) {
	var native []Lease
	query := s.db.Table("leases").Where("account_id = ?", uint(accountID))
	if terminated != nil {
		if *terminated {
			// select only terminated leases
			query = query.Where("terminated_at > ?", time.Time{})
		} else {
			// select only non-terminated leases
			query = query.Where("terminated_at IS NULL")
		}
	}
	err := query.Find(&native).Error
	if err == gorm.ErrRecordNotFound {
		return []Lease{}, err
	}
	return native, err
}

// LeasesForCloudaccount retches from DB all the leases for a specific
// cloudaccount.
func (s *DBService) LeasesForCloudaccount(cloudaccountID int, terminated *bool) ([]Lease, error) {
	var native []Lease
	query := s.db.Table("leases").Where("cloudaccount_id = ?", uint(cloudaccountID))
	if terminated != nil {
		if *terminated {
			// select only terminated leases
			query = query.Where("terminated_at > ?", time.Time{})
		} else {
			// select only non-terminated leases
			query = query.Where("terminated_at IS NULL")
		}
	}
	err := query.Find(&native).Error
	if err == gorm.ErrRecordNotFound {
		return []Lease{}, err
	}
	return native, err
}

// LeasesForOwner retches from DB all the leases for a specific
// owner (by owner email address).
func (s *DBService) LeasesForOwner(ownerEmail string) ([]Lease, error) {
	// TODO: verify email

	var owner Owner
	err := s.db.Table("owners").Where("email = ?", ownerEmail).Find(&owner).Error
	if err == gorm.ErrRecordNotFound {
		return []Lease{}, err
	}

	var native []Lease
	err = s.db.Table("leases").Where("owner_id = ?", owner.ID).Find(&native).Error
	if err == gorm.ErrRecordNotFound {
		return []Lease{}, err
	}
	return native, err
}

// FetchLeaseByID fetches from DB a specific lease selected by ID
func (s *DBService) FetchLeaseByID(leaseID int) (*Lease, error) {

	// TODO: figure out why it always finds one result, even if none are in the db
	// check whether the lease exists
	var leaseCount int64
	var lease Lease
	err := s.db.Table("leases").Where("id = ?", uint(leaseID)).Count(&leaseCount).Find(&lease).Error
	if err == gorm.ErrRecordNotFound {
		return &Lease{}, err
	}

	if uint(leaseID) != lease.ID {
		return &Lease{}, gorm.ErrRecordNotFound
	}

	return &lease, err
}

// InstanceByInstanceID fetches from DB a specific instance selected by instanceID
func (s *DBService) InstanceByInstanceID(accountID uint, cloudaccountID *uint, instanceID string) (*Instance, error) {
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

func (s *DBService) LeaseByInstanceID(accountID uint, cloudaccountID *uint, instanceID string) (*Lease, error) {
	ins, err := s.InstanceByInstanceID(accountID, cloudaccountID, instanceID)
	if err != nil {
		return nil, err
	}

	lease, err := s.LeaseByGroupUID(accountID, cloudaccountID, ins.GroupUID)
	if err != nil {
		return nil, err
	}

	return lease, nil
}

// LeaseByGroupUID fetches from DB a specific lease selected by groupUID
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

// FetchCloudaccountByID fetches a cloudaccount from DB selected by ID.
func (s *DBService) FetchCloudaccountByID(cloudaccountID int) (*Cloudaccount, error) {

	// TODO: figure out why it always finds one result, even if none are in the db
	// check whether the cloudaccount exists
	var cloudaccountCount int64
	var cloudaccount Cloudaccount
	err := s.db.Find(&cloudaccount, uint(cloudaccountID)).Count(&cloudaccountCount).Error
	if err == gorm.ErrRecordNotFound {
		return &Cloudaccount{}, err
	}

	if uint(cloudaccountID) != cloudaccount.ID {
		return &Cloudaccount{}, gorm.ErrRecordNotFound
	}

	return &cloudaccount, err
}

// CloudaccountByAWSIDExists returns true in case a cloudaccount with that AWSID
// exists in the DB.
func (s *DBService) CloudaccountByAWSIDExists(AWSID string) (bool, error) {
	var cloudaccount Cloudaccount
	err := s.db.Where(&Cloudaccount{AWSID: AWSID}).Find(&cloudaccount).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	return true, nil
}

// IsOwnerOf returns true if the account owns the cloudaccount.
func (a *Account) IsOwnerOf(cloudaccount *Cloudaccount) bool {
	return a.ID == cloudaccount.AccountID
}

// IsOwnerOf returns true if the cloudaccount owns the lease.
func (a *Cloudaccount) IsOwnerOf(lease *Lease) bool {
	return a.ID == lease.CloudaccountID
}

// IsOwnerOfLease returns true if the account owns the lease
func (a *Account) IsOwnerOfLease(lease *Lease) bool {
	return a.ID == lease.AccountID
}

// FetchSlackConfig fetches the slack config for an account
func (s *DBService) FetchSlackConfig(accountID uint) (*SlackConfig, error) {
	var slackConf SlackConfig
	err := s.db.Table("slack_configs").Where("account_id = ?", uint(accountID)).Find(&slackConf).Error
	if err == gorm.ErrRecordNotFound {
		return &SlackConfig{}, err
	}

	if uint(accountID) != slackConf.AccountID {
		return &SlackConfig{}, gorm.ErrRecordNotFound
	}

	return &slackConf, err
}

func (s *DBService) FetchMailerConfig(accountID uint) (*MailerConfig, error) {
	var mailerConf MailerConfig
	err := s.db.Table("mailer_configs").Where("account_id = ?", uint(accountID)).Find(&mailerConf).Error
	if err == gorm.ErrRecordNotFound {
		return &MailerConfig{}, err
	}

	if uint(accountID) != mailerConf.AccountID {
		return &MailerConfig{}, gorm.ErrRecordNotFound
	}

	return &mailerConf, nil
}

// LeaseByUUID returns a lease selected by instanceID and leaseUUID
func (s *DBService) LeaseByUUID(leaseUUID uuid.UUID) (*Lease, error) {
	var lease Lease
	err := s.db.Table("leases").Where(&Lease{
		UUID: leaseUUID.String(),
	}).Where("terminated_at IS NULL").First(&lease).Error

	return &lease, err
}

// IsExpired returns true in case the lease has expired
func (l *Lease) IsExpired() bool {
	return l.ExpiresAt.Before(time.Now().UTC())
}
