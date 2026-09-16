//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/server"
	"image"
	"io"
	"math"
	"os"
	"path/filepath"
	"testing"
	"unsafe"

	"github.com/opennox/libs/object"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func (s *mapDrawableStream) light(outer int16, variant int) {
	s.u32(4)
	s.f32(float32(20 + 60*variant))
	s.u32(11)
	s.u32(0)
	s.u32(1 << 16)
	s.u32(2 << 16)
	s.u32(3 << 16)
	s.u16(123)
	s.u16(7)
	s.u32(0x10203040)
	if outer < 2 {
		return
	}
	s.u16(2)
	for i := 0; i < 16; i++ {
		if i%4 == 0 {
			s.fillZero(3)
		} else {
			s.u8(byte(i))
			s.u8(byte(i + 1))
			s.u8(byte(i + 2))
		}
	}
	for i := 0; i < 16; i++ {
		if i%3 == 0 {
			s.u8(0)
		} else {
			s.u8(byte(i + 1))
		}
	}
	for i := 0; i < 16; i++ {
		if i%5 == 0 {
			s.u8(0)
		} else {
			s.u8(byte(16 - i))
		}
	}
	s.u16(1)
	s.u16(2)
	s.u16(3)
	s.u32(4)
	s.u16(5)
	s.u16(6)
	s.u8(127)
	if outer == 41 {
		s.u8(7)
	} else if outer >= 42 {
		s.u32(0x11223344)
	}
}
func (s *mapDrawableStream) fillZero(n int) { s.Write(make([]byte, n)) }

func TestMapDrawableTypedRecords(t *testing.T) {
	mapDrawableTables(t)
	o := newObjectDrawingOwner(t)
	c := o.c
	t.Cleanup(noxflags.PortTestGameFlags(1))
	t.Cleanup(c.srv.PortTestMinimapWalls())
	words, restore := legacy.PortTestMinimapWords()
	t.Cleanup(restore)
	*words["polygons"] = 1
	names := make([]*byte, len(o.mods))
	for i := range names {
		p, free := alloc.CString(fmt.Sprintf("MapFixtureMod%d", i))
		names[i] = p
		t.Cleanup(free)
	}
	t.Cleanup(c.srv.PortTestControlsModifiers(o.mods, names))
	original := cryptfile.Global()
	cryptfile.SetGlobal(nil)
	t.Cleanup(func() { cryptfile.Close(); cryptfile.SetGlobal(original) })
	path := filepath.Join(t.TempDir(), "typed-object.bin")
	typ := c.Things.TypeByInd(4)
	type record struct {
		Kind         string
		Outer        int16
		Variant      int
		Failed       bool
		Return       int
		Position     int64
		LoadError    bool
		Flags        uint32
		Data         []byte
		Counts       [3]byte
		Mods         [4]int
		Sentinels    [2]uint16
		OnMinimap    bool
		Wall         bool
		WallPosition image.Point
	}
	var records []record
	for _, kind := range []string{"light", "box", "door", "team", "plate", "marker"} {
		for _, outer := range []int16{0, 1, 2, 39, 40, 41, 42, 60, 61, 64} {
			for variant := 0; variant < 2; variant++ {
				for _, failed := range []bool{false, true} {
					t.Run(fmt.Sprintf("%s-v%d-case%d-fail%v", kind, outer, variant, failed), func(t *testing.T) {
						typ.ObjClass = 0
						typ.ObjSubClass = 0
						var s mapDrawableStream
						s.u16(uint16(outer))
						if outer < 40 {
							s.oldBase(outer, 1, 0, 0)
						} else {
							s.modernBase(64, false, 0, 0)
						}
						baseBytes := s.Len()
						op := 0
						switch kind {
						case "light":
							op = 2
							s.light(outer, variant)
						case "box":
							op = 1
							typ.ObjClass = object.Class(0x200)
							s.u32(uint32(59 + 2*variant))
							s.u32(uint32(100 - 101*variant))
						case "door":
							op = 1
							typ.ObjClass = object.Class(0x80)
							s.u32(uint32(8 * variant))
							s.u32(uint32(variant))
							if outer >= 41 {
								s.u32(uint32(16 + 8*variant))
							}
						case "team":
							op = 3
							for _, name := range []string{"", "MapFixtureMod0", "missing", "MapFixtureMod3"} {
								s.u8(byte(len(name)))
								s.WriteString(name)
							}
						case "plate":
							op = 4
							s.u32(uint32(59 + 2*variant))
							s.u32(uint32(100 - 101*variant))
							if outer >= 41 {
								s.Write([]byte{1, 2, 3, 4, 5, 6})
							}
						case "marker":
							op = 5
							if outer >= 61 {
								s.u32(0x55667788)
								s.u8(byte(variant))
							}
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
						c.Objs.LoadError = false
						ret := legacy.PortTestMapDrawableRecord(op, 4)
						pos, err := cryptfile.Global().File.Seek(0, io.SeekCurrent)
						if err != nil {
							t.Fatal(err)
						}
						r := record{Kind: kind, Outer: outer, Variant: variant, Failed: failed, Return: ret, Position: pos, LoadError: c.Objs.LoadError}
						if failed {
							want := baseBytes
							if outer >= 40 {
								want = 20
							}
							if ret != 0 || !r.LoadError || pos != int64(want) {
								t.Fatalf("allocation failure %+v want read %d", r, want)
							}
						} else {
							if ret != s.Len() || pos != int64(s.Len()) || r.LoadError {
								t.Fatalf("read %+v want %d", r, s.Len())
							}
							dr := c.Objs.List1
							if dr == nil {
								t.Fatal("missing real drawable")
							}
							defer func() {
								c.srv.Walls.EachWallRaw(func(w *server.Wall) bool { w.Field32 = 0; return true })
								dr.TeamVal.ID = 0
								c.Nox_xxx_spriteDeleteStatic_45A4E0_drawable(dr)
							}()
							raw := unsafe.Slice((*byte)(unsafe.Pointer(dr)), int(unsafe.Sizeof(*dr)))
							r.Flags = binary.LittleEndian.Uint32(raw[120:])
							if r.Flags&4 == 0 {
								t.Fatal("loaded drawable is not active")
							}
							if binary.LittleEndian.Uint32(raw[288:]) != 0 {
								t.Fatal("loaded frame state was not reset")
							}
							switch kind {
							case "light":
								r.Data = append([]byte(nil), raw[136:276]...)
								copy(r.Counts[:], raw[432:435])
								want := [3]byte{}
								if outer >= 2 {
									want = [3]byte{12, 10, 12}
								}
								if r.Counts != want {
									t.Fatalf("light counts=%v want=%v", r.Counts, want)
								}
								intensity := float32(20 + 60*variant)
								if outer < 2 && variant == 1 {
									intensity = 63
								}
								if math.Float32frombits(binary.LittleEndian.Uint32(raw[140:])) != intensity {
									t.Fatal("light intensity/default clamp")
								}
							case "box", "plate":
								r.Data = append([]byte(nil), raw[44:100]...)
								w, h := float32(59+2*variant), float32(100-101*variant)
								if kind == "box" {
									if w > 60 {
										w = 60
									}
									if h > 60 {
										h = 60
									}
								}
								if math.Float32frombits(binary.LittleEndian.Uint32(raw[56:])) != w || math.Float32frombits(binary.LittleEndian.Uint32(raw[60:])) != h {
									t.Fatal("shape width/height clamp")
								}
								if kind == "plate" {
									r.Data = append(r.Data, raw[432:438]...)
									want := []byte{90, 90, 90, 10, 10, 10}
									if outer >= 41 {
										want = []byte{1, 2, 3, 4, 5, 6}
									}
									if !bytes.Equal(raw[432:438], want) {
										t.Fatal("plate defaults/override")
									}
								}
							case "door":
								r.Data = append([]byte(nil), raw[432:436]...)
								r.Data = append(r.Data, raw[299])
								if raw[432] != byte(variant) || raw[433] != byte(variant) || raw[299] != byte(8*variant) {
									t.Fatal("door lock/direction")
								}
								c.srv.Walls.EachWallRaw(func(w *server.Wall) bool {
									if w.Field32 == uint32(uintptr(dr.C())) {
										r.Wall = true
										r.WallPosition = image.Pt(int(w.X5), int(w.Y6))
										if w.Flags4&0x10 == 0 || byte(w.Field8) != 1 {
											t.Fatal("door wall flags/region")
										}
									}
									return true
								})
								wantPos := image.Pt(2+variant, 5)
								if outer >= 41 {
									wantPos = image.Pt(3-variant, 6)
								}
								if r.WallPosition != wantPos {
									t.Fatalf("door wall position=%v want=%v", r.WallPosition, wantPos)
								}
								if !r.Wall {
									t.Fatal("door was not registered with actual wall owner")
								}
							case "team":
								for i := 0; i < 4; i++ {
									p := binary.LittleEndian.Uint32(raw[432+4*i:])
									for j, m := range o.mods {
										if p == uint32(uintptr(unsafe.Pointer(m))) {
											r.Mods[i] = j + 1
										}
									}
									if p != 0 && r.Mods[i] == 0 {
										t.Fatal("modifier pointer outside owner")
									}
								}
								if r.Mods != [4]int{0, 1, 0, 4} {
									t.Fatalf("modifier lookup %v", r.Mods)
								}
								r.Sentinels = [2]uint16{binary.LittleEndian.Uint16(raw[448:]), binary.LittleEndian.Uint16(raw[450:])}
								if r.Sentinels != [2]uint16{65535, 65535} {
									t.Fatal("team sentinels")
								}
							case "marker":
								r.OnMinimap = c.Objs.FirstMinimapList() == dr
								if r.OnMinimap != (outer >= 61 && variant == 1) {
									t.Fatal("minimap membership")
								}
							}
						}
						records = append(records, r)
					})
				}
			}
		}
	}
	spellbookCapture(t, "map-drawable-typed", records, "48492ed7e77735caa052fe9a5d1dddcfa6812cb75ae18da4a3c658b528cb46b5")
}
