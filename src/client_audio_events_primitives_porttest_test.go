//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/timer"
	"reflect"
	"testing"
	"unsafe"
)

func audioEventPointer(p unsafe.Pointer) uint64 { return uint64(uintptr(p)) }
func audioEventSigned(v int32) uint64           { return uint64(int64(v)) }

func TestClientAudioEventsFormats(t *testing.T) {
	p, free := alloc.Make([]uint32{}, 8)
	defer free()
	ptr := audioEventPointer(unsafe.Pointer(&p[0]))
	var rows []map[string]any
	for _, encoding := range []uint32{0, 1, 2, 3, 0xffffffff} {
		for _, channels := range []uint32{0, 1, 2, 3} {
			for _, width := range []uint32{0, 1, 2, 3} {
				for i := range p {
					p[i] = 0x12340000 + uint32(i)
				}
				p[1], p[3], p[4] = encoding, channels, width
				before := append([]uint32(nil), p...)
				want := uint64(0)
				if encoding == 2 {
					want = 5
					if channels == 2 {
						want = 7
					}
				} else if encoding == 0 {
					if width == 2 {
						want = 1
					}
					if channels == 2 {
						want += 2
					}
				}
				got := legacy.PortTestAudioEventCall("sub_43F0E0", ptr)
				if got != want || !reflect.DeepEqual(p, before) {
					t.Fatal("AIL format code/unchanged metadata", encoding, channels, width, got, want)
				}
				rows = append(rows, map[string]any{"encoding": encoding, "channels": channels, "width": width, "format": got})
			}
		}
	}
	spellbookCapture(t, "client-audio-events-formats", rows, "eb436ababe23e799cb8b3706c718d37944e984a76198642cf04c7eb6dabab366")
}

func TestClientAudioEventsVolumePan(t *testing.T) {
	event, free := alloc.Make([]uint32{}, 144)
	defer free()
	meta, free := alloc.Make([]uint32{}, 50)
	defer free()
	ep, mp := audioEventPointer(unsafe.Pointer(&event[0])), audioEventPointer(unsafe.Pointer(&meta[0]))
	event[9] = uint32(mp)
	oldClock := timer.PlatformTicks
	timer.PlatformTicks = func() uint64 { return 100 }
	defer func() { timer.PlatformTicks = oldClock }()
	volume := (*timer.Timer)(unsafe.Pointer(&event[46]))
	pan := (*timer.Timer)(unsafe.Pointer(&event[62]))
	var rows []map[string]any
	for _, input := range []int32{-2147483648, -101, -51, -50, -1, 0, 1, 49, 50, 51, 99, 100, 101, 2147483647} {
		clamped := input
		if clamped < -50 {
			clamped = -50
		}
		if clamped > 50 {
			clamped = 50
		}
		wantPan := uint32(clamped*8192/50 + 8192)
		got := legacy.PortTestAudioEventCall("sub_452FA0", audioEventSigned(input))
		if got != uint64(wantPan) {
			t.Fatal("pan clamp", input, got, wantPan)
		}
		pan.Init(0)
		legacy.PortTestAudioEventCall("sub_452F80", ep, audioEventSigned(input))
		if pan.Target != wantPan<<16 || pan.Current != 0 || pan.Flags&1 == 0 {
			t.Fatal("pan raw setter")
		}
		pan.Init(0)
		legacy.PortTestAudioEventCall("sub_452FE0", ep, audioEventSigned(input))
		if pan.Target != wantPan<<16 || pan.Flags&1 != 0 {
			t.Fatal("pan interpolated setter")
		}
		for _, raw := range []uint32{0, 1, 0xffff, 0x10000, 0x40000000, 0x80000000, 0xffffffff} {
			meta[5] = raw
			percent := input
			if percent < 0 {
				percent = 0
			}
			if percent > 100 {
				percent = 100
			}
			want := (uint32(163) * uint32(percent) * (raw >> 16)) >> 14
			got := legacy.PortTestAudioEventCall("sub_452F10", ep, audioEventSigned(input))
			if got != uint64(want) {
				t.Fatal("volume clamp/product", input, raw, got, want)
			}
			volume.Init(0)
			ret := legacy.PortTestAudioEventCall("sub_452EE0", ep, audioEventSigned(input))
			wantRet := uint64(0)
			if want != 0 {
				wantRet = 1
			}
			if ret != wantRet || volume.Current != want<<16 || volume.Target != want<<16 {
				t.Fatal("volume raw update", ret, wantRet)
			}
			volume.Init(0)
			legacy.PortTestAudioEventCall("sub_452F50", ep, audioEventSigned(input))
			if volume.Target != want<<16 || volume.Current != 0 || volume.Flags&1 != 0 {
				t.Fatal("volume interpolation")
			}
			rows = append(rows, map[string]any{"input": input, "source": raw, "volume": want, "pan": wantPan})
		}
	}
	spellbookCapture(t, "client-audio-events-volume-pan", rows, "fe0eb2ad08d6994b9917a47e3609564f4e3f6df6dca2a5a29490d17ca17436dd")
}

func TestClientAudioEventsSwitches(t *testing.T) {
	words, restore := legacy.PortTestAudioEventGlobals()
	defer restore()
	mapped := serverConfigOwnBytes(t, 0x5D4594, 831252, 1)
	var rows []map[string]any
	for _, driver := range []uint32{0, 1, 0x7fffffff, 0x80000000, 0xffffffff} {
		for _, enabled := range []uint32{0, 1, 17, 0xffffffff} {
			for _, pair := range []struct{ driver, enabled, on, off, get string }{
				{"dword_5d4594_816376", "dword_587000_93156", "sub_43DC10", "sub_43DC00", "sub_43DC30"},
				{"dword_5d4594_831092", "dword_587000_122848", "sub_44D970", "sub_44D960", "sub_44D990"},
			} {
				*words[pair.driver] = driver
				*words[pair.enabled] = enabled
				got := legacy.PortTestAudioEventCall(pair.on)
				want := enabled
				if driver != 0 {
					want = 1
				}
				if got != audioEventSigned(int32(driver)) || *words[pair.enabled] != want || legacy.PortTestAudioEventCall(pair.get) != audioEventSigned(int32(want)) {
					t.Fatal("enable requires driver", pair.on)
				}
				legacy.PortTestAudioEventCall(pair.off)
				if *words[pair.enabled] != 0 || *words[pair.driver] != driver {
					t.Fatal("disable retains driver")
				}
				rows = append(rows, map[string]any{"on": pair.on, "driver": driver, "before": enabled, "enabled": want, "return": got})
			}
		}
	}
	for value := 0; value < 256; value++ {
		mapped[0] = 0
		ret := legacy.PortTestAudioEventCall("sub_450760", uint64(value))
		if ret != uint64(int64(int8(value))) || legacy.PortTestAudioEventCall("sub_450750") != uint64(value) || mapped[0] != byte(value) {
			t.Fatal("byte switch signed return", value)
		}
		rows = append(rows, map[string]any{"byte": value, "return": ret})
	}
	*words["dword_587000_126996"] = 17
	legacy.PortTestAudioEventCall("sub_453050")
	if legacy.PortTestAudioEventCall("sub_453070") != 0 {
		t.Fatal("playback disable")
	}
	legacy.PortTestAudioEventCall("nox_xxx____setargv_9_453060")
	if legacy.PortTestAudioEventCall("sub_453070") != 1 {
		t.Fatal("playback enable")
	}
	spellbookCapture(t, "client-audio-events-switches", rows, "d9a0e015f39022e53fae462ae0516a2a8b40d51667d38eb811d1d7dc4bbbde0c")
}
