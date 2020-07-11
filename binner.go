package main

func binTimestamps(timesFromLines []int64, bucketCount int) []float64 {
	linesPerBucket := make([]float64, bucketCount, bucketCount)

	switch len(timesFromLines) {
	case 0:
		return linesPerBucket
	case 1:
		linesPerBucket[0] = 1
		return linesPerBucket
	}

	firstTime := timesFromLines[0]
	lastTime := timesFromLines[len(timesFromLines)-1]
	spread := lastTime - firstTime + 1
	for _, lineUnixTime := range timesFromLines {
		if lineUnixTime < firstTime {
			continue
		}
		bucket := int64((float64(bucketCount) * float64(lineUnixTime-firstTime)) / float64(spread))
		linesPerBucket[bucket]++
	}
	return linesPerBucket
}
