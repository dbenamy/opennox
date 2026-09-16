//go:build porttest

package opennox

import (
	"github.com/opennox/libs/wall"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
)

func TestWorldWallsVariantsLightingAndClipping(t *testing.T) {
	o := newWorldWallOwner(t)
	type result struct {
		Variant, Scene int
		Clip           image.Rectangle
		State          worldWallResult
	}
	var rows []result
	for _, variant := range []byte{0, 1, 15} {
		for scene := 0; scene < 6; scene++ {
			for _, clip := range []image.Rectangle{image.Rect(0, 0, 96, 96), image.Rect(0, 0, 45, 45), image.Rect(44, 44, 96, 96)} {
				for light := uint32(0); light < 2; light++ {
					o.resetWalls(t)
					o.c.r.Data().SetClip(true)
					o.c.r.Data().SetClipRect(clip)
					end := clip
					end.Max = end.Max.Sub(image.Pt(1, 1))
					o.c.r.Data().SetClipRect2(end)
					*o.words["edgeMinX"], *o.words["edgeMinY"] = uint32(clip.Min.X), uint32(clip.Min.Y)
					*o.words["edgeMaxX"], *o.words["edgeMaxY"] = uint32(clip.Max.X), uint32(clip.Max.Y)
					*o.words["highFront"], *o.words["translucent"], *o.words["highFloors"] = 1, 1, 1
					*memmap.PtrUint32(0x587000, 80816) = light
					for x := range o.c.tiles.nox_arr2_853BC0 {
						for y := range o.c.tiles.nox_arr2_853BC0[x] {
							o.c.tiles.nox_arr2_853BC0[x][y] = noxrender.RGB{R: (40 + (x*13)%180) << 16, G: (70 + (y*17)%150) << 16, B: (90 + (x*7+y*3)%130) << 16}
						}
					}
					flags := byte(3)
					switch scene {
					case 0:
						*memmap.PtrUint32(0x5D4594, 805848) = 1
					case 1:
						o.c.srv.Walls.DefByInd(0).Flags32 = 4
					case 2:
						flags = 11
					case 3:
						flags = 0x41
					case 4:
						flags = 0x43
					case 5:
						*o.words["highFront"] = 0
					}
					w := o.c.srv.Walls.CreateAtGrid(image.Pt(10, 10))
					w.Flags4 = wall.Flags(flags)
					w.Field2 = variant
					w.Dir0 = byte(scene % 2)
					before := append([]uint16(nil), o.pix.Pix...)
					legacy.PortTestWorldWalls(2, o.c.Viewport(), nil, w, image.Point{})
					for y := 0; y < 96; y++ {
						for x := 0; x < 96; x++ {
							i := y*o.pix.Stride + x
							if !image.Pt(x, y).In(clip) && before[i] != o.pix.Pix[i] {
								t.Fatalf("wall changed pixel outside clip at %d,%d", x, y)
							}
						}
					}
					if w.Flags4&3 != 0 || w.Field12 != 1 {
						t.Fatal("normal wall draw did not complete visibility state")
					}
					rows = append(rows, result{int(variant), scene, clip, o.captureWall(t, int(w.Dir0), flags, scene, light)})
				}
			}
		}
	}
	worldWallsCapture(t, "variants-lighting-clipping", rows, "bdbefd26349eb6451f1e2fe22d67e990e76a09d738f422250f7312004bb731ec")
}
