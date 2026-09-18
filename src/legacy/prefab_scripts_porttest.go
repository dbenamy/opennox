//go:build porttest

package legacy

import (
	"math"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func PortTestPrefabScriptsGlobals() ([]*uint32, func()) {
	old := prefabScriptOffsets
	return []*uint32{&prefabScriptOffsets[0], &prefabScriptOffsets[1], &prefabScriptOffsets[2]}, func() { prefabScriptOffsets = old }
}
func PortTestPrefabScriptsCall(op int, a, b, c unsafe.Pointer, arg uint32) uint64 {
	switch op {
	case 0:
		return uint64(prefabScriptReadInt(fileByHandle((*FILE)(a))))
	case 1:
		return math.Float64bits(prefabScriptReadFloat(fileByHandle((*FILE)(a))))
	case 2:
		return uint64(prefabScriptWriteFloat(fileByHandle((*FILE)(a)), math.Float32frombits(arg)))
	case 3:
		return uint64(prefabScriptWriteInt(fileByHandle((*FILE)(a)), arg))
	case 4:
		return uint64(prefabScriptInstructions(fileByHandle((*FILE)(a)), fileByHandle((*FILE)(b)), arg != 0))
	case 5:
		return uint64(prefabScriptStrings(fileByHandle((*FILE)(a)), fileByHandle((*FILE)(b)), fileByHandle((*FILE)(c))))
	case 6:
		return uint64(prefabScriptReserved(fileByHandle((*FILE)(a)), fileByHandle((*FILE)(b)), fileByHandle((*FILE)(c))))
	case 7:
		return uint64(prefabScriptLocals(fileByHandle((*FILE)(a)), fileByHandle((*FILE)(b))))
	case 8:
		return uint64(prefabScriptFunctions(fileByHandle((*FILE)(a)), fileByHandle((*FILE)(b)), fileByHandle((*FILE)(c))))
	case 9:
		return uint64(prefabScriptMerge(fileByHandle((*FILE)(a)), fileByHandle((*FILE)(b)), fileByHandle((*FILE)(c))))
	case 10:
		return uint64(prefabScriptRewrite(alloc.GoString((*byte)(a)), *(*[2]int32)(b)))
	case 11:
		var xy [2]int32
		if a != nil {
			xy = *(*[2]int32)(a)
		}
		prefabScriptObjectNames(int32(arg), xy[0], xy[1])
		return 0
	case 14:
		prefabScriptSetCallback((*server.Object)(a), int(arg), alloc.GoString((*byte)(b)))
		return 0
	case 15:
		return uint64(prefabScriptPending((*[4]int32)(a), arg))
	case 16:
		return uint64(prefabScriptGeneration(alloc.GoString((*byte)(a))))
	case 19:
		prefabScriptBounds((*[8]uint32)(a))
		return uint64(uintptr(a))
	default:
		panic("prefab script operation")
	}
}
func PortTestPrefabScriptsName(name string, instance, x, y int32, objectOnly bool) string {
	return prefabScriptName(name, instance, x, y, objectOnly)
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
