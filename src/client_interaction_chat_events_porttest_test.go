//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"image"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
	"golang.org/x/image/font/basicfont"
)

func TestClientInteractionChatDrawingAndSubmit(t *testing.T) {
	o := sessionDialogResources(t)
	oldNet := o.c.srv.NetList
	o.c.srv.NetList = netlist.New()
	o.c.srv.NetList.Init()
	defer func() { o.c.srv.NetList.Free(); o.c.srv.NetList = oldNet }()
	words, restore := legacy.PortTestClientInteractionWords()
	defer restore()
	scoreboard, restoreScoreboard := legacy.PortTestScoreboardWords()
	defer restoreScoreboard()
	*scoreboard["nox_player_netCode_85319C"] = 0x1234
	oldList := o.c.Objs.List1
	defer func() { o.c.Objs.List1 = oldList }()
	o.c.Objs.List1 = nil
	mode := serverConfigOwnBytes(t, 0x5D4594, 1064872, 4)
	serverConfigOwnBytes(t, 0x5D4594, 1064876, 8)
	restoreFlags := noxflags.PortTestGameFlags(0)
	defer restoreFlags()
	dim := legacy.PortTestBindingDimensions()
	oldw, oldh := *dim[0], *dim[1]
	defer func() { *dim[0], *dim[1] = oldw, oldh }()
	*dim[0], *dim[1] = 801, 601
	fontRestore := o.c.Render().GetFonts().PortTestDefaultFont(basicfont.Face7x13)
	defer fontRestore()
	if interactionCall("sub_46A730") == 0 {
		t.Fatal("chat constructor")
	}
	defer interactionCall("sub_46A860")
	w := (*gui.Window)(unsafe.Pointer(uintptr(*words["dword_5d4594_1064856"])))
	edit := w.ChildByID(9201)
	data := (*gui.EntryFieldData)(edit.WidgetData)
	type drawRow struct {
		Length, Composition, Width int
		Position                   image.Point
		Return                     uint64
	}
	var drawRows []drawRow
	for _, length := range []int{0, 1, 12, 13, 14, 43, 44, 45, 46, 100, 250} {
		for _, comp := range []int{0, 4, 10} {
			clear(data.Text[:])
			for i := 0; i < length; i++ {
				data.Text[i] = 'A'
			}
			for i := 0; i < comp; i++ {
				data.Text[256+i] = 'B'
			}
			data.Field_1052 = uint32(length)
			ret := interactionCall("sub_46A5D0", uintptr(edit.C()), uintptr(unsafe.Pointer(edit.DrawData())))
			width := max(100, min(320, 7*(length+comp)+10))
			want := image.Pt((801-width)/2, 400)
			if edit.SizeVal != image.Pt(width, 20) || w.Off != want || o.c.GUI.Focused() != edit || ret != 1 {
				t.Fatal("chat sizing/focus", length, comp, edit.SizeVal, w.Off, ret)
			}
			drawRows = append(drawRows, drawRow{length, comp, width, w.Off, ret})
		}
	}
	type submitRow struct {
		LengthWord uint32
		Mode       uint32
		Return     uint64
		Message    []byte
	}
	var submitRows []submitRow
	for _, lengthWord := range []uint32{0, 1, 0xffff, 0x10000, 0x10001, 0xffffffff} {
		for _, team := range []uint32{0, 1, 2, 0xffffffff} {
			interactionCall("nox_client_chatStart_46A430", uintptr(team))
			clear(data.Text[:])
			data.Text[0] = 'H'
			data.Text[1] = 'i'
			data.Field_1052 = lengthWord
			binary.LittleEndian.PutUint32(mode, team)
			o.c.srv.NetList.ResetAll()
			ret := interactionCall("sub_46A820", uintptr(w.C()), 16415, uintptr(edit.C()), 0)
			var msg []byte
			packets := 0
			o.c.srv.NetList.ByInd(31, netlist.Kind0).Each(func(b []byte) bool { msg = append(msg, b...); packets++; return false })
			var want []byte
			if uint16(lengthWord) != 0 {
				// No local player object: only the message header, flags and ASCII text.
				want = make([]byte, 14)
				want[0] = 0xa8
				want[3] = 2
				if team != 0 {
					want[3] |= 1
				}
				want[8] = 3
				binary.LittleEndian.PutUint16(want[1:], 0x1234)
				binary.LittleEndian.PutUint16(want[4:], 0xffff)
				binary.LittleEndian.PutUint16(want[6:], 0xffff)
				copy(want[11:], []byte{'H', 'i', 0})
			}
			if ret != 0 || !bytes.Equal(msg, want) || packets != len(want)/14 || !w.GetFlags().IsHidden() || *words["dword_5d4594_1064868"] != 0 {
				t.Fatalf("submit length%x mode%x: message%x want%x ret%d", lengthWord, team, msg, want, ret)
			}
			submitRows = append(submitRows, submitRow{lengthWord, team, ret, msg})
		}
	}
	interactionCapture(t, "chat-drawing", drawRows)
	interactionCapture(t, "chat-submit", submitRows)
}
