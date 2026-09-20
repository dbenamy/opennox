//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy"
	"reflect"
	"testing"
	"unsafe"
)

func TestGameMessageClientSessionGauntletKey(t *testing.T) {
	sessionDialogResources(t)
	words, restore := legacy.PortTestClientInteractionWords()
	t.Cleanup(restore)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	if interactionCall("sub_4BFC90") != 1 {
		t.Fatal("key constructor")
	}
	t.Cleanup(func() { interactionCall("sub_4BFD10") })
	w := (*gui.Window)(unsafe.Pointer(uintptr(*words["dword_5d4594_1319060"])))
	var sounds [][2]int
	t.Cleanup(legacy.PortTestClientSoundObserver(func(id, v int) { sounds = append(sounds, [2]int{id, v}) }))
	type row struct {
		On, Hidden, Value int
		Initial, After    uint32
		AfterHidden       bool
		Sounds            [][2]int
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for hidden := 0; hidden < 2; hidden++ {
			for _, initial := range []uint32{0, 1, 2, 0xffffffff} {
				for value := 0; value < 256; value++ {
					binary.LittleEndian.PutUint32(connected, uint32(on))
					*words["dword_5d4594_1319056"] = initial
					w.Show()
					if hidden != 0 {
						w.Hide()
					}
					sounds = nil
					want, hide := initial, hidden != 0
					var expected [][2]int
					if on != 0 {
						if initial == 0 && value == 1 {
							want = 1
							hide = false
							expected = [][2]int{{1022, 100}}
						} else if initial == 1 && value == 0 {
							want = 0
							hide = true
						}
					}
					data := []byte{240, 24, byte(value)}
					input := bytes.Clone(data)
					n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(240), data)
					if n != 3 || !bytes.Equal(input, data) || *words["dword_5d4594_1319056"] != want || w.GetFlags().IsHidden() != hide || !reflect.DeepEqual(sounds, expected) {
						t.Fatal("quest key visibility", on, hidden, initial, value, n, sounds)
					}
					rows = append(rows, row{on, hidden, value, initial, want, hide, append([][2]int(nil), sounds...)})
				}
			}
		}
	}
	interactionCapture(t, "game-client-session-gauntlet-key", rows)
}
func TestGameMessageClientSessionGauntletStage(t *testing.T) {
	newEntryOwner(t)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	stage := serverConfigOwnBytes(t, 0x5D4594, 527720, 4)
	type row struct {
		On, Value int
		After     uint32
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for value := 0; value < 256; value++ {
			binary.LittleEndian.PutUint32(connected, uint32(on))
			binary.LittleEndian.PutUint32(stage, 0xcafe1234)
			data := []byte{240, 28, byte(value)}
			input := bytes.Clone(data)
			n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(240), data)
			want := uint32(0xcafe1234)
			if on != 0 {
				want = uint32(value)
			}
			if n != 3 || !bytes.Equal(input, data) || binary.LittleEndian.Uint32(stage) != want {
				t.Fatal("quest stage", on, value, n)
			}
			rows = append(rows, row{on, value, want})
		}
	}
	interactionCapture(t, "game-client-session-gauntlet-stage", rows)
}
