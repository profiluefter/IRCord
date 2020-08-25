package main

import (
	"fmt"
	"github.com/profiluefter/IRCord/irc"
)

func main() {
	var motd = "This is the message of the day!\nIf you can see this then the server did not crash yet\nNice."

	server := irc.Server{
		Name: "irc-cord",
		Port: 6667,
		Motd: &motd,
	}

	err := server.Start()
	if err != nil {
		fmt.Println(err)
	}
}
