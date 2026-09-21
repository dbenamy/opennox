//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
	"unsafe"

	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestServerRuntimeModifierIcons(t *testing.T) {
	o := newObjectDrawingOwner(t)
	table := serverConfigOwnBytes(t, 0x587000, 27332, 120)
	ready := serverConfigOwnBytes(t, 0x5D4594, 251624, 4)
	var entries []strman.Entry
	names := make([]string, 6)
	for i := range names {
		names[i] = fmt.Sprintf("RuntimeIcon%d", i)
		entries = append(entries, strman.Entry{ID: strman.ID(fmt.Sprintf("RuntimeModifier:%d", i)), Vals: []strman.Variant{{Str: fmt.Sprintf("Effect %d é😀", i)}}})
	}
	language, restore := o.c.srv.PortTestMeterStrings(entries...)
	t.Cleanup(restore)
	language(0)
	oldLoad := legacy.Nox_xxx_gLoadImg
	t.Cleanup(func() { legacy.Nox_xxx_gLoadImg = oldLoad })
	type row struct {
		Initial                             uint32
		Missing, Pattern, Mask, Icon, Calls int
		Ready                               uint32
		Label                               string
		Found                               bool
	}
	var rows []row
	for _, initial := range []uint32{0, 1, 0xffffffff} {
		for _, missing := range []int{0, 1, 32, 63} {
			for pattern := 0; pattern < 3; pattern++ {
				clear(table)
				binary.LittleEndian.PutUint32(ready, initial)
				keys := []byte{1, 2, 4, 8, 16, 32}
				if pattern == 1 {
					keys = []byte{0, 1, 1, 127, 128, 255}
				}
				if pattern == 2 {
					keys = []byte{255, 128, 127, 1, 1, 0}
				}
				for i, key := range keys {
					table[20*i] = key
					binary.LittleEndian.PutUint32(table[20*i+4:], uint32(uintptr(unsafe.Pointer(alloc.InternCString(names[i])))))
					binary.LittleEndian.PutUint32(table[20*i+8:], uint32(uintptr(o.images[i+6].C())))
					binary.LittleEndian.PutUint32(table[20*i+12:], uint32(uintptr(unsafe.Pointer(alloc.InternCString(fmt.Sprintf("RuntimeModifier:%d", i))))))
					binary.LittleEndian.PutUint32(table[20*i+16:], 0xaabbccdd)
				}
				before := bytes.Clone(table)
				calls := 0
				legacy.Nox_xxx_gLoadImg = func(name string) *noxrender.Image {
					if calls >= 6 || name != names[calls] {
						t.Fatalf("modifier image load order: %d %q", calls, name)
					}
					i := calls
					calls++
					if missing&(1<<i) != 0 {
						return nil
					}
					return o.images[i]
				}
				for mask := 0; mask < 256; mask++ {
					got := legacy.PortTestRuntimeModifierIcon(byte(mask))
					label, found := legacy.PortTestRuntimeModifierLabel(byte(mask))
					index := -1
					for i, key := range keys {
						if int(int8(mask)) == int(key) {
							index = i
							break
						}
					}
					want, id := uint32(0), 0
					wantLabel := ""
					if index >= 0 {
						wantLabel = fmt.Sprintf("Effect %d é😀", index)
						if initial != 0 {
							id = index + 7
						} else if missing&(1<<index) == 0 {
							id = index + 1
						}
						if id != 0 {
							want = uint32(uintptr(o.images[id-1].C()))
						}
					}
					wantReady, wantCalls := initial, 0
					if initial == 0 {
						wantReady, wantCalls = 1, 6
					}
					if got != want || label != wantLabel || found != (index >= 0) || calls != wantCalls || binary.LittleEndian.Uint32(ready) != wantReady {
						t.Fatalf("modifier lookup initial=%d missing=%d pattern=%d mask=%d: icon=%x/%x label=%q/%q found=%v calls=%d/%d", initial, missing, pattern, mask, got, want, label, wantLabel, found, calls, wantCalls)
					}
					rows = append(rows, row{initial, missing, pattern, mask, id, calls, binary.LittleEndian.Uint32(ready), label, found})
				}
				for i := 0; i < 120; i++ {
					if i%20 >= 8 && i%20 < 12 {
						continue
					}
					if table[i] != before[i] {
						t.Fatal("modifier loader changed unrelated table byte", i)
					}
				}
			}
		}
	}
	interactionCapture(t, "server-runtime-modifier-icons", rows)
}
