//go:build porttest

package opennox

import (
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/opennox/noxscript/ns/asm"
	"github.com/opennox/opennox/v1/legacy"
)

func TestQuestProgressVariables(t *testing.T) {
	t.Cleanup(legacy.PortTestQuestProgressOwner())
	type row struct {
		Op, Name string
		Value    uint32
		Return   uint64
		Records  [][37]uint32
		Scratch  []byte
	}
	var rows []row
	run := func(op, name string, v uint32) uint64 {
		r := legacy.PortTestQuestProgress(op, name, v)
		rows = append(rows, row{op, name, v, r, legacy.PortTestQuestProgressSnapshot(), legacy.PortTestQuestProgressScratch()})
		return r
	}
	for _, ns := range []string{"", "MapA", "MAPa", "other", strings.Repeat("n", 63)} {
		run("reset", "*:*", 0)
		run("namespace", ns, 0)
		run("namespace-nil", "", 0)
		for i, name := range []string{"", "count", "COUNT", "prefix:Count", ":empty", "A:b:c", strings.Repeat("k", 63)} {
			if got := run("set-int", name, uint32(i+7)); i == 0 && got != 0 {
				t.Fatal("new first entry should return old nil head")
			}
			if got := run("int", name, 0); got != uint64(i+7) {
				t.Fatal("integer lookup", ns, name, got)
			}
			run("find", name, 0)
			for _, v := range []uint32{0, 1, 0x7fffffff, 0x80000000, 0xffffffff, 0x3f800000, 0x80000001, 0x7f800000, 0xff800000, 0x7fc01234} {
				run("set-float", name, v)
				if got := run("int", name, 0); got != uint64(v) {
					t.Fatal("bits changed")
				}
				run("float", name, 0)
			}
			run("qualify", name, 0)
		}
		if run("int", "missing", 0) != 0 || run("float", "missing", 0) != math.Float64bits(0) {
			t.Fatal("missing default")
		}
	}
	spellbookCapture(t, "quest-progress-variables", rows, "cb92b508cb4c1e0639c62f1e34db005d831d3f677e0a5a9a82d84f1f1154e97e")
}
func TestQuestProgressResetMasks(t *testing.T) {
	t.Cleanup(legacy.PortTestQuestProgressOwner())
	type row struct {
		Mask          string
		Before, After [][37]uint32
		Scratch       []byte
	}
	var rows []row
	names := []string{"a:one", "a:OneMore", "a:TWO", "b:one", "b:two", "a:prefixone", "c:ONE", "a:", "a:middleend", "a:repeatedendend"}
	for _, mask := range []string{"*:*", "a:*", "A:*", "*:one", "*:ONE", "a:*one", "a:*end", "a:one", "A:ONE", "one", "*", "missing", "a:o*e", "a:*missing", "b:*", "*:"} {
		legacy.PortTestQuestProgress("reset", "*:*", 0)
		legacy.PortTestQuestProgress("namespace", "a", 0)
		for i, n := range names {
			legacy.PortTestQuestProgress("set-int", n, uint32(i+1))
		}
		before := legacy.PortTestQuestProgressSnapshot()
		legacy.PortTestQuestProgress("reset", mask, 0)
		after := legacy.PortTestQuestProgressSnapshot()
		if mask == "*:*" && len(after) != 0 {
			t.Fatal("clear all")
		}
		if mask == "a:one" && len(after) != 9 {
			t.Fatal("exact removal")
		}
		rows = append(rows, row{mask, before, after, legacy.PortTestQuestProgressScratch()})
	}
	spellbookCapture(t, "quest-progress-reset", rows, "5232551592b4fd4b6f11622284abea445998560ac416d7f3a7fd4b0022a6a333")
}
func TestQuestProgressBuiltins(t *testing.T) {
	o := newCollisionCoreOwner(t)
	t.Cleanup(legacy.PortTestQuestProgressOwner())
	type row struct {
		Name       string
		Value, Got uint32
		Records    [][37]uint32
	}
	var rows []row
	legacy.PortTestQuestProgress("namespace", "scriptmap", 0)
	call := func(fi asm.Builtin) {
		t.Helper()
		r, ok := legacy.CallScriptBuiltin(fi)
		if !ok || r != 0 {
			t.Fatal("registered builtin", fi, r, ok)
		}
	}
	for _, v := range []uint32{0, 1, 0xffffffff, 0x80000000, 0x3fc00000, 0x80000001, 0x7fc05678} {
		for _, floating := range []bool{false, true} {
			name := fmt.Sprintf("key%08x%v", v, floating)
			o.s.NoxScriptVM.PushU32(v)
			o.s.NoxScriptVM.PushString(name)
			set, get := asm.BuiltinSetQuestStatus, asm.BuiltinGetQuestStatus
			if floating {
				set, get = asm.BuiltinSetQuestStatusFloat, asm.BuiltinGetQuestStatusFloat
			}
			call(set)
			o.s.NoxScriptVM.PushString(name)
			call(get)
			got := o.s.NoxScriptVM.PopU32()
			if got != v {
				t.Fatalf("builtin value %08x != %08x", got, v)
			}
			rows = append(rows, row{name, v, got, legacy.PortTestQuestProgressSnapshot()})
		}
	}
	o.s.NoxScriptVM.PushString("*:*")
	call(asm.BuiltinResetQuestStatus)
	if len(legacy.PortTestQuestProgressSnapshot()) != 0 {
		t.Fatal("builtin reset")
	}
	spellbookCapture(t, "quest-progress-builtins", rows, "ba2fc34bdd1c3977d125b8123b09984caf7b3db582aad39f09ca372750c8d669")
}

func TestQuestProgressNameBoundary(t *testing.T) {
	t.Cleanup(legacy.PortTestQuestProgressOwner())
	var rows [][37]uint32
	for _, name := range []string{":" + strings.Repeat("x", 129), ":" + strings.Repeat("y", 130)} {
		legacy.PortTestQuestProgress("set-int", name, 0x13579bdf)
		if legacy.PortTestQuestProgress("int", name, 0) != 0x13579bdf {
			t.Fatal("maximum defined name lost")
		}
		snap := legacy.PortTestQuestProgressSnapshot()
		if snap[0][33] != 0 || snap[0][34] != 0x13579bdf {
			t.Fatal("name overwrote adjacent kind/value")
		}
		rows = append(rows, snap[0])
	}
	spellbookCapture(t, "quest-progress-name-boundary", rows, "ed1af364fabc04b59099f3241b633564bed485dc722d360fdeca39cde522af7e")
}
