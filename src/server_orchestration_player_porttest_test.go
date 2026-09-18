//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestServerOrchestrationPlayerReset(t *testing.T) {
	o := newMatchRosterOwner(t)
	t.Cleanup(o.s.PortTestRewardTypes([]string{"Glyph", "ResetOther"}, nil, true, 0, 0))
	serverConfigOwnBytes(t, 0x5D4594, 1568268, 4)
	*memmap.PtrUint32(0x5D4594, 1568268) = 0
	var inventory []*server.Object
	for _, name := range []string{"Glyph", "ResetOther", "Glyph"} {
		u := o.s.NewObjectByTypeID(name)
		if u == nil {
			t.Fatal("inventory allocation")
		}
		inventory = append(inventory, u)
		t.Cleanup(func() { o.s.Objs.FreeObject(u) })
	}
	type row struct {
		Flags                              uint32
		Mask                               int
		Glyphs                             byte
		Words                              [7]uint32
		PlayerWords                        [7]uint32
		ReturnMatches, OtherBytesPreserved bool
	}
	var rows []row
	for _, flags := range []uint32{0, 2048, 4096, 6144, 0xffffffff} {
		for mask := 0; mask < 8; mask++ {
			t.Run(fmt.Sprintf("flags%x/mask%d", flags, mask), func(t *testing.T) {
				defer noxflags.PortTestGameFlags(noxflags.GameFlag(flags))()
				u := &o.units[0]
				pl := u.UpdateDataPlayer().Player
				u.InvFirstItem = nil
				glyphs := byte(0)
				for i := len(inventory) - 1; i >= 0; i-- {
					it := inventory[i]
					it.InvNextItem = nil
					if mask&(1<<i) != 0 {
						it.InvNextItem = u.InvFirstItem
						u.InvFirstItem = it
						if i != 1 {
							glyphs++
						}
					}
				}
				defer func() { u.InvFirstItem = nil }()
				ud := unsafe.Slice((*byte)(u.UpdateData), int(unsafe.Sizeof(server.PlayerUpdateData{})))
				pr := unsafe.Slice((*byte)(pl.C()), int(unsafe.Sizeof(server.Player{})))
				offs := []int{116, 120, 124, 128, 308, 264}
				for i, off := range offs {
					binary.LittleEndian.PutUint32(ud[off:], uint32(0x80112233)+uint32(i))
				}
				ud[244] = 73
				*(*uint32)(unsafe.Add(u.CObj(), 520)) = 0x76543210
				for i, off := range []int{4796, 4800, 4804, 4808, 4812, 3660, 3664} {
					binary.LittleEndian.PutUint32(pr[off:], uint32(0x12345600+i))
				}
				wantUD, wantPR := bytes.Clone(ud), bytes.Clone(pr)
				for _, off := range offs {
					binary.LittleEndian.PutUint32(wantUD[off:], 0)
				}
				if flags&2048 == 0 {
					wantUD[244] = glyphs
				}
				for _, off := range []int{4796, 4800, 4804, 4808, 4812} {
					binary.LittleEndian.PutUint32(wantPR[off:], 0)
				}
				for _, off := range []int{3660, 3664} {
					binary.LittleEndian.PutUint32(wantPR[off:], 0xdeadface)
				}
				ret := legacy.PortTestServerOrchestration("reset-player", u, 0)
				r := row{Flags: flags, Mask: mask, Glyphs: ud[244], ReturnMatches: ret == uint32(uintptr(pl.C())), OtherBytesPreserved: bytes.Equal(ud, wantUD) && bytes.Equal(pr, wantPR)}
				for i, off := range offs {
					r.Words[i] = binary.LittleEndian.Uint32(ud[off:])
				}
				r.Words[6] = *(*uint32)(unsafe.Add(u.CObj(), 520))
				for i, off := range []int{4796, 4800, 4804, 4808, 4812, 3660, 3664} {
					r.PlayerWords[i] = binary.LittleEndian.Uint32(pr[off:])
				}
				if !r.ReturnMatches || !r.OtherBytesPreserved || *(*uint32)(unsafe.Add(u.CObj(), 520)) != 0 {
					t.Fatal("reset state", r)
				}
				rows = append(rows, r)
			})
		}
	}
	spellbookCapture(t, "server-orchestration-player-reset", rows, "9e62b5a235e909219f97b671b718a09595b8765cdc82c56704e7db01d4473225")
}
