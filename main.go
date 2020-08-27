package main

import (
	"fmt"
	"github.com/profiluefter/IRCord/irc"
)

func main() {
	var motd = "This is the message of the day!\nIf you can see this then the server did not crash yet\nNice."

	options := irc.ServerOptions{
		Name: "irc-cord",
		Port: 6667,
		Motd: &motd,
	}
	server := irc.NewServer(options)

	channel := server.NewChannel("#testing")
	listener := func(event irc.Event) {
		mre := event.(irc.MessageReceivedEvent)
		fmt.Printf("%s: %s\n", mre.Nickname, mre.Content)
	}
	channel.AddListener((*irc.EventListener)(&listener))

	err := server.Start()
	if err != nil {
		fmt.Println(err)
	}
}
