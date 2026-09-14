//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"math/bits"
	"testing"
)

func TestMapPopulationMarkerRemoval(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for mask := 0; mask < 16; mask++ {
		s := populationCacheBase()
		s.Globals["dword_5d4594_2487656"] = roomValue(6)
		s.Globals["blob2487660"] = roomValue(7)
		s.Globals["blob2487664"] = roomValue(8)
		s.Globals["blob2487668"] = roomValue(9)
		s.Records[0].Refs[80] = roomArg(6)
		s.Records[5].Words[76] = 1
		s.Records[5].Refs[148] = roomArg(2)
		types := []int{}
		for dir := 0; dir < 4; dir++ {
			if mask&(1<<dir) != 0 {
				types = append(types, 6+dir)
			}
		}
		types = append(types, 1)
		s.Globals["dword_5d4594_1599540"] = roomArg(9)
		for i, typ := range types {
			obj := legacy.PortTestPaintObject{Slot: 100 + i, Type: typ, Words: map[int]uint32{56: math.Float32bits(46), 60: math.Float32bits(46)}, Refs: map[int]legacy.PortTestMapRoomArg{}}
			node := roomRecord(12)
			node.Refs[0] = roomArg(100 + i)
			if i+1 < len(types) {
				obj.Refs[444] = roomArg(101 + i)
				node.Refs[4] = roomArg(10 + i)
			}
			if i > 0 {
				obj.Refs[448] = roomArg(99 + i)
				node.Refs[8] = roomArg(8 + i)
			}
			s.Objects = append(s.Objects, obj)
			s.Records = append(s.Records, node)
		}
		s.Actions = []legacy.PortTestPaintAction{paintAction(32, roomArg(1))}
		cases = append(cases, s)
	}
	out := populationCapture(t, "marker-removal", cases)
	for i, r := range out {
		step := r.Steps[0]
		removed := bits.OnesCount(uint(i))
		if step.Return != 1 || step.Objects[0] != 1 {
			t.Fatalf("case %d marker removal admission/count", i)
		}
		deadObjects, deadNodes := 0, 0
		for _, rec := range step.Records {
			if rec.Alive {
				continue
			}
			if rec.Kind == "object" {
				deadObjects++
			}
			if rec.Kind == "input" {
				deadNodes++
			}
		}
		if deadObjects != removed || deadNodes != removed {
			t.Fatalf("case %d removed objects/nodes %d/%d want %d", i, deadObjects, deadNodes, removed)
		}
		if step.Globals["dword_5d4594_1599540"] != step.Slots[9+removed] {
			t.Fatalf("case %d surviving cache head", i)
		}
	}
}
