//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
	"unsafe"
)

func TestSummonLifecycle(t *testing.T) {
	o := newSummonOwner(t)
	var rows []summonResult
	for _, width := range []uint32{640, 1024} {
		for _, state := range []uint32{0, 1, 2, 3, 255} {
			o.reset(t)
			*o.words["nox_win_width"] = width
			o.init(t)
			o.check(t, *o.summonWords["dword_5d4594_1320988"] == width-95, "constructor right edge")
			for i := 0; i < 4; i++ {
				o.check(t, o.record(i)[2] == 0, "constructor clears activity")
			}
			rows = append(rows, o.snapshot(fmt.Sprintf("width%d-state%d-load", width, state), 1))
			*o.word(1321200) = 0xaabbcc00 | state
			box := *o.summonWords["dword_5d4594_1321036"]
			for frame := 0; frame < 10; frame++ {
				ret := o.call("nox_xxx_guiDrawSummonBox_4C1FE0", box)
				o.check(t, *o.word(1321200)&0xffffff00 == 0xaabbcc00, "slider writes state byte only")
				rows = append(rows, o.snapshot(fmt.Sprintf("width%d-state%d-frame%d", width, state, frame), ret))
			}
		}
	}
	spellbookCapture(t, "summon-lifecycle", rows, "44e41e19ed51745591ae2a3fb11730d33abd68eb9b6839700fc087e5ecc3922e")
}
func TestSummonCreatureLifecycle(t *testing.T) {
	o := newSummonOwner(t)
	var rows []summonResult
	for _, quiet := range []uint32{0, 1} {
		for _, state := range []uint32{0, 1, 2, 3} {
			o.reset(t)
			o.init(t)
			*o.word(1321200) = 0x11223300 | state
			typ := uint32(o.c.Things.TypeByID("PortSmallCreature").Index())
			for i := 0; i < 5; i++ {
				o.call("nox_xxx_cliSummonCreat_4C2E50", uint32(71+i), typ, quiet)
				want := i + 1
				if want > 4 {
					want = 4
				}
				o.check(t, *o.summonWords["dword_5d4594_1321196"] == uint32(want), "add caps at four records")
				rows = append(rows, o.snapshot(fmt.Sprintf("quiet%d-state%d-add%d", quiet, state, i), 0))
			}
			o.call("nox_xxx_cliSummonCreat_4C2E50", 71, typ, quiet)
			o.check(t, *o.summonWords["dword_5d4594_1321196"] == 4, "duplicate does not increment")
			rows = append(rows, o.snapshot(fmt.Sprintf("quiet%d-state%d-duplicate", quiet, state), 0))
			for _, id := range []uint32{999, 72, 74, 71, 73, 73} {
				o.call("nox_xxx_cliSummonOnDieOrBanish_4C3140", id, quiet)
				rows = append(rows, o.snapshot(fmt.Sprintf("quiet%d-state%d-remove%d", quiet, state, id), 0))
			}
			o.check(t, *o.summonWords["dword_5d4594_1321196"] == 0, "remove empties count")
			o.check(t, *o.word(1321200) == 0x11223303, "last removal closes using byte write")
		}
	}
	spellbookCapture(t, "summon-creature-lifecycle", rows, "0ecaa4cf21fd47996e365c8613bc92d5df2a5e9e9c24d6f1f596100f07943f73")
}
func TestSummonClassification(t *testing.T) {
	o := newSummonOwner(t)
	var rows []summonResult
	o.reset(t)
	typ := o.c.Things.TypeByID("PortSmallCreature")
	for _, bits := range []uint32{0, 1, 2, 3, 4, 5, 6, 7, 0x80000000, 0xffffffff} {
		typ.ObjSubClass = object.SubClass(bits)
		ret := o.call("sub_4C2EF0", uint32(typ.Index()))
		want := uint32(4)
		if bits&1 != 0 {
			want = 1
		} else if bits&2 != 0 {
			want = 2
		}
		o.check(t, ret == want, "footprint precedence")
		rows = append(rows, o.snapshot(fmt.Sprintf("bits%x", bits), ret))
	}
	plant := uint32(o.c.Things.TypeByID("CarnivorousPlant").Index())
	for mask := 0; mask < 16; mask++ {
		o.reset(t)
		for i := 0; i < 4; i++ {
			r := o.record(i)
			r[1] = plant
			r[2] = 1
			r[6] = uint32(i)
			if mask>>i&1 != 0 {
				r[1] = uint32(typ.Index())
			}
			o.check(t, o.call("sub_4C2DD0", o.ptr(i)) == uint32(summonBool(mask>>i&1 != 0)), "plant eligibility")
		}
		ret := o.call("sub_4C2E00")
		o.check(t, ret == uint32(summonBool(mask != 0)), "any mobile creature")
		rows = append(rows, o.snapshot(fmt.Sprintf("mobile%d", mask), ret))
	}
	spellbookCapture(t, "summon-classification", rows, "ce994ee04f5082bdaa37b786f88e6338823af94c891fad660fc75a5d1feb7278")
}
func TestSummonCommands(t *testing.T) {
	o := newSummonOwner(t)
	var rows []summonResult
	for _, selected := range []bool{false, true} {
		for _, command := range []uint32{0, 1, 2, 3, 4, 5, 255, 256, 0xffffffff} {
			o.reset(t)
			o.init(t)
			r := o.record(0)
			r[0] = 0x98761234
			r[2] = 1
			r[6] = 0
			p := uint32(0)
			if selected {
				p = o.ptr(0)
			}
			o.call("nox_client_orderCreature", p, command)
			got := o.snapshot(fmt.Sprintf("selected%t-command%x", selected, command), 0)
			if !selected && command == 1 {
				o.check(t, len(got.Messages) == 0, "untargeted attack produces no order")
			} else {
				o.check(t, len(got.Messages) == 1, "one creature order")
				msg := got.Messages[0]
				o.check(t, len(msg) == 4 && msg[0] == 0x78 && msg[3] == byte(command), "order encoding")
				if selected {
					o.check(t, msg[1] == 0x34 && msg[2] == 0x12, "low sixteen-bit creature code")
				} else {
					o.check(t, msg[1] == 0 && msg[2] == 0, "group creature code")
				}
			}
			rows = append(rows, got)
		}
	}
	spellbookCapture(t, "summon-commands", rows, "d384c488b52a08e35728f6097bd84e1dac2ab2e20ce66a5b0eb6921cec5ec004")
}
func TestSummonMenus(t *testing.T) {
	o := newSummonOwner(t)
	var rows []summonResult
	pos, free := alloc.New([2]int32{})
	defer free()
	p := uint32(uintptr(unsafe.Pointer(pos)))
	for _, xy := range [][2]int32{{0, 0}, {100, 100}, {639, 479}, {640, 480}, {-1, -1}, {800, 600}} {
		o.reset(t)
		o.init(t)
		*pos = xy
		o.call("nox_xxx_wndSummonCreateList_4C2560", p)
		o.collect()
		o.check(t, *o.summonWords["dword_5d4594_1321044"] != 0, "menu owns a real window")
		rows = append(rows, o.snapshot(fmt.Sprintf("menu%d,%d", xy[0], xy[1]), 0))
		ret := o.call("nox_xxx_guiHideSummonWindow_4C2470")
		o.check(t, *o.summonWords["dword_5d4594_1321044"] == 0 && *o.summonWords["dword_5d4594_1321204"] == 0, "close clears menu and selection")
		rows = append(rows, o.snapshot(fmt.Sprintf("close%d,%d", xy[0], xy[1]), ret))
	}
	spellbookCapture(t, "summon-menus", rows, "c0ec411613a20e3183fb1924952b0520244243e9815982715095ff8c6de2d05a")
}
