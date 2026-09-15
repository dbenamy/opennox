//go:build porttest

package opennox

import (
	"encoding/binary"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
)

func TestClientInventoryNativeStorageContracts(t *testing.T) {
	if got := legacy.PortTestUIInventoryNativeLayout(); got != [6]uintptr{148, 4, 132, 140, 8, 4} {
		t.Fatalf("shared inventory ABI %v", got)
	}
	o := newUIInventoryOwner(t)
	o.reset(t)
	o.stack(83, 32)
	o.grid[83*148+140] = 255
	if legacy.PortTestUIInventoryCall(3, 0x101f, 0, 0) == 0 || memmap.Uint32(0x5D4594, 1049792) != 31 {
		t.Fatal("oversized count lost last stored code")
	}
	if legacy.PortTestUIInventoryCall(22, 0x1000, 0, 0) != 255 {
		t.Fatal("count getter must preserve stored byte")
	}
	if legacy.PortTestUIInventoryCall(3, 0xffffffff, 0, 0) != 0 {
		t.Fatal("oversized count searched outside stored codes")
	}
	o.reset(t)
	o.equipment[0] = uint32(uintptr(o.items[0].C()))
	o.equipment[8] = uint32(uintptr(o.items[1].C()))
	o.items[1].TypeIDVal = o.items[0].TypeIDVal
	if legacy.PortTestUIInventoryCall(1, o.items[0].TypeIDVal, 0, 0) != o.equipment[0] {
		t.Fatal("duplicate equipment type must choose lowest slot")
	}
	o.stack(0, 1)
	binary.LittleEndian.PutUint32(o.grid[4:], 0)
	if legacy.PortTestUIInventoryCall(3, 0, 0, 0) == 0 {
		t.Fatal("stored zero code lookup")
	}
	*(*uint16)(unsafe.Add(o.items[0].C(), 448)) = 123
	if legacy.PortTestUIInventoryCall(30, 0, 55, 66) != 0 || *(*uint16)(unsafe.Add(o.items[0].C(), 448)) != 123 {
		t.Fatal("zero-code charge update must stay ignored")
	}
}
