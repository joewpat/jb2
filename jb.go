package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func jb(query string) string {
	url := "https://joe.surf/jb/raw/" + `"` + query + `"`
	client := &http.Client{Timeout: 20 * time.Second}
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		sendLog(err.Error())
		return ""
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(responseData)
}
