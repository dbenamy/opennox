# Proposed next batch: temporary objects and projectile updates

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
