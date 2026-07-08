package tiff2pdf

/*
#cgo CFLAGS: -D_THREAD_SAFE -pthread -I../libtiff-src/libtiff
#cgo LDFLAGS: -lm
#include <stdio.h>
#include <stdlib.h>
#include <math.h>

// Forward declarations of the Go callbacks exported from hooks.go. They are
// used by the C files included below (tif_golang.c, and the GoOutput* hooks
// spliced into tiff2pdf.c by the Makefile). The generated "_cgo_export.h"
// cannot be included from this package's own preamble, so declare them here
// with cgo's ABI types (Go int -> long long).
extern long long GoTiffReadProc(long long fd, void *ptr, long long size);
extern long long GoTiffWriteProc(long long fd, void *ptr, long long size);
extern long long GoTiffSeekProc(long long fd, long long offset, long long whence);
extern long long GoTiffCloseProc(long long fd);
extern long long GoTiffSizeProc(long long fd);
extern void GoOutputDisable(long long fd);
extern void GoOutputEnable(long long fd);
extern void GoTiffWarningExt(long long fd, char *err);
extern void GoTiffErrorExt(long long fd, char *err);

#include "c/libtiff.h"
#include "c/tiff2pdf.c"
#include "c/tif_golang.c"
*/
import "C"
import (
	"errors"
	"unsafe"
)

//Config represents the command line tiff2pdf configuration
type Config struct {
	// PageSize sets the PDF page size, e.g. legal, letter or A4
	PageSize string
	// FullPage makes the tiff image fill the PDF page
	FullPage bool

	// Creator is the image software used to create the document
	Creator string

	// Author is the document name
	Author string
	// Subject is the document description
	Subject string
	// Title is the document name
	Title string
}

// ConvertTiffToPDFOutput is returned by ConvertTiffToPDF
type ConvertTiffToPDFOutput struct {
	// PageCount is the number of pages converted from TIFF to PDF
	PageCount uint
	// PDF is the output from tiff2pdf
	PDF []byte
	// Errors contains any errors reported by tiff2pdf
	Errors []string
	// Warnings contains any warnings reported by tiff2pdf
	Warnings []string
}

// DefaultConfig creates the default tiff2pdf configuration
func DefaultConfig() *Config {
	return &Config{
		PageSize: "A4",
		FullPage: true,
		Creator:  "go-tiff2pdf",
	}
}

func createTiff(tiff []byte, name, mode string) (*C.TIFF, error) {
	newFd := NewFd(tiff)

	cName := C.CString(name)
	cMode := C.CString(mode)
	tif := C.GoTIFFFdOpen(C.int(newFd.fd), cName, cMode)
	C.free(unsafe.Pointer(cName))
	C.free(unsafe.Pointer(cMode))

	if tif == nil {
		return nil, ErrOpenFailed
	}
	return tif, nil
}

func configureT2p(t2p *C.T2P, config *Config) {
	cPageSize := C.CString(config.PageSize)
	if r := C.tiff2pdf_match_paper_size(&t2p.pdf_defaultpagewidth, &t2p.pdf_defaultpagelength, cPageSize); r != 0 {
		t2p.pdf_overridepagesize = 1
	} else {
		// TODO warning?
	}
	C.free(unsafe.Pointer(cPageSize))

	if config.FullPage {
		t2p.pdf_image_fillpage = 1
	} else {
		t2p.pdf_image_fillpage = 0
	}

	// FIXME if len(config.Creator) == 0, is that "no flag" or "empty string"
	t2p.pdf_creator = stringTo512Cchar(config.Creator)
	t2p.pdf_author = stringTo512Cchar(config.Author)
	t2p.pdf_subject = stringTo512Cchar(config.Subject)
	t2p.pdf_title = stringTo512Cchar(config.Title)
}

func stringTo512Cchar(s string) [512]C.char {
	var cArr [512]C.char
	for i, c := range s {
		cArr[i] = C.char(c)
	}
	cArr[len(s)] = C.char(0)
	return cArr
}

// ConvertTiffToPDF converts an input TIFF byte array to an output PDF byte array
func ConvertTiffToPDF(tiff []byte, config *Config, inputName string, outputName string) (*ConvertTiffToPDFOutput, error) {
	input, err := createTiff(tiff, inputName, "rw")
	if err != nil {
		return nil, err
	}
	input_fd := int(input.tif_fd)
	defer func() {
		fdMap.Delete(input_fd)
	}()

	output, err := createTiff([]byte{}, outputName, "w")
	if err != nil {
		return nil, err
	}
	output_fd := int(output.tif_fd)
	defer func() {
		fdMap.Delete(output_fd)
	}()

	GoTiffSeekProc(output_fd, 0, 0)

	t2p := C.t2p_init()
	defer C.t2p_free(t2p)
	if t2p == nil {
		return nil, errors.New("Error: t2p_init") // FIXME capture FD0
	}

	configureT2p(t2p, config)

	C.t2p_write_pdf(t2p, input, output)
	if t2p.t2p_error != 0 {
		return nil, errors.New("t2p_error") // FIXME capture FD0
	}

	loaded, ok := fdMap.Load(int(output.tif_fd))
	if !ok {
		return nil, errors.New("t2p_error loading from map")
	}
	out := &ConvertTiffToPDFOutput{
		uint(t2p.tiff_pagecount),
		loaded.buffer,
		loaded.errors,
		loaded.warnings,
	}
	return out, nil
}
