//go:build porttest

package opennox

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"image"
	"os"
	"testing"
	"unsafe"
)

type tileCompositionOwner struct {
	*tileRasterOwner
	grid             *[128]*[128][11]uint32
	counter          *uint32
	handles          []noxrender.ImageHandle
	usage, edgeUsage []uint32
}

func newTileCompositionOwner(t *testing.T) *tileCompositionOwner {
	o := &tileCompositionOwner{tileRasterOwner: newTileRasterOwner(t)}
	var undo func()
	o.grid, o.counter, undo = legacy.PortTestTileCompositionGrid()
	t.Cleanup(undo)
	for _, r := range [][3]uintptr{{0x85B3FC, 228, 704}, {0x5D4594, 2523980, 256}, {0x85B3FC, 28644, 3840}} {
		b := unsafe.Slice((*byte)(memmap.PtrOff(r[0], r[1])), r[2])
		old := append([]byte(nil), b...)
		t.Cleanup(func() { copy(b, old) })
		clear(b)
	}
	o.usage = unsafe.Slice(memmap.PtrUint32(0x85B3FC, 228), 176)
	o.edgeUsage = unsafe.Slice(memmap.PtrUint32(0x5D4594, 2523980), 64)
	raw2 := append([]byte(nil), o.raw...)
	for i := range raw2 {
		raw2[i] ^= 0xc5
	}
	imgs, freeImages := o.c.r.GetBag().PortTestWallEdgeImages([][]byte{o.raw, raw2}, []int{0, 0})
	t.Cleanup(freeImages)
	hs, freeHandles := alloc.Make([]noxrender.ImageHandle{}, 4)
	t.Cleanup(freeHandles)
	for i := range hs {
		hs[i] = imgs[i%2].C()
	}
	o.handles = hs
	for i := range o.defs {
		o.defs[i].Data32 = unsafe.Pointer(&hs[0])
		o.defs[i].Field46 = uint16(i % 2)
		o.defs[i].Color48 = uint32(0x1234abcd + i*0x10201)
		o.defs[i].Field58 = 0
	}
	return o
}
func (o *tileCompositionOwner) populate(pattern int) {
	for x, row := range o.grid {
		for y := range row {
			c := &row[y]
			*c = [11]uint32{}
			if pattern != 0 {
				c[0] = uint32((x + 2*y + pattern) % 4)
				c[1] = uint32((x + y) % 3)
				c[2] = uint32(y % 2)
				c[6] = uint32((x + 2*y + 1) % 3)
				c[7] = uint32(x % 2)
			}
		}
	}
}

type tileCompositionState struct {
	Pixels       [32]byte
	Words        map[string]uint32
	Counter      uint32
	Usage, Edges []uint32
}

func (o *tileCompositionOwner) state() tileCompositionState {
	s := tileCompositionState{Pixels: sha256.Sum256(o.bytes()), Words: map[string]uint32{}, Counter: *o.counter, Usage: append([]uint32(nil), o.usage...), Edges: append([]uint32(nil), o.edgeUsage...)}
	for n, p := range o.words {
		s.Words[n] = *p
	}
	return s
}
func tileCompositionViewport(x, y int) *noxrender.Viewport {
	return &noxrender.Viewport{Screen: image.Rect(3, 7, 99, 103), World: image.Rect(x+3, y+7, x+99, y+103), Size: image.Pt(96, 96)}
}
func tileCompositionCapture(t *testing.T, label string, rows any, want string) {
	t.Helper()
	b, e := json.Marshal(rows)
	if e != nil {
		t.Fatal(e)
	}
	if p := os.Getenv("OPENNOX_TILE_COMPOSITION_CAPTURE"); p != "" {
		if e := os.WriteFile(p+"-"+label+".json", b, 0600); e != nil {
			t.Fatal(e)
		}
	}
	got := fmt.Sprintf("%x", sha256.Sum256(b))
	t.Logf("%s: %s", label, got)
	if got != want {
		t.Fatalf("%s got%s want%s", label, got, want)
	}
}

func (o *tileCompositionOwner) gridHash() [32]byte {
	h := sha256.New()
	for _, row := range o.grid {
		h.Write(unsafe.Slice((*byte)(unsafe.Pointer(row)), 128*44))
	}
	var sum [32]byte
	copy(sum[:], h.Sum(nil))
	return sum
}
