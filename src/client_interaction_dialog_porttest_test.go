//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"image"
	"slices"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/client/audio/ail"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/dialog"
	"github.com/opennox/opennox/v1/legacy/timer"
)

func TestClientInteractionConversation(t *testing.T) {
	o := sessionDialogResources(t)
	words, restore := legacy.PortTestClientInteractionWords()
	defer restore()
	inv, restoreInv := legacy.PortTestInventoryWindowWords()
	defer restoreInv()
	dim := legacy.PortTestBindingDimensions()
	oldw, oldh := *dim[0], *dim[1]
	defer func() { *dim[0], *dim[1] = oldw, oldh }()
	*dim[0], *dim[1] = 665, 510
	serverConfigOwnBytes(t, 0x5D4594, 1107052, 4)
	choice := serverConfigOwnBytes(t, 0x5D4594, 1123516, 4)
	filename := serverConfigOwnBytes(t, 0x5D4594, 1115312, 4)
	file, freeFile := alloc.CString("TestVoice.wav")
	defer freeFile()
	binary.LittleEndian.PutUint32(filename, uint32(uintptr(unsafe.Pointer(file))))
	oldNet := o.c.srv.NetList
	o.c.srv.NetList = netlist.New()
	o.c.srv.NetList.Init()
	defer func() { o.c.srv.NetList.Free(); o.c.srv.NetList = oldNet }()
	var state [6]uint32
	state[0] = 1
	var driver ail.Driver
	var timers [4]timer.TimerGroup
	oldDialog := legacy.Dialogs
	defer func() { legacy.Dialogs = oldDialog }()
	legacy.Dialogs = dialog.NewDialog("dialog", &state[0], &state[1], &state[2], &state[3], &state[4], &driver, &state[5], o.c.srv.Strings, &timers[0], &timers[1], &timers[2], &timers[3], nil, nil, nil, nil, nil, nil, nil)
	var sounds [][2]int
	defer legacy.PortTestClientSoundObserver(func(id, v int) { sounds = append(sounds, [2]int{id, v}) })()
	if interactionCall("sub_4799A0") != 1 {
		t.Fatal("conversation constructor")
	}
	defer interactionCall("sub_479D10")
	w := (*gui.Window)(unsafe.Pointer(uintptr(*words["dword_5d4594_1123524"])))
	list := w.ChildByID(3901)
	slider := w.ChildByID(3904)
	data := (*gui.ScrollListBoxData)(list.WidgetData)
	if data.Field_9 != slider.C() || data.Field_7 != w.ChildByID(3903).C() || data.Field_8 != w.ChildByID(3902).C() || slider.DrawData().Window != list || !w.GetFlags().IsHidden() || w.Flags.IsEnabled() || interactionCall("sub_47A260") != 0 {
		t.Fatal("conversation initial owners/links")
	}
	if interactionCall("sub_479D00") != 1 || interactionCall("sub_479950") != 0 {
		t.Fatal("conversation initial predicates")
	}
	type row struct {
		Book, Button uint32
		Choice       []byte
		Message      []byte
		Voice        string
		Sounds       [][2]int
	}
	var rows []row
	for _, book := range []uint32{0, 1, 2, 0xffffffff} {
		for _, id := range []uint32{3906, 3907, 3908, 3909, 3901} {
			*inv["dword_5d4594_1047520"] = book
			o.c.srv.NetList.ResetAll()
			sounds = nil
			legacy.Dialogs.Sub_44D8F0()
			copy(choice, []byte{0xaa, 0xbb, 0xcc, 0xdd})
			w.Show()
			if interactionCall("nox_xxx_guiDialog_479B00", uintptr(w.C()), 16391, uintptr(w.ChildByID(uint(id)).C()), 0) != 0 {
				t.Fatal("conversation callback return")
			}
			var got []byte
			o.c.srv.NetList.ByInd(31, netlist.Kind0).Each(func(b []byte) bool { got = append(got, b...); return false })
			wantChoice := byte(0xaa)
			var want []byte
			voice := ""
			var wantSounds [][2]int
			if book == 0 {
				wantSounds = [][2]int{{766, 100}}
				switch id {
				case 3906:
					wantChoice = 0
					want = []byte{0xd0, 2, 0}
				case 3907:
					voice = "TestVoice.wav"
				case 3908:
					wantChoice = 1
					want = []byte{0xd0, 2, 1}
				case 3909:
					wantChoice = 2
					want = []byte{0xd0, 2, 2}
				}
			}
			if choice[0] != wantChoice || !bytes.Equal(choice[1:], []byte{0xbb, 0xcc, 0xdd}) || !bytes.Equal(got, want) || legacy.Dialogs.FileToRead() != voice || !slices.Equal(sounds, wantSounds) || w.GetFlags().IsHidden() {
				t.Fatal("conversation action", book, id, choice, got, legacy.Dialogs.FileToRead(), sounds)
			}
			rows = append(rows, row{book, id, append([]byte(nil), choice...), got, legacy.Dialogs.FileToRead(), append([][2]int(nil), sounds...)})
		}
	}
	for _, code := range []uintptr{0, 5, 6, 7, 9, 10, 11, 13, 14, 15, 21, 22, 23, 16390} {
		if interactionCall("sub_479BE0", uintptr(w.C()), code, 0xffffffff, 0) != 1 {
			t.Fatal("conversation input", code)
		}
		if code != 16391 && interactionCall("nox_xxx_guiDialog_479B00", uintptr(w.C()), code, 0xffffffff, 0) != 0 {
			t.Fatal("conversation numeric event", code)
		}
	}
	// Resource images use the root image table, independently of the legacy image hook.
	// Supply an owned nonempty image for the anchoring contract.
	w.DrawData().BgImageHnd = o.images[0].C()
	clear(o.pix.Pix)
	if interactionCall("sub_479CB0", uintptr(w.C()), uintptr(w.DrawData().C())) != 1 {
		t.Fatal("conversation draw")
	}
	got := append([]uint16(nil), o.pix.Pix...)
	clear(o.pix.Pix)
	o.c.R2().DrawImageAt(o.images[0], image.Pt(25, 30))
	if !slices.Equal(got, o.pix.Pix) {
		t.Fatal("conversation background anchoring")
	}
	interactionCall("sub_479D10")
	if *words["dword_5d4594_1123524"] != 0 || interactionCall("sub_47A260") != 0 {
		t.Fatal("conversation owner reset")
	}
	interactionCapture(t, "conversation", rows)
}
