//go:build porttest

package opennox

import (
	"bytes"
	"fmt"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/ntype"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
)

func TestReliableReportsDelivery(t *testing.T) {
	o := newReliableReportsOwner(t)
	var rows []struct {
		Name  string
		Bytes []byte
		State legacy.PortTestReliableReportState
	}
	defer func() {
		spellbookCapture(t, "reliable-reports-delivery", rows, "6074dcef71eb0e00fda7637b61d6343f20857032a385aab910e3fff0d58af070")
	}()
	for _, to := range []int{1, 7, 31, 129, 159, 255} {
		for _, recipient := range []int{1, 7, 31} {
			for _, ordered := range []bool{false, true} {
				for _, size := range []int{1, 150} {
					for _, relation := range []int{0, 1, 2, 3} {
						for _, priority := range []int{0, 1} {
							for _, active := range []uint32{0, 0x80000082} {
								name := fmt.Sprintf("to%d/recipient%d/ordered%t/size%d/relation%d/priority%d/active%x", to, recipient, ordered, size, relation, priority, active)
								t.Run(name, func(t *testing.T) {
									o.reset()
									legacy.PortTestReliableReports(0, 0, 0, 0, nil, nil, 0, false)
									data := bytes.Repeat([]byte{77}, size)
									var obj *server.Object
									if relation != 0 {
										obj = &o.units[0]
										objectXferSetWord(obj.CObj(), 148, 0xffffffff)
										objectXferSetWord(obj.CObj(), 16, 0)
										if relation == 2 {
											objectXferSetWord(obj.CObj(), 148, 0)
										}
										if relation == 3 {
											objectXferSetWord(obj.CObj(), 16, 0x20)
										}
									}
									*memmap.PtrUint16(0x5D4594, 1565524+uintptr(recipient)*2) = 65535
									legacy.PortTestReliableReports(8, to, 0, 0, data, obj, priority, ordered)
									*o.words["mask"] = active
									legacy.PortTestReliableReports(19, recipient, 0, 1, nil, nil, 0, false)
									packet := o.s.NetList.CopyPacketsA(ntype.PlayerInd(recipient), netlist.Kind1)
									state := o.state()
									addressed := to == 255 || to == recipient || to >= 128 && recipient != to&127
									sent := addressed && relation != 2 && (priority == 0 || active&(uint32(1)<<uint(recipient)) != 0)
									var want []byte
									if sent {
										want = []byte{170}
										if recipient == 31 {
											want = append(want, 123, 0, 0, 0)
										}
										if ordered {
											want = append(want, 204, 255, 255, byte(size))
										}
										want = append(want, data...)
									}
									if !bytes.Equal(packet, want) {
										t.Fatalf("packet %x want %x", packet, want)
									}
									if sent {
										if len(state.Nodes) != 1 {
											t.Fatal("sent node missing")
										}
										n := state.Nodes[0]
										bit := uint32(1) << uint(recipient)
										if n.Sent != bit || n.LastSent[recipient] != 123 || n.Countdown[recipient] != 60 {
											t.Fatalf("delivery state %+v", n)
										}
									}
									if relation == 3 && len(state.Nodes) > 0 && state.Nodes[0].Related != 0 {
										t.Fatal("destroyed reference retained")
									}
									rows = append(rows, struct {
										Name  string
										Bytes []byte
										State legacy.PortTestReliableReportState
									}{name, packet, state})
								})
							}
						}
					}
				}
			}
		}
	}
}

func TestReliableReportsRetryTiming(t *testing.T) {
	o := newReliableReportsOwner(t)
	var rows []struct {
		Name  string
		Bytes []byte
		State legacy.PortTestReliableReportState
	}
	defer func() {
		spellbookCapture(t, "reliable-reports-retries", rows, "6c01b5a99276b7aa3c9807054dd686c46cad505be1d70f3104cd6e91c7e78e17")
	}()
	for _, recipient := range []int{1, 31} {
		for _, ordered := range []bool{false, true} {
			for _, rate := range []uint32{1, 30, 255} {
				for _, budget := range []byte{0, 1, 2} {
					name := fmt.Sprintf("recipient%d/ordered%t/rate%d/budget%d", recipient, ordered, rate, budget)
					t.Run(name, func(t *testing.T) {
						o.reset()
						*o.words["rate"] = rate
						legacy.PortTestReliableReports(0, 0, 0, 0, nil, nil, 0, false)
						*memmap.PtrUint8(0x5D4594, 1565124+uintptr(recipient)*12) = budget
						for i := 0; i < 3; i++ {
							legacy.PortTestReliableReports(8, recipient, 0, 0, []byte{77, byte(i)}, nil, 0, ordered)
						}
						for step := 0; step < 65; step++ {
							o.s.SetFrame(uint32(123 + step))
							legacy.PortTestReliableReports(19, recipient, 0, 1, nil, nil, 0, false)
							packet := o.s.NetList.CopyPacketsA(ntype.PlayerInd(recipient), netlist.Kind1)
							state := o.state()
							if len(state.Nodes) != 3 {
								t.Fatal("retry removed node")
							}
							if step == 0 {
								for _, n := range state.Nodes {
									if n.LastSent[recipient] != 123 || n.Countdown[recipient] != byte(60/rate) || n.Retries != 0 {
										t.Fatal("first delivery")
									}
								}
							} else if budget == 0 && len(packet) != 0 {
								t.Fatal("zero retry budget sent")
							}
							rows = append(rows, struct {
								Name  string
								Bytes []byte
								State legacy.PortTestReliableReportState
							}{fmt.Sprintf("%s/step%d", name, step), packet, state})
						}
						// ACK matches delivery frame, not the ordered sequence carried on the wire.
						frame := o.state().Nodes[0].LastSent[recipient]
						legacy.PortTestReliableReports(16, recipient, 0, frame, nil, nil, 0, false)
						for _, n := range o.state().Nodes {
							if n.LastSent[recipient] == frame {
								t.Fatal("ack left matching timestamp")
							}
						}
					})
				}
			}
		}
	}
}

func TestReliableReportsSequenceWrap(t *testing.T) {
	o := newReliableReportsOwner(t)
	var states []legacy.PortTestReliableReportState
	for _, start := range []uint16{0, 32767, 65534, 65535} {
		for _, to := range []int{1, 31, 129, 255} {
			o.reset()
			seq := (*[32]uint16)(memmap.PtrOff(0x5D4594, 1565524))
			for i := range seq {
				seq[i] = start
			}
			for i := 0; i < 3; i++ {
				data := []byte{77, byte(i)}
				if legacy.PortTestReliableReports(11, to, 0, 0, data, nil, 0, false) != 1 {
					t.Fatal("ordered wrapper")
				}
				data[0] = 0 // Constructor owns its payload copy.
				state := o.state()
				if state.Nodes[0].Data[0] != 77 {
					t.Fatal("payload aliases caller")
				}
				states = append(states, state)
			}
			for i, n := range o.state().Nodes {
				if n.Data[1] != byte(2-i) || n.Ordered != 1 {
					t.Fatal("insertion order")
				}
				for slot, v := range n.Sequence {
					included := to == slot || to == 255 && *o.words["mask"]&(1<<uint(slot)) != 0 || to >= 128 && to != 255 && slot != to&127 && *o.words["mask"]&(1<<uint(slot)) != 0
					want := uint16(0)
					if included {
						want = start + uint16(2-i)
					}
					if v != want {
						t.Fatalf("sequence %d want %d", v, want)
					}
				}
			}
		}
	}
	spellbookCapture(t, "reliable-reports-sequence-wrap", states, "764a0d28a7601fd19df384a70c81afed3bba52c0b6dc8550b2ad956b26428a9a")
}
