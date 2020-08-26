package irc

type eventType int

const (
	messageReceived eventType = iota
)

type event interface {
	eventType() eventType
}

type messageReceivedEvent struct {
	content string
}

type EventListener func(event)
