# Native modifier dispatch and observers

Status: native modifier dispatch and all four observers qualified against `f77b9d51`.
Qualified parent: `defbdb98`.

Keep the existing native modifier registries and production handler registrations,
replace their four foreign fallback branches with explicit rejection, and retire
CallVoidPtr5/CallVoidPtr6 from the generator and generated dispatcher. Translate
four C observers: generic modifier forwarding, equipment, damage and melee attack.
Preserve frozen assertions, native identity comparisons, mutable fire-effect hooks,
signed results, discarded results, pointer lifetimes and all runtime.KeepAlive calls.
External bindings and allocation domains remain unchanged.

## Producer and fixture scope

All 40 production named registration keys match a native handler of the expected
arity. Calls occur only in legacy/modifier_register.go.
The parser in server/modifiers.go resolves names through those tables and writes
*ModParseTarg.Func. Eight descriptor slots use damage/defend/update/engage tables;
empty strings leave nil and unknown names return an error. Existing native maps
supply the corresponding arity handlers. Review all descriptor writers and raw
field aliases, not only generic dispatcher callers.

The generic observer is not the only foreign-key producer. Equipment installs
engage/disengage/defend observers; damage installs scalar/two-word defend and
pre-damage observers; melee installs an attack/pre-hit observer. Move all four
fixtures together. Modifier item-update tests also reuse the generic void observer.
Metadata-only test keys and fields compared without invocation are distinct from
executed producers. Other callback families retain their fallbacks.

Unknown modifier keys deliberately panic after the producer audit. This reversible
invalid-input correction does not preserve foreign C crash diagnostics. Go array
bounds also make invalid fixture trace overflow/indexing explicit rather than C
out-of-bounds behavior. Valid fixture capacities, output and captures remain exact.

## Observer review

- Generic forwarding retains four distinct identities, six uintptr argument words,
  signed 32-bit result, call count/kind and optional complemented uint32 output.
- Damage preserves event kinds, pointer-word order, scalar versus two-word reads,
  pre-mutation observations and optional raw uint32 output replacement.
- Melee copies nine record words, masking the same padding bytes. Collision sets
  the Front-word mask to zero; observations precede optional first-word mutation.
- Equipment clears only its trace on reset, preserving configured return/output
  arrays. It increments row count even when capture capacity is exceeded, stores
  only rows passing the original bound, and preserves signed return and float bits.
  Stable callback identity normalization uses the original capture IDs.

AST review: 13 existing functions changed, two generated dispatchers removed,
35 unchanged, 14 added native observer helpers/init functions. Generated dispatcher
output matches the draft byte-for-byte. Production handler registrations are
outside the change. The draft helper's initial C-selector check mistakenly scanned
the generator's C code strings; corrected that check before producing the manifest.
No source was installed or compiled from that incomplete draft.

## Qualification

Four shared fixtures underpin many nested world/gameplay owners, so use the prior
broad audited root selection: 2,482 default/highres and 2,471 server roots. Remove
only the retired prefab traversal contract and include the newer animation
completion contract. Preserve the same eight previously excluded compiled roots;
this is not a full-corpus claim. The exact known-suite production gate remains.

Run original focused forwarding/update/equipment contracts separately, followed
by the broad selection on exact-source qualified parent binaries. Native preflight
also includes TestDamageModifiers and TestAttackEffectDispatch, already included
in each original broad selection, to catch observer mismatches early. Validate names,
source, binary hashes and environment. Repeat matching selections after conversion,
plus safe contracts/build/static, three production builds/ABI, exact known-suite
outcomes, fresh save/load and all original asset hashes. No goldens change.

Qualified counts: embedded production C bodies 14→12; fixture C-import files 98→94;
production cgo files remain 5 client / 6 server; standalone C and legacy exports
remain zero; headers remain 157 files / 2,731 lines.

Primary only; Luna remains quota-unavailable. Local artifacts:
`build/port-modifier-dispatch/`. Before starting the baseline, removed seven
superseded callback binaries (387,018,752 allocated bytes), after verifying their
committed source, current replacements, hashes and host use. Rebuild `dcab2e38`
if needed; logs and metadata remain in the prior batch. Cleanup records:
`build/port-modifier-dispatch/callback-cleanup-*`.

Also removed seven superseded animation binaries (387,063,808 allocated bytes)
after committed-source, replacement/hash and host-use checks. Rebuild `5233d084`
when those historical executables are needed. The active original modifier tests
use retained UI completion binaries. Records: `animation-cleanup-*` in this batch.

Original acceptance: all 2,482/2,471/2,482 roots passed without skips; four focused
contracts passed separately in each profile; three safe contracts passed. Verified
exact discovered/run/pass name sets, source and binary hashes, and relevant
environment values. [Original baseline](modifier-dispatch-baseline.json).
Native qualification passed against this accepted baseline.

All original and matching native broad roots passed without skips: 2,482 default,
2,471 server, 2,482 highres. Four focused contracts passed separately in each
profile on both versions; native preflight also passed the damage/melee captures.
Three original and native safe contracts passed. Safe build/static, three production
builds/ABI, exact known-suite outcomes and fresh save/load passed. All 1,654
original asset hashes are unchanged. Retired dispatcher/observer symbols are
absent, remaining C bodies unchanged, and all 40 production registrations unchanged.
No assertions or frozen captures changed. Evidence:
[baseline](modifier-dispatch-baseline.json),
[qualification](modifier-dispatch-qualification.json),
[inventory](modifier-dispatch-inventory-after.json).

After all jobs joined, removed 1,654 verified duplicate scenario assets
(560,107,520 allocated bytes), preserving originals and recovery records.
Restore before replay: `python3 build/port-artifact-cleanup/restore-recent-scenario.py
build/baseline/runs/modifier-dispatch-save/deduplicated-assets.json`.
