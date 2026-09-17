//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestWorldCollisionsSign(t *testing.T) {
	o := newWorldCollisionOwner(t)
	o.s.PortTestAIEmptyMap()
	t.Cleanup(o.s.Map.Free)
	a := newObjectXferSimple(t, o.s)
	b := &o.units[1]
	a.Use, _ = server.PortTestWorldUseRegistry("ReadUse")
	data := o.record(t, 260)
	a.UseData.Ptr = data
	t.Cleanup(func() { a.UseData.Ptr = nil })
	copy(unsafe.Slice((*byte)(data), 256), "world-collision-sign")
	a.PosVec = types.Ptf(100, 100)
	b.PosVec = types.Ptf(101, 100)
	var rows []struct {
		Name  string
		Frame uint32
		Bytes []byte
	}
	defer func() {
		spellbookCapture(t, "world-collisions-sign", rows, "c16ee907299d1075a6fd97f4dfc10546d88ffaf9670c78071834706bc1dc2de6")
	}()
	for _, frame := range []uint32{0, 1, 90, 91, 123, 0xffffffff} {
		for _, previous := range []uint32{0, 1, 33, 0xfffffff0} {
			for _, player := range []bool{false, true} {
				for _, registered := range []bool{false, true} {
					name := fmt.Sprintf("frame%d/previous%d/player%t/registered%t", frame, previous, player, registered)
					t.Run(name, func(t *testing.T) {
						o.reset()
						o.s.SetFrame(frame)
						objectXferSetWord(data, 256, previous)
						b.ObjClass = object.ClassSimple
						if player {
							b.ObjClass = object.ClassPlayer
						}
						if registered {
							a.Collide, _ = server.PortTestWorldCollisionRegistry("SignCollide")
							a.CallCollide(int(uintptr(b.CObj())), 0)
						} else {
							legacy.PortTestWorldCollision(13, a, b, nil)
						}
						fired := player && (previous == 0 || frame-previous > 90)
						want := previous
						if fired {
							want = frame
						}
						packet := o.s.NetList.CopyPacketsA(7, netlist.Kind1)
						if objectXferGetWord(data, 256) != want || (len(packet) > 0) != fired {
							t.Fatalf("frame%d bytes%d fired%t", objectXferGetWord(data, 256), len(packet), fired)
						}
						rows = append(rows, struct {
							Name  string
							Frame uint32
							Bytes []byte
						}{name, objectXferGetWord(data, 256), packet})
					})
				}
			}
		}
	}
	legacy.PortTestWorldCollision(13, a, nil, nil)
}
