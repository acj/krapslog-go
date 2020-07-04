package timefinder

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
	"sync"
	"time"
)

const (
	apacheCommonLogFormatDate = "02/Jan/2006:15:04:05.000"
	goAnsicDateFormat = "Mon Jan 2 15:04:05 2006"
)

type TimeFinder struct {
	parallelism int
	timeFormat  string
	timeRegex   *regexp.Regexp
}

// NewTimeFinder constructs a new TimeFinder instance. It returns an error if the time format is invalid.
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

// ExtractTimestampFromEachLine scans each line of the reader to find a timestamp.  It returns a slice of all the
// timestamps that were found. If no timestamp is found, then the line is skipped.
func (tf *TimeFinder) ExtractTimestampFromEachLine(r io.Reader) []time.Time {
	var wg sync.WaitGroup
	lineChans := make([]chan string, tf.parallelism)
	timeChans := make([]chan time.Time, tf.parallelism)

	for i, _ := range lineChans {
		lineChans[i] = make(chan string, 1)
		timeChans[i] = make(chan time.Time, 1)
	}

	// Stage 1: Produce lines
	wg.Add(1)
	go func() {
		defer wg.Done()
		nextInChanIndex := 0
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			lineChans[nextInChanIndex] <- scanner.Text()
			nextInChanIndex = (nextInChanIndex + 1) % tf.parallelism
		}

		for i := 0; i < tf.parallelism; i++ {
			close(lineChans[i])
		}
	}()

	// Stage 2: Fan out to convert lines into time structs. This is the expensive step.
	for i := 0; i < tf.parallelism; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer close(timeChans[i])

			for line := range lineChans[i] {
				t, err := tf.findFirstTimestamp(line)
				if err != nil {
					// Zero value will be ignored
					timeChans[i] <- time.Time{}
					continue
				}
				timeChans[i] <- t
			}
		}()
	}

	// Stage 3: Fan in to assemble list of timestamps, reading one line per channel to preserve line ordering
	times := make([]time.Time, 0)
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			drainedChannelCount := 0
			for i := 0; i < tf.parallelism; i++ {
				if t, ok := <-timeChans[i]; ok {
					times = append(times, t)
				} else {
					drainedChannelCount++
				}
			}

			// All inbound channels are drained, so we're done
			if drainedChannelCount == tf.parallelism {
				return
			}
		}
	}()

	wg.Wait()

	return times
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

	return strings.NewReplacer(replaceSet...).Replace(format)
}

func (tf *TimeFinder) findFirstTimestamp(s string) (time.Time, error) {
	if dateString := tf.timeRegex.FindString(s); dateString != "" {
		return time.Parse(tf.timeFormat, dateString)
	}

	return time.Time{}, fmt.Errorf("couldn't find time in line '%s'", s)
}
