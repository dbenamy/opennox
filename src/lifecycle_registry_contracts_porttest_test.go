//go:build porttest

package opennox

import (
	"bytes"
	"fmt"
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestLifecycleRegistryNames(t *testing.T) {
	for _, sp := range []struct {
		create bool
		count  int
	}{{true, 8}, {false, 12}} {
		names := legacy.PortTestLifecycleRegistryNames(sp.create)
		if len(names) != sp.count {
			t.Fatalf("registrations %d want %d", len(names), sp.count)
		}
		seen := map[string]bool{}
		for _, name := range names {
			if seen[name] {
				t.Fatal("duplicate name", name)
			}
			seen[name] = true
		}
	}
	a, an := server.PortTestLifecycleRegistry("MonsterInit", false)
	b, bn := server.PortTestLifecycleRegistry("ShopkeeperInit", false)
	if a == nil || a != b || an != 0 || bn != unsafe.Sizeof(server.ShopkeeperInitData{}) {
		t.Fatal("monster/shopkeeper alias and sizes")
	}
}

// Call the unchanged real pending-object owner, with a single isolated object.
// A non-active object deliberately avoids unrelated spatial-index registration;
// callback selection, arguments, list insertion and pending-flag clearing are real.
func lifecyclePending(t *testing.T, u *server.Object) {
	t.Helper()
	s := u.Server()
	if s.Objs.Pending != nil || s.Objs.List != nil {
		t.Fatal("pending fixture requires empty lists")
	}
	u.ObjFlags |= object.FlagPending
	s.Objs.Pending = u
	noxServer.ObjectsAddPending()
	if s.Objs.Pending != nil || s.Objs.List != u || u.ObjFlags.Has(object.FlagPending) {
		t.Fatal("pending ownership transition")
	}
	s.Objs.List = nil
	u.ObjNext = nil
	u.ObjPrev = nil
}
func TestLifecycleRegistryInitRaw(t *testing.T) {
	s := newObjectXferOwner(t)
	u := newObjectXferSimple(t, s)
	cb := legacy.PortTestDeathForwardReset()
	u.CallInit()
	if p, n := legacy.PortTestDeathForwardSnapshot(); p != nil || n != 0 {
		t.Fatal("nil initializer called")
	}
	u.Init = cb
	u.CallInit()
	if p, n := legacy.PortTestDeathForwardSnapshot(); p != u.CObj() || n != 1 {
		t.Fatal("single-pointer raw initializer forwarding")
	}
	u.Init = nil
	u.CallInit()
	lifecyclePending(t, u)
	if _, n := legacy.PortTestDeathForwardSnapshot(); n != 1 {
		t.Fatal("cleared initializer called")
	}
	raw, events := legacy.PortTestSessionEntryInitCallback()
	u.Init = raw
	lifecyclePending(t, u)
	got := events()
	if len(got) != 1 || got[0] != [2]uintptr{uintptr(u.CObj()), 0} {
		t.Fatalf("pending raw two-argument forwarding: %v", got)
	}
}
func TestLifecycleRegistryInitDynamic(t *testing.T) {
	old := legacy.Nox_xxx_unitMonsterInit_4F0040
	defer func() { legacy.Nox_xxx_unitMonsterInit_4F0040 = old }()
	for _, name := range []string{"MonsterInit", "ShopkeeperInit"} {
		t.Run(name, func(t *testing.T) {
			s := newObjectXferOwner(t)
			u := newObjectXferSimple(t, s)
			u.Init = legacy.PortTestLifecycleRegistryPointer(name, false)
			calls := 0
			for i := uint32(1); i <= 2; i++ {
				want := i
				legacy.Nox_xxx_unitMonsterInit_4F0040 = func(got *server.Object) {
					if got != u {
						t.Fatal("initializer object identity")
					}
					calls++
					got.Worth = want
				}
				u.CallInit()
				if u.Worth != want || calls != int(i*2-1) {
					t.Fatal("single-argument replacement")
				}
				u.Worth = 0
				lifecyclePending(t, u)
				if u.Worth != want || calls != int(i*2) {
					t.Fatal("pending replacement")
				}
			}
		})
	}
}
func TestLifecycleRegistryPendingInit(t *testing.T) {
	stage := memmap.PtrUint32(0x5D4594, 2388660)
	oldStage := *stage
	*stage = 0
	t.Cleanup(func() { *stage = oldStage })
	type row struct {
		Name                  string
		Index                 int
		Direction, Direction2 server.Dir16
		Position, Previous    types.Pointf
		Flags                 object.Flags
		Status                object.SubClass
		Init, Update          []byte
	}
	var rows []row
	for _, name := range []string{"SparkInit", "FrogInit", "ChestInit", "BoulderInit", "SkullInit", "DirectionInit", "GoldInit", "BreakInit", "MonsterGeneratorInit"} {
		for i := 0; i < 4; i++ {
			t.Run(fmt.Sprintf("%s/%d", name, i), func(t *testing.T) {
				s := newObjectXferOwner(t)
				u := newObjectXferSimple(t, s)
				init, freeInit := alloc.Make([]byte{}, 2200)
				defer freeInit()
				update, freeUpdate := alloc.Make([]byte{}, 2200)
				defer freeUpdate()
				u.InitData = unsafe.Pointer(&init[0])
				u.UpdateData = unsafe.Pointer(&update[0])
				defer func() { u.InitData = nil; u.UpdateData = nil }()
				u.Init = legacy.PortTestLifecycleRegistryPointer(name, false)
				*(*uint32)(u.InitData) = uint32(1 + i)
				*(*uint32)(unsafe.Add(u.InitData, 4)) = uint32(i - 2)
				u.PosVec = types.Pointf{X: float32(10 + i), Y: float32(20 - i)}
				u.Field5 = uint32(i * 2)
				// Unknown generator level preserves its prior cap without balance dependencies.
				update[83] = 4
				update[87] = 77
				lifecyclePending(t, u)
				switch name {
				case "SparkInit":
					if objectXferGetWord(u.UpdateData, 0) != 32 || objectXferGetWord(u.UpdateData, 4) != 32 {
						t.Fatal("spark extents")
					}
				case "FrogInit":
					if update[0] < 55 || update[0] > 60 || update[1] != 1 || update[2] != 0 || u.Direction2 > 255 {
						t.Fatal("frog initialization")
					}
				case "BoulderInit":
					if u.Pos39 != u.PosVec {
						t.Fatal("boulder previous position")
					}
				case "GoldInit":
					if objectXferGetWord(u.InitData, 0) != uint32(1+i) {
						t.Fatal("explicit gold changed")
					}
				case "SkullInit", "DirectionInit":
					if u.Direction1 != u.Direction2 {
						t.Fatal("direction mismatch")
					}
				case "MonsterGeneratorInit":
					if update[87] != 77 || u.Direction1 != u.Direction2 {
						t.Fatal("generator preserved cap/direction")
					}
				}
				rows = append(rows, row{Name: name, Index: i, Direction: u.Direction1, Direction2: u.Direction2, Position: u.PosVec, Previous: u.Pos39, Flags: u.ObjFlags, Init: bytes.Clone(init), Update: bytes.Clone(update)})
			})
		}
	}
	spellbookCapture(t, "lifecycle-registry-pending-init", rows, "307b29ea5a06fde046f164199d154e6b77eb9dc206779b1324aef6e569265020")
}

func TestLifecycleRegistryInitArguments(t *testing.T) {
	u, freeU := alloc.New(server.Object{})
	defer freeU()
	other, freeOther := alloc.New(server.Object{})
	defer freeOther()
	for _, arg := range []unsafe.Pointer{nil, u.CObj(), other.CObj()} {
		raw, events := legacy.PortTestSessionEntryInitCallback()
		u.Init = raw
		server.PortTestLifecycleInitWithArg(u, arg)
		got := events()
		if len(got) != 1 || got[0] != [2]uintptr{uintptr(u.CObj()), uintptr(arg)} {
			t.Fatalf("raw initializer arguments: %v", got)
		}
	}
}
