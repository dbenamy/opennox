//go:build porttest

package legacy

import (
	"bytes"
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

type PortTestShopLoadSpec struct {
	Names               []string
	Stage, MarkerChance uint32
	Marker              bool
}
type portTestShopLoad struct {
	configure   func(bool, uint32)
	restore     func()
	freeStrings []func()
}

func portTestShopLoadEnvironment(proxy *portTestRoamOwnerServer) *portTestShopLoad {
	configure, freeTypes := proxy.core.PortTestShopLoadTypes()
	s := &portTestShopLoad{configure: configure}
	regions := []struct {
		base, off uintptr
		n         int
	}{{0x587000, 234816, 16}, {0x587000, 207044, 84}, {0x587000, 202028, 4}, {0x5D4594, 1568276, 4}, {0x5D4594, 1564960, 4}}
	var saved [][]byte
	for _, r := range regions {
		saved = append(saved, bytes.Clone(unsafe.Slice((*byte)(memmap.PtrOff(r.base, r.off)), r.n)))
	}
	s.restore = func() {
		for _, f := range s.freeStrings {
			f()
		}
		for i, r := range regions {
			copy(unsafe.Slice((*byte)(memmap.PtrOff(r.base, r.off)), r.n), saved[i])
		}
		freeTypes()
	}
	return s
}
func (s *portTestShopLoad) prepare(sp *PortTestShopLoadSpec) {
	for _, f := range s.freeStrings {
		f()
	}
	s.freeStrings = nil
	s.configure(false, 0)
	if sp == nil {
		return
	}
	s.configure(sp.Marker, sp.MarkerChance)
	if len(sp.Names) > 3 {
		panic("shop name fixture capacity")
	}
	table := unsafe.Slice((*unsafe.Pointer)(memmap.PtrOff(0x587000, 234816)), 4)
	clear(table)
	for i, name := range sp.Names {
		ptr, free := alloc.CString(name)
		table[i] = unsafe.Pointer(ptr)
		s.freeStrings = append(s.freeStrings, free)
	}
	*memmap.PtrUint32(0x587000, 202028) = sp.Stage
	*memmap.PtrUint32(0x5D4594, 1568276) = 0
	*memmap.PtrUint32(0x5D4594, 1564960) = 0
	// Keep the actual generator. Its category table enables spell and ability
	// books, and its stage-dependent spell selector sees one eligible entry.
	clear(unsafe.Slice((*byte)(memmap.PtrOff(0x587000, 207044)), 84))
	*memmap.PtrUint8(0x587000, 207044) = 1
	*memmap.PtrUint8(0x587000, 207052) = 1
	*memmap.PtrUint8(0x587000, 207104) = 1
	*memmap.PtrUint32(0x587000, 207108) = 1
	*memmap.PtrUint32(0x587000, 207112) = 31
}
