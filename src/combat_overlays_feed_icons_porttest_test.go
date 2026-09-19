//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"github.com/opennox/libs/spell"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestCombatOverlayFeedInit(t *testing.T) {
	o := newCombatOverlayOwner(t)
	var rows []struct {
		Cached      bool
		Values      [4]uint32
		Read, Write uint32
	}
	for _, cached := range []bool{false, true} {
		for i := 0; i < 4; i++ {
			*memmap.PtrUint32(0x5D4594, 1203844+4*uintptr(i)) = 0
			if cached {
				*memmap.PtrUint32(0x5D4594, 1203844+4*uintptr(i)) = uint32(700 + i)
			}
		}
		*o.words["feedRead"], *o.words["feedWrite"] = 99, 98
		legacy.PortTestCombatFeedInit()
		var values [4]uint32
		for i, name := range []string{"ArcherBolt", "ArcherArrow", "Bow", "CrossBow"} {
			values[i] = *memmap.PtrUint32(0x5D4594, 1203844+4*uintptr(i))
			want := uint32(o.c.Things.IndByID(name))
			if cached {
				want = uint32(700 + i)
			}
			if values[i] != want || want == 0 {
				t.Fatalf("cache %s: %d want %d", name, values[i], want)
			}
		}
		if *o.words["feedRead"] != 0 || *o.words["feedWrite"] != 0 {
			t.Fatal("initialization resets ring indices")
		}
		rows = append(rows, struct {
			Cached      bool
			Values      [4]uint32
			Read, Write uint32
		}{cached, values, *o.words["feedRead"], *o.words["feedWrite"]})
	}
	spellbookCapture(t, "combat-overlay-feed-init", rows, "15aba33dfe878d04f44f298815689c16ebfcb9a01a287142c416d9e02a44a9ba")
}
func TestCombatOverlayFeedSubstitution(t *testing.T) {
	o := newCombatOverlayOwner(t)
	legacy.PortTestCombatFeedInit()
	type record struct {
		Type          byte
		Input, Output uint16
		Stored        uint32
	}
	var rows []record
	for _, kind := range []byte{0, 1, 2, 255} {
		for _, name := range []string{"ArcherBolt", "ArcherArrow", "Bow", "CrossBow"} {
			id := uint16(o.c.Things.IndByID(name))
			var p [11]byte
			binary.LittleEndian.PutUint16(p[8:], id)
			p[10] = kind
			index := *o.words["feedWrite"]
			legacy.PortTestCombatFeedAdd(unsafe.Pointer(&p[0]))
			got := binary.LittleEndian.Uint16(p[8:])
			want := id
			if kind == 1 {
				if name == "ArcherBolt" {
					want = uint16(o.c.Things.IndByID("CrossBow"))
				}
				if name == "ArcherArrow" {
					want = uint16(o.c.Things.IndByID("Bow"))
				}
			}
			stored := *memmap.PtrUint32(0x5D4594, 1201440+24*uintptr(index))
			if got != want || stored != uint32(want) {
				t.Fatalf("substitution %s kind %d: %d/%d want %d", name, kind, got, stored, want)
			}
			rows = append(rows, record{kind, id, got, stored})
		}
	}
	spellbookCapture(t, "combat-overlay-feed-substitution", rows, "dac6b47f81df2407076c63296225158c978beb1b4f0276959d682fc52af6c8be")
}
func TestCombatOverlayFeedIcons(t *testing.T) {
	type record struct{ Name, Pixels string }
	var rows []record
	absent := make(map[[2]int]string)
	for _, present := range []bool{false, true} {
		o := newCombatOverlayOwner(t)
		ids := []int{5, 130, 60, 43, 56, 16, 15}
		var defs []server.PortTestSpellClassDef
		for _, id := range ids {
			defs = append(defs, server.PortTestSpellClassDef{Index: uint32(id), Valid: true})
		}
		configure, restore := o.c.srv.PortTestAISpellDefs()
		t.Cleanup(restore)
		configure(defs)
		raw := make([][]byte, 9)
		for i := range raw {
			raw[i] = spriteAnimationTestImage(i)
		}
		imgs, free := o.c.r.GetBag().PortTestSpriteImages(raw)
		t.Cleanup(free)
		for i, id := range ids {
			if present {
				o.c.srv.Spells.DefByInd(spell.ID(id)).Icon = unsafe.Pointer(imgs[i])
			}
		}
		oldAbility := o.c.srv.abilities.defs[1]
		t.Cleanup(func() { o.c.srv.abilities.defs[1] = oldAbility })
		o.c.srv.abilities.defs[1] = AbilityDef{}
		if present {
			o.c.srv.abilities.defs[1].icon8 = imgs[7]
		}
		typ := o.c.Things.TypeByID("Bow")
		if present {
			typ.PrettyImage = uint32(uintptr(imgs[8].C()))
		}
		legacy.PortTestCombatFeedInit()
		cases := []struct{ kind, cause, index int }{{0, 0, 6}, {1, 65535, 6}, {1, typ.Index(), 8}, {2, 0, 6}, {2, 1, 0}, {2, 12, 0}, {2, 2, 7}, {2, 4, 1}, {2, 5, 2}, {2, 9, 3}, {2, 17, 3}, {2, 15, 4}, {2, 16, 5}, {2, 255, 6}}
		for _, c := range cases {
			name := fmt.Sprintf("present=%v/type=%d/cause=%d", present, c.kind, c.cause)
			row := [6]uint32{0, 0, 0, uint32(c.cause), uint32(c.kind), 123}
			clear(o.pix.Pix)
			legacy.PortTestCombatFeedRow(&row)
			got := effectsPixelHash(o.pix)
			// Render the same row through the fallback branch with the expected actual image.
			old := *memmap.PtrUint32(0x5D4594, 1203828)
			if present {
				*memmap.PtrUint32(0x5D4594, 1203828) = uint32(uintptr(imgs[c.index].C()))
			}
			clear(o.pix.Pix)
			legacy.PortTestCombatFeedRow(&[6]uint32{})
			want := effectsPixelHash(o.pix)
			*memmap.PtrUint32(0x5D4594, 1203828) = old
			if got != want {
				t.Errorf("%s: icon pixels differ from selected metadata", name)
			}
			key := [2]int{c.kind, c.cause}
			if present {
				if got == absent[key] {
					t.Errorf("%s: actual icon did not change pixels", name)
				}
			} else {
				absent[key] = got
			}
			rows = append(rows, record{name, got})
		}
	}
	spellbookCapture(t, "combat-overlay-feed-icons", rows, "5f430d9f048f645c2ed811f38d3ab04dfb648b75fe2d6350a70dcc44a1f0bb89")
}
