package core

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	uuid "github.com/satori/go.uuid"
)

// FetchAccountByID fetches from DB a specific account selected by ID
func (s *Service) FetchAccountByID(accountID int) (*Account, error) {

	// TODO: figure out why it always finds one result, even if none are in the db
	// check whether the account exists
	var accountCount int64
	var account Account
	err := s.DB.Table("accounts").Where("id = ?", uint(accountID)).Count(&accountCount).Find(&account).Error
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
func (s *Service) AccountByEmailExists(accountEmail string) (*Account, bool, error) {
	var account Account
	err := s.DB.Where(&Account{Email: accountEmail}).Find(&account).Error
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
func (s *Service) LeasesForAccount(accountID int, terminated *bool) ([]Lease, error) {
	var native []Lease
	query := s.DB.Table("leases").Where("account_id = ?", uint(accountID))
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

// LeasesForCloudAccount retches from DB all the leases for a specific
// cloudAccount.
func (s *Service) LeasesForCloudAccount(cloudAccountID int, terminated *bool) ([]Lease, error) {
	var native []Lease
	query := s.DB.Table("leases").Where("cloud_account_id = ?", uint(cloudAccountID))
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
func (s *Service) LeasesForOwner(ownerEmail string) ([]Lease, error) {
	// TODO: verify email

	var owner Owner
	err := s.DB.Table("owners").Where("email = ?", ownerEmail).Find(&owner).Error
	if err == gorm.ErrRecordNotFound {
		return []Lease{}, err
	}

	var native []Lease
	err = s.DB.Table("leases").Where("owner_id = ?", owner.ID).Find(&native).Error
	if err == gorm.ErrRecordNotFound {
		return []Lease{}, err
	}
	return native, err
}

// FetchLeaseByID fetches from DB a specific lease selected by ID
func (s *Service) FetchLeaseByID(leaseID int) (*Lease, error) {

	// TODO: figure out why it always finds one result, even if none are in the db
	// check whether the lease exists
	var leaseCount int64
	var lease Lease
	err := s.DB.Table("leases").Where("id = ?", uint(leaseID)).Count(&leaseCount).Find(&lease).Error
	if err == gorm.ErrRecordNotFound {
		return &Lease{}, err
	}

	if uint(leaseID) != lease.ID {
		return &Lease{}, gorm.ErrRecordNotFound
	}

	return &lease, err
}

// CloudformationHasLease returns a lease in case the specified cloudformation has one
func (s *Service) CloudformationHasLease(accountID int, stackID, stackName string) (*Lease, error) {
	var native Lease
	query := s.DB.Table("leases").
		Where("account_id = ?", uint(accountID)).
		Where("stack_id = ?", stackID).
		Where("stack_name = ?", stackName)

	err := query.Find(&native).Error
	if err == gorm.ErrRecordNotFound {
		return &Lease{}, err
	}
	return &native, err
}

// FetchCloudAccountByID fetches a cloudaccount from DB selected by ID.
func (s *Service) FetchCloudAccountByID(cloudAccountID int) (*CloudAccount, error) {

	// TODO: figure out why it always finds one result, even if none are in the db
	// check whether the cloudAccount exists
	var cloudAccountCount int64
	var cloudAccount CloudAccount
	err := s.DB.Find(&cloudAccount, uint(cloudAccountID)).Count(&cloudAccountCount).Error
	if err == gorm.ErrRecordNotFound {
		return &CloudAccount{}, err
	}

	if uint(cloudAccountID) != cloudAccount.ID {
		return &CloudAccount{}, gorm.ErrRecordNotFound
	}

	return &cloudAccount, err
}

// CloudAccountByAWSIDExists returns true in case a cloudaccount with that AWSID
// exists in the DB.
func (s *Service) CloudAccountByAWSIDExists(AWSID string) (bool, error) {
	var cloudAccount CloudAccount
	err := s.DB.Where(&CloudAccount{AWSID: AWSID}).Find(&cloudAccount).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	return true, nil
}

// IsOwnerOf returns true if the account owns the cloudaccount.
func (a *Account) IsOwnerOf(cloudAccount *CloudAccount) bool {
	return a.ID == cloudAccount.AccountID
}

// IsOwnerOf returns true if the cloudaccount owns the lease.
func (a *CloudAccount) IsOwnerOf(lease *Lease) bool {
	return a.ID == lease.CloudAccountID
}

// IsOwnerOfLease returns true if the account owns the lease
func (a *Account) IsOwnerOfLease(lease *Lease) bool {
	return a.ID == lease.AccountID
}

// FetchSlackConfig fetches the slack config for an account
func (s *Service) FetchSlackConfig(accountID uint) (*SlackConfig, error) {
	var slackConf SlackConfig
	err := s.DB.Table("slack_configs").Where("account_id = ?", uint(accountID)).Find(&slackConf).Error
	if err == gorm.ErrRecordNotFound {
		return &SlackConfig{}, err
	}

	if uint(accountID) != slackConf.AccountID {
		return &SlackConfig{}, gorm.ErrRecordNotFound
	}

	return &slackConf, err
}

func (s *Service) FetchMailerConfig(accountID uint) (*MailerConfig, error) {
	var maileConf MailerConfig
	err := s.DB.Table("mailer_configs").Where("account_id = ?", uint(accountID)).Find(&maileConf).Error
	if err == gorm.ErrRecordNotFound {
		return &MailerConfig{}, err
	}

	if uint(accountID) != maileConf.AccountID {
		return &MailerConfig{}, gorm.ErrRecordNotFound
	}

	return &maileConf, err
}

// CloudAccountByAWSIDExists returns true in case a cloudaccount with that AWSID
// exists in the DB.
func (s *Service) LeaseByIDAndUUID(instanceID string, leaseUUID uuid.UUID) (*Lease, error) {
	var lease Lease
	err := s.DB.Table("leases").Where(&Lease{
		InstanceID: instanceID,
		UUID:       leaseUUID.String(),
	}).Where("terminated_at IS NULL").First(&lease).Error

	return &lease, err
}
