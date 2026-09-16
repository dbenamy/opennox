//go:build porttest && !server

package opennox

import (
	"github.com/opennox/libs/wall"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
)

// Server builds deliberately have no field-of-view clipper. Exercise the real
// client clipper with its actual scanline intervals in both client targets.
func TestWorldWallsFieldOfView(t *testing.T) {
	o := newWorldWallOwner(t)
	type result struct {
		Override, Clip int
		State          worldWallResult
	}
	var rows []result
	for dir := 0; dir < 11; dir++ {
		for override := 0; override < 16; override++ {
			for clip := 0; clip < 3; clip++ {
				o.resetWalls(t)
				*memmap.PtrUint32(0x587000, 80808) = 1
				for y := 0; y < 96; y++ {
					o.c.tiles.nox_arr_956A00[y] = 2
					o.c.tiles.nox_arr_957820[y].arr[0] = 0
					right := 96
					if clip == 1 {
						right = 48
					}
					if clip == 2 {
						right = 10
					}
					o.c.tiles.nox_arr_957820[y].arr[1] = right
				}
				w := o.c.srv.Walls.CreateAtGrid(image.Pt(10, 10))
				w.Dir0 = byte(dir)
				w.Field3 = byte(override)
				w.Flags4 = wall.Flags(1)
				before := effectsPixelHash(o.pix)
				legacy.PortTestWorldWalls(2, o.c.Viewport(), nil, w, image.Point{})
				if clip == 2 && effectsPixelHash(o.pix) != before {
					t.Fatal("wall beyond field of view changed pixels")
				}
				rows = append(rows, result{override, clip, o.captureWall(t, dir, 1, 0, 0)})
			}
		}
	}
	worldWallsCapture(t, "field-of-view", rows, "08bc8250f4ea083893469e27c94b575da89834afec5cbb13925db95ef1f30330")
}
