//go:build porttest

package opennox

import (
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"testing"
)

func TestPlayerStateSuddenDeath(t *testing.T) {
	o := newMatchRosterOwner(t)
	for _, off := range []uintptr{1392, 3476, 3508, 3520, 3536} {
		serverConfigOwnBytes(t, 0x5D4594, off, 4)
	}
	type row struct {
		Frame, ThresholdBits, LifetimeBits uint32
		Threshold, Lifetime                int32
		Status                             [3]uint32
		Flags                              uint32
	}
	var rows []row
	for _, frame := range []uint32{0, 123, 0xfffffff0} {
		for _, threshold := range []float32{-2.75, 0, 3.75, 19.5} {
			for _, lifetime := range []float32{-3.125, 0, 60.75} {
				name := fmt.Sprintf("frame%x/threshold%g/lifetime%g", frame, threshold, lifetime)
				t.Run(name, func(t *testing.T) {
					defer noxflags.PortTestGameFlags(noxflags.GameFlag(0x4000011))()
					o.s.SetFrame(frame)
					o.balance(map[string]float64{"SuddenDeathPlayerThreshold": float64(threshold), "SuddenDeathLifeTime": float64(lifetime)})
					*memmap.PtrUint32(0x5D4594, 3536) = 0xaabbccdd
					before := [3]uint32{0x100, 0x121, 0xfffffeff}
					for i := range o.units {
						o.units[i].UpdateDataPlayer().Player.Field3680 = before[i]
					}
					legacy.PortTestPlayerStateReset()
					r := row{Frame: memmap.Uint32(0x5D4594, 3520), ThresholdBits: math.Float32bits(threshold), LifetimeBits: math.Float32bits(lifetime), Threshold: memmap.Int32(0x5D4594, 3476), Lifetime: memmap.Int32(0x5D4594, 1392), Flags: uint32(noxflags.GetGame())}
					for i := range o.units {
						r.Status[i] = o.units[i].UpdateDataPlayer().Player.Field3680
						if r.Status[i] != before[i]&^256 {
							t.Fatal("status reset", r.Status)
						}
					}
					if r.Frame != frame || r.Threshold != int32(threshold) || r.Lifetime != int32(lifetime) || r.Flags != 0x11 || memmap.Uint32(0x5D4594, 3536) != 0 {
						t.Fatal("sudden death reset", r)
					}
					if legacy.PortTestPlayerStateQuery("threshold", nil, nil) != int(int32(threshold)) {
						t.Fatal("threshold query")
					}
					rows = append(rows, r)
				})
			}
		}
	}
	spellbookCapture(t, "player-state-sudden-death", rows, "4da2d318ac3783dbbd75c0a22ab3500e2c6313b1fe7863d35b5e7bb4559a68f3")
}
func TestPlayerStateElapsed(t *testing.T) {
	o := newMatchRosterOwner(t)
	serverConfigOwnBytes(t, 0x5D4594, 3520, 4)
	type row struct {
		FPS          int
		Start, Delta uint32
		Result       int
	}
	var rows []row
	for _, fps := range []int{1, 30, 60} {
		o.s.SetTickRate(uint32(fps))
		for _, start := range []uint32{0, 123, 0xfffffff0} {
			for _, delta := range []uint32{0, uint32(20*fps - 1), uint32(20 * fps), uint32(20*fps + 1), 0x80000000, 0xffffffff} {
				*memmap.PtrUint32(0x5D4594, 3520) = start
				o.s.SetFrame(start + delta)
				got := legacy.PortTestPlayerStateQuery("elapsed", nil, nil)
				if got != bool2int(delta > uint32(20*fps)) {
					t.Fatal(fps, start, delta, got)
				}
				rows = append(rows, row{fps, start, delta, got})
			}
		}
	}
	spellbookCapture(t, "player-state-elapsed", rows, "40d2325331f8d923e2871f5b2194efafe466d3841045006933216f83e5ecafe1")
}

func TestPlayerStateReentry(t *testing.T) {
	newMatchRosterOwner(t)
	type row struct{ Input, Result, Stored int32 }
	var rows []row
	for _, value := range []int32{0, 1, -1, 2147483647, -2147483648} {
		got := legacy.PortTestPlayerStateReentry(value)
		stored := memmap.Int32(0x5D4594, 3508)
		if got != value || stored != value {
			t.Fatal("reentry storage", value, got, stored)
		}
		rows = append(rows, row{value, got, stored})
	}
	spellbookCapture(t, "player-state-reentry", rows, "b125e602a59a40b515c26d30ca16e8b1f539d38faf15b67959b7d445b26145b9")
}
