//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"math"
	"reflect"
	"testing"
	"unsafe"
)

func TestItemRespawnList(t *testing.T) {
	o := newMatchRosterOwner(t)
	t.Cleanup(legacy.PortTestTeamRuntimeRespawns(nil))
	allow, restore := legacy.PortTestItemRespawnGlobals()
	t.Cleanup(restore)
	ids := map[uint32]uint32{0: 0}
	for i := range o.units {
		ids[uint32(uintptr(o.units[i].CObj()))] = uint32(i + 1)
		o.units[i].ObjClass = 0
	}
	type row struct {
		Name    string
		Allow   uint32
		Records [][15]uint32
	}
	var rows []row
	snapshot := func(name string, want []uint32) {
		t.Helper()
		rs := legacy.PortTestItemRespawnRecords()
		r := row{Name: name, Allow: *allow}
		if len(rs) != len(want) {
			t.Fatalf("%s length %d want%d", name, len(rs), len(want))
		}
		for i, p := range rs {
			x := *p
			id, ok := ids[x[1]]
			if !ok || id != want[i] {
				t.Fatalf("%s item%d=%d", name, i, id)
			}
			prev, next := uint32(0), uint32(0)
			if i > 0 {
				prev = uint32(uintptr(unsafe.Pointer(rs[i-1])))
			}
			if i+1 < len(rs) {
				next = uint32(uintptr(unsafe.Pointer(rs[i+1])))
			}
			if x[13] != next || x[14] != prev {
				t.Fatalf("%s links%d", name, i)
			}
			x[1] = id
			x[13] = 0
			x[14] = 0
			if next != 0 {
				x[13] = uint32(i + 2)
			}
			if prev != 0 {
				x[14] = uint32(i)
			}
			r.Records = append(r.Records, x)
		}
		rows = append(rows, r)
	}
	legacy.PortTestItemRespawn("remove", &o.units[0])
	snapshot("empty-remove", nil)
	*allow = 0
	if legacy.PortTestItemRespawn("add", &o.units[0]) != 0 {
		t.Fatal("disabled return")
	}
	snapshot("disabled", nil)
	legacy.PortTestItemRespawn("reset", nil)
	for i := range o.units {
		u := &o.units[i]
		u.TypeInd = uint16(410 + i)
		u.PosVec = types.Pointf{X: float32(i) + 1.25, Y: -float32(i) - 2.5}
		u.Direction1 = server.Dir16(255 + i)
		prior := legacy.PortTestItemRespawnRecords()
		var want uintptr
		if len(prior) > 0 {
			want = uintptr(unsafe.Pointer(prior[0]))
		}
		if got := legacy.PortTestItemRespawn("add", u); got != want {
			t.Fatalf("insertion return%d", i)
		}
		x := *legacy.PortTestItemRespawnRecords()[0]
		if x[0] != uint32(u.TypeInd) || x[2] != math.Float32bits(u.PosVec.X) || x[3] != math.Float32bits(u.PosVec.Y) || x[4] != uint32(uint16(u.Direction1)) || x[5] != 0 || x[6] != 0 || !reflect.DeepEqual(x[7:13], make([]uint32, 6)) {
			t.Fatal("record fields", x)
		}
	}
	snapshot("three", []uint32{3, 2, 1})
	legacy.PortTestItemRespawn("remove", nil)
	snapshot("missing", []uint32{3, 2, 1})
	legacy.PortTestItemRespawn("remove", &o.units[1])
	snapshot("middle", []uint32{3, 1})
	legacy.PortTestItemRespawn("remove", &o.units[0])
	snapshot("tail", []uint32{3})
	legacy.PortTestItemRespawn("remove", &o.units[2])
	snapshot("head", nil)
	legacy.PortTestItemRespawn("remove", &o.units[2])
	snapshot("repeat-empty", nil)
	for i := 0; i < 2; i++ {
		legacy.PortTestItemRespawn("add", &o.units[0])
	}
	snapshot("duplicates", []uint32{1, 1})
	legacy.PortTestItemRespawn("remove", &o.units[0])
	snapshot("remove-first-duplicate", []uint32{1})
	*allow = 0x80000000
	legacy.PortTestItemRespawn("reset", nil)
	snapshot("reset", nil)
	if *allow != 1 {
		t.Fatal("reset enable")
	}
	// Use C-owned definition data to verify all class-mask branches and ammunition ordering.
	t.Cleanup(o.s.PortTestRewardTypes([]string{"PortTestRewardWeapon"}, nil, true, 0, 0x82))
	for _, class := range []uint32{0, 0x1000, 0x1000000, 0x2000000, 0x10000000, 0x13001000, 0x4000} {
		u := &o.units[0]
		u.TypeInd = uint16(o.s.Types.IndByID("PortTestRewardWeapon"))
		u.ObjClass = object.Class(class)
		init := o.record(t, 20)
		for j := 0; j < 5; j++ {
			*(*uint32)(unsafe.Add(init, j*4)) = uint32(0x123400 + j)
		}
		oldInit, oldUse := u.InitData, u.UseData
		u.InitData = init
		u.UseData.Ptr = o.record(t, 64)
		*(*byte)(u.UseData.Ptr) = 17
		*(*byte)(unsafe.Add(u.UseData.Ptr, 1)) = 29
		legacy.PortTestItemRespawn("add", u)
		x := legacy.PortTestItemRespawnRecords()[0]
		for j := 0; j < 5; j++ {
			want := uint32(0)
			if class&0x13001000 != 0 {
				want = uint32(0x123400 + j)
			}
			if x[7+j] != want {
				t.Fatal("modifier copy", class, j)
			}
		}
		want := uint32(0)
		if class&0x1000000 != 0 {
			want = 17<<8 | 29
		}
		if x[12] != want {
			t.Fatal("ammo order", class, x[12])
		}
		snapshot(fmt.Sprintf("class%x", class), []uint32{1})
		legacy.PortTestItemRespawn("reset", nil)
		u.InitData, u.UseData = oldInit, oldUse
	}
	spellbookCapture(t, "item-respawn-list", rows, "3db2efd780abb907035c60837b212ec4250505bb6e1f1ed56278a76bae67b289")
}

func TestItemRespawnPoolCapacity(t *testing.T) {
	o := newMatchRosterOwner(t)
	t.Cleanup(legacy.PortTestTeamRuntimeRespawns(nil))
	u := &o.units[0]
	u.ObjClass = 0
	type row struct {
		Attempt, Count   int
		ReturnedPrevious bool
	}
	var rows []row
	for i := 0; i < 386; i++ {
		before := legacy.PortTestItemRespawnRecords()
		var previous uintptr
		if len(before) != 0 {
			previous = uintptr(unsafe.Pointer(before[0]))
		}
		got := legacy.PortTestItemRespawn("add", u)
		after := legacy.PortTestItemRespawnRecords()
		want := i + 1
		if want > 384 {
			want = 384
		}
		if len(after) != want || i < 384 && got != previous || i >= 384 && got != 0 {
			t.Fatal("capacity", i, len(after), got)
		}
		rows = append(rows, row{i, len(after), got == previous})
	}
	legacy.PortTestItemRespawn("remove", u)
	if len(legacy.PortTestItemRespawnRecords()) != 383 {
		t.Fatal("remove at capacity")
	}
	legacy.PortTestItemRespawn("add", u)
	if len(legacy.PortTestItemRespawnRecords()) != 384 {
		t.Fatal("reuse freed record")
	}
	legacy.PortTestItemRespawn("reset", nil)
	if len(legacy.PortTestItemRespawnRecords()) != 0 {
		t.Fatal("pool reset")
	}
	legacy.PortTestItemRespawn("add", u)
	if len(legacy.PortTestItemRespawnRecords()) != 1 {
		t.Fatal("reuse reset pool")
	}
	spellbookCapture(t, "item-respawn-pool-capacity", rows, "ed53d98969f543ca3ebee9cf431b0b3eecaea97659534239421ab6e0fad4073c")
}
