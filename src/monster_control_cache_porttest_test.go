//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestMonsterControlScriptCache(t *testing.T) {
	o := newCollisionCoreOwner(t)
	owned := serverConfigOwnBytes(t, 0x5D4594, 2386620, 208)
	clear(owned)
	init, restore := legacy.PortTestMonsterCacheInit()
	t.Cleanup(restore)
	units := make([]*server.Object, 20)
	ids := map[unsafe.Pointer]uint32{}
	for i := range units {
		u := newObjectXferSimple(t, o.s)
		saved := *u
		t.Cleanup(func() { *u = saved })
		u.ScriptIDVal = 100 + i
		u.ObjFlags = 4
		units[i] = u
		ids[u.CObj()] = uint32(1001 + i)
	}
	base := uint32(uintptr(memmap.PtrOff(0x5D4594, 2386628)))
	normalize := func(v uint32) uint32 {
		if v == 0 || v == 1 {
			return v
		}
		if id, ok := ids[unsafe.Pointer(uintptr(v))]; ok {
			return id
		}
		if v >= base && v < base+192 && (v-base)%12 == 0 {
			return 4000 + (v-base)/12
		}
		t.Fatalf("cache pointer outside owner: %x", v)
		return 0
	}
	type row struct {
		Op           string
		Unit         int
		ID           int32
		Return, Init uint32
		Words        [52]uint32
		Active       []uint32
	}
	var rows []row
	call := func(op string, unit int, id int32) uint32 {
		var u *server.Object
		if unit >= 0 {
			u = units[unit]
		}
		rv := legacy.PortTestMonsterCache(op, u, id)
		r := row{Op: op, Unit: unit, ID: id, Return: normalize(rv), Init: *init}
		words := unsafe.Slice((*uint32)(unsafe.Pointer(&owned[0])), 52)
		for i, v := range words {
			r.Words[i] = normalize(v)
		}
		seen := map[uint32]bool{}
		last := uint32(0)
		node := words[50]
		for node != 0 {
			if seen[node] {
				t.Fatal("cycle in active script cache")
			}
			seen[node] = true
			n := normalize(node)
			if n < 4000 || n >= 4016 {
				t.Fatal("non-node in active cache")
			}
			slot := 2 + 3*int(n-4000)
			if words[slot+1] != last {
				t.Fatal("active cache back link mismatch")
			}
			r.Active = append(r.Active, normalize(words[slot]))
			last = node
			node = words[slot+2]
		}
		if last != words[51] {
			t.Fatal("cache tail mismatch")
		}
		for node = words[0]; node != 0; {
			if seen[node] {
				t.Fatal("duplicate or cycle in free cache")
			}
			seen[node] = true
			n := normalize(node)
			if n < 4000 || n >= 4016 {
				t.Fatal("non-node in free cache")
			}
			node = words[2+3*int(n-4000)+2]
		}
		if *init == 0 && len(seen) != 16 {
			t.Fatal("cache lost a node", len(seen))
		}
		rows = append(rows, r)
		return r.Return
	}
	if call("find", -1, 100) != 0 || *init != 0 {
		t.Fatal("cold lookup should initialize empty cache")
	}
	for i := 0; i < 16; i++ {
		call("prepare", i, 0)
	}
	if call("find", -1, 100) != 1001 {
		t.Fatal("oldest live unit missing")
	}
	call("prepare", 16, 0)
	if call("find", -1, 101) != 0 {
		t.Fatal("least recently used unit not evicted")
	}
	if call("find", -1, 100) != 1001 {
		t.Fatal("recently looked-up unit was evicted")
	}
	units[0].ObjFlags |= 0x20
	if call("find", -1, 100) != 0 {
		t.Fatal("destroyed cached unit returned")
	}
	units[0].ObjFlags &^= 0x20
	for _, i := range []int{0, 8, 16, 19, 8} {
		call("remove", i, 0)
	}
	call("prepare", 17, 0)
	call("prepare", 17, 0)
	call("remove", 17, 0)
	if call("find", -1, 117) != 1018 {
		t.Fatal("duplicate cache insertion contract")
	}
	call("clear", -1, 0)
	call("clear", -1, 0)
	for i := 0; i < 20; i++ {
		call("prepare", i, 0)
	}
	call("reset", -1, 0)
	if call("find", -1, 119) != 0 {
		t.Fatal("reset preserved active entry")
	}
	spellbookCapture(t, "monster-control-script-cache", rows, "8827cd1fd10e6e5448945d8e1eda8fe6f02ef17d756f7804be01a8f1a399eab1")
}
