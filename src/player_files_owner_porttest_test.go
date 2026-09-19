//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func playerFileSection(t *testing.T, name string, input []byte, args ...uint32) (uint32, []byte, int64) {
	t.Helper()
	old := cryptfile.Global()
	cryptfile.SetGlobal(nil)
	defer func() { cryptfile.Close(); cryptfile.SetGlobal(old) }()
	path := filepath.Join(t.TempDir(), "section.bin")
	mode := cryptfile.WriteOnly
	if input != nil {
		if err := os.WriteFile(path, input, 0600); err != nil {
			t.Fatal(err)
		}
		mode = cryptfile.ReadOnly
	}
	if err := cryptfile.OpenGlobal(path, mode, -1); err != nil {
		t.Fatal(err)
	}
	ret := legacy.PortTestPlayerFileCall(name, args...)
	pos, err := cryptfile.Global().File.Seek(0, io.SeekCurrent)
	if err != nil {
		t.Fatal(err)
	}
	if err := cryptfile.Close(); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return ret, data, pos
}
