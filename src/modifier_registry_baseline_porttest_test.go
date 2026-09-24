//go:build porttest

package opennox

import (
	"runtime"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func TestModifierRegistryKeysAndLifetime(t *testing.T) {
	type expected struct{ symbol, name, kind, parser string }
	want := []expected{
		{"nox_xxx_effectDamageMultiplier_4E04C0", "DamageMultiplierEffect", "damage", "float"},
		{"nox_xxx_stunEffect_4E04D0", "StunEffect", "damage", "int"},
		{"nox_xxx_fireEffect_4E0550", "FireEffect", "damage", "float"},
		{"nox_xxx_fireRingEffect_4E05B0", "FireRingEffect", "damage", "int"},
		{"nox_xxx_blueFREffect_4E05F0", "BlueFireRingEffect", "damage", "int"},
		{"nullsub_38", "FrostEffect", "damage", "int"},
		{"nox_xxx_recoilEffect_4E0640", "RecoilEffect", "damage", "float"},
		{"nox_xxx_confuseEffect_4E0670", "ConfuseEffect", "damage", "int"},
		{"nox_xxx_lightngEffect_4E06F0", "LightningEffect", "damage", "float"},
		{"nox_xxx_drainMEffect_4E0740", "DrainManaEffect", "damage", "float"},
		{"nox_xxx_vampirismEffect_4E07C0", "VampirismEffect", "damage", "float"},
		{"nox_xxx_poisonEffect_4E0850", "PoisonEffect", "damage", "int"},
		{"nullsub_39", "PanicEffect", "damage", "int"},
		{"nox_xxx_sympathyEffect_4E08E0", "SympathyEffect", "damage", "float"},
		{"nullsub_22", "ReadinessEffect", "damage", "int"},
		{"nox_xxx_effectProjectileSpeed_4E09B0", "ProjectileSpeedEffect", "damage", "float"},
		{"nullsub_36", "ReplenishmentEffect", "damage", "int"},
		{"sub_4E0370", "ArmorMultiplierEffect", "defend", "float"},
		{"sub_4E0380", "DurabilityMultiplierEffect", "defend", "float"},
		{"nullsub_40", "ResilienceEffect", "defend", "float"},
		{"nox_xxx_inversionEffect_4E03D0", "InversionEffect", "defend", "int"},
		{"nox_xxx_gripEffect_4E0480", "GripEffect", "defend", "int"},
		{"nullsub_41", "BreakingEffect", "defend", "float"},
		{"nullsub_42", "PunctureProneEffect", "defend", "float"},
		{"nox_xxx_effectRegeneration_4E01D0", "RegenerationUpdate", "update", "int"},
		{"nullsub_43", "ParasiteUpdate", "update", "int"},
		{"nullsub_44", "AttractionUpdate", "update", "int"},
		{"nox_xxx_attribContinualReplen_4E02C0", "ContinualReplenishmentUpdate", "update", "int"},
		{"sub_4DFB50", "BrillianceEngage", "engage", "int"},
		{"sub_4DFB80", "BrillianceDisengage", "engage", "int"},
		{"nox_xxx_effectSpeedEngage_4DFC30", "SpeedEngage", "engage", "float"},
		{"nox_xxx_effectSpeedDisengage_4DFCA0", "SpeedDisengage", "engage", "float"},
		{"sub_4DFD10", "FireProtectEngage", "engage", "float"},
		{"nox_xxx_modifFireProtection_4DFD40", "FireProtectDisengage", "engage", "float"},
		{"nox_xxx_buff_4DFD80", "LightningProtectEngage", "engage", "float"},
		{"sub_4DFDB0", "LightningProtectDisengage", "engage", "float"},
		{"nox_xxx_checkPoisonProtectEnch_4DFDE0", "PoisonProtectEngage", "engage", "float"},
		{"sub_4DFE10", "PoisonProtectDisengage", "engage", "float"},
		{"sub_4E0140", "RegenerationEngage", "engage", "none"},
		{"sub_4E0170", "RegenerationDisengage", "engage", "none"},
	}
	entries := legacy.PortTestModifierRegistryEntries()
	if len(entries) != len(want) {
		t.Fatalf("registry entries: got %d, want %d", len(entries), len(want))
	}
	seen := make(map[unsafe.Pointer]string, len(entries))
	mods, freeMods := alloc.Make([]server.ModifierEff{}, len(entries))
	defer freeMods()
	for i, e := range entries {
		w := want[i]
		if e.Symbol != w.symbol || e.Name != w.name || e.Kind != w.kind || e.Parser != w.parser {
			t.Fatalf("entry %d metadata: got %+v want %+v", i, e, w)
		}
		if e.Key == nil {
			t.Fatalf("%s has nil C address", e.Symbol)
		}
		if prev, ok := seen[e.Key]; ok {
			t.Fatalf("%s shares original C key with %s", e.Symbol, prev)
		}
		seen[e.Key] = e.Symbol
		got, parser, ok := server.PortTestModifierRegistryLookup(e.Kind, e.Name)
		if !ok || got != e.Key || parser != e.Parser {
			t.Fatalf("%s registry lookup=(%p,%q,%v), table=(%p,%q)", e.Name, got, parser, ok, e.Key, e.Parser)
		}
	}
	for i, e := range entries {
		mod := &mods[i]
		mod.Price20 = -0x1234567
		mod.AllowWeapons28 = 0xfedcba98
		mod.Attack40.Fnc, mod.AttackPreHit52.Fnc, mod.AttackPreDmg64.Fnc = e.Key, e.Key, e.Key
		mod.Defend76.Fnc, mod.DefendCollide88.Fnc, mod.Update100.Fnc = e.Key, e.Key, e.Key
		mod.Engage112, mod.Disengage116 = e.Key, e.Key
	}
	before := append([]server.ModifierEff(nil), mods...)
	runtime.GC()
	fresh := legacy.PortTestModifierRegistryEntries()
	if len(fresh) != len(entries) {
		t.Fatal("modifier table length changed")
	}
	for i, e := range fresh {
		if mods[i] != before[i] {
			t.Fatal("foreign modifier storage changed", i)
		}
		if e != entries[i] {
			t.Fatalf("C key changed after GC/storage churn at %d: %p -> %p", i, entries[i].Key, e.Key)
		}
		got, _, ok := server.PortTestModifierRegistryLookup(e.Kind, e.Name)
		if !ok || got != entries[i].Key {
			t.Fatalf("registry changed after GC/storage churn for %s", e.Name)
		}
	}
}

func TestModifierParserContracts(t *testing.T) {
	entries := legacy.PortTestModifierRegistryEntries()
	fields := []struct {
		offset int
		kind   string
	}{{40, "damage"}, {52, "damage"}, {64, "damage"}, {76, "defend"}, {88, "defend"}, {100, "update"}, {112, "engage"}, {116, "engage"}}
	intCases := []struct {
		arg  string
		want int32
	}{{"", -1}, {"19 rest", 19}, {"-27 tail", -27}, {"19tail", -1}, {"2147483647", 0x7fffffff}, {"-2147483648", -0x80000000}, {"2147483648", -1}, {"bad", -1}}
	floatCases := []struct {
		arg  string
		want float32
	}{{"", -1}, {"1.25 rest", 1.25}, {"-2.5 trailing", -2.5}, {"1.25tail", -1}, {"3.4028234663852886e38", float32(3.4028234663852886e38)}, {"1e100", -1}, {"bad", -1}}
	mods, freeMods := alloc.Make([]server.ModifierEff{}, 8)
	defer freeMods()
	for _, field := range fields {
		for _, entry := range entries {
			if field.kind != entry.Kind {
				continue
			}
			cases := []struct {
				arg string
				int int32
				flt float32
			}{
				{arg: "", int: -1, flt: -1},
			}
			if entry.Parser == "int" {
				for _, c := range intCases {
					cases = append(cases, struct {
						arg string
						int int32
						flt float32
					}{c.arg, c.want, 0})
				}
			} else if entry.Parser == "float" {
				for _, c := range floatCases {
					cases = append(cases, struct {
						arg string
						int int32
						flt float32
					}{c.arg, 0, c.want})
				}
			} else {
				cases = append(cases, struct {
					arg string
					int int32
					flt float32
				}{"ignored", 0, 0})
			}
			for caseIndex, c := range cases {
				mod := &mods[fieldIndex(field.offset)]
				seedModifierCanaries(mod, entries)
				before := *mod
				raw := entry.Name
				if caseIndex != 0 {
					raw += " " + c.arg
				}
				err := server.PortTestModifierParseEffect(field.offset, mod, raw)
				if err != nil {
					t.Fatalf("slot%d %s/%q: %v", field.offset, entry.Name, c.arg, err)
				}
				want := before
				setModifierExpected(&want, field.offset, entry.Key, entry.Parser, c.int, c.flt, true)
				if *mod != want {
					t.Fatalf("slot%d %s/%q got %+v want %+v", field.offset, entry.Name, c.arg, *mod, want)
				}
			}

			// Unknown names and empty values must not change either the callback
			// slot or its associated scalar storage.
			for _, raw := range []string{"", "UnknownModifierEffect 42"} {
				mod := &mods[fieldIndex(field.offset)]
				seedModifierCanaries(mod, entries)
				before := *mod
				err := server.PortTestModifierParseEffect(field.offset, mod, raw)
				if raw == "" && err != nil {
					t.Fatalf("slot%d empty parser value: %v", field.offset, err)
				}
				if raw != "" && err == nil {
					t.Fatalf("slot%d unknown parser name accepted", field.offset)
				}
				if *mod != before {
					t.Fatalf("slot%d %q changed target storage", field.offset, raw)
				}
			}
		}
	}
}

func fieldIndex(offset int) int {
	switch offset {
	case 40:
		return 0
	case 52:
		return 1
	case 64:
		return 2
	case 76:
		return 3
	case 88:
		return 4
	case 100:
		return 5
	case 112:
		return 6
	case 116:
		return 7
	default:
		panic("unknown modifier slot")
	}
}

func seedModifierCanaries(mod *server.ModifierEff, entries []legacy.PortTestModifierRegistryEntry) {
	*mod = server.ModifierEff{Price20: -0x1234567, AllowWeapons28: 0xfedcba98, AllowArmor32: 0x87654321, AllowPos36: 0x76543210}
	for i, e := range entries {
		ptr := e.Key
		value := int32(0x1020304 + i)
		flt := float32(i) + 0.375
		switch e.Kind {
		case "damage":
			slots := []*server.ModifierEffFnc{&mod.Attack40, &mod.AttackPreHit52, &mod.AttackPreDmg64}
			for _, s := range slots {
				s.Fnc, s.Val, s.Valf = ptr, value, flt
			}
		case "defend":
			for _, s := range []*server.ModifierEffFnc{&mod.Defend76, &mod.DefendCollide88} {
				s.Fnc, s.Val, s.Valf = ptr, value, flt
			}
		case "update":
			mod.Update100.Fnc, mod.Update100.Val, mod.Update100.Valf = ptr, value, flt
		case "engage":
			mod.Engage112, mod.EngageInt124, mod.EngageFloat120 = ptr, value, flt
			mod.Disengage116, mod.DisengageInt132, mod.DisengageFloat128 = ptr, value, flt
		}
	}
}

func setModifierExpected(mod *server.ModifierEff, offset int, key unsafe.Pointer, parser string, iv int32, fv float32, parsed bool) {
	var p *unsafe.Pointer
	var i *int32
	var f *float32
	switch offset {
	case 40:
		p, i, f = &mod.Attack40.Fnc, &mod.Attack40.Val, &mod.Attack40.Valf
	case 52:
		p, i, f = &mod.AttackPreHit52.Fnc, &mod.AttackPreHit52.Val, &mod.AttackPreHit52.Valf
	case 64:
		p, i, f = &mod.AttackPreDmg64.Fnc, &mod.AttackPreDmg64.Val, &mod.AttackPreDmg64.Valf
	case 76:
		p, i, f = &mod.Defend76.Fnc, &mod.Defend76.Val, &mod.Defend76.Valf
	case 88:
		p, i, f = &mod.DefendCollide88.Fnc, &mod.DefendCollide88.Val, &mod.DefendCollide88.Valf
	case 100:
		p, i, f = &mod.Update100.Fnc, &mod.Update100.Val, &mod.Update100.Valf
	case 112:
		p, i, f = &mod.Engage112, &mod.EngageInt124, &mod.EngageFloat120
	case 116:
		p, i, f = &mod.Disengage116, &mod.DisengageInt132, &mod.DisengageFloat128
	}
	*p = key
	if parsed && parser == "int" {
		*i = iv
	}
	if parsed && parser == "float" {
		*f = fv
	}
}
