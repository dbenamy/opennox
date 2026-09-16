//go:build porttest

package opennox

import (
	"bytes"
	"crypto/sha256"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"image"
	"testing"
	"unsafe"
)

type tileCompositionRun struct {
	Op, N  int
	Pixels []byte
}

func tileCompositionOverlay(first, last, mode int) ([]byte, [][]tileCompositionRun) {
	b := []byte{byte(first), byte(last)}
	rows := make([][]tileCompositionRun, 46)
	for y := first; y <= last; y++ {
		_, n, _ := tileRasterRow(y)
		for n > 0 {
			op := 1 + mode%3
			size := n
			if mode >= 3 {
				size = min(n, 1+(y+n)%7)
				op = 1 + (y+n+mode)%3
			}
			r := tileCompositionRun{Op: op, N: size}
			b = append(b, byte(op), byte(size))
			if op == 2 {
				r.Pixels = make([]byte, size*2)
				for i := range r.Pixels {
					r.Pixels[i] = byte(y*31 + i*17 + mode*53)
				}
				b = append(b, r.Pixels...)
			}
			rows[y] = append(rows[y], r)
			n -= size
		}
	}
	return b, rows
}
func TestTileCompositionOverlays(t *testing.T) {
	o := newTileCompositionOwner(t)
	type result struct {
		Mode, First, Last, Count int
		Position                 image.Point
		Fill, Nil                bool
		Pixels                   [32]byte
		Usage                    []uint32
	}
	var results []result
	for mode := 0; mode < 6; mode++ {
		for _, span := range [][2]int{{0, 0}, {0, 45}, {1, 44}, {22, 23}, {45, 45}} {
			raw, runs := tileCompositionOverlay(span[0], span[1], mode)
			raw2 := append([]byte(nil), o.raw...)
			for i := range raw2 {
				raw2[i] ^= 0xc5
			}
			images, freeImages := o.c.r.GetBag().PortTestWallEdgeImages([][]byte{o.raw, raw2, raw}, []int{0, 0, 0})
			oldHandles := append([]noxrender.ImageHandle(nil), o.handles...)
			for i := range o.handles {
				o.handles[i] = images[i%2].C()
			}
			hs, freeHandles := alloc.Make([]noxrender.ImageHandle{}, 3)
			hs[2] = images[2].C()
			nodes, freeNodes := alloc.Make([][5]uint32{}, 2)
			for i := range nodes {
				nodes[i][0] = uint32(i)
				nodes[i][1] = 0
				nodes[i][2] = uint32(i)
				nodes[i][3] = 1
			}
			for i := 0; i < 2; i++ {
				*memmap.PtrUint32(0x85B3FC, 28676+60*uintptr(i)) = uint32(uintptr(unsafe.Pointer(&hs[0])))
				*memmap.PtrUint16(0x85B3FC, 28690+60*uintptr(i)) = 1
			}
			for _, pos := range []image.Point{{0, 0}, {31, 17}, {230, 184}, {253, 229}, {0, 230}, {31, 400}} {
				for _, fill := range []bool{false, true} {
					for _, nilImage := range []bool{false, true} {
						for _, count := range []int{1, 2} {
							o.resetBuffer()
							clear(o.edgeUsage)
							nox_client_texturedFloors_154956 = !fill
							nox_xxx_tileSetDrawFn_481420()
							hs[2] = images[2].C()
							if nilImage {
								hs[2] = nil
							}
							nodes[0][4] = 0
							if count == 2 {
								nodes[0][4] = uint32(uintptr(unsafe.Pointer(&nodes[1])))
							}
							want := append([]byte(nil), o.bytes()...)
							if !fill && !nilImage {
								for k := 0; k < count; k++ {
									for y := span[0]; y <= span[1]; y++ {
										left, _, src := tileRasterRow(y)
										at := ((pos.Y+y)*int(*o.words["stride"]) + 2*(pos.X+left)) % len(want)
										for _, r := range runs[y] {
											for i := 0; i < r.N*2; i++ {
												if r.Op == 2 {
													want[(at+i)%len(want)] = r.Pixels[i]
												} else if r.Op == 3 {
													v := o.raw[src+i]
													if k == 1 {
														v ^= 0xc5
													}
													want[(at+i)%len(want)] = v
												}
											}
											at = (at + 2*r.N) % len(want)
											src += 2 * r.N
										}
									}
								}
							}
							legacy.PortTestTileCompositionEdges(pos, unsafe.Pointer(&nodes[0]))
							if !bytes.Equal(want, o.bytes()) {
								t.Fatalf("overlay pixels mode%d span%v pos%v fill%v nil%v count%d", mode, span, pos, fill, nilImage, count)
							}
							for i, n := range o.edgeUsage {
								expected := uint32(0)
								if !fill && i < count {
									expected = 1
								}
								if n != expected {
									t.Fatalf("edge usage %d got%d want%d", i, n, expected)
								}
							}
							results = append(results, result{mode, span[0], span[1], count, pos, fill, nilImage, sha256.Sum256(o.bytes()), append([]uint32(nil), o.edgeUsage...)})
						}
					}
				}
			}
			// No grid cell retains these nodes. Detach all definition pointers first.
			for i := 0; i < 2; i++ {
				*memmap.PtrUint32(0x85B3FC, 28676+60*uintptr(i)) = 0
			}
			freeNodes()
			freeHandles()
			copy(o.handles, oldHandles)
			freeImages()
		}
	}
	tileCompositionCapture(t, "overlays", results, "6fd9c980bfebc961818c9916ec2ca6e14bbb15b87985b45d23375a18b3c232ba")
}
