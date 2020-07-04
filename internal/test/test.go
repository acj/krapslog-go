package test

import "time"

func RepeatFloat(num float64, count int) []float64 {
	a := make([]float64, count, count)
	for idx := range a {
		a[idx] = num
	}
	return a
}

func RepeatInt64(num int64, count int) []int64 {
	a := make([]int64, count, count)
	for idx := range a {
		a[idx] = num
	}
	return a
}

func RepeatTime(t time.Time, count int) []time.Time {
	a := make([]time.Time, count, count)
	for idx := range a {
		a[idx] = t
	}
	return a
}
