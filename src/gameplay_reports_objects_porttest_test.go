//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestGameplayReportsObjectFields(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	type check struct {
		data               []byte
		recipient, ordered byte
		priority           uint32
	}
	var checks []check
	for _, op := range []int{7, 8, 10, 14, 22, 41, 60, 61, 63, 68, 70, 71} {
		for _, recipient := range []uint32{1, 7, 31} {
			for _, id := range []uint32{0, 1, 32767, 32768, 65535} {
				for _, value := range []uint32{0, 1, 256, 0xffffffff} {
					s := gameplayReportsBase(op, reportValue(recipient), reportObject(1), reportValue(value), reportValue(value^0x55))
					a := s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack
					a.ActorWords = map[int]uint32{36: id, 20: value, 132: value, 340: value, 440: value}
					// The normal player owner uses the network ID; out-of-range IDs resolve to zero.
					code := uint16(id)
					if id >= 32768 {
						code = 0
					}
					lo, hi := byte(code), byte(code>>8)
					c := check{recipient: byte(recipient), ordered: 1, priority: 1}
					switch op {
					case 7, 8, 41:
						opcode := byte(107)
						if op == 8 {
							opcode = 101
						}
						if op == 41 {
							opcode = 90
						}
						c.data = []byte{opcode, lo, hi, 0, 0, 0, 0}
						binary.LittleEndian.PutUint32(c.data[3:], value)
						c.ordered = 0
						if op == 41 {
							c.ordered = 1
						}
					case 10:
						c.data = []byte{100, lo, hi, byte(value), byte(value ^ 0x55)}
						c.priority = 0
					case 14:
						c.data = []byte{97, lo, hi}
						c.priority = 0
					case 22:
						c.data = []byte{91, byte(value)}
						c.ordered = 0
					case 60:
						c.data = []byte{224, lo, hi, byte(value)}
						c.priority = 0
					case 61:
						c.data = []byte{225, lo, hi}
						c.priority = 0
					case 63:
						opcode := byte(53)
						if value == 1 {
							opcode = 52
						}
						c.data = []byte{opcode, lo, hi}
					case 68:
						c.data = []byte{240, 1, byte(id), byte(id >> 8)}
					case 70, 71:
						subtype := byte(22)
						if op == 71 {
							subtype = 23
						}
						c.data = []byte{240, subtype, byte(value), byte(id), byte(id >> 8)}
						c.ordered = 0
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
		if len(ps) != 1 {
			t.Fatalf("object case%d message count%d", i, len(ps))
		}
		p := ps[0]
		if !bytes.Equal(p.Data, c.data) || p.Recipient != c.recipient || p.Ordered != c.ordered || p.A4 != 0 || p.A5 != c.priority {
			t.Fatalf("object case%d data %x want %x; recipient %d/%d order %d/%d priority %d/%d", i, p.Data, c.data, p.Recipient, c.recipient, p.Ordered, c.ordered, p.A5, c.priority)
		}
	}
	gameplayReportsCapture(t, "object-fields", out)
}

func TestGameplayReportsHealthFields(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	type check struct {
		data      []byte
		recipient byte
		priority  uint32
	}
	var checks []check
	for _, op := range []int{15, 16, 17, 20} {
		for _, recipient := range []uint32{1, 7, 31} {
			for _, hp := range []uint16{0, 1, 255, 256, 511, 512, 32767, 65535} {
				for _, max := range []uint16{1, 100, 65535} {
					s := gameplayReportsBase(op, reportValue(recipient), reportObject(1))
					p := s.Callbacks.Shop
					p.Resources.HP = hp
					p.Resources.MaxHP = max
					p.Resources.NoHealth = false
					p.TemporaryUpdates.World.Objectives.Attack.ActorWords = map[int]uint32{36: 123}
					c := check{recipient: byte(recipient), priority: 1}
					switch op {
					case 15, 20:
						opcode := byte(221)
						if op == 20 {
							opcode = 68
							c.priority = 0
						}
						c.data = []byte{opcode, 123, 0, byte(hp), byte(hp >> 8), byte(max), byte(max >> 8)}
					case 16:
						c.data = []byte{65, 123, 0, byte(hp >> 1)}
					case 17:
						c.data = []byte{196, 12, 123, 0, byte(uint32(hp) * 100 / uint32(max))}
					}
					cases = append(cases, s)
					checks = append(checks, c)
				}
			}
			s := gameplayReportsBase(op, reportValue(recipient), reportObject(1))
			s.Callbacks.Shop.Resources.NoHealth = true
			cases = append(cases, s)
			checks = append(checks, check{})
		}
	}
	out := controlsRun(t, cases)
	for i, r := range out {
		step := r.Callbacks.Shop.Sequence[0]
		c := checks[i]
		if c.data == nil {
			if len(step.Packets) != 0 || step.Return != 0 {
				t.Fatalf("missing health case%d", i)
			}
			continue
		}
		if len(step.Packets) != 1 {
			t.Fatalf("health case%d message count%d", i, len(step.Packets))
		}
		p := step.Packets[0]
		if !bytes.Equal(p.Data, c.data) || p.Recipient != c.recipient || p.Ordered != 1 || p.A4 != 0 || p.A5 != c.priority {
			t.Fatalf("health case%d defined fields or routing", i)
		}
	}
	gameplayReportsCapture(t, "health-fields", out)
}
