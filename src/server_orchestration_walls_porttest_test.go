//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/sound"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
	"unsafe"
)

func TestServerOrchestrationSecretWalls(t *testing.T) {
	o := newCollisionCoreOwner(t)
	o.s.PortTestAIEmptyMap()
	t.Cleanup(o.s.Map.Free)
	t.Cleanup(o.s.PortTestMinimapWalls())
	head, restore := legacy.PortTestPrefabSecretList()
	t.Cleanup(restore)
	t.Cleanup(o.s.PortTestCombatAudioReset)
	wl := o.s.Walls.CreateAtGrid(image.Pt(6, 4))
	if wl == nil {
		t.Fatal("wall allocation")
	}
	node := o.record(t, 32)
	raw := unsafe.Slice((*byte)(node), 32)
	*head = uint32(uintptr(node))
	binary.LittleEndian.PutUint32(raw[4:], 6)
	binary.LittleEndian.PutUint32(raw[8:], 4)
	binary.LittleEndian.PutUint32(raw[12:], uint32(uintptr(wl.C())))
	u := newObjectXferSimple(t, o.s)
	worldGeometryResetObject(u, 1001, 149.5, 103.5, false)
	u.ObjClass = 0
	u.ObjFlags = 4
	u.Collide = o.callback
	u.Update = nil
	o.s.Map.AddObjectToIndex(u)
	ids := collisionCoreIDs(u)
	type row struct {
		Phase, Flags, Progress, AfterPhase, AfterProgress byte
		Timer, Duration, AfterTimer                       uint32
		Active                                            bool
		Queues                                            [3]uint32
		Sound                                             sound.ID
		SoundPos                                          types.Pointf
	}
	var rows []row
	for _, phase := range []byte{0, 1, 2, 3, 4, 255} {
		for _, flags := range []byte{0, 4, 8, 12} {
			for _, progress := range []byte{0, 1, 22, 23, 255} {
				for _, timer := range []uint32{0, 1, 2} {
					for _, duration := range []uint32{0, 3} {
						name := fmt.Sprintf("phase%d/flags%d/progress%d/timer%d/duration%d", phase, flags, progress, timer, duration)
						t.Run(name, func(t *testing.T) {
							o.resetQueues()
							u.Field115 = 0
							u.Field116 = 0
							o.s.PortTestCombatAudioReset()
							o.s.SetTickRate(30)
							raw[20], raw[21], raw[22] = flags, phase, progress
							binary.LittleEndian.PutUint32(raw[16:], duration)
							binary.LittleEndian.PutUint32(raw[24:], timer)
							legacy.PortTestServerOrchestration("walls", nil, 0)
							wp, wx, wt := phase, progress, timer
							active := false
							var snd sound.ID
							switch phase {
							case 1:
								if flags&12 == 12 {
									wt--
									if wt == 0 {
										wp = 4
										snd = sound.SoundSecretWallOpen
									}
								}
							case 2:
								wx--
								active = true
								if wx == 0 {
									wp = 1
									wt = 30 * duration
								}
							case 3:
								if flags&4 != 0 && flags&8 == 0 {
									wt--
									if wt == 0 {
										wp = 2
										snd = sound.SoundSecretWallClose
									}
								}
							case 4:
								wx++
								active = true
								if wx == 23 {
									wp = 3
									wt = 30 * duration
								}
							}
							r := row{Phase: phase, Flags: flags, Progress: progress, AfterPhase: raw[21], AfterProgress: raw[22], Timer: timer, Duration: duration, AfterTimer: binary.LittleEndian.Uint32(raw[24:]), Active: u.Field116&1 != 0, Queues: o.queues(ids)}
							sounds := o.s.PortTestCombatAudioSnapshot()
							if snd != 0 {
								if len(sounds) != 1 {
									t.Fatal("sound count", sounds)
								}
								r.Sound = sounds[0].ID
								r.SoundPos = sounds[0].Pos
								if !sounds[0].ByPos || r.SoundPos != (types.Pointf{X: 149, Y: 103}) {
									t.Fatal("sound position", sounds)
								}
							} else if len(sounds) != 0 {
								t.Fatal("unexpected sound")
							}
							if r.AfterPhase != wp || r.AfterProgress != wx || r.AfterTimer != wt || r.Active != active || r.Sound != snd || raw[20] != flags {
								t.Fatal("wall transition", r, wp, wx, wt, active, snd)
							}
							rows = append(rows, r)
						})
					}
				}
			}
		}
	}
	// Unknown phases preserve the previous record's scan flag in the existing C loop.
	// Two distant indexed objects distinguish that observable ordering from a global scan.
	wl2 := o.s.Walls.CreateAtGrid(image.Pt(12, 4))
	if wl2 == nil {
		t.Fatal("second wall")
	}
	second := o.record(t, 32)
	r2 := unsafe.Slice((*byte)(second), 32)
	binary.LittleEndian.PutUint32(raw, uint32(uintptr(second)))
	binary.LittleEndian.PutUint32(r2[4:], 12)
	binary.LittleEndian.PutUint32(r2[8:], 4)
	binary.LittleEndian.PutUint32(r2[12:], uint32(uintptr(wl2.C())))
	v := newObjectXferSimple(t, o.s)
	worldGeometryResetObject(v, 1002, 287.5, 103.5, false)
	v.ObjClass = 0
	v.ObjFlags = 4
	v.Collide = o.callback
	v.Update = nil
	o.s.Map.AddObjectToIndex(v)
	for _, first := range []byte{0, 1, 2, 3, 4} {
		o.resetQueues()
		u.Field115 = 0
		u.Field116 = 0
		v.Field115 = 0
		v.Field116 = 0
		raw[20], raw[21], raw[22] = 0, first, 10
		r2[21] = 0
		legacy.PortTestServerOrchestration("walls", nil, 0)
		want := first == 2 || first == 4
		if (v.Field116&1 != 0) != want {
			t.Fatal("sequential scan carry", first, v.Field116)
		}
		rows = append(rows, row{Phase: first, Active: v.Field116&1 != 0, AfterPhase: r2[21]})
	}
	spellbookCapture(t, "server-orchestration-secret-walls", rows, "a3f5a14c8ee0a72b5764ee7e96308824cd461dc55bb2e50a88f1ded0e25f3607")
}
