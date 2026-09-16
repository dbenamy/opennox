//go:build porttest

package opennox

import (
	"fmt"
	"sort"
	"strings"
	"testing"
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

type summonOwner struct {
	*spellbookOwner
	summonWords map[string]*uint32
	raw         []uint32
	constants   []byte
	auxiliary   [][]byte
}
type summonResult struct {
	Book        spellbookResult
	Words       map[string]uint32
	Region      []uint32
	Messages    [][]byte
	SpriteFlags []uint32
}

func newSummonOwner(t *testing.T) *summonOwner {
	o := &summonOwner{spellbookOwner: newSpellbookOwner(t, "CarnivorousPlant", "PortSmallCreature", "PortMediumCreature", "PortLargeCreature")}
	for _, reg := range [][3]uintptr{{0x5D4594, 1200916, 512}, {0x5D4594, 740076, 41 * 28}, {0x587000, 70500, 41 * 4}} {
		b := memmap.BlobByAddr(reg[0]).Data[reg[1] : reg[1]+reg[2]]
		saved := append([]byte(nil), b...)
		t.Cleanup(func() { copy(b, saved) })
		o.auxiliary = append(o.auxiliary, b)
	}
	var restore func()
	o.summonWords, restore = legacy.PortTestSummonWords()
	t.Cleanup(restore)
	raw := memmap.BlobByAddr(0x5D4594).Data[1320988:1321216]
	o.raw = unsafe.Slice((*uint32)(unsafe.Pointer(&raw[0])), len(raw)/4)
	old := append([]uint32(nil), o.raw...)
	t.Cleanup(func() { copy(o.raw, old) })
	o.constants = memmap.BlobByAddr(0x587000).Data[184344:184556]
	saved := append([]byte(nil), o.constants...)
	t.Cleanup(func() { copy(o.constants, saved) })
	for i := range o.raw {
		o.c.dataRefs[uint32(uintptr(unsafe.Pointer(&o.raw[i])))] = 0xecd00000 + uint32(i*4)
	}
	callbacks := legacy.PortTestSummonCallbacks()
	var names []string
	for n := range callbacks {
		names = append(names, n)
	}
	sort.Strings(names)
	for i, n := range names {
		if callbacks[n] != nil {
			o.c.callbackRefs[callbacks[n]] = 0xece00000 + uint32(i)
		}
	}
	oldLoad := legacy.Nox_xxx_gLoadImg
	t.Cleanup(func() { legacy.Nox_xxx_gLoadImg = oldLoad })
	legacy.Nox_xxx_gLoadImg = func(name string) *noxrender.Image {
		if strings.HasPrefix(name, "CreatureCage") {
			o.loads = append(o.loads, name)
			i := 0
			for _, b := range []byte(name) {
				i = (i + int(b)) % len(o.images)
			}
			return o.images[i]
		}
		return oldLoad(name)
	}
	var entries []strman.Entry
	for _, name := range []string{"ToolTipSummon", "ccs:BANISH", "ccs:ORDER", "ccs:UNUSED", "ccs:GUARD", "ccs:ESCORT", "ccs:HUNT"} {
		entries = append(entries, strman.Entry{ID: strman.ID("guisumn.c:" + name), Vals: []strman.Variant{{Str: name}}})
	}
	configure, restore := o.c.srv.Server.PortTestMeterStrings(entries...)
	t.Cleanup(restore)
	configure(0)
	return o
}
func (o *summonOwner) reset(t *testing.T) {
	o.resetBook(t)
	for _, b := range o.auxiliary {
		clear(b)
	}
	for i := 0; i < 41; i++ {
		*memmap.PtrPtr(0x587000, 70500+uintptr(4*i)) = unsafe.Pointer(alloc.InternCString(""))
	}
	for i, n := range []string{"PortSmallCreature", "PortMediumCreature", "PortLargeCreature", "CarnivorousPlant"} {
		*memmap.PtrPtr(0x587000, 70500+uintptr(4*(i+1))) = unsafe.Pointer(alloc.InternCString(n))
	}
	clear(o.raw)
	clear(o.constants)
	for _, p := range o.summonWords {
		*p = 0
	}
	copy(o.constants[184456-184344:], blobdata.PortTestSummonConstants())
	for i, n := range []string{"ccs:BANISH", "ccs:ORDER", "ccs:UNUSED", "ccs:GUARD", "ccs:ESCORT", "ccs:HUNT"} {
		*memmap.PtrPtr(0x587000, 184344+uintptr(4*i)) = unsafe.Pointer(alloc.InternCString(n))
	}
	for i, n := range []string{"PortSmallCreature", "PortMediumCreature", "PortLargeCreature", "CarnivorousPlant"} {
		o.c.Things.TypeByID(n).ObjSubClass = object.SubClass([]int{1, 2, 0, 1}[i])
	}
	o.c.srv.NetList.ResetByInd(31, netlist.Kind0)
}
func (o *summonOwner) call(op string, a ...uint32) uint32 {
	var args [5]uint32
	copy(args[:], a)
	return legacy.PortTestSummonInvoke(op, args)
}
func (o *summonOwner) word(off int) *uint32  { return &o.raw[(off-1320988)/4] }
func (o *summonOwner) record(i int) []uint32 { return o.raw[(1321052-1320988)/4+i*8:][:8] }
func (o *summonOwner) ptr(i int) uint32      { return uint32(uintptr(unsafe.Pointer(&o.record(i)[0]))) }
func (o *summonOwner) check(t *testing.T, ok bool, label string) {
	t.Helper()
	if !ok {
		t.Fatal("summon contract: " + label)
	}
}
func (o *summonOwner) snapshot(label string, ret uint32) summonResult {
	r := summonResult{Book: o.bookSnapshot(label, ret), Words: make(map[string]uint32), Region: append([]uint32(nil), o.raw...)}
	for n, p := range o.summonWords {
		v := *p
		switch n {
		case "dword_5d4594_1321024", "dword_5d4594_1321032", "dword_5d4594_1321036", "dword_5d4594_1321040", "dword_5d4594_1321044", "dword_5d4594_1321204":
			v = o.normalize(v)
		}
		r.Words[n] = v
	}
	for _, off := range []int{1320996, 1321008, 1321012, 1321016, 1321020, 1321028, 1321180, 1321184, 1321188, 1321192} {
		i := (off - 1320988) / 4
		r.Region[i] = o.normalize(r.Region[i])
	}
	o.c.srv.NetList.ByInd(31, netlist.Kind0).Each(func(b []byte) bool { r.Messages = append(r.Messages, append([]byte(nil), b...)); return false })
	return r
}
func (o *summonOwner) init(t *testing.T) {
	o.check(t, o.call("nox_xxx_guiSummonCreatureLoad_4C1D80") == 1, "constructor succeeds")
	o.collect()
}
func (o *summonOwner) label(prefix string, args ...any) string {
	return fmt.Sprint(append([]any{prefix}, args...)...)
}

func summonBool(v bool) int {
	if v {
		return 1
	}
	return 0
}
