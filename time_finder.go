package main

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
	"sync"
	"time"
)

type TimeFinder struct {
	parallelism int
	timeFormat string
	timeRegex  *regexp.Regexp
}

func NewTimeFinder(timeFormat string, parallelism int) (*TimeFinder, error) {
	formatRegexString := convertTimeFormatToRegex(timeFormat)
	formatRegex := regexp.MustCompile(formatRegexString)
	if err := checkDateFormatForErrors(timeFormat); err != nil {
		return nil, err
	}
	return &TimeFinder{
		parallelism: parallelism,
		timeFormat:  timeFormat,
		timeRegex:   formatRegex,
	}, nil
}

func checkDateFormatForErrors(dateFormat string) error {
	canonicalTime, err := time.Parse(time.ANSIC, goAnsicDateFormat)
	if err != nil {
		return fmt.Errorf("couldn't parse canonical time: %v", err)
	}
	t, err := time.Parse(dateFormat, dateFormat)
	if err != nil || t != canonicalTime {
		errorText := fmt.Sprintf("invalid date/time format '%s'", dateFormat)

		if err != nil {
			errorText += fmt.Sprintf(": %v", err)
		}

		return fmt.Errorf("%s\n\nThe format must include year, day, and time. Please follow the format described in https://golang.org/pkg/time/#Time.Format\n", errorText)
	}

	return nil
}

func convertTimeFormatToRegex(format string) string {
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

func (tf *TimeFinder) extractTimestampFromEachLine(r io.Reader) ([]time.Time, error) {
	times := make([]time.Time, 0)

	inC := make(chan string)
	outC := make(chan time.Time)
	var inWg sync.WaitGroup
	for i := 0; i < tf.parallelism; i++ {
		inWg.Add(1)
		go func() {
			defer inWg.Done()

			for line := range inC {
				t, err := tf.findFirstTimestamp(line)
				if err != nil {
					// TODO: Optionally allow exit on error
					//return nil, err
				}
				outC <- t
			}
		}()
	}

	var outWg sync.WaitGroup
	outWg.Add(1)
	go func() {
		defer outWg.Done()

		for t := range outC {
			times = append(times, t)
		}
	}()

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		inC <- scanner.Text()
	}
	close(inC)

	inWg.Wait()
	close(outC)
	outWg.Wait()

	return times, nil
}

func (tf *TimeFinder) findFirstTimestamp(s string) (time.Time, error) {
	if dateString := tf.timeRegex.FindString(s); dateString != "" {
		return time.Parse(tf.timeFormat, dateString)
	}

	return time.Time{}, fmt.Errorf("couldn't find time in line '%s'", s)
}
