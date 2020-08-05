package main

import "time"

func renderHeaderAndFooter(timestampsFromLines []int64, timeMarkerCount int, terminalWidth int) (string, string) {
	if timeMarkerCount == 0 {
		return "", ""
	}

	firstTimestamp := time.Unix(timestampsFromLines[0], 0).UTC()
	lastTimestamp := time.Unix(timestampsFromLines[len(timestampsFromLines)-1], 0).UTC()
	duration := lastTimestamp.Sub(firstTimestamp)
	footerMarkerCount := timeMarkerCount / 2
	if timeMarkerCount%2 != 0 {
		// If we have an odd number of markers, then the footer has one more marker than the header
		footerMarkerCount++
	}

	offsets := timeStemOffsets(timeMarkerCount, terminalWidth)
	durationBetweenOffsets := time.Duration(duration.Nanoseconds() / int64(terminalWidth))

	headerOffsets := offsets[footerMarkerCount:]
	headerCanvas := renderHeader(headerOffsets, terminalWidth, firstTimestamp, durationBetweenOffsets)

	footerOffsets := offsets[0:footerMarkerCount]
	footerCanvas := renderFooter(footerOffsets, terminalWidth, firstTimestamp, durationBetweenOffsets)
	return headerCanvas.String(), footerCanvas.String()
}

func renderHeader(markerOffsets []int, terminalWidth int, firstTimestamp time.Time, durationBetweenOffsets time.Duration) canvas {
	canvas := newCanvas(canvasTypeHeader, terminalWidth, len(markerOffsets)+1)
	needStackedMarkers := (len(firstTimestamp.Format(goAnsicTimeFormat))+1)*len(markerOffsets) >= (terminalWidth / 2)
	for verticalOffset, horizontalOffset := range markerOffsets {
		if needStackedMarkers {
			verticalOffset += 2
		} else {
			verticalOffset = 2
		}

		timeMarker{
			horizontalOffset: horizontalOffset,
			time:             firstTimestamp.Add(time.Duration(horizontalOffset) * durationBetweenOffsets),
		}.render(canvas, verticalOffset, stemAlignmentRight)
	}
	return canvas
}

func renderFooter(markerOffsets []int, terminalWidth int, firstTimestamp time.Time, durationBetweenOffsets time.Duration) canvas {
	canvas := newCanvas(canvasTypeFooter, terminalWidth, len(markerOffsets)+1)
	needStackedMarkers := (len(firstTimestamp.Format(goAnsicTimeFormat))+1)*len(markerOffsets) >= (terminalWidth / 2)
	for verticalOffset, horizontalOffset := range markerOffsets {
		if needStackedMarkers {
			verticalOffset = len(markerOffsets) - verticalOffset + 1
		} else {
			verticalOffset = 2
		}

		timeMarker{
			horizontalOffset: horizontalOffset,
			time:             firstTimestamp.Add(time.Duration(horizontalOffset) * durationBetweenOffsets),
		}.render(canvas, verticalOffset, stemAlignmentLeft)
	}
	return canvas
}
