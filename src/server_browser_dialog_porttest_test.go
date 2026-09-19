//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"reflect"
	"testing"
	"unsafe"
)

func TestServerBrowserConnectionDialogs(t *testing.T) {
	o := newEntryOwner(t)
	words, restore := legacy.PortTestServerBrowserWords()
	defer restore()
	table := serverConfigOwnBytes(t, 0x587000, 87416, 44)
	keys := []string{"ConnError", "Password", "PasswordRequired", "Connected"}
	for i := 0; i < 11; i++ {
		key := fmt.Sprintf("BrowserError%d", i)
		keys = append(keys, key)
		p, free := alloc.CString(key)
		defer free()
		binary.LittleEndian.PutUint32(table[i*4:], uint32(uintptr(unsafe.Pointer(p))))
	}
	entries := make([]strman.Entry, len(keys))
	for i, key := range keys {
		entries[i] = strman.Entry{ID: strman.ID("noxworld.c:" + key), Vals: []strman.Variant{{Str: key}}}
	}
	set, restoreStrings := o.c.srv.PortTestMeterStrings(entries...)
	defer restoreStrings()
	set(0)
	exists, title, message, buttons, enable, back, create := legacy.Sub_44A4A0, legacy.Sub_449E00, legacy.Sub_449E30, legacy.Sub_449EA0, legacy.Sub_44A360, legacy.Sub_4A24C0, legacy.Nox_xxx_dialogMsgBoxCreate_449A10
	draw, drawHook := o.c.DrawFunc, legacy.Nox_xxx_cliDrawConnectedLoop_43B360
	defer func() {
		legacy.Sub_44A4A0 = exists
		legacy.Sub_449E00 = title
		legacy.Sub_449E30 = message
		legacy.Sub_449EA0 = buttons
		legacy.Sub_44A360 = enable
		legacy.Sub_4A24C0 = back
		legacy.Nox_xxx_dialogMsgBoxCreate_449A10 = create
		o.c.DrawFunc = draw
		legacy.Nox_xxx_cliDrawConnectedLoop_43B360 = drawHook
	}()
	var calls []string
	legacy.Sub_449E00 = func(s string) int { calls = append(calls, "title:"+s); return 0 }
	legacy.Sub_449E30 = func(s string) int { calls = append(calls, "text:"+s); return 0 }
	legacy.Sub_449EA0 = func(f gui.DialogFlags) { calls = append(calls, fmt.Sprint("buttons:", int(f))) }
	legacy.Sub_44A360 = func(v int) { calls = append(calls, fmt.Sprint("enable:", v)) }
	legacy.Sub_4A24C0 = func(v int) int { calls = append(calls, fmt.Sprint("back:", v)); return 1 }
	legacy.Nox_xxx_dialogMsgBoxCreate_449A10 = func(w *gui.Window, title, text string, flags gui.DialogFlags, a, b func()) {
		if w != o.parent || title != "" || text != "" || flags != 0 || a != nil || b != nil {
			t.Fatal("dialog construction", w, title, text, flags)
		}
		calls = append(calls, "create")
	}
	for _, present := range []int{0, 1} {
		for code := 0; code < 11; code++ {
			calls = nil
			legacy.Sub_44A4A0 = func() int { return present }
			*words["nox_wol_wnd_world_814980"] = uint32(uintptr(o.parent.C()))
			*words["dword_5d4594_814548"] = 2
			*words["dword_5d4594_815044"] = 99
			*words["nox_client_connError_814552"] = uint32(code)
			want := []string{}
			if present == 0 {
				want = append(want, "create")
			}
			if code < 8 {
				want = append(want, "title:ConnError")
			}
			want = append(want, fmt.Sprintf("text:BrowserError%d", code), "buttons:1", "enable:1", "back:1")
			if legacy.PortTestServerBrowserTick() != 1 || *words["dword_5d4594_814548"] != 1 || *words["dword_5d4594_815044"] != 0 || !reflect.DeepEqual(calls, want) {
				t.Fatal("error transition", present, code, calls, want)
			}
		}
	}
	calls = nil
	*words["dword_5d4594_814548"] = 5
	if legacy.PortTestServerBrowserTick() != 1 || *words["dword_5d4594_814548"] != 6 || !reflect.DeepEqual(calls, []string{"title:Password", "text:PasswordRequired", "buttons:7", "enable:0", "back:1"}) {
		t.Fatal("password transition", calls)
	}
	calls = nil
	*words["dword_5d4594_814548"] = 7
	draws := 0
	legacy.Nox_xxx_cliDrawConnectedLoop_43B360 = func() int { draws++; return 1 }
	if legacy.PortTestServerBrowserTick() != 1 || *words["dword_5d4594_814548"] != 1 || !reflect.DeepEqual(calls, []string{"enable:1", "text:Connected", "buttons:0"}) {
		t.Fatal("connected transition", calls)
	}
	if o.c.DrawFunc == nil || !o.c.DrawFunc() || draws != 1 {
		t.Fatal("connected draw routing", draws)
	}
}
