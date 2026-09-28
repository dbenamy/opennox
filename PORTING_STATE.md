# Porting checkpoint

This is the current resume checkpoint. Workflow and delegation rules live in
[PORT.md](PORT.md); detailed evidence belongs in the linked batch reports.

## Status: resumed; internal C-glue removal

- **Standalone engine C:** 142,665/142,665 original lines ported or retired.
- **Production C imports:** 459/463 client cgo files eliminated on net (4 remain);
  server: 458/463 eliminated (5 remain). Two project packages directly use cgo.
- **Legacy C export bridges:** 1,890/1,890 retired in normal profiles;
  11/11 safe-only allocator/memory bridges retired.
- **Embedded production C bodies:** 79/79 retired (0 remain).
- **Fixture C dependencies:** 2 porttest-tagged files still import C;
  no fixture-only export callbacks remain.
- **Project C headers:** 0 remain; the last 157 headers / 2,731 lines are retired.

Production cgo counts cover selected Linux 386 builds; these metrics are not equal
units of effort. Fixture imports cover all build tags. Standalone production/test C
remain zero.

Latest qualified chunk replaces C allocation observers in four fixture files with
thread-scoped native test dispatch and retires the last 157 unused project headers.
A new original/native contract checks actual other-thread allocation isolation. Production dispatch
stays direct; aligned buffers and the allocator remain for the next phase work.
See [FIXTURE_ALLOCATION_OBSERVERS.md](docs/porting/FIXTURE_ALLOCATION_OBSERVERS.md).

Continue chunk-by-chunk with primary review, qualification, documentation and
commit/push. Use one Luna helper when its quota is available; no substitute model.
Stop at the milestone, usage limits or a substantial question. Latest qualified
artifacts: `build/port-fixture-allocation-observers/`.

## What remains

| Area | Remaining work or dependency |
| --- | --- |
| Embedded C callback glue | Complete: zero production bodies remain in selected builds; pure-Go callback registries remain. |
| Types and declarations | No project C headers remain. Two porttest files still import C for aligned allocation/abort; no fixture C exports remain. |
| Production C imports | alloc/raw.go, legacy/object_xfer_world.go (reachable abort), legacy/cgo_common.go and legacy/video_highres.go (flags). Server adds legacy/cgo_server.go; safe adds cgo_safe.go (sanitizer flags). |
| Memory and layout | C-heap allocator, raw pointers, fixed offsets and 32-bit address assumptions remain. Preserve ownership/lifetime/failure semantics behind the centralized allocator. |
| External libraries | SDL2, OpenGL, OpenAL and similar bindings remain for this phase; later replacement requires a separate discussion. |
| Portability | Qualified target remains Linux 386/SSE2 with cgo. Whole-build cgo-free, 64-bit and other platforms are not qualified. Physical display/audible playback remain manual checks. |

## Latest qualification

Original/native selections passed 80 owners and eight repeats per normal profile,
eleven safe roots and direct allocator/string/clock contracts. The complete native
root corpus passed 2491/2480/2491 roots (default/server/highres), with one expected
opt-in diagnostic skip per profile. Static checks, safe build, production builds/ABI,
exact known-suite outcomes and fresh save/load passed. Frozen assertions/captures
and original asset hashes remain unchanged.

Evidence: [qualification](docs/porting/fixture-allocation-observers-qualification.json),
[inventory](docs/porting/fixture-allocation-observers-inventory-after.json).
Known-suite expectation: [record](docs/porting/internal-callback-expected-suite.jsonl).
Standalone metric/history: [C_LOC.md](docs/porting/C_LOC.md).
This is the latest complete regression across all three profiles; the skipped
population diagnostic's four cases have explicit regression coverage.

## Next work and review items

Next: retire the two aligned-buffer fixture imports, reachable abort/compiler flags
and centralized libc allocator implementation. Preserve normal/safe allocation
domains, alignment, layout and failure behavior, including safe Calloc's nil marker.
External native libraries remain preserved dependencies. Original-path contracts are
frozen under `build/port-final-engine-boundary/`; the committed original
baseline uses libc. Five allocation roots passed twice in normal/safe, and the fatal
contract passed twice in all four profiles after a fixture setup correction.
The original baseline is pushed as `71ff3095`. The native allocator/abort/flag
replacement passed the first functional gates, but its real-owner and parallel
performance regressed. The second revision also passed functional gates but still regressed performance.
A larger bounded reuse cache fixed most owner overhead. The current candidate
also uses a measured 386 byte-clear leaf, with independent guard-byte checks;
performance is accepted with an explicit parallel-churn tradeoff. The full sweep then found four hallway/prefab capture mismatches in both
profiles and was stopped. All 2,512 hallway cases differ only in saved-slot identity after address reuse.
The correction now passes the original backend with unchanged captures in all
three profiles, and its forced-reuse test passes normal/safe. It is installed for
fresh native qualification. All 166 owners and 16 repeats now pass per normal
profile, along with initial safe/library/leaf checks. Corrected-fixture performance
is at parity within about 1.3%; the complete corpus and production gates remain.
status counts above
still describe qualified `49c8ef62`. See
[FINAL_ENGINE_BOUNDARY.md](docs/porting/FINAL_ENGINE_BOUNDARY.md).

The [immediate goal](PORT.md#goal-and-target) remains internal engine C-glue
removal. Follow [INTERNAL_C_GLUE.md](docs/porting/INTERNAL_C_GLUE.md).
External SDL2/OpenGL/OpenAL backends remain out of scope for this phase.

Audited callback families now use native registration maps and reject unknown
keys with explicit panics. Fixture registrations preserve parser names, native
priority, shared callback routes and capture IDs. Bounds panics outside valid
observer capacity are recorded in the corresponding reports.

CString retains raw malloc/free normally and tracked Malloc/FreePtr in safe.
The qualified baseline RawMalloc uses cgo's process-fatal malloc wrapper; the
pending native replacement preserves its fatal disposition. Keep normal
fatal failure distinct from safe recoverable panic; both have qualified contracts.
Placement math preserves libc rounding only in its qualified argument domain;
requalify before expanding callers. Layout aliases preserve allocation sizes;
packed Player alignment and Object/Drawable extensions remain separately defined.
See [DECISIONS.md](docs/porting/DECISIONS.md) and linked batch reports for details.

## Resume and artifact recovery

Read this checkpoint and PORT.md; inspect Git status before editing. Never stage
or delete `nox-iso-from-archive-org.7z` or `build/assets/extracted/drive_c/Nox`.
Source `build/baseline/env.sh` in every Go shell and put `/usr/lib/go-1.26/bin`
first on PATH. Linux 386 execution needs host execution in this VM.
Do not change source consumed by running builds/tests.

Current root/safe/production binaries live under
`build/port-fixture-allocation-observers/{contracts/profiles,safe,production/production/bin}`.
The complete corpus reuses those exact-source root binaries; logs live in `full`.
Original observer binaries and safe/library evidence live under `original/profiles`,
`original-safe` and `original-library`; the new contract's repeated original checks
are in `baseline-preflight`. Exact source and binary hashes are in the baseline.
Retired headers and C observers recover from baseline `333d4590`. Original assets
remain intact; earlier reports record artifact recovery and verified cleanup.
Rebuild 6048add7 or d2df2eb6 for retired type/allocation binaries. Numeric/safe-bridge
native binaries were superseded and removed (774 MB); rebuild dbc4e0be or 8580330d.
Original numeric C binaries and current observer references remain; batch reports
retain source/binary fingerprints. Verified scenario deduplication reclaimed another
560 MB and reproducible cache cleanup 2.11 GB; restoration is recorded in the batch
report. Physical free space after cleanup was 4.40 GB.
Superseded GUI binaries were removed; rebuild `0cf5064c` if needed. Cache archives
rebuild automatically. Details are in [FINAL_CALLBACK_DISPATCH.md](docs/porting/FINAL_CALLBACK_DISPATCH.md).

Check physical free space before large runs. Source, reports and qualification
metadata are committed; ignored build artifacts are not backed up by Git. Current
binaries and useful original references remain locally. Completed scripts are
consumed: never rerun a completed installation, qualification or cleanup script.
Historical removals and lossless log-archive recovery are documented in the linked
batch reports and their ignored manifests. Rebuild superseded binaries from the
recorded source revision; cache archives rebuild automatically. Do not infer
safe deletion merely from age. Preserve original assets and the archive.

Prior GUI cleanup/recovery is in [WINDOW_DISPATCH.md](docs/porting/WINDOW_DISPATCH.md).
