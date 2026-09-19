//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"strings"
	"testing"
)

func TestSessionMOTDLineBoundaries(t *testing.T) {
	// Independent line contract: CR, LF and CRLF terminate a line. An unterminated
	// final line still copies its text but returns NULL; empty input also returns NULL.
	cases := []string{"", "x", "\r", "\n", "\r\n", "\n\r", "a\r\nb", "a\nb", "a\rb", "\x80\xff\t x"}
	for _, n := range []int{1, 2, 254, 255, 256, 511, 1024} {
		for _, ending := range []string{"", "\n", "\r", "\r\n", "\n\r", "\r\r"} {
			cases = append(cases, strings.Repeat("x", n)+ending+"tail")
		}
	}
	alphabet := []byte{'a', '\r', '\n', '\t', 0x80}
	var enumerate func(string, int)
	enumerate = func(s string, n int) {
		cases = append(cases, s)
		if n == 0 {
			return
		}
		for _, b := range alphabet {
			enumerate(s+string([]byte{b}), n-1)
		}
	}
	enumerate("", 4)
	for _, input := range cases {
		text := input
		next := -1
		if i := strings.IndexAny(input, "\r\n"); i >= 0 {
			text = input[:i]
			next = i + 1
			if input[i] == '\r' && next < len(input) && input[next] == '\n' {
				next++
			}
		}
		got := legacy.PortTestMOTDLine(input)
		if got.Text != text || got.Next != next || !got.Guard || !got.InputUnchanged {
			t.Fatalf("input %q: got %+v; want text %q next %d and intact guards/input", input, got, text, next)
		}
	}
}
