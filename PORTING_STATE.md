# Porting checkpoint

This is the current resume checkpoint. Workflow and delegation rules live in
[PORT.md](PORT.md); detailed evidence belongs in the linked batch reports.

## Status: resumed; internal C-glue removal

**Progress: 142,665/142,665 original standalone C lines ported or retired;
internal glue: 458/463 client cgo files eliminated on net (5 remain;
server: 457/463 eliminated, 6 remain).**
Selected legacy C export bridges: **1,890/1,890 retired (0 remain)**.
Embedded production C bodies: **60/79 retired (19 remain)**.

These are selected Linux 386 production files, not equal units of effort.
Three project packages directly use cgo. Standalone production/test C remain zero.
Latest qualified chunk retires an unused spell-call adapter and makes four live
spell scalar helper return types native Go. Source changes only in legacy/spells.go.
All 2,482/2,471/2,482 selected default/server/highres roots passed without skips.
See [SPELL_SCALAR_BOUNDARIES.md](docs/porting/SPELL_SCALAR_BOUNDARIES.md).

Continue chunk-by-chunk with primary review, qualification, documentation and
commit/push. Use one Luna helper when its quota is available; no substitute model.
Stop at the milestone, usage limits or a substantial question. Latest qualified
artifacts: `build/port-spell-scalar-boundaries/`.

## What remains

| Area | Remaining work or dependency |
| --- | --- |
| Embedded C callback glue | 19 generic dispatch bodies; one signature is used only by fixtures. Trace all owners/registrations before removing raw fallbacks. |
| Types and declarations | 157 tracked headers / 2,731 physical lines. 103 porttest-tagged source files still import C across all build tags (not one selected profile). Fixture C observers/types still need retirement or explicit qualification scope. |
| Production C imports | alloc/raw.go, ccall/ccall.go, legacy/object_xfer_world.go (reachable abort), legacy/cgo_common.go and legacy/video_highres.go (flags). Server adds legacy/cgo_server.go. |
| Memory and layout | C-heap allocator, raw pointers, fixed offsets and 32-bit address assumptions remain. Preserve ownership/lifetime/failure semantics behind the centralized allocator. |
| External libraries | SDL2, OpenGL, OpenAL and similar bindings remain for this phase; later replacement requires a separate discussion. |
| Portability | Qualified target remains Linux 386/SSE2 with cgo. Whole-build cgo-free, 64-bit and other platforms are not qualified. Physical display/audible playback remain manual checks. |

## Latest qualification

Original baseline `b6049004` reuses exact-source `32df9553` results after source,
supplemental-source, environment, binary-hash and discovered-name checks. The
selected compiled roots are strict subsets of those previously accepted. New
94-root native spellbook/quickbar/objective preflight passed, followed by matching
three-profile roots, safe build/static, all three production builds/ABI, exact
known-suite failure/package outcomes and fresh headless save/load/resume.
The specialized adapter symbol is absent from all three production binaries.
All 1,654 original asset hashes are unchanged. No source expectations changed.

Selection is conservative through shared initialization, not a claim to every
root in the repository. The opt-in map-population diagnostic and roots outside
the audited graph are excluded; profile-specific uncompiled names are recorded.
Last complete default corpus: `6e9681f2` (2,489 passes plus diagnostic skip).

Evidence: [qualification](docs/porting/spell-scalar-boundaries-qualification.json),
[inventory](docs/porting/spell-scalar-boundaries-inventory-after.json).
Known-suite expectation: [record](docs/porting/mp3-go-expected-suite.jsonl).
Standalone metric/history: [C_LOC.md](docs/porting/C_LOC.md).

## Next work and review items

The [immediate goal](PORT.md#goal-and-target) remains internal engine C-glue
removal. Follow [INTERNAL_C_GLUE.md](docs/porting/INTERNAL_C_GLUE.md).
Metadata discovery is not compilation evidence.

In progress: callback fixture retirement has an accepted original baseline and
a reviewed, uninstalled nine-file draft plus one obsolete test deletion. See
[CALLBACK_FIXTURE_RETIREMENT.md](docs/porting/CALLBACK_FIXTURE_RETIREMENT.md).
29 original roots passed in all three profiles; private adapter/safe radial
checks passed too. Require the 28 surviving roots after conversion.
Artifacts: `build/port-callback-retirement-audit/`.
The object-construction wrapper has only its own fixture as caller; prefab group
traversal has only recursive/fixture callers. Confirm full reference scope before
retirement. An animation callback has no named non-nil writes, but raw-offset
writes still need checking. None of this proves every generic fallback unused.
Then close allocator, reachable abort, compiler flags and fixture dependencies.

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

Current converted root/safe/production binaries live under
`build/port-spell-scalar-boundaries/{contracts/profiles,safe,production/production/bin}`.
The current original baseline uses these qualified spell binaries. Seven older
string binaries were removed after replacement/source/hash and host-use checks
(387,117,056 bytes); their logs/source maps remain, and `32df9553` can be rebuilt.
Cleanup records: `build/port-callback-retirement-audit/string-cleanup-*`.
Spell baseline/qualification metadata and reports are committed; ignored local
binaries/logs/drafts are not backed up by Git. Completed scripts are consumed.

Seven superseded placement binaries were removed after committed-source,
qualified-replacement/hash and host-use checks (387,072,000 allocated bytes).
Rebuild `81cf750f` with retained commands/source maps. Cleanup records:
`build/port-spell-scalar-boundaries/placement-cleanup-{approved.json,deleted.jsonl}`.
Final spell-scenario duplicate assets were verified and removed; saves/results
remain. Restore before replay:
`python3 build/port-artifact-cleanup/restore-recent-scenario.py build/baseline/runs/spell-scalar-boundaries-save/deduplicated-assets.json`.

Check free disk before large runs. Historical artifacts may be gzip-archived or
deduplicated; retain source/rebuild commands and verify manifests before reuse.
Placement original/repeat/converted logs and string original logs restore with
`gzip -dk FILE.jsonl.gz`; exact hashes/paths are recorded in
`build/port-string-boundaries/log-archive.json`.
Earlier cleanup/recovery inventories remain in the corresponding batch reports,
ignored manifests and checkpoint history at `32df9553` (historical instructions
there are not current). Do not rerun old cleanup scripts or infer safety from age.

Six older original-baseline executables from `9dcf1b9a` and `7034a7e4` were also
removed after committed-source, replacement/hash and host-use checks
(400,736,256 allocated bytes). Rebuild those revisions with retained binary
records; current original evidence uses the qualified string binaries.
Records: `build/port-spell-scalar-boundaries/old-baselines-{approved.json,deleted.jsonl}`.
