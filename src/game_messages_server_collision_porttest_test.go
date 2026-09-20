//go:build porttest

package opennox

import (
	"encoding/binary"
	"reflect"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func TestGameMessageServerCollision(t *testing.T) {
	var direct, messages []legacy.PortTestRoamSpec
	for _, flags := range []uint32{0, 1, 2, 3, 0x100} {
		for _, code := range []uint16{7, 0x8017, 0xffff} {
			for _, busy := range []int{0, 280, 284} {
				for _, callback := range []bool{false, true} {
					call := flags&3 == 0 && code != 0xffff && busy == 0 && callback
					makeCase := func(dispatch bool) legacy.PortTestRoamSpec {
						op := 67
						if call {
							op = 78
						}
						s := controlsBase(op)
						p := s.Callbacks.Shop
						o := p.TemporaryUpdates.World.Objectives
						a := o.Attack
						a.Controls.Target = 3
						a.Controls.GameMessageCollision = &call
						a.Controls.GameMessageNullCollision = !callback
						a.UpdateWords = map[int]uint32{280: 0, 284: 0}
						if busy != 0 {
							a.UpdateWords[busy] = 1
						}
						o.PlayerDataWords[0] = map[int]uint32{3680: flags}
						p.Items[0].Flags = 0
						p.TemporaryUpdates.ItemWords[0][36] = 7
						p.TemporaryUpdates.ItemWords[0][40] = 23
						o.ObjectList = []int{3, 4, 5}
						if dispatch {
							data := []byte{123, 0, 0}
							binary.LittleEndian.PutUint16(data[1:], code)
							a.Controls.GameMessage, a.Controls.GameMessageLength = data, 3
						}
						return s
					}
					direct = append(direct, makeCase(false))
					messages = append(messages, makeCase(true))
				}
			}
		}
	}
	want, got := controlsRun(t, direct), controlsRun(t, messages)
	gameMessageStateGuards(t, want)
	gameMessageStateGuards(t, got)
	for i := range got {
		if !reflect.DeepEqual(got[i], want[i]) {
			t.Fatalf("collision request differs from callback: case %d", i)
		}
	}
	interactionCapture(t, "game-server-collision", got)
}
