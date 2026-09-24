package legacy

/*
#include <stdbool.h>
#include <stdio.h>
*/
import "C"
import (
	"io"
	"os"
	"sync"
	"unsafe"

	"github.com/opennox/libs/ifs"

	"github.com/opennox/opennox/v1/internal/binfile"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
)

var files struct {
	sync.RWMutex
	byHandle map[unsafe.Pointer]*binfile.File
}

type FILE = C.FILE

//export nox_fs_set_workdir
func nox_fs_set_workdir(path *C.char) C.bool {
	return ifs.Chdir(GoString(path)) == nil
}

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

//export nox_fs_fseek
func nox_fs_fseek(f *FILE, off C.long, mode int) int {
	fp := fileByHandle(f)
	_, err := fp.Seek(int64(off), convWhence(mode))
	if err != nil {
		return -1
	}
	return 0
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

//export nox_fs_open
func nox_fs_open(path *C.char) *FILE {
	f, err := ifs.Open(GoString(path))
	if err != nil {
		return nil
	}
	return NewFileHandle(binfile.NewFile(f))
}

//export nox_fs_open_text
func nox_fs_open_text(path *C.char) *FILE {
	f, err := ifs.Open(GoString(path))
	if err != nil {
		return nil
	}
	return NewFileHandle(binfile.NewTextFile(f))
}

//export nox_fs_create
func nox_fs_create(path *C.char) *FILE {
	f, err := ifs.Create(GoString(path))
	if err != nil {
		return nil
	}
	return NewFileHandle(binfile.NewFile(f))
}

//export nox_fs_create_text
func nox_fs_create_text(path *C.char) *FILE {
	f, err := ifs.Create(GoString(path))
	if err != nil {
		return nil
	}
	return NewFileHandle(binfile.NewTextFile(f))
}

//export nox_fs_open_rw
func nox_fs_open_rw(path *C.char) *FILE {
	f, err := ifs.OpenFile(GoString(path), os.O_RDWR)
	if err != nil {
		return nil
	}
	return NewFileHandle(binfile.NewFile(f))
}
