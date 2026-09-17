//go:build porttest

package opennox

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
)

func TestObjectXferPlacementContracts(t *testing.T) {
	core := newObjectXferOwner(t)
	wallSize := (*[2]int32)(memmap.PtrOff(0x5D4594, 739980))
	oldWall := *wallSize
	t.Cleanup(func() { *wallSize = oldWall })
	cache := memmap.PtrUint32(0x5D4594, 1563904)
	oldCache := *cache
	t.Cleanup(func() { *cache = oldCache })
	for _, flags := range []uint32{0, 0x200000, 0x400000} {
		for _, allowed := range []bool{false, true} {
			for _, filtered := range []bool{false, true} {
				for _, offset := range []bool{false, true} {
					t.Run(fmt.Sprintf("flags%x-allowed%v-filtered%v-offset%v", flags, allowed, filtered, offset), func(t *testing.T) {
						t.Cleanup(noxflags.PortTestGameFlags(noxflags.GameFlag(flags)))
						editorList, closeList := legacy.PortTestObjectXferEditorList()
						defer closeList()
						cl, fl := object.ClassSimple, object.Flags(0)
						if filtered {
							cl, fl = object.ClassImmobile, object.FlagNoCollide
						}
						core.PortTestObjectXferAdmission(cl, fl, allowed)
						u := core.NewObjectByTypeInd(1)
						if u == nil {
							t.Fatal("placement allocation")
						}
						u.PosVec = types.Pointf{X: 64.25, Y: 128.75}
						*wallSize = [2]int32{3, -2}
						point := [2]int32{97, -31}
						var arg unsafe.Pointer
						if offset {
							arg = unsafe.Pointer(&point)
						}
						beforeLive, beforeList := core.Objs.Alive, len(editorList())
						got := legacy.Nox_xxx_servMapLoadPlaceObj_4F3F50(u, 0, arg)
						accepted := flags == 0x400000 || allowed && (flags == 0x200000 || !filtered)
						if !accepted {
							if got != 0 || core.Objs.Alive != beforeLive-1 || len(editorList()) != beforeList {
								t.Fatal("rejected placement ownership/result")
							}
							return
						}
						trackObjectXferTyped(t, core, u)
						if got != 1 || core.Objs.Alive != beforeLive {
							t.Fatal("accepted placement ownership/result")
						}
						want := types.Pointf{X: 64.25, Y: 128.75}
						if offset {
							want.X = float32(float64(want.X) - 69 + 97 - 11)
							want.Y = float32(float64(want.Y) + 46 - 31 - 11)
						}
						if u.PosVec != want {
							t.Fatalf("offset position=%v want=%v", u.PosVec, want)
						}
						if flags == 0x400000 {
							list := editorList()
							if len(list) != beforeList+1 || list[0] != u {
								t.Fatal("real editor list insertion")
							}
						} else {
							if core.Objs.Pending != u || u.NewPos != want || u.PrevPos != want || u.ObjFlags&(object.FlagActive|object.FlagPending) != (object.FlagActive|object.FlagPending) {
								t.Fatal("real world creation")
							}
							core.Objs.Pending = nil
						}
					})
				}
			}
		}
	}
}
