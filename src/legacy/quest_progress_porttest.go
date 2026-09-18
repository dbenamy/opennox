//go:build porttest

package legacy

/*
#include "GAME4.h"
#include "GAME4_1.h"
extern uint32_t dword_5d4594_1570272;
void sub_500510(const char*);
int sub_51A920(int);
void sub_51A1F0(int);
*/
import "C"

import (
	"bytes"
	"math"
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

// Own the actual C list, mapped namespace/scratch buffers and separator bytes.
func PortTestQuestProgressOwner() func() {
	old := C.dword_5d4594_1570272
	C.dword_5d4594_1570272 = 0
	type saved struct{ dst, data []byte }
	var regions []saved
	for _, r := range [][3]uintptr{{0x5D4594, 1570008, 264}, {0x587000, 217952, 2}, {0x587000, 217960, 2}, {0x5D4594, 2388656, 20}} {
		b := unsafe.Slice((*byte)(memmap.PtrOff(r[0], r[1])), int(r[2]))
		regions = append(regions, saved{b, bytes.Clone(b)})
		clear(b)
	}
	*memmap.PtrUint16(0x587000, 217952) = ':'
	*memmap.PtrUint16(0x587000, 217960) = ':'
	return func() {
		PortTestQuestProgress("reset", "*:*", 0)
		C.dword_5d4594_1570272 = old
		for _, r := range regions {
			copy(r.dst, r.data)
		}
	}
}

func PortTestQuestProgress(op, name string, value uint32) uint64 {
	str, free := alloc.CString(name)
	defer free()
	p := (*C.char)(unsafe.Pointer(str))
	var ptr unsafe.Pointer
	switch op {
	case "namespace":
		C.sub_500510(p)
	case "namespace-nil":
		C.sub_500510(nil)
	case "set-int":
		ptr = unsafe.Pointer(C.nox_xxx_journalQuestSet_500540(p, C.int(value)))
	case "set-float":
		ptr = unsafe.Pointer(C.nox_xxx_journalQuestSetBool_5006B0(p, C.int(value)))
	case "find":
		ptr = unsafe.Pointer(C.nox_xxx_scriptGetJournal_5005E0(p))
	case "int":
		return uint64(uint32(C.sub_500750(p)))
	case "float":
		return math.Float64bits(float64(C.sub_500770(p)))
	case "reset":
		C.sub_5007E0(p) // Every production caller ignores the incidental return.
	case "qualify":
		return uint64(C.sub_5009B0(p))
	case "write":
		return uint64(C.sub_500A60())
	case "read":
		return uint64(C.sub_500B70())
	case "stage-set":
		return uint64(uint32(C.sub_51A920(C.int(value))))
	case "stage":
		return uint64(uint32(C.nox_xxx_getQuestStage_51A930()))
	case "minions-set":
		return uint64(uint32(C.sub_51A940(C.int(value))))
	default:
		panic(op)
	}
	if ptr == nil {
		return 0
	}
	for i, n := uint64(1), uintptr(C.dword_5d4594_1570272); n != 0; i++ {
		if unsafe.Pointer(n) == ptr {
			return i
		}
		n = uintptr(*(*uint32)(unsafe.Pointer(n + 140)))
	}
	panic("quest return outside list")
}
func PortTestQuestProgressScratch() []byte {
	return bytes.Clone(unsafe.Slice((*byte)(memmap.PtrOff(0x5D4594, 1570140)), 132))
}
func PortTestQuestProgressSnapshot() [][37]uint32 {
	var out [][37]uint32
	seen := map[uintptr]bool{}
	previous := uintptr(0)
	for n := uintptr(C.dword_5d4594_1570272); n != 0; {
		if seen[n] {
			panic("quest list cycle")
		}
		seen[n] = true
		w := *(*[37]uint32)(unsafe.Pointer(n))
		if uintptr(w[36]) != previous {
			panic("quest predecessor mismatch")
		}
		next := uintptr(w[35])
		w[35] = 0
		w[36] = 0
		if next != 0 {
			w[35] = uint32(len(out) + 2)
		}
		if previous != 0 {
			w[36] = uint32(len(out))
		}
		out = append(out, w)
		previous = n
		n = next
	}
	return out
}
func PortTestQuestProgressObject(op string, u *server.Object, stage int) uint32 {
	switch op {
	case "generator-type":
		return uint32(C.sub_51A500(C.int(uintptr(unsafe.Pointer(u)))))
	case "generator-init":
		C.sub_51A550()
	case "prepare":
		C.sub_51A1F0(C.int(stage))
	case "hecubah":
		C.nox_xxx_spawnHecubahQuest_51A5A0((*C.int)(unsafe.Pointer(&u.PosVec)))
	case "necro":
		C.nox_xxx_spawnNecroQuest_51A7A0((*C.int)(unsafe.Pointer(&u.PosVec)))
	default:
		panic(op)
	}
	return 0
}
