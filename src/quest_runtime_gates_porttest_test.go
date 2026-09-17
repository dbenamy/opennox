//go:build porttest

package opennox

import (
	"fmt"
	"math"
	"testing"
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

func questRuntimeGate(t *testing.T, o *questRuntimeOwner) *server.Object {
	o.s.PortTestAIEmptyMap()
	t.Cleanup(o.s.Map.Free)
	gate := newObjectXferSimple(t, o.s)
	gate.NetCode = 3001
	o.objects[gate.CObj()] = 3001
	gate.ObjClass = object.Class(0x20)
	gate.ObjSubClass = 2
	gate.CollideData = o.record(t, 88)
	t.Cleanup(func() { gate.CollideData = nil })
	*(*types.Pointf)(unsafe.Add(gate.CollideData, 80)) = types.Ptf(175, 200)
	return gate
}

type questGateSnapshot struct {
	Pos                  [2]uint32
	Gate, Camera         bool
	Status, Flags, Buffs uint32
	Sounds               []int
	Queue                legacy.PortTestReliableReportState
}

func questGateState(o *questRuntimeOwner, u *server.Object) questGateSnapshot {
	pl := u.UpdateDataPlayer().Player
	r := questGateSnapshot{Pos: [2]uint32{math.Float32bits(u.PosVec.X), math.Float32bits(u.PosVec.Y)}, Gate: objectXferGetWord(u.UpdateData, 316) != 0, Camera: pl.CameraFollowObj != nil, Status: pl.Field3680, Flags: uint32(u.ObjFlags), Buffs: u.Buffs, Queue: o.state()}
	for _, e := range o.s.PortTestCombatAudioSnapshot() {
		r.Sounds = append(r.Sounds, int(e.ID))
	}
	return r
}
func questGatePlayer(o *questRuntimeOwner, u, gate *server.Object, participation uint32, entered bool) {
	o.reset()
	o.s.PortTestCombatAudioReset()
	pl := u.UpdateDataPlayer().Player
	pl.Field3680 = 0x111
	pl.CameraFollowObj = gate
	u.ObjFlags = 0x40
	u.Buffs = 1
	u.PosVec = types.Ptf(100, 100)
	u.PrevPos = u.PosVec
	u.NewPos = u.PosVec
	objectXferSetWord(pl.C(), 4792, participation)
	objectXferSetWord(u.UpdateData, 312, 0)
	objectXferSetWord(u.UpdateData, 316, 0)
	if entered {
		*(*unsafe.Pointer)(unsafe.Add(u.UpdateData, 316)) = gate.CObj()
	}
}
func TestQuestRuntimeReturnFromGate(t *testing.T) {
	o := newQuestRuntimeOwner(t)
	gate := questRuntimeGate(t, o)
	u := &o.units[0]
	type row struct {
		Name  string
		State questGateSnapshot
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "quest-runtime-gate-return", rows, "5f6ac2c28b48ba9a4063d7af6042787f3c468f81181a385ba9626ca88aed5846")
	}()
	for _, quest := range []bool{false, true} {
		for _, entered := range []bool{false, true} {
			for _, part := range []uint32{0, 1, 2} {
				name := fmt.Sprintf("quest%t/entered%t/participation%d", quest, entered, part)
				t.Run(name, func(t *testing.T) {
					flags := noxflags.GameFlag(0)
					if quest {
						flags = noxflags.GameModeQuest
					}
					restore := noxflags.PortTestGameFlags(flags)
					defer restore()
					questGatePlayer(o, u, gate, part, entered)
					questRuntimeCall("sub_4D7480", u)
					got := questGateState(o, u)
					want := types.Ptf(100, 100)
					if entered {
						want = types.Ptf(175, 200)
					}
					if u.PosVec != want || got.Gate {
						t.Fatal("gate return position/reference", u.PosVec, got.Gate)
					}
					if entered && (got.Camera || got.Status&0x121 != 0 || got.Flags&0x40 != 0 || got.Buffs&1 != 0 || len(got.Sounds) != 2 || got.Sounds[0] != 0 || got.Sounds[1] != 312) {
						t.Fatal("observer leave state", got)
					}
					if !entered && (got.Status != 0x111 || len(got.Sounds) != 0) {
						t.Fatal("missing gate changed observer")
					}
					rows = append(rows, row{name, got})
				})
			}
		}
	}
	questRuntimeCall("sub_4D7480", nil)
}
func TestQuestRuntimeCloseGate(t *testing.T) {
	o := newQuestRuntimeOwner(t)
	gate := questRuntimeGate(t, o)
	u := &o.units[0]
	oldList := o.s.Objs.First()
	o.s.Objs.SetObjects(gate)
	t.Cleanup(func() { o.s.Objs.SetObjects(oldList) })
	defer noxflags.PortTestGameFlags(noxflags.GameModeQuest)()
	type row struct {
		Name        string
		Return      uint64
		Open, Flags uint32
		Player      questGateSnapshot
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "quest-runtime-gate-close", rows, "11c570b3001b2fc30529a3df452966850f2b19801fc4890dfc60b91a986e6dc1")
	}()
	for _, previous := range []uint32{0, 1, 2} {
		for _, next := range []uint32{0, 1, 2, 255} {
			for _, part := range []uint32{0, 1, 2} {
				for _, entered := range []bool{false, true} {
					name := fmt.Sprintf("previous%d/next%d/participation%d/entered%t", previous, next, part, entered)
					t.Run(name, func(t *testing.T) {
						questGatePlayer(o, u, gate, part, entered)
						*o.quest["1556120"] = previous
						gate.ObjFlags = 0x1000000
						rv := questRuntimeCall("sub_4D7520", nil, next)
						got := questGateState(o, u)
						close := previous == 1 && next == 0
						moved := close && part != 0 && entered
						if questRuntimeCall("sub_4D75E0", nil) != uint64(next) || *o.quest["1556120"] != next || (gate.ObjFlags&0x1000000 == 0) != close || got.Gate != (entered && !moved) {
							t.Fatal("gate transition", close, moved, got.Gate, gate.ObjFlags)
						}
						want := types.Ptf(100, 100)
						if moved {
							want = types.Ptf(175, 200)
						}
						if u.PosVec != want {
							t.Fatal("gate return scope")
						}
						rows = append(rows, row{name, rv, *o.quest["1556120"], uint32(gate.ObjFlags), got})
					})
				}
			}
		}
	}
}
