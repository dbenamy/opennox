//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/types"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestItemRespawnCrownAttribution(t *testing.T) {
	o := newMatchRosterOwner(t)
	t.Cleanup(o.s.PortTestRewardTypes([]string{"Crown", "OtherItem", "Carrier", "Glyph", "Torch", "Lantern"}, nil, true, 0, 0))
	for _, off := range []uintptr{1568248, 1568252, 1568256, 1568244} {
		serverConfigOwnBytes(t, 0x5D4594, off, 4)
		*memmap.PtrUint32(0x5D4594, off) = 0
	}
	oldPending := o.s.Objs.Pending
	t.Cleanup(func() { o.s.Objs.Pending = oldPending })
	type row struct {
		Mask        int
		Stamp       uint32
		Attributed  [3]uint32
		Owned, Held [3]bool
		Pos         [3]types.Pointf
		Packets     [][]byte
	}
	var rows []row
	for mask := 0; mask < 8; mask++ {
		for _, stamp := range []uint32{0, 123, 0xffffffff} {
			t.Run(fmt.Sprintf("mask%d/stamp%x", mask, stamp), func(t *testing.T) {
				o.reset()
				defer noxflags.PortTestGameFlags(2048)()
				o.s.Objs.Pending = nil
				owner := o.s.NewObjectByTypeID("Carrier")
				if owner == nil {
					t.Fatal("carrier")
				}
				owner.ObjClass = 0
				owner.PosVec = types.Pointf{X: 81, Y: 123}
				var units []*server.Object
				defer func() {
					o.s.Objs.Pending = nil
					owner.Field129 = nil
					owner.InvFirstItem = nil
					for _, u := range units {
						u.ObjNext = nil
						u.ObjPrev = nil
						u.ObjOwner = nil
						u.InvHolder = nil
						o.s.Objs.FreeObject(u)
					}
					o.s.Objs.FreeObject(owner)
				}()
				for i, name := range []string{"Crown", "OtherItem", "Crown"} {
					u := o.s.NewObjectByTypeID(name)
					if u == nil {
						t.Fatal("item")
					}
					u.ObjClass = 0
					u.ObjOwner = owner
					u.NetCode = uint32(200 + i)
					*(*uint32)(unsafe.Add(u.UpdateData, 4)) = 0x5678
					units = append(units, u)
				}
				owner.Field129 = units[0]
				for i, u := range units {
					if i+1 < len(units) {
						u.Field128 = units[i+1]
					}
					if mask&(1<<i) != 0 {
						u.InvHolder = owner
						u.InvNextItem = owner.InvFirstItem
						if owner.InvFirstItem != nil {
							owner.InvFirstItem.Field125 = u
						}
						owner.InvFirstItem = u
					}
				}
				legacy.PortTestItemRespawnDropCrown(owner, stamp)
				r := row{Mask: mask, Stamp: stamp, Packets: visibilityEffectsPackets(o.s)}
				for i, u := range units {
					r.Attributed[i] = *(*uint32)(unsafe.Add(u.UpdateData, 4))
					r.Owned[i] = u.ObjOwner == owner
					r.Held[i] = u.InvHolder == owner
					r.Pos[i] = u.PosVec
					wantStamp := stamp
					if i == 1 {
						wantStamp = 0x5678
					}
					dropped := i != 1 && mask&(1<<i) != 0
					if r.Attributed[i] != wantStamp || r.Owned[i] == dropped || r.Held[i] != (i == 1 && mask&2 != 0) {
						t.Fatal("crown result", i, r)
					}
					if dropped && u.PosVec != owner.PosVec {
						t.Fatal("drop position")
					}
				}
				rows = append(rows, r)
			})
		}
	}
	spellbookCapture(t, "item-respawn-crown-attribution", rows, "daf09a95d6809349bf9c6d1cdd5669a61bf0966a39eb2bc40b37b902ccb7cf70")
}
