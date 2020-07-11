package timefinder

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"
)

const (
	apacheCommonLogFormatDate = "02/Jan/2006:15:04:05.000"
	goAnsicDateFormat         = "Mon Jan 2 15:04:05 2006"
)

type TimeFinder struct {
	timeFormat string
	timeRegex  *regexp.Regexp
}

// NewTimeFinder constructs a new TimeFinder instance. It returns an error if the time format is invalid.
func NewTimeFinder(timeFormat string) (*TimeFinder, error) {
	formatRegexString := convertTimeFormatToRegex(timeFormat)
	formatRegex, err := regexp.Compile(formatRegexString)
	if err != nil {
		return nil, err
	}
	if err := checkDateFormatForErrors(timeFormat); err != nil {
		return nil, err
	}
	return &TimeFinder{
		timeFormat: timeFormat,
		timeRegex:  formatRegex,
	}, nil
}

// ExtractTimestampFromEachLine scans each line of the reader to find a timestamp.  It returns a slice of all the
// timestamps that were found. If no timestamp is found, then the line is skipped.
func (tf *TimeFinder) ExtractTimestampFromEachLine(r io.Reader) []int64 {
	times := make([]int64, 0)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		t, err := tf.findFirstTimestamp(scanner.Text())
		if err != nil {
			continue
		}
		times = append(times, t.UTC().Unix())
	}

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
