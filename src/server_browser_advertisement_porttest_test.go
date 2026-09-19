//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
	"unsafe"
)

func TestServerBrowserHostDescription(t *testing.T) {
	type row struct {
		Flags               uint32
		Headless            bool
		Quest, Stage, Count int
		Record, OtherSlot   []byte
		Host                string
	}
	var rows []row
	for _, flags := range []uint32{0, 0x20, 0x1000, 0xffff} {
		for _, headless := range []bool{false, true} {
			for _, quest := range []int{0, 1, 2} {
				for _, count := range []int{0, 3} {
					t.Run(fmt.Sprintf("%x/%t/%d/%d", flags, headless, quest, count), func(t *testing.T) {
						o := newServerOptionsOwner(t)
						words, restore := legacy.PortTestServerBrowserWords()
						defer restore()
						defer noxflags.PortTestGameFlags(noxflags.GameFlag(flags))()
						engine := noxflags.GetEngine()
						noxflags.ResetEngine()
						if headless {
							noxflags.SetEngine(noxflags.EngineNoRendering)
						}
						defer func() { noxflags.ResetEngine(); noxflags.SetEngine(engine) }()
						video := legacy.Sub_43BE50_get_video_mode_id
						legacy.Sub_43BE50_get_video_mode_id = func() int { return 3 }
						defer func() { legacy.Sub_43BE50_get_video_mode_id = video }()
						hostFn := legacy.ClientSetServerHost
						host := ""
						legacy.ClientSetServerHost = func(s string) { host = s }
						defer func() { legacy.ClientSetServerHost = hostFn }()
						questWord := serverConfigOwnBytes(t, 0x5D4594, 1556160, 4)
						binary.LittleEndian.PutUint32(questWord, uint32(quest))
						stageWord := serverConfigOwnBytes(t, 0x587000, 202028, 4)
						binary.LittleEndian.PutUint32(stageWord, 65538)
						*words["dword_5d4594_528256"] = uint32(quest / 2)
						for i := range o.players {
							o.players[i].Active = 0
						}
						for i := 0; i < count; i++ {
							o.players[i].Active = 1
						}
						legacy.PortTestServerConfigScalar("limit-set", count, 0)
						name := serverConfigOwnBytes(t, 0x5D4594, 1324, 16)
						clear(name)
						copy(name, "HostABCDEFGHIJK")
						currentMap := serverConfigOwnBytes(t, 0x85B3FC, 36, 80)
						clear(currentMap)
						copy(currentMap, "arena123")
						coords := serverConfigOwnBytes(t, 0x5D4594, 814916, 4)
						binary.LittleEndian.PutUint32(coords, 0x80007fff)
						slot0 := serverConfigOwnBytes(t, 0x5D4594, 371380, 58)
						slot1 := serverConfigOwnBytes(t, 0x5D4594, 371438, 58)
						record := serverConfigOwnBytes(t, 0x5D4594, 371516, 184)
						for i := range slot0 {
							slot0[i] = byte(17 + i)
						}
						for i := range slot1 {
							slot1[i] = byte(91 + i)
						}
						for i := range record {
							record[i] = 0xa5
						}
						expected := append([]byte(nil), record...)
						copy(expected[111:169], slot0)
						binary.LittleEndian.PutUint16(expected[163:], uint16(flags))
						for i := 135; i < 163; i++ {
							expected[i] = 255
						}
						copy(expected[120:135], name[:15])
						copy(expected[111:120], currentMap[:9])
						stage := int(binary.LittleEndian.Uint16(slot0[54:]))
						if quest != 0 {
							stage = 1
							if quest == 2 {
								stage = 2
							}
							binary.LittleEndian.PutUint16(expected[165:], uint16(stage))
						}
						expected[103], expected[104] = byte(count), byte(count)
						if headless {
							expected[103]--
							expected[104]--
						}
						expected[102] = 0x83
						copy(expected[44:48], coords)
						binary.LittleEndian.PutUint32(expected[48:], uint32(noxProtoVersionHighRes))
						binary.LittleEndian.PutUint16(expected[109:], uint16(o.c.srv.ServerPort()))
						other := append([]byte(nil), slot1...)
						binary.LittleEndian.PutUint16(other[52:], binary.LittleEndian.Uint16(other[52:])&0xe90f|0x100)
						got := legacy.PortTestServerBrowserHostDescription()
						if got != uintptr(unsafe.Pointer(&slot1[0])) || !bytes.Equal(record, expected) || !bytes.Equal(slot1, other) || host != "localhost" || *words["nox_game_createOrJoin_815048"] != 0 || *words["dword_5d4594_815052"] != 1 {
							t.Fatalf("host description mismatch: record=%x expected=%x slot=%x want=%x host=%q", record, expected, slot1, other, host)
						}
						captured := append([]byte(nil), record...)
						// legacy/video_highres.go unconditionally enables the C highres version.
						// Preserve and capture that original behavior on every Go build target.
						rows = append(rows, row{flags, headless, quest, stage, count, captured, append([]byte(nil), slot1...), host})
					})
				}
			}
		}
	}
	spellbookCapture(t, "server-browser-host-description", rows, "173636c317bef2a11ab205a50ccd52181952dfbc65aaf9b4b084a0bb530faf74")
}
