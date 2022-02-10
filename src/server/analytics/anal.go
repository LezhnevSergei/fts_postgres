package analytics

import (
	"fmt"
	"sort"
)

type Anal struct{}

func CreateAnal() Anal {
	return Anal{}
}

func (a Anal) Show(searchTimes []float32) {
	fmt.Printf("Median (100): %v ms\n", a.calcMedian(searchTimes))
	fmt.Printf("Avg (100): %v ms\n", a.calcAvg(searchTimes))
	fmt.Printf("Min (100): %v ms\n", a.calcMin(searchTimes))
	fmt.Printf("Max (100): %v ms\n", a.calcMax(searchTimes))
}

func (a Anal) calcAvg(n []float32) float32 {
	var sum float32 = 0

	for _, t := range n {
		sum += t
	}

	return sum / float32(len(n))
}

func (a Anal) calcMin(n []float32) float32 {
	var min float32 = 1000

	for _, t := range n {
		if t < min {
			min = t
		}
	}

	return min
}

func (a Anal) calcMax(n []float32) float32 {
	var max float32 = 0

	for _, t := range n {
		if t > max {
			max = t
		}
	}

	return max
}

func (a Anal) calcMedian(n []float32) float32 {
	sort.Slice(n, func(i, j int) bool { return n[i] < n[j] })

	mNumber := len(n) / 2

	if len(n)%2 != 0 {
		return n[mNumber]
	}

	return float32(n[mNumber-1]+n[mNumber]) / 2.0
}
