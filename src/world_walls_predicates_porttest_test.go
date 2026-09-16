//go:build porttest

package opennox

import (
	"bytes"
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
	"unsafe"
)

func TestWorldWallsDrawablePasses(t *testing.T) {
	o := newObjectRenderOwner(t)
	dr := o.drawable(7, image.Pt(40, 40))
	callbacks := []unsafe.Pointer{nil, legacy.PortTestSpriteAnimationCallback(2), legacy.PortTestSpriteAnimationCallback(3), legacy.PortTestSpriteAnimationCallback(0)}
	type result struct {
		Callback     int
		Flags, Class uint32
		Returns      [4]int
	}
	var rows []result
	bits := []uint32{1, 8, 64, 0x800, 0x1000, 0x4000, 0x1000000}
	for callback, fn := range callbacks {
		for mask := 0; mask < 1<<len(bits); mask++ {
			flags := uint32(0)
			for i, b := range bits {
				if mask&(1<<i) != 0 {
					flags |= b
				}
			}
			for _, class := range []uint32{0, 4, 0x2000, 0x400000, 0x800000, 0x80000000, 0x80800000, 0x20400000} {
				dr.DrawFuncPtr = fn
				dr.ObjFlags = object.Flags(flags)
				dr.ObjClass = object.Class(class)
				before := append([]byte(nil), unsafe.Slice((*byte)(unsafe.Pointer(dr)), int(unsafe.Sizeof(*dr)))...)
				row := result{Callback: callback, Flags: flags, Class: class}
				for op := 4; op < 8; op++ {
					row.Returns[op-4], _ = legacy.PortTestWorldWalls(op, nil, dr, nil, image.Point{})
				}
				if !bytes.Equal(before, unsafe.Slice((*byte)(unsafe.Pointer(dr)), int(unsafe.Sizeof(*dr)))) {
					t.Fatal("pass predicate mutated drawable")
				}
				eligible := fn != nil && flags&1 != 0 && flags&0x1000 == 0
				if (row.Returns[0]+row.Returns[1] == 1) != eligible {
					t.Fatal("active passes did not partition eligible drawables")
				}
				if row.Returns[3] != 0 && flags&1 != 0 {
					t.Fatal("inactive pass accepted active drawable")
				}
				if fn == nil && row.Returns != [4]int{} {
					t.Fatal("drawable with no callback entered a draw pass")
				}
				rows = append(rows, row)
			}
		}
	}
	worldWallsCapture(t, "drawable-passes", rows, "e479a374b5fb1e9675eeeedbebdf87cfda0ba0acf649074b1ce3aab2e14ae1c9")
}
