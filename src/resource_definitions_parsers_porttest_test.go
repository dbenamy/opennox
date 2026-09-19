//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"math"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

type resourceParserRow struct {
	Kind, Input   string
	Fill          byte
	Return        int
	Data, Scratch []byte
}

func resourceParse(t *testing.T, kind, input string, fill byte) resourceParserRow {
	t.Helper()
	scratch, free := alloc.Make([]byte{}, 512)
	defer free()
	data, freeData := alloc.Make([]byte{}, 256)
	defer freeData()
	for i := range scratch {
		scratch[i] = 0xcc
	}
	copy(scratch, input)
	scratch[len(input)] = 0
	for i := range data {
		data[i] = fill
	}
	ret := legacy.PortTestResourceParser(kind, unsafe.Pointer(&scratch[0]), unsafe.Pointer(&data[0]))
	// Defined returns are zero or one. Never normalize a random stack address.
	if ret != 0 && ret != 1 {
		t.Fatalf("%s returned unexpected pointer/value %x", kind, ret)
	}
	return resourceParserRow{kind, input, fill, ret, append([]byte(nil), data...), append([]byte(nil), scratch...)}
}
func resourcePut32(dst []byte, off int, v uint32) { binary.LittleEndian.PutUint32(dst[off:], v) }
func TestResourceDefinitionsIntegerParsers(t *testing.T) {
	var rows []resourceParserRow
	cases := []struct {
		input string
		value int32
		valid bool
	}{
		{"", 0, false}, {"bad", 0, false}, {" \t\r\n", 0, false},
		{"0", 0, true}, {"-1", -1, true}, {"+17", 17, true}, {"255", 255, true}, {"256", 256, true},
		{"-257", -257, true}, {"2147483647", 2147483647, true}, {"-2147483648", -2147483648, true},
		{"  42tail", 42, true}, {"0x12", 0, true}, {"0019 99", 19, true}, {"\t-23\n", -23, true},
	}
	for _, kind := range []string{"lifetime", "projectile", "spark", "mana", "arrow"} {
		for _, c := range cases {
			// spark scan failure leaves a stack pointer; mana/arrow require a token.
			if kind == "spark" && !c.valid || (kind == "mana" || kind == "arrow") && (c.input == "" || c.input == " \t\r\n") {
				continue
			}
			for _, fill := range []byte{0, 0xa5} {
				row := resourceParse(t, kind, c.input, fill)
				want := bytes.Repeat([]byte{fill}, 256)
				switch kind {
				case "lifetime", "projectile":
					if c.valid {
						resourcePut32(want, 0, uint32(c.value))
					}
				case "spark", "mana":
					want[0] = byte(c.value)
				case "arrow":
					resourcePut32(want, 0, uint32(c.value))
					resourcePut32(want, 4, uint32(c.value))
				}
				if row.Return != 1 || !bytes.Equal(row.Data, want) {
					t.Fatalf("%s %q: got %x want %x", kind, c.input, row.Data[:12], want[:12])
				}
				rows = append(rows, row)
			}
		}
	}
	spellbookCapture(t, "resource-definitions-integers", rows, "d3ff3ccb83d92bece2f7cba7012329bca4641f700d04f3aac1343417397065ac")
}
func TestResourceDefinitionsPartialScans(t *testing.T) {
	var rows []resourceParserRow
	triples := []struct {
		input string
		vals  []int32
	}{
		{"", nil}, {"bad", nil}, {"7", []int32{7}}, {"7 bad 9", []int32{7}},
		{"7 -2", []int32{7, -2}}, {"7 -2 99 trailing", []int32{7, -2, 99}},
		{"+4\t5\n-6", []int32{4, 5, -6}}, {"0x2 3 4", []int32{0}},
	}
	pushes := []struct {
		input string
		vals  []float32
	}{
		{"", nil}, {"bad", nil}, {"1.5", []float32{1.5}}, {"1.5 bad", []float32{1.5}},
		{"-2.25 3.5 extra", []float32{-2.25, 3.5}}, {"0 -0", []float32{0, math.Float32frombits(0x80000000)}},
		{"1e-30 1e30", []float32{1e-30, 1e30}}, {"0x1.8p+1 4", []float32{3, 4}},
	}
	for _, fill := range []byte{0, 0xa5} {
		for _, c := range triples {
			row := resourceParse(t, "triple", c.input, fill)
			want := bytes.Repeat([]byte{fill}, 256)
			for i, v := range c.vals {
				resourcePut32(want, 4*i, uint32(v))
			}
			if row.Return != 1 || !bytes.Equal(row.Data, want) {
				t.Fatal("triple", c.input, row.Data[:12], want[:12])
			}
			rows = append(rows, row)
		}
		for _, c := range pushes {
			row := resourceParse(t, "push", c.input, fill)
			want := bytes.Repeat([]byte{fill}, 256)
			if len(c.vals) > 0 {
				resourcePut32(want, 0, math.Float32bits(c.vals[0]))
			}
			if len(c.vals) > 1 {
				resourcePut32(want, 8, math.Float32bits(c.vals[1]))
			}
			copy(want[4:8], want[:4])
			if row.Return != 1 || !bytes.Equal(row.Data, want) {
				t.Fatal("push", c.input, row.Data[:12], want[:12])
			}
			rows = append(rows, row)
		}
	}
	spellbookCapture(t, "resource-definitions-partial-scans", rows, "fa07e29ada7387a239e3111df79dca23618a1eaaf2fc04637fefbdd6bbbc6232")
}
