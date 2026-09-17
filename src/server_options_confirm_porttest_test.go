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

func TestServerOptionsTeamConfirmation(t *testing.T) {
	type row struct {
		Flags, Mode, Before, Gameplay int
		Dialog                        uint32
		Prompt                        string
		Focus, After                  int
		ResultFlags                   uint
	}
	var rows []row
	for _, flags := range []int{0, 128, 144} {
		for _, mode := range []int{0x20, 0x40, 0x100, 0x10} {
			for _, count := range []int{0, 1, 2, 3, 4} {
				for _, gameplay := range []int{0, 4} {
					t.Run(fmt.Sprintf("%x-%x-%d-%d", flags, mode, count, gameplay), func(t *testing.T) {
						o := newServerOptionsOwner(t)
						t.Cleanup(o.c.srv.PortTestMapDrawableTeamMessages())
						clear(o.c.srv.Teams.Arr)
						o.c.srv.Teams.ActiveCnt = 0
						seed := noxflags.PortTestGameFlags(noxflags.GameModeCoopTeam)
						for i := 0; i < count; i++ {
							if o.c.srv.Teams.Create(0) == nil {
								t.Fatal("seed team")
							}
						}
						seed()
						defer noxflags.PortTestGameFlags(noxflags.GameFlag(flags))()
						oldGameplay := noxflags.GetGamePlay()
						noxflags.UnsetGamePlay(^noxflags.GameplayFlag(0))
						noxflags.SetGamePlay(noxflags.GameplayFlag(gameplay))
						t.Cleanup(func() { noxflags.UnsetGamePlay(^noxflags.GameplayFlag(0)); noxflags.SetGamePlay(oldGameplay) })
						name := legacy.PortTestServerOptionsModeName(uint16(mode))
						o.event(10119, 16385, uintptr(unsafe.Pointer(alloc.InternCString16(name))), 0)
						r := row{Flags: flags, Mode: mode, Before: count, Gameplay: gameplay}
						var confirm func()
						calls := 0
						oldDialog, oldFocus := legacy.Nox_xxx_dialogMsgBoxCreate_449A10, legacy.Sub_44A360
						t.Cleanup(func() { legacy.Nox_xxx_dialogMsgBoxCreate_449A10 = oldDialog; legacy.Sub_44A360 = oldFocus })
						legacy.Nox_xxx_dialogMsgBoxCreate_449A10 = func(parent *gui.Window, title, text string, df gui.DialogFlags, ok, cancel func()) {
							calls++
							r.Dialog = uint32(df)
							r.Prompt = text
							confirm = ok
							if parent != o.options || cancel != nil {
								t.Error("dialog ownership")
							}
						}
						legacy.Sub_44A360 = func(v int) { r.Focus += v }
						if legacy.PortTestServerOptionsEvent(o.options, o.options.ChildByID(10145), 16391, 0) != 1 {
							t.Fatal("apply event result")
						}
						wantDialog := uint32(0)
						if flags&128 != 0 {
							if mode&0x60 != 0 {
								if count < 2 {
									wantDialog = 56
								} else if count > 2 {
									wantDialog = 33
								}
							} else if flags&16 != 0 && gameplay&4 != 0 && count > 2 {
								wantDialog = 33
							}
						}
						wantCalls := 0
						if wantDialog != 0 {
							wantCalls = 1
						}
						if r.Dialog != wantDialog || calls != r.Focus || calls != wantCalls {
							t.Fatalf("dialog %+v want %d calls %d", r, wantDialog, calls)
						}
						if (confirm != nil) != (wantDialog == 56) {
							t.Fatal("confirmation callback")
						}
						if confirm != nil {
							if o.c.srv.Teams.Count() != count {
								t.Fatal("teams created before confirmation")
							}
							confirm()
						}
						wantCount := count
						if wantDialog != 33 && mode&0x60 != 0 && wantCount < 2 {
							wantCount = 2
						}
						r.After = o.c.srv.Teams.Count()
						r.ResultFlags = uint(noxflags.GetGamePlay())
						wantFlags := gameplay
						if wantDialog != 33 && mode&0x60 != 0 {
							wantFlags |= 4
						}
						if r.After != wantCount || r.ResultFlags != uint(wantFlags) || *o.optionWords["root"] == 0 {
							t.Fatalf("apply state %+v want teams %d flags %d", r, wantCount, wantFlags)
						}
						rows = append(rows, r)
					})
				}
			}
		}
	}
	spellbookCapture(t, "server-options-team-confirmation", rows, "5e8631f3ce721d61a3d0b7ce849ec86b173619d703ff2db3a77d335453753dcb")
}
