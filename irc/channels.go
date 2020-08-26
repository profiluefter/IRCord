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
	return c
}

func (channel *channel) SendMessage(sender string, content string) {
	prefix := ":" + sender
	for _, subscriber := range channel.subscriber {
		subscriber := subscriber
		go func() {
			_ = subscriber.sendMessage(message{
				prefix:  &prefix,
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
