//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/music"
	"testing"
)

func TestClientAudioEventsMusicQueue(t *testing.T) {
	words, restore := legacy.PortTestAudioEventGlobals()
	defer restore()
	old := legacy.MusicModule
	legacy.MusicModule = &music.Module{}
	defer func() { legacy.MusicModule = old }()
	memory := serverConfigOwnBytes(t, 0x5D4594, 815772, 320)
	count, level := words["dword_5d4594_816368"], words["dword_5d4594_816372"]
	put := func(dst []byte, s music.MusicState) {
		for i, v := range []uint32{s.MusicIdx, s.Volume, s.Position, s.D} {
			binary.LittleEndian.PutUint32(dst[i*4:], v)
		}
	}
	get := func(src []byte) music.MusicState {
		return music.MusicState{MusicIdx: binary.LittleEndian.Uint32(src), Volume: binary.LittleEndian.Uint32(src[4:]), Position: binary.LittleEndian.Uint32(src[8:]), D: binary.LittleEndian.Uint32(src[12:])}
	}
	var rows []map[string]any
	for l := uint32(0); l < 3; l++ {
		for _, n := range []uint32{0, 1, 5, 6, 7, 0x7fffffff} {
			for i := range memory {
				memory[i] = 0xa5
			}
			*count, *level = n, l
			current := music.MusicState{MusicIdx: 17 + l, Volume: 63, Position: 0x12345678 + n, D: 0xabcdef01}
			legacy.MusicModule.SetNextMusic(current)
			expected := append([]byte(nil), memory...)
			wantRet, wantCount := uint64(0), uint32(6)
			if n < 6 {
				wantRet, wantCount = 1, n+1
				put(expected[16*(n+6*l):], current)
			}
			ret := legacy.PortTestAudioEventCall("sub_43DA80")
			if ret != wantRet || *count != wantCount || *level != l || !bytes.Equal(memory, expected) {
				t.Fatal("music save capacity", l, n)
			}
			rows = append(rows, map[string]any{"op": "save", "level": l, "count": n, "return": ret, "after": wantCount})
		}
	}
	for l := uint32(0); l < 3; l++ {
		for n := uint32(0); n <= 6; n++ {
			for i := range memory {
				memory[i] = 0xa5
			}
			for slot := uint32(0); slot < 6; slot++ {
				put(memory[16*(slot+6*l):], music.MusicState{MusicIdx: 100 + slot, Volume: 90 + slot*7, Position: slot * 1000, D: slot ^ 0x87654321})
			}
			expected := append([]byte(nil), memory...)
			*count, *level = n, l
			want := music.MusicState{MusicIdx: 9, Volume: 42, Position: 123, D: 45}
			legacy.MusicModule.SetNextMusic(want)
			if n > 0 {
				want = get(memory[16*(n-1+6*l):])
				if want.Volume > 100 {
					want.Volume = 100
				}
			}
			legacy.PortTestAudioEventCall("sub_43DAD0")
			if *count != 0 || *level != l || legacy.MusicModule.GetCurrentBlock() != want || !bytes.Equal(memory, expected) {
				t.Fatal("music restore", l, n)
			}
			rows = append(rows, map[string]any{"op": "restore", "level": l, "count": n, "current": want})
		}
	}
	for l := uint32(0); l < 3; l++ {
		for _, n := range []uint32{0, 1, 5, 6} {
			clear(memory)
			*count, *level = n, l
			for slot := uint32(0); slot < 6; slot++ {
				put(memory[16*(slot+6*l):], music.MusicState{MusicIdx: slot + 1, Volume: 50, Position: slot * 10, D: l})
			}
			current := music.MusicState{MusicIdx: 60 + l, Volume: 75, Position: 100 + n, D: 77}
			legacy.MusicModule.SetNextMusic(current)
			want := current
			if n == 6 {
				want = get(memory[16*(5+6*l):])
			}
			if legacy.PortTestAudioEventCall("sub_43DB60") != uint64(l+1) || *level != l+1 || *count != 0 {
				t.Fatal("music nesting enter")
			}
			saved := n + 1
			if saved > 6 {
				saved = 6
			}
			if binary.LittleEndian.Uint32(memory[304+4*l:]) != saved {
				t.Fatal("saved parent count")
			}
			legacy.MusicModule.SetNextMusic(music.MusicState{MusicIdx: 99, Volume: 10})
			legacy.PortTestAudioEventCall("sub_43DBA0")
			if *level != l || *count != 0 || legacy.MusicModule.GetCurrentBlock() != want {
				t.Fatal("music nesting restore", l, n)
			}
			rows = append(rows, map[string]any{"op": "nested", "level": l, "count": n, "saved": saved, "restored": want})
		}
	}
	*level, *count = 3, 4
	if legacy.PortTestAudioEventCall("sub_43DB60") != 3 || *level != 3 || *count != 4 {
		t.Fatal("music nesting limit")
	}
	*level = 0
	legacy.PortTestAudioEventCall("sub_43DBA0")
	if *level != 0 || *count != 4 {
		t.Fatal("music empty nesting")
	}
	for _, n := range []uint32{0, 1, 6, 7, 0x7fffffff, 0x80000000, 0xffffffff} {
		ret := legacy.PortTestAudioEventCall("sub_43DB30", uint64(n))
		if ret != audioEventSigned(int32(n)) || legacy.PortTestAudioEventCall("sub_43DB20") != ret || *count != n {
			t.Fatal("raw count setter")
		}
	}
	spellbookCapture(t, "client-audio-events-music-queue", rows, "f4395e653e79f5fcc1873a02955619b263d6f577f3eb71ea8fc9956b71f15527")
}
