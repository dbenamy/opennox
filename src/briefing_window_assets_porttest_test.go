//go:build porttest

package opennox

import (
	"os"
	"path/filepath"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"golang.org/x/image/font/basicfont"
)

func TestBriefingWindowFadeVoiceIntegration(t *testing.T) {
	o := newBriefingWindowOwner(t)
	var rows []briefingWindowResult
	for _, chapter := range []uintptr{0, 254, 255} {
		o.resetWindow(t)
		o.constructBriefing(t)
		legacy.PortTestBriefingWindow(7, chapter, 1, 2, 0)
		duration := nox_client_getIntroScreenDuration_44E3B0()
		for frame := 0; frame <= duration+1; frame++ {
			o.c.r.DrawFade(true)
			if (memmap.Uint32(0x5D4594, 831248) == 1) != (frame >= duration) {
				t.Fatal("fade callback timing")
			}
			if frame == 0 || frame == duration || frame == duration+1 {
				rows = append(rows, o.windowCapture(t, 9, uint64(frame)))
			}
		}
		if (legacy.Dialogs.FileToRead() != "") != (chapter != 254) {
			t.Fatal("fade must queue chapter/credits voice only")
		}
	}
	briefingWindowCapture(t, "fade-voice", rows, "a3b60c202cb929d8539017723bcbf36eaf9213ecfa82b5d94aff8136cdb07e6b")
}
func TestBriefingWindowShippedLifecycle(t *testing.T) {
	path := os.Getenv("OPENNOX_BRIEFING_ASSETS")
	if path == "" {
		t.Skip("set OPENNOX_BRIEFING_ASSETS to original Nox data")
	}
	o := newBriefingWindowOwner(t)
	font, free := o.c.Render().GetFonts().PortTestWindowFont(basicfont.Face7x13, "large")
	t.Cleanup(free)
	o.c.dataRefs[uint32(txptr(font))] = 0xf0200001
	old := legacy.Nox_new_window_from_file
	t.Cleanup(func() { legacy.Nox_new_window_from_file = old })
	legacy.Nox_new_window_from_file = func(name string, fn gui.WindowFunc) *gui.Window {
		if name != "Briefing.wnd" {
			return old(name, fn)
		}
		f, err := os.Open(filepath.Join(path, "window", name))
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		return newWindowFromReader(o.c.GUI, f, fn)
	}
	var rows []briefingWindowResult
	for _, mode := range []uintptr{1, 2, 4} {
		o.resetWindow(t)
		o.constructBriefing(t)
		legacy.PortTestBriefing(9, 0, 0, 0)
		legacy.PortTestBriefingWindow(7, 254, 1, mode, 0)
		w := o.briefWindow().ChildByID(1010)
		o.c.srv.SetFrame(30)
		w.Draw()
		rows = append(rows, o.windowCapture(t, 4, 1))
		for i := 0; i <= nox_client_getIntroScreenDuration_44E3B0(); i++ {
			o.c.r.DrawFade(true)
		}
		if memmap.Uint32(0x5D4594, 831248) != 1 {
			t.Fatal("shipped fade callback")
		}
		rows = append(rows, o.windowCapture(t, 9, 0))
		parent := (*gui.Window)(unsafe.Pointer(uintptr(*o.briefWords["dword_5d4594_831236"])))
		ret := gui.EventRespInt(parent.Func93(gui.WindowKeyPress{Key: 1, Pressed: true}))
		if ret != 1 || memmap.Uint32(0x5D4594, 832488) != 1 {
			t.Fatal("shipped key dismissal")
		}
		if gui.EventRespInt(parent.Func93(gui.WindowKeyPress{Key: 1, Pressed: true})) != 0 {
			t.Fatal("dismissal must disable input callback")
		}
		rows = append(rows, o.windowCapture(t, 2, uint64(ret)))
	}
	briefingWindowCapture(t, "shipped-lifecycle", rows, "90bc3bace550577c1332484db1a44de7fcfc7c8f1aa935e046c5231a70405c11")
}
