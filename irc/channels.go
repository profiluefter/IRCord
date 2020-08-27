package irc

//never construct this yourself
type Channel struct {
	Name           string
	topic          string
	listeners      []*EventListener
	incomingEvents chan Event
	subscriber     []*client
}

func newChannel(name string, topic string) *Channel {
	c := &Channel{
		Name:           name,
		topic:          topic,
		listeners:      []*EventListener{},
		incomingEvents: make(chan Event),
		subscriber:     []*client{},
	}
	c.AddListener(createLoopbackListener(c))

	go c.worker()
	return c
}

func (channel *Channel) SendMessage(sender string, content string) {
	channel.broadcastMessage(&sender, message{
		prefix:  &sender,
		command: "PRIVMSG",
		parameters: []*string{
			&channel.Name,
			&content,
		},
	})
}

func (channel *Channel) broadcastMessage(sender *string, message message) {
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

func (channel *Channel) AddListener(listener *EventListener) {
	channel.listeners = append(channel.listeners, listener)
}

func (channel *Channel) sendEvent(event Event) {
	channel.incomingEvents <- event
}

func (channel *Channel) worker() {
	for event := range channel.incomingEvents {
		for _, listener := range channel.listeners {
			go (*listener)(event)
		}
	}
}

func (channel *Channel) join(c *client) error {
	channel.subscriber = append(channel.subscriber, c)

	channel.broadcastMessage(nil, message{
		prefix:     c.nickname,
		command:    "JOIN",
		parameters: []*string{&channel.Name},
	})
	//TODO: Additional replies

	return nil
}

func (channel *Channel) part(c *client) error {
	channel.subscriber = removeHolyShitWhyIsThisNotABuiltinLikeAppendAndWhyIsItSoUgly(channel.subscriber, c)

	channel.broadcastMessage(nil, message{
		prefix:     c.nickname,
		command:    "PART",
		parameters: []*string{&channel.Name},
	})

	return nil
}

func (channel *Channel) clientSentMessage(nickname string, content string) {
	channel.sendEvent(MessageReceivedEvent{
		Nickname: nickname,
		Content:  content,
	})
}

func createLoopbackListener(channel *Channel) *EventListener {
	f := func(e Event) {
		event, ok := e.(MessageReceivedEvent)
		if !ok {
			return
		}
		channel.SendMessage(event.Nickname, event.Content)
	}
	return (*EventListener)(&f)
}

func removeHolyShitWhyIsThisNotABuiltinLikeAppendAndWhyIsItSoUgly(slice []*client, toRemove *client) []*client {
	var index = -1
	for i, element := range slice {
		if element == toRemove {
			index = i
		}
	}
	if index == -1 {
		panic("couldn't find element")
	}

	slice[index] = slice[len(slice)-1]
	return slice[:len(slice)-1]
}
