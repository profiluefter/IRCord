package irc

import (
	"net"
	"strconv"
)

type Server struct {
	Name string
	Port uint16
}

func (server *Server) Start() error {
	socket, err := net.Listen("tcp4", ":"+strconv.Itoa(int(server.Port)))
	if err != nil {
		return err
	}
	defer socket.Close()

	for {
		connection, err := socket.Accept()
		if err != nil {
			return err
		}

		client := client{
			connection: connection,
			server:     server,
			username:   nil,
		}

		go client.handle()
	}
}
