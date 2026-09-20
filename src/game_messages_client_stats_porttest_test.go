//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
	"unsafe"

	"github.com/opennox/libs/noxnet/netmsg"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestGameMessageClientStats(t *testing.T) {
	o := newInventoryDisplayOwner(t)
	t.Cleanup(noxflags.PortTestGameFlags(0))
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	oldCode := legacy.ClientPlayerNetCode()
	t.Cleanup(func() { legacy.ClientSetPlayerNetCode(oldCode) })
	legacy.ClientSetPlayerNetCode(7)
	player := unsafe.Slice((*byte)(unsafe.Pointer(&o.players[0])), int(unsafe.Sizeof(o.players[0])))
	oldPlayer := bytes.Clone(player)
	t.Cleanup(func() { copy(player, oldPlayer) })
	sign := unsafe.Slice((*uint16)(memmap.PtrOff(0x5D4594, 1062588)), 256)
	fallback := unsafe.Slice((*uint16)(memmap.PtrOff(0x5D4594, 1063676)), 32)
	values := []uint16{0, 1, 255, 32767, 32768, 65535, 17, 40000}
	levels := []byte{0, 1, 5, 10, 11, 127, 128, 255}
	type row struct {
		On, Host, Player, Class, Code, Value, Return int
		Stats                                        [6]uint32
		Sign                                         []byte
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for host := 0; host < 2; host++ {
			for present := 0; present < 2; present++ {
				for class := 0; class < 3; class++ {
					for _, code := range []uint16{7, 0x8007, 8, 0xffff} {
						for v, health := range values {
							o.reset(t)
							noxflags.ResetGame()
							noxflags.SetGame(noxflags.GameFlag(host))
							binary.LittleEndian.PutUint32(connected, uint32(on))
							player[2251], player[3684] = byte(class), 1
							binary.LittleEndian.PutUint16(player[3652:], 1234)
							o.players[0].SetName("Port Ω")
							if present == 0 {
								*o.displayWords["dword_8531A0_2576"] = 0
							}
							alloc.StrCopyZero16(sign, "untouched")
							alloc.StrCopyZero16(fallback, "Fallback")
							data := make([]byte, 14)
							data[0] = 72
							binary.LittleEndian.PutUint16(data[1:], code)
							binary.LittleEndian.PutUint16(data[3:], health)
							binary.LittleEndian.PutUint16(data[5:], values[(v+1)%len(values)])
							binary.LittleEndian.PutUint16(data[7:], values[(v+2)%len(values)])
							binary.LittleEndian.PutUint16(data[9:], values[(v+3)%len(values)])
							binary.LittleEndian.PutUint16(data[11:], values[(v+4)%len(values)])
							data[13] = levels[v]
							input := bytes.Clone(data)
							want := bytes.Clone(player)
							admitted := on != 0 && code&0x7fff == 7
							if admitted && host == 0 && present != 0 {
								for _, pair := range [][2]int{{2247, 3}, {2243, 5}, {2235, 9}, {2239, 11}} {
									binary.LittleEndian.PutUint32(want[pair[0]:], uint32(binary.LittleEndian.Uint16(data[pair[1]:])))
								}
								copy(want[3652:3654], data[7:9])
								want[3684] = data[13]
							}
							wantText := "untouched"
							if admitted {
								wantText = "Fallback"
								if present != 0 {
									wantText = fmt.Sprintf("Port Ω the Rank%d-%d", class, int(int8(want[3684])))
								}
							}
							n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(72), data)
							gotText := alloc.GoString16(&sign[0])
							if n != 14 || !bytes.Equal(input, data) || !bytes.Equal(player, want) || gotText != wantText || len(o.sounds) != 0 || len(o.displayText) != 0 {
								t.Fatalf("stats on%d host%d player%d class%d code%x value%d: return%d title%q want%q playerEqual%v", on, host, present, class, code, v, n, gotText, wantText, bytes.Equal(player, want))
							}
							r := row{On: on, Host: host, Player: present, Class: class, Code: int(code), Value: v, Return: n, Sign: bytes.Clone(unsafe.Slice((*byte)(unsafe.Pointer(&sign[0])), 512))}
							for i, off := range []int{2247, 2243, 2235, 2239} {
								r.Stats[i] = binary.LittleEndian.Uint32(player[off:])
							}
							r.Stats[4], r.Stats[5] = uint32(binary.LittleEndian.Uint16(player[3652:])), uint32(player[3684])
							rows = append(rows, r)
						}
					}
				}
			}
		}
	}
	interactionCapture(t, "game-progress-stats", rows)
}
