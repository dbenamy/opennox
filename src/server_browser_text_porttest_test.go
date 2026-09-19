//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
	"unsafe"
)

func TestServerBrowserByteText(t *testing.T) {
	o := newListboxOwner(t)
	words, restore := legacy.PortTestServerBrowserWords()
	defer restore()
	o.create(t, 8, 1000, 1000, func(_ *gui.WindowData, d *gui.ScrollListBoxData) { d.Count = 64 })
	*words["dword_5d4594_815004"] = uint32(uintptr(o.win.C()))
	raw, free := alloc.Make([]byte{}, 172)
	defer free()
	for _, sample := range [][]byte{{0x80, 0xff}, {0xc3, 0xa9}, {'A', 0xfe, 'Z'}} {
		clear(raw)
		copy(raw[120:], sample)
		copy(raw[111:], sample)
		legacy.PortTestServerBrowserDetails(unsafe.Pointer(&raw[0]))
		text := serverPanelsListNames(o.win)
		expected := make([]rune, len(sample))
		for i, b := range sample {
			expected[i] = rune(b)
		}
		// %S in the original wide formatter widens bytes individually.
		// Find both displayed name fields among the localized labels.
		count := 0
		for _, s := range text {
			if s == string(expected) {
				count++
			}
		}
		if count != 2 {
			t.Fatalf("byte widening %x: got %#v, want two %q rows", sample, text, string(expected))
		}
	}
}
