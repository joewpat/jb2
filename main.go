package main

import (
	"fmt"
	"html"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func readDiscordKey() string {
	key, err := ioutil.ReadFile("discord.token")
	if err != nil {
		panic(err)
	}
	return string(key)
}

func main() {

	//create a simple listner on port 8080 to satisfy DigitalOcean's health checks
	n := "tcp"
	l, err := net.Listen(n, "localhost:8080")
	if err != nil {
		fmt.Println(err)
	}
	defer l.Close()

	var token = readDiscordKey()
	//removing newline character if added from env variables
	token = strings.TrimSuffix(token, "\n")
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
	for {
		_, err := l.Accept()
		if err != nil {
			fmt.Println(err)
		}
	}
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "jb ") {
		fmt.Println("triggered", m.Content[3:])
		text := m.Content[3:]
		resp := processText(text)
		resp = html.UnescapeString(resp)
		fmt.Println("Final Reply: \n", resp+"\n")
		s.ChannelMessageSend(m.ChannelID, resp)
	}
}
