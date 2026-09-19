//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	flags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"strings"
	"testing"
	"unsafe"
)

func TestPlayerFilesFieldbookRestore(t *testing.T) {
	o := newReliableReportsOwner(t)
	t.Cleanup(flags.PortTestGameFlags(2048))
	u := &o.units[0]
	p := u.UpdateDataPlayer().Player
	p.NetCodeVal = u.NetCode
	raw := unsafe.Slice((*byte)(unsafe.Pointer(p)), int(unsafe.Sizeof(*p)))
	table := unsafe.Slice(memmap.PtrUint32(0x587000, 70500), 41)
	old := append([]uint32(nil), table...)
	names := make([]string, 41)
	for id := range table {
		names[id] = fmt.Sprintf("RestoreGuide%02d", id)
		if id == 1 {
			names[id] = ""
		}
		if id == 24 {
			names[id] = "Café Ω"
		}
		if id == 40 {
			names[id] = strings.Repeat("Z", 255)
		}
		name, free := alloc.CString(names[id])
		t.Cleanup(free)
		table[id] = uint32(uintptr(unsafe.Pointer(name)))
	}
	t.Cleanup(func() { copy(table, old) })
	family, free := alloc.New([4]uint32{})
	*family = [4]uint32{24, 7, 8, 0}
	t.Cleanup(free)
	families := unsafe.Slice(memmap.PtrUint32(0x587000, 216292), 2)
	oldFamily := append([]uint32(nil), families...)
	t.Cleanup(func() { copy(families, oldFamily) })
	families[0] = uint32(uintptr(unsafe.Pointer(family)))
	families[1] = 0
	ids, resetRecords, snapshotRecords, release := legacy.PortTestPlayerFileRecords()
	t.Cleanup(release)
	binary.LittleEndian.PutUint32(raw[4636:], ids[0])
	binary.LittleEndian.PutUint32(raw[4640:], ids[1])
	var rows []map[string]any
	for _, version := range []uint16{0, 1, 0x8000, 0xffff} {
		for _, count := range []int{0, 1, 3, 4, 41} {
			o.reset()
			resetRecords()
			clear(raw[4248:4408])
			input := binary.LittleEndian.AppendUint16(nil, version)
			input = append(input, 1, byte(count))
			expected := bytes.Clone(raw)
			bits := uint32(0)
			var packets [][]byte
			sequence := []int{1, 24, 40, 24}
			for i := 0; i < count; i++ {
				id := sequence[i%len(sequence)]
				input = append(input, byte(len(names[id])))
				input = append(input, names[id]...)
				if binary.LittleEndian.Uint32(expected[4244+4*id:]) != 0 {
					continue
				}
				grants := []int{id}
				if id == 24 {
					grants = append(grants, 7, 8)
				}
				for _, grant := range grants {
					binary.LittleEndian.PutUint32(expected[4244+4*grant:], 1)
					bits |= 1 << uint(grant&31)
				}
				packets = append(packets, []byte{209, byte(id), 0})
			}
			size := len(input)
			input = append(input, 0xde, 0xad, 0xbe, 0xef)
			ret, got, pos := playerFileSection(t, "nox_xxx_guiFieldbook_41B420", input, uint32(uintptr(unsafe.Pointer(u))), 0)
			records := snapshotRecords()
			if ret != 1 || pos != int64(size) || !bytes.Equal(got, input) || !bytes.Equal(raw, expected) || records != [3]uint32{0, bits, 0x2468ace0 ^ bits} {
				t.Fatal("fieldbook restore", version, count, ret, pos, size, records, bits, "family", *family, "guides", raw[4248:4408], "wanted", expected[4248:4408])
			}
			state := o.state()
			if len(state.Nodes) != len(packets) {
				t.Fatal("guide reports", len(state.Nodes), len(packets))
			}
			for i, n := range state.Nodes {
				if !bytes.Equal(n.Data, packets[len(packets)-1-i]) {
					t.Fatal("guide packet", n.Data, packets)
				}
			}
			rows = append(rows, map[string]any{"version": version, "count": count, "return": ret, "position": pos, "guides": bytes.Clone(raw[4248:4408]), "records": records, "queue": state})
		}
	}
	spellbookCapture(t, "player-files-fieldbook-restore", rows, "a1a95b881d8f64baf6a81cf20436b549f65e5aaae7f57f825aec7f1a92b93741")
}
