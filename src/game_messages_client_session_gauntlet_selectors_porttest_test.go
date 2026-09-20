//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestGameMessageClientSessionGauntletSelectors(t *testing.T) {
	newEntryOwner(t)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	type row struct{ On, Kind, Return int }
	var rows []row
	handled := map[int]bool{0: true, 1: true, 2: true, 4: true, 12: true, 13: true, 14: true, 15: true, 16: true, 17: true, 18: true, 19: true, 20: true, 21: true, 22: true, 23: true, 24: true, 25: true, 26: true, 28: true, 29: true, 30: true, 31: true, 32: true, 33: true}
	for on := 0; on < 2; on++ {
		for kind := 0; kind < 256; kind++ {
			if handled[kind] {
				continue
			}
			binary.LittleEndian.PutUint32(connected, uint32(on))
			data := bytes.Repeat([]byte{0xa5}, 16)
			data[0], data[1] = 240, byte(kind)
			input := bytes.Clone(data)
			want := -1
			if kind >= 5 && kind <= 10 {
				want = 4
			}
			if kind == 11 {
				want = 16
			}
			n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(240), data)
			if n != want || !bytes.Equal(input, data) {
				t.Fatal("quest selector", on, kind, n, want)
			}
			rows = append(rows, row{on, kind, n})
		}
	}
	interactionCapture(t, "game-client-session-gauntlet-selectors", rows)
}
