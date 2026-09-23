//go:build porttest

package opennox

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
)

func TestTextAudioDirectory(t *testing.T) {
	t.Cleanup(handles.PortTestInit())
	dir := t.TempDir()
	base := filepath.Join(dir, "catalog")
	// A real on-disk catalog exercises file handles and the directory query's
	// owner lifetime. The production path queries the exact `audio` directory.
	header := make([]byte, 12)
	binary.LittleEndian.PutUint32(header[4:], 1)
	if err := os.WriteFile(base+".idx", header, 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(base+".bag", []byte{1, 2, 3}, 0600); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"audio", "empty", "Ång"} {
		if err := os.Mkdir(filepath.Join(dir, name), 0700); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(dir, "plain"), []byte("x"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(dir, "audio"), filepath.Join(dir, "link")); err != nil {
		t.Fatal(err)
	}
	type row struct {
		Kind    int
		Name    string
		Present uint32
		Stored  string
	}
	var rows []row
	for _, test := range []struct {
		name    string
		present uint32
	}{{"", 0}, {"audio", 1}, {"empty", 1}, {"Ång", 1}, {"link", 1}, {"plain", 0}, {"missing", 0}} {
		for kind := 0; kind < 3; kind++ {
			path := ""
			if test.name != "" {
				path = filepath.Join(dir, test.name)
				if kind == 1 {
					path += "/"
				}
				if kind == 2 {
					path = strings.ReplaceAll(path, "/", "\\")
				}
			}
			p := sub_4866F0(base, path)
			if p == nil {
				t.Fatal("directory catalog owner")
			}
			want := test.present
			stored := alloc.GoStringS(p.path2_8[:])
			expected := path
			if expected != "" {
				expected += "\\"
			}
			if p.field276 != want || stored != expected || p.size4 != 0 || p.bagfile268 == nil {
				t.Fatalf("directory name=%s kind=%d present=%d/%d stored=%q/%q", test.name, kind, p.field276, want, stored, expected)
			}
			// Replace only the identified temporary directory; retain path separators.
			normalized := strings.ReplaceAll(stored, dir, "<tmp>")
			normalized = strings.ReplaceAll(normalized, strings.ReplaceAll(dir, "/", "\\"), "<tmp>")
			rows = append(rows, row{kind, test.name, uint32(p.field276), normalized})
			p.Free()
		}
	}
	interactionCapture(t, "text-audio-directory", rows)
}
