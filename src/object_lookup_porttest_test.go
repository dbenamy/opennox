//go:build porttest

package opennox

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func lookupBase(op int, args ...legacy.PortTestGameplayReportArg) legacy.PortTestRoamSpec {
	s := gameplayReportsBase(200+op, args...)
	textSpec(&s).Lookup = &legacy.PortTestObjectLookupSpec{}
	return s
}
func lookupSpec(s *legacy.PortTestRoamSpec) *legacy.PortTestObjectLookupSpec {
	return textSpec(s).Lookup
}
func lookupObject(ref int) legacy.PortTestGameplayReportArg {
	return legacy.PortTestGameplayReportArg{Kind: "lookup-object", Ref: ref}
}
func lookupName(ref int) legacy.PortTestGameplayReportArg {
	return legacy.PortTestGameplayReportArg{Kind: "lookup-name", Ref: ref}
}
func lookupNode(ref int) legacy.PortTestGameplayReportArg {
	return legacy.PortTestGameplayReportArg{Kind: "lookup-node", Ref: ref}
}
func lookupHead(ref int) legacy.PortTestGameplayReportArg {
	return legacy.PortTestGameplayReportArg{Kind: "lookup-head", Ref: ref}
}
func TestObjectLookupProbe(t *testing.T) {
	s := lookupBase(3, reportValue(123))
	sp := lookupSpec(&s)
	sp.Objects = []legacy.PortTestLookupObject{{Code: 123}}
	sp.Main = []int{1}
	out := controlsRun(t, []legacy.PortTestRoamSpec{s})[0]
	step := out.Callbacks.Shop.Sequence[0]
	if step.Return != 940001 || step.GameplayReports.Lookup == nil {
		t.Fatalf("lookup result %d", step.Return)
	}
	c := step.GameplayReports.Lookup
	if c.Heads != [4]uint32{920014, 920000, 920015, 920015} || c.Nodes[15][0] != 940001 {
		t.Fatalf("cache head/value state %+v", c)
	}
	if !out.Intact || !out.Callbacks.Shop.Intact {
		t.Fatal("lookup guards")
	}
}

var objectLookupHashes = map[string]string{
	"names":            "3177ab533d954f26df0f5ffcfae749b962c03d4b7f95a25060dd3d8aea2d6fc9",
	"name-order":       "48573d73b3e4d0d459a30e2a73b326d3e83656ce63c3a3466a1f517f466cbaa1",
	"search-locations": "a9f8b695ef992058dcf8f90354281530633930c97445457b9e7ebc92a2880c92",
	"cache-lifecycle":  "c02d4f321c45c9aaaef43aadf581680c4cb4f6b4700c98e3cdbcd7ecdf33b0d4",
	"cached-mutation":  "230d18646a696653468f69ca7558fd0a76bde1332f379e116c80098ac905254c",
	"cache-lists":      "26f7e8a18777c721f0a8f8938aa23d74098aac753c255b2647130fb023fdadb8",
	"cold-duplicates":  "91bb62cd18e75dcc9f9c9a17a7b3c230d1ac1602765b1e5855aa5e61b8e28b99",
	"cache-mixed":      "4ac78d7730eb17919bdb7399570c58baadbac3b988a4132f6611e9f5d706c297",
}

func lookupCapture(t *testing.T, label string, out []legacy.PortTestRoamResult) {
	t.Helper()
	for i, r := range out {
		if !r.Intact || !r.Callbacks.Intact || !r.Spells.Intact || !r.Combat.Intact || !r.MonsterState.Intact || !r.Callbacks.Shop.Intact {
			t.Fatalf("lookup %s case%d guards", label, i)
		}
	}
	data, err := json.Marshal(out)
	if err != nil {
		t.Fatal(err)
	}
	if prefix := os.Getenv("OPENNOX_OBJECT_LOOKUP_CAPTURE"); prefix != "" {
		if err := os.WriteFile(prefix+"-"+label+".json", data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	hash := fmt.Sprintf("%x", sha256.Sum256(data))
	t.Logf("%s: %d cases %s", label, len(out), hash)
	want, ok := objectLookupHashes[label]
	if !ok {
		t.Fatalf("missing audited C capture: %s", label)
	}
	if hash != want {
		t.Fatalf("%s hash %s want%s", label, hash, want)
	}
}
func lookupString(s string) *string { return &s }

type lookupCommand struct {
	op   int
	args []legacy.PortTestGameplayReportArg
}

func lookupCmd(op int, args ...legacy.PortTestGameplayReportArg) lookupCommand {
	return lookupCommand{op, args}
}
func lookupCommands(s *legacy.PortTestRoamSpec, cmds ...lookupCommand) {
	s.Callbacks.Shop.Sequence = nil
	textSpec(s).Calls = make(map[uint32][5]legacy.PortTestGameplayReportArg)
	for i, c := range cmds {
		var args [5]legacy.PortTestGameplayReportArg
		copy(args[:], c.args)
		textSpec(s).Calls[uint32(i)] = args
		s.Callbacks.Shop.Sequence = append(s.Callbacks.Shop.Sequence, legacy.PortTestShopAction{Op: 2000 + c.op, Value: uint32(i)})
	}
}
func lookupUsed(t *testing.T, c *legacy.PortTestObjectLookupResult) []uint32 {
	t.Helper()
	if c == nil {
		t.Fatal("missing cache snapshot")
	}
	var out []uint32
	seen := map[uint32]bool{}
	for id := c.Heads[2]; id != 0; {
		if id < 920000 || id >= 920016 || seen[id] {
			t.Fatalf("cache used-list link %d", id)
		}
		seen[id] = true
		n := c.Nodes[id-920000]
		out = append(out, n[0])
		id = n[2]
	}
	return out
}
func TestObjectLookupNames(t *testing.T) {
	rows := []struct {
		name  *string
		query string
		match [3]bool
	}{
		{nil, "", [3]bool{}}, {lookupString(""), "", [3]bool{true, true, true}},
		{lookupString(""), "a", [3]bool{}}, {lookupString("a"), "a", [3]bool{true, true, true}},
		{lookupString("a"), "A", [3]bool{}}, {lookupString("x:a"), "a", [3]bool{true, false, true}},
		{lookupString("x:a"), "x:a", [3]bool{true, true, false}},
		{lookupString("x:a:b"), "a:b", [3]bool{false, false, true}},
		{lookupString("x:a:b"), "b", [3]bool{}}, {lookupString("x:a:b"), "x:a:b", [3]bool{true, true, false}},
		{lookupString(":a"), "a", [3]bool{true, false, true}}, {lookupString(":a"), ":a", [3]bool{true, true, false}},
		{lookupString("x:"), "", [3]bool{true, false, true}}, {lookupString("x::"), "", [3]bool{}},
		{lookupString("x::"), ":", [3]bool{false, false, true}}, {lookupString("x::"), "x::", [3]bool{true, true, false}},
		{lookupString("a\x00z"), "a", [3]bool{true, true, true}}, {lookupString("a"), "a\x00z", [3]bool{true, true, true}},
		{lookupString("\x00a"), "", [3]bool{true, true, true}}, {lookupString("\xff"), "\xff", [3]bool{true, true, true}},
		{lookupString("x:\xff"), "\xff", [3]bool{true, false, true}}, {lookupString("x:é"), "é", [3]bool{true, false, true}},
		{lookupString(" name "), " name ", [3]bool{true, true, true}}, {lookupString("name"), "name ", [3]bool{}},
	}
	var cases []legacy.PortTestRoamSpec
	var wants []uint32
	for _, r := range rows {
		for op := 0; op < 3; op++ {
			for location := 0; location < 4; location++ {
				for _, flags := range []uint32{0, 0x20, 0xffffffff} {
					s := lookupBase(op, lookupName(0))
					if op != 0 {
						textSpec(&s).Args = [5]legacy.PortTestGameplayReportArg{lookupObject(1), lookupName(0)}
					}
					sp := lookupSpec(&s)
					sp.Objects = []legacy.PortTestLookupObject{{Flags: flags}, {Flags: flags}}
					sp.Names = []string{r.query}
					target := 1
					if location%2 != 0 {
						target = 2
						sp.Objects[0].Inventory = []int{2}
					}
					sp.Objects[target-1].Name = r.name
					if location < 2 {
						sp.Main = []int{1}
					} else {
						sp.Pending = []int{1}
					}
					want := uint32(0)
					if r.match[op] {
						want = 940000 + uint32(target)
					}
					cases = append(cases, s)
					wants = append(wants, want)
				}
			}
		}
	}
	out := controlsRun(t, cases)
	for i, r := range out {
		if v := r.Callbacks.Shop.Sequence[0].Return; v != wants[i] {
			t.Fatalf("name case%d return%d want%d", i, v, wants[i])
		}
	}
	lookupCapture(t, "names", out)
}
func TestObjectLookupNameOrder(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	var wants []uint32
	for winner := 1; winner <= 6; winner++ {
		s := lookupBase(0, lookupName(0))
		sp := lookupSpec(&s)
		sp.Names = []string{"target"}
		sp.Objects = make([]legacy.PortTestLookupObject, 6)
		sp.Main = []int{1, 3}
		sp.Pending = []int{4}
		sp.Objects[0].Inventory = []int{2}
		sp.Objects[3].Inventory = []int{5}
		sp.Objects[1].Inventory = []int{6}
		for i := winner - 1; i < 6; i++ {
			sp.Objects[i].Name = lookupString("prefix:target")
			sp.Objects[i].Flags = 0x20
		}
		want := uint32(0)
		if winner <= 5 {
			want = 940000 + uint32(winner)
		}
		cases = append(cases, s)
		wants = append(wants, want)
	}
	out := controlsRun(t, cases)
	for i, r := range out {
		if r.Callbacks.Shop.Sequence[0].Return != wants[i] {
			t.Fatalf("name precedence%d", i)
		}
	}
	lookupCapture(t, "name-order", out)
}
func TestObjectLookupSearchLocations(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	var wants []uint32
	var cached []bool
	for _, op := range []int{3, 10} {
		for loc := 0; loc <= 8; loc++ {
			for _, deleted := range []bool{false, true} {
				for _, value := range []uint32{0, 1, 0x7fffffff, 0x80000000, 0xffffffff} {
					s := lookupBase(op, reportValue(value))
					sp := lookupSpec(&s)
					sp.Objects = make([]legacy.PortTestLookupObject, 7)
					for i := range sp.Objects {
						sp.Objects[i].Code = uint32(10000 + i)
						sp.Objects[i].ScriptID = uint32(20000 + i)
					}
					sp.Main = []int{1}
					sp.Pending = []int{4}
					sp.Missiles = []int{6}
					sp.Objects[0].Inventory = []int{2, 3}
					sp.Objects[1].Inventory = []int{7}
					sp.Objects[3].Inventory = []int{5}
					flags := uint32(0)
					if deleted {
						flags = 0x20
					}
					if loc >= 1 && loc <= 7 {
						if op == 3 {
							sp.Objects[loc-1].Code = value
						} else {
							sp.Objects[loc-1].ScriptID = value
						}
						sp.Objects[loc-1].Flags = flags
					}
					o := s.Callbacks.Shop.TemporaryUpdates.World.Objectives
					o.PlayerWords[0] = map[int]uint32{36: 9000, 44: 9001}
					if loc == 8 {
						o.PlayerWords[0] = map[int]uint32{36: value, 44: value, 16: flags}
					}
					want := uint32(0)
					cache := false
					if !deleted {
						if loc >= 1 && loc <= 4 {
							want = 940000 + uint32(loc)
							cache = op == 3
						} else if loc == 6 && op == 10 {
							want = 940006
						} else if loc == 8 && op == 3 {
							want = 54000 // The resources fixture identifies player zero last.
						}
					}
					cases = append(cases, s)
					wants = append(wants, want)
					cached = append(cached, cache)
				}
			}
		}
	}
	out := controlsRun(t, cases)
	for i, r := range out {
		st := r.Callbacks.Shop.Sequence[0]
		if st.Return != wants[i] {
			t.Fatalf("search case%d return%d want%d", i, st.Return, wants[i])
		}
		used := lookupUsed(t, st.GameplayReports.Lookup)
		want := []uint32(nil)
		if cached[i] {
			want = []uint32{wants[i]}
		}
		if !slices.Equal(used, want) {
			t.Fatalf("search case%d cache %v want%v", i, used, want)
		}
	}
	lookupCapture(t, "search-locations", out)
}
func TestObjectLookupCacheLifecycle(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	var returns [][]uint32
	var states [][][]uint32
	for _, count := range []int{1, 15, 16, 17, 20} {
		for _, cold := range []bool{false, true} {
			s := lookupBase(3, reportValue(1))
			sp := lookupSpec(&s)
			sp.Objects = make([]legacy.PortTestLookupObject, 20)
			sp.Uninitialized = cold
			for i := range sp.Objects {
				sp.Objects[i].Code = uint32(i + 1)
				sp.Main = append(sp.Main, i+1)
			}
			var cmds []lookupCommand
			var expected []uint32
			var order []uint32
			var orders [][]uint32
			add := func(c lookupCommand, ret uint32) {
				cmds = append(cmds, c)
				expected = append(expected, ret)
				orders = append(orders, slices.Clone(order))
			}
			touch := func(id uint32) {
				if i := slices.Index(order, id); i >= 0 {
					order = slices.Delete(order, i, i+1)
				} else if len(order) == 16 {
					order = order[:15]
				}
				order = append([]uint32{id}, order...)
			}
			for i := 1; i <= count; i++ {
				id := 940000 + uint32(i)
				touch(id)
				add(lookupCmd(3, reportValue(uint32(i))), id)
			}
			for _, v := range []uint32{1, 5, 20, 20} {
				id := 940000 + v
				touch(id)
				add(lookupCmd(3, reportValue(v)), id)
			}
			add(lookupCmd(3, reportValue(999)), 0)
			add(lookupCmd(4, reportValue(999)), 0)
			for _, v := range []int{1, 20, 7, 3} {
				id := 940000 + uint32(v)
				if i := slices.Index(order, id); i >= 0 {
					order = slices.Delete(order, i, i+1)
				}
				add(lookupCmd(11, lookupObject(v)), 0)
			}
			order = nil
			add(lookupCmd(12), 0)
			add(lookupCmd(12), 0)
			touch(940009)
			add(lookupCmd(3, reportValue(9)), 940009)
			order = nil
			add(lookupCmd(7), 0)
			add(lookupCmd(4, reportValue(9)), 0)
			lookupCommands(&s, cmds...)
			cases = append(cases, s)
			returns = append(returns, expected)
			states = append(states, orders)
		}
	}
	out := controlsRun(t, cases)
	for i, r := range out {
		for j, st := range r.Callbacks.Shop.Sequence {
			if st.Return != returns[i][j] {
				t.Fatalf("cache case%d step%d return", i, j)
			}
			if got := lookupUsed(t, st.GameplayReports.Lookup); !slices.Equal(got, states[i][j]) {
				t.Fatalf("cache case%d step%d order %v want%v", i, j, got, states[i][j])
			}
		}
	}
	lookupCapture(t, "cache-lifecycle", out)
}
func TestObjectLookupCachedMutation(t *testing.T) {
	s := lookupBase(3, reportValue(7))
	sp := lookupSpec(&s)
	sp.Objects = []legacy.PortTestLookupObject{{Code: 7}, {Code: 7}}
	sp.Main = []int{1, 2}
	lookupCommands(&s,
		lookupCmd(3, reportValue(7)), lookupCmd(13, lookupObject(1), reportValue(16), reportValue(0x20)), lookupCmd(3, reportValue(7)),
		lookupCmd(11, lookupObject(1)), lookupCmd(3, reportValue(7)), lookupCmd(13, lookupObject(2), reportValue(36), reportValue(9)),
		lookupCmd(3, reportValue(7)), lookupCmd(3, reportValue(9)), lookupCmd(13, lookupObject(2), reportValue(16), reportValue(0x20)),
		lookupCmd(3, reportValue(9)), lookupCmd(12), lookupCmd(3, reportValue(9)), lookupCmd(13, lookupObject(2), reportValue(16), reportValue(0)), lookupCmd(3, reportValue(9)))
	out := controlsRun(t, []legacy.PortTestRoamSpec{s})
	want := []uint32{940001, 0, 940001, 0, 940002, 0, 0, 940002, 0, 940002, 0, 0, 0, 940002}
	for i, st := range out[0].Callbacks.Shop.Sequence {
		if st.Return != want[i] {
			t.Fatalf("mutation step%d return%d want%d", i, st.Return, want[i])
		}
	}
	lookupCapture(t, "cached-mutation", out)
}
func TestObjectLookupCacheLists(t *testing.T) {
	s := lookupBase(9)
	lookupCommands(&s,
		lookupCmd(9), lookupCmd(5, lookupHead(0), lookupNode(15)), lookupCmd(9), lookupCmd(5, lookupHead(0), lookupNode(14)), lookupCmd(9), lookupCmd(5, lookupHead(0), lookupNode(13)),
		lookupCmd(6, lookupHead(0), lookupNode(14)), lookupCmd(5, lookupHead(1), lookupNode(14)), lookupCmd(6, lookupHead(0), lookupNode(13)), lookupCmd(5, lookupHead(1), lookupNode(13)), lookupCmd(6, lookupHead(0), lookupNode(15)), lookupCmd(5, lookupHead(1), lookupNode(15)), lookupCmd(7))
	out := controlsRun(t, []legacy.PortTestRoamSpec{s})
	for i, st := range out[0].Callbacks.Shop.Sequence {
		want := uint32(0)
		if i == 0 {
			want = 920015
		} else if i == 2 {
			want = 920014
		} else if i == 4 {
			want = 920013
		}
		if st.Return != want {
			t.Fatalf("list step%d", i)
		}
	}
	lookupCapture(t, "cache-lists", out)
}
func TestObjectLookupColdAndDuplicates(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, flag := range []uint32{1, 2, 0x80000000, 0xffffffff} {
		s := lookupBase(12)
		sp := lookupSpec(&s)
		sp.Objects = []legacy.PortTestLookupObject{{Code: 7}}
		sp.Uninitialized = true
		sp.InitFlag = flag
		sp.CacheValues = []int{1, 1, 1}
		lookupCommands(&s, lookupCmd(12), lookupCmd(11, lookupObject(1)), lookupCmd(7), lookupCmd(8, lookupObject(1)), lookupCmd(8, lookupObject(1)), lookupCmd(4, reportValue(7)), lookupCmd(11, lookupObject(1)), lookupCmd(4, reportValue(7)), lookupCmd(7), lookupCmd(4, reportValue(7)))
		cases = append(cases, s)
	}
	out := controlsRun(t, cases)
	for i, r := range out {
		for j, st := range r.Callbacks.Shop.Sequence {
			if j < 2 && st.GameplayReports.Lookup.NeedInit != []uint32{1, 2, 0x80000000, 0xffffffff}[i] {
				t.Fatal("cold cache mutated before init")
			}
			want := uint32(0)
			if j == 5 || j == 7 {
				want = 940001
			}
			if st.Return != want {
				t.Fatalf("cold/duplicate case%d step%d", i, j)
			}
		}
	}
	lookupCapture(t, "cold-duplicates", out)
}

// The oracle uses an ordered slice, independently of the cache's linked nodes.
func TestObjectLookupCacheMixed(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	var returns [][]uint32
	var states [][][]uint32
	for seed := uint32(1); seed <= 8; seed++ {
		s := lookupBase(3, reportValue(1))
		sp := lookupSpec(&s)
		sp.Objects = make([]legacy.PortTestLookupObject, 24)
		for i := range sp.Objects {
			sp.Objects[i].Code = uint32(i + 1)
			sp.Main = append(sp.Main, i+1)
		}
		var cmds []lookupCommand
		var expected []uint32
		var order []uint32
		var orders [][]uint32
		rng := seed
		for step := 0; step < 160; step++ {
			rng = rng*1664525 + 1013904223
			code := (rng>>8)%28 + 1
			id := 940000 + code
			op := 3
			ret := uint32(0)
			var args []legacy.PortTestGameplayReportArg
			switch (rng >> 24) % 16 {
			case 0:
				op = 12
				order = nil
			case 1:
				op = 7
				order = nil
			case 2, 3:
				op = 11
				code = (code-1)%24 + 1
				id = 940000 + code
				args = []legacy.PortTestGameplayReportArg{lookupObject(int(code))}
				if i := slices.Index(order, id); i >= 0 {
					order = slices.Delete(order, i, i+1)
				}
			default:
				if rng&1 != 0 {
					op = 4
				}
				args = []legacy.PortTestGameplayReportArg{reportValue(code)}
				i := slices.Index(order, id)
				if i >= 0 || (op == 3 && code <= 24) {
					ret = id
					if i >= 0 {
						order = slices.Delete(order, i, i+1)
					} else if len(order) == 16 {
						order = order[:15]
					}
					order = append([]uint32{id}, order...)
				}
			}
			cmds = append(cmds, lookupCmd(op, args...))
			expected = append(expected, ret)
			orders = append(orders, slices.Clone(order))
		}
		lookupCommands(&s, cmds...)
		cases = append(cases, s)
		returns = append(returns, expected)
		states = append(states, orders)
	}
	out := controlsRun(t, cases)
	for i, r := range out {
		for j, st := range r.Callbacks.Shop.Sequence {
			if st.Return != returns[i][j] {
				t.Fatalf("mixed case%d step%d return%d want%d", i, j, st.Return, returns[i][j])
			}
			if got := lookupUsed(t, st.GameplayReports.Lookup); !slices.Equal(got, states[i][j]) {
				t.Fatalf("mixed case%d step%d order%v want%v", i, j, got, states[i][j])
			}
		}
	}
	lookupCapture(t, "cache-mixed", out)
}
