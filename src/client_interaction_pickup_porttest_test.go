//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"slices"
	"strings"
	"testing"

	"github.com/opennox/libs/object"
	"github.com/opennox/libs/strman"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
)

func TestClientInteractionPickupAndSecondary(t *testing.T) {
	o := newInventoryTransactionOwner(t, "Gold", "QuestGoldPile", "QuestGoldChest", "InteractionItem", "InteractionFiller")
	words, restore := legacy.PortTestClientInteractionWords()
	defer restore()
	cache := serverConfigOwnBytes(t, 0x5D4594, 1064928, 12)
	messages := serverConfigOwnBytes(t, 0x5D4594, 823804, 1932)
	configure, restoreStrings := o.c.srv.Server.PortTestMeterStrings(
		strman.Entry{ID: "pickup.c:CarryingTooMuch", Vals: []strman.Variant{{Str: "CarryingTooMuch"}}},
		strman.Entry{ID: "guimsg.c:systemmsg", Vals: []strman.Variant{{Str: "System: %s"}}},
	)
	defer restoreStrings()
	configure(0)
	type row struct {
		Name, Layout          string
		Count                 int
		Class, Code           uint32
		Quest, Secondary, Nil bool
		Return                uint64
		Packet                []byte
		Sounds                [][2]int
		Console               []string
	}
	var captured []row
	run := func(name, layout string, count int, class, code uint32, quest, secondary, nilObject bool, allowed bool) {
		t.Helper()
		o.reset(t)
		clear(cache)
		clear(messages)
		*words["dword_5d4594_825736"] = 0
		if quest {
			noxflags.SetGame(6144)
		}
		dr := o.item(t, name, code)
		dr.ObjClass = 0
		filler := o.item(t, "InteractionFiller", 1)
		filler.ObjClass = 0x4000000
		for i := 0; i < 84; i++ {
			binary.LittleEndian.PutUint32(o.grid[148*i:], uint32(uintptr(filler.C())))
			o.grid[148*i+140] = 1
		}
		switch layout {
		case "empty":
			o.grid[79*148+140] = 0 // row16 of final column is usable.
		case "reserved":
			for col := 0; col < 4; col++ {
				o.grid[(20+21*col)*148+140] = 0
			}
		case "stack":
			stack := o.item(t, name, 2)
			stack.ObjClass = object.Class(class)
			binary.LittleEndian.PutUint32(o.grid[0:], uint32(uintptr(stack.C())))
			o.grid[140] = byte(count)
		}
		beforeGrid := append([]byte(nil), o.grid...)
		var arg uintptr
		if !nilObject {
			arg = uintptr(dr.C())
		}
		op := "nox_xxx_clientPickup_46C140"
		if secondary {
			op = "nox_xxx_clientReportSecondaryWeapon_4BF010"
		}
		ret := interactionCall(op, arg)
		var got []byte
		packets := 0
		o.c.srv.NetList.ByInd(31, netlist.Kind0).Each(func(b []byte) bool { got = append(got, b...); packets++; return false })
		var want []byte
		if allowed {
			unit := uint16(code)
			if nilObject || code >= 0x8000 {
				unit = 0
			}
			opcode := byte(115)
			if secondary {
				opcode = 224
			}
			want = []byte{opcode, byte(unit), byte(unit >> 8)}
		}
		if !bytes.Equal(got, want) || packets != len(want)/3 || !bytes.Equal(o.grid, beforeGrid) {
			t.Fatalf("%s %s count%d class%x code%x secondary%v nil%v: packet%x want%x or changed grid", name, layout, count, class, code, secondary, nilObject, got, want)
		}
		denied := !secondary && !nilObject && !allowed
		if denied {
			if !slices.Equal(o.sounds, [][2]int{{925, 100}}) || len(o.console) != 1 || !strings.HasSuffix(o.console[0], "System: CarryingTooMuch") || *words["dword_5d4594_825736"] != 1 {
				t.Fatal("pickup rejection feedback", o.sounds, o.console)
			}
		} else if len(o.sounds) != 0 || len(o.console) != 0 || *words["dword_5d4594_825736"] != 0 {
			t.Fatal("unexpected pickup feedback")
		}
		if secondary {
			if ret != 1 || !bytes.Equal(cache, make([]byte, 12)) {
				t.Fatal("secondary return/cache", ret)
			}
		} else {
			for i, n := range []string{"Gold", "QuestGoldPile", "QuestGoldChest"} {
				if binary.LittleEndian.Uint32(cache[4*i:]) != uint32(o.c.Things.TypeByID(n).Index()) {
					t.Fatal("pickup caches, including nil input", n)
				}
			}
		}
		captured = append(captured, row{name, layout, count, class, code, quest, secondary, nilObject, ret, got, append([][2]int(nil), o.sounds...), append([]string(nil), o.console...)})
	}
	run("InteractionItem", "full", 1, 0, 7, false, false, true, false)
	for _, name := range []string{"Gold", "QuestGoldPile", "QuestGoldChest", "InteractionItem"} {
		for _, layout := range []string{"full", "empty", "reserved"} {
			run(name, layout, 1, 0, 7, false, false, false, name != "InteractionItem" || layout == "empty")
		}
	}
	for _, class := range []uint32{0, 0x10, 0x4000000, 0x4000010} {
		for _, count := range []int{1, 2, 3, 8, 9, 30, 31, 255} {
			for _, quest := range []bool{false, true} {
				limit := 31
				if class&0x10 != 0 {
					limit = 3
					if quest {
						limit = 9
					}
				}
				run("InteractionItem", "stack", count, class, 7, quest, false, false, class&0x4000000 == 0 && count < limit)
			}
		}
	}
	for _, code := range []uint32{0, 1, 0x7fff, 0x8000, 0xffff, 0xffffffff} {
		run("InteractionItem", "empty", 0, 0, code, false, false, false, true)
		run("InteractionItem", "full", 0, 0, code, false, true, false, true)
	}
	run("InteractionItem", "full", 0, 0, 0, false, true, true, true)
	interactionCapture(t, "pickup-secondary", captured)
}
