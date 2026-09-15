//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy"
)

func TestClientTradeUIAddCapacity(t *testing.T) {
	o := newTradeUIOwner(t)
	var rows []tradeUIResult
	data := o.record(t, 15)
	for side, reportSide := range []byte{1, 0} {
		o.reset(t)
		o.construct(t)
		o.constructTrade(t)
		*o.windowWords["dword_5d4594_1320964"] = 1
		o.tradeCell(t, side, 0, 31, "RedApple", 100)
		cells := legacy.PortTestTradeUICells()[side]
		tradeUIAddRecord(data, reportSide, uint16(cells[0].Drawable.TypeIDVal), 999, 7)
		rows = append(rows, o.tradeCapture(t, len(rows), 27, txptr(unsafe.Pointer(&data[0]))))
		if cells[0].Count != 32 || cells[0].Codes[31] != 999 || cells[0].Value != 224 {
			t.Fatal("31-to-32 trade append")
		}
		full := append([]byte(nil), o.tradeRegions[side][:140]...)
		binary.LittleEndian.PutUint16(data[5:], 1000)
		rows = append(rows, o.tradeCapture(t, len(rows), 27, txptr(unsafe.Pointer(&data[0]))))
		if !bytes.Equal(full, o.tradeRegions[side][:140]) || cells[2].Count != 1 || cells[2].Codes[0] != 1000 || cells[2].Value != 7 {
			t.Fatal("full stack append did not use first empty cell")
		}
		for _, i := range []int{1, 3} {
			o.tradeCell(t, side, i, 32, "RedApple", uint32(2000+100*i))
		}
		cells[2].Count = 32
		before := append([]byte(nil), o.tradeRegions[side]...)
		events := len(o.events)
		r := o.tradeCapture(t, len(rows), 27, txptr(unsafe.Pointer(&data[0])))
		rows = append(rows, r)
		if r.Return != 0 || !bytes.Equal(before, o.tradeRegions[side]) || len(o.events) != events {
			t.Fatal("full trade grid accepted or allocated another item")
		}
	}
	tradeUICapture(t, "add-capacity", rows, "ea40575ba49418799bae43fb292dce7d90cf6426c494312ee9649863cdd4c583")
}
func TestClientTradeUIAddModifiers(t *testing.T) {
	o := newTradeUIOwner(t)
	var rows []tradeUIResult
	data := o.record(t, 15)
	names := make([]*byte, len(o.mods))
	for i, m := range o.mods {
		names[i] = *(**byte)(m.C())
	}
	t.Cleanup(o.c.srv.Server.PortTestControlsModifiers(o.mods, names))
	for side, reportSide := range []byte{1, 0} {
		for _, name := range []string{"RedApple", "Bow"} {
			for _, ids := range [][4]byte{{255, 255, 255, 255}, {1, 2, 3, 4}, {4, 3, 2, 1}, {0, 5, 254, 255}} {
				o.reset(t)
				o.construct(t)
				o.constructTrade(t)
				*o.windowWords["dword_5d4594_1320964"] = 1
				tradeUIAddRecord(data, reportSide, uint16(o.c.Things.TypeByID(name).Index()), 100, 0xffffffff)
				copy(data[11:], ids[:])
				r := o.tradeCapture(t, len(rows), 27, txptr(unsafe.Pointer(&data[0])))
				rows = append(rows, r)
				cell := legacy.PortTestTradeUICells()[side][0]
				if cell.Drawable == nil || cell.Count != 1 || cell.Codes[0] != 100 || cell.Value != 0xffffffff || r.Return != 0xffffffff {
					t.Fatal("trade item metadata/value")
				}
				for i, id := range ids {
					want := uint32(0)
					if name == "Bow" && id >= 1 && id <= 4 {
						want = uint32(txptr(o.mods[id-1].C()))
					}
					if got := *txword(cell.Drawable, 432+4*uintptr(i)); got != want {
						t.Fatalf("%s modifier%d ID%d got%#x want%#x", name, i, id, got, want)
					}
				}
			}
		}
	}
	tradeUICapture(t, "add-modifiers", rows, "9cabfcf5c110e85611014ecd4b590f8dc1f45406e7476a150011bff636e92fd9")
}
func TestClientTradeUIRemoveItems(t *testing.T) {
	o := newTradeUIOwner(t)
	var rows []tradeUIResult
	data := o.record(t, 4)
	data[0], data[1] = 201, 5
	for side := 0; side < 2; side++ {
		for _, count := range []int{1, 2, 3, 32} {
			for _, which := range []int{0, count / 2, count - 1} {
				for _, value := range []uint32{0, 1, uint32(7*count + 1), 0xffffffff} {
					o.reset(t)
					o.construct(t)
					o.constructTrade(t)
					*o.windowWords["dword_5d4594_1320964"] = 1
					dr := o.tradeCell(t, side, 0, count, "RedApple", 100)
					cell := &legacy.PortTestTradeUICells()[side][0]
					cell.Value = value
					binary.LittleEndian.PutUint16(data[2:], uint16(100+which))
					rows = append(rows, o.tradeCapture(t, len(rows), 35, txptr(unsafe.Pointer(&data[0]))))
					if cell.Count != uint32(count-1) || cell.Value != value-value/uint32(count) {
						t.Fatal("trade removal count/value")
					}
					var want []uint32
					for i := 0; i < count; i++ {
						if i != which {
							want = append(want, uint32(100+i))
						}
					}
					for i, code := range want {
						if cell.Codes[i] != code {
							t.Fatal("trade removal reordered surviving codes")
						}
					}
					if cell.Codes[31] != 0 {
						t.Fatal("trade removal did not clear final slot")
					}
					live := o.objects[o.identities[dr]].Live
					if live != (count > 1) || (count == 1 && cell.Drawable != nil) {
						t.Fatal("trade removal drawable ownership")
					}
				}
			}
		}
	}
	tradeUICapture(t, "remove-items", rows, "d4cb79b96e0665e9b36c3845afa92c6c4a899373dc18cae8340021eda5003d14")
}
