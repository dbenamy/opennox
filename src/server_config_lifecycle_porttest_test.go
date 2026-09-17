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

func TestServerConfigRuleMissingResource(t *testing.T) {
	if os.Getenv("OPENNOX_SERVER_CONFIG_RESOURCE_CHILD") == "" {
		cmd := exec.Command(os.Args[0], "-test.run=^TestServerConfigRuleMissingResource$", "-test.count=1")
		cmd.Env = append(os.Environ(), "OPENNOX_SERVER_CONFIG_RESOURCE_CHILD=1")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("rule resource child: %v\n%s", err, out)
		}
		return
	}
	o := newServerOptionsOwner(t)
	o.installSubpanels(t, true)
	old := legacy.Nox_new_window_from_file
	legacy.Nox_new_window_from_file = func(string, gui.WindowFunc) *gui.Window { return nil }
	t.Cleanup(func() { legacy.Nox_new_window_from_file = old })
	if got := legacy.PortTestServerConfigRuleOpen(o.options, unsafe.Pointer(&o.settings[0])); got != 0 {
		t.Fatal("missing resource return", got)
	}
	if *o.optionWords["panel-1523024"] != 0 {
		t.Fatal("missing resource retained root")
	}
	legacy.Nox_new_window_from_file = old
	if legacy.PortTestServerConfigRuleOpen(o.options, unsafe.Pointer(&o.settings[0])) == 0 {
		t.Fatal("retry failed")
	}
}
