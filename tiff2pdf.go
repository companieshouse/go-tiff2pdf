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
import "errors"

var (
	ErrOpenFailed = errors.New("Opening TIFF failed")
)

var fdCount int
var fdMap = make(map[int][]byte)

func ConvertTiffToPDF(tiff []byte) ([]byte, error) {
	name := C.CString("test.tif")
	mode := C.CString("r")

	fdCount += 1
	fd := fdCount
	fdMap[fd] = tiff

	tif := C.TIFFFdOpen(C.int(fd), name, mode)
	if tif == nil {
		return nil, ErrOpenFailed
	}

	return tiff, nil
}
