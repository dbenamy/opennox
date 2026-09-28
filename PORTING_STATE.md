# Porting checkpoint

This is the current resume checkpoint. Workflow and delegation rules live in
[PORT.md](PORT.md); detailed evidence belongs in the linked batch reports.

## Status: resumed; internal C-glue removal

- **Standalone engine C:** 142,665/142,665 original lines ported or retired.
- **Production C imports:** 459/463 client cgo files eliminated on net (4 remain);
  server: 458/463 eliminated (5 remain). Two project packages directly use cgo.
- **Legacy C export bridges:** 1,890/1,890 retired (0 remain).
- **Embedded production C bodies:** 79/79 retired (0 remain).
- **Fixture and declaration dependencies:** 69 porttest-tagged files still import C;
  157 tracked headers / 2,731 physical lines remain.

Production cgo counts cover selected Linux 386 builds; these metrics are not equal
units of effort. Fixture imports cover all build tags. Standalone production/test C
remain zero.

Latest qualified chunk retires the final eight generic C callback dispatch bodies
and migrates their shared fixture observers. Native callback behavior, assertions
and captures are preserved. See
[FINAL_CALLBACK_DISPATCH.md](docs/porting/FINAL_CALLBACK_DISPATCH.md).

Continue chunk-by-chunk with primary review, qualification, documentation and
commit/push. Use one Luna helper when its quota is available; no substitute model.
Stop at the milestone, usage limits or a substantial question. Latest qualified
artifacts: `build/port-final-callback-dispatch/`.

## What remains

| Area | Remaining work or dependency |
| --- | --- |
| Embedded C callback glue | Complete: zero production bodies remain in selected builds; pure-Go callback registries remain. |
| Types and declarations | 157 tracked headers / 2,731 physical lines. 69 porttest-tagged source files still import C across all build tags (not one selected profile). Fixture C observers/types still need retirement or explicit qualification scope. |
| Production C imports | alloc/raw.go, legacy/object_xfer_world.go (reachable abort), legacy/cgo_common.go and legacy/video_highres.go (flags). Server adds legacy/cgo_server.go. |
| Memory and layout | C-heap allocator, raw pointers, fixed offsets and 32-bit address assumptions remain. Preserve ownership/lifetime/failure semantics behind the centralized allocator. |
| External libraries | SDL2, OpenGL, OpenAL and similar bindings remain for this phase; later replacement requires a separate discussion. |
| Portability | Qualified target remains Linux 386/SSE2 with cgo. Whole-build cgo-free, 64-bit and other platforms are not qualified. Physical display/audible playback remain manual checks. |

## Latest qualification

Broad native selections passed 2,482/2,471/2,482 roots without skips. Exact-source
qualified damage/monster broad runs supplied the original baseline; 28 focused
roots per profile, three safe contracts and private legacy adapters in four
profiles ran freshly before/after conversion. Native focused/static and fresh
default save/load passed before the broad sweep. Safe build/static, three
production builds/ABI, exact known-suite outcomes and final fresh save/load passed.
Assertions, native handlers and original assets are unchanged.

Evidence: [qualification](docs/porting/final-callback-dispatch-qualification.json),
[inventory](docs/porting/final-callback-dispatch-inventory-after.json).
Known-suite expectation: [record](docs/porting/internal-callback-expected-suite.jsonl).
Standalone metric/history: [C_LOC.md](docs/porting/C_LOC.md).
Latest broad regression: this remaining-callback batch, 2,482/2,471/2,482 audited roots, no skips.
Last complete default corpus: `6e9681f2`, 2,489 passes plus diagnostic skip.

## Next work and review items

Next: review and qualify the 18-file fixture import cleanup drafted under
`build/port-fixture-import-cleanup/` (15 unused header imports and three scalar
bridges). Its original baseline is accepted (verified prior broad results, eight fresh focused
roots per profile and two safe contracts); it is not installed. See
[FIXTURE_IMPORT_CLEANUP.md](docs/porting/FIXTURE_IMPORT_CLEANUP.md). Then
audit remaining fixture types, allocator ownership, reachable abort and compiler
flags. Production still uses the centralized C heap; fixture allocation and
safe-mode accounting need explicit ownership matching before migration.

The [immediate goal](PORT.md#goal-and-target) remains internal engine C-glue
removal. Follow [INTERNAL_C_GLUE.md](docs/porting/INTERNAL_C_GLUE.md).
External SDL2/OpenGL/OpenAL backends remain out of scope for this phase.

Audited callback families now use native registration maps and reject unknown
keys with explicit panics. Fixture registrations preserve parser names, native
priority, shared callback routes and capture IDs. Bounds panics outside valid
observer capacity are recorded in the corresponding reports.

CString retains raw malloc/free normally and tracked Malloc/FreePtr in safe.
RawMalloc intentionally still uses cgo's process-fatal malloc wrapper. Keep normal
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
`build/port-final-callback-dispatch/{contracts/profiles,safe,production/production/bin}`.
Original focused/safe/private-legacy captures live under
`build/port-final-callback-dispatch/{original,original-safe,original-legacy}`.
Original broad logs/binaries are reused from the preceding damage/monster batch.
Exact paths/hashes are in the committed baseline. Both final-callback scenario
asset copies were verified and deduplicated (1,120,485,376 bytes reclaimed); the
batch report has restoration commands. Original assets and run outputs remain.
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
