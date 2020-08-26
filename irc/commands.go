package irc

import (
	"fmt"
	"strings"
)

type messageHandler func(*client, *message) error

var notFoundHandler messageHandler = func(c *client, m *message) error {
	fmt.Printf("Unknown command: %s\n", m.command)
	return nil
}

var commands = map[string]messageHandler{
	"PASS": func(c *client, m *message) error {
		if c.registered {
			return c.sendNumeric(ERR_ALREADYREGISTRED, "Already registered")
		}
		if m.parameters[0] == nil {
			return c.sendNumeric(ERR_NEEDMOREPARAMS, "No password given")
		}
		//currently there is no auth so all passwords are allowed
		return nil
	},
	"NICK": func(c *client, m *message) error {
		if m.parameters[0] == nil {
			return c.sendNumeric(ERR_NONICKNAMEGIVEN, "No nickname given")
		}
		c.nickname = m.parameters[0]
		c.registered = true

		var errors [5]error
		errors[0] = c.sendNumeric(RPL_WELCOME, fmt.Sprintf("Welcome to the Internet Relay Network %s!%s@%s", *c.nickname, *c.nickname, c.server.options.Name))
		errors[1] = c.sendNumeric(RPL_YOURHOST, fmt.Sprintf("Your host is %s, running version git", c.server.options.Name))
		errors[2] = c.sendNumeric(RPL_CREATED, "This server was created sometime")
		errors[3] = c.sendNumeric(RPL_MYINFO, fmt.Sprintf("%s git", c.server.options.Name))
		errors[4] = c.sendNumeric(RPL_MOTDSTART, fmt.Sprintf(":- %s Message of the day - ", c.server.options.Name))

		for _, err := range errors {
			if err != nil {
				return err
			}
		}

		for _, line := range strings.Split(*c.server.options.Motd, "\n") {
			err := c.sendNumeric(RPL_MOTD, fmt.Sprintf(":- %s", line))
			if err != nil {
				return err
			}
		}

		return c.sendNumeric(RPL_ENDOFMOTD, ":End of MOTD command")
	},
	"USER": func(c *client, m *message) error {
		if c.registered {
			return c.sendNumeric(ERR_ALREADYREGISTRED, "Already registered")
		}
		if m.parameters[3] == nil {
			return c.sendNumeric(ERR_NEEDMOREPARAMS, "Not enough parameters")
		}
		return c.sendNumeric(ERR_USERSDISABLED, "Users are not implemented")
	},
	"QUIT": func(c *client, m *message) error {
		c.connection.Close()
		return nil
	},
}
