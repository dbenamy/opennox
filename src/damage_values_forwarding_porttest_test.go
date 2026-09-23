//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"testing"
)

func TestDamageValueRawForwarding(t *testing.T) {
	target, freeTarget := alloc.New(server.Object{})
	defer freeTarget()
	source, freeSource := alloc.New(server.Object{})
	defer freeSource()
	weapon, freeWeapon := alloc.New(server.Object{})
	defer freeWeapon()
	target.Damage = legacy.PortTestDamageForwardCallback(0)
	words := []int32{-1 << 31, -1, 0, 1, 1<<31 - 1}
	results := []int32{-1 << 31, -256, -1, 0, 1, 255, 256, 1<<31 - 1}
	for _, src := range []*server.Object{nil, source, target} {
		for _, item := range []*server.Object{nil, weapon, target} {
			for _, amount := range words {
				for _, kind := range words {
					for _, result := range results {
						legacy.PortTestDamageForwardReset(result)
						got := legacy.PortTestProjectileDamage(target, src, item, amount, kind)
						if got != result || byte(got) != byte(result) {
							t.Fatalf("result %d/%x want %d/%x", got, byte(got), result, byte(result))
						}
						snap := legacy.PortTestDamageForwardSnapshot()
						want := [5]uint32{uint32(uintptr(target.CObj())), uint32(uintptr(src.CObj())), uint32(uintptr(item.CObj())), uint32(amount), uint32(kind)}
						if snap.Count != 1 || snap.Callback != 1 || snap.Words != want {
							t.Fatalf("snapshot %+v want words%v", snap, want)
						}
					}
				}
			}
		}
	}
}

// Boolean registration promises only truth semantics. Integer callers must keep
// using its original C callback unless an exact-value implementation is registered.
func TestDamageValueBooleanRegistrationFallback(t *testing.T) {
	target, freeTarget := alloc.New(server.Object{})
	defer freeTarget()
	source, freeSource := alloc.New(server.Object{})
	defer freeSource()
	weapon, freeWeapon := alloc.New(server.Object{})
	defer freeWeapon()
	target.Damage = legacy.PortTestDamageForwardCallback(0)
	for _, accepted := range []bool{false, true} {
		t.Run(map[bool]string{false: "false", true: "true"}[accepted], func(t *testing.T) {
			calls := 0
			restore := server.PortTestDamageBoolRegistration(target.Damage, func(obj, src, item *server.Object, amount, kind int32) bool {
				calls++
				if obj != target || src != source || item != weapon || amount != -7 || kind != 9 {
					t.Fatal("boolean registration arguments")
				}
				return accepted
			})
			defer restore()
			legacy.PortTestDamageForwardReset(256)
			if target.CallDamage(source, weapon, -7, 9) != accepted || calls != 1 {
				t.Fatal("boolean registration ignored")
			}
			if snap := legacy.PortTestDamageForwardSnapshot(); snap.Count != 0 {
				t.Fatal("boolean registration unexpectedly called C")
			}
			if got := legacy.PortTestProjectileDamage(target, source, weapon, -7, 9); got != 256 {
				t.Fatalf("integer result %d want256", got)
			}
			if snap := legacy.PortTestDamageForwardSnapshot(); snap.Count != 1 || calls != 1 {
				t.Fatalf("raw fallback count%d boolcalls%d", snap.Count, calls)
			}
		})
	}
}
