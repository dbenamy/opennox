//go:build porttest

package opennox

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/unit/ai"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
)

func TestCreatureXferActionEdges(t *testing.T) {
	s := newCreatureXferOwner(t)
	path := filepath.Join(t.TempDir(), "edges.bin")
	type row struct {
		Case        string
		State       [6]uint32
		Return, CRC uint32
		Position    int64
	}
	var rows []row
	defer func() { spellbookCapture(t, "creature-xfer-action-edges", rows, "") }()
	for _, name := range []string{"", "unknown-action", strings.Repeat("x", 255), ai.ActionType(0).String()} {
		for _, tail := range []bool{false, true} {
			t.Run(fmt.Sprintf("name%d-tail%v", len(name), tail), func(t *testing.T) {
				u := newCreatureXferObject(t, s, "Monster")
				p := new(mapDrawableStream)
				itemXferRewardName(p, name)
				p.u8(0)
				if tail {
					p.u32(0xaabbccdd)
				}
				if err := os.WriteFile(path, p.Bytes(), 0600); err != nil {
					t.Fatal(err)
				}
				if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
					t.Fatal(err)
				}
				defer cryptfile.Close()
				entry := [6]uint32{0, 1, 2, 3, 4, 0x51515151}
				ret := legacy.PortTestCreatureXferHelper(1, u, unsafe.Pointer(&entry[0]), 5)
				want := [6]uint32{0, 1, 2, 3, 4, 0x51515151}
				wantRet := uint32(0)
				if tail {
					want[5] = 0xaabbccdd
					wantRet = 1
				}
				pos, err := cryptfile.Global().File.Seek(0, 1)
				if err != nil || pos != int64(p.Len()) || entry != want || ret != wantRet {
					t.Fatalf("state=%x ret=%d pos=%d, want=%x ret=%d pos=%d", entry, ret, pos, want, wantRet, p.Len())
				}
				rows = append(rows, row{t.Name(), entry, ret, cryptfile.Global().PortTestChecksum(), pos})
			})
		}
	}
	// Kind 6 is supported by the format but unused by the shipped metadata.
	// Unknown kinds stop before touching the argument or trailing word.
	for _, kind := range []uint32{6, 8, 255, 0xffffffff} {
		t.Run(fmt.Sprintf("kind%d", kind), func(t *testing.T) {
			count := memmap.PtrUint32(0x587000, 255604)
			arg := memmap.PtrUint32(0x587000, 255608)
			oldCount, oldArg := *count, *arg
			*count, *arg = 1, kind
			defer func() { *count, *arg = oldCount, oldArg }()
			u := newCreatureXferObject(t, s, "Monster")
			p := new(mapDrawableStream)
			itemXferRewardName(p, ai.ActionType(0).String())
			p.u8(1)
			stop := p.Len()
			p.u32(0xfffffffd)
			p.u32(0xaabbccdd)
			if err := os.WriteFile(path, p.Bytes(), 0600); err != nil {
				t.Fatal(err)
			}
			if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
				t.Fatal(err)
			}
			defer cryptfile.Close()
			entry := [6]uint32{0, 1, 2, 3, 4, 0x51515151}
			want := entry
			ret := legacy.PortTestCreatureXferHelper(1, u, unsafe.Pointer(&entry[0]), 5)
			wantRet := kind
			if kind == 6 {
				want[1] = 0xfffffffd
				want[5] = 0xaabbccdd
				wantRet = 1
				stop = p.Len()
			}
			pos, err := cryptfile.Global().File.Seek(0, 1)
			if err != nil || entry != want || ret != wantRet || pos != int64(stop) {
				t.Fatalf("state=%x ret=%d pos=%d, want=%x ret=%d pos=%d", entry, ret, pos, want, wantRet, stop)
			}
			rows = append(rows, row{t.Name(), entry, ret, cryptfile.Global().PortTestChecksum(), pos})
		})
	}
}
