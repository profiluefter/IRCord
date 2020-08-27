package irc

type Event interface{}

type MessageReceivedEvent struct {
	Nickname string
	Content  string
}

type EventListener func(Event)
