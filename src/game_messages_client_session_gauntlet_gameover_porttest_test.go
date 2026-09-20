//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"image"
	"reflect"
	"testing"
	"unsafe"
)

func TestGameMessageClientSessionGauntletGameOver(t *testing.T) {
	o := sessionDialogResources(t)
	words, restore := legacy.PortTestClientInteractionWords()
	t.Cleanup(restore)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	storage := serverConfigOwnBytes(t, 0x5D4594, 1301852, 1616)
	configure, restoreStrings := o.c.srv.PortTestMeterStrings(strman.Entry{ID: "GGOver.wnd:GeneratorsDestroyed", Vals: []strman.Variant{{Str: "Generators %d"}}}, strman.Entry{ID: "GGOver.wnd:NumSecretsFound", Vals: []strman.Variant{{Str: "Secrets %d"}}}, strman.Entry{ID: "GGOver.wnd:Kills", Vals: []strman.Variant{{Str: "Kills %d"}}})
	t.Cleanup(restoreStrings)
	configure(0)
	oldFrame := o.c.srv.Frame()
	t.Cleanup(func() { o.c.srv.SetFrame(oldFrame) })
	if interactionCall("sub_49B3E0") != 1 {
		t.Fatal("game over constructor")
	}
	t.Cleanup(func() { interactionCall("sub_49B490") })
	w := (*gui.Window)(unsafe.Pointer(uintptr(*words["dword_5d4594_1303452"])))
	var sounds [][2]int
	t.Cleanup(legacy.PortTestClientSoundObserver(func(id, v int) { sounds = append(sounds, [2]int{id, v}) }))
	type state struct {
		Hidden, Enabled bool
		Pos             image.Point
		Storage         []byte
		Text            []string
		Sounds          [][2]int
	}
	snapshot := func() state {
		r := state{Hidden: w.GetFlags().IsHidden(), Enabled: w.Flags.IsEnabled(), Pos: w.Off, Storage: bytes.Clone(storage), Sounds: append([][2]int(nil), sounds...)}
		for _, id := range []uint{10710, 10705, 10706, 10707, 10708, 10711} {
			r.Text = append(r.Text, interactionStaticText(w.ChildByID(id)))
		}
		return r
	}
	type row struct {
		On     int
		Frame  uint32
		Values [3]uint16
		Return int
		State  state
	}
	var rows []row
	data, free := alloc.Make([]byte{}, 14)
	t.Cleanup(free)
	for on := 0; on < 2; on++ {
		for _, frame := range []uint32{0, 100, 0xffffffff} {
			for _, values := range [][3]uint16{{0, 0, 0}, {1, 257, 65535}, {32767, 32768, 1234}} {
				clear(data)
				data[0], data[1] = 240, 2
				for i, v := range values {
					binary.LittleEndian.PutUint16(data[2+2*i:], v)
				}
				for i := 8; i < len(data); i++ {
					data[i] = byte(17 * i)
				}
				setup := func() {
					interactionCall("sub_49B6B0")
					clear(storage)
					for _, r := range []struct {
						off  int
						text string
					}{{1303460, "0"}, {1303464, "-"}, {1302940, "Stages 12"}, {1302684, "High 34"}} {
						alloc.StrCopyZero16(unsafe.Slice((*uint16)(unsafe.Pointer(&storage[r.off-1301852])), len(r.text)+1), r.text)
					}
					for _, id := range []uint{10710, 10705, 10706, 10707, 10708, 10711} {
						w.ChildByID(id).Func94(gui.AsWindowEvent(16385, uintptr(unsafe.Pointer(alloc.InternCString16("Before"))), 0))
					}
					w.SetPos(image.Pt(7, 9))
					binary.LittleEndian.PutUint32(connected, uint32(on))
					o.c.srv.SetFrame(frame)
					sounds = nil
				}
				setup()
				if on != 0 {
					interactionCall("sub_49B4B0", txptr(unsafe.Pointer(&data[0])))
				}
				want := snapshot()
				setup()
				input := bytes.Clone(data)
				n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(240), data)
				got := snapshot()
				if n != 14 || !bytes.Equal(input, data) || !reflect.DeepEqual(got, want) {
					t.Fatal("quest game over", on, frame, values, n)
				}
				rows = append(rows, row{on, frame, values, n, got})
			}
		}
	}
	interactionCapture(t, "game-client-session-gauntlet-gameover", rows)
}
