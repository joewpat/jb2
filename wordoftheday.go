package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var apiKey = os.Getenv("WORDNIK_API_KEY")

const wordOfTheDayURL = "https://api.wordnik.com/v4/words.json/wordOfTheDay"

var definitionText string

// WordOfTheDayResponse represents the structure of the Wordnik API response
type WordOfTheDayResponse struct {
	Word        string `json:"word"`
	Definitions []struct {
		Text string `json:"text"`
	} `json:"definitions"`
}

// GetWordOfTheDay fetches the word of the day from the Wordnik API
func GetWordOfTheDay() string {
	req, err := http.NewRequest("GET", wordOfTheDayURL, nil)
	if err != nil {

		return fmt.Sprint(err)
	}

	query := req.URL.Query()
	query.Add("api_key", apiKey)
	req.URL.RawQuery = query.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Sprint(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("failed to fetch word of the day: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprint(err)
	}

	var wordOfTheDay WordOfTheDayResponse
	if err := json.Unmarshal(body, &wordOfTheDay); err != nil {
		return fmt.Sprint(err)
	}
	if len(wordOfTheDay.Definitions) > 0 {
		definitionText = wordOfTheDay.Definitions[0].Text
	} else {
		definitionText = "No definition available"
	}
	response := fmt.Sprintf("Word of the Day: %s\nDefinition: %s", wordOfTheDay.Word, definitionText)

	fmt.Println(response)
	// Return the word of the day
	// and its definition
	return response
}
