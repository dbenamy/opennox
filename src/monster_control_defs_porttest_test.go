//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/internal/binfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func monsterControlDefsOwner(t *testing.T) *collisionCoreOwner {
	t.Helper()
	o := newCollisionCoreOwner(t)
	t.Cleanup(handles.PortTestInit())
	t.Cleanup(legacy.PortTestMonsterDefsOwner())
	data, relocs := blobdata.PortTestMonsterDefinitionTables()
	copy(serverConfigOwnBytes(t, 0x587000, 247464, len(data)), data)
	for _, r := range relocs {
		*memmap.PtrPtr(0x587000, r[0]) = memmap.PtrOff(0x587000, r[1])
	}
	if memmap.Uint32(0x587000, 248196) != 0 || memmap.Uint32(0x587000, 248200) != 64 {
		t.Fatal("definition schema not supplied")
	}
	return o
}
func monsterControlWriteBin(t *testing.T, text string) {
	t.Helper()
	f, err := binfile.BinfileOpen("monster.bin", binfile.WriteOnly)
	if err != nil {
		t.Fatal(err)
	}
	if err = f.SetKey(23); err != nil {
		t.Fatal(err)
	}
	if _, err = f.Write([]byte(text)); err != nil {
		t.Fatal(err)
	}
	if err = f.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestMonsterControlDefinitionTokens(t *testing.T) {
	monsterControlDefsOwner(t)
	t.Chdir(t.TempDir())
	inputs := []string{"", " \\t", "one two \n", "\t\r\nOne\vTwo\fThree ", "// header\nAlpha // tail\nBeta ", "discard// rest\nkept ", "/ single // comment\nlast ", "// no newline", "tail", strings.Repeat("x", 240) + " "}
	type row struct {
		Text   string
		Tokens []legacy.PortTestMonsterToken
	}
	var rows []row
	for _, text := range inputs {
		monsterControlWriteBin(t, text)
		f, err := binfile.BinfileOpen("monster.bin", binfile.ReadOnly)
		if err != nil {
			t.Fatal(err)
		}
		if err = f.SetKey(23); err != nil {
			t.Fatal(err)
		}
		tokens := legacy.PortTestMonsterTokens(f)
		for _, r := range tokens {
			if !r.Intact {
				t.Fatal("token buffer guard")
			}
		}
		if text == "one two \n" && (len(tokens) != 3 || string(bytes.SplitN(tokens[0].Buffer, []byte{0}, 2)[0]) != "one" || string(bytes.SplitN(tokens[1].Buffer, []byte{0}, 2)[0]) != "two") {
			t.Fatal("whitespace token contract")
		}
		if text == "discard// rest\nkept " && (len(tokens) != 2 || string(bytes.SplitN(tokens[0].Buffer, []byte{0}, 2)[0]) != "kept") {
			t.Fatal("comment resets partial token")
		}
		rows = append(rows, row{text, tokens})
	}
	spellbookCapture(t, "monster-control-tokens", rows, "80f31d5de195968b639510bff7ed8624f72f9f7b5199fd350d862bbf96514ed0")
}

func TestMonsterControlDefinitions(t *testing.T) {
	o := monsterControlDefsOwner(t)
	t.Chdir(t.TempDir())
	t.Cleanup(o.s.PortTestRewardTypes([]string{"PortMonsterA", "PortMonsterB"}, nil, true, 0, 0))
	type fixture struct {
		Name, Text string
		Flags      uint32
	}
	cases := []fixture{
		{"empty", "", 0}, {"single", "PortMonsterA END ", 0}, {"pair", "PortMonsterA HEALTH 123 END PortMonsterB HEALTH 456 END ", 0},
		{"duplicate", "PortMonsterA HEALTH 123 END PortMonsterA HEALTH 456 END ", 0},
		{"partial", "PortMonsterA HEALTH 123 ", 0}, {"unknown", "PortMonsterA UNKNOWN 1 END ", 0},
		{"unknown-after-valid", "PortMonsterA END PortMonsterB UNKNOWN 1 END ", 0},
		{"missing-type", "AbsentMonster END ", 0},
		{"comments", "// header\nPortMonsterA // field\nhealth 123 // value\nEnD ", 0},
	}
	all := []string{"PortMonsterA"}
	kinds := map[uint32]int{}
	for i := uintptr(0); i < 29; i++ {
		p := uintptr(248192) + 12*i
		name := alloc.GoString((*byte)(*memmap.PtrPtr(0x587000, p)))
		kind := *memmap.PtrUint32(0x587000, p+4)
		kinds[kind]++
		value := ""
		switch kind {
		case 0:
			value = "-17"
		case 1:
			value = "-1.375"
		case 3:
			value = "MonsterStrike"
		case 4:
			value = "MonsterDie"
		case 5:
			value = "MonsterDead"
		case 6:
			value = alloc.GoString((*byte)(*memmap.PtrPtr(0x587000, 247536))) + "+" + alloc.GoString((*byte)(*memmap.PtrPtr(0x587000, 247596)))
		case 7:
			value = "Arrow"
		case 8:
			value = alloc.GoString((*byte)(*memmap.PtrPtr(0x587000, 247464)))[7:]
		default:
			t.Fatal("unaccounted shipped property kind", kind)
		}
		// Use the actual first callback names and addresses from the shipped tables.
		if kind >= 3 && kind <= 5 {
			base := map[uint32]uintptr{3: 287096, 4: 287280, 5: 287192}[kind]
			value = alloc.GoString((*byte)(*memmap.PtrPtr(0x587000, base)))
		}
		cases = append(cases, fixture{name, "PortMonsterA " + strings.ToLower(name) + " " + value + " END ", 0})
		all = append(all, name, value)
	}
	if len(kinds) != 8 {
		t.Fatal("missing shipped property kinds", kinds)
	}
	all = append(all, "END ")
	cases = append(cases, fixture{"all-fields", strings.Join(all, " "), 0})
	for _, flags := range []uint32{0, 2048, 0x200000, 0x2000, 2048 | 0x2000} {
		cases = append(cases, fixture{"mode-filter", "PortMonsterA HEALTH 10\nARENA HEALTH 20\nSOLO HEALTH 30\nEND \n", flags})
	}
	for _, prop := range []string{"MELEE_STRIKE_FUNCTION", "DIE_FUNCTION", "DEAD_FUNCTION", "MISSILE_NAME"} {
		for _, value := range []string{"NULL", "null", "unknown"} {
			cases = append(cases, fixture{prop + "-" + value, "PortMonsterA " + prop + " " + value + " END ", 0})
		}
	}
	for _, value := range []string{"0", "2147483647", "-2147483648", "123tail", "invalid"} {
		cases = append(cases, fixture{"integer-" + value, "PortMonsterA HEALTH " + value + " END ", 0})
	}
	for _, value := range []string{"0", "-0", "1.000000059604644775390625", "1e-40", "1e40", "inf", "nan", "invalid"} {
		cases = append(cases, fixture{"float-" + value, "PortMonsterA FLEE_RANGE " + value + " END ", 0})
	}
	for i := uintptr(0); i < 18; i++ {
		name := alloc.GoString((*byte)(*memmap.PtrPtr(0x587000, 247464+4*i)))[7:]
		cases = append(cases, fixture{"damage-" + name, "PortMonsterA MELEE_ATTACK_DAMAGE_TYPE " + strings.ToLower(name) + " END ", 0})
	}
	cases = append(cases, fixture{"bad-damage", "PortMonsterA MELEE_ATTACK_DAMAGE_TYPE unknown END ", 0})
	type row struct {
		Name             string
		Flags            uint32
		Load, Bind, Free int
		Before, After    [][62]uint32
		Lookup           []int
	}
	var rows []row
	for _, c := range cases {
		restore := noxflags.PortTestGameFlags(noxflags.GameFlag(c.Flags))
		monsterControlWriteBin(t, c.Text)
		r := row{Name: c.Name, Flags: c.Flags, Load: legacy.PortTestMonsterDefs("load", 0)}
		r.Before = legacy.PortTestMonsterDefsSnapshot()
		if r.Load != 1 {
			t.Fatal("opened definition file", c.Name)
		}
		if c.Name == "pair" && (len(r.Before) != 2 || r.Before[0][17] != 456 || r.Before[1][17] != 123) {
			t.Fatal("definition prepend order")
		}
		if c.Name == "all-fields" && (len(r.Before) != 1 || r.Before[0][17] != 0xffffffef || r.Before[0][22] != math.Float32bits(-1.375)) {
			t.Fatal("typed numeric fields")
		}
		r.Bind = legacy.PortTestMonsterDefs("bind", 0)
		r.After = legacy.PortTestMonsterDefsSnapshot()
		for _, id := range []int{0, o.s.Types.IndByID("PortMonsterA"), o.s.Types.IndByID("PortMonsterB"), 0x7fffffff} {
			r.Lookup = append(r.Lookup, legacy.PortTestMonsterDefs("lookup", id))
		}
		if c.Name == "missing-type" && (r.Bind != 0 || len(r.After) != 0) {
			t.Fatal("missing type must release complete list")
		}
		if c.Name == "duplicate" && r.Lookup[1] != 1 {
			t.Fatal("lookup must select newest duplicate")
		}
		r.Free = legacy.PortTestMonsterDefs("free", 0)
		if r.Free != 0 || len(legacy.PortTestMonsterDefsSnapshot()) != 0 || legacy.PortTestMonsterDefs("free", 0) != 0 {
			t.Fatal("definition list free")
		}
		rows = append(rows, r)
		restore()
	}
	if err := os.Remove("monster.bin"); err != nil {
		t.Fatal(err)
	}
	if legacy.PortTestMonsterDefs("load", 0) != 0 {
		t.Fatal("missing file accepted")
	}
	spellbookCapture(t, "monster-control-definitions", rows, "05cd0830e8603615794360e3f09c32bf5938c784d33216802f88fe72d608cc46")
}

func TestMonsterControlShippedDefinitions(t *testing.T) {
	path := os.Getenv("NOX_DATA")
	if path == "" {
		t.Skip("set NOX_DATA to original assets")
	}
	data, err := os.ReadFile(filepath.Join(path, "monster.bin"))
	if err != nil {
		t.Fatal(err)
	}
	o := monsterControlDefsOwner(t)
	t.Chdir(t.TempDir())
	if err = os.WriteFile("monster.bin", data, 0600); err != nil {
		t.Fatal(err)
	}
	type row struct {
		Flags       uint32
		Load, Bind  int
		Definitions [][62]uint32
		Lookups     []int
	}
	var rows []row
	for _, flags := range []uint32{2048, 0x200000, 0x2000} {
		restore := noxflags.PortTestGameFlags(noxflags.GameFlag(flags))
		r := row{Flags: flags, Load: legacy.PortTestMonsterDefs("load", 0)}
		r.Definitions = legacy.PortTestMonsterDefsSnapshot()
		if r.Load != 1 || len(r.Definitions) < 40 {
			t.Fatal("shipped definitions incomplete", flags, len(r.Definitions))
		}
		names := make([]string, len(r.Definitions))
		for i, d := range r.Definitions {
			var b [64]byte
			for j := 0; j < 16; j++ {
				binary.LittleEndian.PutUint32(b[4*j:], d[j])
			}
			names[i] = string(bytes.SplitN(b[:], []byte{0}, 2)[0])
		}
		restoreTypes := o.s.PortTestRewardTypes(names, nil, true, 0, 0)
		r.Bind = legacy.PortTestMonsterDefs("bind", 0)
		if r.Bind != 1 {
			t.Fatal("shipped definition binding")
		}
		for i, name := range names {
			v := legacy.PortTestMonsterDefs("lookup", o.s.Types.IndByID(name))
			if v == 0 {
				t.Fatal(fmt.Sprintf("shipped lookup %d %s", i, name))
			}
			r.Lookups = append(r.Lookups, v)
		}
		legacy.PortTestMonsterDefs("free", 0)
		restoreTypes()
		restore()
		rows = append(rows, r)
	}
	spellbookCapture(t, "monster-control-shipped-definitions", rows, "acb2d4bf4067e00c81182fbc2d54ce81e071787a353880060de4ca68ab33002d")
}
