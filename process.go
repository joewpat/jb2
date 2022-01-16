//breaks down commands and routes to the appropriate functions
//handles most of the "brains" :-| of the bot

package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func processText(t string) string {
	rand.Seed(time.Now().UnixNano()) //init random
	if t == "" {
		return "error - blank search query"
	}
	fmt.Println("query: ", t)
	if strings.HasPrefix(t, "gif ") {
		text := t[4:]
		fmt.Println("searching for gifs, query: ", text)
		t = strings.Replace(text, " ", "+", -1)
		tenor := searchTenor(t)
		giphy := searchGiphy(t)
		var responses []string
		if strings.HasPrefix(tenor, "http") {
			responses = append(responses, tenor)
			fmt.Println("response from tenor: ", tenor)
		}
		if strings.HasPrefix(giphy, "http") {
			responses = append(responses, giphy)
			fmt.Println("response from giphy: ", giphy)
		}
		if len(responses) > 0 {
			return responses[rand.Intn(len(responses))]
		}
		return "error - no gifs found"
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
		r := getSurfReport()
		return r
	}
	t = strings.Replace(t, " ", "+", -1)
	var responses []string
	r := getRedditComment(t)
	y := youtube(t)
	if r != "" {
		responses = append(responses, r)
	}
	if y != "" {
		responses = append(responses, y)
	}
	fmt.Println("youtube:", y)
	fmt.Println("reddit:", r)
	if len(responses) > 0 {
		return responses[rand.Intn(len(responses))]
	}
	return "error - no responses found for query: " + t
}
