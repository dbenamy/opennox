//go:build porttest

// Package protectionref preserves the C implementation for differential tests.
package protectionref

/*
#cgo CFLAGS: -fno-strict-aliasing
int reference_checksum(int*, unsigned int);
int reference_checksum_nullable(int*, unsigned int);
*/
import "C"
import "unsafe"

func Checksum(data []byte) uint32 {
	return uint32(C.reference_checksum((*C.int)(unsafe.Pointer(unsafe.SliceData(data))), C.uint(len(data))))
}

func Nullable(data []byte) uint32 {
	return uint32(C.reference_checksum_nullable((*C.int)(unsafe.Pointer(unsafe.SliceData(data))), C.uint(len(data))))
}

func NullWithLength(n uint32) uint32 {
	return uint32(C.reference_checksum_nullable(nil, C.uint(n)))
}
