package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

// CSVFetcher implements HeartRateFetcher using a CSV file
type CSVFetcher struct {
	rows       [][]string
	currentIdx int
}

// NewCSVFetcher creates a new CSVFetcher instance from a CSV file path.
// initialIdx specifies the starting index (defaults to 0 if not provided).
func NewCSVFetcher(csvPath string, initialIdx ...int) (*CSVFetcher, error) {
	file, err := os.Open(csvPath)
	if err != nil {
		return nil, fmt.Errorf("opening CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("reading CSV: %w", err)
	}

	// Skip header row
	if len(rows) > 0 {
		rows = rows[1:]
	}

	idx := 0
	if len(initialIdx) > 0 {
		idx = initialIdx[0]
	}

	return &CSVFetcher{
		rows:       rows,
		currentIdx: idx,
	}, nil
}

// FetchHeartRate fetches the next heart rate data from the CSV file
func (c *CSVFetcher) FetchHeartRate() (*HeartRateResponse, error) {
	if c.currentIdx >= len(c.rows) {
		c.currentIdx = 0
	}

	row := c.rows[c.currentIdx]
	c.currentIdx++

	if len(row) < 3 {
		return nil, fmt.Errorf("invalid CSV row: expected 3 columns, got %d", len(row))
	}

	measuredAtMs, err := strconv.ParseInt(row[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parsing measured_at_ms: %w", err)
	}

	heartRate, err := strconv.Atoi(row[2])
	if err != nil {
		return nil, fmt.Errorf("parsing heart_rate: %w", err)
	}

	hrResp := &HeartRateResponse{
		MeasuredAt: measuredAtMs,
	}
	hrResp.Data.HeartRate = heartRate

	return hrResp, nil
}
