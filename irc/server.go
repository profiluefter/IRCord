package irc

import (
	"fmt"
	"net"
	"strconv"
)

type ServerOptions struct {
	Name string
	Port uint16
	Motd *string
}

func NewServer(options ServerOptions) *server {
	server := new(server)
	server.options = options
	return server
}

type server struct {
	options  ServerOptions
	channels []*channel
}

func (server *server) Start() error {
	socket, err := net.Listen("tcp4", ":"+strconv.Itoa(int(server.options.Port)))
	if err != nil {
		return err
	}
	defer socket.Close()

	fmt.Printf("Now listening on %s\n", socket.Addr().String())
	for {
		connection, err := socket.Accept()
		if err != nil {
			return err
		}

		client := client{
			connection: connection,
			server:     server,
			nickname:   nil,
		}

		go client.handle()
	}
}

func (server *server) NewChannel(name string) *channel {
	channel := newChannel(name)
	go channel.worker()

	server.channels = append(server.channels, channel)
	return channel
}
