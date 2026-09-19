//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	flags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestPlayerFilesInventoryRestore(t *testing.T) {
	var rows []map[string]any
	for _, version := range []uint16{1, 2, 3, 0xffff} {
		for _, gf := range []flags.GameFlag{2048, 4096} {
			t.Run(fmt.Sprintf("v%d/flags%d", version, gf), func(t *testing.T) {
				o := newPlayerFileInventoryReadOwner(t)
				o.s.PortTestAIEmptyMap()
				t.Cleanup(o.s.Map.Free)
				oldPickup := o.s.Objs.DefaultPickup
				o.s.Objs.DefaultPickup = nox_xxx_pickupDefault_4F31E0
				t.Cleanup(func() { o.s.Objs.DefaultPickup = oldPickup })
				eligibility := serverConfigOwnBytes(t, 0x5D4594, 1568328, 28)
				clear(eligibility)
				binary.LittleEndian.PutUint32(eligibility, uint32(o.s.Types.IndByID("PortGold")))
				o.reset()
				flags.ResetGame()
				flags.SetGame(gf)
				u := &o.units[0]
				ud := u.UpdateDataPlayer()
				p := ud.Player
				*u.HealthData = server.HealthData{Cur: 17, Max: 25}
				ud.ManaCur = 11
				ud.ManaMax = 15
				u.Experience = 100
				u.InvFirstItem = nil
				p.GoldVal = 321
				*(*byte)(unsafe.Add(u.UpdateData, 244)) = 0xa7
				var input mapDrawableStream
				input.u16(version)
				input.u8(1)
				input.u32(500)
				input.u32(2)
				for i := 0; i < 2; i++ {
					name := o.s.Types.ByID("PortGold").ID()
					input.u8(byte(len(name)))
					input.WriteString(name)
					input.u16(60)
					input.u16(64)
					input.u32(uint32(500 + i))
					input.u32(uint32(900 + i))
					input.f32(64.25)
					input.f32(128.75)
					input.u8(0)
					input.u32(uint32(70 + i))
				}
				input.u8(0)
				input.u32(900)
				if int16(version) >= 2 {
					input.u32(901)
				}
				if int16(version) >= 3 {
					input.u8(0x5c)
				}
				size := input.Len()
				input.Write([]byte{0xde, 0xad, 0xbe, 0xef})
				ret, got, pos := playerFileSection(t, "sub_41AC30", input.Bytes(), uint32(uintptr(unsafe.Pointer(u))), 0)
				if ret != 1 || pos != int64(size) || !bytes.Equal(got, input.Bytes()) || p.GoldVal != 500 || u.CarryCapacity == 0 {
					t.Fatal("populated inventory result", version, gf, ret, pos, size, p.GoldVal, u.CarryCapacity)
				}
				var items []*server.Object
				for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
					items = append(items, it)
					if len(items) > 2 {
						t.Fatal("inventory cycle/count")
					}
				}

				for _, it := range items {
					trackObjectXferTyped(t, o.s, it)
				}
				t.Cleanup(func() {
					u.InvFirstItem = nil
					for _, it := range items {
						it.InvNextItem = nil
						it.InvHolder = nil
						o.s.ObjSetOwner(nil, it)
					}
				})
				if len(items) != 2 {
					t.Fatal("inventory count", len(items))
				}
				var values []map[string]any
				for i, it := range items {
					// Inventory insertion is at the head, reversing the file order.
					source := 1 - i
					if it.InvHolder != u || it.Owner() != u || *(*uint32)(it.InitData) != uint32(70+source) || it.PosVec.X != 2944 || it.PosVec.Y != 2944 {
						t.Fatal("restored item", i, it.InvHolder == u, it.Owner() == u, it.PosVec, *(*uint32)(it.InitData))
					}
					if gf == 2048 && (it.ScriptIDVal != 900+source || it.Extent != uint32(500+source)) {
						t.Fatal("saved item IDs", i, it.ScriptIDVal, it.Extent)
					}
					if gf == 4096 && (it.Extent != it.NetCode || i == 1 && it.ScriptIDVal != items[0].ScriptIDVal+1) {
						t.Fatal("quest renumbering", i, it.ScriptIDVal, it.Extent, it.NetCode)
					}
					values = append(values, map[string]any{"type": o.s.Types.ByInd(int(it.TypeInd)).ID(), "script": it.ScriptIDVal, "extent": it.Extent, "netcode": it.NetCode, "value": *(*uint32)(it.InitData)})
				}
				wantField := byte(0)
				if int16(version) >= 3 && gf == 2048 {
					wantField = 0x5c
				}
				if *(*byte)(unsafe.Add(u.UpdateData, 244)) != wantField {
					t.Fatal("inventory tail")
				}
				state := o.state()
				loaded := 0
				for _, node := range state.Nodes {
					if bytes.Equal(node.Data, []byte{113}) {
						loaded++
					}
				}
				if loaded != 1 {
					t.Fatal("inventory completion report", state)
				}
				rows = append(rows, map[string]any{"version": version, "flags": uint32(gf), "return": ret, "position": pos, "items": values, "queue": state})
			})
		}
	}
	spellbookCapture(t, "player-files-inventory-restore", rows, "c20170bae6d622d70c4109d2120965045f944b376d5a31fa70ac3248d6dc42d1")
}
