//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"reflect"
	"testing"
	"unsafe"
)

func TestCombatOverlayHealthList(t *testing.T) {
	type record struct {
		Name                   string
		Added, Removed, Reused []combatHealthRecord
	}
	var rows []record
	for _, amount := range []int16{-32768, -1, 0, 1, 32767} {
		for remove := 0; remove < 3; remove++ {
			name := fmt.Sprintf("amount=%d/remove=%d", amount, remove)
			t.Run(name, func(t *testing.T) {
				o := newCombatOverlayOwner(t)
				for code := uint32(1); code <= 3; code++ {
					legacy.PortTestCombatHealthAdd(code, amount)
					o.c.Inp.Tick()
				}
				want := []combatHealthRecord{{3, amount, 2}, {2, amount, 1}, {1, amount, 0}}
				r := record{Name: name, Added: o.healthRecords(t)}
				if !reflect.DeepEqual(r.Added, want) {
					t.Fatal("health insertion/timestamp")
				}
				p := unsafe.Pointer(uintptr(*o.words["health"]))
				for i := 0; i < remove; i++ {
					p = *(*unsafe.Pointer)(unsafe.Add(p, 12))
				}
				legacy.PortTestCombatHealthRemove(p)
				want = append(want[:remove], want[remove+1:]...)
				r.Removed = o.healthRecords(t)
				if !reflect.DeepEqual(r.Removed, want) {
					t.Fatal("health removal")
				}
				legacy.PortTestCombatHealthClear()
				if *o.words["health"] != 0 {
					t.Fatal("clear health head")
				}
				legacy.PortTestCombatHealthAdd(99, amount)
				r.Reused = o.healthRecords(t)
				if !reflect.DeepEqual(r.Reused, []combatHealthRecord{{99, amount, 3}}) {
					t.Fatal("health reuse")
				}
				legacy.PortTestCombatHealthDestroy()
				if *o.pools["health"] != nil || *o.words["health"] != 0 || *o.words["font"] != 0 {
					t.Fatal("destroy health state")
				}
				rows = append(rows, r)
			})
		}
	}
	spellbookCapture(t, "combat-overlay-health-list", rows, "c47dba99defead97ece290a555bf0d67e1d3f2d223a822e7f0e51cf50cbbff6b")
}
func TestCombatOverlayHealthCapacity(t *testing.T) {
	o := newCombatOverlayOwner(t)
	for i := 0; i < 32; i++ {
		legacy.PortTestCombatHealthAdd(uint32(i), int16(i))
	}
	before := o.healthRecords(t)
	if len(before) != 32 {
		t.Fatal("health capacity")
	}
	legacy.PortTestCombatHealthAdd(32, 32)
	if !reflect.DeepEqual(o.healthRecords(t), before) {
		t.Fatal("health capacity overflow")
	}
	legacy.PortTestCombatHealthRemove(unsafe.Pointer(uintptr(*o.words["health"])))
	legacy.PortTestCombatHealthAdd(32, 32)
	after := o.healthRecords(t)
	if len(after) != 32 || after[0].Code != 32 {
		t.Fatal("health slot reuse")
	}
	spellbookCapture(t, "combat-overlay-health-capacity", [][]combatHealthRecord{before, after}, "ad754626403d9d3ffb690ac246d568aea509405f002db3014a35cc4d31072f0e")
}
func TestCombatOverlayHealthDrawing(t *testing.T) {
	type record struct {
		Name, Pixels string
		Live         []combatHealthRecord
	}
	var rows []record
	hashes := make(map[string]string)
	for _, amount := range []int16{-32768, -1, 0, 1, 32767} {
		for _, local := range []bool{false, true} {
			for _, age := range []uint32{0, 1, 30, 31, 32} {
				name := fmt.Sprintf("amount=%d/local=%t/age=%d", amount, local, age)
				t.Run(name, func(t *testing.T) {
					o := newCombatOverlayOwner(t)
					dr := o.drawable(7, image.Pt(320, 300))
					dr.ZVal = 5
					dr.ZSizeMax = 9.75
					if local {
						*memmap.PtrPtr(0x852978, 8) = dr.C()
					}
					legacy.PortTestCombatHealthAdd(7, amount)
					current := unsafe.Pointer(uintptr(*o.words["health"]))
					objectXferSetWord(current, 8, uint32(o.c.GetInputSeq())-age)
					legacy.PortTestCombatHealthAdd(8, 123)
					objectXferSetWord(unsafe.Pointer(uintptr(*o.words["health"])), 8, uint32(o.c.GetInputSeq())-32)
					blank := effectsPixelHash(o.pix)
					legacy.PortTestCombatHealthDraw(o.c.Viewport(), dr)
					r := record{name, effectsPixelHash(o.pix), o.healthRecords(t)}
					if age <= 30 {
						if len(r.Live) != 1 || r.Live[0].Code != 7 || r.Pixels == blank {
							t.Fatalf("live health state %+v; pixels unchanged=%t", r.Live, r.Pixels == blank)
						}
					} else if len(r.Live) != 0 || r.Pixels != blank {
						t.Fatal("expired health text/state")
					}
					hashes[name] = r.Pixels
					rows = append(rows, r)
				})
			}
		}
	}
	for _, local := range []bool{false, true} {
		negative := fmt.Sprintf("amount=-1/local=%t/age=0", local)
		positive := fmt.Sprintf("amount=1/local=%t/age=0", local)
		if hashes[negative] == hashes[positive] {
			t.Fatal("health sign must select distinct text color")
		}
	}
	spellbookCapture(t, "combat-overlay-health-drawing", rows, "841071d01acebf704edbfb14b52bc42dcd4a164f579a8e59938b1ef23e124af5")
}
