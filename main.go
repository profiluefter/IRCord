package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/profiluefter/IRCord/irc"
	"os"
)

func main() {
	discordToken := os.Getenv("DISCORD_TOKEN")
	guildID := os.Getenv("GUILD_ID")

	discord, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		panic(err.Error())
	}

	discord.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	err = discord.Open()
	if err != nil {
		panic(err.Error())
	}

	err = discord.UpdateStatus(0, "")
	if err != nil {
		panic(err.Error())
	}

	channels, err := discord.GuildChannels(guildID)
	if err != nil {
		panic(err.Error())
	}

	var motd = "This is the message of the day!\nIf you can see this then the server did not crash yet\nNice."

	options := irc.ServerOptions{
		Name: "irc-cord",
		Port: 6667,
		Motd: &motd,
	}
	server := irc.NewServer(options)

	var discordIDToIRCChannel = map[string]*irc.Channel{}
	var ircChannelToWebhook = map[string]*discordgo.Webhook{}

	for _, discordChannel := range channels {
		discordChannel := discordChannel
		if discordChannel.Type != discordgo.ChannelTypeGuildText {
			continue
		}

		fmt.Printf("Adding channel %s with topic \"%s\"!\n", discordChannel.Name, discordChannel.Topic)
		channel := server.NewChannel("#"+discordChannel.Name, discordChannel.Topic)

		webhooks, err := discord.ChannelWebhooks(discordChannel.ID)
		if err != nil {
			panic(err.Error())
		}
		var channelWebhook *discordgo.Webhook
		if len(webhooks) > 0 {
			channelWebhook = webhooks[0]
		} else {
			webhook, err := discord.WebhookCreate(discordChannel.ID, "IRCord", "")
			if err != nil {
				panic(err.Error())
			}
			channelWebhook = webhook
		}

		ircChannelToWebhook[channel.Name] = channelWebhook

		listener := func(event irc.Event) {
			mre := event.(irc.MessageReceivedEvent)

			webhook := ircChannelToWebhook[channel.Name]
			_, _ = discord.WebhookExecute(webhook.ID, webhook.Token, false, &discordgo.WebhookParams{
				Content:  mre.Content,
				Username: mre.Nickname,
			})
		}
		channel.AddListener((*irc.EventListener)(&listener))

		discordIDToIRCChannel[discordChannel.ID] = channel
	}

	discord.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.Bot {
			return //Ignored
		}
		if m.GuildID != guildID {
			panic(fmt.Sprintf("Message received from unknown server: %s", m.GuildID))
		}
		channel := discordIDToIRCChannel[m.ChannelID]
		if channel == nil {
			panic(fmt.Sprintf("Unknown channel %s", m.ChannelID))
		}
		channel.SendMessage(m.Author.Username, m.Content)
	})

	err = server.Start()
	if err != nil {
		fmt.Println(err)
	}
}
