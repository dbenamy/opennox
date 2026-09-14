//go:build porttest

package opennox

import (
	"os"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/internal/binfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
	"github.com/opennox/opennox/v1/server"
)

func orchestrationRun(cases []legacy.PortTestPaintSpec) []legacy.PortTestPaintResult {
	return legacy.PortTestMapOrchestration(cases, func(core *server.Server) (legacy.Server, func()) {
		old := noxServer
		wrapped := &Server{Server: core}
		core.ExtServer = unsafe.Pointer(wrapped)
		noxServer = wrapped
		return wrapped, func() { noxServer = old; core.ExtServer = nil }
	})
}
func TestMapOrchestrationStepProbe(t *testing.T) {
	handles.Init()
	defer handles.Release()
	t.Chdir(t.TempDir())
	if err := os.Mkdir("mapgen", 0700); err != nil {
		t.Fatal(err)
	}
	f, err := binfile.BinfileOpen("mapgen/integration.thm", binfile.WriteOnly)
	if err != nil {
		t.Fatal(err)
	}
	if err = f.SetKey(1); err != nil {
		_ = f.Close()
		t.Fatal(err)
	}
	text := "ALGORITHM_DATA mapSize 200 midRoomSize 5 roomVariance 0 recursionLimit 2 seed 17 END DECOR ROOM base WALL_FLOOR PaintWall PaintTile END DECOR HALL hall WALL_FLOOR PaintWall PaintTile END "
	if _, err = f.Write([]byte(text)); err != nil {
		_ = f.Close()
		t.Fatal(err)
	}
	if err = f.Close(); err != nil {
		t.Fatal(err)
	}
	s := growthBase()
	s.Globals["growthInitGrid"] = roomValue(0)
	paintString(&s.Records[1], 0, "integration")
	s.Actions = []legacy.PortTestPaintAction{paintAction(0, roomArg(2)), paintAction(1), paintAction(2)}
	results := orchestrationRun([]legacy.PortTestPaintSpec{s})
	r := results[0]
	if !r.Intact || !r.ControlOK {
		t.Fatal("generator guards/control")
	}
	step := r.Steps[2]
	if step.Return != 1 {
		t.Fatalf("generator return %d", step.Return)
	}
	released, starts := 0, 0
	for _, rec := range step.Records {
		if rec.Kind == "object" && rec.Alive {
			starts++
		}
		if rec.Kind == "orchestrationReleased" && len(rec.Words) == 95 {
			released++
		}
	}
	if starts != 1 {
		t.Fatalf("player start object count %d", starts)
	}
	if released == 0 {
		t.Fatal("no generated room release observed")
	}
	t.Logf("generator success, %d released rooms, %d records", released, len(step.Records))
}
