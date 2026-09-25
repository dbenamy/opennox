//go:build porttest

package opennox

import (
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestClientMetersInstalledTooltip(t *testing.T) {
	o := newMeterOwner(t)
	bow, quiver := o.weapons(t)
	o.plain(t)
	legacy.PortTestMeterCall(16, o.parent, 20, 30, 40, 61)
	w := o.meters.Records[4].Window
	if w == nil || w.TooltipFuncPtr != legacy.PortTestMeterCallback(33) {
		t.Fatal("weapon constructor did not install its tooltip")
	}
	for _, argument := range []uintptr{0, 1, 0x7fffffff, 0x80000000, 0xffffffff} {
		o.equip(bow, quiver, 100)
		w.TooltipFunc(argument)
		if got := alloc.GoString16((*uint16)(memmap.PtrOff(0x5D4594, 1096676))); got != "Bow\nQuiver" {
			t.Fatalf("installed equipped tooltip, argument %#x: %q", argument, got)
		}
		clear(o.equipment)
		w.TooltipFunc(argument)
		if got := alloc.GoString16((*uint16)(memmap.PtrOff(0x5D4594, 1096676))); got != "Current weapon" {
			t.Fatalf("installed empty tooltip, argument %#x: %q", argument, got)
		}
	}
}

func TestGUITooltipForeignBoundary(t *testing.T) {
	o := newMeterOwner(t)
	o.plain(t)
	w := o.parent
	key, words, restore := legacy.PortTestGUITooltipProbe()
	defer restore()
	w.SetTooltipFunc(key)
	for _, argument := range []uintptr{0, 1, 0x7fffffff, 0x80000000, 0xffffffff} {
		*words = [4]uint32{}
		w.TooltipFunc(argument)
		want := [4]uint32{1, uint32(uintptr(w.C())), uint32(uintptr(w.DrawData().C())), uint32(argument)}
		if *words != want {
			t.Fatalf("foreign tooltip words %x, want %x", *words, want)
		}
	}
	*words = [4]uint32{}
	var absent *gui.Window
	absent.TooltipFunc(7)
	w.SetTooltipFunc(nil)
	w.TooltipFunc(8)
	if unsafe.Offsetof(w.TooltipFuncPtr) != 384 {
		t.Fatal("tooltip ABI field moved")
	}
	// This ABI field also admits the dead-word sentinel. Store it as an integer
	// word so no Go write barrier treats the sentinel as a managed pointer.
	field := (*uint32)(unsafe.Add(w.C(), 384))
	*field = 0xacacacac
	w.TooltipFunc(9)
	*field = 0
	w.SetTooltipFunc(key)
	w.Destroy()
	w.TooltipFunc(10)
	if *words != ([4]uint32{}) {
		t.Fatalf("nil, sentinel or destroyed tooltip invoked callback: %x", *words)
	}
}
