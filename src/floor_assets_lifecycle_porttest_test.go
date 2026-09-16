//go:build porttest

package opennox

import (
	"bytes"
	"crypto/sha256"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"strings"
	"testing"
	"unsafe"
)

func TestFloorAssetsFreeOwners(t *testing.T) {
	o := newFloorAssetsOwner(t)
	type record struct {
		Edge          bool
		Count, Repeat int
		Remaining     []int
		Metadata      [32]byte
	}
	var rows []record
	for _, edge := range []bool{false, true} {
		cap := 176
		if edge {
			cap = 64
		}
		for _, count := range []int{0, 1, cap - 1, cap} {
			o.reset()
			original := make([]unsafe.Pointer, cap)
			for _, i := range []int{0, 1, cap - 1} {
				p := legacy.PortTestFloorAssetsAlloc(16)
				if p == nil {
					t.Fatal("allocation")
				}
				original[i] = p
				for j := range unsafe.Slice((*byte)(p), 16) {
					unsafe.Slice((*byte)(p), 16)[j] = byte(17 + i + j)
				}
				if edge {
					*(*unsafe.Pointer)(unsafe.Pointer(&o.edges[60*i+32])) = p
				} else {
					o.defs[i].Data32 = p
				}
			}
			if edge {
				*o.edgeCount = uint32(count)
			} else {
				*o.count = uint32(count)
			}
			beforeDefs, beforeEdges := o.metadata()
			for repeat := 0; repeat < 2; repeat++ {
				if edge {
					legacy.Sub_485F30()
				} else {
					legacy.Sub_485CF0()
				}
				afterDefs, afterEdges := o.metadata()
				if !bytes.Equal(beforeDefs, afterDefs) || !bytes.Equal(beforeEdges, afterEdges) {
					t.Fatal("free changed metadata")
				}
				if edge && *o.edgeCount != uint32(count) || !edge && *o.count != uint32(count) {
					t.Fatal("free changed count")
				}
				var remaining []int
				for i, p := range original {
					actual := o.defs[i].Data32
					if edge {
						actual = *(*unsafe.Pointer)(unsafe.Pointer(&o.edges[60*i+32]))
					}
					want := p
					if i < count {
						want = nil
					}
					if actual != want {
						t.Fatalf("free edge%v count%d slot%d", edge, count, i)
					}
					if actual != nil {
						remaining = append(remaining, i)
						for j, b := range unsafe.Slice((*byte)(actual), 16) {
							if b != byte(17+i+j) {
								t.Fatal("free changed excluded allocation")
							}
						}
					}
				}
				rows = append(rows, record{edge, count, repeat, remaining, sha256.Sum256(append(afterDefs, afterEdges...))})
			}
		}
	}
	floorAssetsCapture(t, "free", rows, "c1b5d4dec07d07c54b189294a39bcb58fe59abe7504d60ebeb190dddd4650d41")
}
func TestFloorAssetsFacadeNames(t *testing.T) {
	o := newFloorAssetsOwner(t)
	names := floorAssetsFacades()
	if len(names) != 6 {
		t.Fatal("actual facade table must have six names")
	}
	cases := append([]string(nil), names...)
	cases = append(cases, "", strings.ToLower(names[0]), "Unknown", "\x80raw", strings.Repeat("Z", 31))
	first := memmap.PtrPtr(0x587000, 26488)
	old := *first
	t.Cleanup(func() { *first = old })
	type record struct {
		Empty  bool
		Name   string
		Return int
	}
	var rows []record
	for _, empty := range []bool{false, true} {
		*first = old
		if empty {
			*first = nil
		}
		for _, name := range cases {
			clear(o.defs[0].NameBuf[:])
			copy(o.defs[0].NameBuf[:], name)
			before := o.defs[0]
			ret := legacy.PortTestFloorAssetsFacade(0)
			want := 0
			if !empty {
				for _, n := range names {
					if name == n {
						want = 1
					}
				}
			}
			if ret != want || o.defs[0] != before {
				t.Fatalf("facade empty%v name%q got%d want%d", empty, name, ret, want)
			}
			rows = append(rows, record{empty, name, ret})
		}
	}
	floorAssetsCapture(t, "facades", rows, "c46dc7c97e5d008e43ea19709c8e79cbc9a72547e421b1e5a02427388ad7ad93")
}
