//go:build porttest

package opennox

import (
	"testing"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

// The registered C address is unchanged; its Go target is a dynamic function
// variable, so each invocation must observe the current override.
func TestGlyphDeathRegistryUsesCurrentOverride(t *testing.T) {
	o := newMatchRosterOwner(t)
	o.reset()
	u := &o.units[0]

	old := legacy.Nox_xxx_dieGlyph_54DF30
	t.Cleanup(func() { legacy.Nox_xxx_dieGlyph_54DF30 = old })

	var first []*server.Object
	legacy.Nox_xxx_dieGlyph_54DF30 = func(got *server.Object) { first = append(first, got) }
	legacy.PortTestDeathRegisteredProjectile(u, "GlyphDie")
	if len(first) != 1 || first[0] != u {
		t.Fatalf("first glyph override calls = %v, want exactly [%p]", first, u)
	}

	var second []*server.Object
	legacy.Nox_xxx_dieGlyph_54DF30 = func(got *server.Object) { second = append(second, got) }
	legacy.PortTestDeathRegisteredProjectile(u, "GlyphDie")
	if len(second) != 1 || second[0] != u {
		t.Fatalf("second glyph override calls = %v, want exactly [%p]", second, u)
	}
	if len(first) != 1 {
		t.Fatalf("second invocation changed first override count: %d", len(first))
	}
}
