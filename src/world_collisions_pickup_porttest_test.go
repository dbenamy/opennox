//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestWorldCollisionsPickupAdmission(t *testing.T) {
	o := newWorldCollisionOwner(t)
	a, b := &o.units[0], &o.units[1]
	b.CarryCapacity = 0
	var rows []struct {
		Name   string
		Return uint32
	}
	defer func() {
		spellbookCapture(t, "world-collisions-pickup-admission", rows, "e2788e5ba8ace0abab722db614e72368ba8e179660c943b9524d26ff8b7bc3fe")
	}()
	for _, frame := range []uint32{0, 1, 14, 15, 16, 0xffffffff} {
		for _, previous := range []uint32{0, 1, 2, 0xfffffff8} {
			for _, fps := range []uint32{0, 1, 30, 60, 0x80000000, 0xffffffff} {
				for _, class := range []object.Class{object.ClassMonster, object.ClassPlayer, object.ClassSimple} {
					for _, flags := range []byte{0, 1, 2, 255} {
						name := fmt.Sprintf("frame%d/previous%d/fps%d/class%x/flags%d", frame, previous, fps, class, flags)
						t.Run(name, func(t *testing.T) {
							o.s.SetFrame(frame)
							o.s.SetTickRate(fps)
							a.Field32 = previous
							b.ObjClass = class
							*(*byte)(unsafe.Add(b.UpdateData, 240)) = flags
							rv := legacy.PortTestWorldCollision(2, a, b, nil)
							admitted := class&object.ClassMonster == 0 && frame-previous >= uint32(int32(fps)>>1) && (class&object.ClassPlayer == 0 || flags&1 != 0)
							want := uint32(2)
							if admitted {
								want = 0
							} // Real inventory admission rejects zero carrying capacity.
							if rv != want {
								t.Fatalf("return%d want%d", rv, want)
							}
							rows = append(rows, struct {
								Name   string
								Return uint32
							}{name, rv})
						})
					}
				}
			}
		}
	}
	if legacy.PortTestWorldCollision(2, a, nil, nil) != 0 {
		t.Fatal("nil pickup receiver")
	}
}

func TestWorldCollisionsPickupIntegration(t *testing.T) {
	o := newWorldCollisionOwner(t)
	receiver := &o.units[0]
	item := newObjectXferSimple(t, o.s)
	t.Cleanup(func() { receiver.InvFirstItem = nil; item.InvHolder = nil; o.s.ObjSetOwner(nil, item) })
	receiver.CarryCapacity = 2
	receiver.InvFirstItem = nil
	*(*byte)(unsafe.Add(receiver.UpdateData, 240)) = 1
	item.Pickup = server.PortTestWorldPickupRegistry("DefaultPickup")
	item.Weight = 1
	item.ObjFlags = 0
	item.ScriptPickup.Func = -1
	item.NetCode = 2001
	item.Field32 = 100
	o.s.SetFrame(115)
	o.s.SetTickRate(30)
	rv := legacy.PortTestWorldCollision(2, item, receiver, nil)
	if rv != 1 || receiver.InvFirstItem != item || item.InvHolder != receiver || item.Owner() != receiver {
		t.Fatalf("pickup return%d inventory/owner mismatch", rv)
	}
	state := o.state()
	if len(state.Nodes) != 1 || state.Nodes[0].To != 1 {
		t.Fatal("pickup report")
	}

	spellbookCapture(t, "world-collisions-pickup-integration", []struct {
		Return uint32
		State  legacy.PortTestReliableReportState
	}{{rv, state}}, "c2d21586a33534cbf04e59214f41f13bbbc699768100784205991a9369b9e3a3")
}
