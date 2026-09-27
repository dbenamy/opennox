# Porting checkpoint

This is the current resume checkpoint. Workflow and delegation rules live in
[PORT.md](PORT.md); detailed evidence belongs in the linked batch reports.

## Status: resumed; internal C-glue removal

**Progress: 142,665/142,665 original standalone C lines ported or retired;
internal glue: 458/463 client cgo files eliminated on net (5 remain;
server: 457/463 eliminated, 6 remain).**
Selected legacy C export bridges: **1,890/1,890 retired (0 remain)**.
Embedded production C bodies: **67/79 retired (12 remain)**.

These are selected Linux 386 production files, not equal units of effort.
Three project packages directly use cgo. Standalone production/test C remain zero.
Latest qualified chunk removes modifier dispatch's foreign fallbacks and translates
four shared observers (modifier forwarding, equipment, damage and melee). Two
C dispatch signatures and four fixture C imports are retired. All 2,482/2,471/2,482
audited roots passed in default/server/highres. See
[MODIFIER_DISPATCH.md](docs/porting/MODIFIER_DISPATCH.md).

Continue chunk-by-chunk with primary review, qualification, documentation and
commit/push. Use one Luna helper when its quota is available; no substitute model.
Stop at the milestone, usage limits or a substantial question. Latest qualified
artifacts: `build/port-modifier-dispatch/`.

## What remains

| Area | Remaining work or dependency |
| --- | --- |
| Embedded C callback glue | 12 generic dispatch bodies. Trace all owners/registrations before removing raw fallbacks. |
| Types and declarations | 157 tracked headers / 2,731 physical lines. 94 porttest-tagged source files still import C across all build tags (not one selected profile). Fixture C observers/types still need retirement or explicit qualification scope. |
| Production C imports | alloc/raw.go, ccall/ccall.go, legacy/object_xfer_world.go (reachable abort), legacy/cgo_common.go and legacy/video_highres.go (flags). Server adds legacy/cgo_server.go. |
| Memory and layout | C-heap allocator, raw pointers, fixed offsets and 32-bit address assumptions remain. Preserve ownership/lifetime/failure semantics behind the centralized allocator. |
| External libraries | SDL2, OpenGL, OpenAL and similar bindings remain for this phase; later replacement requires a separate discussion. |
| Portability | Qualified target remains Linux 386/SSE2 with cgo. Whole-build cgo-free, 64-bit and other platforms are not qualified. Physical display/audible playback remain manual checks. |

## Latest qualification

Original and native broad selections passed 2,482/2,471/2,482 roots without skips
in default/server/highres. Four focused contracts ran separately in each profile
on both versions; native preflight also covered damage/melee captures. Three safe
contracts passed on both versions. Safe build/static, three production builds/ABI,
exact known-suite outcomes and fresh save/load passed. Retired symbols are absent;
remaining C bodies, 40 native registrations, assertions and assets are unchanged.

Evidence: [qualification](docs/porting/modifier-dispatch-qualification.json),
[inventory](docs/porting/modifier-dispatch-inventory-after.json).
Known-suite expectation: [record](docs/porting/mp3-go-expected-suite.jsonl).
Standalone metric/history: [C_LOC.md](docs/porting/C_LOC.md).
Latest broad regression: this modifier batch, 2,482/2,471/2,482 audited roots, no skips.
Last complete default corpus: `6e9681f2`, 2,489 passes plus diagnostic skip.

## Next work and review items

Active batch: player-section/metadata and screen-particle dispatch. Original
baseline accepted: 139 owner roots and six focused contracts in each profile,
plus two safe contracts. Source remains at the qualified modifier conversion;
reviewed eight-file draft is not installed yet. Next: install and qualify. See
[SECTION_PARTICLE_DISPATCH.md](docs/porting/SECTION_PARTICLE_DISPATCH.md).

The [immediate goal](PORT.md#goal-and-target) remains internal engine C-glue
removal. Follow [INTERNAL_C_GLUE.md](docs/porting/INTERNAL_C_GLUE.md).
Metadata discovery is not compilation evidence.

Continue callback removal by tracing all producers, registrations and fixture
observers. Remaining families include GUI/drawable, object/AI, sustained spells,
audio streams and smaller legacy callbacks. Audit by owner before rejecting foreign
keys; migrate shared observers and preserve existing captures. Also remaining:
allocator ownership, reachable abort, compiler flags and fixture dependencies.
Literal allocator calls occur in 40 C-using fixture files; avoid cross-domain frees.

Animation, quantity, image completion and modifier dispatch reject unregistered
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
`build/port-modifier-dispatch/{contracts/profiles,safe,production/production/bin}`.
Original modifier captures live under `build/port-modifier-dispatch/original`,
with safe contracts in `original-safe`. Baseline runs use retained qualified UI
completion binaries at `defbdb98`. Older cleanup records retain source revisions,
hashes and rebuild information for superseded binaries.

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
inactive Linux386 cache archives untouched for six hours (496,467,968 bytes).
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
