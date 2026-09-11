# Object state, geometry and ownership

Ported 44 functions in 43 address blocks / 1,134 C lines from GAME3_3.c. Two adjacent forward declarations remain in C; the unmarked
`nox_objectCollideDefault` no-op is included. Include sync
bit updates, animation/elevation/buffs/item attributes, spawned-object cleanup,
unit classification, distance/direction/front tests, coordinate updates,
on/off/freeze/pet state, ownership notifications and nearby-door collisions.

Reuse the guarded attack/world/owner fixture with an optional object-state spec
and dispatch range 1200+. Identify return pointers (including object+688 and
input+16), record byte versus word results, guard the five-word modifier input,
and capture all 32 synchronization words. Preserve actual type definitions and
player protection records. Install only dependency dispatch/recording seams;
use real native Object methods where they already implement the behavior.

Before conversion, exercise masks/class-specific synchronization, all player
visibility bits, actual defaults versus changed values, signed and floating
boundaries, spatial boundary cases, owner/inventory chains, lifecycle cleanup,
and successful state changes. Repeat full original-C captures, lock hashes and
commit/push first. Preserve every prior capture hash. Convert the connected
batch, qualify once across variants/builds/full-suite/headless gameplay, update
physical tracked C_LOC and recovery docs, commit/push, summarize and continue.

The batch is complete. Original-C baseline `c300fd2d` was committed and pushed
before conversion.

## Candidate address blocks

- 004E44F0: `void nox_xxx_unitNeedSync_4E44F0(nox_object_t* a1)` (5 lines).
- 004E4500: `int* sub_4E4500(nox_object_t* a1p, int a2, int a3, int a4)` (27 lines).
- 004E4670: `int* nox_xxx_unitSetOnOff_4E4670(int a1, int a2)` (33 lines).
- 004E46F0: `void nox_xxx_unitRaise_4E46F0(nox_object_t* obj, float a2)` (27 lines).
- 004E4880: `int* nox_xxx_servMarkObjAnimFrame_4E4880(int a1, int a2)` (25 lines).
- 004E48F0: `int* nox_xxx_setUnitBuffFlags_4E48F0(int a1, int a2)` (30 lines).
- 004E4990: `int* nox_xxx_modifSetItemAttrs_4E4990(nox_object_t* a1p, int* a2)` (54 lines).
- 004E4A70: `double nox_xxx_objectGetMass_4E4A70(int a1)` (3 lines).
- 004E5AD0: `void nox_xxx_playerRemoveSpawnedStuff_4E5AD0(nox_object_t* a1p)` (27 lines).
- 004E5B50: `int nox_xxx_isUnit_4E5B50(nox_object_t* a1p)` (18 lines).
- 004E5B80: `int sub_4E5B80(nox_object_t* a1p)` (27 lines).
- 004E5BF0: `void sub_4E5BF0(int a1)` (36 lines).
- 004E6BD0: `int sub_4E6BD0(int a1)` (5 lines).
- 004E6C00: `double nox_xxx_calcDistance_4E6C00(nox_object_t* a1p, nox_object_t* a2p)` (47 lines).
- 004E6CE0: `int sub_4E6CE0(float2* a1, float2* a2)` (65 lines).
- 004E6E50: `int nox_server_testTwoPointsAndDirection_4E6E50(float2* a1, int a2, float2* a3)` (10 lines).
- 004E7190: `void nox_xxx_teleportToMB_4E7190(uint8_t* a1, float* a2)` (9 lines).
- 004E7290: `int nox_xxx_objectUnkUpdateCoords_4E7290(nox_object_t* a1p)` (35 lines).
- 004E7470: `void nox_xxx_spawnSomeBarrel_4E7470(int a1, int a2)` (50 lines).
- 004E7540: `void sub_4E7540(nox_object_t* a1p, nox_object_t* a2p)` (19 lines).
- 004E75B0: `char nox_xxx_objectSetOn_4E75B0(nox_object_t* obj)` (23 lines).
- 004E7600: `int nox_xxx_objectSetOff_4E7600(nox_object_t* obj)` (22 lines).
- 004E7700: `int sub_4E7700(int a1)` (20 lines).
- 004E7980: `int nox_xxx_inventoryGetFirst_4E7980(int a1)` (3 lines).
- 004E7990: `int nox_xxx_inventoryGetNext_4E7990(int a1)` (12 lines).
- 004E79B0: `int sub_4E79B0(int a1)` (9 lines).
- 004E79C0: `char nox_xxx_unitFreeze_4E79C0(nox_object_t* obj, int a2)` (36 lines).
- 004E7A60: `char nox_xxx_unitUnFreeze_4E7A60(nox_object_t* obj, int a2)` (38 lines).
- 004E7B00: `void nox_xxx_unitBecomePet_4E7B00(int a1, int a2)` (18 lines).
- 004E7B60: `void nox_xxx_monsterRemoveMonitors_4E7B60(nox_object_t* a1p, nox_object_t* a2p)` (20 lines).
- 004E7BC0: `int sub_4E7BC0(int a1)` (11 lines).
- 004E7BE0: `int nox_xxx_unitIsCrown_4E7BE0(int a1)` (23 lines).
- 004E7C30: `int nox_xxx_unitIsGameball_4E7C30(int a1)` (23 lines).
- 004E7CF0: `int nox_xxx_unitCountSlaves_4E7CF0(int a1, int a2, int a3)` (19 lines).
- 004E7DE0: `int sub_4E7DE0(int a1, nox_object_t* item)` (48 lines).
- 004E7F10: `char* nox_xxx_unitPostCreateNotify_4E7F10(nox_object_t* a1p)` (27 lines).
- 004E8110: `char* sub_4E8110(int a1)` (53 lines).
- 004E81D0: `int sub_4E81D0(nox_object_t* a1p)` (19 lines).
- 004E8340: `void nox_xxx_fnFindCloseDoors_4E8340(float* a1, int a2)` (15 lines).
- 004E8390: `int sub_4E8390(int a1)` (10 lines).
- 004E83B0: `unsigned char* nox_xxx_collideMonsterEventProc_4E83B0(int a1, int a2)` (5 lines).
- 004E83D0: `unsigned char* nox_xxx_collideMimic_4E83D0(int a1, int a2)` (25 lines).
- 004E8460: `void nox_xxx_collidePlayer_4E8460(int a1, int a2)` (105 lines).

## Original-C baseline

Before conversion, all 44 functions were original C. **2,757 cases / 53 complete capture groups**
repeat byte-for-byte in separate processes and are locked in
`src/object_state_porttest_test.go`. Local captures: `build/port-object-state/`
with `c-final-*`, `c-confirm-*` and `baseline-hashes.json`.

Coverage includes every synchronization bit and all 32 output words; class and
default-state branches; elevation/animation/buff/modifier updates; all checksum
inputs, signed shorts and health presence; shape distances and direction edges;
teleport admission; inventory/owner chains, cleanup, pet monitoring and special
owned items; freeze/unfreeze; shipped barrel/crate loot tables across 200 seeded
cases; door state and monster/player collision events. Explicit assertions check
sync/on flags and animation. Captures include successful loot creation and
berserker impacts that reduce actor health from 50 to 38.

Fixture findings:

- Freeze returns a signed char containing the low byte of the actual action-stack
  pointer. The fixture verifies that byte against the post-call stack address
  before replacing it with stable identity 89400. Four otherwise-identical cases
  exposed this address dependence during the first repeat. No production return
  behavior is changed.
- Player collision targets receive guarded health records and the existing
  damage recorder. Ability-disable requests are recorded and forwarded to the
  real fixture server's ability-state owner; the application-global network
  transport is outside this caller's scope.
- Loot tables and strings come from the shipped blob, with the same 26 pointer
  relocations as runtime initialization. Actor-name setup preserves the
  case-sensitive Barrel prefix.
- Two forward declarations stay in C. The unmarked `nox_objectCollideDefault`
  no-op is included, making 44 functions and 1,134 C lines to remove.

At the original-C baseline, production C was **120,952 lines / 149 files /
zero reference C**. This baseline was committed and pushed before conversion.

Locked combined regression: **34,380 cases / 315 groups**, passed in 105.687s.

## Native conversion

All 44 functions are implemented in five Go files (sync, geometry, ownership,
collisions and C ABI exports). Existing Object methods own sync updates,
collider bounds and AI stack changes. Original C bodies are removed; the two
adjacent forward declarations remain. No C reference algorithms are retained.
All 2,757 cases / 53 hashes match on the first native run (9.310s).
Physical tracked production C: **119,818 lines / 149 files** (−1,134).

The broad regression caught a C bridge ownership issue in ShopStockLoading:
the no-change modifier return points into a temporary Go input array. Shop,
generator and quest callers now call the native attribute helper directly.
The targeted rerun passes all object-state cases plus the 129-case shop stock
baseline (10.551s); no expected captures changed. Final qualification below includes this caller change.

## Qualification

Accumulated port tests, including 34,380 focused cases / 315 groups, pass in
default/server/highres: 178.994s / 161.251s / 162.972s. All 53 final native captures match the
original-C files byte-for-byte (8.755s). Three production binaries
are verified ELF32/i386, GO386=sse2 and CGO enabled. Full-suite failure identities
and multiplicities match exactly: 1,553 entries; 15 packages pass, 3 fail and
32 skip. Fresh unchanged repeat-a headless gameplay passes in
37.775s using Xvfb and null audio, without updating expectations.

Local evidence: build/port-object-state (native-final captures, variant logs,
binaries and qualification.json), and baseline/runs/object-state-final.
Next: [object initialization and reward generation](REWARD_GENERATION.md),
21 address blocks / 1,862 C lines in one connected batch.
