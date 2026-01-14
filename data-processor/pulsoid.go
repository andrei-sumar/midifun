package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	apiURL = "https://dev.pulsoid.net/api/v1/data/heart_rate/latest?scope=data:heart_rate:read"
)

type HeartRateResponse struct {
	MeasuredAt int64 `json:"measured_at"`
	Data       struct {
		HeartRate int `json:"heart_rate"`
	} `json:"data"`
}

// HeartRateFetcher defines the interface for fetching heart rate data
type HeartRateFetcher interface {
	FetchHeartRate() (*HeartRateResponse, error)
}

// PulsoidFetcher implements HeartRateFetcher using the Pulsoid API
type PulsoidFetcher struct {
	client *http.Client
	token  string
}

// NewPulsoidFetcher creates a new PulsoidFetcher instance
func NewPulsoidFetcher(client *http.Client, token string) *PulsoidFetcher {
	return &PulsoidFetcher{
		client: client,
		token:  token,
	}
}

// FetchHeartRate fetches heart rate data from the Pulsoid API
func (p *PulsoidFetcher) FetchHeartRate() (*HeartRateResponse, error) {
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.token))

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	var hrResp HeartRateResponse
	if err := json.Unmarshal(body, &hrResp); err != nil {
		return nil, fmt.Errorf("parsing JSON: %w", err)
	}
	return &hrResp, nil
}
