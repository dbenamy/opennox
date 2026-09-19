//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"slices"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func onlineCall(op string, a, b uint32) uint32 { return legacy.PortTestOnlineSessionCall(op, a, b) }
func onlineSnapshot(words map[string]*uint32) [6]uint32 {
	return [6]uint32{*words["failed"], *words["status"], *words["pending"], *words["active"], *words["deadline"], *words["origin"]}
}
func TestOnlineSessionRetryTimers(t *testing.T) {
	o := newListboxOwner(t)
	words, restore := legacy.PortTestOnlineSessionWords()
	defer restore()
	frame, rate := o.c.srv.Frame(), o.c.srv.TickRate()
	defer func() { o.c.srv.SetFrame(frame); o.c.srv.SetTickRate(rate) }()
	oldLog := legacy.NetworkLogPrint
	defer func() { legacy.NetworkLogPrint = oldLog }()
	var logs []string
	legacy.NetworkLogPrint = func(s string) { logs = append(logs, s) }
	marker := serverConfigOwnBytes(t, 0x5D4594, 528268, 4)
	scratch := serverConfigOwnBytes(t, 0x85B3FC, 10304, 116)
	type row struct {
		Op             string
		Frame, Rate    uint32
		Before, After  [6]uint32
		Return, Marker uint32
		Cleared        bool
		Logs           []string
	}
	var rows []row
	for _, fps := range []uint32{1, 30, 60, 0x1000000, 0xffffffff} {
		o.c.srv.SetTickRate(fps)
		for _, now := range []uint32{0, 1, 119, 0x7fffffff, 0xfffffffe, 0xffffffff} {
			o.c.srv.SetFrame(now)
			for _, pending := range []uint32{0, 1, 2, 0xffffffff} {
				for _, active := range []uint32{0, 1, 2, 0xffffffff} {
					for _, origin := range []uint32{0, 1, now, now - 3600*fps, now - 3600*fps - 1} {
						for _, deadline := range []uint32{0, 1, now + 120*fps} {
							before := [6]uint32{9, 7, pending, active, deadline, origin}
							for _, op := range []string{"start", "attempt", "retry", "reset"} {
								for i, k := range []string{"failed", "status", "pending", "active", "deadline", "origin"} {
									*words[k] = before[i]
								}
								binary.LittleEndian.PutUint32(marker, 0x12345678)
								for i := range scratch {
									scratch[i] = 0xa5
								}
								want := before
								ret := uint32(0)
								mark := uint32(0x12345678)
								cleared := false
								switch op {
								case "start":
									if pending != 1 && active != 1 && deadline == 0 && origin == 0 {
										want[2] = 1
										want[3] = 0
										want[4] = now + 120*fps
										want[5] = now
									}
								case "reset":
									want[2] = 0
									want[3] = 0
									want[4] = 0
									want[5] = 0
								case "retry":
									want[3] = 0
									want[4] = now + 120*fps
									ret = want[4]
								case "attempt":
									if now-origin > 3600*fps {
										want[2] = 0
										want[3] = 0
										want[4] = 0
										want[5] = 0
										mark = 1
										ret = 1
									} else if pending == 0 {
										ret = 0
									} else if active != 0 {
										ret = active
									} else {
										want[0] = 0
										want[3] = 0
										want[4] = now + 120*fps
										ret = want[4]
										cleared = true
									}
								}
								logs = nil
								got := onlineCall(op, 0, 0)
								var wantLogs []string
								retryLog := fmt.Sprintf("RECON: TryReconnectAgain called on frame (%d)", int32(now))
								if op == "start" && pending != 1 && active != 1 && deadline == 0 && origin == 0 {
									wantLogs = []string{fmt.Sprintf("RECON: Starting reconnection process frame (%d)", int32(now))}
								}
								if op == "retry" {
									wantLogs = []string{retryLog}
								}
								if cleared {
									wantLogs = []string{"RECON: Attempting to re-login", retryLog}
								}
								if !slices.Equal(logs, wantLogs) {
									t.Fatal("log order/format", op, now, logs, wantLogs)
								}
								after := onlineSnapshot(words)
								if got != ret || after != want || binary.LittleEndian.Uint32(marker) != mark {
									t.Fatalf("%s frame=%x fps=%x before=%v: got %x %v want %x %v", op, now, fps, before, got, after, ret, want)
								}
								expected := bytes.Repeat([]byte{0xa5}, 116)
								if cleared {
									clear(expected[4:112])
								}
								if !bytes.Equal(scratch, expected) {
									t.Fatal("login scratch extent", op, before)
								}
								rows = append(rows, row{op, now, fps, before, after, got, mark, cleared, append([]string(nil), logs...)})
							}
						}
					}
				}
			}
		}
	}
	spellbookCapture(t, "online-session-timers", rows, "")
}

func TestOnlineSessionStatusAndEmptyServices(t *testing.T) {
	newListboxOwner(t)
	words, restore := legacy.PortTestOnlineSessionWords()
	defer restore()
	oldCallback := legacy.Sub_41E300
	defer func() { legacy.Sub_41E300 = oldCallback }()
	table := serverConfigOwnBytes(t, 0x587000, 58128, 12*16)
	for i := 0; i < 12; i++ {
		if binary.LittleEndian.Uint32(table[16*i+4:]) != 0 || binary.LittleEndian.Uint32(table[16*i+8:]) != 0 {
			t.Fatal("unexpected live legacy queue", i)
		}
	}
	original := append([]byte(nil), table...)
	type row struct {
		Op                    string
		Status, Reply, Return uint32
		Calls                 []int
	}
	var rows []row
	for status := uint32(0); status < 12; status++ {
		for _, reply := range []uint32{0, 1} {
			for _, op := range []string{"append", "clear", "channels", "users", "status", "player", "players"} {
				*words["status"] = status
				var calls []int
				legacy.Sub_41E300 = func(v int) int { calls = append(calls, v); return int(reply) }
				got := onlineCall(op, status, 0xffff)
				want := uint32(0)
				called := false
				switch op {
				case "append":
					want = reply
					called = true
				case "status":
					want = status
				case "channels", "users":
					want = status
					if status == 7 {
						want = reply
						called = true
					}
				}
				if got != want || len(calls) != bool2int(called) || called && calls[0] != 11 {
					t.Fatal(op, status, reply, got, calls, want)
				}
				if !bytes.Equal(table, original) {
					t.Fatal("empty queue table mutated")
				}
				rows = append(rows, row{op, status, reply, got, append([]int(nil), calls...)})
			}
		}
	}
	briefing := serverConfigOwnBytes(t, 0x5D4594, 527720, 4)
	marker := serverConfigOwnBytes(t, 0x5D4594, 528268, 4)
	transition := serverConfigOwnBytes(t, 0x5D4594, 371700, 4)
	for _, v := range []uint32{0, 1, 2, 0x7fffffff, 0x80000000, 0xffffffff} {
		if onlineCall("briefing", v, 0) != v || binary.LittleEndian.Uint32(briefing) != v {
			t.Fatal("briefing setter", v)
		}
		if onlineCall("marker", v, 0) != v || binary.LittleEndian.Uint32(marker) != v {
			t.Fatal("marker setter", v)
		}
		binary.LittleEndian.PutUint32(transition, v)
		if onlineCall("map", 0, 0) != 0 || binary.LittleEndian.Uint32(transition) != 1 {
			t.Fatal("map transition", v)
		}
	}
	if legacy.Sub_41FA40() != "" {
		t.Fatal("unexpected legacy account name")
	}
	spellbookCapture(t, "online-session-status", rows, "")
}
