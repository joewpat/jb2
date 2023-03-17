package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var openAiKey = readOpenAiKey()

/*
curl https://api.openai.com/v1/chat/completions \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer YOUR_API_KEY' \
  -d '{
  "model": "gpt-3.5-turbo",
  "messages": [{"role": "user", "content": "Hello!"}]
}'
*/

// openAiSearch uses a supplied API key to query openAi with supplied string.
// currently utilizes the davinci model completionrequest function.
func openAiSearch(query string) string {
	type OpenAI struct {
		Choices []struct {
			Message struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"message"`
			FinishReason string `json:"finish_reason"`
			Index        int    `json:"index"`
		} `json:"choices"`
	}
	apiKey := openAiKey
	if apiKey == "" {
		log.Fatalln("Missing API KEY")
	}

	client := &http.Client{Timeout: 60 * time.Second}

	requestBody := fmt.Sprintf(`{
	"model": "gpt-4",
	"messages": [{"role": "user", "content": "%s"}]
	}`, query)

	sendLog(requestBody)

	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer([]byte(requestBody)))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+openAiKey)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer resp.Body.Close()

	post := &OpenAI{}
	openaiErr := json.NewDecoder(resp.Body).Decode(post)
	if openaiErr != nil {
		fmt.Println(openaiErr)
	}

	sendLog(fmt.Sprintln(post))

	if post.Choices[0].Message.Content != "" {
		responseText := "```" + post.Choices[0].Message.Content + "```"
		return responseText
	}

	return "I cannot"
}

/*

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
*/

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
