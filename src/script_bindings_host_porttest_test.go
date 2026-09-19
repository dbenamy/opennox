//go:build porttest

package opennox

import (
	"testing"

	"github.com/opennox/noxscript/ns/asm"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

func TestScriptBindingsHostPredicates(t *testing.T) {
	o := newReliableReportsOwner(t)
	o.s.NoxScriptVM.Init(o.s)
	noxServer.noxScript.Init(noxServer)
	p := o.s.Players.ByInd(31)
	ud := p.PlayerUnit.UpdateDataPlayer()
	savedP, savedUD := *p, *ud
	t.Cleanup(func() { *p = savedP; *ud = savedUD })
	var trade server.TradeSession
	for _, active := range []bool{false, true} {
		for _, talk := range []bool{false, true} {
			for _, trading := range []bool{false, true} {
				for _, root := range []bool{false, true} {
					for _, fi := range []asm.Builtin{asm.BuiltinIsTalking, asm.BuiltinIsTrading} {
						p.Active = 0
						if active {
							p.Active = 1
						}
						ud.DialogWith = nil
						if talk {
							ud.DialogWith = &o.units[0]
						}
						ud.Trade70 = nil
						if trading {
							ud.Trade70 = &trade
						}
						const sentinel = 0x2468ace0
						o.s.NoxScriptVM.PushU32(sentinel)
						if root {
							if err := noxServer.noxScript.callBuiltinNative(fi); err != nil {
								t.Fatal(err)
							}
						} else if r, ok := legacy.CallScriptBuiltin(fi); r != 0 || !ok {
							t.Fatal("host builtin dispatch", r, ok)
						}
						want := uint32(0)
						if active && (fi == asm.BuiltinIsTalking && talk || fi == asm.BuiltinIsTrading && trading) {
							want = 1
						}
						if got := o.s.NoxScriptVM.PopU32(); got != want {
							t.Fatalf("host predicate %d active%v talk%v trade%v root%v: %d != %d", fi, active, talk, trading, root, got, want)
						}
						if o.s.NoxScriptVM.PopU32() != sentinel {
							t.Fatal("host predicate changed surrounding stack")
						}
					}
				}
			}
		}
	}
}
