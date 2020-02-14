package main

import (
	"reflect"
	"testing"
	"time"
)

func Test_BinTimestampsToFitLineWidth(t *testing.T) {
	type args struct {
		timestampsFromLines []time.Time
		terminalWidth       int
	}
	tests := []struct {
		name string
		args args
		want []float64
	}{
		{
			name: "zero timestamps", args: args{
				timestampsFromLines: []time.Time{},
				terminalWidth:       5,
			},
			want: repeatFloat(0, 5),
		},
		{
			name: "one timestamp", args: args{
				timestampsFromLines: []time.Time{{}},
				terminalWidth:       5,
			},
			want: []float64{1, 0, 0, 0, 0},
		},
		{
			name: "10 timestamps", args: args{
				timestampsFromLines: []time.Time{
					time.Time{}.Add(1 * time.Second),
					time.Time{}.Add(2 * time.Second),
					time.Time{}.Add(3 * time.Second),
					time.Time{}.Add(4 * time.Second),
					time.Time{}.Add(5 * time.Second),
					time.Time{}.Add(6 * time.Second),
					time.Time{}.Add(7 * time.Second),
					time.Time{}.Add(8 * time.Second),
					time.Time{}.Add(9 * time.Second),
					time.Time{}.Add(10 * time.Second),
				},
				terminalWidth: 80,
			},
			want: []float64{1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := binTimestampsToFitLineWidth(tt.args.timestampsFromLines, tt.args.terminalWidth); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("binTimestampsToFitLineWidth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func repeatFloat(num float64, count int) []float64 {
	a := make([]float64, count, count)
	for idx := range a {
		a[idx] = num
	}
	return a
}

func repeatInt64(num int64, count int) []int64 {
	a := make([]int64, count, count)
	for idx := range a {
		a[idx] = num
	}
	return a
}

func repeatTime(t time.Time, count int) []time.Time {
	a := make([]time.Time, count, count)
	for idx := range a {
		a[idx] = t
	}
	return a
}
