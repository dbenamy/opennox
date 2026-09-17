//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/client/gui"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"golang.org/x/image/font/basicfont"
	"os"
	"os/exec"
	"testing"
	"unsafe"
)

// Isolate the original constructor's missing-resource path so an invalid C access
// reports a single failing contract without preventing the remaining contracts.
func TestServerOptionsMissingResource(t *testing.T) {
	if os.Getenv("OPENNOX_SERVER_OPTIONS_RESOURCE_CHILD") != "1" {
		cmd := exec.Command(os.Args[0], "-test.run=^TestServerOptionsMissingResource$", "-test.count=1")
		cmd.Env = append(os.Environ(), "OPENNOX_SERVER_OPTIONS_RESOURCE_CHILD=1")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("missing resource child: %v\n%s", err, out)
		}
		return
	}
	o := newServerOptionsOwner(t)
	*o.optionWords["root"] = 0
	old := legacy.Nox_new_window_from_file
	legacy.Nox_new_window_from_file = func(string, gui.WindowFunc) *gui.Window { return nil }
	t.Cleanup(func() { legacy.Nox_new_window_from_file = old })
	if got := o.call("construct", 0, ""); got != 0 {
		t.Fatalf("failed resource returned %d", got)
	}
	if *o.optionWords["root"] != 0 {
		t.Fatal("failed resource retained root")
	}
	legacy.Nox_new_window_from_file = old
	loads := o.installConstructor(t)
	if got := o.call("construct", 0, ""); got != 1 || len(*loads) != 1 {
		t.Fatalf("resource retry returned %d, loads %v", got, *loads)
	}
	o.call("close", 0, "")
	o.options = nil
}

func TestServerOptionsClientLifecycle(t *testing.T) {
	for lang := 0; lang < 9; lang++ {
		for _, height := range []int{8, 10, 11, 13} {
			t.Run(fmt.Sprintf("%d-%d", lang, height), func(t *testing.T) {
				catalog := legacy.PortTestMapCatalogOpen(3)
				t.Cleanup(catalog.Close)
				o := newServerOptionsOwner(t)
				defer noxflags.PortTestGameFlags(0)()
				o.configureLanguage(lang)
				face := *basicfont.Face7x13
				face.Height = height
				face.Ascent = height
				t.Cleanup(o.c.Render().GetFonts().PortTestDefaultFont(&face))
				_, restore := o.c.Render().GetFonts().PortTestWindowFont(&face, "small", "large")
				t.Cleanup(restore)
				loads := o.installConstructor(t)
				for repeat := 0; repeat < 2; repeat++ {
					if got := o.call("construct", 0, ""); got != 1 || *o.optionWords["root"] == 0 {
						t.Fatalf("construct %d", got)
					}
					wantLang := lang
					if height > 10 {
						wantLang = 2
					}
					if got := (*loads)[len(*loads)-1]; got != fmt.Sprintf("server-options-%d.wnd", wantLang) {
						t.Fatalf("resource %s", got)
					}
					if o.options.ChildByID(10119).DrawData().Text() != "Arena battle" {
						t.Fatalf("mode title %q", o.options.ChildByID(10119).DrawData().Text())
					}
					o.call("close", 0, "")
					if *o.optionWords["root"] != 0 {
						t.Fatal("root survived close")
					}
					o.options = nil
					o.c.GUI.FreeDestroyed()
				}
				if len(*loads) != 2 {
					t.Fatal("resource count")
				}
			})
		}
	}
}

func TestServerOptionsKeyboardClose(t *testing.T) {
	type row struct {
		Event, Key, State, Result int
		Open                      bool
	}
	var rows []row
	for _, event := range []int{0, 21, 22, 16391} {
		for _, key := range []int{0, 1, 2, 255} {
			for _, state := range []int{0, 1, 2} {
				t.Run(fmt.Sprintf("%d-%d-%d", event, key, state), func(t *testing.T) {
					o := newServerOptionsOwner(t)
					got := legacy.PortTestServerOptionsKey(o.options, event, key, state)
					want := 1
					if event == 21 && key != 1 {
						want = 0
					}
					closed := event == 21 && key == 1 && state == 2
					if got != want || (*o.optionWords["root"] == 0) != closed {
						t.Fatalf("key return %d root %x", got, *o.optionWords["root"])
					}
					rows = append(rows, row{event, key, state, got, *o.optionWords["root"] != 0})
				})
			}
		}
	}
	spellbookCapture(t, "server-options-keyboard", rows, "877b54ddc7b0f91953c7f3499c957785eff294da345be3db0fbb7b64c68e68de")
}

func TestServerOptionsHostConstruction(t *testing.T) {
	for _, mode := range []int{0x100, 0x20, 0x400, 0x1000, 0x80} {
		t.Run(fmt.Sprintf("%x", mode), func(t *testing.T) {
			catalog := legacy.PortTestMapCatalogOpen(13)
			t.Cleanup(catalog.Close)
			o := newServerOptionsOwner(t)
			defer noxflags.PortTestGameFlags(noxflags.GameFlag(1 | mode))()
			// Exercise synchronized host data as well as the retained mode bits.
			binary.LittleEndian.PutUint16(o.settings[52:], uint16(mode))
			copy(unsafe.Slice(memmap.PtrUint8(0x5D4594, 371438), 58), o.settings)
			loads := o.installConstructor(t)
			if got := o.call("construct", 0, ""); got != 1 || *o.optionWords["root"] == 0 {
				t.Fatalf("host construct %d", got)
			}
			if *memmap.PtrUint32(0x5D4594, 371688) != 1 {
				t.Fatal("host did not select working settings")
			}
			if o.call("dirty-get", 0, "") != 1 {
				t.Fatal("host settings were not marked dirty")
			}
			if got := o.call("construct", 0, ""); got != 1 || *o.optionWords["root"] != 0 {
				t.Fatal("repeat open did not close")
			}
			o.options = nil
			o.c.GUI.FreeDestroyed()
			if len(*loads) != 1 {
				t.Fatal("repeat open reparsed resource")
			}
			if got := o.call("construct", 0, ""); got != 1 || len(*loads) != 2 {
				t.Fatal("host reopen")
			}
		})
	}
}
