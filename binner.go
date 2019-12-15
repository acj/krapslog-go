package main

import (
	"time"
)

func BinTimestampsToFitLineWidth(timesFromLines []time.Time, bucketCount int) []float64 {
	timestampsFromLines := make([]int64, len(timesFromLines), len(timesFromLines))
	linesPerBucket := make([]float64, bucketCount, bucketCount)

	switch len(timesFromLines) {
	case 0:
		return linesPerBucket
	case 1:
		linesPerBucket[0] = 1
		return linesPerBucket
	}

	for index, t := range timesFromLines {
		timestampsFromLines[index] = t.Unix()
	}

	firstTime := timestampsFromLines[0]
	lastTime := timestampsFromLines[len(timestampsFromLines)-1]
	spread := lastTime - firstTime + 1
	for _, lineUnixTime := range timestampsFromLines {
		if lineUnixTime < firstTime {
			continue
		}
		bucket := int64((float64(bucketCount) * float64(lineUnixTime - firstTime)) / float64(spread))
		linesPerBucket[bucket]++
	}
	return linesPerBucket
}
