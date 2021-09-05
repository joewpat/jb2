package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("jb-shell 1.0")

	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)
		if strings.Compare("exit", text) == 0 {
			fmt.Println("goodbye")
			break
		}
		if strings.Compare("hi", text) == 0 {
			fmt.Println("hello, Yourself")
		} else {
			fmt.Println(text)
			resp := processText(text)
			fmt.Println(resp)
		}

	}
}
