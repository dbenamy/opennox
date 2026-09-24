# Porting checkpoint

This file is the current resume checkpoint, not a chronological log. Replace
superseded status when updating it. The workflow and delegation rules live in
[PORT.md](PORT.md); detailed evidence belongs in the linked batch reports.

## Status: resumed; internal C-glue removal

**Progress: 142,665/142,665 original standalone C lines ported or retired;
internal glue: 66/463 cgo files eliminated on net (397 remain).**
Selected legacy C export bridges: **446/1,890 retired (1,444 remain)**.

The glue count uses selected project files in each Linux 386 production profile.
Directly cgo-dependent project packages are down from six to three; 79 embedded
C callback bodies remain. These are dependency counts, not equivalent units of
work or an effort percentage. Production and test-reference standalone `.c`
files both remain at zero.

Latest qualified implementation: **378 unused C export wrappers and 333 header
prototypes retired**, against frozen baseline `b1aadddd`. Two more files no longer
need cgo; live Go declarations, tests, callback IDs and external native bindings
are unchanged. See [REMAINING_UNUSED_EXPORTS.md](docs/porting/REMAINING_UNUSED_EXPORTS.md).

Continue chunk-by-chunk with one Luna helper, qualification, commit/push and
recorded reversible decisions. Stop at the milestone or for a substantial question.
Next: review and install the following 265-export cohort: remove 126 unused
wrappers and preserve 139 Go-referenced functions while retiring their C exports.
Luna's ignored `build/port-go-only-exports/draft-v2/` is uninstalled; import cleanup,
primary AST review, baseline freeze and qualification remain. The earlier
`draft/` and `reviewed/` directive-only versions are superseded.

The complete `^Test` root corpus passes in all three profiles with exact name
sets and only the documented map-population diagnostic skip. The first complete
sweep with at most two prebuilt processes passed; keep the cap at two and all
production/headless work sequential. See [COMPLETE_PORT_CORPUS.md](docs/porting/COMPLETE_PORT_CORPUS.md)
and [PREBUILT_PROFILES.md](docs/porting/PREBUILT_PROFILES.md).
Latest artifacts: `build/port-remaining-unused-exports/`.


## What remains

Counts below describe the qualified unused-export cleanup. Zero `.c` lines is
not a count of all C dependencies or a measure of remaining engineering effort.

| Area | Remaining work or dependency |
| --- | --- |
| Embedded C callback glue | 79 production function bodies in Go preambles: 76 generic function-pointer dispatchers and three specialized adapters. |
| Callback routes | Some Go implementations still call each other through C-compatible addresses. More direct Go dispatch is possible; shared raw fallbacks remain until their users and compatibility requirements are resolved. |
| Declarations and C types | 157 tracked headers / 4,146 physical lines; each production profile selects 397 cgo files in three project packages (alloc, ccall, legacy). Selected-build counts replace the earlier whole-tree text count. These are mostly interface/layout machinery, not unported algorithms. |
| Memory and layout | C-heap allocation, raw pointers, fixed offsets and 32-bit address assumptions remain. Removing them requires ownership/layout changes beyond function translation. |
| External libraries | SDL2, OpenGL, OpenAL and similar native dependencies and their cgo bindings stay for this phase; their future is a subsequent discussion. |
| Portability and release validation | Qualified target is Linux 386/SSE2 with cgo. The complete engine is not qualified as cgo-free, 64-bit, native macOS or browser/WebAssembly. Physical display and audible playback remain manual release checks. |

See [C_LOC.md](docs/porting/C_LOC.md) for the exact standalone-line metric and
history, and [TYPED_CALLBACK_ADAPTERS.md](docs/porting/TYPED_CALLBACK_ADAPTERS.md)
for the remaining embedded callback functions. Recent dispatch conversions do
not reduce the 79-body count because they retain the shared fallback machinery.

## Latest qualification and evidence

- Complete porttest roots: default/highres each 2,425 pass plus one expected
  diagnostic skip; server 2,414 pass plus that skip. Exact discovered, started and
  completed root-name sets match; no failure events or unexpected skips.
- Safe build/static checks and three fresh production binaries/ABI checks pass;
  all 378 newly retired symbols are absent. No safe owner-contract run is claimed.
- Headless character creation and explicit save/load/resume pass.
- Full-suite results match the known baseline exactly: 304 failure events,
  with 17 passing, two failing and 32 skipped packages.
- All phases used identical source fingerprints. All test/fixture inputs remained
  unchanged; all 1,654 original asset hashes checked against the saved manifest
  remain unchanged. The prior player-reset fixture repair remains qualified.

Report: [REMAINING_UNUSED_EXPORTS.md](docs/porting/REMAINING_UNUSED_EXPORTS.md).
Evidence: [qualification](docs/porting/remaining-unused-exports-qualification.json).
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
Other behavior/compatibility findings are recorded in [DECISIONS.md](docs/porting/DECISIONS.md).

Luna drafted bounded export/prototype removals. Primary corrected broad inventory
omissions, independently checked all 65 symbols and retained live globals. A later
selector audit also needed correction to use the full regex rather than literal
alternatives. Keep inventory algorithm design with the primary and prefer bounded,
mechanically verifiable helper assignments.

## Resume and artifact recovery

Read this checkpoint and PORT.md, then inspect Git status before editing. The
untracked `nox-iso-from-archive-org.7z` is an asset archive, not unfinished code;
never stage it. Preserve it and `build/assets/extracted/drive_c/Nox`.

The current Go toolchain is `/usr/lib/go-1.26/bin`. Follow the
[build environment instructions](PORT.md#build-and-test-environment), including
sourcing `build/baseline/env.sh` in every Go shell.

Latest local artifacts are under `build/port-remaining-unused-exports/`:
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
| Current qualified production/safe binaries | Retained under `build/port-remaining-unused-exports/`; preceding original full-corpus evidence is under `build/port-complete-corpus/`. |
| Completed scenario data: `go-memory-save`, `raw-allocation-save`, `string-boundary-save`, `unused-exports-save` | Only SHA256-identical original-asset copies were removed. Saves/comparisons remain. Follow each run's `deduplicated-assets.json`; shared restore tool: `build/port-artifact-cleanup/restore-recent-scenario.py`. |
| Large historical captures in `port-game-messages`, `port-map-sections`, `port-client-interaction`, `port-session-dialogs` | Restore with `gzip -dk` and verify hashes against `build/port-artifact-cleanup/large-historical-20260924/`. Its 116 discarded text logs are not recoverable from these archives. |
| Initial complete-corpus default/server logs | Losslessly compressed; restore commands and hashes in `build/port-complete-corpus/initial-log-archive.json`. Keep the server failure evidence. |
| Qualified complete-corpus logs | Losslessly compressed; restore commands/hashes: `build/port-complete-corpus/qualified-log-archive.json`. |
| Completed pointer-fixture and prebuilt-pilot binaries | Five rebuildable binaries removed; source `4a0ab0dc`, original logs/records retained. Plan/journal: `build/port-artifact-cleanup/completed-corpus-binaries-{plan.json,removed.jsonl}`. |
| Superseded Go-memory/raw-allocation/string-boundary binaries | Rebuild from `4e1e86a6`, `6ff98033`, `92029ddd` respectively. Exact inventory/removal journal: `build/port-artifact-cleanup/superseded-glue-binaries-{plan.json,removed.jsonl}`. |
| Completed-chunk old project cache archives | 24 Sep23 root/legacy archives removed after all Go jobs joined; rebuild normally. Plan/journal: `build/port-artifact-cleanup/remaining-exports-cache-{plan,removed}.json`. |
| Older leaf/transfer binaries and cache entries | Rebuild from recorded revisions as needed; cleanup records are in `build/port-artifact-cleanup/`, including `superseded-leaves-xfer-binaries-20260924/` and `string-boundary-cache-removed.jsonl`. |

Historical batch details belong in `docs/porting/` and Git history. Earlier
checkpoint/archive details are recoverable at `b034c43e` and `69669bcb`; their
old current/next instructions are historical, not the active plan. Local logs,
archives and rebuildable binaries are not backed up by pushing Git.
