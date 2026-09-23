# Reward-item use callback registry

## Scope

Use the existing typed registry for SpellRewardUse, AbilityRewardUse and
FieldGuideUse. Register the existing bookUseSpell/bookUseAbility/bookUseGuide
functions, preserving their C callback addresses and per-object data sizes.
Registration still occurs at package initialization with the same duplicate-name
checks. Object layouts, exported C interfaces, the generic registry implementation
and its fallback for unregistered C callbacks remain unchanged.

The old path calls C through CallIntPtr2, reenters the existing Go export and
recovers the same two object pointers on 386. The new closures call those same
Go implementations and use the same nonzero-result boolean conversion. No gameplay
algorithm or allocation ownership changes. No measured performance claim.

## Baseline and tests

The existing TestBookAwardsItemUse remains on its C-entrypoint route. A sibling,
TestBookAwardsItemUseRegistry, shares its unchanged 576-case matrix and expected
values, looks up each real registration through PortTestWorldUseRegistry, checks
the exact callback address/data size, and invokes it.Use.Get()(u, it). The existing
object/player owners reset levels, flags, deletion lists, audio and message state
between cases. Cases include classes, quest combinations, invalid IDs and level
limits. Both routes retain the original frozen capture hash:
`ce8035fbce4dee00ac94d7250a5ea74c2cd5d70d148b08e760890bded0798a7f`.

Before production changes, both routes match that hash in default/server/highres.
All 20 selected roots per profile pass without skips: BookAwards, UnitGameplay
ReadUse/ReadMessageLimits/RewardOwners and WorldCollisionsSpellAward. These cover
reward effects and adjacent consumers; the new sibling specifically covers the
three registry entries. Original-C and registry captures have separate output
paths. Baseline session32128 joined PASS. See
[baseline qualification](reward-use-registry-c-qualification.json).

Production evidence is reused from `787f3a08` only after checking all source
fingerprints (the sole difference is the porttest test file), runtime/discovery
environments and four binary hashes. All three callback exports exist in those
binaries. After conversion rerun the 20-root selection, safe/static, four fresh
binaries/ABI, exact known-suite comparison and headless creation/save-load.
The unrelated full mixed-forwarding selection is not repeated: no shared registry
implementation changes. All three C callback exports must remain Go-backed.

## Delegation and count

Luna drafted the three registration changes and the sibling test. Primary reviewed
pointer conversion, initialization, callback identity, ownership and the unchanged
matrix/expectations; no behavioral draft correction was needed before baseline.
Primary owns qualification and acceptance. No model cost/speed savings measured.

Standalone C remains zero physical lines/files; production C preamble bodies remain
79. This removes three registered Go→C→Go routes, without changing those body counts.
