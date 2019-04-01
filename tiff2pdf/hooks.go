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

	loaded, ok := fdMap[fd]
	if !ok {
		log.Printf("[%d] GoTiffReadProc load error", fd)
		return -1
	}
	for i := int64(0); i < int64(size); i++ {
		if loaded.offset >= int64(len(loaded.buffer)) {
			return int(i)
		}
		goSlice[i] = loaded.buffer[loaded.offset]
		loaded.offset++
	}

	return size
}

//export GoTiffWriteProc
func GoTiffWriteProc(fd int, ptr unsafe.Pointer, size int) int {
	loaded, ok := fdMap[fd]
	if !ok {
		log.Printf("[%d] GoTiffWriteProc load error", fd)
		return -1
	}
	if loaded.outputdisable == 1 {
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
		if loaded.offset >= int64(len(loaded.buffer)) {
			loaded.buffer = append(loaded.buffer, goSlice[i])
		} else {
			loaded.buffer[loaded.offset] = goSlice[i]
		}
		loaded.offset++
	}

	return size
}

//export GoTiffSeekProc
func GoTiffSeekProc(fd int, offset int64, whence int) int64 {
	loaded, ok := fdMap[fd]
	if !ok {
		log.Printf("[%d] GoTiffSeekProc load error", fd)
		return -1
	}
	if loaded.outputdisable == 1 {
		return offset
	}
	newOffset := loaded.offset
	switch whence {
	case SEEK_SET:
		newOffset = offset
	case SEEK_CUR:
		newOffset += offset
	case SEEK_END:
		newOffset = int64(len(loaded.buffer)) - offset
	}
	if newOffset < 0 {
		return -1
	} else if newOffset > int64(len(loaded.buffer)) {
		for int64(len(loaded.buffer)) < newOffset {
			loaded.buffer = append(loaded.buffer, '\000')
		}
	}
	loaded.offset = newOffset
	return loaded.offset
}

//export GoTiffCloseProc
func GoTiffCloseProc(fd int) int {
	return -1
}

//export GoTiffSizeProc
func GoTiffSizeProc(fd int) int {
	loaded, ok := fdMap[fd]
	if !ok {
		log.Printf("[%d] GoTiffSizeProc load error", fd)
		return -1
	}
	return len(loaded.buffer)
}

//export GoOutputDisable
func GoOutputDisable(fd int) {
	loaded, ok := fdMap[fd]
	if !ok {
		log.Printf("[%d] GoOutputDisable load error", fd)
		return
	}
	loaded.outputdisable = 1
}

//export GoOutputEnable
func GoOutputEnable(fd int) {
	loaded, ok := fdMap[fd]
	if !ok {
		log.Printf("[%d] GoOutputEnable load error", fd)
		return
	}
	loaded.outputdisable = 0
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

	loaded, ok := fdMap[fd]
	if !ok {
		// TODO don't think we care about warnings with fd 0
		log.Printf("[%d] WARNING: %s", fd, s)
		return
	}
	loaded.warnings = append(fdMap[fd].warnings, s)
}

//export GoTiffErrorExt
func GoTiffErrorExt(fd int, err *C.char) {
	s := C.GoString(err)

	loaded, ok := fdMap[fd]
	if !ok {
		// TODO don't think we care about errors with fd 0
		log.Printf("[%d] ERROR: %s", fd, s)
		return
	}
	loaded.errors = append(fdMap[fd].errors, s)
}
