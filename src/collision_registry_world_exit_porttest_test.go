//go:build porttest

package opennox

import (
	"bytes"
	"fmt"
	"github.com/opennox/libs/object"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"strings"
	"testing"
	"unsafe"
)

func TestCollisionRegistryWorldExitAdmission(t *testing.T) {
	o := newWorldCollisionOwner(t)
	a := newObjectXferSimple(t, o.s)
	b := &o.units[1]
	other := newObjectXferSimple(t, o.s)
	data := o.record(t, 88)
	a.CollideData = data
	t.Cleanup(func() { a.CollideData = nil })
	words := worldCollisionExitGlobals(t, o)
	var rows []struct {
		Name             string
		Marked, Switched bool
		Path, Waypoint   string
	}
	defer func() {
		collisionRegistryCapture(t, "collision-registry-world-exit-admission", rows)
	}()
	for _, mode := range []string{"clear", "nil", "nonplayer", "busy", "retry", "update", "dead", "book", "pause", "warp-closed"} {
		for _, path := range []string{"", "WorldTest.map", "WorldTest.map:Entry"} {
			name := fmt.Sprintf("%s/path%q", mode, path)
			t.Run(name, func(t *testing.T) {
				o.reset()
				flags := noxflags.GameFlag(0)
				if mode == "pause" {
					flags = noxflags.GamePause
				}
				defer noxflags.PortTestGameFlags(flags)()
				clear(unsafe.Slice((*byte)(data), 88))
				copy(unsafe.Slice((*byte)(data), 88), path)
				dword_5d4594_1563080 = mode == "busy"
				dword_5d4594_1563092 = 0
				dword_5d4594_1563088 = 0
				dword_5d4594_1563064 = false
				if mode == "retry" {
					dword_5d4594_1563092 = 1
				}
				*words["dword_5d4594_1047520"] = 0
				if mode == "book" {
					*words["dword_5d4594_1047520"] = 1
				}
				b.ObjFlags = 0
				if mode == "dead" {
					b.ObjFlags = object.FlagDead
				}
				objectXferSetWord(b.UpdateData, 284, 0)
				if mode == "update" {
					objectXferSetWord(b.UpdateData, 284, 1)
				}
				a.ObjSubClass = 0
				if mode == "warp-closed" {
					a.ObjSubClass = 2
				}
				*o.globals["warpOpen"] = 0
				legacy.Set_dword_5d4594_1548524(0)
				noxServer.mapSwitchWPName = ""
				clear(unsafe.Slice(memmap.PtrUint8(0x5D4594, 2598188), 80))
				target := b
				if mode == "nil" {
					target = nil
				}
				if mode == "nonplayer" {
					target = other
				}
				collisionRegistryWorld(7, a, target, nil)
				admitted := mode == "clear" && path != ""
				if dword_5d4594_1563064 != admitted || (legacy.Get_dword_5d4594_1548524() != 0) != admitted {
					t.Fatal("exit admission/map request")
				}
				raw := unsafe.Slice(memmap.PtrUint8(0x5D4594, 2598188), 80)
				n := bytes.IndexByte(raw, 0)
				got := string(raw[:n])
				wantPath, wantWP := "", ""
				if admitted {
					pieces := strings.SplitN(path, ":", 2)
					wantPath = strings.ToLower(pieces[0])
					if len(pieces) == 2 {
						wantWP = pieces[1]
					}
				}
				if got != wantPath || noxServer.mapSwitchWPName != wantWP {
					t.Fatal("queued map/waypoint")
				}
				rows = append(rows, struct {
					Name             string
					Marked, Switched bool
					Path, Waypoint   string
				}{name, dword_5d4594_1563064, legacy.Get_dword_5d4594_1548524() != 0, got, noxServer.mapSwitchWPName})
			})
		}
	}
	if collisionRegistryWorld(3, nil, nil, nil) != 0 {
		t.Fatal("pending map buffer")
	}
}

func TestCollisionRegistryWorldQuestExit(t *testing.T) {
	var rows []struct {
		Name                              string
		Stage, Highest, Exit, Warp, Buffs uint32
		Switched, Observer                bool
		Pending                           string
	}
	defer func() {
		collisionRegistryCapture(t, "collision-registry-world-quest-exit", rows)
	}()
	for _, kind := range []uint32{1, 2} {
		for _, stage := range []uint32{0, 19} {
			for _, allReady := range []bool{false, true} {
				name := fmt.Sprintf("kind%d/stage%d/all-ready%t", kind, stage, allReady)
				t.Run(name, func(t *testing.T) {
					o := newWorldCollisionOwner(t)
					words := worldCollisionExitGlobals(t, o)
					*words["dword_5d4594_1047520"] = 0
					dword_5d4594_1563080 = false
					dword_5d4594_1563092 = 0
					dword_5d4594_1563088 = 0
					dword_5d4594_1563064 = false
					legacy.Set_dword_5d4594_1548524(0)
					oldGet := legacy.GetServer
					catalog := legacy.PortTestMapCatalogOpen(7)
					t.Cleanup(catalog.Close)
					catalog.Add(legacy.PortTestMapCatalogEntry{Name: "WorldQuest", Enabled: 1, Flags: 2})
					if catalog.BuildQuest() != 1 {
						t.Fatal("quest map catalog")
					}
					legacy.GetServer = oldGet
					oldAllow, oldInfinite, oldInc := questAllowDefault, questLevelWarpInfinite, questLevelWarpInc
					t.Cleanup(func() { questAllowDefault = oldAllow; questLevelWarpInfinite = oldInfinite; questLevelWarpInc = oldInc })
					questAllowDefault = true
					questLevelWarpInfinite = false
					questLevelWarpInc = 5
					defer noxflags.PortTestGameFlags(noxflags.GameModeQuest)()
					*memmap.PtrUint32(0x587000, 202028) = stage
					*o.globals["warpOpen"] = 2
					a := newObjectXferSimple(t, o.s)
					a.ObjSubClass = object.SubClass(kind)
					a.CollideData = o.record(t, 88)
					t.Cleanup(func() { a.CollideData = nil })
					b := &o.units[1]
					for i := range o.units {
						u := &o.units[i]
						p := u.UpdateDataPlayer().Player.C()
						objectXferSetWord(p, 4792, 1)
						objectXferSetWord(p, 4696, 0)
						objectXferSetWord(u.UpdateData, 312, 0)
						objectXferSetWord(u.UpdateData, 316, 0)
						if allReady && u != b {
							off := 312
							if kind == 2 {
								off = 316
							}
							objectXferSetWord(u.UpdateData, off, uint32(uintptr(a.CObj())))
						}
					}
					before := platformTicks()
					collisionRegistryWorld(7, a, b, nil)
					after := platformTicks()
					exit, warp := objectXferGetWord(b.UpdateData, 312), objectXferGetWord(b.UpdateData, 316)
					if (exit != 0) != (kind == 1) || (warp != 0) != (kind == 2) {
						t.Fatal("quest portal association")
					}
					if exit != 0 {
						if exit != uint32(uintptr(a.CObj())) {
							t.Fatal("exit pointer")
						}
						exit = 1
					}
					if warp != 0 {
						if warp != uint32(uintptr(a.CObj())) {
							t.Fatal("warp pointer")
						}
						warp = 1
					}
					switched := legacy.Get_dword_5d4594_1548524() != 0
					if switched != allReady || !dword_5d4594_1563064 || b.Update != legacy.Get_nox_xxx_updatePlayerObserver_4E62F0() {
						t.Fatal("quest observer/map request")
					}
					wantStage := stage
					if kind == 2 && allReady {
						wantStage = (stage/5+1)*5 - 1
					}
					highest := objectXferGetWord(b.UpdateDataPlayer().Player.C(), 4696)
					wantHighest := uint32(0)
					if kind == 1 {
						wantHighest = stage + 1
					}
					if memmap.Uint32(0x587000, 202028) != wantStage || highest != wantHighest {
						t.Fatal("quest stage progress")
					}
					if kind == 1 && !allReady {
						deadline := memmap.Uint64(0x5D4594, 3468)
						if deadline < before+20000 || deadline > after+20000 {
							t.Fatal("partial exit countdown")
						}
					}
					raw := unsafe.Slice(memmap.PtrUint8(0x5D4594, 1567844), 96)
					pending := string(raw[:bytes.IndexByte(raw, 0)])
					if pending != "WorldQuest.map" {
						t.Fatal("pending quest map")
					}
					rows = append(rows, struct {
						Name                              string
						Stage, Highest, Exit, Warp, Buffs uint32
						Switched, Observer                bool
						Pending                           string
					}{name, memmap.Uint32(0x587000, 202028), highest, exit, warp, uint32(b.Buffs), switched, true, pending})
				})
			}
		}
	}
}

func TestCollisionRegistryWorldCoopExitSave(t *testing.T) {
	o := newWorldCollisionOwner(t)
	words := worldCollisionExitGlobals(t, o)
	*words["dword_5d4594_1047520"] = 0
	oldName, oldPortal, oldFlag := saveName1557900, dword_5d4594_1563084, int(*o.globals["savePortal"])
	t.Cleanup(func() {
		saveName1557900 = oldName
		dword_5d4594_1563084 = oldPortal
		legacy.Set_dword_5d4594_1563096(oldFlag)
	})
	defer noxflags.PortTestGameFlags(noxflags.GameModeCoop)()
	a := newObjectXferSimple(t, o.s)
	b := &o.units[1]
	data := o.record(t, 88)
	a.CollideData = data
	t.Cleanup(func() { a.CollideData = nil })
	var rows []struct {
		Name         string
		Pending      bool
		Frame, Delay uint32
		Save         string
		Portal       bool
	}
	for _, frame := range []uint32{0, 123, 0xffffffff} {
		for _, path := range []string{"", "WorldTest.map", "WorldTest.map:Entry"} {
			name := fmt.Sprintf("frame%d/path%q", frame, path)
			t.Run(name, func(t *testing.T) {
				o.reset()
				o.s.SetFrame(frame)
				dword_5d4594_1563080 = false
				dword_5d4594_1563092 = 0
				dword_5d4594_1563088 = 0
				dword_5d4594_1563084 = nil
				dword_5d4594_1563064 = false
				saveName1557900 = "before"
				legacy.Set_dword_5d4594_1563096(0)
				legacy.Set_dword_5d4594_1548524(0)
				clear(unsafe.Slice((*byte)(data), 88))
				copy(unsafe.Slice((*byte)(data), 88), path)
				collisionRegistryWorld(7, a, b, nil)
				want := path != ""
				if dword_5d4594_1563080 != want || dword_5d4594_1563064 != want || (int(*o.globals["savePortal"]) != 0) != want || legacy.Get_dword_5d4594_1548524() != 0 {
					t.Fatal("cooperative save request")
				}
				if want && (dword_5d4594_1563084 != a.CObj() || saveName1557900 != "WORKING" || dword_5d4594_1563088 != frame || dword_5d4594_1563092 != 0) {
					t.Fatal("save portal/frame/name")
				}
				rows = append(rows, struct {
					Name         string
					Pending      bool
					Frame, Delay uint32
					Save         string
					Portal       bool
				}{name, dword_5d4594_1563080, dword_5d4594_1563088, dword_5d4594_1563092, saveName1557900, dword_5d4594_1563084 == a.CObj()})
			})
		}
	}
	collisionRegistryCapture(t, "collision-registry-world-coop-exit-save", rows)
}

func TestCollisionRegistryWorldExitGlyphCleanup(t *testing.T) {
	var rows []struct {
		Name    string
		Traps   byte
		Deleted [3]bool
	}
	defer func() {
		collisionRegistryCapture(t, "collision-registry-world-exit-glyphs", rows)
	}()
	for _, class := range []byte{0, 1, 2} {
		for _, traps := range []byte{0, 1, 2, 3} {
			for _, carried := range []bool{false, true} {
				name := fmt.Sprintf("class%d/traps%d/carried%t", class, traps, carried)
				t.Run(name, func(t *testing.T) {
					o := newWorldCollisionOwner(t)
					words := worldCollisionExitGlobals(t, o)
					*words["dword_5d4594_1047520"] = 0
					dword_5d4594_1563080 = false
					dword_5d4594_1563092 = 0
					dword_5d4594_1563088 = 0
					dword_5d4594_1563064 = false
					defer noxflags.PortTestGameFlags(0)()
					t.Cleanup(o.s.PortTestAttackTypes(1, nil, "Glyph", "OtherEffect"))
					a := newObjectXferSimple(t, o.s)
					b := &o.units[1]
					data := o.record(t, 88)
					a.CollideData = data
					t.Cleanup(func() { a.CollideData = nil; b.InvFirstItem = nil; o.s.Objs.DeletedList = nil })
					copy(unsafe.Slice((*byte)(data), 88), "WorldGlyph.map")
					*(*byte)(unsafe.Add(b.UpdateDataPlayer().Player.C(), 2251)) = class
					*(*byte)(unsafe.Add(b.UpdateData, 244)) = traps
					var items []*server.Object
					for _, typ := range []string{"Glyph", "Glyph", "OtherEffect"} {
						it := o.s.NewObjectByTypeID(typ)
						if it == nil {
							t.Fatal("effect factory")
						}
						items = append(items, it)
						o.s.ObjSetOwner(b, it)
						t.Cleanup(func() { it.InvHolder = nil; o.s.ObjSetOwner(nil, it); o.s.Objs.FreeObject(it) })
					}
					if carried {
						items[0].InvHolder = b
						b.InvFirstItem = items[0]
					}
					collisionRegistryWorld(7, a, b, nil)
					var deleted [3]bool
					removed := byte(0)
					for i, it := range items {
						deleted[i] = it.ObjFlags.Has(object.FlagDestroyed)
						want := class == 1 && i < 2 && !(i == 0 && carried)
						if deleted[i] != want {
							t.Fatal("exit effect cleanup")
						}
						if want {
							removed++
						}
					}
					wantTraps := byte(0)
					if traps > removed {
						wantTraps = traps - removed
					}
					got := *(*byte)(unsafe.Add(b.UpdateData, 244))
					if got != wantTraps {
						t.Fatal("exit trap count saturation")
					}
					rows = append(rows, struct {
						Name    string
						Traps   byte
						Deleted [3]bool
					}{name, got, deleted})
				})
			}
		}
	}
}
