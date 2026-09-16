//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
)

func TestWorldWallsImageIntervals(t *testing.T) {
	o := newWorldWallOwner(t)
	type result struct {
		Clip          image.Rectangle
		Input         image.Point
		Initial, Flag uint32
		Return        int
		Bounds        [2]int32
	}
	var rows []result
	points := []int{-10, 0, 19, 20, 21, 40, 69, 70, 71, 96, 110}
	for _, clip := range []image.Rectangle{image.Rect(0, 0, 96, 96), image.Rect(20, 10, 70, 80)} {
		for _, a := range points {
			for _, b := range points {
				for _, initial := range []uint32{0, 7} {
					o.resetWalls(t)
					o.c.r.Data().SetClipRect(clip)
					*o.words["imageClip"] = initial
					*memmap.PtrInt32(0x973F18, 52) = 888
					*memmap.PtrInt32(0x973F18, 12) = 777
					ret, _ := legacy.PortTestWorldWalls(8, nil, nil, nil, image.Pt(a, b))
					lo, hi := max(min(a, b), clip.Min.X), min(max(a, b), clip.Max.X)
					visible := lo < hi
					if (ret != 0) != visible {
						t.Fatalf("interval %d,%d in%v returned%d", a, b, clip, ret)
					}
					flag := initial
					bounds := [2]int32{888, 777}
					if visible && (lo != clip.Min.X || hi != clip.Max.X) {
						flag = 1
						bounds = [2]int32{int32(lo), int32(hi)}
					}
					got := [2]int32{memmap.Int32(0x973F18, 52), memmap.Int32(0x973F18, 12)}
					if *o.words["imageClip"] != flag || got != bounds || o.c.r.Data().ClipRect() != clip {
						t.Fatal("image clip interval state mismatch")
					}
					rows = append(rows, result{clip, image.Pt(a, b), initial, *o.words["imageClip"], ret, got})
				}
			}
		}
	}
	worldWallsCapture(t, "image-intervals", rows, "0d151547bb49a85dba116b14df3ecbedb133c2c746611f3d127232a9ad95a96e")
}
