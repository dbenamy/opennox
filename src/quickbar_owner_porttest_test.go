//go:build porttest

package opennox

import (
	"fmt"
	"sort"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

type quickbarOwner struct {
	*spellbookOwner
	bar        []uint32
	raw        []uint32
	quickWords map[string]*uint32
	configure  func([]server.PortTestSpellClassDef)
}
type quickbarResult struct {
	Book       spellbookResult
	Region     []uint32
	Words      map[string]uint32
	Messages   [][]byte
	Timeouts   map[byte]uint32
	LastButton uint32
}

func newQuickbarOwner(t *testing.T) *quickbarOwner {
	o, bar, configure := newSpellbookAdditionOwner(t)
	q := &quickbarOwner{spellbookOwner: o, bar: bar, configure: configure}
	var restore func()
	q.quickWords, restore = legacy.PortTestQuickbarWords()
	t.Cleanup(restore)
	q.raw = unsafe.Slice(memmap.PtrUint32(0x5D4594, 1047548), (1049716-1047548)/4)
	old := append([]uint32(nil), q.raw...)
	t.Cleanup(func() { copy(q.raw, old) })
	timeouts := inputKeyTimeoutsOld
	t.Cleanup(func() { inputKeyTimeoutsOld = timeouts })
	last := memmap.Uint32(0x587000, 133484)
	t.Cleanup(func() { *memmap.PtrUint32(0x587000, 133484) = last })
	callbacks := legacy.PortTestQuickbarCallbacks()
	names := make([]string, 0, len(callbacks))
	for n := range callbacks {
		names = append(names, n)
	}
	sort.Strings(names)
	for i, n := range names {
		// Keep already-established callback identities used by spellbook captures.
		if _, ok := o.c.callbackRefs[callbacks[n]]; !ok && callbacks[n] != nil {
			o.c.callbackRefs[callbacks[n]] = 0xeeb00001 + uint32(i)
		}
	}
	// Register actual row/slot addresses, including the ability loop end pointers.
	for off := uintptr(1047548); off < 1049716; off += 4 {
		p := uint32(uintptr(memmap.PtrOff(0x5D4594, off)))
		if _, ok := o.c.dataRefs[p]; !ok {
			o.c.dataRefs[p] = 0xeec00000 + uint32(off-1047548)
		}
	}
	return q
}
func (q *quickbarOwner) reset(t *testing.T) {
	clear(q.raw)
	for _, p := range q.quickWords {
		*p = 0
	}
	prepareSpellbookAddition(t, q.spellbookOwner, q.bar, 640)
	inputKeyTimeoutsOld = make(map[byte]uint32)
	q.c.srv.NetList.ResetByInd(31, netlist.Kind0)
	*memmap.PtrUint32(0x587000, 133484) = 0xffffffff
	for id := 1; id <= 5; id++ {
		*memmap.PtrUint32(0x5D4594, 1047764+uintptr(24*id)) = uint32(id)
	}
}
func (q *quickbarOwner) call(op string, a ...uint32) uint32 {
	var args [7]uint32
	copy(args[:], a)
	return legacy.PortTestQuickbarInvoke(op, args)
}
func (q *quickbarOwner) snapshot(label string, ret uint32) quickbarResult {
	r := quickbarResult{Book: q.bookSnapshot(label, ret), Region: append([]uint32(nil), q.raw...), Words: make(map[string]uint32), Timeouts: make(map[byte]uint32), LastButton: memmap.Uint32(0x587000, 133484)}
	for n, p := range q.quickWords {
		v := *p
		switch n {
		case "dword_5d4594_1049500", "dword_5d4594_1049504", "dword_5d4594_1049508", "dword_5d4594_1049512", "dword_5d4594_1049516", "dword_5d4594_1049520", "dword_5d4594_1049524", "dword_5d4594_1049532", "dword_5d4594_1049692", "dword_5d4594_1049696":
			v = q.normalize(v)
		}
		r.Words[n] = v
	}
	for _, base := range []int{1047940, 1048196, 1048452, 1048708, 1048964, 1049220} {
		for i := 51; i <= 62; i++ {
			ix := (base + 4*i - 1047548) / 4
			r.Region[ix] = q.normalize(r.Region[ix])
		}
	}
	r.Region[(1049528-1047548)/4] = q.normalize(r.Region[(1049528-1047548)/4])
	r.Region[(1049684-1047548)/4] = q.normalize(r.Region[(1049684-1047548)/4])
	for k, v := range inputKeyTimeoutsOld {
		r.Timeouts[k] = v
	}
	q.c.srv.NetList.ByInd(31, netlist.Kind0).Each(func(b []byte) bool { r.Messages = append(r.Messages, append([]byte(nil), b...)); return false })
	return r
}
func (q *quickbarOwner) check(t *testing.T, ok bool, label string) {
	t.Helper()
	if !ok {
		t.Fatal(fmt.Sprintf("quickbar contract: %s", label))
	}
}
