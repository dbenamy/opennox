//go:build porttest

package opennox

import (
	"bytes"
	"crypto/sha256"
	"testing"
)

// Apply the documented C empty-floor-definition prerequisite before this test.
func TestFloorAssetsFailureContracts(t *testing.T) {
	o := newFloorAssetsOwner(t)
	type record struct {
		Edge                                  bool
		Phase, InitialCount, Return, Consumed int
		Metadata                              [32]byte
		Scratch                               [32]byte
	}
	var rows []record
	for _, edge := range []bool{false, true} {
		o.reset()
		name := "Missing"
		op := 4
		input := floorAssetsFloor(name, [3]byte{0, 0, 0}, [3]byte{1, 1, 1}, 0)
		if edge {
			op = 5
			input = floorAssetsEdge(name, [3]byte{1, 1, 0}, 0, 0, 0x454e4420)
		}
		beforeDefs, beforeEdges := o.metadata()
		ret, consumed := o.invoke(t, op, input)
		defs, edges := o.metadata()
		o.assertOnlyData(t, false, -1)
		if ret != 0 || consumed != 5+len(name) || *o.count != 0 || *o.edgeCount != 0 || !bytes.Equal(beforeDefs, defs) || !bytes.Equal(beforeEdges, edges) || !bytes.Equal(o.scratch, bytes.Repeat([]byte{0xa5}, len(o.scratch))) {
			t.Fatal("empty definition table did not fail without side effects")
		}
		rows = append(rows, record{edge, 0, 0, ret, consumed, sha256.Sum256(append(defs, edges...)), sha256.Sum256(o.scratch)})
	}
	for _, count := range []int{0, 63} {
		o.reset()
		*o.edgeCount = uint32(count)
		for phase, end := range []uint32{0x12345678, 0x454e4420} {
			input := floorAssetsEdge("RetryEdge", [3]byte{1, 1, 0}, 2, 0, end)
			ret, consumed := o.invoke(t, 1, input)
			want := phase
			if ret != want || consumed != len(input) || *o.edgeCount != uint32(count+phase) || *o.count != 0 {
				t.Fatal("edge END failure/retry cursor or count")
			}
			o.assertOnlyData(t, false, -1)
			defs, edges := o.metadata()
			for i := 0; i < 64; i++ {
				if i != count && !bytes.Equal(edges[60*i:][:60], make([]byte, 60)) {
					t.Fatal("edge retry changed another row")
				}
			}
			rows = append(rows, record{true, phase + 1, count, ret, consumed, sha256.Sum256(append(defs, edges...)), sha256.Sum256(o.scratch)})
		}
	}
	floorAssetsCapture(t, "failures", rows, "a3c38136006b82d590e109a0bfc9fe493788f5910042176934f0214cdc6a8528")
}
