//go:build porttest

package opennox

import (
	"math"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func questProgressTable(t *testing.T, o *collisionCoreOwner) [][2]string {
	t.Helper()
	data, relocs := blobdata.PortTestQuestProgressTable()
	copy(serverConfigOwnBytes(t, 0x587000, 249896, len(data)), data)
	for _, r := range relocs {
		*memmap.PtrPtr(0x587000, r[0]) = memmap.PtrOff(0x587000, r[1])
	}
	var pairs [][2]string
	var names []string
	for i := uintptr(0); i < 53; i++ {
		off := uintptr(249896) + 16*i
		a := alloc.GoString((*byte)(*memmap.PtrPtr(0x587000, off)))
		b := alloc.GoString((*byte)(*memmap.PtrPtr(0x587000, off+8)))
		if a == "" || b == "" {
			t.Fatal("empty shipped generator pair", i)
		}
		pairs = append(pairs, [2]string{a, b})
		names = append(names, a, b)
	}
	t.Cleanup(o.s.PortTestRewardTypes(names, nil, true, 0, 0))
	return pairs
}
func TestQuestProgressGeneratorMapping(t *testing.T) {
	o := newCollisionCoreOwner(t)
	t.Cleanup(legacy.PortTestQuestProgressOwner())
	pairs := questProgressTable(t, o)
	u := newCreatureXferObject(t, o.s, "Monster")
	type row struct {
		Name         string
		Type, Return uint32
		Words        []uint32
		Cold         uint32
	}
	var rows []row
	snap := func(name string, id, got uint32) {
		r := row{Name: name, Type: id, Return: got, Cold: memmap.Uint32(0x5D4594, 2388664)}
		for i := uintptr(0); i < 53; i++ {
			r.Words = append(r.Words, memmap.Uint32(0x587000, 249900+16*i), memmap.Uint32(0x587000, 249908+16*i))
		}
		rows = append(rows, r)
	}
	if legacy.PortTestQuestProgressObject("generator-type", nil, 0) != 0 || memmap.Uint32(0x5D4594, 2388664) != 1 {
		t.Fatal("nil query initializes before admission")
	}
	snap("nil", 0, 0)
	for _, pair := range pairs {
		id := o.s.Types.ByID(pair[0]).Ind()
		want := o.s.Types.ByID(pair[1]).Ind()
		u.TypeInd = uint16(id)
		got := legacy.PortTestQuestProgressObject("generator-type", u, 0)
		if got != uint32(want) {
			t.Fatal("shipped mapping", pair, got, want)
		}
		snap(pair[0], uint32(id), got)
	}
	for _, id := range []uint16{0, 1, 65535} {
		u.TypeInd = id
		got := legacy.PortTestQuestProgressObject("generator-type", u, 0)
		if got != 0 {
			t.Fatal("missing mapping")
		}
		snap("missing", uint32(id), got)
	}
	legacy.PortTestQuestProgressObject("generator-init", nil, 0)
	snap("repeat-init", 0, 0)
	spellbookCapture(t, "quest-progress-mapping", rows, "64106ad4fee2d15d6ebd9f62ba0db2320d55a47c0809d6368f04da740e49b44a")
}
func TestQuestProgressStageState(t *testing.T) {
	t.Cleanup(legacy.PortTestQuestProgressOwner())
	var rows [][5]uint64
	for _, v := range []uint32{0, 1, 2, 5, 10, 0x7fffffff, 0x80000000, 0xffffffff, math.Float32bits(1)} {
		a := legacy.PortTestQuestProgress("stage-set", "", v)
		b := legacy.PortTestQuestProgress("stage", "", 0)
		c := legacy.PortTestQuestProgress("minions-set", "", v)
		d := *(*uint32)(unsafe.Pointer(memmap.PtrUint32(0x5D4594, 2388656)))
		if a != uint64(v) || b != uint64(v) || c != uint64(v) || d != v {
			t.Fatal("stage word preservation")
		}
		rows = append(rows, [5]uint64{uint64(v), a, b, c, uint64(d)})
	}
	spellbookCapture(t, "quest-progress-stage-state", rows, "32b3f0186866273402b909965a1a80029caeed49422c51ee8e0d88db9cef14bc")
}

func TestQuestProgressEmptyMapping(t *testing.T) {
	o := newCollisionCoreOwner(t)
	t.Cleanup(legacy.PortTestQuestProgressOwner())
	questProgressTable(t, o)
	u := newCreatureXferObject(t, o.s, "Monster")
	*memmap.PtrPtr(0x587000, 249904) = nil
	var rows [][3]uint32
	for _, cold := range []uint32{0, 1} {
		*memmap.PtrUint32(0x5D4594, 2388664) = cold
		got := legacy.PortTestQuestProgressObject("generator-type", u, 0)
		if got != 0 || memmap.Uint32(0x5D4594, 2388664) != 1 {
			t.Fatal("empty mapping initialization")
		}
		rows = append(rows, [3]uint32{cold, got, memmap.Uint32(0x5D4594, 2388664)})
	}
	spellbookCapture(t, "quest-progress-empty-mapping", rows, "07a2f8651f37050bc2c47b723fc5f4180ce2c111a1d6a1c668f2999803252280")
}
