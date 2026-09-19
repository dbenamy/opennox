//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestSessionFilterConfigurationCopies(t *testing.T) {
	modes := serverConfigOwnBytes(t, 0x5D4594, 1193372, 8)
	settings := serverConfigOwnBytes(t, 0x5D4594, 1193388, 88)
	data, free := alloc.Make([]byte{}, 44)
	defer free()
	for i := range data {
		data[i] = byte(i*37 + 19)
	}
	for slot := 0; slot < 2; slot++ {
		for _, mode := range []uint32{0, 1, 2, 0xffffffff} {
			for i := range modes {
				modes[i] = 0xa5
			}
			for i := range settings {
				settings[i] = 0x5a
			}
			old := append([]byte(nil), data...)
			got := legacy.PortTestSessionDialogCall("filterConfig", unsafe.Pointer(&data[0]), 0, uint32(slot), mode)
			if got != uintptr(11*slot) || binary.LittleEndian.Uint32(modes[slot*4:]) != mode {
				t.Fatal("config return/mode", slot, mode, got)
			}
			if binary.LittleEndian.Uint32(modes[(1-slot)*4:]) != 0xa5a5a5a5 {
				t.Fatal("neighbor mode")
			}
			if !bytes.Equal(settings[slot*44:(slot+1)*44], data) || !bytes.Equal(settings[(1-slot)*44:(2-slot)*44], bytes.Repeat([]byte{0x5a}, 44)) || !bytes.Equal(data, old) {
				t.Fatal("config copy boundaries", slot, mode)
			}
		}
	}
}

func TestSessionFilterActionDispatch(t *testing.T) {
	o := sessionDialogResources(t)
	_, restore := legacy.PortTestSessionDialogWords()
	defer restore()
	modes := serverConfigOwnBytes(t, 0x5D4594, 1193372, 8)
	settings := serverConfigOwnBytes(t, 0x5D4594, 1193388, 88)
	clear(settings)
	clear(modes)
	w := (*gui.Window)(unsafe.Pointer(sessionCall("filterOpen", o.parent, 0, 0, 0)))
	if w == nil {
		t.Fatal("filter constructor")
	}
	defer sessionCall("filterClose", nil, 0, 0, 0)
	controls := w.ChildByID(10012)
	type row struct {
		ID                           uint
		Bits                         uint32
		Mode                         uint32
		Hidden, Ping, RadioA, RadioB bool
	}
	var rows []row
	for _, bits := range []uint32{0, 4, 0xfffffffb, 0xffffffff} {
		for _, id := range []uint{10024, 10025, 10026, 10028, 10015, 10029, 10030, 10014, 10018} {
			b := w.ChildByID(id)
			b.DrawData().Field0 = bits
			ping := w.ChildByID(10031)
			a := w.ChildByID(10016)
			z := w.ChildByID(10017)
			ping.Flags |= gui.StatusEnabled
			a.Flags |= gui.StatusEnabled
			z.Flags |= gui.StatusEnabled
			if sessionCall("filterEvent", w, 16391, uint32(uintptr(b.C())), 0) != 0 {
				t.Fatal("click result", id)
			}
			mode := binary.LittleEndian.Uint32(modes)
			hidden := controls.Flags.Has(gui.StatusHidden)
			if id >= 10024 && id <= 10026 {
				if mode != uint32(id-10024) || hidden != (id != 10026) {
					t.Fatal("mode action", id, mode, hidden)
				}
			}
			if id == 10028 && ping.Flags.Has(gui.StatusEnabled) != (bits&4 == 0) {
				t.Fatal("ping toggle", bits)
			}
			if id == 10015 && (a.Flags.Has(gui.StatusEnabled) != (bits&4 == 0) || z.Flags.Has(gui.StatusEnabled) != (bits&4 == 0)) {
				t.Fatal("resolution toggle", bits)
			}
			if b.DrawData().Field0 != bits {
				t.Fatal("callback changed button before widget handles click")
			}
			rows = append(rows, row{id, bits, mode, hidden, ping.Flags.Has(gui.StatusEnabled), a.Flags.Has(gui.StatusEnabled), z.Flags.Has(gui.StatusEnabled)})
		}
	}
	sessionCapture(t, "filter-actions", rows)
}
