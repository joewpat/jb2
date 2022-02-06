package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"time"

	strip "github.com/grokify/html-strip-tags-go"
)

// play example https://go.dev/play/p/qknP4uBUWFT

type Book struct {
	Data []struct {
		Abbr string `json:"abbreviation"`
		ID   string `json:"id"`
		Name string `json:"name"`
	}
}

type Chapter struct {
	Data []struct {
		Abbr string `json:"abbreviation"`
		ID   string `json:"id"`
		Name string `json:"name"`
	}
}

type Verse struct {
	Data []struct {
		Abbr string `json:"abbreviation"`
		ID   string `json:"id"`
		Name string `json:"name"`
	}
}

type VerseText struct {
	Data []struct {
		Content []struct {
			Text string
		}
	}
}

func readBibleAPIKey() string {
	key, err := ioutil.ReadFile("bible.key")
	if err != nil {
		panic(err)
	}
	return string(key)
}

func getBibleVerse() string {
	key := readBibleAPIKey()
	kjv := "de4e12af7f28f599-01"
	url := "https://api.scripture.api.bible/v1/bibles/" + kjv + "/books"
	client := &http.Client{Timeout: 3 * time.Second}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("api-key", key)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	var r Book

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(responseData, &r)
	if err != nil {
		fmt.Println(err)
	}
	var booklist []string
	for _, v := range r.Data {
		booklist = append(booklist, v.ID)
	}

	rand.Seed(time.Now().Unix())
	randBook := booklist[rand.Intn(len(booklist))]

	url = "https://api.scripture.api.bible/v1/bibles/" + kjv + "/books/" + randBook + "/chapters"
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Set("api-key", key)
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("response from bible api: ", resp.StatusCode)
	defer resp.Body.Close()

	responseData, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var c Chapter
	err = json.Unmarshal(responseData, &c)
	if err != nil {
		fmt.Println(err)
	}
	var chapterlist []string
	for _, v := range c.Data {
		chapterlist = append(chapterlist, v.ID)
	}

	rand.Seed(time.Now().Unix())
	randChapter := chapterlist[rand.Intn(len(chapterlist))]
	fmt.Println(randChapter)

	url = "https://api.scripture.api.bible/v1/bibles/" + kjv + "/chapters/" + randChapter + "/verses"
	fmt.Println(url)
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Set("api-key", key)
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("response from bible api: ", resp.StatusCode)
	defer resp.Body.Close()

	responseData, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var v Verse
	err = json.Unmarshal(responseData, &v)
	if err != nil {
		fmt.Println(err)
	}
	var verselist []string
	for _, v := range v.Data {
		verselist = append(verselist, v.ID)
	}

	rand.Seed(time.Now().Unix())
	randverse := verselist[rand.Intn(len(verselist))]

	url = "https://api.scripture.api.bible/v1/bibles/" + kjv + "/verses/" + randverse
	fmt.Println(url)
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Set("api-key", key)
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("response from bible api: ", resp.StatusCode)
	defer resp.Body.Close()

	responseData, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var f interface{}
	err = json.Unmarshal(responseData, &f)
	if err != nil {
		fmt.Println(err)
	}

	m := f.(map[string]interface{})
	fmt.Println()
	//fmt.Println(m["data"])
	n := m["data"].(map[string]interface{})
	verse := n["content"]
	strverse := fmt.Sprintf("%v", verse)
	fixedVerse := strip.StripTags(strverse)
	reg, err := regexp.Compile("[^a-zA-Z.,!?]+")
	if err != nil {
		log.Fatal(err)
	}
	fixedVerse = reg.ReplaceAllString(fixedVerse, " ")

	fmt.Println(fixedVerse)

	return fixedVerse
}
