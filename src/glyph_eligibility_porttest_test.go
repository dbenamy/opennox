//go:build porttest

package opennox

import (
	"reflect"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func checkGlyphCases(t *testing.T, cc, ic uint32, calls []legacy.PortTestGlyphEligibilityCall) {
	t.Helper()
	got := legacy.PortTestGlyphEligibility(cc, ic, calls)
	if len(got) != len(calls) {
		t.Fatal("missing glyph results")
	}
	for i, c := range calls {
		cache := &cc
		if c.Kind == "item" {
			cache = &ic
		}
		lookups := 0
		if *cache == 0 {
			*cache = c.LookupType
			lookups = 1
		}
		want := 0
		var trace []legacy.PortTestGlyphEligibilityTrace
		glyphAllowed := c.Type != *cache || c.PlayerClass == 1
		if c.Kind == "client" {
			if c.Player && glyphAllowed {
				want = 1
			}
		} else if c.Drawable && c.LocalPlayer && c.Player && glyphAllowed {
			if c.Cheat {
				want = 1
			} else {
				trace = []legacy.PortTestGlyphEligibilityTrace{{Class: c.Class, Subclass: c.Subclass, Type: c.Type}}
				// The current 386 C ABI masks the shift count to five bits, then narrows
				// the result to a byte. Verify even invalid class bytes against that ABI.
				bit := uint(c.PlayerClass) % 32
				if bit < 8 && (uint32(c.ClassMask)/(1<<bit))%2 != 0 {
					want = 1
				}
			}
		}
		g := got[i]
		if g.Return != want || g.ClientCache != cc || g.ItemCache != ic || g.LookupCalls != lookups || !reflect.DeepEqual(g.ClassMaskTrace, trace) || !g.DrawableUntouched || !g.PlayerUntouched {
			t.Fatalf("case %d %+v got=%+v want=%d caches=%x/%x lookups=%d trace=%v", i, c, g, want, cc, ic, lookups, trace)
		}
	}
}

func TestGlyphEligibilityClasses(t *testing.T) {
	var calls []legacy.PortTestGlyphEligibilityCall
	for cl := 0; cl < 256; cl++ {
		for _, mask := range []byte{0, 1, 2, 4, 7, 0x80, 0xaa, 0xff} {
			c := legacy.PortTestGlyphEligibilityCall{Drawable: true, Player: true, LocalPlayer: true, PlayerClass: byte(cl), Type: 77, Class: 0x80000002, Subclass: 0xa5f0aaaa, LookupType: 77, ClassMask: mask, Kind: "client"}
			calls = append(calls, c)
			c.Kind, c.Type = "item", 99
			calls = append(calls, c)
			c.Wrapper = true
			calls = append(calls, c)
		}
	}
	// Every mask for the three legitimate classes, including glyph and non-glyph.
	for cl := 0; cl < 3; cl++ {
		for mask := 0; mask < 256; mask++ {
			for _, typ := range []uint32{77, 0xffffffff} {
				calls = append(calls, legacy.PortTestGlyphEligibilityCall{Kind: "item", Wrapper: mask%2 != 0, Drawable: true, Player: true, LocalPlayer: true, PlayerClass: byte(cl), Type: typ, Class: 0xffffffff, Subclass: 0x80000000, LookupType: 77, ClassMask: byte(mask)})
			}
		}
	}
	checkGlyphCases(t, 0, 0, calls)
}

func TestGlyphEligibilityGates(t *testing.T) {
	for _, lookup := range []uint32{0, 77, 0x80000000} {
		var calls []legacy.PortTestGlyphEligibilityCall
		for state := 0; state < 32; state++ {
			for cl := byte(0); cl < 3; cl++ {
				for _, kind := range []string{"client", "item"} {
					c := legacy.PortTestGlyphEligibilityCall{Kind: kind, Drawable: state&1 != 0, Player: state&2 != 0, LocalPlayer: state&4 != 0, Cheat: state&8 != 0, PlayerClass: cl, LookupType: lookup, Type: lookup, ClassMask: 0xff, Class: 0x80000000, Subclass: 0xffffffff}
					if state&16 != 0 {
						c.Type ^= 0xffffffff
					}
					if kind == "client" && !c.Drawable && c.Player {
						continue
					} // C would dereference nil.
					calls = append(calls, c)
					if kind == "item" {
						c.Wrapper = true
						calls = append(calls, c)
					}
				}
			}
		}
		checkGlyphCases(t, 0, 0, calls)
	}
}

func TestGlyphEligibilityCaches(t *testing.T) {
	for _, initial := range [][2]uint32{{0, 0}, {77, 99}, {0xffffffff, 0x80000000}} {
		var calls []legacy.PortTestGlyphEligibilityCall
		for _, lookup := range []uint32{0, 0, 77, 99, 0x80000000, 0, 0xffffffff} {
			for _, kind := range []string{"client", "item"} {
				// Lookup precedes missing-player and missing-drawable gates. Independent
				// cache fill must survive subsequent registry changes and zero lookups.
				calls = append(calls, legacy.PortTestGlyphEligibilityCall{Kind: kind, LookupType: lookup})
				calls = append(calls, legacy.PortTestGlyphEligibilityCall{Kind: kind, LookupType: lookup, Drawable: true, Player: true, LocalPlayer: true, PlayerClass: 1, Type: lookup, ClassMask: 2})
			}
		}
		checkGlyphCases(t, initial[0], initial[1], calls)
	}
}
