# Player attack, hit dispatch, bow/crossbow and quiver reload

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

## Original-C baseline

Baseline `b94b52c8` was committed and pushed before conversion. The guarded fixture covers 2,100 cases
in 16 groups: unarmed timing, hit filters, nearest target, modifier effects,
projectiles, weapon dispatch, reloads, spatial traces, warcry, positive outcomes,
NPC attacks, abilities, projectile modifiers/failures, chakram depletion and
NPC/player ammo consumption. Two independent runs produced byte-identical full
captures. SHA-256 expectations are committed in player_attack_porttest_test.go.

Independent assertions verify damage amounts, projectile position/velocity,
warcry duration, round-chakram inventory transfer, modifier changes to damage,
ammo consumption even when object creation fails, and fan-chakram deletion
requests plus replacement equipment. The delete hook observes requests; it does
not execute deletion. Fixture teardown detaches already-captured transferred
items before freeing created objects, avoiding duplicate frees.

Only unused attack-record padding and relocated object/function addresses are
normalized. The warcry short return can contain the low 16 bits of an input
pointer; that known pointer relationship is normalized before comparison.
No C gameplay algorithm is copied into the fixture.

Local evidence: build/port-player-attack/{complete-first.log,complete-repeat.log,
c-first-attack-*.json,c-repeat-attack-*.json,baseline-hashes.json,c-regression.log}.
Current production C count remains 124,395 / 149 files / zero reference C.

Accumulated focused regression: all 26,337 cases / 171 groups pass against
original C in 76.006s, with every earlier locked hash unchanged.

## Native conversion

All 12 functions / 1,105 C lines now live in player_attack.go,
player_attack_melee.go, player_attack_ranged.go and thin player_attack_exports.go.
All 16 full native captures match original C byte-for-byte (2,100 cases, 6.631s).
Production C is 123,290 lines / 149 files / zero test-reference C.

The dispatcher shares repeated melee setup but preserves weapon-mask priority,
byte truncation of animation counters, uint32 frame wrap, half-frame gates,
readiness calls, ability state transitions and equipment/inventory ordering.
It calls native equipment, inventory and item-use owners directly. Bow/crossbow
shots preserve clear-ray return values and consume player ammunition even when
projectile allocation fails; fan and round chakrams retain distinct depletion
and inventory-transfer behavior.

The first native run matched 15/16 groups. The trace corpus caught an incorrect
substitution of center-distance EachObjInCircle for the original shape-aware
sub_518040 query. The original query remains a production dependency and invokes
the native hit/nearest callbacks through their ABI exports. This preserves object
extent checks, strict range boundaries and spatial-query visitation state.
The unarmed modifier record is allocated in C-owned memory while the remaining
C damage-calculation dependency reads it, including any callbacks into Go.

## Qualification

All accumulated port tests, including 26,337 focused cases / 171 groups, pass in
default/server/highres: 146.770s / 127.486s / 134.713s. All three production binaries are
verified ELF32/i386, GO386=sse2 and CGO enabled. Full-suite failure identities and
multiplicities match the baseline exactly: 1,553 entries; 15 packages pass,
3 fail and 32 skip. Fresh unchanged repeat-a headless gameplay passes in
39.677s using Xvfb and null audio, without updating expected captures.

Local evidence: build/port-player-attack (native-final captures, variant logs,
binaries and qualification.json), and build/baseline/runs/player-attack-port.
Next: [projectile collisions](PROJECTILE_COLLISIONS.md), one connected batch.
