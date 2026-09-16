//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"image"
	"testing"
	"unsafe"
)

func TestTileCompositionAttachedEdges(t *testing.T) {
	o := newTileCompositionOwner(t)
	raw, _ := tileCompositionOverlay(0, 45, 4)
	raw2 := append([]byte(nil), o.raw...)
	for i := range raw2 {
		raw2[i] ^= 0xc5
	}
	images, freeImages := o.c.r.GetBag().PortTestWallEdgeImages([][]byte{o.raw, raw2, raw}, []int{0, 0, 0})
	for i := range o.handles {
		o.handles[i] = images[i%2].C()
	}
	t.Cleanup(freeImages)
	hs, freeHandles := alloc.Make([]noxrender.ImageHandle{}, 1)
	t.Cleanup(freeHandles)
	hs[0] = images[2].C()
	nodes, freeNodes := alloc.Make([][5]uint32{}, 2)
	t.Cleanup(freeNodes)
	for i := range nodes {
		nodes[i][0] = uint32(i)
		nodes[i][2] = uint32(i)
		*memmap.PtrUint32(0x85B3FC, 28676+60*uintptr(i)) = uint32(uintptr(unsafe.Pointer(&hs[0])))
	}
	t.Cleanup(func() {
		for _, r := range o.grid {
			for y := range r {
				r[y][5] = 0
				r[y][10] = 0
			}
		}
		for i := range nodes {
			*memmap.PtrUint32(0x85B3FC, 28676+60*uintptr(i)) = 0
		}
	})
	type row struct {
		Fill     bool
		Mask     int
		Position image.Point
		State    tileCompositionState
	}
	var rows []row
	for _, fill := range []bool{false, true} {
		for mask := 0; mask < 4; mask++ {
			for _, pos := range []image.Point{{0, 0}, {57, 103}, {230, 276}, {5800, 5800}} {
				o.populate(1)
				for _, r := range o.grid {
					for y := range r {
						r[y][0] = 3
						for half := 0; half < 2; half++ {
							if mask&(1<<half) != 0 {
								r[y][5+5*half] = uint32(uintptr(unsafe.Pointer(&nodes[half])))
							}
						}
					}
				}
				o.resetBuffer()
				nox_client_texturedFloors_154956 = !fill
				nox_xxx_tileSetDrawFn_481420()
				before := o.gridHash()
				legacy.PortTestTileCompositionFull(tileCompositionViewport(pos.X, pos.Y))
				if before != o.gridHash() {
					t.Fatal("composition changed attached edge inputs")
				}
				for i, n := range o.edgeUsage {
					want := uint32(0)
					if !fill && i < 2 && mask&(1<<i) != 0 {
						want = 1
					}
					if n != want {
						t.Fatalf("attached edge usage fill%v mask%d position%v edge%d got%d want%d", fill, mask, pos, i, n, want)
					}
				}
				rows = append(rows, row{fill, mask, pos, o.state()})
			}
		}
	}
	tileCompositionCapture(t, "attached-edges", rows, "e4e2ef4f01efc47182b48d8fa0aea2017a80aaade63599b5f63bceaf2f40006f")
}
