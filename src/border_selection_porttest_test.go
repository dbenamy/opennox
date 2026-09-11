//go:build porttest

package opennox

import (
	"reflect"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func borderI32(bits uint32) int32 { return int32(bits) }
func borderRows() []legacy.PortTestBorderRow {
	return []legacy.PortTestBorderRow{
		legacy.NewPortTestBorderRow(0, []byte("Alpha"), 2),
		legacy.NewPortTestBorderRow(1, []byte("alpha"), 12),
		legacy.NewPortTestBorderRow(2, []byte("Beta"), 20),
		legacy.NewPortTestBorderRow(3, []byte("Alpha"), 4), // first exact duplicate wins
		legacy.NewPortTestBorderRow(7, []byte("Selected"), 20),
		legacy.NewPortTestBorderRow(16, []byte("Variation"), 12),
	}
}
func checkBorder(t *testing.T, initial legacy.PortTestBorderState, specs []legacy.PortTestBorderSpec) []legacy.PortTestBorderResult {
	t.Helper()
	snap := legacy.PortTestBorderSelection(initial, borderRows(), specs)
	if !reflect.DeepEqual(snap.Before, snap.AfterRestore) || len(snap.Results) != len(specs) || !snap.TableRestored || !snap.GuardsRestored {
		t.Fatalf("fixture restore failed: %+v", snap)
	}
	for i, r := range snap.Results {
		if !r.TableUnchanged || !r.InputUnchanged || !r.GuardsUnchanged {
			t.Fatalf("case %d mutated input/table: %+v", i, r)
		}
	}
	return snap.Results
}

func TestBorderSelectionCABI(t *testing.T) {
	initial := legacy.PortTestBorderState{Count: 4, Flag: 0x7f, Primary: 0xaaaa5555, Secondary: 0x11223344}
	got := checkBorder(t, initial, []legacy.PortTestBorderSpec{
		{Mode: 0, NilName: true}, {Mode: 0, Name: []byte("Alpha")}, {Mode: 0, Name: []byte("ALPHA")}, {Mode: 0, Name: []byte("missing")},
		{Mode: 1, Name: []byte("Alpha")}, {Mode: 1, Name: []byte("NONE")}, {Mode: 1, Name: []byte("none")}, {Mode: 1, Name: []byte("missing")},
	})
	wantRet := []int{-1, 0, -1, -1, 1, 1, 1, 0}
	for i, r := range got {
		if r.Return != wantRet[i] {
			t.Fatalf("lookup/select %d got %+v", i, r)
		}
	}
	// Exact lookup is case-sensitive; select's NONE check alone is insensitive.
	if got[4].State.Flag != 1 || got[4].State.Primary != 0 || got[5].State.Flag != 0 || got[5].State.Primary != 255 || got[6].State.Primary != 255 || got[7].State.Primary != 255 || got[7].State.Flag != 0 {
		t.Fatalf("selection state: %+v", got)
	}

	got = checkBorder(t, legacy.PortTestBorderState{Count: 4, Flag: 9, Primary: 0xdeadbeef, Secondary: 7}, []legacy.PortTestBorderSpec{
		{Mode: 2, Value: -1}, {Mode: 2, Value: 0}, {Mode: 2, Value: 3}, {Mode: 2, Value: 4}, {Mode: 2, Value: borderI32(0x80000000)},
	})
	wantRet = []int{0, 1, 1, 0, 0}
	for i, r := range got {
		if r.Return != wantRet[i] {
			t.Fatalf("byte3 %d %+v", i, r)
		}
	}
	if got[1].State.Flag != 1 || got[1].State.Primary != 0 || got[2].State.Primary != 3 || got[3].State.Primary != 3 || got[4].State.Primary != 3 {
		t.Fatalf("byte3 state: %+v", got)
	}

	// Any nonzero flag, rather than only 1, enables Byte4. Its comparison uses
	// the input row's +44 u16 count, as the original C currently does.
	got = checkBorder(t, legacy.PortTestBorderState{Count: 64, Flag: 0x80, Primary: 7, Secondary: 0x55}, []legacy.PortTestBorderSpec{
		{Mode: 3, Value: -1}, {Mode: 3, Value: 0}, {Mode: 3, Value: 1}, {Mode: 3, Value: 2},
	})
	wantRet = []int{0, 1, 1, 1}
	for i, r := range got {
		if r.Return != wantRet[i] {
			t.Fatalf("byte4 %d %+v", i, r)
		}
	}
	short := checkBorder(t, legacy.PortTestBorderState{Count: 64, Flag: 0, Primary: 0xffffffff, Secondary: 0xa5}, []legacy.PortTestBorderSpec{{Mode: 3, Value: borderI32(0x80000000)}})[0]
	if short.Return != 1 || short.State.Secondary != 0xa5 {
		t.Fatalf("byte4 zero-flag short-circuit: %+v", short)
	}
}

func TestBorderSelectionByte4WrongRowOriginalC(t *testing.T) {
	// The loader stores border count at row+44. C currently checks the input
	// variation's row, not selected primary row. Both examples use plausible
	// even counts 20 (5+5) and 12 (3+3).
	rows := []legacy.PortTestBorderRow{
		legacy.NewPortTestBorderRow(7, []byte("Selected"), 20),
		legacy.NewPortTestBorderRow(16, []byte("Variation"), 12),
	}
	run := func(initial legacy.PortTestBorderState) legacy.PortTestBorderResult {
		s := legacy.PortTestBorderSelection(initial, rows, []legacy.PortTestBorderSpec{{Mode: 2, Value: 7}, {Mode: 3, Value: 16}})
		if !reflect.DeepEqual(s.Before, s.AfterRestore) || !s.TableRestored || !s.GuardsRestored {
			t.Fatal("fixture restore")
		}
		if s.Results[0].Return != 1 || s.Results[0].State.Primary != 7 || s.Results[0].State.Flag != 1 {
			t.Fatalf("primary selection failed: %+v", s.Results[0])
		}
		if !s.Results[1].TableUnchanged || !s.Results[1].GuardsUnchanged || !s.Results[1].InputUnchanged {
			t.Fatalf("wrong-row probe mutated table: %+v", s.Results[0])
		}
		return s.Results[1]
	}
	// Selected row7 allows variation16 (<20), but input row16 limit12 rejects.
	if r := run(legacy.PortTestBorderState{Count: 64, Flag: 2, Primary: 7, Secondary: 0xa5}); r.Return != 0 || r.State.Secondary != 0xa5 {
		t.Fatalf("wrong-row reject changed: %+v", r)
	}
	// Selected row7 rejects variation16 (>=12), but input row16 limit20 accepts.
	rows[0] = legacy.NewPortTestBorderRow(7, []byte("Selected"), 12)
	rows[1] = legacy.NewPortTestBorderRow(16, []byte("Variation"), 20)
	if r := run(legacy.PortTestBorderState{Count: 64, Flag: 0xff, Primary: 7, Secondary: 0xa5}); r.Return != 1 || r.State.Secondary != 16 {
		t.Fatalf("wrong-row accept changed: %+v", r)
	}
}

func TestBorderSelectionSignedCount(t *testing.T) {
	for _, count := range []uint32{0, 1, 4, 64, 0x80000000, 0xffffffff} {
		initial := legacy.PortTestBorderState{Count: count, Flag: 9, Primary: 123, Secondary: 456}
		got := checkBorder(t, initial, []legacy.PortTestBorderSpec{{Mode: 0, Name: []byte("Alpha")}, {Mode: 2, Value: 0}})
		lookup, selectRet := -1, 0
		selected := initial
		if int32(count) > 0 {
			lookup, selectRet = 0, 1
			selected.Flag, selected.Primary = 1, 0
		}
		if got[0].Return != lookup || got[0].State != initial || got[1].Return != selectRet || got[1].State != selected {
			t.Fatalf("count %08x: %+v", count, got)
		}
	}
}
