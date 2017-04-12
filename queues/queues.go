package queues

import "github.com/gagliardetto/simpleQueue"

type QueuesGroup struct {
	newInstanceQueue        *simpleQueue.Queue
	terminatorQueue         *simpleQueue.Queue
	instanceTerminatedQueue *simpleQueue.Queue
	extenderQueue           *simpleQueue.Queue
	notifierQueue           *simpleQueue.Queue
}

type QueuesGroupInterface interface {
	NewInstanceQueue() *simpleQueue.Queue
	TerminatorQueue() *simpleQueue.Queue
	InstanceTerminatedQueue() *simpleQueue.Queue
	ExtenderQueue() *simpleQueue.Queue
	NotifierQueue() *simpleQueue.Queue

	SetNewInstanceQueue(q *simpleQueue.Queue)
	SetTerminatorQueue(q *simpleQueue.Queue)
	SetInstanceTerminatedQueue(q *simpleQueue.Queue)
	SetExtenderQueue(q *simpleQueue.Queue)
	SetNotifierQueue(q *simpleQueue.Queue)
}

func (qg *QueuesGroup) NewInstanceQueue() *simpleQueue.Queue {
	return qg.newInstanceQueue
}
func (qg *QueuesGroup) TerminatorQueue() *simpleQueue.Queue {
	return qg.terminatorQueue
}
func (qg *QueuesGroup) InstanceTerminatedQueue() *simpleQueue.Queue {
	return qg.instanceTerminatedQueue
}
func (qg *QueuesGroup) ExtenderQueue() *simpleQueue.Queue {
	return qg.extenderQueue
}
func (qg *QueuesGroup) NotifierQueue() *simpleQueue.Queue {
	return qg.notifierQueue
}

func (qg *QueuesGroup) SetNewInstanceQueue(q *simpleQueue.Queue) {
	qg.newInstanceQueue = q
}
func (qg *QueuesGroup) SetTerminatorQueue(q *simpleQueue.Queue) {
	qg.terminatorQueue = q
}
func (qg *QueuesGroup) SetInstanceTerminatedQueue(q *simpleQueue.Queue) {
	qg.instanceTerminatedQueue = q
}
func (qg *QueuesGroup) SetExtenderQueue(q *simpleQueue.Queue) {
	qg.extenderQueue = q
}
func (qg *QueuesGroup) SetNotifierQueue(q *simpleQueue.Queue) {
	qg.notifierQueue = q
}
