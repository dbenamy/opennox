//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"github.com/opennox/libs/spell"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"os"
	"path/filepath"
	"testing"
	"unsafe"
)

func TestServerOptionsMapSelectionReset(t *testing.T) {
	type row struct {
		Mode       int
		Missing    bool
		Selection  int
		Name, Text string
		Color      uint32
		Dirty      uint32
		Settings   []byte
	}
	var rows []row
	for _, mode := range []int{0x100, 0x20} {
		for _, missing := range []bool{false, true} {
			t.Run(fmt.Sprintf("%x-%t", mode, missing), func(t *testing.T) {
				catalog := legacy.PortTestMapCatalogOpen(19)
				t.Cleanup(catalog.Close)
				catalog.Add(legacy.PortTestMapCatalogEntry{Name: "Arena", Enabled: 1, Flags: 12})
				catalog.Add(legacy.PortTestMapCatalogEntry{Name: "Zeta", Enabled: 1, Flags: 12})
				o := newServerOptionsOwner(t)
				dir := serverOptionsRulesOwner(t, o)
				binary.LittleEndian.PutUint16(o.settings[52:], uint16(mode))
				for _, name := range []string{"Arena", "Zeta"} {
					p := filepath.Join(dir, "maps", name)
					if err := os.MkdirAll(p, 0700); err != nil {
						t.Fatal(err)
					}
					if err := os.WriteFile(filepath.Join(p, name+".rul"), []byte("[COMMON]\nset spell SPELL_FIREBALL on\n"), 0600); err != nil {
						t.Fatal(err)
					}
				}
				user := filepath.Join(dir, "maps", "Arena", "user.rul")
				if !missing {
					if err := os.WriteFile(user, []byte("[COMMON]\nset spell SPELL_FIREBALL off\n"), 0600); err != nil {
						t.Fatal(err)
					}
				}
				legacy.PortTestServerOptionsMaps(mode, "Zeta", true)
				list := o.options.ChildByID(10114)
				// Use a different valid current selection to prove the event payload selects the map.
				if legacy.PortTestServerOptionsEvent(o.options, list, 16400, 0) != 1 {
					t.Fatal("map selection return")
				}
				if alloc.GoString(&o.settings[0]) != "Arena" {
					t.Fatal("map event did not read payload")
				}
				index := int(spell.SPELL_FIREBALL)
				if (o.settings[24+index/8]&(1<<uint(index%8)) != 0) != missing {
					t.Fatal("map rule application")
				}
				o.event(10114, 16403, 0, 0)
				before := alloc.GoString16((*uint16)(unsafe.Pointer(uintptr(o.event(10114, 16406, 0, 0)))))
				if legacy.PortTestServerOptionsEvent(o.options, o.options.ChildByID(10141), 16391, 0) != 1 {
					t.Fatal("reset button")
				}
				if _, err := os.Stat(user); !os.IsNotExist(err) {
					t.Fatalf("reset must remove user rule file: %v", err)
				}
				if o.settings[24+index/8]&(1<<uint(index%8)) == 0 {
					t.Fatal("reset did not restore base rule")
				}
				data := (*gui.ScrollListBoxData)(list.WidgetData)
				after := alloc.GoString16((*uint16)(unsafe.Pointer(uintptr(o.event(10114, 16406, 0, 0)))))
				selection := o.event(10114, 16404, 0, 0)
				if before != after || selection != 0 || data.Field_11_0 != 2 || *o.optionWords["dirty"] != 1 {
					t.Fatal("reset row/selection/count/dirty")
				}
				items := unsafe.Slice(data.Items, int(data.Count))
				rows = append(rows, row{mode, missing, selection, alloc.GoString(&o.settings[0]), after, items[0].Field_129, *o.optionWords["dirty"], append([]byte(nil), o.settings...)})
			})
		}
	}
	spellbookCapture(t, "server-options-map-selection-reset", rows, "de2db54764668fe602a5e582b1caae53ab8b56953f2b8b2f5a1179a39b14b526")
}
