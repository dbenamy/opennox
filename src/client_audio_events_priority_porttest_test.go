//go:build porttest

package opennox

import "testing"

func TestClientAudioEventsPriority(t *testing.T) {
	o := newAudioEventsOwner(t, 1)
	var rows []map[string]any
	for priority := uint32(0); priority < 6; priority++ {
		for _, volume := range []uint64{0, 10, 50, 99, 100} {
			o.eventCall("sub_452010")
			meta := o.metadata(1, 0, 1)
			mw := audioStreamWords(meta, 50)
			mw[12] = priority
			p := o.create(meta)
			o.eventCall("sub_452EE0", uint64(p), volume)
			o.eventCall("sub_452050", uint64(p))
			bucket := uint32(163*volume) / 0x666
			if mw[27] != priority || mw[31] != bucket {
				t.Fatal("bucket", priority, volume, mw[27], mw[31], bucket)
			}
			for limit := uint64(0); limit <= 6; limit++ {
				got := uint32(o.eventCall("sub_4521A0", limit))
				want := uint32(0)
				if limit > uint64(priority) {
					want = meta
				}
				if got != want {
					t.Fatal("priority selection", priority, limit, got, want)
				}
			}
			rows = append(rows, map[string]any{"priority": priority, "volume": volume, "bucket": bucket})
			o.eventCall("sub_452190", uint64(meta))
			if o.eventCall("sub_4521A0", 6) != 0 {
				t.Fatal("removed bucket")
			}
			o.eventCall("sub_4523D0", uint64(p))
			o.eventCall("sub_452230")
		}
	}
	// Metadata with several events has one bucket and is evicted as a group.
	o.eventCall("sub_452010")
	low := o.metadata(1, 0, 1)
	high := o.metadata(2, 0, 1)
	audioStreamWords(low, 50)[12] = 1
	audioStreamWords(high, 50)[12] = 3
	for i := 0; i < 3; i++ {
		p := o.create(low)
		o.eventCall("sub_452EE0", uint64(p), uint64(20+i*30))
		o.eventCall("sub_452050", uint64(p))
	}
	p := o.create(high)
	o.eventCall("sub_452EE0", uint64(p), 100)
	o.eventCall("sub_452050", uint64(p))
	if o.eventCall("sub_452120", uint64(p)) != 1 || o.eventCall("sub_452120", uint64(p)) != 0 {
		t.Fatal("eviction group")
	}
	if got := uint32(o.eventCall("sub_4521A0", 4)); got != high {
		t.Fatal("higher priority retained")
	}
	rows = append(rows, map[string]any{"group_evicted": 3, "retained_priority": 3})
	spellbookCapture(t, "client-audio-events-priority", rows, "b4dbccbbb8183c8193a5421f01b30e90cb222d08b08a4723aad955ffb0437baa")
}
