//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/spell"
	flags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestPlayerFilesQuestSpellbook(t *testing.T) {
	o := newReliableReportsOwner(t)
	t.Cleanup(flags.PortTestGameFlags(4096))
	u := &o.units[0]
	p := u.UpdateDataPlayer().Player
	p.NetCodeVal = u.NetCode
	raw := unsafe.Slice((*byte)(unsafe.Pointer(p)), int(unsafe.Sizeof(*p)))
	configure, restore := o.s.PortTestAISpellDefs()
	t.Cleanup(restore)
	t.Cleanup(o.s.PortTestBookSpellOwner())
	configure([]server.PortTestSpellClassDef{{Index: 1, Flags: 0x100, Valid: true}})
	table := serverConfigOwnBytes(t, 0x587000, 207108, 24)
	ids, resetRecords, snapshotRecords, free := legacy.PortTestPlayerFileRecords()
	t.Cleanup(free)
	binary.LittleEndian.PutUint32(raw[4636:], ids[0])
	binary.LittleEndian.PutUint32(raw[4640:], ids[1])
	var rows []map[string]any
	for _, class := range []byte{1, 2} {
		for _, version := range []uint16{1, 2, 3, 0xffff} {
			for _, level := range []uint32{0, 1, 3, 4, 0xffffffff} {
				for _, enabled := range []bool{false, true} {
					o.reset()
					resetRecords()
					clear(raw[3700:4244])
					raw[2251] = class
					clear(table)
					binary.LittleEndian.PutUint32(table, 1)
					if enabled {
						binary.LittleEndian.PutUint32(table[4:], 1)
					}
					name := spell.ID(1).String()
					input := binary.LittleEndian.AppendUint16(nil, version)
					input = append(input, 1, 1, byte(len(name)))
					input = append(input, name...)
					value := uint32(3)
					if int16(version) >= 2 {
						value = level
						input = binary.LittleEndian.AppendUint32(input, level)
					}
					size := len(input)
					input = append(input, 0xde, 0xad, 0xbe, 0xef)
					wantRet := uint32(0)
					wantValue := uint32(0)
					bits := uint32(0)
					if enabled && int32(value) <= 3 {
						wantRet = 1
						wantValue = value
						if value == 0 {
							wantValue = 1
						}
						bits = 2
					}
					expected := bytes.Clone(raw)
					binary.LittleEndian.PutUint32(expected[3700:], wantValue)
					ret, got, pos := playerFileSection(t, "nox_xxx_guiSpellbook_41B660", input, uint32(uintptr(unsafe.Pointer(u))), 0)
					records := snapshotRecords()
					state := o.state()
					if ret != wantRet || pos != int64(size) || !bytes.Equal(got, input) || !bytes.Equal(raw, expected) || records != [3]uint32{bits, 0, 0x2468ace0 ^ bits} || len(state.Nodes) != int(wantRet) {
						t.Fatal("quest spell", class, version, level, enabled, ret, wantRet, pos, records, state)
					}
					if wantRet == 1 && !bytes.Equal(state.Nodes[0].Data, []byte{111, 1, byte(wantValue), 0}) {
						t.Fatal("quest spell report", state)
					}
					rows = append(rows, map[string]any{"class": class, "version": version, "level": level, "enabled": enabled, "return": ret, "position": pos, "value": wantValue, "records": records, "queue": state})
				}
			}
		}
	}
	spellbookCapture(t, "player-files-quest-spellbook", rows, "85d4baeb3cb3fdaa75dd8c3dca8e508f0e54e4fa6fc4d25ec363ddd13dc93d91")
}

func TestPlayerFilesQuestFieldbook(t *testing.T) {
	o := newReliableReportsOwner(t)
	t.Cleanup(flags.PortTestGameFlags(4096))
	u := &o.units[0]
	p := u.UpdateDataPlayer().Player
	p.NetCodeVal = u.NetCode
	raw := unsafe.Slice((*byte)(unsafe.Pointer(p)), int(unsafe.Sizeof(*p)))
	names := serverConfigOwnBytes(t, 0x587000, 70500, 41*4)
	clear(names)
	for id := 0; id < 41; id++ {
		name := "OtherGuide"
		if id == 1 {
			name = "QuestGuide"
		}
		ptr, free := alloc.CString(name)
		t.Cleanup(free)
		binary.LittleEndian.PutUint32(names[id*4:], uint32(uintptr(unsafe.Pointer(ptr))))
	}
	rewards := serverConfigOwnBytes(t, 0x587000, 207796, 24)
	families := serverConfigOwnBytes(t, 0x587000, 207032, 8)
	clear(families)
	awards := serverConfigOwnBytes(t, 0x587000, 216292, 8)
	clear(awards)
	ids, resetRecords, snapshotRecords, free := legacy.PortTestPlayerFileRecords()
	t.Cleanup(free)
	binary.LittleEndian.PutUint32(raw[4636:], ids[0])
	binary.LittleEndian.PutUint32(raw[4640:], ids[1])
	family, free := alloc.New([3]uint32{})
	t.Cleanup(free)
	*family = [3]uint32{24, 1, 0}
	var rows []map[string]any
	for _, version := range []uint16{0, 1, 0xffff} {
		for _, which := range []string{"disabled", "listed", "family"} {
			o.reset()
			resetRecords()
			clear(raw[4247:4288])
			clear(rewards)
			clear(families)
			binary.LittleEndian.PutUint32(rewards, 1)
			if which == "listed" {
				binary.LittleEndian.PutUint32(rewards[4:], 1)
			}
			if which == "family" {
				binary.LittleEndian.PutUint32(families, uint32(uintptr(unsafe.Pointer(family))))
			}
			input := binary.LittleEndian.AppendUint16(nil, version)
			input = append(input, 1, 1, 10)
			input = append(input, "QuestGuide"...)
			size := len(input)
			input = append(input, 0xde, 0xad, 0xbe, 0xef)
			wantRet := uint32(1)
			bits := uint32(2)
			if which == "disabled" {
				wantRet = 0
				bits = 0
			}
			expected := bytes.Clone(raw)
			expected[4248] = byte(wantRet)
			ret, got, pos := playerFileSection(t, "nox_xxx_guiFieldbook_41B420", input, uint32(uintptr(unsafe.Pointer(u))), 0)
			records := snapshotRecords()
			state := o.state()
			if ret != wantRet || pos != int64(size) || !bytes.Equal(got, input) || !bytes.Equal(raw, expected) || records != [3]uint32{0, bits, 0x2468ace0 ^ bits} || len(state.Nodes) != int(wantRet) {
				t.Fatal("quest guide", version, which, ret, wantRet, pos, records, state)
			}
			rows = append(rows, map[string]any{"version": version, "case": which, "return": ret, "position": pos, "records": records, "queue": state})
		}
	}
	spellbookCapture(t, "player-files-quest-fieldbook", rows, "da0503b71453aa9c202f6327cc8247fed241ec061a68f54e83f1d966992b5934")
}
