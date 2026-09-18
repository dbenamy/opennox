# Quest progress and stage preparation

## Scope and checkpoint

Parent **4a8d73ed is committed and pushed**. This connected batch covers **24 live
C functions / 715 body lines**: script quest-variable list operations, namespace
qualification and wildcard reset, registered builtins, save/load records, generator
type mapping and stage preparation, Hecubah/Necromancer spawning and stage state.
Manual whole-source/header/preamble review confirms sub_51A950 has only its
definition and header declaration; it is an orphan to remove. See [quest-progress-scope.json](quest-progress-scope.json).
Production C is now **38,498 physical lines / 74 files / zero reference C**
(**−695 physical lines**).

## Baseline design

Use the actual C linked list and script VM. Capture every record byte with list
links normalized by owned-list position; check predecessor links independently.
Exercise qualified/unqualified and case-varied names, empty names, bounds within
the original 132-byte storage, numeric bit patterns, existing-record kind retention,
missing values and each wildcard branch. Ignore incidental reset return values
because all actual callers ignore them. Supply both shipped colon separators.

Use the real global file owner for exact serialized bytes/checksum/position and
round trips. Cover game-mode write admission, list order, integer/float kinds,
duplicate names, unknown kinds and the signed 16-bit version comparison. Do not
freeze undefined uninitialized reads as a desired contract.

Audit and install the shipped generator mapping and nonzero type IDs before
capturing stage preparation. Reuse actual object/type allocation, world lists,
RNG, delayed deletion, creation and reward owners. Cover all generator levels and
admission/rate branches, exit pruning, minion thresholds and marker cleanup.
Boss cases require raw health rounding/overflow, definition overrides, spawn
fields, inventory and reward selection with real registered types.

Repeat C captures in separate processes and freeze only after independent contracts
pass. Qualify affected default/server/highres selections; reuse the immediately
preceding production only with production-source and binary identity evidence.
Fresh production, headless gameplay, save/load and flat regeneration follow the
conversion. All eleven focused roots now repeat byte-for-byte in separate processes: **3,156
records**. Expectations are frozen and all-target qualification passes.

## Review notes

The modern noxscript API still has separate quest-status TODO methods. This batch
qualifies and ports the registered legacy VM path; extending that separate API is
not necessary to preserve the current caller behavior.

Existing kind is not replaced when a setter changes an existing entry. Wildcard
prefix comparison is case-insensitive while suffix substring matching is
case-sensitive and uses its first occurrence. Unsupported positive versions clear
the old list before returning failure; high-bit versions pass the signed check.
These are observable original behaviors to retain on defined inputs.

Evidence and proposals: build/port-quest-progress/{proposal.json,references.json,
original-bodies.txt,c-initial-result.json}. The whole-repository reachability
proposal is conservative and requires manual review before deleting candidates.


## Frozen C corpus

The eleven captures contain 1,215 variable-operation records, sixteen reset-mask
sequences, fourteen registered-builtin cases, 38 serialization cases, 58 shipped
mapping snapshots, two empty-map cases, nine stage-word cases, two name-boundary
cases, 1,080 generator/exit preparation cases, 542 direct boss cases and 180 combined
marker-selection cases. All repeat identically. An independent RNG contract
predicts admission, selected authored positions and exact consumption for marker
preparation. Boss fixtures validate fixed attributes, nonzero wrapped health,
spawn position and actual inventory/reward ownership, including missing factories.

The first fixture compile needed an explicit GameFlag conversion; this was test-only.
Static mapped-state checks pass in the configured target environment. The production
reuse report verifies unchanged original source and all three parent binary hashes.
No production fix or translation was needed before freezing. Modern API TODOs and
undefined malformed-input behavior are not silently frozen as desired behavior.


## Qualified C baseline

All three affected target sweeps pass with zero skips and identical captures. The
frozen focused repeat passes; all four gates share one unchanged source manifest.
See [quest-progress-c-qualification.json](quest-progress-c-qualification.json)
for counts/timings and [quest-progress-c-production-reuse.json](quest-progress-c-production-reuse.json)
for verified production reuse (2,272 original sources unchanged, eight new
porttest-only files, all three binary hashes checked).

The initial focused manifest had its pattern argument at the wrong command index;
that driver invocation failed before tests ran. The corrected c-frozen-2 run passes
on exactly the same source as the three successful sweeps. No captures changed.
The C baseline was committed/pushed before installing the Go drafts. Fresh
production qualification after conversion is recorded below.


## Qualified native implementation

C baseline **1c915175 is committed and pushed**. The twenty-four live functions
are replaced, the orphan getter is removed, and the quest-variable list now uses
Go records with the qualified 148-byte layout. Only reset/read/write retain C
exports. Twenty-four obsolete interfaces (including two upstream bridges) and one private
C global are retired;
script builtins now call Go through their registered VM entries.

Compiled-C review confirms the two boss health conversions have distinct widths:
current HP truncates the wide product to int64 before its 16-bit store, while max
HP converts a spilled float32 copy. Preserve both, byte-wrapped generator limits,
signed minion-stage comparison, exact RNG order and first-occurrence suffix matching.

Reversible decisions for review: reject names or composed namespaces exceeding
131 bytes instead of overwriting adjacent fields; reject truncated save records
and stop at EOF even if the record count is oversized. Preserve already accepted
records on a later malformed record and the original clearing-before-version-check
behavior. An independent native test covers these limits, all truncation positions,
large counts and embedded-NUL C-string semantics. This does not change the frozen
captures for defined original input. Full native qualification is complete.


First native focused run passes all twelve roots with eleven unchanged captures /
3,156 records. Final caller review retires the scalar-health and type-health C
bridges whose last C users were the boss factories. The existing scalar fixture
now invokes its same Go implementation directly. Static mapped checks pass;
focused-2 and full native qualification pass. Current C is **38,498 lines /
74 files / zero reference C**; the net physical reduction is **695 lines**.
The raw C-to-Go body count differs because existing address markers and layout
remain. No artificial blank-line cleanup is counted as translated algorithms.


## Completed native qualification

Final focused-2 passes twelve roots in 0.164s; eleven frozen captures / **3,156
records** remain unchanged. All three affected-target sweeps pass without skips
and reproduce all **227 captures / 150,824 records** from C. Fresh production
passes three binaries, ABI/interface checks, the exact known asset-suite failures,
headless gameplay, save/load and forced flat-map regeneration. All four gates share
unchanged source and every session is joined. See
[quest-progress-native-qualification.json](quest-progress-native-qualification.json)
for exact counts, timings, source identity and binary hashes.

No golden expectation changed. All freeze/install/retirement/finalization scripts
are consumed. Current artifacts: build/port-quest-progress/native-{default,server,
highres,production}, native-focused-2, native-interface-audit.json and
upstream-retired.json. Original assets/archive remain intact.
