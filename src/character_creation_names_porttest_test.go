//go:build porttest

package opennox

import (
	"encoding/binary"
	"strings"
	"testing"
	"unicode/utf16"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestCharacterCreationNames(t *testing.T) {
	type row struct {
		Input, Output string
		Return        int
		Data          []byte
	}
	var rows []row
	cases := []struct {
		input, want string
		ret         int
	}{
		{"", "", 0}, {" ", " ", 0}, {"\t\r\n", "\t\r\n", 0},
		{"Hero", "Hero", 1}, {"  Hero  ", "Hero", 1}, {"\tA B\n", "A B", 1},
		{"A*?<>\\/:\"|", "A*?<>\\/:\"|", 1}, {"é英雄", "é英雄", 1},
		{strings.Repeat("A", 25), strings.Repeat("A", 25), 1},
	}
	for _, c := range `*?<>\/:"|` {
		cases = append(cases, struct {
			input, want string
			ret         int
		}{"  " + string(c) + "A" + string(c) + " ", "-A" + string(c), 1})
	}
	for _, c := range cases {
		p, free := alloc.Calloc(1, 80)
		raw := unsafe.Slice((*byte)(p), 80)
		for i := range raw {
			raw[i] = 0xa5
		}
		u := append(utf16.Encode([]rune(c.input)), 0)
		for i, v := range u {
			binary.LittleEndian.PutUint16(raw[2*i:], v)
		}
		before := append([]byte(nil), raw...)
		ret := legacy.PortTestCharacterName(p)
		var got []uint16
		for i := 0; i < 40; i++ {
			v := binary.LittleEndian.Uint16(raw[2*i:])
			if v == 0 {
				break
			}
			got = append(got, v)
		}
		output := string(utf16.Decode(got))
		if ret != c.ret || output != c.want {
			t.Fatalf("%q => %q/%d want %q/%d", c.input, output, ret, c.want, c.ret)
		}
		for i := 2 * len(u); i < len(raw); i++ {
			if raw[i] != before[i] {
				t.Fatalf("write after input terminator at %d", i)
			}
		}
		rows = append(rows, row{c.input, output, ret, append([]byte(nil), raw...)})
		free()
	}
	spellbookCapture(t, "character-creation-names", rows, "0d67a9e55b3853feb5e4f7209f32997cb7733da71dcf8b93521104552c0ad463")
}
