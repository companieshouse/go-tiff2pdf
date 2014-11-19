package tiff2pdf

import "C"
import (
	"log"
	"unsafe"
)

//export CallGo
func CallGo() {
	log.Println("GO CALLED!")
}

//export GoTiffReadProc
func GoTiffReadProc(fd int, ptr unsafe.Pointer, size int64) int64 {
	log.Println("GoTiffReadProc!")
	return -1
}

//export GoTiffWriteProc
func GoTiffWriteProc(fd int, ptr unsafe.Pointer, size int64) int64 {
	log.Println("GoTiffWriteProc!")
	return -1
}

//export GoTiffSeekProc
func GoTiffSeekProc(fd int, offset int64, whence int) int64 {
	log.Println("GoTiffSeekProc!")
	return -1
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
