//go:build porttest

package opennox

import (
	"bytes"
	"fmt"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestObjectReportsQueueCapacity(t *testing.T) {
	s, units, _ := objectReportsPlayers(t)
	a, u := &units[0], &units[1]
	var rows []struct {
		Name   string
		Return int
		Bytes  []byte
	}
	defer func() {
		spellbookCapture(t, "object-reports-queue-capacity", rows, "9da6e3a53b1454df07e8af94e7be74072f7c29e936dbdc9731209614dbcd2e61")
	}()
	for _, op := range []int{4, 5} {
		for _, occupied := range []int{0, 2035, 2036, 2037, 2044, 2045, 2048} {
			name := fmt.Sprintf("op%d/occupied%d", op, occupied)
			t.Run(name, func(t *testing.T) {
				s.NetList.ResetAll()
				u.TypeInd = 200
				u.NetCode = 123
				u.ObjClass = 4
				u.HealthData = nil
				u.UpdateDataPlayer().State = 0
				if op == 4 {
					u.ObjClass = 0x80
				}
				fill := bytes.Repeat([]byte{0xa5}, occupied)
				if !s.NetList.AddToMsgListCli(1, netlist.Kind1, fill) {
					t.Fatal("fill failed")
				}
				rv := legacy.PortTestObjectReports(op, a, u, 1, 0, 0, nil)
				got := s.NetList.CopyPacketsA(1, netlist.Kind1)
				size := 12
				if op == 4 {
					size = 4
				}
				want := 1
				if occupied+size > 2048 {
					want = 0
					size = 0
				}
				if rv != want || len(got) != occupied+size || !bytes.Equal(got[:occupied], fill) {
					t.Fatalf("return%d want%d size%d", rv, want, len(got))
				}
				rows = append(rows, struct {
					Name   string
					Return int
					Bytes  []byte
				}{name, rv, got})
			})
		}
	}
}
