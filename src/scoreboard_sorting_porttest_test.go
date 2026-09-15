//go:build porttest

package opennox

import (
	"encoding/binary"
	"sort"
	"testing"
	"unsafe"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
)

func TestScoreboardPlayerOrdering(t *testing.T) {
	o := newScoreboardOwner(t)
	var rows []scoreboardResult
	for _, mode := range []int{2, 5} {
		for _, elimination := range []bool{false, true} {
			for count := 0; count <= 32; count++ {
				for _, state := range []uint32{0, 1, 0x21} {
					o.resetRank(t)
					*o.rankWords["dword_5d4594_1090120"] = uint32(mode)
					if elimination {
						noxflags.SetGame(noxflags.GameModeElimination)
					}
					want := make([]int, count)
					// Deliberate ties and mixed signed scores in ordinary mode; ascending elimination
					// scores stay in their valid nonnegative domain, below C's sort sentinels.
					for i := 0; i < count; i++ {
						p := &o.players[i]
						p.Active = 1
						p.PlayerInd = byte(i)
						p.Field3680 = state
						p.Lessons = int32(i%7 - 3)
						p.Field2140 = uint32(i % 7)
						p.Field2108 = uint32(i % 7)
						*(*uint16)(unsafe.Add(unsafe.Pointer(p), 2148)) = uint16(i * 2077)
						*(*byte)(unsafe.Add(unsafe.Pointer(p), 2251)) = byte(i % 3)
						want[i] = i
					}
					sort.SliceStable(want, func(a, b int) bool {
						pa, pb := &o.players[want[a]], &o.players[want[b]]
						if mode == 5 {
							a, b := pa.Field2108, pb.Field2108
							if state != 1 {
								if a == 0 {
									a = 0x8000000
								}
								if b == 0 {
									b = 0x8000000
								}
							}
							return a < b
						}
						if elimination {
							return pa.Field2140 < pb.Field2140
						}
						return pa.Lessons > pb.Lessons
					})
					ret := legacy.PortTestScoreboard(11, 0, 0, 0)
					if got := *memmap.PtrUint8(0x5D4594, 1090117); int(got) != count {
						t.Fatalf("mode%d elimination%v count%d got%d", mode, elimination, count, got)
					}
					seen := map[uint32]bool{}
					for row, ind := range want {
						b := unsafe.Slice((*byte)(memmap.PtrOff(0x5D4594, 1084132+uintptr(row*80))), 80)
						id := binary.LittleEndian.Uint32(b[60:])
						p := &o.players[ind]
						if id != p.NetCodeVal || seen[id] {
							t.Fatalf("mode%d elimination%v count%d state%x row%d id%d want%d", mode, elimination, count, state, row, id, p.NetCodeVal)
						}
						seen[id] = true
						score := uint32(p.Lessons)
						if elimination {
							score = p.Field2140
						}
						if mode == 5 {
							score = p.Field2108
						}
						if binary.LittleEndian.Uint32(b[64:]) != score || binary.LittleEndian.Uint32(b[68:]) != uint32(uint16(ind*2077)) || b[56] != byte(ind%3) || binary.LittleEndian.Uint32(b[76:]) != state {
							t.Fatalf("row fields mode%d elimination%v count%d row%d", mode, elimination, count, row)
						}
						if p.Field2108 != uint32(ind%7) || p.Field2140 != uint32(ind%7) || p.Lessons != int32(ind%7-3) {
							t.Fatal("temporary sorting score was not restored")
						}
					}
					rows = append(rows, o.rankCapture(t, 11, ret))
				}
			}
		}
	}
	scoreboardCapture(t, "player-ordering", rows, "aa3027374c1862166c8884080d31b3cf734934956e72a27d2759b152ba5fac8c")
}
