//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/client/keybind"
	"github.com/opennox/opennox/v1/legacy"
	"reflect"
	"sort"
	"testing"
)

func TestBindingEditorAssignments(t *testing.T) {
	var records []bindingResult
	for _, menu := range []bool{false, true} {
		t.Run(fmt.Sprintf("menu=%v", menu), func(t *testing.T) {
			o := newBindingOwner(t, menu)
			keys := append(keybind.ListKeys(), keybind.Key(0), keybind.Key(0xffffffff))
			sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
			for _, key := range keys {
				for _, mouse := range []bool{false, true} {
					title := ""
					if key.IsValid() {
						title = key.Title(o.c.Strings())
					}
					for _, column := range []int{-1, 0, 1} {
						for _, row := range []int{-1, 0, 3} {
							initial := [2][]string{{title, "other Ω", title, " "}, {title, title, "other Ω", " "}}
							effectiveRow := row
							if column >= 0 && row >= 0 && initial[column][row] == "" {
								effectiveRow = -1
							}
							o.resetRows(initial, column, row)
							got := legacy.PortTestBindingAssign(menu, mouse, uint32(key))
							r := o.snapshot(mouse, uint32(key), column, row, got)
							want := bindingResult{Menu: menu, Mouse: mouse, Key: uint32(key), Column: column, Row: row, Result: 1, Selection: [2]int32{-1, -1}, Selected: column}
							for i := range initial {
								want.Text[i] = append([]string(nil), initial[i]...)
							}
							if column >= 0 {
								want.Selection[column] = int32(effectiveRow)
							}
							if !mouse && !key.IsValid() {
								want.Result = 0
							} else if column >= 0 {
								for i := range want.Text {
									for j, s := range want.Text[i] {
										if s == title {
											want.Text[i][j] = " "
										}
									}
								}
								if effectiveRow >= 0 {
									want.Text[column][effectiveRow] = title
								}
								want.Selection[column] = -1
								want.Selected = -1
							}
							if !reflect.DeepEqual(r, want) {
								t.Fatalf("key=%v mouse=%v column=%d row=%d\ngot %#v\nwant %#v", key, mouse, column, row, r, want)
							}
							records = append(records, r)
						}
					}
				}
			}
		})
	}
	spellbookCapture(t, "binding-assignments", records, "90fcb0808475ebc4ff87d1b6d9334cceadb84acb88b9a7d1670ff908b22d37f6")
}
