//go:build porttest

package opennox

import (
	"bytes"
	"fmt"
	"github.com/opennox/libs/types"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/ntype"
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"testing"
	"unsafe"
)

func TestSessionEntryPlayerIncoming(t *testing.T) {
	type row struct {
		Flags               uint32
		Index               int
		Mask, Status, Stage uint32
		Position            [2]uint32
		AliasesCleared      bool
		Queue               legacy.PortTestReliableReportState
	}
	var rows []row
	for _, flags := range []uint32{0, 32, 64, 512, 1024, 4096, 8192} {
		for _, index := range []int{1, 7, 31} {
			t.Run(fmt.Sprintf("flags%x/index%d", flags, index), func(t *testing.T) {
				o := newMatchRosterOwner(t)
				defer noxflags.PortTestGameFlags(noxflags.GameFlag(flags))()
				t.Cleanup(func() {
					for _, slot := range []ntype.PlayerInd{1, 3, 7, 31} {
						for i := range o.units {
							o.s.Players.Nox_xxx_netUnmarkMinimapObj_417300(slot, &o.units[i], 0xffffffff)
						}
					}
				})
				pl := o.s.Players.ByInd(ntype.PlayerInd(index))
				u := pl.PlayerUnit
				callback, events := legacy.PortTestSessionEntryInitCallback()
				u.Init = callback
				u.PosVec = types.Pointf{X: 12.5, Y: -3.25}
				pl.Field3680 = 0
				objectXferSetWord(pl.C(), 3632, 0x11223344)
				objectXferSetWord(pl.C(), 3636, 0x55667788)
				objectXferSetWord(pl.C(), 4700, 0x1234)
				aliases := unsafe.Slice((*byte)(unsafe.Add(pl.C(), 16)), 2040)
				for i := range aliases {
					aliases[i] = 0xa5
				}
				o.reset()
				*o.words["mask"] = 0
				legacy.Nox_xxx_netPlayerIncomingServ_4DDF60(index)
				ev := events()
				if len(ev) != 1 || ev[0][0] != uintptr(u.CObj()) || ev[0][1] != 0 {
					t.Fatal("arrival initializer callback")
				}
				r := row{Flags: flags, Index: index, Mask: *o.words["mask"], Status: pl.Field3680, Stage: uint32(*(*byte)(unsafe.Add(pl.C(), 3676))), Position: [2]uint32{objectXferGetWord(pl.C(), 3632), objectXferGetWord(pl.C(), 3636)}, AliasesCleared: bytes.Equal(aliases, make([]byte, len(aliases))), Queue: o.state()}
				wantPos := [2]uint32{math.Float32bits(12.5), math.Float32bits(-3.25)}
				if flags&512 != 0 {
					wantPos = [2]uint32{0x11223344, 0x55667788}
				}
				if r.Mask != 1<<uint(index) || r.Stage != 3 || r.Position != wantPos || !r.AliasesCleared || objectXferGetWord(pl.C(), 4700) != 0 {
					t.Fatalf("arrival state %+v", r)
				}
				if flags&4096 != 0 && index != 31 && objectXferGetWord(u.UpdateData, 552) != 1 {
					t.Fatal("quest arrival marker")
				}
				rows = append(rows, r)
			})
		}
	}
	spellbookCapture(t, "session-entry-player-incoming", rows, "f694592f46adbfdb6df00cd4deab78083203cf984d90f76c6812a4ab9b1c22fa")
}
