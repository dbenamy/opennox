# Combat AI actions

## Boundary and original-C baseline

This batch owns FIGHT, BLOCK_ATTACK, BLOCK_FINISH, WEAPON_BLOCK, MELEE_ATTACK
and MISSILE_ATTACK, their lifecycle methods, attack-choice/stack builders,
nearby-unit callbacks and private buff cleanup (531B40–532610, 5341D0).
Spell policy, weapon/strike engines, projectile prediction, shield scanning,
vision, damage and audio delivery remain shared callees.

The original C passes 17,408 generated cases and 170 independently asserted
contracts using the shared guarded AI fixture. Generated complete-state hash:
`68c25cbb61de2e5f1603f76f7349e9595ebaae813a266f61b7f8099ead1221fe`.
Contract complete-state hash:
`3375bac303a78fc6473ce0dfe489428def83a47b48146ac3d34224bc55add2b3`.

Coverage includes stack capacity, unsigned frame/deadline boundaries, lifecycle
flags, all strike animation gates, hit/miss sounds, event 13, cooldown waits,
byte-wrapping stamina behavior, friendly-in-way face/wait/flee stacks, matching
dead-target scans, actual shield projectile detection, missing projectile types,
actual wall rejection/delete and successful create, target prediction, directions,
and mixed-precision scan/projectile arithmetic. The playerlike false-return
weapon path runs the retained attack engine with a missing weapon definition;
the melee morph case uses a real C-owned monster/player-data pointer cycle.
The tests do not simulate every retained weapon/spell engine internally.

The fixture records normalized object/update-data mutations, action records,
RNG indices, globals, script calls, real deferred audio queues, indirect C strike
calls and complete projectile C-layout state. A small porttest-only C callback
counts strikes; it contains no copied combat implementation. Supplemental player
and target-data blocks have guards and record changes. Runtime direction and
facing tables come from the embedded blob data and are restored afterward.

Independent checks caught and repaired fixture gaps before recording the baseline:
non-target missile arguments must be zero, spatial objects need a collider,
projectile rejection needs a ray that actually reaches the wall, and shield
facing reads the projectile's previous position. Frozen snapshots and arithmetic
review artifacts are in ignored `build/port-ai-combat`.

## Native conversion

Original-C baseline commit: `e1cf21bf`. The native registered actions and private
helpers match both complete-state hashes. Eighteen C bodies/declarations and the
six old C registrations are removed, without replacement action exports.
Production C: **138,212 physical lines (minus 612)**, 153 files, zero reference C.

An additional 192 lifecycle assertions cover every end/cancel callback across
running masks. Two explicit precision regressions protect the original-C results:
scan minimum `0x3f800348` and projectile velocity words `3221729642/3221726419`.
The former was verified with the original callback in an ignored standalone
386/PC53 probe; the latter is original generated case 15362.

Compiled x87 behavior matters here: scan length squares both full-precision
deltas before the Y spill used by the facing test. Missile length and velocity
also retain both deltas; the denominator is float32. Spawn X rounds before the
ray addition, while spawn Y remains wide through that addition and rounds
separately for object creation. An initial native draft differed in 416 projectile
cases; correcting these points restored the original hash. The tests preserve
these distinctions without retaining an original C combat implementation.

## Qualification

Accumulated tests pass in default/server/highres (35.585/47.419/35.567 seconds),
with the extra combat precision/lifecycle run also passing. All three production
builds pass and report ELF32/Intel80386 with GO386=sse2. The full suite matches
exactly 1,553 known failure entries and the 15 pass/3 fail/32 skip package baseline.
Fresh `ai-combat-port` gameplay exits 0 against both preserved screenshots,
overrides off, under Xvfb/OpenAL null. Artifacts: `build/port-ai-combat`.
