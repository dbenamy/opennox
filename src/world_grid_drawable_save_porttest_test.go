//go:build porttest

package opennox

import (
	"encoding/binary"
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"image"
	"math"
	"testing"
	"unsafe"
)

func TestWorldGridSaveDrawables(t *testing.T) {
	c, _, _ := newEffectsFullOwner(t)
	core := newObjectXferOwner(t)
	typ := c.Things.TypeByInd(1)
	oldName := typ.Name
	name, freeName := alloc.CString("PaintObject")
	typ.Name = name
	t.Cleanup(func() { typ.Name = oldName; freeName() })
	scratch := memmap.PtrUint32(0x5D4594, 1563904)
	oldScratch := *scratch
	*scratch = 1
	t.Cleanup(func() { *scratch = oldScratch })
	poolCodes := make(map[unsafe.Pointer]uint32)
	var pool []*server.Object
	for i := 0; i < 64; i++ {
		u := core.NewObjectByTypeInd(1)
		if u == nil {
			t.Fatal("pool survey")
		}
		pool = append(pool, u)
		poolCodes[u.CObj()] = u.NetCode
	}
	for _, u := range pool {
		core.Objs.FreeObject(u)
	}
	oldSave := legacy.Nox_xxx_xfer_saveObj51DF90
	t.Cleanup(func() { legacy.Nox_xxx_xfer_saveObj51DF90 = oldSave })
	type state struct {
		Pos   [2]uint32
		Words [8]uint32
	}
	type result struct {
		Class, Flags uint32
		Icon         int
		Position     int
		SaveReturn   int
		Saved        []state
	}
	var rows []result
	for _, class := range []object.Class{object.ClassSimple, object.ClassImmobile, object.ClassImmobile | object.ClassMissile, object.ClassImmobile | object.ClassVisibleEnable} {
		for _, flags := range []object.Flags{0, object.FlagNoCollide, object.FlagNoCollide | object.FlagShadow} {
			for _, icon := range []int{-1, 0, 5} {
				for pi, pos := range []image.Point{{0, 0}, {123, -99}, {16777217, -16777217}} {
					for _, saveReturn := range []int{0, 1} {
						core.PortTestObjectXferAdmission(class, flags, true)
						core.Types.ByInd(1).Icon = icon
						dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(1, pos)
						if dr == nil {
							t.Fatal("drawable factory")
						}
						raw := unsafe.Slice((*byte)(dr.C()), 512)
						binary.LittleEndian.PutUint32(raw[120:], 0x10203040)
						binary.LittleEndian.PutUint32(raw[128:], 0xfedc1234)
						binary.LittleEndian.PutUint32(raw[280:], 0x87654321)
						before := append([]byte(nil), raw...)
						r := result{Class: uint32(class), Flags: uint32(flags), Icon: icon, Position: pi, SaveReturn: saveReturn}
						legacy.Nox_xxx_xfer_saveObj51DF90 = func(_ *cryptfile.CryptFile, u *server.Object) int {
							w := unsafe.Slice((*uint32)(unsafe.Add(u.CObj(), 16)), 8)
							s := state{Pos: [2]uint32{math.Float32bits(u.PosVec.X), math.Float32bits(u.PosVec.Y)}}
							copy(s.Words[:], w)
							r.Saved = append(r.Saved, s)
							if u.NetCode != 0xfedc1234 || w[0] != 0x10203040 || w[1] != 0x87654321 || w[6] != 0xfedc1234 || w[7] != 0xfedc1234 {
								t.Fatalf("transferred fields %+v", s)
							}
							if u.PosVec.X != float32(float64(pos.X)+0.5) || u.PosVec.Y != float32(float64(pos.Y)+0.5) {
								t.Fatalf("position %+v", u.PosVec)
							}
							return saveReturn
						}
						if legacy.Sub_51DED0() != 1 {
							t.Fatal("save traversal return")
						}
						eligible := class == object.ClassImmobile && flags == object.FlagNoCollide && icon != -1
						if (len(r.Saved) == 1) != eligible || core.Objs.Alive != 0 {
							t.Fatalf("eligibility class%x flags%x icon%d count%d alive%d", class, flags, icon, len(r.Saved), core.Objs.Alive)
						}
						if eligible {
							pool = nil
							for i := 0; i < 64; i++ {
								u := core.NewObjectByTypeInd(1)
								if u == nil || u.NetCode != poolCodes[u.CObj()] {
									t.Fatal("temporary netcode not restored")
								}
								pool = append(pool, u)
							}
							for _, u := range pool {
								core.Objs.FreeObject(u)
							}
						}
						for i, v := range raw {
							if v != before[i] {
								t.Fatalf("drawable changed byte%d", i)
							}
						}
						c.Nox_xxx_spriteDeleteStatic_45A4E0_drawable(dr)
						rows = append(rows, r)
					}
				}
			}
		}
	}
	if legacy.Sub_51DED0() != 1 {
		t.Fatal("empty traversal")
	}
	drawableStateCapture(t, "world-save-drawables", rows, "4a642d41cd2980a6fa6a950b46817194425606f381cf4ee7dd36f718db5deb94")
}
