package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-co-op/gocron"
)

func main() {
	var token = os.Getenv("DISCORD_TOKEN")
	var channelID = os.Getenv("DISCORD_CHANNEL_ID") // used for daily message
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
	s.Every(1).Day().At("12:30").Do(dailyMessage, token, channelID)
	//s.Every(1).Day().At("11:30").Do(dailyMessage, token, subOptimalChannelID)
	//s.Every(1).Day().At("11:40").Do(dailySurfMessage, token, surfChannelID)
	//893136225152692284
	s.StartAsync()

	// ctrl+c to quit
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	//log chats

	if strings.HasPrefix(strings.ToLower(m.Content), "jb ") {
		fmt.Println("triggered", m.Content[3:])
		text := m.Content[3:]
		go func() {
			s.ChannelTyping(m.ChannelID)
		}()
		resp := processText(s, text, m.Message)
		fmt.Println("Final Reply: \n", resp+"\n")
		s.ChannelMessageSend(m.ChannelID, resp)
	} else if strings.Contains(strings.ToLower(m.Content), "another one") {
		//dj khaled gif search
		go func() {
			s.ChannelTyping(m.ChannelID)
		}()
		resp := djKhaledGif()
		s.ChannelMessageSend(m.ChannelID, resp)
		return
	} else {
		// Random chance to respond to any message (5% probability)
		rand.Seed(time.Now().UnixNano())
		if rand.Float64() < 0.05 {
			fmt.Println("randomly responding to:", m.Content)
			go func() {
				s.ChannelTyping(m.ChannelID)
			}()
			resp := processText(s, m.Content, m.Message)
			fmt.Println("Random Reply: \n", resp+"\n")
			s.ChannelMessageSend(m.ChannelID, resp)
		}
	}
}

// sendLog is used to send log messages to discord for easy/fun debugging
func sendLog(t string) {
	// token := readDiscordKey()
	// dg, _ := discordgo.New("Bot " + token)
	// _ = dg.Open()
	// logChannel := readLogChannelID()
	// dg.ChannelMessageSend(logChannel, t)
	fmt.Println(t)
}

func onReady(s *discordgo.Session, r *discordgo.Ready) {
	sendLog("jb status: ready")
}
