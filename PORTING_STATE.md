# Porting checkpoint

This file is the current resume checkpoint, not a chronological log. Replace
superseded status when updating it. The workflow and delegation rules live in
[PORT.md](PORT.md); detailed evidence belongs in the linked batch reports.

## Status: resumed; internal C-glue removal

**Progress: 142,665/142,665 original standalone C lines ported or retired;
internal glue: 265/463 cgo files eliminated on net (198 remain).**
Selected legacy C export bridges: **1,151/1,890 retired (739 remain)**.

These are selected project files in each Linux 386 production profile, not equal
units of effort. Three project packages directly use cgo; 77 embedded C callback
bodies remain. Production and test-reference standalone `.c` files both remain zero.

Latest qualified chunk replaces 53 named and three special object-update callback
identities, migrating their consumers and direct-call fixtures (56 exports total).
Four production cgo files disappear. See [UPDATE_IDENTITIES.md](docs/porting/UPDATE_IDENTITIES.md).

Continue chunk-by-chunk with one Luna helper, qualification, commit/push and
recorded reversible decisions. Active batch: 53 sustained-spell callback identities.
Scope/getter/owner evidence is under `build/port-after-update-identities/`. Three
original-behavior fixture files are installed; all 85 roots pass in all three
profiles without skips under `build/port-duration-identities/baseline-original/`.
See [DURATION_IDENTITIES.md](docs/porting/DURATION_IDENTITIES.md). Conversion
is not installed; its bounded Luna draft and primary API contract are ignored. Stop at the milestone or a substantial question.

Latest artifacts: `build/port-update-identities/`.

## What remains

Counts below describe the qualified update identity conversion. Zero `.c` lines is
not a count of all C dependencies or a measure of remaining engineering effort.

| Area | Remaining work or dependency |
| --- | --- |
| Embedded C callback glue | 77 production function bodies in Go preambles: 76 generic function-pointer dispatchers and one specialized adapter. |
| Callback routes | Some Go implementations still call each other through C-compatible addresses. More direct Go dispatch is possible; shared raw fallbacks remain until their users and compatibility requirements are resolved. |
| Declarations and C types | 157 tracked headers / 3,471 physical lines; each production profile selects 198 cgo files in three project packages (alloc, ccall, legacy). Selected-build counts replace the earlier whole-tree text count. These are mostly interface/layout machinery, not unported algorithms. |
| Memory and layout | C-heap allocation, raw pointers, fixed offsets and 32-bit address assumptions remain. Removing them requires ownership/layout changes beyond function translation. |
| External libraries | SDL2, OpenGL, OpenAL and similar native dependencies and their cgo bindings stay for this phase; their future is a subsequent discussion. |
| Portability and release validation | Qualified target is Linux 386/SSE2 with cgo. The complete engine is not qualified as cgo-free, 64-bit, native macOS or browser/WebAssembly. Physical display and audible playback remain manual release checks. |

See [C_LOC.md](docs/porting/C_LOC.md) for the exact standalone-line metric and
history, and [TYPED_CALLBACK_ADAPTERS.md](docs/porting/TYPED_CALLBACK_ADAPTERS.md)
for the remaining embedded callback functions. The earlier catalog/effect batch retired one unused specialized body; shared raw fallback
machinery remains until its live callers are migrated.

## Latest qualification and evidence

- All 279 focused roots pass in each default/server/highres profile, no skips.
- Preceding item full default corpus: 2,441 passes and one established diagnostic skip.
- Safe/static and three production/ABI checks pass; all 56 retired exports absent.
- Headless character creation and explicit save/load/resume pass.
- Full asset suite matches known results: 304 failure events; 17 passing, two
  failing and 32 skipped packages. Original 1,654 asset hashes are unchanged.
- Accepted phases have identical source fingerprints. All 29 changed/new/deleted
  source files match review; existing assertions and captures are unchanged.

Report: [UPDATE_IDENTITIES.md](docs/porting/UPDATE_IDENTITIES.md).
Evidence: [qualification](docs/porting/update-identities-qualification.json).
Known-suite expectation: [mp3-go-expected-suite.jsonl](docs/porting/mp3-go-expected-suite.jsonl).

## Goal, next work and open review items

The [immediate goal](PORT.md#goal-and-target) is removal of the engine's internal
C glue, retaining external native-library bindings and x86/32-bit assumptions.
Whole-build `CGO_ENABLED=0` is deferred for subsequent discussion.
Follow [INTERNAL_C_GLUE.md](docs/porting/INTERNAL_C_GLUE.md) for the dependency
removal order and completion criteria. Client rendering/audio backend replacement
is outside this phase.

The dependency inventory tool is `tools/porting/cgo_inventory.py`; the current
qualified inventory is [update-identities-inventory-after.json](docs/porting/update-identities-inventory-after.json).
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

Magic-missile expiry now uses collision result-aware dispatch. Native identity
keys never enter its raw C fallback; unknown C callbacks retain exact return bits.
See [collision identity decisions](docs/porting/COLLISION_IDENTITIES.md).
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

Latest local artifacts are under `build/port-update-identities/`:
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
| Completed update scenario assets | Removed 1,654 verified original-asset duplicates; 559,837,184 allocated bytes reclaimed. Originals, saves/results retained. Restore: `python3 build/port-artifact-cleanup/restore-recent-scenario.py build/baseline/runs/update-identities-save/deduplicated-assets.json`. Plan/result: `build/port-post-update-cleanup/`. |
| Superseded item/transfer binaries | Removed 14 verified unused executables; 785,170,432 allocated bytes reclaimed. Rebuild revisions `00044175` and `1b14d1a3`; qualified update replacements remain. Source/hash/host-use evidence: `build/port-duration-identities/cleanup-production-{approved.json,deleted.jsonl}`. |
| Obsolete pre-transfer root/legacy Go cache | Removed 56 hash/stat/host-use-verified archives older than qualified item commit `00044175`; 2,075,226,112 allocated bytes reclaimed. Current transfer/update caches and all source/assets remain. Rebuild normally; `build/port-update-identities/cache-cleanup-{proposal.json,approved.json,deleted.jsonl}`. |
| Superseded create/init and damage binaries | Removed 14 verified unused executables; 786,116,608 allocated bytes reclaimed. Rebuild revisions `ed82a5a8` and `8c9b3fc6`; item/transfer replacements remain. Source/hash/host-use evidence: `build/port-update-identities/cleanup-production-{approved.json,deleted.jsonl}`. |
| Completed transfer scenario assets | Removed 1,654 verified original-asset duplicates; 559,902,720 allocated bytes reclaimed. Originals, saves/results retained. Restore: `python3 build/port-artifact-cleanup/restore-recent-scenario.py build/baseline/runs/xfer-identities-save/deduplicated-assets.json`. Plan/result: `build/port-post-xfer-cleanup/`. |
| Obsolete pre-catalog-conversion project cache | 44 unused, hash/stat-verified root/legacy archives removed; 2,402,185,216 allocated bytes reclaimed. Current caches, baselines and binaries remain. Rebuild normally; plan/approval/journal: `build/port-collision-identities/cache-cleanup-{proposal.json,approved.json,deleted.jsonl}`. |
| Completed catalog/effect scenario assets | Removed 1,654 verified duplicate originals; 559,931,392 allocated bytes reclaimed. Originals, saves/results retained. Restore: `python3 build/port-artifact-cleanup/restore-recent-scenario.py build/baseline/runs/catalog-effect-owners-save/deduplicated-assets.json`. Plan/result: `build/port-post-catalog-effect-cleanup/`. |
| Superseded inventory/resource, shop/trade and spell/reward binaries | Removed 21 verified, unused executables; 1,183,252,480 allocated bytes reclaimed. Rebuild revisions `9dac58bf`, `279e7867`, `b5d7ee16` with retained commands. Latest map-room/painting replacements remain. Source/hash/host-use checks and deletion journal: `build/port-catalog-effect-owners/cleanup-{approved.json,deleted.jsonl}`. |
| Completed map-room/painting scenario assets | Removed 1,654 verified original-asset duplicates; 559,882,240 allocated bytes reclaimed. Originals, saves and results remain. Restore with `python3 build/port-artifact-cleanup/restore-recent-scenario.py build/baseline/runs/map-room-paint-owners-save/deduplicated-assets.json`. Plan/result: `build/port-post-map-room-paint-cleanup/`. |
| Completed shop/trade and spell/reward scenario assets | Removed 3,308 verified original-asset duplicates; 1,119,744,000 allocated bytes reclaimed. Saves/results and originals remain. Restore each with `python3 build/port-artifact-cleanup/restore-recent-scenario.py build/baseline/runs/SCENARIO/deduplicated-assets.json` for `shop-trade-owners-save` or `spell-reward-owners-save`. Plan/result: `build/port-post-spell-reward-cleanup/`. |
| Obsolete pre-shop/trade Go cache archives | Removed 37 regular root/legacy archives after hash/stat/cutoff and host-use checks; 2,235,387,904 allocated bytes reclaimed. Rebuild normally; source, binaries and current cache retained. Plan/journal: `build/port-spell-reward-owners/cache-cleanup-{approved.json,deleted.jsonl}`. |
| Completed inventory/resource scenario assets | Removed 1,654 verified duplicate assets; 559,857,664 allocated bytes reclaimed. Saves/results and originals retained. Restore with `python3 build/port-artifact-cleanup/restore-recent-scenario.py build/baseline/runs/inventory-resource-owners-save/deduplicated-assets.json`. Plan/result: `build/port-post-inventory-resource-cleanup/`. |
| Completed native-owner/object-state scenario assets | Removed 3,308 verified original-asset duplicates; 1,119,727,616 allocated bytes reclaimed. Saves/results remain. Restore each with `python3 build/port-artifact-cleanup/restore-recent-scenario.py build/baseline/runs/SCENARIO/deduplicated-assets.json`, using `go-native-owner-constants-save` or `object-state-owners-save`. Plan/result: `build/port-post-object-state-cleanup/`. |
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
| Superseded six-batch profile test binaries | Removed 18 regular executables after exact committed-source, SHA256/stat and host fd/maps checks; 1,219,395,584 allocated bytes reclaimed. Rebuild qualified revisions `6274a6f3`, `b505e0ab`, `fe44bab3`, `f2088976`, `5bfa54b5`, `200d459b` using retained profile commands/source maps. Collision corrected contracts and original baseline retained; first-failure executables were subsequently removed as recorded below. Plan/journal: `build/port-collision-identities/cleanup-next-{approved.json,deleted.jsonl}`. |
| Completed collision scenario assets | Removed 1,654 SHA256-identical original-asset duplicates after host-use checks, reclaiming 559,841,280 allocated bytes; saves/results and originals remain. Restore with `python3 build/port-artifact-cleanup/restore-recent-scenario.py build/baseline/runs/collision-identities-save/deduplicated-assets.json`. |
| Superseded six-batch safe/production binaries | Removed 24 executables after exact committed-source/replacement SHA256 and host fd/maps checks; 1,149,173,760 allocated bytes reclaimed. Rebuild revisions `6274a6f3`, `b505e0ab`, `fe44bab3`, `f2088976`, `5bfa54b5`, `200d459b` using retained phase commands. Qualified collision replacements and active death outputs retained. Plan/journal: `build/port-death-identities/cleanup-production-{approved.json,deleted.jsonl}`. |
| Superseded collision/death first-failure test binaries | Removed six executables after recorded SHA256, committed qualification/replacement and host-use checks; 405,254,144 allocated bytes reclaimed. Failed logs/source maps and corrected qualified binaries remain. Plan/journal: `build/port-create-init-identities/cleanup-failed-{approved.json,deleted.jsonl}`. |
| Completed death scenario assets | Removed 1,654 SHA256-identical original-asset duplicates after host-use checks; 559,845,376 allocated bytes reclaimed. Originals, saves and results retained. Restore with `python3 build/port-artifact-cleanup/restore-recent-scenario.py build/baseline/runs/death-identities-save/deduplicated-assets.json`. Plan/result: `build/port-post-death-cleanup/`. |
| Obsolete pre-death-baseline Go cache | Removed 71 hash/stat-verified, unused root/legacy archives older than 2026-09-24T18:29:06Z; 2,416,160,768 allocated bytes reclaimed. Newer death/create-init caches, all binaries/source/assets retained. Rebuild normally. Plan/journal: `build/port-create-init-identities/cache-cleanup-{proposal.json,approved.json,deleted.jsonl}`. |
| Obsolete pre-damage project cache | Removed 62 verified root/legacy archives older than qualified create/init commit; 2,403,794,944 allocated bytes reclaimed after host checks and test completion. Current damage/item caches remain; rebuild old artifacts normally. Plan/journal: `build/port-item-identities/cache-cleanup-{approved.json,deleted.jsonl}`. |
| Completed item scenario assets | Removed 1,654 verified original-asset duplicates; 559,865,856 allocated bytes reclaimed. Saves/results and originals remain. Restore with `python3 build/port-artifact-cleanup/restore-recent-scenario.py build/baseline/runs/item-identities-save/deduplicated-assets.json`. Plan/result: `build/port-post-item-cleanup/`. |
| Completed damage scenario assets | Removed 1,654 verified duplicates; 559,915,008 allocated bytes reclaimed. Originals, saves and results remain. Restore with `python3 build/port-artifact-cleanup/restore-recent-scenario.py build/baseline/runs/damage-identities-save/deduplicated-assets.json`. Plan/result: `build/port-post-damage-cleanup/`. |
| Completed creation/init scenario assets | Removed 1,654 SHA256-identical original-asset duplicates after host-use checks; 559,972,352 allocated bytes reclaimed. Originals, saves/results remain. Restore with `python3 build/port-artifact-cleanup/restore-recent-scenario.py build/baseline/runs/create-init-identities-save/deduplicated-assets.json`. Plan/result: `build/port-post-create-init-cleanup/`. |
| Superseded collision/death qualified binaries | Removed 14 test/safe/production executables after exact committed-source, replacement/hash and host-use checks; 786,751,488 allocated bytes reclaimed. Rebuild qualified revisions `170b594a` and `e7ca9d31` using retained commands. Source, original baseline binaries, logs, manifests and current create/init and damage outputs remain. Plan/journal: `build/port-item-identities/cleanup-production-{approved.json,deleted.jsonl}`. |
| Current qualified production/safe binaries | Retained under `build/port-update-identities/`; superseded transfer/item binaries were removed as recorded above; superseded damage/create-init binaries were removed as recorded above; older collision/death qualified executables were removed as recorded above. |
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
