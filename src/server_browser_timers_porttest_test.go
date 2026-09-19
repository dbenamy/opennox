//go:build porttest

package opennox

import (
	"encoding/binary"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
	"time"
)

func TestServerBrowserTimers(t *testing.T) {
	o := newEntryOwner(t)
	words, restore := legacy.PortTestServerBrowserWords()
	defer restore()
	clocks, restoreClocks := legacy.PortTestServerBrowserClocks()
	defer restoreClocks()
	deadline := serverConfigOwnBytes(t, 0x5D4594, 814972, 8)
	join, refresh, wait := legacy.Nox_client_joinGame_438A90, legacy.Nox_client_refreshServerList_4378B0, legacy.Sub_438770_waitList
	defer func() {
		legacy.Nox_client_joinGame_438A90 = join
		legacy.Nox_client_refreshServerList_4378B0 = refresh
		legacy.Sub_438770_waitList = wait
	}()
	var calls [3]int
	legacy.Nox_client_joinGame_438A90 = func() int { calls[0]++; return 1 }
	legacy.Nox_client_refreshServerList_4378B0 = func() { calls[1]++ }
	legacy.Sub_438770_waitList = func() { calls[2]++ }
	type row struct {
		State           uint32
		Future          bool
		After, Error    uint32
		Calls           [3]int
		Return          int
		DeadlineInRange bool
	}
	var rows []row
	for _, state := range []uint32{1, 3, 6, 8, 9, 11, 0xffffffff} {
		for _, future := range []bool{false, true} {
			for _, p := range words {
				*p = 0
			}
			calls = [3]int{}
			*words["dword_5d4594_814548"] = state
			due := uint64(0)
			if future {
				due = ^uint64(0)
			}
			*clocks["qword_5d4594_814956"] = due
			binary.LittleEndian.PutUint64(deadline, due)
			for platformTicks() == 0 {
				time.Sleep(time.Millisecond)
			}
			before := platformTicks()
			ret := legacy.PortTestServerBrowserTick()
			after := platformTicks()
			r := row{State: state, Future: future, After: *words["dword_5d4594_814548"], Error: *words["nox_client_connError_814552"], Calls: calls, Return: ret, DeadlineInRange: true}
			wantState, wantError := state, uint32(0)
			wantCalls := [3]int{}
			switch state {
			case 3:
				if !future {
					wantState, wantError = 2, 8
				}
			case 8:
				wantState = 9
				due = binary.LittleEndian.Uint64(deadline)
				r.DeadlineInRange = due >= before+1000 && due <= after+1000
			case 9:
				if !future {
					wantCalls[0] = 1
				}
			}
			if r.After != wantState || r.Error != wantError || r.Calls != wantCalls || r.Return != 1 || !r.DeadlineInRange {
				t.Fatalf("timer: %+v", r)
			}
			rows = append(rows, r)
		}
	}
	type refreshRow struct {
		Gate   string
		Future bool
		Calls  [3]int
	}
	var refreshRows []refreshRow
	for _, gate := range []string{"", "nox_game_createOrJoin_815048", "dword_5d4594_815044", "dword_5d4594_815052"} {
		for _, future := range []bool{false, true} {
			for _, p := range words {
				*p = 0
			}
			calls = [3]int{}
			for _, name := range []string{"dword_5d4594_814984", "dword_5d4594_814988", "nox_wol_wnd_gameList_815012", "dword_5d4594_815000"} {
				w := o.c.GUI.NewWindowRaw(o.parent, 8, 0, 0, 10, 10, nil)
				w.SetHidden(name == "dword_5d4594_815000")
				*words[name] = uint32(uintptr(w.C()))
			}
			*clocks["qword_5d4594_815068"] = 0
			if future {
				*clocks["qword_5d4594_815068"] = ^uint64(0)
			}
			if gate != "" {
				*words[gate] = 1
			}
			ret := legacy.PortTestServerBrowserTick()
			want := [3]int{}
			if gate == "" {
				if future {
					want[2] = 1
				} else {
					want[1] = 1
				}
			}
			if calls != want || ret != 1 {
				t.Fatal("refresh gate", gate, future, calls, ret)
			}
			refreshRows = append(refreshRows, refreshRow{gate, future, calls})
		}
	}
	spellbookCapture(t, "server-browser-timers", rows, "5dc27892376e936abcc8b6f23149c34de9b143fb20614295294400e0fe25dd07")
	spellbookCapture(t, "server-browser-refresh-gates", refreshRows, "b0a9da6be4b3bf0edced07a776421d3e62a83ac954cb6f945bf0d138e192d6a8")
}

func TestServerBrowserConstructionNotifications(t *testing.T) {
	o := newEntryOwner(t)
	_, restore := legacy.PortTestServerBrowserWords()
	defer restore()
	for which := 0; which < 2; which++ {
		for _, id := range []uint32{1, 10001, 10002, 10020, 10047, 10057, 10070, 0xffffffff} {
			if got := legacy.PortTestServerBrowserNotification(which, o.parent.C(), id); got != 0 {
				t.Fatal("construction notification", which, id, got)
			}
		}
	}
}
