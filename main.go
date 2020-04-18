package main

import (
	"flag"
	"fmt"
	"github.com/joliv/spark"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"os"
	"runtime"
)

const (
	apacheCommonLogFormatDate = "02/Jan/2006:15:04:05.000"
	goAnsicDateFormat         = "Mon Jan 2 15:04:05 2006"
)

func main() {
	var concurrency = flag.Int("concurrency", runtime.GOMAXPROCS(0), "number of log lines to process concurrently")
	var requestedDateFormat = flag.String("format", apacheCommonLogFormatDate, "date format to look for (see https://golang.org/pkg/time/#Time.Format)")
	var displayProgress = flag.Bool("progress", false, "display progress while scanning the log file")
	var timeMarkerCount = flag.Int("markers", 0, "number of time markers to display")
	flag.Parse()

	filename := flag.Arg(0)
	file, err := os.Open(filename)
	if err != nil {
		exitWithMessage("error opening '%s': %v", filename, err)
	}
	defer file.Close()

	if err := displaySparklineForLog(file, os.Stdout, *requestedDateFormat, *timeMarkerCount, *displayProgress, *concurrency); err != nil {
		exitWithMessage(err.Error())
	}

	os.Exit(0)
}

func displaySparklineForLog(r io.Reader, w io.Writer, dateFormat string, timeMarkerCount int, shouldDisplayProgress bool, concurrency int) error {
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
		fmt.Fprintf(os.Stdout, "couldn't get terminal size (%v); defaulting to 80 characters\n", err)
		terminalWidth = 80
	}

	linesPerCharacter := binTimestampsToFitLineWidth(timestampsFromLines, terminalWidth)
	sparkLine := spark.Line(linesPerCharacter)

	if timeMarkerCount > 0 {
		firstTimestamp := timestampsFromLines[0]
		lastTimestamp := timestampsFromLines[len(timestampsFromLines)-1]
		duration := lastTimestamp.Sub(firstTimestamp)
		headerText := headerText(firstTimestamp.Add(duration/2), lastTimestamp, timeMarkerCount/2, terminalWidth)
		footerText := footerText(firstTimestamp, firstTimestamp.Add(duration/2), timeMarkerCount/2, terminalWidth)

		fmt.Fprint(w, headerText)
		fmt.Fprintln(w, sparkLine)
		fmt.Fprint(w, footerText)
	} else {
		fmt.Fprintln(w, sparkLine)
	}

	return nil
}

func exitWithMessage(m string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, m, args...)
	os.Exit(-1)
}
