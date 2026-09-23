# Porting checkpoint

Read [PORT.md](PORT.md) for the working plan. This is the resume checkpoint.

**Qualified C remaining: 51 physical lines in 4 production `.c` files**, zero
reference C. Latest conversion removes 44 global definitions, eight mapped buffers
and two translation units (−63 lines). libc/CGO remains. See
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

## Next — reviewed orphan cleanup, still unapplied

`build/port-orphan-inline` contains a reviewed six-file patch for eight unused
inline/header helpers, five unused callbacks and two empty nullsub_35 calls/body.
Preserve both NeedSync calls/control flow and all ten live callback identities.
Expected next C count: 45 lines /4 files. The four existing Obelisk tests reach
objectiveObelisk via op808; an earlier draft named the wrong enclosing function.

After the storage commit/push, verify all six source hashes before refreshing the
orphan HEAD guard. The draft installer intentionally refuses --apply. Apply only
the reviewed guarded patch; no other source changes are needed. The AST input
preparer/checker must show exactly two removed calls and no other Go AST changes.
If formatting positions cause a false difference, fix the audit, not game code.

The ignored native manifest uses a one-line four-test regex, its own scenario,
2,379 cumulative retired symbols and ten retained callbacks. run-native.sh,
deduplicate-preflight/save.py and finish-native.py are prepared but NOT RUN.
The finalizer checks the qualified storage baseline, all gates, source/manifest
identity, explicit save/load and distinct callback addresses in baseline/current/
safe binaries. Primary reviewed and strengthened Luna's draft. Inspect before use.

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
