//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/cfg"
	"github.com/opennox/libs/client/keybind"
	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestBindingEditorApply(t *testing.T) {
	type record struct {
		Menu   bool
		Mask   int
		Keys   []uint32
		Events []uint32
	}
	var records []record
	for _, menu := range []bool{false, true} {
		t.Run(fmt.Sprintf("menu=%v", menu), func(t *testing.T) {
			o := newBindingOwner(t, menu)
			oldBinding := keyBinding
			t.Cleanup(func() { keyBinding = oldBinding })
			// This fixture exclusively owns global input state. Preserve the lazy cache,
			// including its initialization state, without invoking or sharing its lock.
			dst := reflect.ValueOf(&keybindTitles).Elem()
			saved := reflect.New(dst.Type()).Elem()
			saved.Set(dst)
			dst.SetZero()
			t.Cleanup(func() { dst.Set(saved) })
			var entries []strman.Entry
			for _, k := range keybind.ListKeys() {
				entries = append(entries, strman.Entry{ID: k.TitleID(), Vals: []strman.Variant{{Str: "Key " + k.String()}}})
			}
			events := keybind.New(o.c.Strings()).Events()
			for _, ev := range events {
				entries = append(entries, strman.Entry{ID: strman.ID("bindevent:" + ev.Name), Vals: []strman.Variant{{Str: "Action " + ev.Name}}})
			}
			configure, restore := o.c.srv.Server.PortTestMeterStrings(entries...)
			t.Cleanup(restore)
			configure(0)
			keyBinding = keybind.New(o.c.Strings())
			o.create(t, 8, 180, 80, func(_ *gui.WindowData, d *gui.ScrollListBoxData) { d.Count = 16 })
			titles := o.win
			o.win = nil
			off := 1321240
			if menu {
				off = 1522620
			}
			*o.words[off] = uint32(uintptr(titles.C()))
			keys := []keybind.Key{30, 48, 1, 0x10000}
			for mask := 0; mask < 256; mask++ {
				bindingEvent(titles, 0x400f, 0, 0)
				var rows [2][]string
				var wantKeys, wantEvents []uint32
				for row := 0; row < 4; row++ {
					ev := keyBinding.Events()[row]
					bindingText(titles, 0x400d, ev.Title, -1)
					for col := 0; col < 2; col++ {
						title := " "
						if mask&(1<<uint(row*2+col)) != 0 {
							title = keys[(row+col)%len(keys)].Title(o.c.Strings())
						}
						rows[col] = append(rows[col], title)
					}
					for _, col := range []int{1, 0} {
						if mask&(1<<uint(row*2+col)) != 0 {
							wantKeys = append(wantKeys, uint32(keys[(row+col)%len(keys)]))
							wantEvents = append(wantEvents, uint32(ev.Event))
						}
					}
				}
				wantKeys = append(wantKeys, 1)
				wantEvents = append(wantEvents, uint32(keybind.EventToggleQuitMenu))
				o.resetRows(rows, -1, 0)
				o.c.ctrl.addBinding(&CtrlEventBinding{keys: []keybind.Key{99}, events: []keybind.Event{keybind.EventToggleQuitMenu}})
				legacy.PortTestBindingApply(menu)
				r := record{Menu: menu, Mask: mask}
				for _, b := range o.c.ctrl.listBindings() {
					if len(b.keys) != 1 || len(b.events) != 1 {
						t.Fatalf("unexpected binding %+v", b)
					}
					r.Keys = append(r.Keys, uint32(b.keys[0]))
					r.Events = append(r.Events, uint32(b.events[0]))
				}
				if !reflect.DeepEqual(r.Keys, wantKeys) || !reflect.DeepEqual(r.Events, wantEvents) {
					t.Fatalf("mask=%d got=%+v want keys=%v events=%v", mask, r, wantKeys, wantEvents)
				}
				if mask == 0 || mask == 85 || mask == 170 || mask == 255 {
					file := &cfg.File{Sections: make([]cfg.Section, 2)}
					writeConfigHotkeys(&file.Sections[1])
					path := filepath.Join(t.TempDir(), "nox.cfg")
					f, err := os.Create(path)
					if err != nil {
						t.Fatal(err)
					}
					err = file.WriteTo(f)
					closeErr := f.Close()
					if err != nil {
						t.Fatal(err)
					}
					if closeErr != nil {
						t.Fatal(closeErr)
					}
					f, err = os.Open(path)
					if err != nil {
						t.Fatal(err)
					}
					o.c.ctrl.Reset()
					err = parseLegacyConfig(f, true)
					closeErr = f.Close()
					if err != nil {
						t.Fatal(err)
					}
					if closeErr != nil {
						t.Fatal(closeErr)
					}
					var actualKeys, actualEvents []uint32
					for _, b := range o.c.ctrl.listBindings() {
						if len(b.keys) != 1 || len(b.events) != 1 {
							t.Fatalf("round-trip binding %+v", b)
						}
						actualKeys = append(actualKeys, uint32(b.keys[0]))
						actualEvents = append(actualEvents, uint32(b.events[0]))
					}
					// Configuration sections keep the first key position and the last
					// value for duplicate keys, including the mandatory Escape override.
					var savedKeys, savedEvents []uint32
					positions := make(map[uint32]int)
					for i, k := range wantKeys {
						if n, ok := positions[k]; ok {
							savedEvents[n] = wantEvents[i]
						} else {
							positions[k] = len(savedKeys)
							savedKeys = append(savedKeys, k)
							savedEvents = append(savedEvents, wantEvents[i])
						}
					}
					if !reflect.DeepEqual(actualKeys, savedKeys) || !reflect.DeepEqual(actualEvents, savedEvents) {
						t.Fatalf("configuration round-trip mask %d: keys %v / %v events %v / %v", mask, actualKeys, savedKeys, actualEvents, savedEvents)
					}
				}
				records = append(records, r)
			}
		})
	}
	spellbookCapture(t, "binding-apply", records, "45353301e424f93c06c27a77b128c3cdfe6888694c4de294e5adae4762bc3020")
}
