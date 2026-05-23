package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type LLMProvider interface {
	Enhance(raw string) (string, error)
}

type GeminiProvider struct {
	APIKey string
	Model  string
}

type geminiRequest struct {
	Contents []geminiContent `json:"contents"`
}

type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiResponse struct {
	Candidates []geminiCandidate `json:"candidates"`
}

type geminiCandidate struct {
	Content geminiContent `json:"content"`
}

type geminiErrorResponse struct {
	Error struct {
		Code    int               `json:"code"`
		Message string            `json:"message"`
		Details []geminiErrorDetail `json:"details"`
	} `json:"error"`
}

type geminiErrorDetail struct {
	RetryDelay string `json:"retryDelay"`
}

const systemPrompt = `Rewrite the following raw text into a structured coding-agent prompt with these sections:

# Task

## Context

## Requirements

## Constraints

## Expected Output

## Acceptance Criteria

Preserve the original intent exactly. Do not add anything the user didn't ask for. Output only the final structured prompt, no preamble or explanation.`

const maxRetries = 5

func (g *GeminiProvider) Enhance(raw string) (string, error) {
	body := buildRequestBody(raw)
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", g.Model, g.APIKey)

	for attempt := 0; attempt <= maxRetries; attempt++ {
		result, retryAfter, err := g.tryEnhance(url, body)
		if err == nil {
			return result, nil
		}
		if retryAfter <= 0 || attempt == maxRetries {
			return "", err
		}
		log.Printf("Rate limited — retrying in %v (attempt %d/%d)", retryAfter, attempt+1, maxRetries)
		time.Sleep(retryAfter)
	}

	return "", fmt.Errorf("max retries exceeded")
}

func (g *GeminiProvider) tryEnhance(url string, body []byte) (string, time.Duration, error) {
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return "", 0, fmt.Errorf("API call failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		rawBody, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == 429 {
			delay := parseRetryDelay(rawBody)
			return "", delay, fmt.Errorf("rate limited: %s", string(rawBody))
		}
		return "", 0, fmt.Errorf("API returned %d: %s", resp.StatusCode, string(rawBody))
	}

	var geminiResp geminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		return "", 0, fmt.Errorf("decode response: %w", err)
	}

	if len(geminiResp.Candidates) == 0 {
		return "", 0, fmt.Errorf("no candidates in Gemini response")
	}

	var result string
	for _, p := range geminiResp.Candidates[0].Content.Parts {
		result += p.Text
	}
	return result, 0, nil
}

func parseRetryDelay(body []byte) time.Duration {
	var errResp geminiErrorResponse
	if json.Unmarshal(body, &errResp) != nil {
		return 0
	}
	for _, d := range errResp.Error.Details {
		if d.RetryDelay != "" {
			dur, err := time.ParseDuration(d.RetryDelay)
			if err == nil {
				return dur
			}
		}
	}
	return 30 * time.Second
}

func buildRequestBody(raw string) []byte {
	req := geminiRequest{
		Contents: []geminiContent{
			{
				Parts: []geminiPart{
					{Text: systemPrompt + "\n\n" + raw},
				},
			},
		},
	}
	body, _ := json.Marshal(req)
	return body
}
