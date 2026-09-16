//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
	"unsafe"

	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
)

func objectXferEmptyBase(s *mapDrawableStream, version int16) {
	if version >= 40 {
		s.u16(64)
	}
	s.u32(0x12345678)
	if version >= 40 {
		s.u32(88)
	} else {
		s.u32(0)
	}
	if version >= 40 {
		s.f32(64.25)
		s.f32(128.75)
		s.u8(0)
		return
	}
	s.u32(64)
	s.u32(128)
	if version >= 10 {
		s.u8(0)
	}
	if version >= 20 {
		s.u8(0)
	}
	if version >= 30 {
		s.u8(0)
	}
}
func objectXferScript(s *mapDrawableStream, flags uint32) { s.u16(1); s.u32(0); s.u32(flags) }
func objectXferHistoricalPayload(s *mapDrawableStream, name string, version int16) {
	switch name {
	case "SpellPagePedestal":
		s.u32(0x12345678)
	case "Readable", "Exit":
		if version >= 2 {
			s.u32(1)
		}
		s.u8(0)
		if name == "Exit" && version >= 31 {
			s.u32(17)
			s.u32(29)
		}
	case "Door":
		s.u32(0)
		s.u32(7)
		if version >= 41 {
			s.u32(0)
		}
	case "Trigger":
		s.u32(20)
		s.u32(30)
		if version < 41 {
			s.Write(make([]byte, 9))
		} else {
			s.Write([]byte{1, 2, 3, 4, 5, 6})
		}
		s.u32(123)
		if version < 3 {
			for i := 0; i < 2; i++ {
				s.u32(0)
				s.u32(uint32(i + 1))
				s.u32(0)
			}
		} else {
			objectXferScript(s, 1)
			objectXferScript(s, 2)
			if version >= 31 {
				objectXferScript(s, 3)
			}
		}
		if version < 31 {
			for i := 0; i < 4; i++ {
				s.u8(1)
				s.u32(0x2468ace0)
			}
		}
		s.u32(17)
		s.u32(29)
		if version >= 21 {
			s.u8(7)
			s.u8(9)
		}
		if version >= 61 {
			s.u8(11)
			s.u8(13)
			s.u32(37)
		}
	case "Hole":
		if version >= 42 {
			s.u32(41)
		}
		if version >= 41 {
			objectXferScript(s, 7)
		}
		s.u32(17)
		s.u32(29)
		if version >= 41 {
			s.u32(31)
			s.u16(33)
		}
	case "Transporter", "ElevatorShaft":
		s.u32(123)
	case "Elevator":
		s.u32(123)
		if version >= 41 {
			s.u32(321)
		}
		if version >= 61 {
			s.u8(7)
		}
	case "Mover":
		s.u32(17)
		s.u32(29)
		s.u32(31)
		if version >= 41 {
			s.u8(7)
			s.u32(41)
			s.u32(43)
		}
		if version >= 42 {
			s.f32(3.5)
			s.f32(2.25)
		}
	case "Glyph":
		if version < 41 {
			s.u32(0)
		}
		s.u8(32)
		s.u32(17)
		s.u32(29)
		s.u8(0)
		if version < 31 {
			s.Write(make([]byte, 20))
		}
	case "Sentry":
		s.u32(17)
		s.u32(29)
		if version >= 61 {
			s.u32(31)
		}
	default:
		panic(name)
	}
}
func TestObjectXferTypedHistoricalRecords(t *testing.T) {
	handles.Init()
	t.Cleanup(handles.Release)
	core := newObjectXferOwner(t)
	mapDrawableTables(t)
	path := filepath.Join(t.TempDir(), "historical.bin")
	for _, sp := range objectXferKinds {
		for _, version := range []int16{-32768, -1, 0, 1, 2, 3, 9, 10, 19, 20, 21, 29, 30, 31, 39, 40, 41, 42, 59, 60, 61} {
			if version > int16(sp.version) {
				continue
			}
			t.Run(fmt.Sprintf("%s-%d", sp.name, version), func(t *testing.T) {
				u := newObjectXferTyped(t, core, sp.name)
				u.ObjFlags = 0
				u.ScriptIDVal = 88
				u.Field34 = 777
				var stream mapDrawableStream
				stream.u16(uint16(version))
				objectXferEmptyBase(&stream, version)
				objectXferHistoricalPayload(&stream, sp.name, version)
				if err := os.WriteFile(path, stream.Bytes(), 0600); err != nil {
					t.Fatal(err)
				}
				if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
					t.Fatal(err)
				}
				defer func() {
					cf := cryptfile.Global()
					if cf != nil && cf.File.File.Handle != nil {
						legacy.Nox_fs_close((*legacy.FILE)(cf.File.File.Handle))
					}
					cryptfile.Close()
				}()
				if err := u.CallXfer(nil); err != nil {
					t.Fatal(err)
				}
				pos, err := cryptfile.Global().File.Seek(0, io.SeekCurrent)
				if err != nil || pos != int64(stream.Len()) {
					t.Fatalf("stream pos=%d want=%d err=%v", pos, stream.Len(), err)
				}
				wantPos := types.Pointf{X: 64, Y: 128}
				if version >= 40 {
					wantPos = types.Pointf{X: 64.25, Y: 128.75}
				}
				if u.PosVec != wantPos || u.NewPos != wantPos || u.Extent != 0x12345678 || u.ScriptIDVal != 88 || u.Field34 != 777 {
					t.Fatal("common fields or saved lifetime")
				}
				word := func(p unsafe.Pointer, off int) uint32 {
					return binary.LittleEndian.Uint32(unsafe.Slice((*byte)(unsafe.Add(p, off)), 4))
				}
				switch sp.name {
				case "SpellPagePedestal":
					if word(u.CollideData, 0) != 0x12345678 {
						t.Fatal("pedestal value")
					}
				case "Readable":
					if word(u.UseData.Ptr, 256) != 0 {
						t.Fatal("readable terminator")
					}
				case "Exit":
					if version >= 31 && (word(u.CollideData, 80) != 17 || word(u.CollideData, 84) != 29) {
						t.Fatal("exit destination")
					}
				case "Door":
					if *(*byte)(unsafe.Add(u.UpdateData, 1)) != 7 || word(u.UpdateData, 12) != 0 {
						t.Fatal("door state")
					}
				case "Trigger":
					if u.Shape.Box.W != 20 || u.Shape.Box.H != 30 || word(u.UpdateData, 0) != 123 || word(u.UpdateData, 44) != 17 || word(u.UpdateData, 48) != 29 {
						t.Fatal("trigger geometry/state")
					}
					if version >= 61 && u.Field33 != 37 {
						t.Fatal("trigger animation frame")
					}
				case "Hole":
					if word(u.CollideData, 8) != 17 || word(u.CollideData, 12) != 29 {
						t.Fatal("hole target")
					}
				case "Transporter":
					if word(u.UpdateData, 16) != 123 {
						t.Fatal("transporter target")
					}
				case "Elevator", "ElevatorShaft":
					if word(u.UpdateData, 8) != 123 {
						t.Fatal("elevator target")
					}
				case "Mover":
					if word(u.UpdateData, 4) != 17 || word(u.UpdateData, 8) != 29 || word(u.UpdateData, 32) != 31 {
						t.Fatal("mover fields")
					}
				case "Glyph":
					if u.Direction1 != 32 || u.Direction2 != 32 || word(u.InitData, 28) != 17 || word(u.InitData, 32) != 29 || word(u.InitData, 24) != 0 {
						t.Fatal("glyph state")
					}
				case "Sentry":
					want := uint32(17)
					if version >= 61 {
						want = 31
					}
					if word(u.UpdateData, 0) != want || word(u.UpdateData, 4) != 17 || word(u.UpdateData, 8) != 29 {
						t.Fatal("sentry state")
					}
				}
			})
		}
	}
}
