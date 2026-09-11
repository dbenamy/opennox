# Grid tile lookup — 2026-09-11

Scope: 411160 and its public Go wrapper. The corrected C baseline includes the
lower-bound repair d765f9f4 (see GRID_BOUNDS.md). Float-to-int C helpers remain
for C callers because exporting each one adds avoidable boundary overhead.
The Go owner will use native conversion without those extra crossings.

## Repaired-C baseline

325,793 cases pass: 196,608 nominal grid-index/local-boundary inputs, 27,508
local-coordinate/list-prefix cases, 34,969 float32 ULP boundary pairs, 484 raw
IEEE pairs, 60,000 finite/random-word pairs, and 6,224 public-wrapper calls.
Inputs cover all nominal indices 0..127, all 46² local positions, every list
prefix 0..12, both diagonal equalities, signed invalid-conversion cases and
neighbor selection. Physical selection and ULP tests disable lists so an
identical list result cannot hide a wrong cell. Finite random cases densely
sample map coordinates instead of mostly testing distant rejected values.

The oracle models PC53 double arithmetic for `(x+11.5)*0.021739131` followed
by the observed float32 spill; unscaled `x+11.5` spills separately before
truncation/remainder. Converter results are independently decoded from IEEE bits.
Geometry uses two diagonal half-planes, with explicit local-coordinate transforms
and last-match list overrides. Engine PC53/round-nearest remains the grid-math
contract; the converter itself has separate PC/RC preservation coverage.

The fixture owns a guarded 128-row pointer array and guarded rows of 128×44-byte
cells, distinct raw fallback values per cell/half, guarded C nodes and input.
The grid and nodes stay immutable during each batch and are compared in full
afterward. Input guards and the installed grid pointer are checked every call.
The EDGE table and adjacent guards are restored and verified, as is the actual
previous grid pointer. Public-wrapper calls are tested with and without lists.

## Caller-boundary benchmark

The repaired C owner costs 18.87–19.16 ns per call from a C loop. Its public Go
wrapper (which allocates a C point, calls C and frees it) costs 386.9–409.2 ns.
Both benchmarks validate checksums and fixture state. They are microbenchmarks,
not measurements of complete game frames. Compare both routes after replacement.

Production C baseline: **140,455 physical lines**, 153 files, zero reference C.
Artifacts: build/port-grid-lookup. Native implementation and qualification remain.
