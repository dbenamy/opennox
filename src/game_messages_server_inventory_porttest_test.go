//go:build porttest

package opennox

import (
	"encoding/binary"
	"math"
	"reflect"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func TestGameMessageServerDropAndUse(t *testing.T) {
	var direct, messages []legacy.PortTestRoamSpec
	for _, op := range []byte{114, 116} {
		for _, flags := range []uint32{0, 1, 2, 3, 4, 0x100} {
			for _, code := range []uint16{7, 32767} {
				coords := [][2]uint16{{512, 512}}
				if op == 114 {
					coords = append(coords, [2]uint16{520, 530}, [2]uint16{530, 550})
				}
				for _, xy := range coords {
					makeCase := func(dispatch bool) legacy.PortTestRoamSpec {
						operation := 68
						kind := uint32(2)
						if op == 116 {
							operation = 69
							kind = 1
						}
						if flags&3 != 0 {
							operation = 67
							kind = 0
						}
						s := controlsBase(operation)
						p := s.Callbacks.Shop
						a := p.TemporaryUpdates.World.Objectives.Attack
						a.Controls.Target = 3
						a.Controls.GameMessageCheckCalls = true
						a.Controls.GameMessageCallKind = kind
						a.Controls.X = int32(math.Float32bits(float32(xy[0])))
						a.Controls.Y = int32(math.Float32bits(float32(xy[1])))
						p.Inventory.Linked = []int{0}
						p.Inventory.Owned = []int{0}
						p.Inventory.UseDelete = op == 116
						p.Items[0].Class = 8
						p.Items[0].Flags = 0
						p.Items[0].Subclass = 0
						p.TemporaryUpdates.ItemWords[0][36] = uint32(code)
						p.TemporaryUpdates.World.Objectives.PlayerDataWords[0] = map[int]uint32{3680: flags}
						if dispatch {
							n := 3
							if op == 114 {
								n = 7
							}
							data := make([]byte, n)
							data[0] = op
							binary.LittleEndian.PutUint16(data[1:], code)
							if op == 114 {
								binary.LittleEndian.PutUint16(data[3:], xy[0])
								binary.LittleEndian.PutUint16(data[5:], xy[1])
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
	}
	want := controlsRun(t, direct)
	got := controlsRun(t, messages)
	gameMessageStateGuards(t, want)
	gameMessageStateGuards(t, got)
	for i := range got {
		if !reflect.DeepEqual(got[i], want[i]) {
			t.Fatalf("drop/use dispatch differs from qualified direct action: case %d", i)
		}
	}
	interactionCapture(t, "game-server-drop-use", got)
}
