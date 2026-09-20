//go:build porttest

package legacy

/*
#include "GAME1_1.h"
*/
import "C"
import (
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
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
	p0 := unsafe.Pointer(uintptr(a[0]))
	u := (*server.Object)(p0)
	switch name {
	case "nox_xxx_cliPlrInfoLoadFromFile_41A2E0":
		return uint32(playerFileServerLoad(alloc.GoString((*byte)(p0)), int(int32(a[1]))))
	case "nox_xxx_plrLoad_41A480":
		return uint32(nox_xxx_plrLoad_41A480((*C.char)(p0)))
	case "sub_41A590":
		return uint32(playerFileAttributes(u, unsafe.Pointer(uintptr(a[1]))))
	case "sub_41AA30":
		return uint32(playerFileStatus(u))
	case "sub_41AC30":
		return uint32(playerFileInventory(u))
	case "sub_41B3B0":
		return uint32(playerFileInventoryCount())
	case "sub_41B3E0":
		return uint32(bool2int(playerFileInventoryAllowed(u)))
	case "nox_xxx_guiFieldbook_41B420":
		return uint32(playerFileGuides(u))
	case "nox_xxx_guiSpellbook_41B660":
		return uint32(playerFileSpells(u))
	case "nox_xxx_guiEnchantment_41B9C0":
		return uint32(playerFileEnchantment(u))
	case "sub_41BEC0":
		return uint32(playerFileJournal(u))
	case "sub_41C080":
		return uint32(playerFileGame(u))
	case "sub_41C280":
		return uint32(C.sub_41C280(p0))
	case "nox_xxx_parseFileInfoData_41C3B0":
		return uint32(C.nox_xxx_parseFileInfoData_41C3B0(C.int(a[0])))
	case "sub_41C780":
		return uint32(C.sub_41C780(C.int(a[0])))
	case "nox_xxx_netSavePlayer_41CE00":
		return uint32(C.nox_xxx_netSavePlayer_41CE00())
	case "sub_41CEE0":
		return uint32(C.sub_41CEE0(p0, C.int(a[1])))
	case "sub_41CAC0":
		playerFileExtract(alloc.GoString((*byte)(p0)), unsafe.Pointer(uintptr(a[1])))
		return 0
	}
	panic("unknown player file operation: " + name)
}

// Own the real client section table and actual exported callbacks, including sentinel.
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
