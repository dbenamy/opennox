//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/opennox/libs/object"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/ntype"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestMatchRosterInventory(t *testing.T) {
	o := newMatchRosterOwner(t)
	t.Cleanup(o.s.PortTestRewardTypes([]string{"Flag", "PortTestRewardWeapon"}, nil, true, 0, 0x12345678))
	flag := uint16(o.s.Types.ByID("Flag").Ind())
	weapon := uint16(o.s.Types.ByID("PortTestRewardWeapon").Ind())
	items := make([]server.Object, 3)
	for i := range items {
		it := &items[i]
		it.ObjClass = object.ClassWeapon
		it.InitData = o.record(t, 16)
		it.InvHolder = &o.units[0]
		if i+1 < len(items) {
			it.InvNextItem = &items[i+1]
		}
	}
	pl := o.units[0].UpdateDataPlayer().Player
	var rows []matchRosterMessageRow
	defer func() {
		spellbookCapture(t, "match-roster-inventory", rows, "cb261c2db82e011afaa7f28458154327db7df772a29dd15aefb61403733694d6")
	}()
	for _, to := range []uint32{1, 31, 159, 255} {
		for _, count := range []int{0, 1, 3} {
			for flags := 0; flags < 8; flags++ {
				for types := 0; types < 8; types++ {
					for _, unit := range []bool{false, true} {
						label := fmt.Sprintf("to%d/count%d/flags%d/types%d/unit%t", to, count, flags, types, unit)
						t.Run(label, func(t *testing.T) {
							o.reset()
							*o.roster["flag-type"] = 0
							pl.PlayerUnit = nil
							if unit {
								pl.PlayerUnit = &o.units[0]
							}
							o.units[0].InvFirstItem = nil
							if count > 0 {
								o.units[0].InvFirstItem = &items[0]
							}
							var want [][]byte
							for i := range items {
								it := &items[i]
								it.TypeInd = weapon
								if types&(1<<i) != 0 {
									it.TypeInd = flag
								}
								it.ObjFlags = 0
								if flags&(1<<i) != 0 {
									it.ObjFlags = 0x100
								}
								it.InvNextItem = nil
								if i+1 < count {
									it.InvNextItem = &items[i+1]
								}
								if !unit || i >= count || it.ObjFlags&0x100 == 0 && it.TypeInd != flag {
									continue
								}
								holderCode := uint16(o.units[0].NetCode) | 0x8000 // player holder marker
								data := []byte{80, byte(holderCode), byte(holderCode >> 8), 0, 0, 0, 0}
								if it.TypeInd == weapon {
									binary.LittleEndian.PutUint32(data[3:], 0x12345678)
								}
								want = append([][]byte{data}, want...)
							}
							rv := matchRosterCall("inventory", nil, pl, nil, to)
							state := o.state()
							if rv != 0 || len(state.Nodes) != len(want) {
								t.Fatal("inventory reports", rv, len(state.Nodes), len(want))
							}
							for i, data := range want {
								if !bytes.Equal(state.Nodes[i].Data, data) {
									t.Fatalf("inventory payload %d: got %x want %x (flag type %d, weapon type %d)", i, state.Nodes[i].Data, data, flag, weapon)
								}
							}
							if unit && *o.roster["flag-type"] != uint32(flag) {
								t.Fatal("flag type cache")
							}
							rows = append(rows, matchRosterMessageRow{label, rv, state})
						})
					}
				}
			}
		}
	}
	// Restore the real player linkage before enclosing-owner cleanup.
	pl.PlayerUnit = &o.units[0]
	o.units[0].InvFirstItem = nil
}

func TestMatchRosterRosterRecipients(t *testing.T) {
	o := newMatchRosterOwner(t)
	var rows []matchRosterMessageRow
	defer func() {
		spellbookCapture(t, "match-roster-roster", rows, "bd4f66463ee78144a20fcd3524b457474703a02afecd27174a8e8eecc8c42685")
	}()
	for _, to := range []int{1, 3, 7, 31, 159, 255} {
		for _, headless := range []bool{false, true} {
			for _, name := range []string{"", "a", "abcdefghijk"} {
				label := fmt.Sprintf("to%d/headless%t/name%q", to, headless, name)
				t.Run(label, func(t *testing.T) {
					o.reset()
					noxflags.UnsetEngine(noxflags.EngineNoRendering)
					if headless {
						noxflags.SetEngine(noxflags.EngineNoRendering)
					}
					var want [][]byte
					for pl := o.s.Players.First(); pl != nil; pl = o.s.Players.Next(pl) {
						clear(pl.Field2096Buf[:])
						pl.SetField2096(name)
						if int(pl.PlayerInd) == to || headless && pl.PlayerInd == 31 {
							continue
						}
						buf := make([]byte, 132)
						legacy.PortTestMatchRosterPlayerPacket(buf, pl)
						want = append([][]byte{bytes.Clone(buf[:129])}, want...)
					}
					legacy.PortTestMatchRosterSendPlayers(to)
					state := o.state()
					if len(state.Nodes) != len(want) {
						t.Fatal("roster recipient count", len(state.Nodes), len(want))
					}
					for i, data := range want {
						if !bytes.Equal(state.Nodes[i].Data, data) {
							t.Fatal("roster payload", i)
						}
					}
					rows = append(rows, matchRosterMessageRow{label, 0, state})
				})
			}
		}
	}
}

func TestMatchRosterObjectiveMinimap(t *testing.T) {
	o := newMatchRosterOwner(t)
	t.Cleanup(o.s.PortTestRewardTypes([]string{"Crown"}, nil, true, 0, 0))
	old := o.s.Objs.First()
	t.Cleanup(func() { o.s.Objs.SetObjects(old) })
	types := []uint16{uint16(o.s.Types.ByID("GameBall").Ind()), uint16(o.s.Types.ByID("Crown").Ind()), uint16(o.s.Types.ByID("PortCreatureMonster").Ind())}
	type row struct {
		Name    string
		Entries []uint32
		Flags   [3]uint32
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "match-roster-minimap", rows, "27d9f0bad30ea73845c4170923eb6fb7029951758b468feafe9d2197cb188565")
	}()
	for _, slot := range []ntype.PlayerInd{1, 3, 31} {
		for mask := 0; mask < 8; mask++ {
			for count := 0; count <= 3; count++ {
				label := fmt.Sprintf("slot%d/mask%d/count%d", slot, mask, count)
				t.Run(label, func(t *testing.T) {
					o.reset()
					for i := range o.units {
						u := &o.units[i]
						u.TypeInd = types[i]
						u.ObjFlags = 0
						if mask&(1<<i) != 0 {
							u.ObjFlags = 4
						}
						u.Field5 = 0
						*(*unsafe.Pointer)(unsafe.Add(u.CObj(), 444)) = nil
						if i+1 < count {
							*(*unsafe.Pointer)(unsafe.Add(u.CObj(), 444)) = o.units[i+1].CObj()
						}
					}
					o.s.Objs.SetObjects(nil)
					if count > 0 {
						o.s.Objs.SetObjects(&o.units[0])
					}
					pl := o.s.Players.ByInd(slot)
					defer func() {
						for i := range o.units {
							o.s.Players.Nox_xxx_netUnmarkMinimapObj_417300(slot, &o.units[i], 0xffffffff)
						}
					}()
					for repeat := 0; repeat < 2; repeat++ {
						if matchRosterCall("objective-minimap", nil, nil, nil, uint32(slot)) != 0 {
							t.Fatal("iteration return")
						}
					}
					var entries []uint32
					if p := pl.Field4580; p != nil {
						for q := p; ; q = q.Field8 {
							if q.Field0 != 1 {
								t.Fatal("minimap flags")
							}
							entries = append(entries, q.Field4.NetCode)
							if q.Field8 == p {
								break
							}
							if len(entries) > 3 {
								t.Fatal("minimap cycle")
							}
						}
					}
					var want []uint32
					for i := 0; i < count && i < 2; i++ {
						if mask&(1<<i) != 0 {
							want = append([]uint32{o.units[i].NetCode}, want...)
						}
					}
					if fmt.Sprint(entries) != fmt.Sprint(want) {
						t.Fatal("minimap membership", entries, want)
					}
					flags := [3]uint32{o.units[0].Field5, o.units[1].Field5, o.units[2].Field5}
					for i, v := range flags {
						expected := uint32(0)
						if i < count && i < 2 && mask&(1<<i) != 0 {
							expected = 1
						}
						if v != expected {
							t.Fatal("object minimap status", i)
						}
					}
					rows = append(rows, row{label, entries, flags})
				})
			}
		}
	}
}
