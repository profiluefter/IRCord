package irc

type channel struct {
	name           string
	listeners      []*EventListener
	incomingEvents chan event
	subscriber     []*client
}

func newChannel(name string) *channel {
	c := new(channel)
	c.name = name
	c.listeners = []*EventListener{createLoopbackListener(c)}
	c.incomingEvents = make(chan event)
	c.subscriber = []*client{}

	go c.worker()
	return c
}

func (channel *channel) SendMessage(sender string, content string) {
	for _, subscriber := range channel.subscriber {
		if *subscriber.nickname == sender {
			continue
		}
		subscriber := subscriber
		go func() {
			_ = subscriber.sendMessage(message{
				prefix:  &sender,
				command: "PRIVMSG",
				parameters: []*string{
					&channel.name,
					&content,
				},
			})
		}()
	}
}

func (channel *channel) AddListener(listener *EventListener) {
	channel.listeners = append(channel.listeners, listener)
}

func (channel *channel) sendEvent(event event) {
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
	err := c.sendMessage(message{
		prefix:     c.nickname,
		command:    "JOIN",
		parameters: []*string{&channel.name},
	})
	if err != nil {
		return err
	}

	channel.subscriber = append(channel.subscriber, c)
	return nil
}

func (channel *channel) clientSentMessage(nickname string, content string) {
	channel.sendEvent(messageReceivedEvent{
		nickname: nickname,
		content:  content,
	})
}

func createLoopbackListener(channel *channel) *EventListener {
	f := func(e event) {
		event, ok := e.(messageReceivedEvent)
		if !ok {
			return
		}
		channel.SendMessage(event.nickname, event.content)
	}
	return (*EventListener)(&f)
}
