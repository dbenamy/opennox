//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
	"unsafe"
)

func serverConfigOwnBytes(t *testing.T, base, off uintptr, n int) []byte {
	t.Helper()
	buf := unsafe.Slice((*byte)(memmap.PtrOff(base, off)), n)
	old := append([]byte(nil), buf...)
	t.Cleanup(func() { copy(buf, old) })
	return buf
}
func TestServerConfigScalarStorage(t *testing.T) {
	o := newServerOptionsOwner(t)
	serverConfigOwnBytes(t, 0x5D4594, 3512, 4)
	serverConfigOwnBytes(t, 0x5D4594, 3584, 8)
	serverConfigOwnBytes(t, 0x5D4594, 371700, 8)
	mode := serverConfigOwnBytes(t, 0x587000, 4652, 4)
	type row struct {
		Op           string
		Value        int
		Result, Read uint64
		Updated      uint32
	}
	var rows []row
	defer noxflags.PortTestGameFlags(0)()
	for _, pair := range [][2]string{{"limit-set", "limit-get"}, {"3512-set", "3512-get"}, {"respawn-set", "respawn-get"}, {"rate-dirty-set", "rate-dirty-get"}, {"acquired-set", "acquired-get"}} {
		for _, value := range []int{-2147483648, -1, 0, 1, 32, 255, 65536, 2147483647} {
			for repeat := 0; repeat < 2; repeat++ {
				*o.optionWords["settings-updated"] = 0
				result := legacy.PortTestServerConfigScalar(pair[0], value, 0)
				read := legacy.PortTestServerConfigScalar(pair[1], 0, 0)
				want := uint64(int64(value))
				if read != want {
					t.Fatal(pair, value, read, want)
				}
				if pair[0] == "3512-set" {
					if result != 0 {
						t.Fatal("void setter")
					}
				} else if result != want {
					t.Fatal("setter return", pair, value, result, want)
				}
				if pair[0] != "limit-set" && *o.optionWords["settings-updated"] != 0 {
					t.Fatal("unrelated settings notification", pair)
				}
				if pair[0] == "limit-set" && repeat == 1 && *o.optionWords["settings-updated"] != 0 {
					t.Fatal("unchanged limit notified")
				}
				rows = append(rows, row{pair[0], value, result, read, *o.optionWords["settings-updated"]})
			}
		}
	}
	for _, value := range []uint32{0, 1, 0x80000000, 0xffffffff} {
		*memmap.PtrUint32(0x5D4594, 371700) = value
		got := legacy.PortTestServerConfigScalar("record-state", 0, 0)
		if got != uint64(int64(int32(value))) {
			t.Fatal("record-state width", value, got)
		}
		rows = append(rows, row{"record-state", int(value), got, got, 0})
	}
	for _, value := range []uint32{0, 0x800, 0x801, 0x2000, 0x7fff, 0x8000, 0xffffffff} {
		binary.LittleEndian.PutUint32(mode, 0xa5a5a5a5)
		result := legacy.PortTestServerConfigScalar("mode-store", int(value), 0)
		want := uint32(0xa5a5a5a5)
		if value > 0x800 && value < 0x8000 {
			want = value
		}
		if result != uint64(value) || binary.LittleEndian.Uint32(mode) != want {
			t.Fatal("mode boundaries", value, result, mode)
		}
		rows = append(rows, row{"mode-store", int(value), result, uint64(binary.LittleEndian.Uint32(mode)), 0})
	}
	for _, initial := range []uint32{0, 1, 0x80000000, 0xffffffff} {
		*o.optionWords["settings-updated"] = initial
		if got := legacy.PortTestServerConfigScalar("updated-get", 0, 0); got != uint64(int64(int32(initial))) {
			t.Fatal("notification getter", initial, got)
		}
		legacy.PortTestServerConfigScalar("updated-clear", 0, 0)
		if *o.optionWords["settings-updated"] != 0 {
			t.Fatal("clear")
		}
		legacy.PortTestServerConfigScalar("updated-set", 0, 0)
		if *o.optionWords["settings-updated"] != 1 {
			t.Fatal("set")
		}
	}
	spellbookCapture(t, "server-config-scalar-storage", rows, "dcb07f58c2da6d8153563850d8852f7a429e3ad8b5555958b86fa09a971b0b3b")
}

func TestServerConfigSlotStorage(t *testing.T) {
	newServerOptionsOwner(t)
	type row struct {
		Selected, Source, Dest int
		Result                 uint64
		Bytes                  []byte
	}
	var rows []row
	slots := unsafe.Slice(memmap.PtrUint8(0x5D4594, 371380), 116)
	seed := make([]byte, len(slots))
	for i := range seed {
		seed[i] = byte(i*37 + 11)
	}
	for selected := 0; selected < 2; selected++ {
		p := legacy.PortTestServerConfigPointer("slot-select", selected, nil)
		want := unsafe.Pointer(&slots[58*selected])
		if p != want || legacy.PortTestServerConfigPointer("slot-current", 0, nil) != want || legacy.PortTestServerConfigScalar("slot-index", 0, 0) != uint64(selected) {
			t.Fatal("selection identity")
		}
		for src := 0; src < 2; src++ {
			for dst := 0; dst < 2; dst++ {
				copy(slots, seed)
				expected := append([]byte(nil), seed...)
				copy(expected[58*dst:58*(dst+1)], seed[58*src:58*(src+1)])
				if legacy.PortTestServerConfigPointer("slot", src, nil) != unsafe.Pointer(&slots[58*src]) {
					t.Fatal("slot identity")
				}
				result := legacy.PortTestServerConfigScalar("slot-copy", src, dst)
				if result != uint64(dst) || string(slots) != string(expected) {
					t.Fatal("slot copy", src, dst, result)
				}
				rows = append(rows, row{selected, src, dst, result, append([]byte(nil), slots...)})
			}
		}
	}
	settings := legacy.PortTestServerConfigPointer("settings", 0, nil)
	admission := legacy.PortTestServerConfigPointer("admission", 0, nil)
	if settings != memmap.PtrOff(0x5D4594, 371516) || admission != unsafe.Add(settings, 100) {
		t.Fatal("settings pointers")
	}
	for n := 0; n < 256; n++ {
		*(*byte)(admission) = byte(n)
		result := legacy.PortTestServerConfigPointer("open-admission", 0, nil)
		if result != settings || *(*byte)(admission) != byte(n)&^0x10 {
			t.Fatal("admission flag", n)
		}
	}
	spellbookCapture(t, "server-config-slot-storage", rows, "01ab2050cd62168bf4f1693d61d2139728b601bc7c8a7e1e81fe9c293ab6981e")
}

func TestServerConfigRespawnModes(t *testing.T) {
	newServerOptionsOwner(t)
	serverConfigOwnBytes(t, 0x5D4594, 3584, 4)
	type row struct {
		Flags  uint32
		Value  int
		Result uint64
	}
	var rows []row
	for _, flags := range []uint32{0, 1, 32, 1024, 4096, 4097, 0xffffffff} {
		for _, value := range []int{-1, 0, 1, 2, 2147483647} {
			t.Run(fmt.Sprintf("%x/%d", flags, value), func(t *testing.T) {
				defer noxflags.PortTestGameFlags(noxflags.GameFlag(flags))()
				legacy.PortTestServerConfigScalar("respawn-set", value, 0)
				result := legacy.PortTestServerConfigScalar("respawn-get", 0, 0)
				want := uint64(int64(value))
				if flags&4096 != 0 {
					want = 0
				}
				if result != want {
					t.Fatal(result, want)
				}
				rows = append(rows, row{flags, value, result})
			})
		}
	}
	spellbookCapture(t, "server-config-respawn-modes", rows, "8724a618281420bee0aea88aea15a6b889b22bfdc26ae7931e8d49f4adfa1ca7")
}

func TestServerConfigDeadlineWidths(t *testing.T) {
	newServerOptionsOwner(t)
	deadline := unsafe.Slice(memmap.PtrUint8(0x5D4594, 3468), 8)
	oldClock := legacy.PlatformTicks
	t.Cleanup(func() { legacy.PlatformTicks = oldClock })
	type row struct {
		Tick                      uint64
		Value                     int
		Stored, Result, Remaining uint64
		Calls                     int
	}
	var rows []row
	for _, tick := range []uint64{0, 1, 0x7fffffff, 0xffffffff, 0x100000000, ^uint64(0) - 100, ^uint64(0)} {
		for _, value := range []int{-2147483648, -1, 0, 1, 60000, 2147483647} {
			calls := 0
			legacy.PlatformTicks = func() uint64 { calls++; return tick }
			result := legacy.PortTestServerConfigScalar("timer-reset", value, 0)
			stored := binary.LittleEndian.Uint64(deadline)
			if stored != uint64(uint32(tick))+uint64(int64(value)) || result != uint64(int64(value)) || calls != 1 {
				t.Fatal("deadline reset", tick, value, result, stored, calls)
			}
			remaining := legacy.PortTestServerConfigScalar("timer-left", 0, 0)
			// The C clock adapter truncates ticks to uint32.
			// Remaining time observes the low 32 bits before converting back to signed int.
			if remaining != uint64(int64(int32(value))) || calls != 2 {
				t.Fatal("remaining width", tick, value, remaining, calls)
			}
			rows = append(rows, row{tick, value, stored, result, remaining, calls})
		}
	}
	spellbookCapture(t, "server-config-deadline-widths", rows, "bae35f4ba40326c7b4d7a55f359fa86e49b3abc52e6c98fabdc01105f9a2a574")
}
