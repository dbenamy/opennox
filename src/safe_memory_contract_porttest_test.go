//go:build safe && porttest

package opennox

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

type safeMemoryCase struct {
	Name string
	Spec legacy.PortTestSafeMemorySpec
}

type safeMemoryCapture struct {
	Name   string
	Spec   legacy.PortTestSafeMemorySpec
	Result legacy.PortTestSafeMemoryResult
}

func safeBytes(n int, seed byte) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = seed + byte(i*37)
	}
	return b
}

func safeCString(body []byte) []byte {
	return append(append([]byte(nil), body...), 0)
}

func safeBuffer(fill byte, payload []byte, off int) []byte {
	b := bytes.Repeat([]byte{fill}, 96)
	copy(b[16+off:], payload)
	return b
}

func safeCStringLen(b []byte) int {
	for i, v := range b {
		if v == 0 {
			return i
		}
	}
	panic("test case has unterminated string")
}

func safeCompare(a, b []byte, n int, strings bool) int {
	if strings {
		n = safeCStringLen(a)
		m := safeCStringLen(b)
		if m < n {
			n = m
		}
		// Compare the common prefix, then the terminating NUL when one
		// string is a prefix of the other.
		for i := 0; i < n; i++ {
			if a[i] < b[i] {
				return -1
			}
			if a[i] > b[i] {
				return 1
			}
		}
		if safeCStringLen(a) < safeCStringLen(b) {
			return -1
		}
		if safeCStringLen(a) > safeCStringLen(b) {
			return 1
		}
		return 0
	}
	for i := 0; i < n; i++ {
		if a[i] < b[i] {
			return -1
		}
		if a[i] > b[i] {
			return 1
		}
	}
	return 0
}

func safeMemoryCases() []safeMemoryCase {
	var out []safeMemoryCase
	add := func(name, op string, left, right []byte, lo, ro, size int) {
		if lo < 0 || ro < 0 || lo > 63 || ro > 63 || len(left) > 64-lo || len(right) > 64-ro {
			panic("invalid safe-memory fixture bounds")
		}
		out = append(out, safeMemoryCase{Name: name, Spec: legacy.PortTestSafeMemorySpec{
			Op: op, Left: left, Right: right, LeftOff: lo, RightOff: ro, Size: size,
		}})
	}

	// Exhaust every memcpy size with a deterministic non-text byte pattern.
	for n := 0; n <= 64; n++ {
		add(fmt.Sprintf("memcpy-size-%02d", n), "memcpy", safeBytes(n, 0x83), safeBytes(n, 0x11), 0, 0, n)
	}
	// Misaligned addresses at useful lengths, including exact remaining capacity.
	for _, off := range []int{1, 3} {
		for _, n := range []int{0, 1, 7, 31, 64 - off} {
			add(fmt.Sprintf("memcpy-off-%d-size-%d", off, n), "memcpy", safeBytes(n, 0xa0), safeBytes(n, 0x80), off, off, n)
		}
	}

	for _, offs := range [][2]int{{1, 3}, {3, 1}} {
		add(fmt.Sprintf("memcpy-mixed-offset-%d-%d", offs[0], offs[1]), "memcpy", safeBytes(61, 0x77), safeBytes(61, 0xc3), offs[0], offs[1], 61)
	}
	for _, n := range []int{1, 2, 4, 7, 16, 31, 60} {
		a := safeBytes(n+1, 0x77)
		add(fmt.Sprintf("memcmp-bounded-%d", n), "memcmp", a, bytes.Clone(a), 1, 3, n)
		seen := make(map[int]bool)
		for _, pos := range []int{0, n / 2, n - 1, n} {
			if seen[pos] {
				continue
			}
			seen[pos] = true
			b := bytes.Clone(a)
			b[pos]++
			add(fmt.Sprintf("memcmp-span-%d-diff-%d", n, pos), "memcmp", a, b, 1, 3, n)
		}
	}

	for _, tc := range []struct {
		name       string
		a, b       []byte
		offA, offB int
		n          int
	}{
		{"memcmp-equal", []byte{0, 0x80, 0xff, 3}, []byte{0, 0x80, 0xff, 3}, 0, 0, 4},
		{"memcmp-prefix", []byte("sameX"), []byte("sameY"), 1, 3, 4},
		{"memcmp-low-high", []byte{0x7f, 0}, []byte{0x80, 0}, 3, 1, 1},
		{"memcmp-high-low", []byte{0x80, 0}, []byte{0x7f, 0}, 1, 3, 1},
		{"memcmp-zero-size", []byte{1}, []byte{2}, 3, 1, 0},
	} {
		add(tc.name, "memcmp", tc.a, tc.b, tc.offA, tc.offB, tc.n)
	}

	for _, tc := range []struct {
		name string
		s    []byte
		off  int
	}{
		{"strlen-empty", []byte{0}, 0},
		{"strlen-last-payload-byte", []byte{0}, 63},
		{"strlen-full-payload", safeCString(bytes.Repeat([]byte{0xff}, 63)), 0},
		{"strlen-embedded-nul", []byte{'A', 0, 0x80, 0xff, 0}, 1},
		{"strlen-high-byte", []byte{0x80, 0}, 3},
		{"strlen-near-capacity", safeCString(bytes.Repeat([]byte{0x7f}, 60)), 3},
	} {
		add(tc.name, "strlen", tc.s, nil, tc.off, 0, 0)
	}

	for _, tc := range []struct {
		name string
		s    []byte
		off  int
	}{
		{"strcpy-empty", []byte{0}, 0},
		{"strcpy-text", safeCString([]byte("copy")), 1},
		{"strcpy-embedded-nul", []byte{'A', 0, 0x80, 0xff, 0}, 3},
		{"strcpy-high-byte", []byte{0x80, 0xff, 0}, 1},
		{"strcpy-near-capacity", safeCString(bytes.Repeat([]byte{0x80}, 60)), 3},
	} {
		add(tc.name, "strcpy", []byte("keep-tail"), tc.s, tc.off, 1, 0)
	}

	for _, tc := range []struct {
		name string
		dst  []byte
		src  []byte
		do   int
		so   int
	}{
		{"strcat-empty-src", safeCString([]byte("A")), []byte{0}, 1, 3},
		{"strcat-empty-dst", []byte{0}, safeCString([]byte("B")), 3, 1},
		{"strcat-ordinary", safeCString([]byte("left")), safeCString([]byte("right")), 1, 3},
		{"strcat-embedded-nul", []byte{'A', 0, 'x', 0}, []byte{'B', 0, 0x80, 0}, 3, 1},
		{"strcat-high-byte", safeCString([]byte{0x80}), []byte{0xff, 0}, 1, 3},
		{"strcat-exact-capacity", safeCString(bytes.Repeat([]byte{'L'}, 58)), []byte{'R', 'S', 0}, 3, 1},
		{"strcat-both-empty", []byte{0}, []byte{0}, 0, 0},
	} {
		add(tc.name, "strcat", tc.dst, tc.src, tc.do, tc.so, 0)
	}

	for _, tc := range []struct {
		name   string
		a, b   []byte
		ao, bo int
	}{
		{"strcmp-equal", safeCString([]byte("same")), safeCString([]byte("same")), 0, 3},
		{"strcmp-empty-equal", []byte{0}, []byte{0}, 63, 63},
		{"strcmp-empty-first", []byte{0}, []byte{0x80, 0}, 1, 3},
		{"strcmp-empty-second", []byte{0xff, 0}, []byte{0}, 3, 1},
		{"strcmp-long-equal", safeCString(bytes.Repeat([]byte{0xff}, 60)), safeCString(bytes.Repeat([]byte{0xff}, 60)), 3, 3},
		{"strcmp-prefix", safeCString([]byte("abc")), safeCString([]byte("abcd")), 1, 0},
		{"strcmp-difference", safeCString([]byte("abx")), safeCString([]byte("aby")), 3, 1},
		{"strcmp-low-high", []byte{0x7f, 0}, []byte{0x80, 0}, 1, 3},
		{"strcmp-high-low", []byte{0x80, 0}, []byte{0x7f, 0}, 3, 1},
		{"strcmp-embedded-nul", []byte{'a', 0, 'x', 0}, []byte{'a', 0, 'y', 0}, 1, 3},
	} {
		add(tc.name, "strcmp", tc.a, tc.b, tc.ao, tc.bo, 0)
	}
	return out
}

func expectedSafeMemory(c safeMemoryCase) (left, right []byte, ret int, dest bool) {
	s := c.Spec
	left = safeBuffer(0xa5, s.Left, s.LeftOff)
	right = safeBuffer(0x5a, s.Right, s.RightOff)
	li, ri := 16+s.LeftOff, 16+s.RightOff
	switch s.Op {
	case "memcpy":
		copy(left[li:li+s.Size], right[ri:ri+s.Size])
		dest = true
	case "memcmp":
		ret = safeCompare(s.Left, s.Right, s.Size, false)
	case "strlen":
		ret = safeCStringLen(s.Left)
	case "strcpy":
		n := safeCStringLen(s.Right) + 1
		copy(left[li:li+n], right[ri:ri+n])
		dest = true
	case "strcat":
		dn := safeCStringLen(s.Left)
		sn := safeCStringLen(s.Right)
		copy(left[li+dn:li+dn+sn+1], right[ri:ri+sn+1])
		dest = true
	case "strcmp":
		ret = safeCompare(s.Left, s.Right, 0, true)
	default:
		panic("unknown safe-memory operation")
	}
	return left, right, ret, dest
}

func TestSafeMemoryBridges(t *testing.T) {
	const want = "2b34e573e3b8f5b62ced16258ec2de832284fbe9d0cecc3a0c5120edf22eb59a" // Freeze from original bridge after comparison-sign normalization.
	var captures []safeMemoryCapture
	for _, c := range safeMemoryCases() {
		wantLeft, wantRight, wantRet, wantDest := expectedSafeMemory(c)
		got := legacy.PortTestSafeMemory(c.Spec)
		if !bytes.Equal(got.Left, wantLeft) {
			t.Errorf("%s: left mutation/guards differ", c.Name)
		}
		if !bytes.Equal(got.Right, wantRight) {
			t.Errorf("%s: source mutation/guards differ", c.Name)
		}
		if wantRet == 0 && (c.Spec.Op == "memcmp" || c.Spec.Op == "strcmp") {
			if got.Return != 0 {
				t.Errorf("%s: compare got %d, want 0", c.Name, got.Return)
			}
		} else if c.Spec.Op == "memcmp" || c.Spec.Op == "strcmp" {
			if wantRet < 0 && got.Return >= 0 || wantRet > 0 && got.Return <= 0 {
				t.Errorf("%s: compare got %d, want sign %d", c.Name, got.Return, wantRet)
			}
		} else if got.Return != int32(wantRet) {
			t.Errorf("%s: return got %d, want %d", c.Name, got.Return, wantRet)
		}
		if got.ReturnedDest != wantDest {
			t.Errorf("%s: returned-destination got %v, want %v", c.Name, got.ReturnedDest, wantDest)
		}
		// libc comparison magnitudes vary with implementation; the contract
		// and engine consumers use only negative, zero or positive ordering.
		if c.Spec.Op == "memcmp" || c.Spec.Op == "strcmp" {
			if got.Return < 0 {
				got.Return = -1
			} else if got.Return > 0 {
				got.Return = 1
			}
		}
		captures = append(captures, safeMemoryCapture{Name: c.Name, Spec: c.Spec, Result: got})
	}
	data, err := json.Marshal(captures)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(data)
	hash := hex.EncodeToString(sum[:])
	t.Logf("safe memory capture sha256=%s cases=%d", hash, len(captures))
	if want != "" && hash != want {
		t.Fatalf("safe memory capture hash got %s want %s", hash, want)
	}
	if path := os.Getenv("OPENNOX_SAFE_BRIDGES_CAPTURE"); path != "" {
		if err := os.WriteFile(path, data, 0o600); err != nil {
			t.Fatal(err)
		}
	}
}
