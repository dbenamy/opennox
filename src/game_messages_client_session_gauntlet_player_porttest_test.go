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

func TestGameMessageClientSessionGauntletPlayer(t *testing.T) {
	c, _, _ := newEffectsFullOwner(t)
	units, configure, _, free := c.srv.PortTestEscortPlayers()
	t.Cleanup(free)
	configure(2)
	pl, other := units[0].UpdateDataPlayer().Player, units[1].UpdateDataPlayer().Player
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	count := serverConfigOwnBytes(t, 0x5D4594, 1050012, 4)
	oldCode := legacy.ClientPlayerNetCode()
	t.Cleanup(func() { legacy.ClientSetPlayerNetCode(oldCode) })
	type row struct {
		On, Present, Local, Kind int
		Code                     uint16
		Value, Count             uint32
		Return                   int
		Tail                     []byte
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for present := 0; present < 2; present++ {
			for _, local := range []int{17, 0x8011, 999, 0x10011} {
				for _, code := range []uint16{17, 0x8011, 65535} {
					for _, kind := range []int{4, 21, 22, 23} {
						for _, value := range []uint32{0, 1, 127, 128, 255, 0x80000000, 0xffffffff} {
							binary.LittleEndian.PutUint32(connected, uint32(on))
							binary.LittleEndian.PutUint32(count, 0xcafe1234)
							legacy.ClientSetPlayerNetCode(local)
							pl.Active = byte(present)
							pl.NetCodeVal = 17
							other.NetCodeVal = 999
							raw := unsafe.Slice((*byte)(pl.C()), int(unsafe.Sizeof(*pl)))
							for i := 4792; i < len(raw); i++ {
								raw[i] = byte(i*17 + 31)
							}
							want := *pl
							otherBefore := *other
							expected := unsafe.Slice((*byte)(unsafe.Pointer(&want)), int(unsafe.Sizeof(want)))
							data := []byte{240, byte(kind), byte(value), byte(code), byte(code >> 8)}
							nWant := 5
							offset := 4816
							countWant := uint32(0xcafe1234)
							switch kind {
							case 21:
								data = binary.LittleEndian.AppendUint32([]byte{240, 21}, value)
								data = binary.LittleEndian.AppendUint16(data, code)
								nWant = 8
								offset = 4820
							case 22:
								offset = 4824
							case 23:
								offset = 4825
							}
							if on != 0 && present != 0 && code == 17 {
								if kind == 21 {
									binary.LittleEndian.PutUint32(expected[offset:], value)
								} else {
									expected[offset] = byte(value)
								}
							}
							// C compares the 16-bit wire code against the full local code, without masking.
							if on != 0 && kind == 4 && int(code) == local {
								countWant = uint32(byte(value))
							}
							input := bytes.Clone(data)
							n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(240), data)
							if n != nWant || !bytes.Equal(input, data) || *pl != want || *other != otherBefore || binary.LittleEndian.Uint32(count) != countWant {
								t.Fatal("quest player state", on, present, local, code, kind, value, n, binary.LittleEndian.Uint32(count), countWant)
							}
							rows = append(rows, row{on, present, local, kind, code, value, countWant, n, bytes.Clone(raw[4792:])})
						}
					}
				}
			}
		}
	}
	interactionCapture(t, "game-client-session-gauntlet-player", rows)
}
