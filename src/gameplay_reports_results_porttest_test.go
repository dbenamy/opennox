//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestGameplayReportsWinners(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	var checks [][]byte
	for _, op := range []int{31, 32, 33, 34} {
		for _, id := range []uint32{0, 1, 127, 255} {
			for _, value := range []uint32{0, 1, 255, 256, 0xffffffff} {
				for _, present := range []bool{false, true} {
					if op == 33 && !present {
						continue
					}
					s := gameplayReportsBase(op)
					s.Owner.Frame = 0xffff1234
					a := s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack
					a.ActorWords = map[int]uint32{36: id}
					a.RecordWords = map[int]uint32{56: id << 8}
					ref := reportValue(0)
					if present {
						ref = legacy.PortTestGameplayReportArg{Kind: "record"}
						if op == 31 {
							ref = reportObject(1)
						}
					}
					a.Controls.Reports.Args = [5]legacy.PortTestGameplayReportArg{ref, reportValue(value)}
					opcode := map[int]byte{31: 88, 32: 89, 33: 86, 34: 87}[op]
					code := uint16(id)
					if !present {
						code = 0
						if op == 34 {
							code = 65535
						}
					}
					v := byte(value)
					if op == 33 {
						v = 0
					}
					want := []byte{opcode, byte(code), byte(code >> 8), v, 0, 0, 0, 0}
					binary.LittleEndian.PutUint32(want[4:], s.Owner.Frame)
					cases = append(cases, s)
					checks = append(checks, want)
				}
			}
		}
	}
	out := controlsRun(t, cases)
	for i, r := range out {
		ps := r.Callbacks.Shop.Sequence[0].Packets
		if len(ps) != 1 {
			t.Fatalf("winner case%d message count%d", i, len(ps))
		}
		p := ps[0]
		if !bytes.Equal(p.Data, checks[i]) || p.Recipient != 255 || p.Ordered != 1 || p.A4 != 0 || p.A5 != 1 {
			t.Fatalf("winner case%d fields or broadcast", i)
		}
	}
	gameplayReportsCapture(t, "winners", out)
}

func TestGameplayReportsScoreChanges(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	type check struct{ score, deaths uint32 }
	var checks []check
	for _, op := range []int{36, 37, 38} {
		for _, value := range []uint32{0, 1, 127, 65535, 0x7fffffff, 0x80000000, 0xffffffff} {
			for _, delta := range []uint32{0, 1, 2, 0x7fffffff, 0x80000000, 0xffffffff} {
				for _, player := range []bool{false, true} {
					s := gameplayReportsBase(op, reportObject(1), reportValue(delta))
					s.Lifecycle.GameFlags = 0
					p := s.Callbacks.Shop
					o := p.TemporaryUpdates.World.Objectives
					a := o.Attack
					a.ActorWords = map[int]uint32{36: 123}
					target := 4
					if player {
						target = 1
					}
					o.PlayerDataWords[0] = map[int]uint32{2136: value, 2140: value}
					a.Controls.Reports.Calls = map[uint32][5]legacy.PortTestGameplayReportArg{0: {reportObject(target), reportValue(delta)}, 1: {reportObject(1)}}
					p.Sequence = []legacy.PortTestShopAction{{Op: 1800 + op}, {Op: 1839, Value: 1}}
					c := check{value, value}
					if player {
						switch op {
						case 36:
							c.deaths++
						case 37:
							c.score += delta
						case 38:
							c.score -= delta
						}
					}
					cases = append(cases, s)
					checks = append(checks, c)
				}
			}
		}
	}
	out := controlsRun(t, cases)
	for i, r := range out {
		steps := r.Callbacks.Shop.Sequence
		c := checks[i]
		if len(steps[0].Packets) != 0 || len(steps[1].Packets) != 1 {
			t.Fatalf("score case%d message counts", i)
		}
		p := steps[1].Packets[0]
		want := []byte{78, 123, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		binary.LittleEndian.PutUint32(want[3:], c.score)
		binary.LittleEndian.PutUint32(want[7:], c.deaths)
		if !bytes.Equal(p.Data, want) || p.Recipient != 255 || p.Ordered != 1 || p.A4 != 0 || p.A5 != 1 {
			t.Fatalf("score case%d resulting score/deaths", i)
		}
	}
	gameplayReportsCapture(t, "score-changes", out)
}
