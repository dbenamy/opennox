//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/client/gui"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
	"unsafe"
)

// Observe the existing dialog API boundary; the real list event handler,
// selected-name extraction, validation, rename and message queue remain active.
func TestTeamUIRenameDialogEvents(t *testing.T) {
	type record struct {
		Name          string
		Title, Prompt string
		Flags         uint32
		Result        int
		Team          string
		Rows          []string
		Queue         []legacy.PortTestShopPacketResult
	}
	var records []record
	for _, name := range []string{"Violet squad", "Blue team", "", "Team Ω"} {
		t.Run(fmt.Sprintf("name=%q", name), func(t *testing.T) {
			o := newTeamUIOwner(t)
			defer noxflags.PortTestGameFlags(1)()
			w := o.openPlayers(t)
			w.ChildByID(10502).Func94(&gui.RawEvent{Event: 16403, Arg1: 0})
			reset, snapshot, free := legacy.PortTestVisibilityReliableOwner()
			t.Cleanup(free)
			reset()
			r := record{Name: name}
			oldDialog, oldText := legacy.Nox_xxx_dialogMsgBoxCreate_449A10, legacy.Sub_449E60
			t.Cleanup(func() { legacy.Nox_xxx_dialogMsgBoxCreate_449A10 = oldDialog; legacy.Sub_449E60 = oldText })
			calls := 0
			legacy.Nox_xxx_dialogMsgBoxCreate_449A10 = func(parent *gui.Window, title, text string, flags gui.DialogFlags, ok, cancel func()) {
				calls++
				if parent != w || ok != nil || cancel != nil {
					t.Error("rename dialog owner/callbacks")
				}
				r.Title, r.Prompt, r.Flags = title, text, uint32(flags)
			}
			legacy.Sub_449E60 = func(code int8) int {
				if code != -88 {
					t.Errorf("dialog text selector %d", code)
				}
				return int(uintptr(unsafe.Pointer(alloc.InternCString16(name))))
			}
			legacy.PortTestTeamUIEvent(w, w.ChildByID(10509), 16391, 0)
			if calls != 1 || r.Title != "Rename team" || r.Prompt != "New team name" || r.Flags != 163 {
				t.Errorf("dialog: %+v calls=%d", r, calls)
			}
			button := o.c.GUI.NewWindowRaw(w, gui.StatusEnabled, 0, 0, 10, 10, nil)
			button.SetID(4001)
			r.Result = legacy.PortTestTeamUIEvent(w, button, 16391, 0)
			r.Team = o.c.srv.Teams.ByID(1).Name()
			r.Rows = teamUIRowNames(w.ChildByID(10502))
			r.Queue = snapshot()
			if name == "Blue team" {
				if r.Team != "Red team" {
					t.Errorf("invalid name accepted %q", r.Team)
				}
			} else if r.Team != name {
				t.Errorf("rename %q -> %q", name, r.Team)
			}
			records = append(records, r)
		})
	}
	spellbookCapture(t, "team-ui-rename-dialog-events", records, "a21b174552bdc7efb0dd4793cbcb09249b7d04e92a1784016a19feb4ed42ac4b")
}
