package tiff2pdf

import "C"
import (
	"log"
	"reflect"
	"unsafe"
)

//export CallGo
func CallGo() {
	log.Println("GO CALLED!")
}

const (
	SEEK_SET = iota
	SEEK_CUR
	SEEK_END
)

//export GoTiffReadProc
func GoTiffReadProc(fd int, ptr unsafe.Pointer, size int) int {
	log.Printf("GoTiffReadProc off[%d] size[%d]!\n", fdMap[fd].offset, size)
	hdr := reflect.SliceHeader{
		Data: uintptr(ptr),
		Len:  size,
		Cap:  size,
	}
	goSlice := *(*[]byte)(unsafe.Pointer(&hdr))

	for i := int64(0); i < int64(size) ; i++ {
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
		// Len:  2*int(unsafe.Sizeof(ptr)),
		// Cap:  2*int(unsafe.Sizeof(ptr)),
	}
	goSlice := *(*[]byte)(unsafe.Pointer(&hdr))
	log.Printf("[%d] GoTiffWriteProc off[%d] capBuff[%d] size[%d] slcLen[%d]!\n", fd, fdMap[fd].offset, cap(fdMap[fd].buffer), size, len(goSlice))

	/*
	size64 := int64(size)
	curCap := int64(cap(fdMap[fd].buffer))
	needCap := fdMap[fd].offset+size64
	if needCap > curCap {
		log.Printf("[%d]                 off[%d] size[%d] slcLen[%d] need-cur[%d-%d]\n", fd, fdMap[fd].offset, size, len(goSlice), needCap, curCap)
		fdMap[fd].buffer = append(fdMap[fd].buffer, goSlice[0:needCap-curCap-1]...)
	}
	*/
	// fdMap[fd].buffer[fdMap[fd].offset:fdMap[fd].offset+size64] = goSlice[0:]
	// fdMap[fd].offset += len(goSlice)
	// fdMap[fd].offset += size64

	// return size

	for i := 0; i < size ; i++ {
		// log.Printf("[%d] GoTiffWriteProc off[%d] size[%d] i[%d] cap[%d] len[%d]!\n", fd, fdMap[fd].offset, size, i, cap(fdMap[fd].buffer), len(fdMap[fd].buffer))
		if i >= len(goSlice) {
			log.Printf("[%d]  DONE--------   off[%d] size[%d] i[%d]!\n", fd, fdMap[fd].offset, size, i)
			return int(i)
		}
		if fdMap[fd].offset >= int64(len(fdMap[fd].buffer)) {
			// log.Printf("[%d]  append         off[%d] size[%d] i[%d]!\n", fd, fdMap[fd].offset, size, i)
			fdMap[fd].buffer = append(fdMap[fd].buffer, goSlice[i])
		} else {
			log.Printf("[%d]  copy           off[%d] size[%d] i[%d]!\n", fd, fdMap[fd].offset, size, i)
			fdMap[fd].buffer[fdMap[fd].offset] = goSlice[i]
		}
		fdMap[fd].offset++
	}

	log.Printf("[%d]                 off[%d] capBuff[%d]\n", fd, fdMap[fd].offset, cap(fdMap[fd].buffer))
	return size
}

//export GoTiffSeekProc
func GoTiffSeekProc(fd int, offset int64, whence int) int64 {
	log.Printf("[%d] GoTiffSeekProc! off[%d] off[%d] wh[%d]", fd, fdMap[fd].offset, offset, whence)
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
		newOffset = int64(len(fdMap[fd].buffer))-offset
	}
	if newOffset < 0 {
		log.Printf("[%d] GoTiffSeekProc off[%d] len[%d]", fd, newOffset, len(fdMap[fd].buffer))
		return -1
	} else if newOffset > int64(len(fdMap[fd].buffer)) {
		log.Printf("[%d] GoTiffSeekProc off[%d] len[%d]", fd, newOffset, len(fdMap[fd].buffer))
		for int64(len(fdMap[fd].buffer)) < newOffset {
			fdMap[fd].buffer = append(fdMap[fd].buffer, '\000')
		}
	}
	fdMap[fd].offset = newOffset
	return fdMap[fd].offset
}

//export GoTiffCloseProc
func GoTiffCloseProc(fd int) int {
	log.Printf("[%d] GoTiffCloseProc!", fd)
	return -1
}

/* These probably aren't needed... */

//export GoTiffMapProc
func GoTiffMapProc(fd int, base unsafe.Pointer, size int64) int {
	log.Printf("[%d] GoTiffMapProc!", fd)
	return 0
}

//export GoTiffUnmapProc
func GoTiffUnmapProc(fd int, base unsafe.Pointer, size int64) {
	log.Printf("[%d] GoTiffUnmapProc!", fd)
}

//export GoOutputDisable
func GoOutputDisable(fd int) {
	fdMap[fd].outputdisable = 1
}

//export GoOutputEnable
func GoOutputEnable(fd int) {
	fdMap[fd].outputdisable = 0
}
