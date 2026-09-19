//go:build porttest

package opennox

import (
	"math"
	"testing"

	"github.com/opennox/libs/strman"
	"github.com/opennox/noxscript/ns/asm"
	"github.com/opennox/opennox/v1/legacy"
)

// The unchanged C experience/level helper remains outside this binding batch.
// Exercise actual player mutations and reports below the next-level threshold.
func TestScriptBindingsExperience(t *testing.T) {
	o := newWorldCollisionOwner(t)
	language, restoreStrings := o.s.PortTestMeterStrings(strman.Entry{ID: "health.c:gainpoints", Vals: []strman.Variant{{Str: "Gained %u points"}}})
	t.Cleanup(restoreStrings)
	language(0)
	t.Cleanup(o.s.PortTestInventoryDisplayBalance())
	serverConfigOwnBytes(t, 0x5D4594, 2386620, 208)
	_, restore := legacy.PortTestMonsterCacheInit()
	t.Cleanup(restore)
	u := &o.units[0]
	p := u.UpdateDataPlayer().Player
	savedU, savedP := *u, *p
	t.Cleanup(func() { *u = savedU; *p = savedP })
	type row struct {
		Root, Missing          bool
		Initial, Delta, Result uint32
		Queue                  legacy.PortTestReliableReportState
	}
	var rows []row
	for _, root := range []bool{false, true} {
		for _, missing := range []bool{false, true} {
			for _, initial := range []float32{10, 1024.25, 8192} {
				for _, delta := range []float32{0, 0.125, -0.5, 1.75, 17.0625} {
					o.reset()
					*u = savedU
					*p = savedP
					p.Level = 9
					u.Experience = initial
					u.ScriptIDVal = 100125
					u.ObjFlags = 0
					legacy.PortTestMonsterCache("reset", nil, 0)
					legacy.PortTestMonsterCache("prepare", u, 0)
					id := uint32(u.ScriptIDVal)
					if missing {
						id = 999
					}
					const sentinel = 0x2468ace0
					o.s.NoxScriptVM.PushU32(sentinel)
					o.s.NoxScriptVM.PushU32(id)
					o.s.NoxScriptVM.PushU32(math.Float32bits(delta))
					if root {
						if err := noxServer.noxScript.callBuiltinNative(asm.BuiltinGiveXp); err != nil {
							t.Fatal(err)
						}
					} else if r, ok := legacy.CallScriptBuiltin(asm.BuiltinGiveXp); r != 0 || !ok {
						t.Fatal("experience builtin dispatch", r, ok)
					}
					want := initial
					if !missing {
						want = float32(float64(initial) + float64(delta))
					}
					if u.Experience != want || p.Level != 9 || o.s.NoxScriptVM.PopU32() != sentinel {
						t.Fatal("experience value/level/stack mismatch", root, missing, initial, delta, u.Experience, p.Level)
					}
					rows = append(rows, row{root, missing, math.Float32bits(initial), math.Float32bits(delta), math.Float32bits(u.Experience), o.state()})
				}
			}
		}
	}
	spellbookCapture(t, "script-bindings-experience", rows, "bdeebff94f19a8d36848066a3ee9f5dc6c94d8e8af11272815ffd338b2c1c4c5")
}
