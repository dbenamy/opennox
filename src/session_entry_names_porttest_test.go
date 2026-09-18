//go:build porttest

package opennox

import (
	"bytes"
	"github.com/opennox/opennox/v1/legacy"
	"strings"
	"testing"
)

func TestSessionEntryMapNames(t *testing.T) {
	filename := serverConfigOwnBytes(t, 0x5D4594, 2598188, 80)
	basename := serverConfigOwnBytes(t, 0x85B3FC, 36, 80)
	selected := serverConfigOwnBytes(t, 0x5D4594, 3452, 12)
	words, restore := legacy.PortTestServerOptionsWords()
	defer restore()
	type row struct {
		Input, Filename, Basename string
		Changed, Dirty            uint32
	}
	var rows []row
	cases := []struct {
		input, file, base string
		changed           uint32
	}{
		{"arena.map", "arena.map", "arena", 1},
		{"ARENA.MAP", "arena.map", "arena", 0},
		{`maps\Arena\Arena.map`, `maps\Arena\Arena.map`, "Arena", 1},
		{`maps\other\next.nxz`, `maps\other\next.nxz`, "next", 1},
		{`maps\x`, `maps\x`, "", 1},
		{`maps\`, `maps\`, "", 1},
		{"x", "x", "", 1},
		{"", "", "", 1},
		{"", "", "", 0},
		{".map", ".map", "", 1},
		{"four", "four", "", 1},
		{strings.Repeat("a", 90) + ".map", strings.Repeat("a", 79), strings.Repeat("a", 75), 1},
		{strings.Repeat("p", 90) + `\short.map`, strings.Repeat("p", 79), "short", 1},
	}
	clear(filename)
	clear(basename)
	for _, c := range cases {
		*words["settings-updated"] = 0
		_, changed := legacy.PortTestSessionEntryMapText("set-path", c.input)
		file, _ := legacy.PortTestSessionEntryMapText("filename", "")
		base, _ := legacy.PortTestSessionEntryMapText("basename", "")
		got := row{c.input, file, base, changed, *words["settings-updated"]}
		want := row{c.input, c.file, c.base, c.changed, c.changed}
		if got != want {
			t.Errorf("map name transition: got%+v want%+v", got, want)
		}
		rows = append(rows, got)
	}
	spellbookCapture(t, "session-entry-map-names", rows, "ee9637d45867c33fc5365ce3e4e985393a8be7db0d0080a3773b1833ecccace2")
	type selectedRow struct {
		Input, Output string
		Count         uint32
		Bytes         []byte
	}
	var sels []selectedRow
	for _, value := range []string{"", "arena", "12345678", "x", ""} {
		// Observe the terminator and untouched tail, rather than assuming zero padding.
		copy(selected, bytes.Repeat([]byte{0xa5}, len(selected)))
		_, count := legacy.PortTestSessionEntryMapText("set-selected", value)
		out, _ := legacy.PortTestSessionEntryMapText("selected", "")
		want := bytes.Repeat([]byte{0xa5}, len(selected))
		copy(want, value)
		want[len(value)] = 0
		if out != value || count != uint32(len(value)+1) || !bytes.Equal(selected, want) {
			t.Errorf("selected map %q: output%q count%d bytes%x", value, out, count, selected)
		}
		sels = append(sels, selectedRow{value, out, count, bytes.Clone(selected)})
	}
	spellbookCapture(t, "session-entry-selected-map", sels, "246dc13a4293d7a2764fe166e77137d4b948900e1d64a1fe974341fbb437479e")
}
