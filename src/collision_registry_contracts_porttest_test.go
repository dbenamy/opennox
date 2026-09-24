//go:build porttest

package opennox

import (
	"bytes"
	"runtime"
	"testing"
	"unsafe"

	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

func TestCollisionRegistryNames(t *testing.T) {
	names := legacy.PortTestCollisionRegistryNames()
	if len(names) != 53 {
		t.Fatalf("registrations: %d", len(names))
	}
	seen := map[string]bool{}
	for _, name := range names {
		if seen[name] {
			t.Fatal("duplicate registration", name)
		}
		seen[name] = true
	}
	a, an := server.PortTestWorldCollisionRegistry("DefaultCollide")
	b, bn := server.PortTestWorldCollisionRegistry("ElevatorCollide")
	c, cn := server.PortTestWorldCollisionRegistry("TelekinesisCollide")
	if a == nil || a != b || b != c || an != 0 || bn != 8 || cn != 0 {
		t.Fatal("default alias identities and sizes")
	}
	p, n := server.PortTestWorldCollisionRegistry("NoCollide")
	if p != nil || n != 0 {
		t.Fatal("nil collision registration")
	}
}

func TestCollisionRegistryRaw(t *testing.T) {
	s := newObjectXferOwner(t)
	u, other := newObjectXferSimple(t, s), newObjectXferSimple(t, s)
	cb, observed, restore := legacy.PortTestCollisionRawObserver()
	t.Cleanup(restore)
	u.Collide = nil
	u.CallCollide(-1, -1)
	if observed[0] != 0 {
		t.Fatal("nil collision invoked")
	}
	u.Collide = cb
	for _, a := range []uint32{0, 1, 0x7fffffff, 0x80000000, 0xffffffff} {
		for _, b := range []uint32{0, 1, 0x7fffffff, 0x80000000, 0xffffffff} {
			*observed = [4]uint32{}
			u.CallCollide(int(a), int(b))
			if *observed != ([4]uint32{1, uint32(uintptr(u.CObj())), a, b}) {
				t.Fatal("raw word forwarding", *observed, a, b)
			}
		}
	}
	for _, target := range []*server.Object{nil, u, other} {
		for _, normal := range []*types.Pointf{nil, {X: 3, Y: -4}} {
			*observed = [4]uint32{}
			server.PortTestCollisionWith(u, target, normal)
			if *observed != ([4]uint32{1, uint32(uintptr(u.CObj())), uint32(uintptr(unsafe.Pointer(target))), uint32(uintptr(unsafe.Pointer(normal)))}) {
				t.Fatal("raw pointer forwarding", *observed)
			}
			runtime.KeepAlive(target)
			runtime.KeepAlive(normal)
		}
	}
	u.Collide = nil
	*observed = [4]uint32{}
	u.CallCollide(1, 1)
	if observed[0] != 0 {
		t.Fatal("cleared collision invoked")
	}
}

// Keep stack growth observable so pointer-argument contracts exercise a callback
// that can relocate a Go stack and collect unreachable heap storage.
//
//go:noinline
func collisionRegistryGrow(depth int) byte {
	var buf [1024]byte
	for i := range buf {
		buf[i] = byte(i + depth)
	}
	if depth == 0 {
		runtime.GC()
		return buf[511]
	}
	v := collisionRegistryGrow(depth - 1)
	return v + buf[depth%len(buf)]
}

func TestCollisionRegistryDynamic(t *testing.T) {
	s := newObjectXferOwner(t)
	u, other := newObjectXferSimple(t, s), newObjectXferSimple(t, s)
	oldDeath, oldHarpoon, oldGlyph := legacy.Nox_xxx_collideDeathBall_4E9E90, legacy.Nox_xxx_collideHarpoon_4EB6A0, legacy.Nox_xxx_collideGlyph_4E9A00
	defer func() {
		legacy.Nox_xxx_collideDeathBall_4E9E90 = oldDeath
		legacy.Nox_xxx_collideHarpoon_4EB6A0 = oldHarpoon
		legacy.Nox_xxx_collideGlyph_4E9A00 = oldGlyph
	}()
	for _, name := range []string{"DeathBallCollide", "HarpoonCollide", "GlyphCollide"} {
		u.Collide = legacy.PortTestCollisionRegistryPointer(name)
		for round := uint32(1); round <= 2; round++ {
			for _, target := range []*server.Object{nil, u, other} {
				for _, pointers := range []bool{false, true} {
					normal := &types.Pointf{X: 3, Y: -4}
					if !pointers {
						normal = nil
					}
					calls := 0
					check := func(got, hit *server.Object) {
						collisionRegistryGrow(128)
						if got != u || hit != target {
							t.Fatal("late handler argument identities", name)
						}
						calls++
						got.Worth = round
					}
					legacy.Nox_xxx_collideDeathBall_4E9E90 = func(got, hit *server.Object, n *types.Pointf) {
						check(got, hit)
						if n != normal {
							t.Fatal("normal pointer identity")
						}
						if n != nil {
							if *n != (types.Pointf{X: 3, Y: -4}) {
								t.Fatal("normal did not survive stack growth/GC")
							}
							n.Y = 9
						}
					}
					legacy.Nox_xxx_collideHarpoon_4EB6A0 = check
					legacy.Nox_xxx_collideGlyph_4E9A00 = check
					u.Worth = 0
					if pointers {
						server.PortTestCollisionWith(u, target, normal)
					} else {
						u.CallCollide(int(uintptr(unsafe.Pointer(target))), 0)
					}
					if calls != 1 || u.Worth != round {
						t.Fatal("replaceable handler dispatch", name, calls)
					}
					if pointers && name == "DeathBallCollide" && normal.Y != 9 {
						t.Fatal("normal mutation not returned")
					}
					runtime.KeepAlive(target)
					runtime.KeepAlive(normal)
				}
			}
		}
	}
}

func TestCollisionRegistryDefaultAliases(t *testing.T) {
	s := newObjectXferOwner(t)
	u, other := newObjectXferSimple(t, s), newObjectXferSimple(t, s)
	for _, name := range []string{"DefaultCollide", "ElevatorCollide", "TelekinesisCollide"} {
		for _, target := range []*server.Object{nil, u, other} {
			for _, normal := range []*types.Pointf{nil, {X: 3, Y: -4}} {
				before := bytes.Clone(unsafe.Slice((*byte)(u.CObj()), int(unsafe.Sizeof(*u))))
				legacy.PortTestRegisteredCollision(u, target, normal, name)
				if !bytes.Equal(before, unsafe.Slice((*byte)(u.CObj()), len(before))) {
					t.Fatal("default collision changed object", name)
				}
			}
		}
	}
	collisionRegistryCounts(t)
}
