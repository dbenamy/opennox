//go:build porttest

package opennox

import (
	"bytes"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
	"testing"
	"unsafe"
)

func TestClientAudioStreamsCache(t *testing.T) {
	t.Cleanup(handles.PortTestInit())
	cat, free := alloc.Make([]uint32{}, 72)
	defer free()
	entries, free := alloc.Make([]uint32{}, 9*4)
	defer free()
	cp := audioStreamPointer(unsafe.Pointer(&cat[0]))
	cat[0] = audioStreamPointer(unsafe.Pointer(&entries[0]))
	cat[1] = 4
	samples := [][]byte{[]byte("abc"), []byte("defgh"), []byte("ijklmno"), []byte("01234567890123456")}
	var bag []byte
	for i, b := range samples {
		entries[i*9+4] = uint32(len(bag))
		entries[i*9+5] = uint32(len(b))
		entries[i*9+6] = 22050
		entries[i*9+7] = 4
		bag = append(bag, b...)
	}
	files, _ := prefabScriptsFiles(t, bag)
	cat[67] = audioStreamPointer(files[0])
	call := legacy.PortTestAudioStreamCall
	cache := call("sub_4BD340", cp, 4*(4+24), 2, 4)
	if cache == 0 {
		t.Fatal("cache allocation")
	}
	defer call("sub_4BD3C0", cache)
	defer call("sub_486E00", cp)
	var rows []map[string]any
	check := func(index uint32, p uint32) {
		t.Helper()
		if p == 0 {
			t.Fatal("missing cache entry", index)
		}
		w := audioStreamWords(p, 21)
		if w[3] != 0 || w[4] != index || w[5] != 1 || w[13] != cache || w[7] != uint32(len(samples[index])) || call("sub_4BD710", p) != p+24 {
			t.Fatal("cache metadata", index)
		}
		var got []byte
		var lengths []uint32
		for chunk := call("sub_487C80", p+24); chunk != 0; {
			cw := audioStreamWords(chunk, 6)
			if cw[5] != p+24 || cw[3] != chunk+24 {
				t.Fatal("cached chunk owner")
			}
			got = append(got, audioStreamMemory(cw[3], int(cw[4]))...)
			lengths = append(lengths, cw[4])
			if len(lengths) > 10 {
				t.Fatal("chunk cycle")
			}
			chunk = cw[0]
			if chunk == p+32 {
				chunk = 0
			}
		}
		if !bytes.Equal(got, samples[index]) || cat[70] != 0 {
			t.Fatal("cached sample contents/close", index)
		}
		rows = append(rows, map[string]any{"loaded": index, "lengths": lengths, "bytes": got})
	}
	a := call("sub_4BD470", cache, 0)
	check(0, a)
	b := call("sub_4BD470", cache, 1)
	check(1, b)
	if call("sub_4BD470", cache, 0) != a {
		t.Fatal("cache hit identity")
	}
	c := call("sub_4BD470", cache, 2)
	check(2, c)
	if call("sub_4BD420", cache, 1) != 0 || call("sub_4BD420", cache, 0) != a {
		t.Fatal("least recently used eviction")
	}
	call("sub_4BD650", a)
	call("sub_4BD650", c)
	if call("sub_4BD470", cache, 1) != 0 || call("sub_4BD600", cache) != 0 || cat[70] != 0 {
		t.Fatal("pinned cache eviction")
	}
	call("sub_4BD660", a)
	b = call("sub_4BD470", cache, 1)
	check(1, b)
	if call("sub_4BD420", cache, 0) != 0 || call("sub_4BD420", cache, 2) != c {
		t.Fatal("pinned entry survives eviction")
	}
	call("sub_4BD660", c)
	if call("sub_4BD470", cache, 3) != 0 {
		t.Fatal("sample exceeding cache capacity should fail")
	}
	// The original failure path leaves the stream open; explicitly close it before reuse.
	rows = append(rows, map[string]any{"oversized_active": cat[70] != 0, "remaining": cat[71]})
	call("sub_486E00", cp)
	a = call("sub_4BD470", cache, 0)
	check(0, a)
	if call("sub_4BD600", cache) != 1 || call("sub_4BD600", cache) != 0 {
		t.Fatal("eviction reclaimed failed fill")
	}
	spellbookCapture(t, "client-audio-streams-cache", rows, "f29bb95cf6491c2cd0a8c4f89ddb1c9068a3f8ad184591beb52920f4ebcbe788")
}
