//go:build porttest

package opennox

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
)

// Released grid rows must keep their capture identities even when later trace
// allocations reuse their storage. Compare complete output under varied layouts.
func TestMapOrchestrationAllocationLayouts(t *testing.T) {
	handles.Init()
	defer handles.Release()
	t.Chdir(t.TempDir())
	if err := os.Mkdir("mapgen", 0700); err != nil {
		t.Fatal(err)
	}
	orchestrationTheme(t, "layout", "ALGORITHM_DATA mapSize 400 midRoomSize 5 roomVariance 0 recursionLimit 2 hallBranchRate 0 seed 0 END DECOR ROOM base WALL_FLOOR PaintWall PaintTile END DECOR HALL hall WALL_FLOOR PaintWall PaintTile END ")
	orchestrationTheme(t, "ring3", "ALGORITHM_DATA skeleton HALL_RING mapSize 1050 recursionLimit 1 seed 3 END DECOR ROOM base WALL_FLOOR PaintWall PaintTile END DECOR HALL hall WALL_FLOOR PaintWall PaintTile END ")
	orchestrationTheme(t, "ring7", "ALGORITHM_DATA skeleton HALL_RING mapSize 1050 recursionLimit 2 seed 7 END DECOR ROOM base WALL_FLOOR PaintWall PaintTile END DECOR HALL hall WALL_FLOOR PaintWall PaintTile END ")
	cases := []legacy.PortTestPaintSpec{orchestrationCase("layout"), orchestrationCase("ring3"), orchestrationCase("ring7")}
	capture := func() []byte {
		out := orchestrationRun(cases)
		for _, r := range out {
			if !r.Intact || !r.ControlOK || r.Steps[2].Return != 1 {
				t.Fatal("allocation-layout generator state")
			}
		}
		data, err := json.Marshal(out)
		if err != nil {
			t.Fatal(err)
		}
		return data
	}
	expected := capture()
	sizes := []int{100, 256, 500, 512, 768, 1300}
	for layout := 0; layout < 8; layout++ {
		func() {
			var release []func()
			defer func() {
				for _, free := range release {
					if free != nil {
						free()
					}
				}
			}()
			for i := 0; i < 16+layout*7; i++ {
				_, free := alloc.Make([]byte{}, sizes[(i+layout)%len(sizes)])
				release = append(release, free)
			}
			for i := layout % 2; i < len(release); i += 2 {
				release[i]()
				release[i] = nil
			}
			if got := capture(); !bytes.Equal(expected, got) {
				t.Fatalf("complete generator capture depends on allocation layout %d", layout)
			}
		}()
	}
}
