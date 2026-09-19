//go:build porttest

package opennox

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"

	"github.com/opennox/libs/ifs"
	"github.com/opennox/libs/spell"
	"github.com/opennox/opennox/v1/legacy"
)

func TestSessionFilterRuleComparison(t *testing.T) {
	restore := legacy.PortTestSessionFilterRules()
	defer restore()
	old, err := ifs.Workdir()
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	if err = ifs.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := ifs.Chdir(old); err != nil {
			t.Error(err)
		}
	}()
	if err = os.MkdirAll(filepath.Join(dir, "maps", "porttest"), 0700); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "maps", "porttest", "porttest.rul")
	type row struct {
		File     bool
		Flags    uint16
		Mutation int
		Result   int
	}
	var rows []row
	for _, withFile := range []bool{false, true} {
		if withFile {
			if err = os.WriteFile(path, []byte("[COMMON]\nset spell SPELL_FIREBALL off\n"), 0600); err != nil {
				t.Fatal(err)
			}
		}
		for _, flags := range []uint16{0x100, 0x40, 0x1040} {
			var record [169]byte
			copy(record[111:], "porttest")
			binary.LittleEndian.PutUint16(record[163:], flags)
			for i := 135; i < 163; i++ {
				record[i] = 255
			}
			if withFile {
				bit := int(spell.SPELL_FIREBALL)
				record[135+bit/8] &^= 1 << uint(bit%8)
			}
			if flags&0x40 != 0 {
				record[135+132/8] &^= 1 << uint(132%8)
			}
			var filters [11]uint32
			filters[5] = 1
			// Every compared byte must affect acceptance; this spans all five spell
			// words, four equipment bytes, and the final armor word.
			for mutation := -1; mutation < 28; mutation++ {
				input := record
				if mutation >= 0 {
					input[135+mutation] ^= 1 << uint(mutation%8)
				}
				got, intact := legacy.PortTestSessionFilter(2, filters, input)
				want := 0
				if mutation == -1 {
					want = 1
				}
				if got != want || !intact {
					t.Fatal("rule comparison", withFile, flags, mutation, got, want, intact)
				}
				rows = append(rows, row{withFile, flags, mutation, got})
			}
		}
	}
	sessionCapture(t, "filter-rules", rows)
}
