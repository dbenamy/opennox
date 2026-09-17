//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"os"
	"path/filepath"
	"testing"
)

func TestCreatureXferActionHistorical(t *testing.T) {
	s := newCreatureXferOwner(t)
	path := filepath.Join(t.TempDir(), "old-action.bin")
	var rows []itemXferCaptureRow
	defer func() {
		spellbookCapture(t, "creature-xfer-action-historical", rows, "0a2fa616be097bbc8d8d9f492fa3a8389d181671e5af85fa8e7e2107895e366d")
	}()
	for _, version := range []uint16{0, 1, 2, 3, 4, 0x8000, 0xffff} {
		for _, frames := range [][2]uint32{{123, 123}, {123, 0}, {0, 123}, {0xffffffff, 0}, {0, 0xffffffff}} {
			for _, value := range []uint32{0, 1, 0x7fffffff, 0xffffffff} {
				t.Run(fmt.Sprintf("v%04x-frame%08x-%08x-value%08x", version, frames[0], frames[1], value), func(t *testing.T) {
					source := newCreatureXferObject(t, s, "Monster")
					creatureXferSetTimers(source, value)
					// Signed historical versions below one use the oldest record layout.
					layout := version
					if int16(version) < 0 {
						layout = 0
					}
					p := creatureXferEmptyActionRecord(source, layout, frames[0])
					p.Bytes()[0] = byte(version)
					p.Bytes()[1] = byte(version >> 8)
					if err := os.WriteFile(path, p.Bytes(), 0600); err != nil {
						t.Fatal(err)
					}
					u := newCreatureXferObject(t, s, "Monster")
					u.Field34 = 999
					for _, off := range creatureXferAlwaysTimers[3:] {
						objectXferSetWord(u.UpdateData, off, 0x51515151)
					}
					s.SetFrame(frames[1])
					if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
						t.Fatal(err)
					}
					defer cryptfile.Close()
					if ret := legacy.PortTestCreatureXferHelper(0, u, nil, 0); ret != 1 {
						t.Fatalf("result=%d", ret)
					}
					pos, err := cryptfile.Global().File.Seek(0, 1)
					if err != nil || pos != int64(p.Len()) || u.Field34 != 999 {
						t.Fatal("historical action stream/lifetime")
					}
					want := creatureXferClampTimer(value, frames[1]-frames[0])
					for _, off := range creatureXferReadTimers {
						if objectXferGetWord(u.UpdateData, off) != want {
							t.Fatalf("read timer %d", off)
						}
					}
					for i, off := range creatureXferAlwaysTimers {
						expected := want
						if i >= 3 && layout < 4 {
							expected = 0x51515151
						}
						if objectXferGetWord(u.UpdateData, off) != expected {
							t.Fatalf("timer %d=%08x want=%08x", off, objectXferGetWord(u.UpdateData, off), expected)
						}
					}
					itemXferCaptureCase(t, &rows, u)
				})
			}
		}
	}
}
