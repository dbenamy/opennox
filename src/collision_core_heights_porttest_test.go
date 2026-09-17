//go:build porttest

package opennox

import (
	"encoding/binary"
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/prand"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"math"
	"testing"
	"unsafe"
)

func TestCollisionCoreHeightContacts(t *testing.T) {
	o := newCollisionCoreOwner(t)
	elevator, shaft, u := newObjectXferSimple(t, o.s), newObjectXferSimple(t, o.s), newObjectXferSimple(t, o.s)
	ids := collisionCoreIDs(elevator, shaft, u)
	eu, su := elevator.UpdateData, shaft.UpdateData
	elevator.UpdateData = collisionCoreGuarded(t, o, 64)
	shaft.UpdateData = collisionCoreGuarded(t, o, 64)
	t.Cleanup(func() { elevator.UpdateData = eu; shaft.UpdateData = su })
	ed, sd := unsafe.Slice((*byte)(elevator.UpdateData), 64), unsafe.Slice((*byte)(shaft.UpdateData), 64)
	originalType := u.TypeInd
	oldRNG := o.s.Rand.Logic
	t.Cleanup(func() { u.TypeInd = originalType; o.s.Rand.Logic = oldRNG })
	type row struct {
		AbsScratch        uint32
		Op, Type          string
		Height, Mode      int
		Gap, Size         float32
		Box, Inside       bool
		Object            worldGeometryObjectState
		Z, Vertical, Sync uint32
		PlayerSync        [32]uint32
		Fall              [4]uint32
		Hits              [][5]uint32
		Cache             [5]uint32
		RNG               int
	}
	var rows []row
	for _, op := range []string{"elevator", "shaft"} {
		for _, height := range []int{-16, 64, 96} {
			for _, gap := range []float32{-11, -10, 0, 10, 11} {
				for _, box := range []bool{false, true} {
					for _, size := range []float32{20, 50} {
						for _, typ := range []string{"normal", "SmallFist", "MediumFist", "LargeFist", "Meteor"} {
							for _, inside := range []bool{false, true} {
								modes := []int{0}
								if op == "elevator" {
									modes = []int{0, 1}
								}
								for _, mode := range modes {
									o.resetHits()
									*o.absScratch = 0x3fc00000
									o.resetQueues()
									o.s.Rand.Logic = prand.New(23)
									clear(ed)
									clear(sd)
									for _, name := range []string{"small", "medium", "large", "meteor", "ready"} {
										*o.words[name] = 0
									}
									worldGeometryResetObject(elevator, 1001, 200, 200, true)
									worldGeometryResetObject(shaft, 1002, 100, 100, true)
									for _, v := range []*server.Object{elevator, shaft} {
										v.Shape.Box.W = 40
										v.Shape.Box.H = 40
										legacy.PortTestCollisionCore("shape", v, nil, nil, 0)
									}
									binary.LittleEndian.PutUint32(ed[16:], uint32(height))
									binary.LittleEndian.PutUint32(sd[4:], uint32(uintptr(elevator.CObj())))
									source, surface := elevator, float32(height)
									if op == "shaft" {
										source, surface = shaft, float32(height-64)
									}
									pos := source.NewPos
									if !inside {
										pos.X += 100
										pos.Y += 100
									}
									worldGeometryResetObject(u, 1003, pos.X, pos.Y, box)
									u.Shape.Circle.R = size * 0.5
									u.Shape.Circle.R2 = u.Shape.Circle.R * u.Shape.Circle.R
									if box {
										u.Shape.Box.W = size
										u.Shape.Box.H = size
										legacy.PortTestCollisionCore("shape", u, nil, nil, 0)
									}
									u.TypeInd = originalType
									if typ != "normal" {
										u.TypeInd = uint16(o.s.Types.IndByID(typ))
									}
									u.ZVal = surface + gap
									u.Field27 = 3.5
									u.Field38 = 0
									u.Field140 = [32]uint32{}
									u.Pos39 = types.Pointf{-7, -8}
									u.Field41 = 0xabcdef01
									u.Field42 = 0xabcdef02
									u.ObjFlags = 0x40000
									legacy.PortTestCollisionCore(op, source, u, nil, int32(mode))
									if *o.words["ready"] != 1 {
										t.Fatal("height contact type cache not initialized")
									}
									shouldRaise := typ == "normal" && inside && math.Abs(float64(gap)) <= 10 && (op == "shaft" && size <= 40 || op == "elevator" && mode == 1)
									if shouldRaise {
										if u.ZVal != surface+4 || u.Field27 != 0 || u.ObjFlags&0x140000 != 0x100000 || u.Field38 != 0xffffffff {
											t.Fatal("height transition", op, height, gap, size, u.ZVal, u.Field27, u.ObjFlags, u.Field38)
										}
									}
									if typ != "normal" || op == "elevator" && mode == 0 {
										if u.ZVal != surface+gap || u.Field27 != 3.5 || u.Pos24 != (types.Pointf{}) || u.ObjFlags != object.Flags(0x40000) {
											t.Fatal("height exclusion", op, typ, mode)
										}
									}
									r := row{AbsScratch: *o.absScratch, Op: op, Type: typ, Height: height, Mode: mode, Gap: gap, Size: size, Box: box, Inside: inside, Object: worldGeometryState(u), Z: math.Float32bits(u.ZVal), Vertical: math.Float32bits(u.Field27), Sync: u.Field38, PlayerSync: u.Field140, Fall: [4]uint32{math.Float32bits(u.Pos39.X), math.Float32bits(u.Pos39.Y), u.Field41, u.Field42}, Hits: o.hits(ids), Cache: [5]uint32{*o.words["small"], *o.words["medium"], *o.words["large"], *o.words["meteor"], *o.words["ready"]}, RNG: o.s.Rand.Logic.Index()}
									rows = append(rows, r)
								}
							}
						}
					}
				}
			}
		}
	}
	// An unlinked shaft cannot change the actor, even on an otherwise valid contact.
	clear(sd)
	u.ZVal = 7
	u.Field27 = 9
	u.ObjFlags = 0
	u.Pos24 = types.Pointf{}
	legacy.PortTestCollisionCore("shaft", shaft, u, nil, 0)
	if u.ZVal != 7 || u.Field27 != 9 || u.ObjFlags != 0 || u.Pos24 != (types.Pointf{}) {
		t.Fatal("unlinked shaft changed actor")
	}
	spellbookCapture(t, "collision-core-height-contacts", rows, "d12c4429d0ec7336766fb2acb135ffbe40807c709d500a3b2d2c0fbbad820b8a")
}
