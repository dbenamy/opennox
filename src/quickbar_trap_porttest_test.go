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

func TestQuickbarTrapActivation(t *testing.T) {
	q := newQuickbarOwner(t)
	q.ownFullQuickbar(t)
	var rows []quickbarResult
	for row := uint32(0); row < 3; row++ {
		for mask := 0; mask < 8; mask++ {
			for _, flag := range []uint32{0, 1, 255} {
				for _, learned := range []uint32{0, 1} {
					q.reset(t)
					q.call("nox_client_trapSetSelect_4604B0", row)
					*(*uint32)(unsafe.Add(unsafe.Pointer(&q.players[0]), 3832)) = learned
					*memmap.PtrUint32(0x5D4594, 1047916) = 3
					*memmap.PtrUint32(0x5D4594, 1047924) = flag
					*memmap.PtrUint32(0x5D4594, 1049480) = 7
					*memmap.PtrUint32(0x5D4594, 1049488) = 0x123456ff
					var ids []uint32
					for slot := 0; slot < 3; slot++ {
						if mask&(1<<slot) != 0 {
							id := uint32(1 + slot)
							*memmap.PtrUint32(0x5D4594, 1047940+uintptr(row*40)+uintptr(slot*8)) = id
							ids = append(ids, id)
						}
					}
					q.call("nox_client_buildTrap_45E040")
					r := q.snapshot(fmt.Sprintf("row%d-mask%d-flag%d-known%d", row, mask, flag, learned), 0)
					if mask != 0 {
						want := make([]byte, 22)
						want[0] = 121
						want[21] = byte(flag)
						ids = append(ids, 34)
						for i, id := range ids {
							binary.LittleEndian.PutUint32(want[1+4*i:], id)
						}
						q.check(t, reflect.DeepEqual(r.Messages, [][]byte{want}), "trap compacts spells then appends glyph id")
						q.check(t, memmap.Uint32(0x5D4594, 1047916) == 0 && memmap.Uint32(0x5D4594, 1049480) == 0, "trap clears pending activation")
						q.check(t, memmap.Uint32(0x5D4594, 1049488) == 0x12345600, "trap resets only low status byte")
					} else {
						q.check(t, len(r.Messages) == 0 && memmap.Uint32(0x5D4594, 1047916) == 3, "empty trap preserves pending activation")
						if learned != 0 {
							q.check(t, len(r.Book.Sounds) > 0 && r.Book.Sounds[len(r.Book.Sounds)-1][0] == 925, "known empty trap reports feedback")
						}
					}
					rows = append(rows, r)
				}
			}
		}
	}
	spellbookCapture(t, "quickbar-trap-activation", rows, "bb288e1b1c2e15dc461f98d98ae29393b7c91751bffea1cb8b3ef95c65d5dade")
}
