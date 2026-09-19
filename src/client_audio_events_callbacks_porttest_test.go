//go:build porttest

package opennox

import "testing"

func TestClientAudioEventsCallbacks(t *testing.T) {
	var rows []map[string]any
	for _, state := range []uint32{1, 3, 4} {
		for _, pending := range []bool{false, true} {
			for _, remaining := range []uint32{0, 2} {
				for _, delay := range []uint32{0, 33, 100} {
					t.Run(string(rune('a'+len(rows))), func(t *testing.T) {
						o := newAudioEventsOwner(t, 1)
						p := o.create(o.metadata(1, 0, 1))
						w := audioStreamWords(p, 144)
						if o.eventCall("sub_452580", uint64(p)) != 1 {
							t.Fatal("reserve")
						}
						vp := w[44]
						w[7] = state
						w[6] = 2
						w[142] = remaining
						w[71] = delay
						if pending {
							w[74] = 123
						}
						o.ticks = 100
						wantState, wantDelay, wantPending := state, delay, w[74]
						deadline := uint64(0)
						if state != 4 {
							next := uint32(4)
							if pending || remaining != 0 {
								next = 1
							} else {
								wantDelay = 0
							}
							if wantDelay != 0 {
								wantState = 2
								deadline = 100 + uint64(wantDelay)
								wantDelay = 0
							} else {
								wantState = next
							}
						}
						if o.eventCall("sub_4526F0", uint64(vp)) != 0 || w[6] != 0 || w[7] != wantState || w[71] != wantDelay || w[74] != wantPending || uint64(w[72])|uint64(w[73])<<32 != deadline {
							t.Fatal("end callback", state, pending, remaining, delay, w[7], wantState, w[71], wantDelay)
						}
						if wantState == 2 && w[8] != 1 {
							t.Fatal("scheduled next state")
						}
						rows = append(rows, map[string]any{"state": state, "pending": pending, "remaining": remaining, "delay": delay, "result_state": w[7], "result_delay": w[71], "deadline": deadline})
						if o.eventCall("sub_4526D0", uint64(vp)) != 0 || w[7] != 4 {
							t.Fatal("stop callback")
						}
						w[74] = 0 // Synthetic pending marker is not a sample buffer.
					})
				}
			}
		}
	}
	spellbookCapture(t, "client-audio-events-callbacks", rows, "ad422347fc17ee50af583c31bb66bb469cb223fbdcb29bfdf7063a309fe7c9f5")
}

func TestClientAudioEventsCacheSelection(t *testing.T) {
	var rows []map[string]any
	for _, flags := range []uint32{0, 2, 4, 6} {
		for _, delay := range []uint32{0, 32, 33, 100} {
			t.Run(string(rune('a'+len(rows))), func(t *testing.T) {
				o := newAudioEventsOwner(t, 1)
				meta := o.metadata(1, flags, 3)
				audioStreamWords(meta, 50)[17] = delay
				p := o.create(meta)
				w := audioStreamWords(p, 144)
				var selected [][]int
				for repeat := 0; repeat < 6; repeat++ {
					got := o.eventCall("sub_451DC0", uint64(p))
					want := uint32(1)
					if flags&4 != 0 && delay < 33 {
						want = 3
					}
					if got != uint64(want) || w[42] != want {
						t.Fatal("loaded count", flags, delay, got, w[42], want)
					}
					var indices []int
					for i := uint32(0); i < w[42]; i++ {
						entry := audioStreamWords(w[10+i], 21)
						index := int(entry[4])
						if index < 0 || index >= 3 {
							t.Fatal("sample index", index)
						}
						indices = append(indices, index)
						if entry[3] != 1 {
							t.Fatal("cache reference")
						}
					}
					if flags&6 == 0 && indices[0] != 0 {
						t.Fatal("first sample")
					}
					if flags&4 != 0 && flags&2 == 0 && delay >= 33 && indices[0] != repeat%3 {
						t.Fatal("reload sequence", indices, repeat)
					}
					selected = append(selected, indices)
				}
				n := w[42]
				if o.eventCall("sub_451F90", uint64(p)) != uint64(n) || w[42] != 0 || o.eventCall("sub_451F90", uint64(p)) != 0 {
					t.Fatal("release")
				}
				for _, v := range w[10:42] {
					if v != 0 {
						t.Fatal("cleared entry")
					}
				}
				rows = append(rows, map[string]any{"flags": flags, "delay": delay, "selected": selected, "remaining": w[142]})
			})
		}
	}
	spellbookCapture(t, "client-audio-events-cache-selection", rows, "15375124c1d2487d11826cc044e17119ae8b245e38e94ec34d922989340428ce")
}
