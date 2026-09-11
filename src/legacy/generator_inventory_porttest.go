//go:build porttest

package legacy

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

// PortTestGeneratorInventorySpec exercises real cloned modifiers and NPC equipment.
type PortTestGeneratorInventorySpec struct {
	Count    int // 0..3
	Subclass uint32
	Equipped [3]bool
}

type PortTestGeneratorInventoryResult struct {
	Source, Destination []uint32
	Init                [][]uint32
	Items               [][]uint32 // raw C-owned source records, then allocator-backed clones
	Intact              bool
}

type portTestGeneratorInventoryFixture struct {
	proxy       *portTestRoamOwnerServer
	raw         [][]byte
	initRaw     [][]byte
	items       []*server.Object
	source, dst *server.Object
	beforeSrc   []byte
	beforeDst   []byte
	idsBefore   map[uint32]uint32
	idsHad      map[uint32]bool
	count       int
	want        int
}

// portTestGeneratorInventoryFixtureNew supplies C-backed source inventory
// records. It is deliberately separate from generator_objects_porttest.go so
// the owner can install it only for Op=54F2B0 inventory cases.
//
// Integration API:
//
//	inv, freeInv := portTestGeneratorInventoryFixtureNew(proxy)
//	defer freeInv()
//	inv.prepare(objects.source, objects.destination, spec.Inventory)
//	C.nox_xxx_unitCreatureCopyUC_54F2B0(combatPtr(objects.source), combatPtr(objects.destination))
//	result := inv.trace(normalize) // captures, frees clones, restores source/destination
func portTestGeneratorInventoryFixtureNew(proxy *portTestRoamOwnerServer) (fixture *portTestGeneratorInventoryFixture, free func()) {
	fixture = &portTestGeneratorInventoryFixture{proxy: proxy, idsBefore: make(map[uint32]uint32), idsHad: make(map[uint32]bool)}
	var releases []func()
	for range 3 {
		b, release := alloc.Make([]byte{}, int(unsafe.Sizeof(server.Object{}))+16)
		fixture.raw = append(fixture.raw, b)
		fixture.items = append(fixture.items, (*server.Object)(unsafe.Pointer(&b[8])))
		releases = append(releases, release)
		init, freeInit := alloc.Make([]byte{}, 36)
		fixture.initRaw = append(fixture.initRaw, init)
		releases = append(releases, freeInit)
	}
	oldCache := *memmap.PtrUint32(0x5D4594, 1564960)
	return fixture, func() {
		*memmap.PtrUint32(0x5D4594, 1564960) = oldCache
		fixture.cleanup()
		fixture.restoreIDs()
		for i := len(releases) - 1; i >= 0; i-- {
			releases[i]()
		}
	}
}

func (f *portTestGeneratorInventoryFixture) prepare(src, dst *server.Object, spec PortTestGeneratorInventorySpec) {
	if spec.Count < 0 || spec.Count > len(f.items) {
		panic("invalid generator inventory source count")
	}
	f.cleanup()
	f.restoreIDs()
	if src.UpdateData == nil {
		panic("generator inventory source has no update data")
	}
	f.source, f.dst, f.count = src, dst, spec.Count
	f.want = 0
	if spec.Subclass&0x10 != 0 {
		f.want = spec.Count
	}
	f.beforeSrc = bytes.Clone(unsafe.Slice((*byte)(src.CObj()), int(unsafe.Sizeof(*src))))
	f.beforeDst = bytes.Clone(unsafe.Slice((*byte)(dst.CObj()), int(unsafe.Sizeof(*dst))))

	// Inventory cases deliberately use zeroed source update data. The patterned
	// update-data contract belongs to the no-inventory 54F2B0 cases.
	clear(unsafe.Slice((*byte)(src.UpdateData), 2200))
	src.ObjSubClass = object.SubClass(spec.Subclass)
	dst.ObjSubClass = 0x10
	src.InvFirstItem = nil
	for i, b := range f.raw {
		clear(b)
		for j := 0; j < 8; j++ {
			b[j], b[len(b)-8+j] = 0xa5, 0x5a
		}
		if i >= spec.Count {
			continue
		}
		it := f.items[i]
		it.TypeInd = uint16(15 + i%2)
		it.ObjClass = f.proxy.core.Types.ByInd(int(it.TypeInd)).Class()
		ib := f.initRaw[i]
		clear(ib)
		for j := 0; j < 8; j++ {
			ib[j], ib[len(ib)-8+j] = 0xa5, 0x5a
		}
		it.InitData = unsafe.Pointer(&ib[8])
		*(*unsafe.Pointer)(it.InitData) = unsafe.Pointer(f.proxy.callbacks.modifiers.WeaponPower1)
		*(*uint32)(unsafe.Add(it.InitData, 16)) = 0x12340000 + uint32(i)
		it.InvHolder = src
		if spec.Equipped[i] {
			it.ObjFlags |= object.FlagEquipped
		}
		if i+1 < spec.Count {
			it.InvNextItem = f.items[i+1]
		}
	}
	for i := 1; i < spec.Count; i++ {
		f.items[i].Field125 = f.items[i-1]
	}
	if spec.Count != 0 {
		src.InvFirstItem = f.items[0]
	}

	f.idsBefore, f.idsHad = make(map[uint32]uint32), make(map[uint32]bool)
	remember := func(p unsafe.Pointer, id uint32) {
		key := uint32(uintptr(p))
		old, had := f.proxy.life.ids[key]
		f.idsBefore[key], f.idsHad[key] = old, had
		f.proxy.life.ids[key] = id
	}
	for i := 0; i < spec.Count; i++ {
		remember(f.items[i].CObj(), uint32(9200+i))
		remember(f.items[i].InitData, uint32(9400+i))
	}
}

func (f *portTestGeneratorInventoryFixture) trace(normalize func(uint32) uint32) (out PortTestGeneratorInventoryResult) {
	if f.source == nil || f.dst == nil {
		panic("generator inventory fixture not prepared")
	}
	words := func(p unsafe.Pointer) []uint32 {
		b := unsafe.Slice((*byte)(p), 772)
		out := make([]uint32, 0, len(b)/4)
		for i := 0; i < len(b); i += 4 {
			out = append(out, normalize(binary.LittleEndian.Uint32(b[i:])))
		}
		return out
	}
	out.Intact = true
	for _, b := range append(append([][]byte{}, f.raw...), f.initRaw[:f.count]...) {
		for i := 0; i < 8; i++ {
			out.Intact = out.Intact && b[i] == 0xa5 && b[len(b)-8+i] == 0x5a
		}
	}
	seen := make(map[*server.Object]bool)
	var clones []*server.Object
	for it := f.dst.InvFirstItem; it != nil; it = it.InvNextItem {
		if len(clones) == f.want || seen[it] {
			panic("generator inventory clone list cycle or count mismatch")
		}
		seen[it] = true
		key := uint32(uintptr(it.CObj()))
		old, had := f.proxy.life.ids[key]
		f.idsBefore[key], f.idsHad[key] = old, had
		f.proxy.life.ids[key] = uint32(9300 + len(clones))
		ik := uint32(uintptr(it.InitData))
		oldInit, hadInit := f.proxy.life.ids[ik]
		f.idsBefore[ik], f.idsHad[ik] = oldInit, hadInit
		f.proxy.life.ids[ik] = uint32(9500 + len(clones))
		clones = append(clones, it)
	}
	if len(clones) != f.want {
		panic(fmt.Sprintf("generator inventory copied %d items, want %d", len(clones), f.want))
	}
	// Register clone IDs before taking destination words, whose FirstItem field
	// points at the first clone.
	out.Source, out.Destination = words(f.source.CObj()), words(f.dst.CObj())
	for i := 0; i < f.count; i++ {
		out.Items = append(out.Items, words(f.items[i].CObj()))
	}
	for _, it := range clones {
		var init []uint32
		for _, v := range unsafe.Slice((*uint32)(it.InitData), 5) {
			init = append(init, normalize(v))
		}
		out.Init = append(out.Init, init)
		out.Items = append(out.Items, words(it.CObj()))
	}
	f.cleanup()
	return out
}

// cleanup frees cloned buffers before returning object records to the allocator.
// It runs before destination bytes are restored, while inventory links survive.
func (f *portTestGeneratorInventoryFixture) cleanup() {
	if f.dst == nil {
		return
	}
	for it := f.dst.InvFirstItem; it != nil; {
		next := it.InvNextItem
		it.InvNextItem, it.Field125, it.InvHolder = nil, nil, nil
		for _, p := range []unsafe.Pointer{it.InitData, it.UseData.Ptr, it.UpdateData, unsafe.Pointer(it.HealthData), it.Field189} {
			if p != nil {
				alloc.FreePtr(p)
			}
		}
		it.InitData = nil
		it.UseData.Ptr = nil
		it.UpdateData = nil
		it.HealthData = nil
		it.Field189 = nil
		f.proxy.core.Objs.FreeObject(it)
		it = next
	}
	f.dst.InvFirstItem = nil
	copy(unsafe.Slice((*byte)(f.source.CObj()), len(f.beforeSrc)), f.beforeSrc)
	copy(unsafe.Slice((*byte)(f.dst.CObj()), len(f.beforeDst)), f.beforeDst)
	f.source, f.dst, f.count = nil, nil, 0
}

func (f *portTestGeneratorInventoryFixture) restoreIDs() {
	for key := range f.idsBefore {
		if f.idsHad[key] {
			f.proxy.life.ids[key] = f.idsBefore[key]
		} else {
			delete(f.proxy.life.ids, key)
		}
	}
	clear(f.idsBefore)
	clear(f.idsHad)
}
