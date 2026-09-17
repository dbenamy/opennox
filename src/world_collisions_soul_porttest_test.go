//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/object"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestWorldCollisionsSoulGate(t *testing.T) {
	o := newWorldCollisionOwner(t)
	a := newObjectXferSimple(t, o.s)
	b := &o.units[1]
	data := o.record(t, 4)
	a.CollideData = data
	t.Cleanup(func() { a.CollideData = nil })
	var rows []struct {
		Name                       string
		Frame, Started, Open, Gate uint32
		Sounds                     []int
		Bytes                      []byte
	}
	defer func() {
		spellbookCapture(t, "world-collisions-soul-gate", rows, "b4fe0e9b7260639b8c3fd4e37f697a3a3b525b8456fd194609d5a3e8a25b5b15")
	}()
	for _, frame := range []uint32{0, 30, 31, 123, 0xffffffff} {
		for _, previous := range []uint32{0, 1, 93, 0xfffffff0} {
			for _, same := range []bool{false, true} {
				for _, participation := range []uint32{0, 1, 2} {
					for _, quest := range []bool{false, true} {
						for _, player := range []bool{false, true} {
							name := fmt.Sprintf("frame%d/previous%d/same%t/participation%d/quest%t/player%t", frame, previous, same, participation, quest, player)
							t.Run(name, func(t *testing.T) {
								o.reset()
								o.s.SetFrame(frame)
								o.s.PortTestCombatAudioReset()
								flags := noxflags.GameFlag(0)
								if quest {
									flags = noxflags.GameModeQuest
								}
								defer noxflags.PortTestGameFlags(flags)()
								*o.globals["soulFrame"] = 42
								*o.globals["warpOpen"] = 2
								for i := range o.units {
									u := &o.units[i]
									u.ObjClass = object.ClassPlayer
									objectXferSetWord(u.UpdateData, 308, 0)
									objectXferSetWord(u.UpdateDataPlayer().Player.C(), 4792, participation)
								}
								// Another participant has already entered a gate.
								objectXferSetWord(o.units[0].UpdateData, 308, uint32(uintptr(a.CObj())))
								b.ObjClass = object.ClassSimple
								if player {
									b.ObjClass = object.ClassPlayer
								}
								if same {
									objectXferSetWord(b.UpdateData, 308, uint32(uintptr(a.CObj())))
								}
								objectXferSetWord(data, 0, previous)
								legacy.PortTestWorldCollision(19, a, b, nil)
								admitted := quest && player
								fired := admitted && (!same || frame-previous > 30)
								wantFrame, wantStarted, wantOpen := previous, uint32(42), uint32(2)
								if admitted {
									wantFrame = frame
									wantOpen = 0
									if participation != 1 {
										wantStarted = frame
									}
								}
								gate := objectXferGetWord(b.UpdateData, 308)
								wantGate := uint32(0)
								if same || admitted {
									wantGate = uint32(uintptr(a.CObj()))
								}
								if objectXferGetWord(data, 0) != wantFrame || *o.globals["soulFrame"] != wantStarted || *o.globals["warpOpen"] != wantOpen || gate != wantGate {
									t.Fatal("soul gate state")
								}
								var sounds []int
								for _, e := range o.s.PortTestCombatAudioSnapshot() {
									if e.Obj != a || e.ID != 1005 {
										t.Fatal("soul gate audio routing")
									}
									sounds = append(sounds, int(e.ID))
								}
								packet := o.s.NetList.CopyPacketsA(7, netlist.Kind1)
								if (len(sounds) == 1) != fired || len(sounds) > 1 || (len(packet) > 0) != fired {
									t.Fatalf("sounds%v packet%d fired%t", sounds, len(packet), fired)
								}
								if gate != 0 {
									gate = 1
								}
								rows = append(rows, struct {
									Name                       string
									Frame, Started, Open, Gate uint32
									Sounds                     []int
									Bytes                      []byte
								}{name, objectXferGetWord(data, 0), *o.globals["soulFrame"], *o.globals["warpOpen"], gate, sounds, packet})
							})
						}
					}
				}
			}
		}
	}
}
