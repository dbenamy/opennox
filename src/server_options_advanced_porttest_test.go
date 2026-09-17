//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/ifs"
	"github.com/opennox/opennox/v1/client/gui"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
	"testing"
	"unsafe"
)

func TestServerOptionsAdvancedModal(t *testing.T) {
	type row struct {
		Host                  bool
		Loads                 []string
		MainOpen, ModalClosed bool
	}
	var rows []row
	for _, host := range []bool{false, true} {
		t.Run(fmt.Sprint(host), func(t *testing.T) {
			o := newServerOptionsOwner(t)
			flags := noxflags.GameFlag(0)
			if host {
				flags = 1
			}
			defer noxflags.PortTestGameFlags(flags)()
			loads := o.installSubpanels(t)
			t.Cleanup(handles.PortTestInit())
			oldDir, err := ifs.Workdir()
			if err != nil {
				t.Fatal(err)
			}
			if err := ifs.Chdir(t.TempDir()); err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() {
				if err := ifs.Chdir(oldDir); err != nil {
					t.Error(err)
				}
			})
			if got := legacy.PortTestServerOptionsEvent(o.options, o.options.ChildByID(10152), 16391, 0); got != 1 {
				t.Fatal("advanced button return")
			}
			modal := (*gui.Window)(unsafe.Pointer(uintptr(*o.optionWords["panel-1316708"])))
			if modal == nil || modal == o.options {
				t.Fatal("advanced modal ownership")
			}
			want := "spelllst.wnd"
			if host {
				want = "rulelist.wnd"
			}
			if len(*loads) != 2 || (*loads)[0] != "advanced.wnd" || (*loads)[1] != want {
				t.Fatalf("advanced resources %v", *loads)
			}
			button := modal.ChildByID(10148)
			modal.Func94(gui.AsWindowEvent(16391, uintptr(button.C()), 0))
			if *o.optionWords["panel-1316708"] != 0 || *o.optionWords["root"] == 0 {
				t.Fatal("advanced modal close changed main owner")
			}
			rows = append(rows, row{host, append([]string(nil), (*loads)...), true, true})
		})
	}
	spellbookCapture(t, "server-options-advanced-modal", rows, "82e4dc7f9bd631652dd96a26a540525c06f7f3f7dfd5e1e6634b78e11a9972d6")
}
