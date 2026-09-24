# Porting checkpoint

This file is the current resume checkpoint, not a chronological log. Replace
superseded status when updating it. The workflow and delegation rules live in
[PORT.md](PORT.md); detailed evidence belongs in the linked batch reports.

## Status: resumed; internal C-glue removal

**Progress: 142,665/142,665 original standalone C lines ported or retired;
internal glue: 192/463 cgo files eliminated on net (271 remain).**
Selected legacy C export bridges: **711/1,890 retired (1,179 remain)**.

These are selected project files in each Linux 386 production profile, not equal
units of effort. Three project packages directly use cgo; 79 embedded C callback
bodies remain. Production and test-reference standalone `.c` files both remain zero.

Latest qualified chunk moves Go callers across internal adapter boundaries,
extracting native helper bodies and preserving C entrypoints, scalar widths,
raw aliases, nil interfaces and server-hook evaluation order. It removes 31
production C imports across 50 reviewed source files. All qualification passes;
no root fixture or frozen expectation changed. See
[GO_NATIVE_CALL_BOUNDARIES.md](docs/porting/GO_NATIVE_CALL_BOUNDARIES.md).

Continue chunk-by-chunk with one Luna helper, qualification, commit/push and
recorded reversible decisions. Active: the native geometry/shared-state batch
has an accepted original baseline from `4214ea8f`, including target layout probes.
No conversion is installed yet. Review the full caller set around the initial
23 files, 26 remaining state scalar fields and 132 orphan-wrapper candidates.
See [GO_LAYOUT_BOUNDARIES.md](docs/porting/GO_LAYOUT_BOUNDARIES.md).
Stop at the milestone or for a substantial question.

Latest artifacts: `build/port-go-native-call-boundaries/`.

## What remains

Counts below describe the qualified native-boundary conversion. Zero `.c` lines is
not a count of all C dependencies or a measure of remaining engineering effort.

| Area | Remaining work or dependency |
| --- | --- |
| Embedded C callback glue | 79 production function bodies in Go preambles: 76 generic function-pointer dispatchers and three specialized adapters. |
| Callback routes | Some Go implementations still call each other through C-compatible addresses. More direct Go dispatch is possible; shared raw fallbacks remain until their users and compatibility requirements are resolved. |
| Declarations and C types | 157 tracked headers / 3,902 physical lines; each production profile selects 271 cgo files in three project packages (alloc, ccall, legacy). Selected-build counts replace the earlier whole-tree text count. These are mostly interface/layout machinery, not unported algorithms. |
| Memory and layout | C-heap allocation, raw pointers, fixed offsets and 32-bit address assumptions remain. Removing them requires ownership/layout changes beyond function translation. |
| External libraries | SDL2, OpenGL, OpenAL and similar native dependencies and their cgo bindings stay for this phase; their future is a subsequent discussion. |
| Portability and release validation | Qualified target is Linux 386/SSE2 with cgo. The complete engine is not qualified as cgo-free, 64-bit, native macOS or browser/WebAssembly. Physical display and audible playback remain manual release checks. |

See [C_LOC.md](docs/porting/C_LOC.md) for the exact standalone-line metric and
history, and [TYPED_CALLBACK_ADAPTERS.md](docs/porting/TYPED_CALLBACK_ADAPTERS.md)
for the remaining embedded callback functions. Recent dispatch conversions do
not reduce the 79-body count because they retain the shared fallback machinery.

## Latest qualification and evidence

- Both storage contracts pass in default/server/highres; safe passes the scalar
  contract separately. All seven captures equal frozen hashes with exact test
  sets and no skips. Raw safe remains intentionally excluded because that fixture
  crosses an address rejected by safe; no guard was weakened.
- Complete root suites: default/highres each 2,426 pass; server 2,415 pass. Each
  has only the expected `TestMapPopulationPrerequisiteProbe` skip. Discovered,
  started and completed root names match the historical inventory plus the existing
  collision regression, with no failure events.
- Safe build/static checks and three fresh production binaries/ABI checks pass.
  No complete safe runtime suite or safe raw-storage pass is claimed.
- Headless character creation and explicit save/load/resume pass.
- Full-suite results exactly match the known baseline: 304 failure events,
  with 17 passing, two failing and 32 skipped packages.
- Every phase uses identical source fingerprints. Reconstruction verifies all
  50 changed source files; no root fixture correction was needed. All frozen
  expectations and all 1,654 original asset hashes remain unchanged.

Report: [GO_NATIVE_CALL_BOUNDARIES.md](docs/porting/GO_NATIVE_CALL_BOUNDARIES.md).
Evidence: [qualification](docs/porting/go-native-call-boundaries-qualification.json).
Known-suite expectation: [mp3-go-expected-suite.jsonl](docs/porting/mp3-go-expected-suite.jsonl).

## Goal, next work and open review items

The [immediate goal](PORT.md#goal-and-target) is removal of the engine's internal
C glue, retaining external native-library bindings and x86/32-bit assumptions.
Whole-build `CGO_ENABLED=0` is deferred for subsequent discussion.
Follow [INTERNAL_C_GLUE.md](docs/porting/INTERNAL_C_GLUE.md) for the dependency
removal order and completion criteria. Client rendering/audio backend replacement
is outside this phase.

The dependency inventory tool is `tools/porting/cgo_inventory.py`; the current
qualified inventory is [go-native-call-boundaries-inventory-after.json](docs/porting/go-native-call-boundaries-inventory-after.json).
The original phase baseline is under `build/port-cgo-leaves/inventory-before/`.
The completed leaf cleanup leaves three project packages directly using cgo in
all profiles, plus OpenGL/SDL2/OpenAL bindings in the clients. Metadata discovery
is not compilation or qualification. The helper's external-review draft is not
accepted evidence: its suggestion that go-gl is residue is contradicted by the
actual dependency graph (`libs/client/seat/opengl` imports it).

The earlier monster-callback candidate remains deferred; its context is in
INTERNAL_C_GLUE.md. The follow-on 25-export object-state proposal under
`build/port-after-go-only-exports/` is **not accepted**: its audit omitted active
C calls inside a test preamble. Those interfaces remain unchanged.

The integer-returning collision owner in `temporaryMagicMissile` remains raw;
its return cannot be replaced with a void-dispatch result. See
[collision compatibility decisions](docs/porting/COLLISION_REGISTRY.md).
Other behavior/compatibility findings are recorded in [DECISIONS.md](docs/porting/DECISIONS.md).

Luna's bounded edit manifests remain useful with primary reconstruction and AST
review. Broad reachability audits needed correction for package scope and C
preambles; keep algorithm design and acceptance with the primary. The scalar
draft's C-prefix replacement error was corrected before installation. Preserve
these explicit checks rather than treating a helper's reconstruction as acceptance.

## Resume and artifact recovery

Read this checkpoint and PORT.md, then inspect Git status before editing. The
untracked `nox-iso-from-archive-org.7z` is an asset archive, not unfinished code;
never stage it. Preserve it and `build/assets/extracted/drive_c/Nox`.

The current Go toolchain is `/usr/lib/go-1.26/bin`. Follow the
[build environment instructions](PORT.md#build-and-test-environment), including
sourcing `build/baseline/env.sh` in every Go shell.

Latest local artifacts are under `build/port-go-native-call-boundaries/`:
`contracts/`, `safe/opennox-safe`, and
`production/production/bin/{opennox,opennox-hd,opennox-server}`.
Source, tests, reports and qualification metadata are committed; ignored local
binaries/logs/drafts are not backed up by pushing Git. Completed finalizers are
consumed and must not be rerun against later sources.

Check free disk before large runs. Older artifacts may be gzip-archived,
hardlinked or deduplicated. Restore scenario data before replay and separate
shared binary inodes before modifying either path. Cleanup scripts are consumed;
do not rerun them or infer deletion safety from age alone.

| Artifact | Recovery or current location |
| --- | --- |
| Superseded Go-only-export binaries | Seven test/safe/production executables removed after primary source/replacement/hash and host-use checks; 395,526,144 allocated bytes reclaimed. Five source maps (3,062 unique files) match `12bc387d`; retain phase commands and rebuild that revision. Plan/journal: `build/port-go-primitive-interfaces/cleanup-go-only-{approved.json,deleted.jsonl}`. |
| Superseded 378-export root test binaries | Three executables removed after host-use, inode and hash checks; 204,169,216 allocated bytes reclaimed. All 3,090 recorded source fingerprints match `f6f5ee4c`; rebuild that revision with the retained profile commands. Plan/journal: `build/port-go-primitive-interfaces/cleanup-binaries-{approved-plan.json,deleted.jsonl}`. |
| Native-boundary full-corpus logs | Losslessly compressed; restore with `gzip -dk`. Hashes and commands: `build/port-go-native-call-boundaries/contract-log-archive.json`; 306,065,119 bytes reclaimed. |
| Completed native-boundary scenario assets | Removed only 1,654 verified original-asset duplicates, reclaiming 559,943,680 allocated bytes. Saves/results remain. Restore with `python3 build/port-artifact-cleanup/restore-recent-scenario.py build/baseline/runs/go-native-call-boundaries-save/deduplicated-assets.json`. |
| Current qualified production/safe binaries | Retained under `build/port-go-native-call-boundaries/`; preceding qualified primitive-interface binaries also remain. |
| Superseded scalar-storage binaries | Seven executables removed after source/replacement/hash and host-use checks; 395,452,416 allocated bytes reclaimed. All 3,062 source fingerprints match `e64ff24e`; qualified replacements match `78ff20f9`. Rebuild the old revision using retained phase commands. Plan/journal: `build/port-go-native-call-boundaries/cleanup-scalar-{approved.json,deleted.jsonl}`. |
| Completed scenario data: `go-memory-save`, `raw-allocation-save`, `string-boundary-save`, `unused-exports-save`, `remaining-unused-exports-save`, `go-only-exports-save` | Only SHA256-identical original-asset copies were removed. Saves/comparisons remain. Follow each run's `deduplicated-assets.json`; shared restore tool: `build/port-artifact-cleanup/restore-recent-scenario.py`. |
| Large historical captures in `port-game-messages`, `port-map-sections`, `port-client-interaction`, `port-session-dialogs` | Restore with `gzip -dk` and verify hashes against `build/port-artifact-cleanup/large-historical-20260924/`. Its 116 discarded text logs are not recoverable from these archives. |
| Initial complete-corpus default/server logs | Losslessly compressed; restore commands and hashes in `build/port-complete-corpus/initial-log-archive.json`. Keep the server failure evidence. |
| Latest scalar full-corpus logs | Losslessly compressed after qualification/commit; restore commands and SHA256s: `build/port-go-scalar-storage/contract-log-archive.json`. |
| Prior 265-export full-corpus logs | Losslessly compressed after qualification/commit; restore commands and SHA256s: `build/port-go-only-exports/contract-log-archive.json`. |
| Remaining old project cache archives | 26 verified old root/legacy archives removed after host checks, reclaiming 1,614,089,922 bytes; rebuild normally. Plan/journal: `build/port-artifact-cleanup/go-scalar-storage-cache-{plan,removed}.json`. |
| 378-export full-corpus logs | Losslessly compressed after qualification/commit; restore commands and SHA256s: `build/port-remaining-unused-exports/contract-log-archive.json`. |
| Additional obsolete project cache archives | 31 old root/legacy archives removed after host checks, reclaiming 2,017,688,956 bytes; rebuild normally. Plan/journal: `build/port-artifact-cleanup/go-only-exports-cache-{plan,removed}.json`. |
| Qualified complete-corpus logs | Losslessly compressed; restore commands/hashes: `build/port-complete-corpus/qualified-log-archive.json`. |
| Completed pointer-fixture and prebuilt-pilot binaries | Five rebuildable binaries removed; source `4a0ab0dc`, original logs/records retained. Plan/journal: `build/port-artifact-cleanup/completed-corpus-binaries-{plan.json,removed.jsonl}`. |
| Superseded Go-memory/raw-allocation/string-boundary binaries | Rebuild from `4e1e86a6`, `6ff98033`, `92029ddd` respectively. Exact inventory/removal journal: `build/port-artifact-cleanup/superseded-glue-binaries-{plan.json,removed.jsonl}`. |
| Completed-chunk old project cache archives | 24 Sep23 root/legacy archives removed after all Go jobs joined; rebuild normally. Plan/journal: `build/port-artifact-cleanup/remaining-exports-cache-{plan,removed}.json`. |
| Older leaf/transfer binaries and cache entries | Rebuild from recorded revisions as needed; cleanup records are in `build/port-artifact-cleanup/`, including `superseded-leaves-xfer-binaries-20260924/` and `string-boundary-cache-removed.jsonl`. |

Historical batch details belong in `docs/porting/` and Git history. Earlier
checkpoint/archive details are recoverable at `b034c43e` and `69669bcb`; their
old current/next instructions are historical, not the active plan. Local logs,
archives and rebuildable binaries are not backed up by pushing Git.

Initial primitive-interface root logs are losslessly archived under
`build/port-go-primitive-interfaces/contracts/profiles/*.jsonl.gz`; restore with
`gzip -dk FILE.jsonl.gz`. The hash/size manifest is
`build/port-go-primitive-interfaces/contract-log-archive.json` (306,057,005 bytes
reclaimed). Failure captures and result metadata remain directly readable.

Final primitive-interface root logs are also losslessly archived under
`contracts-final/profiles/*.jsonl.gz`; `final-contract-log-archive.json` records
hashes and restore commands (306,065,024 bytes reclaimed). Original asset copies
in `go-scalar-storage-save` and `go-primitive-interfaces-save` were removed only
after byte/hash and host-use checks: 3,308 duplicates / 1,119,789,056 allocated
bytes. Saves, changed files and results remain. Restore either scenario with
`python3 build/port-artifact-cleanup/restore-recent-scenario.py` followed by its
`build/baseline/runs/SCENARIO/deduplicated-assets.json` path.
