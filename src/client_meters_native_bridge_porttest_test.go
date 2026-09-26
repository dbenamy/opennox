//go:build porttest

package opennox

import (
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy"
)

func TestGUITooltipNativeBoundary(t *testing.T) {
	o := newMeterOwner(t)
	o.plain(t)
	w := o.parent
	keys := new([2]byte)
	first, second := unsafe.Pointer(&keys[0]), unsafe.Pointer(&keys[1])
	var gotWindow *gui.Window
	var gotDraw *gui.WindowData
	var gotArgument uintptr
	calls := [2]int{}
	callback := func(index int) func(*gui.Window, *gui.WindowData, uintptr) {
		return func(window *gui.Window, draw *gui.WindowData, argument uintptr) {
			calls[index]++
			gotWindow, gotDraw, gotArgument = window, draw, argument
		}
	}
	gui.RegisterTooltipCallbackGo(first, callback(0))
	gui.RegisterTooltipCallbackGo(second, callback(1))
	mustPanic := func(name string, f func()) {
		t.Helper()
		panicked := false
		func() { defer func() { panicked = recover() != nil }(); f() }()
		if !panicked {
			t.Fatalf("%s accepted invalid registration", name)
		}
	}
	mustPanic("nil key", func() { gui.RegisterTooltipCallbackGo(nil, callback(0)) })
	mustPanic("nil function", func() { gui.RegisterTooltipCallbackGo(unsafe.Pointer(new(byte)), nil) })
	mustPanic("duplicate", func() { gui.RegisterTooltipCallbackGo(first, callback(1)) })
	for _, argument := range []uintptr{0, 1, 0x7fffffff, 0x80000000, 0xffffffff} {
		calls = [2]int{}
		w.SetTooltipFunc(first)
		w.TooltipFunc(argument)
		if calls != ([2]int{1, 0}) || gotWindow != w || gotDraw != w.DrawData() || gotArgument != argument {
			t.Fatalf("native arguments/count: %v %p %p %#x", calls, gotWindow, gotDraw, gotArgument)
		}
		w.SetTooltipFunc(second)
		w.TooltipFunc(argument)
		if calls != ([2]int{1, 1}) || gotWindow != w || gotDraw != w.DrawData() || gotArgument != argument {
			t.Fatal("replacement native callback route")
		}
	}
	// Replacing a registered key with a foreign callback must not retain a stale route.
	foreign, words, restore := legacy.PortTestGUITooltipProbe()
	defer restore()
	calls = [2]int{}
	w.SetTooltipFunc(foreign)
	w.TooltipFunc(0xfedcba98)
	if calls != ([2]int{}) || *words != ([4]uint32{1, uint32(uintptr(w.C())), uint32(uintptr(w.DrawData().C())), 0xfedcba98}) {
		t.Fatal("foreign callback after native key replacement", calls, *words)
	}
	w.SetTooltipFunc(first)
	w.Destroy()
	w.TooltipFunc(123)
	if calls != ([2]int{}) {
		t.Fatal("destroyed window invoked registered callback")
	}
}

type meterBridgeOrderedEvent struct {
	code  int
	order []int
}

func (e *meterBridgeOrderedEvent) EventCode() int {
	e.order = append(e.order, 1)
	return e.code
}
func (e *meterBridgeOrderedEvent) EventArgsC() (uintptr, uintptr) {
	e.order = append(e.order, 2)
	e.code = 8 // The callback must use the code read before argument evaluation.
	return 0x89abcdef, 0xfedcba98
}

func TestClientMetersNativeEventOrder(t *testing.T) {
	o := newMeterOwner(t)
	o.plain(t)
	legacy.PortTestMeterCall(16, o.parent, 20, 30, 40, 61)
	w := o.meters.Records[4].Window
	for _, code := range []int{0, 1, 8, 12, 16, 19, -1, 2147483647, -2147483648} {
		event := &meterBridgeOrderedEvent{code: code}
		got := w.Func93(event)
		if len(event.order) != 2 || event.order[0] != 1 || event.order[1] != 2 {
			t.Fatalf("event evaluation order for%d: %v", code, event.order)
		}
		zero := code == 8 || code == 12 || code == 16
		if zero && got != nil || !zero && (got == nil || got.EventRespC() != 1) {
			t.Fatalf("event%d response %v, zero=%t", code, got, zero)
		}
	}
}
