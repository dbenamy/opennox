//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"image"
	"testing"
	"unsafe"

	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/legacy"
)

func TestGameMessageClientSessionBookData(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	shopBytes := serverConfigOwnBytes(t, 0x5D4594, 1098636, 60*140)
	shopWords, restore := legacy.PortTestInventoryWindowWords()
	t.Cleanup(restore)
	// Own the actual exported inventory grid; its address is supplied by the owner.
	cells := legacy.PortTestInventoryCells()
	actualInventory := unsafe.Slice((*byte)(unsafe.Pointer(&cells[0])), len(cells)*int(unsafe.Sizeof(cells[0])))
	oldInventory := bytes.Clone(actualInventory)
	t.Cleanup(func() { copy(actualInventory, oldInventory) })
	type row struct {
		On, Mask, Marked, Present, Value, Selected int
		Fields                                     [4]uint32
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for mask := 0; mask < 16; mask++ {
			for marked := 0; marked < 2; marked++ {
				for present := 0; present < 2; present++ {
					values := []int{0, 1, 127, 128, 255}
					if mask == 15 {
						values = make([]int, 256)
						for i := range values {
							values[i] = i
						}
					}
					for _, value := range values {
						c.resetCase(env, pix, 1, 100)
						clear(shopBytes)
						clear(actualInventory)
						binary.LittleEndian.PutUint32(connected, uint32(on))
						*shopWords["dword_5d4594_1098624"] = uint32(mask & 1)
						var drs [4]*client.Drawable
						var want [4][]byte
						for i := range drs {
							dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(300+i, 400))
							if dr == nil {
								t.Fatal("allocation")
							}
							drs[i] = dr
							dr.ObjClass = 0
							dr.NetCode32 = uint32(100 + i)
							*txword(dr, 432) = 0xcafe0000 + uint32(i)
						}
						shop := &legacy.PortTestShopUICells()[0]
						shop.Drawable, shop.Count, shop.Codes[0] = drs[0], 1, 17
						if mask&2 != 0 {
							cells[0].Drawable, cells[0].Count, cells[0].Codes[0] = drs[1], 1, 17
						}
						if mask&4 != 0 {
							drs[2].NetCode32 = 17
						}
						drs[3].ObjClass = object.Class(0x20400000)
						if mask&8 != 0 {
							drs[3].NetCode32 = 17
						}
						selected := -1
						if present != 0 {
							if mask&1 != 0 {
								selected = 0
							} else if mask&2 != 0 {
								selected = 1
							} else if marked == 0 && mask&4 != 0 {
								selected = 2
							} else if marked != 0 && mask&8 != 0 {
								selected = 3
							}
						}
						for i, dr := range drs {
							want[i] = bytes.Clone(unsafe.Slice((*byte)(dr.C()), int(unsafe.Sizeof(*dr))))
							if i == selected {
								binary.LittleEndian.PutUint32(want[i][432:], uint32(value))
							}
						}
						code := uint16(18)
						if present != 0 {
							code = 17
						}
						if marked != 0 {
							code |= 0x8000
						}
						data := binary.LittleEndian.AppendUint16([]byte{226}, code)
						data = append(data, byte(value))
						input := bytes.Clone(data)
						n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(226), data)
						if n != 4 || !bytes.Equal(data, input) {
							t.Fatal("book data length/input")
						}
						r := row{On: on, Mask: mask, Marked: marked, Present: present, Value: value, Selected: selected}
						for i, dr := range drs {
							if !bytes.Equal(want[i], unsafe.Slice((*byte)(dr.C()), int(unsafe.Sizeof(*dr)))) {
								t.Fatal("book data lookup priority/isolation", on, mask, marked, present, value, i, selected)
							}
							r.Fields[i] = *txword(dr, 432)
						}
						rows = append(rows, r)
					}
				}
			}
		}
	}
	interactionCapture(t, "game-client-session-book-data", rows)
}
