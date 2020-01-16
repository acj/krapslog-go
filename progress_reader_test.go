package main

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestNewProgressReader(t *testing.T) {
	type args struct {
		r            io.Reader
		progressFunc func(float64)
	}
	tests := []struct {
		name    string
		args    args
		want    *ProgressReader
		wantErr bool
	}{
		{
			name: "",
			args: args{
				r:            strings.NewReader("one\ntwo\nthree\n"),
				progressFunc: nil,
			},
			want: &ProgressReader{
				strings.NewReader("one\ntwo\nthree\n"),
				0,
				0,
				nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewProgressReader(tt.args.r, tt.args.progressFunc)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewProgressReader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewProgressReader() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProgressReader_Read(t *testing.T) {
	t.Run("when we don't know the overall size, it doesn't invoke the callback function", func(t *testing.T) {
		called := false
		pr := &ProgressReader{
			strings.NewReader("hi mom"),
			0,
			0,
			func(float64) {
				called = true
			},
		}

		var buf []byte
		_, err := pr.Read(buf)
		if err != nil {
			t.Fatalf("failed to read: %v", err)
		}
		if called {
			t.Errorf("callback function was invoked unexpectedly")
		}
	})

	t.Run("when there's no change in percentage read, it doesn't invoke the callback function", func(t *testing.T) {
		called := false
		pr := &ProgressReader{
			strings.NewReader("hi mom"),
			999990,
			1000000,
			func(float64) {
				called = true
			},
		}

		buf := make([]byte, 6)
		_, err := pr.Read(buf)
		if err != nil {
			t.Fatalf("failed to read: %v", err)
		}
		if called {
			t.Errorf("callback function was invoked unexpectedly")
		}
	})

	t.Run("when there is a change in percentage read, it invokes the callback function with the correct percentage", func(t *testing.T) {
		called := false
		actualPercentage := -1.0
		pr := &ProgressReader{
			strings.NewReader("hi mom"),
			4,
			10,
			func(p float64) {
				called = true
				actualPercentage = p
			},
		}

		buf := make([]byte, 6)
		n, err := pr.Read(buf)
		if err != nil {
			t.Fatalf("failed to read: %v", err)
		}
		if !called {
			t.Errorf("callback function should have been called but wasn't")
		}
		if n != 6 {
			t.Errorf("read error: got %d bytes but want %d", n, 6)
		}
		expectedPercentage := 100.0
		if actualPercentage != expectedPercentage {
			t.Errorf("wrong percentage: got %f but want %f", actualPercentage, expectedPercentage)
		}
	})
}