//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"math"
	"testing"
)

func TestWorldCollisionsTeleport(t *testing.T) {
	o := newWorldCollisionOwner(t)
	o.s.PortTestAIEmptyMap()
	t.Cleanup(o.s.Map.Free)
	a := newObjectXferSimple(t, o.s)
	b := &o.units[1]
	other := newObjectXferSimple(t, o.s)
	data := o.record(t, 8)
	a.CollideData = data
	t.Cleanup(func() { a.CollideData = nil })
	objectXferSetWord(data, 0, 150)
	objectXferSetWord(data, 4, 175)
	var rows []struct {
		Name             string
		Pos, Destination [2]uint32
		Sounds           []int
	}
	defer func() {
		spellbookCapture(t, "world-collisions-teleport", rows, "d9b02b0d050a9c7932178a62c0291b889da1aec57e94a937cf47b98da8decee9")
	}()
	for _, class := range []object.Class{object.ClassPlayer, object.ClassSimple, object.ClassDoor} {
		for _, coop := range []bool{false, true} {
			for _, blocked := range []bool{false, true} {
				for _, registered := range []bool{false, true} {
					name := fmt.Sprintf("class%d/coop%t/blocked%t/registered%t", class, coop, blocked, registered)
					t.Run(name, func(t *testing.T) {
						o.reset()
						o.s.PortTestCombatAudioReset()
						flags := noxflags.GameFlag(0)
						if coop {
							flags = noxflags.GameModeCoop
						}
						defer noxflags.PortTestGameFlags(flags)()
						b = &o.units[1]
						if class != object.ClassPlayer {
							b = other
						}
						b.ObjClass = class
						b.ObjFlags = 0
						b.Buffs = 0
						if blocked {
							b.Buffs = 1 << 14
						}
						b.PosVec = types.Ptf(100, 100)
						b.PrevPos = b.PosVec
						b.NewPos = b.PosVec
						b.Field41 = 0
						b.Field42 = 0
						if registered {
							a.Collide, _ = server.PortTestWorldCollisionRegistry("TeleportCollide")
							a.CallCollide(int(uintptr(b.CObj())), 0)
						} else {
							legacy.PortTestWorldCollision(15, a, b, nil)
						}
						admitted := class != object.ClassDoor
						moved := admitted && !blocked && (coop || class == object.ClassPlayer)
						want := types.Ptf(100, 100)
						if moved {
							want = types.Ptf(150, 175)
						}
						if b.PosVec != want {
							t.Fatalf("position %v want %v", b.PosVec, want)
						}
						dest := [2]uint32{b.Field41, b.Field42}
						wantDest := [2]uint32{}
						if admitted {
							wantDest = [2]uint32{math.Float32bits(150), math.Float32bits(175)}
						}
						if dest != wantDest {
							t.Fatal("teleport destination")
						}
						var sounds []int
						for _, e := range o.s.PortTestCombatAudioSnapshot() {
							if e.Obj != b || e.ID != 147 {
								t.Fatal("teleport sound routing")
							}
							sounds = append(sounds, int(e.ID))
						}
						wantSounds := 0
						if admitted {
							wantSounds = 2
						}
						if len(sounds) != wantSounds {
							t.Fatal("teleport sound count")
						}
						rows = append(rows, struct {
							Name             string
							Pos, Destination [2]uint32
							Sounds           []int
						}{name, [2]uint32{math.Float32bits(b.PosVec.X), math.Float32bits(b.PosVec.Y)}, dest, sounds})
					})
				}
			}
		}
	}
	legacy.PortTestWorldCollision(15, a, nil, nil)
}
