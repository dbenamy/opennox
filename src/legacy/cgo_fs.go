package legacy

import (
	"sync"
	"unsafe"

	"github.com/opennox/opennox/v1/internal/binfile"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
)

var files struct {
	sync.RWMutex
	byHandle map[unsafe.Pointer]*binfile.File
}

// FILE identifies a file in the protected handle region. It is never dereferenced
// and does not contain libc FILE storage. Keep it non-zero-sized for pointer identity.
type FILE struct{ _ byte }

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
