//breaks down commands and routes accordingly

package main

import (
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
	answer := getRedditComment(processedString)
	return answer
}
