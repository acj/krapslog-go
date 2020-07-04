package main

import (
	"github.com/acj/krapslog/internal/test"
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
			want: test.RepeatFloat(0, 5),
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
			if got := binTimestamps(tt.args.timestampsFromLines, tt.args.terminalWidth); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("binTimestamps() = %v, want %v", got, tt.want)
			}
		})
	}
}

