//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/object"
	flags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestPlayerFilesQuestInventoryWrite(t *testing.T) {
	o := newMatchRosterOwner(t)
	t.Cleanup(o.s.PortTestInventoryEnvironment(false, nil, nil))
	serverConfigOwnBytes(t, 0x5D4594, 527704, 4)
	serverConfigOwnBytes(t, 0x5D4594, 527724, 4)
	limits := serverConfigOwnBytes(t, 0x5D4594, 1568356, 52)
	glyphID := o.s.Types.IndByID("Glyph")
	if glyphID == 0 {
		t.Fatal("missing glyph")
	}
	for i := 0; i < 12; i++ {
		binary.LittleEndian.PutUint32(limits[4*i:], uint32(glyphID))
	}
	binary.LittleEndian.PutUint32(limits[48:], uint32(o.s.Types.IndByID("PortGold")))
	configure, restore := o.s.PortTestWorldCollisionBalance()
	t.Cleanup(restore)
	_, _, restore = legacy.PortTestMeterInventory()
	t.Cleanup(restore)
	words, restore := legacy.PortTestUIInventoryWords()
	t.Cleanup(restore)
	u := &o.units[0]
	p := u.UpdateDataPlayer().Player
	a := newItemXferObject(t, o.s, "Gold")
	b := newItemXferObject(t, o.s, "Gold")
	glyph, free := alloc.New(server.Object{})
	t.Cleanup(free)
	glyph.TypeInd = uint16(glyphID)
	a.ObjFlags = 0
	a.Extent = 500
	a.ScriptIDVal = 900
	a.NetCode = 100
	b.ObjFlags = 0
	b.ObjClass = object.Class(0x40)
	b.ScriptIDVal = 901
	b.NetCode = 101
	a.InvNextItem = b
	b.InvNextItem = glyph
	u.InvFirstItem = a
	t.Cleanup(func() { u.InvFirstItem = nil; a.InvNextItem = nil; b.InvNextItem = nil })
	var rows []map[string]any
	for _, gf := range []flags.GameFlag{4096, 4096 | 2048, 4096 | 8192} {
		for _, equipped := range []bool{false, true} {
			for _, limit := range []float64{1, 3} {
				o.reset()
				flags.ResetGame()
				flags.SetGame(gf)
				configure(map[string]float64{"ForceOfNatureStaffLimit": limit})
				p.GoldVal = 500
				*(*byte)(unsafe.Add(u.UpdateData, 244)) = 0xa7
				b.ObjFlags = 0
				if equipped {
					b.ObjFlags = object.Flags(0x100)
				}
				for _, w := range words {
					*w = 0
				}
				*words[3] = b.NetCode
				var want mapDrawableStream
				want.u16(3)
				want.u8(1)
				want.u32(500)
				want.u32(1)
				name := o.s.Types.ByInd(int(a.TypeInd)).ID()
				want.u8(byte(len(name)))
				want.WriteString(name)
				want.u16(60)
				want.u16(64)
				want.u32(500)
				want.u32(900)
				want.f32(0)
				want.f32(0)
				want.u8(0)
				want.u32(0)
				if equipped {
					want.u8(1)
					want.u32(901)
				} else {
					want.u8(0)
				}
				want.u32(0)
				want.u32(901)
				want.u8(0xa7)
				ret, data, pos := playerFileSection(t, "sub_41AC30", nil, uint32(uintptr(unsafe.Pointer(u))), 0)
				wantRet := uint32(1)
				if limit == 1 {
					wantRet = 0
				}
				state := o.state()
				if ret != wantRet || pos != int64(want.Len()) || !bytes.Equal(data, want.Bytes()) || len(state.Nodes) != int(wantRet) {
					t.Fatal("quest inventory write", gf, equipped, limit, ret, wantRet, pos, want.Len(), data, want.Bytes(), state)
				}
				if wantRet == 1 && !bytes.Equal(state.Nodes[0].Data, []byte{113}) {
					t.Fatal("quest inventory report", state)
				}
				rows = append(rows, map[string]any{"flags": uint32(gf), "equipped": equipped, "limit": limit, "return": ret, "position": pos, "bytes": data, "queue": state})
			}
		}
	}
	spellbookCapture(t, "player-files-quest-inventory-write", rows, "a6aef744f32b4d327a5d615760706de043c6dd14d7c00d57395ae9be6547695d")
}
