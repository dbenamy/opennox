//go:build porttest

package opennox

import (
	"encoding/binary"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
)

func TestClientUIRenderIntersection(t *testing.T) {
	type result struct {
		A, B, D int
		OK      bool
		Rects   [12]int32
	}
	values := [][4]int32{{0, 0, 96, 96}, {10, 10, 20, 20}, {20, 10, 30, 20}, {-5, -7, 1, 1}, {5, 5, 5, 10}, {20, 20, 10, 10}, {-2147483648, -2147483648, 2147483647, 2147483647}, {0, 0, 0, 0}}
	var out []result
	for ai, a := range values {
		for bi, b := range values {
			for dst := 0; dst < 3; dst++ {
				r := [12]int32{11, 22, 33, 44}
				copy(r[4:8], a[:])
				copy(r[8:], b[:])
				before := r
				ok := legacy.PortTestUIRenderIntersect(&r, dst)
				want := [4]int32{max(a[0], b[0]), max(a[1], b[1]), min(a[2], b[2]), min(a[3], b[3])}
				valid := want[0] < want[2] && want[1] < want[3]
				if ok != valid {
					t.Fatal("intersection acceptance")
				}
				if valid {
					copy(before[dst*4:], want[:])
				}
				if r != before {
					t.Fatal("intersection aliases or rejected-output mutation")
				}
				out = append(out, result{ai, bi, dst, ok, r})
			}
		}
	}
	effectsCapture(t, "ui-render-intersection", out, len(out), "6bda793def3833cc944df4ba986416119a3e7df94106c4076197f7f7870a4934")
}
func TestClientUIRenderFill(t *testing.T) {
	type result struct {
		Offset, Size, Wrapper int
		Color                 uint32
		Return                int32
		Bytes                 []byte
	}
	var out []result
	for off := 0; off < 8; off++ {
		for _, size := range []int32{-2147483648, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 15, 16, 17, 63, 64, 65, 127} {
			for _, color := range []uint32{0, 1, 0xffffffff, 0x12345678, 0x80007fff} {
				for wrapper := 0; wrapper < 2; wrapper++ {
					data := make([]byte, 160)
					for i := range data {
						data[i] = byte(i*17 + 3)
					}
					want := append([]byte(nil), data...)
					if size > 0 {
						for i := int32(0); i < size/4; i++ {
							binary.LittleEndian.PutUint32(want[off+int(i)*4:], color)
						}
						if size&3 >= 2 {
							binary.LittleEndian.PutUint16(want[off+int(size/4)*4:], uint16(color))
						}
					}
					ret := legacy.PortTestUIRenderFill(data, off, color, size, wrapper != 0)
					if ret != 0 || string(data) != string(want) {
						t.Fatal("pixel fill extent, remainder or return")
					}
					out = append(out, result{off, int(size), wrapper, color, ret, data})
				}
			}
		}
	}
	effectsCapture(t, "ui-render-fill", out, len(out), "44a77c3ae962c59d080bcdc2f6b03604fe9a1e2035448a05eb26ccd7767394e2")
}
func TestClientUIRenderClipAndBounds(t *testing.T) {
	o := newObjectRenderOwner(t)
	env := legacy.PortTestNewScreenEnvironment()
	defer env.Restore()
	type result struct {
		Case, Step int
		Return     int32
		Render     objectRenderResult
		Bounds     []uint32
	}
	var out []result
	id := 0
	rects := []image.Rectangle{image.Rect(0, 0, 96, 96), image.Rect(10, 20, 70, 80), {Min: image.Pt(20, 20), Max: image.Pt(10, 10)}, {Min: image.Pt(-2147483648, -2147483648), Max: image.Pt(2147483647, 2147483647)}}
	args := [][6]int32{{0, 0, 96, 96}, {10, 20, 30, 40}, {-5, -5, 10, 10}, {100, 100, 5, 5}, {20, 30, 0, 1}, {20, 30, -5, 10}, {2147483647, 1, 5, 10}, {-2147483648, -2, -1, 5}}
	for _, rect := range rects {
		for _, a := range args {
			for _, flag := range []int32{0, 1, -1, -2147483648} {
				id++
				o.resetRender(uint32(id), 120)
				env.Reset()
				d := o.c.r.Data()
				d.SetRect3(rect)
				capture := func(step int, ret int32) {
					out = append(out, result{id, step, ret, o.renderResult(t, id, step, int(ret)), env.State()})
				}
				ret := legacy.PortTestUIRenderOp(0, [6]int32{flag})
				if ret != 1 {
					t.Fatalf("initial clipping flag %d", ret)
				}
				capture(0, ret)
				copyRet := legacy.PortTestUIRenderOp(1, a)
				if copyRet != 0 && copyRet != 1 {
					t.Fatal("rectangle copy must return scalar false/true")
				}
				capture(1, copyRet)
				capture(2, legacy.PortTestUIRenderOp(2, [6]int32{a[0], a[2]}))
				legacy.PortTestObjectRenderClip(true)
				capture(3, legacy.PortTestUIRenderOp(1, [6]int32{2, 3, 7, 8}))
				legacy.PortTestObjectRenderClip(false)
				capture(4, 0)
				ret = legacy.PortTestUIRenderOp(0, [6]int32{0})
				if ret != flag {
					t.Fatal("raw clipping flag did not roundtrip")
				}
				capture(5, ret)
				ret = legacy.PortTestUIRenderOp(4, a)
				if ret != a[3] {
					t.Fatal("bounds setter return")
				}
				capture(6, ret)
			}
		}
	}
	effectsCapture(t, "ui-render-clip-bounds", out, len(out), "164c78fa31852e3e5ff07cdbc0b707e12413676d80afbf066fbe85d2e0567bfe")
}
func TestClientUIRenderBorders(t *testing.T) {
	o := newObjectRenderOwner(t)
	type result struct {
		Case   int
		Args   [6]int32
		Render objectRenderResult
	}
	var out []result
	id := 0
	for _, pos := range []image.Point{image.Pt(20, 20), image.Pt(-5, 10), image.Pt(95, 95), image.Pt(110, 110)} {
		for _, size := range []image.Point{image.Pt(0, 5), image.Pt(5, 0), image.Pt(1, 1), image.Pt(10, 15), image.Pt(-3, 7), image.Pt(60, 50)} {
			for _, color := range []int32{0, 0x7fff, 0x12345678} {
				for _, width := range []int32{1, 2, 4, 8} {
					for flag := 0; flag < 2; flag++ {
						// Unclipped particle drawing requires the entire border inside its buffer.
						if flag == 0 && pos != image.Pt(20, 20) {
							continue
						}
						id++
						o.resetRender(uint32(id), 120)
						legacy.PortTestUIRenderOp(0, [6]int32{int32(flag)})
						a := [6]int32{int32(pos.X), int32(pos.Y), int32(size.X), int32(size.Y), color, width}
						legacy.PortTestUIRenderOp(3, a)
						out = append(out, result{id, a, o.renderResult(t, id, 0, 0)})
					}
				}
			}
		}
	}
	effectsCapture(t, "ui-render-borders", out, len(out), "53a826566438b955d4cd3a4e3a0e80b2203c866501d82c440fd33eb5255f00ed")
}
