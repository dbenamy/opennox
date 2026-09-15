//go:build porttest

package opennox

import (
	"testing"
	"unsafe"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
)

func TestScoreboardRanksAndColors(t *testing.T) {
	o := newScoreboardOwner(t)
	var rows []scoreboardResult
	scores := []uint32{0, 1, 1, 0x7fffffff, 0x80000000, 0xfffffffe, 0xffffffff}
	for _, elimination := range []bool{false, true} {
		for count := 0; count <= len(scores); count++ {
			o.resetRank(t)
			if elimination {
				noxflags.SetGame(noxflags.GameModeElimination)
			}
			*memmap.PtrUint8(0x5D4594, 1090116) = byte(count)
			*memmap.PtrUint8(0x5D4594, 1090117) = byte(count)
			for i := 0; i < count; i++ {
				*memmap.PtrUint32(0x5D4594, 1084192+uintptr(i*80)) = uint32(100 + i)
				*memmap.PtrUint32(0x5D4594, 1084196+uintptr(i*80)) = scores[i]
				*memmap.PtrUint32(0x5D4594, 1087252+uintptr(i*56)) = scores[i]
			}
			for _, score := range scores {
				want := uint32(1)
				for _, s := range scores[:count] {
					if (!elimination && s > score) || (elimination && s < score) {
						want++
					}
				}
				got := legacy.PortTestScoreboard(23, uintptr(score), 0, 0)
				if got != want {
					t.Fatalf("team rank count%d elimination%v score%x got%d want%d", count, elimination, score, got, want)
				}
				rows = append(rows, o.rankCapture(t, 23, got))
			}
			// The actual local-player owner uses id 100; absent rows return zero.
			want := uint32(0)
			if count > 0 {
				want = 1
				for _, s := range scores[:count] {
					if (!elimination && s > scores[0]) || (elimination && s < scores[0]) {
						want++
					}
				}
			}
			got := legacy.PortTestScoreboard(22, 0, 0, 0)
			if got != want {
				t.Fatalf("local rank count%d got%d want%d", count, got, want)
			}
			rows = append(rows, o.rankCapture(t, 22, got))
		}
	}
	scoreboardCapture(t, "ranks", rows, "cc51ea6ee9d0c7a3b5cdd5bda483f3c9ef0cdb1ee60550570a9538f8ff1459f4")
	rows = nil
	palette := unsafe.Slice((*byte)(memmap.PtrOff(0x587000, 145584)), 80)
	for index := 0; index < 9; index++ {
		for color := 0; color < 256; color++ {
			*memmap.PtrUint8(0x5D4594, 1087256+uintptr(index*56)) = byte(color)
			got := legacy.PortTestScoreboard(21, uintptr(index), 0, 0)
			want := uint32(palette[(color%10)*8])
			if got != want {
				t.Fatalf("team color index%d color%d got%d want%d", index, color, got, want)
			}
			// Every byte is independently checked; capture one complete palette cycle
			// and wraparound boundaries to keep redundant pixel snapshots bounded.
			if color < 11 || color == 127 || color == 128 || color == 254 || color == 255 {
				rows = append(rows, o.rankCapture(t, 21, got))
			}
		}
	}
	scoreboardCapture(t, "team-colors", rows, "29aee868eb085d296cd71850940fc423df08605261935c7c509ec866ff2bd281")
}
