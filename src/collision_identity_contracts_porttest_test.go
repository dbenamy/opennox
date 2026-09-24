//go:build porttest

package opennox

import (
	"runtime"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy"
)

// Identity means equality and lifetime, not a linker-specific code address.
// Qualify these contracts against the original C keys before moving the slots.
func TestCollisionIdentityLifetime(t *testing.T) {
	o := newCollisionCoreOwner(t)
	u := newObjectXferSimple(t, o.s)
	names := legacy.PortTestCollisionRegistryNames()
	keys := make(map[string]unsafe.Pointer)
	unique := make(map[unsafe.Pointer]string)
	for _, name := range names {
		key := legacy.PortTestCollisionRegistryPointer(name)
		if key == nil {
			t.Fatal("nil registered key", name)
		}
		if prev, ok := unique[key]; ok && !((name == "ElevatorCollide" || name == "TelekinesisCollide") && prev == "DefaultCollide") {
			t.Fatal("unexpected identity alias", prev, name)
		}
		if _, ok := unique[key]; !ok {
			unique[key] = name
		}
		keys[name] = key
		// The real object uses stable non-Go-heap storage. Force stack growth and GC
		// with its key stored in the existing pointer field, then read it back.
		u.Collide = key
		collisionRegistryGrow(128)
		if u.Collide != key || legacy.PortTestCollisionRegistryPointer(name) != key {
			t.Fatal("unstable callback identity", name)
		}
	}
	if len(names) != 53 || len(unique) != 51 {
		t.Fatal("registration/identity counts", len(names), len(unique))
	}
	runtime.GC()
	for name, key := range keys {
		if legacy.PortTestCollisionRegistryPointer(name) != key {
			t.Fatal("identity changed after GC", name)
		}
	}
}

func TestCollisionIdentityActivationReturn(t *testing.T) {
	o := newCollisionCoreOwner(t)
	u := newObjectXferSimple(t, o.s)
	u.Update = nil
	u.ObjFlags = 0x40
	for _, name := range legacy.PortTestCollisionRegistryNames() {
		u.Collide = legacy.PortTestCollisionRegistryPointer(name)
		want := int32(int8(uintptr(u.Collide)))
		if got := legacy.PortTestCollisionCore("activate", u, nil, nil, 0); got != want {
			t.Fatalf("%s activation return=%d want current key low byte %d", name, got, want)
		}
	}
}
