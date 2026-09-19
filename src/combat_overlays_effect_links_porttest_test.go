//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"image"
	"reflect"
	"testing"
	"unsafe"
)

func TestCombatOverlayEffectLinks(t *testing.T) {
	type record struct {
		Name   string
		Global []int
		Actors [2][]int
		Links  [][5]int
	}
	var rows []record
	for remove := -1; remove < 4; remove++ {
		t.Run(fmt.Sprint(remove), func(t *testing.T) {
			o := newCombatOverlayOwner(t)
			actors := []*client.Drawable{o.drawable(7, image.Pt(300, 300)), o.drawable(8, image.Pt(350, 300))}
			var nodes []unsafe.Pointer
			ids := map[unsafe.Pointer]int{nil: 0}
			for i := 0; i < 4; i++ {
				p, free := alloc.New([20]uint32{})
				t.Cleanup(free)
				nodes = append(nodes, unsafe.Pointer(p))
				ids[unsafe.Pointer(p)] = i + 1
			}
			global := memmap.PtrPtr(0x5D4594, 1203872)
			t.Cleanup(func() {
				for i := 0; i < 4 && *global != nil; i++ {
					legacy.PortTestCombatEffectDetach(*global)
				}
				*global = nil
				for _, dr := range actors {
					dr.Field_114 = nil
				}
			})
			ownerIDs := map[unsafe.Pointer]int{nil: 0, actors[0].C(): 1, actors[1].C(): 2}
			ref := func(p unsafe.Pointer) int {
				v, ok := ids[p]
				if !ok {
					t.Fatal("link points outside owned effects")
				}
				return v
			}
			chain := func(p unsafe.Pointer, next uintptr) []int {
				var out []int
				for p != nil {
					if len(out) >= 4 {
						t.Fatal("effect list cycle")
					}
					out = append(out, ref(p))
					p = *(*unsafe.Pointer)(unsafe.Add(p, next))
				}
				return out
			}
			legacy.PortTestCombatEffectDetach(nil)
			legacy.PortTestCombatEffectAttach(nil, actors[0])
			legacy.PortTestCombatEffectAttach(nodes[0], nil)
			if *global != nil || actors[0].Field_114 != nil {
				t.Fatal("nil attach changed heads")
			}
			for i, actor := range []int{0, 0, 1, 0} {
				legacy.PortTestCombatEffectAttach(nodes[i], actors[actor])
			}
			if remove >= 0 {
				legacy.PortTestCombatEffectDetach(nodes[remove])
			}
			r := record{Name: fmt.Sprint(remove), Global: chain(*global, 72)}
			for i, dr := range actors {
				r.Actors[i] = chain(unsafe.Pointer(dr.Field_114), 64)
			}
			var wantGlobal []int
			var wantActors [2][]int
			for _, id := range []int{4, 3, 2, 1} {
				if id == remove+1 {
					continue
				}
				wantGlobal = append(wantGlobal, id)
				actor := 0
				if id == 3 {
					actor = 1
				}
				wantActors[actor] = append(wantActors[actor], id)
			}
			if !reflect.DeepEqual(r.Global, wantGlobal) || !reflect.DeepEqual(r.Actors, wantActors) {
				t.Fatalf("effect lists %+v", r)
			}
			for _, p := range nodes {
				var row [5]int
				who := *(*unsafe.Pointer)(unsafe.Add(p, 60))
				var ok bool
				row[0], ok = ownerIDs[who]
				if !ok {
					t.Fatal("effect owner outside fixture")
				}
				for i, off := range []uintptr{64, 68, 72, 76} {
					row[i+1] = ref(*(*unsafe.Pointer)(unsafe.Add(p, off)))
				}
				r.Links = append(r.Links, row)
			}
			rows = append(rows, r)
		})
	}
	spellbookCapture(t, "combat-overlay-effect-links", rows, "b6741034fb488dd05b7c66266400a790a7d61d9c4b87cd742d0d64b7a47327dc")
}
