//go:build porttest

package opennox

import (
	"strings"
	"testing"
	"unsafe"

	"github.com/opennox/libs/client/keybind"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestClientMetersNativeStorageContracts(t *testing.T) {
	o := newMeterOwner(t)
	bow, quiver := o.weapons(t)
	o.plain(t)
	o.equip(bow, quiver, 100)
	name, freeName := alloc.CString16(strings.Repeat("W", 900))
	t.Cleanup(freeName)
	o.weapon.Desc8 = name
	*memmap.PtrUint32(0x5D4594, 1092992) = 0x1234abcd
	legacy.PortTestMeterCall(33, nil, 0, 0, 0, 0)
	text := alloc.GoString16((*uint16)(memmap.PtrOff(0x5D4594, 1091968)))
	if text != strings.Repeat("W", 511) {
		t.Fatalf("bounded combined tooltip length %d", len(text))
	}
	if *memmap.PtrUint32(0x5D4594, 1092992) != 0x1234abcd {
		t.Fatal("long weapon tooltip changed adjacent meter color")
	}
	if got := alloc.GoString16((*uint16)(memmap.PtrOff(0x5D4594, 1096676))); got != strings.Repeat("W", 255) {
		t.Fatal("cursor tooltip bound")
	}
	label := unsafe.Slice((*uint16)(memmap.PtrOff(0x5D4594, 1090300)), 4)
	copy(label, []uint16{'a', 'b', 'c', 'd'})
	o.c.ctrl.addBinding(&CtrlEventBinding{keys: []keybind.Key{keybind.KeyQ}, events: []keybind.Event{36}})
	legacy.PortTestMeterCall(31, nil, 0, 0, 0, 0)
	if label[0] != 'Q' || label[1] != 0 || label[2] != 0 || label[3] != 0 {
		t.Fatalf("short binding must pad and terminate: %v", label)
	}
}
