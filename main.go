package main

import (
	"fmt"
	"github.com/profiluefter/IRCord/irc"
)

func main() {
	server := irc.Server{
		Name: "irc-cord",
		Port: 6667,
	}

	err := server.Start()
	if err != nil {
		fmt.Println(err)
	}
}
