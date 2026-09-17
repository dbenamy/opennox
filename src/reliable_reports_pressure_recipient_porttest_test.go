//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestReliableReportsPressureRecipient(t *testing.T) {
	o := newReliableReportsOwner(t)
	// The chosen oldest message belongs to the slow player. Removing that player
	// already releases this node; pressure cleanup must not release it a second time.
	for i := 0; i < 3; i++ {
		o.s.SetFrame(uint32(100 + i))
		legacy.PortTestReliableReports(8, 1, 0, 0, []byte{77, byte(i)}, nil, 0, false)
	}
	if rv := legacy.PortTestReliableReports(9, 0, 0, 0, nil, nil, 0, false); rv != 1 {
		t.Fatalf("pressure return%d", rv)
	}
	if len(o.state().Nodes) != 0 {
		t.Fatal("slow player queue retained")
	}
	if objectXferGetWord(o.s.Players.ByInd(1).C(), 3680)&128 == 0 {
		t.Fatal("slow player status missing")
	}
	// The pool remains reusable, including after allocating beyond its old occupancy.
	for i := 0; i < 4; i++ {
		if legacy.PortTestReliableReports(8, 31, 0, 0, []byte{78, byte(i)}, nil, 0, false) != 1 {
			t.Fatal("pool reuse")
		}
	}
	state := o.state()
	if len(state.Nodes) != 4 {
		t.Fatal("queue after reuse")
	}
	spellbookCapture(t, "reliable-reports-pressure-recipient", []legacy.PortTestReliableReportState{state}, "0d3081e3f1bbb257902ebe8e0825310253259ac47f873e8470c3e64e8d30971d")
}

func TestReliableReportsPressureSurvivors(t *testing.T) {
	o := newReliableReportsOwner(t)
	var rows []legacy.PortTestReliableReportState
	for _, oldestTo := range []int{1, 7, 31, 255} {
		o.reset()
		o.s.SetFrame(100)
		legacy.PortTestReliableReports(8, oldestTo, 0, 0, []byte{49}, &o.units[0], 0, false)
		objectXferSetWord(o.units[0].CObj(), 148, 0xffffffff)
		for i := 0; i < 3; i++ {
			o.s.SetFrame(uint32(101 + i))
			legacy.PortTestReliableReports(8, 1, 0, 0, []byte{77, byte(i)}, nil, 0, false)
		}
		if legacy.PortTestReliableReports(9, 0, 0, 0, nil, nil, 0, false) != 1 {
			t.Fatal("pressure")
		}
		state := o.state()
		if len(state.Nodes) != 0 {
			t.Fatal("oldest or slow recipient retained")
		}
		rows = append(rows, state)
		if oldestTo == 1 || oldestTo == 255 {
			if objectXferGetWord(o.units[0].CObj(), 148)&2 != 0 {
				t.Fatal("related known bit not cleared")
			}
		}
	}
	spellbookCapture(t, "reliable-reports-pressure-survivors", rows, "33baf72e33950f23e3fe876d249ff5170b925cd49a3f1e3b6834b85f49ca726d")
}
