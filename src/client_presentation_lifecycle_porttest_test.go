//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
)

func TestClientLifecycleState(t *testing.T) {
	state := serverConfigOwnBytes(t, 0x5D4594, 754052, 4)
	host := serverConfigOwnBytes(t, 0x5D4594, 807172, 97)
	for i := range host {
		host[i] = byte(i*17 + 9)
	}
	before := append([]byte(nil), host...)
	for _, v := range []int{0, 1, -1, 2147483647, -2147483648, 0x12345678} {
		legacy.Sub_42EB90(v)
		if got := int32(binary.LittleEndian.Uint32(state)); got != int32(v) {
			t.Fatal("state", got, v)
		}
		if p := legacy.Nox_xxx_getHostInfoPtr_431770(); unsafe.Pointer(p) != unsafe.Pointer(memmap.PtrOff(0x5D4594, 807172)) {
			t.Fatal("host identity")
		}
	}
	if !bytes.Equal(host, before) {
		t.Fatal("host getter mutated data")
	}
}

func TestClientServerAddressCopy(t *testing.T) {
	dest := serverConfigOwnBytes(t, 0x5D4594, 806060, 32)
	var inputs []string
	for n := 0; n <= 32; n++ {
		inputs = append(inputs, strings.Repeat("a", n))
	}
	inputs = append(inputs, strings.Repeat("z", 70), strings.Repeat("é", 15))
	for _, at := range []int{0, 1, 22, 23, 24} {
		inputs = append(inputs, strings.Repeat("b", at)+"\x00tail")
	}
	type row struct {
		Input  string
		Output []byte
	}
	var rows []row
	for _, input := range inputs {
		for i := range dest {
			dest[i] = 0xa5
		}
		want := bytes.Repeat([]byte{0xa5}, 32)
		clear(want[:23])
		prefix := strings.SplitN(input, "\x00", 2)[0]
		copy(want[:23], prefix)
		legacy.Nox_xxx_copyServerIPAndPort_431790(input)
		if !bytes.Equal(dest, want) {
			t.Fatalf("address %q: %x != %x", input, dest, want)
		}
		rows = append(rows, row{input, append([]byte(nil), dest...)})
	}
	interactionCapture(t, "client-server-address-copy", rows)
}

func TestClientAudioContextLifecycle(t *testing.T) {
	type row struct {
		Voices, Flags    int
		Stage            string
		Events           []audioStreamCallback
		References, Slot uint32
	}
	var rows []row
	for _, n := range []int{0, 1, 3} {
		for _, flags := range []uint32{0, 1, 4, 5} {
			t.Run(fmt.Sprintf("voices%d_flags%d", n, flags), func(t *testing.T) {
				o := newAudioStreamOwner(t, 1, 3)
				old := legacy.Get_dword_5d4594_805984()
				t.Cleanup(func() { legacy.Set_dword_5d4594_805984(old) })
				legacy.Set_dword_5d4594_805984(nil)
				legacy.Sub_431290()
				legacy.Sub_431270()
				if len(o.events) != 0 {
					t.Fatal("nil context called audio")
				}
				device := o.device()
				ctx := o.call("sub_487150", 0, 0)
				if device == 0 || ctx == 0 {
					t.Fatal("context allocation")
				}
				for i := 0; i < n; i++ {
					voice := o.call("sub_487750", ctx)
					if voice == 0 {
						t.Fatal("voice allocation")
					}
					w := audioStreamWords(voice, 78)
					w[3] = uint32(i - 1)
					w[31] = flags
					w[37] = audioStreamPointer(o.callbacks[13])
				}
				legacy.Set_dword_5d4594_805984(unsafe.Pointer(uintptr(ctx)))
				record := func(stage string) {
					dw := audioStreamWords(device, 22)
					slot := uint32(0)
					if dw[6] != 0 {
						slot = 1
					}
					rows = append(rows, row{n, int(flags), stage, append([]audioStreamCallback(nil), o.events...), dw[4], slot})
				}
				o.events = nil
				legacy.Sub_431290()
				wantStop := 0
				if flags&5 != 0 {
					wantStop = n
				}
				if o.callbackCounts[7] != wantStop || o.callbackCounts[13] != n || o.callbackCounts[5] != 0 || o.callbackCounts[3] != 0 {
					t.Fatal("stop callbacks", o.callbackCounts)
				}
				if legacy.Get_dword_5d4594_805984() != unsafe.Pointer(uintptr(ctx)) {
					t.Fatal("stop destroyed context")
				}
				record("stop")
				o.events = nil
				legacy.Sub_431270()
				if legacy.Get_dword_5d4594_805984() != nil || o.callbackCounts[5] != n || o.callbackCounts[3] != 1 || o.callbackCounts[13] != 2*n {
					t.Fatal("destroy callbacks", o.callbackCounts)
				}
				dw := audioStreamWords(device, 22)
				if dw[4] != 0 || dw[6] != 0 {
					t.Fatal("context remained linked to device", dw[4], dw[6])
				}
				record("destroy")
				o.events = nil
				legacy.Sub_431290()
				legacy.Sub_431270()
				if len(o.events) != 0 {
					t.Fatal("repeated cleanup called audio")
				}
				record("repeat")
			})
		}
	}
	interactionCapture(t, "client-audio-context-lifecycle", rows)
}
