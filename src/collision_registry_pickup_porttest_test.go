//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestCollisionRegistryPickup(t *testing.T) {
	type row struct {
		Name   string
		Picked bool
		State  legacy.PortTestReliableReportState
	}
	var rows []row
	for _, sp := range []struct {
		name             string
		capacity         uint16
		frame            uint32
		absent, admitted bool
	}{
		{"accepted", 2, 115, false, true}, {"too-soon", 2, 114, false, false},
		{"capacity", 0, 115, false, false}, {"nil-target", 2, 115, true, false},
	} {
		t.Run(sp.name, func(t *testing.T) {
			o := newWorldCollisionOwner(t)
			receiver := &o.units[0]
			item := newObjectXferSimple(t, o.s)
			t.Cleanup(func() { receiver.InvFirstItem = nil; item.InvHolder = nil; o.s.ObjSetOwner(nil, item) })
			receiver.CarryCapacity = sp.capacity
			receiver.InvFirstItem = nil
			*(*byte)(unsafe.Add(receiver.UpdateData, 240)) = 1
			item.Pickup = server.PortTestWorldPickupRegistry("DefaultPickup")
			item.Weight = 1
			item.ObjFlags = 0
			item.ScriptPickup.Func = -1
			item.NetCode = 2001
			item.Field32 = 100
			o.s.SetFrame(sp.frame)
			o.s.SetTickRate(30)
			target := receiver
			if sp.absent {
				target = nil
			}
			legacy.PortTestRegisteredCollision(item, target, nil, "PickupCollide")
			if sp.admitted {
				if receiver.InvFirstItem != item || item.InvHolder != receiver || item.Owner() != receiver {
					t.Fatal("pickup ownership")
				}
			} else if receiver.InvFirstItem != nil || item.InvHolder != nil || item.Owner() != nil {
				t.Fatal("rejected pickup changed ownership")
			}
			state := o.state()
			if sp.admitted && (len(state.Nodes) != 1 || state.Nodes[0].To != 1) {
				t.Fatal("pickup report")
			}
			rows = append(rows, row{sp.name, item.InvHolder == receiver, state})
		})
	}
	collisionRegistryCapture(t, "collision-registry-pickup", rows)
}
