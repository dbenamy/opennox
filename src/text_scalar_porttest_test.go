//go:build porttest

package opennox

import (
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func textScalarFold(c uint16) int {
	if c >= 'A' && c <= 'Z' {
		c += 'a' - 'A'
	}
	return int(c)
}
func TestTextScalarsWide(t *testing.T) {
	a, freeA := alloc.Make([]uint16{}, 4)
	defer freeA()
	b, freeB := alloc.Make([]uint16{}, 4)
	defer freeB()
	type row struct {
		A, B   uint16
		Prefix bool
		Return int
	}
	var rows []row
	for x := 0; x < 65536; x++ {
		for _, y := range []uint16{0, 'A', uint16(x)} {
			for _, prefix := range []bool{false, true} {
				clear(a)
				clear(b)
				i := 0
				if prefix {
					a[0] = 'X'
					b[0] = 'x'
					i = 1
				}
				a[i] = uint16(x)
				b[i] = y
				got := legacy.PortTestTextCompareWide(&a[0], &b[0])
				want := textScalarFold(uint16(x)) - textScalarFold(y)
				if got != want {
					t.Fatalf("wide x=%x y=%x prefix=%v: %d/%d", x, y, prefix, got, want)
				}
				rows = append(rows, row{uint16(x), y, prefix, got})
			}
		}
	}
	// A terminator ends comparison even if both backing buffers have suffix data.
	copy(a, []uint16{'A', 0, 'x', 0})
	copy(b, []uint16{'a', 0, 'y', 0})
	if legacy.PortTestTextCompareWide(&a[0], &b[0]) != 0 {
		t.Fatal("wide terminator")
	}
	interactionCapture(t, "text-scalars-wide", rows)
}
func TestTextScalarsNarrow(t *testing.T) {
	a, freeA := alloc.Make([]byte{}, 4)
	defer freeA()
	b, freeB := alloc.Make([]byte{}, 4)
	defer freeB()
	type row struct {
		A, B   byte
		Prefix bool
		Return int
	}
	var rows []row
	for x := 0; x < 256; x++ {
		for y := 0; y < 256; y++ {
			for _, prefix := range []bool{false, true} {
				clear(a)
				clear(b)
				i := 0
				if prefix {
					a[0] = 'X'
					b[0] = 'x'
					i = 1
				}
				a[i] = byte(x)
				b[i] = byte(y)
				want := textScalarFold(uint16(x)) - textScalarFold(uint16(y))
				want = max(-1, min(1, want))
				got := legacy.PortTestTextCompareNarrow(&a[0], &b[0])
				if got != want {
					t.Fatalf("narrow %x %x %v: %d/%d", x, y, prefix, got, want)
				}
				rows = append(rows, row{byte(x), byte(y), prefix, got})
			}
		}
	}
	copy(a, []byte{'A', 0, 'x', 0})
	copy(b, []byte{'a', 0, 'y', 0})
	if legacy.PortTestTextCompareNarrow(&a[0], &b[0]) != 0 {
		t.Fatal("narrow terminator")
	}
	interactionCapture(t, "text-scalars-narrow", rows)
}
func TestTextScalarsDecimal(t *testing.T) {
	type row struct {
		Input  []uint16
		Return int32
	}
	var rows []row
	check := func(input []uint16, want int32) {
		p, free := alloc.Make([]uint16{}, len(input)+1)
		copy(p, input)
		got := legacy.PortTestTextDecimal(&p[0])
		free()
		if got != want {
			t.Fatalf("decimal %x: %d/%d", input, got, want)
		}
		rows = append(rows, row{slices.Clone(input), got})
	}
	for _, c := range []struct {
		s string
		v int32
	}{{"", 0}, {" ", 0}, {"+", 0}, {"-", 0}, {"+ 1", 0}, {"- 1", 0}, {"--1", 0}, {"++1", 0}, {"0", 0}, {"-0", 0}, {"00012", 12}, {"0x12", 0}, {"0b11", 0}, {"12tail", 12}, {"2147483647", 2147483647}, {"2147483648", 2147483647}, {"-2147483648", -2147483648}, {"-2147483649", -2147483648}, {strings.Repeat("9", 200), 2147483647}, {"-" + strings.Repeat("9", 200), -2147483648}, {"12\x0034", 12}, {"１２", 0}, {"\u200312", 0}, {"12é34", 12}} {
		check(textFormatWide(c.s), c.v)
	}
	for _, prefix := range []string{" ", "\t", "\n", "\v", "\f", "\r", " \t\n\v\f\r"} {
		check(textFormatWide(prefix+"-123"), -123)
	}
	// Exhaust each UTF-16 code unit inside a decimal token. Non-ASCII is replaced
	// by DEL by C and terminates the token; ASCII digits continue it.
	for ch := 0; ch < 65536; ch++ {
		want := int32(12)
		if ch >= '0' && ch <= '9' {
			v, _ := strconv.ParseInt("12"+string(rune(ch))+"34", 10, 32)
			want = int32(v)
		}
		check([]uint16{'1', '2', uint16(ch), '3', '4'}, want)
	}
	interactionCapture(t, "text-scalars-decimal", rows)
}
func TestTextScalarsCopy(t *testing.T) {
	src, freeS := alloc.Make([]uint16{}, 5)
	defer freeS()
	dst, freeD := alloc.Make([]uint16{}, 9)
	defer freeD()
	type row struct {
		Unit  uint16
		Words []uint16
	}
	var rows []row
	for ch := 0; ch < 65536; ch++ {
		copy(src, []uint16{'x', uint16(ch), 'y', 0, 0xbeef})
		for i := range dst {
			dst[i] = 0xa55a
		}
		got := legacy.PortTestTextCopy(&dst[2], &src[0])
		want := []uint16{0xa55a, 0xa55a, 'x', uint16(ch), 'y', 0, 0xa55a, 0xa55a, 0xa55a}
		if ch == 0 {
			want[4], want[5] = 0xa55a, 0xa55a
		}
		if got != &dst[2] || !slices.Equal(dst, want) || src[4] != 0xbeef {
			t.Fatalf("copy unit=%x: %x/%x", ch, dst, want)
		}
		rows = append(rows, row{uint16(ch), slices.Clone(dst)})
	}
	if legacy.PortTestTextCopy(&src[0], &src[0]) != &src[0] {
		t.Fatal("copy self pointer")
	}
	interactionCapture(t, "text-scalars-copy", rows)
}
