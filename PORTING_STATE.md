# Porting checkpoint

This file holds the current checkpoint. Durable workflow and delegation rules live
in [PORT.md](PORT.md); detailed evidence belongs in the linked batch reports.

## Status: Linux engine C-glue milestone qualified; paused for next-phase discussion

- **Standalone engine C:** 142,665/142,665 original lines ported or retired.
- **Production C imports:** 463/463 original client and server files eliminated on
  net; zero remain in selected Linux builds, across zero direct project cgo packages.
- **Legacy C export bridges:** 1,890/1,890 normal and 11/11 safe-only bridges retired.
- **Embedded production C bodies:** 79/79 retired; zero remain.
- **Fixture C dependencies:** zero C imports or fixture-only export callbacks remain.
- **Project C headers:** zero remain; the last 157 headers / 2,731 lines were retired
  in the preceding chunk.

Counts cover the qualified Linux 386/SSE2 target and exclude external dependencies.
They measure different areas, not equal units of effort. The Windows-only WinSock
binding still imports C outside this target. Standalone production/test C remain
zero; the complete application still uses cgo for external native libraries.

The final chunk replaces libc allocation and fatal-call glue with Go-owned OS
mappings and native termination handling, then removes engine compiler/linker
preambles. SDL2, OpenGL, OpenAL and their bindings remain unchanged. See
[FINAL_ENGINE_BOUNDARY.md](docs/porting/FINAL_ENGINE_BOUNDARY.md).

## Latest qualification

The final source passed 166 owners and sixteen focused repeats per normal profile,
twelve safe roots, thirteen allocator and six legacy fixture/string/clock roots in both normal
and safe modes, and repeated original-behavior allocation/fatal contracts.
The full native root corpus passed 2492/2481/2492 roots in default/server/highres,
with only the expected opt-in population diagnostic skipped in each profile.

Static checks, assembly validation, safe and three production builds, ABI/export
checks, exact known-suite outcomes, fresh save/load and all 1,654 original asset
hashes passed. The allocator's thirteen contracts also pass with cgo disabled;
Windows 386 allocator support is compile-checked only. Frozen assertions/captures
remain unchanged. Physical display and audible playback remain manual release checks.

Evidence: [qualification](docs/porting/final-engine-boundary-qualification.json),
[inventory](docs/porting/final-engine-boundary-inventory-after.json).
Known-suite expectation: [record](docs/porting/internal-callback-expected-suite.jsonl).
Standalone metric/history: [C_LOC.md](docs/porting/C_LOC.md).

## Review items and next phase

The agreed immediate milestone is complete. Pause here to discuss whether and how
to remove cgo from external native bindings; replacing rendering, audio and window
backends is a separate product/architecture decision. No 64-bit or broader platform
support is implied. Existing raw pointers, fixed layouts and 32-bit assumptions remain.

The allocator retains at most 18 MiB of idle payload mappings; this excludes live
fragmentation and Go metadata. Alternating real-owner timing medians are
within about 1.3% of the original, with comparable RSS. Synthetic two-worker allocation churn remains
about 30% slower than the frozen libc baseline; five other benchmark medians improve.
This is an explicit tradeoff to revisit if real concurrent workloads show a problem.

Safe Go allocation tracking and mapped-memory checks remain, but engine ASAN flags
are retired with its C glue. These are not equivalent whole-process ASAN coverage.
RawMalloc remains process-fatal on failure, distinct from safe recoverable panic.
See [DECISIONS.md](docs/porting/DECISIONS.md) for these and earlier behavior decisions.

Continue using at most one Luna helper for suitable bounded work when its quota
is available, with primary review and qualification. This chunk ran locally because
Luna quota was unavailable; no substitute model was used.

## Resume and artifact recovery

Read this checkpoint and PORT.md; inspect Git status before editing. Never stage
or delete `nox-iso-from-archive-org.7z` or `build/assets/extracted/drive_c/Nox`.
Source `build/baseline/env.sh` in every Go shell and put `/usr/lib/go-1.26/bin`
first on PATH. Linux 386 execution needs host execution in this VM. Do not change
source consumed by running builds/tests.

Current evidence lives under `build/port-final-engine-boundary/`. Final root,
safe and production binaries are in `contracts/profiles`, `safe` and
`production/production/bin`; full-corpus logs are in `full`. Accepted performance
uses current normal binaries and logs in `paired-owner-performance`; allocator-only
benchmarks are explicitly reused from `pre-slot-fix/early-performance-results` on
byte-identical allocation-package source. The combined record is in
`early-performance-results`.

The original libc/fatal baseline is committed as `71ff3095`; the additional slot
fixture correction is reproducible from its committed patch and original-backend
proof. Native implementation recovery checkpoint: `33e52f7c`. Retained original
binaries live under `original/profiles`; the corrected-fixture comparison binary
is in `fixture-slot-original`. Baseline production source is identical to
qualified `49c8ef62`. Rejected allocator attempts are explicitly segregated and
are not qualification evidence. Original headers recover from `333d4590`.
Earlier reports retain their artifact recovery and verified-cleanup records.

Cleanup removed 4.56 GB of old reproducible cache archives during this chunk;
source, assets, logs and retained binaries were preserved. Caches rebuild
automatically. Check physical disk headroom before large runs. Git backs up source,
expectations and reports, not ignored build artifacts or original assets.
Completed installation/qualification/cleanup scripts are single-use; never rerun
one merely because it is still present. Match source fingerprints when reusing
binaries rather than relying on their recorded Git HEAD alone.
