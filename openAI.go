package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

var (
	openAiKey  = os.Getenv("OPENAI_API_KEY")
	oaClient   *openai.Client
	clientOnce sync.Once
)

func getOpenAIClient() *openai.Client {
	clientOnce.Do(func() {
		oaClient = openai.NewClient(openAiKey)
	})
	return oaClient
}

// Allow aliases; do not silently downgrade gpt-5
var modelAlias = map[string]string{
	"gpt-5":   "gpt-5",
	"gpt5":    "gpt-5",
	"gpt-4o":  "gpt-4o",
	"gpt4":    "gpt-4o",
	"gpt-4":   "gpt-4o",
	"default": "gpt-4o",
}

func resolveModel(name string) string {
	name = strings.TrimSpace(strings.ToLower(name))
	if name == "" {
		return modelAlias["default"]
	}
	if m, ok := modelAlias[name]; ok {
		return m
	}
	return name // allow raw model ids
}

func runChat(model, userPrompt, systemPrompt string, maxTokens int, temperature float64) string {
	if openAiKey == "" {
		return "OpenAI key missing"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	resolved := resolveModel(model)

	req := openai.ChatCompletionRequest{
		Model: resolved,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
			{Role: openai.ChatMessageRoleUser, Content: userPrompt},
		},
	}

	// Token field: gpt-5 uses MaxCompletionTokens, others use MaxTokens
	if strings.HasPrefix(resolved, "gpt-5") {
		req.MaxCompletionTokens = maxTokens
		// Do NOT set Temperature / TopP / N / PresencePenalty / FrequencyPenalty (fixed by model)
	} else {
		req.MaxTokens = maxTokens
		req.Temperature = float32(temperature)
		// (Optionally set other tuning params here for non–gpt-5 models)
	}

	resp, err := getOpenAIClient().CreateChatCompletion(ctx, req)
	if err != nil {
		return fmt.Sprintf("OpenAI error: %v", err)
	}
	if len(resp.Choices) == 0 {
		return "No completion returned"
	}
	out := resp.Choices[0].Message.Content
	if out == "" {
		return "Empty response"
	}
	if len(out) > 2000 {
		out = out[:2000]
	}
	return out
}

// helper removed: MaxCompletionTokens is an int in this SDK, so a pointer helper is unnecessary.

// Backwards compatible primary function
func gpt(query string, contextMsg string) string {
	return runChat("gpt-4o", query, contextMsg, 5000, 0.7)
}

// Custom model variant
func gptModel(model string, query string, contextMsg string) string {
	return runChat(model, query, contextMsg, 4000, 0.7)
}

// generateImageContent tries URL first (no response_format to avoid 400), then falls back to base64.
// Returns (url, bytes, error). If url != "", bytes will be nil.
func generateImageContent(prompt string) (string, []byte, error) {
	if openAiKey == "" {
		return "", nil, fmt.Errorf("OpenAI key missing")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	client := getOpenAIClient()

	// FIRST ATTEMPT: Do NOT set ResponseFormat (gpt-image-1 rejects response_format=url -> 400)
	req := openai.ImageRequest{
		Model:  "dall-e-3",
		Prompt: prompt,
		N:      1,
		Size:   openai.CreateImageSize1024x1024,
		// Leave ResponseFormat empty (defaults to URL if supported)
	}

	resp, err := client.CreateImage(ctx, req)
	if err != nil {
		return "", nil, fmt.Errorf("image generation error: %w", err)
	}
	if len(resp.Data) > 0 && resp.Data[0].URL != "" {
		return resp.Data[0].URL, nil, nil
	}

	// FALLBACK: Ask specifically for base64
	req.ResponseFormat = openai.CreateImageResponseFormatB64JSON
	resp, err = client.CreateImage(ctx, req)
	if err != nil {
		return "", nil, fmt.Errorf("image (b64) generation error: %w", err)
	}
	if len(resp.Data) == 0 || resp.Data[0].B64JSON == "" {
		return "", nil, fmt.Errorf("no image content returned")
	}

	raw, err := base64.StdEncoding.DecodeString(resp.Data[0].B64JSON)
	if err != nil {
		return "", nil, fmt.Errorf("base64 decode failed: %w", err)
	}
	return "", raw, nil
}

// Backward compatible simple helper
func generateImage(prompt string) string {
	url, _, err := generateImageContent(prompt)
	if err != nil {
		return fmt.Sprintf("Image generation error: %v", err)
	}
	if url == "" {
		return "Image generated (no direct URL); internal upload used."
	}
	return url
}

func dallEText(prompt string) string {
	return generateImage(prompt)
}
