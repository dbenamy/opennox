//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestMonsterControlPendingOwnership(t *testing.T) {
	o := newCollisionCoreOwner(t)
	head, restore := legacy.PortTestMonsterPendingOwner()
	t.Cleanup(restore)
	oldMain, oldPending, oldMissiles := o.s.Objs.List, o.s.Objs.Pending, o.s.Objs.MissileList
	t.Cleanup(func() { o.s.Objs.List = oldMain; o.s.Objs.Pending = oldPending; o.s.Objs.MissileList = oldMissiles })
	units := make([]*server.Object, 6)
	saved := make([]server.Object, 6)
	for i := range units {
		units[i] = newObjectXferSimple(t, o.s)
		saved[i] = *units[i]
	}
	t.Cleanup(func() {
		for i, u := range units {
			*u = saved[i]
		}
	})
	ids := collisionCoreIDs(units...)
	type row struct {
		Op      string
		Mode    int
		Return  uint32
		Pairs   [][2]uint32
		Objects [][3]uint32
	}
	var rows []row
	record := func(op string, mode int, ret uint32) {
		r := row{Op: op, Mode: mode, Return: ret}
		seen := map[uint32]bool{}
		for p := *head; p != 0; {
			if seen[p] {
				t.Fatal("pending owner cycle")
			}
			seen[p] = true
			node := (*[3]uint32)(unsafe.Pointer(uintptr(p)))
			r.Pairs = append(r.Pairs, [2]uint32{node[0], node[1]})
			p = node[2]
		}
		for _, u := range units {
			r.Objects = append(r.Objects, [3]uint32{collisionCoreID(t, ids, uint32(uintptr(u.ObjOwner.CObj()))), collisionCoreID(t, ids, uint32(uintptr(u.Field128.CObj()))), collisionCoreID(t, ids, uint32(uintptr(u.Field129.CObj())))})
		}
		rows = append(rows, r)
	}
	for mode := 0; mode < 5; mode++ {
		for i, u := range units {
			*u = saved[i]
			u.ScriptIDVal = 100 + i
			u.ObjFlags = 4
			u.ObjOwner = nil
			u.Field128 = nil
			u.Field129 = nil
			u.InvFirstItem = nil
			u.InvNextItem = nil
			u.ObjNext = nil
		}
		o.s.Objs.List = units[0]
		units[0].ObjNext = units[1]
		units[1].InvFirstItem = units[2]
		o.s.Objs.Pending = units[3]
		units[3].ObjNext = units[4]
		o.s.Objs.MissileList = units[5]
		if mode == 1 {
			units[1].ObjFlags |= 0x20
		}
		if mode == 2 {
			units[4].ObjFlags |= 0x20
		}
		if mode == 3 {
			units[1].ScriptIDVal = 100
		}
		rv := legacy.PortTestMonsterPending("init", 0, 0)
		if rv != 1 {
			t.Fatal("pending owner pool init")
		}
		record("init", mode, rv)
		pairs := [][2]int32{{100, 101}, {101, 102}, {100, 103}, {100, 104}, {100, 105}, {999, 104}, {100, 999}, {0, 0}}
		for _, p := range pairs {
			rv = legacy.PortTestMonsterPending("add", p[0], p[1])
			if rv == 0 || rv != *head {
				t.Fatal("pending owner prepend")
			}
			record("add", mode, 1)
		}
		if mode == 4 {
			legacy.PortTestMonsterPending("clear", 0, 0)
			record("clear", mode, 0)
		}
		legacy.PortTestMonsterPending("resolve", 0, 0)
		if *head != 0 {
			t.Fatal("resolve did not clear list")
		}
		record("resolve", mode, 0)
		if mode == 0 && (units[1].ObjOwner != units[0] || units[2].ObjOwner != units[1] || units[5].ObjOwner != units[0]) {
			t.Fatal("lookup across main/inventory/pending/missile lists")
		}
		if mode == 4 {
			for _, u := range units {
				if u.ObjOwner != nil {
					t.Fatal("cleared records resolved")
				}
			}
		}
		legacy.PortTestMonsterPending("resolve", 0, 0)
		record("resolve-empty", mode, 0)
		// Exercise pool capacity and reuse, without depending on allocation addresses.
		for i := 0; i < 520; i++ {
			rv = legacy.PortTestMonsterPending("add", int32(i), int32(i+1))
			if rv != 0 {
				rv = 1
			}
			record("capacity", mode, rv)
		}
		legacy.PortTestMonsterPending("clear", 0, 0)
		record("clear-capacity", mode, 0)
		rv = legacy.PortTestMonsterPending("add", 100, 105)
		if rv == 0 {
			t.Fatal("pool failed reuse")
		}
		record("reuse", mode, 1)
		rv = legacy.PortTestMonsterPending("free", 0, 0)
		if rv != 0 || *head != 0 {
			t.Fatal("pool free")
		}
		record("free", mode, rv)
	}
	spellbookCapture(t, "monster-control-pending", rows, "58ceed5d2b22291c749d37f77223c35649dbf0e3a0ac4d8ffa16e52eea4b7d04")
}
