# Porting checkpoint

This file is the current resume checkpoint, not a chronological log. Replace
superseded status when updating it. The workflow and delegation rules live in
[PORT.md](PORT.md); detailed evidence belongs in the linked batch reports.

## Status: resumed; internal C-glue removal

The original **142,665 physical lines of standalone C have been ported or
retired**. Production and test-reference `.c` files now both count zero.
The Go MP3 decoder is integrated and its C implementation header is retired.
The legacy algorithm-port milestone is complete; the engine still requires cgo.
Current work removes libc helpers and redundant **Go → C → Go** callback routes.

Latest qualified implementation: **`4e1e86a6` — Go memory/string helpers**,
following original baselines `5cc27785` and `58f37c6c`.
See [GO_MEMORY.md](docs/porting/GO_MEMORY.md). Six libc helper calls are replaced;
allocation/free ownership and external native bindings are unchanged.

Continue chunk-by-chunk with one Luna helper, qualification, commit/push and
recorded reversible decisions. Stop at the milestone or for a substantial question.
The requested C_LOC cleanup is complete: its historical table is continuous and
stale tail status notes are removed. Next: qualify and integrate allocation-call
centralization. Reviewed ignored drafts and original-path contract drafts are under
`build/port-go-memory/`; they are not installed or qualified.

## What remains

Counts below describe the qualified Go memory-helper conversion. Zero `.c` lines is
not a count of all C dependencies or a measure of remaining engineering effort.

| Area | Remaining work or dependency |
| --- | --- |
| Embedded C callback glue | 79 production function bodies in Go preambles: 76 generic function-pointer dispatchers and three specialized adapters. |
| Callback routes | Some Go implementations still call each other through C-compatible addresses. More direct Go dispatch is possible; shared raw fallbacks remain until their users and compatibility requirements are resolved. |
| Declarations and C types | 157 tracked headers / 4,548 physical lines; each production profile selects 429 cgo files in three project packages (alloc, ccall, legacy). Selected-build counts replace the earlier whole-tree text count. These are mostly interface/layout machinery, not unported algorithms. |
| Memory and layout | C-heap allocation, raw pointers, fixed offsets and 32-bit address assumptions remain. Removing them requires ownership/layout changes beyond function translation. |
| External libraries | SDL2, OpenGL, OpenAL and similar native dependencies and their cgo bindings stay for this phase; their future is a subsequent discussion. |
| Portability and release validation | Qualified target is Linux 386/SSE2 with cgo. The complete engine is not qualified as cgo-free, 64-bit, native macOS or browser/WebAssembly. Physical display and audible playback remain manual release checks. |

See [C_LOC.md](docs/porting/C_LOC.md) for the exact standalone-line metric and
history, and [TYPED_CALLBACK_ADAPTERS.md](docs/porting/TYPED_CALLBACK_ADAPTERS.md)
for the remaining embedded callback functions. Recent dispatch conversions do
not reduce the 79-body count because they retain the shared fallback machinery.

## Latest qualification and evidence

- Direct memory/string and allocator-class contracts pass in all three profiles.
- Safe memory bridges and shop loading pass without skips.
- Accumulated roots: 1,280 default, 1,276 server and 1,280 highres complete, with
  no failures and only the allowed map-population diagnostic skip in each.
- Safe build/static checks and three fresh production binaries/ABI checks pass.
- Headless character creation and explicit save/load/resume pass.
- Full-suite results match the known baseline exactly: 304 failure events,
  with 17 passing, two failing and 32 skipped packages.

Report: [GO_MEMORY.md](docs/porting/GO_MEMORY.md).
Evidence: [qualification](docs/porting/go-memory-qualification.json).
Known-suite expectation: [mp3-go-expected-suite.jsonl](docs/porting/mp3-go-expected-suite.jsonl).

## Goal, next work and open review items

The [immediate goal](PORT.md#goal-and-target) is removal of the engine's internal
C glue, retaining external native-library bindings and x86/32-bit assumptions.
Whole-build `CGO_ENABLED=0` is deferred for subsequent discussion.
Follow [INTERNAL_C_GLUE.md](docs/porting/INTERNAL_C_GLUE.md) for the dependency
removal order and completion criteria. Client rendering/audio backend replacement
is outside this phase.

The dependency inventory tool is `tools/porting/cgo_inventory.py`; the refreshed
baseline is under `build/port-cgo-leaves/inventory-before/`.
The completed leaf cleanup leaves three project packages directly using cgo in
all profiles, plus OpenGL/SDL2/OpenAL bindings in the clients. Metadata discovery
is not compilation or qualification. The helper's external-review draft is not
accepted evidence: its suggestion that go-gl is residue is contradicted by the
actual dependency graph (`libs/client/seat/opengl` imports it).

A preliminary candidate is 26 monster callbacks: 11 strike, five die and ten dead.
Their production owners are `combatMelee`, `lifecycleDyingStart` and
`lifecycleDeadStart`; the tables are at base 0x587000, offsets 287096, 287280 and
287192. Existing matrices exercise 25 entries. BomberDead and built-in identities
through the actual AI owners need targeted coverage review. Inventory and draft
paths are `build/port-xfer-sound-registry/next-monster-*`; they are ignored,
unapplied and unqualified. This is a candidate, not an instruction to start.

The integer-returning collision owner in `temporaryMagicMissile` remains raw;
its return cannot be replaced with a void-dispatch result. See
[collision compatibility decisions](docs/porting/COLLISION_REGISTRY.md).
The Go MP3 decoder remains slower than C in the recorded bounded benchmarks
(1.96× on one mono input, 2.56× on a synthetic stereo input); whole-game impact
has not been established. See [the performance report](docs/porting/MP3_SYNTHESIS_PERFORMANCE.md).
Other behavior/compatibility findings are recorded in [DECISIONS.md](docs/porting/DECISIONS.md).

Luna drafted the six helpers and reviewed primary contracts. Primary corrected
nonportable libc comparison expectations by restoring and recapturing the original
path, and replaced the slow fill loop after measurement. See the batch report.
The next allocation draft centralizes 49 calls across 21 files, preserving profile
semantics. Primary caught unused imports before integration; the draft is corrected
but unqualified. New domain/string ownership contracts and a 96-existing-root owner
selection await baseline qualification, including two additional existing free-owner
tests absent from prior selectors. See `build/port-go-memory/raw-*`.

## Resume and artifact recovery

Read this checkpoint and PORT.md, then inspect Git status before editing. The
untracked `nox-iso-from-archive-org.7z` is an asset archive, not unfinished code;
never stage it. Preserve it and `build/assets/extracted/drive_c/Nox`.

The current Go toolchain is `/usr/lib/go-1.26/bin`. Follow the
[build environment instructions](PORT.md#build-and-test-environment), including
sourcing `build/baseline/env.sh` in every Go shell.

Latest local artifacts are under `build/port-go-memory/`:
`helpers/`, `safe-contracts/`, `contracts/`, `safe/opennox-safe`, and
`production/production/bin/{opennox,opennox-hd,opennox-server}`.
Source, tests, reports and qualification metadata are committed; ignored local
binaries/logs/drafts are not backed up by pushing Git. Completed finalizers are
consumed and must not be rerun against later sources.

Check free disk before large runs. Older artifacts may be gzip-archived, hardlinked
or have duplicate scenario assets removed. Consult each artifact directory's
archive/removal records and each scenario's `deduplicated-assets.json`; use the
recorded restore script before historical replay. Separate shared binary inodes
before modifying either path. Cleanup records also live in
`build/port-artifact-cleanup/`. Some old binaries/cache entries were removed and
need rebuilding; original assets and current qualified binaries were retained.
Do not repeat consumed cleanup scripts or infer deletion safety from age alone.

The large historical captures/results in `port-game-messages`, `port-map-sections`,
`port-client-interaction` and `port-session-dialogs` were losslessly gzip-archived;
116 superseded text logs were discarded. This reclaimed 4.44 GiB. The journal and
restore instructions are in `build/port-artifact-cleanup/large-historical-20260924/`.
For a retained capture, `gzip -dk <original-path>.gz` restores the raw file; verify
its SHA256 against the journal. Discarded text logs cannot be restored this way.
Current qualified binaries, original assets and qualification summaries remain.

Historical batch details belong in `docs/porting/` and Git history. The complete
older checkpoint, including individual archive/restore notes, is recoverable with
`git show b034c43e:PORTING_STATE.md`; its old “current/next” instructions are
historical and must not be followed as the current plan.
