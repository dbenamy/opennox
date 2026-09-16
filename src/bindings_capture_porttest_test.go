//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/client/keybind"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
	"unsafe"
)

func TestBindingEditorCapture(t *testing.T) {
	type captureResult struct {
		Binding                  bindingResult
		Event, State             int
		Hidden, Focused, Stacked bool
	}
	var records []captureResult
	for _, menu := range []bool{false, true} {
		t.Run(fmt.Sprintf("menu=%v", menu), func(t *testing.T) {
			o := newBindingOwner(t, menu)
			for event := 0; event <= 24; event++ {
				for _, key := range []uint32{0, 1, 30, 0xffffffff} {
					for state := 0; state <= 2; state++ {
						for _, column := range []int{-1, 0, 1} {
							o.resetRows([2][]string{{"first", "second"}, {"third", "fourth"}}, column, 0)
							o.modal.SetHidden(false)
							o.modal.Focus()
							o.modal.StackPop()
							o.modal.StackPush()
							got := legacy.PortTestBindingModal(menu, uint32(uintptr(unsafe.Pointer(o.modal))), event, key, state)
							r := captureResult{Binding: o.snapshot(false, key, column, 0, got), Event: event, State: state, Hidden: o.modal.Flags.IsHidden(), Focused: o.c.GUI.Focused() == o.modal, Stacked: o.c.GUI.StackHead() == o.modal}
							close := event == 6 || event == 7 || event == 10 || event == 11 || event == 14 || event == 15 || event == 19 || event == 20 || event == 21 && (key == 1 && state == 2 || key != 1 && state == 1 && keybind.Key(key).IsValid())
							handled := close || event == 21 && key == 1
							want := 0
							if handled {
								want = 1
							}
							if got != want || r.Hidden != close || r.Focused == close || r.Stacked == close {
								t.Fatalf("modal contract: %+v close=%v handled=%v", r, close, handled)
							}
							if event == 21 && key == 1 && state == 2 && column >= 0 {
								if r.Binding.Selection[column] != -1 || r.Binding.Selected != column {
									t.Fatalf("cancel must clear row and preserve selected-list owner: %+v", r)
								}
							}
							records = append(records, r)
						}
					}
				}
			}
		})
	}
	spellbookCapture(t, "binding-capture", records, "664905ce9353cbeff63fae7c6efa1c33bdada83a38ce06b66f3e80b7a9fe242b")
}
