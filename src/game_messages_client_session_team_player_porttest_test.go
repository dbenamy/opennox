//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
	"unsafe"
)

func TestGameMessageClientSessionTeamPlayer(t *testing.T) {
	c, _, _ := newEffectsFullOwner(t)
	units, configure, _, free := c.srv.PortTestEscortPlayers()
	t.Cleanup(free)
	configure(2)
	pl, other := units[0].UpdateDataPlayer().Player, units[1].UpdateDataPlayer().Player
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	type row struct {
		On, Present, Value int
		Code               uint16
		After              byte
		Return             int
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for present := 0; present < 2; present++ {
			for _, code := range []uint16{17, 0x8011, 65535} {
				for value := 0; value < 256; value++ {
					binary.LittleEndian.PutUint32(connected, uint32(on))
					pl.Active = byte(present)
					pl.NetCodeVal = 17
					other.NetCodeVal = 999
					*(*byte)(unsafe.Add(pl.C(), 2282)) = 0x5a
					want := *pl
					otherBefore := *other
					expected := byte(0x5a)
					if on != 0 && present != 0 && code == 17 {
						expected = byte(value)
						*(*byte)(unsafe.Add(unsafe.Pointer(&want), 2282)) = expected
					}
					data := []byte{196, 12, byte(code), byte(code >> 8), byte(value)}
					input := bytes.Clone(data)
					n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(196), data)
					if n != 5 || !bytes.Equal(input, data) || *pl != want || *other != otherBefore {
						t.Fatal("team player flag", on, present, code, value)
					}
					rows = append(rows, row{on, present, value, code, expected, n})
				}
			}
		}
	}
	interactionCapture(t, "game-client-session-team-player", rows)
}

func TestGameMessageClientSessionTeamTradeSelectors(t *testing.T) {
	newEntryOwner(t)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	known := map[int][]int{196: {0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 12}, 201: {1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 13, 27, 29, 31}}
	type row struct{ On, Op, Kind, Return int }
	var rows []row
	for on := 0; on < 2; on++ {
		for _, op := range []int{196, 201} {
			for kind := 0; kind < 256; kind++ {
				skip := false
				for _, v := range known[op] {
					if v == kind {
						skip = true
					}
				}
				if skip {
					continue
				}
				binary.LittleEndian.PutUint32(connected, uint32(on))
				data := []byte{byte(op), byte(kind)}
				input := bytes.Clone(data)
				n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(op), data)
				if n != -1 || !bytes.Equal(data, input) {
					t.Fatal("unknown team/trade selector", on, op, kind, n)
				}
				rows = append(rows, row{on, op, kind, n})
			}
		}
	}
	interactionCapture(t, "game-client-session-team-trade-selectors", rows)
}
