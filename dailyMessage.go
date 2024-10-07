package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
)

// runs daily for the motivational message
func dailyMessage(token string, channelID string) {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		sendLog(err.Error())
		log.Fatal(err)
	}
	err = dg.Open()
	if err != nil {
		sendLog(err.Error())
		log.Fatal(err)
	}
	sendLog("Daily message triggered.")

	dotw := time.Now().Weekday()
	day := time.Now().Day()
	month := time.Now().Month()

	greeting := fmt.Sprintf("Good morning! Today is %s, %s %d.\n", dotw, month, day)
	forecast := weatherForecast()
	bquote := getBibleVerse()
	//deepthought := getDeepThought()

	var choices []string
	choices = append(choices, bquote)
	choice := choices[rand.Intn(len(choices))]
	theme := "\nThe theme of today is:\n" + choice
	message := greeting + "\n" + forecast + theme
	dg.ChannelMessageSend(channelID, message)
}

// func dailySurfMessage(token string, channelID string) {
// 	dg, err := discordgo.New("Bot " + token)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	err = dg.Open()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	message := "sup bras, here's this morning's surf report and forecast :call_me:\n" + getSurflineForecast()
// 	dg.ChannelMessageSend(channelID, message)
// }
