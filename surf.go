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

type BuoyData struct {
	WaveHeight  float64
	SwellPeriod string
	WaterTemp   float64
	WindSpeed   float64
	WindDir     string
}
type TidePrediction struct {
	Time  string `json:"t"`
	Level string `json:"type"`
}

type TideData struct {
	Predictions []TidePrediction `json:"predictions"`
}

func getBuoyData() BuoyData {
	url := "https://www.ndbc.noaa.gov/data/realtime2/41112.txt"
	client := &http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error fetching buoy data:", err)
		return BuoyData{}
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading buoy response:", err)
		return BuoyData{}
	}

	responseString := string(responseData)
	lines := strings.Split(responseString, "\n")
	if len(lines) < 3 {
		fmt.Println("Error: Buoy data file does not contain enough lines")
		return BuoyData{}
	}

	latestReading := lines[2]
	data := strings.Fields(latestReading) // Use Fields to handle multiple spaces

	// Ensure the data slice has enough elements
	if len(data) < 17 {
		fmt.Println("Error: Buoy data line does not contain enough fields")
		return BuoyData{}
	}

	// Parse wave height (WVHT)
	wvht := data[8]
	wvhtFloat, err := strconv.ParseFloat(wvht, 64)
	if err != nil {
		fmt.Println("Error parsing wave height:", err)
		return BuoyData{}
	}
	waveHeightFt := math.Round(wvhtFloat*3.28084*100) / 100 // Convert meters to feet

	// Parse swell period (DPD)
	swellPeriod := data[9]

	// Parse water temperature (WTMP) and convert to Fahrenheit
	wtmp := data[14]
	wtmpFloat, err := strconv.ParseFloat(wtmp, 64)
	if err != nil {
		fmt.Println("Error parsing water temperature:", err)
		return BuoyData{}
	}
	waterTempF := math.Round(wtmpFloat*1.8 + 32) // Convert Celsius to Fahrenheit and round to nearest integer

	// Return the parsed data
	report := BuoyData{
		WaveHeight:  waveHeightFt,
		SwellPeriod: swellPeriod,
		WaterTemp:   waterTempF,
	}
	return report
}

func getWindData() (float64, string) {
	url := "https://api.weather.gov/gridpoints/JAX/57,74/forecast"
	client := &http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Go-http-client/1.1") // Required by the NWS API
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error fetching wind data:", err)
		return 0, "N/A"
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading wind response:", err)
		return 0, "N/A"
	}

	var forecast map[string]interface{}
	err = json.Unmarshal(responseData, &forecast)
	if err != nil {
		fmt.Println("Error parsing wind data:", err)
		return 0, "N/A"
	}

	// Extract wind data from the forecast
	properties := forecast["properties"].(map[string]interface{})
	periods := properties["periods"].([]interface{})
	firstPeriod := periods[0].(map[string]interface{})
	windSpeed := firstPeriod["windSpeed"].(string) // Example: "10 mph"
	windDirection := firstPeriod["windDirection"].(string)

	// Parse wind speed (strip " mph" and convert to float)
	windSpeedValue, err := strconv.ParseFloat(strings.Split(windSpeed, " ")[0], 64)
	if err != nil {
		fmt.Println("Error parsing wind speed:", err)
		return 0, "N/A"
	}

	return windSpeedValue, windDirection
}

// func convertDirection(bearing string) string {
// 	b, err := strconv.Atoi(bearing)
// 	if err != nil {
// 		return "?"
// 	}
// 	switch {
// 	case float64(b) >= 348.75 || float64(b) < 11.25:
// 		return "N"
// 	case float64(b) >= 11.25 && float64(b) < 33.75:
// 		return "NNE"
// 	case float64(b) >= 33.75 && float64(b) < 56.25:
// 		return "NE"
// 	case float64(b) >= 56.25 && float64(b) < 78.75:
// 		return "ENE"
// 	case float64(b) >= 78.75 && float64(b) < 101.25:
// 		return "E"
// 	case float64(b) >= 101.25 && float64(b) < 123.75:
// 		return "ESE"
// 	case float64(b) >= 123.75 && float64(b) < 146.25:
// 		return "SE"
// 	case float64(b) >= 146.25 && float64(b) < 168.75:
// 		return "SSE"
// 	case float64(b) >= 168.75 && float64(b) < 191.25:
// 		return "S"
// 	case float64(b) >= 191.25 && float64(b) < 213.75:
// 		return "SSW"
// 	case float64(b) >= 213.75 && float64(b) < 236.25:
// 		return "SW"
// 	case float64(b) >= 236.25 && float64(b) < 258.75:
// 		return "WSW"
// 	case float64(b) >= 258.75 && float64(b) < 281.25:
// 		return "W"
// 	case float64(b) >= 281.25 && float64(b) < 303.75:
// 		return "WNW"
// 	case float64(b) >= 303.75 && float64(b) < 326.25:
// 		return "NW"
// 	case float64(b) >= 326.25 && float64(b) < 348.75:
// 		return "NNW"
// 	default:
// 		return "?"
// 	}
// }

func getTideData() []TidePrediction {
	stationID := "8720291" // Jacksonville Beach, FL
	now := time.Now()
	beginDate := now.Format("20060102")
	endDate := now.Format("20060102")

	url := fmt.Sprintf("https://api.tidesandcurrents.noaa.gov/api/prod/datagetter?product=predictions&application=NOS.COOPS.TAC.WL&begin_date=%s&end_date=%s&datum=MLLW&station=%s&time_zone=lst_ldt&units=english&interval=hilo&format=json", beginDate, endDate, stationID)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("Error fetching tide data:", err)
		return nil
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading tide response:", err)
		return nil
	}

	var tideData TideData
	err = json.Unmarshal(responseData, &tideData)
	if err != nil {
		fmt.Println("Error parsing tide data:", err)
		return nil
	}

	return tideData.Predictions
}

func getSurfData() string {
	bd := getBuoyData()
	tides := getTideData()
	windSpeed, windDir := getWindData()

	responseText := "```---\n"
	responseText += "Current Surf Conditions:\n"
	responseText += "Wave Height: " + fmt.Sprintf("%.2f", bd.WaveHeight) + " ft"
	responseText += "\nSwell Period: " + bd.SwellPeriod + " seconds"
	responseText += "\nWater Temp: " + fmt.Sprintf("%.0f", bd.WaterTemp) + " °F"
	responseText += "\nWind Speed: " + fmt.Sprintf("%.0f", windSpeed) + " mph"
	responseText += "\nWind Direction: " + windDir
	responseText += "\nTides:\n"

	for _, tide := range tides {
		tideTime, err := time.Parse("2006-01-02 15:04", tide.Time)
		if err != nil {
			fmt.Println("Error parsing tide time:", err)
			continue
		}
		responseText += fmt.Sprintf("%s Tide at %s\n", strings.Title(tide.Level), tideTime.Format("03:04 PM"))
	}

	responseText += "```"
	return responseText
}
