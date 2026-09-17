//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
	"unsafe"
)

func TestServerConfigRateNotifications(t *testing.T) {
	type row struct {
		Flags, Before, After uint32
		Result               uint64
		Packets              []legacy.PortTestShopPacketResult
	}
	var rows []row
	for _, flags := range []uint32{0, 1, 0x20000, 0x20001} {
		for _, before := range []uint32{0, 1, 255, 0xffffffff} {
			for _, value := range []uint32{0, 1, 256, 0xffffffff} {
				t.Run(fmt.Sprintf("%x/%x/%x", flags, before, value), func(t *testing.T) {
					o := newServerOptionsOwner(t)
					rate := serverConfigOwnBytes(t, 0x587000, 4728, 4)
					binary.LittleEndian.PutUint32(rate, before)
					defer noxflags.PortTestGameFlags(noxflags.GameFlag(flags))()
					result := legacy.PortTestServerConfigScalar("rate-set", int(value), 0)
					packets := o.packets()
					expected := 0
					if before != value && flags&0x20000 != 0 {
						expected = 1
					}
					if len(packets) != expected || binary.LittleEndian.Uint32(rate) != value || legacy.PortTestServerConfigScalar("rate-get", 0, 0) != uint64(int64(int32(value))) {
						t.Fatal("rate update", result, packets, rate)
					}
					if expected != 0 {
						p := packets[0]
						if !bytes.Equal(p.Data, []byte{236, byte(before)}) || p.Recipient != 159 || p.Ordered != 1 {
							t.Fatal("notify must use old rate and existing route", p)
						}
					}
					rows = append(rows, row{flags, before, value, result, packets})
				})
			}
		}
	}
	spellbookCapture(t, "server-config-rate-notifications", rows, "8b97496a2a06576ed8e3e9a3e3820544ffed3ffbb210ce47eabb7ef17a9e27df")
}
func TestServerConfigConnectionRates(t *testing.T) {
	newServerOptionsOwner(t)
	table := serverConfigOwnBytes(t, 0x587000, 4664, 40)
	// Distinct nonzero values ensure a missing/default lookup cannot mask a bad index.
	for i := 0; i < 5; i++ {
		binary.LittleEndian.PutUint32(table[8*i:], uint32(11+3*i))
		binary.LittleEndian.PutUint32(table[8*i+4:], uint32(101+7*i))
	}
	type row struct {
		Value  int
		Result uint64
	}
	var rows []row
	for _, value := range []int{-1, 0, 1, 11, 14, 17, 20, 23, 26, 2147483647} {
		got := legacy.PortTestServerConfigScalar("connection-rate", value, 0)
		want := uint64(1)
		for i := 0; i < 5; i++ {
			if value == 11+3*i {
				want = uint64(101 + 7*i)
			}
		}
		if got != want {
			t.Fatal("connection lookup", value, got, want)
		}
		rows = append(rows, row{value, got})
	}
	spellbookCapture(t, "server-config-connection-rates", rows, "5c36aa36295996879bd2a23c71013c7eeb70cc07a884a8a755258fd0a085f606")
}
func TestServerConfigModeLimitsAndTimer(t *testing.T) {
	newServerOptionsOwner(t)
	deadline := unsafe.Slice(memmap.PtrUint8(0x5D4594, 3468), 8)
	old := legacy.PlatformTicks
	t.Cleanup(func() { legacy.PlatformTicks = old })
	type row struct {
		Flags                                    uint32
		Index                                    int
		Tick, Score, Minutes, Duration, Deadline uint64
		Calls                                    int
	}
	var rows []row
	for _, mode := range []struct {
		flags uint32
		index int
	}{{0x100, 0}, {0x101, 0}, {0x400, 1}, {0x20, 2}, {0x10, 3}, {0x40, 4}, {0x80, 0}} {
		for _, tick := range []uint64{0, 0xffffffff, 0x100000000, ^uint64(0)} {
			restore := noxflags.PortTestGameFlags(noxflags.GameFlag(mode.flags))
			calls := 0
			legacy.PlatformTicks = func() uint64 { calls++; return tick }
			score := legacy.PortTestServerConfigScalar("score", int(mode.flags), 0)
			minutes := legacy.PortTestServerConfigScalar("minutes", int(mode.flags), 0)
			if score != uint64(101+11*mode.index) || minutes != uint64(13+7*mode.index) {
				t.Fatal("mode table", mode, score, minutes)
			}
			duration := legacy.PortTestServerConfigScalar("timer-init", 0, 0)
			stored := binary.LittleEndian.Uint64(deadline)
			if duration != minutes*60000 || stored != uint64(uint32(tick))+duration || calls != 1 {
				t.Fatal("timer init", tick, duration, stored, calls)
			}
			rows = append(rows, row{mode.flags, mode.index, tick, score, minutes, duration, stored, calls})
			restore()
		}
	}
	spellbookCapture(t, "server-config-mode-limits-timer", rows, "0be5f95e7e0c6e9031379b0fa4b2cdeabd2dd96bca399dcb10e8d1fda7e2c4b6")
}
func TestServerConfigSelectedModePredicate(t *testing.T) {
	newServerOptionsOwner(t)
	slots := unsafe.Slice(memmap.PtrUint8(0x5D4594, 371380), 116)
	type row struct {
		Flags  uint32
		Slot   int
		Value  byte
		Result uint64
	}
	var rows []row
	for _, flags := range []uint32{0, 1, 128, 129, 0xffffffff} {
		for slot := 0; slot < 2; slot++ {
			for _, value := range []byte{0, 1, 0x7f, 0x80, 0xff} {
				restore := noxflags.PortTestGameFlags(noxflags.GameFlag(flags))
				slots[53], slots[111] = value^0x80, value^0x80
				slots[58*slot+53] = value
				legacy.PortTestServerConfigPointer("slot-select", slot, nil)
				result := legacy.PortTestServerConfigScalar("special-mode", 0, 0)
				if result != uint64(bool2int(flags&128 != 0 && value&0x80 != 0)) {
					t.Fatal("selected mode", flags, slot, value, result)
				}
				rows = append(rows, row{flags, slot, value, result})
				restore()
			}
		}
	}
	spellbookCapture(t, "server-config-selected-mode", rows, "ca3f14e2f4b60a2a9fdbd95aae6aad65dbd16d626073e58fa33a87c5e0ffa6cb")
}

func TestServerConfigTimerNotifications(t *testing.T) {
	type row struct {
		Flags   uint32
		Value   int
		Result  uint64
		Calls   int
		Packets []legacy.PortTestShopPacketResult
	}
	var rows []row
	for _, flags := range []uint32{0, 1, 0x20000, 0x20001} {
		for _, value := range []int{-1, 0, 1, 2147483647} {
			t.Run(fmt.Sprintf("%x/%d", flags, value), func(t *testing.T) {
				o := newServerOptionsOwner(t)
				enabled := serverConfigOwnBytes(t, 0x587000, 4660, 4)
				binary.LittleEndian.PutUint64(unsafe.Slice(memmap.PtrUint8(0x5D4594, 3468), 8), 0x100000014)
				defer noxflags.PortTestGameFlags(noxflags.GameFlag(flags))()
				old := legacy.PlatformTicks
				calls := 0
				legacy.PlatformTicks = func() uint64 { calls++; return 0xfffffff5 }
				t.Cleanup(func() { legacy.PlatformTicks = old })
				result := legacy.PortTestServerConfigScalar("timer-set", value, 0)
				packets := o.packets()
				expected := bool2int(flags&1 != 0)
				if len(packets) != expected || calls != expected || binary.LittleEndian.Uint32(enabled) != uint32(value) || legacy.PortTestServerConfigScalar("timer-get", 0, 0) != uint64(int64(value)) {
					t.Fatal("timer notification", result, calls, packets)
				}
				if expected != 0 {
					p := packets[0]
					want := make([]byte, 13)
					want[0] = 211
					binary.LittleEndian.PutUint32(want[1:], uint32(value))
					binary.LittleEndian.PutUint32(want[5:], 31)
					binary.LittleEndian.PutUint32(want[9:], o.c.srv.Frame())
					if !bytes.Equal(p.Data, want) || p.Recipient != 159 || p.Ordered != 1 {
						t.Fatal("timer packet", p, want)
					}
				}
				rows = append(rows, row{flags, value, result, calls, packets})
			})
		}
	}
	spellbookCapture(t, "server-config-timer-notifications", rows, "29307c94fcef3981716532b7235f4a4c3a5b381ff22b965e8f07e7c42f25ae0b")
}
