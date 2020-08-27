package irc

import (
	"fmt"
	"strconv"
	"strings"
)

type messageHandler func(*client, *message) error

var notFoundHandler messageHandler = func(c *client, m *message) error {
	fmt.Printf("Unknown command: %s\n", m.command)
	return c.sendNumeric(ERR_UNKNOWNCOMMAND, m.command, "Unknown command")
}

var commands = map[string]messageHandler{
	"PASS": func(c *client, m *message) error {
		if c.registered {
			return c.sendNumeric(ERR_ALREADYREGISTRED, "Already registered")
		}
		if len(m.parameters) < 1 {
			return c.sendNumeric(ERR_NEEDMOREPARAMS, "No password given")
		}
		//currently there is no auth so all passwords are allowed
		return nil
	},
	"NICK": func(c *client, m *message) error {
		if len(m.parameters) < 1 {
			return c.sendNumeric(ERR_NONICKNAMEGIVEN, "No nickname given")
		}
		c.nickname = m.parameters[0]
		c.registered = true

		var errors [5]error
		errors[0] = c.sendNumeric(RPL_WELCOME, fmt.Sprintf("Welcome to the Internet Relay Network %s!%s@%s", *c.nickname, *c.nickname, c.server.options.Name))
		errors[1] = c.sendNumeric(RPL_YOURHOST, fmt.Sprintf("Your host is %s, running version git", c.server.options.Name))
		errors[2] = c.sendNumeric(RPL_CREATED, "This server was created sometime")
		errors[3] = c.sendNumeric(RPL_MYINFO, fmt.Sprintf("%s git", c.server.options.Name))
		errors[4] = c.sendNumeric(RPL_MOTDSTART, fmt.Sprintf("- %s Message of the day - ", c.server.options.Name))

		for _, err := range errors {
			if err != nil {
				return err
			}
		}

		for _, line := range strings.Split(*c.server.options.Motd, "\n") {
			err := c.sendNumeric(RPL_MOTD, fmt.Sprintf("- %s", line))
			if err != nil {
				return err
			}
		}

		return c.sendNumeric(RPL_ENDOFMOTD, "End of MOTD command")
	},
	"USER": func(c *client, m *message) error {
		if c.registered {
			return c.sendNumeric(ERR_ALREADYREGISTRED, "Already registered")
		}
		if len(m.parameters) < 4 {
			return c.sendNumeric(ERR_NEEDMOREPARAMS, "Not enough parameters")
		}
		return c.sendNumeric(ERR_USERSDISABLED, "Users are not implemented")
	},
	"QUIT": func(c *client, m *message) error {
		c.connection.Close()
		return nil
	},
	"PING": func(c *client, m *message) error {
		if len(m.parameters) < 1 {
			return c.sendNumeric(ERR_NOORIGIN, "Not enough parameters")
		}
		return c.sendMessage(message{
			prefix:  &c.server.options.Name,
			command: "PONG",
			parameters: []*string{
				&c.server.options.Name,
				m.parameters[0],
			},
		})
	},
	"JOIN": func(c *client, m *message) error {
		if len(m.parameters) < 1 {
			return c.sendNumeric(ERR_NEEDMOREPARAMS, "No channels given")
		}
		if *m.parameters[0] == "0" {
			//TODO: Unsubscribe from all
			return nil
		}

		channelList := strings.Split(*m.parameters[0], ",")
		for _, name := range channelList {
			channel := c.server.channels[name]
			if channel != nil {
				err := channel.join(c)
				if err != nil {
					return err
				}
			} else {
				err := c.sendNumeric(ERR_NOSUCHCHANNEL, name, "No such channel")
				if err != nil {
					return err
				}
			}
		}
		return nil
	},
	"PRIVMSG": func(c *client, m *message) error {
		if len(m.parameters) == 0 {
			return c.sendNumeric(ERR_NORECIPIENT, "No recipient given (PRIVMSG)")
		}
		if len(m.parameters) == 1 {
			return c.sendNumeric(ERR_NOTEXTTOSEND, "No text to send")
		}
		channel := c.server.channels[*m.parameters[0]]
		if channel == nil {
			return c.sendNumeric(ERR_NOSUCHNICK, *m.parameters[0], "No such nick/channel")
		}
		channel.clientSentMessage(*c.nickname, *m.parameters[1])
		return nil
	},
	"LIST": func(c *client, m *message) error {
		if len(m.parameters) == 2 && *m.parameters[1] != c.server.options.Name {
			return c.sendNumeric(ERR_NOSUCHSERVER, *m.parameters[1], "No such server")
		}

		//TODO: Replace hardcoded +100 with a variable that can be set by the library user
		//      This is currently necessary because channels with 0 users don't show up in most IRC clients
		if len(m.parameters) == 1 {
			requestedChannels := strings.Split(*m.parameters[0], ",")
			for _, requestedChannel := range requestedChannels {
				channel := c.server.channels[requestedChannel]
				if channel != nil {
					topic := channel.topic
					if topic == "" {
						topic = "No topic is set"
					}
					err := c.sendNumeric(RPL_LIST, requestedChannel, strconv.Itoa(len(channel.subscriber)+100), topic)
					if err != nil {
						return err
					}
				} else {
					err := c.sendNumeric(ERR_NOSUCHNICK, requestedChannel, "No such nick/channel")
					if err != nil {
						return err
					}
				}
			}
		} else {
			for _, channel := range c.server.channels {
				topic := channel.topic
				if topic == "" {
					topic = "No topic is set"
				}
				err := c.sendNumeric(RPL_LIST, channel.Name, strconv.Itoa(len(channel.subscriber)+100), topic)
				if err != nil {
					return err
				}
			}
		}
		return c.sendNumeric(RPL_LISTEND, "End of LIST")
	},
}
