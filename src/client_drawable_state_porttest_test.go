//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
	"unsafe"
)

func TestClientDrawableStateAccessors(t *testing.T) {
	c, _, _ := newEffectsFullOwner(t)
	dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(300, 400))
	next := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(310, 410))
	type row struct {
		Class, Sub, Anim, Value uint32
		Predicate               bool
		Slave, Previous, Flags  uint32
	}
	var rows []row
	for _, cl := range []uint32{0, 2, 4, 0x400000, 0x400002, 0x80000000, 0xffffffff} {
		for _, sub := range []uint32{0, 8, 0x40000, 0x40008, 0xffffffff} {
			for _, anim := range []uint32{0, 7, 8, 9, 0xffffffff} {
				for _, value := range []uint32{0, 1, 255, 0x80000000, 0xffffffff} {
					dr.ObjClass = object.Class(cl)
					dr.ObjSubClass = object.SubClass(sub)
					dr.AnimInd = anim
					dr.AnimFrameSlave = 0xabcdef01
					dr.Field_78 = 0x12345678
					dr.ObjFlags = object.Flags(0xa5010200)
					raw := unsafe.Slice((*byte)(dr.C()), 512)
					before := append([]byte(nil), raw...)
					got := legacy.PortTestDrawableState(0, dr, 0)
					want := cl&0x400000 != 0 && sub&8 != 0
					if (got == 1) != want || got > 1 || !bytes.Equal(before, raw) {
						t.Fatal("drawable predicate")
					}
					ret := legacy.PortTestDrawableState(2, dr, int(value))
					if !(cl&2 != 0 && sub&0x40000 != 0 && anim == 8) {
						binary.LittleEndian.PutUint32(before[312:], 0xabcdef01)
						binary.LittleEndian.PutUint32(before[308:], value)
					}
					if ret != uintptr(dr.C()) || !bytes.Equal(before, raw) {
						t.Fatal("frame helper writes/return")
					}
					ret = legacy.PortTestDrawableState(1, dr, 0)
					binary.LittleEndian.PutUint32(before[120:], 0xa5010204)
					if ret != uintptr(dr.C()) || !bytes.Equal(before, raw) {
						t.Fatal("active helper writes/return")
					}
					rows = append(rows, row{cl, sub, anim, value, want, dr.AnimFrameSlave, dr.Field_78, uint32(dr.ObjFlags)})
				}
			}
		}
	}
	// The accessors return the specific production links without altering state.
	dr.NextPtr = next
	dr.Field_104 = next
	if legacy.PortTestDrawableState(3, dr, 0) != uintptr(next.C()) || legacy.PortTestDrawableState(4, dr, 0) != uintptr(next.C()) || legacy.PortTestDrawableState(3, nil, 0) != 0 {
		t.Fatal("list accessor")
	}
	dr.NextPtr = nil
	dr.Field_104 = nil
	// Restore class/flags before owner cleanup; this fixture changed only fields,
	// not membership in the corresponding auxiliary lists.
	dr.ObjClass = 0
	dr.ObjSubClass = 0
	dr.ObjFlags = 0
	drawableStateCapture(t, "state", rows, "6fc328fa72079b952f904b5afa1901bbfdb44027757034c2b38d6d0122fc830e")
}

func TestClientDrawableStreamBoundaries(t *testing.T) {
	raw := memmap.Slice(0x5D4594, 1198020)[:255*8+8]
	saved := append([]byte(nil), raw...)
	t.Cleanup(func() { copy(raw, saved) })
	clear(raw)
	binary.LittleEndian.PutUint16(raw[8:], 17)
	binary.LittleEndian.PutUint16(raw[10:], 4)
	c, pix, env := newEffectsFullOwner(t)
	type row struct {
		Start    [2]int32
		Delta    [2]int8
		Absolute bool
		Return   int32
		Pos      [2]int32
		Count    int
	}
	var rows []row
	for _, abs := range []bool{false, true} {
		for _, start := range [][2]int32{{0, 0}, {1, 1}, {5999, 5999}, {6000, 6000}, {-1, 300}, {6001, 300}, {300, -1}, {300, 6001}, {32767, -32768}, {0x7fffffff, -0x80000000}} {
			for _, delta := range [][2]int8{{0, 0}, {1, 0}, {-1, 0}, {0, 1}, {0, -1}, {-128, 127}} {
				c.resetCase(env, pix, 1, 100)
				// Alias zero is a terminator prefix; alias 1 comes from an owned table.
				packet := []byte{1, byte(delta[0]), byte(delta[1]), 0}
				want := [2]int32{int32(uint32(start[0]) + uint32(int32(delta[0]))), int32(uint32(start[1]) + uint32(int32(delta[1])))}
				if abs {
					packet = []byte{0, 1}
					packet = binary.LittleEndian.AppendUint16(packet, uint16(start[0]))
					packet = binary.LittleEndian.AppendUint16(packet, uint16(start[1]))
					packet = append(packet, 0)
					want = [2]int32{int32(int16(start[0])), int32(int16(start[1]))}
				}
				// All positions rejected by C return minus bytes consumed through coordinates.
				valid := want[0] >= 0 && want[0] <= 6000 && want[1] >= 0 && want[1] <= 6000
				if valid {
					continue
				} // Successful real-owner records are covered by the stream matrix.
				ret, pos := legacy.PortTestDrawableStream(2, packet, start)
				consumed := 3
				if abs {
					consumed = 6
				}
				if ret != int32(-consumed) || pos != want || c.Objs.Count != 0 {
					t.Fatalf("boundary abs%v start%v delta%v ret%d pos%v", abs, start, delta, ret, pos)
				}
				rows = append(rows, row{start, delta, abs, ret, pos, c.Objs.Count})
			}
		}
	}
	for _, packet := range [][]byte{{0, 0, 0}, {0, 0, 0, 255, 255}} {
		ret, pos := legacy.PortTestDrawableStream(2, packet, [2]int32{123, 456})
		if ret != -3 || pos != [2]int32{123, 456} {
			t.Fatal("terminator changed coordinates")
		}
	}
	drawableStateCapture(t, "bounds", rows, "fbb13bc9de414585770058fd962f71959b57461c657216e237e9afbacf2fcddc")
}
