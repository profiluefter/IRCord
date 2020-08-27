package irc

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

const connectionClosedByRemote = "An existing connection was forcibly closed by the remote host."
const connectionClosedByServer = "use of closed network connection"

type client struct {
	connection net.Conn
	server     *server
	registered bool
	nickname   *string
	//mode     int8
}

func (client *client) handle() {
	defer client.connection.Close()

	for {
		message, recvError := client.recvMessage()
		if recvError != nil {
			if !strings.Contains(recvError.Error(), connectionClosedByRemote) && !strings.Contains(recvError.Error(), connectionClosedByServer) {
				println(recvError.Error())
			}
			return
		}

		handler := commands[strings.ToUpper(message.command)]
		if handler == nil {
			handler = notFoundHandler
		}
		replyError := handler(client, message)

		if replyError != nil {
			if !strings.Contains(recvError.Error(), connectionClosedByRemote) && !strings.Contains(recvError.Error(), connectionClosedByServer) {
				println(replyError.Error())
			}
			return
		}
	}
}

func (client *client) recvMessage() (*message, error) {
	reader := bufio.NewScanner(client.connection)
	reader.Scan()

	if reader.Err() != nil {
		return nil, reader.Err()
	}

	line := reader.Text()
	message, err := parseMessage(line)
	return message, err
}

func (client *client) sendNumeric(numeric reply, parameters ...string) error {
	var target = client.nickname
	if target == nil {
		target = new(string)
		*target = "*"
	}

	var messageParameters = []*string{target}
	for _, parameter := range parameters {
		parameter := parameter
		messageParameters = append(messageParameters, &parameter)
	}
	return client.sendMessage(message{
		prefix:     &client.server.options.Name,
		command:    fmt.Sprintf("%03d", numeric),
		parameters: messageParameters,
	})
}

func (client *client) sendMessage(message message) error {
	_, err := client.connection.Write(append([]byte(message.serialize()), '\r', '\n'))
	return err
}
