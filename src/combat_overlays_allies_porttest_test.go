//go:build porttest

package opennox

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
	"unsafe"
)

func combatAllyStorage(t *testing.T) []byte {
	t.Helper()
	raw := serverConfigOwnBytes(t, 0x5D4594, 1200916, 512)
	for i := range raw {
		raw[i] = 0xa5
	}
	legacy.PortTestCombatAllyClear()
	return raw
}
func combatAllySlot(t *testing.T, raw []byte, p unsafe.Pointer) int {
	t.Helper()
	if p == nil {
		return -1
	}
	delta := uintptr(p) - uintptr(unsafe.Pointer(&raw[0]))
	if delta >= 512 || delta%16 != 0 {
		t.Fatal("ally pointer outside owned slots")
	}
	return int(delta / 16)
}
func TestCombatOverlayAllyClear(t *testing.T) {
	raw := combatAllyStorage(t)
	for i, b := range raw {
		want := byte(0)
		if i%16 == 5 || i%16 == 10 || i%16 == 11 {
			want = 0xa5
		}
		if b != want {
			t.Fatalf("clear byte %d: %x want %x", i, b, want)
		}
	}
	if legacy.PortTestCombatAllyLookup(0) != nil || legacy.PortTestCombatIsAlly(0) {
		t.Fatal("inactive zero slot found")
	}
	spellbookCapture(t, "combat-overlay-ally-clear", []string{hex.EncodeToString(raw)}, "6647908fd8f09ac378705487bcd8294c74736f8dd31d6a69cabff09b1e934e33")
}
func TestCombatOverlayAllyCapacity(t *testing.T) {
	raw := combatAllyStorage(t)
	type record struct {
		Step string
		Free int
		Raw  string
	}
	var rows []record
	for i := 0; i < 32; i++ {
		if combatAllySlot(t, raw, legacy.PortTestCombatAllyFree()) != i {
			t.Fatal("first free slot")
		}
		if !legacy.PortTestCombatAllyAdd(uint32(i+1), uint16(i), uint16(65535-i)) {
			t.Fatal("add before capacity")
		}
		if combatAllySlot(t, raw, legacy.PortTestCombatAllyLookup(uint32(i+1))) != i {
			t.Fatal("lookup slot")
		}
	}
	full := hex.EncodeToString(raw)
	if legacy.PortTestCombatAllyAdd(33, 1, 2) || legacy.PortTestCombatAllyFree() != nil {
		t.Fatal("capacity overflow")
	}
	if !legacy.PortTestCombatAllyAdd(7, 300, 400) || hex.EncodeToString(raw) != full {
		t.Fatal("duplicate must preserve existing fields even at capacity")
	}
	rows = append(rows, record{"full", -1, full})
	for _, slot := range []int{0, 15, 31} {
		code := uint32(slot + 1)
		if !legacy.PortTestCombatAllyRemove(code) || legacy.PortTestCombatAllyRemove(code) {
			t.Fatal("remove return")
		}
		if legacy.PortTestCombatIsAlly(code) {
			t.Fatal("removed ally found")
		}
		if got := combatAllySlot(t, raw, legacy.PortTestCombatAllyFree()); got != slot {
			t.Fatalf("reuse slot %d want %d", got, slot)
		}
		if !legacy.PortTestCombatAllyAdd(code, 500, 600) {
			t.Fatal("reuse capacity")
		}
		rows = append(rows, record{fmt.Sprint("reuse", slot), combatAllySlot(t, raw, legacy.PortTestCombatAllyFree()), hex.EncodeToString(raw)})
	}
	spellbookCapture(t, "combat-overlay-ally-capacity", rows, "b2a6b68dc8a5a40bee80a07648c873da8635cf1a4bf318a681900bfde18925b4")
}
func TestCombatOverlayAllyUpdates(t *testing.T) {
	type record struct {
		Name          string
		First, Second uint16
		Flag          byte
		Raw           string
	}
	var rows []record
	for _, code := range []uint32{0, 1, 65535, 0x80000001, 0xffffffff} {
		for _, pair := range [][2]uint16{{0, 0}, {1, 65535}, {32768, 32767}, {65535, 65535}} {
			for _, flag := range []byte{0, 1, 127, 128, 255} {
				name := fmt.Sprintf("code=%x/pair=%v/flag=%d", code, pair, flag)
				t.Run(name, func(t *testing.T) {
					raw := combatAllyStorage(t)
					values, oldFlag := [2]uint16{123, 456}, byte(0x77)
					if legacy.PortTestCombatAllyRead(code, &values, &oldFlag) || values != [2]uint16{123, 456} || oldFlag != 0x77 {
						t.Fatal("missing lookup changed outputs")
					}
					if legacy.PortTestCombatAllyFlag(code, flag) || legacy.PortTestCombatAllyPair(code, 1, 2) || legacy.PortTestCombatAllyFirst(code, 3) {
						t.Fatal("missing update succeeded")
					}
					if !legacy.PortTestCombatAllyAdd(code, 55, 66) || !legacy.PortTestCombatIsAlly(code) {
						t.Fatal("add/membership")
					}
					if !legacy.PortTestCombatAllyPair(code, pair[0], pair[1]) || !legacy.PortTestCombatAllyFlag(code, flag) {
						t.Fatal("existing update failed")
					}
					if !legacy.PortTestCombatAllyRead(code, &values, &oldFlag) || values != pair || oldFlag != flag {
						t.Fatal("pair/flag readback")
					}
					first := pair[0] ^ 0xffff
					if !legacy.PortTestCombatAllyFirst(code, first) || !legacy.PortTestCombatAllyRead(code, &values, &oldFlag) || values != [2]uint16{first, pair[1]} || oldFlag != flag {
						t.Fatal("first-only update")
					}
					if binary.LittleEndian.Uint32(raw[:4]) != code || binary.LittleEndian.Uint32(raw[12:16]) != 1 {
						t.Fatal("slot identity/active")
					}
					rows = append(rows, record{name, values[0], values[1], oldFlag, hex.EncodeToString(raw[:16])})
				})
			}
		}
	}
	spellbookCapture(t, "combat-overlay-ally-updates", rows, "70b1d05a12bf2c83257756d44f66426f221b4ec0ead81b4ed66234f7c039e404")
}
