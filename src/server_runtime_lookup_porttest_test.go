//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
	"unsafe"

	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestServerRuntimeEquipmentLookup(t *testing.T) {
	type row struct {
		Armor          bool
		Count, Pattern int
		Query          string
		Mask, Result   uint32
		Label          string
		Found          bool
	}
	var rows []row
	labels := []string{"Sword", "sword", "Épée", "", "Rune😀", "Σ", "σ", "İ", "i"}
	queries := []string{"Sword", "SWORD", "sword", "Épée", "éPÉE", "", "Rune😀", "Σ", "σ", "İ", "i", "missing", "Sword\x00ignored", "ſword"}
	for _, armor := range []bool{false, true} {
		off, maxCount := uintptr(33392), 27
		if armor {
			off, maxCount = 35496, 26
		}
		table := serverConfigOwnBytes(t, 0x587000, off, (maxCount+1)*12)
		for _, count := range []int{0, 1, 3, 9, maxCount} {
			for pattern := 0; pattern < 2; pattern++ {
				clear(table)
				var names []string
				var masks []uint32
				for i := 0; i < count; i++ {
					name := labels[i%len(labels)]
					if pattern == 1 && i == 0 {
						name = ""
					}
					mask := uint32(1) << uint(i)
					if pattern == 1 && i == 0 {
						mask = 0
					}
					names = append(names, name)
					masks = append(masks, mask)
					binary.LittleEndian.PutUint32(table[i*12:], uint32(uintptr(unsafe.Pointer(alloc.InternCString16(name)))))
					binary.LittleEndian.PutUint32(table[i*12+8:], mask)
				}
				before := bytes.Clone(table)
				for _, query := range queries {
					result := legacy.PortTestRuntimeEquipmentMask(armor, query)
					if count == 0 && result != 0 {
						t.Fatal("empty equipment table lookup")
					}
					if count > 0 && pattern == 0 && (query == "SWORD" || query == "Sword\x00ignored") && result != 1 {
						t.Fatal("first ASCII-folded equipment match", query, result)
					}
					rows = append(rows, row{Armor: armor, Count: count, Pattern: pattern, Query: query, Result: result})
				}
				for bit := -1; bit <= 32; bit++ {
					mask := uint32(0)
					if bit >= 0 && bit < 32 {
						mask = uint32(1) << uint(bit)
					} else if bit == 32 {
						mask = 0xffffffff
					}
					label, found := legacy.PortTestRuntimeEquipmentLabel(armor, mask)
					want, exists := "", false
					for i, v := range masks {
						if v == mask {
							want, exists = names[i], true
							break
						}
					}
					if label != want || found != exists {
						t.Fatal("equipment mask lookup", armor, count, pattern, mask, label, found, want, exists)
					}
					rows = append(rows, row{Armor: armor, Count: count, Pattern: pattern, Mask: mask, Label: label, Found: found})
				}
				if !bytes.Equal(before, table) {
					t.Fatal("equipment lookup changed table")
				}
			}
		}
	}
	interactionCapture(t, "server-runtime-equipment-lookup", rows)
}

func TestServerRuntimeEquipmentLoad(t *testing.T) {
	o := newObjectDrawingOwner(t)
	var entries []strman.Entry
	for i := 0; i < 27; i++ {
		entries = append(entries, strman.Entry{ID: strman.ID(fmt.Sprintf("RuntimeEquip:%d", i)), Vals: []strman.Variant{{Str: fmt.Sprintf("Loaded %02d é😀", i)}}})
	}
	language, restore := o.c.srv.PortTestMeterStrings(entries...)
	t.Cleanup(restore)
	language(0)
	type row struct {
		Armor, MissingFirst bool
		Count, Step         int
		Initial, Ready      uint32
		Labels              []string
	}
	var rows []row
	for _, armor := range []bool{false, true} {
		off, ready, maxCount := uintptr(33392), uintptr(371248), 27
		if armor {
			off, ready, maxCount = 35496, 371256, 26
		}
		table := serverConfigOwnBytes(t, 0x587000, off, (maxCount+1)*12)
		flag := serverConfigOwnBytes(t, 0x5D4594, ready, 4)
		for _, count := range []int{0, 1, maxCount} {
			for _, initial := range []uint32{0, 1, 0xffffffff} {
				for _, missingFirst := range []bool{false, true} {
					clear(table)
					binary.LittleEndian.PutUint32(flag, initial)
					for i := 0; i < count; i++ {
						binary.LittleEndian.PutUint32(table[i*12:], uint32(uintptr(unsafe.Pointer(alloc.InternCString16(fmt.Sprintf("Before %02d", i))))))
						binary.LittleEndian.PutUint32(table[i*12+4:], uint32(uintptr(unsafe.Pointer(alloc.InternCString(fmt.Sprintf("RuntimeEquip:%d", i))))))
						binary.LittleEndian.PutUint32(table[i*12+8:], uint32(1)<<uint(i))
					}
					if missingFirst {
						binary.LittleEndian.PutUint32(table[4:], 0)
					}
					for step := 0; step < 2; step++ {
						legacy.PortTestRuntimeEquipmentLoad(armor)
						r := row{Armor: armor, MissingFirst: missingFirst, Count: count, Step: step, Initial: initial, Ready: binary.LittleEndian.Uint32(flag)}
						expectedReady := initial
						if initial == 0 {
							expectedReady = 1
						}
						if r.Ready != expectedReady {
							t.Fatal("equipment ready flag", r)
						}
						for i := 0; i < count; i++ {
							p := (*uint16)(unsafe.Pointer(uintptr(binary.LittleEndian.Uint32(table[i*12:]))))
							got := alloc.GoString16(p)
							want := fmt.Sprintf("Before %02d", i)
							if initial == 0 && !missingFirst {
								want = fmt.Sprintf("Loaded %02d é😀", i)
							}
							if got != want {
								t.Fatal("equipment label loading", armor, count, initial, missingFirst, step, i, got, want)
							}
							r.Labels = append(r.Labels, got)
						}
						rows = append(rows, r)
					}
				}
			}
		}
	}
	interactionCapture(t, "server-runtime-equipment-load", rows)
}
