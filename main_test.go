package main

import (
	"reflect"
	"testing"
)

func Test_convertDateFormatToRegex(t *testing.T) {
	type args struct {
		format string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "common log format",
			args:    args{"2/Jan/2006:15:04:05.000"},
			want:    "\\d/[A-Za-z]{3}/\\d{4}:\\d\\d:\\d\\d:\\d\\d\\.\\d\\d\\d",
		},
		{
			name:    "Go default log format",
			args:    args{"2006/1/2 15:04:05"},
			want:    "\\d{4}/\\d/\\d \\d\\d:\\d\\d:\\d\\d",
		},
		{
			name:    "abbreviated log format",
			args:    args{"Jan 2 15:04:05"},
			want:    "[A-Za-z]{3} \\d \\d\\d:\\d\\d:\\d\\d",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertDateFormatToRegex(tt.args.format)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertDateFormatToRegex() got = %v, want %v", got, tt.want)
			}
		})
	}
}