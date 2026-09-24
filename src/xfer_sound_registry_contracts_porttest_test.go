//go:build porttest

package opennox

import (
	"errors"
	"fmt"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func TestXferSoundRegistryNames(t *testing.T) {
	names := legacy.PortTestXferRegistryNames()
	if len(names) != 28 {
		t.Fatalf("xfer registrations: %d", len(names))
	}
	seen := map[unsafe.Pointer]bool{}
	for _, name := range names {
		p := legacy.PortTestXferRegistryPointer(name)
		if p == nil || seen[p] {
			t.Fatal("xfer identity", name)
		}
		seen[p] = true
	}
	a := legacy.PortTestDamageSoundRegistryPointer("DefaultDamageSound")
	b := legacy.PortTestDamageSoundRegistryPointer("PlayerDamageSound")
	if a == nil || b == nil || a == b {
		t.Fatal("sound identities")
	}
	if server.DefaultXfer != legacy.PortTestXferRegistryPointer("DefaultXfer") || server.DefaultDamageSound != a {
		t.Fatal("default registry identities")
	}
}

func TestXferSoundRegistryRawXfer(t *testing.T) {
	s := newObjectXferOwner(t)
	u, other := newObjectXferSimple(t, s), newObjectXferSimple(t, s)
	cb, words, result, restore := legacy.PortTestXferSoundRawObserver()
	t.Cleanup(restore)
	u.Xfer = cb
	for _, value := range []int32{0, 1, 256, -1, -2147483648, 2147483647} {
		for _, arg := range []unsafe.Pointer{nil, u.CObj(), other.CObj()} {
			for _, outer := range []bool{false, true} {
				*words = [3]uint32{}
				*result = value
				var err error
				if outer {
					err = asObjectS(u).CallXfer(arg)
				} else {
					err = u.CallXfer(arg)
				}
				if (err != nil) != (value == 0) {
					t.Fatal("exact nonzero xfer success", value, err)
				}
				if err != nil && err.Error() != fmt.Sprintf("xfer for %s failed", u.String()) {
					t.Fatal("raw error text", err)
				}
				if *words != ([3]uint32{1, uint32(uintptr(u.CObj())), uint32(uintptr(arg))}) {
					t.Fatal("raw xfer arguments", *words)
				}
			}
		}
	}
}

func TestXferSoundRegistryDefaultDynamic(t *testing.T) {
	s := newObjectXferOwner(t)
	u, other := newObjectXferSimple(t, s), newObjectXferSimple(t, s)
	u.Xfer = legacy.PortTestXferRegistryPointer("DefaultXfer")
	old := legacy.Nox_xxx_XFerDefault4F49A0
	defer func() { legacy.Nox_xxx_XFerDefault4F49A0 = old }()
	original := cryptfile.Global()
	defer cryptfile.SetGlobal(original)
	for round := 0; round < 2; round++ {
		cf := new(cryptfile.CryptFile)
		cryptfile.SetGlobal(cf)
		for _, fail := range []bool{false, true} {
			for _, arg := range []unsafe.Pointer{nil, u.CObj(), other.CObj()} {
				calls := 0
				legacy.Nox_xxx_XFerDefault4F49A0 = func(gotCF *cryptfile.CryptFile, got *server.Object, ctx unsafe.Pointer) error {
					if gotCF != cf || got != u || ctx != arg {
						t.Fatal("default xfer late-bound arguments")
					}
					calls++
					if fail {
						return errors.New("deliberate original-handler failure")
					}
					return nil
				}
				err := u.CallXfer(arg)
				if calls != 1 || (err != nil) != fail {
					t.Fatal("default handler replacement", calls, err)
				}
				if err != nil && err.Error() != fmt.Sprintf("xfer for %s failed", u.String()) {
					t.Fatal("default callback error convention", err)
				}
			}
		}
	}
}

func TestXferSoundRegistryDamageOwner(t *testing.T) {
	s := newObjectXferOwner(t)
	u, source, weapon := newObjectXferSimple(t, s), newObjectXferSimple(t, s), newObjectXferSimple(t, s)
	hp, free := alloc.New(server.HealthData{})
	u.HealthData = hp
	t.Cleanup(func() { u.HealthData = nil; free() })
	// Simple objects avoid unrelated monster/player admission and audio branches;
	// damageDefault still selects the actual source and invokes the live sound owner.
	u.ObjClass = 0
	source.ObjClass = 0
	weapon.ObjClass = 0
	oldDefault, oldPlayer := legacy.Nox_xxx_soundDefaultDamageSound_532E20, legacy.Nox_xxx_soundPlayerDamageSound_5328B0
	oldPointer := server.DefaultDamageSound
	defer func() {
		legacy.Nox_xxx_soundDefaultDamageSound_532E20 = oldDefault
		legacy.Nox_xxx_soundPlayerDamageSound_5328B0 = oldPlayer
		server.DefaultDamageSound = oldPointer
	}()
	cb, words, result, restore := legacy.PortTestXferSoundRawObserver()
	t.Cleanup(restore)
	// A nil object slot must call the fixed default wrapper even if this registry
	// default pointer is changed. The wrapper's Go handler remains replaceable.
	server.DefaultDamageSound = legacy.PortTestDamageSoundRegistryPointer("PlayerDamageSound")
	for round := 0; round < 2; round++ {
		for _, slot := range []string{"nil", "raw", "DefaultDamageSound", "PlayerDamageSound"} {
			for _, ret := range []int32{0, 1, 256, -1} {
				for _, src := range []*server.Object{nil, source} {
					for _, wpn := range []*server.Object{nil, weapon} {
						*hp = server.HealthData{Cur: 80, Field2: 100, Max: 100}
						u.ObjFlags = 0
						*words = [3]uint32{}
						*result = ret
						u.DamageSound = nil
						if slot == "raw" {
							u.DamageSound = cb
						} else if slot != "nil" {
							u.DamageSound = legacy.PortTestDamageSoundRegistryPointer(slot)
						}
						actual := wpn
						if actual == nil {
							actual = src
						}
						called := ""
						calls := 0
						observe := func(name string, got, other *server.Object) int {
							if got != u || other != actual || got.HealthData.Cur != 80 {
								t.Fatal("sound owner arguments/timing", slot, name)
							}
							called = name
							calls++
							return int(ret)
						}
						legacy.Nox_xxx_soundDefaultDamageSound_532E20 = func(a, b *server.Object) int { return observe("default", a, b) }
						legacy.Nox_xxx_soundPlayerDamageSound_5328B0 = func(a, b *server.Object) int { return observe("player", a, b) }
						rv := legacy.PortTestDamageSoundOwner(u, src, wpn)
						if rv != 1 || hp.Cur != 68 {
							t.Fatal("sound result must not alter damage", slot, ret, rv, hp.Cur)
						}
						if slot == "raw" {
							if calls != 0 || *words != ([3]uint32{1, uint32(uintptr(u.CObj())), uint32(uintptr(unsafe.Pointer(actual)))}) {
								t.Fatal("raw sound forwarding", *words, calls)
							}
						} else {
							want := "default"
							if slot == "PlayerDamageSound" {
								want = "player"
							}
							if calls != 1 || called != want || words[0] != 0 {
								t.Fatal("sound selection/late handler", slot, called, calls)
							}
						}
					}
				}
			}
		}
	}
}
