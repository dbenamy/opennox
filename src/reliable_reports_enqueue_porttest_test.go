//go:build porttest

package opennox

import (
	"bytes"
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
)

func TestReliableReportsEnqueue(t *testing.T) {
	o := newReliableReportsOwner(t)
	var rows []struct {
		Name   string
		Return uint32
		State  legacy.PortTestReliableReportState
	}
	defer func() {
		spellbookCapture(t, "reliable-reports-enqueue", rows, "c0f7d315ea2a262fc7b636f8b9f83b0b93bfa2c9a411a059b71bf786120b7686")
	}()
	for _, to := range []int{0, 1, 7, 31, 128, 129, 159, 255} {
		for _, mask := range []uint32{0, 2, 0x80000000, 0x80000082, 0xffffffff} {
			for _, size := range []int{0, 1, 149, 150, 151} {
				for _, ordered := range []bool{false, true} {
					for _, priority := range []int{0, 1} {
						for _, related := range []bool{false, true} {
							name := fmt.Sprintf("to%d/mask%x/size%d/order%t/priority%d/related%t", to, mask, size, ordered, priority, related)
							t.Run(name, func(t *testing.T) {
								o.reset()
								*o.words["mask"] = mask
								data := make([]byte, size)
								for i := range data {
									data[i] = byte(17*i + 73)
								}
								var obj *server.Object
								if related {
									obj = &o.units[0]
								}
								rv := legacy.PortTestReliableReports(8, to, 0, 0, data, obj, priority, ordered)
								state := o.state()
								eligible := to == 255 || to < 128 || mask&^(uint32(1)<<uint(to&127)) != 0
								wantReturn := uint32(1)
								if eligible && size > 150 {
									wantReturn = 0
								}
								count := 0
								if eligible && size <= 150 {
									count = 1
								}
								if rv != wantReturn || len(state.Nodes) != count || state.Pool != (count == 1) {
									t.Fatalf("return%d want%d nodes%d pool%t", rv, wantReturn, len(state.Nodes), state.Pool)
								}
								if count == 1 {
									n := state.Nodes[0]
									if n.Frame != 123 || n.To != byte(to) || n.Recipients != mask || n.Priority != uint32(priority) || !bytes.Equal(n.Data, data) || state.Capacity != 256 {
										t.Fatalf("node %+v capacity%d", n, state.Capacity)
									}
									wantRelated := uint32(0)
									if related {
										wantRelated = 1001
									}
									if n.Related != wantRelated || n.Ordered != byte(boolInt(ordered)) {
										t.Fatal("reference/order")
									}
									if n.Acknowledged != 0 || n.Sent != 0 || n.Retries != 0 || n.LastSent != [32]uint32{} || n.Countdown != [32]byte{} || n.Sequence != [32]uint16{} {
										t.Fatal("new node state")
									}
								}
								for i, v := range state.Sequence {
									want := uint16(0)
									if count == 1 && ordered && (to < 128 && i == to || to == 255 && mask&(uint32(1)<<uint(i)) != 0 || to >= 128 && to != 255 && i != to&127 && mask&(uint32(1)<<uint(i)) != 0) {
										want = 1
									}
									if v != want {
										t.Fatalf("sequence%d=%d want%d", i, v, want)
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
		}
	}
}
func boolInt(v bool) int {
	if v {
		return 1
	}
	return 0
}
