package main

import (
	"strings"
	"time"
)

func headerText(firstTimestamp time.Time, lastTimestamp time.Time, markerCount int, terminalWidth int) string {
	// +1 to allow space for the bottom marker to have a one-line stem
	header := make([][]byte, markerCount+1, markerCount+1)
	for i := 0; i < len(header); i++ {
		header[i] = make([]byte, terminalWidth, terminalWidth)
		for j := 0; j < len(header[i]); j++ {
			header[i][j] = ' '
		}
	}
	segmentDuration := time.Duration(lastTimestamp.Sub(firstTimestamp).Seconds() / float64(markerCount))

	for i := 0; i < markerCount; i++ {
		for j := 0; j < len(header); j++ {
			renderLine(
				header[j],
				firstTimestamp.Add(time.Duration(j*1e9)*segmentDuration),
				(terminalWidth/2-1)+(i+1)*(terminalWidth/2/markerCount),
				i+1,
				j,
				true,
			)
		}
	}

	var displayHeader strings.Builder
	for i := len(header) - 1; i >= 0; i-- {
		displayHeader.Write(header[i])
		displayHeader.WriteByte('\n')
	}

	return displayHeader.String()
}

func footerText(firstTimestamp time.Time, lastTimestamp time.Time, markerCount int, terminalWidth int) string {
	// +1 to allow space for the top marker to have a one-line stem
	header := make([][]byte, markerCount+1, markerCount+1)
	for i := 0; i < len(header); i++ {
		header[i] = make([]byte, terminalWidth, terminalWidth)
		for j := 0; j < len(header[i]); j++ {
			header[i][j] = ' '
		}
	}
	segmentDuration := time.Duration(lastTimestamp.Sub(firstTimestamp).Seconds() / float64(markerCount))

	for i := 0; i < markerCount; i++ {
		for j := 0; j < len(header); j++ {
			renderLine(
				header[j],
				firstTimestamp.Add(time.Duration(i*1e9)*segmentDuration),
				i*(terminalWidth/2/markerCount),
				len(header)-(i+1),
				j,
				false,
			)
		}
	}

	var displayHeader strings.Builder
	for i := 0; i < len(header); i++ {
		displayHeader.Write(header[i])
		displayHeader.WriteByte('\n')
	}

	return displayHeader.String()
}

func renderLine(buf []byte, timestamp time.Time, horizontalOffset int, verticalOffset int, linesOffsetFromSparkline int, alignStemAtEndOfTime bool) {
	if verticalOffset == linesOffsetFromSparkline {
		displayTime := timestamp.Format("Mon Jan 2 15:04:05")
		startingOffset := horizontalOffset
		if alignStemAtEndOfTime {
			startingOffset -= len(displayTime) - 1
		}
		copy(buf[startingOffset:startingOffset+len(displayTime)], displayTime)
	} else if verticalOffset > linesOffsetFromSparkline {
		buf[horizontalOffset] = '|'
	}
}
