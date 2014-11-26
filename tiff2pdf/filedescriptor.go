package tiff2pdf

import (
	"errors"
	"sync"
)

const fdFirst = 10
var ErrOpenFailed = errors.New("Opening TIFF failed")

type fd struct {
	fd            int
	buffer        []byte
	offset        int64
	outputdisable int
}

var fdCount = fdFirst
var mtx sync.Mutex
var fdMap = make(map[int]*fd)

func NewFd(buffer []byte) *fd {
	var thisFd int

	fdo := &fd{
		buffer: buffer,
	}

	mtx.Lock()

	for {
		if fdCount > 5000 {
			fdCount = fdFirst
		}
		if _, ok := fdMap[fdCount]; !ok {
			thisFd = fdCount
			fdMap[thisFd] = fdo
			fdCount++
			break
		}
		fdCount++
	}

	mtx.Unlock()

	fdMap[thisFd].fd = thisFd

	return fdo
}
