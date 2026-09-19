//go:build porttest

package opennox

import (
	"image"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy"
)

func TestSessionDisconnectInitialResourcePosition(t *testing.T) {
	o := sessionDialogResources(t)
	words, restore := legacy.PortTestSessionDialogWords()
	defer restore()
	dim := legacy.PortTestBindingDimensions()
	oldWidth, oldHeight := *dim[0], *dim[1]
	defer func() { *dim[0], *dim[1] = oldWidth, oldHeight }()
	*dim[0], *dim[1] = 640, 480
	load := legacy.Nox_new_window_from_file
	defer func() { legacy.Nox_new_window_from_file = load }()
	var initial image.Point
	legacy.Nox_new_window_from_file = func(name string, fn gui.WindowFunc) *gui.Window {
		w := load(name, fn)
		if name == "discon.wnd" {
			w.SetPos(initial)
		}
		return w
	}
	// C reads offsets24/28 (EndPos), not width/height at8/12. These expectations
	// preserve the legacy placement of customized resources that start off-origin.
	for _, tc := range []struct{ initial, want image.Point }{
		{image.Pt(0, 0), image.Pt(245, 170)},
		{image.Pt(10, 12), image.Pt(240, 164)},
		{image.Pt(-11, 7), image.Pt(251, 167)},
	} {
		initial = tc.initial
		sessionCall("disconnectOpen", nil, 0, 0, 0)
		w := (*gui.Window)(unsafe.Pointer(uintptr(*words["disconnect"])))
		if w.SizeVal != image.Pt(150, 140) {
			t.Fatal("fixture resource size changed", w.SizeVal)
		}
		if w.Off != tc.want {
			t.Errorf("initial %v: got %v want %v", initial, w.Off, tc.want)
		}
		sessionCall("disconnectClose", nil, 0, 0, 0)
		o.c.GUI.FreeDestroyed()
	}
}
