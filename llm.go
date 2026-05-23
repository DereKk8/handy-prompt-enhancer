package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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

const systemPrompt = `Rewrite the following raw text into a structured coding-agent prompt with these sections:

# Task

## Context

## Requirements

## Constraints

## Expected Output

## Acceptance Criteria

Preserve the original intent exactly. Do not add anything the user didn't ask for. Output only the final structured prompt, no preamble or explanation.`

func (g *GeminiProvider) Enhance(raw string) (string, error) {
	reqBody := geminiRequest{
		Contents: []geminiContent{
			{
				Parts: []geminiPart{
					{Text: systemPrompt + "\n\n" + raw},
				},
			},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", g.Model, g.APIKey)
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("API call failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errBody bytes.Buffer
		errBody.ReadFrom(resp.Body)
		return "", fmt.Errorf("API returned %d: %s", resp.StatusCode, errBody.String())
	}

	var geminiResp geminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	if len(geminiResp.Candidates) == 0 {
		return "", fmt.Errorf("no candidates in Gemini response")
	}

	var result string
	for _, p := range geminiResp.Candidates[0].Content.Parts {
		result += p.Text
	}
	return result, nil
}
