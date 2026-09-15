//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"strings"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func (o *tradeUIOwner) record(t *testing.T, size int) []byte {
	t.Helper()
	data, free := alloc.Make([]byte{}, size)
	t.Cleanup(free)
	o.c.dataRefs[uint32(txptr(unsafe.Pointer(&data[0])))] = 0xed510001
	return data
}
func TestClientTradeUIAcceptanceFlags(t *testing.T) {
	o := newTradeUIOwner(t)
	var rows []tradeUIResult
	data := o.record(t, 3)
	data[0], data[1] = 201, 3
	for _, active := range []uint32{0, 1, 2} {
		for bits := 0; bits < 256; bits++ {
			o.reset(t)
			o.construct(t)
			o.constructTrade(t)
			*o.windowWords["dword_5d4594_1320964"] = active
			*o.tradeWords["dword_5d4594_1320944"], *o.tradeWords["dword_5d4594_1320948"] = 9, 10
			data[2] = byte(bits)
			rows = append(rows, o.tradeCapture(t, len(rows), 32, txptr(unsafe.Pointer(&data[0]))))
			a, b := uint32(9), uint32(10)
			if active != 0 {
				a, b = uint32(bits)&1, (uint32(bits)>>1)&1
			}
			if *o.tradeWords["dword_5d4594_1320944"] != a || *o.tradeWords["dword_5d4594_1320948"] != b {
				t.Fatalf("active%d bits%d acceptance flags", active, bits)
			}
		}
	}
	tradeUICapture(t, "acceptance-flags", rows, "b9a29c1f6dad74a3c5c9faaa6b4671bb6d8cffd5f24222f7f3b95432d7d66e81")
}
func TestClientTradeUIMoneyText(t *testing.T) {
	o := newTradeUIOwner(t)
	var rows []tradeUIResult
	data := o.record(t, 14)
	data[0], data[1] = 201, 6
	for _, active := range []uint32{0, 1, 2} {
		for _, v := range [][3]uint32{{0, 0, 0}, {1, 2, 3}, {0x7fffffff, 0, 0xffffffff}, {0x80000000, 0xffffffff, 0x7fffffff}, {0xffffffff, 0x80000000, 1}} {
			o.reset(t)
			o.construct(t)
			o.constructTrade(t)
			*o.windowWords["dword_5d4594_1320964"] = active
			for i, n := range v {
				binary.LittleEndian.PutUint32(data[2+4*i:], n)
			}
			rows = append(rows, o.tradeCapture(t, len(rows), 31, txptr(unsafe.Pointer(&data[0]))))
			rows = append(rows, o.tradeCapture(t, len(rows), 16))
			want := [3]string{}
			if active != 0 {
				want[0], want[2] = fmt.Sprint(int32(v[0])), fmt.Sprint(int32(v[2]))
				if v[1] != 0 {
					want[1] = fmt.Sprintf("(%d)", int32(v[1]))
				}
			}
			for i, id := range []uint{3711, 3712, 3713} {
				if got := tradeQuantityText(o.tradeWindow(), id); got != want[i] {
					t.Fatalf("money field%d %q want%q", i, got, want[i])
				}
			}
		}
	}
	tradeUICapture(t, "money-text", rows, "c0272109cb1f5c1b95f6b5fffc45594622e44418a6a2be837183101b7196240c")
}
func TestClientTradeUIStartGatesAndNames(t *testing.T) {
	o := newTradeUIOwner(t)
	var rows []tradeUIResult
	data := o.record(t, 52)
	data[0], data[1] = 201, 12
	for _, active := range []uint32{0, 1, 2} {
		for _, player := range []bool{false, true} {
			for _, name := range []string{"", "Partner", strings.Repeat("X", 24), strings.Repeat("Ω", 24)} {
				o.reset(t)
				o.construct(t)
				o.constructTrade(t)
				*o.windowWords["dword_5d4594_1320964"] = active
				o.players[0].SetName("Local player")
				if !player {
					*o.displayWords["dword_8531A0_2576"] = 0
				}
				alloc.StrCopyZero16(unsafe.Slice((*uint16)(unsafe.Pointer(&data[2])), 25), name)
				r := o.tradeCapture(t, len(rows), 22, txptr(unsafe.Pointer(&data[0])))
				rows = append(rows, r)
				want := uint32(0)
				if player && active != 1 {
					want = 1
				}
				if r.Return != want {
					t.Fatalf("active%d player%v start%d want%d", active, player, r.Return, want)
				}
				if want == 1 {
					if *o.windowWords["dword_5d4594_1320964"] != 1 || o.tradeWindow().Flags.IsHidden() {
						t.Fatal("start did not activate/show trade")
					}
					if tradeQuantityText(o.tradeWindow(), 3702) != "Local player" || tradeQuantityText(o.tradeWindow(), 3703) != name {
						t.Fatal("start did not copy player/partner names")
					}
					rows = append(rows, o.tradeCapture(t, len(rows), 33))
					if tradeQuantityText(o.tradeWindow(), 3703) != name {
						t.Fatal("prepare lost partner name")
					}
				}
			}
		}
	}
	tradeUICapture(t, "start-gates-names", rows, "0e154940e9e33fea80ed9d127c49b23e13fb9138447b6bfd1b69b8ac472d8460")
}
func TestClientTradeUIResetOwnership(t *testing.T) {
	o := newTradeUIOwner(t)
	var rows []tradeUIResult
	for _, op := range []int{23, 24} {
		o.reset(t)
		o.construct(t)
		o.constructTrade(t)
		*o.windowWords["dword_5d4594_1320964"] = 1
		grids := legacy.PortTestTradeUICells()
		var want []uint32
		for side, cells := range grids {
			for i := range cells {
				cells[i].Drawable = o.item(t, "RedApple", uint32(100+side*4+i))
				cells[i].Count = 1
				cells[i].Codes[0] = uint32(100 + side*4 + i)
				cells[i].Value = 7
			}
		}
		for _, cells := range grids {
			for _, i := range []int{0, 2, 1, 3} {
				want = append(want, o.norm(uint32(txptr(cells[i].Drawable.C()))))
			}
		}
		before := len(o.events)
		rows = append(rows, o.tradeCapture(t, len(rows), op))
		for _, id := range []uint{3711, 3712, 3713} {
			if got := tradeQuantityText(o.tradeWindow(), id); got != "Value 0" {
				t.Errorf("reset label %d: %q want Value 0", id, got)
			}
		}
		var got []uint32
		for _, e := range o.events[before:] {
			if e[0] == 2 {
				got = append(got, e[2])
			}
		}
		if len(got) != len(want) {
			t.Fatalf("reset deleted%d want%d", len(got), len(want))
		}
		for i := range got {
			if got[i] != want[i] {
				t.Fatalf("reset deletion%d got%#x want%#x", i, got[i], want[i])
			}
		}
		for side, cells := range grids {
			for i, cell := range cells {
				if cell.Drawable != nil || cell.Count != 0 || cell.Codes[0] != uint32(100+side*4+i) || cell.Value != 7 {
					t.Fatal("reset changed unexpected cell fields")
				}
			}
		}
		before = len(o.events)
		rows = append(rows, o.tradeCapture(t, len(rows), op))
		if len(o.events) != before {
			t.Fatal("repeated reset deleted again")
		}
	}
	tradeUICapture(t, "reset-ownership", rows, "8ee984ccc5c2adb3c77c8d2d83854fb0ca7b97f2f63ab5b00c4253d98eaa200a")
}
