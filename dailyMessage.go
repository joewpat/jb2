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
	//genAlphaForecast := gpt("Translate this weather forecast in Gen alpha lingo, and shorten it to just two sentences. Emoji use is encouraged.: "+forecast, "Super hip Gen Alpha meteorologist.")
	wotd := GetWordOfTheDay()
	message := greeting + "\n" + forecast + "\n\n" + wotd //+ "\n" + getSurfData()
	dg.ChannelMessageSend(channelID, message)
}
