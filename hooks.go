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
	log.Printf("GoTiffWriteProc off[%d] size[%d]!\n", fdMap[fd].offset, size)
	hdr := reflect.SliceHeader{
		Data: uintptr(ptr),
		Len:  2*int(unsafe.Sizeof(ptr)),
		Cap:  2*int(unsafe.Sizeof(ptr)),
	}
	goSlice := *(*[]byte)(unsafe.Pointer(&hdr))

	if fdMap[fd].offset+size >= int64(cap(fdMap[fd].buffer)) {
		fdMap[fd].buffer = append(fdMap[fd].buffer, goSlice[0:]...)
	// fdMap[fd].offset += len(goSlice)

	for i := 0; i < size ; i++ {
		log.Printf("GoTiffWriteProc off[%d] size[%d] i[%d] cap[%d] len[%d]!\n", fdMap[fd].offset, size, i, cap(fdMap[fd].buffer), len(fdMap[fd].buffer))
		if i >= len(goSlice) {
			log.Printf(" OUT            off[%d] size[%d] i[%d]!\n", fdMap[fd].offset, size, i)
			return int(i)
		}
		log.Printf("                off[%d] size[%d] i[%d]!\n", fdMap[fd].offset, size, i)
		// if fdMap[fd].offset >= int64(cap(fdMap[fd].buffer)) {
			// log.Printf(" a              off[%d] size[%d] i[%d]!\n", fdMap[fd].offset, size, i)
			fdMap[fd].buffer = append(fdMap[fd].buffer, goSlice[i])
		// } else {
			// log.Printf(" b              off[%d] size[%d] i[%d]!\n", fdMap[fd].offset, size, i)
			// fdMap[fd].buffer[fdMap[fd].offset] = goSlice[i]
		// }
		fdMap[fd].offset++
	}

	return size
}

//export GoTiffSeekProc
func GoTiffSeekProc(fd int, offset int64, whence int) int64 {
	log.Println("GoTiffSeekProc!")
	newOffset := fdMap[fd].offset
	switch whence {
	case SEEK_SET:
		newOffset = offset
	case SEEK_CUR:
		newOffset += offset
	case SEEK_END:
		newOffset = int64(len(fdMap[fd].buffer))-offset
	}
	if newOffset < 0 || newOffset > int64(len(fdMap[fd].buffer)) {
		return -1
	}
	fdMap[fd].offset = newOffset
	return fdMap[fd].offset
}

//export GoTiffCloseProc
func GoTiffCloseProc(fd int) int {
	log.Println("GoTiffCloseProc!")
	return -1
}

/* These probably aren't needed... */

//export GoTiffMapProc
func GoTiffMapProc(fd int, base unsafe.Pointer, size int64) int {
	log.Println("GoTiffMapProc!")
	return 0
}

//export GoTiffUnmapProc
func GoTiffUnmapProc(fd int, base unsafe.Pointer, size int64) {
	log.Println("GoTiffUnmapProc!")
}
