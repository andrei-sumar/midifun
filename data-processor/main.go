package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
)

const (
	midiPortName = "IAC Driver Bus 1"
)

func main() {
	defer midi.CloseDriver()

	out, err := midi.FindOutPort(midiPortName)
	if err != nil {
		fmt.Printf("can't find output port: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("found outport: %s\n", out)

	token := os.Getenv("PULSOID_TOKEN")
	if token == "" {
		fmt.Fprintf(os.Stderr, "Error: PULSOID_TOKEN environment variable not set\n")
		os.Exit(1)
	}

	client := &http.Client{Timeout: 10 * time.Second}

	for {
		hrResp, err := fetchHeartRate(client, token)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			time.Sleep(1 * time.Second)
			continue
		}

		hr := hrResp.Data.HeartRate
		fmt.Printf("Measured at: %d, Heart rate: %d\n", hrResp.MeasuredAt, hr)

		cc := mapHRtoCC(hr)

		if err := sendMIDICC(out, 1, 1, cc); err != nil {
			fmt.Fprintf(os.Stderr, "Error sending MIDI CC: %v\n", err)
		} else {
			fmt.Println("MIDI CC sent successfully")
		}
		fmt.Println("MIDI CC sent successfully")

		time.Sleep(1 * time.Second)
	}
}
