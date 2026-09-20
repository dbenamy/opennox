//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"reflect"
	"testing"
)

func TestGameMessageClientSessionSummonCommand(t *testing.T) {
	o := newSummonOwner(t)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	type row struct {
		On, Present, Command int
		State                summonResult
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for present := 0; present < 2; present++ {
			for command := 0; command < 256; command++ {
				setup := func() {
					o.reset(t)
					if present != 0 {
						o.init(t)
					}
					binary.LittleEndian.PutUint32(connected, uint32(on))
					*memmap.PtrUint32(0x587000, 184448) = 0xcafe1234
				}
				setup()
				if on != 0 {
					o.call("sub_4C1CA0", uint32(command))
				}
				want := o.snapshot("command", 0)
				setup()
				data := []byte{237, byte(command)}
				input := bytes.Clone(data)
				n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(237), data)
				got := o.snapshot("command", 0)
				expected := uint32(0xcafe1234)
				if on != 0 {
					expected = uint32(command)
				}
				if n != 2 || !bytes.Equal(data, input) || memmap.Uint32(0x587000, 184448) != expected || !reflect.DeepEqual(got, want) {
					t.Fatal("creature command", on, present, command)
				}
				rows = append(rows, row{on, present, command, got})
			}
		}
	}
	interactionCapture(t, "game-client-session-summon-command", rows)
}
