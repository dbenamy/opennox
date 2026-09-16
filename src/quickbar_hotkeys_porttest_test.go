//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/common/memmap"
	"reflect"
	"testing"
	"unsafe"
)

func TestQuickbarHotkeyActivation(t *testing.T) {
	q := newQuickbarOwner(t)
	var rows []quickbarResult
	for class := uint32(0); class < 3; class++ {
		for slot := 0; slot < 5; slot++ {
			for _, animation := range []uint32{0, 1, 2, 3, 51} {
				for _, cursor := range []uint32{0, 5} {
					q.reset(t)
					*(*byte)(unsafe.Add(unsafe.Pointer(&q.players[0]), 2251)) = byte(class)
					*memmap.PtrPtr(0x852978, 8) = q.items[0].C()
					*(*uint32)(unsafe.Add(q.items[0].C(), 276)) = animation
					*memmap.PtrUint32(0x5D4594, 1096672) = cursor
					if class == 2 {
						*q.words["dword_8531A0_2576"] = 0
					}
					q.bar[2*slot] = uint32(slot + 1)
					q.bar[2*slot+1] = 0xabcdef80 | uint32(slot&1)
					q.call("nox_client_invokeSpellSlot_45DA50", uint32(slot))
					r := q.snapshot(fmt.Sprintf("class%d-slot%d-animation%d-cursor%d", class, slot, animation, cursor), 0)
					eligible := class != 2 && animation != 1 && animation != 2 && animation != 51
					if eligible && cursor == 0 {
						want := []byte{122, byte(slot + 1)}
						if class != 0 {
							want = make([]byte, 22)
							want[0] = 121
							want[21] = byte(slot & 1)
							binary.LittleEndian.PutUint32(want[1:], uint32(slot+1))
						}
						q.check(t, reflect.DeepEqual(r.Messages, [][]byte{want}), "hotkey emits selected slot through actual message owner")
					} else {
						q.check(t, len(r.Messages) == 0, "blocked hotkey emits no message")
					}
					if eligible {
						q.check(t, r.LastButton == uint32(slot), "eligible hotkey records button even if cursor blocks message")
					} else {
						q.check(t, r.LastButton == 0xffffffff, "animation-blocked hotkey leaves last button unchanged")
					}
					rows = append(rows, r)
				}
			}
		}
	}
	spellbookCapture(t, "quickbar-hotkeys", rows, "975f751a3a22f8e17180c7b688976f95392b71ea4f885f1e57d5e0d3ca4ba5ee")
}
