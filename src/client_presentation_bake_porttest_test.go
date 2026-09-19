//go:build porttest

package opennox

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"image"
	"testing"
	"unsafe"
)

func TestClientPresentationBake(t *testing.T) {
	o := newTileCompositionOwner(t)
	raws := [][]byte{spriteAnimationTestImage(0), spriteAnimationTestImage(1)}
	var mixedPixels [5][7]uint16
	var mixedMask [5][7]bool
	mixed := make([]byte, 17)
	binary.LittleEndian.PutUint32(mixed, 7)
	binary.LittleEndian.PutUint32(mixed[4:], 5)
	binary.LittleEndian.PutUint32(mixed[8:], ^uint32(2)) // image origin X = -3
	binary.LittleEndian.PutUint32(mixed[12:], 5)
	for y := 0; y < 5; y++ {
		left, right := y%3, y%2
		if left > 0 {
			mixed = append(mixed, 1, byte(left))
		}
		mixed = append(mixed, 2, byte(7-left-right))
		for x := left; x < 7-right; x++ {
			value := uint16(0x1234 + 17*x + 39*y)
			mixedPixels[y][x] = value
			mixedMask[y][x] = true
			mixed = binary.LittleEndian.AppendUint16(mixed, value)
		}
		if right > 0 {
			mixed = append(mixed, 1, byte(right))
		}
	}
	raws = append(raws, mixed)
	imgs, freeImgs := o.c.r.GetBag().PortTestSpriteImages(raws)
	t.Cleanup(freeImgs)
	frames, freeFrames := alloc.Make([]uint32{}, 3)
	t.Cleanup(freeFrames)
	for i := range frames {
		frames[i] = uint32(uintptr(imgs[i].C()))
	}
	data, freeData := alloc.Make([]uint32{}, 2)
	t.Cleanup(freeData)
	data[0] = 8
	dr := o.drawable(7, image.Point{})
	dr.DrawData = unsafe.Pointer(&data[0])
	*o.words["originX"], *o.words["originY"] = 100, 200
	w, h, stride := int(*o.words["width"]), int(*o.words["height"]), int(*o.words["stride"])
	type record struct {
		Name   string
		Stamp  uint32
		Buffer string
	}
	var rows []record
	positions := []image.Point{{-1, 2}, {0, 2}, {w - 8, 2}, {w - 7, 2}, {2, -1}, {2, 0}, {2, h - 6}, {2, h - 5}}
	for _, anim := range []bool{false, true} {
		for imageIndex := 0; imageIndex < 3; imageIndex++ {
			dr.DrawFuncPtr = legacy.PortTestSpriteAnimationCallback(2)
			data[1] = frames[imageIndex]
			if anim {
				dr.AnimFrameSlave = uint32(imageIndex)
				dr.DrawFuncPtr = legacy.PortTestSpriteAnimationCallback(5)
				data[1] = uint32(uintptr(unsafe.Pointer(&frames[0])))
			}
			for _, clock := range [][2]uint32{{0, 7}, {1, 0}, {2, 0}, {7, 6}, {7, 7}, {7, 8}, {0xffffffff, 7}} {
				for _, pos := range positions {
					for _, scroll := range []image.Point{{0, 0}, {w - 3, 0}, {0, h - 2}, {w - 3, h - 2}} {
						name := fmt.Sprintf("anim=%v/image=%d/clock=%v/pos=%v/scroll=%v", anim, imageIndex, clock, pos, scroll)
						o.resetBuffer()
						*o.counter = clock[0]
						dr.Field_86 = clock[1]
						dr.Field_0, dr.ZVal, dr.ZVal2 = 0, 0, 0
						dr.PosVec = pos.Add(image.Pt(100, 200))
						if imageIndex == 2 {
							// Image origin (-3,5), sprite inset (2,3), height (-3,4).
							// Compensating world position places the first pixel at pos.
							dr.Field_0, dr.ZVal, dr.ZVal2 = 0x0302, 65533, 4
							dr.PosVec = dr.PosVec.Add(image.Pt(5, -1))
						}
						*o.words["scrollX"], *o.words["scrollY"] = uint32(scroll.X), uint32(scroll.Y)
						want := append([]byte(nil), o.bytes()...)
						stamp := clock[1]
						inside := pos.X >= 0 && pos.X+7 < w && pos.Y >= 0 && pos.Y+5 < h
						if !inside {
							stamp = 0
						} else {
							if int32(clock[0]) <= 0 {
								stamp = 0
							}
							if int32(clock[0])-int32(stamp) > 1 || int32(clock[0]) <= 0 {
								// Independent literal pixels placed into the circular floor buffer.
								for y := 0; y < 5; y++ {
									for x := 0; x < 7; x++ {
										if imageIndex == 2 && !mixedMask[y][x] {
											continue
										}
										value := mixedPixels[y][x]
										if imageIndex != 2 {
											src := 17 + y*16 + 2 + 2*x
											value = binary.LittleEndian.Uint16(raws[imageIndex][src:])
										}
										at := ((pos.Y+scroll.Y+y)*stride + 2*(pos.X+scroll.X+x)) % len(want)
										binary.LittleEndian.PutUint16(want[at:], value)
									}
								}
							} else {
								stamp = clock[0]
							}
						}
						legacy.PortTestPresentationBake(o.c.Viewport(), dr)
						if dr.Field_86 != stamp || !bytes.Equal(o.bytes(), want) {
							t.Fatalf("%s: stamp=%d want %d or floor pixels differ", name, dr.Field_86, stamp)
						}
						rows = append(rows, record{name, dr.Field_86, fmt.Sprintf("%x", sha256.Sum256(o.bytes()))})
					}
				}
			}
		}
	}
	for _, mode := range []string{"inactive-static", "missing-image"} {
		o.resetBuffer()
		want := append([]byte(nil), o.bytes()...)
		dr.Field_86 = 99
		dr.DrawFuncPtr = legacy.PortTestSpriteAnimationCallback(2)
		dr.ObjClass, dr.ObjFlags = 0, 0
		data[1] = 0
		if mode == "inactive-static" {
			dr.ObjClass = 0x40000
			data[1] = frames[0]
		}
		legacy.PortTestPresentationBake(o.c.Viewport(), dr)
		if dr.Field_86 != 99 || !bytes.Equal(o.bytes(), want) {
			t.Fatal("bake early exit changed state")
		}
		rows = append(rows, record{mode, dr.Field_86, fmt.Sprintf("%x", sha256.Sum256(o.bytes()))})
	}
	spellbookCapture(t, "client-presentation-bake", rows, "3e5c37ae93a6a95623dc84ab9841ea05208ca8ba220324a4dc4e971611d8bace")
}
