package main

import (
	"io"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"
)

var apacheCommonLogFormatFields = struct {
	timeFormat string
	timeRegex  *regexp.Regexp
}{
	timeFormat: apacheCommonLogFormatDate,
	timeRegex:  regexp.MustCompile(convertTimeFormatToRegex(apacheCommonLogFormatDate)),
}

var sampleApacheCommonLogFormatTimestamp = "23/Nov/2019:06:26:40.781"

func Test_checkDateFormatForErrors(t *testing.T) {
	type args struct {
		dateFormat    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "for garbage value, returns an error",
			args:    args{
				dateFormat: "blah",
			},
			wantErr: true,
		},
		{
			name:    "for valid date value, does not return an error",
			args:    args{
				dateFormat: "2/Jan/2006:15:04:05.000",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkDateFormatForErrors(tt.args.dateFormat); (err != nil) != tt.wantErr {
				t.Errorf("checkDateFormatForErrors() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

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
			got := convertTimeFormatToRegex(tt.args.format)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertTimeFormatToRegex() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func parseTime(s string) time.Time {
	t, err := time.Parse(apacheCommonLogFormatDate, s)
	if err != nil {
		panic("couldn't parse time: " + err.Error())
	}
	return t
}

func TestTimeFinder_extractTimestampFromEachLine(t *testing.T) {
	type fields struct {
		timeFormat string
		timeRegex  *regexp.Regexp
	}
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []time.Time
		wantErr bool
	}{
		{
			name: "for single line, returns single timestamp",
			fields: apacheCommonLogFormatFields,
			args: args{
				r:          strings.NewReader(`haproxy[20128]: 192.168.23.456:57305 [23/Nov/2019:06:26:40.781] public myapp/i-05fa49c0e7db8c328 0/0/0/78/78 206 913/458 - - ---- 9/9/6/0/0 0/0 {} {||1|bytes 0-0/499704} "GET /foobarbaz.html HTTP/1.1\n`),
			},
			want: []time.Time{parseTime(sampleApacheCommonLogFormatTimestamp)},
			wantErr: false,
		},
		{
			name: "for two lines, returns two timestamps",
			fields: apacheCommonLogFormatFields,
			args: args{
				r:          strings.NewReader(strings.Repeat("haproxy[20128]: 192.168.23.456:57305 [23/Nov/2019:06:26:40.781] public myapp/i-05fa49c0e7db8c328 0/0/0/78/78 206 913/458 - - ---- 9/9/6/0/0 0/0 {} {||1|bytes 0-0/499704} \"GET /foobarbaz.html\" HTTP/1.1\n", 2)),
			},
			want: repeatTime(parseTime(sampleApacheCommonLogFormatTimestamp), 2),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tf := &TimeFinder{
				timeFormat: tt.fields.timeFormat,
				timeRegex:  tt.fields.timeRegex,
			}
			got, err := tf.extractTimestampFromEachLine(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractTimestampFromEachLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractTimestampFromEachLine() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeFinder_findFirstTimestamp(t *testing.T) {
	type fields struct {
		timeFormat string
		timeRegex  *regexp.Regexp
	}
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    time.Time
		wantErr bool
	}{
		{
			name: "for garbage input, returns an error",
			fields: apacheCommonLogFormatFields,
			args: args{
				"garbage date",
			},
			want: time.Time{},
			wantErr: true,
		},
		{
			name: "for valid input, returns a date and no error",
			fields: apacheCommonLogFormatFields,
			args: args{
				sampleApacheCommonLogFormatTimestamp,
			},
			want: func() time.Time { t, _ := time.Parse(apacheCommonLogFormatDate, sampleApacheCommonLogFormatTimestamp); return t }(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tf := &TimeFinder{
				timeFormat: tt.fields.timeFormat,
				timeRegex:  tt.fields.timeRegex,
			}
			got, err := tf.findFirstTimestamp(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("findFirstTimestamp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findFirstTimestamp() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewTimeFinder(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		tf, err := NewTimeFinder(apacheCommonLogFormatDate)
		if err != nil {
			t.Errorf("unexpected NewTimeFinder() error = %v", err)
			return
		}

		expected, _ := time.Parse(apacheCommonLogFormatDate, sampleApacheCommonLogFormatTimestamp)
		actual, err := tf.findFirstTimestamp(sampleApacheCommonLogFormatTimestamp)
		if err != nil {
			t.Errorf("unexpected findFirstTimestamp error = %v", err)
			return
		}
		if actual != expected {
			t.Errorf("findFirstTimestamp: got %v, want %v", actual, expected)
			return
		}
	})

	t.Run("for garbage format, returns an error", func(t *testing.T) {
		_, err := NewTimeFinder("garbage")
		if err == nil {
			t.Error("NewTimeFinder: expected an error but didn't get one")
			return
		}
	})
}