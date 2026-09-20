//go:build porttest

package opennox

import (
	"encoding/binary"
	"reflect"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

// Missing target IDs and disabled player states must not mutate gameplay owners.
// Positive item, spell and trade actions have separate contracts.
func TestGameMessageServerMissingTargets(t *testing.T) {
	var direct, messages []legacy.PortTestRoamSpec
	for _, op := range []byte{114, 115, 116, 117, 118, 120, 123} {
		for _, flags := range []uint32{0, 1, 2, 3, 0x100, 0xffffffff} {
			for _, code := range []uint16{0x7fff, 0xffff} {
				makeCase := func(dispatch bool) legacy.PortTestRoamSpec {
					s := controlsBase(67)
					p := s.Callbacks.Shop
					a := p.TemporaryUpdates.World.Objectives.Attack
					p.TemporaryUpdates.World.Objectives.PlayerDataWords[0] = map[int]uint32{3680: flags}
					if dispatch {
						n := 3
						if op == 114 {
							n = 7
						}
						if op == 120 {
							n = 4
						}
						data := make([]byte, n)
						data[0] = op
						binary.LittleEndian.PutUint16(data[1:], code)
						if op == 114 {
							binary.LittleEndian.PutUint16(data[3:], 32768)
							binary.LittleEndian.PutUint16(data[5:], 65535)
						}
						if op == 120 {
							data[3] = 7
						}
						a.Controls.GameMessage = data
						a.Controls.GameMessageLength = n
					}
					return s
				}
				direct = append(direct, makeCase(false))
				messages = append(messages, makeCase(true))
			}
		}
	}
	want := controlsRun(t, direct)
	got := controlsRun(t, messages)
	gameMessageStateGuards(t, want)
	gameMessageStateGuards(t, got)
	for i := range got {
		if !reflect.DeepEqual(got[i], want[i]) {
			t.Fatalf("missing-target message changed gameplay state: case %d", i)
		}
	}
	gameMessageCapture(t, "game-server-missing-targets", got)
}
