//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"testing"
)

func TestGameplayReportsPlayerFields(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	type check struct {
		data               []byte
		recipient, ordered byte
	}
	var checks []check
	for _, op := range []int{5, 6, 9, 21, 23, 24, 25, 69} {
		for _, recipient := range []uint32{1, 7, 31} {
			for _, value := range []uint32{0, 1, 127, 255, 256, 65535} {
				for _, class := range []byte{0, 1, 2} {
					s := gameplayReportsBase(op, reportValue(recipient), reportObject(1), reportValue(value))
					p := s.Callbacks.Shop
					o := p.TemporaryUpdates.World.Objectives
					a := o.Attack
					p.Resources.PlayerClass = class
					a.ActorWords = map[int]uint32{36: 123, 16: value, 28: math.Float32bits(float32(value) + 0.75)}
					a.UpdateWords = map[int]uint32{4: value, 8: 65535 - value, 88: value << 24, 320: value}
					o.PlayerDataWords[0] = map[int]uint32{2064: recipient, 2164: value, 2248: uint32(class) << 24}
					c := check{recipient: byte(recipient)}
					switch op {
					case 5:
						a.Controls.Reports.Args = [5]legacy.PortTestGameplayReportArg{reportObject(1), reportObject(1), reportValue(value)}
						c.data = []byte{218, 123, 0, byte(value)}
					case 6:
						a.Controls.Reports.Args = [5]legacy.PortTestGameplayReportArg{reportObject(1)}
						c.data = make([]byte, 5)
						c.data[0] = 110
						binary.LittleEndian.PutUint32(c.data[1:], value)
					case 9:
						a.Controls.Reports.Args = [5]legacy.PortTestGameplayReportArg{reportObject(1)}
						c.data = make([]byte, 5)
						c.data[0] = 102
						binary.LittleEndian.PutUint32(c.data[1:], value)
					case 21:
						c.data = []byte{71, byte(value)}
						c.ordered = 1
					case 23:
						c.data = make([]byte, 5)
						c.data[0] = 74
						binary.LittleEndian.PutUint32(c.data[1:], value)
					case 24:
						if class != 0 {
							max := 65535 - value
							c.data = []byte{222, 123, 0, byte(value), byte(value >> 8), byte(max), byte(max >> 8)}
						}
					case 25:
						if class != 0 {
							c.data = []byte{69, 123, 0, byte(value), byte(value >> 8)}
						}
					case 69:
						c.data = []byte{240, 4, byte(value), 123, 0}
					}
					cases = append(cases, s)
					checks = append(checks, c)
				}
			}
		}
	}
	out := controlsRun(t, cases)
	for i, r := range out {
		ps := r.Callbacks.Shop.Sequence[0].Packets
		c := checks[i]
		if c.data == nil {
			if len(ps) != 0 {
				t.Fatalf("class-filtered player case%d", i)
			}
			continue
		}
		if len(ps) != 1 {
			t.Fatalf("player case%d message count%d", i, len(ps))
		}
		p := ps[0]
		if !bytes.Equal(p.Data, c.data) || p.Recipient != c.recipient || p.Ordered != c.ordered || p.A4 != 0 || p.A5 != 1 {
			t.Fatalf("player case%d defined fields or routing", i)
		}
	}
	gameplayReportsCapture(t, "player-fields", out)
}

func TestGameplayReportsHeightFields(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	var checks [][]byte
	for _, recipient := range []uint32{1, 7, 31} {
		for _, flags := range []uint32{0, 0x20} {
			for _, height := range []float32{-256.75, -255.5, -1.75, -0.5, 0, 0.5, 1.75, 255.5, 256.75} {
				s := gameplayReportsBase(43, reportValue(recipient), reportObject(1))
				s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.ActorWords = map[int]uint32{36: 123, 20: flags, 104: math.Float32bits(height), 108: math.Float32bits(height + 1.25), 116: math.Float32bits(height - 1.25)}
				var want []byte
				if flags != 0 {
					want = []byte{159, 123, 0, byte(int32(height)), byte(int32(height + 1.25)), byte(int32(height - 1.25))}
				} else {
					opcode := byte(94)
					h := height
					if h < 0 {
						opcode = 95
						h = -h
					}
					want = []byte{opcode, 123, 0, byte(int32(h))}
				}
				cases = append(cases, s)
				checks = append(checks, want)
			}
		}
	}
	out := controlsRun(t, cases)
	for i, r := range out {
		step := r.Callbacks.Shop.Sequence[0]
		index := i / 18
		if len(step.Packets) != 0 || len(step.ResourceMessages) != 3 || !bytes.Equal(step.ResourceMessages[index], checks[i]) {
			t.Fatalf("height case%d defined fields", i)
		}
		for j, b := range step.ResourceMessages {
			if j != index && len(b) != 0 {
				t.Fatalf("height case%d routing", i)
			}
		}
	}
	gameplayReportsCapture(t, "height-fields", out)
}
