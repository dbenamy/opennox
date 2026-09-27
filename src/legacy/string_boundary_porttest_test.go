//go:build porttest

package legacy

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"unicode/utf16"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestStringBoundaryBytes(t *testing.T) {
	if GoStringP(nil) != "" || GoStringNP(nil, 0) != "" {
		t.Fatal("nil/empty string changed")
	}
	for _, input := range []string{"", "a", "\x00", "a\x00tail", "\x80\xff", "héllo", strings.Repeat("xyz", 137)} {
		before := alloc.PortTestAllocationCount()
		p := CString(input)
		ptr := unsafe.Pointer(p)
		buf := unsafe.Slice((*byte)(ptr), len(input)+1)
		if !bytes.Equal(buf, append([]byte(input), 0)) {
			t.Fatalf("CString bytes for %q", input)
		}
		tracked := PortTestAllocationUsesTracker()
		if alloc.PortTestAllocationLive(ptr) != tracked {
			t.Fatal("CString allocation domain changed")
		}
		wantDelta := 0
		if tracked {
			wantDelta = 1
		}
		if alloc.PortTestAllocationCount()-before != wantDelta {
			t.Fatal("CString allocation count changed")
		}
		want, _, _ := strings.Cut(input, "\x00")
		got := GoStringP(ptr)
		if got != want {
			t.Fatalf("GoString=%q want %q", got, want)
		}
		for n := 0; n <= len(input)+1; n++ {
			expect := input
			if n < len(expect) {
				expect = expect[:n]
			}
			expect, _, _ = strings.Cut(expect, "\x00")
			if got := GoStringNP(ptr, n); got != expect {
				t.Fatalf("bounded string n=%d: %q want %q", n, got, expect)
			}
		}
		if len(input) != 0 {
			buf[0] ^= 0x55
		}
		if got != want {
			t.Fatal("GoString did not copy its input")
		}
		StrFree(p)
		if alloc.PortTestAllocationLive(ptr) || alloc.PortTestAllocationCount() != before {
			t.Fatal("CString release unbalanced")
		}
	}
}

func TestStringBoundaryCopy(t *testing.T) {
	for _, input := range []string{"", "abc", "a\x00bc", "\x80\xff", "héllo", "longer than destination"} {
		for n := 1; n <= 12; n++ {
			buf := bytes.Repeat([]byte{0xa5}, n+2)
			want := bytes.Clone(buf)
			count := len(input)
			if count > n-1 {
				count = n - 1
			}
			copy(want[1:1+count], input)
			want[1+count] = 0
			if got := StrCopyP(unsafe.Pointer(&buf[1]), n, input); got != count || !bytes.Equal(buf, want) {
				t.Fatalf("copy n=%d input=%q got=%d bytes=%x want=%x", n, input, got, buf, want)
			}
		}
	}
}

func TestStringBoundaryWide(t *testing.T) {
	for _, input := range []string{"", "a", "a\x00b", "λ界", "😀", "\xff", "a😀界z"} {
		before := alloc.PortTestAllocationCount()
		p, free := CWString(input)
		units := utf16.Encode([]rune(input))
		got := unsafe.Slice((*uint16)(unsafe.Pointer(p)), len(units)+1)
		for i, v := range append(units, 0) {
			if got[i] != v {
				t.Fatalf("wide unit %d = %04x want %04x", i, got[i], v)
			}
		}
		want, _, _ := strings.Cut(string([]rune(input)), "\x00")
		if actual := GoWStringP(unsafe.Pointer(p)); actual != want {
			t.Fatalf("wide=%q want %q", actual, want)
		}
		free()
		if alloc.PortTestAllocationCount() != before {
			t.Fatal("wide allocation release unbalanced")
		}
	}
	for _, units := range [][]uint16{{0}, {0xd800, 0}, {0xdc00, 0}, {0xd800, 'a', 0}, {0xd83d, 0xde00, 0}} {
		want := string(utf16.Decode(units[:len(units)-1]))
		if got := GoWStringP(unsafe.Pointer(&units[0])); got != want {
			t.Fatalf("wide malformed sequence %x = %q want %q", units, got, want)
		}
	}
}

func TestStringMallocFailure(t *testing.T) {
	const childEnv = "OPENNOX_PORT_STRING_MALLOC_CHILD"
	if os.Getenv(childEnv) == "1" {
		defer func() {
			if v := recover(); v != nil {
				fmt.Fprintln(os.Stdout, "recovered allocation failure:", v)
				os.Exit(23)
			}
		}()
		// This cannot fit in the qualified 32-bit address space; no large Go
		// allocation is attempted. Exercise only malloc's failure disposition.
		p := portTestStringMalloc(^uintptr(0))
		fmt.Fprintln(os.Stdout, "unexpected allocation return", p == nil)
		os.Exit(24)
	}
	cmd := exec.Command(os.Args[0], "-test.run=^TestStringMallocFailure$", "-test.count=1")
	cmd.Env = append(os.Environ(), childEnv+"=1", "ASAN_OPTIONS=allocator_may_return_null=1:detect_leaks=0")
	output, err := cmd.CombinedOutput()
	exit, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("allocation child did not exit with failure: %v", err)
	}
	wantCode, wantText := 2, "runtime: C malloc failed"
	if PortTestAllocationUsesTracker() {
		wantCode, wantText = 23, "recovered allocation failure: cannot allocate"
	}
	if exit.ExitCode() != wantCode || !bytes.Contains(output, []byte(wantText)) {
		t.Fatalf("allocation disposition exit=%d want=%d; expected marker %q; output: %.600s", exit.ExitCode(), wantCode, wantText, output)
	}
}
