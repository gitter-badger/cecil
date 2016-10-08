package core

import mailgun "gopkg.in/mailgun/mailgun-go.v1"

type MockMailGun struct {
	MailgunInvocations chan interface{}

	// Embed the Mailgun interface. No idea what will happen if unimplemented methods are called.
	mailgun.Mailgun
}

func NewMockMailGun() *MockMailGun {

	mockMailGun := MockMailGun{
		MailgunInvocations: make(chan interface{}, 100),
	}

	return &mockMailGun
}

func (mmg *MockMailGun) Send(m *mailgun.Message) (string, string, error) {
	defer func() {
		mmg.MailgunInvocations <- m
	}()
	return "", "", nil
}
