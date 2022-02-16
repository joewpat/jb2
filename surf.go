package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type SurflineForecast struct {
	Data struct {
		Conditions []struct {
			Timestamp   int    `json:"timestamp"`
			ForecastDay string `json:"forecastDay"`
			Forecaster  struct {
				Name   string `json:"name"`
				Avatar string `json:"avatar"`
			} `json:"forecaster"`
			Human       bool   `json:"human"`
			Observation string `json:"observation"`
			Am          struct {
				MaxHeight        int         `json:"maxHeight"`
				MinHeight        int         `json:"minHeight"`
				Plus             bool        `json:"plus"`
				HumanRelation    string      `json:"humanRelation"`
				OccasionalHeight interface{} `json:"occasionalHeight"`
				Rating           string      `json:"rating"`
			} `json:"am"`
			Pm struct {
				MaxHeight        int         `json:"maxHeight"`
				MinHeight        int         `json:"minHeight"`
				Plus             bool        `json:"plus"`
				HumanRelation    string      `json:"humanRelation"`
				OccasionalHeight interface{} `json:"occasionalHeight"`
				Rating           string      `json:"rating"`
			} `json:"pm"`
		} `json:"conditions"`
	} `json:"data"`
}

type BuoyData struct {
	WaveHeight  float64
	SwellPeriod string
}

func parseForecast(sf SurflineForecast) string {
	bd := getBuoyData()
	responseText := "```Current Buoy Data:\n"
	responseText += "Wave Height: " + fmt.Sprint(bd.WaveHeight)
	responseText += "\nSwell Period: " + bd.SwellPeriod
	responseText += "\n\nSurf Forecast:\n"
	for _, conditions := range sf.Data.Conditions {
		if conditions.Am.Rating != "" {
			responseText += conditions.ForecastDay + "\n" +
				"am: " + fmt.Sprint(conditions.Am.MinHeight) +
				"-" + fmt.Sprint(conditions.Am.MaxHeight) + "ft " +
				"\trating: " + conditions.Am.Rating + "\n" +
				"pm: " + fmt.Sprint(conditions.Pm.MinHeight) +
				"-" + fmt.Sprint(conditions.Pm.MaxHeight) + "ft " +
				"\trating: " + conditions.Pm.Rating + "\n"
		}
	}
	responseText += "```"
	return responseText
}

func getSurflineForecast() SurflineForecast {
	url := "https://services.surfline.com/kbyg/regions/forecasts/conditions?subregionId=5e556e9231e571b1a21d34a0"
	client := &http.Client{Timeout: 5 * time.Second}
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	var sf SurflineForecast
	err = json.Unmarshal(responseData, &sf)
	if err != nil {
		fmt.Println(err)
	}
	return sf
}

func getBuoyData() BuoyData {
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
	data := strings.Split(latestReading, " ")

	wh := data[15]
	whint, err := strconv.ParseFloat(wh, 64)
	ft := whint * 3.28084
	whft := math.Round(ft*100) / 100
	report := BuoyData{
		WaveHeight:  whft,
		SwellPeriod: data[20],
	}
	return report
}
