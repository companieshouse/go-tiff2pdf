package tiff2pdf

/*
#include <stdarg.h>
*/
import "C"
import (
	"log"
	"reflect"
	"unsafe"
)

const (
	SEEK_SET = iota
	SEEK_CUR
	SEEK_END
)

//export GoTiffReadProc
func GoTiffReadProc(fd int, ptr unsafe.Pointer, size int) int {
	hdr := reflect.SliceHeader{
		Data: uintptr(ptr),
		Len:  size,
		Cap:  size,
	}
	goSlice := *(*[]byte)(unsafe.Pointer(&hdr))

	for i := int64(0); i < int64(size); i++ {
		if fdMap[fd].offset >= int64(len(fdMap[fd].buffer)) {
			return int(i)
		}
		goSlice[i] = fdMap[fd].buffer[fdMap[fd].offset]
		fdMap[fd].offset++
	}

	return size
}

//export GoTiffWriteProc
func GoTiffWriteProc(fd int, ptr unsafe.Pointer, size int) int {
	if fdMap[fd].outputdisable == 1 {
		return size
	}

	hdr := reflect.SliceHeader{
		Data: uintptr(ptr),
		Len:  size,
		Cap:  size,
	}
	goSlice := *(*[]byte)(unsafe.Pointer(&hdr))

	for i := 0; i < size; i++ {
		if i >= len(goSlice) {
			return int(i)
		}
		if fdMap[fd].offset >= int64(len(fdMap[fd].buffer)) {
			fdMap[fd].buffer = append(fdMap[fd].buffer, goSlice[i])
		} else {
			fdMap[fd].buffer[fdMap[fd].offset] = goSlice[i]
		}
		fdMap[fd].offset++
	}

	return size
}

//export GoTiffSeekProc
func GoTiffSeekProc(fd int, offset int64, whence int) int64 {
	if fdMap[fd].outputdisable == 1 {
		return offset
	}
	newOffset := fdMap[fd].offset
	switch whence {
	case SEEK_SET:
		newOffset = offset
	case SEEK_CUR:
		newOffset += offset
	case SEEK_END:
		newOffset = int64(len(fdMap[fd].buffer)) - offset
	}
	if newOffset < 0 {
		return -1
	} else if newOffset > int64(len(fdMap[fd].buffer)) {
		for int64(len(fdMap[fd].buffer)) < newOffset {
			fdMap[fd].buffer = append(fdMap[fd].buffer, '\000')
		}
	}
	fdMap[fd].offset = newOffset
	return fdMap[fd].offset
}

//export GoTiffCloseProc
func GoTiffCloseProc(fd int) int {
	return -1
}

//export GoTiffSizeProc
func GoTiffSizeProc(fd int) int {
	return len(fdMap[fd].buffer)
}

//export GoOutputDisable
func GoOutputDisable(fd int) {
	fdMap[fd].outputdisable = 1
}

//export GoOutputEnable
func GoOutputEnable(fd int) {
	fdMap[fd].outputdisable = 0
}

/* These probably aren't needed... */

//export GoTiffMapProc
func GoTiffMapProc(fd int, base unsafe.Pointer, size int64) int {
	return 0
}

//export GoTiffUnmapProc
func GoTiffUnmapProc(fd int, base unsafe.Pointer, size int64) {
}

//export GoTiffWarningExt
func GoTiffWarningExt(fd int, err *C.char) {
	s := C.GoString(err)
	if _, ok := fdMap[fd]; !ok {
		// TODO don't think we care about warnings with fd 0
		log.Printf("[%d] WARNING: %s", fd, s)
		return
	}
	fdMap[fd].warnings = append(fdMap[fd].warnings, s)
}

//export GoTiffErrorExt
func GoTiffErrorExt(fd int, err *C.char) {
	s := C.GoString(err)
	if _, ok := fdMap[fd]; !ok {
		// TODO don't think we care about errors with fd 0
		log.Printf("[%d] ERROR: %s", fd, s)
		return
	}
	fdMap[fd].errors = append(fdMap[fd].errors, s)
}
