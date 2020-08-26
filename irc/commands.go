package irc

import (
	"fmt"
	"strings"
)

type messageHandler func(*client, message) error

var notFoundHandler messageHandler = func(client *client, message message) error {
	fmt.Printf("Unknown command: %s\n", message.command)
	return nil
}

var commands = map[string]messageHandler{
	"PASS": func(client *client, message message) error {
		if client.registered {
			return client.sendNumeric(ERR_ALREADYREGISTRED, "Already registered")
		}
		if message.parameters[0] == nil {
			return client.sendNumeric(ERR_NEEDMOREPARAMS, "No password given")
		}
		//currently there is no auth so all passwords are allowed
		return nil
	},
	"NICK": func(client *client, message message) error {
		if message.parameters[0] == nil {
			return client.sendNumeric(ERR_NONICKNAMEGIVEN, "No nickname given")
		}
		client.nickname = message.parameters[0]
		client.registered = true

		var errors [5]error
		errors[0] = client.sendNumeric(RPL_WELCOME, fmt.Sprintf("Welcome to the Internet Relay Network %s!%s@%s", *client.nickname, *client.nickname, client.server.Name))
		errors[1] = client.sendNumeric(RPL_YOURHOST, fmt.Sprintf("Your host is %s, running version git", client.server.Name))
		errors[2] = client.sendNumeric(RPL_CREATED, "This server was created sometime")
		errors[3] = client.sendNumeric(RPL_MYINFO, fmt.Sprintf("%s git", client.server.Name))
		errors[4] = client.sendNumeric(RPL_MOTDSTART, fmt.Sprintf(":- %s Message of the day - ", client.server.Name))

		for _, err := range errors {
			if err != nil {
				return err
			}
		}

		for _, line := range strings.Split(*client.server.Motd, "\n") {
			err := client.sendNumeric(RPL_MOTD, fmt.Sprintf(":- %s", line))
			if err != nil {
				return err
			}
		}

		return client.sendNumeric(RPL_ENDOFMOTD, ":End of MOTD command")
	},
	"USER": func(client *client, message message) error {
		if client.registered {
			return client.sendNumeric(ERR_ALREADYREGISTRED, "Already registered")
		}
		if message.parameters[3] == nil {
			return client.sendNumeric(ERR_NEEDMOREPARAMS, "Not enough parameters")
		}
		return client.sendNumeric(ERR_USERSDISABLED, "Users are not implemented")
	},
	"QUIT": func(client *client, message message) error {
		client.connection.Close()
		return nil
	},
}
