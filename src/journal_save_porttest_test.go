//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
	"strings"
	"testing"

	flags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

func TestJournalSaveLoadOwner(t *testing.T) {
	o := newJournalOwner(t)
	var rows []journalResult
	oldFile := cryptfile.Global()
	cryptfile.SetGlobal(nil)
	defer func() { cryptfile.Close(); cryptfile.SetGlobal(oldFile) }()
	oldFlags := flags.GetGame()
	defer func() { flags.ResetGame(); flags.SetGame(oldFlags) }()
	for _, count := range []int{0, 1, 3, 6} {
		o.resetJournal(t)
		flags.ResetGame()
		flags.SetGame(2048)
		names := []string{"Short", "Wrapped", "Multiline", "", strings.Repeat("X", 63), "Café-Ω"}
		kinds := []uint16{1, 2, 4, 8, 0x8000, 0xffff}
		want := []byte{1, 0, 1, byte(count), 0}
		for i := 0; i < count; i++ {
			o.call(t, 0, 31, names[i], kinds[i])
			want = append(want, byte(len(names[i])))
			want = append(want, []byte(names[i])...)
			want = binary.LittleEndian.AppendUint16(want, kinds[i])
		}
		path := filepath.Join(t.TempDir(), "journal.bin")
		if err := cryptfile.OpenGlobal(path, cryptfile.WriteOnly, -1); err != nil {
			t.Fatal(err)
		}
		ret := legacy.PortTestJournal(11, uintptr(o.units[2].CObj()), 0, 0)
		if err := cryptfile.Close(); err != nil {
			t.Fatal(err)
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if ret != 1 || !bytes.Equal(raw, want) {
			t.Fatalf("journal save bytes got%x want%x return%d", raw, want, ret)
		}
		r := o.snapshot(t, 11, ret, false)
		r.Saved = raw
		rows = append(rows, r)
		o.release()
		o.call(t, 0, 31, "OldEntry", 2)
		o.queueReset()
		if err = cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
			t.Fatal(err)
		}
		ret = legacy.PortTestJournal(11, uintptr(o.units[2].CObj()), 0, 0)
		if err = cryptfile.Close(); err != nil {
			t.Fatal(err)
		}
		if ret != 1 {
			t.Fatal("journal load return")
		}
		// Newly allocated loaded nodes receive identities in their serialized order.
		var loaded []*server.PlayerJournal
		for n := o.players[31].Journal; n != nil; n = n.Next {
			loaded = append(loaded, n)
		}
		for i := len(loaded) - 1; i >= 0; i-- {
			o.remember(loaded[i])
		}
		r = o.snapshot(t, 11, ret, false)
		rows = append(rows, r)
		if len(loaded) != count {
			t.Fatal("journal load did not replace prior entries")
		}
		for i, n := range loaded {
			j := count - 1 - i
			if string(bytes.TrimRight(n.EntryBuf[:], "\x00")) != names[j] || n.Field3 != kinds[j] {
				t.Fatal("journal save/load order or entry fields")
			}
		}
	}
	journalCapture(t, "save-load", rows, "a06abea7c494eeb31fb9675e526a4f6fddafcd604edc6e564130609460084c40")
}
