//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestServerOptionsModeNames(t *testing.T) {
	o := newEntryOwner(t)
	titles := []struct {
		key, title string
		mode       uint16
	}{
		{"guiserv.c:CTF", "Capture flag", 0x20},
		{"guiserv.c:Arena", "Arena battle", 0x100},
		{"guiserv.c:Highlander", "Last survivor", 0x400},
		{"guiserv.c:KotR", "King of the realm", 0x10},
		{"guiserv.c:Flagball", "Flag ball", 0x40},
		{"guiserv.c:Quest", "Quest adventure", 0x1000},
		{"Noxworld.c:Chat", "Chat lobby", 0x80},
	}
	var entries []strman.Entry
	for _, x := range titles {
		entries = append(entries, strman.Entry{ID: strman.ID(x.key), Vals: []strman.Variant{{Str: x.title}}})
	}
	language, restore := o.c.srv.PortTestMeterStrings(entries...)
	t.Cleanup(restore)
	language(0)
	t.Cleanup(legacy.PortTestServerOptionsModes())
	if legacy.PortTestServerOptionsModeLoaded() {
		t.Fatal("owner retained mode cache")
	}
	// Exhaust the actual 16-bit ABI input while using explicit independent table expectations.
	for bits := 0; bits < 65536; bits++ {
		want := "Arena battle"
		for _, x := range titles {
			if uint16(bits)&0x17f0 == x.mode {
				want = x.title
				break
			}
		}
		got := legacy.PortTestServerOptionsModeName(uint16(bits))
		if got != want {
			t.Fatalf("mode %04x: %q, want %q", bits, got, want)
		}
	}
	if !legacy.PortTestServerOptionsModeLoaded() {
		t.Fatal("lookup did not initialize mode cache")
	}
	for _, x := range titles {
		t.Run(fmt.Sprintf("reverse-%04x", x.mode), func(t *testing.T) {
			if got := legacy.PortTestServerOptionsModeFromName(x.title); got != int(x.mode) {
				t.Fatalf("reverse: %x", got)
			}
		})
	}
	for _, name := range []string{"", "arena battle", "Arena battle ", "Unknown", "Квест"} {
		if got := legacy.PortTestServerOptionsModeFromName(name); got != 0 {
			t.Fatalf("unknown title %q: %x", name, got)
		}
	}
}
