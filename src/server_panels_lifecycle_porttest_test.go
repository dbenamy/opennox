//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy"
	"os"
	"os/exec"
	"testing"
	"unsafe"
)

func TestServerPanelsMissingResource(t *testing.T) {
	kind := os.Getenv("OPENNOX_SERVER_PANELS_RESOURCE_CHILD")
	if kind == "" {
		for _, kind := range []string{"weapon", "armor", "spell", "access", "general", "advanced", "advserv"} {
			t.Run(kind, func(t *testing.T) {
				cmd := exec.Command(os.Args[0], "-test.run=^TestServerPanelsMissingResource$", "-test.count=1")
				cmd.Env = append(os.Environ(), "OPENNOX_SERVER_PANELS_RESOURCE_CHILD="+kind)
				if out, err := cmd.CombinedOutput(); err != nil {
					t.Fatalf("%s resource child: %v\n%s", kind, err, out)
				}
			})
		}
		return
	}
	o := newServerOptionsOwner(t)
	objectWords, restore := legacy.PortTestServerPanelsObjectWords()
	t.Cleanup(restore)
	o.installSubpanels(t, true)
	var root *uint32
	switch kind {
	case "weapon", "armor":
		root = objectWords["root"]
	case "spell":
		root = o.optionWords["panel-1045484"]
	case "access":
		root = o.optionWords["panel-1045516"]
	case "general":
		root = o.optionWords["panel-1309812"]
	case "advanced":
		root = o.optionWords["panel-1316708"]
	case "advserv":
		root = o.optionWords["panel-1316972"]
	default:
		t.Fatal(kind)
	}
	old := legacy.Nox_new_window_from_file
	legacy.Nox_new_window_from_file = func(string, gui.WindowFunc) *gui.Window { return nil }
	t.Cleanup(func() { legacy.Nox_new_window_from_file = old })
	if got := legacy.PortTestServerPanelsConstruct(kind, o.options, unsafe.Pointer(&o.settings[0])); got != 0 {
		t.Fatalf("missing resource returned %d", got)
	}
	if *root != 0 {
		t.Fatal("missing resource retained root")
	}
	legacy.Nox_new_window_from_file = old
	legacy.PortTestServerPanelsConstruct(kind, o.options, unsafe.Pointer(&o.settings[0]))
	if *root == 0 {
		t.Fatal("retry failed to create root")
	}
	// Standalone advanced owners are not children of the main options window.
	if kind == "advanced" || kind == "advserv" {
		w := (*gui.Window)(unsafe.Pointer(uintptr(*root)))
		w.Destroy()
		*root = 0
	}
}
