package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/PullRequestInc/go-gpt3"
)

var openAiKey = readOpenAiKey()

//openAiSearch uses a supplied API key to query openAi with supplied string.
//currently utilizes the davinci model completionrequest function.
func openAiSearch(query string) string {
	apiKey := openAiKey
	if apiKey == "" {
		log.Fatalln("Missing API KEY")
	}

	ctx := context.Background()
	client := gpt3.NewClient(apiKey, gpt3.WithDefaultEngine("text-davinci-002"))

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

//Dall-E image generation based on text query
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
	"size": "256x256"
	}`, query)

	fmt.Println(requestBody)

	//resp, _ := client.Get("https://api.openai.com/v1/images/generations")
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

	/*
		responseData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return ""
		}

		fmt.Println(responseData)
	*/

	if post.Data[0].URL != "" {
		return post.Data[0].URL
	}

	return "I cannot"
}

//reads openai key from app directory.
//TODO: move to env variables
func readOpenAiKey() string {
	key, err := ioutil.ReadFile("openAi.key")
	if err != nil {
		panic(err)
	}
	return string(key)
}
