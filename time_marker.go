package main

import (
	"log"
	"math"
	"strings"
	"time"
)

type stemAlignment int

const (
	stemAlignmentLeft = iota
	stemAlignmentRight
)

type canvasType int

const (
	canvasTypeHeader = iota
	canvasTypeFooter
)

type canvas struct {
	_type canvasType
	buf   [][]byte
}

func newCanvas(canvasType canvasType, width, height int) canvas {
	buf := make([][]byte, height, height)
	for i := 0; i < len(buf); i++ {
		buf[i] = make([]byte, width, width)
		for j := 0; j < len(buf[i]); j++ {
			buf[i][j] = ' '
		}
	}

	return canvas{
		buf:   buf,
		_type: canvasType,
	}
}

func (c canvas) put(row int, col int, text []byte) {
	copy(c.buf[row][col:col+len(text)], text)
}

func (c canvas) String() string {
	var displayHeader strings.Builder
	switch c._type {
	case canvasTypeHeader:
		// Display timestamps at the top of the stems
		for i := len(c.buf) - 1; i >= 0; i-- {
			displayHeader.Write(c.buf[i])
			displayHeader.WriteByte('\n')
		}
	case canvasTypeFooter:
		// Display timestamps at the bottom of the stems
		for i := 0; i < len(c.buf); i++ {
			displayHeader.Write(c.buf[i])
			displayHeader.WriteByte('\n')
		}
	default:
		log.Fatalf("unrecognized canvas type: %d", c._type)
	}
	return displayHeader.String()
}

type timeMarker struct {
	horizontalOffset int
	time             time.Time
}

func (ts timeMarker) render(canvas canvas, verticalOffset int, alignment stemAlignment) {
	for i := 0; i < verticalOffset; i++ {
		if i == verticalOffset-1 {
			displayTime := ts.time.Format("Mon Jan 2 15:04:05")
			startingOffset := ts.horizontalOffset
			if alignment == stemAlignmentRight {
				startingOffset -= len(displayTime) - 1
			}
			canvas.put(i, startingOffset, []byte(displayTime))
		} else {
			canvas.put(i, ts.horizontalOffset, []byte{'|'})
		}
	}
}

func timeStemOffsets(markerCount, terminalWidth int) []int {
	var offsets []int

	// Always show a marker at the left edge
	offsets = append(offsets, 0)

	// Divide the non-edge offsets into equally sized segments, placing a marker between them
	skip := float64(terminalWidth-2) / float64(markerCount-1)
	currentOffset := skip
	for i := 0; i < markerCount-2; i++ {
		offsets = append(offsets, int(math.Ceil(currentOffset))%terminalWidth)
		currentOffset += skip
	}

	// Always show a marker at the right edge
	offsets = append(offsets, terminalWidth-1)

	return offsets
}
