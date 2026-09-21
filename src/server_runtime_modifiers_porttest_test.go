//go:build porttest

package opennox

import (
	"encoding/binary"
	"math"
	"testing"
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func TestServerRuntimeMaterialCache(t *testing.T) {
	o := newObjectDrawingOwner(t)
	mods, free := alloc.Make([]server.ModifierEff{}, 2)
	t.Cleanup(free)
	names, free := alloc.Make([]byte{}, 32)
	t.Cleanup(free)
	copy(names[16:], "Other")
	t.Cleanup(o.c.srv.PortTestControlsModifiers([]*server.ModifierEff{&mods[0], &mods[1]}, []*byte{&names[0], &names[16]}))
	unit, free := alloc.New(server.Object{})
	t.Cleanup(free)
	slots, free := alloc.Make([]unsafe.Pointer{}, 4)
	t.Cleanup(free)
	unit.InitData = unsafe.Pointer(&slots[0])
	cache := serverConfigOwnBytes(t, 0x5D4594, 251620, 4)
	refs := []unsafe.Pointer{nil, mods[0].C(), mods[1].C()}
	id := func(v uint32) int {
		for i, p := range refs {
			if v == uint32(uintptr(p)) {
				return i
			}
		}
		t.Fatal("unknown material cache pointer")
		return -1
	}
	type row struct {
		Missing           bool
		Prime, Slot, Step int
		Class             uint32
		Result, Cache     int
	}
	var rows []row
	for _, missing := range []bool{false, true} {
		for prime := 0; prime < 3; prime++ {
			for slot := 0; slot < 3; slot++ {
				for _, class := range []uint32{0, 1, 0x1000, 0x1000000, 0x2000000, 0x10000000, 0x80000000, 0xffffffff} {
					clear(names[:16])
					if missing {
						copy(names, "Absent")
					} else {
						copy(names, "Material7")
					}
					binary.LittleEndian.PutUint32(cache, uint32(uintptr(refs[prime])))
					slots[1] = refs[slot]
					unit.ObjClass = object.Class(class)
					expectedCache := prime
					for step := 0; step < 2; step++ {
						if step == 1 {
							clear(names[:16])
							copy(names, "Material7")
						}
						if expectedCache == 0 && (!missing || step == 1) {
							expectedCache = 1
						}
						want := 0
						if class&0x13001000 != 0 && slot == expectedCache {
							want = 1
						}
						got := legacy.Sub_4133D0(unit)
						cached := id(binary.LittleEndian.Uint32(cache))
						if got != want || cached != expectedCache || slots[1] != refs[slot] || uint32(unit.ObjClass) != class {
							t.Fatal("material identity/cache", missing, prime, slot, class, step, got, cached, want, expectedCache)
						}
						rows = append(rows, row{missing, prime, slot, step, class, got, cached})
					}
				}
			}
		}
	}
	interactionCapture(t, "server-runtime-material-cache", rows)
}

func TestServerRuntimeArmorConductivity(t *testing.T) {
	o := newObjectDrawingOwner(t)
	defs, free := alloc.Make([]server.Modifier{}, 2)
	t.Cleanup(free)
	old := o.c.srv.Modif.Dword_5d4594_251608
	t.Cleanup(func() { o.c.srv.Modif.Dword_5d4594_251608 = old })
	defs[0].TypeInd = 7
	defs[0].Next80 = &defs[1]
	defs[1].TypeInd = 65535
	unit, free := alloc.New(server.Object{})
	t.Cleanup(free)
	type row struct {
		Present bool
		Class   uint32
		Type    uint16
		Input   uint32
		Output  uint64
	}
	var rows []row
	for _, present := range []bool{false, true} {
		o.c.srv.Modif.Dword_5d4594_251608 = nil
		if present {
			o.c.srv.Modif.Dword_5d4594_251608 = &defs[0]
		}
		for _, class := range []uint32{0, 1, 0x1000000, 0x2000000, 0xffffffff} {
			for _, typ := range []uint16{0, 7, 8, 65535} {
				for _, bits := range []uint32{0, 0x80000000, 1, 0x80000001, 0x3f800000, 0xbf800000, 0x3eaaaaab, 0x7f7fffff, 0x7f800000, 0xff800000, 0x7fc12345, 0x7f800001} {
					unit.ObjClass = object.Class(class)
					unit.TypeInd = typ
					defs[0].DamageCoeffOrArmor64 = math.Float32frombits(bits)
					defs[1].DamageCoeffOrArmor64 = math.Float32frombits(bits)
					got := legacy.PortTestRuntimeArmorConductivity(unit)
					want := float64(0)
					if present && class&0x2000000 != 0 && (typ == 7 || typ == 65535) {
						want = float64(math.Float32frombits(bits))
					}
					if !math.IsNaN(want) && math.Float64bits(got) != math.Float64bits(want) || math.IsNaN(want) && !math.IsNaN(got) {
						t.Fatal("armor conductivity", present, class, typ, bits, got, want)
					}
					rows = append(rows, row{present, class, typ, bits, math.Float64bits(got)})
				}
			}
		}
	}
	interactionCapture(t, "server-runtime-armor-conductivity", rows)
}
