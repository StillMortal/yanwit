package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type AlternativesRequest struct {
	Text  string `json:"text"`
	Count int    `json:"count"`
	Style string `json:"style"`
}

type AlternativesResponse struct {
	Alternatives []string `json:"alternatives"`
	FromCache    bool     `json:"from_cache"`
	LatencyMs    int      `json:"latency_ms"`
}

type ManipulationRequest struct {
	Text string `json:"text"`
}

type ManipulationResponse struct {
	HasManipulation bool     `json:"has_manipulation"`
	Types           []string `json:"types"`
	Confidence      float64  `json:"confidence"`
	Suggestions     []string `json:"suggestions"`
	RewrittenVersion string  `json:"rewritten_version,omitempty"`
}

var httpClient = &http.Client{
	Timeout: 5 * time.Second,
}

func GetAlternatives(text string, count int, style string) (*AlternativesResponse, error) {
	url := os.Getenv("AI_ALTERNATIVES_URL")
	if url == "" {
		url = "http://localhost:8002"
	}
	
	reqBody := AlternativesRequest{
		Text:  text,
		Count: count,
		Style: style,
	}
	
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	
	resp, err := httpClient.Post(url+"/generate", "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AI service returned status %d", resp.StatusCode)
	}
	
	var result AlternativesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	
	return &result, nil
}

func DetectManipulation(text string) (*ManipulationResponse, error) {
	url := os.Getenv("AI_MANIPULATION_URL")
	if url == "" {
		url = "http://localhost:8003"
	}
	
	reqBody := ManipulationRequest{Text: text}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	
	resp, err := httpClient.Post(url+"/detect", "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AI service returned status %d", resp.StatusCode)
	}
	
	var result ManipulationResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	
	return &result, nil
}