//go:build safe && porttest

package legacy

/*
void* nox_memcpy(void* dst, const void* src, unsigned int size);
int nox_memcmp(const void* a, const void* b, unsigned int size);
unsigned int nox_strlen(const char* s);
char* nox_strcpy(char* dst, const char* src);
char* nox_strcat(char* dst, const char* src);
int nox_strcmp(const char* a, const char* b);
*/
import "C"
import (
	"bytes"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

type PortTestSafeMemorySpec struct {
	Op                      string
	Left, Right             []byte
	LeftOff, RightOff, Size int
}

type PortTestSafeMemoryResult struct {
	Left, Right  []byte
	Return       int32
	ReturnedDest bool
}

// PortTestSafeMemory invokes the actual safe-profile C entrypoints. Both buffers
// are foreign allocations; pointer returns are normalized only by exact equality.
func PortTestSafeMemory(s PortTestSafeMemorySpec) PortTestSafeMemoryResult {
	const payload, guard = 64, 16
	if s.LeftOff < 0 || s.LeftOff >= payload || s.RightOff < 0 || s.RightOff >= payload || s.Size < 0 ||
		len(s.Left)+s.LeftOff > payload || len(s.Right)+s.RightOff > payload {
		panic("invalid safe-memory fixture bounds")
	}
	left, freeLeft := alloc.Make([]byte(nil), payload+2*guard)
	defer freeLeft()
	right, freeRight := alloc.Make([]byte(nil), payload+2*guard)
	defer freeRight()
	for i := range left {
		left[i] = 0xa5
		right[i] = 0x5a
	}
	l, r := guard+s.LeftOff, guard+s.RightOff
	copy(left[l:], s.Left)
	copy(right[r:], s.Right)
	lp, rp := unsafe.Pointer(&left[l]), unsafe.Pointer(&right[r])
	stringLen := func(b []byte) int {
		n := bytes.IndexByte(b, 0)
		if n < 0 {
			panic("unterminated safe-memory fixture string")
		}
		return n
	}
	out := PortTestSafeMemoryResult{}
	switch s.Op {
	case "memcpy", "memcmp":
		if s.Size > payload-s.LeftOff || s.Size > payload-s.RightOff {
			panic("invalid safe-memory fixture span")
		}
		if s.Op == "memcpy" {
			out.ReturnedDest = C.nox_memcpy(lp, rp, C.uint(s.Size)) == lp
		} else {
			out.Return = int32(C.nox_memcmp(lp, rp, C.uint(s.Size)))
		}
	case "strlen":
		stringLen(left[l : guard+payload])
		out.Return = int32(C.nox_strlen((*C.char)(lp)))
	case "strcpy":
		n := stringLen(right[r : guard+payload])
		if n+1 > payload-s.LeftOff {
			panic("insufficient safe-memory fixture copy capacity")
		}
		out.ReturnedDest = unsafe.Pointer(C.nox_strcpy((*C.char)(lp), (*C.char)(rp))) == lp
	case "strcat":
		nl, nr := stringLen(left[l:guard+payload]), stringLen(right[r:guard+payload])
		if nl+nr+1 > payload-s.LeftOff {
			panic("insufficient safe-memory fixture append capacity")
		}
		out.ReturnedDest = unsafe.Pointer(C.nox_strcat((*C.char)(lp), (*C.char)(rp))) == lp
	case "strcmp":
		stringLen(left[l : guard+payload])
		stringLen(right[r : guard+payload])
		out.Return = int32(C.nox_strcmp((*C.char)(lp), (*C.char)(rp)))
	default:
		panic("unknown safe-memory fixture operation")
	}
	out.Left, out.Right = bytes.Clone(left), bytes.Clone(right)
	return out
}
