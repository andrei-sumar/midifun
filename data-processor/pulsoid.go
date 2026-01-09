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

func fetchHeartRate(client *http.Client, token string) (*HeartRateResponse, error) {
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
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
