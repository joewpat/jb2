package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

/*
GET https://youtube.googleapis.com/youtube/v3/search?part=snippet&maxResults=25&q=surfing&key=[YOUR_API_KEY] HTTP/1.1
Authorization: Bearer [YOUR_ACCESS_TOKEN]
Accept: application/json
*/

type YoutubeResponse struct {
	Items []struct {
		Kind string `json:"kind"`
		Etag string `json:"etag"`
		ID   struct {
			Kind    string `json:"kind"`
			VideoID string `json:"videoId"`
		} `json:"id"`
		Snippet struct {
			PublishedAt time.Time `json:"publishedAt"`
			ChannelID   string    `json:"channelId"`
			Title       string    `json:"title"`
			Description string    `json:"description"`
			Thumbnails  struct {
				Default struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"default"`
			} `json:"snippet"`
		} `json:"items"`
	}
}

func youtube() {
	yt := searchYT("butt")
	for _, y := range yt.Items {
		fmt.Println(y.ID.VideoID)
	}
}

func searchYT(query string) YoutubeResponse {
	yt_api_key := "AIzaSyDdTaCfsKc7tfpNMMzP4whQ-BipKN0SVI0"
	url := "https://youtube.googleapis.com/youtube/v3/search?part=snippet&maxResults=25&q=" + query + "&key=" + yt_api_key
	client := &http.Client{Timeout: 3 * time.Second}
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	var y YoutubeResponse
	err = json.NewDecoder(resp.Body).Decode(&y)
	if err != nil {
		panic(err)
	}
	return y
}
