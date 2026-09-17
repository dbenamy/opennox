//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"strings"
	"testing"
)

func TestVisibilityEffectsPrediction(t *testing.T) {
	s, units, configure := visibilityEffectsPlayers(t)
	u := newCreatureXferObject(t, s, "Monster")
	u.NetCode = 0x1234
	var rows []visibilityEffectsRow
	defer func() {
		spellbookCapture(t, "visibility-effects-prediction", rows, "f23426209a392d8ab8decd29218c11f3b5c94b8fef332b7759cee71fd8f4157f")
	}()
	for _, n := range []int{0, 1, 2, 3} {
		configure(n)
		for _, v := range [][5]float32{{0, 0, 0, 0, 0}, {100.75, -200.25, 1.25, -2.125, 16.5}, {65537.5, -65537.5, 4096, -4096, 0.0625}, {2147483648, -2147483648, 0.5, 0.75, -0.75}} {
			for _, dir := range []uint16{0, 255, 256, 65535} {
				name := fmt.Sprintf("players%d/values%v/dir%d", n, v, dir)
				t.Run(name, func(t *testing.T) {
					s.NetList.ResetAll()
					u.PosVec = types.Pointf{v[0], v[1]}
					u.Float28 = v[2]
					u.VelVec = types.Pointf{v[3], v[4]}
					u.Direction1 = server.Dir16(dir)
					rv := legacy.PortTestVisibilityEffects(6, u, nil, nil, nil, [5]int32{}, nil, "")
					want := []byte{181, 0x34, 0x12}
					want = binary.LittleEndian.AppendUint16(want, u.TypeInd)
					want = binary.LittleEndian.AppendUint16(want, uint16(int64(v[0])))
					want = binary.LittleEndian.AppendUint16(want, uint16(int64(v[1])))
					want = binary.LittleEndian.AppendUint16(want, dir)
					want = append(want, byte(int64(float64(v[2])*16)), byte(int64(float64(v[3])*16)), byte(int64(float64(v[4])*16)))
					got := visibilityEffectsPackets(s)
					for i := range got {
						expected := []byte{}
						for j := 0; j < n; j++ {
							if i == int(units[j].UpdateDataPlayer().Player.PlayerInd) {
								expected = want
							}
						}
						if !bytes.Equal(got[i], expected) {
							t.Fatalf("player%d got%x want%x", i, got[i], expected)
						}
					}
					if rv != 0 {
						t.Fatalf("return%x", rv)
					}
					rows = append(rows, visibilityEffectsRow{name, rv, got})
				})
			}
		}
	}
}
func TestVisibilityEffectsChat(t *testing.T) {
	s, units, configure := visibilityEffectsPlayers(t)
	u := newCreatureXferObject(t, s, "Monster")
	u.NetCode = 0x1234
	u.PosVec = types.Pointf{100.75, -200.25}
	var rows []visibilityEffectsRow
	defer func() {
		spellbookCapture(t, "visibility-effects-chat", rows, "2e1bc29f467f305715c0815311a1e892d4809d2493b6882ef22594f8879533f4")
	}()
	for _, n := range []int{0, 1, 3} {
		configure(n)
		for _, length := range []int{0, 1, 10, 254, 255, 256, 508} {
			for _, recipient := range []int{-2, -1, 0, 1, 2} {
				name := fmt.Sprintf("players%d/length%d/to%d", n, length, recipient)
				t.Run(name, func(t *testing.T) {
					s.NetList.ResetAll()
					var source *server.Object
					if recipient == -2 {
						source = u
					} else if recipient >= 0 {
						source = &units[recipient]
					}
					command := strings.Repeat("x", length)
					rv := legacy.PortTestVisibilityEffects(28, u, source, nil, nil, [5]int32{-1234}, nil, command)
					// Length is a byte in the original protocol, including the terminating zero.
					want := []byte{168, 0x34, 0x12, 8, 100, 0, 56, 255, byte(length + 1), 46, 251}
					want = append(want, []byte(command)...)
					want = append(want, 0)
					want = want[:11+int(byte(length+1))]
					got := visibilityEffectsPackets(s)
					for i := range got {
						expected := []byte{}
						if recipient == -1 {
							for j := 0; j < n; j++ {
								if i == int(units[j].UpdateDataPlayer().Player.PlayerInd) {
									expected = want
								}
							}
						} else if recipient >= 0 && i == int(units[recipient].UpdateDataPlayer().Player.PlayerInd) {
							expected = want
						}
						if !bytes.Equal(got[i], expected) {
							t.Fatalf("player%d got%x want%x", i, got[i], expected)
						}
					}
					norm := rv
					if recipient == -2 {
						if rv != uint32(uintptr(u.CObj())) {
							t.Fatalf("return%x", rv)
						}
						norm = 100
					} else if recipient == -1 {
						if rv != 0 {
							t.Fatalf("return%x", rv)
						}
					} else if rv != 1 {
						t.Fatalf("return%x", rv)
					}
					rows = append(rows, visibilityEffectsRow{name, norm, got})
				})
			}
		}
	}
	configure(3)
	s.NetList.ResetAll()
	rv := legacy.PortTestVisibilityEffects(29, nil, nil, nil, nil, [5]int32{}, nil, "")
	got := visibilityEffectsPackets(s)
	for _, slot := range []int{1, 7, 31} {
		if !bytes.Equal(got[slot], []byte{202, 173, 222}) {
			t.Fatalf("chat clear player%d %x", slot, got[slot])
		}
	}
	rows = append(rows, visibilityEffectsRow{"clear", rv, got})
}
