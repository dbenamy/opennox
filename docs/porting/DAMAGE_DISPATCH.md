# Damage dispatch, defense effects and durability

## Original-C baseline — 2026-09-11

All 26 candidate functions remain original C. The guarded fixture covers
**3,582 cases / 42 complete capture groups**, repeated byte-for-byte in separate
processes (12.074s confirmation). Hashes are locked in
`src/damage_dispatch_porttest_test.go`; full local evidence is in
`build/port-damage-dispatch/c-final-*` and `c-confirm-*`.

Coverage includes all eighteen damage kinds; player/simple/monster subjects;
health/lethal boundaries; game modes, damage sources and owner relationships;
dead/invulnerable/shock/reflection/shield behavior; elemental resistances;
armor absorption and conductive armor; fractional accumulation; defense and
pre-damage modifier ordering and changed inputs; weapon/armor durability calls;
blocking; generator health thresholds and clocks; gameball drop thresholds;
projectile reflection geometry/directions; and supported nil-input paths.
Explicit assertions check actual 50→38 health changes, parser results, and
ball ownership on either side of the 30-damage threshold.

Fixture details discovered before locking:

- Durability's slot-one defense callback receives one float; defense slots two
  and three receive an integer damage/type pair. Separate recorders avoid
  reading uninitialized stack words.
- Standalone tests do not initialize the runtime damage-name blob. The fixture
  saves/restores the region, loads the shipped bytes, and relocates the eighteen
  pointers just as runtime initialization does. Exact parser assertions caught
  the earlier empty-string setup. The same region supplies the original 0.15
  conductivity constant.
- `itemDestroyed` takes a player index, not an object pointer. Its dispatcher
  uses an explicit scalar index. `damageArmor` requires a nonnil object; nil
  admission tests cover only original-C-supported paths.
- Recorders only observe dependency calls and optionally return configured
  values. No production algorithm has been copied into test C.

Production C remains **122,284 lines / 149 files / zero reference C**.
All 31,623 accumulated focused cases / 262 capture groups pass with locked
hashes (95.765s), including all 28,041 prior cases unchanged.
Commit/push this baseline before converting the connected batch.

Candidate: 26 connected functions / 1,331 C lines spanning GAME3_2.c and
GAME3_3.c, addresses 004E0A00 through 004E27D0. Group default/player damage,
armor/weapon durability, defense modifiers, special damage wrappers, HP/mana
sharing, projectile reflection and damage-type parsing. Audit callee boundaries
before locking the scope. Reuse native resource/equipment/inventory/objective
owners and the guarded attack/collision fixture.

Create the original-C baseline first; repeat complete captures, lock hashes,
commit and push before converting production. Preserve all earlier hashes.
Use actual player/NPC slabs, health/protection records, modifier definitions,
team/owner chains and callback identities. Add a damage fixture spec and dispatch
range 1100+, guarded in/out damage values, defense-effect recorders and explicit
float-return capture. Type parsing should use a saved/restored pointer table of
real C strings if the original relocated name table needs fixture setup.

Exercise nil/admission returns, dead/zombie/invulnerable states, friendly fire
and game modes, elemental resistances, armor/shield absorption, defense/pre-hit
effects, HP/mana sharing, durability/destruction, damage-source selection,
ball-carrier behavior, NPC scripting and specialized damage wrappers. Include
integer/floating boundaries, clock wrap, low-byte returns and seeded randomness.
Assert actual HP/mana/durability deltas, blocked damage, modified callback values,
destruction requests and successful downstream behavior. No copied C algorithms.

After the connected batch matches the C baseline, qualify once with accumulated
port tests/default/server/highres, production builds and headless gameplay.
Update C_LOC and recovery docs, commit/push, summarize, continue. No new agents.

## Candidate scope

- 004E0A00: `int nox_xxx_parseDamageTypeByName_4E0A00(const char* a1)` (17 lines).
- 004E0A70: `int nox_xxx_projectileReflect_4E0A70(int a1, int a2)` (38 lines).
- 004E0B30: `int nox_xxx_damageDefaultProc_4E0B30(int a1, int a2, int a3, int a4, int a5)` (317 lines).
- 004E1230: `void nox_xxx_gameballOnPlayerDamage_4E1230(int a1, int a2, int a3)` (37 lines).
- 004E1320: `int nox_xxx_itemApplyDefendEffect2_4E1320(int a1, int a2, int a3, int* a4, int a5)` (39 lines).
- 004E13B0: `int nox_xxx_itemApplyPreDamageEffect_4E13B0(int a1, int a2, int a3, int a4)` (26 lines).
- 004E1400: `int sub_4E1400(int a1, uint32_t* a2)` (29 lines).
- 004E1470: `int sub_4E1470(int a1)` (17 lines).
- 004E14A0: `int sub_4E14A0() { return 0; }` (3 lines).
- 004E14B0: `int sub_4E14B0(int a1, int a2, int a3, int a4, int a5)` (12 lines).
- 004E1500: `int nox_xxx_damageArmor_4E1500(int a1, int a2, int a3, int a4, int a5)` (14 lines).
- 004E1560: `void nox_xxx_playerDamageWeapon_4E1560(int a1, int a2, int a3, int a4, float a5, int a6)` (42 lines).
- 004E1650: `int nox_xxx_itemDestroyed_4E1650(int a1, uint32_t* a2, unsigned short a3, unsigned short a4)` (20 lines).
- 004E16D0: `void nox_xxx_equipDamage_4E16D0(int a1, int a2, int a3, int a4, float a5, int a6)` (34 lines).
- 004E17B0: `int nox_server_handler_PlayerDamage_4E17B0(int a1, int a2, int a3, int a4, int a5)` (401 lines).
- 004E20F0: `void nox_xxx_playerDecrementHPMana_4E20F0(int a1, int a2, float a3)` (34 lines).
- 004E2180: `void nox_xxx_playerDamageItems_4E2180(int a1, int a2, int a3, int a4, float a5)` (36 lines).
- 004E2220: `double sub_4E2220(int a1)` (33 lines).
- 004E22A0: `int sub_4E22A0(int a1, int a2, int a3, int a4, float a5, int a6)` (26 lines).
- 004E2330: `int sub_4E2330(int a1, int a2, int a3, int a4, float a5, int a6)` (29 lines).
- 004E23C0: `int sub_4E23C0(int a1, int a2, int a3, int a4, int a5)` (26 lines).
- 004E24B0: `int sub_4E24B0(int a1, int a2, int a3, int a4, int a5) { return nox_xxx_damageDefaultProc_4E0B30(a1, a2, a3, a4, a5); }` (3 lines).
- 004E24E0: `int sub_4E24E0(int a1, int a2, int a3, int a4, int a5)` (12 lines).
- 004E2520: `int nox_xxx_damageFlammable_4E2520(int a1, int a2, int a3, int a4, int a5)` (10 lines).
- 004E2560: `int nox_xxx_damageBlackPowder_4E2560(int a1, int a2, int a3, int a4, int a5)` (12 lines).
- 004E27D0: `int nox_xxx_damageMonsterGen_4E27D0(int a1, int a2, int a3, int a4, int a5)` (64 lines).

Status: scope/fixture audit only. All candidate functions remain C; no
original-C hashes for this batch are locked yet.
