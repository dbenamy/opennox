# Final engine allocation and process boundary

Status: original contracts and broader baseline qualified. No allocator
conversion is installed. The preceding qualified source is `49c8ef62`.

## Scope and acceptance

Retire remaining engine libc allocation and abort calls, aligned fixture imports
and obsolete compiler directives on the qualified Linux 386/SSE2 target. Preserve
raw versus safe tracking, explicit lifetime, alignment, raw-address storage,
zero-size/overflow handling and failure disposition. External SDL2/OpenGL/OpenAL
bindings remain. Keep pointer storage outside Go's managed heap.

Five new allocator roots cover successful and zero-product allocation, overflow,
failed realloc retaining old contents, growing/shrinking prefixes, cross-thread
handoff, raw-address lifetime through GC, mixed-size reuse and explicit alignment.
They passed twice in normal and safe configurations against the original backend.
The new fatal root checks invalid room/painting fixture operations and missing
static/dynamic light drawables in separate processes: exit 2, no return, no deferred
cleanup or recoverable panic. It passed twice in safe/default/server/highres.

The first light test reused an effects fixture and encountered that fixture's
unrelated historical mapped access before the intended branch. The corrected test
supplies the real empty drawable owner directly. Safe runtime checks remain enabled;
production and existing captures are unchanged. Initial/diagnostic preflights are
not the accepted baseline; use the corrected folders in the working evidence.

The original owner selection combines established allocation families, all active
observer roots, aligned numeric fixtures and light transfer contracts: 162 roots
per normal profile, twelve focused repeats and twelve safe roots. Direct allocator
and string/clock contracts run separately. Reuse the preceding production baseline
only because all new source files are test-only; rerun every production gate after
conversion. Native acceptance requires the full compiled root corpus on all three
profiles, exact known-suite comparison, safe/static checks, build/ABI verification,
fresh save/load and original asset integrity.

Record production-dispatch allocator benchmarks at one/two CPUs and three repetitions
of real theme/population/grid owners, including resident memory. Run performance
measurements without another scheduled port build/test job. Compilation or a
microbenchmark alone does not establish acceptable behavior or performance.

## Preliminary allocator evaluation

An isolated original-backend probe and an uninstalled mmap prototype both pass the
new primitive contracts. The prototype uses size-class pools and explicit metadata.
Mapping each block above 32 KiB was rejected: exploratory allocation churn was
roughly 6–7 microseconds/op against libc's 0.3–0.5. Bounded reuse through 1 MiB
reduced the prototype to roughly 0.46–0.62 microseconds/op. These runs overlapped
another regression sweep and are exploratory, not acceptance measurements.
Substantial shrinking releases excess storage when possible; a failed shrinking
move may retain the original block. The prototype still needs complete integration
review, platform scope and qualified performance/memory checks.

Working evidence: `build/port-final-engine-boundary/`; earlier isolated experiments
and full original scratch source are under
`build/port-fixture-allocation-observers/next-drafts/`. Completed scripts are
single-use. Primary handles this batch because Luna quota remains unavailable.

## Accepted original baseline

All 162 owners and twelve focused repeats passed freshly in each normal profile;
twelve safe roots and ten allocator/five string-clock contracts passed. The corrected
new-contract preflights use the exact same source. Production remains unchanged
from `49c8ef62`; the five new files are test-only, including a temporary C alignment
probe. Fixture imports are temporarily three in this baseline; the checkpoint's
qualified conversion counts continue to describe `49c8ef62`.

See [baseline](final-engine-boundary-baseline.json) for exact source, binary hashes,
started/passed names and performance samples. Median RSS in the three owner runs
was 177,104 KiB. Owner timings and allocator medians are recorded individually;
compare matching CPU settings instead of combining them into a single number.
