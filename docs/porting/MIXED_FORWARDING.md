# Mixed-signature Go forwarding

## Scope and baseline

The next batch replaces 94 Go→C→Go calls with their existing Go implementations.
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
and inventory_pickup.go. Total proposed source scope: 38 files, 94 calls.

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

After conversion require the full combined native selection, safe/static checks,
four fresh production binaries with ABI checks, exact known full-suite failure
comparison, headless creation and save/load. Require all 42 selected C-call bridges
and the two retired exports absent; retained callbacks remain Go-backed.

Luna supplied the bounded 92-call draft and an owner-coverage audit. Primary owns
signature/ownership review, the two mismatched interfaces, final test selection,
application and qualification. The draft matched independent reconstruction;
the coverage report explicitly distinguishes adjacent coverage from branch proof.
No model cost or speed savings have been measured.

Standalone C remains zero physical lines/files; production C preamble bodies stay
79. The two removed preambles contain only includes and are not counted bodies.
