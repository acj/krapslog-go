package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/joliv/spark"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	apacheCommonLogFormatDate = "2/Jan/2006:15:04:05.000"
	goAnsicDateFormat = "Mon Jan 2 15:04:05 2006"
)

func main() {
	var requestedDateFormat = flag.String("format", apacheCommonLogFormatDate, "date format to look for (see https://golang.org/pkg/time/#Time.Format)")
	var showProgress = flag.Bool("progress", false, "display progress while scanning the log file")
	var strict = flag.Bool("strict", false, "abort scanning when a line doesn't contain a timestamp")
	flag.Parse()

	filename := flag.Arg(0)
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening '%s': %v", filename, err)
		os.Exit(-1)
	}
	defer file.Close()

	formatRegexString := convertDateFormatToRegex(*requestedDateFormat)
	formatRegex := regexp.MustCompile(formatRegexString)
	canonicalTime, err := time.Parse(time.ANSIC, goAnsicDateFormat)
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't parse canonical time: %v", err)
		os.Exit(-1)
	}
	t, err := time.Parse(*requestedDateFormat, *requestedDateFormat)
	if err != nil || t != canonicalTime {
		fmt.Fprintf(os.Stderr, "Invalid date/time format '%s'", *requestedDateFormat)

		if err != nil {
			fmt.Fprintf(os.Stderr, ": %v", err)
		}

		fmt.Fprint(os.Stderr, "\n\nThe format must include year, day, and time. Please follow the format described in https://golang.org/pkg/time/#Time.Format\n")
		os.Exit(-1)
	}

	dateFinder := func(line string) (time.Time, error) {
		if dateString := formatRegex.FindString(line); dateString != "" {
			return time.Parse(*requestedDateFormat, dateString)
		}

		return time.Time{}, fmt.Errorf("couldn't find time in line '%s'", line)
	}

	times := make([]int64, 0)

	fileStat, err := file.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error calling stat: %v", err)
		os.Exit(-1)
	}
	offset := int64(0)
	totalBytes := fileStat.Size()
	progressByteThreshold := totalBytes / 10

	nextProgressByteThreshold := int64(0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		t, err := dateFinder(line)
		if err != nil {
			if *strict {
				fmt.Fprintf(os.Stderr,"\rerror: %v", err)
				os.Exit(-1)
			}
			continue
		}
		times = append(times, t.Unix())

		if offset >= nextProgressByteThreshold {
			progressPercent := (float64(nextProgressByteThreshold)/float64(totalBytes))*100.0
			if *showProgress {
				fmt.Fprintf(os.Stderr, "\r%.f%%", progressPercent)
			}

			nextProgressByteThreshold += progressByteThreshold
		}
		offset += int64(len(line))
	}
	fmt.Fprintf(os.Stderr, "\r")

	terminalWidth, _, err := terminal.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't get terminal size: %v", err)
		os.Exit(-1)
	}

	if len(times) == 0 {
		fmt.Fprintln(os.Stderr, "didn't find any lines with recognizable dates")
		os.Exit(-1)
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

	sparkline := spark.Line(linesPerBucket)

	fmt.Println(sparkline)

	os.Exit(0)
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
