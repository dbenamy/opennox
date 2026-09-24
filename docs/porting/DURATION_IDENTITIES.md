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


## Qualified conversion

All 53 callback identities now use distinct process-lifetime Go byte keys. A typed
server dispatcher calls the existing Go owners, preserving separate result and
void fallback paths for unknown C addresses. Root create/update/destroy calls,
getters and fixture routes migrate together. Mutable wall hooks remain late-bound;
32-bit record layout, nil guards and cancellation ordering remain unchanged.
The 42 sparse sustained fixture operations and six summon/charm operation keys
retain their exact numeric IDs and frozen snapshots. The char-returning cancel
adapter keeps signed low-byte behavior; historical float arguments preserve the
record address they transported.

All 86 focused roots pass in each of default/server/highres without skips. The
full default accumulated corpus completes exactly 2,452 roots: 2,451 passes and the
established TestMapPopulationPrerequisiteProbe diagnostic skip. Independent
expected-name sets match execution, not merely counts. Safe/static checks, all
three production/ABI binaries, exact known asset-suite outcomes and fresh headless
character creation/save/load/resume pass. All 1,654 original asset hashes remain
unchanged. Accepted phases share exact source fingerprints.

Progress: selected production cgo files 198→195; legacy export bridges 739→686;
157 headers now contain 3,421 physical lines (50 removed). Embedded production
callback bodies remain 77. Standalone production and test-reference C remain 0 LOC.
The immediate phase has eliminated 268/463 cgo files and 1,204/1,890 export bridges
on net. External native bindings are unchanged.

Luna supplied the bounded 19-path conversion overlay. Primary added the native
API test, checked all 53 owner/getter/result mappings and 42 sparse fixture routes,
reviewed retained export bodies/signatures and scanned whole-source references.
Exactly 53 exports retire; 55 remaining exact-name lines are fixture metadata only.
Primary corrected a generated fixture call that passed an argument to a zero-arg
owner before compilation. First compile then found a stale unsafe import in
spell_start.go; it was removed and the corrected run is contracts-fixed/. No test
ran in the failed compile and no existing assertion or capture changed. These two
small corrections support continuing bounded delegation with independent review;
no subscription-cost saving is measured.

Removed 29 obsolete root/legacy Go cache archives after all test jobs exited
and host-use/hash checks passed: 1,695,002,624 allocated bytes reclaimed. Current
duration caches, all binaries, source and original assets remain. Rebuild older
archives normally; cleanup manifests are under the batch directory.

[Qualification](duration-identities-qualification.json),
[updated dependency inventory](duration-identities-inventory-after.json).
Local evidence: build/port-duration-identities/ (contracts-fixed, corpus, safe,
production, inventory-after). Original baseline is commit 12a8ed98.
