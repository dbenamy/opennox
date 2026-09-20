//go:build porttest

package opennox

import (
	"encoding/binary"
	"reflect"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func TestGameMessageServerSecondaryWeapon(t *testing.T) {
	var direct, messages []legacy.PortTestRoamSpec
	for _, netCode := range []uint32{7, 0x80000007} {
		for _, class := range []uint32{8, 0x1000, 0x1000000, 0x2000000} {
			for _, flags := range []uint32{0, 1, 2, 3, 0x100} {
				for _, code := range []uint16{0, 7, 0x8017, 0x8007, 0xffff, 0x8000} {
					for _, previous := range []int{0, 2} {
						makeCase := func(dispatch bool) legacy.PortTestRoamSpec {
							target := 0
							if code == 7 && netCode == 7 || code == 0x8017 {
								target = 3
							}
							s := controlsBase(75)
							p := s.Callbacks.Shop
							a := p.TemporaryUpdates.World.Objectives.Attack
							a.Controls.Target = target
							a.Controls.GameMessageSecondary = &target
							p.Equipment.SecondaryWeapon = previous
							p.Items[0].Class = class
							p.Items[0].Flags = 0
							p.Items[0].Subclass = 4
							for i := 0; i < 3; i++ {
								p.TemporaryUpdates.ItemWords[i][36] = netCode + uint32(i)
								// Marked wire IDs address object extents, not net codes.
								p.TemporaryUpdates.ItemWords[i][40] = uint32(23 + i)
							}
							p.TemporaryUpdates.World.Objectives.ObjectList = []int{3, 4, 5}
							p.TemporaryUpdates.World.Objectives.PlayerDataWords[0] = map[int]uint32{3680: flags}
							if dispatch {
								data := []byte{224, 0, 0}
								binary.LittleEndian.PutUint16(data[1:], code)
								a.Controls.GameMessage = data
								a.Controls.GameMessageLength = 3
							}
							return s
						}
						direct = append(direct, makeCase(false))
						messages = append(messages, makeCase(true))
					}
				}
			}
		}
	}
	want := controlsRun(t, direct)
	got := controlsRun(t, messages)
	gameMessageStateGuards(t, want)
	gameMessageStateGuards(t, got)
	for i := range got {
		if !reflect.DeepEqual(got[i], want[i]) {
			t.Fatalf("secondary weapon dispatch differs from qualified operation: case %d", i)
		}
	}
	interactionCapture(t, "game-server-secondary-weapon", got)
}
