//breaks down commands and routes accordingly

package main

import (
	"fmt"
	"strings"
)

func processText(t string) string {
	if t == "" {
		return "error - blank search query"
	}
	t = strings.Replace(t, " ", "+", -1)
	r := getRedditComment(t)
	y := youtube(t)
	fmt.Println("youtube: ", y)
	fmt.Println("reddit:", r)
	return "done"
}

//strip all special characters from query
//reg, err := regexp.Compile("^[a-zA-Z0-9 ]*$")
//if err != nil {
//	return "error with regex"
//}
//processedString := reg.ReplaceAllString(t, "")
//fmt.Println("processed string: ", processedString)
//r := getRedditComment(processedString)
