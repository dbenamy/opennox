# Map-generation rooms, geometry and decoration selection

## Scope

This batch follows sustained spells (`efe8316c`). It covers 60 connected functions
and 1,405 C lines in GAME4_2.c: coordinate rounding, occupancy grids, room lists,
room/hall allocation, reciprocal connections, corridor trimming, exclusion
rectangles, entry-point deduplication, required/weighted decoration selection,
and the map-generation RNG adapters. Terrain painting and prefab/object
population are separate subsequent batches. After removing this scope virtually,
only platform RNG and nullsub_28 remain as external callees.

Production remains C; its complete repeated baseline is locked. Physical production C is
110,283 lines in 149 files, with no reference C. Candidate source, operations and
ABI inventories are under `build/port-map-rooms/`; the initial audit finds 49
entries with production references and 11 without. Recheck Go getters and
wrappers before retiring any ABI.

## Fixture and coverage

`legacy/map_rooms_porttest.go` invokes the original entry points through a thin
signature adapter. Inputs use aligned C allocations with trailing guards. The
alignment makes char-return pointer low bytes deterministic without masking
results, and the input addresses remain valid free() bases. Room, exclusion,
occupancy-grid and scratch allocations use the real production allocation paths.
The fixture tracks live/freed regions to avoid reading released storage or
freeing it twice, and records all live bytes, links, grid cells, globals, slots,
return bits and a trailing random value after each sequence.

Install a fresh real `platform.New()` RNG owner and restore the previous platform,
so seeded cases cannot contaminate other tests. Preserve and restore all five
named globals and the x87 control word on a locked thread (PC53/nearest for the
fixture). The headless harness initially leaves two required blob tables zero:
supply actual startup data from the embedded blobs, capture all six words, and
restore the prior values. These are epsilon 0.1 and opposite directions 1,0,3,2.
Independent point-deduplication and reciprocal-direction checks distinguish this
from merely finding an arbitrary back-link or hashing an uninitialized table.

The corpus exercises every entry point, half-cell rounding and neighboring
floats, signed coordinates and occupancy edges, zero/negative dimensions where
the individual operation defines them, allocation initialization and cleanup,
list insertion/unlinking, saturated reciprocal connections, exclusion boundaries
and selective removal, all room kinds, corridor trimming, entry-point capacity
and epsilon, decoration flags/counts/required assignments, seeded RNG ranges and
variety. Independent positive contracts accompany full-state hashes.

## RNG compatibility correction before locking

The two remaining consumers of `nox_platform_rand` scale by a 15-bit CRT maximum,
but the Go platform returns a wider integer. With seed 12345, original C requested
a float in [-5,8] and returned 575182.218727404. The random corpus failed its
independent bound check and took 52.35 seconds, largely in rejection sampling.
Constrain only the C compatibility export to the low 15 bits. The same corpus
then passes in 0.05 seconds; the broader Go platform API is unchanged.

This changes generated layouts for a seed and is an intentional correction,
not exact preservation of the broken adapter. It follows the standing policy
for confident reversible decisions and is recorded for later seed-compatibility
review in [DECISIONS.md](DECISIONS.md). Before/after evidence is in
`rng-before-fix.json`, `c-boundaries.log` and `c-rng-fixed.log` under the batch's
build directory. Commit/push the corrected C baseline before replacing algorithms.

## Locked C baseline and next gate

All **2,997 cases / 122 complete captures** repeat byte-for-byte (0.720s / 0.559s).
Hashes are locked in `src/map_rooms_porttest_test.go`. The prior 51,001 focused
cases / 701 groups pass unchanged (root package 220.916s). Enforced new plus
6,442 existing spell cases pass in 31.165s. Commit/push this corrected C
baseline before native conversion.

The native scope retains 49 C ABIs and retires 11 internal helpers. Also remove
the unused one-line nullsub_28 when its only two callers move; this gives 1,406
physical C lines removed, with 60 actual native algorithms and target C total
108,877. Raw C allocation primitives must preserve shared ownership while other
C code allocates/frees these records; common/alloc's tracking registry cannot own
untracked C pointers. No C algorithm will be retained solely for comparison.

After conversion, require unchanged complete captures. Run accumulated default,
server and highres checks, build/verify all three production binaries, compare
the asset-backed full-suite failure multiset, and replay fresh unchanged headless
gameplay. Record the new physical C count, commit/push and continue.
