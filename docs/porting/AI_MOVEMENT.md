# AI movement actions — 2026-09-11

Connected batch: random walk (545020 and private 545090), confusion (545140),
face location/object/angle and set angle (545210/545240/545300/545340/5453E0).
Six action registrations previously went from Go to C. All char returns are discarded
at the action interface; internal returns only forward them. Native registration
removes all eight C bodies/entry points together. Audio and combat capability
callees, and 534120's other C callers, remain outside this batch.

One guarded C-owned object/monster/definition/target fixture and synthetic grid
and direction table cover the entire batch through the real action registry.
The C baseline captures every changed object/monster word, both directions, RNG
indices and stack-change state. Read-only data and guards are checked; global
server binding, grid, direction table and engine flags are restored. Golden
hashes retain those complete per-case state records in compact form. During
conversion the real C dispatcher and native registration also compared every case
directly, reporting the first differing input/state. The temporary C dispatcher has now been removed.

Corpus: all 65,536 current/target angle pairs; all 65,536 low-16-bit set-angle
values; all 65,536 signed starting walk directions with running/terrain variants;
5,120 facing point cases; 8,000 confusion seed/capability/stack combinations;
54 start/end/cancel cases; one explicit full-expression rounding discriminator.
The initial baseline passes in 4.2 seconds, including registry-vs-C comparison.
Finite boundary, zero-distance and nonfinite facing vectors are included. Facing
current/target table indices are valid 0..255; arbitrary outside-table addresses
are outside this contract. Set-angle and random-walk wrap are separately exhaustive.

Precision: PC53 double deltas; float32 length and normalized X; normalized Y
remains double for the location turn but spills for the alignment dot test. The
dot test spills its first product across a C table lookup. Actual mutable table
reads remain authoritative. Rounding direction*30 before adding position changes
the explicit terrain case and must fail. Disassembly and logs: build/port-ai-actions.

Current stage: direct C/native parity passed, then all eight C bodies and their
declarations were removed along with the temporary C test dispatcher. Production
C is now 140,260 physical lines (minus 182), 153 files, zero reference C. Shared
qualification passed in all three configurations. Baseline: 425e9c78; preserved state hashes remain.

Direct comparison now passes all 209,783 cases. It found a one-ULP running-force
difference: the actual C compiler spills running speed to float32 before both
force products, despite the decompiler's double temporary. Native Go now matches
that spill. The location turn also retains normalized Y at double precision.

Initial repeated-update measurement (200,000 updates/path, guarded fixture,
matching checksum and complete final state) under GO386=softfloat: C dispatcher
653 ns / native registry 802 ns with terrain; 287 ns / 384 ns with terrain
skipped. These include the small fixture setup cost and are VM microbenchmarks,
not a whole-game slowdown estimate. Hardware-floating-point qualification was required separately; the user confirmed that pre-SSE2 CPU support is unnecessary. The build driver
and recovery environment now explicitly select GO386=sse2; C x87 flags remain
unchanged. The accumulated tests, all builds, full suite and fresh
gameplay qualify the new Go target setting.

The installed ccache was trialed without changing GCC flags. It reported zero
hits across 2,354 cacheable compiler calls during the changing cgo builds; Go's
compiler invocations include changing work paths/build seeds. No cache speedup
was established. Direct GCC/G++ remain the recovery/default configuration; no
cache validity checks were relaxed to force hits.

SSE2 direct comparison preserves every baseline state hash. In the first paired
measurement, C dispatcher/native registry costs were 337/222 ns per terrain-aware
update and 257/165 ns with terrain skipped. A second pair also favored Go, though
absolute VM timings varied. Both measured 200,000 updates and compared the full
final state plus a direction checksum. Permanent repeated-update tests retain
those original-C state hashes without keeping C bodies for tests.

The SSE2 full suite matches the exact known baseline (1,553 failure entries;
15 passing/3 known failing/32 skipped-no-test packages). Accumulated tests exposed
a latent waypoint fixture defect: typed Go pointer assignments could present an
old poisoned pointer to the write barrier, despite C-owned storage. Pointer slots
now receive zero byte values before typed assignment; padding/guard coverage is
unchanged. A GC-pressure focused run and the accumulated matrix qualify the fix.

Two additional dot-product threshold cases prove that omitting the first-product
float32 spill changes the result in each direction. They compare against 534120,
which remains production C for its other real callers; no C body is retained only
for testing. All AI checks including these cases pass in all three variants.

All accumulated port tests pass for default, server and highres. All production
binaries build as ELF32/80386 with GO386=sse2 recorded in their Go metadata.
The first batch reused one fixture and one qualification matrix for six actions;
it also incurred one-time CPU-target qualification and the fixture repair. The
182-line reduction alone does not establish a sustained throughput multiplier.

Fresh ai-movement-port gameplay exits 0 against both preserved screenshots with
overrides disabled. Artifacts are in build/port-ai-actions and
build/baseline/runs/ai-movement-port.
