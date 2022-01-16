package main

import (
	"fmt"
	"html"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-co-op/gocron"
)

func main() {
	var token = readDiscordKey()
	var channelID = readChannelID()
	//removing newline character if added from env variables
	token = strings.TrimSuffix(token, "\n")
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	//add handlers for discord events
	dg.AddHandler(messageCreate)
	dg.AddHandler(onReady)
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	//go-cron scheduler for daily messasge
	s := gocron.NewScheduler(time.UTC)
	s.Every(1).Day().At("13:00").Do(dailyMessage, token, channelID)
	s.StartAsync()

	// ctrl+c to quit
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}
	if strings.HasPrefix(m.Content, "jb ") {
		fmt.Println("triggered", m.Content[3:])
		text := m.Content[3:]
		go func() {
			s.ChannelTyping(m.ChannelID)
		}()
		resp := processText(text)
		resp = html.UnescapeString(resp)
		fmt.Println("Final Reply: \n", resp+"\n")
		s.ChannelMessageSend(m.ChannelID, resp)
	}
}

func onReady(s *discordgo.Session, r *discordgo.Ready) {
	testChannel := "716083032188256308"
	s.ChannelMessageSend(testChannel, "I have been redeployed.")
}

//runs daily for the motivational message
func dailyMessage(token string, channelID string) {
	dg, _ := discordgo.New("Bot " + token)
	_ = dg.Open()
	fmt.Println("Daily message triggered.")
	bquote := getBibleVerse() + "\n:pray:"
	message := "Good morning friends! The theme of today is:\n" + bquote
	dg.ChannelMessageSend(channelID, message)
}

func readDiscordKey() string {
	key, err := ioutil.ReadFile("discord.token")
	if err != nil {
		panic(err)
	}
	return string(key)
}

func readChannelID() string {
	key, err := ioutil.ReadFile("discord.channelID")
	if err != nil {
		panic(err)
	}
	return string(key)
}
