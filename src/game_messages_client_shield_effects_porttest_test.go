//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"image"
	"reflect"
	"testing"

	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy"
)

func TestGameMessageClientShieldEffects(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t, drawableShieldNames...)
	words, restore := legacy.PortTestDrawableEffectGlobals()
	t.Cleanup(restore)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	type state struct {
		Drawables [][]uint32
		Calls     []effectsSpawnCall
		Globals   []uint32
	}
	snapshot := func() state {
		g := make([]uint32, len(words))
		for i, p := range words {
			g[i] = *p
		}
		return state{c.snapshotDrawables(t), append([]effectsSpawnCall(nil), c.Calls...), g}
	}
	type row struct {
		On, Present, Static, Fail, Direction, Step int
		State                                      state
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for present := 0; present < 2; present++ {
			for static := 0; static < 2; static++ {
				for fail := 0; fail < 2; fail++ {
					for dir := 0; dir < 9; dir++ {
						c.resetCase(env, pix, 1, 100)
						for _, p := range words {
							*p = 0
						}
						binary.LittleEndian.PutUint32(connected, uint32(on))
						target := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(300, 400))
						if target == nil {
							t.Fatal("shield target allocation")
						}
						target.NetCode32 = 17
						code := uint16(17)
						if present == 0 {
							code = 18
						}
						if static != 0 {
							code |= 0x8000
							target.ObjClass = object.Class(0x400000)
						}
						c.Calls = nil
						c.FailEvery = fail
						before := snapshot()
						for step := 0; step < 2; step++ {
							data := []byte{128, byte(code), byte(code >> 8), byte(dir)}
							input := bytes.Clone(data)
							n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(128), data)
							if n != 4 || !bytes.Equal(data, input) {
								t.Fatal("shield message return/input")
							}
							live := 1
							if on != 0 && present != 0 && fail == 0 && dir != 4 {
								live++
							}
							if c.Objs.Count != live {
								t.Fatalf("shield count on%d present%d static%d fail%d dir%d step%d: %d", on, present, static, fail, dir, step, c.Objs.Count)
							}
							if live == 2 {
								dr := c.Objs.List1
								which := dir
								if dir > 4 {
									which--
								}
								if dr == target || dr.TypeIDVal != uint32(c.Things.IndByID(drawableShieldNames[which])) || dr.PosVec != image.Pt(300, 403) || *txword(dr, 432) != uint32(code) {
									t.Fatal("shield type/position/full owner code")
								}
							}
							after := snapshot()
							if on == 0 && !reflect.DeepEqual(before, after) {
								t.Fatal("disconnected shield changed state")
							}
							if live == 2 && len(c.Calls) != 1 {
								t.Fatal("duplicate shield allocated")
							}
							rows = append(rows, row{on, present, static, fail, dir, step, after})
						}
					}
				}
			}
		}
	}
	interactionCapture(t, "game-progress-shield-effects", rows)
}
