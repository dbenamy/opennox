//go:build porttest

package opennox

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"image"
	"testing"
	"unsafe"
)

func TestTileRasterPrimitives(t *testing.T) {
	o := newTileRasterOwner(t)
	type result struct {
		Fill           bool
		Stride, Offset int
		Color          uint32
		Buffer         [32]byte
	}
	var rows []result
	colors := []uint32{0, 0xffffffff, 0x1234abcd, 0xabcd1234, 0x80008000, 0x001f7c00}
	for _, stride := range []int{96, 97, 128, 256} {
		for _, offset := range []int{0, 1, 2, 3, 17} {
			for _, color := range colors {
				for _, fill := range []bool{false, true} {
					*o.words["stride"] = uint32(stride)
					b, free := alloc.Make([]byte{}, offset+stride*46+32)
					defer free()
					for i := range b {
						b[i] = byte(i*17 + 99)
					}
					want := append([]byte(nil), b...)
					source := append([]byte(nil), o.raw...)
					for y := 0; y < 46; y++ {
						x, width, src := tileRasterRow(y)
						dst := want[offset+y*stride+2*x:][:2*width]
						if fill {
							tileRasterFastFill(dst, x, width, color)
						} else {
							copy(dst, source[src:src+2*width])
						}
					}
					legacy.PortTestTileRasterPrimitive(fill, unsafe.Pointer(&b[offset]), unsafe.Pointer(&source[0]), color)
					if !bytes.Equal(b, want) {
						for i := range b {
							if b[i] != want[i] {
								t.Fatalf("fill%v stride%d offset%d color%x byte%d got%x want%x", fill, stride, offset, color, i, b[i], want[i])
							}
						}
					}
					if !bytes.Equal(source, o.raw) {
						t.Fatal("tile primitive changed source")
					}
					rows = append(rows, result{fill, stride, offset, color, sha256.Sum256(b)})
				}
			}
		}
	}
	tileRasterCapture(t, "primitives", rows, "226ad32567debd8455361abf42a4a0e69810d5e163206618479aa548163992b7")
}
func TestTileRasterCallbackAndRing(t *testing.T) {
	o := newTileRasterOwner(t)
	type result struct {
		Fill        bool
		Position    image.Point
		Color, Tile uint32
		Flag        byte
		Buffer      [32]byte
		Flat, Dirty uint32
	}
	var rows []result
	positions := []image.Point{{0, 0}, {31, 17}, {0, 184}, {0, 185}, {230, 184}, {275, 229}, {0, 230}, {31, 450}}
	for _, fill := range []bool{false, true} {
		for _, pos := range positions {
			for _, color := range []uint32{0, 0xffffffff, 0x1234abcd, 0x80008000} {
				for _, flag := range []byte{0, 1, 2, 3} {
					for _, tile := range []uint32{0, 0xbeef0000} {
						o.resetBuffer()
						o.defs[0].Color48 = color
						o.defs[0].Field58 = flag
						nox_client_texturedFloors_154956 = !fill
						nox_xxx_tileSetDrawFn_481420()
						b := o.bytes()
						want := append([]byte(nil), b...)
						stride := int(*o.words["stride"])
						base := (pos.Y*stride + pos.X*2) % len(b)
						fast := base+int(*memmap.PtrUint32(0x973CE0, 376)) < len(b)
						for y := 0; y < 46; y++ {
							x, width, src := tileRasterRow(y)
							start := (base + y*stride + 2*x) % len(b)
							n := 2 * width
							if !fill {
								for i := 0; i < n; i++ {
									want[(start+i)%len(b)] = o.raw[src+i]
								}
							} else if fast {
								tileRasterFastFill(want[start:start+n], x, width, color)
							} else {
								var pattern [4]byte
								binary.LittleEndian.PutUint32(pattern[:], color)
								for i := 0; i < n; i++ {
									index := start + i
									phase := i
									if index >= len(b) {
										index -= len(b)
										phase = i - (len(b) - start)
									}
									want[index] = pattern[phase%4]
								}
							}
						}
						legacy.PortTestTileRasterDispatch(pos, o.images[0].C(), tile)
						if !bytes.Equal(b, want) {
							for i := range b {
								if b[i] != want[i] {
									t.Fatalf("fill%v pos%v fast%v color%x tile%x byte%d got%x want%x", fill, pos, fast, color, tile, i, b[i], want[i])
								}
							}
						}
						expectedFlag := uint32(0)
						if fill && flag&1 != 0 {
							expectedFlag = 1
						}
						if *o.words["flatFlag"] != expectedFlag || *o.words["dirty"] != 1 {
							t.Fatal("callback configuration flags differ")
						}
						if !bytes.Equal(o.images[0].Pixdata(), o.raw) {
							t.Fatal("callback changed image data")
						}
						rows = append(rows, result{fill, pos, color, tile, flag, sha256.Sum256(b), *o.words["flatFlag"], *o.words["dirty"]})
					}
				}
			}
		}
	}
	tileRasterCapture(t, "callback-ring", rows, "8b2d59847b69464a4b12bb8a6966980f920c0bd16a63de3b2a6f7e582df8f857")
}
func TestTileRasterWrapSetup(t *testing.T) {
	o := newTileRasterOwner(t)
	type result struct {
		Stride, Dirty, Threshold uint32
		Shift                    byte
	}
	var rows []result
	for _, stride := range []uint32{0, 1, 92, 552, 0x10000000, 0x7fffffff, 0xffffffff} {
		for _, shift := range []byte{0, 1, 2} {
			*o.words["stride"] = stride
			*o.words["dirty"] = 7
			*memmap.PtrUint32(0x973CE0, 376) = 0xdeadbeef
			*memmap.PtrUint8(0x973F18, 7696) = shift
			legacy.Nox_xxx_tile_486060()
			got := *memmap.PtrUint32(0x973CE0, 376)
			want := uint32(uint64(stride)*45 + (uint64(46) << shift))
			if got != want || *o.words["dirty"] != 1 {
				t.Fatal("tile wrap threshold or dirty flag differs")
			}
			rows = append(rows, result{stride, *o.words["dirty"], got, shift})
		}
	}
	tileRasterCapture(t, "wrap-setup", rows, "50c2681ec7661f5d919bdfe9a7b2e63a98e3a620ec8d2df633846724cf812129")
}
