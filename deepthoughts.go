package main

import (
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func getDeepThought() string {
	url := "https://joe.surf/deepthoughts"
	client := &http.Client{Timeout: 3 * time.Second}
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		sendLog(fmt.Sprintln(err))
	}

	sendLog(fmt.Sprintln("Request--joe.surf/deepthoughts status: ", resp.StatusCode))
	responseBody := fmt.Sprintf("%s", responseData)

	//strip out HTML
	dt := strings.Split(responseBody, "p>")
	dt = strings.Split(dt[1], "</")

	//isolate deep thought from response, and clean up formatting
	deepThought := dt[0]
	deepThought = html.UnescapeString(deepThought)
	deepThought = strings.TrimSpace(deepThought)
	space := regexp.MustCompile(`\s+`)
	deepThought = space.ReplaceAllString(deepThought, " ")

	sendLog(deepThought)

	return deepThought
}
