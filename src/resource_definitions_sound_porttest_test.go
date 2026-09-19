//go:build porttest

package opennox

import (
	"bytes"
	"fmt"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/sound"
	"github.com/opennox/opennox/v1/internal/binfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
	"strings"
	"testing"
	"unsafe"
)

func resourceSoundOwner(t *testing.T) {
	t.Helper()
	t.Cleanup(handles.PortTestInit())
	t.Cleanup(legacy.PortTestResourceSoundOwner())
	t.Chdir(t.TempDir())
}
func resourceSoundFile(t *testing.T, text string) {
	t.Helper()
	f, err := binfile.BinfileOpen("soundset.bin", binfile.WriteOnly)
	if err != nil {
		t.Fatal(err)
	}
	if err = f.SetKey(5); err != nil {
		t.Fatal(err)
	}
	if _, err = f.Write([]byte(text)); err != nil {
		t.Fatal(err)
	}
	if err = f.Close(); err != nil {
		t.Fatal(err)
	}
}
func TestResourceDefinitionsTokens(t *testing.T) {
	resourceSoundOwner(t)
	var rows []map[string]any
	for _, text := range []string{"", " \t\r\n", "one two \n", "\t\r\nOne\vTwo\fThree ", "// header\nAlpha // tail\nBeta ", "discard// rest\nkept ", "/ single // comment\nlast ", "// no newline", "tail", strings.Repeat("x", 240) + " "} {
		resourceSoundFile(t, text)
		f, err := binfile.BinfileOpen("soundset.bin", binfile.ReadOnly)
		if err != nil {
			t.Fatal(err)
		}
		if err = f.SetKey(5); err != nil {
			t.Fatal(err)
		}
		tokens := legacy.PortTestResourceTokens(f)
		for _, r := range tokens {
			if !r.Intact {
				t.Fatal("token buffer guard")
			}
		}
		if text == "one two \n" && (len(tokens) != 3 || string(bytes.SplitN(tokens[0].Buffer, []byte{0}, 2)[0]) != "one" || string(bytes.SplitN(tokens[1].Buffer, []byte{0}, 2)[0]) != "two") {
			t.Fatal("whitespace token contract")
		}
		if text == "discard// rest\nkept " && (len(tokens) != 3 || string(bytes.SplitN(tokens[0].Buffer, []byte{0}, 2)[0]) != "rest" || string(bytes.SplitN(tokens[1].Buffer, []byte{0}, 2)[0]) != "kept") {
			t.Fatal("legacy SkipLine consumes first non-newline byte, leaving rest of comment")
		}
		rows = append(rows, map[string]any{"text": text, "tokens": tokens})
	}
	spellbookCapture(t, "resource-definitions-tokens", rows, "c8627ee90ecb326269a119ab4c087500c09bb298354aa52bbdd1b40f0fad47b0")
}
func TestResourceDefinitionsSoundSets(t *testing.T) {
	resourceSoundOwner(t)
	table := unsafe.Slice(memmap.PtrUint32(0x587000, 64704), 38)
	old := append([]uint32(nil), table...)
	defer copy(table, old)
	for i := 0; i < 18; i++ {
		p, free := alloc.CString(fmt.Sprintf("KEY%d", i))
		defer free()
		table[2*i] = uint32(uintptr(unsafe.Pointer(p)))
		table[2*i+1] = uint32(4 + 4*i)
	}
	table[36] = 0
	table[37] = 0xdeadbeef
	name := sound.ID(226).String()
	val := uint32(sound.ByName(name))
	if val == 0 {
		t.Fatal("sound fixture")
	}
	full := "Full "
	for i := 0; i < 18; i++ {
		full += fmt.Sprintf("KEY%d %s ", i, name)
	}
	full += "END "
	cases := []struct {
		text       string
		ret, count int
	}{
		{"", 1, 0}, {"A END ", 1, 1}, {"A END B END ", 1, 2},
		{"A KEY0 " + name + " END ", 1, 1}, {"A KEY0 " + name + " KEY0 missing END ", 1, 1},
		{"A KEY0 " + name + " ", 1, 1}, {"A UNKNOWN " + name + " END ", 0, 1},
		{"A END B UNKNOWN " + name + " END ", 0, 2}, {"A END A KEY0 " + name + " END ", 1, 2},
		{"// comment\nA KEY0 " + name + " // tail\nEND ", 0, 1}, {full, 1, 1},
	}
	var rows []map[string]any
	for _, c := range cases {
		legacy.PortTestResourceSoundClear()
		resourceSoundFile(t, c.text)
		ret := legacy.Nox_xxx_parseSoundSetBin_424170("soundset.bin")
		state := legacy.PortTestResourceSounds()
		if ret != c.ret || len(state) != c.count {
			t.Fatal("load", c.text, ret, state)
		}
		if c.text == full {
			for i := 0; i < 18; i++ {
				if state[0].Words[i] != val {
					t.Fatal("sound row stride", i)
				}
			}
		}
		lookups := map[string]int{}
		for _, key := range []string{"A", "B", "Full", "missing", "a"} {
			lookups[key] = legacy.PortTestResourceSoundLookup(key)
		}
		if c.text == "A END B END " && (lookups["B"] != 1 || lookups["A"] != 2) {
			t.Fatal("sound list order", lookups)
		}
		if c.text == "A END A KEY0 "+name+" END " && lookups["A"] != 1 {
			t.Fatal("duplicate first match")
		}
		rows = append(rows, map[string]any{"text": c.text, "ret": ret, "state": state, "lookups": lookups})
	}
	if legacy.Nox_xxx_parseSoundSetBin_424170("absent.bin") != 0 {
		t.Fatal("missing sound file")
	}
	spellbookCapture(t, "resource-definitions-soundsets", rows, "a4d95d2ff7116089942401a5392dae43c4d34142be8999ed356d1be4670bc806")
}
