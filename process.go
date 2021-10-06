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
	if strings.HasPrefix(t, "gif ") {
		text := t[5:]
		t = strings.Replace(text, " ", "+", -1)
		tenor := searchTenor(t)
		giphy := searchGiphy(t)
		var responses []string
		responses = append(responses, tenor, giphy)
		return responses[rand.Intn(len(responses))]
	}
	if strings.HasPrefix(t, "time") {
		time := time.Now()
		return time.Local().String()
	}
	if t == "surf" {
		r := getSurfReport()
		return r
	}
	t = strings.Replace(t, " ", "+", -1)
	var responses []string
	r := getRedditComment(t)
	y := youtube(t)
	responses = append(responses, r, y)
	fmt.Println("youtube:", y)
	fmt.Println("reddit:", r)
	return responses[rand.Intn(len(responses))]
}
