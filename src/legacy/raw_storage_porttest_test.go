//go:build porttest

package legacy

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"os"
	"reflect"
	"runtime"
	"sort"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
)

func TestRawStorageOwners(t *testing.T) {
	if len(rawStorageContracts) != 52 {
		t.Fatalf("owner count %d", len(rawStorageContracts))
	}
	layout := []uintptr{148, 84, 0, 4, 8, 128, 132, 136, 140, 141, 142, 143, 144, 12, 0, 4, 8}
	if got := rawStorageLayout(); !reflect.DeepEqual(got, layout) {
		t.Fatalf("typed layout %v want %v", got, layout)
	}
	type region struct {
		lo, hi uintptr
		name   string
	}
	var regions []region
	originals := make([][]byte, len(rawStorageContracts))
	index := make(map[string]int)
	for i, c := range rawStorageContracts {
		if c.size != c.wantSize || c.ptr == nil || uintptr(c.ptr)%c.align != 0 {
			t.Fatalf("%s size/alignment", c.name)
		}
		if !rawStorageInitialZero[i] {
			t.Fatalf("%s not statically zero", c.name)
		}
		if c.alias() != c.ptr {
			t.Fatalf("%s typed owner alias", c.name)
		}
		regions = append(regions, region{uintptr(c.ptr), uintptr(c.ptr) + c.size, c.name})
		originals[i] = append([]byte(nil), unsafe.Slice((*byte)(c.ptr), c.size)...)
		index[c.name] = i
		if c.blobBase != 0 {
			b := memmap.BlobByAddr(c.blobBase)
			if b == nil || b.Addr != c.blobBase || b.Name != c.name || b.Size != c.size || len(b.Data) != int(c.size) || unsafe.Pointer(&b.Data[0]) != c.ptr {
				t.Fatalf("%s registry", c.name)
			}
			if len(b.Initial) != int(c.size) || !bytes.Equal(b.Initial, make([]byte, c.size)) {
				t.Fatalf("%s initial registry bytes", c.name)
			}
			offsets := []uintptr{0, 1, c.size / 2, c.size - 1}
			// All registered fields still within a blob must preserve their exact offsets.
			// Extracted owners outside that blob's physical extent are checked separately.
			for _, v := range memmap.Variables() {
				if v.Addr >= c.blobBase && v.Addr-c.blobBase < c.size && v.Size > 0 && v.Size <= c.size-(v.Addr-c.blobBase) {
					offsets = append(offsets, v.Addr-c.blobBase, v.Addr-c.blobBase+v.Size-1)
				}
			}
			for _, off := range offsets {
				p := unsafe.Add(c.ptr, off)
				if memmap.PtrOff(c.blobBase, off) != p || len(memmap.Slice(c.blobBase, off)) != int(c.size-off) {
					t.Fatalf("%s offset %d", c.name, off)
				}
				got, rel := memmap.BlobByPtr(p)
				if got != b || rel != off {
					t.Fatalf("%s reverse offset %d", c.name, off)
				}
			}
		}
	}
	// Include the already-converted numeric owners in the disjointness check.
	for _, c := range scalarStorageContracts {
		regions = append(regions, region{uintptr(c.ptr), uintptr(c.ptr) + c.size, c.name})
	}
	sort.Slice(regions, func(i, j int) bool { return regions[i].lo < regions[j].lo })
	for i := 1; i < len(regions); i++ {
		if regions[i].lo < regions[i-1].hi {
			t.Fatalf("overlap %s/%s", regions[i-1].name, regions[i].name)
		}
	}
	restore := func() {
		for i, c := range rawStorageContracts {
			copy(unsafe.Slice((*byte)(c.ptr), c.size), originals[i])
		}
	}
	defer restore()
	if grid := uiInventoryGrid(); len(grid) != 84 || unsafe.Pointer(&grid[0]) != rawStorageContracts[index["nox_client_inventory_grid_1050020"]].ptr {
		t.Fatal("production inventory grid alias")
	}
	if equipment := uiInventoryEquipment(); len(equipment) != 9 || unsafe.Pointer(&equipment[0]) != rawStorageContracts[index["array_5D4594_1049872"]].ptr {
		t.Fatal("production equipment alias")
	}
	if Get_nox_common_maplist() != rawStorageContracts[index["nox_common_maplist"]].ptr {
		t.Fatal("production list alias")
	}
	if unsafe.Pointer(scoreboardData.parent) != rawStorageContracts[index["dword_5d4594_1090048"]].ptr || unsafe.Pointer(scoreboardData.rank) != rawStorageContracts[index["dword_5d4594_1090100"]].ptr {
		t.Fatal("production scoreboard aliases")
	}
	neighbors := func(except int) {
		t.Helper()
		for i, c := range rawStorageContracts {
			if i != except && !bytes.Equal(unsafe.Slice((*byte)(c.ptr), c.size), originals[i]) {
				t.Fatalf("write to %s changed %s", rawStorageContracts[except].name, c.name)
			}
		}
	}
	type observation struct {
		Name        string
		Size        uintptr
		InitialZero bool
		Hashes      []string
	}
	observations := make([]observation, 0, len(rawStorageContracts))
	for i, c := range rawStorageContracts {
		data := unsafe.Slice((*byte)(c.ptr), c.size)
		alias := unsafe.Slice((*byte)(c.alias()), c.size)
		if c.blobBase != 0 {
			alias = memmap.BlobByAddr(c.blobBase).Data
		}
		o := observation{Name: c.name, Size: c.size, InitialZero: rawStorageInitialZero[i]}
		expected := make([]byte, c.size)
		// Six fills, walking bits, offset-sensitive bytes, inverse offsets, and a
		// deterministic stream. Every byte, including odd-sized blob tails, is visited.
		for pattern := 0; pattern < 17; pattern++ {
			state := uint32(0x6d2b79f5)
			for off := range expected {
				var value byte
				switch {
				case pattern < 6:
					value = []byte{0, 0xff, 0x55, 0xaa, 1, 0x80}[pattern]
				case pattern < 14:
					value = 1 << uint(pattern-6)
				case pattern == 14:
					value = byte(off*37 + off/251)
				case pattern == 15:
					value = ^byte(off*37 + off/251)
				default:
					state ^= state << 13
					state ^= state >> 17
					state ^= state << 5
					value = byte(state)
				}
				expected[off] = value
			}
			copy(data, expected)
			if pattern == 16 {
				runtime.GC()
			}
			if c.alias() != c.ptr || !bytes.Equal(alias, expected) {
				t.Fatalf("%s pattern %d alias/GC", c.name, pattern)
			}
			h := sha256.Sum256(alias)
			o.Hashes = append(o.Hashes, hex.EncodeToString(h[:]))
			neighbors(i)
			// Reverse the alias direction too, without changing any expected capture.
			for off := range alias {
				alias[off] ^= 0xff
				expected[off] ^= 0xff
			}
			if !bytes.Equal(data, expected) {
				t.Fatalf("%s reverse alias pattern %d", c.name, pattern)
			}
			copy(data, originals[i])
		}
		if c.writeWord != nil && c.writePointer == nil {
			for _, v := range []uint32{0, 1, 0x7fffffff, 0x80000000, 0xffffffff, 0x55aa55aa} {
				c.writeWord(v)
				if binary.LittleEndian.Uint32(data) != v {
					t.Fatalf("%s typed integer store", c.name)
				}
				binary.LittleEndian.PutUint32(data, ^v)
				if c.readWord() != ^v {
					t.Fatalf("%s typed integer load", c.name)
				}
			}
		}
		if c.name == "dword_5d4594_831236" || c.name == "dword_5d4594_1309720" || c.name == "dword_5d4594_1090100" {
			target := rawStorageContracts[index["byte_5D4594"]].ptr
			c.writeWord(uint32(uintptr(target)))
			var got unsafe.Pointer
			switch c.name {
			case "dword_5d4594_831236":
				got = unsafe.Pointer(Get_dword_5d4594_831236())
			case "dword_5d4594_1309720":
				got = unsafe.Pointer(Get_dword_5d4594_1309720())
			case "dword_5d4594_1090100":
				got = unsafe.Pointer(Get_dword_5d4594_1090100())
			}
			if got != target {
				t.Fatalf("%s mixed-extern production getter", c.name)
			}
		}
		if c.writePointer != nil {
			// Only real foreign addresses enter a typed pointer slot. Arbitrary bit
			// patterns above are byte-only and never become GC-visible Go pointers.
			target := rawStorageContracts[index["byte_5D4594"]].ptr
			for _, p := range []unsafe.Pointer{nil, target, unsafe.Add(target, 128)} {
				c.writePointer(p)
				runtime.GC()
				if c.name == "dword_5d4594_1548532" {
					Set_dword_5d4594_1548532(nil)
					if c.readWord() != 0 {
						t.Fatal("game setter nil")
					}
					Set_dword_5d4594_1548532(p)
				}
				if binary.LittleEndian.Uint32(data) != uint32(uintptr(p)) || c.readWord() != uint32(uintptr(p)) {
					t.Fatalf("%s typed pointer alias", c.name)
				}
			}
		}
		neighbors(i)
		copy(data, originals[i])
		observations = append(observations, o)
	}
	// Exercise typed array and record field addressing against the raw byte ABI.
	eq := rawStorageContracts[index["array_5D4594_1049872"]]
	for n := 0; n < 9; n++ {
		v := uint32(n+1) * 0x1020304
		if rawStorageTypedEquipment(n, v) != v || binary.LittleEndian.Uint32(unsafe.Slice((*byte)(eq.ptr), eq.size)[4*n:]) != v {
			t.Fatalf("equipment cell %d", n)
		}
	}
	neighbors(index[eq.name])
	restore()
	grid := rawStorageContracts[index["nox_client_inventory_grid_1050020"]]
	ptr := rawStorageContracts[index["byte_581450"]].ptr
	for n := 0; n < 84; n++ {
		code, tail, count := uint32(n+1)*0x01020304, ^uint32(n), byte(n+1)
		rawStorageTypedInventory(n, ptr, code, tail, count)
		cell := unsafe.Slice((*byte)(unsafe.Add(grid.ptr, n*148)), 148)
		if binary.LittleEndian.Uint32(cell) != uint32(uintptr(ptr)) || binary.LittleEndian.Uint32(cell[4:]) != code || cell[140] != count || binary.LittleEndian.Uint32(cell[144:]) != tail {
			t.Fatalf("inventory cell %d", n)
		}
	}
	neighbors(index[grid.name])
	restore()
	list := rawStorageContracts[index["nox_common_maplist"]]
	rawStorageTypedList(list.ptr)
	for off := uintptr(0); off < 12; off += 4 {
		if *(*uint32)(unsafe.Add(list.ptr, off)) != uint32(uintptr(list.ptr)) {
			t.Fatalf("list field %d", off)
		}
	}
	neighbors(index[list.name])
	restore()
	for i, c := range rawStorageContracts {
		if !bytes.Equal(unsafe.Slice((*byte)(c.ptr), c.size), originals[i]) {
			t.Fatalf("%s restoration", c.name)
		}
	}
	if path := os.Getenv("OPENNOX_RAW_STORAGE_CAPTURE"); path != "" {
		b, err := json.MarshalIndent(observations, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		if err = os.WriteFile(path, append(b, '\n'), 0644); err != nil {
			t.Fatal(err)
		}
	}
	t.Logf("checked52 owners/884 full-region patterns, typed aliases, blob offsets, neighboring storage, GC and restoration")
}
