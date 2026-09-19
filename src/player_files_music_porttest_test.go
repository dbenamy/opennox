//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	flags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/music"
	"testing"
)

func playerFileMusicBytes(dst []byte, s music.MusicState) []byte {
	for _, v := range []uint32{s.D, s.Position, s.MusicIdx, s.Volume} {
		dst = binary.LittleEndian.AppendUint32(dst, v)
	}
	return dst
}
func playerFileMusicStore(dst []byte, s music.MusicState) {
	for i, v := range []uint32{s.MusicIdx, s.Volume, s.Position, s.D} {
		binary.LittleEndian.PutUint32(dst[i*4:], v)
	}
}
func playerFileMusicStates(n int) []music.MusicState {
	var out []music.MusicState
	for i := 0; i < n; i++ {
		out = append(out, music.MusicState{MusicIdx: uint32(10 + i), Volume: uint32(60 + i*3), Position: uint32(123456 + i*1789), D: 0x87654320 + uint32(i)})
	}
	return out
}
func TestPlayerFilesMusicWrite(t *testing.T) {
	words, restore := legacy.PortTestAudioEventGlobals()
	defer restore()
	old := legacy.MusicModule
	defer func() { legacy.MusicModule = old }()
	defer flags.PortTestGameFlags(0)()
	memory := serverConfigOwnBytes(t, 0x5D4594, 815772, 320)
	count, level := words["dword_5d4594_816368"], words["dword_5d4594_816372"]
	var rows []map[string]any
	for _, gameFlags := range []flags.GameFlag{0, 8192} {
		for _, n := range []int{0, 1, 6, 7} {
			for l := uint32(0); l < 3; l++ {
				flags.ResetGame()
				flags.SetGame(gameFlags)
				for i := range memory {
					memory[i] = 0xa5
				}
				*count, *level = uint32(n), l
				legacy.MusicModule = &music.Module{}
				current := music.MusicState{MusicIdx: 17, Volume: 63, Position: 0x12345678, D: 0xabcdef01}
				legacy.MusicModule.SetNextMusic(current)
				states := playerFileMusicStates(n)
				for i, s := range states {
					playerFileMusicStore(memory[16*(uint32(i)+6*l):], s)
				}
				before := append([]byte(nil), memory...)
				want := []byte{11, 0, 0}
				if gameFlags&8192 == 0 {
					want[2] = 1
					want = playerFileMusicBytes(want, current)
					want = binary.LittleEndian.AppendUint32(want, uint32(n))
					for _, s := range states {
						want = playerFileMusicBytes(want, s)
					}
				}
				ret, got, pos := playerFileSection(t, "sub_41C780", nil, 0)
				if ret != 1 || !bytes.Equal(got, want) || pos != int64(len(want)) || *count != uint32(n) || *level != l || !bytes.Equal(memory, before) || legacy.MusicModule.GetCurrentBlock() != current {
					t.Fatal("music write", gameFlags, n, l, ret, pos, got, want)
				}
				rows = append(rows, map[string]any{"flags": uint32(gameFlags), "count": n, "level": l, "bytes": got, "position": pos})
			}
		}
	}
	spellbookCapture(t, "player-files-music-write", rows, "")
}
func TestPlayerFilesMusicRead(t *testing.T) {
	words, restore := legacy.PortTestAudioEventGlobals()
	defer restore()
	old := legacy.MusicModule
	defer func() { legacy.MusicModule = old }()
	defer flags.PortTestGameFlags(0)()
	memory := serverConfigOwnBytes(t, 0x5D4594, 815772, 320)
	count, level := words["dword_5d4594_816368"], words["dword_5d4594_816372"]
	var rows []map[string]any
	for _, version := range []uint16{0, 1, 10, 11, 12, 0xffff} {
		for _, gameFlags := range []flags.GameFlag{0, 8192} {
			for _, present := range []byte{0, 1} {
				for _, n := range []int{0, 1, 6, 7} {
					for _, l := range []uint32{0, 2} {
						if version != 11 && present == 0 {
							continue
						}
						flags.ResetGame()
						flags.SetGame(gameFlags)
						for i := range memory {
							memory[i] = 0xa5
						}
						*count, *level = 5, l
						legacy.MusicModule = &music.Module{}
						current := music.MusicState{MusicIdx: 17, Volume: 63, Position: 0x12345678, D: 0xabcdef01}
						legacy.MusicModule.SetNextMusic(current)
						loaded := music.MusicState{MusicIdx: 37, Volume: 75, Position: 123456789, D: 0xfedcba98}
						states := playerFileMusicStates(n)
						input := binary.LittleEndian.AppendUint16(nil, version)
						if version == 11 {
							input = append(input, present)
						}
						input = playerFileMusicBytes(input, loaded)
						input = binary.LittleEndian.AppendUint32(input, uint32(n))
						for _, s := range states {
							input = playerFileMusicBytes(input, s)
						}
						input = append(input, 0xde, 0xad, 0xbe, 0xef)
						expected := append([]byte(nil), memory...)
						wantRet, wantPos, wantCount, wantCurrent := uint32(1), int64(len(input)-4), uint32(5), current
						if int16(version) > 11 {
							wantRet = 0
							wantPos = 2
						} else if version == 11 && present == 0 {
							wantPos = 3
						} else {
							for i, s := range states {
								playerFileMusicStore(expected[16*(uint32(i)+6*l):], s)
							}
							if gameFlags&8192 == 0 {
								wantCount = uint32(n)
								wantCurrent = loaded
							}
						}
						ret, got, pos := playerFileSection(t, "sub_41C780", input, 0)
						if ret != wantRet || pos != wantPos || !bytes.Equal(got, input) || *count != wantCount || *level != l || !bytes.Equal(memory, expected) || legacy.MusicModule.GetCurrentBlock() != wantCurrent {
							t.Fatal("music read", version, gameFlags, present, n, l, ret, pos, wantPos, *count, wantCount)
						}
						rows = append(rows, map[string]any{"version": version, "flags": uint32(gameFlags), "present": present, "count": n, "level": l, "return": ret, "position": pos, "stored_count": *count, "current": legacy.MusicModule.GetCurrentBlock(), "queue": append([]byte(nil), memory...)})
					}
				}
			}
		}
	}
	spellbookCapture(t, "player-files-music-read", rows, "")
}
