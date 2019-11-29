package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"time"

	"github.com/joliv/spark"
)

func main() {
	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening '%s': %v", filename, err)
		os.Exit(-1)
	}
	defer file.Close()

	haproxyDateTimeRegex := regexp.MustCompile(`[A-Za-z]{3} \d\d? \d{2}:\d{2}:\d{2}`)
	const haproxyDateTimeLayout = "Jan 2 15:04:05"

	dateFinder := func(line string) (time.Time, error) {
		if dateString := haproxyDateTimeRegex.FindString(line); dateString != "" {
			return time.Parse(haproxyDateTimeLayout, dateString)
		}

		return time.Time{}, fmt.Errorf("couldn't find time in line '%s'", line)
	}

	times := make([]int64, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		t, err := dateFinder(scanner.Text())
		if err != nil {
			// TODO: bailing out on these errors should be optional
			panic(err)
		}
		times = append(times, t.Unix())
	}

	// TODO: Make this dynamic -- either based on terminal width or a parameter value
	//  width, height, err := terminal.GetSize(0)
	terminalWidth := 80

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