//go:build porttest

package legacy

/*
#include <stdint.h>
#include <stdio.h>
#include "GAME4.h"
#include "GAME4_3.h"
#include "GAME3_2.h"
#include "server__script__file.h"
int nox_script_readInt_542B70(FILE*);
double nox_script_readFloat_542B90(FILE*);
int nox_script_writeFloat_542BD0(float, FILE*);
int nox_script_writeInt_542BB0(int, FILE*);
int nox_script_readWriteYyy_542380(FILE*, FILE*, int);
int nox_script_readWriteJjj_5418C0(FILE*, FILE*, FILE*);
int nox_script_readWriteVvv_541E40(FILE*, FILE*, FILE*);
int nox_script_readWriteIii_541D80(FILE*, FILE*);
int nox_script_readWriteXxx_541A50(FILE*, FILE*, FILE*);
size_t nox_script_readWriteWww_5417C0(FILE*, FILE*, FILE*);
extern unsigned int dword_5d4594_2489420, dword_5d4594_2489424, dword_5d4594_2489428;
*/
import "C"

import (
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"math"
	"unsafe"
)

func PortTestPrefabScriptsGlobals() ([]*uint32, func()) {
	p := []*uint32{(*uint32)(unsafe.Pointer(&C.dword_5d4594_2489420)), (*uint32)(unsafe.Pointer(&C.dword_5d4594_2489424)), (*uint32)(unsafe.Pointer(&C.dword_5d4594_2489428))}
	old := [3]uint32{*p[0], *p[1], *p[2]}
	return p, func() {
		for i := range p {
			*p[i] = old[i]
		}
	}
}

func PortTestPrefabScriptsCall(op int, a, b, c unsafe.Pointer, arg uint32) uint64 {
	switch op {
	case 0:
		return uint64(uint32(C.nox_script_readInt_542B70((*C.FILE)(a))))
	case 1:
		return math.Float64bits(float64(C.nox_script_readFloat_542B90((*C.FILE)(a))))
	case 2:
		return uint64(uint32(C.nox_script_writeFloat_542BD0(C.float(math.Float32frombits(arg)), (*C.FILE)(a))))
	case 3:
		return uint64(uint32(C.nox_script_writeInt_542BB0(C.int(arg), (*C.FILE)(a))))
	case 4:
		return uint64(uint32(C.nox_script_readWriteYyy_542380((*C.FILE)(a), (*C.FILE)(b), C.int(arg))))
	case 5:
		return uint64(uint32(C.nox_script_readWriteJjj_5418C0((*C.FILE)(a), (*C.FILE)(b), (*C.FILE)(c))))
	case 6:
		return uint64(uint32(C.nox_script_readWriteVvv_541E40((*C.FILE)(a), (*C.FILE)(b), (*C.FILE)(c))))
	case 7:
		return uint64(uint32(C.nox_script_readWriteIii_541D80((*C.FILE)(a), (*C.FILE)(b))))
	case 8:
		return uint64(uint32(C.nox_script_readWriteXxx_541A50((*C.FILE)(a), (*C.FILE)(b), (*C.FILE)(c))))
	case 9:
		return uint64(C.nox_script_readWriteWww_5417C0((*C.FILE)(a), (*C.FILE)(b), (*C.FILE)(c)))
	case 10:
		return uint64(uint32(C.sub_543110((*C.char)(a), (*C.int)(b))))
	case 11:
		var xy [2]int32
		if a != nil {
			xy = *(*[2]int32)(a)
		}
		return uint64(uintptr(unsafe.Pointer(C.sub_542BF0(C.int(arg), C.int(xy[0]), C.int(xy[1])))))
	case 14:
		C.sub_509120((*C.uint32_t)(a), C.int(arg), (*C.char)(b))
		return 0
	case 15:
		return uint64(uint32(C.nox_xxx_interesting_xfer_4D0010((*C.uint32_t)(a), C.int(arg))))
	case 16:
		return uint64(uint32(C.sub_4D39F0((*C.char)(a))))
	case 17:
		C.nox_xxx_tileInitdataClear_4D3C50(a)
		return 0
	case 18:
		return uint64(uintptr(unsafe.Pointer(C.sub_4D3C70())))
	case 19:
		return uint64(uintptr(unsafe.Pointer(C.sub_4D3C80((*C.uint32_t)(a)))))
	default:
		panic("prefab script operation")
	}
}

func PortTestPrefabScriptsName(name string, instance, x, y int32, objectOnly bool) string {
	p, free := alloc.CString(name)
	defer free()
	if objectOnly {
		return GoString(C.sub_543620(C.int(uintptr(unsafe.Pointer(p))), C.int(instance)))
	}
	return GoString(C.sub_5435C0(C.int(uintptr(unsafe.Pointer(p))), C.int(instance), C.int(x), C.int(y)))
}

// Observe the real handle registry, and dispose only fixture-created leftovers.
func PortTestPrefabScriptsHandleBalance() (func() int, func()) {
	files.RLock()
	old := make(map[unsafe.Pointer]bool, len(files.byHandle))
	for p := range files.byHandle {
		old[p] = true
	}
	files.RUnlock()
	added := func() []unsafe.Pointer {
		files.RLock()
		defer files.RUnlock()
		var out []unsafe.Pointer
		for p := range files.byHandle {
			if !old[p] {
				out = append(out, p)
			}
		}
		return out
	}
	return func() int { return len(added()) }, func() {
		for _, p := range added() {
			Nox_fs_close((*FILE)(p))
		}
	}
}

// Only scalar paint settings reset by generation; exclude allocation/list owners.
func PortTestPrefabScriptsPaintGlobals() (map[string]*uint32, func()) {
	all := paintGlobals()
	out := map[string]*uint32{}
	old := map[*uint32]uint32{}
	for _, name := range []string{"dword_5d4594_3835348", "dword_5d4594_3835352", "dword_5d4594_3835356", "dword_5d4594_3835360", "dword_5d4594_3835364", "dword_5d4594_3835368", "dword_5d4594_3835372", "dword_5d4594_3835388", "dword_5d4594_3835392"} {
		p := all[name]
		out[name] = p
		old[p] = *p
	}
	return out, func() {
		for p, v := range old {
			*p = v
		}
	}
}
