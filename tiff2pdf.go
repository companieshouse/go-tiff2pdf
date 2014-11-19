package tiff2pdf

/*
#cgo CFLAGS: -D_THREAD_SAFE -pthread -I../../vadz/libtiff/libtiff
#cgo LDFLAGS: -lm
#include <stdio.h>
#include <stdlib.h>
#include <math.h>
#include "c/libtiff.h"
#include "c/tiff2pdf.c"
#include "c/cgo.c"
*/
import "C"
import "errors"

var (
	ErrOpenFailed = errors.New("Opening TIFF failed")
)

func ConvertTiffToPDF(tiff []byte) ([]byte, error) {
	name := C.CString("name")
	mode := C.CString("mode")

	tif := C.TIFFFdOpen(-1, name, mode)
	if tif == nil {
		return nil, ErrOpenFailed
	}

	return tiff, nil
}
