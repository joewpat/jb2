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

// reads openai key from app directory.
// TODO: move to env variables/kv
func readOpenAiKey() string {
	key, err := os.ReadFile("openAi.key")
	if err != nil {
		panic(err)
	}
	return string(key)
}

// openAiSearch uses a supplied API key to query openAi with supplied string.
// currently utilizes the davinci model completionrequest function.
func gpt3(query string) string {
	type OpenAIGPT3Response struct {
		ID      string `json:"id"`
		Object  string `json:"object"`
		Created int    `json:"created"`
		Model   string `json:"model"`
		Choices []struct {
			Text         string      `json:"text"`
			Index        int         `json:"index"`
			Logprobs     interface{} `json:"logprobs"`
			FinishReason string      `json:"finish_reason"`
		} `json:"choices"`
		Usage struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		} `json:"usage"`
	}
	apiKey := openAiKey
	if apiKey == "" {
		log.Fatalln("Missing API KEY")
	}

	client := &http.Client{Timeout: 120 * time.Second}

	//GPT3 Text completion:
	requestBody := fmt.Sprintf(`{
	"model": "text-davinci-003",
	"prompt": "%s",
	"max_tokens": 4000
	}`, query)

	/*requestBody := fmt.Sprintf(`{
	"model": "gpt-4",
	"messages": [{"role": "user", "content": "%s"}]
	}`, query)*/

	sendLog(requestBody)

	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/completions", bytes.NewBuffer([]byte(requestBody)))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+openAiKey)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer resp.Body.Close()

	post := &OpenAIGPT3Response{}
	openaiErr := json.NewDecoder(resp.Body).Decode(post)
	if openaiErr != nil {
		fmt.Println(openaiErr)
	}

	fmt.Println("Response data: ")
	sendLog(fmt.Sprintln(post))

	if post.Choices != nil {
		responseText := "```" + post.Choices[0].Text + "```"
		return responseText
	}

	return "Error retreiving OpenAI response"
}

// GPT 4:
func gpt4(query string) string {
	type OpenAIGPT4Response struct {
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

	client := &http.Client{Timeout: 120 * time.Second}

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

	post := &OpenAIGPT4Response{}
	openaiErr := json.NewDecoder(resp.Body).Decode(post)
	if openaiErr != nil {
		fmt.Println(openaiErr)
	}

	sendLog(fmt.Sprintln(post))

	if post.Choices != nil {
		responseText := "```" + post.Choices[0].Message.Content + "```"
		return responseText
	}

	return "Error retreiving OpenAI response"
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
