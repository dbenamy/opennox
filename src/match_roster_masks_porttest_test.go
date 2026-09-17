//go:build porttest

package opennox

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/opennox/libs/object"
)

func TestMatchRosterObjectMasks(t *testing.T) {
	o := newMatchRosterOwner(t)
	oldList := o.s.Objs.First()
	t.Cleanup(func() { o.s.Objs.SetObjects(oldList) })
	typ := o.s.Types.ByID("PortCreatureMonster")
	if typ == nil {
		t.Fatal("fixture type")
	}
	type objectState struct {
		Words [4]uint32
		Slots [32]uint32
	}
	snapshot := func(i int) objectState {
		u := &o.units[i]
		return objectState{[4]uint32{u.Field35, u.Field36, u.Field37, u.Field38}, u.Field140}
	}
	type row struct {
		Name    string
		Return  uint32
		Objects [3]objectState
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "match-roster-object-masks", rows, "f7346cce75a8474253069a1638e5c1268a79220fcc4ad252a4901edb22a7be1e")
	}()
	for _, op := range []string{"object-clear-mask", "object-resync-mask", "object-report-mask"} {
		slots := []uint32{0, 1, 7, 31}
		if op != "object-report-mask" {
			slots = append(slots, 32, 127, 128, 255)
		}
		for _, count := range []int{0, 1, 3} {
			for _, slot := range slots {
				for _, class := range []uint32{0, 2, 4, 0x400000, 0x20000000, 0x20400000, 0x1000000, 0xffffffff} {
					for _, flags := range []uint32{0, 0x1000020} {
						for _, seed := range []uint32{0, 0xa5a55a5a, 0xffffffff} {
							name := fmt.Sprintf("%s/count%d/slot%d/class%x/flags%x/seed%x", op, count, slot, class, flags, seed)
							t.Run(name, func(t *testing.T) {
								var before [3]objectState
								var extras [3]uint32
								for i := range o.units {
									u := &o.units[i]
									u.TypeInd = uint16(typ.Ind())
									u.ObjClass = object.Class(class)
									u.ObjFlags = object.Flags(flags)
									u.ObjSubClass = object.SubClass(seed)
									u.Field5 = seed
									u.Field33 = seed
									u.Buffs = seed
									u.ZVal = float32(i)
									u.Field35 = seed
									u.Field36 = seed ^ 0xffffffff
									u.Field37 = seed
									u.Field38 = seed ^ 0xffffffff
									for j := range u.Field140 {
										u.Field140[j] = seed ^ uint32(j*0x1001)
									}
									*(*unsafe.Pointer)(unsafe.Add(u.CObj(), 444)) = nil
									if i+1 < count {
										*(*unsafe.Pointer)(unsafe.Add(u.CObj(), 444)) = o.units[i+1].CObj()
									}
									before[i] = snapshot(i)
									for bit := uint(1); bit < 0x10000; bit <<= 1 {
										if u.Sub_4E4C90(bit) {
											extras[i] |= uint32(bit) << 16
										}
									}
								}
								o.s.Objs.SetObjects(nil)
								if count > 0 {
									o.s.Objs.SetObjects(&o.units[0])
								}
								rv := matchRosterCall(op, nil, nil, nil, slot)
								var after [3]objectState
								for i := range o.units {
									after[i] = snapshot(i)
									want := before[i]
									if i < count {
										mask := uint32(1) << (slot & 31)
										switch op {
										case "object-clear-mask":
											want.Words[0] &^= mask
											want.Words[1] &^= mask
										case "object-resync-mask":
											want.Words[3] |= mask
											if flags&0x20 == 0 && class&0x20400006 == 0 {
												want.Words[2] &^= mask
											}
										case "object-report-mask":
											want.Words[3] |= mask
											want.Slots[slot] &= 0xfff
											if class&0x20400000 == 0 {
												want.Words[2] &^= mask
											} else {
												want.Slots[slot] |= extras[i]
											}
										}
									}
									if after[i] != want {
										t.Fatalf("object %d changed outside its report contract", i)
									}
								}
								if rv != 0 {
									t.Fatal("iteration result", rv)
								}
								rows = append(rows, row{name, rv, after})
							})
						}
					}
				}
			}
		}
	}
}
