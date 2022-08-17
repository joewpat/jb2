package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"sort"
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
	if b < 24 {
		return "N"
		//fmt.Println(b)
	}
	if b < 69 {
		return "NE"
	}
	if b < 114 {
		return "E"
	}
	if b < 159 {
		return "SE"
	}
	if b < 204 {
		return "S"
	}
	if b < 249 {
		return "SW"
	}
	if b < 294 {
		return "W"
	}
	if b < 337 {
		return "NW"
	}
	if b < 361 {
		return "N"
	}
	return "?"
}

//takes data from APIs and formats a response
func parseForecast(wave SurflineWaveForecast, wind SurflineWindForecast, tide TideAndWeatherChart) string {
	bd := getBuoyData()
	responseText := "```------------------------------------------------------------------------\n"
	responseText += "Current Buoy 411112 Data\n"
	responseText += "Wave Height: " + fmt.Sprint(bd.WaveHeight) + " ft"
	responseText += "\nSwell Period: " + bd.SwellPeriod + " seconds"
	responseText += "\nWind: " + fmt.Sprint(math.Round(mphFromKnots(wind.Data.Wind[0].Speed))) + "mph " + convertDirection(wind.Data.Wind[0].Direction)
	responseText += "\n---------------------------------------------------------------------------------------\n"
	//Surf Forecast from Surfline
	for i, conditions := range wave.Data.Conditions {
		//ensure there is a response
		if conditions.Am.Rating != "" {
			if i < 3 {
				//surfline conditions forecast comes in days, wind and tide come in hours
				x := (i * 24) + 7
				y := (i * 24) + 14
				//	z := i*24
				//	p := (2*i)
				//	q := (2*i)
				amWindLow := math.Round(mphFromKnots(wind.Data.Wind[x].Speed))
				amWindHigh := math.Round(mphFromKnots(wind.Data.Wind[x].Gust))
				amWindDir := convertDirection(wind.Data.Wind[x].Direction)
				pmWindLow := math.Round(mphFromKnots(wind.Data.Wind[y].Speed))
				pmWindHigh := math.Round(mphFromKnots(wind.Data.Wind[y].Gust))
				pmWindDir := convertDirection(wind.Data.Wind[y].Direction)
				//filter out non gusty wind so it doesn't display like 7-7mph
				var amWindReport string
				if amWindHigh == amWindLow {
					amWindReport = fmt.Sprint(amWindHigh, "mph ", amWindDir)
				} else {
					amWindReport = fmt.Sprint(amWindLow, "-", amWindHigh, "mph ", amWindDir)
				}
				var pmWindReport string
				if pmWindHigh == pmWindLow {
					pmWindReport = fmt.Sprint(pmWindHigh, "mph ", pmWindDir)
				} else {
					pmWindReport = fmt.Sprint(pmWindLow, "-", pmWindHigh, "mph ", pmWindDir)
				}
				forecastTime, err := time.Parse("2006-01-02", conditions.ForecastDay)
				if err != nil {
					fmt.Println("Could not parse time:", err)
				}
				forecastDayOfWeek := fmt.Sprint(forecastTime.Weekday())

				amPrimarySwells := tide.Data.Forecasts[2*i].Swells
				pmPrimarySwells := tide.Data.Forecasts[5*i].Swells
				//sort by biggest swell to report primary swell/period
				sort.Slice(amPrimarySwells, func(i, j int) bool {
					return amPrimarySwells[i].Height > amPrimarySwells[j].Height
				})
				sort.Slice(pmPrimarySwells, func(i, j int) bool {
					return pmPrimarySwells[i].Height > pmPrimarySwells[j].Height
				})

				//build response
				responseText += forecastDayOfWeek + "\t\t" + conditions.ForecastDay[5:] + "\n\n" +
					fmt.Sprint(conditions.Observation, "\n\n") +
					"am: " + fmt.Sprint(conditions.Am.MinHeight) +
					"-" + fmt.Sprint(conditions.Am.MaxHeight) + "ft  " +
					fmt.Sprint(conditions.Am.HumanRelation) +
					"\tprimary swell: " + fmt.Sprint(math.Floor(amPrimarySwells[0].Height*10)/10, "ft@", amPrimarySwells[0].Period, "sec") +
					"\twind: " + amWindReport +
					"\t\trating: " + conditions.Am.Rating +
					"\npm: " + fmt.Sprint(conditions.Pm.MinHeight) +
					"-" + fmt.Sprint(conditions.Pm.MaxHeight) + "ft  " +
					fmt.Sprint(conditions.Pm.HumanRelation) +
					"\tprimary swell: " + fmt.Sprint(math.Floor(pmPrimarySwells[0].Height*10)/10, "ft@", pmPrimarySwells[0].Period, "sec") +
					"\twind: " + pmWindReport +
					"\t\trating: " + conditions.Pm.Rating + "\n" +
					"--------------------------------------------------------------------------------------------------------------\n"
			}
		}
	}
	responseText += "```"
	return responseText
}

func getSurflineWaveForecast() SurflineWaveForecast {
	url := "https://services.surfline.com/kbyg/regions/forecasts/conditions?subregionId=5e556e9231e571b1a21d34a0"
	client := &http.Client{Timeout: 20 * time.Second}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
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
	fmt.Println("wave forecast OK")
	return sf
}

func getSurflineWindForecast() SurflineWindForecast {
	url := "https://services.surfline.com/kbyg/spots/forecasts/wind?spotId=" + JaxBeachPierSpotID
	client := &http.Client{Timeout: 20 * time.Second}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
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
		if e, ok := err.(*json.SyntaxError); ok {
			fmt.Printf("syntax error at byte offset %d\n", e.Offset)
		}
		fmt.Printf("sakura response: %q\n", responseData)
		fmt.Println(err)
	}
	fmt.Println("wind forecast OK")
	return sf
}

type TideAndWeatherChart struct {
	Data struct {
		SunriseSunsetTimes []struct {
			Midnight int `json:"midnight"`
			Sunrise  int `json:"sunrise"`
			Sunset   int `json:"sunset"`
		} `json:"sunriseSunsetTimes"`
		TideLocation struct {
			Name string  `json:"name"`
			Min  float64 `json:"min"`
			Max  float64 `json:"max"`
			Lon  float64 `json:"lon"`
			Lat  float64 `json:"lat"`
			Mean int     `json:"mean"`
		} `json:"tideLocation"`
		Forecasts []struct {
			Timestamp int `json:"timestamp"`
			Weather   struct {
				Temperature int    `json:"temperature"`
				Condition   string `json:"condition"`
			} `json:"weather"`
			Wind struct {
				Speed     float64 `json:"speed"`
				Direction float64 `json:"direction"`
			} `json:"wind"`
			Surf struct {
				Min float64 `json:"min"`
				Max float64 `json:"max"`
			} `json:"surf"`
			Swells []struct {
				Height       float64 `json:"height"`
				Direction    float64 `json:"direction"`
				DirectionMin float64 `json:"directionMin"`
				Period       int     `json:"period"`
			} `json:"swells"`
		} `json:"forecasts"`
		Tides []struct {
			Timestamp int     `json:"timestamp"`
			Type      string  `json:"type"`
			Height    float64 `json:"height"`
		} `json:"tides"`
	} `json:"data"`
}

func tideAndWeather() TideAndWeatherChart {
	url := "https://services.surfline.com/kbyg/spots/forecasts/?spotId=5842041f4e65fad6a7708aa0"
	client := &http.Client{Timeout: 20 * time.Second}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	var weather TideAndWeatherChart
	err = json.Unmarshal(responseData, &weather)
	if err != nil {
		fmt.Println(err)
	}
	return weather
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
	whint, _ := strconv.ParseFloat(wh, 64)
	ft := whint * 3.28084
	whft := math.Round(ft*100) / 100
	report := BuoyData{
		WaveHeight:  whft,
		SwellPeriod: data[20],
	}
	return report
}

func getSurflineForecast() string {
	waveForecastData := getSurflineWaveForecast()
	windForecastData := getSurflineWindForecast()
	tideAndWeatherData := tideAndWeather()
	parsedResult := parseForecast(waveForecastData, windForecastData, tideAndWeatherData)

	return parsedResult
}
