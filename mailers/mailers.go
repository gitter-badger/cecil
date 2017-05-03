package mailers

import (
	"fmt"
	"sync"

	"github.com/inconshreveable/log15"
	"github.com/jinzhu/gorm"
	"github.com/tleyden/cecil/models"
	mailgun "gopkg.in/mailgun/mailgun-go.v1"
)

// Logger is the logger used in this package; it is initialized by the core package (see core/core-init.go)
var Logger log15.Logger

type CustomMailerService struct {
	db              *models.DBService
	mailerInstances map[uint]*MailerInstance // map account_id to *MailerInstance
	mu              *sync.RWMutex
}

func NewService(db *models.DBService) *CustomMailerService {
	return &CustomMailerService{
		db:              db,
		mailerInstances: make(map[uint]*MailerInstance),
		mu:              &sync.RWMutex{},
	}
}

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
func (s *CustomMailerService) SetupMailers() error {

	// fetch all accounts from DB
	// for each account, fetch MailerConfig from DB
	// and call service.StartMailerInstance(&mailerConfig)

	accounts, err := s.db.GetAllAccounts()
	if err != nil {
		return err
	}
	// start mailer instances for all accounts
	for _, account := range accounts {
		mailerConfig, err := s.db.GetMailerConfigForAccount(account.ID)
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
func (s *CustomMailerService) StartMailerInstance(mailerConfig *models.MailerConfig) error {
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
func (s *CustomMailerService) MailerInstanceByID(accountID uint) (*MailerInstance, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	mailerInst, ok := s.mailerInstances[accountID]
	if !ok {
		return &MailerInstance{}, fmt.Errorf("mailer Instance for account id %v is NOT running", accountID)
	}
	return mailerInst, nil
}

// TerminateMailerInstance terminates an eventual running mailer instance by accountID
func (s *CustomMailerService) TerminateMailerInstance(accountID uint) error {
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

func (s *CustomMailerService) StopAll() {
	// terminate all mailerInstances
	for accountID := range s.mailerInstances {
		err := s.TerminateMailerInstance(accountID)
		if err != nil {
			Logger.Warn("Error terminating mailer instance: %v", err)
		}
	}
	return
}
