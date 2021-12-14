package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"io/ioutil"
)

//	https://stackoverflow.com/questions/17156371/how-to-get-json-response-from-http-get
//	fmt.Println("Permalink: ", r.Data.Children[1].Data.Permalink)
//	RedditResponse struct hold relevant response data from reddit api JSON response
//	get one comment example
//	fmt.Println(r[1].Data.Children[1].Data.Body)

type RedditResponse struct {
	Kind string `json:"kind"`
	Data struct {
		Children []struct {
			Kind string `json:"kind"`
			Data struct {
				SubredditID string `json:"subreddit_id"`
				Permalink   string `json:"permalink"`
				URL         string `json:"URL"`
				NumComments int    `json:"num_comments"`
				ID          string `json:"id"`
			} `json:"data,omitempty"`
		} `json:"children"`
	} `json:"data"`
}

type RedditComment struct {
	Kind string `json:"kind"`
	Data struct {
		Children []struct {
			Kind string `json:"kind"`
			Data struct {
				Body        string `json:"body"`
				Author      string `json:"author"`
				URL         string `json:"URL"`
				NumComments int    `json:"num_comments"`
				ID          string `json:"id"`
			} `json:"data,omitempty"`
		} `json:"children"`
	} `json:"data"`
}

func readUserAgentString() string {
	key, err := ioutil.ReadFile("reddit.ua.string")
	if err != nil {
		panic(err)
	}
	return string(key)
}

//the primary function - searches reddit for a comment based on text query
func getRedditComment(t string) string {
	s := searchReddit(t)
	if len(s.Data.Children) < 1 {
		return "jberror - no reddit content found"
	}
	time.Sleep(time.Second * 1) // delays are in place to satisfy API requirements (max 60req/min)
	rand.Seed(time.Now().Unix())
	randomPost := s.Data.Children[rand.Intn(len(s.Data.Children))] //pull random post from the results
	url := "https://reddit.com" + randomPost.Data.Permalink + ".json"
	rc := getComments(url)
	if len(rc) > 0 {
		return getRandomComment(rc)
	} else {
		return ""
	}
}

// searchreddit searches reddit for content based on text query and returns a RedditResponse struct
func searchReddit(query string) RedditResponse {
	//build http client and request
	url := "https://www.reddit.com/search.json?q=" + query + "&include_over_18=on&limit=50"
	client := &http.Client{Timeout: 5 * time.Second} //set this to 5 seconds due to reddit being slow
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Golang_Reddit_Bot/0.1 by /u/Robert_Arctor")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	//fmt.Println(resp.Status) - log this instead
	var r RedditResponse
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		fmt.Println("ERROR - ", err)
	}
	return r
}

//func getComment takes a reddit post id and returns a slice of comments
func getComments(url string) []RedditComment {
	client := &http.Client{Timeout: 5 * time.Second}
	req, _ := http.NewRequest("GET", url, nil)
	uastring := readUserAgentString()
	req.Header.Set("User-Agent", uastring)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	var r []RedditComment
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		panic(err)
	}
	return r
}

//getRandomComment takes a slice of threads and returns a random RedditComment
func getRandomComment(r []RedditComment) string {
	comment := ""
	for comment == "" {
		if len(r) > 0 {
			thread := r[rand.Intn(len(r))].Data.Children
			if len(thread) > 1 {
				comment = thread[rand.Intn(len(thread))].Data.Body
			}
			comment = thread[0].Data.Body
		} else {
			thread := r[0].Data.Children
			if len(thread) > 1 {
				comment = thread[rand.Intn(len(thread))].Data.Body
			}
			comment = thread[0].Data.Body
		}

	}
	for strings.Contains(comment, "https") {
		fmt.Println("ignored comment for containing hyperlink: ", comment)
		thread := r[rand.Intn(len(r))].Data.Children
		comment = thread[rand.Intn(len(thread))].Data.Body
	}
	return comment
}
