//go:build porttest

package opennox

import (
	"bytes"
	"fmt"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"os"
	"path/filepath"
	"testing"
	"unsafe"
)

func TestQuickbarSavedSpellRows(t *testing.T) {
	q := newQuickbarOwner(t)
	original := cryptfile.Global()
	cryptfile.SetGlobal(nil)
	defer func() { cryptfile.Close(); cryptfile.SetGlobal(original) }()
	names := []string{"SPELL_INVALID", "SPELL_ANCHOR", "SPELL_ARACHNAPHOBIA", "SPELL_BLIND", "SPELL_BLINK", "SPELL_BURN"}
	var records []struct {
		State quickbarResult
		Saved []byte
	}
	for _, rowCount := range []uint32{0xffffffff, 0, 1, 3, 5} {
		for _, slotCount := range []uint32{0xffffffff, 0, 1, 3, 5} {
			q.reset(t)
			for i := 0; i < 25; i++ {
				q.bar[2*i] = uint32(i % 6)
				q.bar[2*i+1] = 0x12345600 + uint32(i)
			}
			var want []byte
			for row := 0; row < int(int32(rowCount)); row++ {
				for slot := 0; slot < int(int32(slotCount)); slot++ {
					i := 5*row + slot
					name := names[i%6]
					want = append(want, byte(len(name)))
					want = append(want, []byte(name)...)
					want = append(want, byte(i))
				}
			}
			path := filepath.Join(t.TempDir(), "quickbar.dat")
			if err := cryptfile.OpenGlobal(path, cryptfile.WriteOnly, -1); err != nil {
				t.Fatal(err)
			}
			ret := q.call("sub_460A10", uint32(uintptr(unsafe.Pointer(&q.bar[0]))), rowCount, slotCount, 1)
			if err := cryptfile.Close(); err != nil {
				t.Fatal(err)
			}
			raw, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			q.check(t, ret == 1 && bytes.Equal(raw, want), "saved slots use row stride, byte-length names and one flag byte")
			label := fmt.Sprintf("rows%d-slots%d", rowCount, slotCount)
			records = append(records, struct {
				State quickbarResult
				Saved []byte
			}{q.snapshot(label+"-save", ret), raw})
			for i := 0; i < 25; i++ {
				q.bar[2*i] = 99
				q.bar[2*i+1] = 0xaabbccff
			}
			if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
				t.Fatal(err)
			}
			ret = q.call("sub_460A10", uint32(uintptr(unsafe.Pointer(&q.bar[0]))), rowCount, slotCount, 1)
			if err := cryptfile.Close(); err != nil {
				t.Fatal(err)
			}
			for row := 0; row < 5; row++ {
				for slot := 0; slot < 5; slot++ {
					i := row*5 + slot
					if row < int(int32(rowCount)) && slot < int(int32(slotCount)) {
						q.check(t, q.bar[2*i] == uint32(i%6) && q.bar[2*i+1] == 0xaabbcc00+uint32(i), "load restores id and preserves upper flag bytes")
					} else {
						q.check(t, q.bar[2*i] == 99 && q.bar[2*i+1] == 0xaabbccff, "load preserves unselected slots")
					}
				}
			}
			records = append(records, struct {
				State quickbarResult
				Saved []byte
			}{q.snapshot(label+"-load", ret), nil})
		}
	}
	spellbookCapture(t, "quickbar-save-spells", records, "9a89dda2d6bc2811cbab65c3339e69375934bc4b1951135a99319bd8b08af81a")
}
