//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestQuickbarActivationMessages(t *testing.T) {
	q := newQuickbarOwner(t)
	cursor := memmap.PtrUint32(0x5D4594, 1096672)
	oldCursor := *cursor
	t.Cleanup(func() { *cursor = oldCursor })
	words, free := alloc.Make([]uint32{}, 5)
	defer free()
	var rows []quickbarResult
	for _, busy := range []uint32{0, 1, 5} {
		for _, id := range []uint32{0, 1, 5, 255, 256, 0xffffffff} {
			q.reset(t)
			*cursor = busy
			ret := q.call("nox_xxx_clientSendAbil_45DAF0", id)
			r := q.snapshot(fmt.Sprintf("ability-cursor%d-id%d", busy, id), ret)
			if busy == 0 && id != 0 {
				q.check(t, reflect.DeepEqual(r.Messages, [][]byte{{122, byte(id)}}), "ability identifier uses one byte")
			} else {
				q.check(t, len(r.Messages) == 0, "guarded ability emits no message")
			}
			rows = append(rows, r)
		}
		for _, count := range []uint32{0xffffffff, 0, 1, 2, 4, 5, 6} {
			for _, flag := range []uint32{0, 1, 0x80, 0xff} {
				for _, first := range []uint32{0, 1, 0xffffffff} {
					q.reset(t)
					*cursor = busy
					copy(words, []uint32{first, 0x12345678, 0x80000000, 136, 34})
					original := append([]uint32(nil), words...)
					ret := q.call("nox_xxx_clientSendSpell_45DB20", uint32(uintptr(unsafe.Pointer(&words[0]))), count, flag)
					r := q.snapshot(fmt.Sprintf("spell-cursor%d-count%d-flag%d-first%d", busy, count, flag, first), ret)
					if busy == 0 && first != 0 {
						want := make([]byte, 22)
						want[0] = 121
						want[21] = byte(flag)
						for i := 0; i < 5; i++ {
							if int32(i) < int32(count) {
								binary.LittleEndian.PutUint32(want[1+4*i:], words[i])
							}
						}
						q.check(t, reflect.DeepEqual(r.Messages, [][]byte{want}), "spell message count, padding and byte order")
					} else {
						q.check(t, len(r.Messages) == 0, "guarded or empty spell emits no message")
					}
					q.check(t, reflect.DeepEqual(original, words), "message assembly leaves input unchanged")
					rows = append(rows, r)
				}
			}
		}
	}
	spellbookCapture(t, "quickbar-activation-messages", rows, "f1f1b054aeff77a483bed7e9f007c55260c324538de777738205df7c31429d22")
}

func TestQuickbarPendingActivation(t *testing.T) {
	q := newQuickbarOwner(t)
	var rows []quickbarResult
	for _, op := range []string{"nox_xxx_guiSpellTargetClickSet_45D9D0", "nox_xxx_guiSpell_45DA10"} {
		for _, pending := range []uint32{0, 2, 0xffffffff} {
			for _, id := range []uint32{0, 1, 136, 0xffffffff} {
				q.reset(t)
				*memmap.PtrUint32(0x5D4594, 1047916) = pending
				*memmap.PtrUint32(0x5D4594, 1047920) = 0xaabbccdd
				*memmap.PtrUint32(0x5D4594, 1047924) = 9
				*memmap.PtrUint32(0x5D4594, 1047928) = 7
				before := append([]uint32(nil), q.raw...)
				ret := q.call(op, id)
				if pending != 0 {
					q.check(t, ret == 0 && reflect.DeepEqual(before, q.raw) && len(inputKeyTimeoutsOld) == 0, "occupied activation is unchanged")
				} else {
					q.check(t, ret == 1 && memmap.Uint32(0x5D4594, 1047916) == id, "pending activation stores id")
					q.check(t, memmap.Uint32(0x5D4594, 1047920) == 0xaabbcc00, "activation resets only low byte")
					_, set := inputKeyTimeoutsOld[5]
					q.check(t, set, "activation sets actual input timeout")
					if op == "nox_xxx_guiSpellTargetClickSet_45D9D0" {
						q.check(t, memmap.Uint32(0x5D4594, 1047924) == 0 && memmap.Uint32(0x5D4594, 1047928) == 0, "aimed activation clears modes")
					} else {
						q.check(t, memmap.Uint32(0x5D4594, 1047924) == 1 && memmap.Uint32(0x5D4594, 1047928) == 7, "instant activation preserves cursor mode")
					}
				}
				rows = append(rows, q.snapshot(fmt.Sprintf("%s-pending%d-id%d", op, pending, id), ret))
				q.c.srv.NetList.ResetByInd(31, netlist.Kind0)
				q.call("nox_xxx_guiSpellTargetClickCheckSend_45DBB0")
				q.check(t, memmap.Uint32(0x5D4594, 1047916) == 0, "sending consumes pending activation")
				rows = append(rows, q.snapshot(fmt.Sprintf("%s-pending%d-id%d-send", op, pending, id), 0))
			}
		}
	}
	spellbookCapture(t, "quickbar-pending-activation", rows, "ac2b95389c86c2584b572dc8b9faa855125c9dbc8bf4e8f997e3f0222079c738")
}
