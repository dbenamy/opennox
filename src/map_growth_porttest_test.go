//go:build porttest

package opennox

import (
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

func growthBase() legacy.PortTestPaintSpec {
	s := legacy.PortTestPaintSpec{Seed: 12345, Globals: map[string]legacy.PortTestMapRoomArg{"gameFlags": roomValue(1 << 22), "growthInitGrid": roomValue(1)}}
	s.Records = []legacy.PortTestMapRoomRecord{roomRecord(1116), roomRecord(376), roomRecord(376)}
	s.Records[0].Words = map[int]uint32{4: 5, 8: 1, 12: 3, 16: 1, 20: 3, 24: 40, 28: 20, 32: 5, 36: 1, 40: 20, 48: 100, 52: 100, 56: 1, 68: 32, 72: 3}
	return s
}
func growthRun(cases []legacy.PortTestPaintSpec) []legacy.PortTestPaintResult {
	return legacy.PortTestMapGrowth(cases, func(core *server.Server) (legacy.Server, func()) {
		old := noxServer
		wrapped := &Server{Server: core}
		core.ExtServer = unsafe.Pointer(wrapped)
		noxServer = wrapped
		return wrapped, func() { noxServer = old; core.ExtServer = nil }
	})
}
func TestMapGrowthBootstrap(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for _, kind := range []uint32{0, 1} {
		for seed := 0; seed < 4; seed++ {
			s := growthBase()
			s.Seed = seed
			s.Records[0].Words[0] = kind
			s.Actions = []legacy.PortTestPaintAction{paintAction(1, roomArg(1))}
			cases = append(cases, s)
		}
	}
	out := growthRun(cases)
	for i, r := range out {
		if !r.Intact || !r.ControlOK || r.Steps[0].Return == 0 {
			t.Fatalf("growth bootstrap %d", i)
		}
		rooms := 0
		for _, record := range r.Steps[0].Records {
			if record.Kind == "growthAllocation" && record.Alive && len(record.Words) == 94 {
				rooms++
			}
		}
		want := 1
		if i >= 4 {
			want = 48
		}
		if rooms != want {
			t.Fatalf("growth bootstrap %d rooms %d want %d", i, rooms, want)
		}
	}
}
