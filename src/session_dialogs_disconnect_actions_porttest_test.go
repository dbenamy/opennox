//go:build porttest

package opennox

import (
	"image"
	"slices"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy"
)

func TestSessionDisconnectActionsAndDrawing(t *testing.T) {
	o := sessionDialogResources(t)
	words, restore := legacy.PortTestSessionDialogWords()
	defer restore()
	browser, restoreBrowser := legacy.PortTestServerBrowserWords()
	defer restoreBrowser()
	oldOnline := legacy.Get_dword_5d4594_2650652()
	defer legacy.Set_dword_5d4594_2650652(oldOnline)
	legacy.Set_dword_5d4594_2650652(0)
	oldWait, oldLeave := legacy.Sub_43CF40, legacy.Sub_446380
	defer func() { legacy.Sub_43CF40 = oldWait; legacy.Sub_446380 = oldLeave }()
	var calls []string
	legacy.Sub_43CF40 = func() { calls = append(calls, "wait") }
	legacy.Sub_446380 = func() { calls = append(calls, "leave") }
	sessionCall("disconnectOpen", nil, 0, 0, 0)
	defer sessionCall("disconnectClose", nil, 0, 0, 0)
	w := (*gui.Window)(unsafe.Pointer(uintptr(*words["disconnect"])))
	icon := (*gui.Window)(unsafe.Pointer(uintptr(*words["disconnectIcon"])))
	sessionCall("disconnectShow", nil, 0, 1, 0)
	if sessionCall("disconnectEvent", w, 16391, uint32(uintptr(w.ChildByID(576).C())), 0) != 0 || !slices.Equal(calls, []string{"wait"}) || w.Flags.Has(gui.StatusHidden) {
		t.Fatal("wait action", calls)
	}
	sessionCall("disconnectEvent", w, 16391, uint32(uintptr(w.ChildByID(577).C())), 0)
	if !slices.Equal(calls, []string{"wait", "leave"}) || !w.Flags.Has(gui.StatusHidden) || *browser["dword_5d4594_815100"] != 1 {
		t.Fatal("leave action", calls)
	}
	// Compare actual C drawing with the renderer at independently summed positions.
	for _, off := range []image.Point{{0, 0}, {7, 9}, {-5, 3}} {
		icon.Off = image.Pt(25, 30)
		icon.DrawData().ImgPtVal = off
		clear(o.pix.Pix)
		if sessionCall("disconnectDraw", icon, 0, 0, 0) != 1 {
			t.Fatal("draw return")
		}
		got := append([]uint16(nil), o.pix.Pix...)
		clear(o.pix.Pix)
		o.c.R2().DrawImageAt(o.images[0], image.Pt(25+off.X, 30+off.Y))
		if !slices.Equal(got, o.pix.Pix) {
			t.Fatal("icon image offset", off)
		}
		nonzero := false
		for _, v := range got {
			nonzero = nonzero || v != 0
		}
		if !nonzero {
			t.Fatal("empty draw fixture")
		}
	}
}
