//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/console"
	"github.com/opennox/opennox/v1/legacy"
	"reflect"
	"testing"
	"unicode/utf16"
)

func consoleCommandWide(t *testing.T, base, off uintptr, s string, n int) {
	t.Helper()
	buf := serverConfigOwnBytes(t, base, off, n)
	clear(buf)
	for i, c := range append(utf16.Encode([]rune(s)), 0) {
		buf[2*i] = byte(c)
		buf[2*i+1] = byte(c >> 8)
	}
}
func TestConsoleCommandsListUsers(t *testing.T) {
	type row struct {
		Client bool
		Status uint32
		Output []consoleCommandLine
	}
	var rows []row
	for _, client := range []bool{false, true} {
		for _, status := range []uint32{0, 4, 8, 12} {
			t.Run(fmt.Sprintf("%t/%d", client, status), func(t *testing.T) {
				o := newConsoleCommandOwner(t)
				for i := range o.players {
					o.players[i].Active = 0
				}
				consoleCommandWide(t, 0x587000, 106604, "%s", 8)
				p := &o.players[0]
				p.Active = 1
				p.PlayerInd = 0
				p.SetName("Alpha Ω")
				p.Field3680 = status
				want := "Alpha Ω"
				if !client && status&4 != 0 {
					want += ", server muted"
				}
				if status&8 != 0 {
					want += ", client muted"
				}
				if !o.call(t, "list users", client) {
					t.Fatal("list return")
				}
				expected := []consoleCommandLine{{console.ColorRed, "users"}, {console.ColorRed, want}}
				if !reflect.DeepEqual(o.printer.lines, expected) {
					t.Fatal(o.printer.lines, expected)
				}
				rows = append(rows, row{client, status, o.printer.lines})
			})
		}
	}
	spellbookCapture(t, "console-commands-list-users", rows, "899ef38fced123931c386587c9681a2ee55e93144c71a82bc0841fbf812d6cd0")
}
func TestConsoleCommandsListMaps(t *testing.T) {
	type row struct {
		Count  int
		Output []consoleCommandLine
	}
	var rows []row
	for _, count := range []int{0, 1, 3, 4, 5, 8, 9} {
		t.Run(fmt.Sprint(count), func(t *testing.T) {
			o := newConsoleCommandOwner(t)
			f := legacy.PortTestMapCatalogOpen(1)
			t.Cleanup(f.Close)
			serverConfigOwnBytes(t, 0x5D4594, 822404, 256)
			consoleCommandWide(t, 0x587000, 103276, "%s", 8)
			consoleCommandWide(t, 0x587000, 103284, "%s", 8)
			suffix := serverConfigOwnBytes(t, 0x587000, 103292, 5)
			copy(suffix, []byte(".map\x00"))
			for i := count - 1; i >= 0; i-- {
				f.Add(legacy.PortTestMapCatalogEntry{Name: fmt.Sprintf("map%02d", i), Enabled: int32(i % 2)})
			}
			if !o.call(t, "list maps", false) {
				t.Fatal("map-list return")
			}
			var expected []consoleCommandLine
			// The current C formatter ignores width/precision on %S; preserve its output.
			line := ""
			for i := 0; i < count; i++ {
				name := fmt.Sprintf("map%02d.map", i)
				line += name + "\t\t"
				if i%4 == 3 || i == count-1 {
					expected = append(expected, consoleCommandLine{console.ColorRed, line})
					line = ""
				}
			}
			if !reflect.DeepEqual(o.printer.lines, expected) {
				t.Fatal(o.printer.lines, expected)
			}
			rows = append(rows, row{count, o.printer.lines})
		})
	}
	spellbookCapture(t, "console-commands-list-maps", rows, "b51a85c3e3ce5a5562dd6706e1fddb2729289e065632ae2a82dd0e38ff4fda2c")
}
