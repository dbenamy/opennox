//go:build porttest

package opennox

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
)

// Encode documented fields independently of the production reader. Script
// handler version 1 with an empty name is valid in both editor and game mode.
func objectXferCommonStream(outer, inner int16, present bool, team byte) *mapDrawableStream {
	s := new(mapDrawableStream)
	if outer >= 40 {
		s.u16(uint16(inner))
	} else {
		inner = 0
	}
	if inner > 64 {
		return s
	}
	modern := outer >= 40 && inner >= 61
	s.u32(0x12345678)
	if modern {
		s.u32(0x2468ace0)
	} else {
		s.u32(0x01408162)
	}
	if outer < 40 || inner < 4 {
		s.u32(64)
		s.u32(128)
	} else {
		s.f32(64.25)
		s.f32(128.75)
	}
	if modern {
		if !present {
			s.u8(0)
			return s
		}
		s.u8(1)
		s.u32(0x01408162)
	}
	if outer >= 10 {
		s.u8(3)
		s.WriteString("obj")
	}
	if outer >= 20 {
		s.u8(team)
	}
	if outer >= 30 {
		s.u8(7)
	}
	if outer >= 40 {
		if !modern {
			s.u32(0x2468ace0)
		}
		if inner >= 2 {
			if inner < 5 {
				s.u32(0x12340002)
			} else {
				s.u16(2)
			}
			s.u32(123)
			s.u32(456)
		}
		if inner >= 3 {
			s.u32(0x5e)
		}
	}
	if modern && inner >= 63 {
		s.u16(1)
		s.u32(0)
		s.u32(0x11223344)
	}
	if modern && inner >= 64 {
		s.u32(99)
	}
	return s
}

func TestObjectXferCommonReadContracts(t *testing.T) {
	var captures []objectXferCaptureRow
	defer func() { spellbookCapture(t, "object-xfer-common-read", captures, "275370459efd1056acaf282f58e260cb43dd006b8dbd1704866ec2985f00dfd7") }()
	core := newObjectXferOwner(t)
	path := filepath.Join(t.TempDir(), "record.bin")
	for _, outer := range []int16{-1, 0, 9, 10, 19, 20, 29, 30, 39, 40, 60, 64, 32767} {
		for _, inner := range []int16{-32768, -1, 0, 1, 2, 3, 4, 5, 60, 61, 62, 63, 64, 65, 32767} {
			if outer < 40 && inner != 0 {
				continue
			}
			for _, present := range []bool{false, true} {
				modern := outer >= 40 && inner >= 61 && inner <= 64
				if !modern && !present {
					continue
				}
				for _, team := range []byte{0, 7, 255} {
					t.Run(fmt.Sprintf("outer%d-inner%d-present%v-team%d", outer, inner, present, team), func(t *testing.T) {
						u := newObjectXferSimple(t, core)
						defer objectXferCaptureCase(t, &captures, u)
						u.ObjFlags = object.Flags(0x20000040)
						u.Field5 = 0x80
						u.Field34 = 999
						u.Field32 = 555
						u.TeamVal.ID = 9
						u.ScriptIDVal = 88
						u.PosVec = types.Pointf{X: 1, Y: 2}
						u.NewPos = types.Pointf{X: 3, Y: 4}
						s := objectXferCommonStream(outer, inner, present, team)
						if err := os.WriteFile(path, s.Bytes(), 0600); err != nil {
							t.Fatal(err)
						}
						if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
							t.Fatal(err)
						}
						defer cryptfile.Close()
						result := legacy.Nox_xxx_mapReadWriteObjData_4F4530(u, int(outer))
						pos, err := cryptfile.Global().File.Seek(0, io.SeekCurrent)
						if err != nil {
							t.Fatal(err)
						}
						if pos != int64(s.Len()) {
							t.Fatalf("stream position %d want %d", pos, s.Len())
						}
						if inner > 64 {
							if result != 0 || u.Field34 != 999 || u.PosVec.X != 1 || u.ScriptIDVal != 88 {
								t.Fatal("rejected version mutated object")
							}
							return
						}
						if result != 1 {
							t.Fatal("valid record rejected")
						}
						wantPos := types.Pointf{X: 64.25, Y: 128.75}
						if outer < 40 || inner < 4 {
							wantPos = types.Pointf{X: 64, Y: 128}
						}
						if u.PosVec != wantPos || u.NewPos != wantPos {
							t.Fatalf("position %v/%v want %v", u.PosVec, u.NewPos, wantPos)
						}
						wantScript := 88
						if outer >= 40 {
							wantScript = 0x2468ace0
						}
						if u.ScriptIDVal != wantScript {
							t.Fatalf("script ID %x want %x", u.ScriptIDVal, wantScript)
						}
						flags, extra, inventory, outTeam, name, lifetime := uint32(0x21408162), uint32(0x80), uint32(0), byte(9), "", uint32(555)
						if modern && !present {
							flags = 0x20000040
						} else {
							if outer >= 10 {
								name = "obj"
							}
							if outer >= 20 {
								outTeam = team
							}
							if outer >= 30 {
								inventory = 7
							}
							if outer >= 40 && inner >= 3 {
								extra |= 0x5e
							}
							if modern && inner >= 64 {
								lifetime = 99
							}
						}
						if uint32(u.ObjFlags) != flags || u.Field5 != extra || u.Field34 != inventory || byte(u.TeamVal.ID) != outTeam || u.ID() != name || u.Field32 != lifetime {
							t.Fatalf("fields flags=%x extra=%x inventory=%d team=%d name=%q lifetime=%d; want %x/%x/%d/%d/%q/%d", u.ObjFlags, u.Field5, u.Field34, u.TeamVal.ID, u.ID(), u.Field32, flags, extra, inventory, outTeam, name, lifetime)
						}
					})
				}
			}
		}
	}
}
