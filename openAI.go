package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/PullRequestInc/go-gpt3"
)

//openAiSearch uses a supplied API key to query openAi with supplied string.
//currently utilizes the davinci model completionrequest function.
func openAiSearch(query string) string {
	//godotenv.Load()

	openAiKey := readOpenAiKey()

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

//reads openai key from app directory.
//TODO: move to env variables
func readOpenAiKey() string {
	key, err := ioutil.ReadFile("openAi.key")
	if err != nil {
		panic(err)
	}
	return string(key)
}
