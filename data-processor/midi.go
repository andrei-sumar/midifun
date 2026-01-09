package main

import (
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
)

func sendMIDICC(out drivers.Out, channel, controller, value uint8) error {
	msg := midi.ControlChange(channel, controller, value)
	sendFunc, err := midi.SendTo(out)
	if err != nil {
		return err
	}
	return sendFunc(msg)
}
