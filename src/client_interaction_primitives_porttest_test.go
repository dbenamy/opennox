//go:build porttest

package opennox

import (
	"encoding/binary"
	"encoding/json"
	"os"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func interactionCall(op string, args ...uintptr) uint64 {
	return legacy.PortTestClientInteractionCall(op, args...)
}
func interactionCapture(t *testing.T, name string, v any) {
	t.Helper()
	if p := os.Getenv("OPENNOX_CLIENT_INTERACTION_CAPTURE"); p != "" {
		b, e := json.MarshalIndent(v, "", "  ")
		if e != nil {
			t.Fatal(e)
		}
		if e = os.WriteFile(p+"-"+name+".json", append(b, '\n'), 0600); e != nil {
			t.Fatal(e)
		}
	}
}
func TestClientInteractionMouseModes(t *testing.T) {
	words, restore := legacy.PortTestClientInteractionWords()
	defer restore()
	type row struct{ Input, Return, Mode, Cursor uint32 }
	var rows []row
	for _, v := range []uint32{0, 1, 2, 3, 4, 0x7fffffff, 0x80000000, 0xffffffff} {
		*words["dword_5d4594_805820"] = 0xfeed
		*words["nox_xxx_useAudio_587000_80772"] = 0xbeef
		got := uint32(interactionCall("sub_430AA0", uintptr(v)))
		mode, cursor, want := uint32(0), uint32(5), v-2
		if v == 1 {
			mode, cursor, want = 1, 9, 0
		} else if v == 2 {
			mode, cursor, want = 2, 13, 0
		}
		if got != want || uint32(interactionCall("nox_client_mousePriKey_430AF0")) != mode || uint32(interactionCall("nox_xxx_cursor_430B00")) != cursor {
			t.Fatal("mouse mode", v, got, mode, cursor)
		}
		rows = append(rows, row{v, got, *words["dword_5d4594_805820"], *words["nox_xxx_useAudio_587000_80772"]})
	}
	interactionCapture(t, "mouse-modes", rows)
}
func TestClientInteractionDrawablePredicates(t *testing.T) {
	dr, free := alloc.New(client.Drawable{})
	defer free()
	ptr := uintptr(unsafe.Pointer(dr))
	type row struct {
		Mask uint32
		Bit  int
		Set  bool
	}
	var rows []row
	for _, mask := range []uint32{0, 1, 2, 0x80000000, 0x55555555, 0xaaaaaaaa, 0xffffffff} {
		dr.Buffs = mask
		for bit := 0; bit < 32; bit++ {
			got := interactionCall("nox_client_drawable_testBuff_4356C0", ptr, uintptr(bit)) != 0
			want := mask&(uint32(1)<<uint(bit)) != 0
			if got != want || dr.Buffs != mask {
				t.Fatal("buff predicate", mask, bit, got, want)
			}
			if interactionCall("nox_client_drawable_testBuff_4356C0", 0, uintptr(bit)) != 0 {
				t.Fatal("nil drawable")
			}
			rows = append(rows, row{mask, bit, got})
		}
	}
	active := serverConfigOwnBytes(t, 0x852978, 8, 4)
	if interactionCall("nox_xxx_playerAnimCheck_4372B0") != 1 {
		t.Fatal("no active player")
	}
	binary.LittleEndian.PutUint32(active, uint32(ptr))
	for _, anim := range []uint32{0, 1, 2, 3, 50, 51, 52, 0x80000000, 0xffffffff} {
		dr.AnimInd = anim
		got := interactionCall("nox_xxx_playerAnimCheck_4372B0") != 0
		if got != (anim == 1 || anim == 2 || anim == 51) {
			t.Fatal("animation", anim, got)
		}
	}
	interactionCapture(t, "drawable-predicates", rows)
}
func TestClientInteractionObserverAndToggle(t *testing.T) {
	words, restore := legacy.PortTestClientInteractionWords()
	defer restore()
	pl, free := alloc.New(server.Player{})
	defer free()
	raw := unsafe.Slice((*byte)(unsafe.Pointer(pl)), int(unsafe.Sizeof(*pl)))
	*words["dword_8531A0_2576"] = 0
	if interactionCall("nox_xxx_clientIsObserver_4372E0") != 0 {
		t.Fatal("nil observer")
	}
	*words["dword_8531A0_2576"] = uint32(uintptr(unsafe.Pointer(pl)))
	type row struct {
		Status, Flags uint32
		Observer      bool
	}
	var rows []row
	for _, status := range []uint32{0, 1, 2, 0xffffffff} {
		for _, flags := range []uint32{0, 1, 2, 3, 4, 0x80000000, 0xffffffff} {
			binary.LittleEndian.PutUint32(raw[2092:], status)
			binary.LittleEndian.PutUint32(raw[3680:], flags)
			got := interactionCall("nox_xxx_clientIsObserver_4372E0") != 0
			if got != (status == 1 && flags&3 != 0) {
				t.Fatal("observer", status, flags, got)
			}
			rows = append(rows, row{status, flags, got})
		}
	}
	for _, old := range []uint32{0, 1, 2, 0x7fffffff, 0x80000000, 0xffffffff} {
		p := words["dword_5d4594_811904"]
		*p = old
		if got := uint32(interactionCall("sub_435F60")); got != 1-old || *p != got {
			t.Fatal("toggle", old, got, *p)
		}
	}
	interactionCapture(t, "observer", rows)
}
