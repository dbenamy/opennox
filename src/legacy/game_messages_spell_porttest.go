//go:build porttest

package legacy

import (
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"unsafe"
)

// The test supplies the expected admission decision independently of the dispatcher.
type PortTestServerSpellSpec struct {
	IDs         [5]int32
	Warnings    byte
	Insert      bool
	Count, Self int32
	WantBook    bool
}

func (p *portTestShopPools) portTestServerSpellDirect() {
	a := p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives.Attack
	sp := a.Controls.GameMessageSpell
	u := p.temporaryRef(a.Actor)
	for i, key := range []string{"GeneralPrint:NoSpellWarningGeneral", "GeneralPrint:ConjureNoSpellWarning1"} {
		if sp.Warnings&(1<<i) != 0 {
			gameplayTextPrivate(u, alloc.InternCString(key), 0)
		}
	}
	if !sp.Insert {
		return
	}
	ids, free := alloc.New([5]int32{})
	defer free()
	*ids = sp.IDs
	if spellLifeInsertBook(u, unsafe.Pointer(ids), sp.Count, 3, sp.Self) == 0 && sp.Count == 1 {
		for _, id := range sp.IDs {
			if id != 0 {
				gameplayReportSpellStat(int(u.UpdateDataPlayer().Player.PlayerInd), uint32(id), 0)
			}
		}
	}
}

func (p *portTestShopPools) portTestServerSpellCheck() {
	a := p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives.Attack
	sp := a.Controls.GameMessageSpell
	b := spellLifeBookHead()
	if (b != nil) != sp.WantBook {
		panic("spell request queue admission differs from contract")
	}
	if b == nil {
		return
	}
	if b.Owner != p.temporaryRef(a.Actor) || b.Next != nil || b.Prev != nil || b.Delay != 3 || b.Self != sp.Self || b.Frame != p.proxy.core.Frame() {
		panic("spell request queue metadata differs from contract")
	}
	for i, id := range b.Spells {
		want := int32(0)
		if int32(i) < sp.Count {
			want = sp.IDs[i]
		}
		if id != want {
			panic("spell request queue spell order differs from contract")
		}
	}
	// Give the newly allocated queue node a stable identity in the existing snapshot.
	p.identify(unsafe.Pointer(b), 9200)
}
