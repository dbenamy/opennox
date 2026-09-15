//go:build porttest

package opennox

import (
	"fmt"
	"testing"
	"unicode/utf16"
	"unsafe"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestClientInventoryDisplayNameSign(t *testing.T) {
	o := newInventoryDisplayOwner(t)
	var rows []inventoryDisplayResult
	for class := 0; class < 3; class++ {
		for _, level := range []uint32{0, 1, 5, 10, 11, 127, 128, 255} {
			for _, special := range []bool{false, true} {
				for _, name := range []string{"Jack", "Jäck Ω"} {
					o.reset(t)
					if special {
						noxflags.SetGame(4096)
					}
					*o.displayWords["dword_5d4594_1049844"] = level
					*(*byte)(unsafe.Add(unsafe.Pointer(&o.players[0]), 2251)) = byte(class)
					*(*byte)(unsafe.Add(unsafe.Pointer(&o.players[0]), 3684)) = byte(level)
					dst := unsafe.Slice((*uint16)(unsafe.Add(unsafe.Pointer(&o.players[0]), 4704)), 25)
					alloc.StrCopyZero16(dst, name)
					r := o.call(t, len(rows), 14, 0, 0, 0)
					rank := int(int8(level))
					if special {
						rank = int(level)
						if rank > 10 {
							rank = 10
						}
					}
					want := fmt.Sprintf("%s the Rank%d-%d", name, class, rank)
					text := alloc.GoString16((*uint16)(memmap.PtrOff(0x5D4594, 1062588)))
					if text != want || r.Return != uint64(len(utf16.Encode([]rune(want)))) {
						t.Fatalf("name sign class%d level%d special%v got%q/%d want%q", class, level, special, text, r.Return, want)
					}
					rows = append(rows, r)
				}
			}
		}
	}
	o.reset(t)
	*o.displayWords["dword_8531A0_2576"] = 0
	r := o.call(t, len(rows), 14, 0, 0, 0)
	if r.Return != 0 {
		t.Fatal("no-player name return")
	}
	rows = append(rows, r)
	inventoryDisplayCapture(t, "names", rows, "8e10769a832cb7a5d54db155809e5545e00f5f7886117a7e1f809701f096649b")
}
