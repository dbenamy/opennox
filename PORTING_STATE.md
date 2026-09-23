# Porting checkpoint

Read [PORT.md](PORT.md) for the working plan. This is the resume checkpoint.

**Qualified C remaining: 45 physical lines in 4 production `.c` files**, zero
reference C. Latest cleanup removes six standalone C lines plus eight unused
header/preamble helpers. The preceding storage conversion removed 44 global
definitions, eight mapped buffers and two translation units (−63 lines). libc/CGO
remains. See
[C_LOC.md](docs/porting/C_LOC.md).

<!-- current-checkpoint -->

## Current — remaining storage qualified

The remaining-storage conversion and fixture identity repair pass full
qualification against pushed baseline `781901ef`: 2,291/2,280/2,291 consumer roots
across default/server/highres without skips, three legacy contracts per profile,
static checks, safe build/symbol checks, three production binaries/ABI, exact
known-suite comparison, headless gameplay and save/load. Both frozen storage
captures remain unchanged: 884 raw patterns and 54,905 numeric patterns.
Safe runtime was not tested; the known full-suite failures remain unchanged.

Pipeline 37056 JOINED PASS. Final artifacts are
`build/port-remaining-storage/lifetime-{default,server,highres,safe,preflight,production}`.
Source fingerprints match across all phases and the preflight binary matches
production. See [RAW_STORAGE.md](docs/porting/RAW_STORAGE.md) and its committed
qualification JSON. The allocator retains foreign storage for process lifetime;
this does not remove libc, C types, callback boundaries or third-party C.

The fixture repair has deterministic red/green evidence. Per-case identity-map
entries outlived freed allocations; cleanup now removes transient keys after all
snapshots while retaining aliases for persistent objects. The initial repair
reset a persistent player alias; complete JSON diffs caught that and the final
repair preserves every frozen capture. Failure-only callback capture is enabled
for future diagnostics. The original intermittent stats digest is not proven to
have that cause. See [FIXTURE_IDENTITIES.md](docs/porting/FIXTURE_IDENTITIES.md).

Earlier native-* and scoped-* failures and the deliberately interrupted observed-*
run remain local evidence. Diagnostic captures are losslessly compressed with
SHA/round-trip records. All Go/build/scenario jobs are joined at this checkpoint.
The remaining-storage migration generators and successful finalizers are CONSUMED;
tracked source supersedes ignored drafts. Do not rerun them.

This checkpoint is included in the qualified storage commit. On resume, check
Git status/log/remote to establish whether its push completed; no commit hash
is inferred from this checkpoint alone. Preserve the untracked asset archive.

## Current follow-up — orphan cleanup qualified

Storage commit `f5970121` is pushed. Its six-file follow-up removes eight unused
header/preamble helpers, five unused callbacks and the empty nullsub_35 function
plus its two Obelisk calls. Both NeedSync calls and all ten live C callback
identities remain. C is 45 lines /4 files (−6). See
[ORPHAN_INLINE.md](docs/porting/ORPHAN_INLINE.md) and its native qualification JSON.

Pipeline 18151 JOINED PASS. All three contracts and four Obelisk roots pass in
all three profiles without skips or changed storage hashes. Static, safe build,
three production binaries/ABI, exact known-suite comparison, gameplay and explicit
save/load pass. Source fingerprints match, preflight matches production, and all
ten callback addresses remain distinct in baseline/current/safe binaries. The
AST comparison shows only two removed calls; its line-position printing artifact
was corrected in the audit. All jobs are joined. Migration and finalizers are
CONSUMED. Check Git log/status/remote for this checkpoint's commit/push status.

Artifacts are build/port-orphan-inline/native-{default,server,highres,safe,
preflight,production}. Scenario asset restoration uses that directory's
`deduplicate-preflight.py --restore orphan-inline-native` and
`deduplicate-save.py --restore orphan-inline-native-save`; inspect per-run records.

Next: consider direct libc calls for the two entry-character forwarding helpers,
with exhaustive uint16 boolean parity and existing UI entry gating tests. Luna's
read-only plan is under build/port-entry-classifier-audit/direct-libc-plan.md and
still requires primary review. Preserve libc locale semantics; no ASCII/Unicode
substitution is authorized by this implementation choice. Live callback exports,
safe adapters and third-party MP3 implementation remain later work.

## Current — entry classifier C baseline qualified

The actual-C baseline passes four legacy contracts and twelve widget roots in all
three profiles, no skips/changed goldens, plus static, safe build, production/ABI,
exact known-suite, gameplay and explicit save/load. Pipeline23563 JOINED PASS.
Source fingerprints agree, preflight matches production, and ten callback
identities remain distinct. Three fresh-process captures froze all 131,072
classifications to SHA256
`f5c39db5e866885891bbb1fe8d0b2244d9fbcea7d232cdd4b7e6f90dd03ff6e0`.
No reference C was added; the test calls the actual production wrappers through
small Go predicates. See ENTRY_CLASSIFIERS.md and the C qualification JSON.
Check Git status/log/remote for this checkpoint's commit/push completion.

Artifacts are build/port-entry-direct/c-{default,server,highres,safe,preflight,
production}. Finalizers/cleanup are CONSUMED. All jobs joined. Restore scenario
assets using deduplicate-c-{preflight,save}.py and per-run restore metadata.

Next: replace only the two predicates' C.entry* calls with C.isw*(C.wint_t(v))
and remove their two C preamble bodies; preserve the caller/branch and all frozen
expectations. Install reviewed ignored native-batch.json, record native-guard.json
with committed HEAD/12 consumer names, then run run.sh native. After all gates
pass use finish.py native, native save deduplication and finish-docs.py native;
commit/push and continue. All scripts require the baseline environment and native
jobs require host execution. Standalone C remains 45; two preamble bodies retire.

## Other prepared evidence

Read-only inline/header, safe-bridge and empty-callback audits are under
build/port-inline-c-audit, port-safe-bridge-audit and port-empty-callbacks.
The entry-classifier audit under port-entry-classifier-audit was corrected by the
primary: live constructor callers include src/gui_console.go and src/gui_widgets.go.
No locale mutation was found in repository/module source, but linked-library
locale behavior is not established. Preserve C classification semantics.
Unimplemented follow-up ideas are in port-empty-callbacks/go-entrypoint-design.md.

One bounded GPT-6 Luna helper remains the default. Current helper is idle;
implementation drafts, caller inventories and finalizer scripts require primary
review. Exact whole-source search commands/path evidence are now required.

## Disk and assets

Safe/preflight cleanup round one is CONSUMED: six obsolete phase outputs,
301,307,320 bytes, after 159 host process checks and fresh file/evidence checks.
Five had observed current hashes rather than hashes in old phase results; one
matched a retained qualified production binary. No qualification linkage was
inferred for the others. Source, logs, manifests and captures remain.

Cleanup rounds five/six are CONSUMED: three old render-helper executables
(145,417,452 bytes, 162 host processes checked) and eight older highres/server
executables (384,739,612 bytes, 164 processes). All evidence remains under
build/port-artifact-cleanup. Current storage/orphan/entry/audio artifacts,
assets, scenarios and caches were excluded.

Cleanup round four is CONSUMED: 12 older qualified production executables,
582,770,236 bytes, removed after verifying metadata/hashes/ABI and checking 158
host processes. Luna prepared the inventory; the primary revalidated and applied
it. Evidence remains in build/port-artifact-cleanup/older-binaries-round4*.json.

Old executable cleanup rounds two/three are CONSUMED: 33 files /1,609,335,708 bytes
and 18 files /876,580,664 bytes, respectively, after host process/open-file checks.
All qualification metadata, logs, scenarios, original assets and caches remain.
Plans/consumed records are under build/port-artifact-cleanup.

The completed latest preflight and save/load runs released 556,388,715 and
556,358,986 bytes of verified duplicate assets; restoration records remain with
their runs. For the
latest preflight use `build/port-remaining-storage/deduplicate-scoped.py --restore
raw-storage-scoped`; for save/load use `deduplicate-save.py --restore
raw-storage-native-save`. Inspect each run's deduplicated-assets.json before reuse.
Preserve build/assets/extracted/drive_c/Nox and nox-iso-from-archive-org.7z.
