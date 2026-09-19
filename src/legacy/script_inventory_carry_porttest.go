//go:build porttest

package legacy

/*
extern int dword_5d4594_2386848;
extern unsigned int dword_5d4594_2386852;
*/
import "C"
import (
	"fmt"
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

type PortTestScriptCarrySpec struct {
	FreeCells int
	Count     byte
	Class     uint32
	SameType  bool
	Reserved  int32
	Notice    uint32
	Capacity  int
	WantItem  int // -1 means no drop attempt
	Repeats   int
}

func (p *portTestShopPools) scriptCarryContract() []uint32 {
	sp := p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives.Attack.Controls.ScriptCarry
	if sp == nil {
		panic("missing carry contract")
	}
	grid := uiInventoryGrid()
	oldGrid := append([]uiInventoryCell(nil), grid...)
	defer func() { copy(grid, oldGrid) }()
	dr, free := alloc.New(client.Drawable{})
	defer free()
	target := p.items[2].u
	dr.TypeIDVal = uint32(target.TypeInd)
	if !sp.SameType {
		dr.TypeIDVal += 1000
	}
	dr.ObjClass = object.Class(sp.Class)
	for i := range grid {
		grid[i] = uiInventoryCell{Drawable: dr, Count: sp.Count}
	}
	left := sp.FreeCells
	for row := 0; row < 20; row++ {
		for col := 0; col < 4; col++ {
			if left > 0 {
				grid[row+21*col].Count = 0
				left--
			}
		}
	}
	if got := int(PortTestInventoryTransaction(25, uintptr(target.TypeInd), 1, 0, 0)); got != sp.Capacity {
		panic(fmt.Sprintf("carry capacity %d want %d", got, sp.Capacity))
	}
	reserved, notice := C.dword_5d4594_2386848, C.dword_5d4594_2386852
	defer func() { C.dword_5d4594_2386848, C.dword_5d4594_2386852 = reserved, notice }()
	C.dword_5d4594_2386848 = C.int(sp.Reserved)
	C.dword_5d4594_2386852 = C.uint(sp.Notice)
	glyph := memmap.PtrUint32(0x5D4594, 2386856)
	oldGlyph := *glyph
	defer func() { *glyph = oldGlyph }()
	*glyph = 0
	var out []uint32
	for repeat := 0; repeat < sp.Repeats; repeat++ {
		before := portTestInventoryDropCalls()
		rng := p.proxy.core.Rand.Logic.Index()
		Nox_xxx_playerCanCarryItem_513B00(p.resources.unit, target)
		after := portTestInventoryDropCalls()
		want := sp.WantItem
		wantNotice := sp.Notice
		if want >= 0 {
			if len(after) != len(before)+6 {
				panic("carry missing drop callback")
			}
			call := after[len(before):]
			if call[0] != 2 || call[1] != uint32(uintptr(p.resources.unit.CObj())) || call[2] != uint32(uintptr(p.items[want].u.CObj())) {
				panic("carry chose wrong item/owner")
			}
			if wantNotice == 0 {
				wantNotice = 1
			}
		} else if len(after) != len(before) || p.proxy.core.Rand.Logic.Index() != rng {
			panic("carry unexpected drop/RNG")
		}
		if uint32(C.dword_5d4594_2386852) != wantNotice || int32(C.dword_5d4594_2386848) != sp.Reserved {
			panic("carry counters")
		}
		out = append(out, uint32(sp.Capacity), uint32(C.dword_5d4594_2386852), uint32(p.proxy.core.Rand.Logic.Index()))
	}
	if *glyph != uint32(p.proxy.core.Types.IndByID("Glyph")) {
		panic("carry Glyph type cache")
	}
	return out
}
