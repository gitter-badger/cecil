package core

import (
	"fmt"

	"github.com/jinzhu/gorm"
	mailgun "gopkg.in/mailgun/mailgun-go.v1"
)

// MailerInstance is an instance of a custom mailer
type MailerInstance struct {
	Client      mailgun.Mailgun
	FromAddress string
}

// NewMailerInstance returns a pointer to a new mailer instance
func NewMailerInstance() *MailerInstance {
	return &MailerInstance{}
}

// SetupMailers initializes all the custom mailers
func (s *Service) SetupMailers() error {

	// fetch all accounts from DB
	// for each account, fetch MailerConfig from DB
	// and call service.StartMailerInstance(&mailerConfig)

	var accounts []Account
	err := s.DB.Table("accounts").Find(&accounts).Error
	if err != nil {
		panic(err)
	}

	// start mailer instances for all accounts
	for _, account := range accounts {
		mailerConfig, err := s.FetchMailerConfig(account.ID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				continue
			}
			Logger.Error("error while fetching mailer config", "account", account.ID, "err", err)
			continue
		}
		err = s.StartMailerInstance(mailerConfig)
		if err != nil {
			Logger.Error("error while starting mailer", "account", account.ID, "err", err)
			continue
		}
	}

	return nil
}

// StartMailerInstance starts a MailerInstance and adds it to tracker in Service
func (s *Service) StartMailerInstance(mailerConfig *MailerConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	mailerInst := NewMailerInstance()

	if _, ok := s.mailerInstances[mailerConfig.AccountID]; ok {
		return fmt.Errorf("mailer Instance for account id %v already running", mailerConfig.AccountID)
	}

	mailerInst.Client = mailgun.NewMailgun(
		mailerConfig.Domain,
		mailerConfig.APIKey,
		mailerConfig.PublicAPIKey,
	)
	mailerInst.FromAddress = fmt.Sprintf("%v <noreply@%v>", mailerConfig.FromName, mailerConfig.Domain)

	// add to list of mailer instances
	s.mailerInstances[mailerConfig.AccountID] = mailerInst

	return nil
}

// MailerInstanceByID selects a MailerInstance by account ID
func (s *Service) MailerInstanceByID(accountID uint) (*MailerInstance, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	mailerInst, ok := s.mailerInstances[accountID]
	if !ok {
		return &MailerInstance{}, fmt.Errorf("mailer Instance for account id %v is NOT running", accountID)
	}
	return mailerInst, nil
}

// TerminateMailerInstance terminates an eventual running mailer instance by accountID
func (s *Service) TerminateMailerInstance(accountID uint) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.mailerInstances[accountID]
	if !ok {
		return fmt.Errorf("mailer Instance for account id %v is NOT running", accountID)
	}

	// remove from s.mailerInstances
	delete(s.mailerInstances, accountID)

	// TODO: better cleanup

	return nil
}
