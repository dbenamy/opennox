//go:build porttest

package opennox

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/opennox/libs/spell"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func TestServerPanelsSharedMasks(t *testing.T) {
	// Guards belong to raw backing storage, not extracted live globals.
	weaponBacking := unsafe.Slice(memmap.PtrUint32(0x5D4594, 1045448), 4)
	spellBacking := unsafe.Slice(memmap.PtrUint32(0x5D4594, 1045488), 6)
	spellBefore := unsafe.Slice((*uint32)(memmap.PtrOff(0x5D4594, uintptr(1045472))), 4)
	for _, region := range [][2]int{{1045448, 16}, {1045484, 28}} {
		b := unsafe.Slice((*byte)(memmap.PtrOff(0x5D4594, uintptr(region[0]))), region[1])
		old := append([]byte(nil), b...)
		t.Cleanup(func() { copy(b, old) })
	}
	weapon := memmap.PtrUint32(0x5D4594, 1045452)
	armor := memmap.PtrUint32(0x5D4594, 1045456)
	spells := memmap.PtrUint32(0x5D4594, 1045488)
	if legacy.PortTestServerPanelsWeaponPointer() != weapon || legacy.PortTestServerPanelsSpellPointer() != spells {
		t.Fatal("shared pointer identity")
	}
	src, free := alloc.Make([]uint32{}, 7)
	defer free()
	type row struct {
		Value, Weapon, Armor uint32
		Spells               [5]uint32
	}
	var rows []row
	for _, v := range []uint32{0, 1, 0x55555555, 0xaaaaaaaa, 0x7fffffff, 0x80000000, 0xffffffff} {
		*memmap.PtrUint32(0x5D4594, 1045448) = 0x13579bdf
		weaponBacking[3] = 0x2468ace0
		spellBefore[3] = 0x13579bdf
		spellBacking[5] = 0x2468ace0
		src[0], src[6] = 0x12345678, 0x87654321
		for i := 1; i <= 5; i++ {
			src[i] = v ^ uint32(i)*0x1020304
		}
		input := src[1]
		if !legacy.PortTestServerPanelsWeaponStore(&src[1]) || *weapon != input {
			t.Fatal("weapon copy/result")
		}
		src[1] ^= 0xffffffff
		if *weapon != input {
			t.Fatal("weapon aliases caller storage")
		}
		if legacy.PortTestServerPanelsArmorStore(v) != v || legacy.PortTestServerPanelsArmorLoad() != v || *armor != v {
			t.Fatal("armor store/load/result")
		}
		legacy.PortTestServerPanelsSpellStore(&src[1])
		var got [5]uint32
		copy(got[:], unsafe.Slice(spells, 5))
		for i := 0; i < 5; i++ {
			if got[i] != src[i+1] {
				t.Fatal("spell copy width")
			}
		}
		src[1] ^= 0xffffffff
		if *spells != got[0] {
			t.Fatal("spell mask aliases caller storage")
		}
		if src[0] != 0x12345678 || src[6] != 0x87654321 || *memmap.PtrUint32(0x5D4594, 1045448) != 0x13579bdf || weaponBacking[3] != 0x2468ace0 || spellBefore[3] != 0x13579bdf || spellBacking[5] != 0x2468ace0 {
			t.Fatal("copy guard")
		}
		rows = append(rows, row{v, *weapon, *armor, got})
	}
	spellbookCapture(t, "server-panels-shared-masks", rows, "c6ae3fa989143e91821aff4b53709b1ddbc74c872019b4edbbb06aff8bc03f00")
}

func TestServerPanelsSpellState(t *testing.T) {
	o := newServerOptionsOwner(t)
	old := o.c.srv.Spells
	t.Cleanup(func() { o.c.srv.Spells = old })
	words, free := alloc.Make([]uint32{}, 7)
	defer free()
	type row struct {
		AllowAll           bool
		Last               int
		Pattern            uint32
		Snapshot           [5]uint32
		Result             uint32
		PointerWord, Apply int
		Enabled            []bool
	}
	var rows []row
	// The last definition covers each distinct scalar/pointer return branch.
	for _, allowAll := range []bool{false, true} {
		for last := 0; last < 5; last++ {
			for _, pattern := range []uint32{0, 0x55555555, 0xaaaaaaaa, 0xffffffff} {
				t.Run(fmt.Sprintf("%t-%d-%x", allowAll, last, pattern), func(t *testing.T) {
					var defs []server.PortTestSpellClassDef
					for i := 1; i <= 137; i++ {
						if i%11 == 0 || i == 136 && last == 0 {
							continue
						}
						flags := uint32(0x01000000) << uint(i%3)
						if i%5 == 0 {
							flags = 0x80
						}
						valid := i%7 != 0
						if i == 136 {
							valid = last != 1
							if last == 2 {
								flags = 0x80
							}
						}
						defs = append(defs, server.PortTestSpellClassDef{Index: uint32(i), Flags: flags, Valid: valid})
					}
					s := server.PortTestSpellClassServer(defs)
					o.c.srv.Spells = s.Spells
					o.c.srv.Spells.AllowAll = allowAll
					for _, d := range defs {
						sp := o.c.srv.Spells.DefByInd(spell.ID(d.Index))
						sp.Title = fmt.Sprintf("Panel spell %d", d.Index)
						sp.Enabled = pattern&(1<<uint(d.Index%32)) != 0
						if d.Index == 136 {
							sp.Enabled = last == 4
						}
					}
					words[0], words[6] = 0x12345678, 0x87654321
					for i := 1; i <= 5; i++ {
						words[i] = 0x12341234
					}
					result, ptr := legacy.PortTestServerPanelsSpellSnapshot(&words[1])
					var want [5]uint32
					for i := range want {
						want[i] = 0xffffffff
					}
					for _, d := range defs {
						if d.Index > 136 {
							continue
						}
						sp := o.c.srv.Spells.DefByInd(spell.ID(d.Index))
						eligible := d.Flags&0x07000000 != 0 || allowAll
						if d.Valid && eligible && !sp.Enabled {
							want[d.Index/32] &^= 1 << uint(d.Index%32)
						}
					}
					var got [5]uint32
					copy(got[:], words[1:6])
					if got != want {
						t.Fatalf("snapshot got %x want %x", got, want)
					}
					wantResult, wantPtr := uint32(0), -1
					if last == 2 && !allowAll {
						wantResult = 0x80
					} else if last >= 2 && last != 4 {
						wantPtr = 4
					} else if last == 4 {
						wantResult = 1
					}
					if result != wantResult || ptr != wantPtr {
						t.Fatalf("snapshot return %x/%d want %x/%d", result, ptr, wantResult, wantPtr)
					}
					for i := 1; i <= 5; i++ {
						words[i] = pattern
					}
					outside := o.c.srv.Spells.DefByInd(137).Enabled
					applied := legacy.PortTestServerPanelsSpellApply(&words[1])
					wantApply := 1
					if last == 0 {
						wantApply = 0
					}
					if applied != wantApply {
						t.Fatal("apply result", applied)
					}
					var enabled []bool
					for _, d := range defs {
						sp := o.c.srv.Spells.DefByInd(spell.ID(d.Index))
						wantEnabled := pattern&(1<<uint(d.Index%32)) != 0
						if d.Index == 137 {
							wantEnabled = outside
						}
						if sp.Enabled != wantEnabled {
							t.Fatalf("apply id %d", d.Index)
						}
						enabled = append(enabled, sp.Enabled)
					}
					if words[0] != 0x12345678 || words[6] != 0x87654321 {
						t.Fatal("spell state guards")
					}
					for _, v := range words[1:6] {
						if v != pattern {
							t.Fatal("apply mutated input")
						}
					}
					rows = append(rows, row{allowAll, last, pattern, got, result, ptr, applied, enabled})
				})
			}
		}
	}
	spellbookCapture(t, "server-panels-spell-state", rows, "8e020422a42b80126570453036a9f846ada6e862185cfad4d2fab2150df6f583")
}
