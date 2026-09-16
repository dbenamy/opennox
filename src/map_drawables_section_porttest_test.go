//go:build porttest

package opennox

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/opennox/libs/object"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

func TestMapDrawableSections(t *testing.T) {
	mapDrawableTables(t)
	names := []string{"ColorLight", "ColorLightMovable", "TeamBase", "PressurePlate", "MapBox", "MapDoor", "MapMarker", "MapIgnored", "MapUnknown"}
	o := newObjectDrawingOwner(t, names...)
	c := o.c
	t.Cleanup(noxflags.PortTestGameFlags(1))
	t.Cleanup(c.srv.PortTestMinimapWalls())
	words, restore := legacy.PortTestMinimapWords()
	t.Cleanup(restore)
	*words["polygons"] = 1
	for _, off := range []uintptr{1309788, 1309792, 1309796, 1309800, 1309804} {
		p := memmap.PtrUint32(0x5D4594, off)
		old := *p
		*p = 0
		t.Cleanup(func() { *p = old })
	}
	original := cryptfile.Global()
	cryptfile.SetGlobal(nil)
	t.Cleanup(func() { cryptfile.Close(); cryptfile.SetGlobal(original) })
	oldCodes, oldLen := objectTypeCode16ByInd, objectTypeCode16ByInd_len
	objectTypeCode16ByInd = make([]uint16, 128)
	objectTypeCode16ByInd_len = 9
	t.Cleanup(func() { objectTypeCode16ByInd = oldCodes; objectTypeCode16ByInd_len = oldLen })
	ids := make([]int, len(names))
	bodies := make([][]byte, len(names))
	for i, name := range names {
		ids[i] = c.Things.IndByID(name)
		typ := c.Things.TypeByInd(ids[i])
		objectTypeCode16ByInd[ids[i]] = uint16(101 + i)
		typ.ObjClass = object.Class(0x20000000)
		var s mapDrawableStream
		s.u16(64)
		s.modernBase(64, false, 0, 0)
		switch i {
		case 0, 1:
			s.light(64, i)
		case 2:
			for j := 0; j < 4; j++ {
				s.u8(0)
			}
		case 3:
			s.u32(70)
			s.u32(80)
			s.Write([]byte{1, 2, 3, 4, 5, 6})
		case 4:
			typ.ObjClass = object.Class(0x400200)
			s.u32(70)
			s.u32(80)
		case 5:
			typ.ObjClass = object.Class(0x400080)
			s.u32(8)
			s.u32(1)
			s.u32(16)
		case 6:
			typ.ObjClass = object.Class(0x400000)
			typ.ObjSubClass = object.SubClass(8)
			s.u32(42)
			s.u8(1)
		case 7:
			typ.ObjClass = 0
			s.fill(7)
		case 8:
			s.fill(9)
		}
		bodies[i] = append([]byte(nil), s.Bytes()...)
	}
	type record struct {
		Name         string
		Warm, Failed bool
		Return       int
		Position     int64
		Types        []int
		LoadError    bool
		Cache        [4]uint32
	}
	var records []record
	path := filepath.Join(t.TempDir(), "section.bin")
	for _, spec := range []struct {
		name                string
		version             uint16
		which               []int
		pad, short, unknown bool
	}{
		{name: "empty", version: 1}, {name: "zero-version", version: 0}, {name: "signed-version", version: 32768}, {name: "future-version", version: 2},
		{name: "all", version: 1, which: []int{0, 1, 2, 3, 4, 5, 6, 7, 8}},
		{name: "padding", version: 1, which: []int{3, 0, 7, 6}, pad: true},
		{name: "short-length", version: 1, which: []int{3, 4, 6}, short: true},
		{name: "ignored", version: 1, which: []int{7, 8}},
		{name: "unknown-toc", version: 1, unknown: true},
	} {
		for _, warm := range []bool{false, true} {
			for _, failed := range []bool{false, true} {
				t.Run(fmt.Sprintf("%s-warm%v-fail%v", spec.name, warm, failed), func(t *testing.T) {
					for i, off := range []uintptr{1309792, 1309796, 1309800, 1309804} {
						v := uint32(0)
						if warm {
							v = uint32(ids[i])
						}
						*memmap.PtrUint32(0x5D4594, off) = v
					}
					var s mapDrawableStream
					s.u16(spec.version)
					for _, i := range spec.which {
						s.u16(uint16(101 + i))
						size := len(bodies[i])
						if spec.pad {
							size += 5
						}
						if spec.short {
							size = 0
						}
						s.u32(uint32(size))
						s.Write(bodies[i])
						if spec.pad {
							s.fill(5)
						}
					}
					if spec.unknown {
						s.u16(65535)
						s.u32(11)
						s.fill(11)
					}
					s.u16(0)
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
					ret := legacy.PortTestMapDrawableRecord(6, 0)
					pos, err := cryptfile.Global().File.Seek(0, io.SeekCurrent)
					if err != nil {
						t.Fatal(err)
					}
					r := record{Name: spec.name, Warm: warm, Failed: failed, Return: ret, Position: pos, LoadError: c.Objs.LoadError}
					for dr := c.Objs.List1; dr != nil; dr = dr.NextPtr {
						r.Types = append(r.Types, int(dr.TypeIDVal))
					}
					defer func() {
						c.srv.Walls.EachWallRaw(func(w *server.Wall) bool { w.Field32 = 0; return true })
						for c.Objs.List1 != nil {
							dr := c.Objs.List1
							dr.TeamVal.ID = 0
							c.Nox_xxx_spriteDeleteStatic_45A4E0_drawable(dr)
						}
					}()
					for i, off := range []uintptr{1309792, 1309796, 1309800, 1309804} {
						r.Cache[i] = *memmap.PtrUint32(0x5D4594, off)
					}
					wantRet, wantPos := 1, int64(s.Len())
					if spec.version == 2 {
						wantRet, wantPos = 0, 2
					}
					if spec.unknown {
						wantRet, wantPos = 0, 8
					}
					// With failed allocation the legacy record length includes the already-read
					// base bytes, so framing may stop at the next misread code. Capture that
					// existing behavior separately rather than assuming successful framing.
					if !failed || len(spec.which) == 0 || spec.name == "ignored" {
						if ret != wantRet || pos != wantPos {
							t.Fatalf("section %+v want return=%d position=%d", r, wantRet, wantPos)
						}
					}
					if !failed && spec.version != 2 {
						var want []int
						for j := len(spec.which) - 1; j >= 0; j-- {
							i := spec.which[j]
							if i < 7 {
								want = append(want, ids[i])
							}
						}
						if !slices.Equal(r.Types, want) {
							t.Fatalf("loaded types=%v want=%v", r.Types, want)
						}
					}
					records = append(records, r)
				})
			}
		}
	}
	t.Run("write-mode", func(t *testing.T) {
		if err := cryptfile.OpenGlobal(path, cryptfile.WriteOnly, -1); err != nil {
			t.Fatal(err)
		}
		ret := legacy.PortTestMapDrawableRecord(6, 0)
		if err := cryptfile.Close(); err != nil {
			t.Fatal(err)
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if ret != 0 || !bytes.Equal(raw, []byte{1, 0}) {
			t.Fatalf("write mode return=%d bytes=%x", ret, raw)
		}
	})
	spellbookCapture(t, "map-drawable-sections", records, "750143749d1d04499fae4c1c5ae5174eebdd8f9cd450ce1095421414834527a5")
}
