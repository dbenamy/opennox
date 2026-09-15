//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
	"testing"
	"unsafe"
)

func TestClientInventoryWindowOpenClose(t *testing.T) {
	o := newInventoryWindowOwner(t)
	var rows []inventoryWindowResult
	for state := 0; state < 4; state++ {
		for _, mode := range []uint32{0, 5, 6} {
			o.reset(t)
			o.construct(t)
			*memmap.PtrUint8(0x5D4594, 1049868) = byte(state)
			*o.windowWords["dword_5d4594_1049864"] = mode
			*o.windowWords["dword_5d4594_1062516"] = 125
			for step, op := range []int{31, 28, 31, 29, 31, 30, 31, 32} {
				rows = append(rows, o.capture(t, state*10+int(mode), step, op, 0, 0, 0, 0))
			}
		}
	}
	inventoryWindowCapture(t, "open-close", rows, "6d4a6ae86b1433099fa8326b949a58c73ec0576f24280a9ae2b536e263ef417a")
}
func TestClientInventoryWindowIdentificationModes(t *testing.T) {
	o := newInventoryWindowOwner(t)
	var rows []inventoryWindowResult
	for _, state := range []byte{0, 1, 2, 3} {
		o.reset(t)
		o.construct(t)
		*memmap.PtrUint8(0x5D4594, 1049868) = state
		for step, op := range []int{12, 0, 0, 26, 29, 30} {
			rows = append(rows, o.capture(t, int(state), step, op, 0, 0, 0, 0))
		}
	}
	inventoryWindowCapture(t, "identify-modes", rows, "0336d6867390693effad63afd9a1d08b2b05f461064191536ed76d46b61e84ce")
}
func TestClientInventoryWindowScrolling(t *testing.T) {
	o := newInventoryWindowOwner(t)
	var rows []inventoryWindowResult
	for _, scroll := range []uint32{0, 1, 24, 25, 49, 50, 849, 850} {
		o.reset(t)
		o.construct(t)
		*o.windowWords["dword_5d4594_1062512"] = scroll
		main := o.mainWindow()
		up, down := main.ChildByID(9102), main.ChildByID(9103)
		for step, control := range []*gui.Window{up, down, down, up} {
			rows = append(rows, o.capture(t, int(scroll), step, 17, txptr(main.C()), 16391, txptr(control.C()), 0))
		}
		slider := (*gui.Window)(unsafe.Pointer(uintptr(*o.windowWords["dword_5d4594_1062508"])))
		for _, value := range []uintptr{0, 1, 425, 850} {
			rows = append(rows, o.capture(t, int(scroll), len(rows), 17, txptr(main.C()), 16393, txptr(slider.C()), value))
		}
	}
	inventoryWindowCapture(t, "scroll", rows, "4df4ade1140e70d54abcf8e09a1e1fc2a5a5f43bec93966aa71f777d711b8f52")
}
func TestClientInventoryWindowPanelSwitches(t *testing.T) {
	o := newInventoryWindowOwner(t)
	var rows []inventoryWindowResult
	for _, mode := range []uint32{0, 5, 6} {
		for _, height := range []uint32{0, 149, 150, 151, 999} {
			o.reset(t)
			o.construct(t)
			*o.windowWords["dword_5d4594_1049864"] = mode
			*memmap.PtrUint32(0x5D4594, 1064848) = height
			*o.windowWords["dword_5d4594_1062512"] = 125
			*o.windowWords["dword_5d4594_1062520"] = 75
			for step, id := range []uint{9105, 9106, 9107, 9108, 9111} {
				control := o.mainWindow().ChildByID(id)
				// Identify mode refuses the transition; use the actual button with its current ID.
				if control == nil {
					if id == 9106 {
						control = o.mainWindow().ChildByID(9105)
					} else if id == 9108 {
						control = o.mainWindow().ChildByID(9107)
					}
				}
				if control == nil {
					t.Fatalf("mode%d missing control%d", mode, id)
				}
				rows = append(rows, o.capture(t, int(mode*1000+height), step, 17, txptr(o.mainWindow().C()), 16391, txptr(control.C()), 0))
			}
		}
	}
	inventoryWindowCapture(t, "panels", rows, "d45ad10e13899b825f8079cbaf1eba3aec998b99565c4c6090864d04b7211bc5")
}
