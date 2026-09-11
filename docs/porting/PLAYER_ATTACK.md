# Next batch: player attack, hit dispatch, bow/crossbow and quiver reload

Candidate: 12 connected functions / 1,105 C lines in GAME4_3.c, addresses
00538290 through 00539FF0, excluding already converted equipment. Keep player
attack dispatch and its melee/ranged callbacks together. Audit callee boundaries
before locking scope; reuse native equipment/effects/inventory owners. Existing
actual projectile definitions, frame clock, player/NPC slabs, damage/collision
recorders and object allocator should supply most of the fixture.

Baseline first: extend the objective/world/equipment fixture only where needed,
with guarded attack records, actual animation data and projectile-class table
entries. Exercise all weapon-mask branches and player/NPC admission, frame and
half-frame/duplicate-hit gates, ability states, ammo/charge boundaries, bow and
crossbow creation/aim/spread, melee hit/miss/wall filtering, damage and effect
callbacks, weapon durability and auto-reload. Assert successful damage, created
shots, ownership/trajectory, consumed ammo and reload/equip changes. Include
seeded random streams, unsigned clock wrap and float rounding boundaries.

Original C owns every candidate during repeat/lock/commit/push. Do not copy C
algorithms into tests. Compare complete guarded state and event captures, then
convert the whole connected batch and qualify once. Update C_LOC/docs, commit,
push, summarize, continue. No new agents and no pending user question.

Fixture audit notes:
- Existing player animation table already supplies deterministic [frames,delay]
  for 64 animation kinds; expose overrides only for explicit boundary probes.
- Projectile-class definitions are actual 88-byte server.Modifier records in
  core.Modif.Dword_5d4594_251600. Equipment fixture already allocates/guards them;
  add optional Range68/DamageMin72/DamageType76 or word overrides.
- Attack record is 36 bytes: float damage0, byte type4, float radius8, owner12,
  origin16/20, hit-immovable24, weapon28, direction-mask32. Use a guarded 64-byte
  allocation and normalized object references; no algorithm in the dispatcher.
- Add ArcherArrow/ArcherBolt/WeakArcherArrow real type definitions. Their
  collision payload needs 8 bytes (shooter pointer at +4), plus 16 guard bytes.
  Existing effects snapshot assumes a 4-byte Spark payload; derive payload
  length from the actual type's CollideDataSize to preserve old captures.
- Objectives fixture zeroes InitData+4 for optional flag color names. Attack
  setup must restore that modifier slot from objectives.initWords before
  assigning effect callbacks, so all four modifier slots remain testable.
- Save/restore hit/nearest globals 2488652/2488656/2488660. Use actual spatial
  indexes and wall queries, with positive coordinates for area searches.

## Audited scope

- 00538290: `int nox_xxx_playerPreAttackEffects_538290(int a1, int a2, int a3, int a4)` (40 lines).
- 00538330: `int nox_xxx_playerTraceAttack_538330(int a1, int a2)` (70 lines).
- 00538510: `void sub_538510(int a1, int a2)` (76 lines).
- 005386A0: `void sub_5386A0(int a3, int a2)` (58 lines).
- 00538840: `int nox_xxx_itemApplyAttackEffect_538840(int a1, int a2, int a3)` (26 lines).
- 00538960: `int nox_xxx_playerAttack_538960(nox_object_t* a1p)` (590 lines).
- 00539B90: `short nox_xxx_warcryStunMonsters_539B90(int a1, int a2)` (16 lines).
- 00539BD0: `int nox_xxx_shootBowCrossbow1_539BD0(int a1, int a2)` (77 lines).
- 00539D80: `uint32_t* nox_xxx_shootBowCrossbow2_539D80(int a1, int a2, int a3, char* a4)` (84 lines).
- 00539F40: `int nox_xxx_shootApplyEffects_539F40(int a1, int a2, int a3)` (34 lines).
- 00539FB0: `int sub_539FB0(uint32_t* a1)` (17 lines).
- 00539FF0: `int nox_xxx_playerTryReloadQuiver_539FF0(uint32_t* a1)` (17 lines).

Production is still C; baseline fixture work begins after conversion `ca0f193c`.
No original-C hashes are locked yet. Current C count: 124,395 / 149 files / zero
reference C.
