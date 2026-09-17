//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/client/gui"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
	"unsafe"
)

func serverOptionsChildren(t *testing.T, w *gui.Window) []uint {
	t.Helper()
	var ids []uint
	seen := map[*gui.Window]bool{}
	for child := w.Field100Ptr; child != nil; child = child.Prev() {
		if seen[child] || child.Parent() != w {
			t.Fatal("invalid child links")
		}
		seen[child] = true
		ids = append(ids, child.ID())
		if len(ids) > 256 {
			t.Fatal("child traversal bound")
		}
	}
	return ids
}
func TestServerOptionsTabOrder(t *testing.T) {
	o := newServerOptionsOwner(t)
	type row struct {
		ID       int
		Children []uint
	}
	var rows []row
	count := len(serverOptionsChildren(t, o.options))
	for _, id := range []int{10161, 10162, 10163, 10163, 10161, 10162} {
		o.call("tab-order", id, "")
		ids := serverOptionsChildren(t, o.options)
		if len(ids) != count || ids[0] != uint(id) {
			t.Fatalf("tab %d did not move to front: %v", id, ids)
		}
		rows = append(rows, row{id, ids})
	}
	spellbookCapture(t, "server-options-tab-order", rows, "f3eac40a80ab835eeb8212f9da25e8f27c74ed744f9f9dd898ea0968e6f75959")
}
func TestServerOptionsTabLifecycle(t *testing.T) {
	type row struct {
		Host                     bool
		Tab                      int
		Loaded                   []string
		General, Access, Players bool
		Children                 []uint
	}
	var rows []row
	for _, host := range []bool{false, true} {
		t.Run(fmt.Sprint(host), func(t *testing.T) {
			catalog := legacy.PortTestMapCatalogOpen(11)
			t.Cleanup(catalog.Close)
			o := newServerOptionsOwner(t)
			flags := noxflags.GameFlag(0)
			if host {
				flags = 1
			}
			defer noxflags.PortTestGameFlags(flags)()
			loads := o.installSubpanels(t)
			*o.optionWords["tabs2"] = uint32(uintptr(o.images[1].C()))
			*o.optionWords["tabs3"] = uint32(uintptr(o.images[2].C()))
			for _, tab := range []int{0, 2, 1, 0, 1, 2, 0} {
				o.call("tab", tab, "")
				general := *o.optionWords["advanced-open"] != 0
				access := *o.optionWords["access"] != 0
				players := *o.optionWords["players"] != 0
				if general != (tab == 0) || access != (tab == 1) || players != (tab == 2) {
					t.Fatalf("panel states: %t/%t/%t", general, access, players)
				}
				var root *gui.Window
				switch tab {
				case 0:
					root = (*gui.Window)(unsafe.Pointer(uintptr(*o.optionWords["advanced-open"])))
				case 1:
					root = (*gui.Window)(unsafe.Pointer(uintptr(*o.optionWords["access"])))
				case 2:
					root = (*gui.Window)(unsafe.Pointer(uintptr(*o.optionWords["players"])))
				}
				if root == nil || root.Parent() != o.options {
					t.Fatal("subpanel parent")
				}
				ids := serverOptionsChildren(t, o.options)
				rows = append(rows, row{host, tab, append([]string(nil), (*loads)...), general, access, players, ids})
				o.c.GUI.FreeDestroyed()
			}
			o.call("close", 0, "")
			for _, key := range []string{"root", "access", "players", "advanced-open", "panel-1309812", "panel-1045516"} {
				if *o.optionWords[key] != 0 {
					t.Fatalf("%s survived close", key)
				}
			}
		})
	}
	spellbookCapture(t, "server-options-tab-lifecycle", rows, "70f05b1f423d74871200d1cbf5da649a4da73564a6e013634f13af324cdfb8a9")
}
