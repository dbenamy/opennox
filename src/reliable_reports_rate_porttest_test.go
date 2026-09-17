//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/ntype"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
	"unsafe"
)

func TestReliableReportsRateAdaptation(t *testing.T) {
	o := newReliableReportsOwner(t)
	var rows []struct {
		Name           string
		Return, Status uint32
		State          legacy.PortTestReliableReportState
	}
	defer func() {
		spellbookCapture(t, "reliable-reports-rate-adaptation", rows, "5e9c04abba360f8eba93e3bc655b8110ffd323fedbd7c761a77d6fa9df9db4c8")
	}()
	for _, count := range []int{0, 1, 2, 3, 4, 5, 6, 11} {
		for _, budget := range []byte{1, 2} {
			for _, delay := range []byte{2, 3, 4, 5} {
				for _, rate := range []uint32{1, 30} {
					for _, mode := range []uint32{0, 1} {
						name := fmt.Sprintf("count%d/budget%d/delay%d/rate%d/mode%d", count, budget, delay, rate, mode)
						t.Run(name, func(t *testing.T) {
							o.reset()
							*o.words["rate"] = rate
							*o.words["rateMode"] = mode
							o.s.SetTickRate(1)
							pl := o.s.Players.ByInd(1)
							objectXferSetWord(pl.C(), 3680, 0)
							for i := 0; i < count; i++ {
								legacy.PortTestReliableReports(8, 1, 0, 0, []byte{77}, nil, 0, false)
							}
							// Other recipients and the explicitly excluded recipient do not add pressure.
							legacy.PortTestReliableReports(8, 7, 0, 0, []byte{78}, nil, 0, false)
							legacy.PortTestReliableReports(8, 129, 0, 0, []byte{79}, nil, 0, false)
							p := unsafe.Slice((*byte)(memmap.PtrOff(0x5D4594, 1565136)), 12)
							p[0] = budget
							p[1] = delay
							p[2] = 0
							legacy.PortTestReliableReports(1, 1, 0, 0, nil, nil, 0, false)
							high, low := binary.LittleEndian.Uint32(p[4:]), binary.LittleEndian.Uint32(p[8:])
							wantBudget, wantDelay := budget, delay
							kick := false
							if uint32(count) > high {
								wantDelay++
								if wantDelay > 5 {
									wantDelay = 5
									wantBudget = 2
									kick = budget == 2
								}
							} else if low > 0 && uint32(count) < low {
								if budget == 2 {
									wantBudget = 1
								} else {
									wantDelay--
									if wantDelay < 2 {
										wantDelay = 2
									}
								}
							}
							rv := legacy.PortTestReliableReports(18, 1, 0, 0, nil, nil, 0, false)
							state := o.state()
							status := objectXferGetWord(pl.C(), 3680)
							if p[0] != wantBudget || p[1] != wantDelay || p[2] != byte(rate) || (status&128 != 0) != kick {
								t.Fatalf("rate %v status%x", p, status)
							}
							divisor := uint32(1)
							if mode == 1 {
								divisor = rate
							}
							want := uint32(wantBudget) * uint32(wantDelay) * (1 / divisor)
							if rv != want {
								t.Fatalf("return%d want%d", rv, want)
							}
							if kick {
								for _, n := range state.Nodes {
									if n.To == 1 {
										t.Fatal("slow recipient retained")
									}
								}
							}
							rows = append(rows, struct {
								Name           string
								Return, Status uint32
								State          legacy.PortTestReliableReportState
							}{name, rv, status, state})
						})
					}
				}
			}
		}
	}
}

func TestReliableReportsPoolPressure(t *testing.T) {
	o := newReliableReportsOwner(t)
	var rows []struct {
		Name   string
		Return uint32
		State  legacy.PortTestReliableReportState
	}
	defer func() {
		spellbookCapture(t, "reliable-reports-pool-pressure", rows, "8a1660c562478ba1e4d028b51749293814a97043e90ed0610f22f97f17e6b7cf")
	}()
	for _, flags := range []uint32{0, 1, 2048, 2049} {
		t.Run(fmt.Sprintf("flags%d", flags), func(t *testing.T) {
			restore := noxflags.PortTestGameFlags(noxflags.GameFlag(flags))
			defer restore()
			o.reset()
			cap := 256
			if flags&2048 != 0 {
				cap = 512
			} else if flags&1 != 0 {
				cap = 3072
			}
			for i := 0; i <= cap; i++ {
				o.s.SetFrame(uint32(i + 1))
				data := []byte{77, byte(i), byte(i >> 8)}
				if legacy.PortTestReliableReports(8, 31, 0, 0, data, nil, 0, false) != 1 {
					t.Fatalf("enqueue%d", i)
				}
			}
			state := o.state()
			wantCount := cap
			if flags&2048 != 0 {
				wantCount++
			}
			if int(state.Capacity) != cap || len(state.Nodes) != wantCount {
				t.Fatalf("capacity%d nodes%d", state.Capacity, len(state.Nodes))
			}
			if state.Nodes[0].Frame != uint32(cap+1) || state.Nodes[len(state.Nodes)-1].Frame != uint32(cap+2-wantCount) {
				t.Fatal("pressure removed wrong end")
			}
			// Store boundary state only; each node and all list links are checked by snapshot.
			rows = append(rows, struct {
				Name   string
				Return uint32
				State  legacy.PortTestReliableReportState
			}{fmt.Sprintf("flags%d", flags), 1, state})
		})
	}
	for _, frame := range []uint32{0, 123, 999999998, 999999999, 1000000000, 0xffffffff} {
		o.reset()
		o.s.SetFrame(frame)
		if legacy.PortTestReliableReports(9, 0, 0, 0, nil, nil, 0, false) != 0 {
			t.Fatal("empty pressure")
		}
		for i := 0; i < 3; i++ {
			legacy.PortTestReliableReports(8, 31, 0, 0, []byte{77, byte(i)}, nil, 0, false)
		}
		rv := legacy.PortTestReliableReports(9, 0, 0, 0, nil, nil, 0, false)
		state := o.state()
		want := uint32(0)
		if frame < 999999999 {
			want = 1
		}
		if rv != want || len(state.Nodes) != 3-int(want) {
			t.Fatal("frame sentinel")
		}
		if want == 1 && state.Nodes[0].Data[1] != 1 {
			t.Fatal("equal timestamps choose newest")
		}
		rows = append(rows, struct {
			Name   string
			Return uint32
			State  legacy.PortTestReliableReportState
		}{fmt.Sprintf("frame%d", frame), rv, state})
	}
}

func TestReliableReportsRemovePlayer(t *testing.T) {
	o := newReliableReportsOwner(t)
	var rows []struct {
		Recipient      int
		Return, Status uint32
		State          legacy.PortTestReliableReportState
	}
	for _, recipient := range []int{0, 1, 7, 31} {
		o.reset()
		pl := o.s.Players.ByIndRaw(ntype.PlayerInd(recipient))
		objectXferSetWord(pl.C(), 3680, 0x400)
		for _, to := range []int{1, 7, 31, 129, 255} {
			legacy.PortTestReliableReports(8, to, 0, 0, []byte{77}, nil, 0, false)
		}
		rv := legacy.PortTestReliableReports(10, recipient, 0, 0, nil, nil, 0, false)
		state := o.state()
		status := objectXferGetWord(pl.C(), 3680)
		if recipient != 0 {
			if status != 0x480 {
				t.Fatal("timestamp status")
			}
			for _, n := range state.Nodes {
				if n.To == byte(recipient) {
					t.Fatal("recipient retained")
				}
			}
		} else if status != 0x400 {
			t.Fatal("inactive recipient changed")
		}
		rows = append(rows, struct {
			Recipient      int
			Return, Status uint32
			State          legacy.PortTestReliableReportState
		}{recipient, rv, status, state})
	}
	spellbookCapture(t, "reliable-reports-remove-player", rows, "178dc3acd210f86f43dd00c143037c47e0b42a22bd25bd0de671d02f168d9541")
}
