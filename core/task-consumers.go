package core

import (
	"fmt"
	"sync/atomic"

	"gopkg.in/mailgun/mailgun-go.v1"
)

// @@@@@@@@@@@@@@@ Task consumers @@@@@@@@@@@@@@@

func (s *Service) NewLeaseQueueConsumer(t interface{}) error {

	return nil
}

func (s *Service) TerminatorQueueConsumer(t interface{}) error {
	if t == nil {
		return fmt.Errorf("%v", "t is nil")
	}
	task := t.(TerminatorTask)
	// TODO: check whether fields are non-null and valid

	_ = task

	atomic.AddInt64(&s.counter, 1)
	fmt.Println(s.counter)
	return nil
}

func (s *Service) LeaseTerminatedQueueConsumer(t interface{}) error {

	return nil
}

func (s *Service) RenewerQueueConsumer(t interface{}) error {
	if t == nil {
		return fmt.Errorf("%v", "t is nil")
	}
	task := t.(RenewerTask)
	// TODO: check whether fields are non-null and valid

	_ = task

	atomic.AddInt64(&s.counter, 1)
	fmt.Println(s.counter)

	return nil
}

func (s *Service) NotifierQueueConsumer(t interface{}) error {
	if t == nil {
		return fmt.Errorf("%v", "t is nil")
	}
	task := t.(NotifierTask)
	// TODO: check whether fields are non-null and valid

	message := mailgun.NewMessage(
		task.From,
		task.Subject,
		task.BodyText,
		task.To,
	)

	message.SetTracking(true)
	//message.SetDeliveryTime(time.Now().Add(24 * time.Hour))
	message.SetHtml(task.BodyHTML)
	_, id, err := s.Mailer.Send(message)
	if err != nil {
		logger.Error("error while sending email", err)
		return err
	}
	_ = id

	return nil
}
