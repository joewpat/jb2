package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type Gif struct {
	Results []struct {
		Media []struct {
			Gif struct {
				URL string `json:"url"`
			} `json:"gif"`
		} `json:"media"`
		URL string `json:"url"`
	} `json:"results"`
}

type Giphy struct {
	Data []struct {
		Type string `json:"type"`
		ID   string `json:"id"`
		URL  string `json:"url"`
	} `json:"data"`
}

func readTenorKey() string {
	key, err := ioutil.ReadFile("tenor.key")
	if err != nil {
		panic(err)
	}
	return string(key)
}

func readGiphyKey() string {
	key, err := ioutil.ReadFile("giphy.key")
	if err != nil {
		panic(err)
	}
	return string(key)
}

func searchGiphy(query string) string {
	giphyKey := readGiphyKey()
	url := "http://api.giphy.com/v1/gifs/search?q=" + query + "&api_key=" + giphyKey + "&limit=50"
	fmt.Println("searching giphy: ", url)
	client := &http.Client{Timeout: 3 * time.Second}
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	fmt.Println("giphy API status code: ", resp.StatusCode)

	if resp.StatusCode == 200 {
		// Transform our response to a []byte
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		}
		var data Giphy
		err = json.Unmarshal(body, &data)
		if err != nil {
			fmt.Println(err)
		}
		rand.Seed(time.Now().Unix()) // initialize pseudo random generator
		fmt.Println("giphy - number of results: ", len(data.Data))
		if len(data.Data) < 1 {
			fmt.Println("giphy returned 0 search results")
			return "error - giphy returned 0 search results"
		}
		randResult := data.Data[rand.Intn(len(data.Data))].URL
		return randResult
	} else {
		fmt.Println("error searching giphy")
		return "error searching giphy"
	}
}

func searchTenor(query string) string {
	tenorKey := readTenorKey()
	url := "https://g.tenor.com/v1/search?q=" + query + "&key=" + tenorKey + "&limit=50"
	fmt.Println("searching tenor: ", url)
	client := &http.Client{Timeout: 3 * time.Second}
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	fmt.Println("tenor API status code: ", resp.StatusCode)

	if resp.StatusCode == 200 {
		// Transform our response to a []byte
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		}
		var data Gif
		err = json.Unmarshal(body, &data)
		if err != nil {
			fmt.Println(err)
		}
		rand.Seed(time.Now().Unix()) // initialize pseudo random generator
		fmt.Println("tenor - number of results: ", len(data.Results))
		if len(data.Results) < 1 {
			fmt.Println("tenor returned 0 search results")
			return "error - tenor returned 0 search results"
		}
		randResult := data.Results[rand.Intn(len(data.Results))].Media[0].Gif.URL
		return randResult
	} else {
		fmt.Println("error searching Tenor")
		return "error searching Tenor"
	}
}

func searchGifs(t string) string {
	tenor := searchTenor(t)
	giphy := searchGiphy(t)
	var responses []string
	if strings.HasPrefix(tenor, "http") {
		responses = append(responses, tenor)
		fmt.Println("response from tenor: ", tenor)
	}
	if strings.HasPrefix(giphy, "http") {
		responses = append(responses, giphy)
		fmt.Println("response from giphy: ", giphy)
	}
	if len(responses) > 0 {
		return responses[rand.Intn(len(responses))]
	}
	return "error - no gifs found"
}
