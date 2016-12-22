package core

import (
	"fmt"
	"reflect"

	mailgun "gopkg.in/mailgun/mailgun-go.v1"
)

type MockMailGun struct {
	SentMessages chan *mailgun.Message

	// Embed the Mailgun interface. If unimplemented methods are called, it will panic
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

func (mmg *MockMailGun) WaitForNotification(nt NotificationType) NotificationMeta {

	notificationMeta := NotificationMeta{}
	message := <-mmg.SentMessages

	Logger.Info("waitForNotification", "message", message)

	// Notification Type
	messageType, err := mmg.getHeaderViaReflection(message, X_CECIL_MESSAGETYPE)
	if err != nil {
		panic(fmt.Sprintf("Error getting header value from mock mailgun msg: %v", err))
	}
	messageNotificationType := NotificationTypeFromString(messageType)
	notificationMeta.NotificationType = messageNotificationType

	// Lease UUID
	leaseUUID, err := mmg.getHeaderViaReflection(message, X_CECIL_LEASE_UUID)
	if err != nil {
		panic(fmt.Sprintf("Error getting header value from mock mailgun msg: %v", err))
	}
	notificationMeta.LeaseUuid = leaseUUID

	// Instance ID
	instanceID, err := mmg.getHeaderViaReflection(message, X_CECIL_INSTANCE_ID)
	if err != nil {
		panic(fmt.Sprintf("Error getting header value from mock mailgun msg: %v", err))
	}
	notificationMeta.InstanceId = instanceID

	// Verification Token
	verificationToken, err := mmg.getHeaderViaReflection(message, X_CECIL_VERIFICATION_TOKEN)
	if err != nil {
		panic(fmt.Sprintf("Error getting header value from mock mailgun msg: %v", err))
	}
	notificationMeta.VerificationToken = verificationToken

	switch messageNotificationType {
	case nt:
		return notificationMeta
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
			keyString := fmt.Sprintf("%s", key)
			if keyString != headerKey {
				continue
			}
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

func (mmg *MockMailGun) ValidateEmail(email string) (mailgun.EmailVerification, error) {
	return mailgun.EmailVerification{
		IsValid: true,
	}, nil
}
