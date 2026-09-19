//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/spell"
	flags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestPlayerFilesSpellbookRestore(t *testing.T) {
	o := newReliableReportsOwner(t)
	t.Cleanup(flags.PortTestGameFlags(2048))
	u := &o.units[0]
	p := u.UpdateDataPlayer().Player
	p.NetCodeVal = u.NetCode
	raw := unsafe.Slice((*byte)(unsafe.Pointer(p)), int(unsafe.Sizeof(*p)))
	configure, restore := o.s.PortTestAISpellDefs()
	t.Cleanup(restore)
	t.Cleanup(o.s.PortTestBookSpellOwner())
	// Explicit ordinary spell definitions: family propagation is a dependency
	// behavior, separately covered by the spell reward tests.
	configure([]server.PortTestSpellClassDef{{Index: 1, Flags: 0x100, Valid: true}, {Index: 136, Flags: 0x100, Valid: true}})
	ids, resetRecords, snapshotRecords, free := legacy.PortTestPlayerFileRecords()
	t.Cleanup(free)
	binary.LittleEndian.PutUint32(raw[4636:], ids[0])
	binary.LittleEndian.PutUint32(raw[4640:], ids[1])
	var rows []map[string]any
	for _, class := range []byte{0, 1, 2} {
		for _, version := range []uint16{0, 1, 2, 3, 0x8000, 0xffff} {
			for _, level := range []uint32{0, 1, 3, 5, 0xffffffff} {
				o.reset()
				resetRecords()
				clear(raw[3700:4244])
				raw[2251] = class
				ability := class == 0 && int16(version) >= 3
				indices := []int{1, 136}
				if ability {
					indices = []int{1, 5}
				}
				input := binary.LittleEndian.AppendUint16(nil, version)
				input = append(input, 1, 2)
				var wantPackets [][]byte
				wantBits := uint32(0)
				for _, id := range indices {
					name := spell.ID(id).String()
					if ability {
						name = server.Ability(id).String()
					}
					input = append(input, byte(len(name)))
					input = append(input, name...)
					if int16(version) >= 2 {
						input = binary.LittleEndian.AppendUint32(input, level)
					}
					wantLevel := level
					if int16(version) < 2 {
						wantLevel = 3
					} else if level == 0 {
						wantLevel = 1
					}
					if ability {
						wantLevel = 5
					}
					packet := []byte{111, byte(id), byte(wantLevel), 0}
					if ability {
						packet = []byte{205, byte(id), byte(wantLevel)}
					}
					wantPackets = append(wantPackets, packet)
					wantBits |= 1 << uint(id&31)
				}
				expected := bytes.Clone(raw)
				for _, id := range indices {
					value := level
					if int16(version) < 2 {
						value = 3
					} else if level == 0 {
						value = 1
					}
					if ability {
						value = 5
					}
					binary.LittleEndian.PutUint32(expected[3696+4*id:], value)
				}
				size := len(input)
				input = append(input, 0xde, 0xad, 0xbe, 0xef)
				ret, got, pos := playerFileSection(t, "nox_xxx_guiSpellbook_41B660", input, uint32(uintptr(unsafe.Pointer(u))), 0)
				records := snapshotRecords()
				if ret != 1 || pos != int64(size) || !bytes.Equal(got, input) || !bytes.Equal(raw, expected) || records != [3]uint32{wantBits, 0, 0x2468ace0 ^ wantBits} {
					t.Fatal("spellbook restore", class, version, level, ret, pos, size, records, wantBits)
				}
				state := o.state()
				if len(state.Nodes) != len(wantPackets) {
					t.Fatal("award reports", len(state.Nodes), len(wantPackets))
				}
				// The reliable queue inserts at the head, reversing award order.
				for i, n := range state.Nodes {
					if !bytes.Equal(n.Data, wantPackets[len(wantPackets)-1-i]) {
						t.Fatal("award packet", n.Data, wantPackets)
					}
				}
				rows = append(rows, map[string]any{"class": class, "version": version, "level": level, "return": ret, "position": pos, "skills": bytes.Clone(raw[3700:4244]), "records": records, "queue": state})
			}
		}
	}
	spellbookCapture(t, "player-files-spellbook-restore", rows, "")
}
