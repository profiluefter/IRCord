package irc

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

const connectionClosedMessage = "An existing connection was forcibly closed by the remote host."

type client struct {
	connection net.Conn
	server     *Server
	registered bool
	nickname   *string
	//mode     int8
}

func (client *client) handle() {
	defer client.connection.Close()

	for {
		message, recvError := client.recvMessage()
		if recvError != nil {
			if !strings.Contains(recvError.Error(), connectionClosedMessage) {
				println(recvError.Error())
			}
			return
		}

		var replyError error

		//TODO: replace this with a map
		switch message.command {
		case "PASS":
			if client.registered {
				replyError = client.sendNumeric(ERR_ALREADYREGISTRED, "Already registered")
				break
			}
			if message.parameters[0] == nil {
				replyError = client.sendNumeric(ERR_NEEDMOREPARAMS, "No password given")
				break
			}
			//currently there is no auth so all passwords are allowed
			break
		case "NICK":
			if message.parameters[0] == nil {
				replyError = client.sendNumeric(ERR_NONICKNAMEGIVEN, "No nickname given")
				break
			}
			client.nickname = message.parameters[0]
			client.registered = true
			replyError = client.sendNumeric(RPL_WELCOME, fmt.Sprintf("Welcome to the Internet Relay Network %s!%s@%s", *client.nickname, *client.nickname, client.server.Name))
			replyError = client.sendNumeric(RPL_YOURHOST, fmt.Sprintf("Your host is %s, running version git", client.server.Name))
			replyError = client.sendNumeric(RPL_CREATED, "This server was created sometime")
			replyError = client.sendNumeric(RPL_MYINFO, fmt.Sprintf("%s git", client.server.Name))
			replyError = client.sendNumeric(RPL_MOTDSTART, fmt.Sprintf(":- %s Message of the day - ", client.server.Name))
			if replyError != nil {
				break
			}

			for _, line := range strings.Split(*client.server.Motd, "\n") {
				replyError = client.sendNumeric(RPL_MOTD, fmt.Sprintf(":- %s", line))
			}

			replyError = client.sendNumeric(RPL_ENDOFMOTD, ":End of MOTD command")
			break
		case "USER":
			if client.registered {
				replyError = client.sendNumeric(ERR_ALREADYREGISTRED, "Already registered")
				break
			}
			if message.parameters[3] == nil {
				replyError = client.sendNumeric(ERR_NEEDMOREPARAMS, "Not enough parameters")
				break
			}
			replyError = client.sendNumeric(ERR_USERSDISABLED, "Users are not implemented")
			break
		default:
			fmt.Printf("Unknown command: %s\n", message.command)
			break
		}

		if replyError != nil {
			if !strings.Contains(replyError.Error(), connectionClosedMessage) {
				println(replyError.Error())
			}
			return
		}
	}
}

func (client *client) recvMessage() (message, error) {
	reader := bufio.NewScanner(client.connection)
	reader.Scan()

	if reader.Err() != nil {
		return message{}, reader.Err()
	}

	line := reader.Text()
	message := parseMessage(line)
	return message, nil
}

func (client *client) sendNumeric(numeric reply, reason string) error {
	var target = client.nickname
	if target == nil {
		target = new(string)
		*target = "*"
	}

	return client.sendMessage(message{
		prefix:     &client.server.Name,
		command:    fmt.Sprintf("%03d", numeric),
		parameters: [15]*string{target, &reason},
	})
}

func (client *client) sendMessage(message message) error {
	_, err := client.connection.Write(append([]byte(message.serialize()), '\r', '\n'))
	return err
}
