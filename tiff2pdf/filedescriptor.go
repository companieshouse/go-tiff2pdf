package tiff2pdf

import (
	"errors"
	"sync"
)

var (
	ErrOpenFailed = errors.New("Opening TIFF failed")
)

type fd struct {
	fd            int
	buffer        []byte
	offset        int64
	outputdisable int
}

var fdCount = 10
var mtx sync.Mutex
var fdMap = make(map[int]*fd)

func NewFd(buffer []byte) *fd {
	mtx.Lock()
	thisFd := fdCount
	fdCount++
	mtx.Unlock()

	fdo := &fd{
		fd:     thisFd,
		buffer: buffer,
	}
	fdMap[thisFd] = fdo

	return fdo
}
