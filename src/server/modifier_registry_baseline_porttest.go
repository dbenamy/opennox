//go:build porttest

package server

import (
	"fmt"
	"reflect"
	"unsafe"
)

// PortTestModifierRegistryLookup exposes the actual startup registry entry for
// an independent original-C key comparison.
func PortTestModifierRegistryLookup(kind, name string) (unsafe.Pointer, string, bool) {
	var table map[string]modFuncs
	switch kind {
	case "damage":
		table = modDamageEffects
	case "defend":
		table = modDefendEffects
	case "update":
		table = modUpdateEffects
	case "engage":
		table = modEngageEffects
	default:
		return nil, "", false
	}
	e, ok := table[name]
	if !ok {
		return nil, "", false
	}
	parser := ""
	switch {
	case e.Parse == nil:
		parser = "none"
	case reflect.ValueOf(e.Parse).Pointer() == reflect.ValueOf(ModifierParseFunc(ModEffectParseInt)).Pointer():
		parser = "int"
	case reflect.ValueOf(e.Parse).Pointer() == reflect.ValueOf(ModifierParseFunc(ModEffectParseFloat)).Pointer():
		parser = "float"
	default:
		parser = "other"
	}
	return e.Func, parser, true
}

// PortTestModifierParseEffect executes the real parser against one of the eight
// actual ModifierEff callback slots. slot offsets are the C layout offsets.
func PortTestModifierParseEffect(slot int, mod *ModifierEff, value string) error {
	var table map[string]modFuncs
	var text string
	var targ ModParseTarg
	switch slot {
	case 40:
		table, text = modDamageEffects, "damage"
		targ = ModParseTarg{&mod.Attack40.Fnc, &mod.Attack40.Valf, &mod.Attack40.Val}
	case 52:
		table, text = modDamageEffects, "damage"
		targ = ModParseTarg{&mod.AttackPreHit52.Fnc, &mod.AttackPreHit52.Valf, &mod.AttackPreHit52.Val}
	case 64:
		table, text = modDamageEffects, "damage"
		targ = ModParseTarg{&mod.AttackPreDmg64.Fnc, &mod.AttackPreDmg64.Valf, &mod.AttackPreDmg64.Val}
	case 76:
		table, text = modDefendEffects, "defend"
		targ = ModParseTarg{&mod.Defend76.Fnc, &mod.Defend76.Valf, &mod.Defend76.Val}
	case 88:
		table, text = modDefendEffects, "defend"
		targ = ModParseTarg{&mod.DefendCollide88.Fnc, &mod.DefendCollide88.Valf, &mod.DefendCollide88.Val}
	case 100:
		table, text = modUpdateEffects, "update"
		targ = ModParseTarg{&mod.Update100.Fnc, &mod.Update100.Valf, &mod.Update100.Val}
	case 112:
		table, text = modEngageEffects, "engage"
		targ = ModParseTarg{&mod.Engage112, &mod.EngageFloat120, &mod.EngageInt124}
	case 116:
		table, text = modEngageEffects, "engage"
		targ = ModParseTarg{&mod.Disengage116, &mod.DisengageFloat128, &mod.DisengageInt132}
	default:
		return fmt.Errorf("unknown modifier callback slot %d", slot)
	}
	return modParseEffect(text, table, mod, targ, value)
}
