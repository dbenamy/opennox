//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"reflect"
	"testing"
)

func TestConsoleCommandsWindow(t *testing.T) {
	o := newConsoleCommandOwner(t)
	var calls [][2]int
	old := legacy.Nox_draw_setCutSize_476700
	legacy.Nox_draw_setCutSize_476700 = func(a, b int) { calls = append(calls, [2]int{a, b}) }
	t.Cleanup(func() { legacy.Nox_draw_setCutSize_476700 = old })
	for _, s := range []string{"75", "+10", "-5", "  +7", "-0", "invalid", "2147483648"} {
		if !o.call(t, "window", false, s) {
			t.Fatal("window return")
		}
	}
	if !o.call(t, "window", false) {
		t.Fatal("empty window return")
	}
	want := [][2]int{{75, 0}, {0, 10}, {0, -5}, {7, 0}, {0, 0}, {0, 0}, {2147483647, 0}}
	if !reflect.DeepEqual(calls, want) {
		t.Fatal(calls, want)
	}
	spellbookCapture(t, "console-commands-window", calls, "eca821348dda99d925d23fd31c196fe46562465447b0780ece2407b9152636d4")
}
