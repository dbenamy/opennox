//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/common"
	"github.com/opennox/libs/datapath"
	"github.com/opennox/opennox/v1/legacy"
	"math/bits"
	"os"
	"path/filepath"
	"testing"
)

func TestSessionEntrySaveSlots(t *testing.T) {
	old := datapath.Data()
	t.Cleanup(func() { datapath.SetData(old) })
	type row struct {
		Mask  uint32
		Extra string
		Count int32
	}
	var rows []row
	masks := []uint32{0, 1, 0x3fff, 0x2003, 0x2020, 0x2222}
	for i := uint32(1); i < 14; i++ {
		masks = append(masks, uint32(1)<<i)
	}
	for _, mask := range masks {
		for _, extra := range []string{"", "SAVE0000", "SAVE0014"} {
			t.Run(fmt.Sprintf("mask%x/extra%s", mask, extra), func(t *testing.T) {
				dir := t.TempDir()
				datapath.SetData(dir)
				add := func(slot string) {
					p := filepath.Join(dir, "Save", slot, common.PlayerFile)
					if err := os.MkdirAll(filepath.Dir(p), 0700); err != nil {
						t.Fatal(err)
					}
					if err := os.WriteFile(p, []byte("presence-only fixture"), 0600); err != nil {
						t.Fatal(err)
					}
				}
				if mask&1 != 0 {
					add(common.SaveAuto)
				}
				// The actual save-selection UI reserves index0 for AUTOSAVE and lists1..13.
				for i := 1; i < NOX_SAVEGAME_XXX_MAX; i++ {
					if mask&(uint32(1)<<i) != 0 {
						add(fmt.Sprintf(common.SaveFormat, i))
					}
				}
				if extra != "" {
					add(extra)
				}
				got := legacy.PortTestSessionEntryScalar("save-count", 0)
				want := int32(bits.OnesCount32(mask))
				if got != want {
					t.Errorf("saved game count=%d want%d", got, want)
				}
				if st, err := os.Stat(filepath.Join(dir, "Save")); err != nil || !st.IsDir() {
					t.Fatal("Save directory not created", err)
				}
				rows = append(rows, row{mask, extra, got})
			})
		}
	}
	spellbookCapture(t, "session-entry-save-slots", rows, "efc4f3ea1a8a0fe0744a31a0f66d21c27f28ded9dd34d64fa87a0a6ac00072e3")
}
