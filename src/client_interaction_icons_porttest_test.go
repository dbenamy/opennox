//go:build porttest

package opennox

import (
	"image"
	"slices"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func TestClientInteractionIndicatorWindows(t *testing.T) {
	o := sessionDialogResources(t)
	words, restore := legacy.PortTestClientInteractionWords()
	defer restore()
	dim := legacy.PortTestBindingDimensions()
	oldw, oldh := *dim[0], *dim[1]
	defer func() { *dim[0], *dim[1] = oldw, oldh }()
	serverConfigOwnBytes(t, 0x5D4594, 825748, 4)
	serverConfigOwnBytes(t, 0x5D4594, 1193716, 4)
	pl, free := alloc.New(server.Player{})
	defer free()
	type row struct {
		Viewport, Chat, Observer image.Point
		HiddenReturns            []uint32
	}
	var rows []row
	for _, size := range []image.Point{{640, 480}, {801, 601}, {1280, 960}} {
		*dim[0], *dim[1] = int32(size.X), int32(size.Y)
		if interactionCall("nox_xxx_guiChatIconLoad_445650") != 1 || interactionCall("sub_48C980") != 1 {
			t.Fatal("indicator constructors")
		}
		chat := (*gui.Window)(unsafe.Pointer(uintptr(*words["dword_5d4594_825744"])))
		obs := (*gui.Window)(unsafe.Pointer(uintptr(*words["dword_5d4594_1193712"])))
		if chat == nil || obs == nil {
			t.Fatal("indicator owners")
		}
		wantChat, wantObs := image.Pt(size.X-50, size.Y/2-50), image.Pt(size.X-50, size.Y/2-100)
		if chat.Off != wantChat || obs.Off != wantObs || chat.SizeVal != image.Pt(50, 50) || obs.SizeVal != image.Pt(50, 50) {
			t.Fatal("indicator layout", size, chat.Off, obs.Off)
		}
		r := row{Viewport: size, Chat: chat.Off, Observer: obs.Off}
		for _, tc := range []struct {
			w          *gui.Window
			draw, hide string
		}{{chat, "nox_xxx_guiChatMode_4456E0", "nox_xxx_guiChatShowHide_445730"}, {obs, "sub_48C9F0", "nox_xxx_showObserverWindow_48CA70"}} {
			for _, v := range []uint32{0, 1, 2, 0xffffffff, 0} {
				ret := uint32(interactionCall(tc.hide, uintptr(v)))
				r.HiddenReturns = append(r.HiddenReturns, ret)
				if tc.w.GetFlags().IsHidden() != (v != 0) {
					t.Fatal("indicator visibility", tc.hide, v)
				}
			}
			tc.w.SetPos(image.Pt(25, 30))
			for _, off := range []image.Point{{0, 0}, {7, 9}, {-5, 3}} {
				tc.w.DrawData().ImgPtVal = off
				for _, present := range []bool{false, true} {
					*words["dword_8531A0_2576"] = 0
					if present {
						*words["dword_8531A0_2576"] = uint32(uintptr(unsafe.Pointer(pl)))
					}
					clear(o.pix.Pix)
					if interactionCall(tc.draw, uintptr(tc.w.C())) != 1 {
						t.Fatal("indicator draw return")
					}
					got := append([]uint16(nil), o.pix.Pix...)
					clear(o.pix.Pix)
					if tc.w == chat || present {
						o.c.R2().DrawImageAt(o.images[0], image.Pt(25, 30).Add(off))
					}
					if !slices.Equal(got, o.pix.Pix) {
						t.Fatal("indicator draw pixels", tc.draw, off, present)
					}
				}
			}
		}
		interactionCall("sub_445770")
		if *words["dword_5d4594_825744"] != 0 {
			t.Fatal("chat icon clear")
		}
		obs.Destroy()
		*words["dword_5d4594_1193712"] = 0
		o.c.GUI.FreeDestroyed()
		rows = append(rows, r)
	}
	interactionCapture(t, "indicators", rows)
}
