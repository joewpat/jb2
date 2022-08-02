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

/*
const regex = `<.*?>`

// stripHtmlRegex uses regular expresion to remove HTML tags.
func stripHtmlRegex(s string) string {
	r := regexp.MustCompile(regex)
	return r.ReplaceAllString(s, "")
}*/

func processText(t string, m *discordgo.Message, session *discordgo.Session) string {
	rand.Seed(time.Now().UnixNano()) //init random
	if t == "" {
		return "error - blank search query"
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
	return jb(t)
}
