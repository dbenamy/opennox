//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/common/memmap"
	"testing"
)

func TestClientAudioEventsPublicPlayback(t *testing.T) {
	var rows []map[string]any
	for _, name := range []string{"nox_xxx_clientPlaySoundSpecial_452D80", "sub_452DC0", "sub_452E10"} {
		for _, volume := range []int32{-1, 0, 50, 100, 101} {
			for _, pan := range []int32{-51, 0, 51} {
				t.Run(string(rune('a'+len(rows))), func(t *testing.T) {
					o := newAudioEventsOwner(t, 1)
					o.metadata(1, 0, 1)
					o.eventCall(name, 1, audioEventSigned(volume), audioEventSigned(pan))
					root := uint32(uintptr(memmap.PtrOff(0x5D4594, 840612)))
					p := audioStreamWords(root, 3)[0]
					if p == root {
						t.Fatal("missing event")
					}
					w := audioStreamWords(p, 144)
					v := volume
					if v < 0 {
						v = 0
					}
					if v > 100 {
						v = 100
					}
					wantVolume := uint32(v*163) << 16
					wantPan := int32(8192)
					if name != "nox_xxx_clientPlaySoundSpecial_452D80" {
						x := pan
						if x < -50 {
							x = -50
						}
						if x > 50 {
							x = 50
						}
						wantPan += x * 8192 / 50
					}
					priority := uint32(0)
					if name == "sub_452E10" {
						priority = 2
					}
					if w[7] != 1 || w[47] != wantVolume || w[63] != 8192<<16 || w[64] != uint32(wantPan)<<16 || w[75] != priority {
						t.Fatal("public playback", name, volume, pan, w[7], w[47], w[63], w[75])
					}
					rows = append(rows, map[string]any{"name": name, "volume": volume, "pan": pan, "volume_raw": w[47], "pan_current": w[63], "pan_target": w[64], "priority_offset": w[75], "state": w[7]})
				})
			}
		}
	}
	spellbookCapture(t, "client-audio-events-public-playback", rows, "e84346df2fb66f8af0b5f5691f8f415244e46440a034a734673f994e55ac69ef")
}
