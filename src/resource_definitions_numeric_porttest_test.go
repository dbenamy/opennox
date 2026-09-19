//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"math"
	"testing"
)

func TestResourceDefinitionsNumericEdges(t *testing.T) {
	var rows []resourceParserRow
	// Linux/386 libc clamps decimal overflow to the signed long endpoints.
	integers := []struct {
		s string
		v int32
	}{
		{"2147483648", 2147483647}, {"4294967295", 2147483647}, {"999999999999999999999999999", 2147483647},
		{"-2147483649", -2147483648}, {"-999999999999999999999999999", -2147483648},
	}
	for _, c := range integers {
		for _, kind := range []string{"lifetime", "projectile", "spark", "mana", "arrow"} {
			row := resourceParse(t, kind, c.s, 0xa5)
			want := bytes.Repeat([]byte{0xa5}, 256)
			if kind == "spark" || kind == "mana" {
				want[0] = byte(c.v)
			} else {
				resourcePut32(want, 0, uint32(c.v))
				if kind == "arrow" {
					resourcePut32(want, 4, uint32(c.v))
				}
			}
			if row.Return != 1 || !bytes.Equal(row.Data, want) {
				t.Fatalf("integer overflow %s %q: %x want%x", kind, c.s, row.Data[:8], want[:8])
			}
			rows = append(rows, row)
		}
	}
	floats := []struct {
		s    string
		bits uint32
	}{
		{"+INF", 0x7f800000}, {"-infinity", 0xff800000}, {"1e39", 0x7f800000}, {"-1e39", 0xff800000},
		{"1e-50", 0}, {"-1e-50", 0x80000000}, {"1.000000059604644775390625", 0x3f800000},
		{"1.000000059604644775390626", 0x3f800001}, {"0x1.000001p0", 0x3f800000},
	}
	for _, c := range floats {
		row := resourceParse(t, "push", c.s+" 2", 0xa5)
		want := bytes.Repeat([]byte{0xa5}, 256)
		resourcePut32(want, 0, c.bits)
		resourcePut32(want, 4, c.bits)
		resourcePut32(want, 8, math.Float32bits(2))
		if row.Return != 1 || !bytes.Equal(row.Data, want) {
			t.Fatalf("float edge %q: %08x want%08x", c.s, binary.LittleEndian.Uint32(row.Data), c.bits)
		}
		rows = append(rows, row)
	}
	spellbookCapture(t, "resource-definitions-numeric-edges", rows, "c4f59d1908e91b3ef89b39716c54d8cc9173e8aa823abba2c224b9d2ff88e149")
}
