package main

import (
	"bytes"
	"testing"
	"time"
)

func Test_renderLine(t *testing.T) {
	t.Run("renders expected date on correct line", func(t *testing.T) {
		width := 80
		buf := make([][]byte, 1, 1)
		buf[0] = make([]byte, width, width)
		for j := 0; j < width; j++ {
			buf[0][j] = ' '
		}

		renderLine(buf[0], time.Time{}, 79, 0, 0, true)

		expected := []byte{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', 'M', 'o', 'n', ' ', 'J', 'a', 'n', ' ', '1', ' ', '0', '0', ':', '0', '0', ':', '0', '0'}
		if bytes.Compare(expected, buf[0]) != 0 {
			t.Errorf("buffers do not match: expected \n\t%s\nbut got\n\t%s", expected, buf[0])
		}
	})

	t.Run("renders single time stem on correct line", func(t *testing.T) {
		width := 80
		buf := make([][]byte, 1, 1)
		buf[0] = make([]byte, width, width)
		for j := 0; j < width; j++ {
			buf[0][j] = ' '
		}

		renderLine(buf[0], time.Time{}, width-1, 1, 0, true)

		expected := []byte{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', '|'}
		if bytes.Compare(expected, buf[0]) != 0 {
			t.Errorf("buffers do not match: expected \n\t%s\nbut got\n\t%s", expected, buf[0])
		}
	})

	t.Run("renders multiple time stems on correct lines", func(t *testing.T) {
		width := 80
		buf := make([][]byte, 1, 1)
		buf[0] = make([]byte, width, width)
		for j := 0; j < width; j++ {
			buf[0][j] = ' '
		}

		renderLine(buf[0], time.Time{}, width-1, 2, 0, true)
		renderLine(buf[0], time.Time{}, width-3, 2, 0, true)

		expected := []byte{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', '|', ' ', '|'}
		if bytes.Compare(expected, buf[0]) != 0 {
			t.Errorf("buffers do not match: expected \n\t%s\nbut got\n\t%s", expected, buf[0])
		}
	})
}
