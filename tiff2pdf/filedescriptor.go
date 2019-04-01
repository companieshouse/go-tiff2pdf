package tiff2pdf

import (
	"errors"
	"log"
	"sync"
)

const fdFirst = 10

var ErrOpenFailed = errors.New("Opening TIFF failed")

type fd struct {
	fd            int
	buffer        []byte
	offset        int64
	outputdisable int
	warnings      []string
	errors        []string
}

var fdCount = fdFirst
var mtx sync.Mutex

// var fdMap = make(map[int]*fd)
// var fdMap = MapWrapper{wrappedMap: make(map[int]*fd)}
var fdMap MapWrapper

func NewFd(buffer []byte) *fd {
	var thisFd int

	fdo := &fd{
		buffer:   buffer,
		warnings: make([]string, 0),
		errors:   make([]string, 0),
	}

	mtx.Lock()

	for {
		if fdCount > 5000 {
			fdCount = fdFirst
		}
		// if _, ok := fdMap[fdCount]; !ok {
		_, ok := fdMap.Load(fdCount)
		if !ok {
			thisFd = fdCount
			// fdMap[thisFd] = fdo
			fdMap.Store(thisFd, fdo)
			fdCount++
			break
		}
		fdCount++
	}

	mtx.Unlock()

	// fdMap[thisFd].fd = thisFd
	// loaded, ok := fdMap[thisFd]
	loaded, ok := fdMap.Load(thisFd)
	if !ok {
		log.Printf("[%d] NewFd load error", thisFd)
		return nil
	}
	loaded.fd = thisFd

	return fdo
}
