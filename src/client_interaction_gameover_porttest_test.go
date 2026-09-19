//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"testing"
	"unsafe"

	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/client/gui"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func interactionStaticText(w *gui.Window) string {
	return alloc.GoString16((*uint16)(unsafe.Pointer(uintptr(gui.EventRespInt(w.Func94(gui.StaticTextGetText{}))))))
}
func TestClientInteractionGameOver(t *testing.T) {
	o := sessionDialogResources(t)
	words, restore := legacy.PortTestClientInteractionWords()
	defer restore()
	dim := legacy.PortTestBindingDimensions()
	oldw, oldh := *dim[0], *dim[1]
	defer func() { *dim[0], *dim[1] = oldw, oldh }()
	*dim[0], *dim[1] = 801, 601
	oldFrame, oldRate := o.c.srv.Frame(), o.c.srv.TickRate()
	defer func() { o.c.srv.SetFrame(oldFrame); o.c.srv.SetTickRate(oldRate) }()
	oldNet := o.c.srv.NetList
	o.c.srv.NetList = netlist.New()
	o.c.srv.NetList.Init()
	defer func() { o.c.srv.NetList.Free(); o.c.srv.NetList = oldNet }()
	resetFlags := noxflags.PortTestGameFlags(1)
	defer resetFlags()
	var sounds [][2]int
	defer legacy.PortTestClientSoundObserver(func(id, v int) { sounds = append(sounds, [2]int{id, v}) })()
	configure, resetStrings := o.c.srv.Server.PortTestMeterStrings(
		strman.Entry{ID: "GGOver.wnd:GeneratorsDestroyed", Vals: []strman.Variant{{Str: "Generators %d"}}},
		strman.Entry{ID: "GGOver.wnd:NumSecretsFound", Vals: []strman.Variant{{Str: "Secrets %d"}}},
		strman.Entry{ID: "GGOver.wnd:Kills", Vals: []strman.Variant{{Str: "Kills %d"}}},
		strman.Entry{ID: "Rules.c:Time", Vals: []strman.Variant{{Str: "Time"}}},
	)
	defer resetStrings()
	configure(0)
	// Control the mapped literal slots and output rows explicitly; exercise nonempty text.
	storage := serverConfigOwnBytes(t, 0x5D4594, 1301852, 1616)
	put := func(off int, text string) {
		p := storage[off-1301852:]
		for i, b := range []byte(text) {
			binary.LittleEndian.PutUint16(p[2*i:], uint16(b))
		}
		binary.LittleEndian.PutUint16(p[2*len(text):], 0)
	}
	put(1303460, "0")
	put(1303464, "-")
	put(1302940, "Stages 12")
	put(1302684, "High 34")
	if interactionCall("sub_49B3E0") != 1 {
		t.Fatal("game over constructor")
	}
	defer interactionCall("sub_49B490")
	w := (*gui.Window)(unsafe.Pointer(uintptr(*words["dword_5d4594_1303452"])))
	if !w.GetFlags().IsHidden() || w.Flags.IsEnabled() {
		t.Fatal("game over initial visibility")
	}
	stats, freeStats := alloc.Make([]uint16{}, 4)
	defer freeStats()
	stats[1], stats[2], stats[3] = 65535, 257, 1234
	o.c.srv.SetFrame(0xfffffff0)
	o.c.srv.SetTickRate(30)
	if uint32(interactionCall("sub_49B4B0", uintptr(unsafe.Pointer(&stats[0])))) != 0xfffffff0 || w.GetFlags().IsHidden() || !w.Flags.IsEnabled() || w.Off != image.Pt(801/2-w.SizeVal.X/2, 601/2-w.SizeVal.Y/2) {
		t.Fatal("game over show")
	}
	for id, want := range map[uint]string{10710: "Generators 65535", 10707: "Secrets 257", 10708: "Kills 1234", 10705: "Stages 12", 10706: "High 34", 10711: "0"} {
		if got := interactionStaticText(w.ChildByID(id)); got != want {
			t.Fatal("game over text", id, got, want)
		}
	}
	pl, freePlayer := alloc.New(server.Player{})
	defer freePlayer()
	type row struct {
		Frame  uint32
		Player int
		Text   string
	}
	var rows []row
	for _, frame := range []uint32{0xfffffff0, 0xffffffff, 0, 884, 885, 1000, 0x80000374} {
		for _, player := range []int{-1, 0, 31, 255} {
			*words["dword_8531A0_2576"] = 0
			if player >= 0 {
				pl.PlayerInd = byte(player)
				*words["dword_8531A0_2576"] = uint32(uintptr(unsafe.Pointer(pl)))
			}
			o.c.srv.SetFrame(frame)
			interactionCall("sub_49B6E0")
			delta := int32(uint32(884) - frame)
			if delta < 0 {
				delta = 0
			}
			want := fmt.Sprintf("Time - %d", uint32(delta)/30)
			if player == 31 {
				want = "-"
			}
			got := interactionStaticText(w.ChildByID(10712))
			if got != want {
				t.Fatal("game over countdown", frame, player, got, want)
			}
			rows = append(rows, row{frame, player, got})
		}
	}
	// Restart sends through the real host message queue, then hides/disables the dialog.
	interactionCall("sub_49B420", uintptr(w.C()), 16391, uintptr(w.ChildByID(10702).C()), 0)
	var message []byte
	o.c.srv.NetList.ByInd(31, netlist.Kind0).Each(func(b []byte) bool { message = append(message, b...); return false })
	if !bytes.Equal(message, []byte{0xf0, 3}) || !w.GetFlags().IsHidden() || w.Flags.IsEnabled() {
		t.Fatal("game over restart", message)
	}
	oldQuit := legacy.Nox_client_quit_4460C0
	defer func() { legacy.Nox_client_quit_4460C0 = oldQuit }()
	quitCalls := 0
	legacy.Nox_client_quit_4460C0 = func() { quitCalls++ }
	w.Show()
	interactionCall("sub_49B420", uintptr(w.C()), 16391, uintptr(w.ChildByID(10701).C()), 0)
	if quitCalls != 1 || !w.GetFlags().IsHidden() {
		t.Fatal("game over quit")
	}
	for _, code := range []uintptr{0, 22, 23} {
		if interactionCall("sub_49B420", uintptr(w.C()), code, 0xffffffff, 0) != 0 {
			t.Fatal("game over numeric notification")
		}
	}
	interactionCall("sub_49B490")
	if *words["dword_5d4594_1303452"] != 0 {
		t.Fatal("game over owner clear")
	}
	interactionCapture(t, "game-over", rows)
}
