package main

import (
	"flag"
	"fmt"
	"github.com/joliv/spark"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"os"
	"runtime"
	"time"
)

const (
	apacheCommonLogFormatDate = "02/Jan/2006:15:04:05.000"
	goAnsicDateFormat         = "Mon Jan 2 15:04:05 2006"
)

func main() {
	var concurrency = flag.Int("concurrency", runtime.GOMAXPROCS(0), "number of log lines to process concurrently")
	var displayProgress = flag.Bool("progress", false, "display progress while scanning the log file")
	var requestedDateFormat = flag.String("format", apacheCommonLogFormatDate, "date format to look for (see https://golang.org/pkg/time/#Time.Format)")
	var timeMarkerCount = flag.Int("markers", 0, "number of time markers to display")
	flag.Parse()

	if flag.NArg() == 0 {
		exitWithErrorMessage("no filename given")
	}

	filename := flag.Arg(0)
	file, err := os.Open(filename)
	if err != nil {
		exitWithErrorMessage("error opening '%s': %v", filename, err)
	}
	defer file.Close()

	if err := displaySparkline(file, os.Stdout, *requestedDateFormat, *timeMarkerCount, *displayProgress, *concurrency); err != nil {
		exitWithErrorMessage("couldn't generate sparkline: %v", err)
	}

	os.Exit(0)
}

func displaySparkline(r io.Reader, w io.Writer, dateFormat string, timeMarkerCount int, shouldDisplayProgress bool, concurrency int) error {
	timeFinder, err := NewTimeFinder(dateFormat, concurrency)
	if err != nil {
		return fmt.Errorf("invalid timestamp format: %v", err)
	}

	if shouldDisplayProgress {
		r, err = NewProgressReader(r, func(progressPercent float64) {
			fmt.Fprintf(os.Stderr, "\r%.f%%", progressPercent)
			if progressPercent == 100.0 {
				fmt.Fprintf(os.Stderr, "\r")
			}
		})
		if err != nil {
			return fmt.Errorf("failed to read log: %v", err)
		}
	}

	timestampsFromLines, err := timeFinder.extractTimestampFromEachLine(r)
	if err != nil {
		return fmt.Errorf("failed to process log: %v", err)
	}
	if len(timestampsFromLines) == 0 {
		return fmt.Errorf("didn't find any lines with recognizable dates")
	}

	terminalWidth, _, err := terminal.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't get terminal size (%v); defaulting to 80 characters\n", err)
		terminalWidth = 80
	}

	logLineCountPerCharacter := binTimestamps(timestampsFromLines, terminalWidth)
	sparkLine := spark.Line(logLineCountPerCharacter)

	header, footer := renderHeaderAndFooter(timestampsFromLines, timeMarkerCount, terminalWidth)

	fmt.Fprint(w, header)
	fmt.Fprintln(w, sparkLine)
	fmt.Fprint(w, footer)

	return nil
}

func renderHeaderAndFooter(timestampsFromLines []time.Time, timeMarkerCount int, terminalWidth int) (string, string) {
	if timeMarkerCount == 0 {
		return "", ""
	}

	firstTimestamp := timestampsFromLines[0]
	lastTimestamp := timestampsFromLines[len(timestampsFromLines)-1]
	duration := lastTimestamp.Sub(firstTimestamp)
	footerMarkerCount := timeMarkerCount / 2
	if timeMarkerCount%2 != 0 {
		// If we have an odd number of markers, then the header has one more marker than the footer
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

func exitWithErrorMessage(m string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, m, args...)
	os.Exit(-1)
}
