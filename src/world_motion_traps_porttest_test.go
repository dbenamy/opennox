//go:build porttest

package opennox

import (
	"bytes"
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"math"
	"testing"
	"unsafe"
)

type worldMotionTrapOwner struct {
	*collisionCoreOwner
	trap      *server.Object
	data      *[16]uint32
	globals   map[string]*uint32
	configure func(int)
	unchanged func() bool
}

func newWorldMotionTrapOwner(t *testing.T) *worldMotionTrapOwner {
	t.Helper()
	o := &worldMotionTrapOwner{collisionCoreOwner: newCollisionCoreOwner(t)}
	// The shared collision owner already owns stateDirection scratch globals.
	// Facing admission additionally reads this shipped 9-by-16 table.
	for off, data := range blobdata.PortTestCombatTables() {
		copy(serverConfigOwnBytes(t, 0x587000, off, len(data)), data)
	}
	names := []string{"ArrowTrap1", "ArrowTrap2", "MercArcherArrow", "MotionTestArrow", "MotionTestTarget"}
	t.Cleanup(o.s.PortTestRewardTypes(names, nil, true, 0, 0))
	for _, name := range names[2:4] {
		typ := o.s.Types.ByID(name)
		typ.Speed = 12.5
		typ.SpeedBase = 12.5
		typ.CollideData = collisionCoreGuarded(t, o.collisionCoreOwner, 8)
		typ.CollideDataSize = 8
		*(*[2]uint32)(typ.CollideData) = [2]uint32{17, 19}
	}
	var restore func()
	o.globals, restore = legacy.PortTestWorldMotionGlobals()
	t.Cleanup(restore)
	o.balance(map[string]float64{"ArrowTrapDamage": 37.75})
	o.s.PortTestAIEmptyMap()
	t.Cleanup(o.s.Map.Free)
	o.configure, o.unchanged, restore = o.s.PortTestPathWalls()
	t.Cleanup(restore)
	o.trap = newObjectXferSimple(t, o.s)
	saved := *o.trap
	o.data = (*[16]uint32)(collisionCoreGuarded(t, o.collisionCoreOwner, 64))
	oldPending := o.s.Objs.Pending
	t.Cleanup(func() { *o.trap = saved; o.s.Objs.Pending = oldPending })
	o.trap.UpdateData = unsafe.Pointer(o.data)
	for i := range o.units {
		u := &o.units[i]
		saved := *u
		pl := u.UpdateDataPlayer().Player
		oldPlayer := *pl
		t.Cleanup(func() { *u = saved; *pl = oldPlayer })
		u.PosVec = types.Pointf{100, 100}
		u.NewPos = u.PosVec
		u.ObjFlags = 4
		u.ObjOwner = nil
		// Type zero aliases absent creature IDs in enemy classification.
		u.TypeInd = uint16(o.s.Types.IndByID("MotionTestTarget"))
		pl.Field10 = 1000
		pl.Field12 = 1000
		pl.Field3680 = 0
		pl.CameraFollowObj = nil
	}
	return o
}

type worldMotionArrow struct {
	Type               uint16
	Flags, Owner       uint32
	Position, Velocity [2]uint32
	Direction          [2]uint16
	Damage             [2]uint32
}

func (o *worldMotionTrapOwner) created(t *testing.T) []worldMotionArrow {
	t.Helper()
	// Capture each real allocation, then release it before the next shot so the
	// fixed object pool is reusable throughout the fixture.
	pending := o.s.Objs.Pending
	defer func() {
		for u := pending; u != nil; {
			next := u.ObjNext
			o.s.ObjClearOwner(u)
			if u.IDPtr != nil {
				legacy.PortTestObjectXferFreeName(u.IDPtr)
				u.IDPtr = nil
			}
			for _, pp := range []*unsafe.Pointer{&u.InitData, &u.CollideData, &u.UpdateData, &u.Field189} {
				if *pp != nil {
					alloc.FreePtr(*pp)
					*pp = nil
				}
			}
			o.s.Objs.FreeObject(u)
			u = next
		}
		o.s.Objs.Pending = nil
	}()
	var rows []worldMotionArrow
	for u := o.s.Objs.Pending; u != nil; u = u.ObjNext {
		if len(rows) > 2 {
			t.Fatal("unexpected projectile chain")
		}
		if u.ObjOwner != o.trap {
			t.Fatal("trap projectile owner")
		}
		rows = append(rows, worldMotionArrow{u.TypeInd, uint32(u.ObjFlags), 1001, motionBits(u.PosVec), motionBits(u.VelVec), [2]uint16{uint16(u.Direction1), uint16(u.Direction2)}, *(*[2]uint32)(u.CollideData)})
	}
	return rows
}

func TestWorldMotionTrapProjectile(t *testing.T) {
	o := newWorldMotionTrapOwner(t)
	type row struct {
		Trap, Projectile string
		Dir              uint16
		Radius           uint32
		Created          []worldMotionArrow
		Sounds           []int
		Cache            [3]uint32
	}
	var rows []row
	for _, trap := range []string{"ArrowTrap1", "ArrowTrap2", "Trigger"} {
		for _, projectile := range []string{"MercArcherArrow", "MotionTestArrow", "missing"} {
			for _, dir := range []server.Dir16{0, 32, 64, 128, 192, 255} {
				for _, radius := range []float32{0, 10, 20.125} {
					worldGeometryResetObject(o.trap, 1001, 100.25, 200.75, false)
					o.trap.TypeInd = uint16(o.s.Types.IndByID(trap))
					o.trap.Shape.Circle.R = radius
					o.trap.Direction1 = dir
					o.trap.ObjFlags = 4
					o.s.Objs.Pending = nil
					o.s.PortTestCombatAudioReset()
					for _, name := range []string{"trap-arrow", "trap-one", "trap-two"} {
						*o.globals[name] = 0
					}
					typ := int32(o.s.Types.IndByID(projectile))
					legacy.PortTestWorldMotionTimed("trap-projectile", o.trap, nil, typ)
					created := o.created(t)
					if (len(created) == 1) != (typ != 0) {
						t.Fatal("trap projectile allocation", projectile)
					}
					if len(created) == 1 {
						a := created[0]
						if a.Type != uint16(typ) || a.Direction != ([2]uint16{uint16(dir), uint16(dir)}) {
							t.Fatal("trap projectile direction/type")
						}
						wantDamage := [2]uint32{17, 19}
						if trap != "Trigger" {
							wantDamage = [2]uint32{37, 37}
						}
						if a.Damage != wantDamage {
							t.Fatal("trap damage balance truncation", a.Damage, wantDamage)
						}
						if dir == 0 && (a.Position != motionBits(types.Pointf{100.25 + radius + 4, 200.75}) || a.Velocity != motionBits(types.Pointf{12.5, 0})) {
							t.Fatal("trap projectile spawn offset and speed", a)
						}
					}
					var sounds []int
					for _, e := range o.s.PortTestCombatAudioSnapshot() {
						if e.Obj != o.trap {
							t.Fatal("trap shot sound owner")
						}
						sounds = append(sounds, int(e.ID))
					}
					if (len(sounds) == 1) != (projectile == "MercArcherArrow") || len(sounds) > 1 || len(sounds) == 1 && sounds[0] != 889 {
						t.Fatal("trap arrow audio", projectile, sounds)
					}
					rows = append(rows, row{trap, projectile, uint16(dir), math.Float32bits(radius), created, sounds, [3]uint32{*o.globals["trap-arrow"], *o.globals["trap-one"], *o.globals["trap-two"]}})
				}
			}
		}
	}
	spellbookCapture(t, "world-motion-trap-projectile", rows, "712362128d45a3c215cf92316c77aa0daeee37e84cb4970bbd3a82d6ae24c581")
}

func TestWorldMotionTrapScan(t *testing.T) {
	o := newWorldMotionTrapOwner(t)
	target := &o.units[0]
	saved := *target
	type row struct {
		Wall, Position, Kind int
		Flags                uint32
		SameOwner            bool
		Direct, Scan, State  uint32
		Data                 [16]uint32
	}
	var rows []row
	for wall := 0; wall < 3; wall++ {
		for position, p := range []types.Pointf{{200, 100}, {50, 100}, {100, 200}, {450, 100}, {451, 100}} {
			for kind := 0; kind < 2; kind++ {
				for _, flags := range []object.Flags{4, 4 | 0x20, 4 | 0x8000} {
					for _, same := range []bool{false, true} {
						o.s.PortTestAIEmptyMap()
						o.configure(wall)
						*target = saved
						worldGeometryResetObject(o.trap, 1001, 100, 100, false)
						o.trap.ObjFlags = 4
						o.trap.Direction1 = 0
						o.trap.TypeInd = uint16(o.s.Types.IndByID("ArrowTrap1"))
						o.trap.ObjOwner = nil
						worldGeometryResetObject(target, 1002, p.X, p.Y, false)
						target.ObjClass = object.ClassPlayer
						if kind == 1 {
							target.ObjClass = object.ClassSimple
						}
						target.ObjFlags = flags
						if same {
							o.trap.ObjOwner = target
						}
						o.s.Map.AddObjectToIndex(target)
						*o.data = [16]uint32{}
						o.data[1] = 77
						*o.globals["trap-reachable"] = 0
						legacy.PortTestWorldMotionTimed("trap-candidate", o.trap, target, 0)
						direct := *o.globals["trap-reachable"]
						scan := uint32(legacy.PortTestWorldMotionTimed("trap-scan", o.trap, nil, 0))
						legacy.PortTestWorldMotionTimed("trap-state", o.trap, nil, 0)
						if (kind == 1 || flags&0x8020 != 0 || same) && (direct != 0 || scan != 0 || o.data[2]&0xff != 0) {
							t.Fatal("trap candidate exclusion")
						}
						if wall == 0 && position == 0 && kind == 0 && flags == 4 && !same && (direct != 1 || scan != 1 || o.data[2]&0xff != 1 || o.data[1] != 0) {
							t.Fatal("trap must acquire visible enemy and reset shot timer", direct, scan, *o.data, "enemy", o.s.IsEnemyTo(o.trap, target), "vision", o.s.MapTraceVision(target, o.trap), "facing", legacy.Nox_server_testTwoPointsAndDirection_4E6E50(o.trap.PosVec, int16(o.trap.Direction1), target.PosVec), "teams", o.trap.TeamVal.ID, target.TeamVal.ID)
						}
						if !o.unchanged() {
							t.Fatal("trap scan changed wall")
						}
						rows = append(rows, row{wall, position, kind, uint32(flags), same, direct, scan, *o.globals["trap-reachable"], *o.data})
					}
				}
			}
		}
	}
	spellbookCapture(t, "world-motion-trap-scan", rows, "62ba07d987e3c2596d82c84ddea181f920df0af1e2877773dfc0d6b5c941ffdd")
}

func TestWorldMotionTrapUpdate(t *testing.T) {
	o := newWorldMotionTrapOwner(t)
	target := &o.units[0]
	saved := *target
	type initial struct {
		Name       string
		Scan, Shot uint32
		State, On  byte
	}
	cases := []initial{{"first-power", 99, 99, 0, 0}, {"scan-now", 0, 0, 0, 1}, {"fire-ready", 5, 0, 1, 1}, {"cooldown", 5, 2, 1, 1}, {"rescan", 1, 30, 1, 1}, {"other-state", 5, 0, 2, 1}, {"wrapped", 0xffffffff, 0xffffffff, 1, 1}}
	type row struct {
		Name, Trap      string
		Target, Powered bool
		Step            int
		Return          int32
		Data            [16]uint32
		Arrows          []worldMotionArrow
		Packets         [][]byte
	}
	var rows []row
	for _, trap := range []string{"ArrowTrap1", "ArrowTrap2", "Trigger"} {
		for _, hasTarget := range []bool{false, true} {
			for _, powered := range []bool{false, true} {
				for _, init := range cases {
					o.s.PortTestAIEmptyMap()
					o.configure(0)
					*target = saved
					worldGeometryResetObject(o.trap, 1001, 100, 100, false)
					o.trap.TypeInd = uint16(o.s.Types.IndByID(trap))
					o.trap.Direction1 = 0
					o.trap.ObjOwner = nil
					o.trap.ObjFlags = 4
					if powered {
						o.trap.ObjFlags |= 0x1000000
					}
					worldGeometryResetObject(target, 1002, 200, 100, false)
					target.ObjClass = object.ClassPlayer
					target.ObjFlags = 4
					if hasTarget {
						o.s.Map.AddObjectToIndex(target)
					}
					*o.data = [16]uint32{}
					o.data[0] = init.Scan
					o.data[1] = init.Shot
					o.data[2] = uint32(init.State)
					o.data[3] = uint32(o.s.Types.IndByID("MotionTestArrow"))
					o.data[12] = uint32(init.On)
					o.s.SetTickRate(30)
					for step := 0; step < 3; step++ {
						before := *o.data
						o.s.NetList.ResetAll()
						o.s.Objs.Pending = nil
						o.s.PortTestCombatAudioReset()
						rv := legacy.PortTestWorldMotionTimed("trap-update", o.trap, nil, 0)
						arrows := o.created(t)
						packets := visibilityEffectsPackets(o.s)
						if !powered {
							want := before
							want[12] &^= 0xff
							if *o.data != want || len(arrows) != 0 || rv != 4 {
								t.Fatal("unpowered trap contract")
							}
						}
						if powered && step == 0 && init.Name == "first-power" {
							if o.data[0] != 29 || (o.data[2]&0xff == 1) != hasTarget || (len(arrows) == 1) != hasTarget {
								t.Fatal("trap power-on reset/acquisition")
							}
						}
						if powered && step == 0 && init.Name == "fire-ready" && (len(arrows) != 1 || o.data[0] != 4 || o.data[1] != 29 || rv != 29) {
							t.Fatal("trap firing cooldown contract")
						}
						wantFX := byte(0)
						if len(arrows) > 0 && trap == "ArrowTrap1" {
							wantFX = 1
						}
						if len(arrows) > 0 && trap == "ArrowTrap2" {
							wantFX = 2
						}
						packetCount := 0
						for _, p := range packets {
							if len(p) > 0 {
								packetCount++
								if wantFX == 0 || !bytes.Equal(p, []byte{161, 100, 0, 100, 0, wantFX}) {
									t.Fatal("trap effect packet", p, wantFX)
								}
							}
						}
						if wantFX != 0 && packetCount == 0 {
							t.Fatal("trap effect not delivered to actual players")
						}
						rows = append(rows, row{init.Name, trap, hasTarget, powered, step, rv, *o.data, arrows, packets})
					}
				}
			}
		}
	}
	spellbookCapture(t, "world-motion-trap-update", rows, "b631cb3eb8114b2067c493c670980eca66afedf3826d2e72f229c17b8b227de4")
}
