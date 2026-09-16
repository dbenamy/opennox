//go:build porttest

package opennox

import (
	"fmt"
	"image"
	"io"
	"os"
	"path/filepath"
	"testing"
	"unsafe"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
)

func TestMapDrawableOldDefaults(t *testing.T) {
	c, _, _ := newEffectsFullOwner(t)
	t.Cleanup(noxflags.PortTestGameFlags(1))
	original := cryptfile.Global()
	cryptfile.SetGlobal(nil)
	t.Cleanup(func() { cryptfile.Close(); cryptfile.SetGlobal(original) })
	path := filepath.Join(t.TempDir(), "old-object.bin")
	for _, outer := range []int16{0, 9, 10, 19, 20, 29, 30, 39, 40} {
		for _, inner := range []int16{1, 2, 3, 4, 5} {
			if outer >= 40 && inner >= 3 {
				continue
			}
			t.Run(fmt.Sprintf("outer%d-inner%d", outer, inner), func(t *testing.T) {
				var s mapDrawableStream
				s.oldBase(outer, inner, 0, 0)
				if err := os.WriteFile(path, s.Bytes(), 0600); err != nil {
					t.Fatal(err)
				}
				if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
					t.Fatal(err)
				}
				defer cryptfile.Close()
				var count uint32
				dr := legacy.PortTestMapDrawableBase(4, outer, inner, true, &count)
				if dr == nil {
					t.Fatal("real drawable allocation failed")
				}
				defer func() { dr.TeamVal.ID = 0; c.Nox_xxx_spriteDeleteStatic_45A4E0_drawable(dr) }()
				pos, err := cryptfile.Global().File.Seek(0, io.SeekCurrent)
				if err != nil {
					t.Fatal(err)
				}
				if int(count) != s.Len() || pos != int64(s.Len()) {
					t.Fatalf("consumed=%d position=%d want=%d", count, pos, s.Len())
				}
				if dr.PosVec != image.Pt(64, 128) {
					t.Fatalf("position=%v", dr.PosVec)
				}
				extra := *(*uint32)(unsafe.Add(unsafe.Pointer(dr), 280))
				if dr.TeamVal.ID != 0 || extra != 0 {
					t.Errorf("absent legacy fields must default to zero: team=%d extra=%08x", dr.TeamVal.ID, extra)
				}
			})
		}
	}
}
