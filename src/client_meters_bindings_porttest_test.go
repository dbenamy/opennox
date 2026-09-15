//go:build porttest

package opennox

import (
	"testing"
	"unsafe"

	"github.com/opennox/libs/client/keybind"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestClientMetersBindingContract(t *testing.T) {
	o := newMeterOwner(t)
	o.plain(t)
	*(*byte)(unsafe.Add(unsafe.Pointer(&o.players[0]), 2251)) = 1
	for i, key := range []keybind.Key{keybind.KeyQ, keybind.KeyF12, keybind.KeyKpEnter} {
		o.c.ctrl.addBinding(&CtrlEventBinding{keys: []keybind.Key{key}, events: []keybind.Event{keybind.Event(36 + i)}})
	}
	legacy.PortTestMeterCall(31, nil, 0, 0, 0, 0)
	for i, want := range []string{"Q", "F12", "Key"} {
		p := (*uint16)(memmap.PtrOff(0x5D4594, 1090300+uintptr(i)*536))
		if got := alloc.GoString16(p); got != want {
			t.Fatalf("binding %d text %q want %q", i, got, want)
		}
		if *(*uint16)(unsafe.Add(unsafe.Pointer(p), 6)) != 0 {
			t.Fatal("binding lacks terminator")
		}
	}
}
