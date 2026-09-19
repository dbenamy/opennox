//go:build porttest

package opennox

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func bookAwardCall(name string, a ...uint32) uint32 { return legacy.PortTestBookAwardCall(name, a...) }
func bookAwardString(t *testing.T, s string) uint32 {
	t.Helper()
	p, free := alloc.CString(s)
	t.Cleanup(free)
	return uint32(uintptr(unsafe.Pointer(p)))
}
func bookAwardWords(t *testing.T, base, off uintptr, n int) []uint32 {
	t.Helper()
	p := unsafe.Slice(memmap.PtrUint32(base, off), n)
	old := append([]uint32(nil), p...)
	t.Cleanup(func() { copy(p, old) })
	clear(p)
	return p
}
func TestBookAwardsNameLookup(t *testing.T) {
	guides := bookAwardWords(t, 0x587000, 70500, 41)
	abilities := bookAwardWords(t, 0x587000, 69736, 7)
	names := make([]string, 41)
	for i := range guides {
		names[i] = fmt.Sprintf("Guide%02d", i)
	}
	names[0] = ""
	names[7] = "Café Ω"
	names[9] = names[7]
	names[40] = strings.Repeat("Q", 255)
	for i, s := range names {
		guides[i] = bookAwardString(t, s)
	}
	anames := []string{"", "AbilityOne", "DUP", "DUP", "é", "AbilityFive"}
	for i, s := range anames {
		abilities[i] = bookAwardString(t, s)
	}
	var rows []map[string]any
	queries := append(append([]string{}, names...), "guide01", "Café ω", "absent", strings.Repeat("Q", 254))
	for _, q := range queries {
		want := uint32(0)
		for i, s := range names {
			if q == s {
				want = uint32(i)
				break
			}
		}
		got := bookAwardCall("nox_xxx_guide_427010", bookAwardString(t, q))
		if got != want {
			t.Fatal("guide lookup", q, got, want)
		}
		rows = append(rows, map[string]any{"kind": "guide", "query": q, "result": got})
	}
	for i := range guides {
		if got := bookAwardCall("nox_xxx_guideNameByN_427230", uint32(i)); got != guides[i] {
			t.Fatal("guide name pointer", i)
		}
	}
	for _, end := range []int{0, 1, 3, 6} {
		old := abilities[end]
		abilities[end] = 0
		for _, q := range append(append([]string{}, anames...), "abilityone", "absent") {
			want := uint32(0)
			for i, s := range anames[:end] {
				if s == q {
					want = uint32(i)
					break
				}
			}
			got := bookAwardCall("nox_xxx_abilityNameToN_424D80", bookAwardString(t, q))
			if got != want {
				t.Fatal("ability lookup", end, q, got, want)
			}
			rows = append(rows, map[string]any{"kind": "ability", "end": end, "query": q, "result": got})
		}
		abilities[end] = old
	}
	spellbookCapture(t, "book-awards-name-lookup", rows, "b4ea2fb4c67750f8f81a0b22dbf94f06b42e7f824db604e61af545797e81847e")
}
func TestBookAwardsGuideCatalog(t *testing.T) {
	table := bookAwardWords(t, 0x5D4594, 740076, 41*7)
	var rows []map[string]any
	for _, pattern := range []string{"none", "first", "last", "sparse", "all"} {
		clear(table)
		enabled := make([]bool, 41)
		// Index zero is populated to distinguish first-entry iteration from next(-1).
		for i := 0; i <= 40; i++ {
			enabled[i] = i == 0 || pattern == "all" || pattern == "first" && i == 1 || pattern == "last" && i == 40 || pattern == "sparse" && (i == 2 || i == 19 || i == 40)
			for j := 0; j < 7; j++ {
				table[7*i+j] = uint32(0x11000000 + i*256 + j*17)
			}
			table[7*i+1] = 0
			if enabled[i] {
				table[7*i+1] = uint32(100 + i%7)
			}
		}
		before := append([]uint32(nil), table...)
		next := func(start int) uint32 {
			for i := start; i <= 40; i++ {
				if enabled[i] {
					return uint32(i)
				}
			}
			return 0
		}
		first := bookAwardCall("nox_xxx_bookGetFirstCreMB_427300")
		if first != next(1) {
			t.Fatal("first", pattern, first)
		}
		for i := -1; i <= 41; i++ {
			got := bookAwardCall("nox_xxx_bookGetNextCre_427320", uint32(i))
			want := uint32(0)
			if i+1 < 41 {
				want = next(i + 1)
			}
			if got != want {
				t.Fatal("next", pattern, i, got, want)
			}
			rows = append(rows, map[string]any{"pattern": pattern, "op": "next", "id": i, "result": got})
		}
		for _, typ := range []uint32{0, 99, 100, 101, 102, 103, 104, 105, 106, 107, 0xffffffff} {
			want := uint32(0)
			for i := 1; i <= 40; i++ {
				if table[7*i+1] != 0 && table[7*i+1] == typ {
					want = uint32(i)
					break
				}
			}
			got := bookAwardCall("nox_xxx_creatureIsCharmableByTT_4272B0", typ)
			if got != want {
				t.Fatal("charm", pattern, typ, got, want)
			}
			rows = append(rows, map[string]any{"pattern": pattern, "op": "charm", "type": typ, "result": got})
		}
		for i := -1; i <= 41; i++ {
			for _, spec := range []struct {
				name    string
				word    int
				enabled bool
			}{{"nox_xxx_guiCreatureGetName_427240", 0, true}, {"nox_xxx_bookGetCreatureImg_427400", 4, false}, {"sub_427430", 3, false}} {
				want := uint32(0)
				if i > 0 && i < 41 && (!spec.enabled || enabled[i]) {
					want = table[7*i+spec.word]
				}
				got := bookAwardCall(spec.name, uint32(i))
				if got != want {
					t.Fatal(spec.name, pattern, i, got, want)
				}
				rows = append(rows, map[string]any{"pattern": pattern, "op": spec.name, "id": i, "result": got})
			}
		}
		for i := 0; i <= 40; i++ {
			desc := bookAwardCall("nox_xxx_guideGetDescById_4272E0", uint32(i))
			size := bookAwardCall("nox_xxx_guideGetUnitSize_427460", uint32(i))
			if desc != table[7*i+2] || size != uint32(byte(table[7*i+6])) {
				t.Fatal("raw fields", i, desc, size)
			}
			rows = append(rows, map[string]any{"pattern": pattern, "op": "fields", "id": i, "desc": desc, "size": size})
		}
		if !bytes.Equal(unsafe.Slice((*byte)(unsafe.Pointer(&table[0])), len(table)*4), unsafe.Slice((*byte)(unsafe.Pointer(&before[0])), len(before)*4)) {
			t.Fatal("catalog mutated")
		}
	}
	spellbookCapture(t, "book-awards-guide-catalog", rows, "9aa02736c9568b5d391ba099e358969c96d8aa7a6e6e40f7070c8892b83ac576")
}
func TestBookAwardsEnchantEnumeration(t *testing.T) {
	u, free := alloc.New(server.Object{})
	defer free()
	var rows []map[string]any
	for _, ids := range [][]uint32{nil, {0}, {26, 0, 28, 5}, {5, 5, 0, 5}, {28, 27, 26, 25, 24, 23, 22, 21, 20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0}} {
		restore := legacy.PortTestPlayerFileEnchants(ids)
		t.Cleanup(restore)
		for _, mask := range []uint32{0, 1, 1 << 26, 1 << 28, 0x15555555, 0xffffffff} {
			u.Buffs = mask
			want := uint32(0)
			for _, id := range ids {
				if mask&(1<<id) != 0 {
					want++
				}
			}
			got := bookAwardCall("sub_424CB0", uint32(uintptr(unsafe.Pointer(u))))
			if got != want {
				t.Fatal("count", ids, mask, got, want)
			}
			rows = append(rows, map[string]any{"ids": ids, "mask": mask, "count": got})
		}
		first := uint32(0xffffffff)
		if len(ids) > 0 {
			first = ids[0]
		}
		if got := bookAwardCall("sub_424D00"); got != first {
			t.Fatal("first enchant", got, first)
		}
		for q := -1; q <= 30; q++ {
			want := uint32(0xffffffff)
			for i := 0; i+1 < len(ids); i++ {
				if ids[i] == uint32(q) {
					want = ids[i+1]
					break
				}
			}
			got := bookAwardCall("sub_424D20", uint32(q))
			if got != want {
				t.Fatal("next enchant", ids, q, got, want)
			}
			rows = append(rows, map[string]any{"ids": ids, "query": q, "next": got})
		}
	}
	spellbookCapture(t, "book-awards-enchant-enumeration", rows, "147ba5b3916727d0829b3e2e8fc18136af8acc708c38d7522e5097b33abfe582")
}

func TestBookAwardsEnchantCountGates(t *testing.T) {
	count := legacy.PortTestBookEnchantCount()
	old := *count
	t.Cleanup(func() { *count = old })
	var rows []map[string]any
	for _, n := range []int32{0, -1, -2147483648} {
		*count = n
		a := bookAwardCall("sub_424CB0", 0)
		b := bookAwardCall("sub_424D00")
		c := bookAwardCall("sub_424D20", 27)
		if a != 0 || b != 0xffffffff || c != 0xffffffff || *count != n {
			t.Fatal("signed count gate", n, a, b, c, *count)
		}
		rows = append(rows, map[string]any{"count": n, "matches": a, "first": b, "next": c})
	}
	spellbookCapture(t, "book-awards-enchant-count-gates", rows, "34696b525798a612cf158b3af3bb1ad4ea4c6b5b883da56c601373d2eab49ac3")
}
