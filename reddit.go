package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type redditOAuthToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

var (
	cachedRedditToken     string
	redditTokenExpiration time.Time
)

type redditSearchResponse struct {
	Data struct {
		Children []struct {
			Data struct {
				ID string `json:"id"`
			} `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

type redditCommentsResponse []struct {
	Data struct {
		Children []struct {
			Kind string `json:"kind"`
			Data struct {
				Body    string `json:"body"`
				Replies any    `json:"replies"`
			} `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

// SearchRedditRandomComment searches Reddit for a query and returns a random comment from a random matching thread.
func SearchRedditRandomComment(query string) string {
	if query == "" {
		return "I need something to search for."
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	client := &http.Client{Timeout: 15 * time.Second}

	// Get OAuth token
	token, err := getRedditOAuthToken(client)
	if err != nil {
		fmt.Printf("Reddit OAuth error: %v\n", err)
		return "Couldn't authenticate with Reddit."
	}

	// Use oauth.reddit.com with OAuth token
	subredditURL := "https://oauth.reddit.com/r/all/hot?limit=100"
	postIDs, err := fetchPostIDs(client, subredditURL, token)
	if err != nil {
		fmt.Printf("Reddit API error: %v\n", err)
		return "Couldn't reach Reddit right now."
	}
	if len(postIDs) == 0 {
		return "No Reddit threads found."
	}

	start := rng.Intn(len(postIDs))
	for i := 0; i < len(postIDs); i++ {
		id := postIDs[(start+i)%len(postIDs)]
		comments, err := fetchCommentsForPost(client, id, token)
		if err != nil || len(comments) == 0 {
			continue
		}
		return comments[rng.Intn(len(comments))]
	}

	return "No comments found in those threads."
}

// getRedditOAuthToken obtains or returns a cached OAuth access token
func getRedditOAuthToken(client *http.Client) (string, error) {
	// Return cached token if still valid
	if cachedRedditToken != "" && time.Now().Before(redditTokenExpiration) {
		return cachedRedditToken, nil
	}

	clientID := os.Getenv("REDDIT_CLIENT_ID")
	clientSecret := os.Getenv("REDDIT_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		return "", fmt.Errorf("REDDIT_CLIENT_ID and REDDIT_CLIENT_SECRET environment variables required")
	}

	// Request OAuth token
	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest(http.MethodPost, "https://www.reddit.com/api/v1/access_token", strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(clientID, clientSecret)
	req.Header.Set("User-Agent", "discord:jb2-bot:v1.0 (by /u/robert_arctor)")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("oauth token request failed: %s", resp.Status)
	}

	var tokenResp redditOAuthToken
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", err
	}

	// Cache token with expiration
	cachedRedditToken = tokenResp.AccessToken
	redditTokenExpiration = time.Now().Add(time.Duration(tokenResp.ExpiresIn-60) * time.Second) // Refresh 60s early

	return cachedRedditToken, nil
}

func fetchPostIDs(client *http.Client, endpoint string, token string) ([]string, error) {
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	// Use OAuth authorization
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", "discord:jb2-bot:v1.0 (by /u/robert_arctor)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("reddit search failed: %s", resp.Status)
	}

	var searchResp redditSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(searchResp.Data.Children))
	for _, c := range searchResp.Data.Children {
		if c.Data.ID != "" {
			ids = append(ids, c.Data.ID)
		}
	}
	return ids, nil
}

func fetchCommentsForPost(client *http.Client, postID string, token string) ([]string, error) {
	endpoint := fmt.Sprintf("https://oauth.reddit.com/comments/%s?limit=500", postID)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	// Use OAuth authorization
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", "discord:jb2-bot:v1.0 (by /u/robert_arctor)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("reddit comments failed: %s", resp.Status)
	}

	var commentsResp redditCommentsResponse
	if err := json.NewDecoder(resp.Body).Decode(&commentsResp); err != nil {
		return nil, err
	}

	if len(commentsResp) < 2 {
		return nil, nil
	}

	var out []string
	collectComments(commentsResp[1].Data.Children, &out)
	return out, nil
}

// containsHyperlink checks if a comment contains URLs or hyperlinks
func containsHyperlink(body string) bool {
	lower := strings.ToLower(body)
	// Check for common URL patterns
	if strings.Contains(lower, "http://") || strings.Contains(lower, "https://") ||
		strings.Contains(lower, "www.") || strings.Contains(body, "](") {
		return true
	}
	return false
}

func collectComments(children []struct {
	Kind string `json:"kind"`
	Data struct {
		Body    string `json:"body"`
		Replies any    `json:"replies"`
	} `json:"data"`
}, out *[]string) {
	for _, child := range children {
		if child.Kind != "t1" {
			continue
		}
		if child.Data.Body != "" && !containsHyperlink(child.Data.Body) {
			*out = append(*out, child.Data.Body)
		}

		if repliesMap, ok := child.Data.Replies.(map[string]any); ok {
			repliesData, ok := repliesMap["data"].(map[string]any)
			if !ok {
				continue
			}
			repliesChildren, ok := repliesData["children"].([]any)
			if !ok {
				continue
			}

			normalized := make([]struct {
				Kind string `json:"kind"`
				Data struct {
					Body    string `json:"body"`
					Replies any    `json:"replies"`
				} `json:"data"`
			}, 0, len(repliesChildren))

			for _, rc := range repliesChildren {
				rcm, ok := rc.(map[string]any)
				if !ok {
					continue
				}
				kind, _ := rcm["kind"].(string)
				rd, _ := rcm["data"].(map[string]any)

				entry := struct {
					Kind string `json:"kind"`
					Data struct {
						Body    string `json:"body"`
						Replies any    `json:"replies"`
					} `json:"data"`
				}{Kind: kind}
				entry.Data.Body, _ = rd["body"].(string)
				entry.Data.Replies = rd["replies"]
				normalized = append(normalized, entry)
			}

			collectComments(normalized, out)
		}
	}
}
