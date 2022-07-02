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

type SurflineWaveForecast struct {
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

type SurflineWindForecast struct {
	Associated struct {
		Units struct {
			Temperature string `json:"temperature"`
			TideHeight  string `json:"tideHeight"`
			SwellHeight string `json:"swellHeight"`
			WaveHeight  string `json:"waveHeight"`
			WindSpeed   string `json:"windSpeed"`
		} `json:"units"`
		UtcOffset int `json:"utcOffset"`
		Location  struct {
			Lon float64 `json:"lon"`
			Lat float64 `json:"lat"`
		} `json:"location"`
	} `json:"associated"`
	Data struct {
		Wind []struct {
			Timestamp     int     `json:"timestamp"`
			UtcOffset     int     `json:"utcOffset"`
			Speed         float64 `json:"speed"`
			Direction     float64 `json:"direction"`
			DirectionType string  `json:"directionType"`
			Gust          float64 `json:"gust"`
			OptimalScore  int     `json:"optimalScore"`
		} `json:"wind"`
	} `json:"data"`
}

type BuoyData struct {
	WaveHeight  float64
	SwellPeriod string
}

const JaxBeachPierSpotID = "5842041f4e65fad6a7708aa0"

func mphFromKnots(knots float64) float64 {
	return knots * 1.151
}

func convertDirection(bearing float64) string {
	b := int(bearing)
	if(b<24) {
		return "N"
		fmt.Println(b)
	}
	if(b<69) {
		return "NE"
	}
	if(b<114){
		return "E"
	}
	if(b<159){
		return "SE"
	}
	if(b<204){
		return "S"
	}
	if(b<249){
		return "SW"
	}
	if(b<294){
		return "W"
	}
	if(b<337){
		return "NW"
	}
	if(b<361){
		return "N"
	}
	return "?"
}

func parseForecast(wave SurflineWaveForecast, wind SurflineWindForecast) string {
	bd := getBuoyData()
	responseText := "```------------------------------------------------------------------------\n"
	responseText += "Current Buoy Data:\n"
	responseText += "Wave Height: " + fmt.Sprint(bd.WaveHeight) + " ft"
	responseText += "\nSwell Period: " + bd.SwellPeriod + " seconds"
	responseText += "\n---------------------------------------------------------------------------------------\n"
	responseText += "Surf Forecast:\n"
	for i, conditions := range wave.Data.Conditions {
		//ensure there is a response
		if conditions.Am.Rating != "" {
			//surfline conditions forecast comes in days, wind comes in hours 
			x := (i*24)+7
			y := (i*24)+14
			amWindLow := math.Round(mphFromKnots(wind.Data.Wind[x].Speed))
			amWindHigh := math.Round(mphFromKnots(wind.Data.Wind[x].Gust))
			amWindDir := convertDirection(wind.Data.Wind[x].Direction)
			pmWindLow := math.Round(mphFromKnots(wind.Data.Wind[y].Speed))
			pmWindHigh := math.Round(mphFromKnots(wind.Data.Wind[y].Gust))
			pmWindDir := convertDirection(wind.Data.Wind[y].Direction)
			responseText += conditions.ForecastDay + "\n\n" +
				fmt.Sprint(conditions.Observation,"\n\n")+
				"am: " + fmt.Sprint(conditions.Am.MinHeight) +
				"-" + fmt.Sprint(conditions.Am.MaxHeight) + "ft  " +
				fmt.Sprint(conditions.Am.HumanRelation) +
				"\twind: " + fmt.Sprint(amWindLow, "-", amWindHigh) + "mph " +
				fmt.Sprint(amWindDir,"  ") +
				"\t\trating: " + conditions.Am.Rating +
				"\npm: " + fmt.Sprint(conditions.Pm.MinHeight) +
				"-" + fmt.Sprint(conditions.Pm.MaxHeight) + "ft  " +
				fmt.Sprint(conditions.Pm.HumanRelation) +
				"\twind: " + fmt.Sprint(pmWindLow, "-", pmWindHigh) + "mph " +
				fmt.Sprint(pmWindDir,"  ") +
				"\t\trating: " + conditions.Pm.Rating + "\n" +
				"------------------------------------------------------------------------------------------\n"
		}
	}
	responseText += "```"
	return responseText
}

func getSurflineWaveForecast() SurflineWaveForecast {
	url := "https://services.surfline.com/kbyg/regions/forecasts/conditions?subregionId=5e556e9231e571b1a21d34a0"
	client := &http.Client{Timeout: 10 * time.Second}
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
	var sf SurflineWaveForecast
	err = json.Unmarshal(responseData, &sf)
	if err != nil {
		fmt.Println(err)
	}
	return sf
}

func getSurflineWindForecast() SurflineWindForecast {
	url := "https://services.surfline.com/kbyg/spots/forecasts/wind?spotId=" + JaxBeachPierSpotID
	fmt.Println(url)
	client := &http.Client{Timeout: 10 * time.Second}
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
	var sf SurflineWindForecast
	err = json.Unmarshal(responseData, &sf)
	if err != nil {
		fmt.Println(err)
	}
	return sf
}

func getBuoyData() BuoyData {
	url := "https://www.ndbc.noaa.gov/data/realtime2/41112.txt"
	client := &http.Client{Timeout: 10 * time.Second}
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

func main() {
	waveForecastData := getSurflineWaveForecast()
	windForecastData := getSurflineWindForecast()
	parsedResult := parseForecast(waveForecastData, windForecastData)
	fmt.Println(parsedResult)
}
