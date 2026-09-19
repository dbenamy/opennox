//go:build porttest

package legacy

import (
	"fmt"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
)

// The fixture observes the delayed-deletion boundary; the selected command owns
// item selection and request order. Existing inventory tests cover deletion.
type PortTestScriptStartupSpec struct {
	Order   []int
	Chapter uint32
	Journal []uint16
}

func (p *portTestShopPools) scriptStartupContract() []uint32 {
	sp := p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives.Attack.Controls.ScriptStartup
	if sp == nil {
		panic("missing startup contract")
	}
	host := p.proxy.core.Players.ByInd(31).PlayerUnit
	oldHead := host.InvFirstItem
	defer func() { host.InvFirstItem = oldHead }()
	for _, u := range p.proxy.life.players {
		if u.UpdateDataPlayer().Player.Journal != nil {
			panic("startup fixture requires empty journals")
		}
	}
	for i := range p.proxy.life.players {
		pl := p.proxy.life.players[i].UpdateDataPlayer().Player
		for n, flags := range sp.Journal {
			journalAdd(pl, fmt.Sprintf("startup-%d", n), flags)
		}
		defer func() {
			for pl.Journal != nil {
				journalUnlink(pl, pl.Journal)
			}
		}()
	}
	host.InvFirstItem = nil
	var prev *server.Object
	var wantDeletes []uint32
	for _, i := range sp.Order {
		u := p.items[i].u
		if prev == nil {
			host.InvFirstItem = u
		} else {
			prev.InvNextItem = u
		}
		u.InvNextItem = nil
		prev = u
		if u.ObjClass&0x40 != 0 {
			wantDeletes = append(wantDeletes, 32, p.proxy.life.ids[uint32(uintptr(u.CObj()))])
		}
	}
	chapter := memmap.PtrUint32(0x5D4594, 2386832)
	oldChapter := *chapter
	defer func() { *chapter = oldChapter }()
	*chapter = sp.Chapter
	start := len(p.proxy.trace)
	Nox_script_StartupScreen_516600_A()
	if *chapter != 1 {
		panic("startup did not set chapter flag")
	}
	gotDeletes := p.proxy.trace[start:]
	if len(gotDeletes) != len(wantDeletes) {
		panic(fmt.Sprintf("startup deletion count %v want %v", gotDeletes, wantDeletes))
	}
	for i, v := range wantDeletes {
		if gotDeletes[i] != v {
			panic("startup deletion order")
		}
	}
	out := []uint32{*chapter, uint32(len(wantDeletes) / 2)}
	for i := range p.proxy.life.players {
		u := &p.proxy.life.players[i]
		var want []uint16
		for n := len(sp.Journal) - 1; n >= 0; n-- {
			f := sp.Journal[n]
			if u != host || f&0xE == 0 {
				want = append(want, f)
			}
		}
		n := u.UpdateDataPlayer().Player.Journal
		out = append(out, uint32(len(want)))
		for _, f := range want {
			if n == nil || n.Field3 != f {
				panic("startup journal mask/other player mismatch")
			}
			out = append(out, uint32(f))
			n = n.Next
		}
		if n != nil {
			panic("startup journal extra entry")
		}
	}
	return out
}
