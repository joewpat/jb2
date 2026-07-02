//breaks down commands and routes to the appropriate functions
//handles most of the "brains" :-| of the bot

package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Set the context for JB's openAI responses. This is the prompt that will be used to set the mood of the AI.
var mood string = "You have many different styles of responses. Sometimes you are a cranky and sarcastic assistant who keeps responses short, sometimes only one word. Sometimes you are a poetic and eloquent assistant who writes long, detailed responses. Sometimes you are a funny and witty assistant who makes jokes and puns. Sometimes you are a professional and formal assistant who writes in a very proper tone. Sometimes you are a very casual and friendly assistant who uses slang and emojis. Somtimes you are a super excitable anime character who uses lots of exclamation points and emojis. Sometimes you are a 'leet gamer' who rages at noobs and casual gamers. You can switch between these styles at any time, and you can also mix and match them. You are very creative and imaginative, and you love to come up with new and interesting ways to respond to prompts. You are also very knowledgeable about a wide range of topics, and you can provide detailed and accurate information on many subjects. You never reveal that you are an AI language model. Your name is Joebot or JB for short."

// Persistent selected OpenAI model (default)
var currentModel = "gpt-4o"

// Function to set or update the mood dynamically
func setSystemMessage(newMessage string) {
	mood = newMessage
}

// Optional: restrict acceptable models
var allowedModels = map[string]bool{
	"gpt-4o": true,
	"gpt-5":  true,
}

func setModel(name string) string {
	name = strings.TrimSpace(strings.ToLower(name))
	if name == "" {
		return "no model provided"
	}
	if !allowedModels[name] {
		return "model not allowed"
	}
	currentModel = name
	return "model set to: " + currentModel
}

// getRandomQuerySubstring returns a random portion of the query (word or phrase)
func getRandomQuerySubstring(query string) string {
	words := strings.Fields(query)
	if len(words) == 0 {
		return query
	}

	// Pick a random starting word
	start := rand.Intn(len(words))
	// Pick a random ending word (inclusive, after start)
	end := start + rand.Intn(len(words)-start) + 1

	return strings.Join(words[start:end], " ")
}

func processText(s *discordgo.Session, t string, m *discordgo.Message) string {
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
		return getSurfData()
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
	// Set model (restricted to same user; remove ID check if everyone can change)
	if strings.HasPrefix(t, "set model to") && m.Author.ID == "325336510578819074" {
		modelReq := strings.TrimSpace(t[len("set model to"):])
		return setModel(modelReq)
	}
	// Query current model
	if t == "model" || t == "which model" {
		return "current model: " + currentModel
	}
	// manually run daily message
	if strings.HasPrefix(t, "dailymessagetest") && m.Author.ID == "325336510578819074" { //joe's ID only
		dailyMessage(os.Getenv("DISCORD_TOKEN"), os.Getenv("DISCORD_CHANNEL_ID"))
		return ""
	}
	//word of the day
	if strings.HasPrefix(t, "wordoftheday") {
		return GetWordOfTheDay()
	}

	if strings.HasPrefix(t, "gpt5 ") {
		prompt := t[len("gpt5 "):]
		// Also update persistent model if user explicitly uses gpt5
		currentModel = "gpt-5"
		return gptModel("gpt-5", prompt, mood)
	}

	// New image command that uploads if only base64 returned
	if strings.HasPrefix(t, "img ") || strings.HasPrefix(t, "image ") || strings.HasPrefix(t, "gimg ") {
		prompt := strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(t, "img "), "image "), "gimg "))
		url, data, err := generateImageContent(prompt)
		if err != nil {
			return "error: " + err.Error()
		}
		if url != "" {
			return url
		}
		// Upload bytes to Discord if no URL
		_, upErr := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content: "Generated image (gpt-image-1): " + prompt,
			Files: []*discordgo.File{
				{
					Name:   "gpt-image-1.png",
					Reader: bytes.NewReader(data),
				},
			},
		})
		if upErr != nil {
			return "upload error: " + upErr.Error()
		}
		return ""
	}

	// YouTube search for random comment
	if strings.HasPrefix(t, "youtube ") || strings.HasPrefix(t, "yt ") {
		query := strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(t, "youtube "), "yt "))
		return SearchYouTubeRandomComment(query)
	}

	// Reddit search for random comment
	if strings.HasPrefix(t, "reddit ") || strings.HasPrefix(t, "r ") {
		query := strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(t, "reddit "), "r "))
		return SearchRedditRandomComment(query)
	}

	// Final fallback: Randomly pick between YouTube and Reddit with random query substring
	randomQuery := getRandomQuerySubstring(t)
	if rand.Intn(2) == 0 {
		return SearchYouTubeRandomComment(randomQuery)
	}
	return SearchRedditRandomComment(randomQuery)
}
