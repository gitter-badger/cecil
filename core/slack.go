package core

// TODO:
// decide whether to use email, slack, or other (configuration flag ???, setup option???)
// setup client
// listen to messages in notifier queue

// send notification: push to NotifierQueue.
// NotifierQueue decides what to do;

// receive command: can be email_action endpoints, or slack websocket.
// each action is verified and then passed to the right queue (renew, terminate, etc.)
