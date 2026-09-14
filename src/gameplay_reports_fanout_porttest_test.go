//go:build porttest

package opennox

import (
	"bytes"
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"testing"
)

func TestGameplayReportsInterestingFanout(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	type check struct {
		players, op int
		pending     bool
	}
	var checks []check
	for _, op := range []int{0, 1, 2} {
		for players := 0; players <= 3; players++ {
			for _, pending := range []bool{false, true} {
				for _, frame := range []uint32{0, 1, 0xffffffff} {
					s := gameplayReportsBase(op, reportObject(1))
					s.Owner.Frame = frame
					p := s.Callbacks.Shop
					o := p.TemporaryUpdates.World.Objectives
					o.Players = players
					for j := 0; j < 3; j++ {
						flag := uint32(0)
						if pending {
							flag = 1
						}
						o.PlayerUpdateWords[j] = map[int]uint32{248: 0xffffffff, 252: 1, 256: 123, 260: flag}
					}
					o.Attack.ActorWords = map[int]uint32{36: 123, 4: 23}
					p.Sequence = []legacy.PortTestShopAction{{Op: 1800 + op}, {Op: 1800 + op}}
					cases = append(cases, s)
					checks = append(checks, check{players, op, pending})
				}
			}
		}
	}
	out := controlsRun(t, cases)
	for i, r := range out {
		c := checks[i]
		first := c.players
		if c.op == 1 && !c.pending {
			first = 0
		}
		if c.op == 2 {
			first = 0
			if c.pending {
				first = c.players * c.players
			}
		}
		for j, step := range r.Callbacks.Shop.Sequence {
			count := first
			if c.op == 0 {
				count *= j + 1
			}
			if len(step.Packets) != count {
				t.Fatalf("fanout case%d step%d count%d want%d", i, j, len(step.Packets), count)
			}
			for _, p := range step.Packets {
				if p.Ordered != 0 || p.A4 != 0 || p.A5 != 1 || len(p.Data) != 7 || p.Data[0] != 210 || p.Data[5] != 2 || p.Data[6] != 2 {
					t.Fatalf("fanout case%d defined fields", i)
				}
				valid := false
				for _, recipient := range []byte{1, 7, 31}[:c.players] {
					if recipient == p.Recipient {
						valid = true
					}
				}
				if !valid {
					t.Fatalf("fanout case%d recipient", i)
				}
				if c.op != 2 && !bytes.Equal(p.Data, []byte{210, 123, 0, 23, 0, 2, 2}) {
					t.Fatalf("fanout case%d source fields", i)
				}
			}
		}
	}
	gameplayReportsCapture(t, "interesting-fanout", out)
}

func TestGameplayReportsEarthquake(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	var checks [][3][]byte
	for players := 0; players <= 3; players++ {
		for _, distance := range []float32{0, 1, 212, 299, 300, 301} {
			for _, amplitude := range []int32{-257, -1, 0, 1, 127, 255, 256} {
				s := gameplayReportsBase(45, legacy.PortTestGameplayReportArg{Kind: "record"}, reportValue(uint32(amplitude)))
				o := s.Callbacks.Shop.TemporaryUpdates.World.Objectives
				o.Players = players
				o.Attack.RecordWords = map[int]uint32{0: 0, 4: 0}
				var expected [3][]byte
				for j := 0; j < 3; j++ {
					x, y := distance, float32(0)
					if j == 1 {
						y = distance
					}
					if j == 2 {
						x = -distance
					}
					o.PlayerDataWords[j] = map[int]uint32{3632: math.Float32bits(x), 3636: math.Float32bits(y)}
					sq := float64(x)*float64(x) + float64(y)*float64(y)
					if j < players && sq < 90000 {
						expected[j] = []byte{151, byte(int64((1 - sq*0.000011111111) * float64(amplitude)))}
					}
				}
				cases = append(cases, s)
				checks = append(checks, expected)
			}
		}
	}
	out := controlsRun(t, cases)
	for i, r := range out {
		step := r.Callbacks.Shop.Sequence[0]
		if len(step.Packets) != 0 {
			t.Fatalf("earthquake case%d unexpected queued messages", i)
		}
		for j, want := range checks[i] {
			if !bytes.Equal(step.ResourceMessages[j], want) {
				t.Fatalf("earthquake case%d recipient%d falloff", i, j)
			}
		}
	}
	gameplayReportsCapture(t, "earthquake", out)
}
