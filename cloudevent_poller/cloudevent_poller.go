package cloudevent_poller

import "log"

type CloudEventPoller struct {
	SQSQueueTopicARN string
	ZeroCloudAPIURL  string
}

func (p *CloudEventPoller) Run() {
	log.Printf("Run() called.  Poller: %+v", p)
	return
}
