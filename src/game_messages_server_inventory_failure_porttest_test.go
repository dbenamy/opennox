//go:build porttest

package opennox

import (
	"encoding/binary"
	"reflect"
	"testing"

	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/legacy"
)

func TestGameMessageServerInventoryFailure(t *testing.T) {
	var direct, messages []legacy.PortTestRoamSpec
	for _, flags := range []uint32{0, 1, 2, 3, 0x100} {
		for _, code := range []uint16{0, 7, 0x8007, 32767, 65535} {
			for _, dropOK := range []bool{false, true} {
				for _, pos := range []types.Pointf{{X: 512, Y: 512}, {X: 520.5, Y: 530.25}, {X: 600, Y: 700}} {
					makeCase := func(dispatch bool) legacy.PortTestRoamSpec {
						op, callKind := 67, uint32(0)
						// This message looks up the full item code without extent conversion.
						if code == 7 {
							op, callKind = 76, 2
						}
						s := controlsBase(op)
						p := s.Callbacks.Shop
						a := p.TemporaryUpdates.World.Objectives.Attack
						a.Controls.Target = 3
						a.Controls.GameMessageCheckCalls = true
						a.Controls.GameMessageCallKind = callKind
						p.Inventory.Linked = []int{0}
						p.Inventory.Owned = []int{0}
						p.Inventory.Position = pos
						p.Inventory.DropResult = dropOK
						p.Items[0].Class = 8
						p.Items[0].Flags = 0
						p.Items[0].Subclass = 0
						p.TemporaryUpdates.ItemWords[0][36] = 7
						p.TemporaryUpdates.ItemWords[0][40] = 7
						p.TemporaryUpdates.World.Objectives.ObjectList = []int{3, 4, 5}
						p.TemporaryUpdates.World.Objectives.PlayerDataWords[0] = map[int]uint32{3680: flags}
						if dispatch {
							data := []byte{241, 0, 0}
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
	want := controlsRun(t, direct)
	got := controlsRun(t, messages)
	gameMessageStateGuards(t, want)
	gameMessageStateGuards(t, got)
	for i := range got {
		if !reflect.DeepEqual(got[i], want[i]) {
			t.Fatalf("inventory-failure dispatch differs from drop and notification: case %d", i)
		}
	}
	interactionCapture(t, "game-server-inventory-failure", got)
}
