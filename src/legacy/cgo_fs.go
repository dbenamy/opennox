package legacy

/*
#include <stdbool.h>
#include <stdio.h>
*/
import "C"
import (
	"io"

	"sync"
	"unsafe"

	"github.com/opennox/opennox/v1/internal/binfile"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
)

var files struct {
	sync.RWMutex
	byHandle map[unsafe.Pointer]*binfile.File
}

type FILE = C.FILE

func convWhence(mode int) int {
	var whence int
	switch C.int(mode) {
	case C.SEEK_SET:
		whence = io.SeekStart
	case C.SEEK_CUR:
		whence = io.SeekCurrent
	case C.SEEK_END:
		whence = io.SeekEnd
	default:
		panic("unsupported seek mode")
	}
	return whence
}

func nox_fs_fread(f *FILE, dst unsafe.Pointer, sz int) int {
	fp := fileByHandle(f)
	n, _ := fp.Read(unsafe.Slice((*byte)(dst), sz))
	return n
}

func fileByHandle(f *FILE) *binfile.File {
	h := unsafe.Pointer(f)
	handles.AssertValidPtr(h)
	files.RLock()
	fp := files.byHandle[h]
	files.RUnlock()
	return fp
}

func nox_fs_close(f *FILE) {
	if f == nil {
		return
	}
	h := unsafe.Pointer(f)
	handles.AssertValidPtr(h)
	files.Lock()
	defer files.Unlock()
	fp := files.byHandle[h]
	if fp != nil {
		_ = fp.Close()
		delete(files.byHandle, h)
	}
}

func Nox_fs_close(f *FILE) {
	nox_fs_close(f)
}

func NewFileHandle(f *binfile.File) *FILE {
	if f.Handle != nil {
		return (*FILE)(f.Handle)
	}
	f.Handle = handles.NewPtr()
	files.Lock()
	defer files.Unlock()
	if files.byHandle == nil {
		files.byHandle = make(map[unsafe.Pointer]*binfile.File)
	}
	files.byHandle[f.Handle] = f
	return (*FILE)(f.Handle)
}
