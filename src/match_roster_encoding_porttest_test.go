//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

func TestMatchRosterPlayerEncoding(t *testing.T) {
	o := newQuestRuntimeOwner(t)
	pl := o.units[0].UpdateDataPlayer().Player
	type row struct {
		Name        string
		Full, Short []byte
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "match-roster-player-encoding", rows, "77d354cb4c415146df721af17f6142a04893cca8288b394acd5d2923523a4c96")
	}()
	for _, seed := range []uint32{0, 1, 0x103, 0x8000ffff, 0xffffffff} {
		for _, status := range []uint32{0, 0x423, 0xffffffff, 0xfffffbdc, 0x80000000} {
			for _, name := range []string{"", "a", "abcdefghi", "abcdefghij", "abcdefghijk"} {
				for _, fill := range []byte{0, 0xa5} {
					label := fmt.Sprintf("seed%x/status%x/name%q/fill%x", seed, status, name, fill)
					t.Run(label, func(t *testing.T) {
						info := unsafe.Slice((*byte)(pl.Info().C()), 97)
						for i := range info {
							info[i] = byte(seed + uint32(i*37))
						}
						vals := map[int]uint32{0: seed ^ 0x12345678, 4: seed ^ 0x87654321, 2060: seed, 2136: seed + 1, 2140: seed - 1, 2152: seed + 7, 2156: seed - 7, 3676: seed, 3680: status}
						for off, v := range vals {
							objectXferSetWord(pl.C(), off, v)
						}
						id := unsafe.Slice((*byte)(unsafe.Add(pl.C(), 2096)), 12)
						clear(id)
						copy(id, name)
						raw := bytes.Repeat([]byte{fill}, 140)
						full := raw[4:136]
						legacy.PortTestMatchRosterPlayerPacket(full, pl)
						want := bytes.Repeat([]byte{fill}, 132)
						want[0] = 45
						binary.LittleEndian.PutUint16(want[1:], uint16(seed))
						copy(want[3:100], info)
						binary.LittleEndian.PutUint16(want[100:], uint16(seed+1))
						binary.LittleEndian.PutUint16(want[102:], uint16(seed-1))
						binary.LittleEndian.PutUint32(want[104:], seed^0x12345678)
						binary.LittleEndian.PutUint32(want[108:], seed^0x87654321)
						binary.LittleEndian.PutUint32(want[112:], status&0x423)
						want[116] = byte(seed + 7)
						want[117] = byte(seed - 7)
						want[118] = 0
						if byte(seed) == 3 {
							want[118] = 1
						}
						copy(want[119:], name)
						want[119+len(name)] = 0
						if !bytes.Equal(full, want) {
							t.Fatalf("C roster encoding %x want %x", full, want)
						}
						if !bytes.Equal(raw[:4], bytes.Repeat([]byte{fill}, 4)) || !bytes.Equal(raw[136:], bytes.Repeat([]byte{fill}, 4)) {
							t.Fatal("caller guard changed")
						}
						// The existing Go caller accepts either a full scratch buffer or its
						// shorter on-wire destination. A full identifier field has no terminator.
						native := bytes.Repeat([]byte{fill}, 132)
						server.EncodePlayerRoster(native, pl)
						if !bytes.Equal(native, want) {
							t.Fatalf("existing Go scratch encoding differs: %x want %x", native, want)
						}
						short := bytes.Repeat([]byte{fill}, 129)
						server.EncodePlayerRoster(short, pl)
						shortWant := bytes.Clone(want[:129])
						if !bytes.Equal(short, shortWant) {
							t.Fatalf("short Go encoding %x want %x", short, shortWant)
						}
						rows = append(rows, row{label, bytes.Clone(full), short})
					})
				}
			}
		}
	}
}
