//go:build porttest

package legacy

import (
	"bytes"
	"math"
	"strings"
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
)

// Own the actual C list, mapped namespace/scratch buffers and separator bytes.
func PortTestQuestProgressOwner() func() {
	old := questProgressHead
	questProgressHead = nil
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
		questProgressHead = old
		for _, r := range regions {
			copy(r.dst, r.data)
		}
	}
}

func PortTestQuestProgress(op, name string, value uint32) uint64 {
	var ptr unsafe.Pointer
	switch op {
	case "namespace":
		questProgressNamespace(name)
	case "namespace-nil":
		/* Preserve the namespace on a nil C input. */
	case "set-int":
		ptr = unsafe.Pointer(questProgressSet(name, value, 0))
	case "set-float":
		ptr = unsafe.Pointer(questProgressSet(name, value, 1))
	case "find":
		ptr = unsafe.Pointer(questProgressFind(name))
	case "int":
		return uint64(questProgressInt(name))
	case "float":
		return math.Float64bits(questProgressFloat(name))
	case "reset":
		questProgressReset(name) // Every production caller ignores the incidental return.
	case "qualify":
		qualified := strings.Contains(questProgressCString(name), ":")
		text, ok := questProgressQualify(name)
		if qualified && ok {
			return uint64(len(text) + 1)
		}
		return 0
	case "write":
		return uint64(questProgressWrite())
	case "read":
		return uint64(questProgressRead())
	case "stage-set":
		return uint64(uint32(questProgressSetStage(value)))
	case "stage":
		return uint64(uint32(questProgressStage()))
	case "minions-set":
		return uint64(uint32(questProgressSetMinions(value)))
	default:
		panic(op)
	}
	if ptr == nil {
		return 0
	}
	for i, n := uint64(1), uintptr(unsafe.Pointer(questProgressHead)); n != 0; i++ {
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
	for n := uintptr(unsafe.Pointer(questProgressHead)); n != 0; {
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
		return questProgressGeneratorType(u)
	case "generator-init":
		questProgressInitMapping()
	case "prepare":
		questProgressPrepare(stage)
	case "hecubah":
		questProgressSpawnBoss(u.PosVec, true)
	case "necro":
		questProgressSpawnBoss(u.PosVec, false)
	default:
		panic(op)
	}
	return 0
}
