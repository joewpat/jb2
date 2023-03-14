package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PullRequestInc/go-gpt3"
)

var openAiKey = readOpenAiKey()

// openAiSearch uses a supplied API key to query openAi with supplied string.
// currently utilizes the davinci model completionrequest function.
func openAiSearch(query string) string {
	apiKey := openAiKey
	if apiKey == "" {
		log.Fatalln("Missing API KEY")
	}

	ctx := context.Background()
	client := gpt3.NewClient(apiKey, gpt3.WithDefaultEngine("text-davinci-003"))

	resp, err := client.Completion(ctx, gpt3.CompletionRequest{
		Prompt:    []string{query},
		MaxTokens: gpt3.IntPtr(512),
	})
	if err != nil {
		log.Fatalln(err)
	}
	ans := resp.Choices[0].Text
	fmt.Println(ans)
	responseText := "```" + ans + "```"
	return responseText

}

// Dall-E image generation based on text query
func dallEText(query string) string {
	type Dalle struct {
		Created int `json:"created"`
		Data    []struct {
			URL string `json:"url"`
		} `json:"data"`
	}

	client := &http.Client{Timeout: 30 * time.Second}

	requestBody := fmt.Sprintf(`{
	"prompt": "%s",
	"n": 1,
	"size": "1024x1024"
	}`, query)

	sendLog(requestBody)

	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/images/generations", bytes.NewBuffer([]byte(requestBody)))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+openAiKey)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer resp.Body.Close()

	post := &Dalle{}
	derr := json.NewDecoder(resp.Body).Decode(post)
	if derr != nil {
		fmt.Println(derr)
	}

	if post.Data[0].URL != "" {
		return post.Data[0].URL
	}

	return "I cannot"
}

// reads openai key from app directory.
// TODO: move to env variables/kv
func readOpenAiKey() string {
	key, err := os.ReadFile("openAi.key")
	if err != nil {
		panic(err)
	}
	return string(key)
}
