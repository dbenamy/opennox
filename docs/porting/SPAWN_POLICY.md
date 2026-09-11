# Spawn ownership, admission and culling

Scope: all 15 GAME4_1 functions from 50D780 through 50E210 (550 physical
C lines), plus the 14-line 50D890 periodic caller in server__system__server.c.
The connected batch totals 16 bodies / 564 C lines. Startup/shutdown, generator
admission/registration, death release and periodic culling share these lists.

## Original-C baseline

All twelve groups repeat exactly. The complete locked baseline passes (4.532s);
generator and adjacent callback/creation/penalty/death regressions pass as well
(16.035s before the final overlapping-view group was added). No C is removed
in the baseline commit.

1,582 locked cases across twelve groups in src/spawn_policy_porttest_test.go:
8 registration sequences, 120 far culls, 400 visible culls, 52 tick schedules,
4 pool exhaustion/reuse sequences, 8 zombie cleanup cases, 24 glyph cases,
315 direct candidate filters, 288 admission rectangles, 180 indexed admission
cases, 48 temporary-list exhaustion cases, and 135 overlapping visibility groups.

Tests reuse the generator's real sparse player list (indices 1, 7, 31), object
allocator, map, ray tracing, balance lookup, and C-owned SpawnClass and
MonsterListClass pools. Up to 100 guarded objects have real Go server handles,
with only the 772-byte C ABI prefix captured. Their full 2200-byte update data,
list links, owner counts, global counters, player records, glyph inventory and
64-byte glyph init records are captured after each action. Deletion calls retain
order. Allocator handles and intrusive node addresses receive stable identities.
All indices, glyph allocations, player memory, globals, and pools are restored.

Registration checks head/tail unlink, idempotence, the real 96-record capacity,
failed allocation, reuse after removal and repeated cleanup. Far culling checks
strict distance700, exact joined==1, busy state and no-player behavior. Visible
culling checks deletion counts, player-specific capacities, overlapping views,
zero distances, sorting permutations and resource exhaustion. Tick checks both
modulo gates across four tick rates and frame boundaries including wraparound.
Candidate tests preserve the C numeric conversions of float views of class and
flags; indexed cases exercise the actual spatial callback path as well.

Glyph tests exercise real allocation/inventory insertion, copied spell words,
cleanup only while associated, and repeated cleanup. Zombie and VileZombie
(type IDs 2 and 3 in the shared fixture) bypass ordinary death release.

The fixture does not inject a nil return into alloc.NewClass: that retained
primitive panics on allocation failure. Fixed-pool NewObject exhaustion is
exercised through the real allocator. Glyph creation uses a synthetic type with
real init storage and no asset callbacks; the inventory engine remains outside
this conversion. Linked records use C-owned memory throughout.
