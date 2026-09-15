//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"strings"
	"testing"
	"unsafe"
)

func TestClientInventoryDisplayStats(t *testing.T) {
	o := newInventoryDisplayOwner(t)
	pos, free := alloc.Make([]int32{}, 2)
	defer free()
	var rows []inventoryDisplayResult
	for class := 0; class < 3; class++ {
		for _, value := range []uint32{0, 1, 25, 50, 99, 100, 150, 255} {
			o.reset(t)
			pos[0], pos[1] = 2, 3
			*(*byte)(unsafe.Add(unsafe.Pointer(&o.players[0]), 2251)) = byte(class)
			for _, off := range []uintptr{2235, 2239, 2243, 2247} {
				*(*uint32)(unsafe.Add(unsafe.Pointer(&o.players[0]), off)) = value
			}
			blank := effectsPixelHash(o.pix)
			r := o.call(t, len(rows), 7, txptr(unsafe.Pointer(&pos[0])), 0, 0)
			if r.Pixels == blank || len(r.Text) == 0 {
				t.Fatal("stats did not render visible text")
			}
			hasMana, hasHealth := false, false
			for _, text := range r.Text {
				hasMana = hasMana || text.Text == "Mana"
				hasHealth = hasHealth || text.Text == "Health"
			}
			if !hasHealth || hasMana != (class != 0) {
				t.Fatalf("class%d stats labels: %v", class, r.Text)
			}
			rows = append(rows, r)
		}
	}
	inventoryDisplayCapture(t, "stats", rows, "bc7c963d695c59a805690d1ec977769a34cd98ddea90d75fa93aba96a8c26225")
}

func TestClientInventoryDisplayCoordinates(t *testing.T) {
	o := newInventoryDisplayOwner(t)
	pos, free := alloc.Make([]int32{}, 2)
	defer free()
	coords, freeCoords := alloc.Make([]int32{}, 2)
	defer freeCoords()
	for _, x := range []int32{-32768, -1, 0, 1, 32767} {
		for _, y := range []int32{-32768, -1, 0, 1, 32767} {
			o.reset(t)
			pos[0], pos[1] = x, y
			r := o.call(t, 0, 4, txptr(o.parent.C()), txptr(unsafe.Pointer(&pos[0])), txptr(unsafe.Pointer(&coords[0])))
			if coords[0] != x-4 || coords[1] != y-6 || uint32(r.Return) != uint32(y-6) {
				t.Fatalf("position %v converted%v return%x", pos, coords, r.Return)
			}
		}
	}
}

func TestClientInventoryDisplayDurability(t *testing.T) {
	o := newInventoryDisplayOwner(t)
	values, free := alloc.Make([]float32{}, 2)
	defer free()
	for _, health := range [][2]uint16{{0, 0}, {0, 100}, {1, 100}, {65535, 65535}, {123, 65535}} {
		o.reset(t)
		dr := o.item(t, "Bow", 123)
		*(*uint16)(unsafe.Add(dr.C(), 292)), *(*uint16)(unsafe.Add(dr.C(), 294)) = health[0], health[1]
		r := o.call(t, 0, 5, txptr(dr.C()), txptr(unsafe.Pointer(&values[0])), txptr(unsafe.Pointer(&values[1])))
		if values[0] != float32(health[0]) || values[1] != float32(health[1]) || uint32(r.Return) != o.norm(uint32(uintptr(dr.C()))) {
			t.Fatalf("raw durability %v became%v return%x", health, values, r.Return)
		}
	}
}

func TestClientInventoryDisplayModeTooltips(t *testing.T) {
	o := newInventoryDisplayOwner(t)
	var rows []inventoryDisplayResult
	for _, id := range []uint32{0, 9105, 9106, 9107, 9108, 9111, 9999} {
		o.reset(t)
		o.parent.SetID(uint(id))
		r := o.call(t, len(rows), 12, txptr(o.parent.C()), 0, 0)
		if id == 9107 {
			// The production tooltip setter owns the UTF16 cursor text.
			text := unsafe.Slice((*uint16)(memmap.PtrOff(0x5D4594, 1096676)), 256)
			if !strings.HasPrefix(alloc.GoString16(&text[0]), "Statistics") {
				t.Fatal("stats mode tooltip")
			}
		}
		rows = append(rows, r)
	}
	inventoryDisplayCapture(t, "tooltips", rows, "fd480747b351cc2ff7384c81ba9d6ec2b7c49551057b6a4f4251e5bef48f2130")
}
