package main

import (
	"flag"
	"fmt"
	"github.com/joliv/spark"
	"os"
)

const (
	apacheCommonLogFormatDate = "2/Jan/2006:15:04:05.000"
	goAnsicDateFormat = "Mon Jan 2 15:04:05 2006"
)

func main() {
	var requestedDateFormat = flag.String("format", apacheCommonLogFormatDate, "date format to look for (see https://golang.org/pkg/time/#Time.Format)")
	var displayProgress = flag.Bool("progress", false, "display progress while scanning the log file")
	var strictParsing = flag.Bool("strict", false, "abort scanning when a line doesn't contain a timestamp")
	flag.Parse()

	filename := flag.Arg(0)
	file, err := os.Open(filename)
	if err != nil {
		exitWithMessage("error opening '%s': %v", filename, err)
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		exitWithMessage("error calling stat: %v", err)
	}

	b := NewBinner()
	b.dateFormat = *requestedDateFormat
	b.displayProgress = *displayProgress
	b.strictLineParsing = *strictParsing
	b.totalLogSize = fileStat.Size()

	linesPerBucket, err := b.BinLinesByTimestamp(file)
	if err != nil {
		exitWithMessage("failed to process log: %v", err)
	}
	sparkline := spark.Line(linesPerBucket)
	fmt.Println(sparkline)

	os.Exit(0)
}

func exitWithMessage(m string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, m, args...)
	os.Exit(-1)
}
