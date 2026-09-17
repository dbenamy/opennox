//go:build porttest

package opennox

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"unsafe"

	"github.com/opennox/libs/ifs"
	"github.com/opennox/libs/spell"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
	"github.com/opennox/opennox/v1/server"
)

func TestServerOptionsMapList(t *testing.T) {
	type row struct {
		Mode      int
		Current   string
		Update    bool
		Names     []string
		Colors    []uint32
		Selection int
		Settings  []byte
	}
	var rows []row
	for _, mode := range []int{0x100, 0x20, 0x40, 0x1000} {
		for _, current := range []string{"", "Arena", "arena", "Zeta", "missing"} {
			for _, update := range []bool{false, true} {
				t.Run(fmt.Sprintf("%x-%s-%t", mode, current, update), func(t *testing.T) {
					catalog := legacy.PortTestMapCatalogOpen(7)
					t.Cleanup(catalog.Close)
					for _, v := range []legacy.PortTestMapCatalogEntry{{Name: "Zeta", Enabled: 1, Flags: 4}, {Name: "Arena", Enabled: 1, Flags: 12}, {Name: "Hidden", Enabled: 0, Flags: 12}, {Name: "Ctf", Enabled: 1, Flags: 8}} {
						catalog.Add(v)
					}
					// Supply nonzero recommended-player metadata on the actual catalog nodes.
					head := legacy.Get_nox_common_maplist()
					for p := *(*unsafe.Pointer)(head); p != head; p = *(*unsafe.Pointer)(p) {
						it := (*legacy.Nox_map_list_item)(p)
						it.Field_8_0 = 2
						it.Field_8_1 = 8
					}
					o := newServerOptionsOwner(t)
					legacy.PortTestServerOptionsModeName(0)
					t.Cleanup(legacy.PortTestServerOptionsRuleState())
					rules, _ := server.PortTestRuleServerSetup()
					oldSpells := o.c.srv.Spells
					o.c.srv.Spells = rules.Spells
					t.Cleanup(func() { o.c.srv.Spells = oldSpells })
					t.Cleanup(handles.PortTestInit())
					dir := t.TempDir()
					oldDir, err := ifs.Workdir()
					if err != nil {
						t.Fatal(err)
					}
					if err := ifs.Chdir(dir); err != nil {
						t.Fatal(err)
					}
					t.Cleanup(func() {
						if err := ifs.Chdir(oldDir); err != nil {
							t.Error(err)
						}
					})
					for _, name := range []string{"Arena", "Zeta", "Ctf"} {
						path := filepath.Join(dir, "maps", name)
						if err := os.MkdirAll(path, 0700); err != nil {
							t.Fatal(err)
						}
						if err := os.WriteFile(filepath.Join(path, name+".rul"), []byte("[COMMON]\n"), 0600); err != nil {
							t.Fatal(err)
						}
						if name == "Arena" {
							if err := os.WriteFile(filepath.Join(path, "user.rul"), []byte("[COMMON]\nset spell SPELL_FIREBALL off\n"), 0600); err != nil {
								t.Fatal(err)
							}
						}
					}
					legacy.PortTestServerOptionsMaps(mode, current, update)
					list := o.options.ChildByID(10114)
					data := (*gui.ScrollListBoxData)(list.WidgetData)
					got := row{Mode: mode, Current: current, Update: update, Selection: o.event(10114, 16404, 0, 0), Settings: append([]byte(nil), o.settings...)}
					items := unsafe.Slice(data.Items, int(data.Count))
					for i := 0; i < int(data.Field_11_0); i++ {
						got.Names = append(got.Names, alloc.GoString16((*uint16)(unsafe.Pointer(uintptr(o.event(10114, 16406, uintptr(i), 0))))))
						got.Colors = append(got.Colors, items[i].Field_129)
					}
					want := []string(nil)
					if mode == 0x100 {
						want = []string{"Arena\t2-8", "Zeta\t2-8"}
					} else if mode == 0x20 {
						want = []string{"Arena\t2-8", "Ctf\t2-8"}
					}
					if fmt.Sprint(got.Names) != fmt.Sprint(want) {
						t.Fatalf("map rows %q want %q", got.Names, want)
					}
					selected := -1
					if len(want) > 0 {
						selected = 0
						if mode == 0x100 && current == "Zeta" {
							selected = 1
						}
					}
					if got.Selection != selected {
						t.Fatalf("selection %d want %d", got.Selection, selected)
					}
					if len(want) > 0 {
						palette := **(**uint32)(memmap.PtrOff(0x85B3FC, 132+4*6))
						if got.Colors[0] != palette {
							t.Fatalf("modified-rule map color %x want %x", got.Colors[0], palette)
						}
					}
					if update {
						name := alloc.GoString(&o.settings[0])
						expected := ""
						if selected >= 0 {
							expected = "Arena"
							if selected == 1 {
								expected = "Zeta"
							}
						}
						if name != expected {
							t.Fatalf("selected map %q want %q", name, expected)
						}
						if expected == "Arena" {
							index := int(spell.SPELL_FIREBALL)
							if o.settings[24+index/8]&(1<<uint(index%8)) != 0 {
								t.Fatal("selected user rule did not disable Fireball")
							}
						}
					}
					rows = append(rows, got)
				})
			}
		}
	}
	spellbookCapture(t, "server-options-map-list", rows, "9a96a79edf9f91051482148afe12cec85c9104e118e5b61fe3c6524759d72fd2")
}
