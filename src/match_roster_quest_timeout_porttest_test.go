//go:build porttest

package opennox

import (
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
	"unsafe"
)

func TestMatchRosterQuestTimeout(t *testing.T) {
	type row struct {
		Name           string
		Return, Timer  uint32
		Path, Waypoint string
		Stage          [3]uint32
		Statistics     [3][11]uint32
		Queue          legacy.PortTestReliableReportState
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "match-roster-quest-timeout", rows, "47006a4b7850544651280080aeb0b7f2d8303b3e2f9bcb0715b64c0a7db7f5ac")
	}()
	for _, name := range []string{"WorldQuest.map", "WorldQuest.map:Start"} {
		for _, parts := range [][3]uint32{{1, 1, 1}, {0, 1, 2}, {0, 0, 0}} {
			for _, exits := range []uint32{0, 5, 7} {
				label := fmt.Sprintf("%s/parts%v/exits%x", name, parts, exits)
				t.Run(label, func(t *testing.T) {
					o := newMatchRosterOwner(t)
					worldCollisionExitGlobals(t, o.worldCollisionOwner)
					restore := noxflags.PortTestGameFlags(noxflags.GameModeQuest)
					defer restore()
					oldWP := noxServer.mapSwitchWPName
					t.Cleanup(func() { noxServer.mapSwitchWPName = oldWP })
					legacy.Set_dword_5d4594_1548524(0)
					pending := unsafe.Slice(memmap.PtrUint8(0x5D4594, 1567844), 96)
					clear(pending)
					copy(pending, name)
					times := unsafe.Slice(memmap.PtrUint8(0x5D4594, 3500), 6)
					for i := range times {
						times[i] = 1
					}
					*memmap.PtrUint64(0x5D4594, 3468) = 0
					*memmap.PtrUint32(0x587000, 4660) = 1
					for i := range o.units {
						u := &o.units[i]
						p := u.UpdateDataPlayer().Player.C()
						objectXferSetWord(p, 4792, parts[i])
						objectXferSetWord(p, 4652, 10)
						objectXferSetWord(p, 4692, 0x10)
						ptr := unsafe.Pointer(nil)
						if exits&(1<<i) != 0 {
							ptr = o.units[i].CObj()
						}
						*(*unsafe.Pointer)(unsafe.Add(u.UpdateData, 312)) = ptr
					}
					if parts == [3]uint32{} { // Exercise the no-player-unit timeout branch.
						for i := range o.units {
							o.units[i].UpdateDataPlayer().Player.PlayerUnit = nil
						}
					}
					rv := matchRosterCall("check-limit", nil, nil, nil)
					if rv != 0 || memmap.Uint32(0x587000, 4660) != 0 || legacy.Get_dword_5d4594_1548524() != 1 {
						t.Fatal("quest timeout transition")
					}
					path := alloc.GoString(memmap.PtrUint8(0x85B3FC, 36))
					wp := ""
					if name != "WorldQuest.map" {
						wp = "Start"
					}
					if path != "worldquest" || alloc.GoString(memmap.PtrUint8(0x5D4594, 2598188)) != "worldquest.map" || noxServer.mapSwitchWPName != wp {
						t.Fatal("quest map request", path, noxServer.mapSwitchWPName)
					}
					r := row{Name: label, Return: rv, Path: path, Waypoint: wp, Queue: o.state()}
					for i := range o.units {
						p := o.units[i].UpdateDataPlayer().Player.C()
						r.Statistics[i] = questRuntimeStats(p)
						r.Stage[i] = objectXferGetWord(p, 4652)
						want := uint32(10)
						flags := uint32(0x10)
						if parts[i] == 1 && exits&(1<<i) == 0 {
							want++
							flags |= 1
						}
						if r.Stage[i] != want || r.Statistics[i][10] != flags {
							t.Fatal("quest completion statistics", i, r.Stage[i], r.Statistics[i][10])
						}
					}
					rows = append(rows, r)
				})
			}
		}
	}
}
