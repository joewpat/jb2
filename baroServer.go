package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/rumblefrog/go-a2s"
	"github.com/yahoo/vssh"
)

type btserverinfo struct {
	Status  bool
	Players int
}

func readBtServerApiKey() string {
	key, err := ioutil.ReadFile("btServerAzureFunction.key")
	if err != nil {
		panic(err)
	}
	return string(key)
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

func btStop() {
	vs := vssh.New().Start()
	config, _ := vssh.GetConfigPEM("azureuser", "btServer_key.pem")
	vs.AddClient("baro.joe.surf:22", config, vssh.SetMaxSessions(4))
	vs.Wait()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cmd := "./btserver-2 stop"
	timeout, _ := time.ParseDuration("6s")
	respChan := vs.Run(ctx, cmd, timeout)

	resp := <-respChan
	if err := resp.Err(); err != nil {
		log.Fatal(err)
	}

	stream := resp.GetStream()
	defer stream.Close()

	for stream.ScanStdout() {
		txt := stream.TextStdout()
		fmt.Println(txt)
	}
}

func btBackup() {
	vs := vssh.New().Start()
	config, _ := vssh.GetConfigPEM("azureuser", "btServer_key.pem")
	vs.AddClient("baro.joe.surf:22", config, vssh.SetMaxSessions(4))
	vs.Wait()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cmd := "./btserver-2 backup"
	timeout, _ := time.ParseDuration("30s")
	respChan := vs.Run(ctx, cmd, timeout)

	resp := <-respChan
	if err := resp.Err(); err != nil {
		log.Fatal(err)
	}

	stream := resp.GetStream()
	defer stream.Close()

	for stream.ScanStdout() {
		txt := stream.TextStdout()
		fmt.Println(txt)
	}
}

func btServerShutdown() {
	vs := vssh.New().Start()
	config, _ := vssh.GetConfigPEM("azureuser", "btServer_key.pem")
	vs.AddClient("baro.joe.surf:22", config, vssh.SetMaxSessions(4))
	vs.Wait()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cmd := "sudo shutdown now"
	timeout, _ := time.ParseDuration("6s")
	respChan := vs.Run(ctx, cmd, timeout)

	resp := <-respChan
	if err := resp.Err(); err != nil {
		log.Fatal(err)
	}

	stream := resp.GetStream()
	defer stream.Close()

	for stream.ScanStdout() {
		txt := stream.TextStdout()
		fmt.Println(txt)
	}
}

func btShutdown() {
	time.Sleep(time.Second * 5)
	btStop()
	time.Sleep(time.Second * 10)
	btServerShutdown()
}

func btServerStart() {
	apiKey := readBtServerApiKey()
	url := "https://fapp-joebot.azurewebsites.net/start/" + apiKey
	client := &http.Client{Timeout: 15 * time.Second}
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(responseData)
}
