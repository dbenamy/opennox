# Sustained-spell callback identities

Scope: 53 actual Create/Update/Destroy callback identities used by duration spells:
42 sustained-spell owners, six charm/summon callbacks, three wall hooks, one
teleport-start callback and the energy-bolt no-op destruction callback. Preserve
all getter identities, 32-bit layout, ownership, callback order, timing, exact
return words, mutable hooks and direct fixture normalization.

Introduce a native duration callback registry with separate result-returning and
result-discarding call paths. Native callbacks preserve the original 32-bit word;
unknown C addresses retain their original integer/void fallback conventions.
Keep the existing pointer fields and New signature. Replace the root duration
service's three raw call sites, all 53 getter keys and the associated fixture
routes. Do not alter spell algorithms or external native-library bindings.

The only char-returning sustained callback is used for destruction, where its
result is ignored. Preserve its signed low-byte result in direct fixtures. Void
destruction callbacks use a native result of zero only as an unused adapter
value. Historical float parameters transport pointer bits, not numeric values;
pass the same DurSpell address directly to the existing native owners.

## Original contracts and testing

Three new porttest files add four independent roots: all 53 keys remain distinct
and stable through foreign storage/GC; independent C observers verify all return
bits and nil/non-nil argument forwarding with int/void calling conventions; wall
hooks remain replaceable after first dispatch; and the real duration allocator,
list and root lifecycle preserve create acceptance, update cancellation, deadline
short-circuiting and destruction across 256 combinations.

Luna drafted the observer/lifetime/hook contracts. Primary made the observer
self-contained, re-read getters after GC, strengthened repeated hook replacement,
and wrote the actual lifecycle owner contract using the established object owner.
The 85-root baseline also covers existing sustained-spell, spell-effect/start/
lifecycle, AI-spell, enchantment/ability restore and server orchestration captures.
Keep existing assertions and goldens unchanged; investigate any mismatch.

Production source is identical to qualified update commit a8d89bda, whose
production baseline is reused. Freeze exact focused names in all three profiles.
After conversion add a native registry contract, rerun focused profiles, then
run the complete default port corpus because this introduces shared duration
dispatch. Also run safe/static checks, three production/ABI builds, exact known
asset-suite comparison and a fresh headless creation/save/load/resume scenario.
The independent C observers test the remaining foreign callback boundary; they
are not retained copies of engine algorithms.

## Frozen original result

All 85 roots pass in default/server/high-resolution profiles without skips in
`build/port-duration-identities/baseline-original/`. Exact discovered/completed
root sets match the audited selection. Production source is identical to
`a8d89bda`; only the three new porttest files differ. No fixture setup correction
was needed during this baseline run.

[Original baseline](duration-identities-baseline.json),
[test selection](duration-identities-tests.txt),
[qualification manifest](duration-identities-batch.json).

Disk housekeeping removed 14 verified superseded item/transfer executables
(785,170,432 allocated bytes) and 1,654 verified duplicate asset files from the
completed update scenario (559,837,184 bytes). Originals, saves and results
remain; recovery paths and manifests are recorded in PORTING_STATE.md.
