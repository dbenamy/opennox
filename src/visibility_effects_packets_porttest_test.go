//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestVisibilityEffectsPointPackets(t *testing.T) {
	s, units, configure := visibilityEffectsPlayers(t)
	configure(3)
	var rows []visibilityEffectsRow
	defer func() {
		spellbookCapture(t, "visibility-effects-point-packets", rows, "dc805c016fe651145abc7127e070e60a846bdc8922b38f8756f3a489f6e6ba75")
	}()
	for _, xy := range [][2]float32{{0, 0}, {100.75, 200.25}, {-1.75, -2.25}, {32767.5, 32768.75}, {65535.5, 65536.5}, {-65537.5, 131073.5}} {
		pos := types.Pointf{xy[0], xy[1]}
		for i := range units {
			units[i].PosVec = pos
			pl := units[i].UpdateDataPlayer().Player
			pl.Field10 = 100
			pl.Field12 = 100
		}
		for _, extra := range []int32{0, 1, 127, 128, 255, 256, -1, 0x12345678} {
			for _, op := range []int{0, 2, 3, 4, 11} {
				name := fmt.Sprintf("xy%v/extra%d/op%d", xy, extra, op)
				t.Run(name, func(t *testing.T) {
					s.NetList.ResetAll()
					args := [5]int32{extra, extra ^ 0x55}
					rv := legacy.PortTestVisibilityEffects(op, nil, nil, &pos, nil, args, nil, "")
					var want []byte
					off := 1
					switch op {
					case 0:
						want = []byte{byte(extra), 0, 0, 0, 0}
					case 2:
						want = []byte{byte(extra), byte(args[1]), 0, 0, 0, 0}
						off = 2
					case 3:
						want = []byte{147, 0, 0, 0, 0, byte(extra)}
					case 4:
						want = []byte{240, 25, 0, 0, 0, 0, byte(extra)}
						off = 2
					case 11:
						want = []byte{161, 0, 0, 0, 0, byte(extra)}
					}
					// Inputs are finite and exactly inside the signed 32-bit conversion domain.
					x, y := uint16(int32(xy[0])), uint16(int32(xy[1]))
					binary.LittleEndian.PutUint16(want[off:], x)
					binary.LittleEndian.PutUint16(want[off+2:], y)
					got := visibilityEffectsPackets(s)
					if rv != 0 {
						t.Fatalf("return %x", rv)
					}
					for i := range got {
						expected := []byte{}
						if i == 1 || i == 7 || i == 31 {
							expected = want
						}
						if !bytes.Equal(got[i], expected) {
							t.Fatalf("player%d got%x want%x", i, got[i], expected)
						}
					}
					rows = append(rows, visibilityEffectsRow{name, rv, got})
				})
			}
		}
	}
}
func TestVisibilityEffectsDestinationDomains(t *testing.T) {
	s, units, configure := visibilityEffectsPlayers(t)
	configure(3)
	var rows []visibilityEffectsRow
	defer func() {
		spellbookCapture(t, "visibility-effects-destinations", rows, "dc0132a00e582242edeefa209dde210e40e86443d7734db52cfede9801e66098")
	}()
	for _, words := range [][4]int32{{1, 2, 100, 200}, {-1, -2, -100, -200}, {0x12345678, -65537, 65536 + 100, 131072 + 200}, {32768, 65535, 32768, 65535}} {
		units[0].PosVec = types.Pointf{float32(uint16(words[2])), float32(uint16(words[3]))}
		units[1].PosVec = types.Pointf{float32(words[2]), float32(words[3])}
		units[2].PosVec = types.Pointf{-1000000, -1000000}
		for i := range units {
			pl := units[i].UpdateDataPlayer().Player
			pl.Field10 = 0
			pl.Field12 = 0
		}
		for _, id := range []int32{0, 255, 32768, 65535, 65536, -1} {
			for _, op := range []int{5, 10} {
				name := fmt.Sprintf("words%v/id%d/op%d", words, id, op)
				t.Run(name, func(t *testing.T) {
					s.NetList.ResetAll()
					args := [5]int32{0x94, id}
					if op == 10 {
						args[0] = id
					}
					rv := legacy.PortTestVisibilityEffects(op, nil, nil, nil, &words, args, nil, "")
					want := []byte{0x94}
					if op == 10 {
						want = []byte{240, 16}
					}
					for _, v := range words {
						want = binary.LittleEndian.AppendUint16(want, uint16(v))
					}
					want = binary.LittleEndian.AppendUint16(want, uint16(id))
					got := visibilityEffectsPackets(s)
					for i := range got {
						expected := []byte{}
						same := words[2] == int32(uint16(words[2])) && words[3] == int32(uint16(words[3]))
						if i == 1 && (op == 5 || same) || i == 7 && (op == 10 || same) {
							expected = want
						}
						if !bytes.Equal(got[i], expected) {
							t.Fatalf("player%d got%x want%x", i, got[i], expected)
						}
					}
					if rv != 0 {
						t.Fatalf("return %x", rv)
					}
					rows = append(rows, visibilityEffectsRow{name, rv, got})
				})
			}
		}
	}
}
