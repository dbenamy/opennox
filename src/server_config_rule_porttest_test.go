//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"unsafe"
)

func serverConfigRuleOwner(t *testing.T) (*serverOptionsOwner, *gui.Window, string) {
	t.Helper()
	o := newServerOptionsOwner(t)
	o.installSubpanels(t, true)
	dir := serverOptionsRulesOwner(t, o)
	loader := legacy.Nox_new_window_from_file
	legacy.Nox_new_window_from_file = func(name string, proc gui.WindowFunc) *gui.Window {
		if name == "rulelist.wnd" {
			resource := strings.Replace(serverOptionsSubpanelResource(name, true), "DATA = 64 -1 0 0;", "DATA = 9 130 0 2;", 1)
			return newWindowFromString(o.c.GUI, resource, proc)
		}
		return loader(name, proc)
	}
	t.Cleanup(func() { legacy.Nox_new_window_from_file = loader })
	for _, off := range []uintptr{191592, 191608, 191624, 191640} {
		b := serverConfigOwnBytes(t, 0x587000, off, 8)
		clear(b)
		copy(b, ".rul")
	}
	clear(serverConfigOwnBytes(t, 0x5D4594, 1523052, 8))
	copy(o.settings, "arena")
	if err := os.MkdirAll(filepath.Join(dir, "maps", "arena"), 0700); err != nil {
		t.Fatal(err)
	}
	raw := legacy.PortTestServerConfigRuleOpen(o.options, unsafe.Pointer(&o.settings[0]))
	if raw == 0 {
		t.Fatal("rule root")
	}
	return o, (*gui.Window)(unsafe.Pointer(uintptr(uint32(raw)))), dir
}
func TestServerConfigRuleEnumeration(t *testing.T) {
	o, w, dir := serverConfigRuleOwner(t)
	list := w.ChildByID(10170)
	if result := legacy.PortTestServerConfigRulePopulate(unsafe.Pointer(&o.settings[0])); result != -1 {
		t.Fatal("empty enumeration return", result)
	}
	if names := serverPanelsListNames(list); len(names) != 0 {
		t.Fatal("empty directory", names)
	}
	for _, name := range []string{"arena.rul", "user.rul", "ArEnA.rul", "USER.rul", "easy.rul", "hard.rul", "café.rul", "notes.txt", strings.Repeat("z", 100) + ".rul"} {
		if err := os.WriteFile(filepath.Join(dir, "maps", "arena", name), []byte("fixture\n"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Mkdir(filepath.Join(dir, "maps", "arena", "folder.rul"), 0700); err != nil {
		t.Fatal(err)
	}
	result := legacy.PortTestServerConfigRulePopulate(unsafe.Pointer(&o.settings[0]))
	names := serverPanelsListNames(list)
	if result != 1 || !reflect.DeepEqual(names, []string{"caf\u00c3\u00a9", "easy", "hard", strings.Repeat("z", 100)}) {
		t.Fatal("rule enumeration should use file basenames and exclude reserved names", result, names)
	}
	if w.ID() != 10169 || list.DrawData().Window != w {
		t.Fatal("root/list ownership")
	}
	for _, id := range []uint{10177, 10178, 10179} {
		if w.ChildByID(id).DrawData().Window != list {
			t.Fatal("scroll owner", id)
		}
	}
	if got := w.ChildByID(10179).Field100Ptr.SizeVal; got.X != 16 || got.Y != 10 {
		t.Fatal("slider thumb", got)
	}
	spellbookCapture(t, "server-config-rule-enumeration", []any{result, names}, "2c09a0db42f53a832d87d24b5591391cccc28f90e35bb02ff1826c9f4d71b4f7")
}
