//go:build porttest

package opennox

import (
	"encoding/binary"
	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestServerBrowserClockWidths(t *testing.T) {
	o := newEntryOwner(t)
	words, restore := legacy.PortTestServerBrowserWords()
	defer restore()
	clocks, restoreClocks := legacy.PortTestServerBrowserClocks()
	defer restoreClocks()
	deadline := serverConfigOwnBytes(t, 0x5D4594, 814972, 8)
	ticks, join, message := legacy.PlatformTicks, legacy.Nox_client_joinGame_438A90, legacy.Sub_449E30
	defer func() {
		legacy.PlatformTicks = ticks
		legacy.Nox_client_joinGame_438A90 = join
		legacy.Sub_449E30 = message
	}()
	set, restoreStrings := o.c.srv.PortTestMeterStrings(strman.Entry{ID: "noxworld.c:TestCon", Vals: []strman.Variant{{Str: "TestCon"}}})
	defer restoreStrings()
	set(0)
	for _, now := range []uint64{0, 1, 0xffffffff - 20001, 0xffffffff - 999, 0xffffffff, 1 << 32, ^uint64(0)} {
		legacy.PlatformTicks = func() uint64 { return now }
		for _, state := range []uint32{3, 4, 8, 9} {
			for _, due := range []uint64{0, uint64(uint32(now)), uint64(uint32(now)) + 1, ^uint64(0)} {
				for _, p := range words {
					*p = 0
				}
				*words["dword_5d4594_814548"] = state
				*clocks["qword_5d4594_814956"] = due
				binary.LittleEndian.PutUint64(deadline, due)
				joins := 0
				messages := []string{}
				legacy.Nox_client_joinGame_438A90 = func() int { joins++; return 1 }
				legacy.Sub_449E30 = func(s string) int { messages = append(messages, s); return 0 }
				ret := legacy.PortTestServerBrowserTick()
				wantState, wantError, wantJoins := state, uint32(0), 0
				switch state {
				case 3:
					if uint64(uint32(now)) >= due {
						wantState, wantError = 2, 8
					}
				case 4:
					wantState = 3
					if *clocks["qword_5d4594_814956"] != uint64(uint32(now)+uint32(20000)) || len(messages) != 1 || messages[0] != "TestCon" {
						t.Fatal("32-bit deadline addition", now, *clocks["qword_5d4594_814956"], messages)
					}
				case 8:
					wantState = 9
					if binary.LittleEndian.Uint64(deadline) != uint64(uint32(now))+1000 {
						t.Fatal("64-bit deadline addition", now, binary.LittleEndian.Uint64(deadline))
					}
				case 9:
					if uint64(uint32(now)) > due {
						wantJoins = 1
					}
				}
				if ret != 1 || *words["dword_5d4594_814548"] != wantState || *words["nox_client_connError_814552"] != wantError || joins != wantJoins {
					t.Fatal("clock transition", now, state, due, ret, *words["dword_5d4594_814548"], *words["nox_client_connError_814552"], joins)
				}
			}
		}
	}
}
