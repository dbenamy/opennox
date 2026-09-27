# Porting checkpoint

This is the current resume checkpoint. Workflow and delegation rules live in
[PORT.md](PORT.md); detailed evidence belongs in the linked batch reports.

## Status: resumed; internal C-glue removal

**Progress: 142,665/142,665 original standalone C lines ported or retired;
internal glue: 458/463 client cgo files eliminated on net (5 remain;
server: 457/463 eliminated, 6 remain).**
Selected legacy C export bridges: **1,890/1,890 retired (0 remain)**.
Embedded production C bodies: **69/79 retired (10 remain)**.

These are selected Linux 386 production files, not equal units of effort.
Three project packages directly use cgo. Standalone production/test C remain zero.
Latest qualified chunk retires unproduced raw GUI event/draw dispatch and converts
the tooltip observer. One embedded C dispatcher and one fixture C import are
retired. All 2,482/2,471/2,482 audited roots passed in default/server/highres. See
[WINDOW_DISPATCH.md](docs/porting/WINDOW_DISPATCH.md).

Continue chunk-by-chunk with primary review, qualification, documentation and
commit/push. Use one Luna helper when its quota is available; no substitute model.
Stop at the milestone, usage limits or a substantial question. Latest qualified
artifacts: `build/port-window-dispatch/`.

## What remains

| Area | Remaining work or dependency |
| --- | --- |
| Embedded C callback glue | 10 generic dispatch bodies. Trace all owners/registrations before removing raw fallbacks. |
| Types and declarations | 157 tracked headers / 2,731 physical lines. 88 porttest-tagged source files still import C across all build tags (not one selected profile). Fixture C observers/types still need retirement or explicit qualification scope. |
| Production C imports | alloc/raw.go, ccall/ccall.go, legacy/object_xfer_world.go (reachable abort), legacy/cgo_common.go and legacy/video_highres.go (flags). Server adds legacy/cgo_server.go. |
| Memory and layout | C-heap allocator, raw pointers, fixed offsets and 32-bit address assumptions remain. Preserve ownership/lifetime/failure semantics behind the centralized allocator. |
| External libraries | SDL2, OpenGL, OpenAL and similar bindings remain for this phase; later replacement requires a separate discussion. |
| Portability | Qualified target remains Linux 386/SSE2 with cgo. Whole-build cgo-free, 64-bit and other platforms are not qualified. Physical display/audible playback remain manual checks. |

## Latest qualification

Original/native broad selections passed 2,482/2,471/2,482 roots without skips.
Six focused contracts passed separately in each profile on both versions, plus
safe deferred cleanup. Native focused preflight and fresh default save/load passed
before the broad sweep. Safe build/static, three production builds/ABI, exact
known-suite outcomes and final fresh save/load passed. Assertions, snapshot IDs,
retained C bodies and original assets are unchanged.

Evidence: [qualification](docs/porting/window-dispatch-qualification.json),
[inventory](docs/porting/window-dispatch-inventory-after.json).
Known-suite expectation: [record](docs/porting/mp3-go-expected-suite.jsonl).
Standalone metric/history: [C_LOC.md](docs/porting/C_LOC.md).
Latest broad regression: this window batch, 2,482/2,471/2,482 audited roots, no skips.
Last complete default corpus: `6e9681f2`, 2,489 passes plus diagnostic skip.

## Next work and review items

Continue with the remaining object/AI, drawable draw/update and spatial/particle
callback families. Active batch: damage/monster dispatch. Original baseline is
accepted at qualified GUI source `0cf5064c`: broad runs reused after exact source,
environment, binary, log and test-name checks; twelve focused roots per profile
and four safe contracts ran freshly. Reviewed ten-file draft is not installed yet.
Next: commit baseline, install, then qualify. See
[DAMAGE_MONSTER_DISPATCH.md](docs/porting/DAMAGE_MONSTER_DISPATCH.md).

The [immediate goal](PORT.md#goal-and-target) remains internal engine C-glue
removal. Follow [INTERNAL_C_GLUE.md](docs/porting/INTERNAL_C_GLUE.md).
Metadata discovery is not compilation evidence.

Continue callback removal by tracing all producers, registrations and fixture
observers. Remaining families include drawable, object/AI and smaller legacy callbacks. Audit by owner before rejecting
foreign keys; migrate shared observers and preserve existing captures. Also remaining:
allocator ownership, reachable abort, compiler flags and fixture dependencies.
Literal allocator calls occur in 40 C-using fixture files; avoid cross-domain frees.

Animation, quantity, image completion, modifier, player-section, screen-particle,
duration-spell, internal audio and GUI tooltip dispatch reject unregistered
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
`build/port-window-dispatch/{contracts/profiles,safe,production/production/bin}`.
Original GUI captures live under `build/port-window-dispatch/original`, with safe
contracts in `original-safe`. Baseline runs use retained duration/audio binaries
at `d3d759ba`. Both completed GUI scenario asset copies were verified and removed
(1,120,215,040 bytes); see the batch report for restoration commands.
Older cleanup records retain source revisions, hashes and rebuild information.

Check physical free space before large runs. Source, reports and qualification
metadata are committed; ignored build artifacts are not backed up by Git. Current
binaries and useful original references remain locally. Completed scripts are
consumed: never rerun a completed installation, qualification or cleanup script.
Historical removals and lossless log-archive recovery are documented in the linked
batch reports and their ignored manifests. Rebuild superseded binaries from the
recorded source revision; cache archives rebuild automatically. Do not infer
safe deletion merely from age. Preserve original assets and the archive.

This batch's recovery details, including prior modifier-binary cleanup and cache
headroom, are in [WINDOW_DISPATCH.md](docs/porting/WINDOW_DISPATCH.md).
