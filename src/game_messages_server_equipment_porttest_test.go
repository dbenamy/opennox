//go:build porttest

package opennox

import (
	"encoding/binary"
	"reflect"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func TestGameMessageServerEquipment(t *testing.T) {
	var direct, messages []legacy.PortTestRoamSpec
	for _, op := range []byte{117, 118} {
		for _, state := range []uint32{0, 1} {
			for _, subclass := range []uint32{4, 8} {
				for _, flags := range []uint32{0, 1, 2, 3, 4, 0x100} {
					for _, code := range []uint16{7, 32767} {
						makeCase := func(dispatch bool) legacy.PortTestRoamSpec {
							operation := 70
							equipped := false
							if op == 118 {
								operation = 71
								equipped = true
							}
							eligible := flags&3 == 0 && !(op == 118 && state == 1 && subclass&8 != 0)
							if eligible {
								equipped = op == 117
							} else {
								operation = 67
							}
							s := controlsBase(operation)
							p := s.Callbacks.Shop
							a := p.TemporaryUpdates.World.Objectives.Attack
							a.Controls.Target = 3
							a.Controls.GameMessageEquipped = &equipped
							a.UpdateWords = map[int]uint32{88: state}
							p.Inventory.Linked = []int{0}
							p.Inventory.Owned = []int{0}
							p.Inventory.WeaponBits = map[uint16]uint32{23: subclass}
							p.Items[0].Class = 0x1000000
							p.Items[0].Flags = 0
							p.Items[0].Subclass = subclass
							p.Equipment.ActiveWeapon = 0
							p.Equipment.WeaponFlags = 0
							if op == 118 {
								p.Items[0].Flags = 0x100
								p.Equipment.ActiveWeapon = 1
								p.Equipment.WeaponFlags = subclass
								a.UpdateRefs = map[int]int{104: 3}
							}
							p.TemporaryUpdates.ItemWords[0][36] = uint32(code)
							p.TemporaryUpdates.World.Objectives.PlayerDataWords[0] = map[int]uint32{3680: flags}
							if dispatch {
								data := []byte{op, 0, 0}
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
			t.Fatalf("equipment dispatch differs from qualified direct action: case %d", i)
		}
	}
	gameMessageCapture(t, "game-server-equipment", got)
}
