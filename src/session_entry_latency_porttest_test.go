//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/common/ntype"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestSessionEntryLatency(t *testing.T) {
	o := newMatchRosterOwner(t)
	type row struct {
		Name     string
		Cursor   int
		Trace    []int
		Messages [][]byte
	}
	var rows []row
	cases := []struct {
		name     string
		slots    []int
		values   map[int]int
		selected [3]int
		queries  [3]int
	}{
		{"empty", nil, nil, [3]int{-1, -1, -1}, [3]int{}},
		{"zero-peer", []int{1}, map[int]int{1: 0}, [3]int{1, 1, 1}, [3]int{33, 33, 33}},
		{"zero-peers", []int{1, 7}, map[int]int{}, [3]int{1, 7, 1}, [3]int{33, 33, 33}},
		{"skip-zero", []int{1, 7, 31}, map[int]int{7: 15}, [3]int{7, 31, 7}, [3]int{3, 1, 3}},
		{"host", []int{31}, map[int]int{31: -1}, [3]int{31, 31, 31}, [3]int{1, 1, 1}},
		{"narrow", []int{1, 7}, map[int]int{1: 65536, 7: -1}, [3]int{1, 7, 1}, [3]int{2, 2, 2}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			restore := legacy.PortTestSessionEntryLatencyOwner()
			defer restore()
			for i := 0; i < 32; i++ {
				p := o.s.Players.ByIndRaw(ntype.PlayerInd(i))
				p.Active = 0
				p.PlayerInd = byte(i)
				p.NetCodeVal = 0x12340000 + uint32(i*9)
			}
			for _, i := range c.slots {
				o.s.Players.ByIndRaw(ntype.PlayerInd(i)).Active = 1
			}
			old := legacy.Sub_554240
			defer func() { legacy.Sub_554240 = old }()
			var trace []int
			legacy.Sub_554240 = func(index ntype.PlayerInd) int { trace = append(trace, int(index)); return c.values[int(index)] }
			for call := 0; call < 3; call++ {
				o.s.NetList.ResetAll()
				trace = nil
				legacy.PortTestSessionEntryScalar("latency", 0)
				cursor := legacy.PortTestSessionEntryLatencyCursor()
				if cursor != c.selected[call] || len(trace) != c.queries[call] {
					t.Fatal("latency rotation", call, cursor, trace)
				}
				messages := visibilityEffectsPackets(o.s)
				payload := []byte(nil)
				if cursor >= 0 {
					payload = make([]byte, 5)
					payload[0] = 215
					binary.LittleEndian.PutUint16(payload[1:], uint16(o.s.Players.ByIndRaw(ntype.PlayerInd(cursor)).NetCodeVal))
					binary.LittleEndian.PutUint16(payload[3:], uint16(c.values[cursor]))
				}
				for i, data := range messages {
					active := false
					for _, slot := range c.slots {
						if slot == i {
							active = true
						}
					}
					want := []byte(nil)
					if active {
						want = payload
					}
					if !bytes.Equal(data, want) {
						t.Fatal("latency fanout", i, data, want)
					}
				}
				rows = append(rows, row{fmt.Sprintf("%s/%d", c.name, call), cursor, append([]int(nil), trace...), messages})
			}
		})
	}
	spellbookCapture(t, "session-entry-latency", rows, "f44c3f8658640f0f74f8e173ef53b445ca3ef9cacb0f2c9133e1bceed4dad989")
}
