//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
	"unsafe"
)

func TestBindingEditorPromptText(t *testing.T) {
	type record struct {
		Menu           bool
		Title, Visible string
	}
	var records []record
	for _, menu := range []bool{false, true} {
		t.Run(fmt.Sprintf("menu=%v", menu), func(t *testing.T) {
			o := newBindingOwner(t, menu)
			configure, restore := o.c.srv.Server.PortTestMeterStrings(strman.Entry{ID: "InputCfg.wnd:PressKey", Vals: []strman.Variant{{Str: "Press a key"}}})
			t.Cleanup(restore)
			configure(0)
			label := o.c.GUI.NewStaticText(o.modal, 981, 0, 0, 90, 35, false, false, "")
			base, buffer, op := 1321236, 1321256, "sub_4C3CD0"
			if menu {
				base, buffer, op = 1522616, 1522636, "sub_4CBF60"
			}
			var lists [2]*gui.Window
			for i := range lists {
				o.create(t, 8, 180, 80, nil)
				lists[i] = o.win
				o.win = nil
				*o.words[base+4*i] = uint32(uintptr(lists[i].C()))
			}
			buf := memmap.BlobByAddr(0x5D4594).Data[buffer : buffer+512]
			old := append([]byte(nil), buf...)
			t.Cleanup(func() { copy(buf, old) })
			clear(buf)
			// Match constructor initialization with an initially empty dynamic buffer.
			bindingEvent(label, 16385, uintptr(unsafe.Pointer(&buf[0])), 0)
			for _, title := range []string{"Move forward", "Action Ω"} {
				bindingEvent(lists[1], 0x400f, 0, 0)
				bindingText(lists[1], 0x400d, title, -1)
				o.resetRows([2][]string{{"A"}, {"B"}}, 0, 0)
				legacy.PortTestBindingInvoke(op, [4]uint32{uint32(uintptr(lists[0].C())), 16400, uint32(uintptr(o.columns[0].C())), 0})
				want := "Press a key\n'" + title + "'"
				got := legacy.GoWStringP(unsafe.Pointer((*gui.StaticTextData)(label.WidgetData).Text))
				if got != want {
					t.Fatalf("visible prompt %q, want %q; backing buffer %q", got, want, legacy.GoWStringP(unsafe.Pointer(&buf[0])))
				}
				records = append(records, record{menu, title, got})
				o.modal.StackPop()
			}
		})
	}
	spellbookCapture(t, "binding-prompt-text", records, "4b3a6545560ef5bd2bed35b0a904da80755703d165aac00661d891f6f059974a")
}
