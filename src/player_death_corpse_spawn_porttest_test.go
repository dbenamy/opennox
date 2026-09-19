//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"github.com/opennox/libs/prand"
	"github.com/opennox/libs/types"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"math"
	"testing"
)

func TestPlayerDeathCorpseSpawn(t *testing.T) {
	o := newMatchRosterOwner(t)
	t.Cleanup(o.s.PortTestRewardTypes([]string{"CorpseFixture"}, nil, true, 0, 0))
	globals, restore := legacy.PortTestWorldMotionListGlobals()
	t.Cleanup(restore)
	vectors := blobdata.PortTestWorldGeometryTables()[192088]
	copy(serverConfigOwnBytes(t, 0x587000, 192088, len(vectors)), vectors)
	points := blobdata.PortTestPlayerCorpsePoints()
	copy(serverConfigOwnBytes(t, 0x587000, 280376, len(points)), points)
	cache := serverConfigOwnBytes(t, 0x5D4594, 2488736, 400)
	typ := uint32(o.s.Types.IndByID("CorpseFixture"))
	if typ == 0 {
		t.Fatal("fixture type")
	}
	type spawned struct {
		Type            uint16
		Position        types.Pointf
		Deadline, Flags uint32
	}
	type record struct {
		Name         string
		Objects      []spawned
		Logic, Other int
	}
	var rows []record
	for _, angle := range []int32{0, 32, 64, 96, 128, 160, 192, 224} {
		for _, count := range []int{0, 1, 5, 10, 11} {
			for mode := 0; mode < 4; mode++ {
				name := fmt.Sprintf("angle=%d/count=%d/mode=%d", angle, count, mode)
				t.Run(name, func(t *testing.T) {
					o.reset()
					noxflags.ResetGame()
					noxflags.SetGame(noxflags.GameHost)
					if mode&1 != 0 {
						noxflags.SetGame(noxflags.GameOnline)
					}
					*o.words["rateMode"] = uint32(mode >> 1)
					for i := 0; i < 99; i++ {
						binary.LittleEndian.PutUint32(cache[4+4*i:], typ)
					}
					binary.LittleEndian.PutUint32(cache, 1)
					sign := func(x int32) int {
						if x > 6 {
							return 1
						}
						if x < -6 {
							return -1
						}
						return 0
					}
					x := int32(binary.LittleEndian.Uint32(vectors[8*angle:]))
					y := int32(binary.LittleEndian.Uint32(vectors[8*angle+4:]))
					dir := 4 + sign(x) + 3*sign(y)
					if dir == 4 {
						t.Fatal("direction table must be nontrivial")
					}
					if count < 11 {
						binary.LittleEndian.PutUint32(cache[4+4*(11*dir+count):], 0)
					}
					oldPending := o.s.Objs.Pending
					o.s.Objs.Pending = nil
					t.Cleanup(func() { o.s.Objs.Pending = oldPending })
					oldLogic, oldOther := o.s.Rand.Logic, o.s.Rand.Other
					t.Cleanup(func() { o.s.Rand.Logic, o.s.Rand.Other = oldLogic, oldOther })
					o.s.Rand.Logic, o.s.Rand.Other = prand.New(17), prand.New(29)
					expectedRNG := prand.New(17)
					pos := types.Pointf{X: 123.5, Y: 247.25}
					legacy.PortTestPlayerCorpseCreate(pos, angle)
					var objects []*server.Object
					for u := o.s.Objs.Pending; u != nil; u = u.ObjNext {
						if len(objects) > 11 {
							t.Fatal("pending list cycle")
						}
						objects = append(objects, u)
					}
					for _, u := range objects {
						trackObjectXferTyped(t, o.s, u)
					}
					t.Cleanup(func() {
						o.s.Objs.Pending = nil
						for _, u := range objects {
							legacy.PortTestWorldMotionList("decay-remove", u, 0)
							u.ObjNext = nil
							u.ObjPrev = nil
						}
					})
					if len(objects) != count {
						t.Fatalf("created %d want %d", len(objects), count)
					}
					r := record{Name: name}
					for i := 0; i < count; i++ {
						u := objects[count-1-i]
						p := points[8*(11*dir+i):]
						wantPos := types.Pointf{X: pos.X + math.Float32frombits(binary.LittleEndian.Uint32(p)), Y: pos.Y + math.Float32frombits(binary.LittleEndian.Uint32(p[4:]))}
						deadline := o.s.Frame() + o.s.TickRate()*uint32(expectedRNG.IntClamp(10, 20))
						if u.TypeInd != uint16(typ) || u.PosVec != wantPos || u.Field34 != deadline {
							t.Fatalf("object %d type/position/deadline %d/%+v/%d want %d/%+v/%d", i, u.TypeInd, u.PosVec, u.Field34, typ, wantPos, deadline)
						}
						if uint32(u.ObjFlags)&0x40 != uint32(bool2int(mode == 3))*0x40 {
							t.Fatalf("online corpse flag %x", u.ObjFlags)
						}
						r.Objects = append(r.Objects, spawned{u.TypeInd, u.PosVec, u.Field34, uint32(u.ObjFlags)})
					}
					r.Logic, r.Other = o.s.Rand.Logic.Index(), o.s.Rand.Other.Index()
					if r.Logic != expectedRNG.Index() || r.Other != 29 {
						t.Fatalf("RNG %d/%d", r.Logic, r.Other)
					}
					if count == 0 && *globals["decay"] != 0 {
						t.Fatal("empty creation changed decay list")
					}
					rows = append(rows, r)
				})
			}
		}
	}
	spellbookCapture(t, "player-death-corpse-spawn", rows, "")
}
