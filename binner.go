package main

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"math"
	"os"
	"regexp"
	"strings"
	"time"
)

type Binner struct {
	dateFormat string
	exitFunc func(m string, args ...interface{})
	displayProgress bool
	strictLineParsing bool
	totalLogSize int64
}

func NewBinner() Binner {
	return Binner{
		exitFunc: exitWithMessage,
	}
}

func (b Binner) BinLinesByTimestamp(r io.Reader) ([]float64, error) {
	formatRegexString := convertDateFormatToRegex(b.dateFormat)
	formatRegex := regexp.MustCompile(formatRegexString)
	canonicalTime, err := time.Parse(time.ANSIC, goAnsicDateFormat)
	if err != nil {
		exitWithMessage("couldn't parse canonical time: %v", err)
	}
	t, err := time.Parse(b.dateFormat, b.dateFormat)
	if err != nil || t != canonicalTime {
		fmt.Fprintf(os.Stderr, "Invalid date/time format '%s'", b.dateFormat)

		if err != nil {
			fmt.Fprintf(os.Stderr, ": %v", err)
		}

		exitWithMessage("\n\nThe format must include year, day, and time. Please follow the format described in https://golang.org/pkg/time/#Time.Format\n")
	}

	dateFinder := func(line string) (time.Time, error) {
		if dateString := formatRegex.FindString(line); dateString != "" {
			return time.Parse(b.dateFormat, dateString)
		}

		return time.Time{}, fmt.Errorf("couldn't find time in line '%s'", line)
	}

	times := make([]int64, 0)

	offset := int64(0)
	progressByteThreshold := b.totalLogSize / 10

	nextProgressByteThreshold := int64(0)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		t, err := dateFinder(line)
		if err != nil {
			if b.strictLineParsing {
				exitWithMessage("\rerror: %v", err)
			}
			continue
		}
		times = append(times, t.Unix())

		if offset >= nextProgressByteThreshold {
			progressPercent := (float64(nextProgressByteThreshold)/float64(b.totalLogSize))*100.0
			if b.displayProgress && b.totalLogSize > 0 {
				fmt.Fprintf(os.Stderr, "\r%.f%%", progressPercent)
			}

			nextProgressByteThreshold += progressByteThreshold
		}
		offset += int64(len(line))
	}
	fmt.Fprintf(os.Stderr, "\r")

	terminalWidth, _, err := terminal.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		exitWithMessage("couldn't get terminal size: %v", err)
	}

	if len(times) == 0 {
		exitWithMessage("didn't find any lines with recognizable dates")
	}
	firstTime := times[0]
	lastTime := times[len(times)-1]
	spread := lastTime - firstTime
	secondsPerBucket := int64(math.Ceil(float64(spread) / float64(terminalWidth)))
	linesPerBucket := make([]float64, terminalWidth, terminalWidth)

	for _, lineUnixTime := range times {
		bucket := (lineUnixTime - firstTime) / secondsPerBucket
		linesPerBucket[bucket]++
	}

	return linesPerBucket, nil
}

func convertDateFormatToRegex(format string) string {
	replaceSet := []string{
		".", "\\.",
		"2006", "\\d{4}",
		"06", "\\d{2}",
		"Jan", "[A-Za-z]{3}",
		"January", "[A-Za-z]{3,4,5,6,7,8,9}",
		"0", "\\d",
		"1", "\\d",
		"2", "\\d",
		"3", "\\d",
		"4", "\\d",
		"5", "\\d",
		"6", "\\d",
		"7", "\\d",
		"8", "\\d",
		"9", "\\d",
	}

	replacer := strings.NewReplacer(replaceSet...)
	regex := replacer.Replace(format)

	return regex
}