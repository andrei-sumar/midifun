package main

// clampInt clamps value between min and max
func clampInt(min, value, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// mapHRtoCC maps heart rate to MIDI CC value
func mapHRtoCC(hr int, minTempo uint8) uint8 {
	cc := hr - int(minTempo)
	cc = clampInt(0, cc, 127)
	return uint8(cc)
}
