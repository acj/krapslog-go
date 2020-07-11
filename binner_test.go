package main

import (
	"github.com/acj/krapslog/internal/test"
	"reflect"
	"testing"
)

func Test_BinTimestampsToFitLineWidth(t *testing.T) {
	type args struct {
		timestampsFromLines []int64
		terminalWidth       int
	}
	tests := []struct {
		name string
		args args
		want []float64
	}{
		{
			name: "zero timestamps", args: args{
				timestampsFromLines: []int64{},
				terminalWidth:       5,
			},
			want: test.RepeatFloat(0, 5),
		},
		{
			name: "one timestamp", args: args{
				timestampsFromLines: []int64{0},
				terminalWidth:       5,
			},
			want: []float64{1, 0, 0, 0, 0},
		},
		{
			name: "10 timestamps", args: args{
				timestampsFromLines: []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				terminalWidth:       80,
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
