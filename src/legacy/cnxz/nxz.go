package cnxz

/*
#include <stdint.h>
#include <stdlib.h>

void* nxz_compress_new();
void nxz_compress_free(void* p);
int nxz_compress(void* a1p, uint8_t* a2p, uint8_t* a3p, int a4p);
*/
import "C"
import (
	"encoding/binary"
	"errors"
	"io"
	"unsafe"

	"github.com/opennox/libs/ifs"

	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

// only works on 32bit
var _ = [1]struct{}{}[unsafe.Sizeof(int(0))-4]

func DecompressFile(src, dst string) error {
	if src == "" {
		return errors.New("empty source path")
	}
	if dst == "" {
		return errors.New("empty destination path")
	}
	r, err := ifs.Open(src)
	if err != nil {
		return err
	}
	defer r.Close()
	fi, err := r.Stat()
	if err != nil {
		return err
	}
	var buf [4]byte
	if _, err = io.ReadFull(r, buf[:]); err != nil {
		return err
	}
	dstSize := binary.LittleEndian.Uint32(buf[:])
	if uint64(dstSize) > uint64(^uint(0)>>1) || fi.Size()-4 > int64(^uint(0)>>1) {
		return errors.New("nxz: file exceeds addressable size")
	}
	sbuf := make([]byte, int(fi.Size()-4))
	if _, err = io.ReadFull(r, sbuf); err != nil {
		return err
	}
	dbuf := make([]byte, int(dstSize))
	if err := newMapDecoder(sbuf).decode(dbuf); err != nil {
		return err
	}

	w, err := ifs.Create(dst)
	if err != nil {
		return err
	}
	defer w.Close()
	if _, err = w.Write(dbuf); err != nil {
		return err
	}
	return w.Close()
}

func compBufferSize(sz int) int {
	return sz + sz/2 + 32
}

func CompressFile(src, dst string) error {
	if src == "" {
		return errors.New("empty source path")
	}
	if dst == "" {
		return errors.New("empty destination path")
	}
	r, err := ifs.Open(src)
	if err != nil {
		return err
	}
	defer r.Close()
	fi, err := r.Stat()
	if err != nil {
		return err
	}

	srcSz := int(fi.Size())
	sbuf, sfree := alloc.Make([]byte{}, srcSz)
	defer sfree()
	if _, err = io.ReadFull(r, sbuf); err != nil {
		return err
	}

	dbuf, dfree := alloc.Make([]byte{}, compBufferSize(srcSz))
	defer dfree()

	ptr := C.nxz_compress_new()
	cnt := 0
	for i := 0; i < srcSz; i += 500000 {
		v := srcSz - i
		if v > 500000 {
			v = 500000
		}
		cnt += int(C.nxz_compress(ptr,
			(*C.uchar)(unsafe.Pointer(&dbuf[cnt])),
			(*C.uchar)(unsafe.Pointer(&sbuf[i])),
			C.int(v),
		))
	}
	C.nxz_compress_free(ptr)

	w, err := ifs.Create(dst)
	if err != nil {
		return err
	}
	defer w.Close()

	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:4], uint32(srcSz))
	if _, err := w.Write(buf[:4]); err != nil {
		return err
	}
	if _, err := w.Write(dbuf[:cnt]); err != nil {
		return err
	}
	return w.Close()
}
