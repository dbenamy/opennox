//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/opennox/v1/client/gui"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/internal/netlist"
	"testing"

	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
)

func (o *meterOwner) items(t *testing.T) []*client.Drawable {
	var out []*client.Drawable
	for _, name := range []string{"RedPotion", "BluePotion", "CurePoisonPotion", "RedApple", "Meat"} {
		dr := o.newUnlinked(t)
		o.c.DrawableLinkThing(dr, o.c.Things.TypeByID(name).Index())
		out = append(out, dr)
	}
	return out
}
func (o *meterOwner) slot(index int, dr *client.Drawable, count byte, code uint16) {
	cell := o.grid[index*148:][:148]
	binary.LittleEndian.PutUint32(cell, uint32(uintptr(dr.C())))
	binary.LittleEndian.PutUint32(cell[4:], uint32(code))
	cell[140] = count
}
func TestClientMetersPotionPriorityContract(t *testing.T) {
	o := newMeterOwner(t)
	items := o.items(t)
	o.constructorStart(t, 1, 0)
	legacy.PortTestMeterCall(35, nil, 0, 0, 0, 0)
	for i, dr := range items {
		o.slot(i, dr, byte(i+1), uint16(100+i))
	}
	for _, idx := range []int{0, 4, 3} {
		legacy.PortTestMeterCall(32, nil, 0, 0, 0, 0)
		if got := *memmap.PtrUint32(0x5D4594, 1090308); got != items[idx].TypeIDVal {
			t.Fatalf("health potion/food priority: type %d want %d", got, items[idx].TypeIDVal)
		}
		if got := *memmap.PtrUint16(0x5D4594, 1090312); got != uint16(idx+1) {
			t.Fatalf("count %d", got)
		}
		o.grid[idx*148+140] = 0
	}
	legacy.PortTestMeterCall(32, nil, 0, 0, 0, 0)
	if *memmap.PtrUint32(0x5D4594, 1090296) != 0 || *memmap.PtrUint32(0x5D4594, 1090308) != 0 || *memmap.PtrUint16(0x5D4594, 1090312) != 0 {
		t.Fatal("empty health slot retains potion")
	}
}
func TestClientMetersPotionMatrix(t *testing.T) {
	o := newMeterOwner(t)
	items := o.items(t)
	var out []meterResult
	id := 0
	for mask := 0; mask < 32; mask++ {
		for _, count := range []byte{1, 99, 255} {
			o.constructorStart(t, 1, 0)
			legacy.PortTestMeterCall(35, nil, 0, 0, 0, 0)
			for i, dr := range items {
				if mask&(1<<i) != 0 {
					o.slot(i*17, dr, count, uint16(100+i))
				}
			}
			out = append(out, o.invoke(t, id, 0, 32, nil, 0, 0, 0, 0))
			for i, name := range []string{"dword_5d4594_1090292", "dword_5d4594_1090828", "dword_5d4594_1091364"} {
				o.collect()
				var wIndex int
				for n, w := range o.windows {
					if uint32(uintptr(w.C())) == *o.meters.NamedWord(name) {
						wIndex = n
						break
					}
				}
				w := o.windows[wIndex]
				if int(uintptr(w.WidgetData)) != i {
					t.Fatal("slot window ownership")
				}
				out = append(out, o.invoke(t, id, i+1, 20, w, 0, 0, 0, 0))
				out = append(out, o.invoke(t, id, i+4, 21, w, 5, 0, 0, 0))
			}
			for i, op := range []int{28, 29, 30} {
				out = append(out, o.invoke(t, id, i+7, op, nil, 0, 0, 0, 0))
			}
			id++
		}
	}
	meterCapture(t, "meters-potions", out, len(out), "3c64b313409284743ff5834b5f76834a31ae39b2530fc70b315eac0d49f8f1a3")
}

func TestClientMetersPotionUseContract(t *testing.T) {
	o := newMeterOwner(t)
	items := o.items(t)
	o.constructorStart(t, 1, 0)
	legacy.PortTestMeterCall(35, nil, 0, 0, 0, 0)
	for i, dr := range items {
		o.slot(i, dr, 1, uint16(0x120+i))
	}
	legacy.PortTestMeterCall(32, nil, 0, 0, 0, 0)
	for i, op := range []int{28, 29, 30} {
		for _, cursor := range []gui.Cursor{0, 1, 2} {
			for _, pause := range []bool{false, true} {
				o.c.srv.NetList.ResetAll()
				o.c.Cursor = cursor
				*memmap.PtrUint32(0x5D4594, 1096672) = uint32(cursor)
				noxflags.ResetGame()
				if pause {
					noxflags.SetGame(noxflags.GamePause)
				}
				legacy.PortTestMeterCall(op, nil, 0, 0, 0, 0)
				var messages [][]byte
				o.c.srv.NetList.ByInd(31, netlist.Kind0).Each(func(b []byte) bool { messages = append(messages, append([]byte(nil), b...)); return false })
				if cursor != 0 || pause {
					if len(messages) != 0 {
						t.Fatalf("quick potion op%d ignored cursor%d/pause%v: %x", op, cursor, pause, messages)
					}
				} else if len(messages) != 1 || !bytes.Equal(messages[0], []byte{116, byte(0x20 + i), 1}) {
					t.Fatalf("potion use request %x", messages)
				}
			}
		}
	}
}
