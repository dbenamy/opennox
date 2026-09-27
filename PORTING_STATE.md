# Porting checkpoint

This is the current resume checkpoint. Workflow and delegation rules live in
[PORT.md](PORT.md); detailed evidence belongs in the linked batch reports.

## Status: resumed; internal C-glue removal

**Progress: 142,665/142,665 original standalone C lines ported or retired;
internal glue: 458/463 client cgo files eliminated on net (5 remain;
server: 457/463 eliminated, 6 remain).**
Selected legacy C export bridges: **1,890/1,890 retired (0 remain)**.
Embedded production C bodies: **64/79 retired (15 remain)**.

These are selected Linux 386 production files, not equal units of effort.
Three project packages directly use cgo. Standalone production/test C remain zero.
Latest qualified chunk routes animation completion through the native registry,
retires two generic C dispatch signatures and replaces two fixture C callbacks.
All 69 selected roots passed in each production profile; the completion contract
passed twice in four profiles. See
[ANIMATION_DISPATCH.md](docs/porting/ANIMATION_DISPATCH.md).

Continue chunk-by-chunk with primary review, qualification, documentation and
commit/push. Use one Luna helper when its quota is available; no substitute model.
Stop at the milestone, usage limits or a substantial question. Latest qualified
artifacts: `build/port-animation-dispatch/`.

## What remains

| Area | Remaining work or dependency |
| --- | --- |
| Embedded C callback glue | 15 generic dispatch bodies. Trace all owners/registrations before removing raw fallbacks. |
| Types and declarations | 157 tracked headers / 2,731 physical lines. 100 porttest-tagged source files still import C across all build tags (not one selected profile). Fixture C observers/types still need retirement or explicit qualification scope. |
| Production C imports | alloc/raw.go, ccall/ccall.go, legacy/object_xfer_world.go (reachable abort), legacy/cgo_common.go and legacy/video_highres.go (flags). Server adds legacy/cgo_server.go. |
| Memory and layout | C-heap allocator, raw pointers, fixed offsets and 32-bit address assumptions remain. Preserve ownership/lifetime/failure semantics behind the centralized allocator. |
| External libraries | SDL2, OpenGL, OpenAL and similar bindings remain for this phase; later replacement requires a separate discussion. |
| Portability | Qualified target remains Linux 386/SSE2 with cgo. Whole-build cgo-free, 64-bit and other platforms are not qualified. Physical display/audible playback remain manual checks. |

## Latest qualification

Baseline `6a30a30d` captured 69 original owner roots in each profile and repeated
the new completion contract twice in default/server/highres/safe. The matching
native runs passed without skips. Focused preflight/static, safe build, three
production builds/ABI, exact known-suite outcomes and fresh headless save/load
also passed. Retired dispatcher and fixture-export symbols are absent. Retained
C bodies and all 25 production callback mappings are unchanged. All 1,654 original
asset hashes are unchanged; no frozen expectations changed.

Evidence: [qualification](docs/porting/animation-dispatch-qualification.json),
[inventory](docs/porting/animation-dispatch-inventory-after.json).
Known-suite expectation: [record](docs/porting/mp3-go-expected-suite.jsonl).
Standalone metric/history: [C_LOC.md](docs/porting/C_LOC.md).
Latest broad regression: `b2c975dc`, 2,482/2,471/2,482 audited roots, no skips.
Last complete default corpus: `6e9681f2`, 2,489 passes plus diagnostic skip.

## Next work and review items

The [immediate goal](PORT.md#goal-and-target) remains internal engine C-glue
removal. Follow [INTERNAL_C_GLUE.md](docs/porting/INTERNAL_C_GLUE.md).
Metadata discovery is not compilation evidence.

In progress: quantity-dialog and image/book completion dispatch. Producer, storage
and handler review is complete; the seven-file draft is not installed. Original
133/132/133-root selections and six focused repeat contracts passed on exact-source
qualified parent binaries. The baseline is accepted; native qualification is pending. See [UI_COMPLETION.md](docs/porting/UI_COMPLETION.md).
Then continue allocator ownership, reachable abort, compiler flags and fixture
work. Literal allocator calls occur in 40 C-using fixture files; avoid incompatible
cross-domain frees when the centralized allocator changes.

Animation now rejects unregistered keys with an explicit panic. All 33 production
field writes resolve to nil or 25 registered identities; this reversible correction
does not authorize the same conclusion for other callback families.

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
`build/port-animation-dispatch/{contracts/profiles,safe,production/production/bin}`.
Original animation captures live under `build/port-animation-dispatch/`;
production baseline reuse is recorded in its committed baseline manifest.
Older executable cleanup records remain under the preceding callback batch;
rebuild the recorded revisions when those historical binaries are needed.

Six completed string/spell regression logs were losslessly archived, reclaiming
613,769,216 bytes. Restore with `gzip -dk FILE.jsonl.gz`; exact hashes/paths and
commands: `build/port-callback-retirement-audit/log-archive.json`.
Older placement and string-original log archive records remain in
`build/port-string-boundaries/log-archive.json`.

Final animation scenario duplicate assets were verified and removed; originals,
saves/results remain. Restore before replay:
`python3 build/port-artifact-cleanup/restore-recent-scenario.py build/baseline/runs/animation-dispatch-save/deduplicated-assets.json`.

Check disk before large runs. Source, reports and qualification metadata are
committed; ignored artifacts are not backed up by Git. Completed scripts are
consumed. Historical cleanup/recovery details live in batch reports, ignored
manifests and checkpoint history (`b2c975dc` / `32df9553`); historical current/next
instructions are not the active plan. Do not infer deletion safety from age.

After all qualification jobs joined, removed 31 inactive Linux386 Go
cache archives untouched for six hours (1,413,988,352 allocated bytes),
with path/stat/hash and host compiler/fd/maps checks. Caches rebuild normally;
module sources and current artifacts remain. Records:
`build/port-callback-retirement-audit/cache-headroom-*`.

After the animation commit, removed seven superseded spell binaries (387,080,192
bytes) with committed-source, replacement and host-use checks. Rebuild `b2c975dc`
if those historical binaries are needed; metadata/logs remain. Also removed 12
inactive Linux386 cache archives untouched for six hours (496,467,968 bytes).
Records: `build/port-animation-dispatch/{spell-cleanup-*,cache-headroom-*}`.
