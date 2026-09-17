//go:build porttest

package opennox

import (
	"encoding/binary"
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	"github.com/opennox/libs/wall"
	"github.com/opennox/opennox/v1/common/sound"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"image"
	"testing"
	"unsafe"
)

func TestCollisionCoreWallOpening(t *testing.T) {
	o := newCollisionCoreOwner(t)
	o.s.PortTestAIEmptyMap()
	t.Cleanup(o.s.Map.Free)
	configure, guard, restore := o.s.PortTestGeometryWalls()
	t.Cleanup(restore)
	u := newObjectXferSimple(t, o.s)
	data := unsafe.Slice((*byte)(collisionCoreGuarded(t, o, 32)), 32)
	type row struct {
		Present, Secret bool
		Class           uint32
		State, Options  byte
		Data            []byte
		Sounds          []server.PortTestCombatAudioEvent
	}
	var rows []row
	for _, present := range []bool{false, true} {
		for _, secret := range []bool{false, true} {
			for _, class := range []uint32{8, 2, 4, 6, 0x2000} {
				for _, state := range []byte{0, 1, 4} {
					for _, options := range []byte{0, 1, 2, 3} {
						dir := -1
						if present {
							dir = 0
						}
						configure(dir, false, false)
						o.s.PortTestCombatAudioReset()
						clear(data)
						binary.LittleEndian.PutUint32(data[4:], 6)
						binary.LittleEndian.PutUint32(data[8:], 4)
						data[20] = options
						data[21] = state
						data[22] = 7
						if wl := o.s.Walls.GetWallAtGrid(image.Pt(6, 4)); wl != nil {
							wl.Data = unsafe.Pointer(&data[0])
							if secret {
								wl.Flags4 |= wall.FlagSecret
							}
						}
						u.ObjClass = object.Class(class)
						grid := [2]int32{6, 4}
						legacy.PortTestCollisionCoreWallOpen(&grid, u)
						opened := present && secret && class&6 != 0 && state == 1 && options&2 != 0
						sounds := o.s.PortTestCombatAudioSnapshot()
						if opened {
							if data[21] != 4 || data[22] != 0 || len(sounds) != 1 {
								t.Fatal("secret wall transition", data, sounds)
							}
							e := sounds[0]
							if e.ID != sound.SoundSecretWallOpen || !e.ByPos || e.Obj != nil || e.Pos != (types.Pointf{149.5, 103.5}) || e.Kind != 0 || e.Code != 0 {
								t.Fatal("secret wall sound", e)
							}
							legacy.PortTestCollisionCoreWallOpen(&grid, u)
							if len(o.s.PortTestCombatAudioSnapshot()) != 1 {
								t.Fatal("repeated wall opening played sound twice")
							}
						} else if data[21] != state || data[22] != 7 || len(sounds) != 0 {
							t.Fatal("guarded wall contact changed state", present, secret, class, state, options)
						}
						if !guard() {
							t.Fatal("wall record guards")
						}
						rows = append(rows, row{present, secret, class, state, options, append([]byte(nil), data...), sounds})
					}
				}
			}
		}
	}
	u.ObjClass = object.ClassSimple
	spellbookCapture(t, "collision-core-wall-opening", rows, "ace614179a7ee41ca22630a3bc6064f700c0bf9e61c1176781d3fdb9b2454397")
}
