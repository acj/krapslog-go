package main

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func Test_timeStemOffsets(t *testing.T) {
	type args struct {
		stemCount     int
		terminalWidth int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{"one stem", args{1, 5}, []int{0, 4}},
		{"two stems", args{2, 5}, []int{0, 4}},
		{"three stems", args{3, 5}, []int{0, 2, 4}},
		{"four stems", args{4, 5}, []int{0, 1, 2, 4}},
		{"five stems", args{5, 5}, []int{0, 1, 2, 3, 4}},
		{"six stems", args{6, 5}, []int{0, 1, 2, 2, 3, 4}},

		{"one stem", args{1, 10}, []int{0, 9}},
		{"two stems", args{2, 10}, []int{0, 9}},
		{"three stems", args{3, 10}, []int{0, 4, 9}},
		{"four stems", args{4, 10}, []int{0, 3, 6, 9}},
		{"five stems", args{5, 10}, []int{0, 2, 4, 6, 9}},
		{"six stems", args{6, 10}, []int{0, 2, 4, 5, 7, 9}},
		{"seven stems", args{7, 10}, []int{0, 2, 3, 4, 6, 7, 9}},
		{"eight stems", args{8, 10}, []int{0, 2, 3, 4, 5, 6, 7, 9}},
		{"nine stems", args{9, 10}, []int{0, 1, 2, 3, 4, 5, 6, 7, 9}},
		{"ten stems", args{10, 10}, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}},
		{"eleven stems", args{11, 10}, []int{0, 1, 2, 3, 4, 4, 5, 6, 7, 8, 9}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := timeStemOffsets(tt.args.stemCount, tt.args.terminalWidth); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("timeStemOffsets() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_timeStem(t *testing.T) {
	t.Run("header type, left align", func(t *testing.T) {
		ansicTime, _ := time.Parse(time.ANSIC, goAnsicDateFormat)
		ts := timeMarker{
			horizontalOffset: 0,
			time:             ansicTime,
		}
		canvas := newCanvas(canvasTypeHeader, 20, 5)

		ts.render(canvas, 5, stemAlignmentLeft)

		actual := canvas.String()
		expected := goAnsicTimeFormat + strings.Repeat(" ", 2) + "\n"
		expected += "|" + strings.Repeat(" ", 19) + "\n"
		expected += "|" + strings.Repeat(" ", 19) + "\n"
		expected += "|" + strings.Repeat(" ", 19) + "\n"
		expected += "|" + strings.Repeat(" ", 19) + "\n"

		if actual != expected {
			t.Errorf("incorrect canvas: expected '%s' (%d), got '%s' (%d)", expected, len(expected), actual, len(actual))
		}
	})

	t.Run("header type, right align", func(t *testing.T) {
		ansicTime, _ := time.Parse(time.ANSIC, goAnsicDateFormat)
		ts := timeMarker{
			horizontalOffset: 17,
			time:             ansicTime,
		}
		canvas := newCanvas(canvasTypeHeader, 20, 5)

		ts.render(canvas, 5, stemAlignmentRight)

		actual := canvas.String()
		expected := "Mon Jan 2 15:04:05  " + "\n"
		expected += "                 |  " + "\n"
		expected += "                 |  " + "\n"
		expected += "                 |  " + "\n"
		expected += "                 |  " + "\n"

		if actual != expected {
			t.Errorf("incorrect canvas: expected '%s' (%d), got '%s' (%d)", expected, len(expected), actual, len(actual))
		}
	})

	t.Run("footer type", func(t *testing.T) {
		ansicTime, _ := time.Parse(time.ANSIC, goAnsicDateFormat)
		ts := timeMarker{
			horizontalOffset: 0,
			time:             ansicTime,
		}
		canvas := newCanvas(canvasTypeFooter, 20, 5)

		ts.render(canvas, 5, stemAlignmentLeft)

		actual := canvas.String()
		expected := "|" + strings.Repeat(" ", 19) + "\n"
		expected += "|" + strings.Repeat(" ", 19) + "\n"
		expected += "|" + strings.Repeat(" ", 19) + "\n"
		expected += "|" + strings.Repeat(" ", 19) + "\n"
		expected += goAnsicTimeFormat + strings.Repeat(" ", 2) + "\n"

		if actual != expected {
			t.Errorf("incorrect canvas: expected '%s' (%d), got '%s' (%d)", expected, len(expected), actual, len(actual))
		}
	})

	t.Run("two time stems", func(t *testing.T) {
		ansicTime, _ := time.Parse(time.ANSIC, goAnsicDateFormat)
		ts1 := timeMarker{
			horizontalOffset: 0,
			time:             ansicTime,
		}
		ts2 := timeMarker{
			horizontalOffset: 5,
			time:             ansicTime,
		}
		width := 25
		height := 5
		canvas := newCanvas(canvasTypeHeader, width, height)

		ts1.render(canvas, height, stemAlignmentLeft)
		ts2.render(canvas, height-1, stemAlignmentLeft)

		actual := canvas.String()
		expected := goAnsicTimeFormat + strings.Repeat(" ", 7) + "\n"
		expected += "|    " + goAnsicTimeFormat + strings.Repeat(" ", width-18-5) + "\n"
		expected += "|    |" + strings.Repeat(" ", width-5-1) + "\n"
		expected += "|    |" + strings.Repeat(" ", width-5-1) + "\n"
		expected += "|    |" + strings.Repeat(" ", width-5-1) + "\n"

		if actual != expected {
			t.Errorf("incorrect canvas: expected '%s' (%d), got '%s' (%d)", expected, len(expected), actual, len(actual))
		}
	})
}
