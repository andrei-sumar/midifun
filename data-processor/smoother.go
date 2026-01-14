package main

// HeartRateSmoother smooths heart rate values using a sliding window mean
type HeartRateSmoother struct {
	window    []int
	windowSize int
	maxSize   int
}

// NewHeartRateSmoother creates a new smoother with the specified window width
func NewHeartRateSmoother(windowWidth int) *HeartRateSmoother {
	if windowWidth < 1 {
		windowWidth = 1
	}
	return &HeartRateSmoother{
		window:    make([]int, 0, windowWidth),
		windowSize: 0,
		maxSize:   windowWidth,
	}
}

// AddValue adds a new heart rate value to the window and returns the smoothed value
// The smoothed value is the mean of all values in the current window
func (s *HeartRateSmoother) AddValue(value int) int {
	if s.windowSize < s.maxSize {
		s.window = append(s.window, value)
		s.windowSize++
	} else {
		copy(s.window, s.window[1:])
		s.window[s.maxSize-1] = value
	}

	return s.GetSmoothedValue()
}

// GetSmoothedValue returns the mean of all values in the current window
func (s *HeartRateSmoother) GetSmoothedValue() int {
	if s.windowSize == 0 {
		return 0
	}

	sum := 0
	for i := 0; i < s.windowSize; i++ {
		sum += s.window[i]
	}

	mean := float64(sum) / float64(s.windowSize)
	return int(mean + 0.5)
}
