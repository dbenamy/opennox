//go:build porttest

package opennox

import "testing"

func TestClientDrawableUnknownTypeDoesNotExhaustPool(t *testing.T) {
	c, _ := newEffectsTestOwner(t)
	// This owner has512 slots. Failed lookups must leave capacity for valid types.
	for i := 0; i < 1024; i++ {
		if dr := c.Nox_new_drawable_for_thing(0); dr != nil {
			t.Fatal("unknown type produced a drawable")
		}
	}
	if c.Objs.Count != 0 || c.Objs.List1 != nil {
		t.Fatal("failed lookups changed live drawable state")
	}
	dr := c.Nox_new_drawable_for_thing(4)
	if dr == nil {
		t.Fatal("failed type lookups exhausted the drawable pool")
	}
	c.Nox_xxx_spriteDeleteStatic_45A4E0_drawable(dr)
	if c.Objs.Count != 0 {
		t.Fatal("valid drawable cleanup")
	}
}
