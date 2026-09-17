//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/object"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
	"unsafe"
)

func TestWorldCollisionsClockThrottle(t *testing.T) {
	o := newWorldCollisionOwner(t)
	a, b := &o.units[0], &o.units[1]
	oldData := a.UpdateData
	a.UpdateData = o.record(t, 64)
	t.Cleanup(func() { a.UpdateData = oldData })
	b.UpdateDataPlayer().Player.PlayerInd = 7
	var rows []struct {
		Name     string
		Throttle uint64
		Sounds   []int
		Bytes    []byte
	}
	defer func() {
		spellbookCapture(t, "world-collisions-clock-throttle", rows, "d58bc83c541e2bf0a765aa9d100848a498590eef2ccf4b424db812c526efa55b")
	}()
	for _, op := range []int{1, 9} {
		for _, ticks := range []uint64{0, 1500, 1501, 10000, 0xffffffff, 0x100000001} {
			for _, last := range []uint64{0, 1, 8500, 8501, 8499, 0xffffffff, 0x100000001} {
				for _, gate := range []bool{false, true} {
					name := fmt.Sprintf("op%d/ticks%d/last%d/gate%t", op, ticks, last, gate)
					t.Run(name, func(t *testing.T) {
						restore := noxflags.PortTestGameFlags(noxflags.GameModeQuest)
						defer restore()
						o.reset()
						o.s.PortTestCombatAudioReset()
						o.ticks = ticks
						*o.throttle = last
						a.ObjFlags = 0
						a.ObjSubClass = 0
						if gate {
							a.ObjSubClass = 4
						}
						b.ObjClass = object.ClassPlayer
						b.InvFirstItem = nil
						objectXferSetWord(a.CObj(), 508, 0)
						*(*byte)(unsafe.Add(a.UpdateData, 1)) = 5
						if op == 9 {
							a.ObjSubClass = 0x100
						}
						legacy.PortTestWorldCollision(op, a, b, nil)
						var sounds []int
						for _, e := range o.s.PortTestCombatAudioSnapshot() {
							sounds = append(sounds, int(e.ID))
						}
						packet := o.s.NetList.CopyPacketsA(7, netlist.Kind1)
						now := uint64(uint32(ticks))
						fired := now-last > 1500
						wantTick := last
						if fired {
							wantTick = now
						}
						if *o.throttle != wantTick || (len(sounds) == 1) != fired || len(sounds) > 1 || (len(packet) > 0) != fired {
							t.Fatalf("throttle%d sounds%v bytes%d want%t", *o.throttle, sounds, len(packet), fired)
						}
						if fired {
							wantSound := 240
							if gate {
								wantSound = 244
							}
							if op == 9 {
								wantSound = 1012
							}
							if sounds[0] != wantSound {
								t.Fatal("message sound")
							}
						}
						rows = append(rows, struct {
							Name     string
							Throttle uint64
							Sounds   []int
							Bytes    []byte
						}{name, *o.throttle, sounds, packet})
					})
				}
			}
		}
	}
}
