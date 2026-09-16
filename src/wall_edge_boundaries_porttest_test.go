//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
)

func TestWallEdgeBoundaries(t *testing.T) {
	o := newWorldWallOwner(t)
	raw := wallEdgeImage(0, 12, 8, image.Pt(-2, -3))
	imgs := wallEdgeImages(t, o, [][]byte{raw}, []int{3})
	type result struct {
		Position       image.Point
		Clip           image.Rectangle
		Crop, Interval int
		Pixels         string
	}
	var rows []result
	positions := []image.Point{{0, 0}, {2, 3}, {-10, 3}, {3, -5}, {96, 96}, {98, 99}, {30, 30}, {18, 33}, {32, 18}}
	clips := []image.Rectangle{image.Rect(0, 0, 96, 96), image.Rect(20, 20, 30, 30), image.Rect(20, 20, 20, 20)}
	for _, pos := range positions {
		for _, clip := range clips {
			for _, crop := range []int{0, 7, 8, 9} {
				for interval := 0; interval < 4; interval++ {
					wallEdgePrepare(t, o, clip, 1, 0)
					before := append([]uint16(nil), o.pix.Pix...)
					left, right := 0, 96
					switch interval {
					case 1:
						left, right = 20, 25
					case 2:
						left, right = 25, 25
					case 3:
						left, right = 31, 20
					}
					a, b := [3]uint32{256, 256, 256}, [3]uint32{256, 256, 256}
					legacy.PortTestWallEdge(imgs[0].C(), pos, &a, &b, 96, left, right, crop, 1)
					for y := 0; y < 96; y++ {
						for x := 0; x < 96; x++ {
							got, old := o.pix.Pix[y*o.pix.Stride+x], before[y*o.pix.Stride+x]
							if (crop >= 8 || interval >= 2 || clip.Empty() || !image.Pt(x, y).In(clip)) && got != old {
								t.Fatalf("pos%v clip%v crop%d interval%d changed pixel%d,%d", pos, clip, crop, interval, x, y)
							}
						}
					}
					rows = append(rows, result{pos, clip, crop, interval, effectsPixelHash(o.pix)})
				}
			}
		}
	}
	wallEdgeCapture(t, "boundaries", rows, "019ff8ac60dcb8d7f55ab32ab32741d047d323013b82ce1311a5a1b24845652a")
}
func TestWallEdgeAlternateRows(t *testing.T) {
	o := newWorldWallOwner(t)
	imgs := wallEdgeImages(t, o, [][]byte{wallEdgeImage(0, 12, 8, image.Point{})}, []int{3})
	type result struct {
		Y, Parity, High int
		Pixels          string
	}
	var rows []result
	for _, y0 := range []int{20, 21} {
		for _, parity := range []int{-3, -2, -1, 0, 1, 2, 3} {
			for high := 0; high < 2; high++ {
				wallEdgePrepare(t, o, image.Rect(0, 0, 96, 96), high, parity)
				before := append([]uint16(nil), o.pix.Pix...)
				a, b := [3]uint32{256, 256, 256}, [3]uint32{256, 256, 256}
				legacy.PortTestWallEdge(imgs[0].C(), image.Pt(20, y0), &a, &b, 96, 0, 96, 0, 4)
				for y := 0; y < 96; y++ {
					for x := 0; x < 96; x++ {
						want := before[y*o.pix.Stride+x]
						if x >= 20 && x < 32 && y >= y0 && y < y0+8 {
							sy := y - y0
							if high == 0 && (y+parity)&1 != 0 {
								sy--
							}
							if sy >= 0 {
								want = wallEdgePixel(x-20, sy)
							}
						}
						if got := o.pix.Pix[y*o.pix.Stride+x]; got != want {
							t.Fatalf("y%d parity%d high%d at%d,%d got%x want%x", y0, parity, high, x, y, got, want)
						}
					}
				}
				rows = append(rows, result{y0, parity, high, effectsPixelHash(o.pix)})
			}
		}
	}
	wallEdgeCapture(t, "alternate-rows", rows, "b142b0464cb1c7c6ac78e03df23c86705863f95256708bb863a7d1357a4f4650")
}
