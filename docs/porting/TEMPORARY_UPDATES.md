# Temporary objects and projectile updates

31 connected functions / 980 C lines (scope.json addresses audited): Spark creation,
Spark/trail/lifetime, homing spell/anti-spell/magic missile, temporary barrels,
moonlight/telekinesis/fist/flame/meteor, toxic clouds, spider spawning, expire,
break/open/remove and returning chakram. Keep map-area damage, target search and
spell implementations as real production dependencies. Reuse EffectsUse fixture;
new optional TemporaryUpdates spec/actions 600+ must not alter older hashes.

Before conversion: repeat full original-C captures, lock contracts and push baseline.

Fixture additions needed:
- guarded item update storage + pointer relocation for owner/target/projectile links;
  distinct sources and targets, player cursor/health/buffs, full return and mutation
  snapshots, exact created/deleted/updatable object lists and existing damage observer.
- actual original shared cache words for anti-spell nearest selection, water flag,
  flame/spider type IDs; real type definitions for Spark/flames/Meteor/MeteorExplode/
  SmallSpider and balance lifetimes. Save/restore every touched global.
- filtered spatial-index fixture objects for rectangle/circle/missile queries;
  direct callback tests plus owner integration with nearest/tie/order boundaries.
- retain actual map ray and collision paths; use safe positive coordinates. Do not
  suppress negative-coordinate map-ray issue already documented in INVENTORY.md.
- lifecycle owns creations and checks exact zero leftovers. Guard size must cover
  each actual template; no blanket disabling of existing Spark guard checks.

Contracts: signed/unsigned frame wrap, exact vs strict deadlines, branch priority,
life count zero/one/negative, death callback vs delayed delete, legal state-bit
combinations for break/open, owner/target disappearance, direction wrap and dot
boundary, nearest-candidate filters/ties, collision and expiry returns, target
search timing, all Spark stages, eight-trail creation RNG order, seeded random
spawns and missing name lookups, blocked/clear rays, actual poison/area damage,
force/velocity rounding and radius boundaries. Tests should independently assert
positive effects and exact creation counts as well as captured C hashes.

Once native: route native Go owners/registrations while preserving callback address
identity, compare full captures, all prior contracts, default/server/highres,
production builds and fresh gameplay once for the connected family. Full-suite
comparison at this projectile/shared-fixture milestone. Record C_LOC and push.

## Audited candidate scope

- src/legacy/GAME4_3.c: 0053ADC0 — void nox_xxx_updateSpark_53ADC0.
- src/legacy/GAME4_3.c: 0053AEC0 — float* nox_xxx_updateProjTrail_53AEC0.
- src/legacy/GAME4_3.c: 0053B8F0 — void nox_xxx_updateLifetime_53B8F0.
- src/legacy/GAME4_3.c: 0053B940 — void nox_xxx_spellFlyUpdate_53B940.
- src/legacy/GAME4_3.c: 0053BB00 — void nox_xxx_updateAntiSpellProj_53BB00.
- src/legacy/GAME4_3.c: 0053BD10 — void sub_53BD10.
- src/legacy/GAME4_3.c: 0053BDA0 — int nox_xxx_updateMagicMissile_53BDA0.
- src/legacy/GAME4_3.c: 0053C9A0 — void nox_xxx_updateBlackPowderBarrel_53C9A0.
- src/legacy/GAME4_3.c: 0053CB60 — void nox_xxx_updateOneSecondDie_53CB60.
- src/legacy/GAME4_3.c: 0053CB90 — void nox_xxx_updateWaterBarrel_53CB90.
- src/legacy/GAME4_3.c: 0053CC30 — void nox_xxx_waterBarrel_53CC30.
- src/legacy/GAME4_3.c: 0053CC90 — void nox_xxx_updateSelfDestruct_53CC90.
- src/legacy/GAME4_3.c: 0053CCB0 — void nox_xxx_updateBlackPowderBurn_53CCB0.
- src/legacy/GAME4_3.c: 0053D220 — void nox_xxx_updateDeathBallFragment_53D220.
- src/legacy/GAME4_3.c: 0053D270 — void nox_xxx_updateMoonglow_53D270.
- src/legacy/GAME4_3.c: 0053D330 — void nox_xxx_updateTelekinesis_53D330.
- src/legacy/GAME4_3.c: 0053D400 — void nox_xxx_updateFist_53D400.
- src/legacy/GAME4_3.c: 0053D510 — void nox_xxx_updateFlameCleanse_53D510.
- src/legacy/GAME4_3.c: 0053D5A0 — void nox_xxx_updateMeteorShower_53D5A0.
- src/legacy/GAME4_3.c: 0053D6E0 — void nox_xxx_meteorExplode_53D6E0.
- src/legacy/GAME4_3.c: 0053D850 — void nox_xxx_updateToxicCloud_53D850.
- src/legacy/GAME4_3.c: 0053D8C0 — void sub_53D8C0.
- src/legacy/GAME4_3.c: 0053D960 — void nox_xxx_updateSmallToxicCloud_53D960.
- src/legacy/GAME4_3.c: 0053D9D0 — void nox_xxx_toxicCloudPoison_53D9D0.
- src/legacy/GAME4_3.c: 0053DA60 — void nox_xxx_updateArachnaphobia_53DA60.
- src/legacy/GAME4_3.c: 0053DB00 — void nox_xxx_updateExpire_53DB00.
- src/legacy/GAME4_3.c: 0053DB30 — int* nox_xxx_updateBreak_53DB30.
- src/legacy/GAME4_3.c: 0053DBB0 — int* nox_xxx_updateOpen_53DBB0.
- src/legacy/GAME4_3.c: 0053DC30 — void nox_xxx_updateBreakAndRemove_53DC30.
- src/legacy/GAME4_3.c: 0053DCC0 — void nox_xxx_updateChakramInMotion_53DCC0.
- src/legacy/GAME5.c: 0054FD80 — float* nox_xxx_createSpark_54FD80.

## Original-C baseline

The effects conversion and timing correction are pushed as `86fb7552` and
`96b89ce7`. This 31-function family remains original C; the starting count is
**127,104 C lines / 149 files / zero reference C**. The reusable optional fixture
now exercises **2,571 cases / 16 groups**:

- lifetime: 480; break/open states: 450; Spark creation: 56; trails: 28.
- homing: 288; owner effects: 144; spawners: 144.
- area callbacks: 224; area owners: 210; independent positive callbacks: 5.
- frame wrap: 216; missing definitions/owners: 15; RNG streams: 35.
- acquisition/filter/ties: 96; collision/impact boundaries: 20; tick rates: 160.

The clock helper synchronizes Owner and MonsterState inputs and asserts the
actual frame/FPS after setup and after each action. Tests cover signed/unsigned
frame-wrap behavior, exact/strict deadlines, tick rates 1/4/30/60, callback returns,
state combinations, owner/target loss, direction wrap, nearest-target selection,
parent/class/dead filtering, ties and scan order, all Spark stages, seven RNG
seeds, and clear/blocked rays. Independent checks verify water extinguishing,
nearest candidate identity, actual damage callback delivery, exact eight-Spark
trails, exact random-spawn counts and preserved collision return words.

The fixture uses actual retained movement, ray, target-search, area-damage,
scorch and deletion paths. It supplies core.ExtServer and the outer server only
for the temporary tests; previous contracts are unaffected. It saves/restores
spatial nodes, updatable/deleted lists, caches and type tables. Scorch uses valid
already-initialized type IDs. Flame, meteor, spider and scorch definitions use
real allocator templates with checked update/collision guards. Missing named
Spark lookup is covered; invalid indexed definitions are not used as an allocator
failure simulation. Pointer inputs retain the original nonnull preconditions.

Creation tracing records exact owner identity. Spatial-list nodes have stable
identities across ASLR. Lifecycle retains ownership of every created object,
including objects later scheduled for deletion, with exact zero-leak checks.
The fixture's C code consists only of dispatch and event recorders; no production
algorithm is copied into a test reference implementation.

Prior **17,743** effects/timing/equipment/inventory/resource/shop/trade contracts
remain unchanged (52.968s). The full original-C captures pass and match byte-for-byte (6.126s / 5.901s);
all 16 hashes are locked. Local evidence is in
build/port-temporary-updates/c-locked-{source,repeat}-temporary-*.json and associated
logs. Earlier 15 groups repeated byte-for-byte and remain unchanged after adding
the tick-rate group. Commit/push the locked baseline before production conversion.

## Native conversion

The 31 implementations now live in legacy/temporary_projectiles.go and
legacy/temporary_effects.go. Thin exports retain all original entry points for
C registrations and callers; the projectile callbacks and Spark-producing Go
callers invoke Go directly. No production C implementation remains as a test
reference. Removed 980 physical C lines: **126,124 / 149 files / zero reference C**.

All 2,571 cases pass (5.742s), and the 16 native-second full captures are
byte-identical to the locked original C. The first comparison detected a typed-nil
owner interface in Spark creation, corrected before qualification (also audited
for spider creation). Meteor impact retains the original nonpositive-height
comparison, including its NaN branch semantics. Integer wrap, raw callback return
words, RNG order, and original float spills remain explicit in the implementation.

The owner-loss test also exposed a fixture lifetime problem: SetPos can queue an
object in the retained C collision-update list. Both list globals are now
saved, cleared for the case, and restored before freeing its objects. Original C
at 14eb6ed2 with this cleanup passes all existing hashes unchanged (5.647s), with
full c-isolated captures byte-identical to c-locked-source. No hash was relaxed
or regenerated to match native behavior.

Qualification evidence is in build/port-temporary-updates. All three production
binaries built and were checked as ELF32/i386, GO386=sse2, CGO enabled. The full
suite exactly matches the baseline failure multiset: 1,553 failure entries,
15 packages pass, 3 fail, 32 skip; no added or removed failures. Fresh unchanged
repeat-a gameplay passed in 52.442s under Xvfb with null audio and overrides
disabled. Accumulated port tests, including all 20,314 focused family cases, passed in
default/server/highres (153.927s / 112.748s / 112.020s).
