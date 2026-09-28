# Final engine allocation and process boundary

Status: original baseline `71ff3095` is pushed; native conversion is installed
and its performance is accepted. Full qualification found a fixture slot-identity
mismatch after address reuse; the correction passed the original backend and
is installed for fresh native qualification. The preceding qualified conversion remains `49c8ef62`.

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

## Native allocator revision in progress

The first integrated allocator passed primitive, ownership, fatal, safe and all
162 selected owner contracts per profile. Acceptance stopped before the full
corpus: real-owner median times regressed 17–30%, and the two-CPU parallel
allocation benchmark was about four times slower. Alternating original/native
runs confirmed the owner regression. CPU profiles showed redundant clearing of
fresh anonymous mappings and extra mapping work; one global lock also serialized
concurrent allocations. These results remain under `rejected-v1/` in the working
evidence, with a relocation manifest preserving the original recorded paths.

Revision two uses per-size-class locks, registries updated only when mappings
change, and the operating system's zero-filled pages for first-use calloc slots.
Every allocation consumes that first-use guarantee; reused calloc storage is
explicitly cleared. Retaining at most two idle spans per class bounds cached
payload storage below 6 MiB and avoids two-worker mapping churn. Large direct
allocations are unmapped on release. Original contract tests and goldens are
unchanged; the native-only cache-budget contract reflects the revised policy.
The isolated revised implementation passed all seven primitive/storage roots ten
times. Integrated qualification and performance acceptance remain pending.

## Implementation and scope decisions

Use explicit OS mappings and size-class pools behind the existing raw allocator
API. Keep allocation metadata in Go-owned registries and payloads outside the Go
heap, preserving integer-address lifetime, alignment and explicit free semantics.
The allocator remains thread-safe for independently owned blocks and cross-thread
handoff; concurrent operations on the same block still require caller ownership.
Failure paths preserve the existing distinction between fatal RawMalloc and safe
wrapper panics. Realloc does not release the original block on failed growth.

Linux fatal paths send SIGABRT to the pinned current thread and unconditionally
exit if signal delivery returns. The non-Linux fallback exits but has no runtime
qualification. Windows allocation uses VirtualAlloc/VirtualFree and has a 386
compile-only check. The existing Windows-only WinSock C binding remains outside
the selected Linux milestone. No broader portability or whole-build cgo-free
claim is intended.

Retire legacy compiler/linker glue, including safe AddressSanitizer flags. Safe
Go allocation tracking and mapped-memory runtime checks remain enabled and tested;
these checks do not provide equivalent whole-process ASAN coverage. External
native dependencies and their own build directives remain unchanged.

Before the revision-two owner sweep, verified host process/open-file checks and
archive hashes allowed removal of 38 old, rebuildable Linux 386 Go cache archives
(1,886,093,312 allocated bytes). Source, assets, logs and retained binaries were
preserved; physical free space increased to about 3.79 GB. Recovery is automatic
cache rebuilding. The working `cache-headroom-*` records retain the exact list.

Revision two also passed all 162 owners and twelve focused repeats per normal
profile, but did not pass performance review. Real-owner medians were still
18–35% slower, and parallel allocation had unstable slow samples. Its profile
located most additional clearing cost in the painting fixture's explicit buffer
initialization, rather than calloc. Each fixture needs about 1 MiB of row buffers
in one class; retaining only two 64 KiB spans repeatedly discarded that working
set, causing mapping and page-fault costs on subsequent initialization.

Revision three retains up to 1 MiB of idle payload per size class, or two spans
for the largest class, with an 18 MiB total upper bound. It also separates class
locks to avoid false sharing. The native storage contract now exceeds each cache
limit before releasing blocks. Frozen original contracts/captures remain unchanged.
Run a normal-profile real-owner performance preflight before repeating broad
qualification; isolated allocation benchmarks did not predict these fixture costs.
Revision-two results and source remain under `rejected-v2/`.

The bounded-cache revision brought owner medians within 1–9% of the original,
with similar RSS, but two-worker allocation churn remained slower. Its CPU profile
was dominated by clearing reused buffers, with little lock-wait time. Uninstalled
thread-ID and Go pool-hint prototypes did not justify their overhead/complexity
and are not part of the selected implementation.

A small 386 exact-byte `REP STOSB` leaf improved measured clearing throughput
compared with Go 1.26's 386 vector loop on this VM. It needs no optional CPU
instructions or runtime internals; other architectures keep Go clear. The clear-only candidate combined revision three's bounded allocator with this
leaf. Native guard-byte contracts cover zero length, unaligned starts and
vector/page-size boundaries. Original raw allocation contracts remain unchanged.
Integrated performance review and milestone qualification are still pending.

A fresh isolated original-backend run kept two-worker churn near 220–264 ns/op,
so the native variability cannot simply be dismissed as general VM noise. The
clear-only integrated candidate matched real-owner timings and improved serial
churn, but still had slow two-worker samples. The current combined candidate uses
Go `sync.Pool` tokens as advisory large-buffer reuse hints, together with the
byte-clear leaf. It neither obtains OS thread IDs nor uses private runtime APIs.
The hint is not ownership: GC may discard it and threads may migrate. Live spans
remain rooted by the allocator registry, and an added storage-contract case checks
that GC and subsequent reuse leave live payload bytes unchanged. Cache eviction
rotates through bounded idle slots so stale hints cannot monopolize them.

## Performance decision on the final source

Accept the combined bounded allocator, advisory reuse hints and byte-clear leaf.
Alternating verified original/native owner binaries gave median times of
1.25→1.17 s for population ordering, 1.10→1.03 s for equipment sets, and
1.62→1.57 s for theme algorithms. Grid ownership is below useful timing
resolution. The separate three-run native RSS median was 173,612 KiB versus
177,104 KiB originally. No competing port jobs ran during measurements.

Five of six allocation microbenchmark medians improved. **Two-worker raw churn
remains slower:** 336.1 versus 257.6 ns/op (+30%), with native samples spanning
207.8–350.2 ns/op. The fresh isolated libc probe also stays faster. This is an
explicit residual tradeoff, not attributed away to VM noise. Accept it because
the real affected owners are at parity or better, memory use is comparable, and
serial/batched allocation improves. Revisit if a real concurrent engine workload
shows allocator cost; these samples do not establish universal hardware/gameplay
speed. Final regression, build and scenario qualification remains required.

The accepted timing run precedes broad functional qualification but has exactly
the same source and supplemental fingerprints. Reuse it without redundant timing
reruns while the source remains frozen. Full records are incorporated into the
final qualification report.

## Full-corpus fixture diagnosis in progress

The final-source full sweep found identical mismatches in hallway routes,
obstructions, candidate fallback and prefab connections in default/server. Stop
the failed sweep rather than continue the queued high-resolution run. The original
obstructions capture matches the current frozen hash; comparing all 72 cases
shows only `Slots[2]` changed, with live records and other captured state identical.
Painting snapshots re-normalize a saved slot's raw address even after its record
was released. The new allocator can reuse it for a new record, changing this
diagnostic identity. The base room fixture already snapshots the bound record ID.
An independent forced-address-reuse regression is being checked against verified
`71ff3095` source before correcting painting's slot snapshot. Frozen expectations
remain unchanged. Large failure captures were losslessly compressed with verified
SHA-256; `failure-capture-archive.json` maps their historical paths (478 MB reclaimed).

The forced-reuse regression fails on the original fixture for both a released
record and an interior alias, then passes with bound-record identity in normal
and safe modes. The corrected original also preserves all frozen hallway/prefab
hashes: default repeated twice, server/highres once. Comparing all 2,512 retained
hallway cases finds only saved-slot identity differences. See
[fixture evidence](final-engine-slot-identity.json) and the
[original-source patch](final-engine-slot-baseline.patch), applicable to `71ff3095`.
The native correction is installed, and focused owner/repeat coverage expands
from 162/12 to 166/16. No allocator or production source changed for this fix.

After installing the fix, all 166 owners and 16 focused repeats passed per normal
profile; thirteen allocator and six legacy roots passed in normal/safe, with the
initial fatal/safe/leaf gates repeated on the corrected source. Alternating owner
measurements now use the corrected fixture on both backends: population 1.21→1.22 s,
equipment 0.99→0.99 s, algorithms 1.59→1.61 s. Treat these as parity within about
1.3%, superseding the earlier owner timing comparison. The allocator-only benchmark
package and production source are byte-identical, so their measured samples are
reused explicitly. The remaining +30% synthetic parallel-churn tradeoff stands.
The full corpus and production gates are still pending; this is a recovery
checkpoint, not milestone acceptance.
