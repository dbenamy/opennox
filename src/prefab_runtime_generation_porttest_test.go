//go:build porttest

package opennox

import (
	"fmt"
	"math"
	"testing"
	"unsafe"

	"github.com/opennox/libs/types"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func TestPrefabRuntimeWaypointGeneration(t *testing.T) {
	s := newObjectXferOwner(t)
	defer noxflags.PortTestGameFlags(0)()
	words, restore := legacy.PortTestPrefabRuntimeGlobals()
	defer restore()
	points, free := alloc.Make([]types.Pointf{}, 3)
	defer free()
	ptr := func(i int) uint32 { return uint32(uintptr(unsafe.Pointer(&points[i]))) }
	type wpRecord struct {
		ID, X, Y, Flags, Next, Prev uint32
		Links                       [][2]uint32
	}
	type row struct {
		Name      string
		FlagWord  uint32
		Returns   []uint64
		Waypoints []wpRecord
	}
	var rows []row
	transform := func(p types.Pointf) types.Pointf {
		x := float32((float64(p.Y)+float64(p.X))*0.70710677 + 2957)
		y := float32((float64(p.Y)-float64(p.X))*0.70710677 + 2956)
		if x <= 80.5 {
			x = 82.5
		}
		if y <= 80.5 {
			y = 81.5
		}
		if x >= 5853.5 {
			x = 5851.5
		}
		if y >= 5853.5 {
			y = 5852.5
		}
		return types.Ptf(x, y)
	}
	for _, kind := range []uint32{0, 1, 2, 127, 128, 255, 0x12345678} {
		for _, auto := range []uint32{0, 1, 2, math.MaxUint32} {
			for _, base := range []float32{0, -10000, 10000} {
				points[0] = types.Ptf(base, base)
				points[1] = types.Ptf(base+100, base+200)
				points[2] = types.Ptf(1234, -1234)
				*words["autoConnect"] = auto
				*words["waypointKind"] = 0xaabbccdd
				if legacy.PortTestPrefabCall(36, [6]uint32{kind}) != 1 || *words["waypointKind"] != 0xaabbcc00|uint32(byte(kind)) {
					t.Fatal("flag setter width")
				}
				if legacy.PortTestPrefabCall(37, [6]uint32{}) != 0 || *words["lastWaypoint"] != 0 {
					t.Fatal("nil creation changed last waypoint")
				}
				var created []*server.Waypoint
				r := row{Name: fmt.Sprintf("kind%x/auto%x/base%g", kind, auto, base), FlagWord: *words["waypointKind"]}
				for i := 0; i < 2; i++ {
					raw := uint32(legacy.PortTestPrefabCall(37, [6]uint32{ptr(i)}))
					wp := s.WPs.List
					if wp == nil || raw != uint32(uintptr(unsafe.Pointer(wp))) || *words["lastWaypoint"] != raw || wp.Index != uint32(i+1) || wp.PosVec != transform(points[i]) {
						t.Fatal("creation, transform or last pointer")
					}
					created = append(created, wp)
				}
				linked := byte(0)
				if auto == 1 {
					linked = 1
				}
				for i, w := range created {
					if w.PointsCnt != linked {
						t.Fatal("auto-connect admission")
					}
					if linked != 0 && (w.Points[0].Waypoint != created[1-i] || w.Points[0].Ind != byte(kind)) {
						t.Fatal("auto-connect link")
					}
				}
				// A reset only breaks the creation chain; existing waypoints remain owned.
				legacy.PortTestPrefabCall(35, [6]uint32{})
				if *words["lastWaypoint"] != 0 || s.WPs.List != created[1] {
					t.Fatal("reset altered world")
				}
				raw := uint32(legacy.PortTestPrefabCall(37, [6]uint32{ptr(2)}))
				third := s.WPs.List
				if raw == 0 || third.Index != 3 || third.PointsCnt != 0 {
					t.Fatal("post-reset creation")
				}
				created = append(created, third)
				ids := map[uint32]uint32{0: 0}
				for _, w := range created {
					ids[uint32(uintptr(unsafe.Pointer(w)))] = w.Index
				}
				norm := func(raw uint32) uint32 {
					v, ok := ids[raw]
					if !ok {
						t.Fatalf("unknown waypoint %x", raw)
					}
					return v
				}
				for i := range points {
					got := uint32(legacy.PortTestPrefabCall(38, [6]uint32{ptr(i)}))
					// Coincident clamped points select the most recently prepended waypoint.
					want := created[i].Index
					if i < 2 && created[0].PosVec == created[1].PosVec {
						want = 2
					}
					if norm(got) != want {
						t.Fatalf("find %d returned %d want %d", i, norm(got), want)
					}
					r.Returns = append(r.Returns, uint64(norm(got)))
				}
				if legacy.PortTestPrefabCall(38, [6]uint32{}) != 0 {
					t.Fatal("nil find")
				}
				for _, args := range [][6]uint32{{}, {ptr(0), 0}, {0, ptr(1)}} {
					if legacy.PortTestPrefabCall(39, args) != 0 {
						t.Fatal("nil connect")
					}
				}
				// Explicit connection is directed and uses the current raw kind byte.
				for attempt := 0; attempt < 2; attempt++ {
					ret := legacy.PortTestPrefabCall(39, [6]uint32{ptr(2), ptr(0)})
					want := uint64(1)
					if attempt == 1 && byte(kind) < 128 {
						want = 0
					}
					if ret != want {
						t.Fatalf("connect return %d want %d", ret, want)
					}
					r.Returns = append(r.Returns, ret)
				}
				if third.PointsCnt != 1+byte(bool2int(byte(kind) >= 128)) {
					t.Fatal("duplicate link signedness")
				}
				id := func(w *server.Waypoint) uint32 {
					if w == nil {
						return 0
					}
					return w.Index
				}
				for _, w := range created {
					rec := wpRecord{ID: w.Index, X: math.Float32bits(w.PosVec.X), Y: math.Float32bits(w.PosVec.Y), Flags: w.Flags, Next: id(w.WpNext), Prev: id(w.WpPrev)}
					for i := byte(0); i < w.PointsCnt; i++ {
						rec.Links = append(rec.Links, [2]uint32{id(w.Points[i].Waypoint), uint32(w.Points[i].Ind)})
					}
					r.Waypoints = append(r.Waypoints, rec)
				}
				rows = append(rows, r)
				s.Nox_xxx_waypointDeleteAll_579DD0()
				legacy.PortTestPrefabCall(35, [6]uint32{})
			}
		}
	}
	spellbookCapture(t, "prefab-runtime-waypoint-generation", rows, "399afa48848d51d2932880d3f82caeefbda71a815127c16a170fc6e59a0cb9bb")
}
