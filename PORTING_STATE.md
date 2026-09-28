# Porting checkpoint

This is the current resume checkpoint. Workflow and delegation rules live in
[PORT.md](PORT.md); detailed evidence belongs in the linked batch reports.

## Status: resumed; internal C-glue removal

- **Standalone engine C:** 142,665/142,665 original lines ported or retired.
- **Production C imports:** 458/463 client cgo files eliminated on net (5 remain);
  server: 457/463 eliminated (6 remain). Three project packages directly use cgo.
- **Legacy C export bridges:** 1,890/1,890 retired (0 remain).
- **Embedded production C bodies:** 71/79 retired (8 remain).
- **Fixture and declaration dependencies:** 83 porttest-tagged files still import C;
  157 tracked headers / 2,731 physical lines remain.

Production cgo counts cover selected Linux 386 builds; these metrics are not equal
units of effort. Fixture imports cover all build tags. Standalone production/test C
remain zero.
Latest qualified chunk retires foreign object-damage and monster dispatch and
converts five shared observer fixtures, including collision/update/use routes.
Two embedded dispatchers and five fixture C imports are retired. See
[DAMAGE_MONSTER_DISPATCH.md](docs/porting/DAMAGE_MONSTER_DISPATCH.md).

Continue chunk-by-chunk with primary review, qualification, documentation and
commit/push. Use one Luna helper when its quota is available; no substitute model.
Stop at the milestone, usage limits or a substantial question. Latest qualified
artifacts: `build/port-damage-monster-dispatch/`.

## What remains

| Area | Remaining work or dependency |
| --- | --- |
| Embedded C callback glue | 8 generic dispatch bodies. Trace all owners/registrations before removing raw fallbacks. |
| Types and declarations | 157 tracked headers / 2,731 physical lines. 83 porttest-tagged source files still import C across all build tags (not one selected profile). Fixture C observers/types still need retirement or explicit qualification scope. |
| Production C imports | alloc/raw.go, ccall/ccall.go, legacy/object_xfer_world.go (reachable abort), legacy/cgo_common.go and legacy/video_highres.go (flags). Server adds legacy/cgo_server.go. |
| Memory and layout | C-heap allocator, raw pointers, fixed offsets and 32-bit address assumptions remain. Preserve ownership/lifetime/failure semantics behind the centralized allocator. |
| External libraries | SDL2, OpenGL, OpenAL and similar bindings remain for this phase; later replacement requires a separate discussion. |
| Portability | Qualified target remains Linux 386/SSE2 with cgo. Whole-build cgo-free, 64-bit and other platforms are not qualified. Physical display/audible playback remain manual checks. |

## Latest qualification

Broad native selections passed 2,482/2,471/2,482 roots without skips. Exact-source
qualified GUI broad runs were reused for the original baseline; twelve focused
contracts per profile and four safe contracts ran freshly before/after conversion.
Native preflight, safe build/static, three production builds/ABI, exact known-suite
outcomes and fresh save/load passed. Assertions, native handlers, retained C
bodies and original assets are unchanged.

Evidence: [qualification](docs/porting/damage-monster-dispatch-qualification.json),
[inventory](docs/porting/damage-monster-dispatch-inventory-after.json).
Known-suite expectation: [record](docs/porting/mp3-go-expected-suite.jsonl).
Standalone metric/history: [C_LOC.md](docs/porting/C_LOC.md).
Latest broad regression: this damage/monster batch, 2,482/2,471/2,482 audited roots, no skips.
Last complete default corpus: `6e9681f2`, 2,489 passes plus diagnostic skip.

## Next work and review items

Continue with remaining object lifecycle/use/collision and drawable draw/update,
spatial/particle callbacks. A provisional 24-file draft plus three deletions lives
under `build/port-final-callback-dispatch/`; it targets the eight remaining generic
dispatch bodies together. It still needs semantic/producer review, an accepted
original baseline (including legacy-package contracts), installation and full
qualification. Allocator/abort/flags and fixture dependencies remain afterward.

The [immediate goal](PORT.md#goal-and-target) remains internal engine C-glue
removal. Follow [INTERNAL_C_GLUE.md](docs/porting/INTERNAL_C_GLUE.md).
Metadata discovery is not compilation evidence.

Continue callback removal by tracing all producers, registrations and fixture
observers. Remaining families include drawable, object/AI and smaller legacy callbacks. Audit by owner before rejecting
foreign keys; migrate shared observers and preserve existing captures. Also remaining:
allocator ownership, reachable abort, compiler flags and fixture dependencies.
Literal allocator calls occur in 40 C-using fixture files; avoid cross-domain frees.

Animation, quantity, image completion, modifier, player-section, screen-particle,
duration-spell, internal audio, GUI tooltip, damage and monster dispatch reject unregistered
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
`build/port-damage-monster-dispatch/{contracts/profiles,safe,production/production/bin}`.
Original damage/monster focused captures and safe runs live under
`build/port-damage-monster-dispatch/{original,original-safe}`. Original broad logs
and binaries are reused from the preceding GUI batch; exact paths/hashes are in
the committed baseline. Final damage/monster scenario asset copies were verified
and removed (560,435,200 bytes); restoration is documented in the batch report.
Older cleanup records retain source revisions, hashes and rebuild information.

Check physical free space before large runs. Source, reports and qualification
metadata are committed; ignored build artifacts are not backed up by Git. Current
binaries and useful original references remain locally. Completed scripts are
consumed: never rerun a completed installation, qualification or cleanup script.
Historical removals and lossless log-archive recovery are documented in the linked
batch reports and their ignored manifests. Rebuild superseded binaries from the
recorded source revision; cache archives rebuild automatically. Do not infer
safe deletion merely from age. Preserve original assets and the archive.

Prior GUI cleanup/recovery is in [WINDOW_DISPATCH.md](docs/porting/WINDOW_DISPATCH.md).
