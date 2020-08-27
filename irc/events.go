package irc

type event interface{}

type messageReceivedEvent struct {
	nickname string
	content  string
}

type EventListener func(event)
