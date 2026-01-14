package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
)

func main() {
	defer midi.CloseDriver()

	config, err := LoadConfig("config.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	out, err := midi.FindOutPort(config.MIDI.PortName)
	if err != nil {
		fmt.Printf("can't find output port: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("found outport: %s\n", out)

	windowWidth := config.Smoothing.WindowWidth
	rapidGrowthThreshold := config.Smoothing.RapidGrowthThreshold

	processor := NewHeartRateProcessor(windowWidth, rapidGrowthThreshold)
	fmt.Printf("Initialized heart rate processor with window width: %d, rapid growth threshold: %d\n", windowWidth, rapidGrowthThreshold)

	var fetcher HeartRateFetcher

	switch config.DataSource.Type {
	case "csv":
		fetcher, err = NewCSVFetcher(config.DataSource.CSV.Path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing CSV fetcher: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Using sample data from CSV file: %s\n", config.DataSource.CSV.Path)
	case "pulsoid":
		token := os.Getenv("PULSOID_TOKEN")
		if token == "" {
			fmt.Fprintf(os.Stderr, "Error: PULSOID_TOKEN environment variable not set\n")
			os.Exit(1)
		}
		client := &http.Client{Timeout: 10 * time.Second}
		fetcher = NewPulsoidFetcher(client, token)
		fmt.Println("Using Pulsoid API")
	default:
		fmt.Fprintf(os.Stderr, "Error: invalid data_source.type in config. Must be 'csv' or 'pulsoid'\n")
		os.Exit(1)
	}

	for {
		hrResp, err := fetcher.FetchHeartRate()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			time.Sleep(1 * time.Second)
			continue
		}

		hr := hrResp.Data.HeartRate
		smoothedHR := processor.AddValue(hr)
		fmt.Printf("Measured at: %d, Heart rate: %d (smoothed: %d)\n", hrResp.MeasuredAt, hr, smoothedHR)

		// Check for rapid growth and send appropriate MIDI CC
		if processor.HasRapidGrowth() {
			// Send rapid growth MIDI CC message
			if err := sendMIDICC(out, config.MIDI.Channel, config.MIDI.RapidGrowthCC, 127); err != nil {
				fmt.Fprintf(os.Stderr, "Error sending rapid growth MIDI CC: %v\n", err)
			} else {
				fmt.Println("Rapid growth detected - MIDI CC sent successfully")
			}
		}
		// Send normal tempo change MIDI CC
		cc := mapHRtoCC(smoothedHR)
		if err := sendMIDICC(out, config.MIDI.Channel, config.MIDI.TempoChangeCC, cc); err != nil {
			fmt.Fprintf(os.Stderr, "Error sending MIDI CC: %v\n", err)
		} else {
			fmt.Println("MIDI CC sent successfully")
		}

		time.Sleep(1 * time.Second)
	}
}
