# Map-generation rooms, geometry and decoration selection

## Scope

This batch follows sustained spells (`efe8316c`). It covers 60 connected functions
and 1,405 C lines in GAME4_2.c: coordinate rounding, occupancy grids, room lists,
room/hall allocation, reciprocal connections, corridor trimming, exclusion
rectangles, entry-point deduplication, required/weighted decoration selection,
and the map-generation RNG adapters. Terrain painting and prefab/object
population are separate subsequent batches. After removing this scope virtually,
only platform RNG and nullsub_28 remain as external callees.

Baseline `c470a93d` was committed and pushed before conversion. All 60 routines
now have native implementations; 49 C ABIs remain for callers, and 11 internal
helpers plus the obsolete no-op are retired. The source audit finds no remaining
references to those 12 symbols. Physical production C is 108,877 lines in 149
files, with zero reference C. Qualification is complete.

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
build directory. The corrected baseline was pushed before replacing algorithms.

## Locked C baseline

All **2,997 cases / 122 complete captures** repeat byte-for-byte (0.720s / 0.559s).
Hashes are locked in `src/map_rooms_porttest_test.go`. The prior 51,001 focused
cases / 701 groups pass unchanged (root package 220.916s). Enforced new plus
6,442 existing spell cases pass in 31.165s. This corrected C baseline was
pushed as `c470a93d` before native conversion.

The conversion retains 49 C ABIs and retires 11 internal helpers plus the
unused one-line nullsub_28. This removes 1,406 physical C lines and leaves
108,877. Raw C allocation primitives preserve shared ownership while remaining
C callers allocate/free these records; common/alloc's tracking registry cannot
own untracked C pointers. No C algorithm is retained solely for comparison.

## Native precision checks

The first capture run found one-bit differences in random point placement. The
x87 build retains room-center intermediates in double stack slots across RNG
calls, despite their decompiled float declarations. Preserve those intermediates
in float64; only RNG arguments and the output coordinates round to float32.
All original 2,997 cases / 122 complete captures then match byte-for-byte.

Assembly review also found that exclusion comparisons retain half-unit insets
at double precision. Four additional independent cases at 2^23 distinguish
contact from overlap; the original C produces false/true/false/true. The draft's
early float32 stores incorrectly admitted both contact cases, and the new tests
reproduced those failures before the correction. These four contracts supplement
the locked capture corpus; no baseline hashes were changed. Diagnostic C and
assembly evidence stays in the ignored build directory, with no C test reference
added to the repository.

## Qualification

Native final: all 2,997 cases / 122 captures match C unchanged in 0.657s; all four
additional precision contracts pass. Accumulated 53,998 captured cases / 823
groups plus those contracts pass in default/server/highres (228.594s / 298.913s / 233.356s wall time).
All three production binaries are ELF32/i386, SSE2, CGO enabled; all 12 retired
symbols are absent. Asset-backed full-suite failure multiset matches exactly:
1,553 entries, 15 passing packages, 3 failing, 32 skipped. Fresh unchanged
repeat-a headless gameplay passes in 35.817s. See `qualification.json`,
`variants.json`, `builds.json`, `binary-verification.json`, `full-suite-comparison.json`
and `build/baseline/runs/map-rooms-port` for local evidence.

Production C is **108,877 lines / 149 files / zero test-reference C**, down 1,406
physical lines. The next scope/test plan is [MAP_PAINTING.md](MAP_PAINTING.md).
