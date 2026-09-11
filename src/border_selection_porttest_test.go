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
	// the selected row's +44 u16 count.
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

func TestBorderSelectionByte4UsesSelectedRow(t *testing.T) {
	// The loader stores border count at row+44. Validate the selected primary
	// row, independently of the input variation's row. Both examples use plausible
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
	// Selected row7 allows variation16 (<20), despite input row16 limit12.
	if r := run(legacy.PortTestBorderState{Count: 64, Flag: 2, Primary: 7, Secondary: 0xa5}); r.Return != 1 || r.State.Secondary != 16 {
		t.Fatalf("selected-row accept: %+v", r)
	}
	// Selected row7 rejects variation16 (>=12), despite input row16 limit20.
	rows[0] = legacy.NewPortTestBorderRow(7, []byte("Selected"), 12)
	rows[1] = legacy.NewPortTestBorderRow(16, []byte("Variation"), 20)
	if r := run(legacy.PortTestBorderState{Count: 64, Flag: 0xff, Primary: 7, Secondary: 0xa5}); r.Return != 0 || r.State.Secondary != 0xa5 {
		t.Fatalf("selected-row reject: %+v", r)
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

// These boundary calls include inputs unsafe to evaluate with the old C lookup.
func TestBorderSelectionByte4RepairBoundaries(t *testing.T) {
	checks := 0
	for _, count := range []uint32{0, 1, 8, 64, 65, 0x7fffffff, 0x80000000, 0xffffffff} {
		for _, primary := range []uint32{0, 7, 63, 64, 255, 0x7fffffff, 0x80000000, 0xffffffff} {
			for _, flag := range []uint32{0, 1, 2, 0x80000000, 0xffffffff} {
				for _, limit := range []uint16{0, 1, 12, 20, 64, 65, 1020, 65535} {
					values := []int32{-2147483648, -1, 0, int32(limit) - 1, int32(limit), int32(limit) + 1, 63, 64, 65535, 2147483647}
					initial := legacy.PortTestBorderState{Count: count, Flag: flag, Primary: primary, Secondary: 0xa5a5a5a5}
					var rows []legacy.PortTestBorderRow
					if primary < 64 {
						rows = []legacy.PortTestBorderRow{legacy.NewPortTestBorderRow(int(primary), []byte("Selected"), limit)}
					}
					specs := make([]legacy.PortTestBorderSpec, len(values))
					for i, value := range values {
						specs[i] = legacy.PortTestBorderSpec{Mode: 3, Value: value}
					}
					snap := legacy.PortTestBorderSelection(initial, rows, specs)
					if snap.Before != snap.AfterRestore || !snap.TableRestored || !snap.GuardsRestored {
						t.Fatal("fixture restore")
					}
					wantState := initial
					for i, value := range values {
						wantRet := 0
						if flag == 0 {
							wantRet = 1
						} else if int32(count) > 0 && primary < count && primary < 64 && value >= 0 && uint32(value) < uint32(limit) {
							wantRet = 1
							wantState.Secondary = uint32(value)
						}
						r := snap.Results[i]
						if r.Return != wantRet || r.State != wantState || !r.TableUnchanged || !r.InputUnchanged || !r.GuardsUnchanged {
							t.Fatalf("count=%08x primary=%08x flag=%08x limit=%d value=%d: got %+v want return=%d state=%+v", count, primary, flag, limit, value, r, wantRet, wantState)
						}
						checks++
					}
				}
			}
		}
	}
	t.Logf("%d repaired validation boundary calls", checks)
}

func TestBorderSelectionLookupMatrix(t *testing.T) {
	checks := 0
	for count := uint32(0); count <= 64; count++ {
		for index := 0; index < 64; index++ {
			for _, name := range [][]byte{[]byte("Target"), []byte(""), {0xff, 0x80, 'x'}, []byte("Target\x00suffix")} {
				rows := make([]legacy.PortTestBorderRow, 64)
				for i := range rows {
					rows[i] = legacy.NewPortTestBorderRow(i, []byte("Other"), 12)
				}
				rows[index] = legacy.NewPortTestBorderRow(index, name, 20)
				if index+1 < 64 {
					rows[index+1] = legacy.NewPortTestBorderRow(index+1, name, 12)
				}
				initial := legacy.PortTestBorderState{Count: count, Flag: 0x80000000, Primary: 0xa5a5a5a5, Secondary: 0x5a5a5a5a}
				specs := []legacy.PortTestBorderSpec{{Mode: 0, Name: name}, {Mode: 1, Name: name}}
				snap := legacy.PortTestBorderSelection(initial, rows, specs)
				if snap.Before != snap.AfterRestore || !snap.TableRestored || !snap.GuardsRestored {
					t.Fatal("fixture restore")
				}
				lookup, accepted := -1, 0
				selected := initial
				selected.Flag = 0
				if uint32(index) < count {
					lookup, accepted = index, 1
					selected.Flag, selected.Primary = 1, uint32(index)
				}
				for i, r := range snap.Results {
					wantRet, wantState := lookup, initial
					if i == 1 {
						wantRet, wantState = accepted, selected
					}
					if r.Return != wantRet || r.State != wantState || !r.TableUnchanged || !r.InputUnchanged || !r.GuardsUnchanged {
						t.Fatalf("count=%d index=%d name=%x mode=%d got=%+v expected=%d %+v", count, index, name, i, r, wantRet, wantState)
					}
					checks++
				}
			}
		}
	}
	t.Logf("%d exact lookup/name selection checks", checks)
}

func TestBorderSelectionPrimaryMatrix(t *testing.T) {
	for _, count := range []uint32{0, 1, 64, 65, 0x7fffffff, 0x80000000, 0xffffffff} {
		initial := legacy.PortTestBorderState{Count: count, Flag: 2, Primary: 255, Secondary: 0xa5}
		for _, value := range []int32{-2147483648, -1, 0, 1, 63, 64, 65, 2147483646, 2147483647} {
			r := checkBorder(t, initial, []legacy.PortTestBorderSpec{{Mode: 2, Value: value}})[0]
			wantRet, want := 0, initial
			if value >= 0 && value < int32(count) {
				wantRet, want.Flag, want.Primary = 1, 1, uint32(value)
			}
			if r.Return != wantRet || r.State != want {
				t.Fatalf("count=%08x value=%d got=%+v want=%d %+v", count, value, r, wantRet, want)
			}
		}
	}
}
