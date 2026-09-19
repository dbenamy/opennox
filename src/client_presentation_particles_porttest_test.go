//go:build porttest

package opennox

import (
	"image"
	"math"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/legacy"
)

func TestClientPresentationTurnUndead(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t, "UndeadKiller")
	for _, r := range blobdata.PortTestClientPresentationTables() {
		copy(serverConfigOwnBytes(t, r.Base, r.Offset, len(r.Data)), r.Data)
	}
	cache := serverConfigOwnBytes(t, 0x5D4594, 1217508, 4)
	type particle struct {
		Angle    uint16
		Velocity [2]uint32
		Origin   image.Point
		Frame    uint32
	}
	type record struct {
		Position  [2]int16
		Frame     uint32
		Failure   int
		Cached    bool
		Calls     []effectsSpawnCall
		Particles []particle
	}
	var rows []record
	var callback unsafe.Pointer
	for _, pos := range [][2]int16{{0, 0}, {123, -456}, {-32768, 32767}} {
		for _, frame := range []uint32{0, 120, 0xffffffff} {
			for _, fail := range []int{0, 1, 2, 7, 43, 44} {
				for _, cached := range []bool{false, true} {
					c.resetCase(env, pix, 31, frame)
					clear(cache)
					typ := c.Things.IndByID("UndeadKiller")
					if cached {
						*memmap.PtrUint32(0x5D4594, 1217508) = uint32(typ)
					}
					c.FailEvery = fail
					legacy.PortTestPresentationTurnUndead(&pos)
					if len(c.Calls) != 43 || c.srv.Rand.Logic.Index() != 31 || c.srv.Rand.Other.Index() != 32 {
						t.Fatal("particle count or unexpected RNG consumption")
					}
					if memmap.Uint32(0x5D4594, 1217508) != uint32(typ) {
						t.Fatal("particle type cache")
					}
					byRef := make(map[uint32]particle)
					for dr := c.Objs.List1; dr != nil; dr = dr.NextPtr {
						a := uint16(dr.Field_127)
						vx := math.Float32bits(*memmap.PtrFloat32(0x587000, 194136+8*uintptr(a)) * 4)
						vy := math.Float32bits(*memmap.PtrFloat32(0x587000, 194140+8*uintptr(a)) * 4)
						if dr.Field_117 != vx || dr.Field_118 != vy || dr.Field_119 != 0 || dr.AnimStart != frame || int32(dr.Field_81) != int32(pos[0]) || int32(dr.Field_82) != int32(pos[1]) || dr.PosVec != image.Pt(int(pos[0]), int(pos[1])) {
							t.Fatal("particle trajectory/origin/frame")
						}
						if dr.Field_115 == nil || dr.InClientUpdateList != 1 || dr.ObjFlags&0x200000 == 0 {
							t.Fatal("particle callback/list ownership")
						}
						if callback == nil {
							callback = dr.Field_115
						} else if callback != dr.Field_115 {
							t.Fatal("unstable particle callback")
						}
						byRef[c.refs[dr]] = particle{a, [2]uint32{vx, vy}, dr.PosVec, frame}
					}
					var particles []particle
					for i, call := range c.Calls {
						wantFailure := fail != 0 && (i+1)%fail == 0
						if call.Type != typ || call.Position != image.Pt(int(pos[0]), int(pos[1])) || (call.Ref == 0) != wantFailure {
							t.Fatal("particle allocation contract")
						}
						if call.Ref != 0 {
							p, ok := byRef[call.Ref]
							if !ok || p.Angle != uint16(6*i) {
								t.Fatal("particle angle sequence")
							}
							particles = append(particles, p)
						}
					}
					rows = append(rows, record{pos, frame, fail, cached, append([]effectsSpawnCall(nil), c.Calls...), particles})
				}
			}
		}
	}
	spellbookCapture(t, "client-presentation-turn-undead", rows, "2f10703654404cc7870ce5d166d9b362a53e619ace1aeaa913d490203c401f56")
}
