//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"testing"
)

func TestClientInventoryWindowDrawAnimation(t *testing.T) {
	o := newInventoryWindowOwner(t)
	var rows []inventoryWindowResult
	for _, mode := range []uint32{0, 5} {
		for _, cfg := range [][2]int{{0, -225}, {1, -225}, {1, -161}, {1, -1}, {1, 0}, {2, 0}, {3, 0}, {3, -193}, {3, -224}} {
			o.reset(t)
			o.construct(t)
			*o.windowWords["dword_5d4594_1049864"] = mode
			*memmap.PtrUint8(0x5D4594, 1049870) = 1 // real stats owner; paper-doll composition is separate
			*memmap.PtrUint8(0x5D4594, 1049868) = byte(cfg[0])
			*o.windowWords["dword_587000_136184"] = uint32(int32(cfg[1]))
			for step := 0; step < 5; step++ {
				rows = append(rows, o.capture(t, len(rows), step, 5, txptr(o.mainWindow().C()), 0, 0, 0))
			}
			top := int32(*o.windowWords["dword_587000_136184"])
			if top < -225 || top > 0 {
				t.Fatalf("animation out of bounds: %d", top)
			}
		}
	}
	inventoryWindowCapture(t, "draw-animation", rows, "dda84f2662a3fc107c31f85cda9cdd1fee7cf76d83d71f525cc04fdc7a2cc402")
}
func TestClientInventoryWindowJournalDraw(t *testing.T) {
	o := newInventoryWindowOwner(t)
	var rows []inventoryWindowResult
	node, free := alloc.New(server.PlayerJournal{})
	defer free()
	copy(node.EntryBuf[:], "PortInventoryEntry")
	for _, kind := range []uint16{1, 2, 4, 8} {
		for _, scroll := range []uint32{0, 1, 13, 100, 400} {
			o.reset(t)
			o.construct(t)
			node.Field3 = kind
			o.players[0].Journal = node
			*memmap.PtrUint8(0x5D4594, 1049868) = 2
			*memmap.PtrUint8(0x5D4594, 1049869) = 1
			*memmap.PtrUint8(0x5D4594, 1049870) = 1
			*o.windowWords["dword_587000_136184"] = 0
			*o.windowWords["dword_5d4594_1062512"] = scroll
			rows = append(rows, o.capture(t, len(rows), 0, 5, txptr(o.mainWindow().C()), 0, 0, 0))
			o.players[0].Journal = nil
		}
	}
	inventoryWindowCapture(t, "journal-draw", rows, "99fa2f4dbe4d2e3c0bec3c8d69e270e787593ce4b2bed9ca2dd86daafea4272b")
}
