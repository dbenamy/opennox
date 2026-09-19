//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"reflect"
	"testing"
)

func TestClientPresentationTransientRays(t *testing.T) {
	o := newObjectRenderOwner(t)
	state := serverConfigOwnBytes(t, 0x5D4594, 1303540, 860)
	type record struct {
		Count      int
		Deleted    []uint32
		Persistent uint32
		Beam       uint32
	}
	var rows []record
	for _, count := range []int{0, 1, 2, 95, 96} {
		o.resetRender(31, 120)
		clear(state)
		survivor := o.drawable(99, image.Pt(48, 48))
		*memmap.PtrPtr(0x5D4594, 1303924) = survivor.C()
		var refs []uint32
		for i := 0; i < count; i++ {
			dr := o.drawable(100+i, image.Pt(i, i))
			*memmap.PtrPtr(0x5D4594, 1303540+4*uintptr(i)) = dr.C()
			refs = append(refs, o.c.refs[dr])
		}
		*memmap.PtrUint32(0x5D4594, 1304308) = uint32(count)
		for dr := o.c.Objs.List1; dr != nil; dr = dr.NextPtr {
			if !legacy.PortTestPresentationRayContains(dr) {
				t.Fatal("transient/persistent membership")
			}
		}
		if legacy.PortTestObjectRenderBeam(1, o.c.Viewport(), [4]int32{1, 2, 3, 4}) != 1 {
			t.Fatal("beam setup")
		}
		legacy.PortTestPresentationSpareClear()
		if memmap.Uint32(0x5D4594, 1304308) != 0 || o.c.Objs.Count != 1 || !reflect.DeepEqual(o.c.Deleted, refs) || !legacy.PortTestPresentationRayContains(survivor) {
			t.Fatal("transient cleanup/order/persistent isolation")
		}
		beam := o.renderEnv.State()[1]
		if beam != 0 {
			t.Fatal("transient cleanup did not reset beam list")
		}
		legacy.PortTestPresentationSpareClear()
		if !reflect.DeepEqual(o.c.Deleted, refs) {
			t.Fatal("transient cleanup repeated deletion")
		}
		rows = append(rows, record{count, append([]uint32(nil), o.c.Deleted...), o.c.refs[survivor], beam})
		legacy.PortTestPresentationRayClear()
		if o.c.Objs.Count != 0 {
			t.Fatal("persistent cleanup")
		}
	}
	spellbookCapture(t, "client-presentation-transient-rays", rows, "46467824e7795de28a5cd2a58d6b6e59b0c6a08acdfb37665a47a2d9b087f62e")
}
