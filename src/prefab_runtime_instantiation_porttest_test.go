//go:build porttest

package opennox

import (
	"math"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func TestPrefabRuntimeInstantiation(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for _, width := range []float32{0, 32.526913, 130.10765} {
		for _, position := range []float32{0, -100, -10000, 10000} {
			for _, wall := range []bool{false, true} {
				s := legacy.PortTestPaintSpec{Seed: 12345, Globals: map[string]legacy.PortTestMapRoomArg{}, Records: []legacy.PortTestMapRoomRecord{roomRecord(76), roomRecord(8), roomRecord(12), roomRecord(36)}}
				s.Globals["metadata"] = roomArg(1)
				s.Globals["count"] = roomValue(1)
				s.Records[0].Words[64] = math.Float32bits(width)
				s.Records[0].Words[68] = math.Float32bits(width * 2)
				s.Records[1].Words[0] = math.Float32bits(position)
				s.Records[1].Words[4] = math.Float32bits(position)
				if wall {
					s.Globals["walls"] = roomArg(3)
					s.Records[2].Refs[0] = roomArg(4)
					s.Records[3].Words[0] = 1
					s.Records[3].Words[4] = 50<<8 | 50<<16
				}
				s.Actions = []legacy.PortTestPaintAction{paintAction(15, roomArg(2))}
				cases = append(cases, s)
			}
		}
	}
	nilCase := legacy.PortTestPaintSpec{Globals: map[string]legacy.PortTestMapRoomArg{}, Actions: []legacy.PortTestPaintAction{paintAction(15)}}
	cases = append(cases, nilCase)
	out := prefabPaintingRun(cases)
	for i, r := range out {
		if !r.Intact || !r.ControlOK {
			t.Fatalf("case%d owner guard", i)
		}
		want := uint32(1)
		if i == len(cases)-1 || i%8 == 7 {
			want = 0
		}
		step := r.Steps[0]
		if step.Return != want || step.Globals["placed"] != want || step.Globals["instance"] != want {
			t.Fatalf("case%d returned%d placed%d instance%d want%d", i, step.Return, step.Globals["placed"], step.Globals["instance"], want)
		}
		if want != 0 && i%2 == 1 && len(step.Walls.ByPos) != 1 {
			t.Fatalf("case%d missing instantiated wall", i)
		}
	}
	spellbookCapture(t, "prefab-runtime-instantiation", out, "a23c80e5fc62fa92f886f9079e041c1a5dea413dd4ef05154e175b102169af8c")
}
