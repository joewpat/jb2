package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

type YoutubeCommentThread struct {
	Kind          string `json:"kind"`
	Etag          string `json:"etag"`
	NextPageToken string `json:"nextPageToken"`
	PageInfo      struct {
		TotalResults   int `json:"totalResults"`
		ResultsPerPage int `json:"resultsPerPage"`
	} `json:"pageInfo"`
	Items []struct {
		Kind    string `json:"kind"`
		Etag    string `json:"etag"`
		ID      string `json:"id"`
		Snippet struct {
			VideoID         string `json:"videoId"`
			TopLevelComment struct {
				Kind    string `json:"kind"`
				Etag    string `json:"etag"`
				ID      string `json:"id"`
				Snippet struct {
					VideoID               string `json:"videoId"`
					TextDisplay           string `json:"textDisplay"`
					TextOriginal          string `json:"textOriginal"`
					AuthorDisplayName     string `json:"authorDisplayName"`
					AuthorProfileImageURL string `json:"authorProfileImageUrl"`
					AuthorChannelURL      string `json:"authorChannelUrl"`
					AuthorChannelID       struct {
						Value string `json:"value"`
					} `json:"authorChannelId"`
					CanRate      bool      `json:"canRate"`
					ViewerRating string    `json:"viewerRating"`
					LikeCount    int       `json:"likeCount"`
					PublishedAt  time.Time `json:"publishedAt"`
					UpdatedAt    time.Time `json:"updatedAt"`
				} `json:"snippet"`
			} `json:"topLevelComment"`
			CanReply        bool `json:"canReply"`
			TotalReplyCount int  `json:"totalReplyCount"`
			IsPublic        bool `json:"isPublic"`
		} `json:"snippet"`
		Replies struct {
			Comments []struct {
				Kind    string `json:"kind"`
				Etag    string `json:"etag"`
				ID      string `json:"id"`
				Snippet struct {
					VideoID               string `json:"videoId"`
					TextDisplay           string `json:"textDisplay"`
					TextOriginal          string `json:"textOriginal"`
					ParentID              string `json:"parentId"`
					AuthorDisplayName     string `json:"authorDisplayName"`
					AuthorProfileImageURL string `json:"authorProfileImageUrl"`
					AuthorChannelURL      string `json:"authorChannelUrl"`
					AuthorChannelID       struct {
						Value string `json:"value"`
					} `json:"authorChannelId"`
					CanRate      bool      `json:"canRate"`
					ViewerRating string    `json:"viewerRating"`
					LikeCount    int       `json:"likeCount"`
					PublishedAt  time.Time `json:"publishedAt"`
					UpdatedAt    time.Time `json:"updatedAt"`
				} `json:"snippet"`
			} `json:"comments"`
		} `json:"replies,omitempty"`
	} `json:"items"`
}

func readKey() string {
	key, err := ioutil.ReadFile("youtube.key")
	if err != nil {
		panic(err)
	}
	return string(key)
}

func youtube(query string) string {
	yt := searchYT(query)
	for _, y := range yt.Items {
		comments := getYTcomment(y.ID.VideoID)
		for _, x := range comments.Items {
			return x.Snippet.TopLevelComment.Snippet.TextDisplay
		}
	}
	return query
}

func searchYT(query string) YoutubeResponse {
	key := readKey()
	url := "https://youtube.googleapis.com/youtube/v3/search?part=snippet&maxResults=25&q=" + query + "&key=" + string(key)
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
		fmt.Println(y)
		panic(err)
	}
	return y
}

func getYTcomment(id string) YoutubeCommentThread {
	key := readKey()
	url := "https://youtube.googleapis.com/youtube/v3/commentThreads?part=snippet%2Creplies&videoId=" + id + "&key=" + string(key)
	client := &http.Client{Timeout: 3 * time.Second}
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	var y YoutubeCommentThread
	err = json.NewDecoder(resp.Body).Decode(&y)
	if err != nil {
		panic(err)
	}
	//fmt.Println(y)
	return y
}

func IsJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}
