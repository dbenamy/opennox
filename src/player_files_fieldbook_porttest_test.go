//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	flags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"strings"
	"testing"
	"unsafe"
)

func TestPlayerFilesFieldbookWrite(t *testing.T) {
	u, raw := playerFileBookOwner(t)
	defer flags.PortTestGameFlags(0)()
	table := unsafe.Slice(memmap.PtrUint32(0x587000, 70500), 41)
	old := append([]uint32(nil), table...)
	names := make([]string, 41)
	for id := range table {
		names[id] = fmt.Sprintf("Guide%02d", id)
		if id == 1 {
			names[id] = ""
		}
		if id == 20 {
			names[id] = "Café-Ω"
		}
		if id == 40 {
			names[id] = strings.Repeat("Z", 255)
		}
		p, free := alloc.CString(names[id])
		t.Cleanup(free)
		table[id] = uint32(uintptr(unsafe.Pointer(p)))
	}
	t.Cleanup(func() { copy(table, old) })
	var rows []map[string]any
	for _, gf := range []flags.GameFlag{0, 2048, 4096, 8192, 8192 | 4096, 8192 | 2048} {
		for _, pattern := range []string{"empty", "edges", "all"} {
			flags.ResetGame()
			flags.SetGame(gf)
			clear(raw[4248:4408])
			// The preceding word is guide zero and must never be serialized.
			binary.LittleEndian.PutUint32(raw[4244:], 99)
			for id := 1; id <= 40; id++ {
				if pattern == "all" || pattern == "edges" && (id == 1 || id == 20 || id == 40) {
					binary.LittleEndian.PutUint32(raw[4244+4*id:], uint32(id*7))
				}
			}
			before := bytes.Clone(raw)
			present := byte(1)
			if gf&8192 != 0 && gf&4096 == 0 {
				present = 0
			}
			want := []byte{1, 0, present}
			wantRet := uint32(1)
			if present != 0 {
				if gf&(2048|4096) == 0 {
					wantRet = 0
				} else {
					count := byte(0)
					var entries []byte
					for id := 1; id <= 40; id++ {
						if binary.LittleEndian.Uint32(raw[4244+4*id:]) == 0 {
							continue
						}
						count++
						entries = append(entries, byte(len(names[id])))
						entries = append(entries, names[id]...)
					}
					want = append(want, count)
					want = append(want, entries...)
				}
			}
			ret, got, pos := playerFileSection(t, "nox_xxx_guiFieldbook_41B420", nil, uint32(uintptr(unsafe.Pointer(u))), 0)
			if ret != wantRet || pos != int64(len(want)) || !bytes.Equal(got, want) || !bytes.Equal(raw, before) {
				t.Fatal("fieldbook write", gf, pattern, ret, pos, len(want))
			}
			rows = append(rows, map[string]any{"flags": uint32(gf), "pattern": pattern, "return": ret, "bytes": got, "position": pos})
		}
	}
	spellbookCapture(t, "player-files-fieldbook-write", rows, "7aed107c6b8bb79fff6519189fc6da13cff8d72310e619133641dfadbecbbfa7")
}
