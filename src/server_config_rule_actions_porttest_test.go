//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/spell"
	"github.com/opennox/opennox/v1/client/gui"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"unsafe"
)

func TestServerConfigRuleActions(t *testing.T) {
	o, w, dir := serverConfigRuleOwner(t)
	defer noxflags.PortTestGameFlags(1)()
	for i := 24; i < 52; i++ {
		o.settings[i] = 255
	}
	binary.LittleEndian.PutUint16(o.settings[52:], 0x100)
	legacy.PortTestServerPanelsSpellStore((*uint32)(unsafe.Pointer(&o.settings[24])))
	legacy.PortTestServerPanelsWeaponStore((*uint32)(unsafe.Pointer(&o.settings[44])))
	legacy.PortTestServerPanelsArmorStore(0xffffffff)
	o.options.ChildByID(10119).DrawData().SetText(legacy.PortTestServerOptionsModeName(0x100))
	mapList := o.options.ChildByID(10114)
	mapList.Func94(gui.AsWindowEvent(16397, uintptr(unsafe.Pointer(alloc.InternCString16("arena"))), ^uintptr(0)))
	mapList.Func94(gui.AsWindowEvent(16403, 0, 0))
	o.text(10101, "Test server")
	o.text(10134, "17")
	o.text(10135, "23")
	event := func(id uint) {
		if got := gui.EventRespInt(w.Func94(gui.AsWindowEvent(16391, uintptr(w.ChildByID(id).C()), 0))); got != 1 {
			t.Fatal("rule event result", id, got)
		}
	}
	list, entry := w.ChildByID(10170), w.ChildByID(10171)
	type row struct {
		Action, File, Entry string
		Names               []string
		Selection           int
		Dirty               uint32
		Settings            []byte
	}
	var rows []row
	record := func(action, file string) {
		rows = append(rows, row{action, file, serverPanelsGetText(entry), serverPanelsListNames(list), gui.EventRespInt(list.Func94(gui.AsWindowEvent(16404, 0, 0))), *o.optionWords["dirty"], append([]byte(nil), o.settings...)})
	}
	if err := os.WriteFile(filepath.Join(dir, "maps", "arena", "keep.txt"), []byte("keep"), 0600); err != nil {
		t.Fatal(err)
	}
	serverPanelsSetText(entry, "Custom")
	event(10172)
	path := filepath.Join(dir, "maps", "arena", "Custom.rul")
	first, err := os.ReadFile(path)
	if err != nil || string(first) != "[DEATHMATCH]\n" {
		t.Fatal("save as", err, string(first))
	}
	if !reflect.DeepEqual(serverPanelsListNames(list), []string{"Custom"}) || serverPanelsGetText(entry) != "" || w.ChildByID(10172).Flags.IsEnabled() {
		t.Fatal("save as list/entry/button")
	}
	record("create", string(first))
	serverPanelsSetText(entry, "CUSTOM")
	event(10172)
	if !reflect.DeepEqual(serverPanelsListNames(list), []string{"Custom"}) || serverPanelsGetText(entry) != "" {
		t.Fatal("duplicate row handling")
	}
	record("duplicate", string(first))
	if err := os.WriteFile(path, []byte("replace me\n"), 0600); err != nil {
		t.Fatal(err)
	}
	list.Func94(gui.AsWindowEvent(16403, 0, 0))
	event(10173)
	replaced, err := os.ReadFile(path)
	if err != nil || !bytes.Equal(replaced, first) {
		t.Fatal("overwrite", err, string(replaced))
	}
	record("overwrite", string(replaced))
	loaded := "[DEATHMATCH]\nset spell \"SPELL_FIREBALL\" off\n"
	if err := os.WriteFile(path, []byte(loaded), 0600); err != nil {
		t.Fatal(err)
	}
	*o.optionWords["dirty"] = 0
	list.Func94(gui.AsWindowEvent(16403, 0, 0))
	event(10174)
	index := int(spell.SPELL_FIREBALL)
	want := bytes.Repeat([]byte{255}, 28)
	want[index/8] &^= 1 << uint(index&7)
	if !bytes.Equal(o.settings[24:52], want) || *o.optionWords["dirty"] != 1 {
		t.Fatal("loaded masks and dirty state", o.settings[24:52], want)
	}
	shared := unsafe.Slice((*byte)(unsafe.Pointer(legacy.PortTestServerPanelsSpellPointer())), 20)
	if !bytes.Equal(shared, want[:20]) || *legacy.PortTestServerPanelsWeaponPointer() != 0xffffffff || legacy.PortTestServerPanelsArmorLoad() != 0xffffffff {
		t.Fatal("shared masks")
	}
	record("load", loaded)
	list.Func94(gui.AsWindowEvent(16403, 0, 0))
	event(10175)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("delete selected rule", err)
	}
	if len(serverPanelsListNames(list)) != 0 {
		t.Fatal("deleted row retained")
	}
	kept, err := os.ReadFile(filepath.Join(dir, "maps", "arena", "keep.txt"))
	if err != nil || string(kept) != "keep" {
		t.Fatal("unrelated file changed", err)
	}
	record("delete", "")
	spellbookCapture(t, "server-config-rule-actions", rows, "3ef1252337f23c20b0249894fe12d7bc39844f955bccc2247a12abb6c4626c64")
}
