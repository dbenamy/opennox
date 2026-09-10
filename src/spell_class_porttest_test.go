//go:build porttest

package opennox

import (
	"testing"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

func TestSpellClassEligibilityABI(t *testing.T) {
	var defs []server.PortTestSpellClassDef
	for mask := uint32(0); mask < 8; mask++ {
		for i, noise := range []uint32{0, 1, 0xffffff, 0x80000000, 0x88000000, 0xf8ffffff, 0x5005a5a5, 0x2a55ffff} {
			flags := noise &^ 0x07000000
			for bit := uint(0); bit < 3; bit++ {
				if mask&(1<<bit) != 0 {
					flags |= 1 << (24 + bit)
				}
			}
			defs = append(defs, server.PortTestSpellClassDef{Index: uint32(len(defs) + 1), Flags: flags, Valid: i%2 == 0})
		}
	}
	defs = append(defs, server.PortTestSpellClassDef{Index: 0x7fffffff, Flags: 0x01000000}, server.PortTestSpellClassDef{Index: 0, Flags: 0x07000000, Valid: true}, server.PortTestSpellClassDef{Index: 0xffffffff, Flags: 0x07000000, Valid: true})
	var classes []uint32
	for i := uint32(0); i < 256; i++ {
		classes = append(classes, i)
	}
	classes = append(classes, 256, 257, 258, 0x10001, 0xffffff01, 0xffffffff, 0x80000000, 0x7fffffff)
	ids := []uint32{0, 65, 136, 137, 0x10000, 0xffffffff, 0x80000000}
	for _, d := range defs {
		if int32(d.Index) > 0 {
			ids = append(ids, d.Index)
		}
	}
	var calls []legacy.PortTestSpellClassCall
	for _, class := range classes {
		for _, id := range ids {
			calls = append(calls, legacy.PortTestSpellClassCall{PlayerClass: class, Spell: id})
		}
	}
	for _, allowAll := range []bool{false, true} {
		got := legacy.PortTestSpellClass(defs, calls, allowAll)
		if len(got) != len(calls) {
			t.Fatal("missing class results")
		}
		for i, c := range calls {
			flags, exists := uint32(0), false
			if int32(c.Spell) > 0 {
				for _, d := range defs {
					if d.Index == c.Spell {
						flags, exists = d.Flags, true
						break
					}
				}
			}
			any := flags/(1<<24)%2 == 1 || (allowAll && exists)
			own := false
			if c.PlayerClass == 1 {
				own = flags/(1<<25)%2 == 1
			}
			if c.PlayerClass == 2 {
				own = flags/(1<<26)%2 == 1
			}
			want := 9
			if (c.PlayerClass == 1 || c.PlayerClass == 2) && (any || own) {
				want = 0
			}
			if got[i].PlayerClass != c.PlayerClass || got[i].Spell != c.Spell || got[i].Result != want {
				t.Fatalf("class=%08x spell=%08x flags=%08x allowAll=%v got=%+v want=%d", c.PlayerClass, c.Spell, flags, allowAll, got[i], want)
			}
		}
	}
}
