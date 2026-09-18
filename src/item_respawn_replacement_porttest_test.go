//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestItemRespawnReplacement(t *testing.T) {
	o := newMatchRosterOwner(t)
	t.Cleanup(o.s.PortTestRewardTypes([]string{"Crown", "PortTestRewardWeapon"}, nil, true, 0, 0x82))
	_, restore := legacy.PortTestItemRespawnGlobals()
	t.Cleanup(restore)
	oldPending, oldDeleted := o.s.Objs.Pending, o.s.Objs.DeletedList
	t.Cleanup(func() { o.s.Objs.Pending = oldPending; o.s.Objs.DeletedList = oldDeleted })
	t.Cleanup(o.s.PortTestCombatAudioReset)
	type row struct {
		OldClass   uint32
		Deleted    bool
		DeletedAt  uint32
		Ammo       [2]byte
		Direction  [2]uint16
		Pending    uint32
		AudioCount int
		Packets    [][]byte
	}
	var rows []row
	for _, class := range []uint32{0, 2, 0x1000000} {
		t.Run(fmt.Sprintf("class%x", class), func(t *testing.T) {
			o.reset()
			o.s.PortTestCombatAudioReset()
			o.s.Objs.Pending = nil
			o.s.Objs.DeletedList = nil
			u := o.s.NewObjectByTypeID("PortTestRewardWeapon")
			if u == nil {
				t.Fatal("allocation")
			}
			u.ObjClass = object.Class(class)
			u.ObjFlags = 0
			u.PosVec = types.Pointf{X: 32, Y: 64}
			u.Direction1 = 0x8123
			// Nonzero ammunition is retained from the original weapon definition.
			data := unsafe.Slice((*byte)(u.UseData.Ptr), 64)
			data[0] = 13
			data[1] = 27
			cleanup := legacy.PortTestTeamRuntimeRespawns([]*server.Object{u})
			defer cleanup()
			r := legacy.PortTestItemRespawnRecords()[0]
			r[5] = o.s.Frame()
			r[6] = 1
			// Select the same real weapon type regardless of the old object's current class.
			r[12] = 13<<8 | 27
			legacy.PortTestItemRespawn("tick", nil)
			got := (*server.Object)(unsafe.Pointer(uintptr(r[1])))
			if got == nil || got == u || r[6] != 0 {
				t.Fatal("replacement")
			}
			gotData := unsafe.Slice((*byte)(got.UseData.Ptr), 64)
			state := row{OldClass: class, Deleted: u.Flags().Has(object.FlagDestroyed), DeletedAt: u.DeletedAt, Ammo: [2]byte{gotData[0], gotData[1]}, Direction: [2]uint16{uint16(got.Direction1), uint16(got.Direction2)}, Pending: r[6], AudioCount: len(o.s.PortTestCombatAudioSnapshot()), Packets: visibilityEffectsPackets(o.s)}
			if state.Deleted != (class&2 != 0) || class&2 != 0 && state.DeletedAt != 123 || state.Ammo != [2]byte{13, 27} || state.Direction != [2]uint16{0x8123, 0x8123} || state.AudioCount != 1 {
				t.Fatal("replacement effects", state)
			}
			rows = append(rows, state)
			o.s.PortTestCombatAudioReset()
			o.s.Objs.Pending = nil
			o.s.Objs.DeletedList = nil
			got.ObjNext = nil
			got.ObjPrev = nil
			u.DeletedNext = nil
			o.s.Objs.FreeObject(got)
			o.s.Objs.FreeObject(u)
		})
	}
	spellbookCapture(t, "item-respawn-replacement", rows, "49bd4d214a3a47f4b83b0c42c91c7c9eeb99b6bedd1b2a5efdcdf3b3ed91fd5a")
}
