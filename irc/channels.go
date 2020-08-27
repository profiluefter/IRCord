package irc

type channel struct {
	name           string
	listeners      []*EventListener
	incomingEvents chan Event
	subscriber     []*client
}

func newChannel(name string) *channel {
	c := new(channel)
	c.name = name
	c.listeners = []*EventListener{createLoopbackListener(c)}
	c.incomingEvents = make(chan Event)
	c.subscriber = []*client{}

	go c.worker()
	return c
}

func (channel *channel) SendMessage(sender string, content string) {
	channel.broadcastMessage(&sender, message{
		prefix:  &sender,
		command: "PRIVMSG",
		parameters: []*string{
			&channel.name,
			&content,
		},
	})
}

func (channel *channel) broadcastMessage(sender *string, message message) {
	for _, subscriber := range channel.subscriber {
		if sender != nil && *subscriber.nickname == *sender {
			continue
		}
		subscriber := subscriber
		go func() {
			_ = subscriber.sendMessage(message)
		}()
	}
}

func (channel *channel) AddListener(listener *EventListener) {
	channel.listeners = append(channel.listeners, listener)
}

func (channel *channel) sendEvent(event Event) {
	channel.incomingEvents <- event
}

func (channel *channel) worker() {
	for event := range channel.incomingEvents {
		for _, listener := range channel.listeners {
			go (*listener)(event)
		}
	}
}

func (channel *channel) join(c *client) error {
	channel.subscriber = append(channel.subscriber, c)

	channel.broadcastMessage(nil, message{
		prefix:     c.nickname,
		command:    "JOIN",
		parameters: []*string{&channel.name},
	})
	//TODO: Additional replies

	return nil
}

func (channel *channel) clientSentMessage(nickname string, content string) {
	channel.sendEvent(MessageReceivedEvent{
		Nickname: nickname,
		Content:  content,
	})
}

func createLoopbackListener(channel *channel) *EventListener {
	f := func(e Event) {
		event, ok := e.(MessageReceivedEvent)
		if !ok {
			return
		}
		channel.SendMessage(event.Nickname, event.Content)
	}
	return (*EventListener)(&f)
}
