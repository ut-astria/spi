package sgp4

// Source: github.com/elliotchance/c2go

import (
	"unsafe"
)

func Strncpy64(dest, src *byte, len int64) *byte {
	// Copy up to the len or first NULL bytes - whichever comes first.
	var (
		pSrc  = src
		pDest = dest
		i     int64
	)
	for i < len && *pSrc != 0 {
		*pDest = *pSrc
		i++
		pSrc = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(src)) + uintptr(i)))
		pDest = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(dest)) + uintptr(i)))
	}

	// The rest of the dest will be padded with zeros to the len.
	for i < len {
		*pDest = 0
		i++
		pDest = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(dest)) + uintptr(i)))
	}

	return dest
}
