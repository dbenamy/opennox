package legacy

import (
	"os"
	"strings"
	"unsafe"

	"github.com/opennox/libs/datapath"
	"github.com/opennox/libs/ifs"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
)

func mapOrchestrationSetName(name string) uint32 {
	if i := strings.IndexByte(name, 0); i >= 0 {
		name = name[:i]
	}
	dst := unsafe.Slice(memmap.PtrUint8(0x587000, 197860), len(name)+1)
	copy(dst, name)
	dst[len(name)] = 0
	return uint32(len(dst))
}
func mapOrchestrationName() unsafe.Pointer { return memmap.PtrOff(0x587000, 197860) }
func mapOrchestrationStart() uint32 {
	Nox_xxx_mapSwitchLevel_4D12E0(true)
	noxflags.SetGame(0x400000)
	saved := memmap.PtrUint32(0x5D4594, 1550924)
	*saved = uint32(uintptr(unsafe.Pointer(GetServer().S().Objs.GetAndZeroObjects())))
	clear(unsafe.Slice(memmap.PtrUint8(0x973F18, 2408), 1464))
	for tries := 100; ; {
		result := mapOrchestrationStep()
		if result == 1 {
			break
		}
		if result == 0 {
			return 0
		}
		tries--
		if tries == 0 {
			return 0
		}
	}
	current := datapath.Data() + "\\nc.obj"
	blend := datapath.Data() + "\\blend.obj"
	if _, err := ifs.Stat(current); !os.IsNotExist(err) {
		if ifs.Remove(current) != nil {
			return 0
		}
	}
	if _, err := ifs.Stat(blend); !os.IsNotExist(err) {
		if ifs.Rename(blend, current) != nil {
			return 0
		}
	}
	for obj := GetServer().S().Objs.First(); obj != nil; obj = obj.ObjNext {
		obj.ScriptIDVal = 0
	}
	Nox_xxx_mapGenMakeInfo_4D5DB0(memmap.PtrOff(0x973F18, 2408))
	name := GoStringP(mapOrchestrationName())
	directory := datapath.Data() + "\\Maps\\$" + name
	_ = ifs.Mkdir(directory)
	if !Nox_xxx_mapSaveMap_51E010(directory+"\\$"+name+".map", 1) {
		return 0
	}
	Nox_xxx_mapSwitchLevel_4D12E0(true)
	GetServer().S().Objs.SetObjects((*server.Object)(unsafe.Pointer(uintptr(*saved))))
	*saved = 0
	noxflags.UnsetGame(0x400000)
	return 1
}
func mapOrchestrationStartAlt() uint32 {
	noxflags.SetGame(0x400000)
	clear(unsafe.Slice(memmap.PtrUint8(0x973F18, 2408), 1464))
	for tries := 100; ; {
		result := mapOrchestrationStep()
		if result == 1 {
			break
		}
		if result == 0 {
			return 0
		}
		tries--
		if tries == 0 {
			return 0
		}
	}
	current := datapath.Data() + "\\nc.obj"
	blend := datapath.Data() + "\\blend.obj"
	if ifs.Remove(current) != nil {
		return 0
	}
	if ifs.Rename(blend, current) != nil {
		return 0
	}
	Nox_xxx_mapGenMakeInfo_4D5DB0(memmap.PtrOff(0x973F18, 2408))
	noxflags.UnsetGame(0x400000)
	return 1
}
