package main

import (
	"fmt"
	"log"
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
	genZforecast := gpt("Translate this forecast to Gen Z lingo: "+forecast, "Super hip Gen Z meteorologist.")
	//bquote := getBibleVerse()
	wotd := GetWordOfTheDay()
	//dayGif := searchGifs(dotw.String() + "blessings")
	//deepthought := getDeepThought()

	// var choices []string
	// choices = append(choices, bquote)
	// choice := choices[rand.Intn(len(choices))]
	// theme := "\nToday's bible verse is:\n" + choice
	message := greeting + "\n" + genZforecast + "\n\n" + wotd + "\n" + getSurfData()
	dg.ChannelMessageSend(channelID, message)
}
