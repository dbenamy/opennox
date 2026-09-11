# Guard, escort, and sound investigation — 2026-09-11

One connected batch: 545DA0, 546010, 546410/546420/546430, 546600,
5466B0 and 5466F0 (about 397 physical C lines with separators). Include updates
and their private helpers together; reroute existing native roam/idle sound calls.
Keep healing, combat, pathfinding and generic script-object lookup outside scope.

Original-C baseline: 22,723 scenarios using the existing guarded AI object,
monster, target, definition and normalized-pointer fixture. Added sparse C-owned
player records/units, script-name object and pending lists, read-only terrain
fixtures, real empty-map/ray traversal and a recorded path-precheck response.
Existing 12,289 roaming-owner and 98,304 history checks continue to pass.

Coverage: guard/escort lifecycle and updates, mimic/plant/ordinary kinds, aggression
threshold neighbors, full/partial action stacks, damager response, sound absence/
freshness/expiry, uint32 frame/deadline wrap, zero and normal tick rates, dry/water
and subclass bypass, precheck rejection and clear/out-of-map rejected rays. Resolver
covers owner/nil owner, zero to three sparse players plus active nil-unit slot,
script suffix/exact hits, pending-list hit and case-sensitive miss. Independent
assertions check resolver selection/RNG, sound freshness and injured return even
on a full stack. The full normalized changed-word state is hashed before porting.

Three independent precision cases distinguish PC53 double sums from float32 sums
at guard-home distance, escort radius and guard sight range. Disassembly confirms
that guard facing retains double deltas and length despite a redundant float32
memory store of X; only the normalized point passed to the dot helper is float32.
The original dot helper's first-product spill is already covered by AI tests.

Two 200,000-call real-registry idle-update measurements preserve final state;
initial C times were 752 ns/guard update and 160 ns/escort update. These deliberately
hold frame and quiet state fixed to isolate caller cost, not whole-game speed.
Debug logging and deeper healing/combat behavior remain outside this fixture.
Artifacts: build/port-guard-escort. Source C remains 139,951 physical lines,
153 files, zero test-reference C. A bounded helper supplied audit and a sound-only
draft; primary owns integration, arithmetic/state review and qualification.

Native conversion (baseline commit 54614975): all eight C bodies and their unused
entry declarations are removed. Guard/escort register native actions, their
lifecycle and resolution helpers are private Go, and native roam/idle call native
sound investigation. All 22,723 state hashes and independent checks match, as do
the existing roaming corpus and repeated-update hashes. Private type checks read
the same authoritative mimic/plant caches; the C predicates remain for their real
C callers. Native guard uses the already-tested Go facing-dot helper, while its
C counterpart remains for other production callers. No test-only C body remains.

Primary review retained the live heard-point pointer at the precheck boundary
and lazy frame access/deadline reads from the C caller. Disassembly-guided guard
normalization and explicit double aggression/distance arithmetic preserve the
baseline. First native quiet-update measurements were 211 ns/guard and 148
ns/escort, versus 752/160 in the original run; guard improved substantially,
while the escort difference is small relative to observed VM timing variation.

Production C: **139,549 physical lines (minus 402)**, 153 files, zero reference C.
One accumulated three-configuration test/build cycle, full-suite comparison and
fresh gameplay qualify this entire connected batch. Accumulated tests now pass
for default/server/highres. The full suite exactly matches the known baseline
(1,553 failure entries; 15 passing/3 failing/32 skipped-no-test packages).
All three production builds pass and identify as ELF32/80386 with GO386=sse2.
Fresh guard-escort-port gameplay exits 0 against both preserved screenshots,
overrides off. Qualification is complete.

Next grouping: move-to/far-move/dodge/flee/home and retreat/master lifecycle and
private policy helpers (about 347 C lines). Keep actual path/movement engines,
cast/heal policy and generic food lookup outside scope. Extend the shared fixture
with health, preceding action, one-shot movement flag, retreat-generator results
and edible object/map state. Audit: build/port-ai-navigation/audit.md (ignored).
Preserve negated resume eligibility for NaNs, byte-only roam-mask initialization,
and actual compiler spill order for dodge/movement radii.
