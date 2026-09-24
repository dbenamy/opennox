# Porting checkpoint

This file is the current resume checkpoint, not a chronological log. Replace
superseded status when updating it. The workflow and delegation rules live in
[PORT.md](PORT.md); detailed evidence belongs in the linked batch reports.

## Status: resumed; internal C-glue removal

**Progress: 142,665/142,665 original standalone C lines ported or retired;
internal glue: 233/463 cgo files eliminated on net (230 remain).**
Selected legacy C export bridges: **748/1,890 retired (1,142 remain)**.

These are selected project files in each Linux 386 production profile, not equal
units of effort. Three project packages directly use cgo; 78 embedded C callback
bodies remain. Production and test-reference standalone `.c` files both remain zero.

Latest qualified chunk retires 37 fixture-only object-state C exports and
prototypes, moves 19 production caller files to native owners, and replaces the
fixture's C dispatcher. Four live collision callbacks retain their identities.
All 22 changed files qualify with frozen expectations unchanged.
See [OBJECT_STATE_OWNERS.md](docs/porting/OBJECT_STATE_OWNERS.md).

Continue chunk-by-chunk with one Luna helper, qualification, commit/push and
recorded reversible decisions. Next: audit the connected equipment/inventory/
resource owners and test-only exports under `build/port-inventory-owner-audit/`.
That audit is a candidate, not an accepted conversion. Stop at the milestone or
for a substantial question.

Latest artifacts: `build/port-object-state-owners/`.

## What remains

Counts below describe the qualified object-state conversion. Zero `.c` lines is
not a count of all C dependencies or a measure of remaining engineering effort.

| Area | Remaining work or dependency |
| --- | --- |
| Embedded C callback glue | 78 production function bodies in Go preambles: 76 generic function-pointer dispatchers and two specialized adapters. |
| Callback routes | Some Go implementations still call each other through C-compatible addresses. More direct Go dispatch is possible; shared raw fallbacks remain until their users and compatibility requirements are resolved. |
| Declarations and C types | 157 tracked headers / 3,865 physical lines; each production profile selects 230 cgo files in three project packages (alloc, ccall, legacy). Selected-build counts replace the earlier whole-tree text count. These are mostly interface/layout machinery, not unported algorithms. |
| Memory and layout | C-heap allocation, raw pointers, fixed offsets and 32-bit address assumptions remain. Removing them requires ownership/layout changes beyond function translation. |
| External libraries | SDL2, OpenGL, OpenAL and similar native dependencies and their cgo bindings stay for this phase; their future is a subsequent discussion. |
| Portability and release validation | Qualified target is Linux 386/SSE2 with cgo. The complete engine is not qualified as cgo-free, 64-bit, native macOS or browser/WebAssembly. Physical display and audible playback remain manual release checks. |

See [C_LOC.md](docs/porting/C_LOC.md) for the exact standalone-line metric and
history, and [TYPED_CALLBACK_ADAPTERS.md](docs/porting/TYPED_CALLBACK_ADAPTERS.md)
for the remaining embedded callback functions. The preceding batch retired one unused specialized body; shared raw fallback
machinery remains until its live callers are migrated.

## Latest qualification and evidence

- Exact focused root-name sets pass: 283 each in default/server/highres, no skips.
- Safe build/static checks and three fresh production binaries/ABI checks pass.
  All 37 retired C symbols are absent; the four live callbacks remain.
- Headless character creation and explicit save/load/resume pass.
- Full-suite results match the known baseline exactly: 304 failure events,
  with 17 passing, two failing and 32 skipped packages.
- All phases use identical source fingerprints. All 22 changed files match the
  reviewed draft; assertions, frozen expectations and 1,654 asset hashes are unchanged.
- Before conversion, two identity probes exercised all 44 object-state operations
  without consuming a registered state-function address. Existing reservation
  accounting preserves snapshot IDs after 37 address registrations are removed.
- The preceding shared-record milestone passed all seven storage captures and the
  complete root corpus (2,426 client / 2,415 server passes plus one expected skip).

Report: [OBJECT_STATE_OWNERS.md](docs/porting/OBJECT_STATE_OWNERS.md).
Evidence: [qualification](docs/porting/object-state-owners-qualification.json).
Known-suite expectation: [mp3-go-expected-suite.jsonl](docs/porting/mp3-go-expected-suite.jsonl).

## Goal, next work and open review items

The [immediate goal](PORT.md#goal-and-target) is removal of the engine's internal
C glue, retaining external native-library bindings and x86/32-bit assumptions.
Whole-build `CGO_ENABLED=0` is deferred for subsequent discussion.
Follow [INTERNAL_C_GLUE.md](docs/porting/INTERNAL_C_GLUE.md) for the dependency
removal order and completion criteria. Client rendering/audio backend replacement
is outside this phase.

The dependency inventory tool is `tools/porting/cgo_inventory.py`; the current
qualified inventory is [object-state-owners-inventory-after.json](docs/porting/object-state-owners-inventory-after.json).
The original phase baseline is under `build/port-cgo-leaves/inventory-before/`.
The completed leaf cleanup leaves three project packages directly using cgo in
all profiles, plus OpenGL/SDL2/OpenAL bindings in the clients. Metadata discovery
is not compilation or qualification. The helper's external-review draft is not
accepted evidence: its suggestion that go-gl is residue is contradicted by the
actual dependency graph (`libs/client/seat/opengl` imports it).

The earlier monster-callback candidate remains deferred; its context is in
INTERNAL_C_GLUE.md. The earlier 25-export object-state proposal under
`build/port-after-go-only-exports/` was rejected because it omitted fixture C calls.
The qualified 37-export batch supersedes it with a complete caller/identity audit.

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

Latest local artifacts are under `build/port-object-state-owners/`:
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
| Old owner/constants project cache | 65 hash/stat-verified, unused root/legacy archives removed; 4,718,157,824 allocated bytes reclaimed. Rebuild normally. Plan/journal: `build/port-go-native-owner-constants/cache-cleanup-{approved.json,deleted.jsonl}`. |
| Native-record full-corpus logs | Losslessly compressed; restore with `gzip -dk`. Hashes: `build/port-go-native-record-storage/contract-log-archive.json`; 306,063,181 bytes reclaimed. |
| Completed native-record scenario assets | Removed 1,654 verified original-asset duplicates; 559,874,048 allocated bytes reclaimed. Saves/results remain. Restore with `python3 build/port-artifact-cleanup/restore-recent-scenario.py build/baseline/runs/go-native-record-storage-save/deduplicated-assets.json`. |
| Superseded Go-only-export binaries | Seven test/safe/production executables removed after primary source/replacement/hash and host-use checks; 395,526,144 allocated bytes reclaimed. Five source maps (3,062 unique files) match `12bc387d`; retain phase commands and rebuild that revision. Plan/journal: `build/port-go-primitive-interfaces/cleanup-go-only-{approved.json,deleted.jsonl}`. |
| Superseded 378-export root test binaries | Three executables removed after host-use, inode and hash checks; 204,169,216 allocated bytes reclaimed. All 3,090 recorded source fingerprints match `f6f5ee4c`; rebuild that revision with the retained profile commands. Plan/journal: `build/port-go-primitive-interfaces/cleanup-binaries-{approved-plan.json,deleted.jsonl}`. |
| Layout-boundary full-corpus logs | Losslessly compressed; restore with `gzip -dk`. Hashes/commands: `build/port-go-layout-boundaries/contract-log-archive.json`; 306,061,277 bytes reclaimed. |
| Native-boundary full-corpus logs | Losslessly compressed; restore with `gzip -dk`. Hashes and commands: `build/port-go-native-call-boundaries/contract-log-archive.json`; 306,065,119 bytes reclaimed. |
| Completed layout-boundary scenario assets | Removed only 1,654 verified original-asset duplicates, reclaiming 559,890,432 allocated bytes. Saves/results remain. Restore with `python3 build/port-artifact-cleanup/restore-recent-scenario.py build/baseline/runs/go-layout-boundaries-save/deduplicated-assets.json`. |
| Completed native-boundary scenario assets | Removed only 1,654 verified original-asset duplicates, reclaiming 559,943,680 allocated bytes. Saves/results remain. Restore with `python3 build/port-artifact-cleanup/restore-recent-scenario.py build/baseline/runs/go-native-call-boundaries-save/deduplicated-assets.json`. |
| Superseded native-boundary binaries | Seven executables removed after committed-source/replacement/hash and host-use checks; 395,464,704 allocated bytes reclaimed. Old phase fingerprints match `4214ea8f`, replacements match `fe44bab3`; recorded build HEADs are earlier baseline commits. Rebuild those qualified revisions. Plan/journal: `build/port-go-native-record-storage/cleanup-native-{approved.json,deleted.jsonl}`. |
| Current qualified production/safe binaries | Retained under `build/port-object-state-owners/`; preceding native-owner, native-record and layout-boundary binaries also remain. |
| Superseded primitive-interface binaries | Seven executables removed after source/replacement/hash and host-use checks; 395,362,304 allocated bytes reclaimed. All 3,062 source fingerprints match `78ff20f9`; current qualified replacements match `4214ea8f`. Rebuild the old revision using retained phase commands. Plan/journal: `build/port-go-layout-boundaries/cleanup-primitive-{approved.json,deleted.jsonl}`. |
| Superseded scalar-storage binaries | Seven executables removed after source/replacement/hash and host-use checks; 395,452,416 allocated bytes reclaimed. All 3,062 source fingerprints match `e64ff24e`; qualified replacements match `78ff20f9`. Rebuild the old revision using retained phase commands. Plan/journal: `build/port-go-native-call-boundaries/cleanup-scalar-{approved.json,deleted.jsonl}`. |
| Completed scenario data: `go-memory-save`, `raw-allocation-save`, `string-boundary-save`, `unused-exports-save`, `remaining-unused-exports-save`, `go-only-exports-save` | Only SHA256-identical original-asset copies were removed. Saves/comparisons remain. Follow each run's `deduplicated-assets.json`; shared restore tool: `build/port-artifact-cleanup/restore-recent-scenario.py`. |
| Large historical captures in `port-game-messages`, `port-map-sections`, `port-client-interaction`, `port-session-dialogs` | Restore with `gzip -dk` and verify hashes against `build/port-artifact-cleanup/large-historical-20260924/`. Its 116 discarded text logs are not recoverable from these archives. |
| Initial complete-corpus default/server logs | Losslessly compressed; restore commands and hashes in `build/port-complete-corpus/initial-log-archive.json`. Keep the server failure evidence. |
| Latest scalar full-corpus logs | Losslessly compressed after qualification/commit; restore commands and SHA256s: `build/port-go-scalar-storage/contract-log-archive.json`. |
| Prior 265-export full-corpus logs | Losslessly compressed after qualification/commit; restore commands and SHA256s: `build/port-go-only-exports/contract-log-archive.json`. |
| Obsolete project cache archives removed before layout qualification | 26 verified old archives removed after host fd/maps checks, reclaiming 2,116,055,040 allocated bytes; rebuild with the documented Go environment. Plan/journal: `build/port-go-layout-boundaries/cache-cleanup-{approved,deleted}.json`. Current binaries, assets and source retained. |
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
