package slackbot

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	"github.com/jinzhu/gorm"
	"github.com/nlopes/slack"
	. "github.com/tleyden/cecil/commrouter"
	"github.com/tleyden/cecil/models"
	"github.com/tleyden/cecil/tasks"
)

type SlackCommandCtx struct {
	Ctx
	Slack *SlackBotInstance
	Msg   *slack.MessageEvent
	RTM   *slack.RTM
}

// NewSlackCommandCtx creates a new context used by the functions that handle the messages
func (si *SlackBotInstance) NewSlackCommandCtx(rtm *slack.RTM, msg *slack.MessageEvent) *SlackCommandCtx {
	return &SlackCommandCtx{
		Slack: si,
		Msg:   msg,
		RTM:   rtm,
	}
}

// HandleMessage handles a command from slack
func (si *SlackBotInstance) HandleMessage(rtm *slack.RTM, msg *slack.MessageEvent) error {
	newCtx := si.NewSlackCommandCtx(rtm, msg)
	return slackCommandRouter.Execute(msg.Text, newCtx)
}

// StartListenAndServer starts listening to incoming commands
func (si *SlackBotInstance) StartListenAndServer() {
	api := slack.New(si.Token)
	api.SetDebug(false)
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

				// TODO: accept messages only from the group specified in SlackConfig

				incomingMessage.Text = strings.Replace(incomingMessage.Text, fmt.Sprintf("<@%v>", botIdentity.User.ID), "", 1)
				rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("Hey <@%v>, you mentioned me!", incomingMessage.User), incomingMessage.Channel))
				err := si.HandleMessage(rtm, incomingMessage)
				if err != nil {
					rtm.SendMessage(rtm.NewOutgoingMessage(
						fmt.Sprintf(
							"<@%v>, an error occured while executing your query:\n %v",
							incomingMessage.User,
							err.Error(),
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

var slackCommandRouter = CommRouter(func() {
	Description("this is the description of the CommRouter")

	Subject(
		"usage",
		func() {
			Description("")
			Command(
				[]string{"show"},
				func() {
					Description("Show list of available commands of the Cecil slack interface")
					Examples("show usage")
					Controller(ShowUsage)
				},
			)
		},
	)

	Subject(
		"leases",
		func() {
			Description("")
			Command(
				[]string{"list", "all", "show"},
				func() {
					Description("List all leases")
					Examples("list leases")
					Controller(ListLeases)
				},
			)
		},
	)

	Subject(
		"lease",
		func() {
			Description("this is the description of the subject")
			Command(
				[]string{"show", "display", "get"},
				func() {
					Description("Show info about a specific lease")
					Examples("show lease instance-id:i-000000", "show lease id:23")
					Controller(ShowLease)

					// TODO: add a way to add Args and their validation conditions
					Params(func() {
						Param("id", Int, func() {
							//Required()
							MinValue(1)
						})
						Param("instance-id", String, func() {
							//Required()
							MinLength(1)
							MustRegex(regexp.MustCompile("i-([a-z0-9]+)"))
						})
					})

				},
			)

			Command(
				[]string{"terminate", "kill", "anihilate"},
				func() {
					Description("Terminate a specific lease")
					Examples("kill lease instance-id:i-000000", "kill lease id:42")
					Controller(TerminateLease)

					// TODO: add a way to add Args and their validation conditions
					Params(func() {
						Param("id", Int, func() {
							//Required()
							MinValue(1)
						})
						Param("instance-id", String, func() {
							//Required()
							MinLength(1)
							MustRegex(regexp.MustCompile("i-([a-z0-9]+)"))
						})
					})

				},
			)

			Command(
				[]string{"extend"},
				func() {
					Description("Extend a specific lease")
					Examples("extend lease instance-id:i-000000", "extend lease id:42")
					Controller(ExtendLease)

					// TODO: add a way to add Args and their validation conditions
					Params(func() {
						Param("id", Int, func() {
							//Required()
							MinValue(1)
						})
						Param("instance-id", String, func() {
							//Required()
							MinLength(1)
							MustRegex(regexp.MustCompile("i-([a-z0-9]+)"))
						})
					})

				},
			)

			Command(
				[]string{"approve"},
				func() {
					Description("Approve a specific lease")
					Examples("approve lease instance-id:i-000000", "approve lease id:42")
					Controller(ApproveLease)

					// TODO: add a way to add Args and their validation conditions
					Params(func() {
						Param("id", Int, func() {
							//Required()
							MinValue(1)
						})
						Param("instance-id", String, func() {
							//Required()
							MinLength(1)
							MustRegex(regexp.MustCompile("i-([a-z0-9]+)"))
						})
					})

				},
			)

		},
	)

})

// ShowUsage returns the usage of the Cecil slack bot
func ShowUsage(rawCtx interface{}) error {
	ctx := rawCtx.(*SlackCommandCtx)

	message := ctx.RouterUsage()
	response := fmt.Sprintf(
		"<@%v>:\n %v",
		ctx.Msg.User,
		message,
	)

	// send response to the channel from which the command has been sent from
	ctx.RTM.SendMessage(
		ctx.RTM.NewOutgoingMessage(
			response,
			ctx.Msg.Channel,
		),
	)
	return nil
}

// ListLeases returns the usage of the Cecil slack interface
func ListLeases(rawCtx interface{}) error {
	ctx := rawCtx.(*SlackCommandCtx)

	// fetch leases for account
	leases, err := ctx.Slack.s.LeasesForAccount(int(ctx.Slack.AccountID), nil)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("no leases found")
		}
		return errors.New("internal error")
	}

	var response bytes.Buffer

	responseHeader := fmt.Sprintf(
		"<@%v>:\n",
		ctx.Msg.User,
	)
	response.WriteString(responseHeader)

	for leaseIndex := range leases {
		lease := leases[leaseIndex]
		var terminated string
		if lease.TerminatedAt != nil {
			terminated = fmt.Sprint("true :skull:")
		} else {
			terminated = fmt.Sprint("false :green_heart:")
		}

		leaseHeader := fmt.Sprintf("*LEASE* on AWS %v:",
			lease.AWSAccountID,
		)
		response.WriteString(leaseHeader)

		line := fmt.Sprintf("\nlease_id=*%v* \ngroup_type=%v \ngroup_uid=*%v*",
			lease.ID,
			lease.GroupType.String(),
			lease.GroupUID,
		)
		response.WriteString(line + "\n")

		leaseInfo := fmt.Sprintf(
			"expires_at=%v \nlaunched_at=%v \nterminated=%v\n",

			lease.ExpiresAt.Format(time.RFC3339),
			lease.LaunchedAt.Format(time.RFC3339),
			terminated,
		)
		if lease.TerminatedAt != nil {
			terminatedAt := fmt.Sprintf(
				"terminated_at=%v ",
				lease.TerminatedAt.Format(time.RFC3339),
			)
			response.WriteString(terminatedAt)
		}
		response.WriteString(leaseInfo)
		response.WriteString("-------------------\n\n")

		if response.Len() >= 2000 {
			// flush buffer; slack has a limit for message length
			ctx.RTM.SendMessage(
				ctx.RTM.NewOutgoingMessage(
					response.String(),
					ctx.Msg.Channel,
				),
			)

			response.Reset()
		}
	}

	if response.Len() > 0 {
		// if anything is left, send it
		ctx.RTM.SendMessage(
			ctx.RTM.NewOutgoingMessage(
				response.String(),
				ctx.Msg.Channel,
			),
		)
	}

	return nil
}

// ShowLease returns info about a specific lease
func ShowLease(rawCtx interface{}) error {
	ctx := rawCtx.(*SlackCommandCtx)

	instanceID, err1 := ctx.Params().GetString("instance-id")
	leaseID, err2 := ctx.Params().GetInt("id")

	notEnoughInfo := err1 != nil && err2 != nil
	if notEnoughInfo {
		return errors.New("no way to select a lease")
	}

	var err error
	var lease *models.Lease

	// fetch lease
	if instanceID != "" {
		lease, err = ctx.Slack.s.LeaseByInstanceID(ctx.Slack.AccountID, nil, instanceID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.New("lease not found")
			}
			return errors.New("internal error")
		}
	} else {
		lease, err = ctx.Slack.s.GetLeaseByID(leaseID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.New("lease not found")
			}
			return errors.New("internal error")
		}
	}
	if lease.AccountID != ctx.Slack.AccountID {
		return errors.New("lease not found")
	}

	var response bytes.Buffer

	responseHeader := fmt.Sprintf(
		"<@%v>:\n",
		ctx.Msg.User,
	)
	response.WriteString(responseHeader)

	leaseBytes, err := yaml.Marshal(lease)
	if err != nil {
		return err
	}

	response.Write(leaseBytes)

	if response.Len() > 0 {
		// if anything is left, send it
		ctx.RTM.SendMessage(
			ctx.RTM.NewOutgoingMessage(
				response.String(),
				ctx.Msg.Channel,
			),
		)
	}

	return nil
}

// TerminateLease terminates a lease
func TerminateLease(rawCtx interface{}) error {
	ctx := rawCtx.(*SlackCommandCtx)

	instanceID, err1 := ctx.Params().GetString("instance-id")
	leaseID, err2 := ctx.Params().GetInt("id")

	notEnoughInfo := err1 != nil && err2 != nil
	if notEnoughInfo {
		return errors.New("no way to select a lease")
	}

	var err error
	var lease *models.Lease

	// fetch lease
	if instanceID != "" {
		lease, err = ctx.Slack.s.LeaseByInstanceID(ctx.Slack.AccountID, nil, instanceID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.New("lease not found")
			}
			return errors.New("internal error")
		}
	} else {
		lease, err = ctx.Slack.s.GetLeaseByID(leaseID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.New("lease not found")
			}
			return errors.New("internal error")
		}
	}
	if lease.AccountID != ctx.Slack.AccountID {
		return errors.New("lease not found")
	}

	var response bytes.Buffer

	responseHeader := fmt.Sprintf(
		"<@%v>:\n",
		ctx.Msg.User,
	)
	response.WriteString(responseHeader)

	ctx.Slack.s.Queues().TerminatorQueue().PushTask(tasks.TerminatorTask{Lease: *lease})

	response.WriteString(fmt.Sprintf("termination of lease %v initiated", lease.ID))

	if response.Len() > 0 {
		// if anything is left, send it
		ctx.RTM.SendMessage(
			ctx.RTM.NewOutgoingMessage(
				response.String(),
				ctx.Msg.Channel,
			),
		)
	}

	return nil
}

// ExtendLease extends a lease
func ExtendLease(rawCtx interface{}) error {
	ctx := rawCtx.(*SlackCommandCtx)

	instanceID, err1 := ctx.Params().GetString("instance-id")
	leaseID, err2 := ctx.Params().GetInt("id")

	notEnoughInfo := err1 != nil && err2 != nil
	if notEnoughInfo {
		return errors.New("no way to select a lease")
	}

	var err error
	var lease *models.Lease

	// fetch lease
	if instanceID != "" {
		lease, err = ctx.Slack.s.LeaseByInstanceID(ctx.Slack.AccountID, nil, instanceID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.New("lease not found")
			}
			return errors.New("internal error")
		}
	} else {
		lease, err = ctx.Slack.s.GetLeaseByID(leaseID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.New("lease not found")
			}
			return errors.New("internal error")
		}
	}
	if lease.AccountID != ctx.Slack.AccountID {
		return errors.New("lease not found")
	}

	var response bytes.Buffer

	responseHeader := fmt.Sprintf(
		"<@%v>:\n",
		ctx.Msg.User,
	)
	response.WriteString(responseHeader)

	ctx.Slack.s.Queues().ExtenderQueue().PushTask(tasks.ExtenderTask{
		Lease:     *lease,
		Approving: false,
	})

	response.WriteString(fmt.Sprintf("extension of lease %v initiated", lease.ID))

	if response.Len() > 0 {
		// if anything is left, send it
		ctx.RTM.SendMessage(
			ctx.RTM.NewOutgoingMessage(
				response.String(),
				ctx.Msg.Channel,
			),
		)
	}

	return nil
}

// ApproveLease approves a lease
func ApproveLease(rawCtx interface{}) error {
	ctx := rawCtx.(*SlackCommandCtx)

	instanceID, err1 := ctx.Params().GetString("instance-id")
	leaseID, err2 := ctx.Params().GetInt("id")

	notEnoughInfo := err1 != nil && err2 != nil
	if notEnoughInfo {
		return errors.New("no way to select a lease")
	}

	var err error
	var lease *models.Lease

	// fetch lease
	if instanceID != "" {
		lease, err = ctx.Slack.s.LeaseByInstanceID(ctx.Slack.AccountID, nil, instanceID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.New("lease not found")
			}
			return errors.New("internal error")
		}
	} else {
		lease, err = ctx.Slack.s.GetLeaseByID(leaseID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.New("lease not found")
			}
			return errors.New("internal error")
		}
	}
	if lease.AccountID != ctx.Slack.AccountID {
		return errors.New("lease not found")
	}

	var response bytes.Buffer

	responseHeader := fmt.Sprintf(
		"<@%v>:\n",
		ctx.Msg.User,
	)
	response.WriteString(responseHeader)

	ctx.Slack.s.Queues().ExtenderQueue().PushTask(tasks.ExtenderTask{
		Lease:     *lease,
		Approving: true,
	})

	response.WriteString(fmt.Sprintf("approval of lease %v initiated", lease.ID))

	if response.Len() > 0 {
		// if anything is left, send it
		ctx.RTM.SendMessage(
			ctx.RTM.NewOutgoingMessage(
				response.String(),
				ctx.Msg.Channel,
			),
		)
	}

	return nil
}
