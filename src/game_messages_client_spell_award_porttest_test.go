//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/libs/prand"
	"github.com/opennox/libs/things"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"reflect"
	"testing"
	"unsafe"
)

func TestGameMessageClientSpellAward(t *testing.T) {
	o, bar, configure := newSpellbookAdditionOwner(t)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	serverConfigOwnBytes(t, 0x5D4594, 1217504, 4)
	type state struct {
		Book            spellbookResult
		Particles       [][]uint32
		RNG             [2]int
		Slot, Timestamp uint32
	}
	type row struct {
		On, Rank, Flags int
		State           state
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for _, rank := range []byte{0, 1, 255} {
			for _, flags := range []byte{0, 1, 2, 63, 127, 128, 129, 130, 191, 255} {
				var expected state
				for phase := 0; phase < 2; phase++ {
					prepareSpellbookAddition(t, o, bar, 640)
					o.c.srv.Rand.Logic, o.c.srv.Rand.Other = prand.New(31), prand.New(32)
					*memmap.PtrUint32(0x5D4594, 1217504) = 0
					binary.LittleEndian.PutUint32(connected, uint32(on))
					configure([]server.PortTestSpellClassDef{{Index: 1, Flags: uint32(things.SpellClassAny), Valid: true}})
					*(*byte)(unsafe.Add(unsafe.Pointer(&o.players[0]), 2251)) = 1
					particles, free := legacy.PortTestEffectsScreenParticles(768)
					data := []byte{111, 1, rank, flags}
					input := bytes.Clone(data)
					if phase == 0 {
						if on != 0 {
							o.bookCall("nox_xxx_netSpellRewardCli_45CFE0", 1, uint32(rank), uint32(flags&127), uint32(flags>>7))
						}
					} else {
						if n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(111), data); n != 4 {
							t.Fatal("spell award length")
						}
					}
					got := state{o.bookSnapshot("award", 0), particles(), [2]int{o.c.srv.Rand.Logic.Index(), o.c.srv.Rand.Other.Index()}, bar[0], memmap.Uint32(0x5D4594, 1217504)}
					free()
					if !bytes.Equal(input, data) {
						t.Fatal("spell award input")
					}
					wantRank := uint32(0)
					if on != 0 {
						wantRank = uint32(rank)
					}
					if got.Book.Known[1] != wantRank {
						t.Fatal("spell award rank/connection gate")
					}
					if phase == 0 {
						expected = got
						continue
					}
					if !reflect.DeepEqual(expected, got) {
						t.Fatalf("spell award dispatch on%d rank%d flags%d", on, rank, flags)
					}
					rows = append(rows, row{on, int(rank), int(flags), got})
				}
			}
		}
	}
	interactionCapture(t, "game-progress-spell-award", rows)
}
