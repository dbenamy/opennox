//go:build porttest

package opennox

import (
	"fmt"
	"math"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

func prefabPaintingRun(cases []legacy.PortTestPaintSpec) []legacy.PortTestPaintResult {
	return legacy.PortTestPrefabPainting(cases, func(core *server.Server) (legacy.Server, func()) {
		old := noxServer
		wrapped := &Server{Server: core}
		core.ExtServer = unsafe.Pointer(wrapped)
		noxServer = wrapped
		return wrapped, func() { noxServer = old; core.ExtServer = nil }
	})
}
func TestPrefabRuntimePlacement(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	var names []string
	for _, offset := range []int32{-24, -23, -1, 0, 1, 22, 23, 24, 460} {
		for _, existing := range []bool{false, true} {
			for _, compose := range []uint32{0, 1} {
				for _, variation := range []uint32{0, 3, 4, 255} {
					s := legacy.PortTestPaintSpec{Seed: 12345, Globals: map[string]legacy.PortTestMapRoomArg{}, Records: []legacy.PortTestMapRoomRecord{roomRecord(12), roomRecord(36)}}
					s.Records[0].Refs[0] = roomArg(2)
					s.Records[1].Words[0] = 1 | 1<<8 | variation<<16
					s.Records[1].Words[4] = 0xc0 | 50<<8 | 50<<16 | 79<<24
					s.Globals["walls"] = roomArg(1)
					s.Globals["dword_5d4594_3835368"] = roomValue(int32(compose))
					if existing {
						s.Walls = []legacy.PortTestPaintWall{{X: int(offset+1150) / 23, Y: int(offset+1150) / 23, Words: map[int]uint32{0: 2, 4: 0x10}}}
					}
					s.Actions = []legacy.PortTestPaintAction{paintAction(21, roomValue(offset), roomValue(offset))}
					cases = append(cases, s)
					names = append(names, fmt.Sprintf("wall/off%d/existing%t/compose%d/variation%d", offset, existing, compose, variation))
				}
			}
		}
	}
	wallCount := len(cases)
	for _, offset := range []int32{-24, -1, 0, 1, 23, 460} {
		for _, transparent := range []uint32{0, 1, math.MaxUint32} {
			for _, overlay := range []bool{false, true} {
				s := legacy.PortTestPaintSpec{Seed: 12345, Globals: map[string]legacy.PortTestMapRoomArg{}, Records: []legacy.PortTestMapRoomRecord{roomRecord(24), roomRecord(20), roomRecord(24)}}
				s.Records[0].Refs[0] = roomArg(2)
				s.Records[0].Words[4] = math.Float32bits(2300)
				s.Records[0].Words[8] = math.Float32bits(2323)
				s.Records[1].Words[0] = 1
				s.Records[1].Words[4] = 2
				if overlay {
					s.Records[1].Refs[16] = roomArg(3)
					s.Records[2].Words[0] = 0
					s.Records[2].Words[4] = 1
					s.Records[2].Words[8] = 0
					s.Records[2].Words[12] = 1
				}
				s.Globals["tiles"] = roomArg(1)
				s.Globals["transparentFloor"] = roomValue(int32(transparent))
				// Selection is restored after each placement, even with an overlay.
				s.Globals["tile"] = roomValue(0)
				s.Globals["selection35920"] = roomValue(7)
				s.Actions = []legacy.PortTestPaintAction{paintAction(18, roomValue(offset), roomValue(offset))}
				cases = append(cases, s)
				names = append(names, fmt.Sprintf("tile/off%d/transparent%x/overlay%t", offset, transparent, overlay))
			}
		}
	}
	out := prefabPaintingRun(cases)
	for i, r := range out {
		if !r.Intact || !r.ControlOK || len(r.Steps) != 1 || r.Steps[0].Return != 1 {
			t.Fatalf("%s guard or placement failure", names[i])
		}
		if i < wallCount {
			if len(r.Steps[0].Walls.ByPos) != 1 {
				t.Fatalf("%s did not produce a wall", names[i])
			}
		} else {
			step := r.Steps[0]
			if step.Globals["tile"] != 0 || step.Globals["selection35920"] != 7 {
				t.Fatalf("%s did not restore selection", names[i])
			}
			transparent := (i-wallCount)%6/2 == 1
			overlay := (i-wallCount)%2 == 1
			if transparent {
				if len(step.Cells) != 0 {
					t.Fatalf("%s painted transparent base", names[i])
				}
			} else {
				if len(step.Cells) != 1 {
					t.Fatalf("%s missing painted cell", names[i])
				}
				cell := step.Cells[0].Words
				off := 1
				if cell[6] == 1 {
					off = 6
				}
				if cell[off] != 1 || cell[off+1] != 2 {
					t.Fatalf("%s wrong base tile/variation", names[i])
				}
				if (cell[off+4] != 0) != overlay {
					t.Fatalf("%s overlay allocation differs", names[i])
				}
			}
		}
	}
	spellbookCapture(t, "prefab-runtime-placement", out, "e165d6659adb39a9b0fa159554c9ee1bbe90d52db5663a101a7e29c928e5970d")
}
