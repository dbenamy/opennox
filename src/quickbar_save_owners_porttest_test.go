//go:build porttest

package opennox

import (
	"bytes"
	"fmt"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"os"
	"path/filepath"
	"testing"
	"unsafe"
)

func TestQuickbarSaveOwners(t *testing.T) {
	q := newQuickbarOwner(t)
	original := cryptfile.Global()
	cryptfile.SetGlobal(nil)
	defer func() { cryptfile.Close(); cryptfile.SetGlobal(original) }()
	table := unsafe.Slice(memmap.PtrUint32(0x587000, 69736), 7)
	oldTable := append([]uint32(nil), table...)
	t.Cleanup(func() { copy(table, oldTable) })
	clear(table)
	for i, name := range server.AbilityNames {
		p, free := alloc.CString(name)
		t.Cleanup(free)
		table[i] = uint32(uintptr(unsafe.Pointer(p)))
	}
	spells := []string{"SPELL_INVALID", "SPELL_ANCHOR", "SPELL_ARACHNAPHOBIA", "SPELL_BLIND", "SPELL_BLINK", "SPELL_BURN"}
	var rows []struct {
		State quickbarResult
		Saved []byte
	}
	for class := uint32(0); class < 3; class++ {
		for _, pending := range []uint32{0, 1} {
			q.reset(t)
			*(*byte)(unsafe.Add(unsafe.Pointer(&q.players[0]), 2251)) = byte(class)
			q.call("sub_461440", pending)
			trap := unsafe.Slice(memmap.PtrUint32(0x5D4594, 1047940), 50)
			for i := 0; i < 25; i++ {
				q.bar[2*i] = uint32(i % 6)
				q.bar[2*i+1] = 0xaabbcc00 | uint32(i)
				*memmap.PtrUint32(0x5D4594, 1047564+uintptr(8*i)) = uint32((25 - i) % 6)
				*memmap.PtrUint32(0x5D4594, 1047568+uintptr(8*i)) = 0x11223380 | uint32(i)
				trap[2*i] = uint32(i % 6)
				trap[2*i+1] = 0x33445540 | uint32(i)
			}
			want := []byte{byte(class)}
			var expected [25][2]uint32
			appendSlot := func(id, flag uint32) {
				name := spells[id]
				if class == 0 {
					name = server.AbilityNames[id]
				}
				want = append(want, byte(len(name)))
				want = append(want, []byte(name)...)
				want = append(want, byte(flag))
			}
			for i := 0; i < 25; i++ {
				id, flag := uint32(i%6), uint32(i)
				if pending != 0 {
					id = uint32((25 - i) % 6)
					flag = 0x80 | uint32(i)
				}
				expected[i] = [2]uint32{id, flag}
				appendSlot(id, flag)
			}
			if class != 0 {
				for row := 0; row < 3; row++ {
					for slot := 0; slot < 3; slot++ {
						i := row*5 + slot
						appendSlot(uint32(i%6), 0x40|uint32(i))
					}
				}
			}
			path := filepath.Join(t.TempDir(), "slots.dat")
			if err := cryptfile.OpenGlobal(path, cryptfile.WriteOnly, -1); err != nil {
				t.Fatal(err)
			}
			ret := q.call("sub_460940")
			if err := cryptfile.Close(); err != nil {
				t.Fatal(err)
			}
			raw, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			q.check(t, ret == 1 && bytes.Equal(raw, want), "save owner writes class byte and class-specific row layouts")
			q.check(t, q.call("sub_461450") == 0, "saving consumes pending slot restoration")
			label := fmt.Sprintf("class%d-pending%d", class, pending)
			rows = append(rows, struct {
				State quickbarResult
				Saved []byte
			}{q.snapshot(label+"-write", ret), raw})
			for i := 0; i < 25; i++ {
				q.bar[2*i] = 99
				q.bar[2*i+1] = 0xabcdef00
				trap[2*i] = 99
				trap[2*i+1] = 0x12345600
			}
			if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
				t.Fatal(err)
			}
			ret = q.call("sub_460940")
			if err := cryptfile.Close(); err != nil {
				t.Fatal(err)
			}
			for i, pair := range expected {
				q.check(t, q.bar[2*i] == pair[0] && q.bar[2*i+1] == 0xabcdef00|pair[1], "save owner reloads all five rows preserving flag padding")
			}
			for row := 0; row < 5; row++ {
				for slot := 0; slot < 5; slot++ {
					i := row*5 + slot
					id, flag := uint32(99), uint32(0x12345600)
					if class != 0 && row < 3 && slot < 3 {
						id = uint32(i % 6)
						flag |= 0x40 | uint32(i)
					}
					q.check(t, trap[2*i] == id && trap[2*i+1] == flag, "only spell classes load the three-by-three trap rows")
				}
			}
			rows = append(rows, struct {
				State quickbarResult
				Saved []byte
			}{q.snapshot(label+"-read", ret), nil})
		}
	}
	spellbookCapture(t, "quickbar-save-owners", rows, "c2cc5eaacb5cfcd2f0a557bfa221f5e1d62e82982435c34dddbeefdb79281335")
}
