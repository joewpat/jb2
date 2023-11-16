package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Forecast structure to hold forecast data
type Forecast struct {
	Properties struct {
		Periods []struct {
			Name             string `json:"name"`
			Temperature      int    `json:"temperature"`
			WindSpeed        string `json:"windSpeed"`
			DetailedForecast string `json:"detailedForecast"`
			// Add more fields as needed
		} `json:"periods"`
	} `json:"properties"`
}

func weatherForecast() string {

	// JAX forecast endpoint
	const apiEndpoint = "https://api.weather.gov/gridpoints/JAX/66,65/forecast"
	// Make the HTTP request
	response, err := http.Get(apiEndpoint)
	if err != nil {
		fmt.Println("Error making the request:", err)
		return "error with noaa api"
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading the response body:", err)
		return "error with noaa api"
	}

	// Unmarshal JSON into a Forecast struct
	var forecast Forecast
	err = json.Unmarshal(body, &forecast)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return "error with forecast data"
	}

	//compile daily forecast into string response
	var completeForecast string

	completeForecast += string(("Forecast for" + forecast.Properties.Periods[0].Name + ":\n"))
	completeForecast += forecast.Properties.Periods[0].DetailedForecast
	completeForecast += string(("\nForecast for" + forecast.Properties.Periods[1].Name + ":\n"))
	completeForecast += forecast.Properties.Periods[1].DetailedForecast + "\n"

	return completeForecast
}
