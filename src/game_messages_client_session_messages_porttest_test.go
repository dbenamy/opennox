//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"testing"

	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/legacy"
)

func TestGameMessageClientSessionChatClear(t *testing.T) {
	o := newChatBubbleOwner(t)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	type row struct {
		On, Code, Count int
		Bubbles         []chatBubbleRecord
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for _, code := range []uint16{0, 1, 2, 3, 0x8001, 57005, 65535} {
			for _, count := range []int{0, 1, 3, 64} {
				legacy.PortTestChatBubbleClear()
				binary.LittleEndian.PutUint32(connected, uint32(on))
				for i := 1; i <= count; i++ {
					o.create(t, uint16(i), fmt.Sprint("bubble", i), 7, uint32(i*100), 20)
				}
				before := o.snapshot(t)
				var want []chatBubbleRecord
				for _, r := range before {
					if on != 0 && (code == 57005 || r.Code == uint32(code)) {
						continue
					}
					want = append(want, r)
				}
				for i := range want {
					want[i].Previous = 0
					want[i].Next = 0
					if i > 0 {
						want[i].Previous = want[i-1].Code
					}
					if i+1 < len(want) {
						want[i].Next = want[i+1].Code
					}
				}
				data := binary.LittleEndian.AppendUint16([]byte{202}, code)
				input := bytes.Clone(data)
				n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(202), data)
				got := o.snapshot(t)
				if n != 3 || !bytes.Equal(input, data) || !reflect.DeepEqual(got, want) {
					t.Fatal("chat clear/full ID/sentinel", on, code, count)
				}
				rows = append(rows, row{on, int(code), count, got})
			}
		}
	}
	interactionCapture(t, "game-client-session-chat-clear", rows)
}

func TestGameMessageClientSessionMessagesClear(t *testing.T) {
	words, restore := legacy.PortTestClientInteractionWords()
	t.Cleanup(restore)
	raw := serverConfigOwnBytes(t, 0x5D4594, 823804, 1932)
	blank := serverConfigOwnBytes(t, 0x5D4594, 825740, 2)
	clear(blank)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	type row struct {
		On, Head, Fill int
		AfterHead      uint32
		Raw            []byte
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for head := 0; head < 3; head++ {
			for _, fill := range []byte{0, 1, 127, 128, 255} {
				binary.LittleEndian.PutUint32(connected, uint32(on))
				for i := range raw {
					raw[i] = fill
				}
				*words["dword_5d4594_825736"] = uint32(head)
				want := bytes.Clone(raw)
				wantHead := uint32(head)
				if on != 0 {
					wantHead = 0
					for i := 0; i < 3; i++ {
						binary.LittleEndian.PutUint16(want[644*i:], 0)
						binary.LittleEndian.PutUint32(want[644*i+636:], 0)
						want[644*i+640] = 0
					}
				}
				data := []byte{203}
				n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(203), data)
				if n != 1 || data[0] != 203 || !bytes.Equal(want, raw) || *words["dword_5d4594_825736"] != wantHead {
					t.Fatal("message clear fields/gate", on, head, fill)
				}
				rows = append(rows, row{on, head, int(fill), wantHead, bytes.Clone(raw)})
			}
		}
	}
	interactionCapture(t, "game-client-session-messages-clear", rows)
}

func TestGameMessageClientSessionConsoleResult(t *testing.T) {
	o := newInventoryTransactionOwner(t)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	configure, restore := o.c.srv.PortTestMeterStrings(strman.Entry{ID: "cdecode.c:sysopAccess", Vals: []strman.Variant{{Str: "Access granted"}}}, strman.Entry{ID: "cdecode.c:invalidPass", Vals: []strman.Variant{{Str: "Invalid password"}}})
	t.Cleanup(restore)
	configure(0)
	type row struct {
		On, Value int
		Console   []string
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for value := 0; value < 256; value++ {
			o.reset(t)
			binary.LittleEndian.PutUint32(connected, uint32(on))
			data := []byte{189, byte(value)}
			input := bytes.Clone(data)
			n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(189), data)
			var want []string
			if on != 0 {
				text := "Invalid password"
				if value == 1 {
					text = "Access granted"
				}
				want = []string{"6:" + text}
			}
			if n != 2 || !bytes.Equal(input, data) || !reflect.DeepEqual(o.console, want) {
				t.Fatal("console result", on, value, o.console, want)
			}
			rows = append(rows, row{on, value, append([]string(nil), o.console...)})
		}
	}
	interactionCapture(t, "game-client-session-console-result", rows)
}
