//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/client/gui"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestServerOptionsNumericEdit(t *testing.T) {
	type row struct {
		ID      int
		Text    string
		Result  int
		Score   uint16
		Minutes byte
		Dirty   uint32
		Shown   string
	}
	var rows []row
	cases := []struct {
		text    string
		score   uint16
		minutes byte
		shown   string
		dirty   bool
	}{
		{"", 0x7654, 0x43, "", false}, {"0", 0, 0, "0", true}, {"-1", 0, 0, "-1", true},
		{"1", 1, 1, "1", true}, {"255", 255, 255, "255", true}, {"256", 256, 255, "255", true},
		{"65535", 65535, 255, "255", true}, {"65536", 0, 255, "255", true},
		{" 12", 12, 12, " 12", true}, {"17tail", 17, 17, "17tail", true}, {"no digits", 0, 0, "no digits", true},
	}
	for _, id := range []int{10134, 10135} {
		for _, tc := range cases {
			t.Run(fmt.Sprintf("%d-%q", id, tc.text), func(t *testing.T) {
				o := newServerOptionsOwner(t)
				binary.LittleEndian.PutUint16(o.settings[54:], 0x7654)
				o.settings[56] = 0x43
				o.text(uint(id), tc.text)
				result := legacy.PortTestServerOptionsEdit(o.options, 0, id)
				wantScore, wantTime := uint16(0x7654), byte(0x43)
				shown := tc.text
				if id == 10134 {
					wantScore = tc.score
				} else {
					wantTime = tc.minutes
					shown = tc.shown
				}
				score, minutes := binary.LittleEndian.Uint16(o.settings[54:]), o.settings[56]
				dirty := *o.optionWords["dirty"]
				if result != 1 || score != wantScore || minutes != wantTime || (dirty != 0) != tc.dirty || o.entry(uint(id)) != shown {
					t.Fatalf("edit result %d score %d minutes %d dirty %d shown %q", result, score, minutes, dirty, o.entry(uint(id)))
				}
				rows = append(rows, row{id, tc.text, result, score, minutes, dirty, o.entry(uint(id))})
			})
		}
	}
	spellbookCapture(t, "server-options-numeric-edits", rows, "db0a103974a9b8bce247557b5cbe72bb79c32a9ccb5fb4bbd07a3e14ab4abbf8")
}

func TestServerOptionsQuestControls(t *testing.T) {
	for _, value := range []int{-1, 0, 1, 2, 255} {
		t.Run(fmt.Sprint(value), func(t *testing.T) {
			o := newServerOptionsOwner(t)
			o.call("quest", value, "")
			for _, id := range []uint{10152, 10141, 10134, 10135} {
				enabled := o.options.ChildByID(id).GetFlags().Has(gui.StatusEnabled)
				if enabled != (value != 1) {
					t.Fatalf("control %d enabled %t", id, enabled)
				}
			}
			if o.options.ChildByID(10183).GetFlags().IsHidden() != (value != 0) {
				t.Fatal("map panel visibility")
			}
			if o.options.ChildByID(10122).GetFlags().IsHidden() != (value == 1) {
				t.Fatal("checkbox visibility")
			}
		})
	}
}

func TestServerOptionsSettingsLabels(t *testing.T) {
	type row struct {
		Flags                                               uint32
		Mode                                                uint16
		Locked                                              byte
		Limit                                               string
		ScoreEnabled, TimeEnabled, MapEnabled, PanelEnabled bool
		Checked                                             uint32
	}
	var rows []row
	for _, flags := range []uint32{0, 1, 128, 129, 16385, 32769} {
		for _, mode := range []uint16{0, 0x20, 0x100, 0x400, 0x420, 0x1000} {
			for _, locked := range []byte{0, 1, 255} {
				t.Run(fmt.Sprintf("%x-%x-%d", flags, mode, locked), func(t *testing.T) {
					o := newServerOptionsOwner(t)
					defer noxflags.PortTestGameFlags(noxflags.GameFlag(flags))()
					legacy.PortTestServerOptionsModeName(0)
					copy(o.settings, "Arena")
					copy(o.settings[9:], "Test server")
					binary.LittleEndian.PutUint16(o.settings[52:], mode)
					o.settings[57] = locked
					// Start disabled to detect whether each branch actively enables a control.
					for _, id := range []uint{10134, 10135, 10183, 10197} {
						o.options.ChildByID(id).Flags &^= gui.StatusEnabled
					}
					o.call("settings-labels", 0, "")
					want := "Kill limit"
					if mode&0x20 != 0 {
						want = "Capture limit"
					} else if mode&0x400 != 0 {
						want = "Death limit"
					}
					score, time := o.options.ChildByID(10134), o.options.ChildByID(10135)
					if got := score.DrawData().Text(); got != want {
						t.Fatalf("limit %q, want %q", got, want)
					}
					enabled := flags&1 != 0 && locked == 0 && (mode&0x420 != 0 || flags&49152 == 0)
					if score.GetFlags().Has(gui.StatusEnabled) != enabled || time.GetFlags().Has(gui.StatusEnabled) != enabled {
						t.Fatal("limit control enablement")
					}
					mapEnabled := flags&1 != 0 && flags&49152 == 0 && locked == 0
					if o.options.ChildByID(10183).GetFlags().Has(gui.StatusEnabled) != mapEnabled || o.options.ChildByID(10197).GetFlags().Has(gui.StatusEnabled) != mapEnabled {
						t.Fatal("map/settings panel enablement")
					}
					checked := o.options.ChildByID(10122).DrawData().Field0 & 4
					if (checked != 0) != (locked != 0) {
						t.Fatal("map cycle checkbox")
					}
					rows = append(rows, row{flags, mode, locked, want, score.GetFlags().Has(gui.StatusEnabled), time.GetFlags().Has(gui.StatusEnabled), mapEnabled, mapEnabled, checked})
				})
			}
		}
	}
	spellbookCapture(t, "server-options-settings-labels", rows, "d565249b8851298f3206a2eedca5a70e1cad0ccce85c75f1915c17969a6b8832")
}
