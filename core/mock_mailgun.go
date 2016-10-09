package core

import (
	"fmt"
	"reflect"

	mailgun "gopkg.in/mailgun/mailgun-go.v1"
)

type MockMailGun struct {
	SentMessages chan *mailgun.Message

	// Embed the Mailgun interface. No idea what will happen if unimplemented methods are called.
	mailgun.Mailgun
}

func NewMockMailGun() *MockMailGun {

	mockMailGun := MockMailGun{
		SentMessages: make(chan *mailgun.Message, 100),
	}

	return &mockMailGun
}

func (mmg *MockMailGun) Send(m *mailgun.Message) (string, string, error) {
	defer func() {
		mmg.SentMessages <- m
	}()
	return "", "", nil
}

func (mmg *MockMailGun) waitForNotification(nt NotificationType) {

	message := <-mmg.SentMessages
	messageType, err := mmg.getHeaderViaReflection(message, "X-ZeroCloud-MessageType")
	if err != nil {
		panic(fmt.Sprintf("Error getting header value from mock mailgun msg: %v", err))
	}

	messageNotificationType := NotificationTypeFromString(messageType)
	switch messageNotificationType {
	case nt:
		return
	default:
		panic(
			fmt.Sprintf(
				"Unexpected notification type.  Expected: %s, got: %s",
				nt,
				messageNotificationType,
			),
		)
	}

}

func (mmg *MockMailGun) getHeaderViaReflection(message *mailgun.Message, headerKey string) (string, error) {

	val := reflect.ValueOf(*message)
	headers := val.FieldByName("headers")

	switch headers.Kind() {
	case reflect.Map:
		for _, key := range headers.MapKeys() {
			mapVal := headers.MapIndex(key)
			switch mapVal.Kind() {
			case reflect.String:
				return fmt.Sprintf("%s", mapVal), nil
			default:
				return "error", fmt.Errorf("Unexpected type for value at key: %s", headerKey)
			}
		}

	default:
		return "error", fmt.Errorf("Unexpected type for mailgun message headers field")
	}

	return "error", fmt.Errorf("Did not find headerKey")

}
