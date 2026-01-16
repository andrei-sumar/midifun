package main

// HeartRateProcessor processes heart rate values using a sliding window mean
// and detects rapid growth patterns
type HeartRateProcessor struct {
	window               []int
	windowSize           int
	maxSize              int
	rapidGrowthThreshold int
}

// NewHeartRateProcessor creates a new processor with the specified window width and rapid growth threshold
func NewHeartRateProcessor(windowWidth int, rapidGrowthThreshold int) *HeartRateProcessor {
	if windowWidth < 1 {
		windowWidth = 1
	}
	return &HeartRateProcessor{
		window:               make([]int, 0, windowWidth),
		windowSize:           0,
		maxSize:              windowWidth,
		rapidGrowthThreshold: rapidGrowthThreshold,
	}
}

// AddValue adds a new heart rate value to the window and returns the smoothed value
// The smoothed value is the mean of all values in the current window
func (p *HeartRateProcessor) AddValue(value int) int {
	if p.windowSize < p.maxSize {
		p.window = append(p.window, value)
		p.windowSize++
	} else {
		copy(p.window, p.window[1:])
		p.window[p.maxSize-1] = value
	}

	return p.GetSmoothedValue()
}

// GetSmoothedValue returns the mean of all values in the current window
func (p *HeartRateProcessor) GetSmoothedValue() int {
	if p.windowSize == 0 {
		return 0
	}

	sum := 0
	for i := 0; i < p.windowSize; i++ {
		sum += p.window[i]
	}

	mean := float64(sum) / float64(p.windowSize)
	return int(mean + 0.5)
}

// HasRapidGrowth checks if there was rapid growth in the window
// Only returns true if the window is fully filled and the growth exceeds the threshold
func (p *HeartRateProcessor) HasRapidGrowth() bool {
	if p.windowSize < p.maxSize {
		return false
	}

	// Compare beginning and end values
	beginningValue := p.window[0]
	endValue := p.window[p.maxSize-1]
	growth := endValue - beginningValue

	return growth > p.rapidGrowthThreshold
}

// HasDecrease checks if there was steady decrease in the window
// Only returns true if the window is fully filled
func (p *HeartRateProcessor) HasDecrease() bool {
	if p.windowSize < p.maxSize {
		return false
	}

	// Compare beginning and end values
	beginningValue := p.window[0]
	endValue := p.window[p.maxSize-1]
	growth := beginningValue - endValue

	return growth > p.rapidGrowthThreshold
}
