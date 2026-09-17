//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
	"unsafe"
)

func TestObjectReportsTypeCache(t *testing.T) {
	s := newCreatureXferOwner(t)
	names := []string{"TeleportPentagram", "PressurePlate", "Spike", "PeriodicSpike"}
	caches := unsafe.Slice(memmap.PtrUint32(0x5D4594, 2386972), 4)
	old := append([]uint32(nil), caches...)
	t.Cleanup(func() { copy(caches, old) })
	var rows []struct {
		Name   string
		Return int
		Words  [4]uint32
	}
	defer func() {
		spellbookCapture(t, "object-reports-cache", rows, "58d52abc1990f0f8d8fac7c463705b73287b8e17f6b70d00bad2f0fb71cb51ac")
	}()
	for missing := -1; missing < 4; missing++ {
		t.Run(fmt.Sprint(missing), func(t *testing.T) {
			var absent []string
			if missing >= 0 {
				absent = []string{names[missing]}
			}
			t.Cleanup(s.PortTestAttackTypes(0, absent, names...))
			for i := range caches {
				caches[i] = 0xaabbccdd + uint32(i)
			}
			rv := legacy.PortTestObjectReports(0, nil, nil, 0, 0, 0, nil)
			want := 1
			if missing >= 0 {
				want = 0
			}
			if rv != want {
				t.Fatalf("return%d want%d", rv, want)
			}
			for i, name := range names {
				expected := uint32(s.Types.IndByID(name))
				if missing >= 0 && i > missing {
					expected = 0xaabbccdd + uint32(i)
				}
				if caches[i] != expected {
					t.Fatalf("cache%d got%x want%x", i, caches[i], expected)
				}
			}
			rows = append(rows, struct {
				Name   string
				Return int
				Words  [4]uint32
			}{fmt.Sprint(missing), rv, [4]uint32(caches)})
		})
	}
}
