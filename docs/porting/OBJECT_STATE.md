# Candidate next batch: object state, geometry and ownership

43 candidate address blocks / 1,136 physical C lines in GAME3_3.c. Audit adjacent
forward declarations before locking the final function-only count. Include sync
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

No production object-state conversion has started. Current damage qualification
is independent and must finish before changing source for this next fixture.

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
