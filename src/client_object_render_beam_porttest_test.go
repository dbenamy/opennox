//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
)

var objectRenderBeamLines = [][4]int32{{20, 48, 76, 48}, {48, 20, 48, 76}, {10, 10, 86, 86}, {86, 10, 10, 86}, {48, 48, 48, 48}, {0, 0, 95, 95}, {96, 48, 160, 48}, {65535, 48, 1, 48}, {25, 75, 75, 25}}

func TestClientObjectRenderBeamQueueAndSight(t *testing.T) {
	o := newObjectRenderOwner(t)
	c := o.c
	polygons := [][]image.Point{nil, {{0, 0}, {96, 0}, {96, 96}, {0, 96}}, {{20, 20}, {76, 20}, {76, 76}, {20, 76}}, {{48, 1}, {95, 48}, {48, 95}, {1, 48}}, {{10, 10}, {86, 10}, {86, 86}, {60, 86}, {60, 35}, {35, 35}, {35, 86}, {10, 86}}}
	type result struct {
		Polygon, Present, Origin, Count, Step int
		Return                                uint32
		Intersections                         []image.Point
		Render                                objectRenderResult
	}
	var out []result
	id := 0
	for poly, points := range polygons {
		for present := 0; present < 2; present++ {
			for origin := 0; origin < 2; origin++ {
				for _, count := range []int{0, 1, 2, 31, 32, 33, 40} {
					id++
					o.resetRender(uint32(id), 120)
					local := o.drawable(7, image.Pt(48, 48))
					if present != 0 {
						*memmap.PtrPtr(0x852978, 8) = local.C()
					}
					if origin != 0 {
						c.Viewport().World.Min = image.Pt(7, 9)
						c.Viewport().Screen.Min = image.Pt(3, 5)
					}
					c.r.Data().SetAlpha(128)
					intersections, restore := c.Cli().PortTestObjectRenderSight(points)
					capture := func(step int, ret uint32) {
						out = append(out, result{poly, present, origin, count, step, ret, intersections(), o.renderResult(t, id, step, int(ret))})
					}
					capture(-1, legacy.PortTestObjectRenderBeam(0, c.Viewport(), [4]int32{}))
					for i := 0; i < count; i++ {
						ret := legacy.PortTestObjectRenderBeam(1, c.Viewport(), objectRenderBeamLines[i%len(objectRenderBeamLines)])
						want := i + 1
						if want > 32 {
							want = 32
						}
						if ret != uint32(want) {
							t.Fatal("beam queue capacity/append result")
						}
						capture(i, ret)
					}
					before := effectsPixelHash(o.pix)
					ret := legacy.PortTestObjectRenderBeam(3, c.Viewport(), [4]int32{})
					want := count
					if want > 32 {
						want = 32
					}
					if present == 0 {
						want = 0
					}
					if ret != uint32(want) {
						t.Fatal("beam walk count")
					}
					if present == 0 && effectsPixelHash(o.pix) != before {
						t.Fatal("beam rendered without local player")
					}
					if poly == 1 && present == 1 && origin == 0 && count == 1 && effectsPixelHash(o.pix) == before {
						t.Fatal("visible beam did not reach actual rasterizer")
					}
					capture(count, ret)
					capture(count+1, legacy.PortTestObjectRenderBeam(2, c.Viewport(), [4]int32{}))
					before = effectsPixelHash(o.pix)
					ret = legacy.PortTestObjectRenderBeam(3, c.Viewport(), [4]int32{})
					if ret != 0 || effectsPixelHash(o.pix) != before {
						t.Fatal("reset beam queue remained active")
					}
					capture(count+2, ret)
					restore()
				}
			}
		}
	}
	effectsCapture(t, "object-render-beam-queue", out, len(out), "992c65aab7901000f68e229c0cea1104e9009b11a740c51153fdb34a20ec2a77")
}
func TestClientObjectRenderBeamRaster(t *testing.T) {
	o := newObjectRenderOwner(t)
	c := o.c
	type result struct {
		Line, Alpha, Enable, Clip int
		Return                    uint32
		Render                    objectRenderResult
	}
	var out []result
	id := 0
	for li, line := range objectRenderBeamLines {
		for _, alpha := range []int{0, 128, 255} {
			for enable := 0; enable < 2; enable++ {
				for clip := 0; clip < 3; clip++ {
					id++
					o.resetRender(uint32(id), 120)
					legacy.PortTestObjectRenderBeam(0, c.Viewport(), [4]int32{})
					c.r.Data().SetAlpha(byte(alpha))
					c.r.Data().SetAlphaEnabled(enable != 0)
					if clip == 1 {
						c.r.Data().SetClipRect(image.Rect(20, 20, 76, 76))
						c.r.Data().SetClipRect2(image.Rect(20, 20, 75, 75))
					}
					if clip == 2 {
						c.r.Data().SetClipRect(image.Rect(1, 1, 2, 2))
						c.r.Data().SetClipRect2(image.Rect(1, 1, 1, 1))
					}
					ret := legacy.PortTestObjectRenderBeam(4, c.Viewport(), line)
					out = append(out, result{li, alpha, enable, clip, ret, o.renderResult(t, id, 0, int(ret))})
				}
			}
		}
	}
	effectsCapture(t, "object-render-beam-raster", out, len(out), "f82a5f22319c14b9e2b164c5d8b25ce82d03941c521d93bf5e47be9f1c2fe0fd")
}
