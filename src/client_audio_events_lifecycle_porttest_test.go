//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
	"unsafe"
)

func TestClientAudioEventsLifecycle(t *testing.T) {
	o := newAudioEventsOwner(t, 2)
	meta := o.metadata(1, 4, 3)
	handle, free := alloc.Make([]uint32{}, 3)
	defer free()
	hp := audioEventPointer(unsafe.Pointer(&handle[0]))
	var rows []map[string]any
	for _, gate := range []string{"dword_5d4594_1045432", "dword_587000_126996"} {
		*o.words[gate] = 0
		if o.create(meta) != 0 {
			t.Fatal("disabled creation", gate)
		}
		*o.words[gate] = 1
	}
	mw := audioStreamWords(meta, 50)
	mw[0] = 0
	if o.create(meta) != 0 {
		t.Fatal("disabled metadata")
	}
	mw[0] = 1
	for i := 0; i < 3; i++ {
		p := o.create(meta)
		if p == 0 {
			t.Fatal("create", i)
		}
		w := audioStreamWords(p, 144)
		if w[9] != meta || w[7] != 0 || w[6] != 0 || w[70] != uint32(i) {
			t.Fatal("initial event", i, w[6:10], w[70])
		}
		if uint32(o.eventCall("sub_452E90", hp, uint64(p))) != meta || uint32(o.eventCall("sub_452EB0", hp)) != p {
			t.Fatal("live handle")
		}
		o.eventCall("sub_452EE0", uint64(p), 100)
		o.eventCall("sub_452510", uint64(p))
		if w[7] != 1 || w[44] == 0 {
			t.Fatal("reserve voice", w[7], w[44])
		}
		voicePtr := w[44]
		voice := audioStreamWords(voicePtr, 39)
		if voice[38] != p || voice[28] != p+184 || voice[35] == 0 || voice[36] == 0 || voice[37] == 0 {
			t.Fatal("voice ownership")
		}
		if o.eventCall("sub_451DC0", uint64(p)) != 3 || w[42] != 3 {
			t.Fatal("load samples", w[42])
		}
		for j := 0; j < 3; j++ {
			if w[10+j] == 0 {
				t.Fatal("missing cache entry", j)
			}
		}
		loaded := w[42]
		if o.eventCall("sub_451DC0", uint64(p)) != 3 || w[42] != 3 {
			t.Fatal("idempotent load")
		}
		w[74] = uint32(o.eventCall("sub_451CA0", uint64(p)))
		pending := w[74]
		start := o.eventCall("sub_452490", uint64(p))
		if pending == 0 || start != 1 || w[7] != 3 || w[6]&2 == 0 || w[74] != 0 {
			t.Fatal("start playback", w[7], w[6], "pending", pending, "return", start, "voice flags", voice[31], "callbacks", o.events)
		}
		if o.eventCall("sub_4523D0", uint64(p))&1 == 0 || w[7] != 4 || w[44] != 0 || w[42] != 0 || w[70] != 0 {
			t.Fatal("stop event")
		}
		if voice[38] != 0 || voice[28] != 0 || voice[35] != 0 || voice[36] != 0 || voice[37] != 0 {
			t.Fatal("detach voice")
		}
		// The recording device does not deliver completion asynchronously.
		// Deliver the real stream end notification before reusing the voice.
		o.call("sub_4BD9B0", voicePtr)
		if voice[31]&5 != 0 {
			t.Fatal("completed voice still active")
		}
		if o.eventCall("sub_4523D0", uint64(p)) != 0 {
			t.Fatal("repeat stop")
		}
		// Serial zero is also used by the first event: preserve the existing handle convention.
		valid := uint32(o.eventCall("sub_452EB0", hp)) == p
		if valid != (i == 0) {
			t.Fatal("stopped handle serial", i, valid)
		}
		rows = append(rows, map[string]any{"serial": i, "loaded": loaded, "state": w[7], "flags": w[6], "stopped_handle_valid": valid})
		o.eventCall("sub_452230")
	}
	spellbookCapture(t, "client-audio-events-lifecycle", rows, "f4a1d69b5487f298ae3956c12ff445616b8a4abbdd30536b01a265b7bd5dbbb9")
}

func TestClientAudioEventsScheduling(t *testing.T) {
	o := newAudioEventsOwner(t, 1)
	meta := o.metadata(1, 0, 1)
	var rows []map[string]any
	for _, now := range []uint64{0, 100, 0xffffffff, 0x100000005} {
		for _, delay := range []uint64{0, 1, 33, 0xffffffff, 0xffffffffffffffff} {
			p := o.create(meta)
			if p == 0 {
				t.Fatal("create")
			}
			w := audioStreamWords(p, 144)
			o.ticks = now
			want := uint64(uint32(now)) + delay
			if got := o.eventCall("sub_452690", uint64(p), delay, 1); got != want || w[7] != 2 || w[8] != 1 || uint64(w[72])|uint64(w[73])<<32 != want {
				t.Fatal("deadline", now, delay, got, want)
			}
			o.eventCall("sub_452510", uint64(p))
			state := uint32(2)
			if uint64(uint32(now)) > want {
				state = 1
			}
			if w[7] != state {
				t.Fatal("deadline comparison", now, delay, w[7], state)
			}
			rows = append(rows, map[string]any{"now": now, "delay": delay, "deadline": want, "state": w[7]})
			o.eventCall("sub_4523D0", uint64(p))
			o.eventCall("sub_452230")
		}
	}
	// The deadline itself still waits; the following tick advances.
	p := o.create(meta)
	w := audioStreamWords(p, 144)
	o.ticks = 100
	o.eventCall("sub_452690", uint64(p), 33, 1)
	for _, now := range []uint64{132, 133, 134} {
		o.ticks = now
		o.eventCall("sub_452510", uint64(p))
		want := uint32(2)
		if now > 133 {
			want = 1
		}
		if w[7] != want {
			t.Fatal("strict deadline", now, w[7])
		}
		rows = append(rows, map[string]any{"boundary_tick": now, "state": w[7]})
	}
	*o.words["dword_587000_126996"] = 0
	o.eventCall("sub_452510", uint64(p))
	if w[7] != 4 || w[6]&1 == 0 {
		t.Fatal("disabled playback cleanup")
	}
	*o.words["dword_587000_126996"] = 1
	spellbookCapture(t, "client-audio-events-scheduling", rows, "4c4dd7452358569dd09ee7f4ff20bfda1414311150f39742bd859134c5a61b60")
}
