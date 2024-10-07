package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

//new openAI function - gpt-4o-mini

// OpenAI API endpoint
const openAIURL = "https://api.openai.com/v1/chat/completions"

// OpenAIRequest struct defines the input to OpenAI API
type OpenAIChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature"`
}

// Message struct for chat messages
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIResponse struct defines the structure of the response from OpenAI API
type OpenAIChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func gpt(query string, context string) string {
	apiKey := openAiKey
	fmt.Println("API KEY: ", apiKey)
	fmt.Println(query)

	// Create the request payload
	reqBody := OpenAIChatRequest{
		Model: "gpt-4",
		Messages: []Message{
			{
				Role:    "system",
				Content: context,
			},
			{
				Role:    "user",
				Content: query,
			},
		},
		MaxTokens:   1000,
		Temperature: 0.7,
	}

	fmt.Println("Using mood: ", context)

	// Convert struct to JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		log.Fatalf("Failed to marshal request: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", openAIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	//fmt.Println("Response body: ", string(body))

	// Unmarshal the response into OpenAIResponse struct
	var apiResponse OpenAIChatResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		log.Fatalf("Failed to unmarshal response: %v", err)
	}

	//fmt.Println("Response: ", apiResponse)

	// Return the result
	return apiResponse.Choices[0].Message.Content
}
