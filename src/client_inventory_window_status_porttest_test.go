//go:build porttest

package opennox

import (
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"testing"
)

func TestClientInventoryWindowStatusDrawing(t *testing.T) {
	o := newInventoryWindowOwner(t)
	var rows []inventoryWindowResult
	for _, mask := range []uint32{0, 1, 2, 1 << 29, 1 << 30, 1 << 31, 0xffffffff} {
		for _, extra := range []byte{0, 1, 2, 32, 63} {
			o.reset(t)
			o.construct(t)
			*memmap.PtrUint32(0x5D4594, 1062540) = mask
			*memmap.PtrUint8(0x5D4594, 1062536) = extra
			rows = append(rows, o.capture(t, len(rows), 0, 5, txptr(o.mainWindow().C()), 0, 0, 0))
			for _, x := range []int{0, 39, 40, 74, 75, 249, 500, 65535} {
				rows = append(rows, o.capture(t, len(rows), 1, 4, txptr(o.mainWindow().C()), 0, inventoryWindowPoint(x, 20), 0))
			}
		}
	}
	inventoryWindowCapture(t, "status-icons", rows, "9baa05b7d009349399e10ea1b9e00663931160098e80bf25b4837a650fb562a1")
}
func TestClientInventoryWindowQuestDrawing(t *testing.T) {
	o := newInventoryWindowOwner(t)
	var rows []inventoryWindowResult
	for _, frame := range []uint32{0, 1, 3, 4, 0xffffffff} {
		for _, keys := range []uint32{0, 1, 2} {
			o.reset(t)
			o.construct(t)
			noxflags.SetGame(4096)
			o.c.srv.SetFrame(frame)
			*o.windowWords["dword_5d4594_1319056"] = keys
			*memmap.PtrUint32(0x5D4594, 1050012) = frame
			rows = append(rows, o.capture(t, len(rows), 0, 5, txptr(o.mainWindow().C()), 0, 0, 0))
			for _, xy := range [][2]int{{0, 0}, {0, 10}, {10, 10}, {100, 10}, {500, 300}} {
				rows = append(rows, o.capture(t, len(rows), 1, 4, txptr(o.mainWindow().C()), 0, inventoryWindowPoint(xy[0], xy[1]), 0))
			}
		}
	}
	inventoryWindowCapture(t, "quest-icons", rows, "8543f19763833e24dded95a04e3b55b3f4bab56c452216694025b1c38ae08200")
}
func TestClientInventoryWindowAmountConstruction(t *testing.T) {
	o := newInventoryWindowOwner(t)
	o.reset(t)
	o.construct(t)
	o.initAmount(t)
	rows := []inventoryWindowResult{o.capture(t, 0, 0, 7, 0, 0, 0, 0)}
	inventoryWindowCapture(t, "amount-construction", rows, "e11ba507b15a98007dca8fe500c37bacbf11f1e544008f1a60ea041dd347d962")
}
