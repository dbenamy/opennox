# Projectile collisions

27 connected handlers/helpers, 1,006 physical C lines in GAME3_3.c. Keep arrow
and chakram collisions with their helpers; include generic/spark/fireball,
sulphur/death-ball/pixie, webbing, fist/teleport wake and damage/trap callbacks.
Candidate source/scope are local artifacts; review callee boundaries before
locking scope. No production change until an original-C baseline is repeated,
locked, committed and pushed.

Reuse objective/world/attack fixture for real players, guarded object slabs,
modifier definitions, callback damage returns, inventory transfer and projectile
creation. Add a nested collision spec and operation range 1000+, plus optional
normal vector. Existing world collision-data words and temporary update-data
words should cover most inputs. Save/restore touched type caches and wall-contact
scratch globals; guard/norm every new region. Preserve prior capture hashes.

Cover target admission (nil/player/NPC/object, dead/immovable/friendly/owner),
damage accepted/rejected, frame gates, energy/charges, invulnerability and
reflection, wall vs unit collision, normals and velocity, projectile deletion,
owner reassignment, ammo/inventory transfer, trap buffs and mana drain. Positive
assertions for actual damage, bounce vectors, pickup/return, buff timing and
resource deltas. Callback recorder must allow controlled return values without
retaining original algorithms. Keep exact event and full guarded-state hashes.

Qualify once after the connected batch is complete. No extra agents. Update C
LOC/doc checkpoint, commit/push, summarize and continue.

Fixture audit:
- ai_callbacks_porttest.go's pt_callback_hit already records all five damage
  arguments and returns 1. Add optional callback DamageResult with reset default
  1, preserving all existing captures; set target/player handlers to this actual
  recorder only after setup. Do not build another algorithm-shaped damage stub.
- Save memmap words 1567836/40, 1567924/32, 1567948/52, 1567964/68/72/76/80/84,
  1568000; inspect whether some contain callback object pointers and normalize
  them. Most are type caches, but retain exact pre/post behavior.
- Add actual ThrowingStone/ImpShot/ArcherArrow/ArcherBolt and trap spawn type
  definitions needed for named branches. Reuse guarded init/collision/update
  templates, including actual payload sizes when captures read created objects.
- Chakram impact builds a 36-byte attack record (v22[9]); the attack callback
  recorder can capture it. It does not initialize Front byte32: only attach
  modifier recorders with an explicitly justified mask for fields unused by
  the impact owner, rather than locking undefined stack bytes.
- Original sub_518040 shape-aware query remains a dependency of native attack;
  do not substitute center-only EachObjInCircle in collision target selection.
- sub_4EB250 uses rectangle visitor state and helpers sub_4EB340/sub_4EB3E0.
  Exercise target selection and seeded direction fallback as a connected path.
- Keep stable wall-contact scratch pointers until collision snapshots finish;
  restore all globals before releasing fixture-owned buffers.

## Candidate scope

- 004E87B0: `void nox_xxx_collideProjectileGeneric_4E87B0(int a1, int a2)` (46 lines).
- 004E8880: `void nox_xxx_collideProjectileSpark_4E8880(int a1, int a2)` (29 lines).
- 004E9430: `void nox_xxx_collideDamage_4E9430(int a1, int a2)` (22 lines).
- 004E9490: `void nox_xxx_collideManadrain_4E9490(int a1, int a2)` (15 lines).
- 004E96F0: `void nox_xxx_collideBomb_4E96F0(int a1, int a2)` (19 lines).
- 004E9770: `void nox_xxx_collideBoom_4E9770(int a1, int a2, float* a3)` (72 lines).
- 004E99B0: `void nox_xxx_collideDie_4E99B0(int unit, int a2)` (18 lines).
- 004E9A30: `int sub_4E9A30(nox_object_t* a1p, nox_object_t* a2p)` (28 lines).
- 004E9AC0: `void nox_xxx_fireballCollide_4E9AC0(int a1, int a2)` (60 lines).
- 004E9D80: `void nox_xxx_collideSulphurShot2_4E9D80(int a1, int a2, float* a3)` (36 lines).
- 004E9E50: `void nox_xxx_collideSulphurShot_4E9E50(int a1, int a2, int a3)` (8 lines).
- 004E9FE0: `void nox_xxx_collideDeathBallFragment_4E9FE0(int a1, int a2, float* a3)` (24 lines).
- 004EA080: `void nox_xxx_collidePixie_4EA080(int a1, int a2, float* a3)` (73 lines).
- 004EA200: `void nox_xxx_collideWallReflectSpark_4EA200(int a1, int a2, float2* a3)` (37 lines).
- 004EA2C0: `void sub_4EA2C0(int a1, int a2)` (15 lines).
- 004EA300: `void nox_xxx_collideSpark_4EA300(int a1, int a2, float* a3)` (22 lines).
- 004EA380: `void nox_xxx_collideWebbing_4EA380(int a1, int a2)` (19 lines).
- 004EADF0: `void nox_xxx_collideFist_4EADF0(int a1, int a2)` (14 lines).
- 004EAE30: `void nox_xxx_collideTeleportWake_4EAE30(int a1, int a2)` (26 lines).
- 004EAF00: `void nox_xxx_collideChakram_4EAF00(int a1, int a2, float* a3)` (148 lines).
- 004EB250: `int sub_4EB250(int a1)` (37 lines).
- 004EB340: `void sub_4EB340(float* a1, int a2)` (22 lines).
- 004EB3E0: `void sub_4EB3E0(int a1)` (36 lines).
- 004EB490: `void nox_xxx_collideArrow_4EB490(int a1, int a2)` (115 lines).
- 004EB800: `void nox_xxx_collideMonsterArrow_4EB800(int a1, int a2)` (30 lines).
- 004EB890: `void nox_xxx_collideBearTrap_4EB890(int* a1, int a2)` (16 lines).
- 004EB910: `void nox_xxx_collidePoisonGasTrap_4EB910(int* a1, int a2)` (19 lines).

## Original-C baseline

Original-C baseline `42fdbc86` was committed and pushed before conversion. A 1,704-case / 49-group corpus
passes against them; first and repeat full captures are byte-identical (5.401s
and 5.841s). Hashes are locked in projectile_collisions_porttest_test.go. The accumulated
regression passes all 28,041 focused cases / 220 groups in 81.843s, including
every earlier locked hash. Every earlier locked hash was unchanged.

Coverage includes accepted/rejected damage (including low-byte versus full-word
return tests), named projectile damage, wall normals and contact, game modes,
buffs/reflection, arrow/bolt definitions and depleted target health, explosion
splash, mana/damage frame boundaries, script-triggered bombs, death overrides,
chakram bounce/owner/return/drop/selection, and trap creation. Positive assertions
verify hit/deletion behavior, bounced velocity, webbing, bear/gas trap creation,
mana subtraction, and restored/re-equipped chakram inventory. Inspection also
confirms successful target selection at the inclusive 400-unit boundary and
multi-target splash damage.

The fixture uses actual guarded object/type/modifier records, a controlled
return on the existing damage recorder, and guarded normal vectors. It saves
and restores wall-contact state, type caches and selection globals. Collision
attack records leave the Front word unused and uninitialized; only that word,
previously documented padding, and relocated addresses are normalized. Existing
attack owners retain their original Front-byte capture mask.

Local evidence: build/port-projectile-collisions/{candidate-scope.json,
c-first.log,c-repeat.log,c-first-projectile-*.json,c-repeat-projectile-*.json,
baseline-hashes.json,c-regression.log}. No C algorithm is retained for tests.
Current production C remains 123,290 lines / 149 files / zero reference C.

## Native conversion

All 27 functions / 1,006 C lines now live in projectile_collisions.go,
projectile_collisions_effects.go, projectile_collisions_weapons.go and thin
projectile_collisions_exports.go. All 49 full native captures match original C
on the first run (1,704 cases, 6.175s). Production C is 122,284 lines /
149 files / zero test-reference C.

Handlers call native attack/equipment/inventory/resource/reflection helpers
directly. Arrow and chakram damage records preserve the original callback and
mutation order; rejected-hit low-byte/full-word distinctions remain explicit.
Chakram candidate filtering retains the original numeric conversion of float
interpretations of class/flag words, including its unusual behavior. Spatial
selection, fallback RNG, signed-short mana timestamps, splash falloff, trap
creation and inventory-return effects match the captured C behavior.

## Qualification

All accumulated port tests, including 28,041 focused cases / 220 groups, pass in
default/server/highres: 151.165s / 137.223s / 142.355s. All three production binaries are
verified ELF32/i386, GO386=sse2 and CGO enabled. Full-suite failure identities and
multiplicities match the baseline exactly: 1,553 entries; 15 packages pass,
3 fail and 32 skip. Fresh unchanged repeat-a headless gameplay passes in
36.421s using Xvfb and null audio, without updating expected captures.

Local evidence: build/port-projectile-collisions (native-first captures,
variant logs, binaries and qualification.json), and
build/baseline/runs/projectile-collisions-port. Next:
[damage dispatch and durability](DAMAGE_DISPATCH.md), one connected batch.
