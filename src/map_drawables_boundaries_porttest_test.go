//go:build porttest

package opennox

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
)

func TestMapDrawableCountBoundaries(t *testing.T) {
	c, _, _ := newEffectsFullOwner(t)
	t.Cleanup(noxflags.PortTestGameFlags(1))
	original := cryptfile.Global()
	cryptfile.SetGlobal(nil)
	t.Cleanup(func() { cryptfile.Close(); cryptfile.SetGlobal(original) })
	path := filepath.Join(t.TempDir(), "counts.bin")
	type record struct {
		Inner             int16
		Count             uint16
		Initial, Consumed uint32
		Position          int64
		Bytes             int
	}
	var records []record
	for _, inner := range []int16{2, 4, 5, 60, 61, 63, 64} {
		for _, count := range []uint16{0, 1, 2, 16383, 16384, 16385, 32767, 65535} {
			for _, initial := range []uint32{0, 0xfffffff0, 0xffffffff} {
				t.Run(fmt.Sprintf("v%d-count%d-initial%x", inner, count, initial), func(t *testing.T) {
					var s mapDrawableStream
					if inner >= 61 {
						s.modernBaseCount(inner, true, 0, 0, count)
					} else {
						s.u16(uint16(inner))
						s.oldBaseCount(40, inner, 0, 0, count)
					}
					if err := os.WriteFile(path, s.Bytes(), 0600); err != nil {
						t.Fatal(err)
					}
					if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
						t.Fatal(err)
					}
					defer cryptfile.Close()
					consumed := initial
					dr := legacy.PortTestMapDrawableBase(4, 40, inner, false, &consumed)
					if dr == nil {
						t.Fatal("actual drawable allocation failed")
					}
					defer func() { dr.TeamVal.ID = 0; c.Nox_xxx_spriteDeleteStatic_45A4E0_drawable(dr) }()
					pos, err := cryptfile.Global().File.Seek(0, io.SeekCurrent)
					if err != nil {
						t.Fatal(err)
					}
					if consumed != initial+uint32(s.Len()) || pos != int64(s.Len()) {
						t.Fatalf("consumed=%x position=%d want count=%x position=%d", consumed, pos, initial+uint32(s.Len()), s.Len())
					}
					records = append(records, record{inner, count, initial, consumed, pos, s.Len()})
				})
			}
		}
	}
	spellbookCapture(t, "map-drawable-counts", records, "3e7317fc0efe59950be7af5719ffa0b0bb87d47b49b89e3ce12c93ce5c8b6b38")
}
