package main

import (
	"flag"
	"fmt"
	"github.com/joliv/spark"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"os"
)

const (
	apacheCommonLogFormatDate = "02/Jan/2006:15:04:05.000"
	goAnsicDateFormat         = "Mon Jan 2 15:04:05 2006"
)

func main() {
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

	timeFinder, err := NewTimeFinder(*requestedDateFormat)
	if err != nil {
		exitWithMessage("invalid timestamp format: %v", err)
	}

	var logReader io.Reader = file
	if *displayProgress {
		logReader, err = NewProgressReader(file, func(progressPercent float64) {
			fmt.Fprintf(os.Stderr, "\r%.f%%", progressPercent)
			if progressPercent == 100.0 {
				fmt.Fprintf(os.Stderr, "\r")
			}
		})
		if err != nil {
			exitWithMessage("failed to read log: %v", err)
		}
	}

	timestampsFromLines, err := timeFinder.extractTimestampFromEachLine(logReader)
	if err != nil {
		exitWithMessage("failed to process log: %v", err)
	}
	if len(timestampsFromLines) == 0 {
		exitWithMessage("didn't find any lines with recognizable dates")
	}

	terminalWidth, _, err := terminal.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		exitWithMessage("couldn't get terminal size: %v", err)
	}
	
	linesPerCharacter := BinTimestampsToFitLineWidth(timestampsFromLines, terminalWidth)
	sparkLine := spark.Line(linesPerCharacter)

	if *timeMarkerCount > 0 {
		firstTimestamp := timestampsFromLines[0]
		lastTimestamp := timestampsFromLines[len(timestampsFromLines)-1]
		duration := lastTimestamp.Sub(firstTimestamp)
		headerText := headerText(firstTimestamp.Add(duration/2), lastTimestamp, *timeMarkerCount/2, terminalWidth)
		footerText := footerText(firstTimestamp, firstTimestamp.Add(duration/2), *timeMarkerCount/2, terminalWidth)

		fmt.Print(headerText)
		fmt.Println(sparkLine)
		fmt.Print(footerText)
	} else {
		fmt.Println(sparkLine)
	}

	os.Exit(0)
}

func exitWithMessage(m string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, m, args...)
	os.Exit(-1)
}
