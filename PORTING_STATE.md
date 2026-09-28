# Porting checkpoint

This is the current resume checkpoint. Workflow and delegation rules live in
[PORT.md](PORT.md); detailed evidence belongs in the linked batch reports.

## Status: resumed; internal C-glue removal

- **Standalone engine C:** 142,665/142,665 original lines ported or retired.
- **Production C imports:** 459/463 client cgo files eliminated on net (4 remain);
  server: 458/463 eliminated (5 remain). Two project packages directly use cgo.
- **Legacy C export bridges:** 1,890/1,890 retired (0 remain).
- **Embedded production C bodies:** 79/79 retired (0 remain).
- **Fixture and declaration dependencies:** 11 porttest-tagged files still import C;
  157 tracked headers / 2,731 physical lines remain.

Production cgo counts cover selected Linux 386 builds; these metrics are not equal
units of effort. Fixture imports cover all build tags. Standalone production/test C
remain zero.

Latest qualified chunk routes fixture string/byte copies through existing native
helpers in three files, retiring three fixture C imports. Byte contents, allocation
ownership and failure behavior are preserved in normal and safe builds. See
[FIXTURE_STRING_COPIES.md](docs/porting/FIXTURE_STRING_COPIES.md).

Continue chunk-by-chunk with primary review, qualification, documentation and
commit/push. Use one Luna helper when its quota is available; no substitute model.
Stop at the milestone, usage limits or a substantial question. Latest qualified
artifacts: `build/port-fixture-string-copies/`.

## What remains

| Area | Remaining work or dependency |
| --- | --- |
| Embedded C callback glue | Complete: zero production bodies remain in selected builds; pure-Go callback registries remain. |
| Types and declarations | 157 tracked headers / 2,731 physical lines. 11 porttest-tagged source files still import C across all build tags (not one selected profile). Fixture C observers/types still need retirement or explicit qualification scope. |
| Production C imports | alloc/raw.go, legacy/object_xfer_world.go (reachable abort), legacy/cgo_common.go and legacy/video_highres.go (flags). Server adds legacy/cgo_server.go. |
| Memory and layout | C-heap allocator, raw pointers, fixed offsets and 32-bit address assumptions remain. Preserve ownership/lifetime/failure semantics behind the centralized allocator. |
| External libraries | SDL2, OpenGL, OpenAL and similar bindings remain for this phase; later replacement requires a separate discussion. |
| Portability | Qualified target remains Linux 386/SSE2 with cgo. Whole-build cgo-free, 64-bit and other platforms are not qualified. Physical display/audible playback remain manual checks. |

## Latest qualification

Original and converted selections passed all six affected roots twice per profile,
plus six safe contracts. Two private string/allocation contracts passed in each of
default/server/highres/safe. Every target compiled all root tests. Static checks,
safe build, production builds/ABI, exact known-suite outcomes and fresh save/load
passed. Root assertions, captures and original asset hashes remain unchanged.

Evidence: [qualification](docs/porting/fixture-string-copies-qualification.json),
[inventory](docs/porting/fixture-string-copies-inventory-after.json).
Known-suite expectation: [record](docs/porting/internal-callback-expected-suite.jsonl).
Standalone metric/history: [C_LOC.md](docs/porting/C_LOC.md).
Latest broad regression: default in fixture-import-cleanup (2,482); all three
profiles in final-callback-dispatch (2,482/2,471/2,482), no skips.
Last complete default corpus: `6e9681f2`, 2,489 passes plus diagnostic skip.

## Next work and review items

Next: install and qualify the native tile-grid fixture helpers after committing
the accepted original baseline. The new allocation contract passed normal/safe
preflight against C, checking all 129 failure positions plus two success cases.
The 564-owner baseline, nine focused repeats per profile and six safe contracts
passed without skips. See [FIXTURE_TILE_GRID.md](docs/porting/FIXTURE_TILE_GRID.md).
Then retire numeric-state helpers and C allocation observers, reachable abort/
compiler flags and the centralized internal allocator. Preserve normal/safe
allocation domains, layout and failure behavior; external native libraries remain.
Only the original-path fixture extension is installed, not the conversion.
Baseline completion reclaimed 4.42 GB of verified, reproducible build-cache
archives; see the batch report for checks and recovery.

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
`build/port-fixture-string-copies/{contracts/profiles,safe,production/production/bin}`.
Original owners/focused/safe/private captures live under
`build/port-fixture-string-copies/{original,original-safe,original-legacy}`; root
binaries come from the preceding fixture-constants batch. Exact paths/hashes are
recorded in the baseline. Earlier batch reports record disk cleanup and scenario
restoration commands; original assets remain intact. Current cleanup reclaimed
about 1.47 GB from scenario copies and seventeen superseded binaries. Rebuild
6048add7 or d2df2eb6 for those older executables; see
[recovery details](docs/porting/FIXTURE_STRING_COPIES.md#artifact-cleanup-and-recovery).
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
