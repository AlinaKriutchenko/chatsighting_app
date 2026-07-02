package analytics

import (
	"math"
	"sort"
)

type Spike struct {
	Window
	Rank     int
	TopWords []string
}

func FindTopSpikes(windows []Window, n int) []Spike {
	sorted := make([]Window, len(windows))
	copy(sorted, windows)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].MessageCount > sorted[j].MessageCount
	})
	if n > len(sorted) {
		n = len(sorted)
	}
	spikes := make([]Spike, n)
	for i := 0; i < n; i++ {
		spikes[i] = Spike{
			Window:   sorted[i],
			Rank:     i + 1,
			TopWords: TopWords(sorted[i].Messages, 5),
		}
	}
	return spikes
}

type SlimWindow struct {
	S float64 `json:"s"` // start seconds
	C int     `json:"c"` // message count
}

func SlimWindows(windows []Window) []SlimWindow {
	out := make([]SlimWindow, len(windows))
	for i, w := range windows {
		out[i] = SlimWindow{S: w.StartSeconds, C: w.MessageCount}
	}
	return out
}

func StreamAverage(windows []Window) int {
	if len(windows) == 0 {
		return 0
	}
	total := 0
	for _, w := range windows {
		total += w.MessageCount
	}
	return int(math.Round(float64(total) / float64(len(windows))))
}
