//go:build porttest

package opennox

import (
	"math"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestBriefingWindowResourceFailure(t *testing.T) {
	o := newBriefingWindowOwner(t)
	load := legacy.Nox_new_window_from_file
	legacy.Nox_new_window_from_file = func(string, gui.WindowFunc) *gui.Window { return nil }
	t.Cleanup(func() { legacy.Nox_new_window_from_file = load })
	var rows []briefingWindowResult
	value := legacy.PortTestBriefingWindow(0, 0, 0, 0, 0)
	if value != 0 || *o.briefWords["dword_5d4594_831236"] == 0 || *o.briefWords["nox_wnd_briefing_831232"] != 0 {
		t.Fatal("resource failure must leave parent for cleanup")
	}
	rows = append(rows, o.windowCapture(t, 0, value))
	legacy.PortTestBriefingWindow(8, 0, 0, 0, 0)
	if *o.briefWords["dword_5d4594_831236"] != 0 {
		t.Fatal("cleanup must clear surviving parent")
	}
	rows = append(rows, o.windowCapture(t, 8, 0))
	briefingWindowCapture(t, "resource-failure", rows, "bc2787510b727cac7db6dd110c7d622c545ed25abb6782172d2ce117a5b5f6fc")
}
func TestBriefingWindowScrollAndEventHelpers(t *testing.T) {
	o := newBriefingWindowOwner(t)
	var rows []briefingWindowResult
	for _, mode := range []uint32{0, 1, 254, 255, 256, 0xffffffff} {
		o.resetWindow(t)
		*o.briefWords["dword_5d4594_831220"] = mode
		v := legacy.PortTestBriefingWindow(5, 0, 0, 0, 0)
		want := float64(0)
		if mode == 255 {
			want = 1
		}
		if v != math.Float64bits(want) {
			t.Fatal("scroll speed")
		}
		rows = append(rows, o.windowCapture(t, 5, v))
	}
	for _, ev := range []uintptr{0, 1, 17, 18, 22, 23, 24, 0xffffffff} {
		v := legacy.PortTestBriefingWindow(3, 0, ev, 0, 0)
		want := uint64(0)
		if ev == 23 {
			want = 1
		}
		if v != want {
			t.Fatal("background event")
		}
		rows = append(rows, o.windowCapture(t, 3, v))
	}
	v := legacy.PortTestBriefingWindow(6, 0, 0, 0, 0)
	if v != 1 {
		t.Fatal("completed draw")
	}
	rows = append(rows, o.windowCapture(t, 6, v))
	briefingWindowCapture(t, "helpers", rows, "b37e15d030aa36eb1b9a4d85d542ea79e9836a9748e46f5bf392d9d11ce4d1d4")
}
func TestBriefingWindowVoiceCallback(t *testing.T) {
	o := newBriefingWindowOwner(t)
	var rows []briefingWindowResult
	for _, enabled := range []uint32{0, 1} {
		for _, pending := range []uint32{0, 1, 2} {
			o.resetWindow(t)
			o.dialogState[0] = enabled
			*o.briefWords["dword_5d4594_831224"] = pending
			voice := "voice_Warrior_Begin_1"
			*o.briefWords["dword_5d4594_831240"] = uint32(txptr(unsafe.Pointer(alloc.InternCString(voice))))
			v := legacy.PortTestBriefingWindow(1, 0, 0, 0, 0)
			want := uint64(0)
			if pending != 0 {
				want = 1
			}
			if v != want || memmap.Uint32(0x5D4594, 831248) != 1 {
				t.Fatal("voice callback state")
			}
			if (legacy.Dialogs.FileToRead() == voice) != (enabled != 0 && pending != 0) {
				t.Fatal("voice queue")
			}
			if *o.briefWords["dword_5d4594_831244"] != uint32(want) {
				t.Fatal("voice completion flag")
			}
			rows = append(rows, o.windowCapture(t, 1, v))
		}
	}
	briefingWindowCapture(t, "voice", rows, "6fc5c0c7dfcd819bccda0bff30fd032c7540401eaf097c0a0d0b72d8210ad40c")
}
