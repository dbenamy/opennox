package legacy

/*
_Static_assert(sizeof(int) == 4, "protection ABI requires 32-bit int");
_Static_assert(sizeof(unsigned int) == 4, "protection ABI requires 32-bit unsigned int");
extern unsigned int dword_5d4594_2516356;
int nox_xxx_protectionCreateStructForInt_56F280(int a1, int a2);
int nox_xxx_protectionCreateStructForFloat_56F480(int a1, float a2);
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/opennox/v1/internal/protection"
)

func Nox_xxx_protectionCreateStructForInt_56F280(a1 int, a2 int) int {
	return int(C.nox_xxx_protectionCreateStructForInt_56F280(C.int(a1), C.int(a2)))
}
func Nox_xxx_protectionCreateStructForFloat_56F480(a1 int, a2 float32) int {
	return int(C.nox_xxx_protectionCreateStructForFloat_56F480(C.int(a1), C.float(a2)))
}

//export nox_xxx_protectionStringCRC_56FAC0
func nox_xxx_protectionStringCRC_56FAC0(data *C.int, size C.uint) C.int {
	return nox_xxx_protectionStringCRCLen_56FAE0(data, size)
}

//export nox_xxx_protectionStringCRCLen_56FAE0
func nox_xxx_protectionStringCRCLen_56FAE0(data *C.int, size C.uint) C.int {
	if data == nil {
		return 0
	}
	// Process complete words only. Chunking avoids converting an unsigned C
	// byte count greater than MaxInt into a negative Go slice length on 386.
	remaining := uint32(size) &^ 3
	p := unsafe.Pointer(data)
	var sum uint32
	for remaining != 0 {
		n := remaining
		if n > 1<<20 {
			n = 1 << 20
		}
		sum ^= protection.Checksum(unsafe.Slice((*byte)(p), int(n)))
		remaining -= n
		if remaining != 0 {
			p = unsafe.Add(p, n)
		}
	}
	return C.int(sum)
}
