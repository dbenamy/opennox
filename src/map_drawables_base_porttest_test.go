//go:build porttest

package opennox

import (
	"encoding/binary"
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

func TestMapDrawableBaseRecords(t *testing.T) {
	c, _, _ := newEffectsFullOwner(t)
	t.Cleanup(noxflags.PortTestGameFlags(1))
	original := cryptfile.Global()
	cryptfile.SetGlobal(nil)
	t.Cleanup(func() { cryptfile.Close(); cryptfile.SetGlobal(original) })
	path := filepath.Join(t.TempDir(), "object.bin")
	type record struct {
		Outer, Inner     int16
		Present, Failed  bool
		Team             byte
		Consumed         uint32
		Position         int64
		Created          bool
		XY               image.Point
		Flags, Extra, ID uint32
		TeamOut          byte
	}
	var records []record
	for _, outer := range []int16{-32768, -1, 0, 9, 10, 19, 20, 29, 30, 39, 40, 41, 60, 61, 62, 63, 64, 32767} {
		for _, inner := range []int16{1, 2, 3, 4, 5, 60, 61, 62, 63, 64} {
			if outer < 40 && inner != 1 {
				continue
			}
			modern := outer >= 40 && inner >= 61
			for _, present := range []bool{false, true} {
				if !modern && !present {
					continue
				}
				for _, team := range []byte{0, 7, 255} {
					for _, failed := range []bool{false, true} {
						name := fmt.Sprintf("outer%d-inner%d-present%v-team%d-fail%v", outer, inner, present, team, failed)
						t.Run(name, func(t *testing.T) {
							var s mapDrawableStream
							if modern {
								s.modernBase(inner, present, team, 0x12345678)
							} else {
								if outer >= 40 {
									s.u16(uint16(inner))
								}
								s.oldBase(outer, inner, team, 0x12345678)
							}
							if err := os.WriteFile(path, s.Bytes(), 0600); err != nil {
								t.Fatal(err)
							}
							if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
								t.Fatal(err)
							}
							defer cryptfile.Close()
							c.FailEvery = 0
							if failed {
								c.FailEvery = 1
							}
							count := uint32(17)
							dr := legacy.PortTestMapDrawableBase(4, outer, inner, false, &count)
							pos, err := cryptfile.Global().File.Seek(0, io.SeekCurrent)
							if err != nil {
								t.Fatal(err)
							}
							wantRead := s.Len()
							if modern && failed {
								wantRead = 18
							}
							if count != 17+uint32(wantRead) || pos != int64(wantRead) {
								t.Fatalf("consumed=%d position=%d want=%d+17/%d", count, pos, wantRead, wantRead)
							}
							r := record{Outer: outer, Inner: inner, Present: present, Failed: failed, Team: team, Consumed: count, Position: pos, Created: dr != nil}
							if (dr == nil) != failed {
								t.Fatalf("creation=%v failed=%v", dr != nil, failed)
							}
							if dr != nil {
								defer func() { dr.TeamVal.ID = 0; c.Nox_xxx_spriteDeleteStatic_45A4E0_drawable(dr) }()
								raw := unsafe.Slice((*byte)(unsafe.Pointer(dr)), int(unsafe.Sizeof(*dr)))
								r.XY = dr.PosVec
								r.Flags = binary.LittleEndian.Uint32(raw[120:])
								r.Extra = binary.LittleEndian.Uint32(raw[280:])
								r.ID = binary.LittleEndian.Uint32(raw[128:])
								r.TeamOut = byte(dr.TeamVal.ID)
								xy := image.Pt(64, 128)
								wantTeam := team
								if (!modern && outer < 20) || (modern && !present) {
									wantTeam = 0
								}
								extra := uint32(0)
								if modern && present || !modern && outer >= 40 && inner >= 3 {
									extra = 0x12345678
								}
								// The actual drawable factory marks new objects active (bit 4).
								flags := uint32(0x545)
								if modern && !present {
									flags = 0x01000004
								}
								if r.XY != xy || r.ID != 0x12345678 || r.TeamOut != wantTeam || r.Extra != extra || r.Flags != flags {
									t.Fatalf("state=%+v want xy=%v team=%d extra=%x flags=%x", r, xy, wantTeam, extra, flags)
								}
							}
							records = append(records, r)
						})
					}
				}
			}
		}
	}
	spellbookCapture(t, "map-drawable-base", records, "25470a33a9b23facda83bf28da0e678c95746b57404602b6e33d046182c78590")
}
