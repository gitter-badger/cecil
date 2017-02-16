package core

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// SlackBotInstance is an instance of the Slack bot running for a specific user
type SlackBotInstance struct {
	OutgoingMessages chan string
	quit             chan bool
	ChannelID        string
	Token            string
	AccountID        uint
	s                *Service
}

// NewSlackBotInstance returns a pointer to a new slack bot instance
func NewSlackBotInstance() *SlackBotInstance {
	return &SlackBotInstance{
		OutgoingMessages: make(chan string, 20),
		quit:             make(chan bool),
	}
}

// SetupSlack setups slack bots for all the accounts that added slack configuration
func (s *Service) SetupSlack() error {

	// fetch all accounts from DB
	// for each account, fetch SlackConfig from DB
	// and call s.StartSlackBotInstance(&slackConfig)

	var accounts []Account
	err := s.DB.Table("accounts").Find(&accounts).Error
	if err != nil {
		Logger.Error("Error while SetupSlack()", "err", err)
	}

	// start Slack instances for all accounts
	for _, account := range accounts {
		slackConfig, err := s.FetchSlackConfig(account.ID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				continue
			}
			Logger.Error("error while fetching slack config", "account", account.ID, "err", err)
			continue
		}
		err = s.StartSlackBotInstance(slackConfig)
		if err != nil {
			Logger.Error("error while starting slack", "account", account.ID, "err", err)
			continue
		}
	}

	return nil
}

// SlackBotInstanceByID selects a SlackBotInstance by account ID
func (s *Service) SlackBotInstanceByID(accountID uint) (*SlackBotInstance, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	slackIns, ok := s.slackBotInstances[accountID]
	if !ok {
		return &SlackBotInstance{}, fmt.Errorf("Slack instance for account id %v is NOT running", accountID)
	}
	return slackIns, nil
}

// StartSlackBotInstance starts a SlackBotInstance and adds it to tracker in Service
func (s *Service) StartSlackBotInstance(slackConfig *SlackConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	slackInst := NewSlackBotInstance()
	slackInst.Token = slackConfig.Token
	slackInst.ChannelID = slackConfig.ChannelID
	slackInst.AccountID = slackConfig.AccountID
	slackInst.s = s

	if _, ok := s.slackBotInstances[slackConfig.AccountID]; ok {
		return fmt.Errorf("Slack instance for account id %v already running", slackConfig.AccountID)
	}

	// add to list of Slack instances
	s.slackBotInstances[slackConfig.AccountID] = slackInst

	// start Slack instance listening
	go slackInst.StartListenAndServer()
	return nil
}

// Send sends a single message to the Slack channel specifies in the config
func (si *SlackBotInstance) Send(message string) error {
	var err error
	defer func() {
		if x := recover(); x != nil {
			err = fmt.Errorf("Unable to send: %v", x)
		}
	}()
	si.OutgoingMessages <- message
	return err
}

// TerminateSlackBotInstance terminates a slack bot running for a certain account
func (s *Service) TerminateSlackBotInstance(accountID uint) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	slackIns, ok := s.slackBotInstances[accountID]
	if !ok {
		return fmt.Errorf("Slack instance for account id %v was NOT running", accountID)
	}

	// quit goroutine
	slackIns.quit <- true

	// remove from s.slackBotInstances
	delete(s.slackBotInstances, accountID)

	// TODO: better cleanup

	return nil
}
