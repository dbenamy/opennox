//go:build porttest

package opennox

import (
	"bytes"
	"fmt"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
)

func TestVisibilityEffectsShield(t *testing.T) {
	s, units, configure := visibilityEffectsPlayers(t)
	configure(3)
	u := newCreatureXferObject(t, s, "Monster")
	u.NetCode = 0x1234
	u.PosVec = types.Pointf{100, 100}
	for i := range units {
		units[i].PosVec = u.PosVec
	}
	var rows []visibilityEffectsRow
	defer func() {
		spellbookCapture(t, "visibility-effects-shield", rows, "238ad11972b983a12a44b8064555f2e2f772f361c2eda61b90926c1fe060acca")
	}()
	for direction := 0; direction < 256; direction++ {
		t.Run(fmt.Sprintf("direction%d", direction), func(t *testing.T) {
			s.NetList.ResetAll()
			u.Direction1 = server.Dir16(direction)
			rv := legacy.PortTestVisibilityEffects(7, u, nil, nil, nil, [5]int32{}, nil, "")
			got := visibilityEffectsPackets(s)
			for _, slot := range []int{1, 7, 31} {
				if len(got[slot]) != 4 || !bytes.Equal(got[slot][:3], []byte{128, 0x34, 0x12}) {
					t.Fatalf("packet%x", got[slot])
				}
				if direction%32 == 0 {
					want := []byte{5, 8, 7, 6, 3, 0, 1, 2}[direction/32]
					if got[slot][3] != want {
						t.Fatalf("direction byte%d want%d", got[slot][3], want)
					}
				}
			}
			if rv != 0 {
				t.Fatalf("return%x", rv)
			}
			rows = append(rows, visibilityEffectsRow{fmt.Sprintf("direction%d", direction), rv, got})
		})
	}
	for i, delta := range []types.Pointf{{1, 0}, {1, 1}, {0, 1}, {-1, 1}, {-1, 0}, {-1, -1}, {0, -1}, {1, -1}, {0, 0}} {
		t.Run(fmt.Sprintf("vector%d", i), func(t *testing.T) {
			s.NetList.ResetAll()
			p := types.Pointf{u.PosVec.X - delta.X, u.PosVec.Y - delta.Y}
			rv := legacy.PortTestVisibilityEffects(7, u, nil, &p, nil, [5]int32{}, nil, "")
			got := visibilityEffectsPackets(s)
			want := []byte{128, 0x34, 0x12, []byte{5, 8, 7, 6, 3, 0, 1, 2, 5}[i]}
			for _, slot := range []int{1, 7, 31} {
				if !bytes.Equal(got[slot], want) {
					t.Fatalf("packet%x want%x", got[slot], want)
				}
			}
			rows = append(rows, visibilityEffectsRow{fmt.Sprintf("vector%d", i), rv, got})
		})
	}
}
