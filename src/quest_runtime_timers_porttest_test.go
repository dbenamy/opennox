//go:build porttest

package opennox

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/opennox/libs/types"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
)

func TestQuestRuntimeObserverDeadline(t *testing.T) {
	o := newQuestRuntimeOwner(t)
	gate := questRuntimeGate(t, o)
	u := &o.units[0]
	pl := u.UpdateDataPlayer().Player
	defer noxflags.PortTestGameFlags(noxflags.GameModeQuest)()
	type row struct {
		Name   string
		Return uint64
		Player questGateSnapshot
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "quest-runtime-observer-deadline", rows, "563b92333c85d29e67859ded69a5a58b512f57c4dbbf46c97a5e2644535b1c1c")
	}()
	for _, deadline := range []uint32{0, 100, 0x80000000, 0xffffffff} {
		for _, frame := range []uint32{0, 100, 101, 0x80000000, 0xffffffff} {
			for _, status := range []uint32{0, 1, 16, 17, 0x111} {
				for _, part := range []uint32{0, 1, 2} {
					for _, waiting := range []int{0, 1, 2} {
						name := fmt.Sprintf("deadline%x/frame%x/status%x/part%d/wait%d", deadline, frame, status, part, waiting)
						t.Run(name, func(t *testing.T) {
							questGatePlayer(o, u, gate, part, false)
							pl.Field3680 = status
							*o.quest["observerDeadline"] = deadline
							o.s.SetFrame(frame)
							if waiting == 1 {
								*(*unsafe.Pointer)(unsafe.Add(u.UpdateData, 312)) = gate.CObj()
							}
							if waiting == 2 {
								*(*unsafe.Pointer)(unsafe.Add(u.UpdateData, 316)) = gate.CObj()
							}
							rv := questRuntimeCall("sub_4D7150", nil)
							got := questGateState(o, u)
							released := deadline != 0 && frame > deadline && status&17 == 17 && part == 1 && waiting == 0
							want := status
							if released {
								want &^= 0x121
							}
							if got.Status != want || got.Camera == released || (got.Flags&0x40 == 0) != released {
								t.Fatal("timed observer release", released, got.Status, want, got.Camera, got.Flags)
							}
							rows = append(rows, row{name, rv, got})
						})
					}
				}
			}
		}
	}
}
func TestQuestRuntimeSoulTimeout(t *testing.T) {
	o := newQuestRuntimeOwner(t)
	gate := questRuntimeGate(t, o)
	defer noxflags.PortTestGameFlags(noxflags.GameModeQuest)()
	old := o.s.Doors.Sub_4D72C0()
	t.Cleanup(func() { o.s.Doors.Sub_4D72B0(old) })
	type row struct {
		Name   string
		Return uint64
		Frame  uint32
		Shared bool
		Queue  legacy.PortTestReliableReportState
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "quest-runtime-soul-timeout", rows, "7820e2bcfc13e3245969253c894bc1929e6c8d1a24f058b748c75173402e8361")
	}()
	for count := 0; count <= 3; count++ {
		for _, soul := range []bool{false, true} {
			for _, shared := range []bool{false, true} {
				for _, start := range []uint32{0, 100, 0xfffffff0} {
					for _, age := range []uint32{0, 8999, 9000, 9001, 0xffffffff} {
						name := fmt.Sprintf("count%d/soul%t/shared%t/start%x/age%x", count, soul, shared, start, age)
						t.Run(name, func(t *testing.T) {
							o.reset()
							o.s.SetFrame(start + age)
							*o.globals["soulFrame"] = start
							o.s.Doors.Sub_4D72B0(shared)
							for i := range o.units {
								p := uint32(0)
								if i < count {
									p = 1
								}
								objectXferSetWord(o.units[i].UpdateDataPlayer().Player.C(), 4792, p)
								objectXferSetWord(o.units[i].UpdateData, 308, 0)
							}
							if soul {
								*(*unsafe.Pointer)(unsafe.Add(o.units[2].UpdateData, 308)) = gate.CObj()
							}
							rv := questRuntimeCall("sub_4D71F0", nil)
							triggered := start != 0 && age >= 9000 && soul && count > 1
							wantFrame := start
							if triggered {
								wantFrame = 0
							}
							gotShared := o.s.Doors.Sub_4D72C0()
							st := o.state()
							if *o.globals["soulFrame"] != wantFrame || gotShared != (shared || triggered) {
								t.Fatal("soul deadline/state")
							}
							packets := 0
							if triggered && !shared {
								packets = 1
							}
							if len(st.Nodes) != packets {
								t.Fatal("soul notification count", len(st.Nodes), packets)
							}
							if packets == 1 {
								data := st.Nodes[0].Data
								if len(data) != 3 || data[0] != 240 || data[1] != 24 || data[2] != 1 {
									t.Fatal("soul notification payload")
								}
							}
							rows = append(rows, row{name, rv, *o.globals["soulFrame"], gotShared, st})
						})
					}
				}
			}
		}
	}
}
func TestQuestRuntimeWarpAdmission(t *testing.T) {
	o := newQuestRuntimeOwner(t)
	gate := questRuntimeGate(t, o)
	oldAllow, oldInfinite, oldInc := questAllowDefault, questLevelWarpInfinite, questLevelWarpInc
	t.Cleanup(func() { questAllowDefault, questLevelWarpInfinite, questLevelWarpInc = oldAllow, oldInfinite, oldInc })
	questLevelWarpInfinite = false
	questLevelWarpInc = 5
	defer noxflags.PortTestGameFlags(noxflags.GameModeQuest)()
	type row struct {
		Name    string
		Players [3]questGateSnapshot
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "quest-runtime-warp-admission", rows, "3f2b5fb233b028d6fff5559886bbc2ecf8c27773a5249a02de7b27e1799d7526")
	}()
	for count := 0; count <= 3; count++ {
		for mask := 0; mask < 8; mask++ {
			for _, age := range []uint32{0, 29, 30, 31, 0xffffffff} {
				for _, allow := range []bool{false, true} {
					name := fmt.Sprintf("count%d/mask%d/age%x/allow%t", count, mask, age, allow)
					t.Run(name, func(t *testing.T) {
						for i := range o.units {
							part := uint32(0)
							if i < count {
								part = 1
							}
							questGatePlayer(o, &o.units[i], gate, part, mask&(1<<uint(i)) != 0)
							objectXferSetWord(o.units[i].UpdateDataPlayer().Player.C(), 4696, 0)
						}
						*o.quest["202028"] = 0
						*o.quest["1556108"] = 100
						o.s.SetFrame(100 + age)
						questAllowDefault = allow
						questRuntimeCall("nox_server_checkWarpGate_4D7600", nil)
						all := count > 0
						for i := 0; i < count; i++ {
							if mask&(1<<uint(i)) == 0 {
								all = false
							}
						}
						rejected := all && age >= 30 && !allow
						var got row
						got.Name = name
						for i := range o.units {
							moved := rejected && i < count
							got.Players[i] = questGateState(o, &o.units[i])
							want := types.Ptf(100, 100)
							if moved {
								want = types.Ptf(175, 200)
							}
							if o.units[i].PosVec != want || got.Players[i].Gate != (mask&(1<<uint(i)) != 0 && !moved) {
								t.Fatal("warp admission/return scope", i, rejected, got.Players[i])
							}
						}
						rows = append(rows, got)
					})
				}
			}
		}
	}
}
