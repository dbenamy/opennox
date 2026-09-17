//go:build porttest

package opennox

import (
	"bytes"
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/ntype"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestReliableReportsDeliveryCapacity(t *testing.T) {
	o := newReliableReportsOwner(t)
	var rows []struct {
		Name  string
		Bytes []byte
		State legacy.PortTestReliableReportState
	}
	defer func() {
		spellbookCapture(t, "reliable-reports-delivery-capacity", rows, "2e2597e525418bbb3117ec89641f80fee0300e87158a8cd8e2b64f2173c5679b")
	}()
	for _, recipient := range []int{1, 31} {
		for _, kind := range []int{0, 1} {
			if recipient != 31 && kind == 0 {
				continue
			}
			for _, ordered := range []bool{false, true} {
				for _, occupied := range []int{0, 2038, 2039, 2040, 2041, 2042, 2043, 2044, 2045, 2047, 2048} {
					for _, reserved := range []int{0, 7} {
						name := fmt.Sprintf("recipient%d/kind%d/ordered%t/occupied%d/reserved%d", recipient, kind, ordered, occupied, reserved)
						t.Run(name, func(t *testing.T) {
							o.reset()
							legacy.PortTestReliableReports(0, 0, 0, 0, nil, nil, 0, false)
							netPlayerBufSize = reserved
							fill := bytes.Repeat([]byte{0x44}, occupied)
							if !o.s.NetList.AddToMsgListCli(ntype.PlayerInd(recipient), netlist.Kind(kind), fill) {
								t.Fatal("fill")
							}
							legacy.PortTestReliableReports(8, recipient, 0, 0, []byte{77, 88, 99}, nil, 0, ordered)
							legacy.PortTestReliableReports(19, recipient, 0, uint32(kind), nil, nil, 0, false)
							packet := o.s.NetList.CopyPacketsA(ntype.PlayerInd(recipient), netlist.Kind(kind))
							state := o.state()
							want := bytes.Clone(fill)
							ts := []byte{170}
							if recipient == 31 {
								ts = append(ts, 123, 0, 0, 0)
							}
							extra := reserved
							if recipient == 31 {
								extra = 0
							}
							timestampFits := occupied+len(ts)+extra <= 2048
							messageFits := false
							if timestampFits {
								want = append(want, ts...)
								msg := []byte{77, 88, 99}
								if ordered {
									msg = append([]byte{204, 0, 0, 3}, msg...)
								}
								messageFits = len(want)+len(msg)+extra <= 2048
								if messageFits {
									want = append(want, msg...)
								}
							}
							if !bytes.Equal(packet, want) {
								t.Fatalf("bytes%d want%d", len(packet), len(want))
							}
							if (state.Nodes[0].Sent != 0) != messageFits {
								t.Fatal("send state after capacity failure")
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

func TestReliableReportsClientWrapper(t *testing.T) {
	o := newReliableReportsOwner(t)
	var rows []struct {
		Name   string
		Return uint32
		Bytes  []byte
		State  legacy.PortTestReliableReportState
	}
	for _, flags := range []uint32{0, 1} {
		for _, to := range []int{1, 31, 129, 255} {
			if flags == 1 && to != 31 {
				continue
			}
			for _, full := range []bool{false, true} {
				name := fmt.Sprintf("flags%d/to%d/full%t", flags, to, full)
				t.Run(name, func(t *testing.T) {
					restore := noxflags.PortTestGameFlags(noxflags.GameFlag(flags))
					defer restore()
					o.reset()
					if flags == 1 && full {
						o.s.NetList.AddToMsgListCli(31, netlist.Kind0, bytes.Repeat([]byte{0x44}, 2048))
					}
					rv := legacy.PortTestReliableReports(12, to, 0, 0, []byte{77, 88}, nil, 0, false)
					var packet []byte
					if flags == 1 {
						packet = o.s.NetList.CopyPacketsA(31, netlist.Kind0)
					}
					state := o.state()
					wantReturn := uint32(1)
					if flags == 0 && to >= 128 {
						wantReturn = 0
					}
					if rv != wantReturn {
						t.Fatal("wrapper return")
					}
					if flags == 1 {
						if len(state.Nodes) != 0 {
							t.Fatal("host wrapper enqueued reliable")
						}
						wantLen := 2
						if full {
							wantLen = 2048
						}
						if len(packet) != wantLen {
							t.Fatal("host wrapper data")
						}
					} else if (len(state.Nodes) == 1) != (to < 128) {
						t.Fatal("client recipient restriction")
					}
					rows = append(rows, struct {
						Name   string
						Return uint32
						Bytes  []byte
						State  legacy.PortTestReliableReportState
					}{name, rv, packet, state})
				})
			}
		}
	}
	spellbookCapture(t, "reliable-reports-client-wrapper", rows, "4471530029de3c5425ef0cbcc7840114c6bde86aa3d04212b798161ab8997b25")
}

func TestReliableReportsDeliveryFlags(t *testing.T) {
	o := newReliableReportsOwner(t)
	engine := noxflags.GetEngine()
	t.Cleanup(func() { noxflags.ResetEngine(); noxflags.SetEngine(engine) })
	var rows []struct {
		Name  string
		Bytes []byte
		State legacy.PortTestReliableReportState
	}
	for _, flags := range []uint32{0, 1} {
		for _, status := range []uint32{0, 16} {
			for _, replay := range []bool{false, true} {
				name := fmt.Sprintf("flags%d/status%d/replay%t", flags, status, replay)
				t.Run(name, func(t *testing.T) {
					restore := noxflags.PortTestGameFlags(noxflags.GameFlag(flags))
					defer restore()
					o.reset()
					legacy.PortTestReliableReports(0, 0, 0, 0, nil, nil, 0, false)
					noxflags.UnsetEngine(noxflags.EngineReplayRead)
					if replay {
						noxflags.SetEngine(noxflags.EngineReplayRead)
					}
					objectXferSetWord(o.s.Players.ByInd(1).C(), 3680, status)
					legacy.PortTestReliableReports(13, 1, 0, 0, []byte{77}, nil, 0, false)
					legacy.PortTestReliableReports(19, 1, 0, 1, nil, nil, 0, false)
					packet := o.s.NetList.CopyPacketsA(1, netlist.Kind1)
					state := o.state()
					sent := flags == 0 || status&16 != 0
					if (len(packet) > 0) != sent || (len(state.Nodes) == 0) != (sent && replay) {
						t.Fatal("delivery gate/replay acknowledgement")
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
	spellbookCapture(t, "reliable-reports-delivery-flags", rows, "d084fafe1d4e276eb1dec1d600350d6de39315ed270eb469f8a6ad3271583ef5")
}
