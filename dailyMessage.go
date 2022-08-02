package main

import (
	"log"
	"math/rand"

	"github.com/bwmarrin/discordgo"
)

//runs daily for the motivational message
func dailyMessage(token string, channelID string) {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}
	err = dg.Open()
	if err != nil {
		log.Fatal(err)
	}
	sendLog("Daily message triggered.")

	//dotw := time.Now().Weekday()
	bquote := getBibleVerse()
	deepthought := getDeepThought()

	var choices []string
	choices = append(choices, bquote, deepthought)
	choice := choices[rand.Intn(len(choices))]
	message := "Good morning friends! The theme of today is:\n" + choice
	dg.ChannelMessageSend(channelID, message)
}

func dailySurfMessage(token string, surfChannelID string) {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}
	err = dg.Open()
	if err != nil {
		log.Fatal(err)
	}
	message := "sup bras, here's this morning's surf report and forecast :call_me:\n"+getSurflineForecast()
	dg.ChannelMessageSend(surfChannelID, message)
}