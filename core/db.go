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

// LeasesForCloudaccount retches from DB all the leases for a specific
// cloudaccount.
func (s *Service) LeasesForCloudaccount(cloudaccountID int, terminated *bool) ([]Lease, error) {
	var native []Lease
	query := s.DB.Table("leases").Where("cloudaccount_id = ?", uint(cloudaccountID))
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

// FetchLeaseByInstanceID fetches from DB a specific lease selected by instance id
func (s *Service) FetchLeaseByInstanceID(instanceID string) (*Lease, error) {
	var err error
	var resourceID uint

	var instance InstanceResource
	instance, err = s.InstanceByInstanceID(instanceID)
	if err != nil {
		return nil, err
	}
	resourceID = instance.ID

	var lease Lease
	err = s.DB.Table("leases").Where(&Lease{
		ResourceID:   resourceID,
		ResourceType: InstanceResourceType,
	}).Where("terminated_at IS NULL").Find(&lease).Error

	if err != nil {
		return nil, err
	}

	return &lease, err
}

// FetchLeaseByStackID fetches from DB a specific lease selected by stack id
func (s *Service) FetchLeaseByStackID(stackID string) (*Lease, error) {
	var err error
	var resourceID uint

	var stack StackResource
	stack, err = s.StackByStackID(stackID)
	if err != nil {
		return nil, err
	}
	resourceID = stack.ID

	var lease Lease
	err = s.DB.Table("leases").Where(&Lease{
		ResourceID:   resourceID,
		ResourceType: StackResourceType,
	}).Where("terminated_at IS NULL").Find(&lease).Error

	if err != nil {
		return nil, err
	}

	return &lease, err
}

// CloudformationHasLease returns a lease in case the specified cloudformation has one
func (s *Service) CloudformationHasLease(accountID int, stackID, stackName string) (*Lease, error) {
	native, err := s.FetchLeaseByStackID(stackID)
	if err == gorm.ErrRecordNotFound {
		return &Lease{}, err
	}
	return native, err
}

// FetchCloudaccountByID fetches a cloudaccount from DB selected by ID.
func (s *Service) FetchCloudaccountByID(cloudaccountID int) (*Cloudaccount, error) {

	// TODO: figure out why it always finds one result, even if none are in the db
	// check whether the cloudaccount exists
	var cloudaccountCount int64
	var cloudaccount Cloudaccount
	err := s.DB.Find(&cloudaccount, uint(cloudaccountID)).Count(&cloudaccountCount).Error
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
func (s *Service) CloudaccountByAWSIDExists(AWSID string) (bool, error) {
	var cloudaccount Cloudaccount
	err := s.DB.Where(&Cloudaccount{AWSID: AWSID}).Find(&cloudaccount).Error
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

// LeaseByUUID returns a lease selected by instanceID and leaseUUID
func (s *Service) LeaseByUUID(leaseUUID uuid.UUID) (*Lease, error) {
	var lease Lease
	err := s.DB.Table("leases").Where(&Lease{
		UUID: leaseUUID.String(),
	}).Where("terminated_at IS NULL").First(&lease).Error

	return &lease, err
}

// IsStack returns true in case the lease is for a stack
func (l *Lease) IsStack() bool {
	return l.ResourceType == StackResourceType
}

// IsInstance returns true in case the lease is for an instance
func (l *Lease) IsInstance() bool {
	return l.ResourceType == InstanceResourceType
}

// IsExpired returns true in case the lease has expired
func (l *Lease) IsExpired() bool {
	return l.ExpiresAt.Before(time.Now().UTC())
}

// ResourceOf returns the resource for which the lease is
func (s *Service) ResourceOf(l *Lease) (interface{}, error) {
	var stack StackResource
	var instance InstanceResource

	if l.ResourceType == StackResourceType {
		err := s.DB.Table(l.ResourceType).Where("lease_id = ?", l.ID).First(&stack).Error
		return stack, err
	}

	if l.ResourceType == InstanceResourceType {
		err := s.DB.Table(l.ResourceType).Where("lease_id = ?", l.ID).First(&instance).Error
		return instance, err
	}

	return nil, gorm.ErrRecordNotFound
}

// StackByStackID returns the stack resource for which the lease is
func (s *Service) StackByStackID(AWSStackID string) (StackResource, error) {
	var stack StackResource
	err := s.DB.Table(StackResourceType).Where("stack_id = ?", AWSStackID).First(&stack).Error
	return stack, err
}

// InstanceByInstanceID returns the instance resource for which the lease is
func (s *Service) InstanceByInstanceID(AWSInstanceID string) (InstanceResource, error) {
	var instance InstanceResource
	err := s.DB.Table(InstanceResourceType).Where("instance_id = ?", AWSInstanceID).First(&instance).Error
	return instance, err
}
