//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"image"
	"reflect"
	"testing"
	"unsafe"

	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/legacy"
)

func TestGameMessageClientDurationEffects(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t, "CharmRay", "DrainManaRay", "HealRay", "HarpoonRope")
	region := serverConfigOwnBytes(t, 0x5D4594, 1303540, 860)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	type state struct {
		Drawables [][]uint32
		Calls     []effectsSpawnCall
		Deleted   []uint32
		Cache     []uint32
	}
	snapshot := func() state {
		cache := append([]uint32(nil), unsafe.Slice((*uint32)(unsafe.Pointer(&region[0])), len(region)/4)...)
		for i, v := range cache[:192] {
			if v != 0 {
				r := c.refs[(*client.Drawable)(unsafe.Pointer(uintptr(v)))]
				if r == 0 {
					t.Fatal("duration cache owner")
				}
				cache[i] = r
			}
		}
		return state{c.snapshotDrawables(t), append([]effectsSpawnCall(nil), c.Calls...), append([]uint32(nil), c.Deleted...), cache}
	}
	type row struct {
		On, Kind, Existing, Fail, Return int
		State                            state
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for kind := 0; kind < 256; kind++ {
			for existing := 0; existing < 2; existing++ {
				for fail := 0; fail < 2; fail++ {
					c.resetCase(env, pix, 1, 100)
					clear(region)
					binary.LittleEndian.PutUint32(connected, uint32(on))
					for i, pos := range []image.Point{image.Pt(100, 200), image.Pt(200, 300)} {
						dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, pos)
						if dr == nil {
							t.Fatal("duration endpoint allocation")
						}
						dr.NetCode32 = uint32(7 + i)
						dr.ObjClass = 0
					}
					if existing != 0 {
						event := [7]byte{158, 1, 17, 7, 0, 8, 0}
						legacy.PortTestPresentationRayAdd(&event)
					}
					c.Calls = nil
					c.Deleted = nil
					c.FailEvery = fail
					before := snapshot()
					data := []byte{158, byte(kind), 255, 7, 0, 8, 0}
					input := bytes.Clone(data)
					n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(158), data)
					valid := kind >= 1 && kind <= 14
					wantN := -1
					if valid {
						wantN = 7
					}
					if n != wantN || !bytes.Equal(data, input) {
						t.Fatal("duration subtype return/input", kind, n)
					}
					wantCount := 2 + existing
					wantCalls, wantDeleted := 0, 0
					if on != 0 && kind >= 1 && kind <= 7 {
						wantCalls = 1
						if fail == 0 {
							wantCount++
						}
					}
					if on != 0 && kind >= 8 && kind <= 14 && existing != 0 {
						wantCount--
						wantDeleted = 1
					}
					if c.Objs.Count != wantCount || len(c.Calls) != wantCalls || len(c.Deleted) != wantDeleted {
						t.Fatalf("duration effects on%d kind%d existing%d fail%d count%d calls%d deleted%d", on, kind, existing, fail, c.Objs.Count, len(c.Calls), len(c.Deleted))
					}
					if on != 0 && kind >= 1 && kind <= 7 && fail == 0 {
						dr := c.Objs.List1
						b := unsafe.Slice((*byte)(unsafe.Add(dr.C(), 432)), 13)
						if b[0] != 1 || binary.LittleEndian.Uint32(b[1:]) != 255 || binary.LittleEndian.Uint32(b[5:]) != 7 || binary.LittleEndian.Uint32(b[9:]) != 8 || dr.PosVec != image.Pt(150, 250) {
							t.Fatal("duration packed fields/midpoint")
						}
					}
					after := snapshot()
					if (!valid || on == 0) && !reflect.DeepEqual(before, after) {
						t.Fatal("duration rejected/gated state changed")
					}
					rows = append(rows, row{on, kind, existing, fail, n, after})
				}
			}
		}
	}
	interactionCapture(t, "game-progress-duration-effects", rows)
}
