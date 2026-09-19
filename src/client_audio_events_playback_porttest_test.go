//go:build porttest

package opennox

import "testing"

func TestClientAudioEventsPlaybackFailures(t *testing.T) {
	var rows []map[string]any
	for _, mode := range []string{"no-samples", "no-context", "higher-priority", "start-failure", "wrong-owner", "initial-delay"} {
		t.Run(mode, func(t *testing.T) {
			o := newAudioEventsOwner(t, 1)
			meta := o.metadata(1, 0, 1)
			mw := audioStreamWords(meta, 50)
			p := o.create(meta)
			w := audioStreamWords(p, 144)
			switch mode {
			case "no-samples":
				mw[48] = 0
			case "no-context":
				*o.words["dword_5d4594_1045428"] = 0
			case "higher-priority":
				v := o.call("sub_487810", o.ctx, 1)
				vw := audioStreamWords(v, 39)
				vw[30] = 5
				vw[31] = 16
			case "initial-delay":
				mw[1] = 8
				mw[17] = 34
				mw[18] = 34
			}
			o.eventCall("sub_452510", uint64(p))
			if mode == "no-samples" || mode == "no-context" || mode == "higher-priority" {
				if w[7] != 4 || w[6]&1 == 0 {
					t.Fatal("failed reservation", w[7], w[6])
				}
			} else if mode == "initial-delay" {
				if w[7] != 2 || w[72] != 134 || w[8] != 1 {
					t.Fatal("initial delay", w[7], w[72])
				}
			} else {
				if w[7] != 1 || w[44] == 0 {
					t.Fatal("reserve")
				}
				o.eventCall("sub_451DC0", uint64(p))
				w[74] = uint32(o.eventCall("sub_451CA0", uint64(p)))
				pending := w[74]
				vp := w[44]
				vw := audioStreamWords(vp, 39)
				if mode == "start-failure" {
					o.failOp = 6
				} else {
					vw[38] = 0
				}
				if o.eventCall("sub_452490", uint64(p)) != 0 || w[7] != 1 || w[6]&2 != 0 || w[74] != pending {
					t.Fatal("start rollback")
				}
				vw[38] = p
				o.failOp = -1
				if o.eventCall("sub_452490", uint64(p)) != 1 || w[7] != 3 {
					t.Fatal("retry")
				}
			}
			rows = append(rows, map[string]any{"mode": mode, "state": w[7], "flags": w[6], "next_state": w[8], "deadline": uint64(w[72]) | uint64(w[73])<<32})
		})
	}
	spellbookCapture(t, "client-audio-events-playback-failures", rows, "ffbfab9bb7a6c9c741e2a51b86e35847521835be5038fbab93e7234b4afad0ae")
}

func TestClientAudioEventsLoop(t *testing.T) {
	var rows []map[string]any
	for _, flags := range []uint32{0, 1, 8, 9} {
		for _, delay := range []uint32{0, 32, 33, 100} {
			t.Run(string(rune('a'+len(rows))), func(t *testing.T) {
				o := newAudioEventsOwner(t, 1)
				meta := o.metadata(1, flags, 1)
				mw := audioStreamWords(meta, 50)
				mw[17] = delay
				mw[18] = delay
				p := o.create(meta)
				w := audioStreamWords(p, 144)
				if o.eventCall("sub_452580", uint64(p)) != 1 {
					t.Fatal("reserve")
				}
				o.eventCall("sub_451DC0", uint64(p))
				first := uint32(o.eventCall("sub_451CA0", uint64(p)))
				vp := w[44]
				o.call("sub_4BDB90", vp, first)
				if o.eventCall("sub_452770", uint64(vp)) != 0 {
					t.Fatal("loop result")
				}
				vw := audioStreamWords(vp, 78)
				looping := flags&1 != 0
				wantBuffer := looping && delay < 33
				wantPending := looping && delay >= 33
				wantDelay := uint32(0)
				if delay >= 33 && (flags&8 == 0 || looping) {
					wantDelay = delay
				}
				if (vw[72] != 0) != wantBuffer || (w[74] != 0) != wantPending || w[71] != wantDelay {
					t.Fatal("loop scheduling", flags, delay, vw[72] != 0, w[74] != 0, w[71])
				}
				rows = append(rows, map[string]any{"flags": flags, "delay": delay, "buffer": vw[72] != 0, "pending": w[74] != 0, "scheduled_delay": w[71]})
			})
		}
	}
	spellbookCapture(t, "client-audio-events-loop", rows, "682112c39a1d22f526674356e36b70460dbc5b42d47ed3f47f90f55c71ddb8cf")
}
