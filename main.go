package main

import (
	"fmt"
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
	var channelID = readChannelID() // used for daily message
	//removing newline character if added from env variables
	token = strings.TrimSuffix(token, "\n")
	fmt.Println("Discord Token: ", token)
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
	s.Every(1).Day().At("11:30").Do(dailyMessage, token, channelID)
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
		resp := processText(text, m.Message)
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

func readDiscordKey() string {
	key, err := os.ReadFile("discord.token")
	if err != nil {
		panic(err)
	}
	return string(key)
}

func readChannelID() string {
	key, err := os.ReadFile("discord.channelID")
	if err != nil {
		panic(err)
	}
	return string(key)
}

// func readLogChannelID() string {
// 	key, err := os.ReadFile("discord.logChannelID")
// 	if err != nil {
// 		panic(err)
// 	}
// 	return string(key)
// }

// func readSurfChannelID() string {
// 	key, err := os.ReadFile("discord.surfChannelID")
// 	if err != nil {
// 		panic(err)
// 	}
// 	return string(key)
// }

// func readSubOptimalChannelID() string {
// 	key, err := os.ReadFile("discord.suboptimalchannelID")
// 	if err != nil {
// 		panic(err)
// 	}
// 	return string(key)
// }
