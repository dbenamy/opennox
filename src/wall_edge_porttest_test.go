//go:build porttest

package opennox

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	noxcolor "github.com/opennox/libs/color"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"os"
	"testing"
)

func wallEdgeCapture(t *testing.T, label string, rows any, want string) {
	t.Helper()
	raw, err := json.Marshal(rows)
	if err != nil {
		t.Fatal(err)
	}
	if prefix := os.Getenv("OPENNOX_WALL_EDGE_CAPTURE"); prefix != "" {
		if err := os.WriteFile(prefix+"-"+label+".json", raw, 0600); err != nil {
			t.Fatal(err)
		}
	}
	got := fmt.Sprintf("%x", sha256.Sum256(raw))
	t.Logf("%s: %s", label, got)
	// Expectations captured from original production C.
	if got != want {
		t.Fatalf("%s hash %s want %s", label, got, want)
	}
}
func wallEdgePixel(x, y int) uint16 {
	return noxcolor.RGB5551Color(byte(40+x*9), byte(70+y*23), byte(220-x*3-y*7)).Color16()
}
func wallEdgeOpaque(pattern, x, y int) bool {
	switch pattern {
	case 0:
		return true
	case 1:
		return false
	case 2:
		return (x/3+y)%2 == 0
	case 3:
		return x >= 5 && x < 19
	default:
		return x < 4 || x >= 17
	}
}
func wallEdgeImage(pattern, w, h int, off image.Point) []byte {
	b := make([]byte, 17)
	binary.LittleEndian.PutUint32(b, uint32(w))
	binary.LittleEndian.PutUint32(b[4:], uint32(h))
	binary.LittleEndian.PutUint32(b[8:], uint32(off.X))
	binary.LittleEndian.PutUint32(b[12:], uint32(off.Y))
	for y := 0; y < h; y++ {
		for x := 0; x < w; {
			opaque := wallEdgeOpaque(pattern, x, y)
			end := x + 1
			for end < w && end-x < 255 && wallEdgeOpaque(pattern, end, y) == opaque {
				end++
			}
			op := byte(1)
			if opaque {
				op = 2
			}
			b = append(b, op, byte(end-x))
			if opaque {
				for i := x; i < end; i++ {
					b = binary.LittleEndian.AppendUint16(b, wallEdgePixel(i, y))
				}
			}
			x = end
		}
	}
	return b
}
func wallEdgeImages(t *testing.T, o *worldWallOwner, raw [][]byte, types []int) []*noxrender.Image {
	imgs, restore := o.c.r.GetBag().PortTestWallEdgeImages(raw, types)
	t.Cleanup(restore)
	old := partViewportOff
	t.Cleanup(func() { partViewportOff = old })
	return imgs
}
func wallEdgePrepare(t *testing.T, o *worldWallOwner, clip image.Rectangle, high, parity int) {
	o.resetWalls(t)
	partViewportOff = image.Pt(0, parity)
	*o.words["highFront"] = uint32(high)
	*o.words["edgeMinX"], *o.words["edgeMinY"] = uint32(clip.Min.X), uint32(clip.Min.Y)
	*o.words["edgeMaxX"], *o.words["edgeMaxY"] = uint32(clip.Max.X), uint32(clip.Max.Y)
	// Distinct rows reveal accidental copies across transparent runs.
	for i := range o.pix.Pix {
		o.pix.Pix[i] = uint16(0x8000 | ((i*29 + 17) & 0x7fff))
	}
}
func TestWallEdgeMatrix(t *testing.T) {
	o := newWorldWallOwner(t)
	var raw [][]byte
	var types []int
	for pattern := 0; pattern < 5; pattern++ {
		raw = append(raw, wallEdgeImage(pattern, 24, 7, image.Pt(-2, 1)))
		types = append(types, 3)
	}
	imgs := wallEdgeImages(t, o, raw, types)
	type result struct {
		Pattern, Scene, Light, High, Parity int
		Pixels                              string
		State                               objectRenderResult
	}
	var rows []result
	lights := [][2][3]uint32{{{256, 256, 256}, {256, 256, 256}}, {{255, 7, 128}, {0, 249, 73}}, {{0, 0, 0}, {255, 255, 255}}}
	for pattern := range imgs {
		for scene := 0; scene < 12; scene++ {
			for li, light := range lights {
				for high := 0; high < 2; high++ {
					for parity := 0; parity < 2; parity++ {
						clip := image.Rect(0, 0, 96, 96)
						pos := image.Pt(22, 19)
						left, right, bottom, crop := 0, 96, 96, 0
						switch scene {
						case 1:
							clip = image.Rect(25, 0, 96, 96)
						case 2:
							clip = image.Rect(0, 0, 37, 96)
						case 3:
							clip = image.Rect(26, 22, 39, 25)
						case 4:
							left, right = 27, 37
						case 5:
							clip = image.Rect(23, 0, 96, 96)
							left, right = 31, 39
						case 6:
							bottom = 24
						case 7:
							crop = 3
						case 8:
							crop = 7
						case 9:
							pos = image.Pt(99, 19)
						case 10:
							bottom = 20
						case 11:
							clip = image.Rect(0, 25, 96, 96)
						}
						wallEdgePrepare(t, o, clip, high, parity)
						before := append([]uint16(nil), o.pix.Pix...)
						first, second := light[0], light[1]
						legacy.PortTestWallEdge(imgs[pattern].C(), pos, &first, &second, bottom, left, right, crop, 0x0f)
						if first != light[0] || second != light[1] {
							t.Fatal("edge draw changed input light colors")
						}
						if !bytes.Equal(raw[pattern], imgs[pattern].Pixdata()) {
							t.Fatal("edge draw changed image data")
						}
						if scene == 8 || scene == 9 || scene == 10 || pattern == 1 && high == 1 {
							for i, p := range before {
								if o.pix.Pix[i] != p {
									t.Fatalf("no-op case changed pixel: pattern%d scene%d high%d at%d", pattern, scene, high, i)
								}
							}
						}
						if high == 1 {
							for y := 0; y < 96; y++ {
								for x := 0; x < 96; x++ {
									if !image.Pt(x, y).In(clip) && o.pix.Pix[y*o.pix.Stride+x] != before[y*o.pix.Stride+x] {
										t.Fatalf("outside clip pattern%d scene%d at%d,%d", pattern, scene, x, y)
									}
								}
							}
						}
						rows = append(rows, result{pattern, scene, li, high, parity, effectsPixelHash(o.pix), o.renderResult(t, pattern, scene, 0)})
					}
				}
			}
		}
	}
	wallEdgeCapture(t, "matrix", rows, "9c5fe9841d7d4c0d6e7df731292b21a4b81c81ee5fb18f5f726fb9a612907e40")
}
func TestWallEdgePixelContracts(t *testing.T) {
	o := newWorldWallOwner(t)
	var raw [][]byte
	var types []int
	for pattern := 0; pattern < 5; pattern++ {
		raw = append(raw, wallEdgeImage(pattern, 24, 7, image.Point{}))
		types = append(types, 3)
	}
	imgs := wallEdgeImages(t, o, raw, types)
	type result struct {
		Pattern, Light int
		Pixels         string
	}
	var rows []result
	lights := [][2][3]uint32{{{256, 256, 256}, {256, 256, 256}}, {{0, 0, 0}, {0, 0, 0}}, {{255, 7, 128}, {0, 249, 73}}, {{17, 199, 33}, {251, 0, 240}}}
	for pattern, img := range imgs {
		for li, light := range lights {
			wallEdgePrepare(t, o, image.Rect(0, 0, 96, 96), 1, 0)
			before := append([]uint16(nil), o.pix.Pix...)
			first, second := light[0], light[1]
			legacy.PortTestWallEdge(img.C(), image.Pt(20, 20), &first, &second, 96, 0, 96, 0, 1)
			for y := 0; y < 96; y++ {
				for x := 0; x < 96; x++ {
					want := before[y*o.pix.Stride+x]
					sx, sy := x-20, y-20
					if sx >= 0 && sx < 24 && sy >= 0 && sy < 7 && wallEdgeOpaque(pattern, sx, sy) {
						c := noxcolor.RGBA5551(wallEdgePixel(sx, sy)).ColorNRGBA()
						channel := func(v byte, i int) byte {
							start := int64(light[0][i]) * 256
							step := (int64(light[1][i])*256 - start) / 24
							return byte(((start + int64(sx)*step) * int64(v)) >> 16)
						}
						want = noxcolor.RGB5551Color(channel(c.R, 0), channel(c.G, 1), channel(c.B, 2)).Color16()
					}
					if got := o.pix.Pix[y*o.pix.Stride+x]; got != want {
						t.Fatalf("pattern%d light%d pixel%d,%d got%04x want%04x", pattern, li, x, y, got, want)
					}
				}
			}
			rows = append(rows, result{pattern, li, effectsPixelHash(o.pix)})
		}
	}
	wallEdgeCapture(t, "pixel-contracts", rows, "52fbaff56448c3f3f29fe0786ca4458ede98aabf7fb6acca8012492815a5c507")
}
func TestWallEdgeGuardsAndRuns(t *testing.T) {
	o := newWorldWallOwner(t)
	var raw [][]byte
	types := []int{3, 0, 2, 4, 7, 0x43, 3, 3, 3}
	widths := []int{24, 24, 24, 24, 24, 24, 1, 255, 260}
	for _, w := range widths {
		raw = append(raw, wallEdgeImage(0, w, 3, image.Point{}))
	}
	imgs := wallEdgeImages(t, o, raw, types)
	type result struct {
		Image, Flags, High int
		Pixels             string
	}
	var rows []result
	for n := -1; n < len(imgs); n++ {
		for _, flags := range []int{0, 1, 2, 4, 8, 15, -1} {
			for high := 0; high < 2; high++ {
				wallEdgePrepare(t, o, image.Rect(0, 0, 96, 96), high, 0)
				before := effectsPixelHash(o.pix)
				var img noxrender.ImageHandle
				if n >= 0 {
					img = imgs[n].C()
				}
				first, second := [3]uint32{256, 256, 256}, [3]uint32{256, 256, 256}
				legacy.PortTestWallEdge(img, image.Pt(4, 4), &first, &second, 96, 0, 90, 0, flags)
				got := effectsPixelHash(o.pix)
				if (n < 0 || types[n]&63 != 3) && got != before {
					t.Fatalf("nil/unsupported image%d drew pixels", n)
				}
				rows = append(rows, result{n, flags, high, got})
			}
		}
	}
	// The original flags argument is unused: compare every flag value per image/mode.
	for i, r := range rows {
		if i%14 >= 2 && r.Pixels != rows[i-i%14+i%2].Pixels {
			t.Fatalf("unused flags changed image%d high%d", r.Image, r.High)
		}
	}
	wallEdgeCapture(t, "guards-runs", rows, "2fa1919602f31736fcddce1e22c52f23a594ae9b22d713e6af3b4810bd60e577")
}
