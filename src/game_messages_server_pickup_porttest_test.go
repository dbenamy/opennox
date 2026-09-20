//go:build porttest

package opennox

import (
	"encoding/binary"
	"reflect"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func TestGameMessageServerPickup(t *testing.T) {
	key := "pickup.c:ObjectEquipClassFail"
	literal := serverConfigOwnBytes(t, 0x587000, 215732, len(key)+1)
	copy(literal, key)
	literal[len(key)] = 0
	type input struct {
		flags        uint32
		code         uint16
		weight, held byte
		capacity     uint16
		class        byte
		special      bool
		typ          uint16
	}
	var inputs []input
	for _, flags := range []uint32{0, 1, 2, 3, 4, 0x100} {
		for _, code := range []uint16{7, 32767} {
			for _, w := range [][3]uint16{{0, 0, 0}, {1, 0, 0}, {1, 0, 1}, {1, 1, 1}, {255, 255, 510}, {255, 255, 509}, {255, 1, 256}, {255, 1, 65535}} {
				inputs = append(inputs, input{flags: flags, code: code, weight: byte(w[0]), held: byte(w[1]), capacity: w[2], typ: 23})
			}
		}
	}
	for class := byte(0); class < 3; class++ {
		for _, typ := range []uint16{0x6a, 0x6b, 0x6d, 0x70} {
			for _, flags := range []uint32{0, 1} {
				for _, capacity := range []uint16{0, 1} {
					inputs = append(inputs, input{flags: flags, code: 7, weight: 1, capacity: capacity, class: class, special: true, typ: typ})
				}
			}
		}
	}
	var direct, messages []legacy.PortTestRoamSpec
	for _, in := range inputs {
		makeCase := func(dispatch bool) legacy.PortTestRoamSpec {
			operation, calls := 72, 1
			if in.flags&3 != 0 {
				operation, calls = 67, 0
			} else if uint32(in.weight)+uint32(in.held) > uint32(in.capacity) {
				operation, calls = 74, 0
			} else if in.special && ((in.typ == 0x6a && in.class != 1) || ((in.typ == 0x6b || in.typ == 0x6d) && in.class != 0)) {
				operation, calls = 73, 0
			}
			s := controlsBase(operation)
			p := s.Callbacks.Shop
			a := p.TemporaryUpdates.World.Objectives.Attack
			a.Controls.Target = 3
			a.Controls.GameMessagePickupCalls = &calls
			p.Resources.PlayerClass = in.class
			p.Inventory.Linked = []int{1}
			p.Inventory.Owned = []int{1}
			p.Inventory.Weights = []byte{in.weight, in.held, 0}
			p.Inventory.Carry = in.capacity
			p.Items[0].Class = 8
			if in.special {
				p.Items[0].Class = 0x110000
			}
			p.Items[0].Flags = 0
			p.Items[0].Type = in.typ
			p.TemporaryUpdates.ItemWords[0][36] = uint32(in.code)
			p.TemporaryUpdates.World.Objectives.ObjectList = []int{3, 4, 5}
			p.TemporaryUpdates.World.Objectives.PlayerDataWords[0] = map[int]uint32{3680: in.flags}
			if dispatch {
				data := []byte{115, 0, 0}
				binary.LittleEndian.PutUint16(data[1:], in.code)
				a.Controls.GameMessage = data
				a.Controls.GameMessageLength = 3
			}
			return s
		}
		direct = append(direct, makeCase(false))
		messages = append(messages, makeCase(true))
	}
	want := controlsRun(t, direct)
	got := controlsRun(t, messages)
	gameMessageStateGuards(t, want)
	gameMessageStateGuards(t, got)
	for i := range got {
		if !reflect.DeepEqual(got[i], want[i]) {
			t.Fatalf("pickup dispatch differs from independent weight/class decision: case %d input %+v", i, inputs[i])
		}
	}
	gameMessageCapture(t, "game-server-pickup", got)
}
