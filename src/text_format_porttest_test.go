//go:build porttest

package opennox

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
	"testing"
	"unicode/utf16"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

type textFormatCase struct {
	Kind   int
	Format string
	A, B   uint32
	Bits   uint64
	W1, W2 []uint16
	Narrow []byte
	Want   []uint16
}

func textFormatWide(s string) []uint16 { return utf16.Encode([]rune(s)) }
func textFormatRun(t *testing.T, name string, cases []textFormatCase) {
	t.Helper()
	type row struct {
		Case, Count int
		Nil         bool
		Return      int
		Words       []uint16
	}
	var rows []row
	for i, sp := range cases {
		func() {
			wide := func(s []uint16) (*uint16, func()) {
				if s == nil {
					return nil, func() {}
				}
				v, free := alloc.Make([]uint16{}, len(s)+1)
				copy(v, s)
				return &v[0], free
			}
			f, freeF := wide(append(textFormatWide(sp.Format), 0))
			defer freeF()
			w1, free1 := wide(sp.W1)
			defer free1()
			w2, free2 := wide(sp.W2)
			defer free2()
			var narrow *byte
			if sp.Narrow != nil {
				v, free := alloc.Make([]byte{}, len(sp.Narrow)+1)
				defer free()
				copy(v, sp.Narrow)
				narrow = &v[0]
			}
			// Capacity covers every formatted unit plus terminator and guards. Count
			// describes the exact historical bounded-write contract, not string length.
			counts := []int{0, 1, max(0, len(sp.Want)-1), len(sp.Want), len(sp.Want) + 1, len(sp.Want) + 9}
			for _, count := range counts {
				for _, nilDst := range []bool{false, true} {
					out, free := alloc.Make([]uint16{}, count+4)
					for j := range out {
						out[j] = 0xa55a
					}
					dst := &out[2]
					if nilDst {
						dst = nil
					}
					got := legacy.PortTestTextFormat(sp.Kind, dst, uint32(count), f, sp.A, sp.B, sp.Bits, w1, w2, narrow)
					expected := make([]uint16, len(out))
					for j := range expected {
						expected[j] = 0xa55a
					}
					if !nilDst {
						want := append(slices.Clone(sp.Want), 0)
						copy(expected[2:2+count], want)
					}
					if got != len(sp.Want) || !slices.Equal(out, expected) {
						t.Fatalf("format case=%d kind=%d fmt=%q a=%x b=%x float=%x count=%d nil=%v return=%d/%d words=%x/%x", i, sp.Kind, sp.Format, sp.A, sp.B, sp.Bits, count, nilDst, got, len(sp.Want), out, expected)
					}
					rows = append(rows, row{i, count, nilDst, got, slices.Clone(out)})
					free()
				}
			}
		}()
	}
	interactionCapture(t, name, rows)
	t.Logf("%d formats /%d bounded-output cases", len(cases), len(rows))
}
func TestTextFormatInteger(t *testing.T) {
	var cases []textFormatCase
	for _, verb := range []byte{'d', 'i', 'u', 'x', 'X', 'o'} {
		for _, v := range []uint32{0, 1, 9, 10, 15, 16, 255, 0x7fffffff, 0x80000000, 0xffffffff} {
			base := 10
			if verb == 'x' || verb == 'X' {
				base = 16
			}
			if verb == 'o' {
				base = 8
			}
			digits := strconv.FormatInt(int64(int32(v)), base)
			if verb == 'X' {
				digits = strings.ToUpper(digits)
			}
			for _, style := range []struct {
				format           string
				width, precision int
				zero             bool
			}{{"%", -1, -1, false}, {"%+", -1, -1, false}, {"%-12", 12, -1, false}, {"%012", 12, -1, true}, {"%15.10", 15, 10, false}, {"%015.10", 15, 10, true}, {"%.0", -1, 0, false}, {"%.15", -1, 15, false}, {"%2.1", 2, 1, false}} {
				// Frozen C's numeric adapter is signed even for u/x/o. Padding counts
				// include the sign, and a positive precision controls outer padding.
				used := len(digits)
				if style.precision > 0 {
					used = style.precision
				}
				pad := " "
				if style.zero {
					pad = "0"
				}
				want := strings.Repeat(pad, max(0, style.width-used)) + strings.Repeat("0", max(0, style.precision-len(digits))) + digits
				cases = append(cases, textFormatCase{Kind: 1, Format: "[" + style.format + string(verb) + "]", A: v, Want: textFormatWide("[" + want + "]")})
			}
		}
	}
	cases = append(cases, textFormatCase{Kind: 1, Format: "%d/%08X", A: 0xfffffffd, B: 255, Want: textFormatWide("-3/000000FF")})
	textFormatRun(t, "text-format-integer", cases)
}
func TestTextFormatStrings(t *testing.T) {
	var cases []textFormatCase
	for _, w := range [][]uint16{nil, {}, {'A'}, {'x', 0, 'y'}, {0x00ff, 0x0100, 0xd800, 0xdc00, 0xffff}, textFormatWide("Ång 😀")} {
		want := slices.Clone(w)
		if w == nil {
			want = textFormatWide("(null)")
		}
		if n := slices.Index(want, 0); n >= 0 {
			want = want[:n]
		}
		for _, f := range []string{"%s", "%20s", "%.1s", "%-20.1s", "%+s"} {
			cases = append(cases, textFormatCase{Kind: 3, Format: f, W1: w, Want: want})
		}
	}
	for _, b := range [][]byte{nil, {}, {'A'}, {'x', 0, 'y'}, {0x7f, 0x80, 0xff}, []byte("Ång 😀")} {
		var want []uint16
		if b == nil {
			want = textFormatWide("(null)")
		} else {
			for _, v := range b {
				if v == 0 {
					break
				}
				want = append(want, uint16(v))
			}
		}
		for _, f := range []string{"%S", "%20S", "%.1S", "%-20.1S", "%+S"} {
			cases = append(cases, textFormatCase{Kind: 4, Format: f, Narrow: b, Want: want})
		}
	}
	for _, v := range []uint32{0, 1, 65, 255, 256, 0xd800, 0xffff, 0x10000, 0xffffffff} {
		cases = append(cases, textFormatCase{Kind: 1, Format: "[%10.2c]", A: v, Want: []uint16{'[', uint16(v), ']'}})
	}
	for _, s := range []string{"", "hello", "Ång 😀", "100%%", "a%!b", "%%s"} {
		want := strings.NewReplacer("%%", "%", "%!", "!").Replace(s)
		cases = append(cases, textFormatCase{Kind: 0, Format: s, Want: textFormatWide(want)})
	}
	cases = append(cases, textFormatCase{Kind: 3, Format: "%s|%s", W1: []uint16{}, W2: nil, Want: textFormatWide("|(null)")})
	textFormatRun(t, "text-format-strings", cases)
}
func TestTextFormatFloat(t *testing.T) {
	var cases []textFormatCase
	values := []uint64{0, 1, 0x8000000000000000, 0x8000000000000001, math.Float64bits(1), math.Float64bits(-1), math.Float64bits(1.25), math.Float64bits(1.75), math.Float64bits(-1.25), math.Float64bits(1.000005), math.Float64bits(12345.6789), math.Float64bits(1e30), math.Float64bits(1e100), math.Float64bits(math.MaxFloat64), 0x7ff0000000000000, 0xfff0000000000000, 0x7ff8000000000000, 0xfff8000000000001}
	for _, bits := range values {
		for _, precision := range []int{-1, 0, 1, 2, 5, 10, 40} {
			for _, prefix := range []string{"%", "%+", "%-25", "%025"} {
				format := prefix
				if precision >= 0 {
					format += fmt.Sprintf(".%d", precision)
				}
				format += "f"
				digits := precision
				if digits <= 0 {
					digits = 5
				}
				value := math.Float64frombits(bits)
				want := strconv.FormatFloat(value, 'f', digits, 64)
				if math.IsInf(value, 1) {
					want = "inf"
				}
				if math.IsInf(value, -1) {
					want = "-inf"
				}
				if math.IsNaN(value) {
					want = "nan"
					if math.Signbit(value) {
						want = "-nan"
					}
				}
				if len(want) > 31 {
					want = want[:31]
				}
				cases = append(cases, textFormatCase{Kind: 2, Format: format, Bits: bits, Want: textFormatWide(want)})
			}
		}
	}
	cases = append(cases, textFormatCase{Kind: 5, Format: "[%s][%S][%03d][%.2f]", W1: []uint16{'W', 0xd800}, Narrow: []byte{'N', 0xff}, A: 0xfffffffd, Bits: math.Float64bits(1.25), Want: []uint16{'[', 'W', 0xd800, ']', '[', 'N', 0xff, ']', '[', '0', '-', '3', ']', '[', '1', '.', '2', '5', ']'}})
	textFormatRun(t, "text-format-float", cases)
}
