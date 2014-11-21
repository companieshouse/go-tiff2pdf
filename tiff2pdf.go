package tiff2pdf

/*
#cgo CFLAGS: -D_THREAD_SAFE -pthread -I../../vadz/libtiff/libtiff
#cgo LDFLAGS: -lm
#include <stdio.h>
#include <stdlib.h>
#include <math.h>
#include "c/libtiff.h"
#include "c/tiff2pdf.c"
#include "c/tif_golang.c"
*/
import "C"
import (
	"errors"
	"log"
	"sync"
)

var (
	ErrOpenFailed = errors.New("Opening TIFF failed")
)

type fd struct {
	fd int
	buffer []byte
	offset int64
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
		fd: thisFd,
		buffer: buffer,
	}
	fdMap[thisFd] = fdo

	return fdo
}

func createTiff(tiff []byte, name, mode string) (*C.TIFF, error) {
	newFd := NewFd(tiff)
	tif := C.TIFFFdOpen(C.int(newFd.fd), C.CString(name), C.CString(mode))
	if tif == nil {
		return nil, ErrOpenFailed
	}
	return tif, nil
}

func ConvertTiffToPDF(tiff []byte, inputName string, outputName string) ([]byte, error) {
	input, err := createTiff(tiff, inputName, "rw")
	if err != nil {
		return nil, err
	}

	output, err := createTiff([]byte{}, outputName, "w")
	if err != nil {
		return nil, err
	}
	GoTiffSeekProc(int(output.tif_fd), 0, 0)

	t2p := C.t2p_init()
	if t2p == nil {
		return nil, errors.New("Error: t2p_init!")
	}
	// t2p.outputfile = C.FILE(output.tif_fd)
	C.t2p_write_pdf(t2p, input, output)
	if t2p.t2p_error != 0 {
		C.t2p_free(t2p)
		return nil, errors.New("t2p_error")
	}

	C.t2p_free(t2p)

	return fdMap[int(output.tif_fd)].buffer, nil
}
