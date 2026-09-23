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
