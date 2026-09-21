//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"strings"
	"testing"

	"github.com/opennox/libs/noxnet/discover"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
)

func TestServerTextListing(t *testing.T) {
	o := newObjectDrawingOwner(t)
	players, free := o.c.srv.PortTestObjectRenderPlayers()
	t.Cleanup(free)
	words, restore := legacy.PortTestServerBrowserWords()
	t.Cleanup(restore)
	settings := serverConfigOwnBytes(t, 0x5D4594, 371516, 184)
	slot := serverConfigOwnBytes(t, 0x5D4594, 371380, 58)
	name := serverConfigOwnBytes(t, 0x5D4594, 1324, 16)
	mapName := serverConfigOwnBytes(t, 0x85B3FC, 36, 64)
	limit := serverConfigOwnBytes(t, 0x5D4594, 3464, 4)
	quest := serverConfigOwnBytes(t, 0x5D4594, 1556160, 4)
	stage := serverConfigOwnBytes(t, 0x587000, 202028, 4)
	oldBlock, oldVideo := legacy.Sub_4D6F30, legacy.Sub_43BE50_get_video_mode_id
	t.Cleanup(func() { legacy.Sub_4D6F30, legacy.Sub_43BE50_get_video_mode_id = oldBlock, oldVideo })
	t.Cleanup(noxflags.PortTestGameFlags(0))
	engine := noxflags.GetEngine()
	t.Cleanup(func() {
		noxflags.ResetEngine()
		noxflags.SetEngine(engine)
	})
	type row struct {
		Mode, Count, Map, SourceLength, DestinationLength, Return int
		Bytes                                                     []byte
	}
	var rows []row
	maps := []string{"", "a", "war01a", "ABCDEFGH", "ABCDEFGHI", strings.Repeat("x", 63), "Épée😀"}
	names := []string{"", "Nox", "FifteenByteName!", "é😀 server"}
	limits := []uint32{0, 1, 32, 255, 256, 0xffffffff, 0x80000000, 0x7fffffff}
	for mode := 0; mode < 64; mode++ {
		enabled, blocked, dedicated, isQuest := mode&1 != 0, mode&2 != 0, mode&4 != 0, mode&8 != 0
		*words["dword_5d4594_815052"] = 0
		if enabled {
			*words["dword_5d4594_815052"] = []uint32{1, 2, 0xffffffff}[mode%3]
		}
		legacy.Sub_4D6F30 = func() int {
			if blocked {
				return 1
			}
			return 0
		}
		video := []int{0, 1, 16, 128, 255, -1}[mode%6]
		legacy.Sub_43BE50_get_video_mode_id = func() int { return video }
		noxflags.ResetEngine()
		if dedicated {
			noxflags.SetEngine(noxflags.EngineNoRendering)
		}
		flags := []uint32{0, 1, 0x80, 0x1000, 0x87654321, 0xffffffff}[mode%6]
		noxflags.PortTestGameFlags(noxflags.GameFlag(flags))
		binary.LittleEndian.PutUint32(quest, 0)
		if isQuest {
			binary.LittleEndian.PutUint32(quest, []uint32{1, 2, 0xffffffff}[mode%3])
		}
		binary.LittleEndian.PutUint32(stage, uint32(mode)*0x1234567)
		binary.LittleEndian.PutUint32(limit, limits[mode%len(limits)])
		for i := range settings {
			settings[i] = byte(i*37 + mode*11)
		}
		for i := range slot {
			slot[i] = byte(i*19 + mode*43)
		}
		clear(name)
		copy(name[:15], names[mode%len(names)])
		serverName := string(name[:bytes.IndexByte(name, 0)])
		for _, count := range []int{0, 1, 31, 32} {
			for i := range players {
				players[i].Active = 0
				players[i].PlayerInd = byte(i)
				if i < count {
					players[i].Active = 1
				}
			}
			for mapIndex, m := range maps {
				clear(mapName)
				copy(mapName, m)
				for _, shape := range [][2]int{{0, 100}, {11, 100}, {12, 0}, {12, 72}, {12, 72 + len(serverName)}, {12, 73 + len(serverName)}, {32, 100}} {
					source := bytes.Repeat([]byte{0x35}, shape[0])
					if len(source) >= 12 {
						binary.LittleEndian.PutUint32(source[8:], 0xfedcba98^uint32(mode))
					}
					beforeSource := bytes.Clone(source)
					out := bytes.Repeat([]byte{0xa5}, shape[1]+16)
					n := legacy.Nox_server_makeServerInfoPacket_554040(source, out[:shape[1]])
					want := bytes.Repeat([]byte{0xa5}, len(out))
					size := 0
					if enabled && !blocked && shape[0] >= 12 && shape[1] >= 73+len(serverName) {
						size = 73 + len(serverName)
						clear(want[:size])
						want[2] = 13
						adjustment := 0
						if dedicated {
							adjustment = 1
						}
						want[3] = byte(count - adjustment)
						want[4] = byte(limits[mode%len(limits)] - uint32(adjustment))
						want[5], want[6] = settings[101]&15, settings[101]>>4
						copy(want[7:10], slot[44:47])
						copy(want[10:18], m)
						want[19] = settings[102] | byte(video)
						want[20] = settings[100]
						want[21] = settings[100] & 16
						copy(want[24:28], settings[48:52])
						ff := flags
						if isQuest {
							ff = ff&^0x80 | 0x1000
							copy(want[68:70], stage[:2])
						}
						binary.LittleEndian.PutUint32(want[28:], ff)
						copy(want[32:36], slot[48:52])
						copy(want[36:38], settings[105:107])
						copy(want[38:40], settings[107:109])
						copy(want[40:44], settings[44:48])
						copy(want[44:48], source[8:12])
						copy(want[48:68], slot[24:44])
						copy(want[72:], serverName)
					}
					if n != size || !bytes.Equal(out, want) || !bytes.Equal(source, beforeSource) {
						t.Fatalf("listing mode=%d count=%d map=%d shape=%v: size=%d/%d bytes=%x/%x", mode, count, mapIndex, shape, n, size, out, want)
					}
					if n != 0 {
						var decoded discover.MsgServerInfo
						if _, err := decoded.Decode(out[3:n]); err != nil || decoded.ServerName != serverName || decoded.MapName != m[:min(len(m), 8)] || decoded.Token != binary.LittleEndian.Uint32(source[8:]) {
							t.Fatal("listing decoded fields", decoded, err)
						}
					}
					rows = append(rows, row{mode, count, mapIndex, shape[0], shape[1], n, bytes.Clone(out)})
				}
			}
		}
	}
	interactionCapture(t, "server-text-listing", rows)
	t.Logf("%d listing cases", len(rows))
}
