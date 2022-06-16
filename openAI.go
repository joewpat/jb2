package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/PullRequestInc/go-gpt3"
)

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
		Prompt:    []string{"Write a limerick about Lee being gay"},
		MaxTokens: gpt3.IntPtr(256),
	})
	if err != nil {
		log.Fatalln(err)
	}
	ans := resp.Choices[0].Text
	fmt.Println(ans)
	responseText := "```" + ans + "```"
	return responseText

}

func readOpenAiKey() string {
	key, err := ioutil.ReadFile("discord.token")
	if err != nil {
		panic(err)
	}
	return string(key)
}
