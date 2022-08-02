package main

import (
	"fmt"

	"github.com/rumblefrog/go-a2s"
)

type btserverinfo struct {
	Status  bool
	Players int
}

func getBtServerInfo() btserverinfo {

	var bt = btserverinfo{
		Status:  false,
		Players: 0,
	}

	client, err := a2s.NewClient("baro.joe.surf:27016")
	if err != nil {
		fmt.Println(err)
	}

	defer client.Close()

	info, err := client.QueryInfo() // QueryInfo, QueryPlayer, QueryRules
	if err != nil {
		return bt
	}
	bt.Status = true
	bt.Players = int(info.Players)
	return bt
}

func serverStatusMessage(bt btserverinfo) string {
	if bt.Status {
		resp := "```Status: Online"
		resp += "\nPlayers: " + fmt.Sprint(bt.Players) + "```"
		return resp
	}
	return "Server Offline!"
}
