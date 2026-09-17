//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func TestCreatureXferAutomaticSpells(t *testing.T) {
	s := newCreatureXferOwner(t)
	cache := unsafe.Slice(memmap.PtrUint32(0x5d4594, 2491624), 9)
	old := append([]uint32(nil), cache...)
	defer copy(cache, old)
	var rows []struct {
		Case string
		Data []byte
	}
	defer func() {
		spellbookCapture(t, "creature-xfer-auto-spells", rows, "d52054751fa0c381c9ddd5c3a4b8abd66ee08afef8cdbe95539d443ff4adf1f7")
	}()
	for _, defined := range []bool{false, true} {
		for _, matched := range []bool{false, true} {
			for _, flag := range []byte{0, 1, 2, 255} {
				t.Run(fmt.Sprintf("defined%v-matched%v-flag%d", defined, matched, flag), func(t *testing.T) {
					u := newCreatureXferObject(t, s, "Monster")
					for i := range cache {
						cache[i] = uint32(0x7000 + i)
					}
					if matched {
						cache[1] = uint32(u.TypeInd)
					}
					def, free := alloc.New(server.MonsterDef{})
					defer free()
					def.TypeInd240 = uint32(u.TypeInd)
					var head unsafe.Pointer
					if defined {
						head = unsafe.Pointer(def)
					}
					restore := legacy.PortTestCreatureXferDefinitions(head)
					defer restore()
					*(*byte)(unsafe.Add(u.UpdateData, 1445)) = 1
					*(*byte)(unsafe.Add(u.UpdateData, 2036)) = flag
					data := unsafe.Slice((*byte)(u.UpdateData), int(unsafe.Sizeof(server.MonsterUpdateData{})))
					want := append([]byte(nil), data...)
					if defined && matched && flag == 1 {
						for off, value := range map[int]uint32{2040: 3, 1540: 0x8000000, 1640: 0x8000000, 1644: 0x10000000, 1776: 0x20000000, 1596: 0x40000000, 1688: 0x40000000, 1660: 0x40000000, 1584: 0x40000000, 1504: 0x80000000} {
							binary.LittleEndian.PutUint32(want[off:], value)
						}
					}
					hp := *u.HealthData
					legacy.PortTestCreatureXferHelper(4, u, nil, 0)
					if !bytes.Equal(data, want) || *u.HealthData != hp {
						t.Fatal("automatic spell gate/defaults or unrelated health changed")
					}
					rows = append(rows, struct {
						Case string
						Data []byte
					}{t.Name(), append([]byte(nil), data...)})
				})
			}
		}
	}
}
