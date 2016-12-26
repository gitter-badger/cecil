package core

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/nlopes/slack"
)

type SlackInstance struct {
	OutgoingMessages chan string
	quit             chan bool
	ChannelID        string
	Token            string
	AccountID        uint
	s                *Service
}

func NewSlackInstance() *SlackInstance {
	return &SlackInstance{
		OutgoingMessages: make(chan string, 20),
		quit:             make(chan bool),
	}
}

func (service *Service) SetupSlack() error {

	// fetch all accounts from DB
	// for each account, fetch SlackConfig from DB
	// and call service.StartSlackInstance(&slackConfig)

	var accounts []Account
	err := service.DB.Table("accounts").Find(&accounts).Error
	if err != nil {
		Logger.Error("Error while SetupSlack()", "err", err)
	}

	// start Slack instances for all accounts
	for _, account := range accounts {
		slackConfig, err := service.FetchSlackConfig(account.ID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				continue
			}
			Logger.Error("error while fetching slack config", "account", account.ID, "err", err)
			continue
		}
		err = service.StartSlackInstance(slackConfig)
		if err != nil {
			Logger.Error("error while starting slack", "account", account.ID, "err", err)
			continue
		}
	}

	// debug:
	go func() {
		time.Sleep(time.Second * 10)
		slackInst, err := service.SlackInstanceByID(1)
		if err != nil {
			Logger.Warn("SlackInstanceByID", "err", err)
			return
		}
		slackInst.Send("Hello world from Cecil")
	}()

	return nil
}

// SlackInstanceByID selects a SlackInstance by account ID
func (s *Service) SlackInstanceByID(accountID uint) (*SlackInstance, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	slackIns, ok := s.slackInstances[accountID]
	if !ok {
		return &SlackInstance{}, fmt.Errorf("Slack instance for account id %v is NOT running", accountID)
	}
	return slackIns, nil
}

// StartSlackInstance starts a SlackInstance and adds it to tracker in Service
func (s *Service) StartSlackInstance(slackConfig *SlackConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	slackInst := NewSlackInstance()
	slackInst.Token = slackConfig.Token
	slackInst.ChannelID = slackConfig.ChannelID
	slackInst.AccountID = slackConfig.AccountID
	slackInst.s = s

	if _, ok := s.slackInstances[slackConfig.AccountID]; ok {
		return fmt.Errorf("Slack instance for account id %v already running", slackConfig.AccountID)
	}

	// add to list of Slack instances
	s.slackInstances[slackConfig.AccountID] = slackInst

	// start Slack instance listening
	go slackInst.StartListen()
	return nil
}

func (si *SlackInstance) Send(message string) error {
	var err error
	defer func() {
		if x := recover(); x != nil {
			err = fmt.Errorf("Unable to send: %v", x)
		}
	}()
	si.OutgoingMessages <- message
	return err
}

func (s *Service) TerminateSlackInstance(accountID uint) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	slackIns, ok := s.slackInstances[accountID]
	if !ok {
		return fmt.Errorf("Slack instance for account id %v was NOT running", accountID)
	}

	// quit goroutine
	slackIns.quit <- true

	// remove from s.slackInstances
	delete(s.slackInstances, accountID)

	// TODO: better cleanup

	return nil
}

func (si *SlackInstance) StartListen() {
	api := slack.New(si.Token)
	api.SetDebug(true)
	api.SetUserAsActive()

	rtm := api.NewRTM()
	go rtm.ManageConnection()

Loop:
	for {
		select {
		// received quit signal
		case <-si.quit:
			{
				// TODO: cleanup
				// TODO: say bye to group
				Logger.Info("quitting Slack Listen Loop")
				return
			}
		case outgoingMessage := <-si.OutgoingMessages:
			{
				params := slack.PostMessageParameters{}
				params.Attachments = []slack.Attachment{}
				params.AsUser = true
				channelID, timestamp, err := rtm.PostMessage(si.ChannelID, outgoingMessage, params)
				// TODO: handle {"ok":false,"error":"not_in_channel"}
				if err != nil {
					Logger.Error("Error while posting message", "err", err)
					return
				}
				_, _ = channelID, timestamp
			}

		case incomingEvent := <-rtm.IncomingEvents:
			switch incomingEvent.Data.(type) {

			case *slack.MessageEvent:
				incomingMessage := incomingEvent.Data.(*slack.MessageEvent)

				botIdentity := rtm.GetInfo()

				thisBotHasBeenTaggedInMesage := strings.Contains(incomingMessage.Text, fmt.Sprintf("<@%v>", botIdentity.User.ID))
				if !thisBotHasBeenTaggedInMesage {
					continue
				}

				rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("Hey <@%v>, you mentioned me!", incomingMessage.User), incomingMessage.Channel))
				response, err := si.HandleMessage(incomingMessage)
				if err != nil {
					rtm.SendMessage(rtm.NewOutgoingMessage(
						fmt.Sprintf(
							"<@%v>, an error occured while executing your query:\n\n %v",
							incomingMessage.User,
							err.Error(),
						),
						incomingMessage.Channel))
				} else {
					rtm.SendMessage(rtm.NewOutgoingMessage(
						fmt.Sprintf(
							"<@%v>, response:\n\n %v",
							incomingMessage.User,
							response,
						),
						incomingMessage.Channel))
				}

			case *slack.InvalidAuthEvent:
				Logger.Error("Invalid credentials")
				break Loop

			}

		}

	}
}

func (si *SlackInstance) HandleMessage(message *slack.MessageEvent) (string, error) {
	var err error

	// TODO:
	// help
	// list leases
	// terminate leaseid
	// renew leaseid

	switch {
	case strings.Contains(message.Text, "help"):
		{
			return si.Usage(), nil
		}
	case strings.Contains(message.Text, "list leases"):
		{
			return si.ListLeases()
		}
	default:
		{
			return "", errors.New("command not found")
		}

	}

	return "", err
}

func (si *SlackInstance) Usage() string {
	return `
Available commands:

*help* - This command.
*list leases* - List all leases for this account.

`
}

func (si *SlackInstance) ListLeases() (string, error) {

	// fetch leases for account
	leases, err := si.s.LeasesForAccount(int(si.AccountID), nil)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", errors.New("no leases found")
		} else {
			return "", errors.New("internal error")
		}
	}

	var response bytes.Buffer

	for _, lease := range leases {
		leaseLine := fmt.Sprintf(
			"*%v* (AWS %v): \ntype=*%v* \naz=%v \nexpires_at=%v \nlaunched_at=%v \nterminate=%v \nterminated_at=%v \n\n\n",
			lease.InstanceID,
			lease.AWSAccountID,

			lease.InstanceType,
			lease.AvailabilityZone,
			lease.ExpiresAt.Format(time.RFC3339),
			lease.LaunchedAt.Format(time.RFC3339),
			lease.Terminated,
			lease.TerminatedAt.Format(time.RFC3339),
		)
		_, err := response.WriteString(leaseLine + "\n")
		if err != nil {
			return "", err
		}
	}
	return response.String(), nil
}
