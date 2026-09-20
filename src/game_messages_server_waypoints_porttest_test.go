//go:build porttest

package opennox

import (
	"encoding/binary"
	"math"
	"reflect"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func TestGameMessageServerWaypoints(t *testing.T) {
	var direct, messages []legacy.PortTestRoamSpec
	for index := uint32(0); index < 3; index++ {
		for _, present := range []bool{false, true} {
			for _, flags := range []uint32{0, 1, 2, 3, 4, 0x100} {
				for _, xy := range [][2]uint16{{0, 0}, {1, 65535}, {65535, 1}, {32767, 32768}, {520, 530}} {
					makeCase := func(dispatch bool) legacy.PortTestRoamSpec {
						s := controlsBase(24)
						p := s.Callbacks.Shop
						a := p.TemporaryUpdates.World.Objectives.Attack
						a.UpdateWords = map[int]uint32{180: index | index<<8}
						a.Controls.X = int32(math.Float32bits(float32(xy[0])))
						a.Controls.Y = int32(math.Float32bits(float32(xy[1])))
						p.TemporaryUpdates.World.Objectives.PlayerDataWords[0] = map[int]uint32{3680: flags}
						if present {
							a.UpdateRefs = map[int]int{168 + 4*int(index): 4}
						}
						if dispatch {
							data := []byte{64, 0x5a, 0xa5, 0, 0, 0, 0}
							binary.LittleEndian.PutUint16(data[3:], xy[0])
							binary.LittleEndian.PutUint16(data[5:], xy[1])
							a.Controls.GameMessage = data
							a.Controls.GameMessageLength = 7
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
			t.Fatalf("waypoint dispatch state differs from qualified coordinate operation: case %d", i)
		}
	}
	gameMessageCapture(t, "game-server-waypoints", got)
}
