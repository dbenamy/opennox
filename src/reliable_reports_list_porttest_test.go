//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"reflect"
	"testing"
)

func TestReliableReportsAcknowledge(t *testing.T) {
	o := newReliableReportsOwner(t)
	var rows []struct {
		Name  string
		Known uint32
		State legacy.PortTestReliableReportState
	}
	defer func() {
		spellbookCapture(t, "reliable-reports-acknowledge", rows, "d467840257f4656b65ee2a2e4ee2afe15a67e9a47bf5f1124e87942f1ee5fa58")
	}()
	for _, to := range []int{0, 1, 7, 31, 128, 129, 159, 255} {
		for _, recipient := range []int{0, 1, 7, 31} {
			for _, initial := range []uint32{0, 2, 0x80000082, 0xffffffff} {
				for _, active := range []uint32{0, 2, 0x80000082, 0xffffffff} {
					for _, opcode := range []byte{48, 49, 50, 51, 52} {
						name := fmt.Sprintf("to%d/recipient%d/initial%x/active%x/opcode%d", to, recipient, initial, active, opcode)
						t.Run(name, func(t *testing.T) {
							o.reset()
							*o.words["mask"] = initial
							objectXferSetWord(o.units[0].CObj(), 148, 0xffffffff)
							if legacy.PortTestReliableReports(8, to, 0, 0, []byte{opcode, 1, 2}, &o.units[0], 0, false) != 1 {
								t.Fatal("enqueue")
							}
							if len(o.state().Nodes) == 0 {
								return
							} // An empty exclusion audience creates no node.
							*o.words["mask"] = active
							bit := uint32(1) << uint(recipient)
							legacy.PortTestReliableReports(15, recipient, 0, bit, nil, nil, 0, false)
							state := o.state()
							known := objectXferGetWord(o.units[0].CObj(), 148)
							wantKnown := uint32(0xffffffff)
							if opcode >= 49 && opcode <= 51 {
								wantKnown &^= bit
							}
							if known != wantKnown {
								t.Fatalf("known%x want%x", known, wantKnown)
							}
							audience := initial & active
							removed := to == recipient
							if to == 255 {
								removed = audience&bit == audience
							} else if to >= 128 {
								audience &^= uint32(1) << uint(to&127)
								removed = audience&bit == audience
							}
							if (len(state.Nodes) == 0) != removed {
								t.Fatalf("removed%t nodes%d", removed, len(state.Nodes))
							}
							if !removed {
								wantAck := uint32(0)
								if to >= 128 {
									wantAck = bit
								}
								if state.Nodes[0].Acknowledged != wantAck {
									t.Fatal("ack mask")
								}
							}
							rows = append(rows, struct {
								Name  string
								Known uint32
								State legacy.PortTestReliableReportState
							}{name, known, state})
						})
					}
				}
			}
		}
	}
}

func TestReliableReportsListEdits(t *testing.T) {
	o := newReliableReportsOwner(t)
	var rows []struct {
		Name   string
		Return uint32
		State  legacy.PortTestReliableReportState
	}
	defer func() {
		spellbookCapture(t, "reliable-reports-list-edits", rows, "a3f4d9d897beeba50650702e976e7a2d51bca3def42ba8fc04c833d7c2e298eb")
	}()
	for _, op := range []int{6, 7, 14, 16, 17} {
		for _, to := range []int{1, 7, 31, 129, 255} {
			for _, index := range []int{0, 2, 4} {
				name := fmt.Sprintf("op%d/to%d/index%d", op, to, index)
				t.Run(name, func(t *testing.T) {
					o.reset()
					for i, code := range []byte{48, 49, 50, 51, 52} {
						if legacy.PortTestReliableReports(8, 255, 0, 0, []byte{code, byte(i)}, nil, 0, false) != 1 {
							t.Fatal("enqueue")
						}
					}
					arg := uint32(0)
					if op == 16 && index != 0 {
						arg = 123
					}
					callTo := to
					if (op == 16 || op == 17) && to >= 128 {
						callTo = to & 31
					}
					rv := legacy.PortTestReliableReports(op, callTo, index, arg, []byte{50, 99}, nil, 0, false)
					state := o.state()
					var codes []byte
					for _, n := range state.Nodes {
						codes = append(codes, n.Data[0])
					}
					want := []byte{52, 51, 50, 49, 48}
					switch op {
					case 6:
						want = []byte{52, 48}
					case 7:
						want = append(want[:index], want[index+1:]...)
					case 14:
						if to >= 128 {
							want = []byte{50, 52, 51, 49, 48}
						} else {
							want = []byte{50, 52, 51, 50, 49, 48}
						}
					}
					if !reflect.DeepEqual(codes, want) {
						t.Fatalf("codes%v want%v", codes, want)
					}
					if op == 16 || op == 17 {
						for _, n := range state.Nodes {
							wantMask := uint32(0)
							if op == 17 || arg == 0 {
								wantMask = uint32(1) << uint(callTo)
							}
							if n.Acknowledged != wantMask {
								t.Fatal("acknowledgement mutation")
							}
						}
					}
					rows = append(rows, struct {
						Name   string
						Return uint32
						State  legacy.PortTestReliableReportState
					}{name, rv, state})
				})
			}
		}
	}
}
