package main

import (
	"flag"
	"fmt"
	"github.com/acj/krapslog/timefinder"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"os"
)

const (
	apacheCommonLogFormatDate = "02/Jan/2006:15:04:05.000"
	goAnsicDateFormat         = "Mon Jan 2 15:04:05 2006"
)

func main() {
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

	if err := displaySparkline(file, os.Stdout, *requestedDateFormat, *timeMarkerCount, *displayProgress); err != nil {
		exitWithErrorMessage("couldn't generate sparkline: %v", err)
	}

	os.Exit(0)
}

func displaySparkline(r io.Reader, w io.Writer, dateFormat string, timeMarkerCount int, shouldDisplayProgress bool) error {
	timeFinder, err := timefinder.NewTimeFinder(dateFormat)
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

	timestampsFromLines := timeFinder.ExtractTimestampFromEachLine(r)
	if len(timestampsFromLines) == 0 {
		return fmt.Errorf("didn't find any lines with recognizable dates")
	}

	terminalWidth, _, err := terminal.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't get terminal size (%v); defaulting to 80 characters\n", err)
		terminalWidth = 80
	}

	logLineCountPerCharacter := binTimestamps(timestampsFromLines, terminalWidth)
	sparkLine := Line(logLineCountPerCharacter)

	header, footer := renderHeaderAndFooter(timestampsFromLines, timeMarkerCount, terminalWidth)

	fmt.Fprint(w, header)
	fmt.Fprintln(w, sparkLine)
	fmt.Fprint(w, footer)

	return nil
}

func exitWithErrorMessage(m string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, m, args...)
	os.Exit(-1)
}
