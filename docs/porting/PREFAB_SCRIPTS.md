# Prefab script merging and generation setup

## Result

**Eighteen live C algorithms are native Go.** Two unused selection helpers,
twenty C implementations/interfaces, four newly unused Go-backed C adapters and
five C globals are retired. Go callers invoke Go directly. No C algorithms remain
solely for tests. The repaired C baseline is **5c83d11a**, committed and pushed
before conversion; [selection](prefab-scripts-selection.json) retains the original
20 bodies / 1,363 body lines and hashes.

Production C is **35,521 physical lines in 72 files**, zero reference C:
**−1,403** from the repaired baseline, or **−1,396 net** including the prerequisite
repairs since the preceding runtime conversion. The three native implementation
files contain 621 Go lines.

## Coverage and qualification

The baseline's sixteen focused groups passed in two separate processes, producing
12 frozen captures / 4,182 records. One capture (16 selection-copy records) is
historical only: reachability found no production caller, callback or registration
for `nox_xxx_tileInitdataClear_4D3C50` or `sub_4D3C70`. Remove both orphan bodies and
their dedicated fixture instead of translating code solely to retain a test.

The native focused run passes **17 groups**. Its **11 retained captures / 4,166
records** match the committed C hashes. Added independent contracts cover existing
builtin-predicate dispatch, empty pending lists, embedded-NUL name termination,
incomplete instruction streams, invalid generation directories and failed source
copies with a stale backup.

All **156 selected root tests / 2,057 including subtests** pass without skips on
default/server/highres, with package times **23.784s / 22.969s / 23.907s**. All
**90 connected captures / 34,943 records** match both C and one another. All runs
use the same unchanged 2,341 source files. Static mapped-memory checks pass.
The corpus exercises every instruction, wide/signed operands and wrapping indices,
complete script merges and rewrites, real VM readers, file ownership, callback
storage, object/waypoint naming, pending references, group references, mixed-sign
bounds ordering, generation reset, library metadata and related runtime owners.

Fresh qualification passes all three ELF32/i386/SSE2/CGO binaries and ABI checks;
24 retired function symbols and five retired C globals are absent. The full asset
suite exactly matches the known result: 1,553 failure entries, with 15 passing,
3 failing and 32 skipped packages. Headless gameplay, explicit save/load and
flat-map regeneration comparisons pass. These checks do not claim a green legacy
asset suite or replace manual physical display/audio checks.

- [Native qualification](prefab-scripts-native-qualification.json)
- [Repaired C qualification](prefab-scripts-c-qualification.json)
- [Frozen captures](prefab-scripts-captures.json)
- [Batch manifest](prefab-scripts-batch.json) and [test selection](prefab-scripts-tests.txt)

Authoritative native evidence is under
`build/port-prefab-scripts/native-final-{default,server,highres,production}`,
with `native-fourth` focused results and `native-capture-audit.json`. An earlier
native sweep passed, but its production run was intentionally stopped after final
review identified predicate dispatch duplication; the final runs qualify that
correction too. No frozen expectations were regenerated to hide differences.

## Prerequisite corrections for review

Independent contracts reproduced these existing defects before freezing C:

1. Several serialized 32-bit operands used one-byte C locals. For example,
   copying integer 128 with relocation disabled produced −128. Preserve the
   full word, matching the real assembly decoder and VM format.
2. Both suffix helpers formatted a name and then unconditionally replaced it with
   `ERROR_NAME_TOO_LONG!`. Preserve the distinct object/function suffix when it
   fits the actual 256-byte limit; use the placeholder only when it does not.
3. Reserved-function merging read variable-length names into a four-byte scalar.
   Actual generated scripts contain the six-byte name `GLOBAL`. Restore a proper
   temporary buffer and bounds in the C baseline; Go uses bounded byte storage.
4. Root file merging failed to create a missing destination, retained stale tails
   (4,096 bytes instead of the independently expected 146), and left three
   temporary file registrations. Create/truncate output and preserve caller-owned
   files/handles. The native merger needs no temporary C registrations at all.
5. `PendingByScriptID` called `Next()` without assigning it. A missing target
   looped forever; the isolated process stack confirmed the real engine loop.
   Assign the next pointer. Missing/duplicate targets and actual group remapping
   now have contracts using real pending-object owners.
6. The dependency script reader assigned both coordinate suffixes to X.
   `OnEnter%7%-46%92` became `(92,0)`. The local VM consumer restores both signed
   axes without editing dependencies. Complete rewritten scripts pass the VM.

These are deliberate, reversible correctness repairs under the standing policy.
The repaired C baseline and native conversion each received fresh production
qualification.

## Native compatibility and decisions

Preserve the C clock bridge's **32-bit truncation** before its timestamp is widened
to 64-bit storage: `0x10000002a` stores 42. Preserve explicit float64 addition before
float32 coordinate storage, mixed signed/unsigned bounds comparisons, first-match
pending lookup, wrapping counters and the distinction between nil and allocated
empty object names. Type-owned transfer callbacks select name/reference handling;
an instance callback field is not a substitute. Repeat passes leave names intact
once their marker flags are cleared.

The merger calls the existing Go builtin-remapping predicates. An eight-case
independent contract changes predicate answers and verifies invocation order,
conditional dispatch and resulting operands; the ordinary C captures alone would
not distinguish duplicated default IDs from the actual dispatch owner.

The native implementation returns failure for incomplete instruction streams and
empty/oversized generation directories, whose C behavior was undefined. The
maximum directory length is 2,035 bytes, preserving the 2,048-byte C path capacity
including the longest shipped suffix and terminator. Failed source copying stops
before opening a stale backup. Contracts exercise these local failure paths;
they do not assert that every malformed script format has been redesigned.

## Disk maintenance and recovery

Verified deduplication reclaimed about 1.55 GiB from three completed older scenario
copies; each has a restoration manifest. Original assets/archive, changed files,
saves, reports and binaries were preserved. Older disposable Go cache entries were
also removed after all gates joined, taking free space to about 13 GiB before the
new builds. `cache-cleanup.json` records the completed second pass (6,839,608,799
bytes); an earlier pass stopped at a cache directory after some regular-file
removals. These are disposable caches, not source or qualification artifacts.

The capture freezer, native source installer, finalizer and deletion passes under
`build/port-prefab-scripts` are consumed. Do not rerun them. Native drafts are stale;
actual source and committed expectations are authoritative.
