//go:build porttest

package opennox

import (
	"bytes"
	"github.com/opennox/libs/wall"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
	"unsafe"
)

type worldWallResult struct {
	Direction     int
	Flags         byte
	Configuration int
	LightMode     uint32
	Wall          [16]byte
	Words         map[string]uint32
	Render        objectRenderResult
}

func (o *worldWallOwner) captureWall(t *testing.T, dir int, flags byte, config int, light uint32) worldWallResult {
	w := o.c.srv.Walls.GetWallAtGrid(image.Pt(10, 10))
	var data [16]byte
	copy(data[:], unsafe.Slice((*byte)(unsafe.Pointer(w)), len(data)))
	words := map[string]uint32{}
	for n, p := range o.words {
		words[n] = *p
	}
	return worldWallResult{dir, flags, config, light, data, words, o.renderResult(t, dir, 0, 0)}
}
func TestWorldWallsRendering(t *testing.T) {
	o := newWorldWallOwner(t)
	var rows []worldWallResult
	for dir := 0; dir < 11; dir++ {
		for _, flags := range []byte{0, 1, 3, 9, 0x43} {
			for config := 0; config < 8; config++ {
				for light := uint32(0); light < 2; light++ {
					o.resetWalls(t)
					*o.words["highFront"] = uint32(config & 1)
					*o.words["highFloors"] = uint32((config >> 1) & 1)
					*o.words["translucent"] = uint32((config >> 2) & 1)
					*memmap.PtrUint32(0x5D4594, 805848) = 1
					*memmap.PtrUint32(0x587000, 80816) = light
					w := o.c.srv.Walls.CreateAtGrid(image.Pt(10, 10))
					w.Dir0 = byte(dir)
					w.Flags4 = wall.Flags(flags)
					before := append([]byte(nil), unsafe.Slice((*byte)(unsafe.Pointer(w)), int(unsafe.Sizeof(*w)))...)
					pixels := effectsPixelHash(o.pix)
					legacy.PortTestWorldWalls(2, o.c.Viewport(), nil, w, image.Point{})
					if flags&1 == 0 {
						if !bytes.Equal(before, unsafe.Slice((*byte)(unsafe.Pointer(w)), int(unsafe.Sizeof(*w)))) || pixels != effectsPixelHash(o.pix) {
							t.Fatal("inactive wall draw changed wall or pixels")
						}
					} else {
						if w.Flags4&3 != 0 || w.Field3 != 0 || w.Field12 != 1 {
							t.Fatal("wall draw did not consume flags and mark explored")
						}
						if pixels == effectsPixelHash(o.pix) {
							t.Fatalf("visible wall dir%d flags%x config%d light%d drew no pixels", dir, flags, config, light)
						}
					}
					rows = append(rows, o.captureWall(t, dir, flags, config, light))
				}
			}
		}
	}
	worldWallsCapture(t, "rendering", rows, "8f4eb579e0681cf61f69161f1ca8ffcf55f00466c11eae41ee7694b69c654702")
}

func TestWorldWallsNilInput(t *testing.T) {
	o := newWorldWallOwner(t)
	var rows []objectRenderResult
	for _, nilViewport := range []bool{false, true} {
		o.resetWalls(t)
		vp := o.c.Viewport()
		if nilViewport {
			vp = nil
		}
		before := effectsPixelHash(o.pix)
		ret, _ := legacy.PortTestWorldWalls(2, vp, nil, nil, image.Point{})
		if ret != 0 || effectsPixelHash(o.pix) != before {
			t.Fatal("nil wall changed drawing")
		}
		rows = append(rows, o.renderResult(t, 0, 0, ret))
	}
	worldWallsCapture(t, "nil-input", rows, "6898cad19c98a98f8f3dc515f15393d46c4ba518eb03b2fece4023cab3d04a45")
}
