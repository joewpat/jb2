package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type surfReport struct {
	WaveHeight  float64
	SwellPeriod string
	//swellDirection string
}

func getSurfReport() string {
	url := "https://www.ndbc.noaa.gov/data/realtime2/41112.txt"
	client := &http.Client{Timeout: 3 * time.Second}
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

	responseString := string(responseData)
	lines := strings.Split(responseString, "\n")
	latestReading := lines[2]
	//fmt.Println(latestReading)
	data := strings.Split(latestReading, " ")

	wh := data[15]
	whint, err := strconv.ParseFloat(wh, 64)
	ft := whint * 3.28084
	whft := math.Round(ft*100) / 100
	report := surfReport{
		WaveHeight:  whft,
		SwellPeriod: data[20],
	}

	response := fmt.Sprintf("Wave Height: %v"+"ft"+"\nSwell Period: %v"+" seconds", report.WaveHeight, report.SwellPeriod)
	return response

}
