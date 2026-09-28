# Final engine allocation and process boundary

The Linux 386/SSE2 internal engine C-glue milestone is qualified. The original
baseline is `71ff3095`; external SDL2/OpenGL/OpenAL bindings remain unchanged.
See [qualification](final-engine-boundary-qualification.json) and
[inventory](final-engine-boundary-inventory-after.json) for exact evidence.

## Result and implementation

Selected default/highres production C imports fall from four to zero, and server
imports from five to zero. Direct project cgo packages fall from two to zero;
fixture C imports fall from two to zero. The baseline temporarily added one C
alignment probe, also retired. Standalone C, project headers, embedded production
C bodies and legacy C exports remain zero. The Windows-only WinSock binding is
outside the qualified Linux target and remains unchanged.

The raw allocator keeps payloads in explicitly owned OS mappings, outside Go's
heap. Go metadata records sizes, live ownership and free slots. Size-class locks
permit cross-thread handoff without a global allocator lock; realloc never holds
two class locks while moving data. Failed growth preserves the original block.
Substantial shrinking moves to smaller storage when possible. Zero-size/product,
overflow, alignment, lifetime and raw-versus-safe failure contracts are preserved.

Pools retain up to 1 MiB per class, or two spans in the largest class: an 18 MiB
upper bound on idle payload mappings. This excludes live fragmentation and Go
metadata. Excess idle spans and released direct allocations are unmapped.
Advisory `sync.Pool` tokens favor worker-local reuse of large buffers; they are not
ownership. GC may discard tokens and scheduling may change them without affecting
live allocations. Rotating eviction prevents obsolete tokens monopolizing the cache.

New anonymous storage supplies calloc's first-use zero guarantee; every allocation
consumes that guarantee. Reused calloc blocks are cleared explicitly. A small 386
`REP STOSB` leaf improves measured throughput on this VM; it requires no optional
CPU instructions or private runtime APIs. Other architectures use Go clear.
Guard-byte contracts check exact ranges, unaligned starts and zero length.

Fatal engine paths send Linux SIGABRT on the pinned current thread, then
unconditionally exit if signal handling returns. RawMalloc separately preserves
its original process-fatal failure disposition and diagnostic; safe wrapper panic
and nil-marker behavior remain distinct. Non-Linux fatal fallback exits but is
not runtime-qualified. Windows VirtualAlloc/VirtualFree support is compile-checked
on 386 only; broader portability is not claimed.

Engine compiler/linker preambles, including safe AddressSanitizer flags, are retired.
Safe Go allocation tracking and mapped-memory runtime checks remain enabled.
They do **not** provide equivalent whole-process ASAN coverage. External native
libraries retain their own binding/build directives.

## Original baseline and final qualification

Five new allocator roots were frozen against libc: boundaries/failure, thread
handoff, raw-address lifetime through GC, mixed lifetimes and aligned allocation.
They passed twice in normal and safe processes. The fatal root checks invalid
room/painting operations and missing static/dynamic light drawables: exit 2,
no return, no deferred cleanup and no recoverable panic. It passed twice in each
of safe/default/server/highres, on both original and final implementations.

The first original light fixture hit an unrelated historical mapped access before
its intended fatal branch. The corrected fixture supplies the real empty drawable owner
without disabling safe checks. Initial diagnostic runs are not the accepted
baseline. Production source at `71ff3095` is identical to qualified `49c8ef62`;
only original tests were added, so the original production gates were reused.
Every production gate was rerun after conversion.

Final-source checks passed:

- 166 selected owners and sixteen focused repeats per normal profile; twelve safe roots.
- Thirteen allocator/storage/clear roots and six legacy fixture/string/clock roots in both normal
  and safe configurations; thirteen allocator roots also passed with cgo disabled.
- All compiled root tests: 2492/2481/2492 passed in default/server/highres, with only
  the expected opt-in population diagnostic skipped in each profile.
- Static memory checks, 386 assembly declaration/disassembly review, safe build,
  three production builds, ABI/export checks and exact known-suite outcomes.
- Fresh save/load with reference comparison and all 1,654 original asset hashes.
- Windows 386 allocator compile-only check. Physical display/audible playback
  remain manual release checks.

Frozen original assertions and captures are unchanged. Native-only tests cover
cache release bounds, substantial shrinking, exact-byte clearing and live payload
survival when GC discards reuse hints. An explicit library selector initially
omitted the new clear root; its twelve selected tests passed, then the corrected
thirteen-root selection passed. This was a launcher omission, not a source failure.

## Address-reuse fixture correction

The first full native sweep exposed four hallway/prefab capture mismatches. All
2,512 hallway cases differed only in saved slot IDs: painting re-normalized a
released record's address, which the new allocator could reuse for another record.
The corrected fixture preserves the bound record identity, as the base room fixture already does.
A forced-reuse regression fails on the original fixture for both a released record
and an interior alias, then passes after correction in normal and safe modes.
Corrected original-backend hallway/prefab captures remain byte-identical in all
three profiles, including two default repetitions. Frozen expectations are unchanged.

See [original-backend proof](final-engine-slot-identity.json) and the
[reproducible patch](final-engine-slot-baseline.patch), applicable to `71ff3095`.
Focused coverage expanded to 166 owners and sixteen repeats. Recovery commit
`33e52f7c` saved the implementation and its explicit pending full gates before
qualification restarted. Failure captures were losslessly compressed with SHA-256
verification, reclaiming 478 MB; the proof record maps their historical paths.

## Performance and accepted tradeoff

Measurements ran without competing port jobs. Alternating binaries are verified
against their own recorded original/native source and binary hashes. Median real
owner timings were:

| Owner | Original | Native |
| --- | ---: | ---: |
| Population ordering | 1.21 s | 1.22 s |
| Equipment sets | 0.99 s | 0.99 s |
| Theme algorithms | 1.59 s | 1.61 s |

Grid ownership is below useful timing resolution. The corresponding native
RSS median was 172,500 KiB versus 171,724 KiB originally.

| Allocation benchmark | Original median | Native median |
| --- | ---: | ---: |
| Serial, one CPU | 618.5 ns/op | 293.6 ns/op |
| Serial, two CPUs available | 343.0 ns/op | 318.1 ns/op |
| Batch 256, one CPU | 791.6 ns/op | 418.6 ns/op |
| Batch 256, two CPUs available | 604.7 ns/op | 422.1 ns/op |
| Parallel, one CPU | 613.4 ns/op | 302.1 ns/op |
| Parallel, two CPUs | 257.6 ns/op | 336.1 ns/op |

**Two-worker raw churn remains about 30% slower.** Native samples ranged from
207.8 to 350.2 ns/op; a fresh isolated original probe measured 220.5–264.3 ns/op.
Treat the residual overhead as real. Accept this tradeoff because actual affected
owner medians are within about 1.3% of the original, memory use is comparable and five other benchmark
medians improve. Revisit if a real concurrent game workload shows allocator cost.
These samples do not establish universal hardware/gameplay speed.

Owner timing was repeated with the corrected fixture on both backends. Allocator-only
benchmark samples were reused after verifying that the complete allocation package,
including its benchmarks, was byte-identical; only legacy porttest fixture source
changed. Their provenance is explicit in the qualification record.

## Rejected attempts and process lessons

The first integrated backend passed functional checks but slowed actual owners
17–30% and parallel churn about fourfold. Profiling exposed clearing/mapping costs
and the global lock. Per-class locks and first-use zero handling alone did not
recover owner performance: the painting fixture repeatedly needs about 1 MiB of
row buffers in one class, while the initial cache retained only 128 KiB there.
Expanding the bounded cache recovered owner timings.

Parallel profiling then showed clearing reused buffers dominated CPU time, with
little lock wait. OS thread-ID hints added too much overhead. Go reuse hints or
the byte-clear leaf alone were insufficiently consistent; the selected combination
improved both serial throughput and concurrent reuse. Earlier source/results are
segregated, not accepted as final evidence. No frozen expectations were changed to
hide differences. Measure real working sets before repeating broad qualification;
isolated allocation probes did not predict all integrated costs.

Primary handled implementation and review because Luna quota was unavailable;
no replacement model was used.

## Local evidence and recovery

Working evidence: `build/port-final-engine-boundary/`. Final normal binaries are
in `contracts/profiles`, complete corpus logs in `full`, safe/production artifacts
in `safe` and `production/production/bin`. Accepted owner timing uses `paired-owner-performance` and the current normal
binaries; allocator-only samples are retained in `pre-slot-fix/early-performance-results`.
The combined record is in `early-performance-results`. Original reference
binaries remain in `original/profiles`; their source recovers from `71ff3095`.
The corrected-fixture original binary and source are in `fixture-slot-original`;
these were used for the final owner timing comparisons.
Rejected integrated attempts live in `rejected-v1`, `rejected-v2`, `rejected-v3`
and `rejected-v7`, with relocation manifests preserving historical recorded paths.
`pre-slot-fix` retains the earlier native gates and the stopped full sweep.
Other isolated probes remain under `revision*`; none substitutes for final gates.
Completed scripts are single-use.

Verified host process/open-file checks and archive hashes allowed two removals
of old rebuildable Linux 386 Go cache archives: 92 files, totaling 4,564,062,208
allocated bytes. Source, assets, logs and retained binaries were preserved. Recovery is automatic cache
rebuilding. See [cleanup record](final-engine-boundary-cache-cleanup.json).
