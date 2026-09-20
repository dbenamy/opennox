//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"

	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/opennox/v1/legacy"
)

func TestGameMessageClientFriends(t *testing.T) {
	o := newCombatOverlayOwner(t)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	seeds := [][]uint32{nil, {1, 1, 2}, nil, nil}
	for i := 0; i < 127; i++ {
		seeds[2] = append(seeds[2], uint32(i))
	}
	for i := 0; i < 128; i++ {
		seeds[3] = append(seeds[3], uint32(i))
	}
	type row struct {
		On, Kind, Seed, Code, Return int
		Friends                      []uint32
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for _, kind := range []int{52, 53, 54} {
			for seed, values := range seeds {
				codes := []uint16{0, 1, 32767, 32768, 32769, 65535}
				if kind == 54 {
					codes = []uint16{0}
				}
				for _, code := range codes {
					legacy.PortTestCombatFriendClear()
					binary.LittleEndian.PutUint32(connected, uint32(on))
					var want []uint32
					for _, v := range values {
						if legacy.PortTestCombatFriendAdd(v) == nil {
							t.Fatal("seed allocation")
						}
						want = append([]uint32{v}, want...)
					}
					data := []byte{byte(kind), byte(code), byte(code >> 8)}
					if kind == 54 {
						data = data[:1]
					}
					before := bytes.Clone(data)
					if on != 0 {
						switch kind {
						case 52:
							if len(want) < 128 {
								want = append([]uint32{uint32(code & 0x7fff)}, want...)
							}
						case 53:
							for i, v := range want {
								if v == uint32(code&0x7fff) {
									want = append(want[:i], want[i+1:]...)
									break
								}
							}
						case 54:
							want = nil
						}
					}
					n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(kind), data)
					got := o.friendCodes(t)
					if n != len(data) || !bytes.Equal(data, before) || !reflect.DeepEqual(got, want) {
						t.Fatalf("friends on%d kind%d seed%d code%x return%d got%v want%v", on, kind, seed, code, n, got, want)
					}
					rows = append(rows, row{on, kind, seed, int(code), n, got})
				}
			}
		}
	}
	interactionCapture(t, "game-client-friends", rows)
}
