//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"os"
	"path/filepath"
	"testing"
	"unsafe"
)

func TestMapDrawableLightBoundaries(t *testing.T) {
	mapDrawableTables(t)
	c, _, _ := newEffectsFullOwner(t)
	t.Cleanup(noxflags.PortTestGameFlags(1))
	original := cryptfile.Global()
	cryptfile.SetGlobal(nil)
	t.Cleanup(func() { cryptfile.Close(); cryptfile.SetGlobal(original) })
	path := filepath.Join(t.TempDir(), "light.bin")
	type record struct {
		Outer          int16
		Intensity      uint32
		Fixed          int32
		Angle          uint16
		Output, Radius uint32
		AngleOut       uint16
		Bytes          int
	}
	var records []record
	for _, outer := range []int16{1, 2, 41, 42} {
		for _, intensity := range []uint32{0, math.Float32bits(20), math.Float32bits(63), math.Float32bits(80), 0x7fc00001, 0x7f800000, 0xff800000} {
			for _, fixed := range []int32{-1, 0, (63 << 16) - 1, 63 << 16, (63 << 16) + 1, math.MaxInt32} {
				for _, angle := range []uint16{0, 1, 181, 182, 32767, 32768, 65535} {
					t.Run(fmt.Sprintf("v%d-float%x-fixed%d-angle%d", outer, intensity, fixed, angle), func(t *testing.T) {
						var s mapDrawableStream
						s.u16(uint16(outer))
						if outer < 40 {
							s.oldBase(outer, 1, 0, 0)
						} else {
							s.modernBase(64, false, 0, 0)
						}
						base := s.Len()
						s.light(outer, 0)
						raw := s.Bytes()
						binary.LittleEndian.PutUint32(raw[base+4:], intensity)
						binary.LittleEndian.PutUint32(raw[base+12:], uint32(fixed))
						binary.LittleEndian.PutUint16(raw[base+28:], angle)
						if err := os.WriteFile(path, raw, 0600); err != nil {
							t.Fatal(err)
						}
						if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
							t.Fatal(err)
						}
						defer cryptfile.Close()
						ret := legacy.PortTestMapDrawableRecord(2, 4)
						if ret != s.Len() {
							t.Fatalf("consumed=%d want=%d", ret, s.Len())
						}
						dr := c.Objs.List1
						if dr == nil {
							t.Fatal("missing actual drawable")
						}
						defer func() { dr.TeamVal.ID = 0; c.Nox_xxx_spriteDeleteStatic_45A4E0_drawable(dr) }()
						p := unsafe.Slice((*byte)(dr.C()), int(unsafe.Sizeof(*dr)))
						out := binary.LittleEndian.Uint32(p[140:])
						radius := binary.LittleEndian.Uint32(p[144:])
						a := binary.LittleEndian.Uint16(p[268:])
						clamp := outer < 2 && (!(math.Float32frombits(intensity) <= 63) || fixed > 63<<16)
						want := intensity
						if clamp {
							want = math.Float32bits(63)
						}
						if out != want {
							t.Fatalf("intensity=%x want=%x clamp=%v", out, want, clamp)
						}
						if !clamp && radius != 11 {
							t.Fatal("unclamped serialized radius changed")
						}
						// The on-disk 16-bit turn value maps to truncated whole degrees.
						wantAngle := uint16(uint32(angle) * 360 / 65536)
						if a != wantAngle {
							t.Fatalf("angle=%d want=%d", a, wantAngle)
						}
						records = append(records, record{outer, intensity, fixed, angle, out, radius, a, ret})
					})
				}
			}
		}
	}
	spellbookCapture(t, "map-drawable-lights", records, "f86e84c8f103bc7ed3dc12739c33d4bb886a33e0d8d931fb58fcdd87dfe949c4")
}
