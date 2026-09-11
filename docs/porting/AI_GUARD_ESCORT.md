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
