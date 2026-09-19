//go:build porttest

package opennox

import (
	"fmt"
	"testing"
)

func TestClientAudioStreamsLifecycle(t *testing.T) {
	var rows []map[string]any
	for _, slots := range []int{1, 2, 16} {
		for _, voices := range []int{0, 1, 4} {
			t.Run(fmt.Sprintf("slots%d/voices%d", slots, voices), func(t *testing.T) {
				o := newAudioStreamOwner(t, slots, voices)
				o.descriptor[2] = 2
				device := o.device()
				if device == 0 || o.descriptor[3] != 1 || o.count() != 1 {
					t.Fatal("device not registered")
				}
				dw := audioStreamWords(device, 22)
				for slot := 0; slot < slots; slot++ {
					ctx := o.call("sub_487150", uint32(slot), 0)
					if ctx == 0 {
						t.Fatal("missing context", slot)
					}
					cw := audioStreamWords(ctx, 66)
					if cw[4] != 1 || cw[5] != device || cw[6] != uint32(slot) || cw[47] != uint32(slot) || cw[49] != uint32(voices) || cw[56] != 33 || dw[6+slot] != ctx {
						t.Fatal("context ownership", slot)
					}
					if o.call("sub_487150", uint32(slot), 0) != ctx || cw[4] != 2 {
						t.Fatal("context reuse", slot)
					}
					if slot == 0 && (o.call("sub_487150", ^uint32(0), 0) != ctx || cw[4] != 3) {
						t.Fatal("default context alias")
					}
					for i := 0; i < voices; i++ {
						voice := o.call("sub_487750", ctx)
						if voice == 0 {
							t.Fatal("missing voice", i)
						}
						vw := audioStreamWords(voice, 78)
						if vw[33] != ctx || vw[3] != 1 || cw[48] != uint32(i+1) || cw[53] != 0 {
							t.Fatal("voice ownership", i)
						}
					}
					if o.call("sub_487750", ctx) != 0 || cw[48] != uint32(voices) {
						t.Fatal("voice capacity exceeded")
					}
				}
				if dw[4] != uint32(slots) || o.call("sub_487150", uint32(slots), 0) != 0 {
					t.Fatal("context capacity")
				}
				o.call("sub_4875F0")
				if dw[4] != 0 {
					t.Fatal("contexts still owned")
				}
				for _, p := range dw[6 : 6+slots] {
					if p != 0 {
						t.Fatal("stale context slot")
					}
				}
				o.call("sub_4870A0")
				if o.count() != 0 || o.descriptor[3] != 0 || o.root[6] != 0 {
					t.Fatal("cleanup state")
				}
				if o.callbackCounts[0] != 1 || o.callbackCounts[1] != 1 || o.callbackCounts[2] != slots || o.callbackCounts[3] != slots || o.callbackCounts[4] != slots*voices || o.callbackCounts[5] != slots*voices {
					t.Fatal("callback lifetime", o.callbackCounts)
				}
				rows = append(rows, map[string]any{"slots": slots, "voices": voices, "events": o.events, "counts": o.callbackCounts})
			})
		}
	}
	spellbookCapture(t, "client-audio-streams-lifecycle", rows, "9b993b6ff2a581ab0fac6b509b5fd9787415941c333faf3c62c20558691aa38a")
}

func TestClientAudioStreamsConstructorFailures(t *testing.T) {
	var rows []map[string]any
	for _, op := range []int{0, 2, 4} {
		t.Run(fmt.Sprint(op), func(t *testing.T) {
			o := newAudioStreamOwner(t, 2, 3)
			o.failOp = op
			device := o.device()
			if op == 0 {
				if device != 0 || o.descriptor[3] != 0 {
					t.Fatal("failed device retained")
				}
			} else {
				if device == 0 {
					t.Fatal("device failed")
				}
				ctx := o.call("sub_487150", 0, 0)
				if op == 2 {
					if ctx != 0 || audioStreamWords(device, 22)[4] != 0 || audioStreamWords(device, 22)[6] != 0 {
						t.Fatal("failed context retained")
					}
				} else {
					if ctx == 0 {
						t.Fatal("context failed")
					}
					if o.call("sub_487750", ctx) != 0 || audioStreamWords(ctx, 66)[48] != 0 {
						t.Fatal("failed voice retained")
					}
				}
			}
			if o.callbackCounts[op] != 1 || o.callbackCounts[op+1] != 1 {
				t.Fatal("failed initializer cleanup", o.callbackCounts)
			}
			o.failOp = -1
			o.call("sub_4875F0")
			o.call("sub_4870A0")
			rows = append(rows, map[string]any{"fail_op": op, "events": o.events, "counts": o.callbackCounts})
		})
	}
	spellbookCapture(t, "client-audio-streams-constructor-failures", rows, "757a0cd5cc6ed7234a6ae3031c05d007d221db239c7591e365abb19202deae16")
}
