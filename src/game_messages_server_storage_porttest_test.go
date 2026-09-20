//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy"
)

func TestGameMessageServerUnknownKinds(t *testing.T) {
	known := map[byte]bool{0x40: true, 0x72: true, 0x73: true, 0x74: true, 0x75: true, 0x76: true, 0x78: true, 0x79: true, 0x7b: true, 0xa5: true, 0xe0: true, 0xe2: true, 0xee: true, 0xc9: true, 0xf0: true, 0xf1: true}
	trade := map[byte]bool{14: true, 15: true, 16: true, 17: true, 18: true, 21: true, 22: true, 23: true, 24: true, 25: true, 26: true, 28: true, 30: true}
	type row struct{ Kind, Subtype, Result int }
	var rows []row
	run := func(kind, sub int) {
		data := bytes.Repeat([]byte{0xa5}, 56)
		data[0], data[1] = byte(kind), byte(sub)
		before := bytes.Clone(data)
		n := legacy.Nox_xxx_netOnPacketRecvServ_51BAD0_net_sdecode_switch(0, data, nil, nil, nil)
		if n != -1 || !bytes.Equal(data, before) {
			t.Fatalf("unknown game message %d/%d return/input %d", kind, sub, n)
		}
		rows = append(rows, row{kind, sub, n})
	}
	for kind := 0; kind < 256; kind++ {
		if !known[byte(kind)] {
			run(kind, 0)
		}
	}
	for sub := 0; sub < 256; sub++ {
		if !trade[byte(sub)] {
			run(201, sub)
		}
		if sub > 5 {
			run(238, sub)
		}
		if sub != 3 && sub != 27 {
			run(240, sub)
		}
	}
	interactionCapture(t, "game-server-unknown", rows)
}

func TestGameMessageServerAliasStorage(t *testing.T) {
	o := newReliableReportsOwner(t)
	u := &o.units[0]
	pl := u.UpdateDataPlayer().Player
	raw := unsafe.Slice((*byte)(pl.C()), unsafe.Sizeof(*pl))
	original := bytes.Clone(raw)
	defer copy(raw, original)
	type row struct {
		Slot   int
		A, B   uint16
		Frame  uint32
		Return int
		Alias  []byte
	}
	var rows []row
	for slot := 0; slot < 255; slot++ {
		for _, values := range [][3]uint32{{0, 0, 0}, {65535, 32768, 0xffffffff}, {123, 456, 0x80000000}} {
			copy(raw, original)
			before := bytes.Clone(raw)
			data := []byte{165, byte(slot), 0, 0, 0, 0, 0, 0, 0, 0}
			binary.LittleEndian.PutUint16(data[2:], uint16(values[0]))
			binary.LittleEndian.PutUint16(data[4:], uint16(values[1]))
			binary.LittleEndian.PutUint32(data[6:], values[2])
			input := bytes.Clone(data)
			n := legacy.Nox_xxx_netOnPacketRecvServ_51BAD0_net_sdecode_switch(pl.PlayerIndex(), data, pl, u, u.UpdateData)
			want := bytes.Clone(before)
			copy(want[16+8*slot:16+8*(slot+1)], data[2:])
			if n != 10 || !bytes.Equal(raw, want) || !bytes.Equal(data, input) {
				t.Fatal("alias storage/input/neighbor contract", slot, values, n)
			}
			rows = append(rows, row{slot, uint16(values[0]), uint16(values[1]), values[2], n, bytes.Clone(raw[16+8*slot : 16+8*(slot+1)])})
		}
	}
	interactionCapture(t, "game-server-aliases", rows)
}
