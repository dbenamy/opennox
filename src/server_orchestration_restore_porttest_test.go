//go:build porttest

package opennox

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/opennox/libs/object"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

// Exercise the real script-ID lookup and the two distinct list failure contracts.
// Postload AI and saved-world cleanup receive separate fixtures.
func TestServerOrchestrationRestoreReferences(t *testing.T) {
	o := newMatchRosterOwner(t)
	globals, restore := legacy.PortTestServerOrchestrationGlobals()
	t.Cleanup(restore)
	_, restorePending := legacy.PortTestMonsterPendingOwner()
	t.Cleanup(restorePending)
	legacy.PortTestMonsterPending("init", 0, 0)
	for _, off := range []uintptr{1563128, 1563132} {
		serverConfigOwnBytes(t, 0x5D4594, off, 4)
		*memmap.PtrUint32(0x5D4594, off) = 0
	}
	t.Cleanup(o.s.PortTestRewardTypes([]string{"SaveGameLocation", "Glyph", "RestoreMonster", "RestoreTarget"}, nil, true, 0, 0))
	oldList, oldPending, oldMissiles := o.s.Objs.List, o.s.Objs.Pending, o.s.Objs.MissileList
	t.Cleanup(func() { o.s.Objs.List = oldList; o.s.Objs.Pending = oldPending; o.s.Objs.MissileList = oldMissiles })
	oldTicks, oldFrame, oldSolo := nox_gameTicks_371764, nox_gameFrame_371772, gameIsSwitchToSolo
	t.Cleanup(func() { nox_gameTicks_371764 = oldTicks; nox_gameFrame_371772 = oldFrame; gameIsSwitchToSolo = oldSolo })
	makeObj := func(name string) *server.Object {
		u := o.s.NewObjectByTypeID(name)
		if u == nil {
			t.Fatal(name)
		}
		t.Cleanup(func() { o.s.Objs.FreeObject(u) })
		return u
	}
	monster, a, b := makeObj("RestoreMonster"), makeObj("RestoreTarget"), makeObj("RestoreTarget")
	savedMonster, savedA, savedB := *monster, *a, *b
	t.Cleanup(func() { *monster = savedMonster; *a = savedA; *b = savedB })
	data := o.record(t, 2560)
	monster.UpdateData = data
	monster.ObjClass = 2
	monster.ObjFlags = 0x80000000
	monster.ScriptIDVal = 99
	a.ScriptIDVal = 100
	a.NetCode = 501
	b.ScriptIDVal = 101
	b.NetCode = 502
	type row struct {
		Arg               int32
		Location, Missing int
		Destroyed         bool
		Pointers          [7]uint32
		Codes             [5]uint32
		Counts            [2]byte
		Reset             bool
	}
	var rows []row
	for _, arg := range []int32{0, 1, 2, -1} {
		for location := 0; location < 4; location++ {
			for missing := -1; missing < 3; missing++ {
				for _, destroyed := range []bool{false, true} {
					t.Run(fmt.Sprintf("arg%d/location%d/missing%d/destroyed%t", arg, location, missing, destroyed), func(t *testing.T) {
						clear(unsafe.Slice((*byte)(data), 2560))
						a.ObjFlags = 0
						b.ObjFlags = 0
						if destroyed {
							b.ObjFlags = object.FlagDestroyed
						}
						monster.ObjNext = nil
						a.ObjNext = b
						b.ObjNext = nil
						monster.InvFirstItem = nil
						o.s.Objs.List = monster
						o.s.Objs.Pending = nil
						o.s.Objs.MissileList = nil
						switch location {
						case 0:
							monster.ObjNext = a
						case 1:
							o.s.Objs.Pending = a
						case 2:
							o.s.Objs.MissileList = a
						case 3:
							monster.InvFirstItem = a
							a.InvNextItem = b
							b.InvNextItem = nil
						}
						defer func() { monster.InvFirstItem = nil; a.InvNextItem = nil }()
						inputs := [3]uint32{100, 101, 100}
						if missing >= 0 {
							inputs[missing] = 999
						}
						for i, v := range inputs {
							objectXferSetWord(data, 1132+4*i, v)
							objectXferSetWord(data, 2140+4*i, v)
						}
						*(*byte)(unsafe.Add(data, 1129)) = 3
						*(*byte)(unsafe.Add(data, 2172)) = 3
						for i, off := range []int{1196, 1216, 392, 1200} {
							objectXferSetWord(data, off, []uint32{100, 101, 100, 999}[i])
						}
						gameIsSwitchToSolo = true
						*globals["restore-cleanup"] = 0
						nox_gameFrame_371772 = 0xffffffff
						noxflags.SetEngine(noxflags.EnginePause)
						before := platformTicks()
						legacy.PortTestServerOrchestration("restore", nil, arg)
						after := platformTicks()
						r := row{Arg: arg, Location: location, Missing: missing, Destroyed: destroyed, Reset: !gameIsSwitchToSolo && *globals["restore-cleanup"] == 0 && nox_gameFrame_371772 == o.s.Frame() && nox_gameTicks_371764 >= before && nox_gameTicks_371764 <= after && !noxflags.HasEngine(noxflags.EnginePause)}
						if !r.Reset {
							t.Fatal("restore epilogue", r)
						}
						lookup := func(id uint32) *server.Object {
							if id == 100 {
								return a
							}
							if id == 101 && !destroyed {
								return b
							}
							return nil
						}
						normalize := func(v uint32) uint32 {
							if v == uint32(uintptr(a.CObj())) {
								return 1
							}
							if v == uint32(uintptr(b.CObj())) {
								return 2
							}
							return v
						}
						wantPtr, wantCode := inputs, inputs
						count := byte(3)
						if arg == 1 {
							for i, id := range inputs {
								u := lookup(id)
								if u == nil {
									wantPtr[i] = 0
									count = 0
									break
								}
								wantPtr[i] = uint32(uintptr(u.CObj()))
							}
							for i, id := range inputs {
								u := lookup(id)
								if u == nil {
									break
								}
								wantCode[i] = u.NetCode
							}
						}
						for i := 0; i < 3; i++ {
							p, c := objectXferGetWord(data, 1132+4*i), objectXferGetWord(data, 2140+4*i)
							if p != wantPtr[i] || c != wantCode[i] {
								t.Fatalf("entry %d: %x/%d want %x/%d", i, p, c, wantPtr[i], wantCode[i])
							}
							r.Pointers[i] = normalize(p)
							r.Codes[i] = c
						}
						for i, off := range []int{1196, 1216, 392, 1200} {
							id := []uint32{100, 101, 100, 999}[i]
							want := id
							if arg == 1 {
								want = 0
								if u := lookup(id); u != nil {
									if i < 2 {
										want = uint32(uintptr(u.CObj()))
									} else {
										want = u.NetCode
									}
								}
							}
							got := objectXferGetWord(data, off)
							if got != want {
								t.Fatalf("field %d: %x want %x", off, got, want)
							}
							if i < 2 {
								r.Pointers[3+i] = normalize(got)
							} else {
								r.Codes[3+i-2] = got
							}
						}
						r.Counts = [2]byte{*(*byte)(unsafe.Add(data, 1129)), *(*byte)(unsafe.Add(data, 2172))}
						if r.Counts != [2]byte{count, count} {
							t.Fatal("list counts", r.Counts, count)
						}
						rows = append(rows, r)
					})
				}
			}
		}
	}
	spellbookCapture(t, "server-orchestration-restore-references", rows, "b67b978535164aab1445f86e1e15d6ddc26c7fb02eaac125dfce421d21a9b7ea")
}

func TestServerOrchestrationRestoreScheduling(t *testing.T) {
	o := newMatchRosterOwner(t)
	_, restore := legacy.PortTestServerOrchestrationGlobals()
	t.Cleanup(restore)
	_, restorePending := legacy.PortTestMonsterPendingOwner()
	t.Cleanup(restorePending)
	legacy.PortTestMonsterPending("init", 0, 0)
	for _, off := range []uintptr{1563128, 1563132} {
		serverConfigOwnBytes(t, 0x5D4594, off, 4)
		*memmap.PtrUint32(0x5D4594, off) = 0
	}
	t.Cleanup(o.s.PortTestRewardTypes([]string{"SaveGameLocation", "Glyph", "RestoreScheduled"}, nil, true, 0, 0))
	oldList, oldUpdates := o.s.Objs.List, o.s.Objs.UpdatableList
	oldTicks, oldFrame, oldSolo := nox_gameTicks_371764, nox_gameFrame_371772, gameIsSwitchToSolo
	t.Cleanup(func() {
		o.s.Objs.List = oldList
		o.s.Objs.UpdatableList = oldUpdates
		nox_gameTicks_371764 = oldTicks
		nox_gameFrame_371772 = oldFrame
		gameIsSwitchToSolo = oldSolo
	})
	u := o.s.NewObjectByTypeID("RestoreScheduled")
	target := o.s.NewObjectByTypeID("RestoreScheduled")
	if u == nil || target == nil {
		t.Fatal("allocation")
	}
	savedU, savedTarget := *u, *target
	t.Cleanup(func() { *u = savedU; *target = savedTarget; o.s.Objs.FreeObject(u); o.s.Objs.FreeObject(target) })
	type row struct {
		Arg       int32
		Class     uint32
		Mode      int
		Sync      uint32
		Scheduled bool
	}
	var rows []row
	for _, arg := range []int32{0, 1, 2} {
		for _, class := range []uint32{0, 0x80, 0x4000, 0x8000, 0x4080, 0x8080, 0xc000} {
			for mode := 0; mode < 4; mode++ {
				t.Run(fmt.Sprintf("arg%d/class%x/mode%d", arg, class, mode), func(t *testing.T) {
					*u = savedU
					*target = savedTarget
					clear(unsafe.Slice((*byte)(u.UpdateData), 64))
					clear(unsafe.Slice((*byte)(target.UpdateData), 64))
					u.ObjClass = object.Class(class)
					u.Field38 = 123
					u.ObjNext = nil
					o.s.Objs.List = u
					o.s.Objs.UpdatableList = nil
					objectXferSetWord(u.UpdateData, 16, uint32(mode&1))
					objectXferSetWord(target.UpdateData, 16, uint32(mode&1))
					if class&0x8000 != 0 {
						if mode&2 != 0 {
							objectXferSetWord(u.UpdateData, 4, uint32(uintptr(target.CObj())))
						}
					} else {
						objectXferSetWord(u.UpdateData, 4, 17)
						objectXferSetWord(u.UpdateData, 12, uint32(17+(mode&2)))
					}
					legacy.PortTestServerOrchestration("restore", nil, arg)
					wantSync := uint32(123)
					wantScheduled := false
					if arg == 1 {
						switch {
						case class&0x4000 != 0:
							if mode&1 != 0 {
								wantSync = 0xffffffff
							}
						case class&0x8000 != 0:
							if mode == 3 {
								wantSync = 0xffffffff
							}
						case class&0x80 != 0:
							wantScheduled = mode&2 != 0
						}
					}
					r := row{arg, class, mode, u.Field38, o.s.Objs.UpdatableList == u}
					if r.Sync != wantSync || r.Scheduled != wantScheduled || (u.IsUpdatable != 0) != wantScheduled {
						t.Fatal("scheduling", r, wantSync, wantScheduled)
					}
					if wantScheduled && (u.UpdatablePrev != nil || u.UpdatableNext != nil) {
						t.Fatal("queue links")
					}
					rows = append(rows, r)
				})
			}
		}
	}
	spellbookCapture(t, "server-orchestration-restore-scheduling", rows, "d403b0c816b97a573960c1e385629836caa0ccb6fddc9cc0962cb61c78d4eadd")
}

func TestServerOrchestrationRestoreCleanup(t *testing.T) {
	o := newMatchRosterOwner(t)
	globals, restore := legacy.PortTestServerOrchestrationGlobals()
	t.Cleanup(restore)
	_, restorePending := legacy.PortTestMonsterPendingOwner()
	t.Cleanup(restorePending)
	legacy.PortTestMonsterPending("init", 0, 0)
	for _, off := range []uintptr{1563128, 1563132, 1565592} {
		serverConfigOwnBytes(t, 0x5D4594, off, 4)
		*memmap.PtrUint32(0x5D4594, off) = 0
	}
	t.Cleanup(o.s.PortTestRewardTypes([]string{"SaveGameLocation", "Glyph", "Pixie", "RestoreCleanup", "OtherMissile"}, nil, true, 0, 0))
	oldList, oldMissiles, oldDeleted := o.s.Objs.List, o.s.Objs.MissileList, o.s.Objs.DeletedList
	oldTicks, oldFrame, oldSolo := nox_gameTicks_371764, nox_gameFrame_371772, gameIsSwitchToSolo
	t.Cleanup(func() {
		o.s.Objs.List = oldList
		o.s.Objs.MissileList = oldMissiles
		o.s.Objs.DeletedList = oldDeleted
		nox_gameTicks_371764 = oldTicks
		nox_gameFrame_371772 = oldFrame
		gameIsSwitchToSolo = oldSolo
	})
	var units []*server.Object
	for _, name := range []string{"RestoreCleanup", "Glyph", "RestoreCleanup", "Pixie", "OtherMissile"} {
		u := o.s.NewObjectByTypeID(name)
		if u == nil {
			t.Fatal(name)
		}
		units = append(units, u)
	}
	saved := make([]server.Object, len(units))
	for i, u := range units {
		saved[i] = *u
	}
	t.Cleanup(func() {
		for i, u := range units {
			*u = saved[i]
			o.s.Objs.FreeObject(u)
		}
	})
	data := o.record(t, 2560)
	type row struct {
		Arg            int32
		Flags, Cleanup uint32
		Protected      bool
		Deleted        [5]bool
		Order          []int
	}
	var rows []row
	for _, arg := range []int32{0, 1, 2} {
		for _, flags := range []uint32{0, 2048, 8192, 10240} {
			for _, cleanup := range []uint32{0, 1, 0xffffffff} {
				for _, protected := range []bool{false, true} {
					t.Run(fmt.Sprintf("arg%d/flags%x/cleanup%x/protected%t", arg, flags, cleanup, protected), func(t *testing.T) {
						defer noxflags.PortTestGameFlags(noxflags.GameFlag(flags))()
						for i, u := range units {
							*u = saved[i]
							u.ObjFlags = 0
							u.ScriptIDVal = 100 + i
						}
						clear(unsafe.Slice((*byte)(data), 2560))
						*(*byte)(unsafe.Add(data, 544)) = 255
						m, g, other, pixie, missile := units[0], units[1], units[2], units[3], units[4]
						m.UpdateData = data
						m.ObjClass = 2
						m.ObjSubClass = 0x2100
						m.InvFirstItem = g
						g.InvNextItem = other
						m.ObjNext = nil
						pixie.ObjClass = 1
						missile.ObjClass = 1
						pixie.ObjOwner = &o.units[0]
						missile.ObjOwner = &o.units[0]
						pixie.ObjNext = missile
						missile.ObjNext = nil
						if protected {
							m.ObjFlags = 0x80000000
							pixie.ObjFlags = 0x80000000
						}
						o.s.Objs.List = m
						o.s.Objs.MissileList = pixie
						o.s.Objs.DeletedList = nil
						*globals["restore-cleanup"] = cleanup
						legacy.PortTestServerOrchestration("restore", nil, arg)
						r := row{Arg: arg, Flags: flags, Cleanup: cleanup, Protected: protected}
						monsterDeleted := arg == 1 && cleanup != 0 && !protected && flags&8192 == 0
						pixieDeleted := arg == 1 && cleanup != 0 && !protected && flags&2048 != 0
						want := [5]bool{monsterDeleted, monsterDeleted, false, pixieDeleted, false}
						for i, u := range units {
							r.Deleted[i] = u.ObjFlags.Has(object.FlagDestroyed)
							if r.Deleted[i] != want[i] {
								t.Fatalf("object %d deletion %v want %v", i, r.Deleted, want)
							}
							if want[i] && u.DeletedAt != o.s.Frame() {
								t.Fatal("deletion frame")
							}
						}
						for u := o.s.Objs.DeletedList; u != nil; u = u.DeletedNext {
							index := -1
							for i, v := range units {
								if u == v {
									index = i
								}
							}
							if index < 0 || len(r.Order) >= 5 {
								t.Fatal("deletion queue")
							}
							r.Order = append(r.Order, index)
						}
						var wantOrder []int
						if pixieDeleted {
							wantOrder = append(wantOrder, 3)
						}
						if monsterDeleted {
							wantOrder = append(wantOrder, 0, 1)
						}
						if fmt.Sprint(r.Order) != fmt.Sprint(wantOrder) {
							t.Fatal("deletion order", r.Order, wantOrder)
						}
						if *globals["restore-cleanup"] != 0 {
							t.Fatal("cleanup state not cleared")
						}
						rows = append(rows, r)
					})
				}
			}
		}
	}
	spellbookCapture(t, "server-orchestration-restore-cleanup", rows, "34a83850ba357e5e09632fca92fe90cedcb6af31273181d33613243af694d33d")
}

func TestServerOrchestrationRestoreWithoutHost(t *testing.T) {
	o := newMatchRosterOwner(t)
	globals, restore := legacy.PortTestServerOrchestrationGlobals()
	t.Cleanup(restore)
	host := o.units[2].UpdateDataPlayer().Player
	old := *host
	t.Cleanup(func() { *host = old })
	oldTicks, oldFrame, oldSolo := nox_gameTicks_371764, nox_gameFrame_371772, gameIsSwitchToSolo
	t.Cleanup(func() { nox_gameTicks_371764 = oldTicks; nox_gameFrame_371772 = oldFrame; gameIsSwitchToSolo = oldSolo })
	var rows [][4]uint32
	for _, active := range []uint32{0, 1} {
		for _, arg := range []int32{0, 1, 2, -1} {
			*host = old
			host.Active = byte(active)
			if active != 0 {
				host.PlayerUnit = nil
			}
			*globals["restore-cleanup"] = 73
			gameIsSwitchToSolo = true
			nox_gameTicks_371764 = 987
			nox_gameFrame_371772 = 654
			noxflags.SetEngine(noxflags.EnginePause)
			legacy.PortTestServerOrchestration("restore", nil, arg)
			if *globals["restore-cleanup"] != 73 || !gameIsSwitchToSolo || nox_gameTicks_371764 != 987 || nox_gameFrame_371772 != 654 || !noxflags.HasEngine(noxflags.EnginePause) {
				t.Fatal("missing host mutated restore state", active, arg)
			}
			rows = append(rows, [4]uint32{active, uint32(arg), *globals["restore-cleanup"], nox_gameFrame_371772})
		}
	}
	spellbookCapture(t, "server-orchestration-restore-without-host", rows, "5b92921676d1f7a7a3fc00b61832ba60a6ce24bb2d8bf634386ed5519b5bf6ca")
}
