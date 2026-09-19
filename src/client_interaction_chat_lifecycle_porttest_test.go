//go:build porttest

package opennox

import (
	"encoding/binary"
	"image"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
)

func TestClientInteractionChatLifecycle(t *testing.T) {
	o := sessionDialogResources(t)
	words, restore := legacy.PortTestClientInteractionWords()
	defer restore()
	mode := serverConfigOwnBytes(t, 0x5D4594, 1064872, 4)
	serverConfigOwnBytes(t, 0x5D4594, 1064876, 8)
	dim := legacy.PortTestBindingDimensions()
	oldw, oldh := *dim[0], *dim[1]
	defer func() { *dim[0], *dim[1] = oldw, oldh }()
	*dim[0], *dim[1] = 801, 601
	restoreFlags := noxflags.PortTestGameFlags(0)
	defer restoreFlags()
	if interactionCall("sub_46A730") == 0 {
		t.Fatal("chat constructor")
	}
	defer interactionCall("sub_46A860")
	w := (*gui.Window)(unsafe.Pointer(uintptr(*words["dword_5d4594_1064856"])))
	edit := w.ChildByID(9201)
	data := (*gui.EntryFieldData)(edit.WidgetData)
	if *words["dword_5d4594_1064860"] != uint32(uintptr(edit.C())) || *words["dword_5d4594_1064864"] != uint32(uintptr(edit.WidgetData)) || w.Off != image.Pt(400, 400) {
		t.Fatal("chat resource owners/position")
	}
	if !w.GetFlags().IsHidden() {
		t.Fatal("resource starts hidden")
	}
	type row struct {
		Blocked              bool
		Mode, Length, Active uint32
		Flags                uint32
	}
	var rows []row
	for _, blocked := range []bool{true, false} {
		for _, value := range []uint32{0, 1, 2, 0x80000000, 0xffffffff} {
			*words["dword_5d4594_1064868"] = 0
			data.Text[0] = 'x'
			data.Field_1052 = 0x12340006
			binary.LittleEndian.PutUint32(mode, 0xfeed)
			flags := noxflags.GameFlag(0)
			if blocked {
				flags = 2048
			}
			reset := noxflags.PortTestGameFlags(flags)
			interactionCall("nox_client_chatStart_46A430", uintptr(value))
			if blocked {
				if data.Text[0] != 'x' || data.Field_1052 != 0x12340006 || interactionCall("sub_46A4A0") != 0 || binary.LittleEndian.Uint32(mode) != 0xfeed {
					t.Fatal("blocked chat mutated")
				}
			} else {
				if data.Text[0] != 0 || data.Field_1052 != 0x12340000 || interactionCall("sub_46A4A0") != 1 || binary.LittleEndian.Uint32(mode) != value || w.GetFlags().IsHidden() || o.c.GUI.Focused() != edit {
					t.Fatal("chat start", value, data.Field_1052)
				}
				data.Text[0] = 'z'
				data.Field_1052 = 0x12340007
				interactionCall("nox_client_chatStart_46A430", 123)
				if data.Text[0] != 'z' || data.Field_1052 != 0x12340007 || binary.LittleEndian.Uint32(mode) != value {
					t.Fatal("active chat reset")
				}
			}
			rows = append(rows, row{blocked, binary.LittleEndian.Uint32(mode), data.Field_1052, *words["dword_5d4594_1064868"], uint32(w.Flags)})
			if !blocked {
				if interactionCall("sub_46A6A0") != 1 || !w.GetFlags().IsHidden() || interactionCall("sub_46A4A0") != 0 || o.c.GUI.Focused() == edit || o.c.GUI.ValYYY != 1 || w.Flags&gui.StatusEnabled != 0 || edit.Flags&gui.StatusEnabled != 0 {
					t.Fatal("chat close")
				}
			}
			if interactionCall("sub_46A6A0") != 0 {
				t.Fatal("hidden close")
			}
			reset()
		}
	}
	// Parser notifications carry numeric IDs; callbacks must not decode them as pointers.
	for _, code := range []uintptr{0, 22, 23, 16391} {
		if interactionCall("sub_46A820", uintptr(w.C()), code, 0xffffffff, 0) != 0 {
			t.Fatal("chat unrelated event", code)
		}
	}
	interactionCall("sub_46A860")
	for _, name := range []string{"dword_5d4594_1064856", "dword_5d4594_1064860", "dword_5d4594_1064864", "dword_5d4594_1064868"} {
		if *words[name] != 0 {
			t.Fatal("chat owner after destroy", name)
		}
	}
	if binary.LittleEndian.Uint32(mode) != 0 {
		t.Fatal("chat mode after destroy")
	}
	interactionCapture(t, "chat-lifecycle", rows)
}
