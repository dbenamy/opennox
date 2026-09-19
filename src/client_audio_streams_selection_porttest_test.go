//go:build porttest

package opennox

import (
	"testing"
)

func TestClientAudioStreamsVoiceSelection(t *testing.T) {
	type spec struct{ kind, flags, priority, level uint32 }
	cases := []struct {
		name   string
		voices []spec
		kind   uint32
		want   int
	}{
		{"empty", nil, 1, -1},
		{"idle", []spec{{1, 0, 99, 0}, {1, 0, 0, 0}}, 1, 0},
		{"other-kind", []spec{{2, 0, 0, 0}, {1, 0, 99, 0}}, 1, 1},
		{"default-kind", []spec{{2, 0, 0, 0}, {1, 0, 99, 0}}, 0xffffffff, 1},
		{"lower-priority", []spec{{1, 1, 30, 10000}, {1, 1, 20, 10000}}, 1, 1},
		{"equal-level", []spec{{1, 1, 20, 10000}, {1, 1, 20, 10000}}, 1, 0},
		{"level-threshold-minus-one", []spec{{1, 1, 20, 10000}, {1, 1, 20, 10000 - 0x665}}, 1, 0},
		{"level-threshold", []spec{{1, 1, 20, 10000}, {1, 1, 20, 10000 - 0x666}}, 1, 1},
		{"level-higher", []spec{{1, 1, 20, 10000}, {1, 1, 20, 20000}}, 1, 0},
		{"busy-tie", []spec{{1, 1, 20, 10000}, {1, 0x10, 20, 0}}, 1, 1},
		{"busy-higher", []spec{{1, 1, 20, 10000}, {1, 4, 21, 0}}, 1, 0},
		{"busy-lower", []spec{{1, 1, 20, 10000}, {1, 4, 19, 0}}, 1, 1},
		{"busy-equal-order", []spec{{1, 4, 20, 0}, {1, 0x10, 20, 0}}, 1, 0},
		{"priority-limit", []spec{{1, 4, 127, 0}, {1, 1, 128, 0}}, 1, -1},
		{"signed-priority", []spec{{1, 4, 0xffffffff, 0}, {1, 1, 0, 0}}, 1, 0},
		{"unrelated-flags", []spec{{1, 2, 127, 0}, {1, 1, 0, 0}}, 1, 0},
	}
	var rows []map[string]any
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			o := newAudioStreamOwner(t, 1, len(c.voices))
			if o.device() == 0 {
				t.Fatal("device")
			}
			ctx := o.call("sub_487150", 0, 0)
			var voices []uint32
			for _, s := range c.voices {
				v := o.call("sub_487750", ctx)
				if v == 0 {
					t.Fatal("voice")
				}
				w := audioStreamWords(v, 78)
				w[3], w[31], w[30], w[45] = s.kind, s.flags, s.priority, s.level
				voices = append(voices, v)
			}
			got := o.call("sub_487810", ctx, c.kind)
			want := uint32(0)
			if c.want >= 0 {
				want = voices[c.want]
			}
			if got != want {
				t.Fatalf("selected %d, want index %d", o.refs[got], c.want)
			}
			rows = append(rows, map[string]any{"name": c.name, "selected": c.want})
		})
	}
	spellbookCapture(t, "client-audio-streams-voice-selection", rows, "bb084550a6f8ed2f6339034c52808dfe092bfbe4bb37694d8b672621730ee925")
}

func TestClientAudioStreamsTickBounds(t *testing.T) {
	cases := []struct {
		last, now, period, max uint64
		busy                   uint32
	}{
		{0, 32, 33, 0, 0}, {0, 33, 33, 0, 0}, {100, 132, 33, 70, 0}, {100, 133, 33, 70, 0},
		{100, 133, 33, 331, 0}, {0xffffffff, 0x100000020, 33, 0, 0},
		{0xfffffffffffffff0, 17, 33, 0, 0}, {100, 90, 33, 0, 0},
		{0, 0, 0, 99, 0}, {0, 0x100000001, 0x100000000, 0, 0}, {0, 100, 33, 0, 1},
	}
	var rows []map[string]any
	for _, c := range cases {
		func() {
			o := newAudioStreamOwner(t, 1, 0)
			if o.device() == 0 {
				t.Fatal("device")
			}
			ctx := o.call("sub_487150", 0, 0)
			w := audioStreamWords(ctx, 66)
			put := func(i int, v uint64) { w[i], w[i+1] = uint32(v), uint32(v>>32) }
			get := func(i int) uint64 { return uint64(w[i]) | uint64(w[i+1])<<32 }
			put(56, c.period)
			put(58, 1234)
			put(60, c.max)
			put(62, c.last)
			w[53] = c.busy
			o.ticks = c.now
			beforeClock := o.clockCalls
			ret := o.call("sub_4873C0", ctx)
			elapsed, last, max := uint64(1234), c.last, c.max
			if c.busy != 0 {
				if ret != 0x80120000 || o.clockCalls != beforeClock {
					t.Fatal("busy tick")
				}
			} else {
				if ret != 0 {
					t.Fatal("tick return")
				}
				now := uint64(uint32(c.now)) // Existing C clock ABI narrows before promotion.
				delta := now - c.last
				if delta >= c.period {
					elapsed, last = delta, now
					if max > 10*c.period {
						max = 0
					}
					if delta > max {
						max = delta
					}
				}
			}
			if get(58) != elapsed || get(60) != max || get(62) != last {
				t.Fatalf("tick state: %+v got %d/%d/%d want %d/%d/%d", c, get(58), get(60), get(62), elapsed, max, last)
			}
			rows = append(rows, map[string]any{"last": c.last, "now": c.now, "period": c.period, "max_before": c.max, "busy": c.busy, "elapsed": elapsed, "max": max, "last_after": last, "return": ret})
			w[53] = 0
			o.call("sub_4875F0")
			o.call("sub_4870A0")
		}()
	}
	spellbookCapture(t, "client-audio-streams-tick-bounds", rows, "143e11db0194c3588315a61d5b466291ccbfc6f3441423be97a2b4b414d4992b")
}
