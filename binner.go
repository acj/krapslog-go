package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"regexp"
	"strings"
	"time"
)

type Binner struct {
	dateFinderFunc func(line, dateFormat string, dateFormatRegex *regexp.Regexp) (time.Time, error)
	dateFormat string
	exitFunc func(m string, args ...interface{})
	displayProgress bool
	strictLineParsing bool
	totalLogSize int64
}

func NewBinner() Binner {
	return Binner{
		dateFinderFunc: findFirstTimestamp,
		exitFunc:       exitWithMessage,
	}
}

func (b Binner) Bin(r io.Reader, maxWidth int) ([]float64, error) {
	formatRegexString := convertDateFormatToRegex(b.dateFormat)
	formatRegex := regexp.MustCompile(formatRegexString)
	canonicalTime, err := time.Parse(time.ANSIC, goAnsicDateFormat)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse canonical time: %v", err)
	}
	if err := checkDateFormatForErrors(b.dateFormat, canonicalTime); err != nil {
		return nil, err
	}

	timestampsFromLines, err := b.extractTimestampsFromLines(r, formatRegex)
	if err != nil {
		return nil, err
	}
	if len(timestampsFromLines) == 0 {
		return nil, errors.New("didn't find any lines with recognizable dates")
	}

	linesPerBucket := binTimestampsToFitLineWidth(timestampsFromLines, maxWidth)

	return linesPerBucket, nil
}

func (b Binner) extractTimestampsFromLines(r io.Reader, formatRegex *regexp.Regexp) ([]int64, error) {
	times := make([]int64, 0)
	offset := int64(0)
	progressByteThreshold := b.totalLogSize / 10
	nextProgressByteThreshold := int64(0)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		t, err := b.dateFinderFunc(line, b.dateFormat, formatRegex)
		if err != nil {
			if b.strictLineParsing {
				return nil, fmt.Errorf("\rerror: %v", err)
			}
			continue
		}
		times = append(times, t.Unix())

		if offset >= nextProgressByteThreshold {
			progressPercent := (float64(nextProgressByteThreshold) / float64(b.totalLogSize)) * 100.0
			if b.displayProgress && b.totalLogSize > 0 {
				fmt.Fprintf(os.Stderr, "\r%.f%%", progressPercent)
			}

			nextProgressByteThreshold += progressByteThreshold
		}
		offset += int64(len(line))
	}
	fmt.Fprintf(os.Stderr, "\r")
	return times, nil
}

func binTimestampsToFitLineWidth(timestampsFromLines []int64, terminalWidth int) []float64 {
	firstTime := timestampsFromLines[0]
	lastTime := timestampsFromLines[len(timestampsFromLines)-1]
	spread := lastTime - firstTime
	secondsPerBucket := int64(math.Ceil(float64(spread) / float64(terminalWidth)))
	linesPerBucket := make([]float64, terminalWidth, terminalWidth)
	for _, lineUnixTime := range timestampsFromLines {
		bucket := (lineUnixTime - firstTime) / secondsPerBucket
		linesPerBucket[bucket]++
	}
	return linesPerBucket
}

func checkDateFormatForErrors(dateFormat string, canonicalTime time.Time) error {
	t, err := time.Parse(dateFormat, dateFormat)
	if err != nil || t != canonicalTime {
		fmt.Fprintf(os.Stderr, "Invalid date/time format '%s'", dateFormat)

		if err != nil {
			fmt.Fprintf(os.Stderr, ": %v", err)
		}

		return errors.New("\n\nThe format must include year, day, and time. Please follow the format described in https://golang.org/pkg/time/#Time.Format\n")
	}

	return nil
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

func findFirstTimestamp(s, timestampFormat string, timestampFormatRegex *regexp.Regexp) (time.Time, error) {
	if dateString := timestampFormatRegex.FindString(s); dateString != "" {
		return time.Parse(timestampFormat, dateString)
	}

	return time.Time{}, fmt.Errorf("couldn't find time in line '%s'", s)
}