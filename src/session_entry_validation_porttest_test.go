//go:build porttest

package opennox

import (
	"encoding/binary"
	"github.com/opennox/libs/datapath"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"os"
	"path/filepath"
	"testing"
)

func TestSessionEntryMapValidation(t *testing.T) {
	original := cryptfile.Global()
	cryptfile.SetGlobal(nil)
	defer func() { cryptfile.Close(); cryptfile.SetGlobal(original) }()
	old := datapath.Data()
	dir := t.TempDir()
	datapath.SetData(dir)
	defer datapath.SetData(old)
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(cwd)
	sep := serverConfigOwnBytes(t, 0x587000, 191672, 2)
	binary.LittleEndian.PutUint16(sep, '\\')
	const checksum = uint32(0x89abcdef)
	writeMap := func(name string, magic int32, crc uint32) {
		t.Helper()
		path := filepath.Join(dir, "maps", name, name+".map")
		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			t.Fatal(err)
		}
		f, err := cryptfile.OpenFile(path, cryptfile.WriteOnly, 19)
		if err != nil {
			t.Fatal(err)
		}
		if err = f.WriteU32(uint32(magic)); err != nil {
			t.Fatal(err)
		}
		off, err := f.File.FileFlush()
		if err != nil {
			t.Fatal(err)
		}
		if err = f.File.WriteUint32At(crc, off); err != nil {
			t.Fatal(err)
		}
		if err = f.Close(); err != nil {
			t.Fatal(err)
		}
	}
	writeMap("current", -86050098, checksum)
	writeMap("old", -86065425, checksum)
	writeMap("other", 12345, checksum)
	type row struct {
		Name   string
		Nil    bool
		CRC    uint32
		Result int
	}
	var rows []row
	for _, c := range []struct {
		name    string
		nilName bool
		crc     uint32
		want    int
	}{
		{"", true, 0, 6}, {"", true, checksum, 0}, {"missing.map", false, 0, 6},
		{"missing.map", false, checksum, 0}, {"current.map", false, checksum, 6},
		{`maps\current\current.map`, false, checksum, 6}, {"CURRENT.MAP", false, checksum, 6},
		{"current.map", false, checksum + 1, 2}, {"old.map", false, checksum, 2},
		{"other.map", false, checksum, 2},
	} {
		got := legacy.PortTestSessionEntryValidate(c.name, c.nilName, c.crc)
		if got != c.want {
			t.Errorf("%+v: got%d want%d", c, got, c.want)
		}
		if cryptfile.Global() != nil {
			t.Fatal("map file left open")
		}
		rows = append(rows, row{c.name, c.nilName, c.crc, got})
	}
	spellbookCapture(t, "session-entry-map-validation", rows, "d0632f95c01a1f456c5995271f6195f45f9513c7492ddf57b11a553f64ba7432")
}
