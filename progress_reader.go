package main

import (
	"io"
	"math"
	"os"
)

type ProgressReader struct {
	io.Reader
	currentOffset int64
	totalBytes int64
	progressFunc func(currentProgress float64)
}

func NewProgressReader(r io.Reader, progressFunc func(float64)) (*ProgressReader, error) {
	pr := &ProgressReader{
		Reader:        r,
		currentOffset: 0,
		totalBytes:    0,
		progressFunc:  progressFunc,
	}

	if r, ok := r.(*os.File); ok {
		stat, err := r.Stat()
		if err != nil {
			return nil, err
		}
		pr.totalBytes = stat.Size()
	}

	return pr, nil
}

func (p *ProgressReader) Read(buf []byte) (int, error) {
	n, err := p.Reader.Read(buf)

	lastProgressPercent := math.Floor(100.0 * float64(p.currentOffset) / float64(p.totalBytes))
	nextProgressPercent := math.Floor(100.0 * (float64(p.currentOffset) + float64(n)) / float64(p.totalBytes))
	if p.totalBytes > 0 && nextProgressPercent != lastProgressPercent {
		p.progressFunc(nextProgressPercent)
	}

	p.currentOffset += int64(n)

	return n, err
}
