//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	flags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"strings"
	"testing"
	"unsafe"
)

func playerFileGameOwner(t *testing.T) (*reliableReportsOwner, []byte, []byte, *byte) {
	o := newReliableReportsOwner(t)
	t.Cleanup(legacy.PortTestQuestProgressOwner())
	p := o.units[0].UpdateDataPlayer().Player
	raw := unsafe.Slice((*byte)(unsafe.Pointer(p)), int(unsafe.Sizeof(*p)))
	current := unsafe.Slice(memmap.PtrUint8(0x85B3FC, 36), 32)
	oldMap := bytes.Clone(current)
	t.Cleanup(func() { copy(current, oldMap) })
	audio := memmap.PtrUint8(0x5D4594, 831252)
	oldAudio := *audio
	t.Cleanup(func() { *audio = oldAudio })
	return o, raw, current, audio
}
func TestPlayerFilesGameWrite(t *testing.T) {
	o, raw, current, audio := playerFileGameOwner(t)
	defer flags.PortTestGameFlags(0)()
	var rows []map[string]any
	for _, gf := range []flags.GameFlag{0, 2048, 8192} {
		for _, name := range []string{"", "map", strings.Repeat("M", 31)} {
			flags.ResetGame()
			flags.SetGame(gf)
			legacy.PortTestQuestProgress("reset", "*:*", 0)
			legacy.PortTestQuestProgress("namespace", "map", 0)
			legacy.PortTestQuestProgress("set-int", "map:key", 0x87654321)
			o.s.Objs.SetLastObjectScriptID(server.ObjectScriptID(0x87654321))
			clear(current)
			copy(current, name)
			for i := 4760; i < len(raw); i++ {
				raw[i] = byte(i*11 + 9)
			}
			expected := bytes.Clone(raw)
			*audio = 0x89
			var want []byte
			if gf&8192 == 0 {
				copy(expected[4760:], name)
				expected[4760+len(name)] = 0
				want = []byte{5, 0}
				want = binary.LittleEndian.AppendUint32(want, 0x87654321)
				want = binary.LittleEndian.AppendUint16(want, uint16(len(name)))
				// Legacy format stores twice the byte-string length, including following
				// record bytes. Preserve this established format rather than treating UTF-8
				// as UTF-16 or silently shortening the section.
				want = append(want, expected[4760:4760+2*len(name)]...)
				if gf&2048 != 0 {
					want = append(want, questProgressWire(1, []string{"map:key"}, []uint32{0}, []uint32{0x87654321})...)
				} else {
					want = append(want, questProgressWire(1, nil, nil, nil)...)
				}
				want = append(want, 1, 0, 0, 0x89) // Empty magic-wall section, audio byte.
			}
			ret, got, pos := playerFileSection(t, "sub_41C080", nil, uint32(uintptr(unsafe.Pointer(&o.units[0]))), 0)
			if ret != 1 || pos != int64(len(want)) || !bytes.Equal(got, want) || !bytes.Equal(raw, expected) || *audio != 0x89 || o.s.Objs.LastObjectScriptID() != 0x87654321 {
				t.Fatal("game write", gf, len(name), ret, pos, len(want))
			}
			rows = append(rows, map[string]any{"flags": uint32(gf), "name": name, "return": ret, "position": pos, "bytes": got, "record_tail": bytes.Clone(raw[4760:])})
		}
	}
	spellbookCapture(t, "player-files-game-write", rows, "")
}
func TestPlayerFilesGameRead(t *testing.T) {
	o, raw, _, audio := playerFileGameOwner(t)
	defer flags.PortTestGameFlags(0)()
	var rows []map[string]any
	for _, version := range []uint16{0, 1, 2, 3, 4, 5, 6, 0x8000, 0xffff} {
		for _, gf := range []flags.GameFlag{0, 8192} {
			for _, size := range []int{0, 1, 15, 16, 31} {
				flags.ResetGame()
				flags.SetGame(gf)
				legacy.PortTestQuestProgress("reset", "*:*", 0)
				legacy.PortTestQuestProgress("namespace", "map", 0)
				legacy.PortTestQuestProgress("set-int", "map:old", 77)
				o.s.Objs.SetLastObjectScriptID(12345)
				for i := 4760; i < len(raw); i++ {
					raw[i] = 0xa5
				}
				raw[4760] = 0
				expected := bytes.Clone(raw)
				*audio = 0x89
				wantAudio := byte(0x89)
				input := binary.LittleEndian.AppendUint16(nil, version)
				if int16(version) >= 5 {
					input = binary.LittleEndian.AppendUint32(input, 0xfedcba98)
				}
				input = binary.LittleEndian.AppendUint16(input, uint16(size))
				payload := bytes.Repeat([]byte{0xe7}, 2*size)
				for i := 0; i < size; i++ {
					payload[i] = 'A' + byte(i%26)
				}
				input = append(input, payload...)
				if int16(version) >= 2 {
					input = append(input, questProgressWire(1, []string{"map:new"}, []uint32{0}, []uint32{91})...)
				}
				if int16(version) >= 3 {
					input = append(input, 1, 0, 0)
				}
				if int16(version) >= 4 {
					input = append(input, 0xfe)
				}
				wantRet, wantPos := uint32(1), int64(len(input))
				applied := false
				if gf&8192 != 0 {
					wantPos = 0
				} else if int16(version) > 5 {
					wantRet = 0
					wantPos = 2
				} else {
					applied = true
					copy(expected[4760:], payload)
					expected[4760+size] = 0
					wantAudio = 0
					if int16(version) >= 4 {
						wantAudio = 0xfe
					}
				}
				input = append(input, 0xde, 0xad, 0xbe, 0xef)
				ret, got, pos := playerFileSection(t, "sub_41C080", input, uint32(uintptr(unsafe.Pointer(&o.units[0]))), 0)
				wantOld, wantNew := uint64(77), uint64(0)
				if applied && int16(version) >= 2 {
					wantOld = 0
					wantNew = 91
				}
				if ret != wantRet || pos != wantPos || !bytes.Equal(got, input) || !bytes.Equal(raw, expected) || *audio != wantAudio || o.s.Objs.LastObjectScriptID() != 12345 || legacy.PortTestQuestProgress("int", "map:old", 0) != wantOld || legacy.PortTestQuestProgress("int", "map:new", 0) != wantNew {
					t.Fatal("game read", version, gf, size, ret, pos, wantPos, *audio, wantAudio)
				}
				rows = append(rows, map[string]any{"version": version, "flags": uint32(gf), "length": size, "return": ret, "position": pos, "audio": *audio, "record_tail": bytes.Clone(raw[4760:]), "variables": legacy.PortTestQuestProgressSnapshot()})
			}
		}
	}
	spellbookCapture(t, "player-files-game-read", rows, "")
}
