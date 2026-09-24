//go:build porttest

package opennox

import (
	"testing"
	"unsafe"

	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

var collisionNativeContractSlot uint32

func TestCollisionNativeResultDispatch(t *testing.T) {
	o := newCollisionCoreOwner(t)
	u := newObjectXferSimple(t, o.s)
	key := unsafe.Pointer(&collisionNativeContractSlot)
	var want, calls uint32
	var args [2]uintptr
	fn := func(obj *server.Object, a2, a3 uintptr) uint32 {
		collisionRegistryGrow(128)
		if obj != u {
			t.Fatal("native result object identity")
		}
		calls++
		args = [2]uintptr{a2, a3}
		return want
	}
	restore := server.PortTestRegisterNativeCollision("PortTestNativeCollision", key, fn, 12)
	defer restore()
	p, size := server.PortTestWorldCollisionRegistry("PortTestNativeCollision")
	if p != key || size != 12 {
		t.Fatal("native identity/data size")
	}
	u.Collide = key
	for _, value := range []uint32{0, 1, 0x7fffffff, 0x80000000, 0xffffffff} {
		want = value
		before := calls
		if got := u.CallCollideResult(uintptr(value), uintptr(^value)); got != value {
			t.Fatalf("native result %x want %x", got, value)
		}
		if calls != before+1 || args != ([2]uintptr{uintptr(value), uintptr(^value)}) {
			t.Fatal("native result args/count")
		}
		u.CallCollide(int(value), int(^value))
		if calls != before+2 || args != ([2]uintptr{uintptr(value), uintptr(^value)}) {
			t.Fatal("native void args/count")
		}
	}
	normal := &types.Pointf{X: 3, Y: -4}
	before := calls
	u.CallCollideWith(u, normal)
	if calls != before+1 || args != ([2]uintptr{uintptr(unsafe.Pointer(u)), uintptr(unsafe.Pointer(normal))}) {
		t.Fatal("native pointer args/count")
	}
	// All three aliases keep the default collision's declared zero result.
	for _, name := range []string{"DefaultCollide", "ElevatorCollide", "TelekinesisCollide"} {
		u.Collide = legacy.PortTestCollisionRegistryPointer(name)
		if u.CallCollideResult(0, 0) != 0 {
			t.Fatal("default alias result", name)
		}
	}
}
