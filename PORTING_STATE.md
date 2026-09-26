# Porting checkpoint

This file is the current resume checkpoint, not a chronological log. Replace
superseded status when updating it. The workflow and delegation rules live in
[PORT.md](PORT.md); detailed evidence belongs in the linked batch reports.

## Status: resumed; internal C-glue removal

**Progress: 142,665/142,665 original standalone C lines ported or retired;
internal glue: 344/463 client cgo files eliminated on net (119 remain;
server: 343/463 eliminated, 120 remain).**
Selected legacy C export bridges: **1,721/1,890 retired (169 remain)**.

These are selected project files in Linux 386 production profiles, not equal
units of effort. Three project packages directly use cgo; 77 embedded C callback
bodies remain. Production and test-reference standalone `.c` files both remain zero.

Latest qualified chunk retires five unused GUI exports, their C-only conversion
helpers and a private C color return type. It removes seven production cgo imports,
including two leftover includes-only imports. Native GUI owners, live callbacks,
shared layouts and compiler settings are preserved. See
[GUI_ADAPTERS.md](docs/porting/GUI_ADAPTERS.md).

Continue chunk-by-chunk with one Luna helper, primary review, qualification,
documentation, commit/push and recorded reversible decisions. Stop at the milestone,
usage limits or a substantial question.

Latest qualified artifacts: `build/port-gui-adapters/`.

## What remains

Counts below describe the qualified GUI-adapter conversion. Zero `.c` lines
is not a count of all C dependencies or remaining engineering effort.

| Area | Remaining work or dependency |
| --- | --- |
| Embedded C callback glue | 77 production function bodies in Go preambles: 76 generic function-pointer dispatchers and one specialized adapter. |
| Callback routes | Remaining Go owners still use C-compatible addresses. Continue migrating identities and every field/alias consumer before removing shared raw fallbacks. |
| Declarations and C types | 157 tracked headers / 2,924 physical lines; client profiles select 119 cgo files and server selects 120 in three project packages (alloc, ccall, legacy). These are mostly interface/layout machinery, not unported algorithms. |
| Memory and layout | C-heap allocation, raw pointers, fixed offsets and 32-bit address assumptions remain. Ownership/lifetime work remains behind the centralized allocator. |
| External libraries | SDL2, OpenGL, OpenAL and similar native dependencies/bindings stay for this phase; their future is a subsequent discussion. |
| Portability and release validation | Qualified target is Linux 386/SSE2 with cgo. The complete engine is not qualified as cgo-free, 64-bit, native macOS or browser/WebAssembly. Physical display and audible playback remain manual checks. |

See [C_LOC.md](docs/porting/C_LOC.md) for the standalone-line metric and history,
and [TYPED_CALLBACK_ADAPTERS.md](docs/porting/TYPED_CALLBACK_ADAPTERS.md) for the
remaining embedded callbacks. Shared fallback machinery remains until its live
users are migrated; test-only C observers separately qualify that boundary.

## Latest qualification and evidence

- All 36 affected roots pass in each of three profiles with exact original names
  and no failures/skips. The baseline reused 26 meter roots from exact qualified
  `46f07aba` source/environment and ran ten additional GUI-owner roots twice per
  profile. The complete union was rerun after conversion.
- Safe/static and three production/ABI checks pass; all five retired exports absent.
- Fresh preview and final headless character creation with explicit save/load/
  resume pass on the same final production binary.
- Full asset suite matches known results: 304 failure events; 17 passing, two
  failing and 32 skipped packages. All 1,654 original asset hashes are unchanged.
- Accepted phases share source fingerprints and primary-reviewed changes.
  Existing independent assertions and state/pixel captures remain unchanged.
- Last full default corpus was drawable-update `024b2632`: 2,461 passes and one
  established prerequisite skip among 2,462 roots. That is an earlier-source result;
  this callback-interface batch uses affected coverage plus production gates.

Report: [GUI_ADAPTERS.md](docs/porting/GUI_ADAPTERS.md).
Evidence: [qualification](docs/porting/gui-adapters-qualification.json).
Known-suite expectation: [mp3-go-expected-suite.jsonl](docs/porting/mp3-go-expected-suite.jsonl).

## Goal, next work and open review items

The [immediate goal](PORT.md#goal-and-target) is removal of the engine's internal
C glue, retaining external native-library bindings and x86/32-bit assumptions.
Whole-build `CGO_ENABLED=0` is deferred for subsequent discussion.
Follow [INTERNAL_C_GLUE.md](docs/porting/INTERNAL_C_GLUE.md) for the dependency
removal order and completion criteria. Client rendering/audio backend replacement
is outside this phase.

The dependency inventory tool is `tools/porting/cgo_inventory.py`; the current
qualified inventory is [gui-adapters-inventory-after.json](docs/porting/gui-adapters-inventory-after.json).
The original phase baseline is under `build/port-cgo-leaves/inventory-before/`.
The completed leaf cleanup leaves three project packages directly using cgo in
all profiles, plus OpenGL/SDL2/OpenAL bindings in the clients. Metadata discovery
is not compilation or qualification. The helper's external-review draft is not
accepted evidence: its suggestion that go-gl is residue is contradicted by the
actual dependency graph (`libs/client/seat/opengl` imports it).

GUI-adapter conversion is qualified. Next audit the connected animation callback
fields, their getter aliases and saved-pointer consumers across menus, character
creation, options and browser owners. Preserve the 68-byte animation layout and
foreign callback behavior. The scout is under `build/port-after-gui-scout/`;
no later conversion is installed or qualified.

The earlier 25-export object-state proposal under
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

Latest local artifacts are under `build/port-gui-adapters/`:
`contracts/`, `preview/`, `safe/opennox-safe`, and
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
| GUI-adapter cleanup | Removed 9 obsolete project cache archives (420,384,768 allocated bytes), four superseded window production/safe binaries (187,895,808 bytes; rebuild `3838462b`), and three superseded rendering test binaries (200,912,896 bytes; rebuild `46f07aba`). Current replacements, source and logs retained. Journals under `build/port-gui-adapters/`. Preview/final each removed 1,654 verified asset copies; original assets, saves and results retained. Restore with `build/port-artifact-cleanup/restore-recent-scenario.py` and `build/baseline/runs/gui-adapters[-preview]-save/deduplicated-assets.json`. |
| Rendering/image cleanup | Removed 9 obsolete project cache archives (420,610,048 allocated bytes), four superseded quantity production/safe binaries (188,002,304 bytes; rebuild `5585ab65`), and three superseded window test binaries (200,974,336 bytes; rebuild `3838462b`). Current replacements, source and logs retained. Journals under `build/port-render-image-bridges/`. Preview/final each removed 1,654 verified asset copies; original assets, saves and results retained. Restore with `build/port-artifact-cleanup/restore-recent-scenario.py` and `build/baseline/runs/render-image-bridges[-preview]-save/deduplicated-assets.json`. |
| Window bridge cleanup | Removed 9 unused project cache archives older than `5585ab65` (420,876,288 allocated bytes), four superseded inventory production/safe binaries (188,059,648 bytes; rebuild `785c2d54`), and three superseded quantity test binaries (201,060,352 bytes; rebuild `5585ab65`). Current replacements and all source/logs retained. Plans/journals under `build/port-window-bridges/`. Preview/final scenarios each removed 1,654 verified asset copies; original assets and saves/results retained. Restore each with `build/port-artifact-cleanup/restore-recent-scenario.py` and `build/baseline/runs/window-bridges[-preview]-save/deduplicated-assets.json`. |
| Superseded audio/UI/inventory previews | Removed three hash/source/host-use-verified preview executables (142,897,152 allocated bytes); rebuild `9eac418e`, `541585fa` or `785c2d54` from retained source/commands. Qualified quantity and current window replacements remain. Proposal and journal: `build/port-window-bridges/old-preview-cleanup/`. |
| Quantity callback cleanup | Removed nine unused project cache archives older than `785c2d54` (421,044,224 allocated bytes) and four superseded UI-meter production/safe binaries (188,129,280 bytes; rebuild `541585fa`). Also removed three superseded inventory test binaries (201,101,312 bytes; rebuild `785c2d54`) after their replacements passed. Current replacements and source/logs retained. Plans/journals: `build/port-quantity-identities/cache-before-conversion/` and `old-production-cleanup-*` / `original-tests-cleanup-*`. Preview/final scenarios each removed 1,654 verified asset copies after qualification; originals and saves/results retained. Restore using `build/port-artifact-cleanup/restore-recent-scenario.py` with each `build/baseline/runs/quantity-identities[-preview]-save/deduplicated-assets.json`. |
| Inventory GUI preview duplicate assets | Removed 1,654 verified copies (559,869,952 allocated bytes); originals and saves/results remain. Restore using `python3 build/port-artifact-cleanup/restore-recent-scenario.py build/baseline/runs/inventory-gui-identities-preview-save/deduplicated-assets.json`. |
| Inventory GUI final duplicate assets | Removed 1,654 verified copies (559,960,064 allocated bytes); originals and saves/results remain. Restore using `python3 build/port-artifact-cleanup/restore-recent-scenario.py build/baseline/runs/inventory-gui-identities-save/deduplicated-assets.json`. |
| Superseded UI meter and inventory test-failure binaries | Removed six verified executables (402,345,984 allocated bytes) after corrected baseline replacement and host-use checks. Rebuild UI meter `541585fa`; the first geometry-failure test source is saved in `build/port-inventory-gui-identities/primary/failed-geometry-contract.go`. Current converted binaries and all records/logs/source remain; original baseline retention is described separately. Journal: `build/port-inventory-gui-identities/old-tests-cleanup-{approved.json,deleted.jsonl}`. |
| Superseded original inventory binaries and cache | Removed three original test executables (201,183,232 allocated bytes) after converted replacement passed; rebuild original baseline `44a780d9`. Also removed eight unused older project archives (589,676,544 bytes). Current converted tests and production binaries, all source/log records and original assets remain. Evidence: `build/port-inventory-gui-identities/original-tests-cleanup-{approved.json,deleted.jsonl}` and `cache-before-production/`. |
| Inventory baseline cleanup | Removed six superseded UI meter / first geometry-failure test binaries (402,345,984 allocated bytes) after accepted replacement and host checks. Removed seven unused project archives predating `541585fa` (339,312,640 bytes). Current corrected inventory baseline binaries, original assets and all source/log records remain. Rebuild UI meter `541585fa`; failed new-test source is `build/port-inventory-gui-identities/primary/failed-geometry-contract.go`. Journals: `old-tests-cleanup-{approved.json,deleted.jsonl}` and `cache-before-conversion/` under that batch directory. |
| Pre-inventory disk cleanup | Removed five unused project cache archives predating original UI baseline `4792e80b` (335,773,696 allocated bytes), and four superseded audio production/safe binaries (188,211,200 bytes), after source/hash/replacement and host-use checks. Current UI binaries/caches, original assets and all logs/source records remain. Rebuild audio `9eac418e` as needed. Evidence: `build/port-inventory-gui-identities/cache-before-inventory/` and `old-production-cleanup-{approved.json,deleted.jsonl}`. |
| UI meter preview/final duplicate assets | Removed 1,654 verified copies per scenario (559,853,568 / 559,964,160 allocated bytes). Originals, saves/results remain. Restore using `python3 build/port-artifact-cleanup/restore-recent-scenario.py` with `build/baseline/runs/ui-meter-identities-preview-save/deduplicated-assets.json` or `build/baseline/runs/ui-meter-identities-save/deduplicated-assets.json`. |
| Obsolete pre-qualified-audio cache | Removed nine unused project archives predating `9eac418e` (421,228,544 allocated bytes), after hash/stat/source-family and host-use checks. Current UI baseline caches, binaries, source and assets preserved at cleanup. Evidence: `build/port-ui-meter-identities/cache-before-conversion/`. |
| Superseded audio and original UI test binaries | Removed six source/hash/host-verified executables (402,395,136 allocated bytes), retaining current converted UI contract binaries. Rebuild audio `9eac418e` or original UI `4792e80b` normally; source/reports/logs remain. Evidence: `build/port-ui-meter-identities/old-tests-cleanup-{approved.json,deleted.jsonl}`. |
| Original UI meter baseline | Committed as `4792e80b`; 363 client/360 server roots plus separate-process repeats. Records/logs remain under `build/port-ui-meter-identities/`; executable retention follows the superseded-binary row above. |
| Obsolete pre-original-audio cache | Removed 14 unused root/legacy archives predating `403efa06`, reclaiming 757,334,016 allocated bytes after source-family/hash/stat and host-use checks. Current audio caches/binaries, all source and assets preserved. Evidence: `build/port-ui-meter-identities/cache-before-meters/`. Cleanup script consumed. |
| Superseded list production/safe and older preview binaries | Removed four list production/safe executables (188,317,696 allocated bytes) after current audio builds/ABI and host-use checks, plus three monster/player/list preview executables (143,073,280 bytes) after source/hash/replacement checks. All logs/reports/source remain; rebuild the recorded revisions normally. Current audio production/safe/preview retained. Evidence: `build/port-audio-bridge-identities/old-{production,preview}-cleanup-{approved.json,deleted.jsonl}`. |
| Audio preview/final duplicate assets | Removed 1,654 verified copies per scenario (559,976,448 / 560,005,120 allocated bytes). Originals, saves/results remain. Restore using `python3 build/port-artifact-cleanup/restore-recent-scenario.py` with `build/baseline/runs/audio-bridge-identities-preview-save/deduplicated-assets.json` or `build/baseline/runs/audio-bridge-identities-save/deduplicated-assets.json`. |
| Superseded list and original audio test binaries | Removed six source/hash/host-verified executables (402,460,672 allocated bytes); converted audio contracts retained. Rebuild original audio `403efa06` or list `6f981385` normally; logs/records/source remain. Evidence: `build/port-audio-bridge-identities/old-tests-cleanup-{approved.json,deleted.jsonl}`. |
| Superseded player-action binaries and older cache archives | Removed seven verified unused player binaries (389,644,288 allocated bytes) and nine older root/legacy cache archives (421,584,896 bytes). Source/hash/replacement and host-use checks passed; current list and original audio binaries remain. Rebuild player `1ca13c2a` normally. Evidence: `build/port-audio-bridge-identities/binary-cleanup-{approved.json,deleted.jsonl}` and `cache-before-audio/`. |
| Completed list preview/final scenario assets | Removed 1,654 identical asset copies per run, reclaiming 559,976,448 / 560,001,024 allocated bytes. Originals, saves and results remain. Restore with `python3 build/port-artifact-cleanup/restore-recent-scenario.py` followed by `build/baseline/runs/list-glue-identities-preview-save/deduplicated-assets.json` or `build/baseline/runs/list-glue-identities-save/deduplicated-assets.json`. Evidence: `build/port-list-glue-identities/{preview,final}-cleanup/`. |
| Superseded monster binaries and old project caches | Removed seven verified unused monster executables (389,943,296 allocated bytes) and 14 root/legacy cache archives predating `44a9a065` (758,870,016 bytes), with committed source/replacement/hash and host-use checks. Rebuild monster `456bdf2d` and older caches normally. Current player/list baseline binaries, current caches, source and assets remain. Evidence: `build/port-list-glue-identities/binary-cleanup-{approved.json,deleted.jsonl}` and `cache-before-list/`. |
| Completed player-action preview/final scenario assets | Removed 1,654 identical asset copies per run, reclaiming 559,988,736 / 559,992,832 allocated bytes. Originals, saves and results remain. Restore with `python3 build/port-artifact-cleanup/restore-recent-scenario.py` followed by `build/baseline/runs/player-action-identities-preview-save/deduplicated-assets.json` or `build/baseline/runs/player-action-identities-save/deduplicated-assets.json`. Evidence: `build/port-player-action-identities/{preview,final}-cleanup/`. |
| Superseded original monster test binaries | Removed three unused original-probe executables; 201,441,280 allocated bytes reclaimed after committed source/hash/replacement and host-use checks. Rebuild original `28642775`; its baseline records/logs and qualified `456bdf2d` replacement binaries remain. Evidence: `build/port-player-action-identities/cleanup-original-probe-{approved.json,deleted.jsonl}`. |
| Superseded drawable-update binaries | Removed seven unused, source/hash/host-verified executables; 390,037,504 allocated bytes reclaimed. Rebuild `024b2632` with retained commands/source maps. Current qualified monster and player baseline binaries remain. Evidence: `build/port-player-action-identities/binary-cleanup-{approved.json,deleted.jsonl}`. |
| Completed monster preview/final scenario assets | Removed 1,654 identical asset copies per run, reclaiming 559,919,104 / 559,996,928 allocated bytes. Originals, saves and results remain. Restore with `python3 build/port-artifact-cleanup/restore-recent-scenario.py` followed by `build/baseline/runs/monster-callback-identities-preview-save/deduplicated-assets.json` or `build/baseline/runs/monster-callback-identities-save/deduplicated-assets.json`. Evidence: `build/port-monster-callback-identities/{preview,final}-cleanup/`. |
| Superseded remaining-draw binaries | Removed seven unused, hash/stat/committed-source-verified executables; 390,156,288 allocated bytes reclaimed. Rebuild `b03fb880` with retained phase commands/source maps. Qualified drawable-update replacements and monster originals remain. Evidence: `build/port-monster-callback-identities/binary-cleanup-{approved.json,deleted.jsonl}`. |
| Superseded pre-monster project cache | Removed nine unused, hash/stat/host-verified root/legacy archives predating `024b2632`; 422,182,912 allocated bytes reclaimed. Current monster baseline caches and all binaries/source/assets remain. Rebuild older cache entries normally. Evidence: `build/port-monster-callback-identities/cache-before-monster/`. |
| Completed drawable-update final scenario assets | Removed 1,654 verified duplicates, reclaiming 559,894,528 allocated bytes. Originals, saves/results retained. Restore: `python3 build/port-artifact-cleanup/restore-recent-scenario.py build/baseline/runs/drawable-update-identities-save/deduplicated-assets.json`. Evidence: `build/port-drawable-update-identities/final-cleanup/`. |
| Completed server-fixture scenario assets | Removed 1,654 verified duplicates, reclaiming 559,919,104 allocated bytes. Originals, saves/results retained. Restore: `python3 build/port-artifact-cleanup/restore-recent-scenario.py build/baseline/runs/server-fixture-bridges-save/deduplicated-assets.json`. Evidence: `build/port-post-server-fixture-cleanup/`. |
| Obsolete pre-server-fixture-baseline cache | Removed nine hash/stat/host-use-verified root/legacy archives older than 9a833ff4; 425,369,600 allocated bytes reclaimed. Current qualified caches, all binaries/source/assets retained. Evidence: `build/port-client-ui-fixture-bridges/cache-cleanup-{proposal.json,approved.json,deleted.jsonl}`. |
| Completed modifier scenario assets | Removed 1,654 verified duplicates, reclaiming 559,865,856 allocated bytes. Originals, saves/results retained. Restore: `python3 build/port-artifact-cleanup/restore-recent-scenario.py build/baseline/runs/modifier-identities-save/deduplicated-assets.json`. Evidence: `build/port-post-modifier-cleanup/`. |
| Obsolete pre-modifier-baseline cache | Removed five hash/stat/host-use-verified root/legacy archives older than 09f15464; 340,205,568 allocated bytes reclaimed. Qualified modifier conversion caches, all binaries/source/assets retained. Evidence: `build/port-server-fixture-bridges/cache-cleanup-{proposal.json,approved.json,deleted.jsonl}`. |
| Completed duration scenario assets | Removed 1,654 verified original-asset duplicates; 559,857,664 allocated bytes reclaimed. Originals, saves/results retained. Restore: `python3 build/port-artifact-cleanup/restore-recent-scenario.py build/baseline/runs/duration-identities-save/deduplicated-assets.json`. Plan/result: `build/port-post-duration-cleanup/`. |
| Obsolete pre-modifier root/legacy Go cache | Removed 14 hash/stat/host-use-verified archives older than qualified duration commit d5d80c42; 764,968,960 allocated bytes reclaimed. Current modifier caches and all source/assets/binaries retained. Rebuild normally; `build/port-modifier-identities/cache-cleanup-{proposal.json,approved.json,deleted.jsonl}`. |
| Obsolete pre-duration root/legacy Go cache | Removed 29 hash/stat/host-use-verified archives older than qualified update commit (2026-09-24T21:12:00Z); 1,695,002,624 allocated bytes reclaimed. Current duration caches, source/assets/binaries retained. Rebuild normally; `build/port-duration-identities/cache-cleanup-{proposal.json,approved.json,deleted.jsonl}`. |
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
| Obsolete pre-audio project cache | Removed 21 hash/stat-verified root/legacy Linux 386 archives older than original audio baseline `9f6b2b46`, after host-use checks; 1,100,709,888 allocated bytes reclaimed. Newer audio caches, all source/assets/binaries remain. Rebuild normally. Plan/journal: `build/port-after-audio/cache-luna/cache-cleanup-{approved.json,deleted.jsonl}`. |
| Current qualified production/safe binaries | Retained under `build/port-gui-adapters/`; superseded outputs removed as recorded below. |
| Superseded UI-fixture binaries | Removed seven verified test/safe/production executables after committed audio replacement and host-use checks; 390,643,712 allocated bytes reclaimed. Rebuild revision `3d47a346` using retained commands/source maps. Logs/manifests and current audio replacements remain. Plan/journal: `build/port-audio-stream-callbacks/binary-cleanup-{approved.json,deleted.jsonl}`. |
| Completed audio-stream scenario assets | Removed 1,654 verified original-asset duplicates; 559,931,392 allocated bytes reclaimed. Originals, saves and results remain. Restore with `python3 build/port-artifact-cleanup/restore-recent-scenario.py build/baseline/runs/audio-stream-callbacks-save/deduplicated-assets.json`. Plan/result: `build/port-post-audio-stream-cleanup/`. |
| Superseded modifier/server-fixture/duration/update binaries | Removed 28 verified test/safe/production executables after source/replacement hashes and host-use checks; 1,567,293,440 allocated bytes reclaimed. Rebuild qualified revisions `99b65896`, `a798ad1c`, `d5d80c42`, `a8d89bda` using retained commands and source maps. Current UI replacements, old logs/manifests and baseline evidence remain. Journals: `build/port-after-client-ui/binary-cleanup-deleted.jsonl` and `binary-cleanup-addendum-deleted.jsonl`. |
| Completed UI-fixture scenario assets | Removed 1,654 verified original-asset duplicates; 559,943,680 allocated bytes reclaimed. Originals, saves and results remain. Restore with `python3 build/port-artifact-cleanup/restore-recent-scenario.py build/baseline/runs/client-ui-fixture-bridges-save/deduplicated-assets.json`. Plan/result: `build/port-post-client-ui-cleanup/`. |
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

Ten completed UI/audio-baseline logs are losslessly compressed; restore with
`gzip -dk`. Hashes and commands: `build/port-audio-stream-callbacks/completed-log-archive.json`
(107,732,992 allocated bytes reclaimed). Current audio qualification logs remain readable.

The current drawing conversion's completed preview asset copy is deduplicated:
1,654 identical files / 559,869,952 allocated bytes reclaimed. Original assets,
saves and logs remain. Restore with
`python3 build/port-artifact-cleanup/restore-recent-scenario.py build/baseline/runs/client-draw-identities-preview-save/deduplicated-assets.json`.
Verification and deletion journal: `build/port-client-draw-identities/preview-cleanup/`.

Seven superseded audio executables were removed after committed replacement,
source/hash and host-use checks (390,561,792 allocated bytes). Rebuild `c22c0edc`
using retained commands and source maps. Evidence:
`build/port-client-draw-identities/binary-cleanup-{approved.json,deleted.jsonl}`.
Current drawing binaries remain. The final drawing scenario also had 1,654 verified
duplicate assets removed (559,861,760 allocated bytes); saves/results and originals
remain. Restore:
`python3 build/port-artifact-cleanup/restore-recent-scenario.py build/baseline/runs/client-draw-identities-save/deduplicated-assets.json`.
Cleanup evidence: `build/port-post-client-draw-cleanup/`.

The completed remaining-draw preview had 1,654 verified duplicate assets removed
(559,861,760 allocated bytes); originals, saves and results remain. Restore:
`python3 build/port-artifact-cleanup/restore-recent-scenario.py build/baseline/runs/remaining-draw-identities-preview-save/deduplicated-assets.json`.
Evidence: `build/port-remaining-draw-identities/preview-cleanup/`.

The completed remaining-draw final scenario also had 1,654 verified duplicate
assets removed (559,943,680 allocated bytes). Restore:
`python3 build/port-artifact-cleanup/restore-recent-scenario.py build/baseline/runs/remaining-draw-identities-save/deduplicated-assets.json`.
Evidence: `build/port-update-identities/drawing-cleanup/`.

Removed 21 verified root/legacy Linux 386 cache archives older than baseline
`9f91acda`, reclaiming 1,099,063,296 allocated bytes. Host checks confirmed no
compiler or target open-file/mapping use; active Go processes only wrapped
prebuilt tests. Current source, binaries and assets remain; rebuild caches normally.
Evidence: `build/port-update-identities/cache-cleanup/`.

Seven superseded drawing executables from `5687bc68` were removed after committed
source, qualified replacement and host-use checks (390,184,960 allocated bytes).
Rebuild using retained commands/source maps; current qualified remaining-draw
binaries and original baseline binaries remain. Evidence:
`build/port-drawable-update-identities/binary-cleanup-{approved.json,deleted.jsonl}`.

Removed 17 more verified root/legacy Linux 386 cache archives older than original
drawable-update baseline `ae17a051`, reclaiming 1,012,924,416 allocated bytes after
host checks. Current conversion caches/binaries and all assets remain; rebuild
old caches normally. Evidence: `build/port-drawable-update-identities/cache-cleanup/`.

The drawable-update preview passed and its 1,654 verified duplicate assets were
removed (559,894,528 allocated bytes). Originals, saves and results remain. Restore:
`python3 build/port-artifact-cleanup/restore-recent-scenario.py build/baseline/runs/drawable-update-identities-preview-save/deduplicated-assets.json`.
Evidence: `build/port-drawable-update-identities/preview-cleanup/`.
