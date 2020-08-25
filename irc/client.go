package irc

import (
	"bufio"
	"net"
	"strings"
)

const connectionClosedMessage = "An existing connection was forcibly closed by the remote host."

type client struct {
	connection net.Conn
	server     *Server
	username   *string
}

func (client *client) handle() {
	defer client.connection.Close()

	for {
		message, err := client.recvMessage()
		if err != nil {
			if strings.Contains(err.Error(), connectionClosedMessage) {
				return
			}
			println(err.Error())
			return
		}
		println(message.serialize())
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
