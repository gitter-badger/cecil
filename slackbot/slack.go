package slackbot

import (
	"fmt"
	"sync"

	"github.com/inconshreveable/log15"
	"github.com/jinzhu/gorm"
	"github.com/tleyden/cecil/interfaces"
	"github.com/tleyden/cecil/models"
)

var Logger log15.Logger

type SlackBotService struct {
	mu                *sync.RWMutex
	slackBotInstances map[uint]*SlackBotInstance // map account_id to *SlackBotInstance
	s                 interfaces.CoreServiceInterface
}

// NewService returns a pointer to a new slack bot service
func NewService(s interfaces.CoreServiceInterface) *SlackBotService {
	return &SlackBotService{
		mu:                &sync.RWMutex{},
		slackBotInstances: make(map[uint]*SlackBotInstance),
		s:                 s,
	}
}

// SlackBotInstance is an instance of the Slack bot running for a specific user
type SlackBotInstance struct {
	OutgoingMessages chan string
	quit             chan bool
	ChannelID        string
	Token            string
	AccountID        uint
	s                interfaces.CoreServiceInterface
}

// NewSlackBotInstance returns a pointer to a new slack bot instance
func NewSlackBotInstance() *SlackBotInstance {
	return &SlackBotInstance{
		OutgoingMessages: make(chan string, 20),
		quit:             make(chan bool),
	}
}

// SetupSlack setups slack bots for all the accounts that added slack configuration
func (s *SlackBotService) SetupSlack() error {

	// fetch all accounts from DB
	// for each account, fetch SlackConfig from DB
	// and call s.StartSlackBotInstance(&slackConfig)

	accounts, err := s.s.FetchAllAccounts()
	if err != nil {
		return err
	}

	// start Slack instances for all accounts
	for _, account := range accounts {
		slackConfig, err := s.s.FetchSlackConfig(account.ID)
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
func (s *SlackBotService) SlackBotInstanceByID(accountID uint) (*SlackBotInstance, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	slackIns, ok := s.slackBotInstances[accountID]
	if !ok {
		return &SlackBotInstance{}, fmt.Errorf("Slack instance for account id %v is NOT running", accountID)
	}
	return slackIns, nil
}

// StartSlackBotInstance starts a SlackBotInstance and adds it to tracker in Service
func (s *SlackBotService) StartSlackBotInstance(slackConfig *models.SlackConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	slackInst := NewSlackBotInstance()
	slackInst.Token = slackConfig.Token
	slackInst.ChannelID = slackConfig.ChannelID
	slackInst.AccountID = slackConfig.AccountID
	slackInst.s = s.s

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
func (s *SlackBotService) TerminateSlackBotInstance(accountID uint) error {
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

func (s *SlackBotService) StopAll() {
	// terminate all slackBotInstances
	for accountID := range s.slackBotInstances {
		err := s.TerminateSlackBotInstance(accountID)
		if err != nil {
			Logger.Warn("Error terminating Slack instance: %v", err)
		}
	}
	return
}
