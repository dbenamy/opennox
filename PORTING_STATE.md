# Porting checkpoint

This is the current resume checkpoint. Workflow and delegation rules live in
[PORT.md](PORT.md); detailed evidence belongs in the linked batch reports.

## Status: resumed; internal C-glue removal

**Progress: 142,665/142,665 original standalone C lines ported or retired;
internal glue: 458/463 client cgo files eliminated on net (5 remain;
server: 457/463 eliminated, 6 remain).**
Selected legacy C export bridges: **1,890/1,890 retired (0 remain)**.
Embedded production C bodies: **68/79 retired (11 remain)**.

These are selected Linux 386 production files, not equal units of effort.
Three project packages directly use cgo. Standalone production/test C remain zero.
Latest qualified chunk removes duration-spell and internal audio-stream foreign
fallbacks and converts both observer fixtures. Two fixture C imports and four C
call sites are retired; generic dispatch bodies remain shared with other families.
All 118 affected owner roots passed in each profile. See
[DURATION_AUDIO_DISPATCH.md](docs/porting/DURATION_AUDIO_DISPATCH.md).

Continue chunk-by-chunk with primary review, qualification, documentation and
commit/push. Use one Luna helper when its quota is available; no substitute model.
Stop at the milestone, usage limits or a substantial question. Latest qualified
artifacts: `build/port-single-pointer-dispatch/`.

## What remains

| Area | Remaining work or dependency |
| --- | --- |
| Embedded C callback glue | 11 generic dispatch bodies. Trace all owners/registrations before removing raw fallbacks. |
| Types and declarations | 157 tracked headers / 2,731 physical lines. 89 porttest-tagged source files still import C across all build tags (not one selected profile). Fixture C observers/types still need retirement or explicit qualification scope. |
| Production C imports | alloc/raw.go, ccall/ccall.go, legacy/object_xfer_world.go (reachable abort), legacy/cgo_common.go and legacy/video_highres.go (flags). Server adds legacy/cgo_server.go. |
| Memory and layout | C-heap allocator, raw pointers, fixed offsets and 32-bit address assumptions remain. Preserve ownership/lifetime/failure semantics behind the centralized allocator. |
| External libraries | SDL2, OpenGL, OpenAL and similar bindings remain for this phase; later replacement requires a separate discussion. |
| Portability | Qualified target remains Linux 386/SSE2 with cgo. Whole-build cgo-free, 64-bit and other platforms are not qualified. Physical display/audible playback remain manual checks. |

## Latest qualification

Original and native owner selections passed 118 roots without skips in each of
three profiles; six focused contracts ran separately in each profile on both
versions. Three safe contracts passed on both versions. Safe build/static, three
production builds/ABI, exact known-suite outcomes and fresh save/load passed.
Retired fixture symbols are absent; native handlers, remaining C bodies,
assertions and assets unchanged.

Evidence: [qualification](docs/porting/single-pointer-dispatch-qualification.json),
[inventory](docs/porting/single-pointer-dispatch-inventory-after.json).
Known-suite expectation: [record](docs/porting/mp3-go-expected-suite.jsonl).
Standalone metric/history: [C_LOC.md](docs/porting/C_LOC.md).
Latest broad regression: `5ba6fce6`, 2,482/2,471/2,482 audited roots, no skips.
Last complete default corpus: `6e9681f2`, 2,489 passes plus diagnostic skip.

## Next work and review items

Active batch: GUI window dispatch and tooltip observation. Original baseline
accepted at `d3d759ba`: 2,482/2,471/2,482 broad roots, six focused repeats per
profile, safe deferred cleanup. Reviewed seven-file draft not installed yet.
Next: install, focused/default-client scenario preflight, then complete qualification.
See [WINDOW_DISPATCH.md](docs/porting/WINDOW_DISPATCH.md). Remaining object/AI
observer connections are inventoried in `build/port-object-dispatch-audit/`.

The [immediate goal](PORT.md#goal-and-target) remains internal engine C-glue
removal. Follow [INTERNAL_C_GLUE.md](docs/porting/INTERNAL_C_GLUE.md).
Metadata discovery is not compilation evidence.

Continue callback removal by tracing all producers, registrations and fixture
observers. Remaining families include GUI/drawable, object/AI and smaller legacy callbacks. Audit by owner before rejecting foreign
keys; migrate shared observers and preserve existing captures. Also remaining:
allocator ownership, reachable abort, compiler flags and fixture dependencies.
Literal allocator calls occur in 40 C-using fixture files; avoid cross-domain frees.

Animation, quantity, image completion, modifier, player-section, screen-particle,
duration-spell and internal audio dispatch reject unregistered
keys with an explicit panic after producer audits. These reversible corrections
apply only to those families. Native fixture arrays also make invalid overflow
explicit; valid-domain captures remain unchanged.

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
`build/port-single-pointer-dispatch/{contracts/profiles,safe,production/production/bin}`.
Original duration/audio captures live under
`build/port-single-pointer-dispatch/original`, with safe contracts in
`original-safe`. Baseline runs use retained section/particle binaries at
`6f03cc39`. Completed duration/audio scenario duplicates were verified and removed
(560,214,016 bytes). Restore before replay: `python3
build/port-artifact-cleanup/restore-recent-scenario.py
build/baseline/runs/single-pointer-dispatch-save/deduplicated-assets.json`.
Older cleanup records retain source revisions, hashes and rebuild information.

Six completed string/spell regression logs were losslessly archived, reclaiming
613,769,216 bytes. Restore with `gzip -dk FILE.jsonl.gz`; exact hashes/paths and
commands: `build/port-callback-retirement-audit/log-archive.json`.
Older placement and string-original log archive records remain in
`build/port-string-boundaries/log-archive.json`.

Final modifier scenario duplicate assets were verified and removed; originals,
saves/results remain. Restore before replay:
`python3 build/port-artifact-cleanup/restore-recent-scenario.py build/baseline/runs/modifier-dispatch-save/deduplicated-assets.json`.

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
inactive Linux 386 cache archives untouched for six hours (496,467,968 bytes).
Records: `build/port-animation-dispatch/{spell-cleanup-*,cache-headroom-*}`.

During modifier preparation, removed seven superseded callback binaries
(387,018,752 bytes; rebuild `dcab2e38`) and seven animation binaries
(387,063,808 bytes; rebuild `5233d084`) after source/replacement/hash/host-use checks.
Logs and records remain; manifests: `build/port-modifier-dispatch/*-cleanup-*`.

Section/particle preparation also reclaimed 1,580,584,960 bytes from nine verified
superseded binaries and twelve losslessly gzipped historical logs. Current modifier
artifacts remain. Exact paths/hashes and restore commands are in
`build/port-section-particle-dispatch/{old-binaries-result,log-archive}.json`;
see the batch report for source revisions and recovery.

Duration/audio preparation removed seven superseded UI-completion binaries
(387,067,904 bytes; rebuild `defbdb98`) after source/replacement/hash/host-use
checks. Logs/records remain: `build/port-single-pointer-dispatch/ui-cleanup-*`.

GUI preparation removed seven superseded modifier binaries (387,018,752 bytes;
rebuild `5ba6fce6`) and 39 inactive Linux 386 cache archives (1,941,274,624 bytes)
after source/replacement/host-use or cache metadata/hash checks. Records and
recovery: [WINDOW_DISPATCH.md](docs/porting/WINDOW_DISPATCH.md).
