//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"image"
	"testing"
)

func combatEffectPool(t *testing.T) alloc.ClassT[client.DrawableFX] {
	t.Helper()
	clear(serverConfigOwnBytes(t, 0x5D4594, 1203868, 4))
	if nox_xxx_allocArrayDrawableFX_495AB0() != 1 {
		t.Fatal("effect allocation class initialization")
	}
	cl := alloc.AsClassT[client.DrawableFX](*memmap.PtrPtr(0x5D4594, 1203868))
	head := memmap.PtrPtr(0x5D4594, 1203872)
	t.Cleanup(func() {
		for i := 0; i < 128 && *head != nil; i++ {
			fx := (*client.DrawableFX)(*head)
			legacy.PortTestCombatEffectDetach(fx.C())
			cl.FreeObjectFirst(fx)
		}
		cl.Free()
	})
	return cl
}
func TestCombatOverlayEffectGoCallers(t *testing.T) {
	type record struct {
		Operation string
		Complete  bool
	}
	var rows []record
	for _, op := range []string{"membership", "cleanup"} {
		t.Run(op, func(t *testing.T) {
			o := newCombatOverlayOwner(t)
			cl := combatEffectPool(t)
			dr := o.drawable(7, image.Pt(300, 300))
			for _, kind := range []uint32{1, 2} {
				fx := cl.NewObject()
				if fx == nil {
					t.Fatal("effect allocation")
				}
				fx.Field0 = kind
				legacy.PortTestCombatEffectAttach(fx.C(), dr)
			}
			switch op {
			case "membership":
				if !dr.HasFX(1) || !dr.HasFX(2) || dr.HasFX(3) {
					t.Fatal("Go membership must follow the C attachment chain")
				}
			case "cleanup":
				o.c.sub_495B00(dr)
				if dr.Field_114 != nil || *memmap.PtrPtr(0x5D4594, 1203872) != nil {
					t.Fatal("Go cleanup left an attached effect")
				}
			}
			rows = append(rows, record{op, true})
		})
	}
	spellbookCapture(t, "combat-overlay-effect-go-callers", rows, "c7168259e8edf0fe765a6cb133766455e0d9b4f9b74946d25b058614e835a0ac")
}
