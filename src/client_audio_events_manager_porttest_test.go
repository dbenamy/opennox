//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/common/memmap"
	"testing"
)

func TestClientAudioEventsPool(t *testing.T) {
	o := newAudioEventsOwner(t, 1)
	meta := o.metadata(1, 0, 1)
	seen := make(map[uint32]bool)
	var events []uint32
	for i := 0; i < 200; i++ {
		p := o.create(meta)
		if p == 0 || seen[p] {
			t.Fatal("pool allocation", i, p)
		}
		seen[p] = true
		events = append(events, p)
		if audioStreamWords(p, 144)[70] != uint32(i) {
			t.Fatal("serial", i)
		}
	}
	if o.create(meta) != 0 {
		t.Fatal("pool capacity")
	}
	for i, p := range events {
		if i%3 == 0 {
			o.eventCall("sub_4523D0", uint64(p))
			delete(seen, p)
		}
	}
	for i := 0; i < 67; i++ {
		p := o.create(meta)
		if p == 0 || seen[p] {
			t.Fatal("automatic reclaim", i, p)
		}
		seen[p] = true
		if audioStreamWords(p, 144)[70] != uint32(200+i) {
			t.Fatal("reused serial", i)
		}
	}
	if o.create(meta) != 0 {
		t.Fatal("refilled pool capacity")
	}
	o.eventCall("sub_4521F0")
	for p := range seen {
		w := audioStreamWords(p, 144)
		if w[7] != 4 || w[6]&1 == 0 {
			t.Fatal("stop all")
		}
	}
	o.eventCall("sub_452230")
	root := uint32(uintptr(memmap.PtrOff(0x5D4594, 840612)))
	if audioStreamWords(root, 3)[0] != root || audioStreamWords(root, 3)[1] != root {
		t.Fatal("empty active list")
	}
	spellbookCapture(t, "client-audio-events-pool", []map[string]any{{"capacity": 200, "reclaimed": 67, "next_serial": *memmap.PtrUint32(0x587000, 127000), "empty": true}}, "ce0e76b6349c21ffb2dd4353f72f6b6d46178b433b2929eec6b98bc5856e2272")
}

func TestClientAudioEventsManager(t *testing.T) {
	var rows []map[string]any
	for _, count := range []int{1, 2, 4} {
		for _, cap := range []uint32{0, 1, 2} {
			t.Run(string(rune('0'+count))+"-"+string(rune('0'+cap)), func(t *testing.T) {
				o := newAudioEventsOwner(t, count)
				meta := o.metadata(1, 4, 3)
				mw := audioStreamWords(meta, 50)
				mw[14] = cap
				var events []uint32
				for i := 0; i < count; i++ {
					p := o.create(meta)
					if p == 0 {
						t.Fatal("create")
					}
					events = append(events, p)
					o.eventCall("sub_452EE0", uint64(p), uint64(100-i*20))
				}
				frame := memmap.PtrUint32(0x5D4594, 1045440)
				busy := memmap.PtrUint32(0x5D4594, 1045448)
				*busy = 1
				o.eventCall("sub_4519C0")
				if *frame != 0 || *busy != 1 {
					t.Fatal("busy guard")
				}
				*busy = 0
				*o.words["dword_5d4594_1045432"] = 0
				o.eventCall("sub_4519C0")
				if *frame != 0 {
					t.Fatal("enabled guard")
				}
				*o.words["dword_5d4594_1045432"] = 1
				o.eventCall("sub_4519C0")
				want := count
				if cap > 0 && int(cap) < want {
					want = int(cap)
				}
				root := uint32(uintptr(memmap.PtrOff(0x5D4594, 840612)))
				live := 0
				for p := audioStreamWords(root, 3)[0]; p != root; p = audioStreamWords(p, 3)[0] {
					live++
					if live > count {
						t.Fatal("list cycle")
					}
					w := audioStreamWords(p, 144)
					if w[7] != 3 || w[6]&2 == 0 || w[42] != 3 {
						t.Fatal("manager playback", w[7], w[6], w[42])
					}
				}
				if live != want || mw[13] != uint32(want) || *frame != 1 || *busy != 0 {
					t.Fatal("manager counts", live, want, mw[13], *frame, *busy)
				}
				o.ticks += 100
				o.eventCall("sub_4519C0")
				if *frame != 2 || *busy != 0 {
					t.Fatal("second frame")
				}
				rows = append(rows, map[string]any{"events": count, "cap": cap, "live": live, "frame": *frame, "active_count": mw[13]})
			})
		}
	}
	spellbookCapture(t, "client-audio-events-manager", rows, "a008535113388e3e3698aef54bcadc10b40e846bee4a76484720f8005b5fa98a")
}
