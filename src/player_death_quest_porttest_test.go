//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/opennox/libs/prand"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/ccall"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestPlayerDeathQuestLives(t *testing.T) {
	type record struct {
		Name               string
		Lives, Gold, Frame uint32
		Slots              []byte
		Stats              [11]uint32
		Logic, Other       int
		Reports            legacy.PortTestReliableReportState
	}
	var rows []record
	for _, lives := range []uint32{0, 1, 3, 0xffffffff} {
		for _, starting := range []float64{0, 2.75, 257} {
			name := fmt.Sprintf("lives=%x/starting=%g", lives, starting)
			t.Run(name, func(t *testing.T) {
				o := newMatchRosterOwner(t)
				oldAbilities := noxServer.abilities
				noxServer.abilities.Init(noxServer)
				t.Cleanup(func() { noxServer.abilities = oldAbilities })
				o.balance(map[string]float64{"QuestGameStartingExtraLives": starting})
				// With no carried gems, the shared penalty service must leave its caches
				// unchanged; these real definitions are absent in this owner.
				for _, name := range []string{"Diamond", "Emerald", "Ruby"} {
					if o.s.Types.IndByID(name) != 0 {
						t.Fatal("unexpected gem definition")
					}
				}
				serverConfigOwnBytes(t, 0x5D4594, 2491680, 12)
				o.reset()
				noxflags.ResetGame()
				noxflags.SetGame(noxflags.GameHost | noxflags.GameModeQuest)
				o.s.Rand.Logic, o.s.Rand.Other = prand.New(7), prand.New(11)
				expectedRNG := prand.New(7)
				*memmap.PtrUint32(0x587000, 202028) = 17
				u := &o.units[0]
				u.ObjFlags = 0
				ud := u.UpdateDataPlayer()
				pl := ud.Player
				objectXferSetWord(u.CObj(), 520, 0)
				objectXferSetWord(u.CObj(), 524, 0)
				objectXferSetWord(pl.C(), 3600, 0)
				objectXferSetWord(unsafe.Pointer(ud), 280, 0)
				objectXferSetWord(unsafe.Pointer(ud), 320, lives)
				objectXferSetWord(unsafe.Pointer(ud), 548, 0x11111111)
				*(*byte)(unsafe.Add(pl.C(), 2251)) = 0
				pl.GoldVal = 101
				clear(pl.SpellLvl[:])
				clear(pl.BeastScrollLvl[:])
				slots := unsafe.Slice((*byte)(unsafe.Add(unsafe.Pointer(ud), 452)), 32)
				for i := range slots {
					slots[i] = 0x99
				}
				for i := 0; i < 11; i++ {
					objectXferSetWord(pl.C(), 4652+4*i, 0x10000+uint32(i))
				}
				wantStats := questRuntimeStats(pl.C())
				wantSlots := bytes.Clone(slots)
				wantLives, wantGold, wantFrame := lives-1, uint32(101), uint32(0x11111111)
				var questPacket []byte
				if lives == 0 {
					questPacket = make([]byte, 14)
					questPacket[0], questPacket[1] = 240, 2
					binary.LittleEndian.PutUint16(questPacket[2:], uint16(wantStats[4]))
					binary.LittleEndian.PutUint16(questPacket[4:], uint16(wantStats[5]))
					binary.LittleEndian.PutUint16(questPacket[6:], uint16(wantStats[3]))
					binary.LittleEndian.PutUint16(questPacket[8:], uint16(wantStats[9]))
					wantStats = [11]uint32{0, 0, 0, 0, 0, 0, 0, 0, 0, 17, 63}
					wantLives = uint32(int32(starting))
					wantGold = 51
					wantFrame = o.s.Frame()
					wantSlots[pl.PlayerInd] = byte(wantLives)
					expectedRNG.IntClamp(1, 0) // One warrior-ability penalty selection, even when empty.
				} else {
					wantStats[2]++
					wantStats[10] |= 2
				}
				ccall.CallVoidPtr(server.PortTestPlayerDeathCallback(), u.CObj())
				r := record{Name: name, Lives: objectXferGetWord(unsafe.Pointer(ud), 320), Gold: pl.GoldVal, Frame: objectXferGetWord(unsafe.Pointer(ud), 548), Slots: bytes.Clone(slots), Stats: questRuntimeStats(pl.C()), Logic: o.s.Rand.Logic.Index(), Other: o.s.Rand.Other.Index(), Reports: o.state()}
				if r.Lives != wantLives || r.Gold != wantGold || r.Frame != wantFrame || r.Stats != wantStats || !bytes.Equal(r.Slots, wantSlots) {
					t.Fatalf("quest lives/gold/frame/stats/slots %d/%d/%d/%v/%x want %d/%d/%d/%v/%x", r.Lives, r.Gold, r.Frame, r.Stats, r.Slots, wantLives, wantGold, wantFrame, wantStats, wantSlots)
				}
				if r.Logic != expectedRNG.Index() || r.Other != 11 {
					t.Fatalf("quest RNG %d/%d want %d/11", r.Logic, r.Other, expectedRNG.Index())
				}
				found := 0
				for _, n := range r.Reports.Nodes {
					if len(n.Data) >= 2 && n.Data[0] == 240 && n.Data[1] == 2 {
						found++
						if !bytes.Equal(n.Data, questPacket) {
							t.Fatalf("quest report %x want %x", n.Data, questPacket)
						}
					}
				}
				if found != bool2int(lives == 0) {
					t.Fatalf("quest report count %d", found)
				}
				rows = append(rows, r)
			})
		}
	}
	spellbookCapture(t, "player-death-quest-lives", rows, "759bfb0a0387cda1ce34303575437c0b1e1368638dab1e25301945a390ef8ddf")
}
