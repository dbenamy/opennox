//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"
	"unsafe"

	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy"
)

func TestGameMessageClientSessionDisconnect(t *testing.T) {
	o := newEntryOwner(t)
	words, restore := legacy.PortTestServerBrowserWords()
	t.Cleanup(restore)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	set, restoreStrings := o.c.srv.PortTestMeterStrings(
		strman.Entry{ID: "noxworld.c:Notification", Vals: []strman.Variant{{Str: "Notice title"}}},
		strman.Entry{ID: "noxworld.c:Kicked", Vals: []strman.Variant{{Str: "Removed from game"}}},
		strman.Entry{ID: "noxworld.c:Timeout", Vals: []strman.Variant{{Str: "Connection expired"}}},
	)
	t.Cleanup(restoreStrings)
	set(0)
	oldCreate, oldEnable := legacy.Nox_xxx_dialogMsgBoxCreate_449A10, legacy.Sub_44A360
	t.Cleanup(func() { legacy.Nox_xxx_dialogMsgBoxCreate_449A10 = oldCreate; legacy.Sub_44A360 = oldEnable })
	var calls []string
	legacy.Nox_xxx_dialogMsgBoxCreate_449A10 = func(w *gui.Window, title, message string, flags gui.DialogFlags, a, b func()) {
		if w != nil || title != "Notice title" || flags != 33 || a != nil || b != nil {
			t.Fatal("notification dialog arguments", title, message, flags)
		}
		calls = append(calls, message)
	}
	legacy.Sub_44A360 = func(v int) {
		if v != 1 {
			t.Fatal("notification enable", v)
		}
		calls = append(calls, "enable")
	}
	type row struct {
		On, Present, Op          int
		Initial, Kicked, Timeout uint32
		Return                   int
		Calls                    []string
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for present := 0; present < 2; present++ {
			for _, op := range []int{197, 198} {
				for _, initial := range []uint32{0, 1, 0x80000000, 0xffffffff} {
					binary.LittleEndian.PutUint32(connected, uint32(on))
					*words["nox_wol_wnd_world_814980"] = 0
					if present != 0 {
						*words["nox_wol_wnd_world_814980"] = uint32(uintptr(unsafe.Pointer(o.parent)))
					}
					*words["dword_5d4594_815096"], *words["dword_5d4594_815100"] = initial, initial
					calls = nil
					data := []byte{byte(op), 0xaa, 0x55}
					input := bytes.Clone(data)
					n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(op), data)
					kicked, timeout := initial, initial
					wantPending := uint32(1)
					var wantCalls []string
					if present != 0 {
						wantPending = 0
						message := "Removed from game"
						if op == 198 {
							message = "Connection expired"
						}
						wantCalls = []string{message, "enable"}
					}
					if op == 197 {
						kicked = wantPending
					} else {
						timeout = wantPending
					}
					if n != 1 || !bytes.Equal(data, input) || *words["dword_5d4594_815096"] != kicked || *words["dword_5d4594_815100"] != timeout || !reflect.DeepEqual(calls, wantCalls) {
						t.Fatal("disconnect notification", on, present, op, initial, n, calls)
					}
					rows = append(rows, row{on, present, op, initial, kicked, timeout, n, append([]string(nil), calls...)})
				}
			}
		}
	}
	interactionCapture(t, "game-client-session-disconnect", rows)
}
