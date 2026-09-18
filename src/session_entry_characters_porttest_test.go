//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/datapath"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"os"
	"path/filepath"
	"testing"
)

func TestSessionEntryCharacterCount(t *testing.T) {
	defer noxflags.PortTestGameFlags(0)()
	original := cryptfile.Global()
	cryptfile.SetGlobal(nil)
	defer func() { cryptfile.Close(); cryptfile.SetGlobal(original) }()
	old := datapath.Data()
	defer datapath.SetData(old)
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(cwd)
	serverConfigOwnBytes(t, 0x85B3FC, 10980, 1278)
	restore := legacy.PortTestSessionEntryMetadataTable()
	defer restore()
	type row struct {
		Kind   string
		Count  int32
		AtRoot bool
	}
	var rows []row
	for _, kind := range []string{"empty", "valid", "other", "mixed", "valid-then-empty", "valid-then-rejected", "directory"} {
		t.Run(kind, func(t *testing.T) {
			dir := t.TempDir()
			datapath.SetData(dir)
			save := filepath.Join(dir, "Save")
			if err := os.MkdirAll(save, 0700); err != nil {
				t.Fatal(err)
			}
			write := func(name string, version int, flags uint32) {
				t.Helper()
				f, err := cryptfile.OpenFile(filepath.Join(save, name), cryptfile.WriteOnly, 27)
				if err != nil {
					t.Fatal(err)
				}
				if version != 0 {
					if err := f.WriteU32(1); err != nil {
						t.Fatal(err)
					}
					f.SectionStart()
					if err := f.WriteU16(uint16(version)); err != nil {
						t.Fatal(err)
					}
					if err := f.WriteU32(flags); err != nil {
						t.Fatal(err)
					}
					// Version10: empty path/tag, timestamp, appearance, empty player name, class/state bytes.
					if _, err := f.Write(make([]byte, 2+1+16+20+1+3)); err != nil {
						t.Fatal(err)
					}
					f.SectionEnd()
				}
				if err := f.WriteU32(0); err != nil {
					t.Fatal(err)
				}
				if err := f.Close(); err != nil {
					t.Fatal(err)
				}
			}
			want := int32(0)
			switch kind {
			case "valid":
				write("a.plr", 10, 2)
				want = 1
			case "other":
				write("a.plr", 10, 1)
			case "mixed":
				write("a.plr", 10, 2)
				write("b.plr", 10, 4)
				write("c.plr", 10, 6)
				write("ignored.txt", 10, 2)
				want = 2
			case "valid-then-empty":
				write("a.plr", 10, 2)
				write("b.plr", 0, 0)
				want = 1
			case "valid-then-rejected":
				write("a.plr", 10, 2)
				write("b.plr", 13, 2)
				want = 1
			case "directory":
				if err := os.Mkdir(filepath.Join(save, "fake.plr"), 0700); err != nil {
					t.Fatal(err)
				}
			}
			got := legacy.PortTestSessionEntryScalar("character-count", 0)
			here, err := os.Getwd()
			if err != nil {
				t.Fatal(err)
			}
			if got != want || here != dir {
				t.Error(fmt.Sprintf("%s count%d want%d restored root%v", kind, got, want, here == dir))
			}
			rows = append(rows, row{kind, got, here == dir})
			// Leave a stable directory before TempDir cleanup.
			if err := os.Chdir(cwd); err != nil {
				t.Fatal(err)
			}
		})
	}
	spellbookCapture(t, "session-entry-character-count", rows, "3cc2c73b75ab94b3c711765d20cac55ea2210397d87068b7ce8b94c08b343ef6")
}
