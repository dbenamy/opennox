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

## Conversion status

Original-C baseline recorded; native implementation is drafted but not installed.
Production C remains **138,824 physical lines**, 153 files, zero reference C.
Run `TestAICombat` with `porttest` from `src` using `build/baseline/env.sh`.
