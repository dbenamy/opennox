//go:build porttest

package opennox

import (
	"testing"
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

type damageForwardAdapter struct {
	obj    *server.Object
	order  *[]string
	name   string
	mutate func()
}

func (a damageForwardAdapter) SObj() *server.Object {
	*a.order = append(*a.order, a.name)
	if a.mutate != nil {
		a.mutate()
	}
	return a.obj
}

type damageForwardMode int

const (
	damageForwardNil damageForwardMode = iota
	damageForwardTypedNil
	damageForwardObject
	damageForwardAdapterObject
	damageForwardAdapterNil
)

func damageForwardArg(mode damageForwardMode, obj *server.Object, order *[]string, name string) server.Obj {
	switch mode {
	case damageForwardNil:
		return nil
	case damageForwardTypedNil:
		var p *server.Object
		return p
	case damageForwardObject:
		return obj
	case damageForwardAdapterObject:
		return damageForwardAdapter{obj: obj, order: order, name: name}
	case damageForwardAdapterNil:
		return damageForwardAdapter{order: order, name: name}
	default:
		panic("unknown damage forwarding mode")
	}
}

func TestDamageCallRawForwarding(t *testing.T) {
	receiver, freeReceiver := alloc.New(server.Object{})
	defer freeReceiver()
	sourceObject, freeSource := alloc.New(server.Object{})
	defer freeSource()
	weaponObject, freeWeapon := alloc.New(server.Object{})
	defer freeWeapon()
	callbackA := legacy.PortTestDamageForwardCallback(0)
	receiver.Damage = callbackA

	values := []int32{-1 << 31, -1, 0, 1, 1<<31 - 1}
	for _, sourceMode := range []damageForwardMode{damageForwardNil, damageForwardTypedNil, damageForwardObject, damageForwardAdapterObject, damageForwardAdapterNil} {
		for _, weaponMode := range []damageForwardMode{damageForwardNil, damageForwardTypedNil, damageForwardObject, damageForwardAdapterObject, damageForwardAdapterNil} {
			for _, amount := range values {
				for _, kind := range values {
					for _, result := range values {
						order := make([]string, 0, 2)
						source := damageForwardArg(sourceMode, sourceObject, &order, "source")
						weapon := damageForwardArg(weaponMode, weaponObject, &order, "weapon")
						legacy.PortTestDamageForwardReset(result)
						got := receiver.CallDamage(source, weapon, int(amount), object.DamageType(kind))
						if got != (result != 0) {
							t.Fatalf("result route source=%d weapon=%d amount=%d kind=%d result=%d got=%v", sourceMode, weaponMode, amount, kind, result, got)
						}
						snap := legacy.PortTestDamageForwardSnapshot()
						wantSource, wantWeapon := uint32(0), uint32(0)
						if sourceMode == damageForwardObject || sourceMode == damageForwardAdapterObject {
							wantSource = uint32(uintptr(unsafe.Pointer(sourceObject)))
						}
						if weaponMode == damageForwardObject || weaponMode == damageForwardAdapterObject {
							wantWeapon = uint32(uintptr(unsafe.Pointer(weaponObject)))
						}
						want := [5]uint32{uint32(uintptr(unsafe.Pointer(receiver))), wantSource, wantWeapon, uint32(amount), uint32(kind)}
						if snap.Count != 1 || snap.Callback != 1 || snap.Words != want {
							t.Fatalf("raw call source=%d weapon=%d amount=%d kind=%d result=%d snapshot=%+v want=%+v", sourceMode, weaponMode, amount, kind, result, snap, want)
						}
						var wantOrder []string
						if sourceMode == damageForwardAdapterObject || sourceMode == damageForwardAdapterNil {
							wantOrder = append(wantOrder, "source")
						}
						if weaponMode == damageForwardAdapterObject || weaponMode == damageForwardAdapterNil {
							wantOrder = append(wantOrder, "weapon")
						}
						if len(order) != len(wantOrder) {
							t.Fatalf("adapter calls source=%d weapon=%d order=%v want=%v", sourceMode, weaponMode, order, wantOrder)
						}
						for i := range order {
							if order[i] != wantOrder[i] {
								t.Fatalf("adapter order source=%d weapon=%d order=%v want=%v", sourceMode, weaponMode, order, wantOrder)
							}
						}
					}
				}
			}
		}
	}
}

func TestDamageCallNilAndTargetSelection(t *testing.T) {
	receiver, freeReceiver := alloc.New(server.Object{})
	defer freeReceiver()
	sourceObject, freeSource := alloc.New(server.Object{})
	defer freeSource()
	weaponObject, freeWeapon := alloc.New(server.Object{})
	defer freeWeapon()
	order := make([]string, 0, 2)
	source := damageForwardAdapter{obj: sourceObject, order: &order, name: "source"}
	weapon := damageForwardAdapter{obj: weaponObject, order: &order, name: "weapon"}

	receiver.Damage = nil
	legacy.PortTestDamageForwardReset(1)
	if receiver.CallDamage(source, weapon, 7, 9) {
		t.Fatal("nil callback returned true")
	}
	if snap := legacy.PortTestDamageForwardSnapshot(); snap.Count != 0 || len(order) != 0 {
		t.Fatalf("nil callback evaluated arguments or dispatched: snapshot=%+v order=%v", snap, order)
	}

	callbackA := legacy.PortTestDamageForwardCallback(0)
	callbackB := legacy.PortTestDamageForwardCallback(1)
	receiver.Damage = callbackA
	source.mutate = func() { receiver.Damage = nil }
	weapon.mutate = func() { receiver.Damage = callbackB }
	legacy.PortTestDamageForwardReset(-1)
	if !receiver.CallDamage(source, weapon, -1, object.DamageType(-1)) {
		t.Fatal("selected callback B nonzero result normalized to false")
	}
	snap := legacy.PortTestDamageForwardSnapshot()
	want := [5]uint32{uint32(uintptr(unsafe.Pointer(receiver))), uint32(uintptr(unsafe.Pointer(sourceObject))), uint32(uintptr(unsafe.Pointer(weaponObject))), 0xffffffff, 0xffffffff}
	if snap.Count != 1 || snap.Callback != 2 || snap.Words != want || receiver.Damage != callbackB {
		t.Fatalf("callback selection/order snapshot=%+v want=%+v final damage=%p", snap, want, receiver.Damage)
	}
	if len(order) != 2 || order[0] != "source" || order[1] != "weapon" {
		t.Fatalf("adapter order %v", order)
	}
}
