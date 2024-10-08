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

// Set the context for JB's openAI responses. This is the prompt that will be used to set the mood of the AI.
var mood string = "You are a grouchy, inept assistant who frequently curses."

// Function to set or update the mood dynamically
func setSystemMessage(newMessage string) {
	mood = newMessage
}

func processText(t string, m *discordgo.Message) string {
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
	if strings.HasPrefix(t, "8ball") {
		return roll8ball()
	}
	if t == "whoami" {
		return fmt.Sprintf(m.Author.ID, " - ", m.Author.Username)
	}
	if strings.HasPrefix(t, "draw") {
		text := t[5:]
		fmt.Println("dall-e search for: ", text)
		return dallEText(text)
	}
	//set mood
	if strings.HasPrefix(t, "set mood to") && m.Author.ID == "325336510578819074" { //joe's ID only
		text := t[11:]
		setSystemMessage(text)
		return "mood set to: " + text
	}
	// manually run daily message
	if strings.HasPrefix(t, "dailymessagetest") && m.Author.ID == "325336510578819074" { //joe's ID only
		dailyMessage(readDiscordKey(), readChannelID())
	}

	return gpt(t, mood)
}
