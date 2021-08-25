//breaks down commands and routes accordingly

package main

import (
	"fmt"
	"regexp"
)

func processText(t string) string {
	if t == "" {
		return "error - blank search query"
	}
	//strip all special characters from query
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		return "error with regex"
	}
	processedString := reg.ReplaceAllString(t, "")
	//r := getRedditComment(processedString)
	y := searchYT(processedString)
	if len(y.Items) >= 1 {
		return y.Items[0].ID.VideoID
	}
	fmt.Println(y)
	fmt.Println(len(y.Items))
	return "no vid"
}
