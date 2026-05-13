package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

type youtubeSearchResponse struct {
	Items []struct {
		ID struct {
			VideoID string `json:"videoId"`
		} `json:"id"`
	} `json:"items"`
}

type youtubeCommentsResponse struct {
	Items []struct {
		Snippet struct {
			TopLevelComment struct {
				Snippet struct {
					TextDisplay string `json:"textDisplay"`
				} `json:"snippet"`
			} `json:"topLevelComment"`
		} `json:"snippet"`
	} `json:"items"`
}

// SearchYouTubeRandomComment searches YouTube for a query and returns a random comment from a random matching video.
func SearchYouTubeRandomComment(query string) string {
	apiKey := os.Getenv("YOUTUBE_API_KEY")
	if apiKey == "" {
		return "YouTube API key not configured."
	}

	// Check if the API key looks valid (should start with AIza)
	if !strings.HasPrefix(apiKey, "AIza") {
		fmt.Printf("Warning: YouTube API key may be invalid (should start with AIza)\n")
	}

	if query == "" {
		return "I need something to search for."
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	client := &http.Client{Timeout: 15 * time.Second}

	videoIDs, err := fetchYouTubeVideoIDs(client, apiKey, query)
	if err != nil {
		fmt.Printf("YouTube API error: %v\n", err)
		return "Couldn't search YouTube right now."
	}
	if len(videoIDs) == 0 {
		return "No YouTube videos found for that topic."
	}

	start := rng.Intn(len(videoIDs))
	for i := 0; i < len(videoIDs); i++ {
		videoID := videoIDs[(start+i)%len(videoIDs)]
		comments, err := fetchYouTubeComments(client, apiKey, videoID)
		if err != nil || len(comments) == 0 {
			continue
		}
		return comments[rng.Intn(len(comments))]
	}

	return "No comments found in those videos."
}

func fetchYouTubeVideoIDs(client *http.Client, apiKey, query string) ([]string, error) {
	searchURL := fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?part=id&type=video&maxResults=25&q=%s&key=%s",
		url.QueryEscape(query), apiKey)

	req, err := http.NewRequest(http.MethodGet, searchURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Read error body for more details
		var errorBody map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorBody)
		return nil, fmt.Errorf("youtube search failed: %s - %v", resp.Status, errorBody)
	}

	var searchResp youtubeSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, err
	}

	videoIDs := make([]string, 0, len(searchResp.Items))
	for _, item := range searchResp.Items {
		if item.ID.VideoID != "" {
			videoIDs = append(videoIDs, item.ID.VideoID)
		}
	}
	return videoIDs, nil
}

func fetchYouTubeComments(client *http.Client, apiKey, videoID string) ([]string, error) {
	commentsURL := fmt.Sprintf("https://www.googleapis.com/youtube/v3/commentThreads?part=snippet&videoId=%s&maxResults=100&key=%s",
		videoID, apiKey)

	req, err := http.NewRequest(http.MethodGet, commentsURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("youtube comments failed: %s", resp.Status)
	}

	var commentsResp youtubeCommentsResponse
	if err := json.NewDecoder(resp.Body).Decode(&commentsResp); err != nil {
		return nil, err
	}

	var comments []string
	for _, item := range commentsResp.Items {
		comment := item.Snippet.TopLevelComment.Snippet.TextDisplay
		if comment != "" && !containsYouTubeHyperlink(comment) {
			// Strip HTML tags
			comment = stripHTMLTags(comment)
			// Strip HTML entities
			comment = strings.ReplaceAll(comment, "&quot;", "\"")
			comment = strings.ReplaceAll(comment, "&#39;", "'")
			comment = strings.ReplaceAll(comment, "&amp;", "&")
			comment = strings.ReplaceAll(comment, "&lt;", "<")
			comment = strings.ReplaceAll(comment, "&gt;", ">")
			comments = append(comments, comment)
		}
	}
	return comments, nil
}

// stripHTMLTags removes HTML tags from a string
func stripHTMLTags(s string) string {
	// Remove HTML tags using regex
	re := regexp.MustCompile(`<[^>]*>`)
	return re.ReplaceAllString(s, "")
}

// containsYouTubeHyperlink checks if a comment contains URLs or hyperlinks
func containsYouTubeHyperlink(comment string) bool {
	lower := strings.ToLower(comment)
	// Check for common URL patterns
	if strings.Contains(lower, "http://") || strings.Contains(lower, "https://") ||
		strings.Contains(lower, "www.") || strings.Contains(comment, "](") ||
		strings.Contains(lower, ".com") || strings.Contains(lower, ".org") ||
		strings.Contains(lower, ".net") {
		return true
	}
	return false
}
