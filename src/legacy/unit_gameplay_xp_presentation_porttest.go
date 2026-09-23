//go:build porttest

package legacy

import (
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"unsafe"
)

func (p *portTestShopPools) unitExperiencePresentation(sp *PortTestUnitExperienceSpec) (func(), func()) {
	if !sp.Presentation {
		return func() {}, func() {}
	}
	a, b, c := dword_5d4594_2523804, dword_5d4594_2523780, dword_5d4594_2523776
	dword_5d4594_2523804, dword_5d4594_2523780, dword_5d4594_2523776 = 0, 0, 0
	offsets := []uintptr{2523772, 2523796, 2523800}
	old := make([]uint32, len(offsets))
	for i, off := range offsets {
		old[i] = *memmap.PtrUint32(0x5D4594, off)
		*memmap.PtrUint32(0x5D4594, off) = 0
	}
	ticks := memmap.PtrUint64(0x5D4594, 2523788)
	oldTicks := *ticks
	*ticks = 0
	game := noxflags.GetGame()
	noxflags.UnsetGame(noxflags.GamePause)
	var missing []string
	if sp.MissingPresentation {
		missing = []string{"LevelUp"}
	}
	restoreTypes := p.proxy.core.PortTestRewardTypes([]string{"LevelUp"}, missing, true, 0, 0)
	before := len(p.proxy.life.created)
	check := func() {
		want := 1
		if sp.MissingPresentation {
			want = 0
		}
		if len(p.proxy.life.created)-before != want {
			panic("XP presentation factory count")
		}
		if dword_5d4594_2523804 != 1 || uintptr(dword_5d4594_2523780) != uintptr(p.resources.unit.CObj()) || !noxflags.HasGame(noxflags.GamePause) || *ticks != 10000 || memmap.Uint32(0x5D4594, 2523796) != 5000 {
			panic("XP presentation state")
		}
		if want == 1 && uintptr(dword_5d4594_2523776) != uintptr(unsafe.Pointer(p.proxy.life.created[before])) {
			panic("XP presentation object")
		}
	}
	restore := func() {
		restoreTypes()
		noxflags.ResetGame()
		noxflags.SetGame(game)
		dword_5d4594_2523804, dword_5d4594_2523780, dword_5d4594_2523776 = a, b, c
		for i, off := range offsets {
			*memmap.PtrUint32(0x5D4594, off) = old[i]
		}
		*ticks = oldTicks
	}
	return check, restore
}
