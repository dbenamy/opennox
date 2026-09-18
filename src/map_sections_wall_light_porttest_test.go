//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestMapSectionsWallPolygonLight(t *testing.T) {
	var out []legacy.PortTestMapSectionResult
	for _, where := range []string{"inside", "edge", "outside"} {
		for _, level := range []byte{0, 173} {
			for _, region := range []string{"none", "inside", "outside"} {
				t.Run(fmt.Sprintf("%s/level%d/region%s", where, level, region), func(t *testing.T) {
					owner := newMapPolygonsOwner(t)
					x0, x1 := float32(180), float32(210)
					y0, y1 := float32(230), float32(260)
					if where == "edge" {
						x0, x1 = 170, 190
					}
					if where == "outside" {
						x0, x1 = 10, 20
						y0, y1 = 10, 20
					}
					owner.construct(t, [][2]float32{{x0, y0}, {x1, y0}, {x1, y1}, {x0, y1}})
					owner.polygon(1)[130] = level
					sp := legacy.PortTestMapSectionSpec{Name: t.Name(), Paint: legacy.PortTestPaintSpec{Seed: 83, Globals: map[string]legacy.PortTestMapRoomArg{"section-magic-wall": {Value: 2}}, Walls: []legacy.PortTestPaintWall{{X: 8, Y: 10, Words: map[int]uint32{0: 0x00030102}}}, Actions: []legacy.PortTestPaintAction{{Op: 0}}}, IO: []legacy.PortTestMapSectionIO{{Function: "walls"}}}
					if region != "none" {
						x := uint32(8)
						if region == "outside" {
							x = 30
						}
						quad := [8]uint32{(x + 2) * 23, 8 * 23, x * 23, 10 * 23, (x + 4) * 23, 10 * 23, (x + 2) * 23, 12 * 23}
						words := map[int]uint32{}
						for i, v := range quad {
							words[4*i] = v
						}
						sp.Paint.Records = []legacy.PortTestMapRoomRecord{{Size: 32, Words: words}}
						sp.Paint.Actions[0].Args[0] = legacy.PortTestMapRoomArg{Slot: 1}
					}
					r := mapSectionsRun(t, []legacy.PortTestMapSectionSpec{sp})[0]
					out = append(out, r)
					wire := r.IO[0]
					if wire.Return != 1 {
						t.Fatal("wall writer result")
					}
					if region == "outside" {
						if len(wire.Data) != 19 {
							t.Fatalf("excluded wall bytes %x", wire.Data)
						}
						return
					}
					want := level
					if where == "outside" {
						want = 100
					}
					if len(wire.Data) != 26 || wire.Data[23] != want {
						t.Fatalf("wall polygon light bytes %x want%d", wire.Data, want)
					}
				})
			}
		}
	}
	spellbookCapture(t, "map-sections-wall-polygon-light", out, "e7a5286a956195936797f543b3612617af81de5d862303c8e2ce8e79f8fb3dff")
}
