//go:build porttest

package legacy

import (
	"bytes"
	"encoding/json"
	"math"
	"os"
	"strconv"
	"testing"
	"time"
	"unsafe"
)

// Expected values were captured twice from the original 386 libc, independently
// of the replacement helpers. Preserve their exact integer and floating bits.
func TestLibcThemeNumbers(t *testing.T) {
	var cases []struct {
		Input     string `json:"input"`
		Integer   int32  `json:"integer"`
		FloatBits string `json:"float_bits"`
	}
	data, err := os.ReadFile("testdata/libc_theme_numbers.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &cases); err != nil {
		t.Fatal(err)
	}
	if len(cases) != 58 {
		t.Fatalf("capture cases = %d, want 58", len(cases))
	}
	for _, c := range cases {
		want, err := strconv.ParseUint(c.FloatBits, 16, 64)
		if err != nil {
			t.Fatal(err)
		}
		if got := mapThemeInt(c.Input); got != c.Integer {
			t.Errorf("integer %q = %d, want %d", c.Input, got, c.Integer)
		}
		if got := math.Float64bits(mapThemeFloat(c.Input)); got != want {
			t.Errorf("float %q = %016x, want %016x", c.Input, got, want)
		}
	}
}

func TestLibcBrowserAddresses(t *testing.T) {
	var cases []struct {
		Input string `json:"input"`
		Word  string `json:"word"`
	}
	data, err := os.ReadFile("testdata/libc_browser_addresses.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &cases); err != nil {
		t.Fatal(err)
	}
	if len(cases) != 5615 {
		t.Fatalf("capture cases = %d, want 5615", len(cases))
	}
	oldSelected, oldHas := browserUI.selected, browserUI.hasSelection
	t.Cleanup(func() { browserUI.selected, browserUI.hasSelection = oldSelected, oldHas })
	browserUI.selected, browserUI.hasSelection = nil, 0
	if got := uint32(nox_client_getServerAddr_43B300()); got != 0 {
		t.Fatalf("no selection = %08x, want zero", got)
	}
	for _, c := range cases {
		want, err := strconv.ParseUint(c.Word, 16, 32)
		if err != nil {
			t.Fatal(err)
		}
		// The real browser entrypoint reads its selected record at byte 12.
		// Own the complete string and guards, including bytes after embedded NULs.
		record := bytes.Repeat([]byte{0xa5}, 12+len(c.Input)+1+16)
		copy(record[12:], c.Input)
		record[12+len(c.Input)] = 0
		before := bytes.Clone(record)
		browserUI.selected, browserUI.hasSelection = unsafe.Pointer(&record[0]), 1
		if got := uint32(nox_client_getServerAddr_43B300()); got != uint32(want) {
			t.Errorf("address %q = %08x, want %08x", c.Input, got, want)
		}
		if !bytes.Equal(record, before) {
			t.Fatalf("address parser mutated record %q", c.Input)
		}
	}
}

func TestLibcThemeClockScope(t *testing.T) {
	defer themeObserve(false, 0)
	for _, epoch := range []uint32{0, 1, 0x7fffffff, 0x80000000, 0xffffffff} {
		themeObserve(true, epoch)
		if got := portTestThemeClock(); got != epoch {
			t.Fatalf("observed clock = %08x, want %08x", got, epoch)
		}
		next := epoch ^ 0x13579bdf
		themeObserve(true, next)
		if got := portTestThemeClock(); got != next {
			t.Fatalf("repeated activation = %08x, want %08x", got, next)
		}
		themeObserve(false, 0)
		themeObserve(false, 0)
		// The system clock is live again. Allow two seconds of wall-clock drift;
		// every injected value is deliberately far from the current epoch.
		delta := int64(int32(portTestThemeClock() - uint32(time.Now().Unix())))
		if delta < -2 || delta > 2 {
			t.Fatalf("clock not restored after idempotent cleanup: delta=%d", delta)
		}
	}
}
