//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/client/gui"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"image"
	"testing"
	"unsafe"
)

func TestTeamUIClientJoin(t *testing.T) {
	type record struct {
		OldTeam, Observer int
		Missing           bool
		Pending           uint32
		Data              []byte
		Team              byte
	}
	var records []record
	for _, oldTeam := range []int{0, 1, 2} {
		for _, observer := range []int{0, 1, 2, 3} {
			for _, missing := range []bool{false, true} {
				t.Run(fmt.Sprintf("%d/%d/%v", oldTeam, observer, missing), func(t *testing.T) {
					o := newTeamUIOwner(t)
					defer noxflags.PortTestGameFlags(0)()
					clear(o.players)
					oldCode := legacy.ClientPlayerNetCode()
					t.Cleanup(func() { legacy.ClientSetPlayerNetCode(oldCode) })
					legacy.ClientSetPlayerNetCode(777)
					pl := &o.players[0]
					if observer == 1 {
						*(*uint32)(unsafe.Add(pl.C(), 4)) = 1
					}
					if observer == 2 {
						pl.Field3680 = 1
					}
					if observer == 3 {
						legacy.Set_dword_8531A0_2576(nil)
					} else {
						legacy.Set_dword_8531A0_2576(pl)
					}
					dr := o.drawable(777, image.Pt(20, 20))
					dr.ObjClass &^= 0x20400000
					if oldTeam != 0 {
						legacy.Nox_xxx_createAtImpl_4191D0(server.TeamID(oldTeam), dr.TeamPtr(), 0, 777, 0)
					}
					if missing {
						dr.NetCode32 = 778
					}
					reset, snapshot, freeQueue := legacy.PortTestVisibilityReliableOwner()
					t.Cleanup(freeQueue)
					reset()
					legacy.PortTestTeamUI("team-join", nil, o.c.srv.Teams.ByID(2), 0, 0, "")
					pending := *memmap.PtrUint32(0x5D4594, 1045696)
					packets := snapshot()
					var got []byte
					if len(packets) > 1 {
						t.Fatalf("multiple requests: %v", packets)
					}
					if len(packets) == 1 {
						got = packets[0].Data
						if packets[0].Recipient != 31 {
							t.Error("request recipient")
						}
					}
					allowed := !missing && (observer == 0 || observer == 3)
					wantPending := uint32(0)
					var want []byte
					if allowed {
						wantPending = 1
						if oldTeam != 2 {
							want = make([]byte, 10)
							want[0] = 196
							want[1] = 10
							if oldTeam != 0 {
								want[1] = 11
							}
							binary.LittleEndian.PutUint32(want[2:], 2)
							binary.LittleEndian.PutUint16(want[6:], 777)
						}
					}
					if pending != wantPending || !bytes.Equal(got, want) {
						t.Errorf("pending=%d packet=%x want %d/%x", pending, got, wantPending, want)
					}
					if int(dr.TeamVal.ID) != oldTeam {
						t.Error("request changed membership before server response")
					}
					teamUICall("requests-reset", 0, 0)
					if *memmap.PtrUint32(0x5D4594, 1045696) != 0 {
						t.Error("pending reset")
					}
					records = append(records, record{oldTeam, observer, missing, pending, got, byte(dr.TeamVal.ID)})
				})
			}
		}
	}
	spellbookCapture(t, "team-ui-client-join", records, "da45c37d6188cc02a45bc9b15c5c83f2d160df377135477d4c88e4f315abffed")
}
func TestTeamUIHostAssign(t *testing.T) {
	type record struct {
		Selection int
		Observer  bool
		IDs       []byte
		Queue     []legacy.PortTestShopPacketResult
	}
	var records []record
	for _, selection := range []int{-1, 0, 1, 3} {
		for _, observer := range []bool{false, true} {
			t.Run(fmt.Sprintf("%d/%v", selection, observer), func(t *testing.T) {
				o := newTeamUIOwner(t)
				defer noxflags.PortTestGameFlags(1)()
				w := o.openPlayers(t)
				units, configure, _, free := o.c.srv.PortTestEscortPlayers()
				t.Cleanup(free)
				configure(3)
				o.c.srv.Players.ByInd(3).Active = 0
				t.Cleanup(server.PortTestAttachAI(o.c.srv.Server, &units[0], &units[1], &units[2]))
				t.Cleanup(legacy.PortTestCreatureXferLookupOwner())
				reset, snapshot, freeQueue := legacy.PortTestVisibilityReliableOwner()
				t.Cleanup(freeQueue)
				for i := range units {
					u := &units[i]
					u.NetCode = uint32(1001 + i)
					p := u.UpdateDataPlayer().Player
					p.NetCodeVal = u.NetCode
					p.SetName(fmt.Sprintf("Player %d", i))
					legacy.Nox_xxx_createAtImpl_4191D0(1, u.TeamPtr(), 0, int(u.NetCode), 0)
				}
				if observer {
					units[0].UpdateDataPlayer().Player.Field3680 = 1
				}
				teamUICall("players-refresh", 0, 0)
				w.ChildByID(10502).Func94(&gui.RawEvent{Event: 16403, Arg1: 1})
				list := w.ChildByID(10501)
				list.Func94(&gui.RawEvent{Event: 16403, Arg1: ^uintptr(0)})
				for i := 0; i < 3; i++ {
					if selection >= 0 && selection&(1<<i) != 0 {
						list.Func94(&gui.RawEvent{Event: 16405, Arg1: uintptr(i)})
					}
				}
				reset()
				legacy.PortTestTeamUIEvent(w, w.ChildByID(10503), 16391, 0)
				var ids []byte
				for i := range units {
					want := byte(1)
					if selection >= 0 && selection&(1<<i) != 0 && !(observer && i == 0) {
						want = 2
					}
					if got := byte(units[i].TeamVal.ID); got != want {
						t.Errorf("player %d team %d want %d", i, got, want)
					}
					ids = append(ids, byte(units[i].TeamVal.ID))
				}
				if int32((*gui.ScrollListBoxData)(w.ChildByID(10502).WidgetData).Field_12) != -1 {
					t.Error("team selection not cleared")
				}
				data := (*gui.ScrollListBoxData)(list.WidgetData)
				if *(*int32)(unsafe.Pointer(uintptr(data.Field_12))) != -1 {
					t.Error("player selection not cleared")
				}
				records = append(records, record{selection, observer, ids, snapshot()})
			})
		}
	}
	spellbookCapture(t, "team-ui-host-assign", records, "beaaf8b014031d31de2f1e6d0f09584b72dea71ae1bac30e20028af19fb02920")
}
