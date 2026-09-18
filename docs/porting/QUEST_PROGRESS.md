# Quest progress and stage preparation

## Scope and checkpoint

Parent **4a8d73ed is committed and pushed**. This connected batch covers **24 live
C functions / 715 body lines**: script quest-variable list operations, namespace
qualification and wildcard reset, registered builtins, save/load records, generator
type mapping and stage preparation, Hecubah/Necromancer spawning and stage state.
Manual whole-source/header/preamble review confirms sub_51A950 has only its
definition and header declaration; it is an orphan to remove. See [quest-progress-scope.json](quest-progress-scope.json).
Production C remains **39,193 physical lines / 74 files / zero reference C**.

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
Source remains entirely C for this batch. Commit/push this baseline before installing
the reviewed Go drafts. Fresh production is required after conversion.
