//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"github.com/opennox/libs/ifs"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
	"github.com/opennox/opennox/v1/server"
	"testing"
)

func serverOptionsRulesOwner(t *testing.T, o *serverOptionsOwner) string {
	t.Helper()
	t.Cleanup(legacy.PortTestServerOptionsRuleState())
	t.Cleanup(handles.PortTestInit())
	rules, _ := server.PortTestRuleServerSetup()
	old := o.c.srv.Spells
	o.c.srv.Spells = rules.Spells
	t.Cleanup(func() { o.c.srv.Spells = old })
	dir := t.TempDir()
	prev, err := ifs.Workdir()
	if err != nil {
		t.Fatal(err)
	}
	if err := ifs.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := ifs.Chdir(prev); err != nil {
			t.Error(err)
		}
	})
	return dir
}
func TestServerOptionsModeSelection(t *testing.T) {
	type row struct {
		Before            uint16
		Selected, Payload int
		After             uint16
		Caption           string
		Hidden, Captured  bool
		Dirty             uint32
		Settings          []byte
	}
	var rows []row
	modes := []uint16{0x20, 0x100, 0x400, 0x10, 0x40, 0x80}
	for _, before := range []uint16{0x100, 0x8021, 0xc401} {
		for _, selected := range []int{-1, 0, 1, 2, 3, 4, 5, 6} {
			for _, payload := range []int{0, 1, 5} {
				t.Run(fmt.Sprintf("%x-%d-%d", before, selected, payload), func(t *testing.T) {
					catalog := legacy.PortTestMapCatalogOpen(17)
					t.Cleanup(catalog.Close)
					o := newServerOptionsOwner(t)
					serverOptionsRulesOwner(t, o)
					binary.LittleEndian.PutUint16(o.settings[52:], before)
					legacy.PortTestServerOptionsEvent(o.options, o.options.ChildByID(10119), 16391, 0)
					list := o.options.ChildByID(10120)
					o.event(10120, 16403, uintptr(selected), 0)
					actual := o.event(10120, 16404, 0, 0)
					// The handler gates on current selection but reads the row from the event payload.
					if legacy.PortTestServerOptionsEvent(o.options, list, 16400, payload) != 1 {
						t.Fatal("selection return")
					}
					valid := actual >= 0 && actual < 6
					want := before
					if valid {
						want = (before & 0xe80f) | modes[payload]
					}
					after := binary.LittleEndian.Uint16(o.settings[52:])
					if after != want || list.Flags.IsHidden() != valid || (o.c.GUI.Captured() == list) == valid || (*o.optionWords["dirty"] != 0) != valid {
						t.Fatalf("selection %d actual %d after %x want %x hidden %t", selected, actual, after, want, list.Flags.IsHidden())
					}
					if valid {
						index := map[uint16]int{0x100: 0, 0x400: 1, 0x20: 2, 0x10: 3, 0x40: 4}[modes[payload]]
						if binary.LittleEndian.Uint16(o.settings[54:]) != uint16(101+11*index) || o.settings[56] != byte(13+7*index) {
							t.Fatal("mode-specific limits were not restored")
						}
					}
					caption := o.options.ChildByID(10119).DrawData().Text()
					if valid && caption != legacy.PortTestServerOptionsModeName(modes[payload]) {
						t.Fatal("selected caption")
					}
					rows = append(rows, row{before, actual, payload, after, caption, list.Flags.IsHidden(), o.c.GUI.Captured() == list, *o.optionWords["dirty"], append([]byte(nil), o.settings...)})
				})
			}
		}
	}
	spellbookCapture(t, "server-options-mode-selection", rows, "b7c9ed21b88915ada194abff56dabf2057478c3c9f9026378fcd599f7ed954df")
}
