//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"image"
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/opennox/libs/object"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

func TestMapDrawableDoorCoordinates(t *testing.T) {
	mapDrawableTables(t)
	c, _, _ := newEffectsFullOwner(t)
	t.Cleanup(noxflags.PortTestGameFlags(1))
	t.Cleanup(c.srv.PortTestMinimapWalls())
	words, restore := legacy.PortTestMinimapWords()
	t.Cleanup(restore)
	*words["polygons"] = 1
	original := cryptfile.Global()
	cryptfile.SetGlobal(nil)
	t.Cleanup(func() { cryptfile.Close(); cryptfile.SetGlobal(original) })
	path := filepath.Join(t.TempDir(), "door.bin")
	c.Things.TypeByInd(4).ObjClass = object.Class(0x80)
	type record struct {
		Position  image.Point
		Direction int
		Found     bool
		Grid      image.Point
		Bytes     int
	}
	var records []record
	for _, pos := range []image.Point{{0, 0}, {0, 64}, {64, 0}, {10, 10}, {11, 11}, {12, 12}, {22, 22}, {23, 23}, {64, 128}} {
		for _, dir := range []int{0, 8, 16, 24} {
			t.Run(fmt.Sprintf("%d-%d-dir%d", pos.X, pos.Y, dir), func(t *testing.T) {
				var s mapDrawableStream
				s.u16(41)
				s.modernBase(64, false, 0, 0)
				s.u32(0)
				s.u32(0)
				s.u32(uint32(dir))
				binary.LittleEndian.PutUint32(s.Bytes()[12:], math.Float32bits(float32(pos.X)))
				binary.LittleEndian.PutUint32(s.Bytes()[16:], math.Float32bits(float32(pos.Y)))
				if err := os.WriteFile(path, s.Bytes(), 0600); err != nil {
					t.Fatal(err)
				}
				if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
					t.Fatal(err)
				}
				defer cryptfile.Close()
				ret := legacy.PortTestMapDrawableRecord(1, 4)
				if ret != s.Len() {
					t.Fatalf("bytes=%d want=%d", ret, s.Len())
				}
				dr := c.Objs.List1
				if dr == nil {
					t.Fatal("missing real drawable")
				}
				r := record{Position: pos, Direction: dir, Bytes: ret}
				c.srv.Walls.EachWallRaw(func(w *server.Wall) bool {
					if w.Field32 == uint32(uintptr(dr.C())) {
						r.Found = true
						r.Grid = image.Pt(int(w.X5), int(w.Y6))
						w.Field32 = 0
					}
					return true
				})
				c.Nox_xxx_spriteDeleteStatic_45A4E0_drawable(dr)
				dx, dy := -11, -11
				if dir == 8 || dir == 16 {
					dx = 11
				}
				if dir == 16 || dir == 24 {
					dy = 11
				}
				// X truncates as signed; Y converts to uint32 before division. A negative
				// adjusted Y is outside the wall grid, while adjusted X -11..-1 becomes 0.
				wantFound := pos.Y+dy >= 0
				want := image.Pt((pos.X+dx)/23, (pos.Y+dy)/23)
				if r.Found != wantFound || r.Found && r.Grid != want {
					t.Fatalf("wall %+v want found=%v grid=%v", r, wantFound, want)
				}
				records = append(records, r)
			})
		}
	}
	spellbookCapture(t, "map-drawable-doors", records, "bfbff3b9940128653de3aeeb077be18d7fd96730325f9174329a4f4076dac9f2")
}
