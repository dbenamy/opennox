//go:build porttest

package legacy

/*
#include <stdint.h>
#include "GAME1_1.h"
int nox_xxx_cliPlrInfoLoadFromFile_41A2E0(char* path, int pind);
int nox_xxx_plrLoad_41A480(char* a1);
int sub_41A590(void* a1p, void* a2p);
int sub_41AA30(void* a1p, void* a2p);
int sub_41AC30(void* a1p, void* a2p);
int sub_41B3B0();
int sub_41B3E0(int a1);
int nox_xxx_guiFieldbook_41B420(void* a1p, void* a2p);
int nox_xxx_guiSpellbook_41B660(void* a1p, void* a2p);
int nox_xxx_guiEnchantment_41B9C0(void* a1p, void* a2p);
int sub_41BEC0(void* a1p, void* a2p);
int sub_41C080(void* a1p, void* a2p);
int sub_41C280(void* a1);
int nox_xxx_parseFileInfoData_41C3B0(int a1);
int sub_41C780(int a1);
void sub_41CAC0(char* a1, void* a2);
int nox_xxx_netSavePlayer_41CE00();
int sub_41CEE0(void* a1p, int a2);
static uint32_t nox_porttest_player_files_call(int op, uint32_t a0,uint32_t a1){switch(op){
case 0: return (uint32_t)nox_xxx_cliPlrInfoLoadFromFile_41A2E0((char*)a0,(int)a1);
case 1: return (uint32_t)nox_xxx_plrLoad_41A480((char*)a0);
case 2: return (uint32_t)sub_41A590((void*)a0,(void*)a1);
case 3: return (uint32_t)sub_41AA30((void*)a0,(void*)a1);
case 4: return (uint32_t)sub_41AC30((void*)a0,(void*)a1);
case 5: return (uint32_t)sub_41B3B0();
case 6: return (uint32_t)sub_41B3E0((int)a0);
case 7: return (uint32_t)nox_xxx_guiFieldbook_41B420((void*)a0,(void*)a1);
case 8: return (uint32_t)nox_xxx_guiSpellbook_41B660((void*)a0,(void*)a1);
case 9: return (uint32_t)nox_xxx_guiEnchantment_41B9C0((void*)a0,(void*)a1);
case 10: return (uint32_t)sub_41BEC0((void*)a0,(void*)a1);
case 11: return (uint32_t)sub_41C080((void*)a0,(void*)a1);
case 12: return (uint32_t)sub_41C280((void*)a0);
case 13: return (uint32_t)nox_xxx_parseFileInfoData_41C3B0((int)a0);
case 14: return (uint32_t)sub_41C780((int)a0);
case 15: sub_41CAC0((char*)a0,(void*)a1); return 0;
case 16: return (uint32_t)nox_xxx_netSavePlayer_41CE00();
case 17: return (uint32_t)sub_41CEE0((void*)a0,(int)a1);
}return 0;}
*/
import "C"
import (
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"unsafe"
)

var portTestPlayerFileNames = []string{
	"nox_xxx_cliPlrInfoLoadFromFile_41A2E0",
	"nox_xxx_plrLoad_41A480",
	"sub_41A590",
	"sub_41AA30",
	"sub_41AC30",
	"sub_41B3B0",
	"sub_41B3E0",
	"nox_xxx_guiFieldbook_41B420",
	"nox_xxx_guiSpellbook_41B660",
	"nox_xxx_guiEnchantment_41B9C0",
	"sub_41BEC0",
	"sub_41C080",
	"sub_41C280",
	"nox_xxx_parseFileInfoData_41C3B0",
	"sub_41C780",
	"sub_41CAC0",
	"nox_xxx_netSavePlayer_41CE00",
	"sub_41CEE0",
}

func PortTestPlayerFileCall(name string, args ...uint32) uint32 {
	if len(args) > 2 {
		panic("too many player file arguments")
	}
	var a [2]uint32
	copy(a[:], args)
	for i, n := range portTestPlayerFileNames {
		if name == n {
			return uint32(C.nox_porttest_player_files_call(C.int(i), C.uint32_t(a[0]), C.uint32_t(a[1])))
		}
	}
	panic("unknown player file operation: " + name)
}

// Own the real client section table and actual C callbacks, including sentinel.
func PortTestPlayerFileClientSections() func() {
	table := unsafe.Slice(memmap.PtrUint32(0x587000, 55936), 12)
	old := append([]uint32(nil), table...)
	clear(table)
	names := []string{"GUI Data", "File Info Data", "Music Data"}
	ids := []uint32{7, 1, 12}
	callbacks := []unsafe.Pointer{unsafe.Pointer(C.sub_41C280), unsafe.Pointer(C.nox_xxx_parseFileInfoData_41C3B0), unsafe.Pointer(C.sub_41C780)}
	var frees []func()
	for i, name := range names {
		p, free := alloc.CString(name)
		frees = append(frees, free)
		table[3*i] = uint32(uintptr(unsafe.Pointer(p)))
		table[3*i+1] = ids[i]
		table[3*i+2] = uint32(uintptr(callbacks[i]))
	}
	return func() {
		copy(table, old)
		for _, free := range frees {
			free()
		}
	}
}
