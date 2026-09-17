//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
)

func TestWorldCollisionsAudioFrames(t *testing.T) {
	o := newWorldCollisionOwner(t)
	a, b := &o.units[0], &o.units[1]
	data := o.record(t, 4)
	objectXferSetWord(data, 0, 281)
	a.CollideData = data
	var rows []struct {
		Name   string
		Frame  uint32
		Sounds []int
	}
	defer func() {
		spellbookCapture(t, "world-collisions-audio-frames", rows, "d5ee2d5f4a469647a4ee0b2de44cc21eec3f945125d0dbff81218882278a4e0c")
	}()
	for _, op := range []int{10, 11} {
		for _, frame := range []uint32{0, 1, 3, 4, 30, 31, 123, 0xfffffffe, 0xffffffff} {
			for _, previous := range []uint32{0, 1, 3, 30, 100, 0xffffffe0, 0xfffffffe} {
				for _, player := range []bool{false, true} {
					for _, registered := range []bool{false, true} {
						name := fmt.Sprintf("op%d/frame%d/previous%d/player%t/registered%t", op, frame, previous, player, registered)
						t.Run(name, func(t *testing.T) {
							o.s.SetFrame(frame)
							a.Field34 = previous
							b.ObjClass = object.ClassMonster
							if player {
								b.ObjClass = object.ClassPlayer
							}
							o.s.PortTestCombatAudioReset()
							if registered {
								name := "BarrelCollide"
								if op == 11 {
									name = "AudioEventCollide"
								}
								a.Collide, _ = server.PortTestWorldCollisionRegistry(name)
								a.CallCollide(int(uintptr(b.CObj())), 0)
							} else {
								legacy.PortTestWorldCollision(op, a, b, nil)
							}
							delta := uint32(3)
							if op == 11 {
								delta = 30
							}
							fired := frame > previous+delta && (op == 10 || player)
							wantFrame := previous
							if fired {
								wantFrame = frame
							}
							if a.Field34 != wantFrame {
								t.Fatal("audio throttle frame")
							}
							var sounds []int
							for _, e := range o.s.PortTestCombatAudioSnapshot() {
								if e.Obj != a || e.Kind != 0 || e.Code != 0 || e.ByPos {
									t.Fatal("audio routing")
								}
								sounds = append(sounds, int(e.ID))
							}
							if (len(sounds) == 1) != fired || len(sounds) > 1 {
								t.Fatal("audio event count")
							}
							if fired && sounds[0] != 281 {
								t.Fatal("sound id")
							}
							rows = append(rows, struct {
								Name   string
								Frame  uint32
								Sounds []int
							}{name, a.Field34, sounds})
						})
					}
				}
			}
		}
	}
}

func TestWorldCollisionsPentagram(t *testing.T) {
	o := newWorldCollisionOwner(t)
	a := &o.units[0]
	var rows [][3]uint32
	for _, value := range []uint32{0, 1, 0xffffffff} {
		for _, registered := range []bool{false, true} {
			objectXferSetWord(a.UpdateData, 4, value)
			before := objectXferGetWord(a.UpdateData, 0)
			rv := uint32(1)
			if registered {
				a.Collide, _ = server.PortTestWorldCollisionRegistry("PentagramCollide")
				a.CallCollide(0, 0)
			} else {
				rv = legacy.PortTestWorldCollision(12, a, nil, nil)
			}
			if rv != 1 || objectXferGetWord(a.UpdateData, 4) != 1 || objectXferGetWord(a.UpdateData, 0) != before {
				t.Fatal("pentagram activation")
			}
			rows = append(rows, [3]uint32{value, rv, objectXferGetWord(a.UpdateData, 4)})
		}
	}
	spellbookCapture(t, "world-collisions-pentagram", rows, "c2e746a185434d3dffcf7990cd5cf416c6b9094abdb9642e6c21fc906eae870d")
}
