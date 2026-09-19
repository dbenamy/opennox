//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/opennox/libs/object"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/common/ccall"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestPlayerDeathInformation(t *testing.T) {
	o := newMatchRosterOwner(t)
	oldAbilities := noxServer.abilities
	noxServer.abilities.Init(noxServer)
	t.Cleanup(func() { noxServer.abilities = oldAbilities })
	type record struct {
		Name    string
		Source  uint32
		Packets [][]byte
		Reports legacy.PortTestReliableReportState
	}
	var rows []record
	for _, source := range []string{"none", "self", "player", "monster", "object", "player-child", "monster-child"} {
		for _, weapon := range [][2]uint32{{0x1234, 0}, {0xabcd, 5}, {2, 2}} {
			for _, online := range []bool{false, true} {
				name := fmt.Sprintf("source=%s/weapon=%x,%d/online=%t", source, weapon[0], weapon[1], online)
				t.Run(name, func(t *testing.T) {
					o.reset()
					o.s.PortTestCombatAudioReset()
					noxflags.ResetGame()
					noxflags.SetGame(noxflags.GameHost)
					if online {
						noxflags.SetGame(noxflags.GameOnline)
					}
					for i := range o.units {
						u := &o.units[i]
						u.ObjFlags = 0
						u.UpdateDataPlayer().State = 0
						u.UpdateDataPlayer().Player.PlayerUnit = u
						objectXferSetWord(unsafe.Pointer(u.UpdateDataPlayer()), 280, 0)
					}
					u := &o.units[0]
					ud := u.UpdateDataPlayer()
					pl := ud.Player
					objectXferSetWord(pl.C(), 3600, 0)
					objectXferSetWord(u.CObj(), 524, 0)
					objectXferSetWord(unsafe.Pointer(ud), 300, weapon[0])
					objectXferSetWord(unsafe.Pointer(ud), 304, weapon[1])
					monster, freeMonster := alloc.New(server.Object{})
					defer freeMonster()
					monster.ObjClass = object.ClassMonster
					monster.TypeInd = 0xbeef
					child, freeChild := alloc.New(server.Object{})
					defer freeChild()
					child.ObjClass = object.ClassSimple
					defer o.s.ObjClearOwner(child)
					var cause *server.Object
					killerCode := uint16(0)
					kind := byte(weapon[1])
					which := uint16(weapon[0])
					switch source {
					case "self":
						cause = u
						killerCode = uint16(u.NetCode)
					case "player":
						cause = &o.units[1]
						killerCode = uint16(cause.NetCode)
					case "monster":
						cause = monster
						kind = 1
						which = cause.TypeInd
					case "object":
						cause = child
					case "player-child", "monster-child":
						cause = child
						o.s.ObjSetOwner(&o.units[1], cause)
						if source == "player-child" {
							killerCode = uint16(o.units[1].NetCode)
						} else {
							o.s.ObjSetOwner(monster, cause)
							kind = 1
							which = monster.TypeInd
						}
					}
					objectXferSetWord(u.CObj(), 520, uint32(uintptr(unsafe.Pointer(cause))))
					ccall.CallVoidPtr(server.PortTestPlayerDeathCallback(), u.CObj())
					r := record{Name: name, Source: objectXferGetWord(unsafe.Pointer(ud), 304), Reports: o.state()}
					want := []byte{169, 14, 0, 0, 0, 0, byte(u.NetCode), byte(u.NetCode >> 8), byte(which), byte(which >> 8), kind}
					binary.LittleEndian.PutUint16(want[2:], killerCode)
					for i := range o.units {
						p := o.units[i].UpdateDataPlayer().Player
						packet := o.s.NetList.CopyPacketsA(p.PlayerIndex(), netlist.Kind1)
						if online && !bytes.HasPrefix(packet, want) {
							t.Fatalf("player %d information prefix %x want %x", i, packet, want)
						}
						if !online && bytes.Contains(packet, []byte{169, 14}) {
							t.Fatalf("offline death information %x", packet)
						}
						r.Packets = append(r.Packets, packet)
					}
					expectedSource := weapon[1]
					if online {
						expectedSource = 0
					}
					if r.Source != expectedSource {
						t.Fatalf("source reset %d want %d", r.Source, expectedSource)
					}
					rows = append(rows, r)
				})
			}
		}
	}
	spellbookCapture(t, "player-death-information", rows, "7a1e35d219a88681678478726557941f3b251a95a5e27809816c5f6c0cacef49")
}
