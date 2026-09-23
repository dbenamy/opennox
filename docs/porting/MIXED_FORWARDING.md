# Mixed-signature Go forwarding

## Scope and baseline

The qualified batch replaces 94 Go→C→Go calls with their existing Go implementations.
Luna drafted 92 calls across 40 functions and 36 caller files. Primary independent
reconstruction matched every edit, including 37 explicit int/int32 argument
adaptations. Existing C.int narrowing expressions remain where present; literals
fit their destinations. Fresh generated metadata at `9b548e8f` confirms the 40
signatures after primitive-width projection on the supported 386/SSE2 target.
Pointer/typedef aliases match; there are no selected-name macros. This is not a
64-bit compatibility claim. Algorithms, argument order and return uses stay intact.

Primary separately adds two calls and retires their unused C interfaces:

- `sub_50B510`: the local C declaration returns int, while the Go function returns
  nothing. The sole caller ignores the result. Call that unchanged Go body directly.
- `nox_xxx_castCounterSpell_52BBB0`: the local C declaration takes integer object
  addresses, while the Go export takes pointers. The sole caller passes the same
  object three times. On 386, `inventoryInt(u)` and `asObjectC(u)` carry the same
  address bits; use typed pointers and preserve spell ID 13 and evaluation order.

Whole-source reference review found no other C callers, callback addresses or
macros for those two interfaces. All four preceding qualified binaries contain
both exports; the new manifest requires their absence. This reversible cleanup
is recorded for review. Do not describe the old mismatched signatures as equal.
Remove two newly unused include-only C preambles/imports in effects_weapon_use.go
and inventory_pickup.go. Qualified source scope: 38 files, 94 calls.

## Ownership and coverage

Message-list storage copies caller bytes. String lookup copies name/file inputs,
returns interned strings and optionally writes an interned narrow-string pointer.
Audio/wall positions and rectangles are copied by value; map-trace outputs are
synchronous. Script blocks belong to stable object update data and may be queued
by the existing VM. Creation/free and state changes keep existing server owners.
Spell acceptance uses the existing typed argument path, also called directly from
Go with scoped argument storage. No caller buffer expression or function body is
changed by the mechanical draft.

Reuse the preceding 686/685/686 affected roots only after exact source, test,
runtime/discovery environment and four-binary hash checks. Six additional existing
roots pass fresh in default/server/highres on unchanged production:
WorldCollisionsQuestExit, QuestReadiness, AudioFrames, TrapAbility, ScriptTriggers
and ProtectionSetABI (each prefixed `Test`). Combined selection: **692/691/692**.
No skips or expectation changes. Baseline session 54569 joined PASS; see
[mixed-forwarding-c-qualification.json](mixed-forwarding-c-qualification.json).

The affected tests cover gameplay owners, not every individual branch. QuestExit
does not install active abilities, and AudioFrames calls the server audio owner
directly; neither alone proves its corresponding selected wrapper branch. Existing
SustainedSpellsTeleport exercises the audio wrapper path. ScriptTriggers is
adjacent coverage. Existing AttackAbilities exercises the counterspell caller;
ObjectStateFreeze exercises the reset caller but does not independently assert
private path-cache fields. Frozen expected values remain unchanged.

## Acceptance and delegation

All 692/691/692 native roots pass default/server/highres, without skips or changed
expectations. Exact discovered names match the baseline union. Safe/static checks,
four fresh binaries/ABI, exact known full-suite comparison (304 existing failure
events; 17 pass/2 fail/32 skip packages), headless creation and save/load all pass.
All 42 selected C-call bridges and both retired exports are absent from all four
binaries. All 40 other selected exports remain Go-backed, including the default
damage-sound address still stored in two callback slots. Production has no PortTest
symbols. See [qualification evidence](mixed-forwarding-qualification.json).

Luna supplied the bounded 92-call draft and an owner-coverage audit. Primary owns
signature/ownership review, the two mismatched interfaces, final test selection,
application and qualification. The draft matched independent reconstruction;
the coverage report explicitly distinguishes adjacent coverage from branch proof.
The first discovery build caught one missing result adaptation in respawn modifier
lookup: the direct Go function returns int32, while its local desc helper accepts
C.int. Primary retained that result boundary with an explicit C.int conversion.
Luna’s draft and primary text reconstruction both missed the contextual Go type
requirement; compilation caught it before tests ran. The failed discovery is
retained under contracts-build-failed. No model cost or speed savings have been measured.

Standalone C remains zero physical lines/files; production C preamble bodies stay
79. The two removed preambles contain only includes and are not counted bodies.

## Disk recovery during qualification

Five reproducible obsolete compiler archives were removed after checking retired
markers, hashes/stat metadata and host open-file references (226,475,872 bytes).
Six completed logs were archived losslessly: 75,474,366 → 4,236,945 bytes. Ten older
scalar/typed qualification binaries were archived: 488,736,004 → 230,483,100 bytes.
A second nine-binary pass archived 438,720,288 → 206,129,910 bytes.
Their exact bytes match committed qualification hashes; host open FD, executable
and mapped-file references were checked before removing raw copies. Gzip
round-trip checks and restore commands are recorded under build/port-mixed-forwarding.
Restore archived binaries/logs before rerunning historical finalizers. Current
pointer baseline binaries and original assets remain available uncompressed.
Archival briefly overlaps testing; these durations are not performance benchmarks.

Two identical pointer-baseline binaries now share an inode, preserving both paths
and verified bytes while reclaiming 49,077,784 bytes. The host-reference checks
passed; pointer-binary-shared.json records the pair. Separate the inode before any
intentional mutation of either path.
