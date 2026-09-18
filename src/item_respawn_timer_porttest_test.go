//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestItemRespawnScheduling(t *testing.T) {
	o := newMatchRosterOwner(t)
	allow, restore := legacy.PortTestItemRespawnGlobals()
	t.Cleanup(restore)
	type row struct {
		Name                     string
		Allow, Deadline, Pending uint32
		Cleared                  bool
	}
	var rows []row
	oldRules := int(legacy.PortTestServerConfigScalar("flags-get", 0, 0))
	t.Cleanup(func() { legacy.PortTestServerConfigScalar("flags-set", oldRules, 0) })
	for _, rules := range []int{0, 2} {
		legacy.PortTestServerConfigScalar("flags-set", rules, 0)
		for _, allowed := range []bool{false, true} {
			rt := o.s.PortTestRewardTypes([]string{"Crown", "RespawnFixture"}, nil, allowed, 0, 0)
			for _, class := range []uint32{0, 2, 0x1000, 0x1000000, 0x2000000, 0x10000000} {
				for _, flags := range []uint32{0, 0x20, 0x8000, 0x8020} {
					for _, held := range []bool{false, true} {
						for _, crown := range []bool{false, true} {
							name := fmt.Sprintf("rules%d/allowed%v/class%x/flags%x/held%v/crown%v", rules, allowed, class, flags, held, crown)
							t.Run(name, func(t *testing.T) {
								u := &o.units[0]
								u.ObjClass = object.Class(class)
								u.ObjFlags = object.Flags(flags)
								u.PosVec = types.Pointf{}
								u.Field32 = 123
								u.TypeInd = uint16(o.s.Types.IndByID("RespawnFixture"))
								if crown {
									u.TypeInd = uint16(o.s.Types.IndByID("Crown"))
								}
								u.InvHolder = nil
								if held {
									u.InvHolder = &o.units[1]
								}
								oldInit, oldUse := u.InitData, u.UseData
								u.InitData = o.record(t, 20)
								u.UseData.Ptr = o.record(t, 64)
								defer func() { u.InitData, u.UseData = oldInit, oldUse; u.InvHolder = nil }()
								cleanup := legacy.PortTestTeamRuntimeRespawns([]*server.Object{u})
								defer cleanup()
								o.s.SetFrame(123)
								o.s.SetTickRate(30)
								defer noxflags.PortTestGameFlags(0)()
								legacy.PortTestItemRespawn("tick", nil)
								r := legacy.PortTestItemRespawnRecords()[0]
								pending, cleared := false, false
								if class&2 != 0 {
									pending = flags&0x8020 != 0
									cleared = flags&0x20 != 0
								} else if flags&0x20 != 0 {
									cleared = true
									pending = allowed
								} else if class&0x3001000 != 0 || crown {
									pending = held && allowed && !crown && rules&2 != 0
								} else {
									pending = held
								}
								wantPending := uint32(0)
								wantDeadline := uint32(0)
								if pending {
									wantPending = 1
									wantDeadline = 1023
								}
								if *allow != 0 || r[5] != wantDeadline || r[6] != wantPending || (r[1] == 0) != cleared {
									t.Fatalf("schedule got allow%d record%v want%d,%d,%v", *allow, *r, wantDeadline, wantPending, cleared)
								}
								rows = append(rows, row{name, *allow, r[5], r[6], r[1] == 0})
							})
						}
					}
				}
			}
			rt()
		}
	}
	legacy.PortTestServerConfigScalar("flags-set", 0, 0)
	// Gated modes must leave both insertion and record state untouched; missing-object
	// scheduling is independent of definition eligibility until the deadline.
	rt := o.s.PortTestRewardTypes([]string{"Crown", "RespawnFixture"}, nil, false, 0, 0)
	defer rt()
	u := &o.units[0]
	u.TypeInd = uint16(o.s.Types.IndByID("RespawnFixture"))
	u.ObjClass = 0
	u.InvHolder = nil
	for _, flags := range []uint32{0, 512, 4096, 4608, 8192} {
		for _, frame := range []uint32{0, 123, 0x7fffffff, 0xfffffff0} {
			cleanup := legacy.PortTestTeamRuntimeRespawns([]*server.Object{u})
			r := legacy.PortTestItemRespawnRecords()[0]
			r[1] = 0
			o.s.SetFrame(frame)
			o.s.SetTickRate(30)
			unflags := noxflags.PortTestGameFlags(noxflags.GameFlag(flags))
			legacy.PortTestItemRespawn("tick", nil)
			gated := flags&4608 != 0
			want := uint32(1)
			deadline := frame + 900
			enabled := uint32(0)
			if gated {
				want = 0
				deadline = 0
				enabled = 1
			}
			if r[6] != want || r[5] != deadline || *allow != enabled {
				t.Fatal("missing/gated timing", flags, frame, *r, *allow)
			}
			rows = append(rows, row{fmt.Sprintf("missing/flags%x/frame%x", flags, frame), *allow, r[5], r[6], true})
			unflags()
			cleanup()
		}
	}
	spellbookCapture(t, "item-respawn-scheduling", rows, "52d93cc47d8292a8c650d6014bf902a3d41c604d021e2b7ab7208501dda12bec")
}

func TestItemRespawnRecreation(t *testing.T) {
	o := newMatchRosterOwner(t)
	t.Cleanup(o.s.PortTestRewardTypes([]string{"Crown", "RespawnFixture"}, nil, true, 0, 0))
	_, restore := legacy.PortTestItemRespawnGlobals()
	t.Cleanup(restore)
	oldPending := o.s.Objs.Pending
	t.Cleanup(func() { o.s.Objs.Pending = oldPending })
	type row struct {
		Frame      uint32
		Pending    uint32
		Created    bool
		Position   types.Pointf
		Dir1, Dir2 uint16
		Flags      uint32
	}
	var rows []row
	for _, frame := range []uint32{499, 500, 501, 0x7fffffff, 0x80000000, 0xffffffff} {
		t.Run(fmt.Sprint(frame), func(t *testing.T) {
			o.reset()
			o.s.SetFrame(frame)
			o.s.Objs.Pending = nil
			u := o.s.NewObjectByTypeID("RespawnFixture")
			if u == nil {
				t.Fatal("source object allocation")
			}
			t.Cleanup(func() { o.s.Objs.FreeObject(u) })
			u.TypeInd = uint16(o.s.Types.IndByID("RespawnFixture"))
			u.ObjClass = 0
			u.PosVec = types.Pointf{X: 42.25, Y: 81.5}
			u.Direction1 = 237
			cleanup := legacy.PortTestTeamRuntimeRespawns([]*server.Object{u})
			defer cleanup()
			r := legacy.PortTestItemRespawnRecords()[0]
			r[1] = 0
			r[5] = 500
			r[6] = 1
			legacy.PortTestItemRespawn("tick", nil)
			got := (*server.Object)(unsafe.Pointer(uintptr(r[1])))
			want := frame >= 500
			if (got != nil) != want || r[6] != uint32(boolInt(!want)) {
				t.Fatal("deadline creation", frame, *r)
			}
			state := row{Frame: frame, Pending: r[6], Created: got != nil}
			if got != nil {
				if got != o.s.Objs.Pending || got.TypeInd != u.TypeInd || got.PosVec != u.PosVec || got.Direction1 != 237 || got.Direction2 != 237 {
					t.Fatal("created fields")
				}
				state.Position = got.PosVec
				state.Dir1 = uint16(got.Direction1)
				state.Dir2 = uint16(got.Direction2)
				state.Flags = uint32(got.ObjFlags)
				o.s.Objs.Pending = nil
				got.ObjNext = nil
				got.ObjPrev = nil
				o.s.Objs.FreeObject(got)
			}
			rows = append(rows, state)
		})
	}
	spellbookCapture(t, "item-respawn-recreation", rows, "9749924edadc73189643338ebfe7643f86b8adbe65b14fccfdfcc6190443f2cd")
}
