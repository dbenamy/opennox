# Main roaming update — 2026-09-11

Scope: 5457E0 and its registration; retire the three history exports once this
last C owner is native. Pathfinding, actual movement, move audio and sound
investigation stay outside this conversion. Reuse ai_roam_porttest.go storage,
guards, normalized IDs and full object/monster change snapshots.

Original C passes 12,289 owner scenarios and the existing 98,304 history cases.
Coverage includes threshold neighbors, action-stack capacity, random idle pushes,
current enemy interruption, absent/invalid waypoint, real nearest-waypoint lookup,
scripted fallback/path boundary responses, arrival and path failure. Real C
movement runs against an empty wall map, including reached-path completion.
Independent assertions check aggression gating, fight interruption, idle RNG/capacity
and an arrival discriminator where float32 summation would incorrectly round a
squared distance above 64 down to 64. Disassembly confirms PC53 deltas and sums
without intermediate float32 spills. Debug logging remains disabled in fixtures.
Retaliation's deeper enemy processing and sound investigation remain existing
callees; this fixture exercises their early returns, not their entire behavior.

Two 200,000-update runs preserve complete final state and callback traces, through
the real action registry. The scripted detailed-path boundary records every call;
reported timings include that tracing overhead. C timings varied by run: roughly
423–563 ns/update for the empty-path response and 520–689 ns/update for reached
path completion. These VM microbenchmarks are not whole-game performance claims.
Logs/disassembly: build/port-roam-owner. Permanent hashes retain the baseline
without keeping test-only C bodies. C is still 140,082 lines at this baseline.

Native conversion passes every original-C state hash and the independent checks.
Both repeated-update state/trace hashes also match. First native measurements:
166 ns/update (empty path) and 297 ns/update (reached path), including identical
trace collection; both improve on the observed C range. No speedup is claimed
for path construction or whole-game execution. The existing Go aggression helpers
are unchanged for their current callers; the new owner preserves C's double
thresholds privately. All three now-unused history C bridges and the owner C
registration/entry are removed. Production C is **139,951 physical lines (minus
131)**, 153 files, zero test-reference C. Original-C baseline: 3a9480b2.

Accumulated tests pass in default/server/highres and all three production
binaries build as ELF32/80386 with GO386=sse2. The full suite matches the exact
known baseline: 1,553 failure entries, 15 passing/3 failing/32 skipped-no-test
packages. Fresh roam-update-port gameplay exits 0 against both preserved
screenshots with overrides disabled. All qualification is complete.
