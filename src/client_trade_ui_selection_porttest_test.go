//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
	"unsafe"
)

func TestClientTradeUISelectedCell(t *testing.T) {
	o := newTradeUIOwner(t)
	data := o.record(t, 15)
	var rows []tradeUIResult
	for side, reportSide := range []byte{1, 0} {
		for _, active := range []uint32{0, 1, 2} {
			for _, mode := range []string{"empty", "same", "different", "equipment", "full"} {
				o.reset(t)
				o.construct(t)
				o.constructTrade(t)
				*o.windowWords["dword_5d4594_1320964"] = active
				cells := legacy.PortTestTradeUICells()[side]
				selected := &cells[3]
				if mode != "empty" {
					name := "RedApple"
					count := 1
					if mode == "different" || mode == "equipment" {
						name = "Bow"
					}
					if mode == "full" {
						count = 32
					}
					o.tradeCell(t, side, 3, count, name, 200)
				}
				incoming := "RedApple"
				if mode == "equipment" {
					incoming = "Bow"
				}
				key := []string{"dword_5d4594_1320932", "dword_5d4594_1320936"}[side]
				*o.tradeWords[key] = uint32(txptr(unsafe.Pointer(selected)))
				*o.tradeWords["dword_5d4594_1320944"] = 1
				*o.tradeWords["dword_5d4594_1320948"] = 1
				tradeUIAddRecord(data, reportSide, uint16(o.c.Things.TypeByID(incoming).Index()), 999, 7)
				before := *selected
				rows = append(rows, o.tradeCapture(t, len(rows), 27, txptr(unsafe.Pointer(&data[0]))))
				if active == 0 {
					if *selected != before || *o.tradeWords[key] == 0 || *o.tradeWords["dword_5d4594_1320944"] != 1 {
						t.Fatal("inactive add changed state")
					}
					continue
				}
				target := selected
				if mode == "different" || mode == "equipment" || mode == "full" {
					target = &cells[0]
				}
				if target.Codes[target.Count-1] != 999 || *o.tradeWords[key] != 0 || *o.tradeWords["dword_5d4594_1320944"] != 0 || *o.tradeWords["dword_5d4594_1320948"] != 0 {
					t.Fatalf("side%d active%d mode%s selected-cell routing", side, active, mode)
				}
			}
		}
	}
	tradeUICapture(t, "selected-cell", rows, "91db8c5ae614129ca68b409c59e789a29c3d852dfc058d3866e47ecde189b0cf")
}

func TestClientTradeUIQuantityModifiers(t *testing.T) {
	o := newTradeUIOwner(t)
	mods, free := alloc.Make([]uint32{}, 5)
	defer free()
	for i := 0; i < 4; i++ {
		mods[i] = uint32(txptr(o.mods[i].C()))
	}
	mods[4] = 0x12345678
	for _, present := range []bool{false, true} {
		o.reset(t)
		o.construct(t)
		o.initAmount(t)
		args := o.quantityArgs(200, 100, 32)
		if present {
			args[5] = txptr(unsafe.Pointer(&mods[0]))
		}
		legacy.PortTestTradeUI(7, args...)
		dr := unsafe.Pointer(uintptr(*o.windowWords["nox_gui_itemAmount_item_1319256"]))
		got := unsafe.Slice((*uint32)(unsafe.Add(dr, 432)), 5)
		for i, v := range got {
			want := uint32(0)
			if present {
				want = mods[i]
			}
			if v != want {
				t.Fatalf("present%v modifier%d got%#x want%#x", present, i, v, want)
			}
		}
	}
}
