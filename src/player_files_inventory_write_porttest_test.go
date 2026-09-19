//go:build porttest

package opennox

import (
	"bytes"
	"github.com/opennox/libs/types"
	flags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestPlayerFilesInventoryWrite(t *testing.T) {
	o := newMatchRosterOwner(t)
	t.Cleanup(o.s.PortTestInventoryEnvironment(false, nil, nil))
	serverConfigOwnBytes(t, 0x5D4594, 527704, 4)
	grid, _, restore := legacy.PortTestMeterInventory()
	t.Cleanup(restore)
	words, restore := legacy.PortTestUIInventoryWords()
	t.Cleanup(restore)
	cells := legacy.PortTestInventoryCells()
	u := &o.units[0]
	p := u.UpdateDataPlayer().Player
	items := []*server.Object{newItemXferObject(t, o.s, "Gold"), newItemXferObject(t, o.s, "Ammo")}
	names := []string{"Gold", "Ammo"}
	for i, item := range items {
		item.ObjFlags = 0
		item.Extent = uint32(500 + i)
		item.ScriptIDVal = 900 + i
		item.NetCode = uint32(100 + i)
		item.PosVec = types.Pointf{X: 64.25, Y: 128.75}
		item.InvHolder = u
	}
	t.Cleanup(func() {
		u.InvFirstItem = nil
		for _, item := range items {
			item.InvNextItem = nil
			item.InvHolder = nil
		}
	})
	var rows []map[string]any
	for _, gf := range []flags.GameFlag{0, 2048, 8192} {
		for _, which := range []string{"empty", "linked", "grid", "missing-grid", "unmatched-selection"} {
			o.reset()
			flags.ResetGame()
			flags.SetGame(gf)
			clear(grid)
			for _, w := range words {
				*w = 0
			}
			p.GoldVal = 0x89abcdef
			*(*byte)(unsafe.Add(u.UpdateData, 244)) = 0xa7
			u.InvFirstItem = nil
			items[0].InvNextItem = nil
			items[1].InvNextItem = nil
			order := []int{}
			if which != "empty" {
				u.InvFirstItem = items[0]
				items[0].InvNextItem = items[1]
				order = []int{0, 1}
			}
			if which == "grid" || which == "missing-grid" {
				cells[0].Count = 2
				cells[0].Codes[0] = items[1].NetCode
				cells[0].Codes[1] = items[0].NetCode
				if which == "missing-grid" {
					cells[0].Codes[1] = 7777
				}
				if gf == 2048 {
					order = []int{1, 0}
				}
			}
			cells[83].Codes[0] = items[1].NetCode
			*words[2] = uint32(uintptr(unsafe.Pointer(&cells[83])))
			*words[3] = items[0].NetCode
			if which == "unmatched-selection" {
				cells[83].Codes[0] = 7777
				*words[3] = 8888
			}
			var want mapDrawableStream
			want.u16(3)
			present := byte(1)
			if gf == 8192 {
				present = 0
			}
			want.u8(present)
			wantRet := uint32(1)
			if gf == 0 {
				wantRet = 0
			} else if gf == 2048 {
				want.u32(p.GoldVal)
				want.u32(uint32(len(order)))
				for at, idx := range order {
					if which == "missing-grid" && at == 1 {
						wantRet = 0
						break
					}
					item := items[idx]
					name := o.s.Types.ByInd(int(item.TypeInd)).ID()
					want.u8(byte(len(name)))
					want.WriteString(name)
					want.u16(60)
					want.u16(64)
					want.u32(item.Extent)
					want.u32(uint32(item.ScriptIDVal))
					want.f32(64.25)
					want.f32(128.75)
					want.u8(0)
					want.Write(itemXferDefaultPayload(names[idx]))
				}
				if wantRet == 1 {
					want.u8(0)
					secondary, quiver := uint32(0), uint32(0)
					if which != "empty" && which != "unmatched-selection" {
						secondary = 901
						quiver = 900
					}
					want.u32(secondary)
					want.u32(quiver)
					want.u8(0xa7)
				}
			}
			beforeGrid := bytes.Clone(grid)
			ret, got, pos := playerFileSection(t, "sub_41AC30", nil, uint32(uintptr(unsafe.Pointer(u))), 0)
			if ret != wantRet || pos != int64(want.Len()) || !bytes.Equal(got, want.Bytes()) || !bytes.Equal(grid, beforeGrid) || p.GoldVal != 0x89abcdef || *(*byte)(unsafe.Add(u.UpdateData, 244)) != 0xa7 {
				t.Fatalf("inventory writer %d/%s ret=%d/%d pos=%d/%d bytes=%x want=%x", gf, which, ret, wantRet, pos, want.Len(), got, want.Bytes())
			}
			state := o.state()
			if len(state.Nodes) != int(wantRet) || wantRet == 1 && !bytes.Equal(state.Nodes[0].Data, []byte{113}) {
				t.Fatal("inventory writer loaded report", gf, which, state)
			}
			rows = append(rows, map[string]any{"flags": uint32(gf), "case": which, "return": ret, "bytes": got, "position": pos, "queue": o.state()})
		}
	}
	spellbookCapture(t, "player-files-inventory-write", rows, "")
}
