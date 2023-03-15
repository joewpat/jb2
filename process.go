//breaks down commands and routes to the appropriate functions
//handles most of the "brains" :-| of the bot

package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func processText(t string, m *discordgo.Message, session *discordgo.Session) string {
	rand.Seed(time.Now().UnixNano()) //init random
	if t == "" {
		return "what?"
	}
	fmt.Println("query: ", t)
	if strings.HasPrefix(t, "gif ") {
		text := t[4:]
		fmt.Println("searching for gifs, query: ", text)
		t = strings.Replace(text, " ", "+", -1)
		return searchGifs(t)
	}
	if strings.HasPrefix(t, "time") {
		loc := time.FixedZone("UTC-5", -5*60*60)
		time := time.Now().In(loc)
		return time.String()
	}
	if strings.HasPrefix(t, "bible") {
		return getBibleVerse()
	}
	if t == "surf" {
		return getSurflineForecast()
	}
	if strings.HasPrefix(t, "roulette") {
		return roulette(m, session)
	}
	if strings.HasPrefix(t, "-ai") {
		text := t[4:]
		fmt.Println("openAI search for ", text)
		return openAiSearch(text)
	}
	if t == "baro status" {
		return serverStatusMessage(getBtServerInfo())
	}
	if strings.HasPrefix(t, "8ball") {
		return roll8ball()
	}
	if t == "whoami" {
		return fmt.Sprintf(m.Author.ID, " - ", m.Author.Username)
	}
	if t == "baro shutdown" {
		if m.Author.ID == "325336510578819074" {
			session.ChannelMessageSend(m.ChannelID, "initiating baro shutdown")
			btShutdown()
			return "bt server shutdown complete"
		} else {
			return "YOU ARE UNAUTHORIZED TO PERFORM THIS COMMAND. THIS INCIDENT HAS BEEN LOGGED"
		}
	}
	if t == "baro start" {
		btServerStart()
		return "starting BT server..."
	}
	if strings.HasPrefix(t, "draw") {
		text := t[5:]
		fmt.Println("dall-e search for: ", text)
		return dallEText(text)
	}
	return jb(t)
}
