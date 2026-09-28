//go:build porttest

package legacy

import (
	"github.com/opennox/opennox/v1/legacy/common/fpenv"
	"math"
	"runtime"
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
)

type PortTestWorldNumericResult struct {
	Input            uint32
	Redirect         bool
	Abs              uint64
	Scratch, Pointed uint32
	Rounded          uint32
	RoundWords       [3]uint32
	PointerOK        bool
}

// PortTestWorldNumeric observes the original helpers' shared scratch state.
// Only the known scratch pointer is normalized; all numeric bits are retained.
func PortTestWorldNumeric(inputs []uint32) (out []PortTestWorldNumericResult) {
	restoreEnv := portTestWorldNumericEnv()
	defer restoreEnv()
	words := unsafe.Slice(memmap.PtrUint32(0x5D4594, 527668), 4)
	old := append([]uint32(nil), words...)
	// The redirect target is fixture-owned raw blob storage, not a live global.
	tail := (*uint32)(unsafe.Pointer(&memmap.BlobByAddr(0x5D4594).Data[527684]))
	oldTail := *tail
	slot := memmap.PtrPtr(0x587000, 55744)
	oldSlot := *slot
	defer func() { copy(words, old); *tail = oldTail; *slot = oldSlot }()
	for _, bits := range inputs {
		for _, redirect := range []bool{false, true} {
			words[0], words[1], words[2], words[3] = 0x13572468, 0x24681357, 0x87654321, 0xfedcba98
			*tail = 0xff812345
			*slot = unsafe.Pointer(&words[1])
			if redirect {
				*slot = unsafe.Pointer(tail)
			}
			r := PortTestWorldNumericResult{Input: bits, Redirect: redirect}
			r.Abs = math.Float64bits(worldAbsScratch(math.Float32frombits(bits)))
			r.Scratch = words[1]
			r.Pointed = *(*uint32)(*slot)
			r.Rounded = worldRoundScratch(math.Float32frombits(bits))
			negative := math.Float32frombits(bits) < 0
			r.PointerOK = words[0] == uint32(uintptr(unsafe.Pointer(&words[2])))
			if negative {
				r.PointerOK = words[0] == 0x13572468
			}
			marker := uint32(1)
			if negative {
				marker = words[0]
			}
			r.RoundWords = [3]uint32{marker, words[2], words[3]}
			out = append(out, r)
		}
	}
	return out
}

func portTestWorldNumericEnv() func() {
	runtime.LockOSThread()
	saved := fpenv.Begin()
	return func() { fpenv.Restore(saved); runtime.UnlockOSThread() }
}
