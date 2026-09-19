//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
	"unsafe"
)

func TestPlayerDeathLifecycle(t *testing.T) {
	type record struct {
		Name    string
		Flags   uint32
		State   byte
		Mana    uint16
		Casting [32]byte
		Reports legacy.PortTestReliableReportState
		Audio   [][4]uint32
	}
	var rows []record
	for _, cause := range []uint32{0, 16} {
		for _, female := range []byte{0, 1} {
			name := fmt.Sprintf("cause=%d/female=%d", cause, female)
			t.Run(name, func(t *testing.T) {
				o := newMatchRosterOwner(t)
				oldAbilities := noxServer.abilities
				noxServer.abilities.Init(noxServer)
				t.Cleanup(func() { noxServer.abilities = oldAbilities })
				o.reset()
				noxflags.ResetGame()
				noxflags.SetGame(noxflags.GameHost)
				u := &o.units[0]
				ud := u.UpdateDataPlayer()
				pl := ud.Player
				u.ObjFlags = 0
				objectXferSetWord(u.CObj(), 520, 0)
				objectXferSetWord(u.CObj(), 524, cause)
				objectXferSetWord(pl.C(), 3600, 0)
				*(*byte)(unsafe.Add(pl.C(), 2252)) = female
				objectXferSetWord(unsafe.Pointer(ud), 280, 0)
				ud.ManaCur = 127
				casting := unsafe.Slice((*byte)(unsafe.Add(unsafe.Pointer(ud), 188)), 32)
				for i := range casting {
					casting[i] = 0x7d
				}
				o.s.PortTestCombatAudioReset()
				legacy.PortTestPlayerDeath("death", u, nil, nil, nil)
				r := record{Name: name, Flags: uint32(u.ObjFlags), State: byte(ud.State), Mana: ud.ManaCur, Reports: o.state()}
				copy(r.Casting[:], casting)
				deathReports := 0
				for _, node := range r.Reports.Nodes {
					if len(node.Data) > 0 && node.Data[0] == 232 {
						deathReports++
						if len(node.Data) != 3 || binary.LittleEndian.Uint16(node.Data[1:]) != uint16(u.NetCode) {
							t.Fatal("death notification bytes")
						}
					}
				}
				if deathReports != 1 {
					t.Fatalf("death notifications: %d", deathReports)
				}
				for _, event := range o.s.PortTestCombatAudioSnapshot() {
					if event.Obj != u {
						t.Fatal("unexpected audio owner")
					}
					r.Audio = append(r.Audio, [4]uint32{uint32(event.ID), u.NetCode, uint32(event.Kind), event.Code})
				}
				expectedSound := uint32(321)
				if female != 0 {
					expectedSound = 331
				}
				if cause == 16 {
					expectedSound = 299
				}
				if len(r.Audio) != 1 || r.Audio[0] != [4]uint32{expectedSound, u.NetCode, 0, 0} {
					t.Fatalf("death sound %v", r.Audio)
				}
				if r.Flags&0x8010 != 0x8010 || r.State != 3 || r.Mana != 0 {
					t.Fatalf("death flags/state/mana %x/%d/%d", r.Flags, r.State, r.Mana)
				}
				for _, off := range []int{0, 24} {
					if r.Casting[off] != 0 {
						t.Fatalf("casting byte %d", off)
					}
				}
				for _, off := range []int{4, 8, 12, 16, 20, 28} {
					if binary.LittleEndian.Uint32(r.Casting[off:]) != 0 {
						t.Fatalf("casting word %d", off)
					}
				}
				for _, off := range []int{1, 2, 3, 25, 26, 27} {
					if r.Casting[off] != 0x7d {
						t.Fatalf("casting padding byte %d changed", off)
					}
				}
				rows = append(rows, r)
			})
		}
	}
	spellbookCapture(t, "player-death-lifecycle", rows, "")
}
