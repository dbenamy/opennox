//go:build porttest

package opennox

import (
	"github.com/opennox/libs/prand"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
	"unsafe"
)

func TestCharacterCreationQuickbar(t *testing.T) {
	q := newQuickbarOwner(t)
	player, free := alloc.Calloc(1, 128)
	defer free()
	defer legacy.PortTestCharacterClassOwner(player)()
	data, free := alloc.Make([]byte{}, 1280)
	defer free()
	table := unsafe.Slice(memmap.PtrPtr(0x587000, 170156), 3)
	old := append([]unsafe.Pointer(nil), table...)
	defer copy(table, old)
	for i := range table {
		table[i] = unsafe.Pointer(&data[0])
	}
	type row struct {
		Class, Count, Mode, Seed, Random int
		Slots                            []uint32
		Selected                         uint32
	}
	var rows []row
	for class := 0; class < 3; class++ {
		for _, count := range []int{0, 1, 3, 255} {
			for _, mode := range []int{0, 1, 2} {
				for _, seed := range []int{1, 17, 63} {
					q.reset(t)
					*(*byte)(unsafe.Add(player, 66)) = byte(class)
					clear(data)
					for i := 0; i < count; i++ {
						for j := 0; j < 5; j++ {
							data[5*i+j] = byte(1 + (i+j)%5)
						}
					}
					for i := 0; i < 25; i++ {
						q.bar[2*i] = 99
						q.bar[2*i+1] = 0xa5a5a5a5
					}
					q.c.srv.Rand.Other.Reset(seed)
					logic := q.c.srv.Rand.Logic.Index()
					rng := prand.New(seed)
					selected := 0
					if count > 0 {
						selected = rng.Int(0, count-1)
					}
					legacy.PortTestCharacterQuickbar(mode)
					if q.c.srv.Rand.Other.Index() != rng.Index() || q.c.srv.Rand.Logic.Index() != logic {
						t.Fatal("RNG", class, count, mode, seed)
					}
					for i := 0; i < 25; i++ {
						want := uint32(99)
						if count > 0 && (class != 0 || i < 5) {
							want = uint32(data[5*selected+i%5])
							if mode == 1 {
								want = 0
							}
						}
						if q.bar[2*i] != want {
							t.Fatalf("class%d count%d mode%d slot%d=%d want%d", class, count, mode, i, q.bar[2*i], want)
						}
					}
					if q.bar[50] != 0 {
						t.Fatal("selected row must end at zero")
					}
					rows = append(rows, row{class, count, mode, seed, q.c.srv.Rand.Other.Index(), append([]uint32(nil), q.bar[:50]...), q.bar[50]})
				}
			}
		}
	}
	spellbookCapture(t, "character-creation-quickbar", rows, "953a75e4c3e882a572759feca49cbfca8472ae3a88d853ec3572abd77a802e14")
}
