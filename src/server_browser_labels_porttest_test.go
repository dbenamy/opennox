//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
)

func TestServerBrowserLabelLifetime(t *testing.T) {
	o := newEntryOwner(t)
	words, restore := legacy.PortTestServerBrowserWords()
	defer restore()
	*words["nox_wol_wnd_world_814980"] = uint32(uintptr(o.parent.C()))
	for _, name := range []string{"nox_wol_wnd_gameList_815012", "dword_5d4594_814984", "dword_5d4594_814988"} {
		w := o.c.GUI.NewWindowRaw(o.parent, 8, 0, 0, 10, 10, nil)
		*words[name] = uint32(uintptr(w.C()))
	}
	label := o.c.GUI.NewStaticText(o.parent, 10011, 0, 0, 100, 30, false, false, "")
	*words["dword_5d4594_814996"] = uint32(uintptr(label.C()))
	keys := []string{"JoinServer", "CreateMsg", "FilterMsg", "ListJoinServer"}
	entries := make([]strman.Entry, len(keys))
	for i, key := range keys {
		entries[i] = strman.Entry{ID: strman.ID("noxworld.c:" + key), Vals: []strman.Variant{{Str: fmt.Sprintf("Browser label %d: join or choose a server", i)}}}
	}
	set, restoreStrings := o.c.srv.PortTestMeterStrings(entries...)
	defer restoreStrings()
	set(0)
	for round := 0; round < 3; round++ {
		for i, key := range keys {
			legacy.PortTestServerBrowserLabel(key, key == "ListJoinServer")
			// Static-text widgets retain the pointer after dispatch. Reusing temporary
			// allocations must not overwrite the label that will be drawn next frame.
			scratch, free := alloc.Make([]uint16{}, len(entries[i].Vals[0].Str)+1)
			for j := range scratch {
				scratch[j] = 'X'
			}
			scratch[len(scratch)-1] = 0
			got := alloc.GoString16((*gui.StaticTextData)(label.WidgetData).Text)
			free()
			if got != entries[i].Vals[0].Str {
				t.Fatalf("retained %s label: got %q, want %q", key, got, entries[i].Vals[0].Str)
			}
		}
	}
}
