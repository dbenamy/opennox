//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/music"
)

func TestGameMessageClientSessionMusic(t *testing.T) {
	words, restore := legacy.PortTestAudioEventGlobals()
	t.Cleanup(restore)
	old := legacy.MusicModule
	legacy.MusicModule = &music.Module{}
	t.Cleanup(func() { legacy.MusicModule = old })
	memory := serverConfigOwnBytes(t, 0x5D4594, 815772, 320)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	count, level := words["dword_5d4594_816368"], words["dword_5d4594_816372"]
	put := func(dst []byte, s music.MusicState) {
		for i, v := range []uint32{s.MusicIdx, s.Volume, s.Position, s.D} {
			binary.LittleEndian.PutUint32(dst[4*i:], v)
		}
	}
	get := func(src []byte) music.MusicState {
		return music.MusicState{MusicIdx: binary.LittleEndian.Uint32(src), Volume: binary.LittleEndian.Uint32(src[4:]), Position: binary.LittleEndian.Uint32(src[8:]), D: binary.LittleEndian.Uint32(src[12:])}
	}
	type row struct {
		On, Kind, ID, Volume, Return int
		Level, Count, After          uint32
		Current                      music.MusicState
		Memory                       []byte
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for _, kind := range []int{229, 230, 231} {
			for l := uint32(0); l < 3; l++ {
				for _, n := range []uint32{0, 1, 5, 6} {
					for _, id := range []int{0, 1, 255} {
						for _, volume := range []int{0, 100, 101, 255} {
							for i := range memory {
								memory[i] = 0xa5
							}
							for slot := uint32(0); slot < 6; slot++ {
								put(memory[16*(slot+6*l):], music.MusicState{MusicIdx: 100 + slot, Volume: 90 + slot*7, Position: slot * 1000, D: slot ^ 0x87654321})
							}
							*count, *level = n, l
							binary.LittleEndian.PutUint32(connected, uint32(on))
							current := music.MusicState{MusicIdx: 17, Volume: 63, Position: 0x12345678, D: 0xabcdef01}
							legacy.MusicModule.SetNextMusic(current)
							want, wantCount, expected := current, n, bytes.Clone(memory)
							if on != 0 {
								switch kind {
								case 229:
									want = music.MusicState{MusicIdx: uint32(id), Volume: uint32(volume)}
									if want.Volume > 100 {
										want.Volume = 100
									}
								case 230:
									if n < 6 {
										put(expected[16*(n+6*l):], current)
										wantCount = n + 1
									} else {
										wantCount = 6
									}
								case 231:
									if n > 0 {
										want = get(memory[16*(n-1+6*l):])
										if want.Volume > 100 {
											want.Volume = 100
										}
									}
									wantCount = 0
								}
							}
							data := []byte{byte(kind), byte(id), byte(volume)}
							input := bytes.Clone(data)
							ret := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(kind), data)
							if ret != 3 || !bytes.Equal(data, input) || *count != wantCount || *level != l || legacy.MusicModule.GetCurrentBlock() != want || !bytes.Equal(memory, expected) {
								t.Fatalf("music on%d kind%d level%d count%d id%d volume%d ret%d", on, kind, l, n, id, volume, ret)
							}
							rows = append(rows, row{on, kind, id, volume, ret, l, n, *count, legacy.MusicModule.GetCurrentBlock(), bytes.Clone(memory)})
						}
					}
				}
			}
		}
	}
	interactionCapture(t, "game-client-session-music", rows)
}
