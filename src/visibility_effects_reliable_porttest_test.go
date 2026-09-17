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

func TestVisibilityEffectsReliablePackets(t *testing.T) {
	s, _, _ := visibilityEffectsPlayers(t)
	reset, snapshot, free := legacy.PortTestVisibilityReliableOwner()
	t.Cleanup(free)
	u := newCreatureXferObject(t, s, "Monster")
	u.NetCode = 0x1234
	var rows []struct {
		Name    string
		Returns []uint32
		Packets []legacy.PortTestShopPacketResult
	}
	defer func() {
		spellbookCapture(t, "visibility-effects-reliable", rows, "b57afb84055aa2c143781abecf3842677096024be5f1a3c5532424cb9cb623c1")
	}()
	for _, id := range []int32{0, 1, 32767, 32768, 65535, 65536, -1, 0x12345678} {
		for _, xy := range [][2]float32{{0, 0}, {100.75, -200.25}, {65537.5, -65537.5}} {
			name := fmt.Sprintf("id%d/xy%v", id, xy)
			t.Run(name, func(t *testing.T) {
				reset()
				p := types.Pointf{xy[0], xy[1]}
				args := [5]int32{id, id ^ 0x55, id ^ 0xAAAA, id ^ 0x1234}
				ops := []int{8, 9, 26, 27}
				var returns []uint32
				var expected [][]byte
				start := []byte{126}
				start = binary.LittleEndian.AppendUint16(start, uint16(int64(xy[0])))
				start = binary.LittleEndian.AppendUint16(start, uint16(int64(xy[1])))
				start = binary.LittleEndian.AppendUint16(start, uint16(id))
				start = binary.LittleEndian.AppendUint16(start, uint16(args[2]))
				start = append(start, byte(args[1]))
				start = binary.LittleEndian.AppendUint16(start, uint16(args[3]))
				expected = append(expected, start, []byte{127, byte(id), byte(id >> 8)}, []byte{50, 0x34, 0x12}, []byte{51, 0x34, 0x12})
				for _, op := range ops {
					a := args
					if op == 26 || op == 27 {
						a[0] = 7
					}
					rv := legacy.PortTestVisibilityEffects(op, nil, u, &p, nil, a, nil, "")
					if rv != 1 {
						t.Fatalf("op%d return%d", op, rv)
					}
					returns = append(returns, rv)
				}
				got := snapshot()
				if len(got) != 4 {
					t.Fatalf("packet count%d", len(got))
				}
				for i, p := range got {
					j := 3 - i
					recipient := byte(255)
					if j >= 2 {
						recipient = 7
					}
					if p.Recipient != recipient || p.Ordered != 0 || p.A4 != 0 || p.A5 != 1 || !bytes.Equal(p.Data, expected[j]) {
						t.Fatalf("packet%d got%+v want%x", i, p, expected[j])
					}
					for _, seq := range p.Sequence {
						if seq != 0 {
							t.Fatal("unordered sequence changed")
						}
					}
				}
				rows = append(rows, struct {
					Name    string
					Returns []uint32
					Packets []legacy.PortTestShopPacketResult
				}{name, returns, got})
			})
		}
	}
}
