//go:build porttest

package opennox

import (
	"runtime"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

// Preserve callback equality, distinctness and lifetime in the existing foreign
// object/type storage without asserting linker-specific address bits.
func TestDeathIdentityLifetime(t *testing.T) {
	typ, freeType := alloc.New(server.ObjectType{})
	defer freeType()
	obj, freeObject := alloc.New(server.Object{})
	defer freeObject()
	oldGlyph := legacy.Nox_xxx_dieGlyph_54DF30
	defer func() { legacy.Nox_xxx_dieGlyph_54DF30 = oldGlyph }()
	calls := 0
	legacy.Nox_xxx_dieGlyph_54DF30 = func(got *server.Object) {
		if got != obj {
			t.Fatal("stored death key object identity")
		}
		calls++
	}
	names := legacy.PortTestDeathRegistryNames()
	seen := make(map[unsafe.Pointer]string)
	for _, name := range names {
		key, size := server.PortTestDeathRegistry(name)
		if key == nil {
			t.Fatal("nil death key", name)
		}
		if previous, ok := seen[key]; ok {
			t.Fatal("aliased death keys", previous, name)
		}
		seen[key] = name
		wantSize := uintptr(0)
		if name == "CreateObjectDie" || name == "SpawnObjectDie" {
			wantSize = 132
		}
		if size != wantSize {
			t.Fatal("death data size", name, size, wantSize)
		}
		typ.Death = key
		obj.Death = typ.Death
		collisionRegistryGrow(128)
		if name == "GlyphDie" {
			obj.CallDeath()
			if calls != 1 {
				t.Fatal("stored key dispatch after GC", calls)
			}
		}
		got, _ := server.PortTestDeathRegistry(name)
		if typ.Death != key || obj.Death != key || got != key {
			t.Fatal("unstable death key", name)
		}
	}
	if len(names) != 14 || len(seen) != 14 {
		t.Fatal("death identity count", len(names), len(seen))
	}
	runtime.GC()
	for key, name := range seen {
		got, _ := server.PortTestDeathRegistry(name)
		if got != key {
			t.Fatal("death key changed after GC", name)
		}
	}
}
