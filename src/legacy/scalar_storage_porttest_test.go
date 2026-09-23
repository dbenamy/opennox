//go:build porttest

package legacy

import (
	"encoding/binary"
	"encoding/json"
	"os"
	"runtime"
	"sort"
	"testing"
	"unsafe"
)

func TestScalarStorageOwners(t *testing.T) {
	type region struct {
		lo, hi uintptr
		name   string
	}
	var regions []region
	raw := func(c scalarStorageContract) uint64 {
		b := unsafe.Slice((*byte)(c.ptr), c.size)
		if c.size == 4 {
			return uint64(binary.LittleEndian.Uint32(b))
		}
		return binary.LittleEndian.Uint64(b)
	}
	put := func(c scalarStorageContract, v uint64) {
		b := unsafe.Slice((*byte)(c.ptr), c.size)
		if c.size == 4 {
			binary.LittleEndian.PutUint32(b, uint32(v))
		} else {
			binary.LittleEndian.PutUint64(b, v)
		}
	}
	typed := func(c scalarStorageContract, v uint64) uint64 {
		if c.size == 4 {
			if c.signed {
				return uint64(int64(int32(v)))
			}
			return uint64(uint32(v))
		}
		return v
	}
	mask := func(c scalarStorageContract, v uint64) uint64 {
		if c.size == 4 {
			return uint64(uint32(v))
		}
		return v
	}
	originals := make([]uint64, len(scalarStorageContracts))
	for i, c := range scalarStorageContracts {
		if c.size != c.wantSize || (c.size != 4 && c.size != 8) {
			t.Fatalf("%s size=%d want=%d", c.name, c.size, c.wantSize)
		}
		if scalarStorageInitial[i] != typed(c, c.wantBits) {
			t.Fatalf("%s initial=%x want=%x", c.name, scalarStorageInitial[i], typed(c, c.wantBits))
		}
		lo := uintptr(c.ptr)
		regions = append(regions, region{lo, lo + c.size, c.name})
		originals[i] = raw(c)
	}
	sort.Slice(regions, func(i, j int) bool { return regions[i].lo < regions[j].lo })
	for i := 1; i < len(regions); i++ {
		if regions[i].lo < regions[i-1].hi {
			t.Fatalf("overlapping scalar owners %s/%s", regions[i-1].name, regions[i].name)
		}
	}
	defer func() {
		for i, c := range scalarStorageContracts {
			put(c, originals[i])
		}
	}()
	patterns := []uint64{0, 1, 0x7fffffff, 0x80000000, 0xffffffff, 0xaaaaaaaa55555555, 0x55555555aaaaaaaa, 0x0123456789abcdef, 0xffffffffffffffff, 0x8000000000000000, 0xdeadbeef}
	for bit := uint(0); bit < 64; bit++ {
		v := uint64(1) << bit
		patterns = append(patterns, v, ^v)
	}
	type observation struct {
		Name    string
		Size    uintptr
		Signed  bool
		Initial uint64
		Values  [][2]uint64
	}
	observations := make([]observation, 0, len(scalarStorageContracts))
	for index, c := range scalarStorageContracts {
		o := observation{Name: c.name, Size: c.size, Signed: c.signed, Initial: scalarStorageInitial[index]}
		for _, pattern := range patterns {
			c.set(pattern)
			fromTyped := raw(c)
			if fromTyped != mask(c, pattern) {
				t.Fatalf("%s typed store %x -> %x", c.name, pattern, fromTyped)
			}
			next := pattern ^ 0xff00ff00ff00ff00
			put(c, next)
			if pattern == 0xdeadbeef {
				runtime.GC()
			}
			fromRaw := c.get()
			if fromRaw != typed(c, next) {
				t.Fatalf("%s raw store %x -> %x", c.name, next, fromRaw)
			}
			for other, d := range scalarStorageContracts {
				if other != index && raw(d) != originals[other] {
					t.Fatalf("%s write changed %s", c.name, d.name)
				}
			}
			o.Values = append(o.Values, [2]uint64{fromTyped, fromRaw})
			put(c, originals[index])
		}
		observations = append(observations, o)
	}
	if path := os.Getenv("OPENNOX_SCALAR_STORAGE_CAPTURE"); path != "" {
		data, err := json.MarshalIndent(observations, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		if err = os.WriteFile(path, append(data, '\n'), 0644); err != nil {
			t.Fatal(err)
		}
	}
	t.Logf("checked %d owners, %d patterns each, including raw/typed aliases and GC", len(scalarStorageContracts), len(patterns))
}
