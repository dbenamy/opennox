//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestGameplayReportsNotifications(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	type check struct {
		data               []byte
		recipient, ordered byte
		priority           uint32
	}
	var checks []check
	for _, op := range []int{13, 27, 30, 47, 49, 51} {
		for _, recipient := range []uint32{1, 7, 31} {
			for _, value := range []uint32{0, 1, 255, 256, 32767, 32768, 0xffffffff} {
				for _, mode := range []uint32{0, 0x8000000} {
					s := gameplayReportsBase(op, reportValue(recipient), reportObject(1), reportValue(value), reportValue(value^0xa5a5a5a5))
					s.Lifecycle.GameFlags = mode
					a := s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack
					a.ActorWords = map[int]uint32{36: value, 4: 23}
					id := uint16(value)
					if value >= 32768 {
						id = 0
					}
					c := check{recipient: byte(recipient), ordered: 1, priority: 1}
					switch op {
					case 13:
						c.data = []byte{96, byte(id), byte(id >> 8)}
						c.priority = 0
					case 27:
						a.Controls.Reports.Args = [5]legacy.PortTestGameplayReportArg{reportValue(recipient), reportValue(value)}
						c.data = make([]byte, 5)
						c.data[0] = 73
						binary.LittleEndian.PutUint32(c.data[1:], value)
						c.ordered = 0
					case 30:
						c.data = []byte{77, byte(id), byte(id >> 8), 23, 0}
						c.priority = 0
					case 47:
						code := uint16(value)
						if mode != 0 {
							code |= 0x8000
						}
						c.data = []byte{109, byte(code), byte(code >> 8)}
					case 49:
						c.data = []byte{220, byte(id), byte(id >> 8)}
					case 51:
						c.data = []byte{104, byte(id), byte(id >> 8), 0, 0, 0, 0, byte(value)}
						binary.LittleEndian.PutUint32(c.data[3:], value^0xa5a5a5a5)
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
			t.Fatalf("notification case%d message count%d", i, len(ps))
		}
		p := ps[0]
		if !bytes.Equal(p.Data, c.data) || p.Recipient != c.recipient || p.Ordered != c.ordered || p.A4 != 0 || p.A5 != c.priority {
			t.Fatalf("notification case%d fields or routing", i)
		}
	}
	gameplayReportsCapture(t, "notifications", out)
}

func TestGameplayReportsCreatureResources(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	type check struct {
		opcode     byte
		code, kind uint16
		health     bool
	}
	var checks []check
	for _, op := range []int{46, 48} {
		for _, id := range []uint32{0, 1, 32767, 32768, 65535} {
			for _, kind := range []uint32{0, 1, 32767, 32768, 65535} {
				for _, mode := range []uint32{0, 0x8000000} {
					for _, health := range []bool{false, true} {
						s := gameplayReportsBase(op, reportValue(7), reportObject(1))
						s.Lifecycle.GameFlags = mode
						p := s.Callbacks.Shop
						p.Resources.NoHealth = !health
						p.Resources.HP = 23
						p.Resources.MaxHP = 99
						p.TemporaryUpdates.World.Objectives.Attack.ActorWords = map[int]uint32{36: id, 4: kind}
						code := uint16(id)
						if id >= 32768 {
							code = 0
						}
						k := uint16(kind)
						opcode := byte(108)
						if op == 48 {
							opcode = 219
						} else if mode != 0 {
							k |= 0x8000
						}
						cases = append(cases, s)
						checks = append(checks, check{opcode, code, k, health})
					}
				}
			}
		}
	}
	out := controlsRun(t, cases)
	for i, r := range out {
		ps := r.Callbacks.Shop.Sequence[0].Packets
		c := checks[i]
		n := 1
		if c.health {
			n = 2
		}
		if len(ps) != n {
			t.Fatalf("creature case%d count%d", i, len(ps))
		}
		// Queue snapshots are newest first: the health report follows acquisition.
		for j, p := range ps {
			want := []byte{c.opcode, byte(c.code), byte(c.code >> 8), byte(c.kind), byte(c.kind >> 8)}
			if c.health && j == 0 {
				want = []byte{221, byte(c.code), byte(c.code >> 8), 23, 0, 99, 0}
			}
			if !bytes.Equal(p.Data, want) || p.Recipient != 7 || p.Ordered != 1 || p.A4 != 0 || p.A5 != 1 {
				t.Fatalf("creature case%d message%d fields or order", i, j)
			}
		}
	}
	gameplayReportsCapture(t, "creature-resources", out)
}
