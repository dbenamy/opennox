//go:build porttest

package opennox

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"

	"github.com/opennox/libs/spell"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

// Exercise the actual duration allocator/list and root create/update/destroy
// consumer with replaceable wall hooks, preserving the original C callback path.
func TestDurationIdentityLifecycle(t *testing.T) {
	for mask := 0; mask < 8; mask++ {
		for _, createWord := range []uint32{0, 1, 0x80000000, 0xffffffff} {
			for _, updateWord := range []uint32{0, 1, 0x80000000, 0xffffffff} {
				for _, deadline := range []bool{false, true} {
					t.Run(fmt.Sprintf("mask%d/create%x/update%x/deadline%v", mask, createWord, updateWord, deadline), func(t *testing.T) {
						core := newObjectXferOwner(t)
						u := newObjectXferSimple(t, core)
						u.ObjClass = 0
						u.ObjFlags = 0
						u.PosVec = types.Pointf{X: 12, Y: 34}
						core.SetFrame(100)
						t.Cleanup(core.PortTestSpellLifecycle([]server.PortTestSpellLifecycleDef{{Index: int(spell.SPELL_SHIELD), Valid: true, Enabled: true}}, nil))
						var duration spellsDuration
						duration.Init(noxServer)
						oldCreate, oldUpdate, oldDestroy := legacy.Nox_xxx_spellWallCreate_4FFA90, legacy.Nox_xxx_spellWallUpdate_500070, legacy.Nox_xxx_spellWallDestroy_500080
						t.Cleanup(func() {
							legacy.Nox_xxx_spellWallCreate_4FFA90 = oldCreate
							legacy.Nox_xxx_spellWallUpdate_500070 = oldUpdate
							legacy.Nox_xxx_spellWallDestroy_500080 = oldDestroy
						})
						var create, update, destroy unsafe.Pointer
						if mask&1 != 0 {
							create = legacy.Get_nox_xxx_spellWallCreate_4FFA90()
						}
						if mask&2 != 0 {
							update = legacy.Get_nox_xxx_spellWallUpdate_500070()
						}
						if mask&4 != 0 {
							destroy = legacy.Get_nox_xxx_spellWallDestroy_500080()
						}
						var live, created *server.DurSpell
						var calls []string
						argumentsOK := true
						legacy.Nox_xxx_spellWallCreate_4FFA90 = func(p *server.DurSpell) int { calls = append(calls, "create"); created = p; return int(createWord) }
						legacy.Nox_xxx_spellWallUpdate_500070 = func(p *server.DurSpell) int {
							calls = append(calls, "update")
							argumentsOK = argumentsOK && p == live && p.Flags88&1 == 0
							return int(updateWord)
						}
						legacy.Nox_xxx_spellWallDestroy_500080 = func(p *server.DurSpell) {
							calls = append(calls, "destroy")
							argumentsOK = argumentsOK && p == live && p.Flags88&1 != 0
						}
						dt := uint32(0)
						if deadline {
							dt = 3
						}
						sa := server.SpellAcceptArg{Pos: types.Pointf{X: 23, Y: 45}}
						accepted := duration.New(spell.SPELL_SHIELD, u, u, nil, &sa, 3, create, update, destroy, dt)
						wantAccepted := create == nil || createWord == 0
						live = duration.List
						if accepted != wantAccepted || live == nil {
							t.Fatal("create/list result", accepted, wantAccepted, live)
						}
						if create != nil && created != live {
							t.Fatal("create callback record")
						}
						if live.Create != create || live.Update != update || live.Destroy != destroy || live.Caster16 != u || live.Obj12 != u || live.Target48 != nil || live.Level != 3 || live.Spell != uint32(spell.SPELL_SHIELD) || live.Frame60 != 100 || live.Frame68 != 100+dt || live.Pos != u.PosVec || live.Pos2 != sa.Pos {
							t.Fatal("duration initialization")
						}
						var want []string
						if create != nil {
							want = append(want, "create")
						}
						if !reflect.DeepEqual(calls, want) {
							t.Fatal("creation dispatch", calls, want)
						}
						if (live.Flags88&1 != 0) != (!wantAccepted) {
							t.Fatal("creation cancellation")
						}
						if deadline {
							core.SetFrame(103)
						}
						collisionRegistryGrow(128)
						duration.spellCastByPlayer()
						if wantAccepted {
							if !deadline && update != nil {
								want = append(want, "update")
							}
							cancelled := deadline || (update != nil && updateWord != 0)
							if duration.List != live || (live.Flags88&1 != 0) != cancelled {
								t.Fatal("update cancellation/list", cancelled)
							}
							if !cancelled {
								duration.CancelSpell(live)
							}
							duration.spellCastByPlayer()
						}
						if destroy != nil {
							want = append(want, "destroy")
						}
						if duration.List != nil || !argumentsOK || !reflect.DeepEqual(calls, want) {
							t.Fatal("duration lifecycle", argumentsOK, calls, want)
						}
						duration.spellCastByPlayer()
						duration.onNewSpell()
						if !reflect.DeepEqual(calls, want) {
							t.Fatal("callback repeated after removal")
						}
					})
				}
			}
		}
	}
}
