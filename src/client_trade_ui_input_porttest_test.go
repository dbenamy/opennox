//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
)

func (o *tradeUIOwner) tradeCell(t *testing.T, side, index, count int, name string, base uint32) *client.Drawable {
	dr := o.item(t, name, base)
	cell := &legacy.PortTestTradeUICells()[side][index]
	cell.Drawable, cell.Count, cell.Value = dr, uint32(count), uint32(count*7)
	for i := 0; i < count; i++ {
		cell.Codes[i] = base + uint32(i)
	}
	return dr
}
func tradeUIRemoveRequests(r tradeUIResult) [][]byte {
	var out [][]byte
	for _, m := range r.Window.Inventory.Messages {
		if len(m) == 4 && m[0] == 201 && m[1] == 16 {
			out = append(out, m)
		}
	}
	return out
}
func TestClientTradeUIDragRelease(t *testing.T) {
	o := newTradeUIOwner(t)
	var rows []tradeUIResult
	for _, fps := range []uint32{30, 60} {
		for _, count := range []int{1, 2, 32} {
			for _, origin := range []uint32{0, 1, 2} {
				for _, elapsed := range []uint32{0, fps/3 - 1, fps / 3, fps/3 + 1} {
					for _, end := range [][2]int{{35, 65}, {44, 65}, {45, 65}, {46, 65}, {205, 65}, {310, 210}} {
						o.reset(t)
						o.construct(t)
						o.constructTrade(t)
						o.tradeWindow().Show()
						*o.windowWords["dword_5d4594_1320964"] = 1
						dr := o.tradeCell(t, 0, 0, count, "RedApple", 100)
						pool := o.c.Objs.Alloc.PortTestCounts()
						o.c.srv.SetTickRate(fps)
						o.c.srv.SetFrame(120)
						rows = append(rows, o.tradeCapture(t, len(rows), 11, txptr(o.tradeWindow().C()), 5, inventoryWindowPoint(35, 65)))
						cell := &legacy.PortTestTradeUICells()[0][0]
						if cell.Count != uint32(count-1) || dr.NetCode32 != uint32(100+count-1) || o.c.dragndropItem != dr {
							t.Fatal("trade press did not borrow/remove the last code")
						}
						*memmap.PtrUint32(0x5D4594, 1320304) = origin
						o.c.srv.SetFrame(120 + elapsed)
						r := o.tradeCapture(t, len(rows), 11, txptr(o.tradeWindow().C()), 6, inventoryWindowPoint(end[0], end[1]))
						rows = append(rows, r)
						if cell.Count != uint32(count) || cell.Drawable != dr || cell.Codes[count-1] != uint32(100+count-1) || o.c.dragndropItem != nil || *o.tradeWords["dword_5d4594_1320968"] != 0 || o.c.GUI.Captured() != nil {
							t.Fatal("release did not restore/clear borrowed drag")
						}
						if o.c.Objs.Alloc.PortTestCounts() != pool {
							t.Fatal("trade drag/release changed live allocations")
						}
						dx, dy := end[0]-35, end[1]-65
						click := elapsed <= fps/3 && dx*dx+dy*dy < 100
						want := false
						if origin == 0 {
							inside := end[0] >= 10 && end[0] <= 110 && end[1] >= 40 && end[1] <= 140
							want = !inside || click
						}
						if origin == 1 {
							inside := end[0] >= 180 && end[0] <= 280 && end[1] >= 40 && end[1] <= 140
							want = !inside || click
						}
						messages := tradeUIRemoveRequests(r)
						n := 0
						if want {
							n = 1
						}
						if len(messages) != n {
							t.Fatalf("origin%d elapsed%d end%v requests%d want%d", origin, elapsed, end, len(messages), n)
						}
						if want && binary.LittleEndian.Uint16(messages[0][2:]) != uint16(100+count-1) {
							t.Fatal("remove request has wrong borrowed code")
						}
					}
				}
			}
		}
	}
	tradeUICapture(t, "drag-release", rows, "6479816beda8fab0576a7047e92053615749bc445119924b7fd46fef909ef316")
}
func TestClientTradeUIResetDuringDrag(t *testing.T) {
	for _, count := range []int{1, 2, 32} {
		for _, op := range []int{23, 24} {
			for _, replacement := range []bool{false, true} {
				t.Run(fmt.Sprintf("count%d/op%d/replacement%v", count, op, replacement), func(t *testing.T) {
					o := newTradeUIOwner(t)
					o.reset(t)
					o.construct(t)
					o.constructTrade(t)
					o.tradeWindow().Show()
					*o.windowWords["dword_5d4594_1320964"] = 1
					pool := o.c.Objs.Alloc.PortTestCounts()
					dr := o.tradeCell(t, 0, 0, count, "RedApple", 100)
					legacy.PortTestTradeUI(11, txptr(o.tradeWindow().C()), 5, inventoryWindowPoint(35, 65))
					if o.c.dragndropItem != dr {
						t.Fatal("trade drag did not begin")
					}
					if replacement {
						data := o.record(t, 15)
						tradeUIAddRecord(data, 1, uint16(dr.TypeIDVal), 999, 7)
						if legacy.PortTestTradeUI(27, txptr(unsafe.Pointer(&data[0]))) == 0 {
							t.Fatal("replacement report failed")
						}
					}
					legacy.PortTestTradeUI(op)
					if o.c.dragndropItem != nil || o.c.GUI.Captured() != nil {
						t.Error("reset retained cursor or capture for the ended trade drag")
					}
					if o.objects[o.identities[dr]].Live || o.c.Objs.Alloc.PortTestCounts() != pool {
						t.Error("reset retained the sole dragged drawable")
					}
					before := len(o.events)
					legacy.PortTestTradeUI(op)
					if len(o.events) != before {
						t.Fatal("repeated reset deleted again")
					}
				})
			}
		}
	}
}
func TestClientTradeUIButtonsAndRequests(t *testing.T) {
	o := newTradeUIOwner(t)
	var rows []tradeUIResult
	for _, event := range []uintptr{0, 16390, 16391, 16392} {
		for _, id := range []uint{3708, 3709, 3710, 3711} {
			o.reset(t)
			o.construct(t)
			o.constructTrade(t)
			r := o.tradeCapture(t, len(rows), 14, txptr(o.tradeWindow().C()), event, txptr(o.tradeWindow().ChildByID(id).C()), 0)
			rows = append(rows, r)
			want := byte(0)
			if event == 16391 {
				if id == 3708 {
					want = 17
				}
				if id == 3710 {
					want = 14
				}
			}
			messages := r.Window.Inventory.Messages
			if want == 0 {
				if len(messages) != 0 {
					t.Fatal("unhandled button sent trade request")
				}
			} else if len(messages) != 1 || len(messages[0]) != 2 || messages[0][0] != 201 || messages[0][1] != want {
				t.Fatalf("button%d messages%v", id, messages)
			}
			sounds := r.Window.Inventory.Sounds
			n := 0
			if event == 16391 {
				n = 1
			}
			if len(sounds) != n || n == 1 && sounds[0] != [2]int{766, 100} {
				t.Fatal("button sound contract")
			}
		}
	}
	for _, op := range []int{10, 12, 15} {
		o.reset(t)
		o.construct(t)
		o.constructTrade(t)
		dr := o.item(t, "RedApple", 0x1234)
		r := o.tradeCapture(t, len(rows), op, txptr(dr.C()))
		rows = append(rows, r)
		m := r.Window.Inventory.Messages
		if len(m) != 1 {
			t.Fatal("direct request count")
		}
		if op == 12 {
			if len(m[0]) != 4 || m[0][0] != 201 || m[0][1] != 16 || binary.LittleEndian.Uint16(m[0][2:]) != 0x1234 {
				t.Fatal("remove request encoding")
			}
		}
	}
	tradeUICapture(t, "buttons-requests", rows, "bfd75f4e604d7c826b5a9a41d6ebf14f9e6fb5bf697f33967839568b172f2269")
}
