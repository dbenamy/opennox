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

func TestClientAudioEventsDefaults(t *testing.T) {
	data, free := alloc.Make([]uint32{}, 50)
	defer free()
	ptr := audioEventPointer(unsafe.Pointer(&data[0]))
	old := timer.PlatformTicks
	calls := 0
	timer.PlatformTicks = func() uint64 { calls++; return 0x123456789 }
	defer func() { timer.PlatformTicks = old }()
	var rows []map[string]any
	for _, fill := range []uint32{0, 0xaabb0000, 0x40000000} {
		for i := range data {
			data[i] = fill + uint32(i)
		}
		data[5] = fill
		want := append([]uint32(nil), data...)
		for _, i := range []int{0, 1, 2, 14, 15, 19, 20, 48, 18, 17, 25, 26} {
			want[i] = 0
		}
		want[12] = 1
		want[16] = 600
		flags, ret := uint32(3), uint64(1)
		if fill == 0x40000000 {
			flags, ret = 1, 0
		}
		copy(want[4:12], []uint32{flags, 0x40000000, 0x40000000, 0x40000000 / 1000, 1000, 0, 0x23456789, 1})
		calls = 0
		got := legacy.PortTestAudioEventCall("sub_451920", ptr)
		if got != ret || calls != 1 || !reflect.DeepEqual(data, want) {
			t.Fatal("event defaults/timer/guards", fill, got, ret, calls)
		}
		rows = append(rows, map[string]any{"fill": fill, "return": got, "words": append([]uint32(nil), data...), "clock_calls": calls})
	}
	spellbookCapture(t, "client-audio-events-defaults", rows, "56af29b0eb48e725aabe552edee2e0a323adcc263fbc4ce2e3cc9536dd75c24a")
}
